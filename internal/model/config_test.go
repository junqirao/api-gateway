package model

import (
	"testing"
)

func TestConfigChanged(t *testing.T) {
	v1 := ReverseProxyConfig{
		TrimRoutingKeyPrefix: false,
		RetryCount:           2,
		DialTimeout:          "",
		TlsHandshakeTimeout:  "",
		Scheme:               "",
	}
	v2 := v1

	if ValueChanged(v1, v2) {
		t.Fatal("value did not change")
	}

	v2.RetryCount = 3
	if !ValueChanged(v1, v2) {
		t.Fatal("value changed")
	}
}
