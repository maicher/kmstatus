package main

import (
	"fmt"
	"time"

	"github.io/maicher/stbar/pkg/parsers/cpu"
)

type CPU struct {
	Load cpu.Load
	Freq cpu.Freq
	Temp cpu.Temp
}

func main() {
	c := CPU{}

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

	_, err = loadParser.Parse()
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(500 * time.Millisecond)

	c.Load, err = loadParser.Parse()
	if err != nil {
		fmt.Println(err)
	}

	c.Freq, err = freqParser.Parse()
	if err != nil {
		fmt.Println(err)
	}

	c.Temp, err = tempParser.Parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("CPU: %0.1f%%\n", c.Load)
	fmt.Printf("Freq: %dMHz\n", c.Freq)
	fmt.Printf("Temp: %v\n", c.Temp)
}
