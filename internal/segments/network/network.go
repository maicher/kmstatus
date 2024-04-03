package network

import (
	"bytes"
	"fmt"

	"github.com/maicher/kmst/internal/segments"
)

type Network struct {
	segments.Segment
	segments.Template

	Data   []Data
	Parser *NetworkParser
}

func New(conf segments.Config) (segments.RefreshReader, error) {
	var n Network
	var err error

	n.Data = make([]Data, 10)

	n.Parser, err = NewNetworkParser()
	if err != nil {
		return &n, err
	}

	err = n.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &n, fmt.Errorf("Unable to parse Network template: %s", err)
	}

	n.Segment = segments.NewSegment(n.read, n.parse, conf.RefreshInterval)

	return &n, nil
}

func (n *Network) Refresh() {
}

func (n *Network) read(b *bytes.Buffer) error {
	var err error

	for i := range n.Data {
		err = n.Tmpl.Execute(b, n.Data[i])
		if err != nil {
			break
		}
	}

	return err
}

func (n *Network) parse() error {
	return n.Parser.Parse(n.Data)
}
