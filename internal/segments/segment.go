package segments

import (
	"bytes"
	"strings"
	"time"
)

type NewSegmentFunc func(SegmentConfig) (Segment, error)

type Segment interface {
	Read(*bytes.Buffer)
}

type SegmentConfig struct {
	ParserName    string
	ParseInterval time.Duration
	ParseOnSig    bool
	Template      string
}

func (c SegmentConfig) StrippedTemplate() string {
	return strings.ReplaceAll(c.Template, "\n", "")
}
