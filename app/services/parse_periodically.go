package services

import (
	"time"

	"github.io/maicher/kmstatus/app/config"
	"github.io/maicher/kmstatus/pkg/parsers"
)

type ParsePeriodically struct {
	Ch            chan<- any
	ParsersConfig []config.ParserConfig
}

func NewParsePeriodically(ch chan<- any, c []config.ParserConfig) *ParsePeriodically {
	return &ParsePeriodically{
		Ch:            ch,
		ParsersConfig: c,
	}
}

func (pp *ParsePeriodically) Loop() error {
	for _, pc := range pp.ParsersConfig {
		p, err := pc.NewParserFunc()
		if err != nil {
			return err
		}

		go pp.parsePeriodically(p, pc.Interval)
	}

	return nil
}

// Parse periodically
func (pp *ParsePeriodically) parsePeriodically(p parsers.Parser, interval time.Duration) {
	onTick(interval, func() {
		parse(p, pp.Ch)
	})
}
