package mirror

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/grace"
	"api-gateway/internal/components/proxy"
)

/**
# config example yaml

mirror:
  # working mode "","client","server"
  mode: "server"
  server:
    # channel buffer size, default 100
    ch_buffer_size: 100
    # async worker count, default 10
    worker_count: 10
    # client heartbeat interval, default 30
    heartbeat_interval: 30
  client:
    server_address: "127.0.0.1:8001"
    # service_name going to replicate to,
    # if not set will replicate to all
    filter:
      - "service_name_1"
      - "service_name_2"
*/

const (
	ModeServer = "server"
	ModeClient = "client"
	ModeNone   = ""
)

func Init(ctx context.Context, newFunc func(ins *registry.Instance) (proxy.ReverseProxyHandler, error)) {
	mode := g.Cfg().MustGet(ctx, "mirror.mode", ModeNone).String()
	if mode == ModeNone {
		return
	}

	g.Log().Infof(ctx, "mirror mode: %s", mode)
	if mode == ModeServer {
		srv = newServer(ctx, newFunc)
		grace.Register(ctx, "stop_mirror_server", func() {
			_ = srv.Close()
		})
	} else if mode == ModeClient {
		var err error
		cli, err = newClient(ctx)
		if err != nil {
			g.Log().Fatalf(ctx, "create mirror client failed: %v", err)
			return
		}
		grace.Register(ctx, "stop_mirror_client", func() {
			_ = cli.Close()
		})
	}
}
