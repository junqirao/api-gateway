package components

import (
	"context"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/registry"
	"api-gateway/internal/components/upstream"
)

func Init(ctx context.Context) {
	// registry
	registry.Init(ctx)
	// biz config init
	config.Init(ctx)
	// upstream management
	upstream.Init(ctx)
}
