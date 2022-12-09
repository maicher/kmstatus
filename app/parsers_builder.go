package app

import (
	"time"

	"github.io/maicher/stbar/pkg/parsers"
)

type newParserFunc func() (parsers.Parser, error)

type parser struct {
	parser   parsers.Parser
	interval time.Duration
}

type parsersBuilder struct {
	controller *Controller
}

func newParsersBuilder(c *Controller) *parsersBuilder {
	return &parsersBuilder{
		controller: c,
	}
}

func (pl *parsersBuilder) mustInit(newParser newParserFunc, interval time.Duration, sensitiveToSig bool) {
	p, err := newParser()
	if err != nil {
		panic(err)
	}

	pi := parser{
		parser:   p,
		interval: interval,
	}

	pl.controller.parsers = append(pl.controller.parsers, pi)

	if sensitiveToSig {
		pl.controller.parsersSensitiveToSig = append(pl.controller.parsersSensitiveToSig, pi)
	}
}
