package segments

import (
	"strings"
	"time"
)

type Config struct {
	ParserName    string
	ParseInterval time.Duration
	ParseOnSig    bool
	Template      string
}

func (c Config) StrippedTemplate() string {
	return strings.ReplaceAll(c.Template, "\n", "")
}
