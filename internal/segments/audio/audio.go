package audio

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maicher/kmstatus/internal/segments/common"
	"github.com/maicher/kmstatus/internal/types"
)

type Audio struct {
	common.PeriodicParser
	common.Template

	Data   Data
	Parser *AudioParser
}

func New(tmpl string, refreshInterval time.Duration) (types.Segment, error) {
	var a Audio
	var err error

	a.Parser, err = NewAudioParser()
	if err != nil {
		return &a, err
	}

	a.PeriodicParser = common.NewPeriodicParser(a.read, a.parse, refreshInterval)

	err = a.NewTemplate(tmpl, helpers)
	if err != nil {
		return &a, fmt.Errorf("Unable to parse Audio template: %s", err)
	}

	return &a, nil
}

func (a *Audio) Refresh() {
	a.PeriodicParser.Parse()
}

func (a *Audio) read(b *bytes.Buffer) error {
	return a.Tmpl.Execute(b, a.Data)
}

func (a *Audio) parse() error {
	return a.Parser.Parse(&a.Data)
}
