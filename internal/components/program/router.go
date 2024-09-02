package program

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func Router(group *ghttp.RouterGroup) *ghttp.RouterGroup {
	return group.Group("/program", func(group *ghttp.RouterGroup) {
		group.Group("/variable", func(group *ghttp.RouterGroup) {
			group.GET("/", management.GetGlobalVariables)
			group.PUT("/", management.SetGlobalVariables)
			group.DELETE("/", management.DeleteGlobalVariables)
		})
		group.Group("/info", func(group *ghttp.RouterGroup) {
			group.GET("/", management.GetProgram)
			group.DELETE("/", management.DeleteProgram)
			group.PUT("/", management.SetProgram)
		})
	})
}
