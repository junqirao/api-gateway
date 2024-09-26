package upstream

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"api-gateway/internal/components/breaker"
	"api-gateway/internal/components/limiter"
	"api-gateway/internal/consts"
	"api-gateway/internal/model"
)

func (u *Upstream) updateConfig(ctx context.Context, module string) {
	// local caches in registry always updated before this function called
	cfg := u.Parent.Config
	switch module {
	case consts.ModuleNameRateLimiter:
		if !model.ValueChanged(cfg.RateLimiter, u.Parent.Config.RateLimiter) {
			g.Log().Infof(ctx, "upstream %s rate limiter not changed", u.Identity())
			return
		}

		u.limiter.SetLimit(limiter.NewLimit(cfg.RateLimiter.Rate))
		u.limiter.SetBurst(cfg.RateLimiter.Peak)
		g.Log().Infof(ctx, "upstream %s limiter -> %+v", u.Identity(), cfg.RateLimiter)
	case consts.ModuleNameBreaker:
		if !model.ValueChanged(cfg.Breaker, u.Parent.Config.Breaker) {
			g.Log().Infof(ctx, "upstream %s breaker not changed", u.Identity())
			return
		}

		u.breaker = breaker.New(cfg.Breaker.Setting(ctx))
		g.Log().Infof(ctx, "upstream %s breaker -> %+v", u.Identity(), cfg.Breaker)
	}
}
