package clock

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/maicher/kmst/internal/segments"
)

type Clock struct {
	Template *template.Template
}

func New(conf segments.Config) (segments.Reader, error) {
	var c Clock
	var err error

	c.Template, err = template.New("").Funcs(helpers).Parse(conf.StrippedTemplate())
	if err != nil {
		return &c, fmt.Errorf("Unable to parse CPU template: %s", err)
	}

	return &c, nil
}

func (c *Clock) Read(b *bytes.Buffer) {
	c.Template.Execute(b, time.Now())
}
