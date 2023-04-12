package network

import (
	"os"
	"strings"
)

type Interface struct {
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

func (i *Interface) IsUp() bool {
	s, _ := os.ReadFile("/sys/class/net/" + i.Name + "/operstate")

	return strings.TrimSpace(string(s)) == "up"
}
