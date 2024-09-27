package registry

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/grace"
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

	// inject grpc config
	ins.WithMetaData(map[string]interface{}{
		"grpc": g.Cfg().MustGet(ctx, "grpc").Map(),
	})
	// init registry and register
	if err := registry.Init(ctx, cfg, ins); err != nil {
		g.Log().Fatal(ctx, err)
	}
	grace.Register(ctx, "deregister_registry", func() {
		_ = registry.Registry.Deregister(ctx)
	})
}
