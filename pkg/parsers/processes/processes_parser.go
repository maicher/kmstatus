package processes

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.io/maicher/kmstatus/pkg/parsers"
)

type ProcessesParser struct {
	path string
}

func NewProcessesParser() (parsers.Parser, error) {
	parser := &ProcessesParser{}

	path, err := exec.LookPath("ps")
	if err != nil {
		return parser, fmt.Errorf("Processes parser: %s\n", err)
	}

	parser.path = path

	return parser, nil
}

func (p *ProcessesParser) Parse() (any, error) {
	var buf bytes.Buffer
	var ignored string
	var item string

	ps := PS{}
	ps.m = make(map[string]int)

	cmd := exec.Command(p.path, "-e")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return ps, fmt.Errorf("Processes parser: %s\n", err)
	}

	scan(&buf, func(line string) {
		r := strings.NewReader(line)
		fmt.Fscanf(r, "%s %s %s %s", &ignored, &ignored, &ignored, &item)
		name := firstSegment(item)
		_, exist := ps.m[name]
		if exist {
			ps.m[name] = ps.m[name] + 1
		} else {
			ps.m[name] = 1
		}
	})

	return ps, nil
}

func scan(buf io.Reader, f func(line string)) {
	s := bufio.NewScanner(buf)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		f(s.Text())
	}
}
func firstSegment(s string) string {
	segments := strings.SplitN(s, "/", 2)
	return segments[0]
}
