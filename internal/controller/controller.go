package controller

import (
	"github.io/maicher/kmstatus/internal/view"
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/filesystem"
	"github.io/maicher/kmstatus/pkg/parsers/gpu"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
	"github.io/maicher/kmstatus/pkg/parsers/network"
	"github.io/maicher/kmstatus/pkg/parsers/processes"
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
		case processes.PS:
			d.PS = val
		case cpu.Freq:
			d.CPU.Freq = val
		case cpu.Load:
			d.CPU.Load = val
		case cpu.Temp:
			d.CPU.Temp = val
		case gpu.GPU:
			d.GPU = val
		case mem.Mem:
			d.Mem = val
		case filesystem.FS:
			d.FS = val
		case network.Net:
			d.Net = val
		case view.RenderView:
			c.View.Render(d)
		case error:
			c.View.RenderErr(val)
		}
	}
}
