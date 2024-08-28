package reverse

import (
	"testing"

	"api-gateway/internal/components/config"
	"api-gateway/internal/model"
)

func BenchmarkParse(b *testing.B) {
	// BenchmarkParse-20    	230034314	         5.157 ns/op
	// cpu: 12th Gen Intel(R) Core(TM) i7-12700F
	config.Gateway = &model.GatewayConfig{
		Prefix: "/api/",
	}
	var s = ""
	for i := 0; i < b.N; i++ {
		s = parseRoutingKey("/api/routing_key/test1/test2/test3?a=1&b=2")
	}
	b.Log(s)
}
