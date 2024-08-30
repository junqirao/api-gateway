package loadbalance

import (
	"fmt"
	"sync"

	"api-gateway/internal/components/config"
)

var (
	m = sync.Map{} // service_name(routing_key):Balancer
)

func GetOrCreate(routingKey string) Balancer {
	var (
		cfg, _ = config.GetServiceConfig(routingKey)
	)

	b, ok := m.Load(key(routingKey, cfg.LoadBalance.Strategy))
	if !ok || b == nil {
		b = New(cfg.LoadBalance.Strategy)
		m.Store(key(routingKey, cfg.LoadBalance.Strategy), b)
	}
	return b.(Balancer)
}

func key(routingKey string, strategy string) string {
	return fmt.Sprintf("%s_%s", routingKey, strategy)
}
