package network

import (
	"bufio"
	"fmt"
	"math"
	"os"
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

func (n *Parser) Parse(data []data) error {
	_, err := n.file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("Network parser: %s", err)
	}

	s := bufio.NewScanner(n.file)
	s.Split(bufio.ScanLines)

	// Drop first 2 lines since they contain headers
	s.Scan()
	s.Scan()

	var ign int
	var i int
	var r *strings.Reader

	for s.Scan() {
		if i > len(data)-1 {
			continue
		}

		r = strings.NewReader(s.Text())
		fmt.Fscanln(r, &data[i].Name, &data[i].RxTotal, &ign, &ign, &ign, &ign, &ign, &ign, &ign, &data[i].TxTotal)

		data[i].Name = strings.TrimSuffix(data[i].Name, ":")
		i++
	}

	n.calculateSpeed(data)
	n.parsedAt = time.Now()

	// Buffer to calculate speed in the next cycle.
	for _, d := range data {
		n.dataBuf[d.Name] = d
	}

	return nil
}

func (n *Parser) calculateSpeed(data []data) {
	mul := time.Since(n.parsedAt)

	for index, i := range data {
		data[index].Rx = n.speed(i.RxTotal-n.dataBuf[i.Name].RxTotal, mul)
		data[index].Tx = n.speed(i.TxTotal-n.dataBuf[i.Name].TxTotal, mul)
	}
}

func (n *Parser) speed(val int, t time.Duration) int {
	return int(math.Round(float64(val) * float64(time.Second) / float64(t)))
}
