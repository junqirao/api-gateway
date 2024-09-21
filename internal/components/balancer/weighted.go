package balancer

import (
	"sync/atomic"
)

type (
	// Weighable interface, used in weighted balancer
	// use NewWeighable to create, do not implement this
	Weighable interface {
		// Weight returns the weight of the object,
		// only used in weighted balancer
		Weight() int64
		// AddWeight adds weight to the object
		AddWeight(delta int64)
		// Current returns the current weight of the object,
		// only used in wrr balancer
		Current() *atomic.Int64
		// Effective returns the effective weight of the object,
		// only used in wrr balancer
		Effective() *atomic.Int64
	}

	weighted struct {
		weight          *atomic.Int64
		effectiveWeight *atomic.Int64 // round-robin only
		currentWeight   *atomic.Int64 // round-robin only
	}
)

// NewWeighable creates a new Weighable
func NewWeighable(weight int64) Weighable {
	wt := &weighted{
		weight:          &atomic.Int64{},
		effectiveWeight: &atomic.Int64{},
		currentWeight:   &atomic.Int64{},
	}
	wt.weight.Store(weight)
	return wt
}

// Weight implements Weighable
func (w *weighted) Weight() int64 {
	return w.weight.Load()
}

// AddWeight implements Weighable
func (w *weighted) AddWeight(weight int64) {
	w.weight.Add(weight)
}

// Current implements Weighable
func (w *weighted) Current() *atomic.Int64 {
	return w.currentWeight
}

// Effective implements Weighable
func (w *weighted) Effective() *atomic.Int64 {
	return w.effectiveWeight
}
