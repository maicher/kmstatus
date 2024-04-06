package segments

import (
	"bytes"

	"github.com/maicher/kmstatus/internal/types"
)

type Collection struct {
	segments []types.Segment
	builder  *builder
}

func New() *Collection {
	return &Collection{
		builder: newBuilder(),
	}
}

func (c *Collection) AppendNewSegment(conf Config) error {
	segment, err := c.builder.Build(conf)
	if err != nil {
		return err
	}
	c.segments = append(c.segments, segment)

	return nil
}

// Read reads data from each segment into the *bytes.Buffer param.
// The data will be formatted according to segments' templates.
func (c *Collection) Read(buf *bytes.Buffer) {
	for i := range c.segments {
		c.segments[i].Read(buf)
	}
}

// Refresh forces segments to get it's data from system files or shell programs.
// It does refresh cpu and network segments, as their data can only be get
// in equal time intervals.
func (c *Collection) Refresh() {
	for i := range c.segments {
		c.segments[i].Refresh()
	}
}
