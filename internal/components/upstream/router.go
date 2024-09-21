package upstream

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func Router(group *ghttp.RouterGroup) *ghttp.RouterGroup {
	return group.Group("/upstream", func(group *ghttp.RouterGroup) {
		group.Group("/service", func(group *ghttp.RouterGroup) {
			group.GET("/names", management.GetServices)
			group.GET("/state", management.ServiceState)
		})
	})
}
