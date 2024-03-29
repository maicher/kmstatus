package cpu

import (
	"fmt"

	"github.com/maicher/kmst/internal/segments"
)

type CPU struct {
	segments.Segment
	segments.Template

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

	err = c.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &c, fmt.Errorf("Unable to parse CPU template: %s", err)
	}

	go c.OnTick(conf.ParseInterval, func() {
		c.MsgQueue <- segments.ParseMsg{}
	})

	return &c, nil
}

func (c *CPU) handleMsg(msg any) error {
	var err error

	switch msg := msg.(type) {
	case segments.ReadMsg:
		err = c.Tmpl.Execute(msg.Buffer, c.Data)
		c.Sync <- struct{}{}
	case segments.ParseMsg:
		err = c.LoadParser.Parse(&c.Data.Load)
		if err != nil {
			return nil
		}
		err = c.FreqParser.Parse(&c.Data.Freq)
	default:
		panic("Invalid message")
	}

	if err != nil {
		return err
	}

	return nil
}
