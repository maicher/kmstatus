package main

import (
	_ "embed"
	"os/signal"
	"syscall"
	"time"

	"fmt"
	"os"

	"github.com/maicher/kmstatus/internal/config"
	"github.com/maicher/kmstatus/internal/ipc"
	"github.com/maicher/kmstatus/internal/kmstatus"
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
	segs := segments.New()
	for _, p := range c.Segments {
		err = segs.AppendNewSegment(p)
		if err != nil {
			return err
		}
	}

	// Initialize view.
	view, err := ui.NewView(xWindow)
	if err != nil {
		return err
	}

	// Initialize.
	kmst := kmstatus.New(view, segs)
	defer kmst.Terminate()

	kmst.SetGreeting("Starting...")
	time.Sleep(50 * time.Millisecond)
	kmst.Render() // Render for the first time

	errCh := make(chan error, 1)
	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, syscall.SIGINT, syscall.SIGTERM)

	// Listen.
	go func() {
		err := ipc.Listen(func(cmd string) {
			switch cmd {
			case "cmd:refresh":
				kmst.Refresh()
			case "cmd:unsetText":
				kmst.SetText("")
			default:
				kmst.SetText(" " + cmd + " ")
			}

			kmst.Render() // Render after receiving IPC command
		})

		if err != nil {
			errCh <- err

			return
		}
	}()

	// Main loop.
	ticker := time.NewTicker(c.MinInterval())
	for {
		select {
		case <-ticker.C:
			kmst.Render() // Render periodically
		case err = <-errCh:
			return err
		case <-terminate:
			return nil
		}
	}
}

func sendCmdToMainProcess(cmd, socketPath string) error {
	ipc := ipc.IPC{SocketPath: socketPath}

	return ipc.Send(cmd)
}
