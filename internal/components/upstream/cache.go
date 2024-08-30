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
		ignoreServiceName string
		m                 sync.Map // routing_key(service_name) : *model.Service
	}
)

func newUpstreamCache(ctx context.Context) *cacheHandler {
	h := &cacheHandler{
		// ignore self
		ignoreServiceName: g.Cfg().MustGet(ctx, "registry.service_name").String(),
		m:                 sync.Map{},
	}
	h.build(ctx)
	h.registerEvent()
	return h
}

func (h *cacheHandler) registerEvent() {
	registry.Registry.RegisterEventHandler(func(instance *registry.Instance, e registry.EventType) {
		ctx := context.Background()
		if instance == nil {
			g.Log().Infof(ctx, "instance not found, skip event.")
			return
		}

		if instance.ServiceName == h.ignoreServiceName {
			return
		}

		switch e {
		case registry.EventTypeUpdate, registry.EventTypeCreate:
			g.Log().Infof(ctx, "service[%s] %s upstreams , instance=%v", instance.ServiceName, e, instance.String())
			h.getOrCreateService(ctx, instance.ServiceName).Set(NewUpstream(ctx, instance))
		case registry.EventTypeDelete:
			g.Log().Infof(ctx, "service[%s] delete upstream instance=%s", instance.ServiceName, instance.Identity())
			srv, ok := h.getService(instance.ServiceName)
			if !ok {
				g.Log().Warningf(ctx, "upstream cache not found service=%s", instance.ServiceName)
				return
			}
			srv.Delete(instance.Identity())
			if len(srv.Ups) == 0 {
				h.m.Delete(instance.ServiceName)
			}
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
		if sName == h.ignoreServiceName {
			continue
		}
		srv := h.getOrCreateService(ctx, sName)
		instances.Range(func(instance *registry.Instance) bool {
			id := instance.Identity()
			delete(current, id)
			// upsert
			srv.Set(NewUpstream(ctx, instance))
			g.Log().Infof(ctx, "build service[%s] upstream set instance=%v, ups_length=%d", sName, instance.String(), len(srv.Ups))
			return true
		})
	}
	for k, _ := range current {
		g.Log().Infof(ctx, "remove service=%s", k)
		h.m.Delete(k)
	}
	g.Log().Infof(ctx, "upstream cache build done.")
}

func (h *cacheHandler) getService(routingKey string) (srv *Service, ok bool) {
	var v interface{}
	if v, ok = h.m.Load(routingKey); ok {
		srv = v.(*Service)
	}
	return
}

func (h *cacheHandler) getOrCreateService(ctx context.Context, routingKey string) (srv *Service) {
	var ok bool
	srv, ok = h.getService(routingKey)
	if !ok {
		cfg, ok := config.GetServiceConfig(routingKey)
		g.Log().Infof(ctx, "upstream cache create service=%s, defaultServiceConfig=%v", routingKey, !ok)
		srv = NewService(routingKey, *cfg)
		h.setService(routingKey, srv)
	}
	return
}

func (h *cacheHandler) setService(routingKey string, srv *Service) {
	h.m.Store(routingKey, srv)
}
