package segments

import (
	"bytes"

	"github.com/maicher/kmst/internal/types"
)

type Collection struct {
	segments []types.Segment
	builder  *Builder
}

func NewCollection() *Collection {
	return &Collection{
		builder: NewBuilder(),
	}

}

func (c *Collection) Build(conf types.Config) error {
	segment, err := c.builder.Build(conf)
	if err != nil {
		return err
	}
	c.segments = append(c.segments, segment)

	return nil
}

func (c *Collection) Read(buf *bytes.Buffer) {
	for i := range c.segments {
		c.segments[i].Read(buf)
	}
}

func (c *Collection) Refresh() {
	for i := range c.segments {
		c.segments[i].Refresh()
	}
}
