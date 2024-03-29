package segments

import (
	"bytes"
	"fmt"
	"text/template"
	"time"
)

type ReadMsg struct{ Buffer *bytes.Buffer }
type ParseMsg struct{}

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

func (s Segment) Loop(f HandleMsgFunc) {
	var err error

	for msg := range s.MsgQueue {
		err = f(msg)

		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}

func (s *Segment) Read(b *bytes.Buffer) {
	s.MsgQueue <- ReadMsg{Buffer: b}
	<-s.Sync
}

func (s *Segment) OnTick(interval time.Duration, f func()) {
	f()

	ticker := time.NewTicker(interval)
	for range ticker.C {
		f()
	}
}
