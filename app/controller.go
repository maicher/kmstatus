package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.io/maicher/kmstatus/app/config"
	"github.io/maicher/kmstatus/app/view"
	"github.io/maicher/kmstatus/pkg/parsers"
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
)

const sig = syscall.SIGUSR1

type renderView struct{}

type Controller struct {
	view            *view.View
	renderInterval  time.Duration
	periodicParsers []*periodicParser
	onSigParsers    []parsers.Parser
}

func (c *Controller) Loop() {
	ch := make(chan any)

	for _, p := range c.periodicParsers {
		go c.parsePeriodically(p, ch)
	}

	go c.parseOnSig(ch)
	go c.generateRenderCommands(ch)

	c.aggregateDataAndRenderView(ch)
}

func NewController(conf *config.Config) (c *Controller, err error) {
	defer func() {
		if r := recover(); r != nil {
			c = &Controller{}
			err = r.(error)
		}
	}()

	v, err := view.New(conf)
	if err != nil {
		return c, err
	}

	c = &Controller{
		view:           v,
		renderInterval: time.Second,
	}

	pb := newParsersBuilder(c)
	pb.mustInit(cpu.NewFreqParser, time.Second, true)
	pb.mustInit(cpu.NewLoadParser, time.Second, true)
	pb.mustInit(cpu.NewTempParser, time.Second, false)
	pb.mustInit(mem.NewMemParser, time.Second, false)

	return c, nil
}

// Parse periodically
func (c *Controller) parsePeriodically(p *periodicParser, ch chan<- any) {
	onTick(p.interval, func() {
		parse(p.parser, ch)
	})
}

// Parse data on signal.
func (c *Controller) parseOnSig(ch chan<- any) {
	onSig(sig, func() {
		for _, p := range c.onSigParsers {
			parse(p, ch)
		}

		ch <- renderView{}
	})
}

// Render
func (c *Controller) generateRenderCommands(ch chan<- any) {
	onTick(c.renderInterval, func() {
		ch <- renderView{}
	})
}

func (c *Controller) aggregateDataAndRenderView(ch <-chan any) {
	d := &view.Data{}

	for {
		switch val := (<-ch).(type) {
		case cpu.Freq:
			d.CPU.Freq = val
		case cpu.Load:
			d.CPU.Load = val
		case cpu.Temp:
			d.CPU.Temp = val
		case mem.Mem:
			d.Mem = val
		case renderView:
			fmt.Println(c.view.Render(d))
		case error:
			fmt.Fprintf(os.Stderr, "%s\n", val)
		}
	}
}

func parse(p parsers.Parser, ch chan<- any) {
	v, err := p.Parse()
	if err != nil {
		ch <- err
	}

	ch <- v
}

func onTick(interval time.Duration, f func()) {
	f()

	ticker := time.NewTicker(interval)
	for range ticker.C {
		f()
	}
}

func onSig(s os.Signal, f func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, s)

	for {
		<-sigs
		f()
	}
}
