package mem

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type MemParser struct {
	file *os.File
}

var meminfoFilePath = "/proc/meminfo"

func (p *MemParser) Parse(d *data) error {
	var free, buffers, cached, swapFree int

	p.scanByLines(func(line string) {
		if strings.HasPrefix(line, "MemTotal") {
			readValue(line, &d.Total)
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
			readValue(line, &d.SwapTotal)
		}
		if strings.HasPrefix(line, "SwapFree") {
			readValue(line, &swapFree)
		}
	})

	d.Used = d.Total - free - buffers - cached
	d.SwapUsed = d.SwapTotal - swapFree

	return nil
}

func NewMemParser() (*MemParser, error) {
	var p MemParser
	var err error

	p.file, err = os.Open(meminfoFilePath)
	if err != nil {
		return &p, fmt.Errorf("Mem parser error: %s", err)
	}

	return &p, nil
}

func (p MemParser) scanByLines(f func(line string)) {
	p.file.Seek(0, 0)
	s := bufio.NewScanner(p.file)
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
