package services

import (
	"time"

	"github.io/maicher/kmstatus/pkg/parsers"
)

func onTick(interval time.Duration, f func()) {
	f()

	ticker := time.NewTicker(interval)
	for range ticker.C {
		f()
	}
}

func parse(p parsers.Parser, ch chan<- any) {
	v, err := p.Parse()
	if err != nil {
		ch <- err
	}

	ch <- v
}
