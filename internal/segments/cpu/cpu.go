package cpu

import (
	"bytes"
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

	c.LoadParser, err = NewLoadParser()
	if err != nil {
		return &c, err
	}

	c.FreqParser, err = NewFreqParser()
	if err != nil {
		return &c, err
	}

	c.Segment = segments.NewSegment(c.readMsg, c.parseMsg, conf.ParseInterval)

	err = c.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &c, fmt.Errorf("Unable to parse CPU template: %s", err)
	}

	return &c, nil
}

func (c *CPU) readMsg(b *bytes.Buffer) error {
	return c.Tmpl.Execute(b, c.Data)
}

func (c *CPU) parseMsg() error {
	err := c.LoadParser.Parse(&c.Data.Load)
	if err != nil {
		return err
	}

	return c.FreqParser.Parse(&c.Data.Freq)
}
