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

	"api-gateway/internal/consts"
	"api-gateway/internal/model"
)

type (
	proxy2httpHandler struct {
		upstream     *Upstream
		cfg          *model.ReverseProxyConfig
		scheme       string
		host         string
		prefixLength int
		routingKey   string
		dialer       *net.Dialer
		httputil.ReverseProxy
	}
	resultCallback func(err error)
)

func newHTTPHandler(upstream *Upstream, cfg *model.ReverseProxyConfig) *proxy2httpHandler {
	var (
		ins                 = &upstream.Instance
		scheme              = cfg.Scheme
		dialTimeout         = time.Second * 1
		tlsHandshakeTimeout = time.Second * 1
	)

	if scheme == "" {
		scheme = "http"
	}
	if cfg.DialTimeout != "" {
		duration, err := time.ParseDuration(cfg.DialTimeout)
		if err == nil {
			dialTimeout = duration
		}
	}
	if cfg.TlsHandshakeTimeout != "" {
		duration, err := time.ParseDuration(cfg.TlsHandshakeTimeout)
		if err == nil {
			tlsHandshakeTimeout = duration
		}
	}

	targetHost := fmt.Sprintf("%s://%s", scheme, ins.Host)
	if ins.Port > 0 {
		targetHost = targetHost + fmt.Sprintf(":%d", ins.Port)
	}
	target, _ := url.Parse(targetHost)
	handler := &proxy2httpHandler{
		upstream:     upstream,
		cfg:          cfg,
		scheme:       target.Scheme,
		host:         target.Host,
		prefixLength: len(ins.ServiceName),
		routingKey:   ins.ServiceName,
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
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          300,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   tlsHandshakeTimeout,
			ExpectContinueTimeout: 1 * time.Second,
		},
		ErrorHandler:   handler.errorHandler,
		ModifyResponse: handler.responseModifier,
	}
	return handler
}

func (h *proxy2httpHandler) director(req *http.Request) {
	req.URL.Scheme = h.scheme
	req.URL.Host = h.host
	req.Host = h.host
	if h.cfg.TrimRoutingKeyPrefix {
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
	// if config.Gateway.Debug {
	// 	// add server id
	// 	resp.Header.Set(consts.HeaderKeyServerId, h.upstream.Instance.Id)
	// 	resp.Header.Set(consts.HeaderKeyServerAddr, fmt.Sprintf("%s:%d", h.upstream.Instance.Host, h.upstream.Instance.Port))
	// 	resp.Header.Set(consts.HeaderKeyServerHostName, h.upstream.Instance.HostName)
	// 	resp.Header.Set(consts.HeaderKeyServiceUpstreamCount, fmt.Sprintf("%d", h.upstream.Parent.CountUpstream()))
	// 	resp.Header.Set(consts.HeaderKeyServiceAvailableUpstreamCount, fmt.Sprintf("%d", h.upstream.Parent.CountAvailableUpstream()))
	// }
	// if return err!=nil will call h.errorHandler
	return nil
}

// Do proxy request and report error
func (h *proxy2httpHandler) Do(ctx context.Context, req *ghttp.Request) (err error) {
	// prepare
	var (
		bs []byte
		cb resultCallback = func(e error) { err = e }
	)

	if req.ContentLength != 0 {
		// read and close request body
		if bs, err = io.ReadAll(req.Body); err != nil {
			return
		}
		if err = req.Body.Close(); err != nil {
			return
		}

		// set request body
		req.Body = io.NopCloser(bytes.NewBuffer(bs))
	}

	// serve proxy
	h.ServeHTTP(
		// use unbuffered response raw writer, make sure response body write properly
		req.Response.RawWriter(),
		// ctx from req.Request, processed by goframe at webservice entrance
		req.Request.WithContext(context.WithValue(ctx, consts.CtxKeyResultCallback, cb)),
	)

	// when err not nil may cause retry and old body may be closed
	if err != nil && bs != nil {
		req.Request.Body = io.NopCloser(bytes.NewBuffer(bs))
	}
	return
}
