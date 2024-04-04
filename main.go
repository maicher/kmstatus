package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maicher/kmst/internal/config"
	"github.com/maicher/kmst/internal/ipc"
	"github.com/maicher/kmst/internal/options"
	"github.com/maicher/kmst/internal/segments"
	"github.com/maicher/kmst/internal/segments/audio"
	"github.com/maicher/kmst/internal/segments/bluetooth"
	"github.com/maicher/kmst/internal/segments/clock"
	"github.com/maicher/kmst/internal/segments/cpu"
	"github.com/maicher/kmst/internal/segments/mem"
	"github.com/maicher/kmst/internal/segments/network"
	"github.com/maicher/kmst/internal/segments/processes"
	"github.com/maicher/kmst/internal/segments/temperature"
	"github.com/maicher/kmst/internal/ui"
)

var version string

//go:embed doc.txt
var doc string

type NewSegmentFunc func(segments.Config) (segments.RefreshReader, error)

type SegmentsBuilder struct {
	builders map[string]NewSegmentFunc
}

func NewSegmentsBuilder() *SegmentsBuilder {
	return &SegmentsBuilder{
		builders: map[string]NewSegmentFunc{
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

func (b *SegmentsBuilder) New(c segments.Config) (segments.RefreshReader, error) {
	newParserFunc, ok := b.builders[c.ParserName]
	if !ok {
		return nil, errors.New("Invalid parser name: " + c.ParserName)
	}

	return newParserFunc(c)
}

func main() {
	var text string
	opts := options.Parse()

	// Print version and exit.
	if opts.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	// Print docs and exit.
	if opts.Doc {
		fmt.Println(doc)
		os.Exit(0)
	}

	ipc := ipc.IPC{
		SocketPath: opts.SocketPath,
	}

	// Send control command to already running main process and exit.
	if opts.ControlCmd != "" {
		err := ipc.Send(opts.ControlCmd)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	// Init main process.
	c, err := config.New(opts.ConfigPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize segments
	segmentsBuilder := NewSegmentsBuilder()
	var segment segments.RefreshReader
	var segments []segments.RefreshReader

	for _, p := range c.Segments {
		segment, err = segmentsBuilder.New(p)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}

		segments = append(segments, segment)
	}

	// Initialize view and write start message.
	view, err := ui.NewView(opts.XWindow)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	buf := bytes.Buffer{}
	buf.WriteString("Starting...")
	view.Flush(&buf)
	time.Sleep(10 * time.Millisecond)

	// Listen
	// check if socket file already exist and display error
	refresh := make(chan struct{})
	render := make(chan struct{})
	go func() {
		err := ipc.Listen(func(cmd string) {
			switch cmd {
			case "cmd:refresh":
				refresh <- struct{}{}
			case "cmd:unsetText":
				text = ""
			default:
				text = " " + cmd + " "
			}

			render <- struct{}{}
		})

		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(2)
		}
	}()
	defer ipc.Close()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(time.Second)
	go func() {
		for range ticker.C {
			render <- struct{}{}
		}
	}()

mainLoop:
	for {
		select {
		case <-render:
			buf.WriteString(text)

			for i := range segments {
				segments[i].Read(&buf)
			}

			view.Flush(&buf)
		case <-refresh:
			for i := range segments {
				segments[i].Refresh()
			}
		case <-terminate:
			break mainLoop
		}
	}
}
