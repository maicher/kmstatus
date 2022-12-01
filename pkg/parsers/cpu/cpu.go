package cpu

type CPU struct {
	Stat     Stat
	PrevStat Stat
}

// Load calculates average CPU load between the previous and current stats.
// The result is combined for all cores and is expressed in percentages.
func (c *CPU) Load() float64 {
	activeDiff := float64(c.Stat.Active - c.PrevStat.Active)
	totalDiff := float64(c.Stat.Total() - c.PrevStat.Total())

	return (activeDiff / totalDiff) * 100
}
