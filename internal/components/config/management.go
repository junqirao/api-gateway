package config

import (
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/response"
	"api-gateway/internal/consts"
	"api-gateway/internal/model"
)

type (
	managementHandler struct {
	}
	ManagementGetConfigRequest struct {
		ServiceName string `json:"service_name"`
	}
	ManagementGetConfigResponse struct {
		Default bool                 `json:"default"`
		Config  *model.ServiceConfig `json:"config"`
	}
	ManagementUpdateConfigRequest struct {
		ServiceName string                 `json:"service_name"`
		Config      map[string]interface{} `json:"config"`
	}
)

func (m *ManagementUpdateConfigRequest) convert(ptr interface{}) error {
	bytes, err := json.Marshal(m.Config)
	if err != nil {
		return err
	}
	return gconv.Struct(bytes, &ptr)
}

func (h *managementHandler) Get(r *ghttp.Request) {
	req := new(ManagementGetConfigRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	config, ok := GetServiceConfig(req.ServiceName)
	response.WriteData(r, response.CodeDefaultSuccess, &ManagementGetConfigResponse{Default: !ok, Config: config})
}

func (h *managementHandler) SetLoadBalanceConfig(r *ghttp.Request) {
	h.setConfig(r, new(model.LoadBalanceConfig), consts.ModuleNameLoadBalance)
}

func (h *managementHandler) SetBreakerConfig(r *ghttp.Request) {
	h.setConfig(r, new(model.BreakerConfig), consts.ModuleNameBreaker)
}

func (h *managementHandler) SetRateLimiterConfig(r *ghttp.Request) {
	h.setConfig(r, new(model.RateLimiterConfig), consts.ModuleNameRateLimiter)
}
func (h *managementHandler) SetReverseProxyConfig(r *ghttp.Request) {
	h.setConfig(r, new(model.ReverseProxyConfig), consts.ModuleNameReverseProxy)
}

func (h *managementHandler) setConfig(r *ghttp.Request, ptr interface{}, module string) {
	req := new(ManagementUpdateConfigRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r,
			response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}
	if req.ServiceName == "" {
		response.WriteJSON(r,
			response.CodeInvalidParameter.WithDetail("service_name can not be empty"))
	}
	if err := req.convert(&ptr); err != nil {
		response.WriteJSON(r,
			response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	err := registry.Storages.GetStorage(StorageNameServiceConfig).
		Set(r.Context(), fmt.Sprintf("%s%s%s", req.ServiceName, StorageSeparator, module), ptr)
	if err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}

	response.WriteJSON(r, response.CodeDefaultSuccess)
}
