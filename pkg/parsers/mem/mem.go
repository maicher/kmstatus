package mem

type SpaceKB int

type Mem struct {
	MemTotal  SpaceKB
	MemUsed   SpaceKB
	SwapTotal SpaceKB
	SwapUsed  SpaceKB
}

func (m *Mem) MemUsedPercentage() float64 {
	return 100 * float64(m.MemUsed) / float64(m.MemTotal)
}

func (m *Mem) SwapUsedPercentage() float64 {
	return 100 * float64(m.SwapUsed) / float64(m.SwapTotal)
}
