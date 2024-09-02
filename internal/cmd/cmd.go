package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/program"
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
			if g.Cfg().MustGet(ctx, "server.debug", false).Bool() {
				g.Log().Info(ctx, "pprof enabled")
				s.EnablePProf()
			}
			s.BindHandler(fmt.Sprintf("%s*", pattern), reverse.New().Proxy)

			// management
			s.Group("/management", func(group *ghttp.RouterGroup) {
				// config
				config.Router(group)
				// program
				program.Router(group)
			})
			s.Run()
			return nil
		},
	}
)
