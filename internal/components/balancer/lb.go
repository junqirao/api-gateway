package balancer

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/components/upstream"
)

const (
	StrategyRandom     = "random"
	StrategyRoundRobin = "round_robin"
)

type (
	Balancer interface {
		Selector(r *ghttp.Request, ss *upstream.Service) (*upstream.Upstream, bool)
	}
)

func New(strategy ...string) (b Balancer) {
	st := StrategyRandom
	if len(strategy) > 0 && strategy[0] != "" {
		st = strategy[0]
	}
	switch st {
	case StrategyRandom:
		b = newRandom()
	case StrategyRoundRobin:
		b = newRoundRobin()
	}
	return
}
