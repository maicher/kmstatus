package temperature

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Thermal zones and hwmon sensors are both available on Intel and AMD CPUs,
// thermal zones are preferred when present.
const globThermal = "/sys/devices/virtual/thermal/thermal_zone*/temp"
const globHwmon = "/sys/class/hwmon/hwmon*/temp1_input"

type sensor struct {
	name string
	file *os.File
}

type Parser struct {
	sensors []sensor
}

func NewParser() (*Parser, error) {
	var p Parser

	paths, err := filepath.Glob(globThermal)
	if err == nil && len(paths) > 0 {
		return newParser(paths, "type")
	}

	paths, err = filepath.Glob(globHwmon)
	if err == nil && len(paths) > 0 {
		return newParser(paths, "name")
	}

	return &p, fmt.Errorf("temp parser: no files matching pattern %s nor %s", globThermal, globHwmon)
}

// Parse reads all sensors. A sensor which can not be read (e.g. a wifi card
// that is turned off) keeps its previous value and does not prevent
// the remaining sensors from being updated.
func (p *Parser) Parse(data []data) error {
	var val int
	var errs []error

	for i, sensor := range p.sensors {
		sensor.file.Seek(0, 0)
		_, err := fmt.Fscanf(sensor.file, "%d", &val)
		if err != nil {
			errs = append(errs, fmt.Errorf("temp parser %s: %w", sensor.file.Name(), err))
			continue
		}

		data[i].Value = val / 1000
	}

	return errors.Join(errs...)
}

func (p *Parser) Names() (names []string) {
	for _, sensor := range p.sensors {
		names = append(names, sensor.name)
	}

	return names
}

// newParser opens the temperature files. The sensor's name is read
// from the nameFile located in the same directory.
func newParser(paths []string, nameFile string) (*Parser, error) {
	var p Parser

	sortNumerically(paths)

	p.sensors = make([]sensor, len(paths))
	for i, path := range paths {
		file, err := os.Open(path)
		if err != nil {
			return &p, fmt.Errorf("temp parser: %s", err)
		}

		name, err := os.ReadFile(filepath.Join(filepath.Dir(path), nameFile))
		if err != nil {
			return &p, fmt.Errorf("temp parser: %s", err)
		}

		p.sensors[i].name = strings.TrimSpace(string(name))
		p.sensors[i].file = file
	}

	return &p, nil
}

var number = regexp.MustCompile(`\d+`)

// sortNumerically sorts paths by the number in their directory name,
// so that thermal_zone10 goes after thermal_zone2.
func sortNumerically(paths []string) {
	key := func(path string) int {
		n, _ := strconv.Atoi(number.FindString(filepath.Base(filepath.Dir(path))))
		return n
	}

	sort.SliceStable(paths, func(i, j int) bool {
		return key(paths[i]) < key(paths[j])
	})
}
