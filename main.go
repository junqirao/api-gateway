package main

import (
	_ "api-gateway/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"api-gateway/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
