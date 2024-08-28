// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"api-gateway/internal/model"
	"context"
)

type (
	IProxy interface {
		// Proxy handles balancer, rate limiter, circuit breaker, retry, response write;
		// 1. fetch upstreams by using RoutingKey in model.ReverseProxyInput,
		// 2. balancer select node
		// 3. rate limiter limits
		// 4. circuit breaker
		// 5. proxy request to next layer -> upstream.Upstream
		// 6. control retry and response write
		Proxy(ctx context.Context, input *model.ReverseProxyInput)
	}
)

var (
	localProxy IProxy
)

func Proxy() IProxy {
	if localProxy == nil {
		panic("implement not found for interface IProxy, forgot register?")
	}
	return localProxy
}

func RegisterProxy(i IProxy) {
	localProxy = i
}
