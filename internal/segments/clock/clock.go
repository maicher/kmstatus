package clock

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maicher/kmst/internal/segments"
)

type Clock struct {
	segments.Template
}

func New(conf segments.Config) (segments.Reader, error) {
	var c Clock
	var err error

	err = c.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &c, fmt.Errorf("Unable to parse CPU template: %s", err)
	}

	return &c, nil
}

func (c *Clock) Read(b *bytes.Buffer) {
	c.Tmpl.Execute(b, time.Now())
}
