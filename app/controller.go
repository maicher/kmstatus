package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.io/maicher/kmstatus/pkg/parsers"
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
)

const sig = syscall.SIGUSR1

type renderCommand struct{}

type Controller struct {
	view            *view
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

func NewController() *Controller {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal error: %s\n", err)
			os.Exit(1)
		}
	}()

	c := &Controller{
		view:           newView("basic.txt.tmpl"),
		renderInterval: time.Second,
	}

	pb := newParsersBuilder(c)
	pb.mustInit(cpu.NewFreqParser, time.Second, true)
	pb.mustInit(cpu.NewLoadParser, time.Second, true)
	pb.mustInit(cpu.NewTempParser, time.Second, false)
	pb.mustInit(mem.NewMemParser, time.Second, false)

	return c
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

		ch <- renderCommand{}
	})
}

// Render
func (c *Controller) generateRenderCommands(ch chan<- any) {
	onTick(c.renderInterval, func() {
		ch <- renderCommand{}
	})
}

func (c *Controller) aggregateDataAndRenderView(ch <-chan any) {
	d := &data{}

	for {
		switch data := (<-ch).(type) {
		case cpu.Freq:
			d.CPU.Freq = data
		case cpu.Load:
			d.CPU.Load = data
		case cpu.Temp:
			d.CPU.Temp = data
		case mem.Mem:
			d.Mem = data
		case renderCommand:
			fmt.Println(c.view.render(d))
		case error:
			fmt.Fprintf(os.Stderr, "%s\n", data)
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
