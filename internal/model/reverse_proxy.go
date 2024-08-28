package model

import (
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	// ReverseProxyHandler interface of reverse proxy
	ReverseProxyHandler interface {
		// Do reverse proxy
		Do(ctx context.Context, req *ghttp.Request) (err error)
	}
	// ReverseProxyInput of reverse proxy
	ReverseProxyInput struct {
		RoutingKey string `json:"routing_key"`
		Request    *ghttp.Request
	}
)
