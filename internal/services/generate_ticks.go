package services

import (
	"sort"
	"time"

	"github.io/maicher/kmstatus/internal/config"
	"github.io/maicher/kmstatus/internal/view"
)

type GenerateTicks struct {
	Ch       chan<- any
	Interval time.Duration
}

func NewGenerateTicks(ch chan<- any, pc []config.ParserSettings) *GenerateTicks {
	return &GenerateTicks{
		Ch:       ch,
		Interval: calculateMinInterval(pc),
	}
}

func (g *GenerateTicks) GenerateTicks() {
	go func() {
		onTick(g.Interval, func() {
			g.Ch <- view.RenderView{}
		})
	}()
}

func calculateMinInterval(pc []config.ParserSettings) time.Duration {
	var intervals []int

	for _, v := range pc {
		if v.Interval > 0 {
			intervals = append(intervals, int(v.Interval))
		}
	}

	sort.Ints(intervals)

	return time.Duration(intervals[0])
}
