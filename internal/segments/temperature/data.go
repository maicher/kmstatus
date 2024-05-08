package temperature

import "math"

type data struct {
	Name  string
	Value int
}

func (d data) Celsius() int {
	return d.Value
}

func (d data) Fahrenheit() int {
	return int(math.Round((float64(d.Value) * 1.8) + 32))
}
