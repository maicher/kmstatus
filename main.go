package main

import (
	"fmt"
	"os"

	"github.io/maicher/stbar/app"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("Fatal error: %s\n", err)
			os.Exit(1)
		}
	}()

	c := app.NewController()
	c.Loop()
}
