package model

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type (

	// ReverseProxyInput of reverse proxy
	ReverseProxyInput struct {
		RoutingKey string `json:"routing_key"`
		Request    *ghttp.Request
	}
)
