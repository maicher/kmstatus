package filesystem

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.io/maicher/kmstatus/pkg/parsers"
)

const cmdString = "df"

type FSParser struct {
	path string
}

func (p *FSParser) Parse() (any, error) {
	var buf bytes.Buffer
	fs := NewFS()

	cmd := exec.Command(p.path)
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return fs, fmt.Errorf("FS parser: %w", err)
	}

	s := bufio.NewScanner(&buf)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()

		if strings.HasPrefix(line, "/dev") {
			fs.Drives = append(fs.Drives, parseDrive(line))
		}

		if strings.HasPrefix(line, "encfs") {
			fs.ENCFS = true
		}
	}

	return fs, nil
}

func NewFSParser() (parsers.Parser, error) {
	parser := FSParser{}

	path, err := exec.LookPath(cmdString)
	if err != nil {
		return &parser, fmt.Errorf("FS parser: %s", err)
	}

	parser.path = path

	return &parser, nil
}

func parseDrive(s string) Drive {
	var ignored string
	d := Drive{}
	r := strings.NewReader(s)

	fmt.Fscanf(r, "%s %d %d %s %s %s", &(d.Name), &(d.Total), &(d.Used), &ignored, &ignored, &(d.MountedOn))
	i := strings.LastIndex(d.Name, "/")
	d.Name = d.Name[i+1:]
	d.Free = d.Total - d.Used

	return d
}
