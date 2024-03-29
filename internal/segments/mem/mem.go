package mem

import (
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

	m.Segment = segments.NewSegment(m.handleMsg)

	m.Parser, err = NewMemParser()
	if err != nil {
		return &m, err
	}

	err = m.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &m, fmt.Errorf("Unable to parse Mem template: %s", err)
	}

	go m.OnTick(conf.ParseInterval, func() {
		m.MsgQueue <- segments.ParseMsg{}
	})

	return &m, nil
}

func (m *Mem) handleMsg(msg any) error {
	var err error

	switch msg := msg.(type) {
	case segments.ReadMsg:
		err = m.Tmpl.Execute(msg.Buffer, m.Data)
		m.Sync <- struct{}{}
	case segments.ParseMsg:
		err = m.Parser.Parse(&m.Data)
	default:
		panic("Invalid message")
	}

	if err != nil {
		return err
	}

	return nil
}
