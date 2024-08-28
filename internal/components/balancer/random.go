package balancer

import (
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/components/upstream"
)

type random struct {
	r *rand.Rand
}

func newRandom() Balancer {
	return &random{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (s random) Selector(_ *ghttp.Request, ss *upstream.Service) (*upstream.Upstream, bool) {
	return ss.Ups[s.r.Intn(len(ss.Ups))], true
}
