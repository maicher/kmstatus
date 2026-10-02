package processes

import (
	"bufio"
	"bytes"
	"strings"

	"github.com/maicher/kmstatus/internal/segments/common"
)

type Parser struct {
	command string
}

func NewParser() (*Parser, error) {
	var p Parser
	p.command = "ps"

	return &p, nil
}

func (p *Parser) Parse(data []data) error {
	var buf bytes.Buffer

	err := common.RunCommand(&buf, p.command, "-e", "-o", "comm=")
	if err != nil {
		return err
	}

	for i := range data {
		data[i].active = false
	}

	s := bufio.NewScanner(&buf)
	s.Split(bufio.ScanLines)
	for s.Scan() {
		for i, d := range data {
			if strings.HasPrefix(s.Text(), d.phrase) {
				data[i].active = true
				break
			}
		}
	}

	return nil
}
