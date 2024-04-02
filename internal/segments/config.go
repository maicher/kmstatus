package segments

import (
	"strings"
	"time"
)

type Config struct {
	ParserName      string
	RefreshInterval time.Duration
	Template        string
}

func (c Config) StrippedTemplate() string {
	return strings.ReplaceAll(c.Template, "\n", "")
}
