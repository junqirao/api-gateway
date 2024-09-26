package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"api-gateway/internal/components/upstream"
)

type GetServiceNamesReq struct {
	g.Meta `path:"/upstream/service/names" method:"get" tags:"Upstream" summary:"Get Upstream Service Names"`
}

type GetServiceNamesRes []string

type GetServiceStateReq struct {
	g.Meta `path:"/upstream/service/state" method:"get" tags:"Upstream" summary:"Get Upstream State By Service"`

	ServiceName string `json:"service_name"`
}

type GetServiceStateRes struct {
	Upstreams      int                  `json:"upstreams"`
	UpstreamStates []*upstream.UpsState `json:"upstream_states"`
}
