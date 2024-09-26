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
