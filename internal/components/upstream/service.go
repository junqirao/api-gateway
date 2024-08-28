package upstream

import (
	"math/rand"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/model"
)

func NewService(routingKey string, cfg model.ServiceConfig) *Service {
	return &Service{
		Ups:        make([]*Upstream, 0),
		Config:     cfg,
		RoutingKey: routingKey,
	}
}

func (s *Service) Set(u *Upstream) {
	for i, upstream := range s.Ups {
		if upstream.Identity() == u.Identity() {
			s.Ups[i] = u
			return
		}
	}
	s.Ups = append(s.Ups, u)
}

func (s *Service) Delete(id string) {
	for i, upstream := range s.Ups {
		if upstream.Identity() == id {
			s.Ups = append(s.Ups[:i], s.Ups[i+1:]...)
			return
		}
	}
}

func (s *Service) Select(r *ghttp.Request, selector ...Selector) (u *Upstream, ok bool) {
	if s == nil || len(s.Ups) == 0 {
		return
	}
	if len(s.Ups) == 1 {
		ok = true
		u = s.Ups[0]
		return
	}

	var se Selector
	if len(selector) > 0 && selector[0] != nil {
		se = selector[0]
	}
	if se == nil {
		se = func(_ *ghttp.Request, ss *Service) (*Upstream, bool) {
			return ss.Ups[rand.New(rand.NewSource(time.Now().UnixNano())).Intn(len(ss.Ups))], true
		}
	}

	return se(r, s)
}
