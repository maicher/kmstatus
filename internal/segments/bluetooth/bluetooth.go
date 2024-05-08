package bluetooth

import (
	"bytes"
	"fmt"
	"time"

	"github.com/maicher/kmstatus/internal/segments/common"
	"github.com/maicher/kmstatus/internal/types"
)

type Bluetooth struct {
	common.PeriodicParser
	common.Template

	data   data
	parser *Parser
}

func New(tmpl string, refreshInterval time.Duration) (types.Segment, error) {
	var bt Bluetooth

	bt.PeriodicParser = common.NewPeriodicParser(bt.read, bt.parse, refreshInterval)

	err := bt.NewTemplate(tmpl, helpers)
	if err != nil {
		return &bt, fmt.Errorf("unable to parse Bluetooth template: %s", err)
	}

	return &bt, nil
}

func (bt *Bluetooth) Refresh() {
	bt.PeriodicParser.Parse()
}

func (bt *Bluetooth) read(b *bytes.Buffer) error {
	return bt.Tmpl.Execute(b, bt.data)
}

func (bt *Bluetooth) parse() error {
	return bt.parser.Parse(&bt.data)
}
