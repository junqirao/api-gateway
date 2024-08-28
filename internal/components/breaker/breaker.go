package breaker

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/sony/gobreaker"
)

type Breaker struct {
	*gobreaker.TwoStepCircuitBreaker
}

var (
	defaultConfig = gobreaker.Settings{
		Name:        "default",
		MaxRequests: 1,
		Interval:    time.Second * 60,
		Timeout:     time.Second * 60,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= 5
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			g.Log().Infof(context.TODO(), "breaker %s state changed from %s to %s", name, from, to)
		},
	}
)

func New(st ...gobreaker.Settings) *Breaker {
	setting := defaultConfig
	if len(st) > 0 {
		setting = st[0]
	}
	return &Breaker{TwoStepCircuitBreaker: gobreaker.NewTwoStepCircuitBreaker(setting)}
}
