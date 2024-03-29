package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/maicher/kmst/internal/config"
	"github.com/maicher/kmst/internal/options"
	"github.com/maicher/kmst/internal/segments"
	"github.com/maicher/kmst/internal/segments/cpu"
	"github.com/maicher/kmst/internal/segments/mem"
	"github.com/maicher/kmst/internal/segments/temperature"
	"github.com/maicher/kmst/internal/ui"
)

var version string

//go:embed doc.txt
var doc string

type NewSegmentFunc func(segments.Config) (segments.Reader, error)

type SegmentsBuilder struct {
	builders map[string]NewSegmentFunc
}

func NewSegmentsBuilder() *SegmentsBuilder {
	return &SegmentsBuilder{
		builders: map[string]NewSegmentFunc{
			"cpu":         cpu.New,
			"temperature": temperature.New,
			"mem":         mem.New,
		},
	}
}

func (b *SegmentsBuilder) New(c segments.Config) (segments.Reader, error) {
	newParserFunc, ok := b.builders[c.ParserName]
	if !ok {
		return nil, errors.New("Invalid parser name: " + c.ParserName)
	}

	return newParserFunc(c)
}

func main() {
	opts := options.Parse()

	if opts.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if opts.Doc {
		fmt.Println(doc)
		os.Exit(0)
	}

	c, err := config.New(opts.ConfigPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	view, err := ui.NewView(opts.XWindow)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	buf := bytes.Buffer{}
	segmentsBuilder := NewSegmentsBuilder()
	var segment segments.Reader
	var segments []segments.Reader

	for _, p := range c.Segments {
		segment, err = segmentsBuilder.New(p)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}

		segments = append(segments, segment)
	}

	buf.WriteString("Starting...")
	view.Flush(&buf)
	time.Sleep(100 * time.Millisecond)

	for {
		for i := range segments {
			segments[i].Read(&buf)
		}
		view.Flush(&buf)
		time.Sleep(time.Second)
	}
}
