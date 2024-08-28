package components

import (
	"context"

	"api-gateway/internal/components/grace"
	"api-gateway/internal/components/registry"
)

func Init(ctx context.Context) {
	// grace exit
	grace.Init(ctx)
	// registry
	registry.Init(ctx)
}
