package temperature

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const globIntel = "/sys/devices/virtual/thermal/thermal_zone*/temp"
const globAMD = "/sys/class/hwmon/hwmon*/temp1_input"

type sensor struct {
	name string
	file *os.File
}

type Parser struct {
	sensors []sensor
}

func NewParser() (*Parser, error) {
	var p Parser

	paths, err := filepath.Glob(globIntel)
	if err == nil && len(paths) > 0 {
		return newTemperatureParserIntel(paths)
	}

	paths, err = filepath.Glob(globAMD)
	if err == nil && len(paths) > 0 {
		return newTemperatureParserAMD(paths)
	}

	return &p, fmt.Errorf("temp parser: no files matching pattern %s nor %s", globIntel, globAMD)
}

func (p *Parser) Parse(data []data) error {
	var val int

	for i, sensor := range p.sensors {
		sensor.file.Seek(0, 0)
		_, err := fmt.Fscanf(sensor.file, "%d", &val)
		if err != nil {
			return fmt.Errorf("temp parser %s: %w", sensor.file.Name(), err)
		}

		data[i].Value = val / 1000
	}

	return nil
}

func (p *Parser) Names() (names []string) {
	for _, sensor := range p.sensors {
		names = append(names, sensor.name)
	}

	return names
}

func newTemperatureParserIntel(paths []string) (*Parser, error) {
	var p Parser

	p.sensors = make([]sensor, len(paths))
	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return &p, fmt.Errorf("temp parser: %s", err)
		}

		p.sensors[i].name = fmt.Sprintf("CPU%d", i)
		p.sensors[i].file = file
	}

	return &p, nil
}

func newTemperatureParserAMD(paths []string) (*Parser, error) {
	var p Parser

	p.sensors = make([]sensor, len(paths))
	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return &p, fmt.Errorf("temp parser: %s", err)
		}

		name, err := os.ReadFile(strings.ReplaceAll(path, "temp1_input", "name"))
		if err != nil {
			return &p, err
		}

		p.sensors[i].name = strings.TrimSpace(string(name))
		p.sensors[i].file = file
	}

	return &p, nil
}
