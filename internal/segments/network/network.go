package network

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/maicher/kmst/internal/segments"
)

type parse struct{}
type read struct{ buffer *bytes.Buffer }

type Network struct {
	segments.Segment

	Data []Data

	Parser *NetworkParser
}

func New(conf segments.Config) (segments.Reader, error) {
	var n Network
	var err error

	n.Segment = segments.NewSegment(n.handleMsg)
	n.Data = make([]Data, 2)

	n.Parser, err = NewNetworkParser()
	if err != nil {
		return &n, err
	}

	n.Template, err = template.New("").Parse(conf.StrippedTemplate())
	if err != nil {
		return &n, fmt.Errorf("Unable to parse Network template: %s", err)
	}

	go n.OnTick(conf.ParseInterval, func() {
		n.MsgQueue <- parse{}
	})

	return &n, nil
}

func (n *Network) Read(b *bytes.Buffer) {
	n.MsgQueue <- read{buffer: b}
	<-n.Sync
}

func (n *Network) handleMsg(msg any) error {
	var err error

	switch msg := msg.(type) {
	case read:
		for i := range n.Data {
			err = n.Template.Execute(msg.buffer, n.Data[i])
			if err != nil {
				break
			}
		}

		n.Sync <- struct{}{}
	case parse:
		for i := range n.Data {
			err = n.Parser.Parse(&n.Data[i])
		}
	default:
		panic("Invalid message")
	}

	if err != nil {
		return err
	}

	return nil
}
