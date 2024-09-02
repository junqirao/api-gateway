package config

import (
	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/consts"
)

func Router(group *ghttp.RouterGroup) *ghttp.RouterGroup {
	return group.Group("/config", func(group *ghttp.RouterGroup) {
		// get
		group.GET("/", management.Get)
		// update
		group.Group("/", func(group *ghttp.RouterGroup) {
			group.PUT(consts.ModuleNameLoadBalance, management.SetLoadBalanceConfig)
			group.PUT(consts.ModuleNameBreaker, management.SetBreakerConfig)
			group.PUT(consts.ModuleNameRateLimiter, management.SetRateLimiterConfig)
		})
	})
}
