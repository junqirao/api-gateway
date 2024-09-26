// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"api-gateway/internal/model"
	"context"
)

type (
	IConfigManagement interface {
		GetConfig(_ context.Context, input model.GetConfigInput) (model.GetConfigOutput, error)
		UpdateConfig(ctx context.Context, input model.UpdateConfigInput) error
	}
)

var (
	localConfigManagement IConfigManagement
)

func ConfigManagement() IConfigManagement {
	if localConfigManagement == nil {
		panic("implement not found for interface IConfigManagement, forgot register?")
	}
	return localConfigManagement
}

func RegisterConfigManagement(i IConfigManagement) {
	localConfigManagement = i
}
