package parsers

import (
	"os"
	"os/signal"
	"syscall"
	"time"
)

type NewParserFunc func() (Parser, error)

type Parser interface {
	Parse() (any, error)
}

type ParserWithInterval struct {
	parser   Parser
	interval time.Duration
}

type Controller struct {
	parsers               []ParserWithInterval
	parsersSensitiveToSig []ParserWithInterval
}

func NewParsersList() *Controller {
	return &Controller{
		parsers: []ParserWithInterval{},
	}
}

func (c *Controller) MustInit(newParser NewParserFunc, interval time.Duration, sensitiveToSig bool) {
	p, err := newParser()
	if err != nil {
		panic(err)
	}

	pi := ParserWithInterval{
		parser:   p,
		interval: interval,
	}

	c.parsers = append(c.parsers, pi)
	if sensitiveToSig {
		c.parsersSensitiveToSig = append(c.parsersSensitiveToSig, pi)
	}
}

func (c *Controller) StartParsing(ch chan<- any) {
	// Parse periodically
	for _, p := range c.parsers {
		go func(p ParserWithInterval) {
			onTick(p.interval, func() {
				parse(p.parser, ch)
			})
		}(p)
	}

	// Parse data on signal.
	go onSignal(syscall.SIGUSR1, func() {
		for _, p := range c.parsersSensitiveToSig {
			parse(p.parser, ch)
		}

		ch <- "render"
	})

	go onTick(time.Second, func() {
		ch <- "render"
	})
}

func parse(p Parser, ch chan<- any) {
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
