package bluetooth

import (
	"bytes"
	"fmt"

	"github.com/maicher/kmst/internal/segments/common"
	"github.com/maicher/kmst/internal/types"
)

type Bluetooth struct {
	common.PeriodicParser
	common.Template

	Data   Data
	Parser *BluetoothParser
}

func New(conf types.Config) (types.Segment, error) {
	var bt Bluetooth
	var err error

	bt.Parser, err = NewBluetoothParser()
	if err != nil {
		return &bt, err
	}

	bt.PeriodicParser = common.NewPeriodicParser(bt.read, bt.parse, conf.RefreshInterval)

	err = bt.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &bt, fmt.Errorf("Unable to parse Bluetooth template: %s", err)
	}

	return &bt, nil
}

func (bt *Bluetooth) Refresh() {
	bt.PeriodicParser.Parse()
}

func (bt *Bluetooth) read(b *bytes.Buffer) error {
	return bt.Tmpl.Execute(b, bt.Data)
}

func (bt *Bluetooth) parse() error {
	return bt.Parser.Parse(&bt.Data)
}
