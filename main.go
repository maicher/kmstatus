package main

import (
	"fmt"
	"os"

	"github.io/maicher/kmstatus/internal/config"
	"github.io/maicher/kmstatus/internal/controller"
	"github.io/maicher/kmstatus/internal/services"
	"github.io/maicher/kmstatus/internal/view"
)

func main() {
	conf := config.Parse()

	if conf.PrintTemplate {
		fmt.Println(string(config.DefaultTemplate))
		os.Exit(0)
	}

	if conf.PrintConfig {
		fmt.Printf("%s %t\n", "xwindow", conf.XWindow)
		for _, p := range conf.ParsersSettings {
			fmt.Printf("%s %s\n", p.Name, p.Interval)
			fmt.Printf("%s-sig %t\n", p.Name, p.OnSig)
		}
		os.Exit(0)
	}

	// Init.
	ch := make(chan any)
	parsePeriodically := services.NewParsePeriodically(ch, conf.ParsersSettings)
	parseOnSig := services.NewParseOnSig(ch, conf.ParsersSettings)
	generate := services.NewGenerateTicks(ch, conf.ParsersSettings)
	v, err := view.New(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	controller := controller.NewController(ch, v)

	// Start.
	err = parsePeriodically.Loop()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = parseOnSig.Loop()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	generate.GenerateTicks()

	controller.AggregateDataAndRenderView()
}
