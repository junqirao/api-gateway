package loadbalance

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type weightedRoundRobin struct {
}

func newRoundRobin() Balancer {
	return &weightedRoundRobin{}
}

func (w *weightedRoundRobin) Selector(_ *ghttp.Request, ups []Weighted) (ref int, ok bool) {
	if len(ups) == 0 {
		return
	}
	return nextWeighted(ups).Ref(), true
}

func nextWeighted(servers []Weighted) (selected Weighted) {
	var total int64 = 0
	var best *weighted
	for i := 0; i < len(servers); i++ {
		w := servers[i].(*weighted)
		if w == nil {
			continue
		}
		w.currentWeight.Add(w.effectiveWeight.Load())
		total += w.effectiveWeight.Load()
		if w.effectiveWeight.Load() < w.weight.Load() {
			w.effectiveWeight.Add(1)
		}
		if best == nil || w.currentWeight.Load() > best.currentWeight.Load() {
			best = w
		}
	}
	if best == nil {
		return nil
	}
	best.currentWeight.Add(-total)
	return best
}
