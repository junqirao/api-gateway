package program

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/program"
	"api-gateway/internal/consts"
	"api-gateway/internal/model"
	"api-gateway/internal/service"
)

type (
	sProgramManagement struct{}
)

func init() {
	service.RegisterProgramManagement(&sProgramManagement{})
}

func (s sProgramManagement) GetProgramInfo(ctx context.Context, serviceName string) (res map[string][]*program.Info, err error) {
	kvs, err := registry.Storages.GetStorage(consts.StorageNameProgram).Get(ctx, serviceName)
	if err != nil {
		return
	}

	res = make(map[string][]*program.Info)
	for _, kv := range kvs {
		info := new(program.Info)
		if err := kv.Value.Scan(&info); err != nil {
			g.Log().Warningf(ctx, "scan program info failed: %v", err)
			continue
		}
		res[info.ServiceName] = append(res[info.ServiceName], info)
	}
	return
}

func (s sProgramManagement) SetProgramInfo(ctx context.Context, info *program.Info) (err error) {
	// make sure no syntax error
	if err = info.TryCompile(ctx); err != nil {
		return
	}
	sto := registry.Storages.GetStorage(consts.StorageNameProgram)
	err = sto.Set(ctx, fmt.Sprintf("%s%s%s", info.ServiceName, config.StorageSeparator, info.Name), info)
	return
}

func (s sProgramManagement) DeleteProgramInfo(ctx context.Context, input model.DeleteProgramInfoInput) (err error) {
	sto := registry.Storages.GetStorage(consts.StorageNameProgram)
	err = sto.Delete(ctx, fmt.Sprintf("%s%s%s", input.ServiceName, config.StorageSeparator, input.Name))
	return
}

// GetGlobalVariables get global variable
func (s sProgramManagement) GetGlobalVariables(ctx context.Context) map[string]interface{} {
	return program.Variables.GetGlobalVariables(ctx)
}

// SetGlobalVariables set global variable
func (s sProgramManagement) SetGlobalVariables(ctx context.Context, input model.SetGlobalVariablesInput) error {
	return program.Variables.SetGlobalVariable(ctx, input.Key, input.Value)
}

// DeleteGlobalVariables delete global variable
func (s sProgramManagement) DeleteGlobalVariables(ctx context.Context, key string) error {
	return program.Variables.DeleteGlobalVariable(ctx, key)
}
