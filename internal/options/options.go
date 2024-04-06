package options

import (
	"flag"
	"fmt"
)

var help = `NAME
  kmstatus - dynamic status bar

SYNOPSIS
  kmstatus [OPTION...]   // start printing status

DESCRIPTION
  Once started, kmstatus will print output every n seconds,
  where n is the lowest +refreshinterval+ from parser options (default: 1s).
  To trigger an additional refresh run a control command:
    kmstatus -r
 (communicates with the main process via sockets.)

OPTIONS
  --config PATH,    -c   path to a kmstatusrc file
                         if not set, kmstatus will try to look up following paths:
                         $XDG_CONFIG_HOME/kmstatus/kmstatusrc.toml
                         $HOME/.config/kmstatus/kmstatusrc.toml
  --xwindow,        -x   print output to default's window WM_NAME (instead stdout)
                         (to use this option kmstatusatus needs to be build with -tag X)
  --doc                  print documentation
  --version,        -v   print version
  --socketpath,     -s   a custom path to a socket file
                         (default: /tmp/kmstatus.sock)
  --text TEXT,      -t   set text control command
  --text-unset,     -u   unset text control command
  --refresh,        -r   refresh now control command

CONFIG
  See the below link for example config:
    https://github.com/maicher/kmstatus/blob/master/internal/config/kmstatusrc.example.toml
`

type Options struct {
	ConfigPath string
	Doc        bool
	Version    bool
	XWindow    bool
	SocketPath string
	Text       string
	UnsetText  bool
	Refresh    bool

	ControlCmd string
}

func Parse() Options {
	var opts Options

	flag.StringVar(&opts.ConfigPath, "config", "", "")
	flag.StringVar(&opts.ConfigPath, "c", "", "")

	flag.BoolVar(&opts.Doc, "doc", false, "")

	flag.BoolVar(&opts.Version, "version", false, "")
	flag.BoolVar(&opts.Version, "v", false, "")

	flag.BoolVar(&opts.XWindow, "xwindow", false, "")
	flag.BoolVar(&opts.XWindow, "x", false, "")

	// todo change default to: /tmp/kmstatus.$USER.$DISPLAY.sock
	flag.StringVar(&opts.SocketPath, "socketpath", "/tmp/kmstatus.sock", "")
	flag.StringVar(&opts.SocketPath, "s", "/tmp/kmstatus.sock", "")

	flag.StringVar(&opts.Text, "text", "", "")
	flag.StringVar(&opts.Text, "t", "", "")

	flag.BoolVar(&opts.UnsetText, "text-unset", false, "")
	flag.BoolVar(&opts.UnsetText, "u", false, "")

	flag.BoolVar(&opts.Refresh, "refresh", false, "")
	flag.BoolVar(&opts.Refresh, "r", false, "")

	f := flag.CommandLine.Output()
	flag.Usage = func() { fmt.Fprintf(f, help) }
	flag.Parse()

	opts.ControlCmd = opts.buildControlCmd()

	return opts
}

func (opts Options) buildControlCmd() string {
	if opts.Text != "" {
		return opts.Text
	}

	if opts.Refresh {
		return "cmd:refresh"
	}

	if opts.UnsetText {
		return "cmd:unsetText"
	}

	return ""
}
