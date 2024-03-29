package mem

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/maicher/kmst/internal/segments"
)

type parse struct{}
type read struct{ buffer *bytes.Buffer }

type Mem struct {
	segments.Segment

	Data Data

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

	m.Template, err = template.New("").Funcs(helpers).Parse(conf.StrippedTemplate())
	if err != nil {
		return &m, fmt.Errorf("Unable to parse Mem template: %s", err)
	}

	go m.OnTick(conf.ParseInterval, func() {
		m.MsgQueue <- parse{}
	})

	return &m, nil
}

func (t *Mem) Read(b *bytes.Buffer) {
	t.MsgQueue <- read{buffer: b}
	<-t.Sync
}

func (m *Mem) handleMsg(msg any) error {
	var err error

	switch msg := msg.(type) {
	case read:
		err = m.Template.Execute(msg.buffer, m.Data)
		if err != nil {
			break
		}

		m.Sync <- struct{}{}
	case parse:
		err = m.Parser.Parse(&m.Data)
	default:
		panic("Invalid message")
	}

	if err != nil {
		return err
	}

	return nil
}
