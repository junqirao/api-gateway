package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"api-gateway/internal/components/program"
)

// GetProgramInfoReq get program
type GetProgramInfoReq struct {
	g.Meta `path:"/program/info" method:"get" tags:"Program" summary:"Get Program"`

	ServiceName string `json:"service_name"`
}
type GetProgramInfoRes map[string][]*program.Info

// SetProgramInfoReq set program
type SetProgramInfoReq struct {
	g.Meta `path:"/program/info" method:"put" tags:"Program" summary:"Set Program"`

	program.Info
}
type SetProgramInfoRes struct {
}

// DeleteProgramInfoReq delete program
type DeleteProgramInfoReq struct {
	g.Meta      `path:"/program/info" method:"delete" tags:"Program" summary:"Delete Program"`
	ServiceName string `json:"service_name"`
	Name        string `json:"name"`
}
type DeleteProgramInfoRes struct {
}

// GetGlobalVariablesReq set global
type GetGlobalVariablesReq struct {
	g.Meta `path:"/program/variable" method:"get" tags:"Program" summary:"Get Global Variables"`
}
type GetGlobalVariablesRes map[string]interface{}

// SetGlobalVariablesReq set global
type SetGlobalVariablesReq struct {
	g.Meta `path:"/program/variable" method:"put" tags:"Program" summary:"Set Global Variables"`

	Key   string `json:"key"`
	Value string `json:"value"`
}
type SetGlobalVariablesRes struct {
}

// DeleteGlobalVariablesReq delete global
type DeleteGlobalVariablesReq struct {
	g.Meta `path:"/program/variable" method:"delete" tags:"Program" summary:"Delete Global Variables"`

	Key string `json:"key"`
}
type DeleteGlobalVariablesRes struct {
}
