package program

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/config"
	"api-gateway/internal/consts"
)

func Init(ctx context.Context) {
	// Variables
	Variables = &variableHandler{global: make(map[string]interface{})}
	Variables.build(ctx)
	registry.Storages.SetEventHandler(consts.StorageNameVariable, Variables.eventHandler)
	// program
	buildCache(ctx)
	registry.Storages.SetEventHandler(consts.StorageNameProgram, func(t registry.EventType, key string, value interface{}) {
		g.Log().Infof(ctx, "program change event: type=%s key=%s", t, key)
		parts := strings.Split(key, config.StorageSeparator)
		if len(parts) < 2 {
			// drop invalid key
			return
		}
		serviceName := parts[0]
		programName := parts[1]
		switch t {
		case registry.EventTypeUpdate, registry.EventTypeCreate:
			if _, err := buildCacheByService(serviceName); err != nil {
				g.Log().Errorf(ctx, "build cache failed: %v", err)
				return
			}
		case registry.EventTypeDelete:
			if v, ok := m.Load(serviceName); ok && v != empty {
				p := v.(*Programs)
				p.Delete(programName)
				if len(p.ps) == 0 {
					m.Delete(serviceName)
				}
			}
		}
	})
}
