package processes

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/maicher/kmstatus/internal/segments/common"
	"github.com/maicher/kmstatus/internal/types"
)

type Processes struct {
	common.PeriodicParser
	common.Template

	data   []data
	parser *Parser
}

func New(tmpl string, refreshInterval time.Duration) (types.Segment, error) {
	var p Processes
	var err error

	p.parser, err = NewParser()
	if err != nil {
		return &p, err
	}

	p.data, err = parseTemplate(tmpl)
	if err != nil {
		return &p, err
	}

	p.PeriodicParser = common.NewPeriodicParser(p.read, p.parse, refreshInterval)

	return &p, err
}

func (p *Processes) Refresh() {
	p.PeriodicParser.Parse()
}

func (p *Processes) read(b *bytes.Buffer) (err error) {
	for _, d := range p.data {
		if d.active {
			_, err = b.WriteString(d.icon)
			if err != nil {
				break
			}
		}
	}

	return err
}

func (p *Processes) parse() error {
	return p.parser.Parse(p.data)
}

// parseTemplate parses lines in the format: ICON PROCESS NAME.
// The process name may contain spaces. Blank lines are skipped.
func parseTemplate(tmpl string) ([]data, error) {
	var ds []data
	var n int

	s := bufio.NewScanner(strings.NewReader(tmpl))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		n++

		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}

		icon := strings.Fields(line)[0]
		phrase := strings.TrimSpace(line[len(icon):])
		if phrase == "" {
			return nil, fmt.Errorf("invalid Processes template, line %d: %q, want: ICON PROCESS", n, s.Text())
		}

		ds = append(ds, data{icon: icon, phrase: phrase})
	}

	return ds, nil
}
