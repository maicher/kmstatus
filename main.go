package main

import (
	"fmt"
	"os"

	"github.io/maicher/kmstatus/app"
	"github.io/maicher/kmstatus/app/config"
	"github.io/maicher/kmstatus/app/view"
)

func main() {
	conf := config.Parse()

	if conf.PrintTemplate {
		fmt.Println(view.DefaultTemplate)
		os.Exit(0)
	}

	if conf.PrintConfig {
		fmt.Printf("%+v\n", conf)
		os.Exit(0)
	}

	c, err := app.NewController(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c.Loop()
}
