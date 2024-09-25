package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/program"
	"api-gateway/internal/components/prometheus"
	"api-gateway/internal/components/response"
	"api-gateway/internal/components/upstream"
	"api-gateway/internal/consts"
	"api-gateway/internal/controller/reverse"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			// prepare
			pattern := config.Gateway.Prefix
			if !strings.HasSuffix(pattern, "/") {
				pattern += "/"
			}
			pattern = fmt.Sprintf("%s*", pattern)

			s := g.Server()
			if g.Cfg().MustGet(ctx, "server.debug", false).Bool() {
				g.Log().Info(ctx, "pprof enabled")
				s.EnablePProf()
			}

			// max body size, default 512M
			s.SetClientMaxBodySize(g.Cfg().MustGet(ctx, "server.max_body_size", consts.DefaultMaxBodySize).Int64())

			// middleware
			s.BindMiddleware(pattern, prometheus.Middleware)
			// reverse
			s.BindHandler(pattern, reverse.New().Proxy)

			// management
			if g.Cfg().MustGet(ctx, "gateway.management.enable", true).Bool() {
				entrance := g.Cfg().MustGet(ctx, "gateway.management.entrance", "").String()
				if entrance != "" && !strings.HasPrefix(entrance, "/") {
					entrance = "/" + entrance
				}
				s.Group(fmt.Sprintf("%s/management", entrance), func(group *ghttp.RouterGroup) {
					// auth middleware
					if pwd := g.Cfg().MustGet(ctx, "gateway.management.password").String(); pwd != "" {
						group.Middleware(func(r *ghttp.Request) {
							if gmd5.MustEncryptString(r.GetHeader("Authorization")) != pwd {
								response.WriteJSON(r, response.CodeNotFound)
								return
							}
							r.Middleware.Next()
						})
					}

					// config
					config.Router(group)
					// program
					program.Router(group)
					// upstream
					upstream.Router(group)
					// prometheus
					group.ALL("/metrics", func(r *ghttp.Request) {
						promhttp.Handler().ServeHTTP(r.Response.RawWriter(), r.Request)
					})
				})
			}

			s.Run()
			return nil
		},
	}
)
