package config

import (
	"fmt"

	"github.com/gogf/gf/v2/net/ghttp"
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
	ManagementSetLoadBalanceConfigRequest struct {
		ServiceName string                   `json:"service_name"`
		Config      *model.LoadBalanceConfig `json:"config"`
	}
	ManagementSetBreakerConfigRequest struct {
		ServiceName string               `json:"service_name"`
		Config      *model.BreakerConfig `json:"config"`
	}
	ManagementSetRateLimiterConfigRequest struct {
		ServiceName string                   `json:"service_name"`
		Config      *model.RateLimiterConfig `json:"config"`
	}
)

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
	req := new(ManagementSetLoadBalanceConfigRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	h.setConfig(r, req.ServiceName, consts.ModuleNameLoadBalance, req.Config)
}

func (h *managementHandler) SetBreakerConfig(r *ghttp.Request) {
	req := new(ManagementSetBreakerConfigRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	h.setConfig(r, req.ServiceName, consts.ModuleNameBreaker, req.Config)
}

func (h *managementHandler) SetRateLimiterConfig(r *ghttp.Request) {
	req := new(ManagementSetRateLimiterConfigRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	h.setConfig(r, req.ServiceName, consts.ModuleNameRateLimiter, req.Config)
}

func (h *managementHandler) setConfig(r *ghttp.Request, serviceName, module string, config any) {
	err := registry.Storages.GetStorage(StorageNameServiceConfig).
		Set(r.Context(), fmt.Sprintf("%s%s%s", serviceName, StorageSeparator, module), config)
	if err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}

	response.WriteJSON(r, response.CodeDefaultSuccess)
}
