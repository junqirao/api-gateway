package proxy

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"

	"api-gateway/internal/components/config"
	"api-gateway/internal/components/loadbalance"
	"api-gateway/internal/components/response"
	"api-gateway/internal/components/upstream"
	"api-gateway/internal/model"
	"api-gateway/internal/service"
)

type sProxy struct {
}

func init() {
	service.RegisterProxy(&sProxy{})
}

// Proxy handles balancer, rate limiter, circuit breaker, retry, response write;
// 1. fetch upstreams by using RoutingKey in model.ReverseProxyInput,
// 2. balancer select node
// 3. rate limiter limits
// 4. circuit breaker
// 5. proxy request to next layer -> upstream.Upstream
// 6. control retry and response write
func (s sProxy) Proxy(ctx context.Context, input *model.ReverseProxyInput) {
	var (
		upstreams, ok = upstream.GetService(input.RoutingKey)
	)

	if !ok || upstreams == nil {
		// service not found response 503
		response.WriteJSON(input.Request,
			response.CodeUnavailable.WithDetail(input.RoutingKey))
		return
	}

	// proxy with retry
	retryCount := upstreams.Config.ReverseProxy.RetryCount
	canRetry, code := s.doProxy(ctx, upstreams, input)
	if code == nil {
		// already response by upstream.Upstream
		return
	}
	// break retry if no other upstream to select
	if upstreams.CountUpstream() <= 1 {
		response.WriteJSON(input.Request, code)
		return
	}

	// retry loop
	for canRetry && retryCount > 0 {
		g.Log().Infof(ctx, "retry proxy, count: %d, reason: %v", retryCount, code)
		canRetry, code = s.doProxy(ctx, upstreams, input)
		retryCount--
	}

	if code == nil {
		// retry succeeded, response by upstream.Upstream
		return
	}

	response.WriteJSON(input.Request, code)
}

func (s sProxy) doProxy(ctx context.Context, upstreams *upstream.Service, input *model.ReverseProxyInput) (canRetry bool, code *response.Code) {
	ups, ok := upstreams.Select(input.Request, loadbalance.GetOrCreate(input.RoutingKey).Selector)
	if !ok {
		// 503
		code = response.CodeUnavailable.WithDetail(input.RoutingKey)
		return
	}
	cb, code := ups.Allow(ctx)
	if code != nil {
		// 429,500,503
		if code.Code() != http.StatusTooManyRequests && ups.Parent.CountUpstream() > 1 {
			// if not 429 and has more than 1 upstream, can retry
			canRetry = true
		}
		return
	}
	e := ups.Do(ctx, input.Request, cb)
	if e != nil {
		// can retry
		canRetry = true
		g.Log().Warningf(ctx, "error caused during proxy: %s", e.Error())
		// 502
		if config.Gateway.Debug {
			// response detail in debug mode
			code = response.CodeBadGateway.WithMessage(e.Error())
		} else {
			code = response.CodeBadGateway.WithDetail(input.RoutingKey)
		}
	}
	return
}
