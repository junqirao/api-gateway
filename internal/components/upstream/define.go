package upstream

import (
	"github.com/gogf/gf/v2/net/ghttp"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/breaker"
	"api-gateway/internal/components/limiter"
	"api-gateway/internal/model"
)

type (
	// Service contains Upstream list
	Service struct {
		Ups        []*Upstream
		Config     model.ServiceConfig
		RoutingKey string
	}
	// Upstream is a reverse proxy target
	Upstream struct {
		registry.Instance
		Handler model.ReverseProxyHandler
		limiter *limiter.Limiter
		breaker *breaker.Breaker
	}
	// Selector selects an Upstream from Service
	Selector func(r *ghttp.Request, ss *Service) (*Upstream, bool)
)
