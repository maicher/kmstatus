package app

import (
	"time"

	"github.io/maicher/stbar/pkg/parsers"
)

type newParserFunc func() (parsers.Parser, error)

type parsersBuilder struct {
	controller *Controller
}

func (pl *parsersBuilder) mustInit(newParser newParserFunc, interval time.Duration, sensitiveToSig bool) {
	p, err := newParser()
	if err != nil {
		panic(err)
	}

	pi := &periodicParser{
		parser:   p,
		interval: interval,
	}

	pl.controller.periodicParsers = append(pl.controller.periodicParsers, pi)

	if sensitiveToSig {
		pl.controller.onSigParsers = append(pl.controller.onSigParsers, p)
	}
}

func newParsersBuilder(c *Controller) *parsersBuilder {
	return &parsersBuilder{
		controller: c,
	}
}
