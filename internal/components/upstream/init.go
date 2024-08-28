package upstream

import (
	"context"
)

func Init(ctx context.Context) {
	cache = newUpstreamCache(ctx)
}
