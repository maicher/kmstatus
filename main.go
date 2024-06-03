package main

import (
	"bytes"
	_ "embed"
	"os/signal"
	"syscall"
	"time"

	"fmt"
	"os"

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
	var err error

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

	if opts.ControlCmd == "" {
		err = runMainProcess(opts.ConfigPath, opts.SocketPath, opts.XWindow)
	} else {
		err = sendCmdToMainProcess(opts.ControlCmd, opts.SocketPath)
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runMainProcess(configPath, socketPath string, xWindow bool) error {
	c, err := config.New(configPath, kmstatusrcExample)
	if err != nil {
		return err
	}

	ipc := ipc.IPC{SocketPath: socketPath}
	defer ipc.CloseListener()

	// Initialize segments
	segmentsCollection := segments.New()
	for _, p := range c.Segments {
		err = segmentsCollection.AppendNewSegment(p)
		if err != nil {
			return err
		}
	}

	// Initialize view and write start message.
	view, err := ui.NewView(xWindow)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	buf.WriteString("Starting...")
	view.Flush(&buf)
	time.Sleep(50 * time.Millisecond)

	// Listen
	refresh := make(chan struct{})
	render := make(chan struct{})
	textCh := make(chan string)
	errCh := make(chan error)
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
			errCh <- err
			return
		}
	}()

	// Main loop
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(c.MinInterval())
	go func() {
		render <- struct{}{}
		for range ticker.C {
			render <- struct{}{}
		}
	}()

	var text string
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
		case err = <-errCh:
			fmt.Println(err)
			break mainLoop
		case <-terminate:
			break mainLoop
		}
	}

	return nil
}

func sendCmdToMainProcess(cmd, socketPath string) error {
	ipc := ipc.IPC{SocketPath: socketPath}

	return ipc.Send(cmd)
}
