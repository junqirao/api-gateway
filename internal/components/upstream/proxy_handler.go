package upstream

import (
	"context"

	"api-gateway/internal/components/config"
)

func NewHandler(_ context.Context, ups *Upstream, cfg config.ReverseProxyConfig) (handler ReverseProxyHandler) {
	// todo distinguish handler type by registry.Instance
	handler = newHTTPHandler(ups, &cfg)
	return
}
