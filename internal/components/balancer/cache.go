package balancer

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"

	"api-gateway/internal/components/config"
)

var (
	m = sync.Map{} // service_name(routing_key):Balancer
)

// GetOrCreate get or create load balancer
func GetOrCreate(routingKey string) Balancer {
	b, ok := m.Load(routingKey)
	if !ok || b == nil {
		return Update(routingKey)
	}
	return b.(Balancer)
}

// Update load balancer instance with latest config
func Update(routingKey string) Balancer {
	cfg, _ := config.GetServiceConfig(routingKey)
	b := New(Strategy(cfg.LoadBalance.Strategy))
	m.Store(routingKey, b)
	g.Log().Infof(context.Background(), "service %s load-balancer updated: %s", routingKey, cfg.LoadBalance.Strategy)
	return b
}
