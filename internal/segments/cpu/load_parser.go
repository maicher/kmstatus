package cpu

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
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
		return &p, fmt.Errorf("error initializing CPU load parser: %s", err)
	}

	// Take the first sample, so that the first load is not the average since boot.
	p.stat, err = p.readStat()
	if err != nil {
		return &p, err
	}

	return &p, nil
}

func (p *LoadParser) Parse(load *float64) error {
	stat, err := p.readStat()
	if err != nil {
		return err
	}

	// No time passed between samples, keep the previous load.
	if stat.total() == p.stat.total() {
		return nil
	}

	*load = p.calculateLoad(stat, p.stat)
	p.stat = stat

	return nil
}

// readStat reads the aggregated "cpu" line of /proc/stat:
// cpu user nice system idle iowait irq softirq steal guest guest_nice
// Guest time is already included in user and nice.
func (p *LoadParser) readStat() (stat, error) {
	var s stat

	p.statFile.Seek(0, 0)
	line, err := bufio.NewReader(p.statFile).ReadString('\n')
	if err != nil && line == "" {
		return s, fmt.Errorf("CPU load parser error: %s", err)
	}

	fields := strings.Fields(line)
	if len(fields) < 5 || fields[0] != "cpu" {
		return s, fmt.Errorf("CPU load parser error: unexpected format: %q", line)
	}

	var v [8]int
	for i := 1; i < len(fields) && i <= len(v); i++ {
		v[i-1], err = strconv.Atoi(fields[i])
		if err != nil {
			return s, fmt.Errorf("CPU load parser error: %s", err)
		}
	}

	user, nice, system, idle, iowait, irq, softirq, steal := v[0], v[1], v[2], v[3], v[4], v[5], v[6], v[7]
	s.active = user + nice + system + irq + softirq + steal
	s.idle = idle + iowait

	return s, nil
}

// load calculates average CPU load between the previous and current stats.
// The result is combined for all cores and is expressed in percentages.
func (p *LoadParser) calculateLoad(stat, prevStat stat) float64 {
	activeDiff := float64(stat.active - prevStat.active)
	totalDiff := float64(stat.total() - prevStat.total())

	return (activeDiff / totalDiff) * 100
}
