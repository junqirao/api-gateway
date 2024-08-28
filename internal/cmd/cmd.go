package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"

	"api-gateway/internal/components/config"
	"api-gateway/internal/controller/reverse"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			pattern := config.Gateway.Prefix
			if !strings.HasSuffix(pattern, "/") {
				pattern += "/"
			}
			s := g.Server()
			s.BindHandler(fmt.Sprintf("%s*", pattern), reverse.New().Proxy)
			s.Run()
			return nil
		},
	}
)
