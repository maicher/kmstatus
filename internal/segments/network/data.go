package network

import (
	"os"
	"strings"
)

type Data struct {
	Name string

	// Total received data in bytes
	RxTotal int

	// Total transmitted data in bytes
	TxTotal int

	// Received speed in bytes per second
	Rx int

	// Transmitted speed in bytes per second
	Tx int
}

func (d *Data) IsUp() bool {
	s, _ := os.ReadFile("/sys/class/net/" + d.Name + "/operstate")

	return strings.TrimSpace(string(s)) == "up"
}
