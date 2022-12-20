package app

import (
	"fmt"
	"os"

	"github.io/maicher/kmstatus/app/view"
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
)

type Controller struct {
	Ch   <-chan any
	View *view.View
}

func NewController(ch <-chan any, v *view.View) *Controller {
	return &Controller{
		Ch:   ch,
		View: v,
	}
}

func (c *Controller) AggregateDataAndRenderView() {
	d := &view.Data{}

	for {
		switch val := (<-c.Ch).(type) {
		case cpu.Freq:
			d.CPU.Freq = val
		case cpu.Load:
			d.CPU.Load = val
		case cpu.Temp:
			d.CPU.Temp = val
		case mem.Mem:
			d.Mem = val
		case view.RenderView:
			fmt.Println(c.View.Render(d))
		case error:
			fmt.Fprintf(os.Stderr, "%s\n", val)
		}
	}
}
