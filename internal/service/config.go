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
	IConfig interface {
		GetConfig(ctx context.Context, input model.GetConfigInput) (model.GetConfigOutput, error)
		UpdateConfig(ctx context.Context, input model.UpdateConfigInput) error
	}
)

var (
	localConfig IConfig
)

func Config() IConfig {
	if localConfig == nil {
		panic("implement not found for interface IConfig, forgot register?")
	}
	return localConfig
}

func RegisterConfig(i IConfig) {
	localConfig = i
}
