package network

type NetworkParser struct {
}

func NewNetworkParser() (*NetworkParser, error) {
	var n NetworkParser

	return &n, nil

}

func (n *NetworkParser) Parse(d *Data) error {
	d.Name = "test"

	return nil
}
