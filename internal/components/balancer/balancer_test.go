package balancer

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

/**
cpu: 12th Gen Intel(R) Core(TM) i7-12700F
BenchmarkHashBalancer_Pick-20           13574460                88.19 ns/op           45 B/op          3 allocs/op
BenchmarkRandomBalancer_Pick-20         40934953                27.83 ns/op           16 B/op          1 allocs/op
BenchmarkWeightedRandom_Pick-20           570837              1973 ns/op            1000 B/op          5 allocs/op
BenchmarkWeightedRoundRobin_Pick-20       665128              1671 ns/op              16 B/op          1 allocs/op
BenchmarkRoundRobin_Pick-20             44536816                26.76 ns/op           16 B/op          1 allocs/op
BenchmarkLessLoad_Pick-20                2586336               464.7 ns/op            16 B/op          1 allocs/op
PASS

*/

type (
	testObject struct {
		Measurable
		Weighable
		val int64
	}
	wantObject struct {
		want, diff int64
	}
)

func (w wantObject) judge(v int64) bool {
	// not compare
	if w.diff == -1 {
		return true
	}
	if w.diff == 0 {
		return v == w.want
	}
	return v >= w.want-w.diff && v <= w.want+w.diff
}

func newTestObject(val int64, weight int64) *testObject {
	return &testObject{
		Measurable: NewMeasurable(val),
		Weighable:  NewWeighable(weight),
		val:        val,
	}
}

func TestBalancer_Pick(t *testing.T) {
	var (
		weight       = 10
		objs         []any
		objectSize   = 10
		batch        = 1000000
		workers      = 10
		deviation    = int64(workers)
		allowedError = workers / 3

		allowedDiff = int64(batch / objectSize / workers)

		allowedMistakeMap = map[Strategy]int64{
			StrategyRoundRobin:         allowedDiff,
			StrategyWeightedRoundRobin: 2,
			StrategyWeightedRandom:     allowedDiff,
			StrategyLessLoad:           deviation,
			// not compare diff
			StrategyRandom: -1,
			StrategyHash:   -1,
		}
		strategies = []Strategy{
			StrategyRoundRobin,
			StrategyRandom,
			StrategyWeightedRoundRobin,
			StrategyWeightedRandom,
			StrategyLessLoad,
			StrategyHash,
		}
	)

	// test same weight
	for i := 0; i < objectSize; i++ {
		objs = append(objs, newTestObject(int64(i), int64(weight)))
	}

	for _, strategy := range strategies {
		fmt.Println()
		testConcurrentPick(t, strategy, batch, workers, objs, allowedMistakeMap[strategy], allowedError)
		fmt.Println()
	}

	t.Logf("test different weight")
	objs2 := objs[:2]
	objs2[0].(*testObject).AddWeight(30) // weight=40
	// objs2[1] weight=10

	allowedMistakeMap[StrategyWeightedRoundRobin] = int64(len(objs2) * 2)
	for _, strategy := range []Strategy{
		StrategyWeightedRoundRobin,
		StrategyWeightedRandom,
	} {
		fmt.Println()
		testConcurrentPick(t, strategy, batch, workers, objs2, allowedMistakeMap[strategy], 1)
		fmt.Println()
	}
}

func testConcurrentPick(t *testing.T,
	st Strategy,
	batchSize, workerCount int,
	objs []any,
	allowedMistakeCount int64,
	allowedError int,
) {

	t.Logf("testing [%s] balancer, objs_size=%d, batchSize=%d, workerCount=%d, total_pick=%d", st, len(objs), batchSize, workerCount, batchSize*workerCount)

	// init data
	var (
		wg                = &sync.WaitGroup{}
		cnt               = atomic.Int64{}
		timeCost          = atomic.Int64{}
		wants             = map[int64]*wantObject{}
		weights           = make(map[int]float64)
		totalWeight int64 = 0

		balancer = New(st)
		res      = make(map[int64]*atomic.Int64)
	)

	for _, obj := range objs {
		totalWeight += obj.(*testObject).Weight()
	}
	for i, obj := range objs {
		weights[i] = float64(obj.(*testObject).Weight()) / float64(totalWeight)
	}

	t.Logf("totalWeight=%d, weights=%v", totalWeight, weights)

	for i, obj := range objs {
		res[obj.(*testObject).val] = &atomic.Int64{}
		if _, ok := wants[obj.(*testObject).val]; !ok {
			weight := weights[i]
			want := int64(float64(batchSize*workerCount) * weight)
			wants[obj.(*testObject).val] = &wantObject{
				want: want,
				diff: allowedMistakeCount,
			}
			t.Logf("obj=%d, want=%d, weight=%f, total=%d", obj.(*testObject).val, want, weight, batchSize*workerCount)
		}
	}
	hashData := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("10.99.88.77"),
		net.ParseIP("223.5.5.5"),
		net.ParseIP("114.114.114.114"),
		net.ParseIP("1.1.1.1"),
	}

	// exec
	for x := 0; x < workerCount; x++ {
		wg.Add(1)
		// sleep x*10ms to avoid all worker start at the same time
		time.Sleep(time.Millisecond * time.Duration(x*10))
		go func(x int, objs []any) {
			start := time.Now().UnixNano()
			for i := 0; i < batchSize; i++ {
				// arg only used in st == StrategyHash
				arg := hashData[i%len(hashData)].String()
				obj, err := balancer.Pick(objs, arg)
				if err != nil || obj == nil {
					t.Error(err)
					break
				}
				// t.Logf("worker %d-%d pick %d", x, i, obj.(*testObject).val)
				res[obj.(*testObject).val].Add(1)
				cnt.Add(1)
			}
			wg.Done()
			timeCost.Add(time.Now().UnixNano() - start)
		}(x, objs)
	}
	wg.Wait()
	if int64(batchSize*workerCount) != cnt.Load() {
		t.Fatalf("execute count mismatch. expected %d, got %d", batchSize*workerCount, cnt.Load())
	}
	t.Logf("execute finished. count=%d, time_cost=%dns, avg_cost=%fns", cnt.Load(), timeCost.Load(), float64(timeCost.Load()/cnt.Load()))

	// check
	errCount := 0
	for v, a := range res {
		if !wants[v].judge(a.Load()) {
			t.Logf("[error] strategy=%s, v=%d, expected %d, got %d, diff=%d, allowed_diff=%d", st, v, wants[v].want, a.Load(), a.Load()-wants[v].want, wants[v].diff)
			errCount++
			continue
		}
		t.Logf("[ok] strategy=%s, v=%d, expected %d, got %d, diff=%d, allowed_diff=%d", st, v, wants[v].want, a.Load(), a.Load()-wants[v].want, wants[v].diff)
	}

	if errCount > allowedError {
		t.Fatalf("[error] errorCount > allowedError:  errorCount= %d/%d | allowed=%d", errCount, len(res), allowedError)
	} else {
		t.Logf("[ok] errorCount: %d/%d | allowed=%d", errCount, len(res), allowedError)
	}

}

func TestRandConcurrency(t *testing.T) {
	var slice []int
	for i := 0; i < 100; i++ {
		slice = append(slice, i)
	}
	// r := rand.New(rand.NewSource(time.Now().UnixMilli()))
	r := newRand()

	for x := 0; x < 100000; x++ {
		if len(slice) != 100 {
			t.Logf("not equal")
		}
		// panics if len(slice) == 0.
		r.IntN(len(slice))
	}

	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			for x := 0; x < 100000; x++ {
				if len(slice) != 100 {
					t.Logf("not equal")
				}
				// panics if len(slice) == 0.
				r.IntN(len(slice))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestIpHashBalancer_Calc(t *testing.T) {
	b := New(StrategyHash)
	hashData := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("10.99.88.77"),
		net.ParseIP("223.5.5.5"),
		net.ParseIP("114.114.114.114"),
		net.ParseIP("1.1.1.1"),
		net.ParseIP("2001:db8::68"),
	}
	hb := b.(*wrapper).b.(*hashBalancer)
	for _, datum := range hashData {
		t.Logf("s=%s, v=%d", datum.String(), hb.calc(datum.String()))
	}
}

/**
Benchmarks
*/

var (
	benchmarkObjectSize = 100
	benchmarkWeight     = 10
)

func BenchmarkHashBalancer_Pick(b *testing.B) {
	balancer := New(StrategyHash)
	hashData := []net.IP{
		net.ParseIP("192.168.1.1"),
		net.ParseIP("10.99.88.77"),
		net.ParseIP("223.5.5.5"),
		net.ParseIP("114.114.114.114"),
		net.ParseIP("1.1.1.1"),
		net.ParseIP("2001:db8::68"),
	}
	var objs []any
	// test same weight
	for i := 0; i < benchmarkObjectSize; i++ {
		objs = append(objs, newTestObject(int64(i), int64(benchmarkWeight)))
	}
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Pick(objs, hashData[i%len(hashData)].String())
	}
}

func BenchmarkRandomBalancer_Pick(b *testing.B) {
	balancer := New(StrategyRandom)
	var objs []any
	// test same weight
	for i := 0; i < benchmarkObjectSize; i++ {
		objs = append(objs, newTestObject(int64(i), int64(benchmarkWeight)))
	}
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Pick(objs, nil)
	}
}

func BenchmarkWeightedRandom_Pick(b *testing.B) {
	balancer := New(StrategyWeightedRandom)
	var objs []any
	// test same weight
	for i := 0; i < benchmarkObjectSize; i++ {
		objs = append(objs, newTestObject(int64(i), int64(benchmarkWeight)))
	}
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Pick(objs, nil)
	}
}

func BenchmarkWeightedRoundRobin_Pick(b *testing.B) {
	balancer := New(StrategyWeightedRoundRobin)
	var objs []any
	// test same weight
	for i := 0; i < benchmarkObjectSize; i++ {
		objs = append(objs, newTestObject(int64(i), int64(benchmarkWeight)))
	}
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Pick(objs, nil)
	}
}

func BenchmarkRoundRobin_Pick(b *testing.B) {
	balancer := New(StrategyRoundRobin)
	var objs []any
	// test same weight
	for i := 0; i < benchmarkObjectSize; i++ {
		objs = append(objs, newTestObject(int64(i), int64(benchmarkWeight)))
	}
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Pick(objs, nil)
	}
}

func BenchmarkLessLoad_Pick(b *testing.B) {
	balancer := New(StrategyLessLoad)
	var objs []any
	// test same weight
	for i := 0; i < benchmarkObjectSize; i++ {
		objs = append(objs, newTestObject(int64(i), int64(benchmarkWeight)))
	}
	for i := 0; i < b.N; i++ {
		_, _ = balancer.Pick(objs, nil)
	}
}
