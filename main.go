package main

import (
	"fmt"
	"os"

	"github.com/maicher/kmst/internal/options"
	"github.com/maicher/kmst/internal/ui"
)

var version string

func main() {
	var err error
	opts := options.Parse()

	if opts.Version {
		fmt.Println(version)
		os.Exit(0)
	}

	if opts.Doc {
		fmt.Println(doc)
		os.Exit(0)
	}

	c, err := NewConfig(opts.ConfigPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%+v\n", c)

	view, err := ui.NewView(opts.XWindow)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	view.Render("test")
}
