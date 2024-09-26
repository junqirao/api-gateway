package program

import (
	"context"

	"api-gateway/api/program/v1"
	"api-gateway/internal/model"
	"api-gateway/internal/service"
)

func (c *ControllerV1) GetProgramInfo(ctx context.Context, req *v1.GetProgramInfoReq) (res *v1.GetProgramInfoRes, err error) {
	v, err := service.ProgramManagement().GetProgramInfo(ctx, req.ServiceName)
	if err != nil {
		return
	}
	vv := v1.GetProgramInfoRes(v)
	res = &vv
	return
}
func (c *ControllerV1) SetProgramInfo(ctx context.Context, req *v1.SetProgramInfoReq) (res *v1.SetProgramInfoRes, err error) {
	err = service.ProgramManagement().SetProgramInfo(ctx, &req.Info)
	return
}
func (c *ControllerV1) DeleteProgramInfo(ctx context.Context, req *v1.DeleteProgramInfoReq) (res *v1.DeleteProgramInfoRes, err error) {
	err = service.ProgramManagement().DeleteProgramInfo(ctx, model.DeleteProgramInfoInput{
		ServiceName: req.ServiceName,
		Name:        req.Name,
	})
	return
}
func (c *ControllerV1) GetGlobalVariables(ctx context.Context, _ *v1.GetGlobalVariablesReq) (res *v1.GetGlobalVariablesRes, err error) {
	v := v1.GetGlobalVariablesRes(service.ProgramManagement().GetGlobalVariables(ctx))
	res = &v
	return
}
func (c *ControllerV1) SetGlobalVariables(ctx context.Context, req *v1.SetGlobalVariablesReq) (res *v1.SetGlobalVariablesRes, err error) {
	err = service.ProgramManagement().SetGlobalVariables(ctx, model.SetGlobalVariablesInput{
		Key:   req.Key,
		Value: req.Value,
	})
	return
}
func (c *ControllerV1) DeleteGlobalVariables(ctx context.Context, req *v1.DeleteGlobalVariablesReq) (res *v1.DeleteGlobalVariablesRes, err error) {
	err = service.ProgramManagement().DeleteGlobalVariables(ctx, req.Key)
	return
}
