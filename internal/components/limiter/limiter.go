package limiter

import (
	"golang.org/x/time/rate"

	"api-gateway/internal/model"
)

type Limiter struct {
	cfg model.RateLimiterConfig
	*rate.Limiter
}

func NewLimiter(cfg model.RateLimiterConfig) *Limiter {
	r := rate.Inf
	if cfg.Rate > 0 {
		r = rate.Limit(cfg.Rate)
	}
	return &Limiter{
		Limiter: rate.NewLimiter(r, cfg.Peak),
	}
}
