package network

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const fileWithNetworkInfo = "/proc/net/dev"

type Parser struct {
	dataBuf  map[string]data
	parsedAt time.Time
	file     *os.File
}

func NewParser() (*Parser, error) {
	var n Parser
	var err error

	n.dataBuf = make(map[string]data)
	n.file, err = os.Open(fileWithNetworkInfo)
	if err != nil {
		return &n, fmt.Errorf("Network parser: %s", err)
	}

	return &n, nil
}

// Parse replaces the content of ifaces with the interfaces currently present in the system.
func (n *Parser) Parse(ifaces *[]data) error {
	_, err := n.file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("Network parser: %s", err)
	}

	s := bufio.NewScanner(n.file)
	s.Split(bufio.ScanLines)

	// Drop first 2 lines since they contain headers
	s.Scan()
	s.Scan()

	*ifaces = (*ifaces)[:0]
	for s.Scan() {
		d, ok := parseLine(s.Text())
		if ok {
			*ifaces = append(*ifaces, d)
		}
	}

	n.calculateSpeed(*ifaces)
	n.parsedAt = time.Now()

	// Buffer to calculate speed in the next cycle.
	// Interfaces which disappeared are dropped.
	clear(n.dataBuf)
	for _, d := range *ifaces {
		n.dataBuf[d.Name] = d
	}

	return nil
}

// parseLine parses a line of /proc/net/dev:
// name: rx_bytes rx_packets rx_errs rx_drop rx_fifo rx_frame rx_compressed rx_multicast tx_bytes ...
func parseLine(line string) (d data, ok bool) {
	name, counters, found := strings.Cut(line, ":")
	if !found {
		return d, false
	}

	fields := strings.Fields(counters)
	if len(fields) < 9 {
		return d, false
	}

	var err1, err2 error
	d.Name = strings.TrimSpace(name)
	d.RxTotal, err1 = strconv.Atoi(fields[0])
	d.TxTotal, err2 = strconv.Atoi(fields[8])

	return d, err1 == nil && err2 == nil
}

func (n *Parser) calculateSpeed(ifaces []data) {
	mul := time.Since(n.parsedAt)

	for i, d := range ifaces {
		prev, ok := n.dataBuf[d.Name]
		if !ok {
			// No previous sample, e.g. the interface just appeared.
			continue
		}

		ifaces[i].Rx = n.speed(d.RxTotal-prev.RxTotal, mul)
		ifaces[i].Tx = n.speed(d.TxTotal-prev.TxTotal, mul)
	}
}

func (n *Parser) speed(val int, t time.Duration) int {
	return int(math.Round(float64(val) * float64(time.Second) / float64(t)))
}
