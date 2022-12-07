package cpu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.io/maicher/stbar/pkg/parsers"
)

const srcFiles = "/sys/devices/virtual/thermal/thermal_zone*/temp"

// Temp holds temperature values of the CPUs.
// Each value is for one CPU socket and is expressed in Celsius degrees.
type Temp []int

type TempParser struct {
	files     []*os.File
	formatter string
}

func (p *TempParser) Parse() (any, error) {
	var val int
	t := Temp(make([]int, len(p.files)))

	for i, f := range p.files {
		f.Seek(0, 0)
		_, err := fmt.Fscanf(f, "%d", &val)
		if err != nil {
			return t, fmt.Errorf("Temp Parser %s: %w", f.Name(), err)
		}
		t[i] = val / 1000
	}

	return t, nil
}

func NewTempParser() (parsers.Parser, error) {
	paths, err := filepath.Glob(srcFiles)
	if err != nil {
		return nil, fmt.Errorf("Temp Parser: %w", err)
	}

	parser := TempParser{
		files: make([]*os.File, len(paths)),
	}

	if len(paths) == 0 {
		return &parser, fmt.Errorf("Temp Parser no files matching the pattern %s", srcFiles)
	}

	for i, p := range paths {
		file, err := os.Open(p)
		if err != nil {
			return &parser, fmt.Errorf("Temp Parser: %w", err)
		}

		parser.files[i] = file
	}

	return &parser, nil
}
