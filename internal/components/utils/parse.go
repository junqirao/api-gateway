package utils

import (
	"strings"

	"api-gateway/internal/components/config"
)

func ParseRoutingKey(u string) string {
	if config.Gateway.Prefix == "" || u == "" {
		return u
	}
	rk := strings.TrimPrefix(u, config.Gateway.Prefix)
	if i := strings.Index(rk, "/"); i != -1 {
		rk = rk[:i]
	}
	return rk
}
