package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"api-gateway/internal/components/config"
)

type GetConfigReq struct {
	g.Meta        `path:"/config" method:"get" tags:"Config" summary:"Get Config"`
	Authorization string `json:"Authorization" in:"header"`

	ServiceName string `json:"service_name"`
}
type GetConfigRes struct {
	Default bool                  `json:"default"`
	Config  *config.ServiceConfig `json:"config"`
}
type UpdateConfigReq struct {
	g.Meta        `path:"/config/:module" method:"put" tags:"Config" summary:"Update Config"`
	Authorization string `json:"Authorization" in:"header"`

	Module      string                 `json:"module" in:"path"`
	ServiceName string                 `json:"service_name"`
	Config      map[string]interface{} `json:"config"`
}
type UpdateConfigRes struct{}
