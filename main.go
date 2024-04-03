package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maicher/kmst/internal/config"
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
	// $USER $DISPLAY
	socketPath := "/tmp/kmst.sock"
	var text string

	opts := options.Parse()

	if opts.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if opts.Doc {
		fmt.Println(doc)
		os.Exit(0)
	}

	if opts.Text != "" {
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			fmt.Println("Main process is not running: Error connecting to socket:", err)
			os.Exit(1)
		}

		_, err = conn.Write([]byte(opts.Text))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

		conn.Close()
		os.Exit(0)
	}

	if opts.Refresh {
		conn, err := net.Dial("unix", socketPath)
		if err != nil {
			fmt.Println("Error connecting to socket:", err)
			return
		}
		defer conn.Close()

		_, err = conn.Write([]byte("cmd:refresh"))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}

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

	// TODO: check here if socket file exist
	// Initialize segments
	buf := bytes.Buffer{}
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

	buf.WriteString("Starting...")
	view.Flush(&buf)
	time.Sleep(10 * time.Millisecond)

	// Listen
	refresh := make(chan struct{})
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer os.Remove(socketPath)
	defer listener.Close()

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				break
			}

			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if err != nil {
				conn.Close()
				fmt.Println("Error reading:", err)
				return
			}

			payload := string(buffer[:n])

			if payload == "cmd:refresh" {
				refresh <- struct{}{}
			} else {
				text = " " + payload + " "
				refresh <- struct{}{}
			}
			conn.Close()
		}
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

	// Start rendering
	ticker := time.NewTicker(time.Second)

mainLoop:
	for {
		select {
		case <-ticker.C:
			buf.WriteString(text)

			for i := range segments {
				segments[i].Read(&buf)
			}
			view.Flush(&buf)
		case <-refresh:
			buf.WriteString(text)

			for i := range segments {
				segments[i].Refresh()
				segments[i].Read(&buf)
			}
			view.Flush(&buf)
		case <-terminate:
			break mainLoop
		}
	}
}
