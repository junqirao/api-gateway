package extra

import (
	"context"

	"api-gateway/internal/components/program/extra/ipgeo"
)

func Init(ctx context.Context) {
	ipgeo.Init(ctx)
}
