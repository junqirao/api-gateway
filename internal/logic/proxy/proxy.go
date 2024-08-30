package proxy

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"

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

	retryCount := upstreams.Config.ReverseProxy.RetryCount
	retry, err := s.doProxy(ctx, upstreams, input)
	if err == nil {
		return
	}
	// break retry if no other upstream to select
	if len(upstreams.Ups) <= 1 {
		response.WriteJSON(input.Request, err)
		return
	}
	for retry && retryCount > 0 {
		g.Log().Infof(ctx, "retry count: %d", retryCount)
		retry, err = s.doProxy(ctx, upstreams, input)
		time.Sleep(time.Millisecond * 100)
		retryCount--
	}
	if err == nil && !retry {
		return
	}
	if err == nil {
		g.Log().Errorf(ctx, "proxy error: %v", err)
		err = response.CodeBadGateway.WithDetail(input.RoutingKey)
	}
	response.WriteJSON(input.Request, err)
}

func (s sProxy) doProxy(ctx context.Context, upstreams *upstream.Service, input *model.ReverseProxyInput) (retry bool, err *response.Code) {
	ups, ok := upstreams.Select(input.Request, loadbalance.GetOrCreate(input.RoutingKey).Selector)
	if !ok {
		// 503
		err = response.CodeUnavailable.WithDetail(input.RoutingKey)
		return
	}
	allow, cb := ups.Allow(ctx)
	if !allow {
		// 429
		err = response.CodeTooManyRequests.WithDetail(input.RoutingKey)
		return
	}
	retry, e := ups.Do(ctx, input.Request, cb)
	if e != nil {
		// 502
		err = response.CodeBadGateway
	}
	return
}
