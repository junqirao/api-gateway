package registry

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/junqirao/gocomponents/kvdb"
	"github.com/junqirao/gocomponents/registry"

	"api-gateway/internal/components/grace"
)

var (
	CurrentInstance *registry.Instance
)

func Init(ctx context.Context) {
	var (
		cfg = registry.Config{}
		v   = g.Cfg().MustGet(ctx, "registry")
	)

	// parse config
	if err := v.Scan(&cfg); err != nil {
		g.Log().Fatal(ctx, err)
	}
	// init registry and register
	if err := registry.Init(ctx, kvdb.MustGetDatabase(ctx)); err != nil {
		g.Log().Fatal(ctx, err)
	}
	CurrentInstance = registry.Current()
	grace.Register(ctx, "deregister_registry", func() {
		_ = registry.Registry.Deregister(ctx)
	})
}
