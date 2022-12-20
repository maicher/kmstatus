package services

import (
	"os"
	"os/signal"
	"syscall"

	"github.io/maicher/kmstatus/app/config"
	"github.io/maicher/kmstatus/app/view"
	"github.io/maicher/kmstatus/pkg/parsers"
)

const sig = syscall.SIGUSR1

type ParseOnSig struct {
	Ch            chan<- any
	ParsersConfig []config.ParserConfig
}

func NewParseOnSig(ch chan<- any, c []config.ParserConfig) *ParseOnSig {
	return &ParseOnSig{
		Ch:            ch,
		ParsersConfig: c,
	}
}

func (pos *ParseOnSig) Loop() error {
	var onSigParsers []parsers.Parser

	for _, pc := range pos.ParsersConfig {
		p, err := pc.NewParserFunc()
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
