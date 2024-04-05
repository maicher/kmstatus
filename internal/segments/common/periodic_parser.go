package common

import (
	"bytes"
	"fmt"
	"time"
)

type readMsg struct{ Buffer *bytes.Buffer }
type parseMsg struct{}

type handleReadMsgFunc func(*bytes.Buffer) error
type handleParseMsgFunc func() error

type PeriodicParser struct {
	msgQueue chan any
	sync     chan any
}

func NewPeriodicParser(r handleReadMsgFunc, p handleParseMsgFunc, refreshInterval time.Duration) (s PeriodicParser) {
	s.msgQueue = make(chan any)
	s.sync = make(chan any)

	go s.Loop(r, p)

	go s.onTick(refreshInterval, func() {
		s.msgQueue <- parseMsg{}
		<-s.sync
	})

	return s
}

func (s PeriodicParser) Loop(r handleReadMsgFunc, p handleParseMsgFunc) {
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

func (s *PeriodicParser) Read(b *bytes.Buffer) {
	s.msgQueue <- readMsg{Buffer: b}
	<-s.sync
}

func (s *PeriodicParser) Parse() {
	s.msgQueue <- parseMsg{}
	<-s.sync
}

func (s *PeriodicParser) onTick(interval time.Duration, f func()) {
	f()

	ticker := time.NewTicker(interval)
	for range ticker.C {
		f()
	}
}
