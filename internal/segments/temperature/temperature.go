package temperature

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/maicher/kmst/internal/segments"
)

type parse struct{}
type read struct{ buffer *bytes.Buffer }

type Temperature struct {
	segments.Segment

	Data []Data

	Parser *TemperatureParser
}

func New(conf segments.Config) (segments.Reader, error) {
	var t Temperature
	var err error

	t.Segment = segments.NewSegment(t.handleMsg)

	t.Parser, err = NewTemperatureParser()
	if err != nil {
		return &t, err
	}

	t.Template, err = template.New("").Parse(conf.StrippedTemplate())
	if err != nil {
		return &t, fmt.Errorf("Unable to parse Temperature template: %s", err)
	}

	go t.OnTick(conf.ParseInterval, func() {
		t.MsgQueue <- parse{}
	})

	for _, name := range t.Parser.Names() {
		t.Data = append(t.Data, Data{Name: name})
	}

	return &t, nil
}

func (t *Temperature) Read(b *bytes.Buffer) {
	t.MsgQueue <- read{buffer: b}
	<-t.Sync
}

func (t *Temperature) handleMsg(msg any) error {
	var err error

	switch msg := msg.(type) {
	case read:
		for i := range t.Data {
			err = t.Template.Execute(msg.buffer, t.Data[i])
			if err != nil {
				break
			}
		}

		t.Sync <- struct{}{}
	case parse:
		err = t.Parser.Parse(t.Data)
	default:
		panic("Invalid message")
	}

	if err != nil {
		return err
	}

	return nil
}
