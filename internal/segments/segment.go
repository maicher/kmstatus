package segments

import (
	"bytes"
	"fmt"
	"time"
)

type readMsg struct{ Buffer *bytes.Buffer }
type parseMsg struct{}

type handleReadMsgFunc func(*bytes.Buffer) error
type handleParseMsgFunc func() error

type Segment struct {
	msgQueue chan any
	sync     chan any
}

func NewSegment(r handleReadMsgFunc, p handleParseMsgFunc, parseInterval time.Duration) (s Segment) {
	s.msgQueue = make(chan any)
	s.sync = make(chan any)

	go s.Loop(r, p)

	go s.onTick(parseInterval, func() {
		s.msgQueue <- parseMsg{}
		<-s.sync
	})

	return s
}

func (s Segment) Loop(r handleReadMsgFunc, p handleParseMsgFunc) {
	var err error

	for msg := range s.msgQueue {
		switch msg := msg.(type) {
		case readMsg:
			err = r(msg.Buffer)
			s.sync <- struct{}{}
		case parseMsg:
			err = p()
			s.sync <- struct{}{}
		default:
			panic("Invalid message")
		}

		if err != nil {
			fmt.Printf("%+v\n", err)
		}
	}
}

func (s *Segment) Read(b *bytes.Buffer) {
	s.msgQueue <- readMsg{Buffer: b}
	<-s.sync
}

func (s *Segment) Parse() {
	s.msgQueue <- parseMsg{}
	<-s.sync
}

func (s *Segment) onTick(interval time.Duration, f func()) {
	f()

	ticker := time.NewTicker(interval)
	for range ticker.C {
		f()
	}
}
