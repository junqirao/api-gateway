package upstream

import (
	"context"

	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/mirror"
	"api-gateway/internal/components/proxy"
)

func Init(ctx context.Context) {
	Cache = newUpstreamCache(ctx)
	mirror.Init(ctx, func(ins *registry.Instance) (proxy.ReverseProxyHandler, error) {
		return newHTTPHandler(ins.Host, ins.Port, ins.ServiceName, func() *config.ReverseProxyConfig {
			return &config.ReverseProxyConfig{
				DialTimeout: "500ms",
			}
		}), nil
	})
}
