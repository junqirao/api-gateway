package balancer

import (
	"sync/atomic"
)

type roundRobinBalancer struct {
	idx *atomic.Uint64
}

func newRoundRobin() Balancer {
	return &roundRobinBalancer{
		idx: &atomic.Uint64{},
	}
}

func (r *roundRobinBalancer) Pick(objects []any, _ ...any) (o any, err error) {
	idx := r.idx.Add(1) % uint64(len(objects))
	r.idx.Store(idx)
	return objects[idx], nil
}
