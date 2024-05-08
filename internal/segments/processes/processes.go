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
	var r *strings.Reader
	var d data

	p.parser, err = NewParser()
	if err != nil {
		return &p, err
	}

	s := bufio.NewScanner(strings.NewReader(tmpl))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		r = strings.NewReader(s.Text())
		fmt.Fscanf(r, "%s %s", &d.icon, &d.phrase)
		p.data = append(p.data, d)
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
