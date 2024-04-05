package clock

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maicher/kmst/internal/segments/common"
	"github.com/maicher/kmst/internal/types"
)

type Clock struct {
	common.Template
}

func New(conf types.Config) (types.Segment, error) {
	var c Clock
	var err error

	err = c.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &c, fmt.Errorf("Unable to parse CPU template: %s", err)
	}

	return &c, nil
}

func (c *Clock) Refresh() {
}

func (c *Clock) Read(b *bytes.Buffer) {
	c.Tmpl.Execute(b, time.Now())
}

func (c *Clock) Parse() {
}
