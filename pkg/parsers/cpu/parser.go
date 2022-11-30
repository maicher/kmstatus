package cpu

import (
	"fmt"
	"os"
)

var defaultFilePath = "/proc/stat"

type Parser struct {
	filePath string
	stat     Stat
}

func (p *Parser) Parse() (CPU, error) {
	var user, nice, system, idle int

	file, err := os.Open(p.filePath)
	if err != nil {
		return CPU{}, fmt.Errorf("cpu Parser: %w", err)
	}
	defer file.Close()

	_, err = fmt.Fscanf(file, "cpu %d %d %d %d", &user, &nice, &system, &idle)
	if err != nil {
		return CPU{}, fmt.Errorf("cpu Parser %s: %w", p.filePath, err)
	}

	stat := Stat{
		active: user + nice + system,
		idle:   idle,
	}

	c := CPU{stat, p.stat}
	p.stat = stat

	return c, nil
}

func NewParser() *Parser {
	return &Parser{
		filePath: defaultFilePath,
	}
}
