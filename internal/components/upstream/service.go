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
		Ups        []*Upstream
		Config     model.ServiceConfig
		RoutingKey string
	}
	// Selector selects an Upstream from Service
	Selector func(r *ghttp.Request, ups []loadbalance.Weighted) (ref int, ok bool)
)

func NewService(routingKey string, cfg model.ServiceConfig) *Service {
	return &Service{
		mu:         sync.RWMutex{},
		Ups:        make([]*Upstream, 0),
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

	for i, upstream := range s.Ups {
		if upstream.Identity() == u.Identity() {
			s.Ups[i] = u
			u.SetRef(i)
			return
		}
	}
	s.Ups = append(s.Ups, u)
	u.SetRef(len(s.Ups) - 1)
}

func (s *Service) Delete(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, upstream := range s.Ups {
		if upstream.Identity() == id {
			s.Ups = append(s.Ups[:i], s.Ups[i+1:]...)
			// reset ref from i+1 -> len(s.Ups
			for j := i; j < len(s.Ups); j++ {
				s.Ups[j].SetRef(j)
			}
			return
		}
	}
}

func (s *Service) Select(r *ghttp.Request, selector Selector) (u *Upstream, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s == nil || len(s.Ups) == 0 {
		return
	}
	if len(s.Ups) == 1 {
		ok = true
		u = s.Ups[0]
		return
	}

	idx, ok := selector(r, s.availableWeighted())
	if ok {
		u = s.Ups[idx]
	}
	return
}

func (s *Service) availableWeighted() []loadbalance.Weighted {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var ups []loadbalance.Weighted
	for _, u := range s.Ups {
		if u.healthy() {
			ups = append(ups, u)
		}
	}
	return ups
}
