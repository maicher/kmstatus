package main

import (
	"github.io/maicher/kmstatus/app"
)

func main() {
	c := app.NewController()
	c.Loop()
}
