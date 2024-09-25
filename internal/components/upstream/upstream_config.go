package upstream

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"api-gateway/internal/components/breaker"
	"api-gateway/internal/components/config"
	"api-gateway/internal/components/limiter"
	"api-gateway/internal/consts"
	"api-gateway/internal/model"
)

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
