package upstream

import (
	"context"
	"net/http/httputil"
	"sync"
	"sync/atomic"

	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/utils"
	"api-gateway/internal/consts"
)

type (
	retryableProxyHandler struct {
		upstream *Upstream
		bufPool  *sync.Pool
		next     ReverseProxyHandler
	}
	bufferPool struct {
		p *sync.Pool
	}
)

func NewHandler(_ context.Context, ups *Upstream, cfg config.ReverseProxyConfig) (handler ReverseProxyHandler) {
	return newRetryableHandler(
		newHTTPHandler(ups.Instance.Host, ups.Instance.Port, ups.ServiceName, &cfg),
		ups)
}

func newRetryableHandler(h ReverseProxyHandler, upstream *Upstream) ReverseProxyHandler {
	return &retryableProxyHandler{
		upstream: upstream,
		next:     h,
		bufPool: &sync.Pool{New: func() any {
			return utils.NewNopCloseBuf()
		}},
	}
}

func (h *retryableProxyHandler) Do(ctx context.Context, req *ghttp.Request) (err error) {
	var (
		retried  = ctx.Value(consts.CtxKeyRetriedTimes).(*atomic.Int64).Load()
		canRetry = ctx.Value(consts.CtxKeyCanRetry).(*atomic.Bool).Load()
	)

	// when canRetry, replace request body with
	// nop closer at first place to avoid read
	// closed body during retry.
	if canRetry && retried == 0 {
		buf := h.bufPool.Get().(*utils.NopCloseBuf)
		// original request body will be closed when read EOF
		buf.SetOrigin(req.Request.Body)
		req.Request.Body = buf
	}

	recovery := func() {
		retryCount := h.retryCount()
		if buf, ok := req.Request.Body.(*utils.NopCloseBuf); ok {
			// when proxy success or reached retry limit
			// put back buffer, only work if canRetry == true
			if canRetry && err == nil || retryCount-retried <= 0 {
				buf.Reset()
				h.bufPool.Put(buf)
			} else {
				// reset buffer index for retry or mirror request
				buf.ResetIndex()
			}
		}
	}

	err = h.next.Do(ctx, req)
	recovery()
	return
}

func (h *retryableProxyHandler) retryCount() int64 {
	return int64(h.upstream.Parent.Config.ReverseProxy.RetryCount)
}

func newBufferPool() httputil.BufferPool {
	return &bufferPool{
		p: &sync.Pool{
			New: func() interface{} {
				return make([]byte, 32*1024)
			},
		},
	}
}

func (b *bufferPool) Get() []byte {
	return b.p.Get().([]byte)
}

func (b *bufferPool) Put(bytes []byte) {
	b.p.Put(bytes)
}
