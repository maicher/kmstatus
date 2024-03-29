package cpu

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/maicher/kmst/internal/segments"
)

type parseLoad struct{}
type parseFreq struct{}
type read struct{ buffer *bytes.Buffer }

type CPU struct {
	segments.Segment

	Data Data

	LoadParser *LoadParser
	FreqParser *FreqParser
}

func New(conf segments.Config) (segments.Reader, error) {
	var c CPU
	var err error

	c.Segment = segments.NewSegment(c.handleMsg)

	c.LoadParser, err = NewLoadParser()
	if err != nil {
		return &c, err
	}

	c.FreqParser, err = NewFreqParser()
	if err != nil {
		return &c, err
	}

	c.Template, err = template.New("").Funcs(helpers).Parse(conf.StrippedTemplate())
	if err != nil {
		return &c, fmt.Errorf("Unable to parse CPU template: %s", err)
	}

	go c.OnTick(conf.ParseInterval, func() {
		c.MsgQueue <- parseLoad{}
		c.MsgQueue <- parseFreq{}
	})

	return &c, nil
}

func (c *CPU) Read(b *bytes.Buffer) {
	c.MsgQueue <- read{buffer: b}
	<-c.Sync
}

func (c *CPU) handleMsg(msg any) error {
	var err error

	switch msg := msg.(type) {
	case read:
		err = c.Template.Execute(msg.buffer, c.Data)

		c.Sync <- struct{}{}
	case parseLoad:
		err = c.LoadParser.Parse(&c.Data.Load)
	case parseFreq:
		err = c.FreqParser.Parse(&c.Data.Freq)
	default:
		panic("Invalid message")
	}

	if err != nil {
		return err
	}

	return nil
}
