package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

type GetServiceNamesReq struct {
	g.Meta `path:"/upstream/service/names" method:"get" tags:"Upstream" summary:"Get Upstream Service Names"`
}

type GetServiceNamesRes []string

type GetServiceStateReq struct {
	g.Meta `path:"/upstream/service/state" method:"get" tags:"Upstream" summary:"Get Upstream State By Service"`

	ServiceName string `json:"service_name"`
}

type GetServiceStateRes map[string][]any
