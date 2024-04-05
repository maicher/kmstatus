package types

import (
	"bytes"
	"strings"
	"time"
)

type Segment interface {
	Refresh()
	Read(*bytes.Buffer)
}

type Config struct {
	ParserName      string
	RefreshInterval time.Duration
	Template        string
}

func (c Config) StrippedTemplate() string {
	return strings.ReplaceAll(c.Template, "\n", "")
}
