package upstream

import (
	"context"
	"errors"
	"sync/atomic"

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
	"api-gateway/internal/consts"
	"api-gateway/internal/model"
)

type (
	// Upstream is a reverse proxy target
	Upstream struct {
		registry.Instance

		balancer.Measurable
		balancer.Weighable

		Parent       *Service
		proxyHandler model.ReverseProxyHandler
		limiter      *limiter.Limiter
		breaker      *breaker.Breaker
		highLoad     *atomic.Bool
	}
)

func NewUpstream(ctx context.Context, instance *registry.Instance, cfg model.ServiceConfig) *Upstream {
	var weight int64 = 0
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
		Measurable: balancer.NewMeasurable(),
		Weighable:  balancer.NewWeighable(weight),
	}
	u.proxyHandler = NewHandler(ctx, u, cfg.ReverseProxy)
	g.Log().Infof(ctx, "upstream %s created. weight=%d, breaker=%+v, limiter=%+v", u.Identity(), u.Weight(), cfg.Breaker, cfg.RateLimiter)
	return u
}

// Allow is a combined entrance of rate limiter and circuit breaker,
// returns limiter allow flag and circuit breaker callback
func (u *Upstream) Allow(_ context.Context) (cb func(success bool), code *response.Code) {
	if ok := u.limiter.Allow(); !ok {
		// 429
		code = response.CodeTooManyRequests
		u.highLoad.Store(true)
		return
	}
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

func (u *Upstream) updateConfig(module string) {
	ctx := context.Background()
	// local caches in registry always updated before this function called
	cfg, _ := config.GetServiceConfig(u.ServiceName)
	switch module {
	case consts.ModuleNameRateLimiter:
		if !model.ValueChanged(cfg.RateLimiter, u.Parent.Config.RateLimiter) {
			g.Log().Infof(ctx, "upstream %s rate limiter not changed", u.Identity())
			return
		}
		u.limiter = limiter.NewLimiter(cfg.RateLimiter)
		g.Log().Infof(ctx, "upstream %s limiter -> %+v", u.Identity(), cfg.RateLimiter)
		u.Parent.Config.RateLimiter = cfg.RateLimiter
	case consts.ModuleNameBreaker:
		if !model.ValueChanged(cfg.Breaker, u.Parent.Config.Breaker) {
			g.Log().Infof(ctx, "upstream %s breaker not changed", u.Identity())
			return
		}
		u.breaker = breaker.New(cfg.Breaker.Setting(ctx))
		g.Log().Infof(ctx, "upstream %s breaker -> %+v", u.Identity(), cfg.Breaker)
		u.Parent.Config.Breaker = cfg.Breaker
	}
}
