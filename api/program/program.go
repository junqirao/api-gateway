// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package program

import (
	"context"

	"api-gateway/api/program/v1"
)

type IProgramV1 interface {
	GetProgramInfo(ctx context.Context, req *v1.GetProgramInfoReq) (res *v1.GetProgramInfoRes, err error)
	SetProgramInfo(ctx context.Context, req *v1.SetProgramInfoReq) (res *v1.SetProgramInfoRes, err error)
	DeleteProgramInfo(ctx context.Context, req *v1.DeleteProgramInfoReq) (res *v1.DeleteProgramInfoRes, err error)
	GetGlobalVariables(ctx context.Context, req *v1.GetGlobalVariablesReq) (res *v1.GetGlobalVariablesRes, err error)
	SetGlobalVariables(ctx context.Context, req *v1.SetGlobalVariablesReq) (res *v1.SetGlobalVariablesRes, err error)
	DeleteGlobalVariables(ctx context.Context, req *v1.DeleteGlobalVariablesReq) (res *v1.DeleteGlobalVariablesRes, err error)
}
