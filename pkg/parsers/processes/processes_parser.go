package processes

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"

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
	ps := NewPS()

	cmd := exec.Command(p.path, "-e")
	cmd.Stdout = &buf
	err := cmd.Run()
	if err != nil {
		return ps, fmt.Errorf("Processes parser: %s\n", err)
	}

	s := bufio.NewScanner(&buf)
	s.Split(eachProcess)
	for s.Scan() {
		ps.Add(s.Text())
	}

	return ps, nil
}

// Split function for parsing the "ps -e" output line.
// It drops first 26 chars from each line, leaving only the cmd name.
// And also drops everything after "/" in the cmd name.
func eachProcess(data []byte, atEOF bool) (advance int, token []byte, err error) {
	advance, token, err = bufio.ScanLines(data, atEOF)

	if len(token) > 0 {
		token = token[26:]
		i := bytes.Index(token, []byte("/"))
		if i > 0 {
			token = token[0:i]
		}
	}

	return advance, token, err
}
