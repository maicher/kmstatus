package cpu

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maicher/kmstatus/internal/segments/common"
	"github.com/maicher/kmstatus/internal/types"
)

type CPU struct {
	common.PeriodicParser
	common.Template

	Data Data

	LoadParser *LoadParser
	FreqParser *FreqParser
}

func New(tmpl string, refreshInterval time.Duration) (types.Segment, error) {
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

	c.PeriodicParser = common.NewPeriodicParser(c.read, c.parse, refreshInterval)

	err = c.NewTemplate(tmpl, helpers)
	if err != nil {
		return &c, fmt.Errorf("Unable to parse CPU template: %s", err)
	}

	return &c, nil
}

func (c *CPU) Refresh() {
}

func (c *CPU) read(b *bytes.Buffer) error {
	return c.Tmpl.Execute(b, c.Data)
}

func (c *CPU) parse() error {
	err := c.LoadParser.Parse(&c.Data.Load)
	if err != nil {
		return err
	}

	return c.FreqParser.Parse(&c.Data.Freq)
}
