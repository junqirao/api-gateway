package config

import (
	"context"

	"api-gateway/api/config/v1"
	"api-gateway/internal/model"
	"api-gateway/internal/service"
)

func (c *ControllerV1) GetConfig(ctx context.Context, req *v1.GetConfigReq) (res *v1.GetConfigRes, err error) {
	output, err := service.ConfigManagement().GetConfig(ctx, model.GetConfigInput{
		ServiceName: req.ServiceName,
	})
	if err != nil {
		return
	}
	return &v1.GetConfigRes{Config: output.Config, Default: output.Default}, nil
}
func (c *ControllerV1) UpdateConfig(ctx context.Context, req *v1.UpdateConfigReq) (res *v1.UpdateConfigRes, err error) {
	err = service.ConfigManagement().UpdateConfig(ctx, model.UpdateConfigInput{
		Model:       req.Module,
		ServiceName: req.ServiceName,
		Config:      req.Config,
	})
	return
}
