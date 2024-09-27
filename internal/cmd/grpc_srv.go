package cmd

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"

	"api-gateway/internal/controller/inner"
)

func runGRPCServer(_ context.Context) {
	s := grpcx.Server.New()
	// register inner service
	inner.Register(s)
	s.Run()
	return
}
