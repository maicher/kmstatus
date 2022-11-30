package main

import (
	"fmt"
	"time"

	"github.io/maicher/stbar/pkg/parsers/cpu"
)

func main() {
	parser := cpu.NewParser()
	_, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(200 * time.Millisecond)

	c, err := parser.Parse()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("CPU: %f%%", c.Load())
}
