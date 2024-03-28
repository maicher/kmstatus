package temperature

import "math"

type Data struct {
	Name  string
	Value int
}

func (d Data) Celsius() int {
	return d.Value
}

func (d Data) Fahrenheit() int {
	return int(math.Round((float64(d.Value) * 1.8) + 32))
}
