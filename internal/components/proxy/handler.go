package proxy

import (
	"context"

	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/model"
)

func NewHandler(_ context.Context, ins *registry.Instance, cfg model.ReverseProxyConfig) (handler model.ReverseProxyHandler) {
	// todo distinguish handler type by registry.Instance
	handler = newHTTPHandler(ins, &cfg)
	return
}
