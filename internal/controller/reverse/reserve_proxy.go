package reverse

import (
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/consts"
	"api-gateway/internal/model"
	"api-gateway/internal/service"
)

// Proxy parse routing key
func (c Controller) Proxy(r *ghttp.Request) {
	// set enter time
	r.SetCtxVar(consts.CtxKeyEnterTime, time.Now().UnixNano())
	// reverse proxy
	service.Proxy().Proxy(
		r.GetCtx(),
		&model.ReverseProxyInput{
			RoutingKey: parseRoutingKey(r.Request.RequestURI),
			Request:    r,
		},
	)
}
