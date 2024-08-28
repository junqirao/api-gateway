package consts

const (
	GatewayConfigPattern         = "gateway"
	GatewayConfigPrefixPattern   = "gateway.prefix"
	GatewayConfigBalancerPattern = "gateway.load_balance"
	GatewayConfigLimiterPattern  = "gateway.limiter"
	ProxyTypeHTTP                = "http"
)

const (
	ConfigScopeReverseProxy = "gateway_config_reverse_proxy"
	ConfigScopeBalancer     = "gateway_config_balancer"
	ConfigScopeLimiter      = "gateway_config_limiter"
	ConfigScopeBreaker      = "gateway_config_breaker"
)
