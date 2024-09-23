package upstream

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/components/balancer"
	"api-gateway/internal/components/config"
	"api-gateway/internal/consts"
	"api-gateway/internal/model"
)

const (
	defaultWeight = 1
)

type (
	// Service contains Upstream list
	Service struct {
		mu         sync.RWMutex
		ups        []*Upstream
		available  int
		Config     model.ServiceConfig
		RoutingKey string
	}
)

func NewService(routingKey string, cfg model.ServiceConfig) *Service {
	s := &Service{
		mu:         sync.RWMutex{},
		ups:        make([]*Upstream, 0),
		Config:     cfg,
		RoutingKey: routingKey,
	}
	config.RegisterConfigChangeEventHandler(routingKey, s.configEventHandler)
	return s
}

func GetService(routingKey string) (*Service, bool) {
	return cache.getService(routingKey)
}

// Set upstream by upsert
func (s *Service) Set(u *Upstream) {
	s.mu.Lock()
	defer s.mu.Unlock()

	u.Parent = s
	for i, upstream := range s.ups {
		if upstream.Identity() == u.Identity() {
			s.ups[i] = u
			return
		}
	}
	s.ups = append(s.ups, u)
}

func (s *Service) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, upstream := range s.ups {
		if upstream.Identity() == id {
			s.ups = append(s.ups[:i], s.ups[i+1:]...)
			return
		}
	}
}

// SelectOne selects an Upstream from Service
func (s *Service) SelectOne(r *ghttp.Request, balancer balancer.Balancer, filters balancer.Filters) (u *Upstream, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var available []any
	for _, uu := range s.ups {
		if uu.healthy() {
			available = append(available, uu)
		}
	}

	// using client ip as hash key, only for balancer.StrategyHash
	v, err := balancer.Pick(available, r.GetClientIp(), filters)
	if err != nil {
		g.Log().Warningf(r.Context(), "select upstream failed: %v", err)
		return
	}
	u = v.(*Upstream)
	return
}

func (s *Service) CountUpstream() int {
	return len(s.ups)
}

func (s *Service) CountAvailableUpstream() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.available == 0 && len(s.ups) > 0 {
		cnt := 0
		for _, u := range s.ups {
			if u.healthy() {
				cnt++
			}
		}
		s.available = cnt
	}
	return s.available
}

func (s *Service) configEventHandler(t config.EventType, module, key string, value interface{}) {
	g.Log().Infof(context.Background(), "config event: type=%s, key=%s, value=%v", t, key, value)
	if module == consts.ModuleNameLoadBalance {
		// update lb instance
		balancer.Update(s.RoutingKey)
		return
	}
	// maybe has some way better than update all at once
	for _, up := range s.ups {
		// async update
		go up.updateConfig(module)
	}
}
