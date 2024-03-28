package cpu

import (
	"fmt"
	"os"
	"path/filepath"
)

const filesWithFreqInfoGlob = "/sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq"

type FreqParser struct {
	files []*os.File
}

func NewFreqParser() (*FreqParser, error) {
	var p FreqParser
	paths, err := filepath.Glob(filesWithFreqInfoGlob)
	if err != nil {
		return &p, fmt.Errorf("CPU freq parser error: %s", err)
	}

	p.files = make([]*os.File, len(paths))
	if len(paths) == 0 {
		return &p, fmt.Errorf("CPU freq parser: no files matching the pattern %s", filesWithFreqInfoGlob)
	}

	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return &p, fmt.Errorf("CPU freq parser: %w", err)
		}

		p.files[i] = file
	}

	return &p, nil
}

func (p *FreqParser) Parse(freq *int) error {
	var val, sum int

	for _, f := range p.files {
		f.Seek(0, 0)
		_, err := fmt.Fscanf(f, "%d", &val)
		if err != nil {
			return fmt.Errorf("Freq parser %s: %w", f.Name(), err)
		}
		sum = sum + val
	}

	*freq = sum / len(p.files)

	return nil
}
