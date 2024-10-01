package program

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	registry "github.com/junqirao/simple-registry"

	"api-gateway/internal/components/program/extra/ipgeo"
)

type (
	requestWrapper struct {
		ClientIP   string         `json:"client_ip"`
		RemoteAddr string         `json:"'remote_addr'"`
		Host       string         `json:"'host'"`
		URL        string         `json:"'url'"`
		Method     string         `json:"'method'"`
		Header     *headerWrapper `json:"header"`
	}
	responseWrapper struct {
		Header *headerWrapper `json:"header"`
	}
	logWrapper struct {
		ctx context.Context
	}
)

func exprMultilineWrapper(lines ...*resultWrapper) (errMsg string) {
	for _, line := range lines {
		if b, reason := line.Ok(); !b {
			return reason
		}
	}
	return
}

func (w logWrapper) Info(v ...interface{}) bool {
	g.Log().Info(w.ctx, v...)
	return true
}

func (w logWrapper) Warn(v ...interface{}) bool {
	g.Log().Warning(w.ctx, v...)
	return true
}

func (w logWrapper) Error(v ...interface{}) bool {
	g.Log().Error(w.ctx, v...)
	return true
}

func (w logWrapper) Infof(format string, v ...interface{}) bool {
	g.Log().Infof(w.ctx, format, v...)
	return true
}

func (w logWrapper) Warnf(format string, v ...interface{}) bool {
	g.Log().Warningf(w.ctx, format, v...)
	return true
}

func (w logWrapper) Errorf(format string, v ...interface{}) bool {
	g.Log().Errorf(w.ctx, format, v...)
	return true
}

func BuildEnvFromRequest(ctx context.Context, r *ghttp.Request, ups registry.Instance) map[string]interface{} {
	clientIp := r.GetClientIp()
	return map[string]interface{}{
		// base
		envKeyNewResultWrapper:     newResultWrapper,
		envKeyExprMultilineWrapper: exprMultilineWrapper,
		envKeyCtx:                  ctx,
		// logger
		envKeyLogger: logWrapper{ctx: ctx},
		// common
		envKeyTerminateIf: func(flag bool, reason ...string) error {
			if !flag {
				return nil
			}
			reasonStr := "request terminated"
			if len(reason) > 0 && reason[0] != "" {
				reasonStr = reason[0]
			}
			return errors.New(reasonStr)
		},
		envKeyIPGEO: &ipgeo.Wrapper{Address: clientIp},
		// variables
		envKeyGlobalVariable: Variables.GetGlobalVariables(ctx),
		envKeySetGlobalVariable: func(key string, value interface{}) error {
			return Variables.SetGlobalVariable(ctx, key, value)
		},
		// runtime
		envKeyUpstream: ups,
		envKeyRequest: &requestWrapper{
			ClientIP:   clientIp,
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
