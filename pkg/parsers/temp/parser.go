package temp

import (
	"fmt"
	"os"
	"path/filepath"
)

var srcFiles = "/sys/devices/virtual/thermal/thermal_zone*/temp"

type Parser struct {
	files []*os.File
}

func (p *Parser) Parse() (Temp, error) {
	var val int
	t := Temp{Values: make([]int, len(p.files))}

	for i, f := range p.files {
		f.Seek(0, 0)
		_, err := fmt.Fscanf(f, "%d", &val)
		if err != nil {
			return Temp{}, fmt.Errorf("Temp Parser %s: %w", f.Name(), err)
		}
		t.Values[i] = val / 1000
	}

	return t, nil
}

func NewParser() (*Parser, error) {
	paths, err := filepath.Glob(srcFiles)
	if err != nil {
		return nil, fmt.Errorf("Temp Parser: %w", err)
	}

	parser := Parser{
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
