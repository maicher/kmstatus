package factory

import (
	"errors"

	"github.io/maicher/kmstatus/pkg/parsers"
	"github.io/maicher/kmstatus/pkg/parsers/cpu"
	"github.io/maicher/kmstatus/pkg/parsers/filesystem"
	"github.io/maicher/kmstatus/pkg/parsers/mem"
	"github.io/maicher/kmstatus/pkg/parsers/processes"
)

var constructors = map[string]parsers.NewParserFunc{
	"cpu-freq": cpu.NewFreqParser,
	"cpu-temp": cpu.NewTempParser,
	"cpu-load": cpu.NewLoadParser,
	"mem":      mem.NewMemParser,
	"fs":       filesystem.NewFSParser,
	"ps":       processes.NewProcessesParser,
}

func NewParser(name string) (parsers.Parser, error) {
	newParserFunc, ok := constructors[name]
	if !ok {
		return nil, errors.New("Invalid parser name: " + name)
	}

	return newParserFunc()
}
