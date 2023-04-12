package network

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.io/maicher/kmstatus/pkg/parsers"
)

const fileWithNetworkInfo = "/proc/net/dev"

type NetParser struct {
	interfacesBuf map[string]Interface
	parsedAt      time.Time

	file *os.File
}

func NewNetParser() (parsers.Parser, error) {
	var err error

	parser := NetParser{}
	parser.interfacesBuf = make(map[string]Interface)
	parser.file, err = os.Open(fileWithNetworkInfo)
	if err != nil {
		return &parser, fmt.Errorf("Net parser: %w", err)
	}

	return &parser, nil
}

func (p *NetParser) Parse() (any, error) {
	var n Net

	_, err := p.file.Seek(0, 0)
	if err != nil {
		return n, fmt.Errorf("Net parser: %w", err)
	}

	s := bufio.NewScanner(p.file)
	s.Split(bufio.ScanLines)

	// Drop first 2 lines since they contain headers
	s.Scan()
	s.Scan()

	for s.Scan() {
		n.Interfaces = append(n.Interfaces, p.parseInterface(s.Text()))
	}

	n.calculateSpeed(p.interfacesBuf, p.parsedAt)

	// Remember for calculating speed in the next cycle.
	for _, i := range n.Interfaces {
		p.interfacesBuf[i.Name] = i
	}
	p.parsedAt = time.Now()

	return n, nil
}

func (p *NetParser) parseInterface(s string) Interface {
	var (
		name        string
		rx, tx, ign int
	)

	r := strings.NewReader(s)
	fmt.Fscanln(r, &name, &rx, &ign, &ign, &ign, &ign, &ign, &ign, &ign, &tx)

	return Interface{
		Name:    strings.TrimSuffix(name, ":"),
		RxTotal: rx,
		TxTotal: tx,
	}
}
