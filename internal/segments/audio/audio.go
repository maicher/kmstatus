package audio

import (
	"bytes"
	"fmt"

	"github.com/maicher/kmst/internal/segments"
)

type Audio struct {
	segments.Segment
	segments.Template

	Data   Data
	Parser *AudioParser
}

func New(conf segments.Config) (segments.ParseReader, error) {
	var a Audio
	var err error

	a.Parser, err = NewAudioParser()
	if err != nil {
		return &a, err
	}

	a.Segment = segments.NewSegment(a.read, a.parse, conf.ParseInterval)

	err = a.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &a, fmt.Errorf("Unable to parse Audio template: %s", err)
	}

	return &a, nil
}

func (a *Audio) read(b *bytes.Buffer) error {
	return a.Tmpl.Execute(b, a.Data)
}

func (a *Audio) parse() error {
	return a.Parser.Parse(&a.Data)
}
