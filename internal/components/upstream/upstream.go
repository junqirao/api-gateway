package upstream

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	registry "github.com/junqirao/simple-registry"
	"github.com/sony/gobreaker"

	"api-gateway/internal/components/balancer"
	"api-gateway/internal/components/breaker"
	"api-gateway/internal/components/config"
	"api-gateway/internal/components/limiter"
	"api-gateway/internal/components/response"
)

type (
	// Upstream is a reverse proxy target
	Upstream struct {
		registry.Instance

		// use query per second as measure target
		balancer.Measurable
		balancer.Weighable

		Parent       *Service
		proxyHandler ReverseProxyHandler
		limiter      *limiter.Limiter
		breaker      *breaker.Breaker
		highLoad     *atomic.Bool
	}

	// UpsState upstream state
	UpsState struct {
		HostName     string  `json:"hostname"`
		InstanceId   string  `json:"instance_id"`
		Healthy      bool    `json:"healthy"`
		Weight       int64   `json:"weight"`
		WeightRatio  float64 `json:"weight_ratio"`
		Load         int64   `json:"load"`
		BreakerState string  `json:"breaker_state"`
	}

	// ReverseProxyHandler interface of reverse proxy
	ReverseProxyHandler interface {
		// Do reverse proxy
		Do(ctx context.Context, req *ghttp.Request) (err error)
	}
)

func NewUpstream(ctx context.Context, instance *registry.Instance, cfg config.ServiceConfig) *Upstream {
	var (
		weight int64 = 0
	)

	if w, ok := instance.Meta["weight"]; ok {
		weight = gconv.Int64(w)
	} else {
		weight = defaultWeight
	}

	breakerSetting := cfg.Breaker.Setting(ctx)
	breakerSetting.Name = instance.Identity("_")

	u := &Upstream{
		Instance:   *instance,
		breaker:    breaker.New(breakerSetting),
		limiter:    limiter.NewLimiter(cfg.RateLimiter),
		highLoad:   &atomic.Bool{},
		Measurable: balancer.NewMeasurable(time.Second),
		Weighable:  balancer.NewWeighable(weight),
	}
	u.proxyHandler = NewHandler(ctx, u, cfg.ReverseProxy)
	g.Log().Infof(ctx, "upstream %s created. weight=%d, breaker=%+v, limiter=%+v", u.Identity(), u.Weight(), cfg.Breaker, cfg.RateLimiter)
	return u
}

// Allow is a combined entrance of rate limiter and circuit breaker,
// returns limiter allow flag and circuit breaker callback
func (u *Upstream) Allow(_ context.Context) (cb func(success bool), code *response.Code) {
	// add measurable
	u.AddLoad(1)

	// rate limiter
	if ok := u.limiter.Allow(); !ok {
		// 429
		code = response.CodeTooManyRequests
		u.highLoad.Store(true)
		return
	}

	// circuit breaker
	u.highLoad.Store(false)
	cb, err := u.breaker.Allow()
	switch {
	case errors.Is(err, gobreaker.ErrTooManyRequests):
		// 429
		code = response.CodeTooManyRequests
	case errors.Is(err, gobreaker.ErrOpenState):
		// 503
		code = response.CodeUnavailable
	}
	return
}

// Do proxy request to next layer -> model.ReverseProxyHandler
func (u *Upstream) Do(ctx context.Context, req *ghttp.Request, cb func(success bool)) (err error) {
	if err = u.proxyHandler.Do(ctx, req); err == nil {
		cb(true)
		return
	}

	cb(errors.Is(err, context.Canceled))
	return
}

func (u *Upstream) healthy() bool {
	return u.breaker.State() != gobreaker.StateOpen && !u.highLoad.Load()
}

// State of upstream
func (u *Upstream) State() *UpsState {
	return &UpsState{
		HostName:     u.HostName,
		InstanceId:   u.Instance.Id,
		Healthy:      u.healthy(),
		Weight:       u.Weight(),
		WeightRatio:  float64(u.Weight()) / float64(u.Parent.totalWeight()),
		Load:         u.Load(),
		BreakerState: u.breaker.State().String(),
	}
}
