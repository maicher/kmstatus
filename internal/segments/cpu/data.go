package cpu

type Data struct {
	Load float64
	Freq int
}

func (d Data) FreqMHz() float64 {
	return float64(d.Freq) / 1024
}

func (d Data) FreqGHz() float64 {
	return float64(d.Freq) / (1024 * 1024)
}
