package services

import (
	"sort"
	"time"

	"github.io/maicher/kmstatus/app/config"
	"github.io/maicher/kmstatus/app/view"
)

type Generate struct {
	Ch       chan<- any
	Interval time.Duration
}

func NewGenerate(ch chan<- any, pc []config.ParserConfig) *Generate {
	return &Generate{
		Ch:       ch,
		Interval: calculateMinInterval(pc),
	}
}

func (g *Generate) Loop() {
	go func() {
		onTick(g.Interval, func() {
			g.Ch <- view.RenderView{}
		})
	}()
}

func calculateMinInterval(pc []config.ParserConfig) time.Duration {
	var intervals []int

	for _, v := range pc {
		if v.Interval > 0 {
			intervals = append(intervals, int(v.Interval))
		}
	}

	sort.Ints(intervals)

	return time.Duration(intervals[0])
}
