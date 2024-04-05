package segments

import (
	"errors"

	"github.com/maicher/kmst/internal/segments/audio"
	"github.com/maicher/kmst/internal/segments/bluetooth"
	"github.com/maicher/kmst/internal/segments/clock"
	"github.com/maicher/kmst/internal/segments/cpu"
	"github.com/maicher/kmst/internal/segments/mem"
	"github.com/maicher/kmst/internal/segments/network"
	"github.com/maicher/kmst/internal/segments/processes"
	"github.com/maicher/kmst/internal/segments/temperature"
	"github.com/maicher/kmst/internal/types"
)

type newSegmentFunc func(types.Config) (types.Segment, error)

type Builder struct {
	builders map[string]newSegmentFunc
}

func NewBuilder() *Builder {
	return &Builder{
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

func (b *Builder) Build(c types.Config) (types.Segment, error) {
	newParserFunc, ok := b.builders[c.ParserName]
	if !ok {
		return nil, errors.New("Invalid parser name: " + c.ParserName)
	}

	return newParserFunc(c)
}
