package balancer

import (
	"sync/atomic"

	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/components/upstream"
)

type roundRobin struct {
	idx *atomic.Int32
}

func newRoundRobin() *roundRobin {
	i := &atomic.Int32{}
	i.Store(0)
	return &roundRobin{idx: i}
}

func (r *roundRobin) Selector(_ *ghttp.Request, ss *upstream.Service) (*upstream.Upstream, bool) {
	var (
		i      = int(r.idx.Load())
		length = len(ss.Ups)
	)

	if i >= length {
		i = length - 1
	} else {
		if i == length-1 {
			i = 0
		} else {
			i++
		}
	}
	r.idx.Store(int32(i))
	return ss.Ups[i], true
}
