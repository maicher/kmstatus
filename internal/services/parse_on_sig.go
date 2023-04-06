package services

import (
	"os"
	"os/signal"
	"syscall"

	"github.io/maicher/kmstatus/internal/config"
	"github.io/maicher/kmstatus/internal/view"
	"github.io/maicher/kmstatus/pkg/parsers"
	"github.io/maicher/kmstatus/pkg/parsers/factory"
)

const sig = syscall.SIGUSR1

type ParseOnSig struct {
	Ch              chan<- any
	ParsersSettings []config.ParserSettings
}

func NewParseOnSig(ch chan<- any, c []config.ParserSettings) *ParseOnSig {
	return &ParseOnSig{
		Ch:              ch,
		ParsersSettings: c,
	}
}

func (pos *ParseOnSig) Loop() error {
	var onSigParsers []parsers.Parser

	for _, pc := range pos.ParsersSettings {
		p, err := factory.NewParser(pc.Name)
		if err != nil {
			return err
		}

		if pc.OnSig {
			onSigParsers = append(onSigParsers, p)
		}
	}

	go pos.parseOnSig(onSigParsers)

	return nil
}

// Parse on signal
func (pos *ParseOnSig) parseOnSig(onSigParsers []parsers.Parser) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, sig)

	for {
		<-sigs

		for _, p := range onSigParsers {
			parse(p, pos.Ch)
		}

		pos.Ch <- view.RenderView{}
	}
}
