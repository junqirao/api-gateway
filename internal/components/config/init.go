package config

import (
	"context"
	"strings"
	"sync"

	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/consts"
)

var (
	evs sync.Map // service_name(routing_key):ChangeEvent
)

type (
	ChangeEvent func(t EventType, module, key string, value interface{})
	// EventType alias of registry.EventType
	EventType = registry.EventType
)

func Init(ctx context.Context) {
	// load configs from file
	loadConfigs(ctx)
	// init config update event bus
	initConfigUpdateEventBus()
}

func initConfigUpdateEventBus() {
	registry.Storages.SetEventHandler(consts.StorageNameServiceConfig, func(t registry.EventType, key string, value interface{}) {
		parts := strings.Split(key, StorageSeparator)
		if len(parts) < 2 {
			// drop invalid key
			return
		}

		serviceName := parts[0]
		if handler, ok := evs.Load(serviceName); ok {
			handler.(ChangeEvent)(t, parts[1], key, value)
		}
	})
}

func RegisterConfigChangeEventHandler(serviceName string, handler ChangeEvent) {
	evs.Store(serviceName, handler)
}
