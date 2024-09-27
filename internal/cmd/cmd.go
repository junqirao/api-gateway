package cmd

import (
	"context"

	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "server",
		Usage: "server",
		Brief: "start http and grpc server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			go runGRPCServer(ctx)
			runHttpSrvBlock(ctx)
			return nil
		},
	}
)
