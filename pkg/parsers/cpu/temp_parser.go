package cpu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.io/maicher/kmstatus/pkg/parsers"
)

const filesWithTemeratureInfoGlob = "/sys/devices/virtual/thermal/thermal_zone*/temp"

var nullTemp = Temp(make([]int, 1))

type NullTempParser struct {
}

func (p *NullTempParser) Parse() (any, error) {
	return nullTemp, nil
}

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
			return t, fmt.Errorf("Temp parser %s: %w", f.Name(), err)
		}
		t[i] = val / 1000
	}

	return t, nil
}

func NewTempParser() (parsers.Parser, error) {
	paths, err := filepath.Glob(filesWithTemeratureInfoGlob)
	if err != nil {
		return nil, fmt.Errorf("Temp parser: %w", err)
	}

	parser := TempParser{
		files: make([]*os.File, len(paths)),
	}

	if len(paths) == 0 {
		fmt.Printf("Temp parser: no files matching the pattern %s", filesWithTemeratureInfoGlob)
		return &NullTempParser{}, nil
	}

	for i, p := range paths {
		file, err := os.Open(p)
		if err != nil {
			return &parser, fmt.Errorf("Temp parser: %w", err)
		}

		parser.files[i] = file
	}

	return &parser, nil
}
