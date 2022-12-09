package cpu

import (
	"fmt"
	"os"
	"path/filepath"

	"github.io/maicher/kmstatus/pkg/parsers"
)

const freqSrcfiles = "/sys/devices/system/cpu/cpu*/cpufreq/scaling_cur_freq"

// Freq holds CPU frequency in kHz.
type Freq int

type FreqParser struct {
	files []*os.File
}

func (p *FreqParser) Parse() (any, error) {
	var val, sum int

	for _, f := range p.files {
		f.Seek(0, 0)
		_, err := fmt.Fscanf(f, "%d", &val)
		if err != nil {
			return Freq(0), fmt.Errorf("freq Parser %s: %w", f.Name(), err)
		}
		sum = sum + val
	}

	f := Freq(sum / len(p.files))

	return f, nil
}

func NewFreqParser() (parsers.Parser, error) {
	parser := FreqParser{}
	paths, err := filepath.Glob(freqSrcfiles)
	if err != nil {
		return &parser, fmt.Errorf("cpu Parser: %w", err)
	}

	parser.files = make([]*os.File, len(paths))

	if len(paths) == 0 {
		return &parser, fmt.Errorf("cpu Parser no files matching the pattern %s", freqSrcfiles)
	}

	for i, p := range paths {
		file, err := os.Open(p)
		if err != nil {
			return &parser, fmt.Errorf("cpu Parser: %w", err)
		}

		parser.files[i] = file
	}

	return &parser, nil
}
