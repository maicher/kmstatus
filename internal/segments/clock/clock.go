package clock

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maicher/kmstatus/internal/segments/common"
	"github.com/maicher/kmstatus/internal/types"
)

type Clock struct {
	common.Template
}

func New(tmpl string, refreshInterval time.Duration) (types.Segment, error) {
	var c Clock
	var err error

	err = c.NewTemplate(tmpl, helpers)
	if err != nil {
		return &c, fmt.Errorf("Unable to parse Clock template: %s", err)
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
