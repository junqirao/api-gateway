package loadbalance

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

func TestNextWeighted(t *testing.T) {
	ws := []Weighted{
		NewWeighted(1, 0),
		NewWeighted(1, 1),
		NewWeighted(1, 2),
	}

	for i := 0; i < 8; i++ {
		w := nextWeighted(ws)
		t.Logf("%v: %+v", w.Ref(), w)
	}
	fmt.Println("----------")
	ws[0].SetWeight(3)
	for i := 0; i < 15; i++ {
		w := nextWeighted(ws)
		t.Logf("%v: %+v", w.Ref(), w)
	}
}

func BenchmarkNextWeighted(b *testing.B) {
	size := 20
	var ws []Weighted
	for i := 0; i < size; i++ {
		ws = append(ws, NewWeighted(1, i))
	}
	for i := 0; i < b.N; i++ {
		nextWeighted(ws)
	}
}

func TestNextWeightedConcurrence(t *testing.T) {
	ws := []Weighted{
		NewWeighted(3, 0),
		NewWeighted(1, 1),
		NewWeighted(1, 2),
	}
	batch := 1000
	worker := 10
	cnt := sync.Map{}
	cnt.Store("0", new(atomic.Int64))
	cnt.Store("1", new(atomic.Int64))
	cnt.Store("2", new(atomic.Int64))

	wg := sync.WaitGroup{}
	for i := 0; i < worker; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < batch; j++ {
				ref := nextWeighted(ws).Ref()
				v, _ := cnt.Load(ref)
				value := v.(*atomic.Int64)
				value.Add(1)
			}
			wg.Done()
		}()
	}
	wg.Wait()

	if v, ok := cnt.Load("0"); !ok || v.(*atomic.Int64).Load() != 6000 {
		t.Fatalf("expected %d, got %d", 6000, v.(*atomic.Int64).Load())
	}
	if v, ok := cnt.Load("1"); !ok || v.(*atomic.Int64).Load() != 2000 {
		t.Fatalf("expected %d, got %d", 2000, v.(*atomic.Int64).Load())
	}
	if v, ok := cnt.Load("2"); !ok || v.(*atomic.Int64).Load() != 2000 {
		t.Fatalf("expected %d, got %d", 2000, v.(*atomic.Int64).Load())
	}
}
