package upstream

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"

	"api-gateway/internal/components/config"
	"api-gateway/internal/consts"
)

type (
	proxy2httpHandler struct {
		cfg          func() *config.ReverseProxyConfig
		scheme       string
		host         string
		prefixLength int
		routingKey   string
		dialer       *net.Dialer

		httputil.ReverseProxy
	}
	resultCallback func(err error)
)

func newHTTPHandler(host string, port int, serviceName string, cfg func() *config.ReverseProxyConfig) *proxy2httpHandler {
	var (
		scheme              = cfg().Scheme
		dialTimeout         = time.Second * 1
		tlsHandshakeTimeout = time.Second * 1
	)

	if scheme == "" {
		scheme = "http"
	}
	if cfg().DialTimeout != "" {
		duration, err := time.ParseDuration(cfg().DialTimeout)
		if err == nil {
			dialTimeout = duration
		}
	}
	if cfg().TlsHandshakeTimeout != "" {
		duration, err := time.ParseDuration(cfg().TlsHandshakeTimeout)
		if err == nil {
			tlsHandshakeTimeout = duration
		}
	}

	targetHost := fmt.Sprintf("%s://%s", scheme, host)
	if port > 0 {
		targetHost = targetHost + fmt.Sprintf(":%d", port)
	}
	target, _ := url.Parse(targetHost)
	handler := &proxy2httpHandler{
		cfg:          cfg,
		scheme:       target.Scheme,
		host:         target.Host,
		prefixLength: len(serviceName),
		routingKey:   serviceName,
		dialer: &net.Dialer{
			Timeout:   dialTimeout,
			KeepAlive: 60 * time.Second,
		},
	}
	handler.ReverseProxy = httputil.ReverseProxy{
		Director: handler.director,
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           handler.dialer.DialContext,
			TLSHandshakeTimeout:   tlsHandshakeTimeout,
			IdleConnTimeout:       90 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
		ErrorHandler:   handler.errorHandler,
		ModifyResponse: handler.responseModifier,
		BufferPool:     newBufferPool(),
	}
	return handler
}

func (h *proxy2httpHandler) director(req *http.Request) {
	req.URL.Scheme = h.scheme
	req.URL.Host = h.host
	req.Host = h.host
	if h.cfg().TrimRoutingKeyPrefix {
		if i := strings.Index(req.URL.Path, h.routingKey); i != -1 {
			req.URL.Path = req.URL.Path[i+h.prefixLength:]
		}
	}
	if _, ok := req.Header["User-Agent"]; !ok {
		// explicitly disable User-Agent, so it's not set to default value
		req.Header.Set("User-Agent", "")
	}
}

func (h *proxy2httpHandler) errorHandler(_ http.ResponseWriter, request *http.Request, err error) {
	if v := request.Context().Value(consts.CtxKeyResultCallback); v != nil {
		if cb, ok := v.(resultCallback); ok && cb != nil {
			cb(err)
		}
	}
}

func (h *proxy2httpHandler) responseModifier(resp *http.Response) error {
	if !h.cfg().ResponseBuffering ||
		// avoid large response body
		resp.ContentLength > consts.BufferingResponseMaxBodySize {
		return nil
	}

	bb := &bytes.Buffer{}
	if _, err := bb.ReadFrom(resp.Body); err != nil {
		return err
	}
	resp.Body = io.NopCloser(bb)
	return nil
}

// Do proxy request and report error
func (h *proxy2httpHandler) Do(ctx context.Context, req *ghttp.Request) (err error) {
	// prepare
	var (
		cb resultCallback = func(e error) { err = e }
	)

	// serve proxy

	// err will be assigned if caused by ServeHTTP,
	// and ServeHTTP will call errorHandler, try to
	// pass back error by resultCallback in ctx.
	h.ServeHTTP(
		// use unbuffered response raw writer,
		// make sure response body write properly.
		req.Response.RawWriter(),
		// ctx from req.Request, processed by
		// goframe at webservice entrance.
		req.Request.WithContext(context.WithValue(ctx, consts.CtxKeyResultCallback, cb)),
	)

	return
}
