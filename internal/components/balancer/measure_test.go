package balancer

import (
	"testing"
	"time"
)

func TestMeasurable(t *testing.T) {
	var size int64 = 3
	measure := NewMeasurable(time.Second * time.Duration(size))
	for i := 0; i < int(size); i++ {
		measure.AddLoad(1)
	}

	t.Logf("wait %d seconds", size)
	time.Sleep(time.Second * time.Duration(size))
	for i := 0; i < int(size)+1; i++ {
		v := measure.Load()
		expect := size - int64(i)
		if v != expect {
			t.Fatalf("[%d] expect %d, got %d", i, expect, v)
		}
		t.Logf("[%d] %d", i, v)
		time.Sleep(time.Second)
	}
}

func BenchmarkMeasurable_Load(b *testing.B) {
	measure := NewMeasurable(time.Second)
	measure.AddLoad(1000000)
	for i := 0; i < b.N; i++ {
		measure.Load()
	}
}

func BenchmarkMeasurable_AddLoad(b *testing.B) {
	measure := NewMeasurable(time.Second)
	for i := 0; i < b.N; i++ {
		if i%2 > 0 {
			measure.AddLoad(1)
		} else {
			measure.AddLoad(-1)
		}
	}
}
