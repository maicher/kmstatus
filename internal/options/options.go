package options

import (
	"flag"
	"fmt"
)

var help = `NAME
  kmst - km status bar

SYNOPSIS
  kmst [OPTION...]

OPTIONS
  --config PATH,    -c   path to a kmstrc file
                         if not set, kmst will try to look up following paths:
                         $XDG_CONFIG_HOME/kmst/kmstrc.toml
                         $HOME/.config/kmst/kmstrc.toml
  --xwindow,        -x   print output to default's window WM_NAME (instead stdout)
                         (to use this option kmstatus needs to be build with -tag X)
  --doc                  print documentation
  --version,        -v   print version

CONFIG
  See the below link for example config:
    https://github.com/maicher/kmst/blob/master/internal/config/kmstrc.example.toml

SIGNALS
  By default kmst will print to output every n seconds,
  where n is the lowest interval from parser options (default 1s).
  A signal can be sent to kmst to trigger an additional refresh:
    pkill -USR1 kmst$
`

type Options struct {
	ConfigPath string
	Doc        bool
	Version    bool
	XWindow    bool
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

	f := flag.CommandLine.Output()
	flag.Usage = func() { fmt.Fprintf(f, help) }
	flag.Parse()

	return opts
}
