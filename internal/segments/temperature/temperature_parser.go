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

type TemperatureParser struct {
	sensors   []sensor
	formatter string
}

func NewTemperatureParser() (*TemperatureParser, error) {
	var p TemperatureParser

	paths, err := filepath.Glob(globIntel)
	if err == nil && len(paths) > 0 {
		return newTemperatureParserIntel(paths)
	}

	paths, err = filepath.Glob(globAMD)
	if err == nil && len(paths) > 0 {
		return newTemperatureParserAMD(paths)
	}

	return &p, fmt.Errorf("Temp parser: no files matching pattern %s nor %s", globIntel, globAMD)
}

func (p *TemperatureParser) Parse(data []Data) error {
	var val int

	for i, sensor := range p.sensors {
		sensor.file.Seek(0, 0)
		_, err := fmt.Fscanf(sensor.file, "%d", &val)
		if err != nil {
			return fmt.Errorf("Temp parser %s: %w", sensor.file.Name(), err)
		}

		data[i].Value = val / 1000
	}

	return nil
}

func (p *TemperatureParser) Names() (names []string) {
	for _, sensor := range p.sensors {
		names = append(names, sensor.name)
	}

	return names
}

func newTemperatureParserIntel(paths []string) (*TemperatureParser, error) {
	var p TemperatureParser

	p.sensors = make([]sensor, len(paths))
	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return &p, fmt.Errorf("Temp parser: %s", err)
		}

		p.sensors[i].name = fmt.Sprintf("CPU%d", i)
		p.sensors[i].file = file
	}

	return &p, nil
}

func newTemperatureParserAMD(paths []string) (*TemperatureParser, error) {
	var p TemperatureParser

	p.sensors = make([]sensor, len(paths))
	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return &p, fmt.Errorf("Temp parser: %s", err)
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
