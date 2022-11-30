package cpu

type CPU struct {
	stat     Stat
	prevStat Stat
}

func (c *CPU) Load() float64 {
	activeDiff := float64(c.stat.active - c.prevStat.active)
	totalDiff := float64(c.stat.Total() - c.prevStat.Total())

	return (activeDiff / totalDiff) * 100
}
