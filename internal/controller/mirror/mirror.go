package mirror

import (
	"context"

	"github.com/gogf/gf/v2/util/gconv"
	registry "github.com/junqirao/simple-registry"

	v1 "api-gateway/api/mirror/v1"
	"api-gateway/internal/components/mirror"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
)

type Controller struct {
	v1.UnimplementedMirrorServer
}

func Register(s *grpcx.GrpcServer) {
	v1.RegisterMirrorServer(s.Server, &Controller{})
}

func (c *Controller) Register(ctx context.Context, in *v1.RegisterReq) (res *v1.RegisterRes, err error) {
	ins := in.GetInstance()
	ci := mirror.ClientInfo{
		Instance: c.convert2ins(ins),
		Filter:   in.GetFilter(),
	}
	ttl, err := mirror.Register(ctx, ci)
	if err != nil {
		return
	}
	res = &v1.RegisterRes{Ttl: int32(ttl), Success: true}
	return
}

func (c *Controller) UnRegister(ctx context.Context, in *v1.UnRegisterReq) (res *v1.UnRegisterRes, err error) {
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
