package inner

import (
	"context"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"

	"api-gateway/api/inner/upstream"
	"api-gateway/internal/components/authentication"
	"api-gateway/internal/components/response"
	ups "api-gateway/internal/components/upstream"
)

type Controller struct {
	upstream.UnimplementedManagementServer
}

func Register(s *grpcx.GrpcServer) {
	upstream.RegisterManagementServer(s.Server, &Controller{})
}

func (*Controller) GetServiceStates(_ context.Context, req *upstream.GetServiceStatesReq) (res *upstream.GetServiceStatesResp, err error) {
	if !authentication.L.Compare(req.GetAuthentication(), req.GetInstanceId()) {
		err = response.CodePermissionDeny
		return
	}
	srv, ok := ups.Cache.GetService(req.ServiceName)
	if !ok {
		return
	}
	var sts []*ups.UpsState
	for _, up := range srv.Upstreams() {
		sts = append(sts, up.State())
	}
	res = &upstream.GetServiceStatesResp{}
	for _, state := range sts {
		res.States = append(res.States, &upstream.State{
			Hostname:     state.HostName,
			InstanceId:   state.InstanceId,
			Healthy:      state.Healthy,
			Weight:       state.Weight,
			WeightRatio:  float32(state.WeightRatio),
			Load:         state.Load,
			BreakerState: state.BreakerState,
		})
	}
	return
}
