package main

import (
	"bytes"
	_ "embed"

	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/maicher/kmstatus/internal/config"
	"github.com/maicher/kmstatus/internal/ipc"
	"github.com/maicher/kmstatus/internal/options"
	"github.com/maicher/kmstatus/internal/segments"
	"github.com/maicher/kmstatus/internal/ui"
)

var version string

//go:embed doc.txt
var doc string

//go:embed internal/config/kmstatusrc.example.toml
var kmstatusrcExample string

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
	c, err := config.New(opts.ConfigPath, kmstatusrcExample)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Initialize segments
	segmentsCollection := segments.New()
	for _, p := range c.Segments {
		err = segmentsCollection.AppendNewSegment(p)
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
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
	time.Sleep(50 * time.Millisecond)

	// Listen
	// check if socket file already exists
	refresh := make(chan struct{})
	render := make(chan struct{})
	textCh := make(chan string)
	go func() {
		err := ipc.Listen(func(cmd string) {
			switch cmd {
			case "cmd:refresh":
				refresh <- struct{}{}
			case "cmd:unsetText":
				textCh <- ""
			default:
				textCh <- " " + cmd + " "
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

	ticker := time.NewTicker(c.MinInterval())
	go func() {
		render <- struct{}{}
		for range ticker.C {
			render <- struct{}{}
		}
	}()

mainLoop:
	for {
		select {
		case text = <-textCh:
		case <-render:
			buf.WriteString(text)
			segmentsCollection.Read(&buf)
			view.Flush(&buf)
		case <-refresh:
			segmentsCollection.Refresh()
		case <-terminate:
			break mainLoop
		}
	}
}
