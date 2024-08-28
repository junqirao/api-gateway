package registry

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"
)

func Init(ctx context.Context) {
	var (
		cfg = registry.Config{}
		ins = &registry.Instance{}
		v   = g.Cfg().MustGet(ctx, "registry")
	)

	if err := v.Scan(&cfg); err != nil {
		g.Log().Fatal(ctx, err)
	}
	if err := v.Scan(&ins); err != nil {
		g.Log().Fatal(ctx, err)
	}
	if err := registry.Init(ctx, cfg, ins); err != nil {
		g.Log().Fatal(ctx, err)
	}
}
