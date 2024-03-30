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

func New(conf segments.Config) (segments.Reader, error) {
	var m Mem
	var err error

	m.Parser, err = NewMemParser()
	if err != nil {
		return &m, err
	}

	m.Segment = segments.NewSegment(m.readMsg, m.parseMsg, conf.ParseInterval)

	err = m.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &m, fmt.Errorf("Unable to parse Mem template: %s", err)
	}

	return &m, nil
}

func (m *Mem) readMsg(b *bytes.Buffer) error {
	return m.Tmpl.Execute(b, m.Data)
}

func (m *Mem) parseMsg() error {
	return m.Parser.Parse(&m.Data)
}
