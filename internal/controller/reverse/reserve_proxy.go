package reverse

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/components/utils"
	"api-gateway/internal/model"
	"api-gateway/internal/service"
)

// Proxy parse routing key
func (c Controller) Proxy(r *ghttp.Request) {
	// reverse proxy
	service.Proxy().Proxy(
		r.GetCtx(),
		&model.ReverseProxyInput{
			RoutingKey: utils.ParseRoutingKey(r.Request.RequestURI),
			Request:    r,
		},
	)
}
