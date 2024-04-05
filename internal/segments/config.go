package segments

import "time"

type Config struct {
	ParserName      string
	RefreshInterval time.Duration
	Template        string
}
