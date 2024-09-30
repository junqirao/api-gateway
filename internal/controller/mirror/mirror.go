package mirror

import (
	"context"

	"github.com/gogf/gf/v2/util/gconv"
	registry "github.com/junqirao/simple-registry"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"

	v1 "api-gateway/api/mirror/v1"
	"api-gateway/internal/components/authentication"
	"api-gateway/internal/components/mirror"
	"api-gateway/internal/components/response"
)

type Controller struct {
	v1.UnimplementedMirrorServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterMirrorServer(s.Server, &Controller{})
}

func (c *Controller) Register(ctx context.Context, in *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	addr := in.GetInstance().GetHost()
	if !mirror.Allow(addr) {
		err = response.CodePermissionDeny.WithDetail(addr)
		return
	}
	ins := in.GetInstance()
	if !authentication.L.Compare(ins.GetId(), in.GetAuthentication()) {
		err = response.CodePermissionDeny
		return
	}

	ttl, err := mirror.Register(ctx, mirror.ClientInfo{
		Instance: c.convert2ins(ins),
		Filter:   in.GetFilter(),
	})
	if err != nil {
		return
	}
	res = &v1.RegisterRes{Ttl: int32(ttl), Success: true}
	return
}

func (c *Controller) UnRegister(ctx context.Context, in *v1.UnRegisterReq) (res *v1.UnRegisterRes, err error) {
	addr := in.GetInstance().GetHost()
	if !mirror.Allow(addr) {
		err = response.CodePermissionDeny.WithDetail(addr)
		return
	}
	if !authentication.L.Compare(in.GetInstance().GetId(), in.GetAuthentication()) {
		err = response.CodePermissionDeny
		return
	}
	err = mirror.UnRegister(ctx, c.convert2ins(in.GetInstance()))
	if err != nil {
		return
	}

	res = &v1.UnRegisterRes{Success: true}
	return
}

func (c *Controller) convert2ins(ins *v1.Instance) *registry.Instance {
	res := &registry.Instance{
		Id:          ins.GetId(),
		Host:        ins.GetHost(),
		HostName:    ins.GetHostName(),
		Port:        int(ins.GetPort()),
		ServiceName: ins.GetServiceName(),
	}
	if ins.GetMeta() != "" {
		res.Meta = gconv.Map(ins.GetMeta())
	}
	return res
}
