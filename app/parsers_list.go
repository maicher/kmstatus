package app

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.io/maicher/stbar/pkg/parsers"
)

type Render struct{}

type ParserWithInterval struct {
	parser   parsers.Parser
	interval time.Duration
}

type ParsersList struct {
	parsers               []ParserWithInterval
	parsersSensitiveToSig []ParserWithInterval
}

func NewParsersList() *ParsersList {
	return &ParsersList{
		parsers: []ParserWithInterval{},
	}
}

func (pl *ParsersList) MustInit(newParser parsers.NewParserFunc, interval time.Duration, sensitiveToSig bool) {
	p, err := newParser()
	if err != nil {
		panic(err)
	}

	pi := ParserWithInterval{
		parser:   p,
		interval: interval,
	}

	pl.parsers = append(pl.parsers, pi)
	if sensitiveToSig {
		pl.parsersSensitiveToSig = append(pl.parsersSensitiveToSig, pi)
	}
}

func (pl *ParsersList) StartParsing(ch chan<- any) {
	// Parse periodically
	for _, p := range pl.parsers {
		go func(p ParserWithInterval) {
			onTick(p.interval, func() {
				parse(p.parser, ch)
			})
		}(p)
	}

	// Parse data on signal.
	go onSignal(syscall.SIGUSR1, func() {
		for _, p := range pl.parsersSensitiveToSig {
			parse(p.parser, ch)
		}

		ch <- Render{}
	})

	go onTick(time.Second, func() {
		ch <- Render{}
	})
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
