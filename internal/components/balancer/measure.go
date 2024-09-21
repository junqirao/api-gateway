package balancer

import (
	"sync/atomic"
)

type (
	// Measurable interface, used in StrategyLessLoad
	// use NewMeasurable to create
	Measurable interface {
		// Load returns the value representing the load state of the object
		Load() int64
		// AddLoad adds load to the object
		AddLoad(delta int64)
		// SetLoad sets the load of the object
		SetLoad(val int64)
	}

	measurable struct {
		v *atomic.Int64
	}
)

// NewMeasurable of balance target
func NewMeasurable(initVal ...int64) Measurable {
	v := &atomic.Int64{}
	if len(initVal) > 0 {
		v.Store(initVal[0])
	}
	return &measurable{v: v}
}

// Load implements Measurable
func (m *measurable) Load() int64 {
	return m.v.Load()
}

// AddLoad implements Measurable
func (m *measurable) AddLoad(delta int64) {
	m.v.Add(delta)
}

// SetLoad implements Measurable
func (m *measurable) SetLoad(val int64) {
	m.v.Store(val)
}
