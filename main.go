package main

import (
	"fmt"
	"os"
	"time"

	"github.io/maicher/stbar/pkg/parsers/cpu"
	"github.io/maicher/stbar/pkg/parsers/mem"
)

func main() {
	var (
		f   cpu.Freq
		l   cpu.Load
		t   cpu.Temp
		m   mem.Mem
		err error
	)

	ch := make(chan any)

	loadParser, err := cpu.NewLoadParser()
	if err != nil {
		fmt.Println(err)
		return
	}

	freqParser, err := cpu.NewFreqParser()
	if err != nil {
		fmt.Println(err)
		return
	}

	tempParser, err := cpu.NewTempParser()
	if err != nil {
		fmt.Println(err)
		return
	}

	memParser, err := mem.NewMemParser()
	if err != nil {
		fmt.Println(err)
		return
	}

	go onTick(time.Second, func() {
		loadParser.Run(ch)
	})
	go onTick(time.Second, func() {
		freqParser.Run(ch)
	})
	go onTick(time.Second, func() {
		tempParser.Run(ch)
	})
	go onTick(time.Second, func() {
		memParser.Run(ch)
	})
	go onTick(time.Second, func() {
		ch <- "render"
	})

	for {
		switch data := (<-ch).(type) {
		case cpu.Freq:
			f = data
		case cpu.Load:
			l = data
		case cpu.Temp:
			t = data
		case mem.Mem:
			m = data
		case string:
			fmt.Printf("%dMHz | ", f)
			fmt.Printf("%0.1f%% | ", l)
			fmt.Printf("%v | ", t)
			fmt.Printf("M: %d(%d) Swap: %d(%d)\n", m.MemUsed, m.MemTotal, m.SwapUsed, m.SwapTotal)
		case error:
			err = data
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
	}
}

func onTick(interval time.Duration, f func()) {
	f()

	ticker := time.NewTicker(interval)
	for range ticker.C {

		f()
	}
}
