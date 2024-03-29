package segments

import (
	"fmt"
	"text/template"
	"time"
)

type HandleMsgFunc func(any) error

type Segment struct {
	MsgQueue chan any
	Sync     chan any

	Template *template.Template
}

func NewSegment(handleMsg HandleMsgFunc) (s Segment) {
	s.MsgQueue = make(chan any)
	s.Sync = make(chan any)

	go s.Loop(handleMsg)

	return s
}

func (a Segment) Loop(f HandleMsgFunc) {
	var err error

	for msg := range a.MsgQueue {
		err = f(msg)

		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}

func (a *Segment) OnTick(interval time.Duration, f func()) {
	f()

	ticker := time.NewTicker(interval)
	for range ticker.C {
		f()
	}
}
