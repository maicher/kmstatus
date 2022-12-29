package app

import (
	"github.io/maicher/kmstatus/app/view"
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
)

type viewRenderer interface {
	Render(*view.Data)
	RenderErr(error)
}

type Controller struct {
	Ch   <-chan any
	View viewRenderer
}

func NewController(ch <-chan any, v viewRenderer) *Controller {
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
			c.View.Render(d)
		case error:
			c.View.RenderErr(val)
		}
	}
}
