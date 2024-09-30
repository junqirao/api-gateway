package registry

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/grace"
)

var (
	CurrentInstance *registry.Instance
)

func Init(ctx context.Context) {
	var (
		cfg = registry.Config{}
		ins = &registry.Instance{}
		v   = g.Cfg().MustGet(ctx, "registry")
	)

	// parse config
	if err := v.Scan(&cfg); err != nil {
		g.Log().Fatal(ctx, err)
	}
	if err := v.Scan(&ins); err != nil {
		g.Log().Fatal(ctx, err)
	}

	// overwrite port and address if exists
	if addr := g.Cfg().MustGet(ctx, "server.address").String(); addr != "" {
		parts := strings.Split(addr, ":")
		if len(parts) == 2 {
			if parts[0] != "" {
				ins.Host = parts[0]
			}
			if parts[1] != "" {
				ins.Port = gconv.Int(parts[1])
			}
		}
	}

	// inject grpc config
	ins.WithMetaData(map[string]interface{}{
		"grpc": g.Cfg().MustGet(ctx, "grpc").Map(),
	})
	// init registry and register
	if err := registry.Init(ctx, cfg, ins); err != nil {
		g.Log().Fatal(ctx, err)
	}
	CurrentInstance = ins
	grace.Register(ctx, "deregister_registry", func() {
		_ = registry.Registry.Deregister(ctx)
	})
}
