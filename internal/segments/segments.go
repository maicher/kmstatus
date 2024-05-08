package segments

import (
	"bytes"
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

type Config struct {
	ParserName      string
	RefreshInterval time.Duration
	Template        string
}

type newSegmentFunc func(string, time.Duration) (types.Segment, error)

type Segments struct {
	segments     []types.Segment
	builderFuncs map[string]newSegmentFunc
}

func New() *Segments {
	return &Segments{
		builderFuncs: map[string]newSegmentFunc{
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

func (s *Segments) AppendNewSegment(conf Config) error {
	segment, err := s.build(conf)
	if err != nil {
		return err
	}
	s.segments = append(s.segments, segment)

	return nil
}

// Read reads data from each segment into the *bytes.Buffer param.
// The data will be formatted according to segments' templates.
func (s *Segments) Read(buf *bytes.Buffer) {
	for i := range s.segments {
		s.segments[i].Read(buf)
	}
}

// Refresh forces segments to get it's data from system files or shell programs.
// It does refresh cpu and network segments, as their data can only be get
// in equal time intervals.
func (s *Segments) Refresh() {
	for i := range s.segments {
		s.segments[i].Refresh()
	}
}

func (s *Segments) build(c Config) (types.Segment, error) {
	newSegmentFunc, ok := s.builderFuncs[c.ParserName]
	if !ok {
		return nil, errors.New("Invalid parser name: " + c.ParserName)
	}

	return newSegmentFunc(c.Template, c.RefreshInterval)
}
