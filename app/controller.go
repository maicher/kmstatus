package app

import (
	"fmt"
	"os"
	"time"

	"github.io/maicher/stbar/pkg/parsers/cpu"
	"github.io/maicher/stbar/pkg/parsers/mem"
)

type Controller struct {
	ch          chan any
	parsersList *ParsersList
	view        *View
}

func (c *Controller) Loop() {
	d := &Data{}

	c.parsersList.StartParsing(c.ch)

	for {
		switch data := (<-c.ch).(type) {
		case cpu.Freq:
			d.CPU.Freq = data
		case cpu.Load:
			d.CPU.Load = data
		case cpu.Temp:
			d.CPU.Temp = data
		case mem.Mem:
			d.Mem = data
		case Render:
			fmt.Println(c.view.Render(d))
		case error:
			fmt.Fprintf(os.Stderr, "%s\n", data)
		}
	}

}

func NewController() *Controller {
	pl := NewParsersList()
	pl.MustInit(cpu.NewFreqParser, time.Second, true)
	pl.MustInit(cpu.NewLoadParser, time.Second, true)
	pl.MustInit(cpu.NewTempParser, time.Second, false)
	pl.MustInit(mem.NewMemParser, time.Second, false)

	return &Controller{
		ch:          make(chan any),
		parsersList: pl,
		view:        NewView("basic.txt.tmpl"),
	}
}
