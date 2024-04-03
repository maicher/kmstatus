package mem

import (
	"bytes"
	"fmt"

	"github.com/maicher/kmst/internal/segments"
)

type Mem struct {
	segments.Segment
	segments.Template

	Data   Data
	Parser *MemParser
}

func New(conf segments.Config) (segments.RefreshReader, error) {
	var m Mem
	var err error

	m.Parser, err = NewMemParser()
	if err != nil {
		return &m, err
	}

	m.Segment = segments.NewSegment(m.read, m.parse, conf.RefreshInterval)

	err = m.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &m, fmt.Errorf("Unable to parse Mem template: %s", err)
	}

	return &m, nil
}

func (m *Mem) Refresh() {
	m.Segment.Parse()
}

func (m *Mem) read(b *bytes.Buffer) error {
	return m.Tmpl.Execute(b, m.Data)
}

func (m *Mem) parse() error {
	return m.Parser.Parse(&m.Data)
}
