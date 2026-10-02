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

	err := bt.NewTemplate(tmpl, helpers)
	if err != nil {
		return &bt, fmt.Errorf("unable to parse Bluetooth template: %s", err)
	}

	// Start parsing once the segment is fully initialized.
	bt.PeriodicParser = common.NewPeriodicParser(bt.read, bt.parse, refreshInterval)

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
