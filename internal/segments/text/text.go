package text

import (
	"bytes"
	"fmt"

	"github.com/maicher/kmst/internal/segments"
)

type Text struct {
	text string

	segments.Segment
	segments.Template
}

func New(conf segments.Config) (segments.ParseReader, error) {
	var t Text
	var err error

	t.text = ""

	err = t.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &t, fmt.Errorf("Unable to parse Text template: %s", err)
	}

	t.Segment = segments.NewSegment(t.read, t.parse, conf.ParseInterval)

	return &t, nil
}

func (t *Text) read(b *bytes.Buffer) error {
	return t.Tmpl.Execute(b, t.text)
}

func (t *Text) parse() error {
	return nil
}
