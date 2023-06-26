package bluetooth

import "github.io/maicher/kmstatus/pkg/parsers"

type BluetoothParser struct {
}

func (p *BluetoothParser) Parse() (any, error) {
	b := Bluetooth{}

	return b, nil
}

func NewBluetothParser() (parsers.Parser, error) {
	parser := BluetoothParser{}

	return &parser, nil
}
