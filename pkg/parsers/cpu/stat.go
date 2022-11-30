package cpu

type Stat struct {
	active int
	idle   int
}

func (s *Stat) Total() int {
	return s.active + s.idle
}
