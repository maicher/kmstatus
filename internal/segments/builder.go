package segments

import (
	"errors"
	"time"

	"github.com/maicher/kmstatus/internal/segments/audio"
	"github.com/maicher/kmstatus/internal/segments/bluetooth"
	"github.com/maicher/kmstatus/internal/segments/clock"
	"github.com/maicher/kmstatus/internal/segments/cpu"
	"github.com/maicher/kmstatus/internal/segments/mem"
	"github.com/maicher/kmstatus/internal/segments/network"
	"github.com/maicher/kmstatus/internal/segments/processes"
	"github.com/maicher/kmstatus/internal/segments/temperature"
	"github.com/maicher/kmstatus/internal/types"
)

type newSegmentFunc func(string, time.Duration) (types.Segment, error)

type builder struct {
	builders map[string]newSegmentFunc
}

func newBuilder() *builder {
	return &builder{
		builders: map[string]newSegmentFunc{
			"cpu":         cpu.New,
			"temperature": temperature.New,
			"mem":         mem.New,
			"network":     network.New,
			"clock":       clock.New,
			"bluetooth":   bluetooth.New,
			"audio":       audio.New,
			"processes":   processes.New,
		},
	}
}

func (b *builder) Build(c Config) (types.Segment, error) {
	newParserFunc, ok := b.builders[c.ParserName]
	if !ok {
		return nil, errors.New("Invalid parser name: " + c.ParserName)
	}

	return newParserFunc(c.Template, c.RefreshInterval)
}
