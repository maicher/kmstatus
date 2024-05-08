package cpu

type data struct {
	Load float64
	Freq int
}

func (d data) FreqMHz() float64 {
	return float64(d.Freq) / 1024
}

func (d data) FreqGHz() float64 {
	return float64(d.Freq) / (1024 * 1024)
}
