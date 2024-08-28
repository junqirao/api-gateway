package config

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"

	"api-gateway/internal/model"
)

var (
	// Gateway global config
	Gateway *model.GatewayConfig
	cache   sync.Map // service_name:*model.ServiceConfig
)

// default value define
var (
	defaultGatewayPrefix = "/api/"
	defaultRPConfig      = model.ReverseProxyConfig{
		TrimRoutingKeyPrefix: true,
		RetryCount:           1,
	}
	defaultLBConfig = model.LoadBalanceConfig{
		Strategy: "random",
	}
	defaultBreakerConfig = model.BreakerConfig{
		MaxFailures:         5,
		HalfOpenMaxRequests: 1,
		OpenTimeout:         "1m",
		Interval:            "1m",
	}
	defaultServiceConfig = &model.ServiceConfig{
		ReverseProxy: defaultRPConfig,
		LoadBalance:  defaultLBConfig,
		Breaker:      defaultBreakerConfig,
	}
)

func GetServiceConfig(serviceName string) (*model.ServiceConfig, bool) {
	if v, ok := cache.Load(serviceName); ok {
		return v.(*model.ServiceConfig), true
	}
	return defaultServiceConfig.Clone(), false
}

func loadConfigs(ctx context.Context) {
	if err := g.Cfg().MustGet(ctx, "gateway",
		&model.GatewayConfig{Prefix: defaultGatewayPrefix}).Scan(&Gateway); err != nil {
		g.Log().Errorf(ctx, "load gateway config error: %v", err)
	}
	services := make(map[string]*model.ServiceConfig)
	if err := g.Cfg().MustGet(ctx, "services",
		&services).Scan(&services); err != nil {
		g.Log().Errorf(ctx, "load service config error: %v", err)
	}
	for k, v := range services {
		cache.Store(k, v)
	}
	g.Log().Infof(ctx, "service config loaded: %d", len(services))
}
