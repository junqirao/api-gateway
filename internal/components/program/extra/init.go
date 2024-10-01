package extra

import (
	"context"

	"api-gateway/internal/components/program/extra/ipgeo"
	"api-gateway/internal/components/program/extra/jwt"
)

func Init(ctx context.Context) {
	ipgeo.Init(ctx)
	jwt.Init(ctx)
}
