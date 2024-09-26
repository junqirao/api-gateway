package model

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/google/uuid"
	"github.com/sony/gobreaker"
)

type (
	// ReverseProxyConfig reverse proxy config
	ReverseProxyConfig struct {
		// TrimRoutingKeyPrefix trim routing key prefix, default true
		TrimRoutingKeyPrefix bool `json:"trim_routing_key_prefix"` // default true
		// RetryCount retry count, default 1
		RetryCount int `json:"retry_count"` // default 1
		// DialTimeout, default 1s
		DialTimeout string `json:"dial_timeout,omitempty"`
		// TlsHandshakeTimeout tls handshake timeout, default 1s
		TlsHandshakeTimeout string `json:"tls_handshake_timeout,omitempty"`
		// Scheme http or https
		Scheme string `json:"scheme,omitempty"`
	}
	// LoadBalanceConfig load balance config
	LoadBalanceConfig struct {
		// Strategy load balance strategy: random, round_robin
		Strategy string `json:"strategy"` // default random
	}
	// RateLimiterConfig rate limiter config
	RateLimiterConfig struct {
		// Rate rate
		Rate float64 `json:"rate"`
		// Peak of burst
		Peak int `json:"peak"`
	}
	// GatewayConfig gateway config
	GatewayConfig struct {
		// Prefix route prefix
		Prefix string `json:"prefix"`
		// Debug debug mode, default false, equals server.debug
		Debug bool
	}
	// BreakerConfig breaker config
	BreakerConfig struct {
		// Name breaker name, if not set auto generate
		Name string `json:"name"`
		// MaxFailures max failures of closed state to half open
		MaxFailures int `json:"max_failures"`
		// HalfOpenMaxRequests max requests of half open state to open or close state
		HalfOpenMaxRequests int `json:"half_open_max_requests"`
		// OpenTimeout the period of the open state, default 1m
		OpenTimeout string `json:"open_timeout"`
		// Interval the cyclic period of the closed state
		// for the CircuitBreaker to clear the internal Counts, default 1m
		Interval string `json:"interval"`
	}
	// ServiceConfig service config
	ServiceConfig struct {
		// ReverseProxy reverse proxy config
		ReverseProxy ReverseProxyConfig `json:"reverse_proxy"`
		// LoadBalance load balance config
		LoadBalance LoadBalanceConfig `json:"load_balance"`
		// RateLimiter rate limiter config
		RateLimiter RateLimiterConfig `json:"rate_limiter"`
		// Breaker breaker config
		Breaker BreakerConfig `json:"breaker"`
	}
)

// Setting get gobreaker.Settings
func (c BreakerConfig) Setting(ctx context.Context) gobreaker.Settings {
	var (
		name        = c.Name
		maxFailures = uint32(c.MaxFailures)
		maxRequests = uint32(c.HalfOpenMaxRequests)
		interval    = time.Second * 60
		timeout     = time.Second * 60
	)

	if name == "" {
		name = fmt.Sprintf("breaker_%s", uuid.New().String())
	}
	if maxFailures == 0 {
		maxFailures = 5
	}
	if maxRequests == 0 {
		maxRequests = 1
	}
	if c.Interval != "" {
		duration, err := time.ParseDuration(c.Interval)
		if err == nil {
			interval = duration
		} else {
			g.Log().Warningf(ctx, "breaker parse interval '%s' error: %v", c.Interval, err)
		}
	}
	if c.OpenTimeout != "" {
		duration, err := time.ParseDuration(c.OpenTimeout)
		if err == nil {
			timeout = duration
		} else {
			g.Log().Warningf(ctx, "breaker parse timeout '%s' error: %v", c.OpenTimeout, err)
		}
	}

	return gobreaker.Settings{
		Name:        name,
		MaxRequests: maxRequests,
		Interval:    interval,
		Timeout:     timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			return counts.ConsecutiveFailures >= maxFailures
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			g.Log().Infof(ctx, "breaker %s state changed from %s to %s", name, from, to)
		},
	}
}

func (s *ServiceConfig) Clone() *ServiceConfig {
	return &ServiceConfig{
		ReverseProxy: s.ReverseProxy,
		LoadBalance:  s.LoadBalance,
		RateLimiter:  s.RateLimiter,
		Breaker:      s.Breaker,
	}
}

func ValueChanged(v1, v2 any) bool {
	return gmd5.MustEncryptString(gconv.String(v1)) != gmd5.MustEncryptString(gconv.String(v2))
}
