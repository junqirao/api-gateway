package upstream

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/config"
)

var (
	cache *cacheHandler
)

type (
	cacheHandler struct {
		mu                sync.RWMutex
		ignoreServiceName string
		m                 sync.Map // routing_key(service_name) : *model.Service
	}
)

func newUpstreamCache(ctx context.Context) *cacheHandler {
	h := &cacheHandler{
		// ignore self
		ignoreServiceName: g.Cfg().MustGet(ctx, "registry.service_name").String(),
		m:                 sync.Map{},
		mu:                sync.RWMutex{},
	}
	h.registerEvent()
	h.build(ctx)
	return h
}

func (h *cacheHandler) registerEvent() {
	registry.Registry.RegisterEventHandler(func(ctx context.Context, instance *registry.Instance, e registry.EventType) {
		h.mu.Lock()
		defer h.mu.Unlock()

		if instance == nil {
			g.Log().Infof(ctx, "instance not found, skip event.")
			return
		}

		if instance.ServiceName == h.ignoreServiceName {
			return
		}

		ups, _ := h.GetService(instance.ServiceName)
		defer func() {
			h.setService(instance.ServiceName, ups)
		}()

		switch e {
		case registry.EventTypeUpdate, registry.EventTypeCreate:
			g.Log().Infof(ctx, "upstream cache %s service=%s instance=%v", e, instance.ServiceName, instance.String())
			ups.Set(NewUpstream(ctx, instance))
		case registry.EventTypeDelete:
			g.Log().Infof(ctx, "upstream cache delete instance=%s", instance.Identity())
			ups.Delete(instance.Identity())
		}
	})
}

func (h *cacheHandler) build(ctx context.Context) {
	services, err := registry.Registry.GetServices(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "upstream cache failed to get services: %v", err)
		return
	}
	current := make(map[string]struct{})
	h.m.Range(func(key, value interface{}) bool {
		current[key.(string)] = struct{}{}
		return true
	})
	for sName, instances := range services {
		ups, _ := h.GetService(sName)
		instances.Range(func(instance *registry.Instance) bool {
			if instance.ServiceName == h.ignoreServiceName {
				return true
			}
			id := instance.Identity()
			delete(current, id)
			g.Log().Infof(ctx, "upstream cache add service=%s instance=%v", sName, instance.String())
			ups.Set(NewUpstream(ctx, instance))
			return true
		})
		h.setService(sName, ups)
	}
	for k, _ := range current {
		g.Log().Infof(ctx, "upstream cache delete instance=%s", k)
		h.m.Delete(k)
	}
	g.Log().Infof(ctx, "upstream cache build done.")
}

func (h *cacheHandler) GetService(routingKey string) (ups *Service, ok bool) {
	v, ok := h.m.Load(routingKey)
	if !ok {
		cfg, has := config.GetServiceConfig(routingKey)
		ups = NewService(routingKey, *cfg)
		if has {
			h.m.Store(routingKey, ups)
		}
	} else {
		ups = v.(*Service)
	}
	return
}

func (h *cacheHandler) setService(routingKey string, ups *Service) {
	h.m.Store(routingKey, ups)
}
