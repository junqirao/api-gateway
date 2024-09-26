package limiter

import (
	"golang.org/x/time/rate"

	"api-gateway/internal/components/config"
)

type Limiter struct {
	cfg config.RateLimiterConfig
	*rate.Limiter
}

func NewLimiter(cfg config.RateLimiterConfig) *Limiter {
	return &Limiter{
		Limiter: rate.NewLimiter(NewLimit(cfg.Rate), cfg.Peak),
	}
}

func NewLimit(rateValue float64) rate.Limit {
	r := rate.Inf
	if rateValue > 0 {
		r = rate.Limit(rateValue)
	}
	return r
}
