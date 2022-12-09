package app

import (
	"time"

	"github.io/maicher/kmstatus/pkg/parsers"
)

// periodicParser wraps a parser with information about the parsing interval.
type periodicParser struct {
	parser   parsers.Parser
	interval time.Duration
}
