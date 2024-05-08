package network

type data struct {
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
