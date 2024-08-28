package reverse

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

type (
	// ProxyController ...
	ProxyController interface {
		Proxy(r *ghttp.Request)
	}
	Controller struct {
	}
)

// New ...
func New() ProxyController {
	return &Controller{}
}
