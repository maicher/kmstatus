package cpu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.io/maicher/kmstatus/pkg/parsers"
)

const filesWithFreqInfoGlob = "/sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq"

type FreqParser struct {
	files []*os.File
}

func (p *FreqParser) Parse() (any, error) {
	var val, sum int

	for _, f := range p.files {
		f.Seek(0, 0)
		_, err := fmt.Fscanf(f, "%d", &val)
		if err != nil {
			return Freq(0), fmt.Errorf("Freq parser %s: %w", f.Name(), err)
		}
		sum = sum + val
	}

	f := Freq(sum / len(p.files))

	return f, nil
}

func NewFreqParser() (parsers.Parser, error) {
	parser := FreqParser{}
	paths, err := filepath.Glob(filesWithFreqInfoGlob)
	if err != nil {
		return &parser, fmt.Errorf("Freq parser: %w", err)
	}

	parser.files = make([]*os.File, len(paths))
	if len(paths) == 0 {
		return &parser, fmt.Errorf("Freq parser: no files matching the pattern %s", filesWithFreqInfoGlob)
	}

	for i, p := range paths {
		file, err := os.Open(p)
		if err != nil {
			return &parser, fmt.Errorf("Freq parser: %w", err)
		}

		parser.files[i] = file
	}

	return &parser, nil
}
