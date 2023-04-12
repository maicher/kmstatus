package network

import (
	"math"
	"strings"
	"time"
)

type Net struct {
	Interfaces []Interface
}

func (n *Net) ByPrefix(name string) (is []Interface) {
	for _, i := range n.Interfaces {
		if strings.HasPrefix(i.Name, name) {
			is = append(is, i)
		}
	}

	return is
}

func (n *Net) FirstByName(name string) (Interface, bool) {
	for _, i := range n.Interfaces {
		if i.Name == name {
			return i, true
		}
	}

	return Interface{}, false
}

func (n *Net) calculateSpeed(buf map[string]Interface, lastParsedAt time.Time) {
	mul := time.Since(lastParsedAt)

	for index, i := range n.Interfaces {
		n.Interfaces[index].Rx = n.speed(i.RxTotal-buf[i.Name].RxTotal, mul)
		n.Interfaces[index].Tx = n.speed(i.TxTotal-buf[i.Name].TxTotal, mul)
	}
}

func (n *Net) speed(val int, t time.Duration) int {
	return int(math.Round(float64(val) * float64(time.Second) / float64(t)))
}
