package upstream

import (
	"sort"

	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/components/response"
)

var (
	management = &managementHandler{}
)

type (
	managementHandler struct {
	}
	ManagementGetServiceStateRequest struct {
		ServiceName string `json:"service_name"`
	}
	ManagementGetServiceStateResponse struct {
		Upstreams      int         `json:"upstreams"`
		UpstreamStates []*UpsState `json:"upstream_states"`
	}
)

func (m managementHandler) GetServices(r *ghttp.Request) {
	var ss []string
	cache.m.Range(func(key, _ any) bool {
		ss = append(ss, key.(string))
		return true
	})
	sort.Strings(ss)
	response.WriteData(r, response.CodeDefaultSuccess, ss)
}

func (m managementHandler) ServiceState(r *ghttp.Request) {
	req := new(ManagementGetServiceStateRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	srv, ok := cache.getService(req.ServiceName)
	if !ok {
		response.WriteJSON(r, response.CodeNotFound.WithDetail(req.ServiceName))
		return
	}
	res := &ManagementGetServiceStateResponse{
		Upstreams: len(srv.ups),
	}
	for _, up := range srv.ups {
		res.UpstreamStates = append(res.UpstreamStates, up.State())
	}

	response.WriteData(r, response.CodeDefaultSuccess, res)
}
