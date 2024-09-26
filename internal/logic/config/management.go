package config

import (
	"context"
	"fmt"

	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/response"
	"api-gateway/internal/model"
	"api-gateway/internal/service"
)

type (
	sConfigManagement struct{}
)

func init() {
	service.RegisterConfigManagement(&sConfigManagement{})
}

func (s sConfigManagement) GetConfig(_ context.Context, input model.GetConfigInput) (model.GetConfigOutput, error) {
	res := model.GetConfigOutput{}
	res.Config, res.Default = config.GetServiceConfig(input.ServiceName)
	return res, nil
}

func (s sConfigManagement) UpdateConfig(ctx context.Context, input model.UpdateConfigInput) error {
	var ptr = input.PtrByModule()
	if ptr == nil {
		return response.CodeInvalidParameter.WithDetail(fmt.Sprintf("invalid module: %s", input.Model))
	}
	if input.ServiceName == "" {
		return response.CodeInvalidParameter.WithDetail("service_name can not be empty")
	}
	if err := input.Convert(&ptr); err != nil {
		return response.CodeInvalidParameter.WithDetail(err.Error())
	}

	return registry.Storages.GetStorage(config.StorageNameServiceConfig).
		Set(ctx, fmt.Sprintf("%s%s%s", input.ServiceName, config.StorageSeparator, input.Model), ptr)
}
