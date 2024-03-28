package cpu

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/maicher/kmst/internal/segments"
)

type parseLoad struct{}
type parseFreq struct{}
type write struct{ buffer *bytes.Buffer }

type CPU struct {
	Data Data

	LoadParser *LoadParser
	FreqParser *FreqParser

	msgQueue chan any
	sync     chan any

	tmpl *template.Template
}

func New(conf segments.SegmentConfig) (segments.Segment, error) {
	var c CPU
	var err error

	c.LoadParser, err = NewLoadParser()
	if err != nil {
		return &c, err
	}

	c.FreqParser, err = NewFreqParser()
	if err != nil {
		return &c, err
	}

	c.tmpl, err = template.New("").Parse(conf.StrippedTemplate())
	if err != nil {
		return &c, fmt.Errorf("Unable to parse CPU template: %s", err)
	}

	c.msgQueue = make(chan any)
	c.sync = make(chan any)

	go c.loop()

	go onTick(conf.ParseInterval, func() {
		c.msgQueue <- parseLoad{}
		c.msgQueue <- parseFreq{}
	})

	return &c, nil
}

func (c *CPU) Read(b *bytes.Buffer) {
	c.msgQueue <- write{buffer: b}
	<-c.sync
}

func onTick(interval time.Duration, f func()) {
	f()

	ticker := time.NewTicker(interval)
	for range ticker.C {
		f()
	}
}

func (c *CPU) loop() {
	var err error

	for msg := range c.msgQueue {
		switch msg := msg.(type) {
		case write:
			err = c.tmpl.Execute(msg.buffer, c.Data)

			c.sync <- struct{}{}
		case parseLoad:
			err = c.LoadParser.Parse(&c.Data.Load)
		case parseFreq:
			err = c.FreqParser.Parse(&c.Data.Freq)
		default:
			panic("Invalid message")
		}

		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}
