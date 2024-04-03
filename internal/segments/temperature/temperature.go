package temperature

import (
	"bytes"
	"fmt"

	"github.com/maicher/kmst/internal/segments"
)

type Temperature struct {
	segments.Segment
	segments.Template

	Data   []Data
	Parser *TemperatureParser
}

func New(conf segments.Config) (segments.RefreshReader, error) {
	var t Temperature
	var err error

	t.Parser, err = NewTemperatureParser()
	if err != nil {
		return &t, err
	}

	t.Segment = segments.NewSegment(t.read, t.parse, conf.RefreshInterval)

	err = t.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &t, fmt.Errorf("Unable to parse Temperature template: %s", err)
	}

	for _, name := range t.Parser.Names() {
		t.Data = append(t.Data, Data{Name: name})
	}

	return &t, nil
}

func (t *Temperature) Refresh() {
	t.Segment.Parse()
}

func (t *Temperature) read(b *bytes.Buffer) error {
	var err error

	for i := range t.Data {
		err = t.Tmpl.Execute(b, t.Data[i])
		if err != nil {
			break
		}
	}

	return err
}

func (t *Temperature) parse() error {
	return t.Parser.Parse(t.Data)
}
