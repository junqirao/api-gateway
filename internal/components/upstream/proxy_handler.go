package upstream

import (
	"context"

	"api-gateway/internal/model"
)

func NewHandler(_ context.Context, ups *Upstream, cfg model.ReverseProxyConfig) (handler model.ReverseProxyHandler) {
	// todo distinguish handler type by registry.Instance
	handler = newHTTPHandler(ups, &cfg)
	return
}
