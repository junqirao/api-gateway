package upstream

import (
	"context"
)

func Init(ctx context.Context) {
	Cache = newUpstreamCache(ctx)
}
