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
		Interval: calculateMinIntercal(pc),
	}
}

func (g *Generate) Loop() {
	onTick(g.Interval, func() {
		g.Ch <- view.RenderView{}
	})
}

func calculateMinIntercal(pc []config.ParserConfig) time.Duration {
	var intervals []int

	for _, v := range pc {
		if v.Interval > 0 {
			intervals = append(intervals, int(v.Interval))
		}
	}

	sort.Ints(intervals)

	return time.Duration(intervals[0])
}
