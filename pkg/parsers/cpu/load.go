package cpu

import (
	"fmt"
	"os"
)

// Load holds the average CPU load, calculated between the previous and current stats.
// The result is combined for all cores and is expressed in percentages.
type Load float64

var statSrcFile = "/proc/stat"

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

func (p *LoadParser) Parse() (Load, error) {
	var user, nice, system, idle int
	s := stat{}

	p.statFile.Seek(0, 0)
	_, err := fmt.Fscanf(p.statFile, "cpu %d %d %d %d", &user, &nice, &system, &idle)
	if err != nil {
		return Load(0), fmt.Errorf("cpu Parser %s: %w", p.statFile.Name(), err)
	}

	s.active = user + nice + system
	s.idle = idle

	load := calculateLoad(s, p.stat)
	p.stat = s

	return load, nil
}

func NewLoadParser() (*LoadParser, error) {
	parser := LoadParser{}
	file, err := os.Open(statSrcFile)

	if err != nil {
		return &parser, fmt.Errorf("CPU Parser: %w", err)
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
