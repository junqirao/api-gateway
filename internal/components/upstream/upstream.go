package upstream

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	registry "github.com/junqirao/simple-registry"
	"github.com/sony/gobreaker"

	"api-gateway/internal/components/breaker"
	"api-gateway/internal/components/config"
	"api-gateway/internal/components/limiter"
	"api-gateway/internal/components/proxy"
)

func GetService(routingKey string) (s *Service) {
	s, _ = cache.GetService(routingKey)
	return
}

func NewUpstream(ctx context.Context, instance *registry.Instance) *Upstream {
	cfg, _ := config.GetServiceConfig(instance.ServiceName)
	u := &Upstream{
		Instance: *instance,
		breaker:  breaker.New(cfg.Breaker.Setting(ctx)),
		limiter:  limiter.NewLimiter(cfg.RateLimiter),
	}
	u.Handler = proxy.NewHandler(ctx, instance, cfg.ReverseProxy)
	return u
}

// Allow is a combined entrance of rate limiter and circuit breaker,
// returns limiter allow flag and circuit breaker callback
func (u *Upstream) Allow(ctx context.Context) (ok bool, cb func(success bool)) {
	ok = u.limiter.Allow()
	if !ok {
		return
	}
	cb, err := u.breaker.Allow()
	if err != nil {
		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			g.Log().Infof(ctx, "upstream %s breaker: %v", u.Identity(), err.Error())
		} else {
			g.Log().Errorf(ctx, "upstream %s breaker error: %v", u.Identity(), err)
		}
		return false, nil
	}
	return
}

// Do proxy request to next layer -> model.ReverseProxyHandler
func (u *Upstream) Do(ctx context.Context, req *ghttp.Request, cb func(success bool)) (retry bool, err error) {
	if err = u.Handler.Do(ctx, req); err == nil {
		cb(true)
		return
	}

	success := errors.Is(err, context.Canceled)
	cb(success)
	retry = !success
	return
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
