package cpu

import (
	"fmt"
	"os"
)

const statFilePath = "/proc/stat"

type stat struct {
	active int
	idle   int
}

func (s *stat) total() int {
	return s.active + s.idle
}

type LoadParser struct {
	statFile *os.File
	stat     stat
}

func NewLoadParser() (*LoadParser, error) {
	var p LoadParser
	var err error

	p.statFile, err = os.Open(statFilePath)

	if err != nil {
		return &p, fmt.Errorf("Error initializing CPU load parser: %s", err)
	}

	return &p, nil
}

func (p *LoadParser) Parse(load *float64) error {
	var user, nice, system, idle int
	var stat stat

	p.statFile.Seek(0, 0)
	_, err := fmt.Fscanf(p.statFile, "cpu %d %d %d %d", &user, &nice, &system, &idle)
	if err != nil {
		return fmt.Errorf("CPU load parser error: %s", err)
	}

	stat.active = user + nice + system
	stat.idle = idle

	*load = p.calculateLoad(stat, p.stat)
	p.stat = stat

	return nil
}

// load calculates average CPU load between the previous and current stats.
// The result is combined for all cores and is expressed in percentages.
func (p *LoadParser) calculateLoad(stat, prevStat stat) float64 {
	activeDiff := float64(stat.active - prevStat.active)
	totalDiff := float64(stat.total() - prevStat.total())

	return (activeDiff / totalDiff) * 100
}
