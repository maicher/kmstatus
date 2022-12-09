package app

import (
	"time"

	"github.io/maicher/stbar/pkg/parsers"
)

type ParserWithInterval struct {
	parser   parsers.Parser
	interval time.Duration
}

type ParsersBuilder struct {
	controller *Controller
}

func NewParsersBuilder(c *Controller) *ParsersBuilder {
	return &ParsersBuilder{
		controller: c,
	}
}

func (pl *ParsersBuilder) MustInit(newParser parsers.NewParserFunc, interval time.Duration, sensitiveToSig bool) {
	p, err := newParser()
	if err != nil {
		panic(err)
	}

	pi := ParserWithInterval{
		parser:   p,
		interval: interval,
	}

	pl.controller.parsers = append(pl.controller.parsers, pi)

	if sensitiveToSig {
		pl.controller.parsersSensitiveToSig = append(pl.controller.parsersSensitiveToSig, pi)
	}
}
