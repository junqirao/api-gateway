package mirror

import (
	"context"
	"net/http"
	"sync"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/proxy"
	r "api-gateway/internal/components/registry"
	"api-gateway/internal/components/utils"
	"api-gateway/internal/consts"
)

var (
	srv                     *server
	chBufferSize            = g.Cfg().MustGet(gctx.GetInitCtx(), "mirror.server.ch_buffer_size", 100).Int()
	workerCount             = g.Cfg().MustGet(gctx.GetInitCtx(), "mirror.server.worker_count", 10).Int()
	clientHeartbeatInterval = g.Cfg().MustGet(gctx.GetInitCtx(), "mirror.server.heartbeat_interval", 30).Int64()
	whiteList               = g.Cfg().MustGet(gctx.GetInitCtx(), "mirror.server.white_list", "").Strings()
)

type (
	server struct {
		mu        sync.RWMutex
		clients   map[string]*client // key: client_id, value: proxy.ReverseProxyHandler
		new       func(ins *registry.Instance) (proxy.ReverseProxyHandler, error)
		ch        chan *request
		ctx       context.Context
		cancel    context.CancelFunc
		wg        sync.WaitGroup
		whiteList map[string]struct{}
	}
	request struct {
		ctx   context.Context
		req   *ghttp.Request
		after func()
	}
	nopWriter struct{}
)

func newServer(ctx context.Context, newFunc func(ins *registry.Instance) (proxy.ReverseProxyHandler, error)) *server {
	s := &server{
		clients: make(map[string]*client),
		new:     newFunc,
		ch:      make(chan *request, chBufferSize),
	}

	s.whiteList = make(map[string]struct{})
	if len(whiteList) > 0 {
		for _, v := range whiteList {
			s.whiteList[v] = struct{}{}
		}
		g.Log().Infof(ctx, "mirror server white list: %v", whiteList)
	}

	s.buildClients(ctx)
	g.Log().Infof(ctx, "mirror server build clients: %d", len(s.clients))
	s.ctx, s.cancel = context.WithCancel(ctx)
	registry.Storages.SetEventHandler(consts.StorageNameMirror, s.eventHandler(ctx))
	s.wg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go s.runWorker()
	}
	g.Log().Infof(s.ctx, "mirror server worker started: %d", workerCount)
	return s
}

func (s *server) eventHandler(ctx context.Context) registry.StorageEventHandler {
	return func(t registry.EventType, key string, value interface{}) {
		// avoid loop
		if key == r.CurrentInstance.Id {
			return
		}

		s.mu.Lock()
		defer s.mu.Unlock()

		switch t {
		case registry.EventTypeCreate:
			s.buildClients(ctx, key)
		case registry.EventTypeDelete:
			delete(s.clients, key)
			g.Log().Infof(ctx, "mirror server deregistered client key=%s", key)
		}
	}
}

func (s *server) buildClients(ctx context.Context, key ...string) {
	kvs, err := registry.Storages.GetStorage(consts.StorageNameMirror).Get(ctx, key...)
	if err != nil {
		g.Log().Errorf(ctx, "mirror server failed to build clients: %v", err)
		return
	}

	for _, kv := range kvs {
		var (
			c = &client{}
		)

		if err = kv.Value.Struct(&c); err != nil {
			g.Log().Errorf(ctx, "unexpected value from mirror storage: %v", err)
			return
		}
		c.handler, err = s.new(c.Instance)
		if err != nil {
			g.Log().Errorf(ctx, "mirror server failed to create client: %v", err)
			return
		}
		s.clients[c.Instance.Id] = c
		g.Log().Infof(ctx, "mirror server registered client key=%s target=%s:%d", c.Instance.Id, c.Instance.Host, c.Instance.Port)
	}
}

func (s *server) runWorker() {
	for {
		select {
		case rr := <-s.ch:
			s.doRequest(rr)
		case <-s.ctx.Done():
			s.wg.Done()
			return
		}
	}
}

func (s *server) Close() error {
	s.cancel()
	close(s.ch)
	s.wg.Wait()
	g.Log().Info(s.ctx, "mirror server worker all stopped.")
	for rr := range s.ch {
		if rr == nil {
			break
		}
		s.doRequest(rr)
	}
	return nil
}

func (s *server) doRequest(request *request) {
	srv.mu.Lock()
	defer func() {
		srv.mu.Unlock()
		request.after()
	}()

	for _, c := range srv.clients {
		request.ResetRequestBody()
		if err := c.handler.Do(request.ctx, request.req); err != nil {
			g.Log().Warningf(request.ctx, "mirror server failed to push: %v", err)
			continue
		}
	}
}

func Push(ctx context.Context, req *ghttp.Request, after func()) bool {
	if srv == nil {
		return false
	}

	req.Header.Add(consts.HeaderKeyReplicateFrom, r.CurrentInstance.Id)
	req.Response.ResponseWriter = &nopWriter{}
	ctxId := gctx.CtxId(ctx)
	ctx, _ = gtrace.WithTraceID(context.Background(), ctxId)

	srv.ch <- &request{ctx: ctx, req: req, after: after}
	return true
}

func (n *nopWriter) Header() http.Header {
	return http.Header{}
}

func (n *nopWriter) Write(bytes []byte) (int, error) {
	return len(bytes), nil
}

func (n *nopWriter) WriteHeader(_ int) {
	return
}

func (r *request) ResetRequestBody() {
	// reset buf index
	if buf, ok := r.req.Request.Body.(*utils.NopCloseBuf); ok {
		buf.ResetIndex()
	}
}

func Register(ctx context.Context, info ClientInfo) (ttl int64, err error) {
	return clientHeartbeatInterval, registry.Storages.GetStorage(consts.StorageNameMirror).
		SetTTL(ctx, info.Instance.Id, info, clientHeartbeatInterval)
}

func UnRegister(ctx context.Context, ins *registry.Instance) (err error) {
	return registry.Storages.GetStorage(consts.StorageNameMirror).Delete(ctx, ins.Id)
}

func Allow(ip string) bool {
	if len(srv.whiteList) == 0 {
		return true
	}
	_, ok := srv.whiteList[ip]
	return ok
}
