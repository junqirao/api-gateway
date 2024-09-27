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
	IUpstreamManagement interface {
		GetServiceNames(_ context.Context) []string
		GetServiceState(ctx context.Context, serviceName string) (output *model.GetServiceStateOutput, err error)
	}
)

var (
	localUpstreamManagement IUpstreamManagement
)

func UpstreamManagement() IUpstreamManagement {
	if localUpstreamManagement == nil {
		panic("implement not found for interface IUpstreamManagement, forgot register?")
	}
	return localUpstreamManagement
}

func RegisterUpstreamManagement(i IUpstreamManagement) {
	localUpstreamManagement = i
}
