package upstream

import (
	"sync"

	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/components/loadbalance"
	"api-gateway/internal/model"
)

const (
	basicLoadBalanceWeight = 1
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
	// Selector selects an Upstream from Service
	Selector func(r *ghttp.Request, ups []loadbalance.Weighted) (ref int, ok bool)
)

func NewService(routingKey string, cfg model.ServiceConfig) *Service {
	return &Service{
		mu:         sync.RWMutex{},
		ups:        make([]*Upstream, 0),
		Config:     cfg,
		RoutingKey: routingKey,
	}
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
			u.SetRef(i)
			return
		}
	}
	s.ups = append(s.ups, u)
	u.SetRef(len(s.ups) - 1)
}

func (s *Service) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, upstream := range s.ups {
		if upstream.Identity() == id {
			s.ups = append(s.ups[:i], s.ups[i+1:]...)
			// reset ref from i+1 -> len(s.ups
			for j := i; j < len(s.ups); j++ {
				s.ups[j].SetRef(j)
			}
			return
		}
	}
}

func (s *Service) Select(r *ghttp.Request, selector Selector) (u *Upstream, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	available := s.availableWeighted()
	if s == nil || len(available) == 0 {
		return
	}

	if len(available) == 1 {
		ok = true
		u = s.ups[available[0].Ref()]
		return
	}

	idx, ok := selector(r, available)
	if ok {
		u = s.ups[idx]
	}
	return
}

func (s *Service) availableWeighted() []loadbalance.Weighted {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var ups []loadbalance.Weighted
	for _, u := range s.ups {
		if u.healthy() {
			ups = append(ups, u)
		}
	}
	s.available = len(ups)
	return ups
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
