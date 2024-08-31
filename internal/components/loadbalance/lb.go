package loadbalance

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	StrategyRandom     = "random"
	StrategyRoundRobin = "round_robin"
)

type (
	Balancer interface {
		Selector(r *ghttp.Request, ups []Weighted) (ref int, ok bool)
	}
)

func New(strategy ...string) (b Balancer) {
	st := StrategyRoundRobin
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
