package balancer

import (
	"sort"
)

type weightedRandom struct {
	r *rand
}

func newWeightedRandom() Balancer {
	return &weightedRandom{
		r: newRand(),
	}
}

func (w *weightedRandom) Pick(objects []any, _ ...any) (o any, err error) {
	sort.Slice(objects, func(i, j int) bool {
		return objects[i].(Weighable).Weight() < objects[j].(Weighable).Weight()
	})
	totals := make([]int, len(objects))
	runningTotal := 0
	for i, c := range objects {
		weight := int(c.(Weighable).Weight())
		if (maxInt - runningTotal) <= weight {
			return
		}
		runningTotal += weight
		totals[i] = runningTotal
	}
	if runningTotal < 1 {
		return
	}

	o = objects[w.search(totals, w.r.IntN(runningTotal)+1)]
	return
}

func (w *weightedRandom) search(a []int, x int) int {
	// Possible further future optimization for search via SIMD if we want
	// to write some Go assembly code: http://0x80.pl/articles/simd-search.html
	i, j := 0, len(a)
	for i < j {
		h := int(uint(i+j) >> 1) // avoid overflow when computing h
		if a[h] < x {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
