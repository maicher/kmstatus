package main

import (
	"fmt"
	"time"

	"github.io/maicher/stbar/pkg/parsers/cpu"
	"github.io/maicher/stbar/pkg/parsers/temp"
)

func main() {
	parser, err := cpu.NewParser()
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = parser.Parse()
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(100 * time.Millisecond)

	c, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("CPU: %f%%\n", c.Load())

	p, err := temp.NewParser()
	if err != nil {
		fmt.Println(err)
		return
	}

	t, err := p.Parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Temp: %v\n", t.Values)
}
