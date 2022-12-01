package cpu

type Stat struct {
	Active int
	Idle   int
}

func (s *Stat) Total() int {
	return s.Active + s.Idle
}
