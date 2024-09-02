package program

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	registry "github.com/junqirao/simple-registry"
)

type (
	requestWrapper struct {
		RemoteAddr string         `json:"'remote_addr'"`
		Host       string         `json:"'host'"`
		URL        string         `json:"'url'"`
		Method     string         `json:"'method'"`
		Header     *headerWrapper `json:"header"`
	}
	responseWrapper struct {
		Header *headerWrapper `json:"header"`
	}
	logWrapper struct{}
)

var (
	// logWrap logger wrapper
	logWrap = &logWrapper{}
)

func (w logWrapper) Info(ctx context.Context, format string, v ...interface{}) bool {
	g.Log().Infof(ctx, format, v...)
	return true
}

func (w logWrapper) Warn(ctx context.Context, format string, v ...interface{}) bool {
	g.Log().Warningf(ctx, format, v...)
	return true
}

func (w logWrapper) Error(ctx context.Context, format string, v ...interface{}) bool {
	g.Log().Errorf(ctx, format, v...)
	return true
}

func BuildEnvFromRequest(r *ghttp.Request, ups registry.Instance) map[string]interface{} {
	return map[string]interface{}{
		envKeyLogger:   logWrap,
		envKeyUpstream: ups,
		envKeyRequest: &requestWrapper{
			RemoteAddr: r.GetRemoteIp(),
			Host:       r.GetHost(),
			URL:        r.GetUrl(),
			Method:     r.Method,
			Header:     newHeaderWrapper(r.Request.Header),
		},
		envKeyResponse: &responseWrapper{
			Header: newHeaderWrapper(r.Response.Header()),
		},
	}
}
