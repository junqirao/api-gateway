package program

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/response"
	"api-gateway/internal/consts"
)

var (
	management = &managementHandler{}
)

type (
	managementHandler struct {
	}

	// ManagementGetProgramRequest get program
	ManagementGetProgramRequest struct {
		ServiceName string `json:"service_name"`
	}
	// ManagementDeleteProgramRequest delete program
	ManagementDeleteProgramRequest struct {
		ServiceName string `json:"service_name"`
		Name        string `json:"name"`
	}
	// ManagementSetGlobalVariablesRequest set global
	ManagementSetGlobalVariablesRequest struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	// ManagementDeleteGlobalVariablesRequest delete global
	ManagementDeleteGlobalVariablesRequest struct {
		Key string `json:"key"`
	}
)

// GetProgram get
func (h *managementHandler) GetProgram(r *ghttp.Request) {
	req := new(ManagementGetProgramRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err))
		return
	}

	kvs, err := registry.Storages.GetStorage(consts.StorageNameProgram).Get(r.Context(), req.ServiceName)
	if err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}

	res := make(map[string][]*Info)
	for _, kv := range kvs {
		info := new(Info)
		if err = kv.Value.Scan(&info); err != nil {
			g.Log().Warningf(r.Context(), "scan program info failed: %v", err)
			continue
		}
		res[info.ServiceName] = append(res[info.ServiceName], info)
	}
	response.WriteData(r, response.CodeDefaultSuccess, res)
}

// SetProgram set
func (h *managementHandler) SetProgram(r *ghttp.Request) {
	req := new(Info)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}
	sto := registry.Storages.GetStorage(consts.StorageNameProgram)
	err := sto.Set(r.Context(), fmt.Sprintf("%s%s%s", req.ServiceName, config.StorageSeparator, req.Name), req)
	if err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}

	response.WriteJSON(r, response.CodeDefaultSuccess)
}

// DeleteProgram delete
func (h *managementHandler) DeleteProgram(r *ghttp.Request) {
	req := new(ManagementDeleteProgramRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	sto := registry.Storages.GetStorage(consts.StorageNameProgram)
	err := sto.Delete(r.Context(), fmt.Sprintf("%s%s%s", req.ServiceName, config.StorageSeparator, req.Name))
	if err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}
	response.WriteJSON(r, response.CodeDefaultSuccess)
}

// GetGlobalVariables get global variable
func (h *managementHandler) GetGlobalVariables(r *ghttp.Request) {
	response.WriteData(r, response.CodeDefaultSuccess, Variables.GetGlobalVariables(r.Context()))
}

// SetGlobalVariables set global variable
func (h *managementHandler) SetGlobalVariables(r *ghttp.Request) {
	req := new(ManagementSetGlobalVariablesRequest)
	if err := r.Parse(&req); err != nil {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err.Error()))
		return
	}

	if err := Variables.SetGlobalVariable(r.Context(), req.Key, req.Value); err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}
	response.WriteJSON(r, response.CodeDefaultSuccess)
}

// DeleteGlobalVariables delete global variable
func (h *managementHandler) DeleteGlobalVariables(r *ghttp.Request) {
	req := new(ManagementDeleteGlobalVariablesRequest)
	if err := r.Parse(&req); err != nil || req.Key == "" {
		response.WriteJSON(r, response.CodeInvalidParameter.WithDetail(err))
		return
	}
	if err := Variables.DeleteGlobalVariable(r.Context(), req.Key); err != nil {
		response.WriteJSON(r, response.CodeDefaultFailure.WithDetail(err.Error()))
		return
	}
	response.WriteJSON(r, response.CodeDefaultSuccess)
}
