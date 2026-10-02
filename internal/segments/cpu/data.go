package cpu

type data struct {
	Load float64

	// Average frequency in kHz
	Freq int
}

func (d data) FreqMHz() float64 {
	return float64(d.Freq) / 1000
}

func (d data) FreqGHz() float64 {
	return float64(d.Freq) / (1000 * 1000)
}
