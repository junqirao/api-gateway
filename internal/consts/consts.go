package consts

// context key define
const (
	CtxKeyResultCallback = "___result_callback"
	CtxKeyRetriedTimes   = "___retried_times"
	CtxKeyCanRetry       = "___can_retry"
)

// header key define
const (
	HeaderKeyContentType                   = "Content-Type"
	HeaderKeyServerId                      = "X-Srv-Instance-Id"
	HeaderKeyServerHostName                = "X-Srv-Host-Name"
	HeaderKeyServerAddr                    = "X-Srv-Address"
	HeaderKeyServiceUpstreamCount          = "X-Svc-Upstream-Count"
	HeaderKeyServiceAvailableUpstreamCount = "X-Svc-Available-Upstream-Count"
)

const (
	HeaderValueContentTypeJSON = "application/json"
)

// module name define
const (
	ModuleNameReverseProxy = "reverse_proxy"
	ModuleNameLoadBalance  = "load_balance"
	ModuleNameRateLimiter  = "rate_limiter"
	ModuleNameBreaker      = "breaker"
)

// default values
const (
	DefaultMaxBodySize    = 512 * 1024 * 1024 // 512M
	RetryMaxContentLength = 10 * 1024 * 1024  // 10M
)

// storage name define
const (
	StorageNameServiceConfig = "service_config"
	StorageNameProgram       = "program"
	StorageNameVariable      = "program_variable"
)
