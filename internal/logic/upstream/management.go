package upstream

import (
	"context"
	"sort"

	"api-gateway/internal/components/upstream"
	"api-gateway/internal/service"
)

type (
	sUpstreamManagement struct{}
)

func init() {
	service.RegisterUpstreamManagement(&sUpstreamManagement{})
}

func (s sUpstreamManagement) GetServiceNames(_ context.Context) []string {
	var ss []string
	upstream.Cache.Range(func(key, _ any) bool {
		ss = append(ss, key.(string))
		return true
	})
	sort.Strings(ss)
	return ss
}

func (s sUpstreamManagement) GetServiceState(_ context.Context, serviceName string) []*upstream.UpsState {
	srv, ok := upstream.Cache.GetService(serviceName)
	if !ok {
		return nil
	}
	var sts []*upstream.UpsState
	for _, up := range srv.Upstreams() {
		sts = append(sts, up.State())
	}
	return sts
}
