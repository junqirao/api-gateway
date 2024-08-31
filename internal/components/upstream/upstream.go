package upstream

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	registry "github.com/junqirao/simple-registry"
	"github.com/sony/gobreaker"

	"api-gateway/internal/components/breaker"
	"api-gateway/internal/components/config"
	"api-gateway/internal/components/limiter"
	"api-gateway/internal/components/loadbalance"
	"api-gateway/internal/components/response"
	"api-gateway/internal/model"
)

type (
	// Upstream is a reverse proxy target
	Upstream struct {
		registry.Instance
		loadbalance.Weighted

		Parent       *Service
		proxyHandler model.ReverseProxyHandler
		limiter      *limiter.Limiter
		breaker      *breaker.Breaker
		highLoad     *atomic.Bool
	}
)

func NewUpstream(ctx context.Context, instance *registry.Instance) *Upstream {
	cfg, _ := config.GetServiceConfig(instance.ServiceName)
	u := &Upstream{
		Instance: *instance,
		breaker:  breaker.New(cfg.Breaker.Setting(ctx)),
		limiter:  limiter.NewLimiter(cfg.RateLimiter),
		Weighted: loadbalance.NewWeighted(basicLoadBalanceWeight),
		highLoad: &atomic.Bool{},
	}
	u.proxyHandler = NewHandler(ctx, u, cfg.ReverseProxy)
	return u
}

// Allow is a combined entrance of rate limiter and circuit breaker,
// returns limiter allow flag and circuit breaker callback
func (u *Upstream) Allow(ctx context.Context) (cb func(success bool), code *response.Code) {
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
	case err == nil:
	default:
		// 500
		g.Log().Errorf(ctx, "upstream %s breaker error: %v", u.Identity(), err)
		code = response.CodeInternalError.WithDetail(err.Error())
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

// func (u *Upstream) updateConfig(ctx context.Context, scope string, op config.Operation, cfg model.ServiceConfig) {
// 	g.Log().Infof(ctx, "upstream %s update config: scope=%s, op=%v", u.Identity(), scope, op)
// 	switch scope {
// 	case consts.ConfigScopeLimiter:
// 		g.Log().Infof(ctx, "upstream %s limiter -> %v", u.Identity(), cfg.RateLimiter)
// 		u.limiter = limiter.NewLimiter(cfg.RateLimiter)
// 	case consts.ConfigScopeBreaker:
// 		g.Log().Infof(ctx, "upstream %s breaker -> %v", u.Identity(), cfg.Breaker)
// 		u.breaker = breaker.New(cfg.Breaker.Setting(ctx))
// 	}
// }
