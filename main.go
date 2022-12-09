package main

import (
	"github.io/maicher/stbar/app"
)

func main() {
	c := app.NewController()
	c.Loop()
}
