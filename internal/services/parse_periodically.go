package services

import (
	"time"

	"github.io/maicher/kmstatus/internal/config"
	"github.io/maicher/kmstatus/pkg/parsers"
	"github.io/maicher/kmstatus/pkg/parsers/factory"
)

type ParsePeriodically struct {
	Ch              chan<- any
	ParsersSettings []config.ParserSettings
}

func NewParsePeriodically(ch chan<- any, c []config.ParserSettings) *ParsePeriodically {
	return &ParsePeriodically{
		Ch:              ch,
		ParsersSettings: c,
	}
}

func (pp *ParsePeriodically) Loop() error {
	for _, pc := range pp.ParsersSettings {
		p, err := factory.NewParser(pc.Name)
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
