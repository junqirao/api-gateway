package limiter

import (
	"fmt"
	"testing"

	"api-gateway/internal/components/config"
)

func TestLimiter(t *testing.T) {
	l := NewLimiter(config.RateLimiterConfig{
		Rate: 1,
		Peak: 2,
	})
	for i := 0; i < 3; i++ {
		ok := l.Allow()
		t.Log(fmt.Sprintf("ok: %v, i: %d", ok, i))
		if !ok && i == 2 {
			return
		}
	}
	t.Fatal("fail")
}
