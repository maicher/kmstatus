package main

import (
	"fmt"
	"os"
	"time"

	"github.io/maicher/stbar/app"
	"github.io/maicher/stbar/pkg/parsers"
	"github.io/maicher/stbar/pkg/parsers/cpu"
	"github.io/maicher/stbar/pkg/parsers/mem"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal error: %s\n", err)
			os.Exit(1)
		}
	}()

	s := &app.System{}
	v := app.NewView("basic.txt.tmpl")
	ch := make(chan any)

	pl := parsers.NewParsersList()
	pl.MustInit(cpu.NewFreqParser, time.Second, true)
	pl.MustInit(cpu.NewLoadParser, time.Second, true)
	pl.MustInit(cpu.NewTempParser, time.Second, false)
	pl.MustInit(mem.NewMemParser, time.Second, false)
	pl.StartParsing(ch)

	for {
		switch data := (<-ch).(type) {
		case cpu.Freq:
			s.CPU.Freq = data
		case cpu.Load:
			s.CPU.Load = data
		case cpu.Temp:
			s.CPU.Temp = data
		case mem.Mem:
			s.Mem = data
		case string:
			fmt.Println(v.Render(s))
		case error:
			fmt.Fprintf(os.Stderr, "%s\n", data)
		}
	}
}
