package mem

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maicher/kmstatus/internal/segments/common"
	"github.com/maicher/kmstatus/internal/types"
)

type Mem struct {
	common.PeriodicParser
	common.Template

	data   data
	parser *Parser
}

func New(tmpl string, refreshInterval time.Duration) (types.Segment, error) {
	var m Mem
	var err error

	m.parser, err = NewParser()
	if err != nil {
		return &m, err
	}

	m.PeriodicParser = common.NewPeriodicParser(m.read, m.parse, refreshInterval)

	err = m.NewTemplate(tmpl, helpers)
	if err != nil {
		return &m, fmt.Errorf("unable to parse Mem template: %s", err)
	}

	return &m, nil
}

func (m *Mem) Refresh() {
	m.PeriodicParser.Parse()
}

func (m *Mem) read(b *bytes.Buffer) error {
	return m.Tmpl.Execute(b, m.data)
}

func (m *Mem) parse() error {
	return m.parser.Parse(&m.data)
}
