package mem

type data struct {
	Used      int
	Total     int
	SwapUsed  int
	SwapTotal int
}

func (d data) UsedPercentage() float64 {
	return 100 * float64(d.Used) / float64(d.Total)
}

func (d data) SwapUsedPercentage() float64 {
	return 100 * float64(d.SwapUsed) / float64(d.SwapTotal)
}
