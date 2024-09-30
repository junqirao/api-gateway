package mirror

import (
	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	registry "github.com/junqirao/simple-registry"

	mir "api-gateway/api/mirror/v1"
	"api-gateway/internal/components/proxy"
	r "api-gateway/internal/components/registry"
	"api-gateway/internal/consts"
)

var (
	cli           *client
	serverAddress = g.Cfg().MustGet(gctx.GetInitCtx(), "mirror.client.server_address", "127.0.0.1:8001").String()
)

type (
	ClientInfo struct {
		Instance *registry.Instance `json:"instance"`
		Filter   []string           `json:"filter"`
	}
	client struct {
		// only in server side
		handler        proxy.ReverseProxyHandler
		ctx            context.Context
		cancel         context.CancelFunc
		acceptInterval int64
		// client info
		ClientInfo
	}
)

func newClient(ctx context.Context) (c *client, err error) {
	c = &client{
		ClientInfo: ClientInfo{
			Instance: r.CurrentInstance,
			Filter:   g.Cfg().MustGet(ctx, "mirror.client.filter", []string{}).Strings(),
		},
	}

	c.ctx, c.cancel = context.WithCancel(ctx)

	if c.acceptInterval, err = c.register(ctx); err != nil {
		g.Log().Errorf(ctx, "mirror client failed to register: %v", err)
		return
	}

	interval := c.acceptInterval - 5
	if interval <= 0 {
		g.Log().Errorf(c.ctx, "mirror server heartbeat interval is too small: %v", interval)
		return
	}

	go c.keepalive(time.NewTicker(time.Duration(interval) * time.Second))
	g.Log().Infof(ctx, "mirror client registered key=%s", c.Instance.Id)
	return
}

func (c *client) register(ctx context.Context) (interval int64, err error) {
	cc, err := grpcx.Client.NewGrpcClientConn(serverAddress)
	if err != nil {
		return
	}
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*1)
	defer cancelFunc()
	res, err := mir.NewMirrorClient(cc).Register(ctx, &mir.RegisterReq{Instance: convert2ins(c.Instance), Filter: c.Filter})
	if err != nil {
		return
	}
	if c.acceptInterval != 0 {
		if c.acceptInterval != gconv.Int64(res.Ttl) {
			err = fmt.Errorf("server interval config changed %d->%d. client stop", c.acceptInterval, res.Ttl)
			_ = c.Close()
			return
		}
	} else {
		c.acceptInterval = gconv.Int64(res.Ttl)
	}
	return int64(res.Ttl), nil
}

func (c *client) keepalive(ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			if _, err := c.register(c.ctx); err != nil {
				g.Log().Errorf(c.ctx, "mirror client failed to register: %v", err)
				continue
			}
			// g.Log().Info(c.ctx, "mirror client keepalive success.")
		case <-c.ctx.Done():
			g.Log().Infof(c.ctx, "mirror client stop keepalive")
			return
		}
	}
}

func (c *client) deregister(ctx context.Context) (err error) {
	cc, err := grpcx.Client.NewGrpcClientConn(serverAddress)
	if err != nil {
		return
	}
	ctx, cancelFunc := context.WithTimeout(ctx, time.Second*1)
	defer cancelFunc()
	_, err = mir.NewMirrorClient(cc).UnRegister(ctx, &mir.UnRegisterReq{Instance: convert2ins(c.Instance)})
	return
}

func (c *client) Do(ctx context.Context, req *ghttp.Request) (err error) {
	v := ctx.Value(consts.CtxKeyRoutingKey)
	if v == nil {
		return
	}

	if len(c.Filter) > 0 {
		key := v.(string)
		for _, s := range c.Filter {
			if key == s {
				return c.handler.Do(ctx, req)
			}
		}
		return
	}

	return c.handler.Do(ctx, req)
}

func (c *client) Close() error {
	c.cancel()
	return c.deregister(c.ctx)
}

func convert2ins(instance *registry.Instance) *mir.Instance {
	return &mir.Instance{
		Id:          instance.Id,
		Host:        instance.Host,
		HostName:    instance.HostName,
		Port:        int32(instance.Port),
		ServiceName: instance.ServiceName,
		Meta:        gconv.String(instance.Meta),
	}
}
