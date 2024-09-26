package upstream

import (
	"context"

	"api-gateway/api/upstream/v1"
	"api-gateway/internal/service"
)

func (c *ControllerV1) GetServiceNames(ctx context.Context, _ *v1.GetServiceNamesReq) (res *v1.GetServiceNamesRes, err error) {
	ss := v1.GetServiceNamesRes(service.UpstreamManagement().GetServiceNames(ctx))
	res = &ss
	return
}
func (c *ControllerV1) GetServiceState(ctx context.Context, req *v1.GetServiceStateReq) (res *v1.GetServiceStateRes, err error) {
	res = new(v1.GetServiceStateRes)
	res.UpstreamStates = service.UpstreamManagement().GetServiceState(ctx, req.ServiceName)
	res.Upstreams = len(res.UpstreamStates)
	return
}
