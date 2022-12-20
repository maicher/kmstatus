package main

import (
	"fmt"
	"os"

	"github.io/maicher/kmstatus/app"
	"github.io/maicher/kmstatus/app/config"
	"github.io/maicher/kmstatus/app/services"
	"github.io/maicher/kmstatus/app/view"
)

func main() {
	conf := config.Parse()

	if conf.PrintTemplate {
		fmt.Println(view.DefaultTemplate)
		os.Exit(0)
	}

	if conf.PrintConfig {
		for _, p := range conf.ParserConfigs {
			fmt.Printf("%s %s\n", p.Name, p.Interval)
			fmt.Printf("%s-sig %t\n", p.Name, p.OnSig)
		}
		os.Exit(0)
	}

	ch := make(chan any)
	parsePeriodically := services.NewParsePeriodically(ch, conf.ParserConfigs)
	parseOnSig := services.NewParseOnSig(ch, conf.ParserConfigs)
	generate := services.NewGenerate(ch, conf.ParserConfigs)

	v, err := view.New(conf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	c := app.NewController(ch, v)

	// Start
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

	go generate.Loop()

	c.AggregateDataAndRenderView()
}
