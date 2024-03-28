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
	"github.com/maicher/kmst/internal/segments/temperature"
	"github.com/maicher/kmst/internal/ui"
)

var version string

//go:embed doc.txt
var doc string

var constructors = map[string]segments.NewSegmentFunc{
	"cpu":         cpu.New,
	"temperature": temperature.New,
}

func NewSegment(c segments.SegmentConfig) (segments.Segment, error) {
	newParserFunc, ok := constructors[c.ParserName]
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
	var component segments.Segment
	var components []segments.Segment

	for _, p := range c.Segments {
		component, err = NewSegment(p)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}

		components = append(components, component)
	}

	buf.WriteString("Starting...")
	view.Flush(&buf)
	time.Sleep(100 * time.Millisecond)

	for {
		for i := range components {
			components[i].Read(&buf)
		}
		view.Flush(&buf)
		time.Sleep(time.Second)
	}
}
