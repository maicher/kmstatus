package network

import (
	"fmt"
	"text/template"

	"github.com/maicher/kmst/internal/segments"
)

type Network struct {
	segments.Segment

	Data   []Data
	Parser *NetworkParser
}

func New(conf segments.Config) (segments.Reader, error) {
	var n Network
	var err error

	n.Segment = segments.NewSegment(n.handleMsg)
	n.Data = make([]Data, 10)

	n.Parser, err = NewNetworkParser()
	if err != nil {
		return &n, err
	}

	n.Template, err = template.New("").Funcs(helpers).Parse(conf.StrippedTemplate())
	if err != nil {
		return &n, fmt.Errorf("Unable to parse Network template: %s", err)
	}

	go n.OnTick(conf.ParseInterval, func() {
		n.MsgQueue <- segments.ParseMsg{}
	})

	return &n, nil
}

func (n *Network) handleMsg(msg any) error {
	var err error

	switch msg := msg.(type) {
	case segments.ReadMsg:
		for i := range n.Data {
			err = n.Template.Execute(msg.Buffer, n.Data[i])
			if err != nil {
				break
			}
		}

		n.Sync <- struct{}{}
	case segments.ParseMsg:
		err = n.Parser.Parse(n.Data)
	default:
		panic("Invalid message")
	}

	if err != nil {
		return err
	}

	return nil
}
