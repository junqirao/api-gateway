package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/util/gconv"
	registry "github.com/junqirao/simple-registry"
	"google.golang.org/grpc"
)

func ClientConnFromInstance(ins *registry.Instance, opt ...grpc.DialOption) (cc *grpc.ClientConn, err error) {
	if ins == nil {
		err = errors.New("empty instance")
		return
	}
	if cfg, ok := ins.Meta["grpc"]; ok {
		target := ins.Host
		part := strings.Split(gconv.String(gconv.Map(cfg)["address"]), ":")
		if len(part) > 1 {
			target = fmt.Sprintf("%s:%s", ins.Host, part[1])
		}

		cc, err = grpcx.Client.NewGrpcClientConn(target, opt...)
	} else {
		err = fmt.Errorf("no grpc config at instance: %s", ins.Identity())
	}
	return
}
