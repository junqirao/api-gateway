package config

import (
	"context"
	"errors"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/consts"
)

var (
	// Gateway global config
	Gateway *GatewayConfig
	// StorageSeparator same as registry.storage.separator
	StorageSeparator = g.Cfg().MustGet(context.Background(), "registry.storage.separator", "/").String()
	// StorageNameServiceConfig name
	StorageNameServiceConfig = "service_config"
)

// default value define
var (
	defaultGatewayPrefix = "/api/"
	defaultRPConfig      = ReverseProxyConfig{
		TrimRoutingKeyPrefix: true,
		RetryCount:           1,
	}
	defaultLBConfig = LoadBalanceConfig{
		Strategy: "random",
	}
	defaultBreakerConfig = BreakerConfig{
		MaxFailures:         5,
		HalfOpenMaxRequests: 1,
		OpenTimeout:         "1m",
		Interval:            "1m",
	}
	defaultServiceConfig = &ServiceConfig{
		ReverseProxy: defaultRPConfig,
		LoadBalance:  defaultLBConfig,
		Breaker:      defaultBreakerConfig,
	}
)

func GetServiceConfig(serviceName string) (*ServiceConfig, bool) {
	ctx := context.Background()
	kvs, err := registry.Storages.GetStorage(StorageNameServiceConfig).Get(ctx, serviceName)
	switch {
	case errors.Is(err, registry.ErrStorageNotFound):
	case err == nil:
	default:
		g.Log().Infof(ctx, "failed to get service %s config, using default. result=%v", StorageNameServiceConfig, err)
		return defaultServiceConfig.Clone(), false
	}
	if len(kvs) == 0 {
		return defaultServiceConfig.Clone(), false
	}

	cfg := defaultServiceConfig.Clone()
	for _, kv := range kvs {
		parts := strings.Split(kv.Key, StorageSeparator)
		if len(parts) < 1 {
			// drop invalid key
			continue
		}
		moduleName := parts[len(parts)-1]
		switch moduleName {
		case consts.ModuleNameLoadBalance:
			err = kv.Value.Scan(&cfg.LoadBalance)
		case consts.ModuleNameBreaker:
			err = kv.Value.Scan(&cfg.Breaker)
		case consts.ModuleNameRateLimiter:
			err = kv.Value.Scan(&cfg.RateLimiter)
		case consts.ModuleNameReverseProxy:
			err = kv.Value.Scan(&cfg.ReverseProxy)
		}
		if err != nil {
			g.Log().Errorf(ctx, "failed to scan service %s config, using default config: %s", serviceName, err.Error())
			return defaultServiceConfig.Clone(), false
		}
	}

	return cfg, true
}

func loadConfigs(ctx context.Context) {
	// load gateway config
	if err := g.Cfg().MustGet(ctx, "gateway",
		&GatewayConfig{Prefix: defaultGatewayPrefix}).Scan(&Gateway); err != nil {
		g.Log().Errorf(ctx, "load gateway config error: %v", err)
	}
	// gateway.debug follow server.debug
	if !Gateway.Debug && g.Cfg().MustGet(ctx, "server.debug", false).Bool() {
		Gateway.Debug = true
	}
	if Gateway.Debug {
		g.Log().Info(ctx, "gateway debug mode enabled")
	}
}
