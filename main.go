package main

import (
	"api-gateway/internal/components"
	_ "api-gateway/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"api-gateway/internal/cmd"
)

func main() {
	ctx := gctx.GetInitCtx()
	// load components
	components.Init(ctx)
	// run
	cmd.Main.Run(ctx)
}
