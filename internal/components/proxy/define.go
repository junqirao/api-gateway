package proxy

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
)
