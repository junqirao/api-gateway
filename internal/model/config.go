package model

import (
	"encoding/json"

	"github.com/gogf/gf/v2/util/gconv"

	"api-gateway/internal/components/config"
	"api-gateway/internal/consts"
)

type (
	GetConfigInput struct {
		ServiceName string `json:"service_name"`
	}
	GetConfigOutput struct {
		Default bool                  `json:"default"`
		Config  *config.ServiceConfig `json:"config"`
	}
	UpdateConfigInput struct {
		Model       string                 `json:"model"`
		ServiceName string                 `json:"service_name"`
		Config      map[string]interface{} `json:"config"`
	}
)

func (m *UpdateConfigInput) Convert(ptr interface{}) error {
	bytes, err := json.Marshal(m.Config)
	if err != nil {
		return err
	}
	return gconv.Struct(bytes, &ptr)
}

func (m *UpdateConfigInput) PtrByModule() interface{} {
	var ptr interface{}
	switch m.Model {
	case consts.ModuleNameBreaker:
		ptr = &config.BreakerConfig{}
	case consts.ModuleNameRateLimiter:
		ptr = &config.RateLimiterConfig{}
	case consts.ModuleNameLoadBalance:
		ptr = &config.LoadBalanceConfig{}
	case consts.ModuleNameReverseProxy:
		ptr = &config.ReverseProxyConfig{}
	default:
	}
	return ptr
}
