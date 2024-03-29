package mem

type Data struct {
	Used      int
	Total     int
	SwapUsed  int
	SwapTotal int
}

func (d Data) UsedPercentage() float64 {
	return 100 * float64(d.Used) / float64(d.Total)
}

func (d Data) SwapUsedPercentage() float64 {
	return 100 * float64(d.SwapUsed) / float64(d.SwapTotal)
}
