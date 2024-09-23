package consts

const (
	CtxKeyResultCallback = "___result_callback"
	CtxKeyRetriedTimes   = "___retried_times"
)

const (
	HeaderKeyServerId                      = "X-Srv-Instance-Id"
	HeaderKeyServerHostName                = "X-Srv-Host-Name"
	HeaderKeyServerAddr                    = "X-Srv-Address"
	HeaderKeyServiceUpstreamCount          = "X-Svc-Upstream-Count"
	HeaderKeyServiceAvailableUpstreamCount = "X-Svc-Available-Upstream-Count"
)

const (
	ModuleNameLoadBalance = "load_balance"
	ModuleNameRateLimiter = "rate_limiter"
	ModuleNameBreaker     = "breaker"
)
