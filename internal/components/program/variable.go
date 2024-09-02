package program

import (
	"context"
	"errors"
	"maps"
	"strings"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/config"
)

var (
	storageNameProgram  = "program"
	storageNameVariable = "program_variable"
	variables           *variableHandler
)

type (
	variableHandler struct {
		load   sync.Once
		mu     sync.RWMutex
		global map[string]interface{}
	}
)

func (h *variableHandler) GetGlobalVariables(_ context.Context) map[string]interface{} {
	h.mu.Lock()
	defer h.mu.Unlock()

	return maps.Clone(h.global)
}

func (h *variableHandler) SetGlobalVariable(ctx context.Context, key string, value interface{}) (err error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	err = registry.Storages.GetStorage(storageNameVariable).Set(ctx, key, value)
	if err == nil {
		h.global[key] = gconv.String(value)
	}
	return
}

func (h *variableHandler) build(ctx context.Context) {
	h.mu.Lock()
	defer h.mu.Unlock()

	kvs, err := registry.Storages.GetStorage(storageNameVariable).Get(ctx)
	switch {
	case err == nil:
	case errors.Is(err, registry.ErrStorageNotFound):
	default:
		g.Log().Warningf(ctx, "build global variables failed: %v", err)
		return
	}

	for _, k := range kvs {
		key := k.Key
		parts := strings.Split(k.Key, config.StorageSeparator)
		if len(parts) > 0 {
			key = parts[len(parts)-1]
		}
		h.global[key] = k.Value.String()
	}
}

func (h *variableHandler) eventHandler(t registry.EventType, key string, value interface{}) {
	g.Log().Infof(context.Background(), "global variable change event: type=%s key=%s", t, key)
	h.mu.Lock()
	defer h.mu.Unlock()
	switch t {
	case registry.EventTypeDelete:
		delete(h.global, key)
	default:
		h.global[key] = gconv.String(value)
	}
}
