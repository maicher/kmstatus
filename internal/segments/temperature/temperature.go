package temperature

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maicher/kmstatus/internal/segments/common"
	"github.com/maicher/kmstatus/internal/types"
)

type Temperature struct {
	common.PeriodicParser
	common.Template

	Data   []Data
	Parser *TemperatureParser
}

func New(tmpl string, refreshInterval time.Duration) (types.Segment, error) {
	var t Temperature
	var err error

	t.Parser, err = NewTemperatureParser()
	if err != nil {
		return &t, err
	}

	t.PeriodicParser = common.NewPeriodicParser(t.read, t.parse, refreshInterval)

	err = t.NewTemplate(tmpl, helpers)
	if err != nil {
		return &t, fmt.Errorf("Unable to parse Temperature template: %s", err)
	}

	for _, name := range t.Parser.Names() {
		t.Data = append(t.Data, Data{Name: name})
	}

	return &t, nil
}

func (t *Temperature) Refresh() {
	t.PeriodicParser.Parse()
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
