package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"api-gateway/internal/components/authentication"
	cfg "api-gateway/internal/components/config"
	"api-gateway/internal/components/prometheus"
	"api-gateway/internal/components/response"
	"api-gateway/internal/consts"
	"api-gateway/internal/controller/config"
	"api-gateway/internal/controller/program"
	"api-gateway/internal/controller/reverse"
	"api-gateway/internal/controller/upstream"
)

func runHttpSrvBlock(ctx context.Context) {
	// prepare
	pattern := cfg.Gateway.Prefix
	if !strings.HasSuffix(pattern, "/") {
		pattern += "/"
	}
	pattern = fmt.Sprintf("%s*", pattern)

	s := g.Server()
	s.SetServerAgent(g.Cfg().MustGet(ctx, "server.agent", consts.DefaultServerAgent).String())
	debug := g.Cfg().MustGet(ctx, "server.debug", false).Bool()
	s.SetDumpRouterMap(debug)
	if debug {
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
			group.Middleware(authentication.L.Middleware)

			group.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(response.Middleware)
				group.Bind(
					config.NewV1(),   // config
					program.NewV1(),  // program
					upstream.NewV1(), // upstream
				)
			})
			// prometheus
			group.ALL("/metrics", func(r *ghttp.Request) {
				promhttp.Handler().ServeHTTP(r.Response.RawWriter(), r.Request)
			})
		})
	}

	s.Run()
}
