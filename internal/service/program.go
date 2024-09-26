// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"api-gateway/internal/components/program"
	"api-gateway/internal/model"
	"context"
)

type (
	IProgramManagement interface {
		GetProgramInfo(ctx context.Context, serviceName string) (res map[string][]*program.Info, err error)
		SetProgramInfo(ctx context.Context, info *program.Info) (err error)
		DeleteProgramInfo(ctx context.Context, input model.DeleteProgramInfoInput) (err error)
		// GetGlobalVariables get global variable
		GetGlobalVariables(ctx context.Context) map[string]interface{}
		// SetGlobalVariables set global variable
		SetGlobalVariables(ctx context.Context, input model.SetGlobalVariablesInput) error
		// DeleteGlobalVariables delete global variable
		DeleteGlobalVariables(ctx context.Context, key string) error
	}
)

var (
	localProgramManagement IProgramManagement
)

func ProgramManagement() IProgramManagement {
	if localProgramManagement == nil {
		panic("implement not found for interface IProgramManagement, forgot register?")
	}
	return localProgramManagement
}

func RegisterProgramManagement(i IProgramManagement) {
	localProgramManagement = i
}
