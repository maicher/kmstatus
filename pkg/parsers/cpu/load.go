package cpu

import (
	"fmt"
	"os"

	"github.io/maicher/kmstatus/pkg/parsers"
)

const statSrcFile = "/proc/stat"

// Load holds the average CPU load, calculated between the previous and current stats.
// The result is combined for all cores and is expressed in percentages.
type Load float64

type stat struct {
	active int
	idle   int
}

func (s *stat) total() int {
	return s.active + s.idle
}

type LoadParser struct {
	stat     stat
	statFile *os.File
}

func (p *LoadParser) Parse() (any, error) {
	var user, nice, system, idle int
	stat := stat{}

	p.statFile.Seek(0, 0)
	_, err := fmt.Fscanf(p.statFile, "cpu %d %d %d %d", &user, &nice, &system, &idle)
	if err != nil {
		return Load(0), fmt.Errorf("Load parser %s: %w", p.statFile.Name(), err)
	}

	stat.active = user + nice + system
	stat.idle = idle

	load := calculateLoad(stat, p.stat)
	p.stat = stat

	return load, nil
}

func NewLoadParser() (parsers.Parser, error) {
	parser := LoadParser{}
	file, err := os.Open(statSrcFile)

	if err != nil {
		return &parser, fmt.Errorf("Load parser: %w", err)
	}

	parser.statFile = file

	return &parser, nil
}

// CalculateLoad calculates average CPU load between the previous and current stats.
// The result is combined for all cores and is expressed in percentages.
func calculateLoad(currentStat, prevStat stat) Load {
	activeDiff := float64(currentStat.active - prevStat.active)
	totalDiff := float64(currentStat.total() - prevStat.total())

	return Load((activeDiff / totalDiff) * 100)
}
