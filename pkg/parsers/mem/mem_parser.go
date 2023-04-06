package mem

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.io/maicher/kmstatus/pkg/parsers"
)

var meminfoFilePath = "/proc/meminfo"

type MemParser struct {
	File *os.File
}

func (p *MemParser) Parse() (any, error) {
	var total, free, buffers, cached, swapTotal, swapFree int

	p.scanByLines(func(line string) {
		if strings.HasPrefix(line, "MemTotal") {
			readValue(line, &total)
		}
		if strings.HasPrefix(line, "MemFree") {
			readValue(line, &free)
		}
		if strings.HasPrefix(line, "Buffer") {
			readValue(line, &buffers)
		}
		if strings.HasPrefix(line, "Cached") {
			readValue(line, &cached)
		}
		if strings.HasPrefix(line, "SwapTotal") {
			readValue(line, &swapTotal)
		}
		if strings.HasPrefix(line, "SwapFree") {
			readValue(line, &swapFree)
		}
	})

	m := Mem{
		MemTotal:  total,
		MemUsed:   total - free - buffers - cached,
		SwapTotal: swapTotal,
		SwapUsed:  swapTotal - swapFree,
	}

	return m, nil
}

func NewMemParser() (parsers.Parser, error) {
	parser := MemParser{}
	file, err := os.Open(meminfoFilePath)

	if err != nil {
		return &parser, fmt.Errorf("Mem parser: %w", err)
	}

	parser.File = file

	return &parser, nil
}

func (p MemParser) scanByLines(f func(line string)) {
	p.File.Seek(0, 0)
	s := bufio.NewScanner(p.File)
	s.Split(bufio.ScanLines)

	for s.Scan() {
		f(s.Text())
	}
}

func readValue(line string, v *int) {
	var ignored string

	r := strings.NewReader(line)
	fmt.Fscanf(r, "%s %d", &ignored, v)
}
