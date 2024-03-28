package temperature

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/maicher/kmst/internal/segments"
)

type parse struct{}
type write struct{ buffer *bytes.Buffer }

type Temperature struct {
	Data []Data

	Parser *TemperatureParser

	msgQueue chan any
	sync     chan any

	tmpl *template.Template
}

func New(conf segments.SegmentConfig) (segments.Segment, error) {
	var t Temperature
	var err error

	t.Parser, err = NewTemperatureParser()
	if err != nil {
		return &t, err
	}

	t.tmpl, err = template.New("").Parse(conf.StrippedTemplate())
	if err != nil {
		return &t, fmt.Errorf("Unable to parse Temperature template: %s", err)
	}

	t.msgQueue = make(chan any)
	t.sync = make(chan any)

	go t.loop()

	go onTick(conf.ParseInterval, func() {
		t.msgQueue <- parse{}
	})

	for _, name := range t.Parser.Names() {
		t.Data = append(t.Data, Data{Name: name})
	}

	return &t, nil
}

func (t *Temperature) Read(b *bytes.Buffer) {
	t.msgQueue <- write{buffer: b}
	<-t.sync
}

func onTick(interval time.Duration, f func()) {
	f()

	ticker := time.NewTicker(interval)
	for range ticker.C {
		f()
	}
}

func (t *Temperature) loop() {
	var err error

	for msg := range t.msgQueue {
		switch msg := msg.(type) {
		case write:
			for i := range t.Data {
				err = t.tmpl.Execute(msg.buffer, t.Data[i])
				if err != nil {
					break
				}
			}

			t.sync <- struct{}{}
		case parse:
			err = t.Parser.Parse(t.Data)
		default:
			panic("Invalid message")
		}

		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}
