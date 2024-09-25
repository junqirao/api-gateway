package balancer

import (
	"math"
	"sync/atomic"
	"time"
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
		interval time.Duration
		v        *atomic.Int64
		old      *atomic.Int64
		last     time.Time
		lastOp   time.Time
	}
)

// NewMeasurable of balance target
func NewMeasurable(resetInterval time.Duration, initVal ...int64) Measurable {
	v := &atomic.Int64{}
	if len(initVal) > 0 {
		v.Store(initVal[0])
	}
	now := time.Now()
	return &measurable{v: v, old: &atomic.Int64{}, interval: resetInterval, last: now, lastOp: now}
}

// Load implements Measurable
func (m *measurable) Load() int64 {
	var res int64 = 0
	m.checkAndDo(time.Now(), func() {
		res = m.v.Load()
	})
	return res
}

// AddLoad implements Measurable
func (m *measurable) AddLoad(delta int64) {
	m.checkAndDo(time.Now(), func() {
		m.v.Add(delta)
	})
}

// SetLoad implements Measurable
func (m *measurable) SetLoad(val int64) {
	m.checkAndDo(time.Now(), func() {
		m.v.Store(val)
	})
}

func (m *measurable) checkAndDo(now time.Time, do func()) {
	if m.interval > 0 {
		// decrease by last operation,
		// m.v - m.old * (now-m.lastOp/m.interval)
		m.decr()
		// reset when time pass by m.interval
		if now.Sub(m.last) > m.interval {
			// reset
			m.last = time.Now()
			if m.last.Sub(m.lastOp) > m.interval {
				m.lastOp = m.last
			}
			m.old.Store(m.v.Load())
		}
	}

	do()
}

func (m *measurable) decr() {
	old := float64(m.old.Load())
	cur := m.v.Load()
	if old <= 0 {
		return
	}
	delta := old * float64(time.Now().Sub(m.lastOp)) / float64(m.interval)
	if delta > 0 {
		v := int64(math.Round(delta))
		if v > cur {
			v = cur
		}
		m.v.Add(v * -1)
	}
	m.lastOp = time.Now()
}
