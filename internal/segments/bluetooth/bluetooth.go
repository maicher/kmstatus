package bluetooth

import (
	"bytes"
	"fmt"

	"github.com/maicher/kmst/internal/segments"
)

type Bluetooth struct {
	segments.Segment
	segments.Template

	Data   Data
	Parser *BluetoothParser
}

func New(conf segments.Config) (segments.ParseReader, error) {
	var bt Bluetooth
	var err error

	bt.Parser, err = NewBluetoothParser()
	if err != nil {
		return &bt, err
	}

	bt.Segment = segments.NewSegment(bt.read, bt.parse, conf.ParseInterval)

	err = bt.NewTemplate(conf.StrippedTemplate(), helpers)
	if err != nil {
		return &bt, fmt.Errorf("Unable to parse Bluetooth template: %s", err)
	}

	return &bt, nil
}

func (bt *Bluetooth) read(b *bytes.Buffer) error {
	return bt.Tmpl.Execute(b, bt.Data)
}

func (bt *Bluetooth) parse() error {
	return bt.Parser.Parse(&bt.Data)
}
