package loadbalance

import (
	"sync/atomic"
)

type (
	Weighted interface {
		Weight() int64
		AddWeight(weight int64)
		SetWeight(weight int64)
		Ref() int
		SetRef(ref int)
	}

	// concurrency safe dynamic weighted
	weighted struct {
		ref             int
		weight          *atomic.Int64
		effectiveWeight *atomic.Int64 // round-robin only
		currentWeight   *atomic.Int64 // round-robin only
	}
)

func NewWeighted(weight int64, ref ...int) Weighted {
	wt := &weighted{
		weight:          &atomic.Int64{},
		effectiveWeight: &atomic.Int64{},
		currentWeight:   &atomic.Int64{},
	}
	if len(ref) > 0 {
		wt.ref = ref[0]
	}
	wt.weight.Store(weight)
	return wt
}

func (w *weighted) Weight() int64 {
	return w.weight.Load()
}

func (w *weighted) AddWeight(weight int64) {
	w.weight.Add(weight)
}

func (w *weighted) Ref() int {
	return w.ref
}

func (w *weighted) SetWeight(weight int64) {
	w.weight.Store(weight)
}

func (w *weighted) SetRef(ref int) {
	w.ref = ref
}
