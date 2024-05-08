package network

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maicher/kmstatus/internal/segments/common"
	"github.com/maicher/kmstatus/internal/types"
)

type Network struct {
	common.PeriodicParser
	common.Template

	data   []data
	parser *NetworkParser
}

func New(tmpl string, refreshInterval time.Duration) (types.Segment, error) {
	var n Network
	var err error

	n.data = make([]data, 10)

	n.parser, err = NewNetworkParser()
	if err != nil {
		return &n, err
	}

	err = n.NewTemplate(tmpl, helpers)
	if err != nil {
		return &n, fmt.Errorf("Unable to parse Network template: %s", err)
	}

	n.PeriodicParser = common.NewPeriodicParser(n.read, n.parse, refreshInterval)

	return &n, nil
}

func (n *Network) Refresh() {
}

func (n *Network) read(b *bytes.Buffer) error {
	var err error

	for i := range n.data {
		err = n.Tmpl.Execute(b, n.data[i])
		if err != nil {
			break
		}
	}

	return err
}

func (n *Network) parse() error {
	return n.parser.Parse(n.data)
}
