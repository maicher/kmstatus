package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.io/maicher/stbar/pkg/parsers"
	"github.io/maicher/stbar/pkg/parsers/cpu"
	"github.io/maicher/stbar/pkg/parsers/mem"
)

type Render struct{}

type Controller struct {
	view                  *View
	interval              time.Duration
	parsers               []ParserWithInterval
	parsersSensitiveToSig []ParserWithInterval
}

func (c *Controller) Loop() {
	ch := make(chan any)

	c.startParsingPeriodically(ch)
	c.startParsingOnSig(ch)
	c.startRendering(ch)

	c.aggregateAndRender(ch)
}

func NewController() *Controller {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal error: %s\n", err)
			os.Exit(1)
		}
	}()

	c := &Controller{
		view:                  NewView("basic.txt.tmpl"),
		interval:              time.Second,
		parsers:               []ParserWithInterval{},
		parsersSensitiveToSig: []ParserWithInterval{},
	}

	pb := NewParsersBuilder(c)
	pb.MustInit(cpu.NewFreqParser, time.Second, true)
	pb.MustInit(cpu.NewLoadParser, time.Second, true)
	pb.MustInit(cpu.NewTempParser, time.Second, false)
	pb.MustInit(mem.NewMemParser, time.Second, false)

	return c
}

// Parse periodically
func (c *Controller) startParsingPeriodically(ch chan<- any) {
	for _, p := range c.parsers {
		go func(p ParserWithInterval) {
			onTick(p.interval, func() {
				parse(p.parser, ch)
			})
		}(p)
	}
}

// Parse data on signal.
func (c *Controller) startParsingOnSig(ch chan<- any) {
	go onSignal(syscall.SIGUSR1, func() {
		for _, p := range c.parsersSensitiveToSig {
			parse(p.parser, ch)
		}

		ch <- Render{}
	})
}

// Render
func (c *Controller) startRendering(ch chan<- any) {
	go onTick(c.interval, func() {
		ch <- Render{}
	})
}

func (c *Controller) aggregateAndRender(ch <-chan any) {
	d := &Data{}

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
		case Render:
			fmt.Println(c.view.Render(d))
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

func onSignal(s os.Signal, f func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, s)

	for {
		<-sigs
		f()
	}
}
