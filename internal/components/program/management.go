package program

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/response"
	"api-gateway/internal/model"
)

var (
	Management = &ManagementHandler{}
)

type (
	ManagementHandler struct {
	}

	ManagementGetProgramRequest struct {
		ServiceName string `json:"service_name"`
	}
	ManagementDeleteProgramRequest struct {
		ServiceName string `json:"service_name"`
		Name        string `json:"name"`
	}
	ManagementSetGlobalVariablesRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
)

func (h *ManagementHandler) GetProgram(r *ghttp.Request) {
	req := new(ManagementGetProgramRequest)
	if err := r.Parse(&req); err != nil || req.ServiceName == "" {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err))
		return
	}

	kvs, err := registry.Storages.GetStorage(storageNameProgram).Get(r.Context())
	if err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}

	infos := make([]*model.ProgramInfo, 0)
	for _, kv := range kvs {
		info := new(model.ProgramInfo)
		if err = kv.Value.Scan(&info); err != nil {
			g.Log().Warningf(r.Context(), "scan program info failed: %v", err)
			continue
		}
		infos = append(infos, info)
	}
	response.WriteData(r, response.CodeDefaultSuccess, infos)
}

func (h *ManagementHandler) SetProgram(r *ghttp.Request) {
	req := new(model.ProgramInfo)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}
	sto := registry.Storages.GetStorage(storageNameProgram)
	err := sto.Set(r.Context(), fmt.Sprintf("%s%s%s", req.ServiceName, config.StorageSeparator, req.Name), req)
	if err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}

	response.WriteJSON(r, response.CodeDefaultSuccess)
}

func (h *ManagementHandler) DeleteProgram(r *ghttp.Request) {
	req := new(ManagementDeleteProgramRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	sto := registry.Storages.GetStorage(storageNameProgram)
	err := sto.Delete(r.Context(), fmt.Sprintf("%s%s%s", req.ServiceName, config.StorageSeparator, req.Name))
	if err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}
	response.WriteJSON(r, response.CodeDefaultSuccess)
}

func (h *ManagementHandler) GetGlobalVariables(r *ghttp.Request) {
	response.WriteData(r, response.CodeDefaultSuccess, variables.GetGlobalVariables(r.Context()))
}

func (h *ManagementHandler) SetGlobalVariables(r *ghttp.Request) {
	req := new(ManagementSetGlobalVariablesRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	if err := variables.SetGlobalVariable(r.Context(), req.Key, req.Value); err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}
	response.WriteJSON(r, response.CodeDefaultSuccess)
}
