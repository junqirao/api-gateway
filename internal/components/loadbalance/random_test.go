package loadbalance

import (
	"testing"
)

func TestRandom_Selector(t *testing.T) {
	balancer := newRandom()
	cnt := map[int]int{}
	for i := 0; i < 100; i++ {
		ref, ok := balancer.Selector(nil, []Weighted{
			NewWeighted(3, 0), // ≈ 60
			NewWeighted(1, 1), // ≈ 20
			NewWeighted(1, 2), // ≈ 20
		})
		if !ok {
			t.Fatalf("expected ok")
		}
		cnt[ref]++
		t.Logf("ref: %d, ok: %v", ref, ok)
	}

	t.Logf("cnt: %+v", cnt)
}

func BenchmarkRandom_Selector(b *testing.B) {
	balancer := newRandom()
	size := 20
	var ws []Weighted
	for i := 0; i < size; i++ {
		ws = append(ws, NewWeighted(1, i))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, ok := balancer.Selector(nil, ws); !ok {
			b.Fatalf("expected ok")
		}
	}
}
