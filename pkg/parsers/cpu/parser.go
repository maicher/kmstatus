package cpu

import (
	"fmt"
	"os"
)

var srcFile = "/proc/stat"

type Parser struct {
	stat Stat
	file *os.File
}

func (p *Parser) Parse() (CPU, error) {
	var user, nice, system, idle int
	c := CPU{}

	p.file.Seek(0, 0)
	_, err := fmt.Fscanf(p.file, "cpu %d %d %d %d", &user, &nice, &system, &idle)
	if err != nil {
		return c, fmt.Errorf("cpu Parser %s: %w", p.file.Name(), err)
	}

	c.Stat = Stat{
		Active: user + nice + system,
		Idle:   idle,
	}
	c.PrevStat = p.stat

	p.stat = c.Stat

	return c, nil
}

func NewParser() (*Parser, error) {
	parser := Parser{}
	file, err := os.Open(srcFile)

	if err != nil {
		return &parser, fmt.Errorf("CPU Parser: %w", err)
	}

	parser.file = file

	return &parser, nil
}
