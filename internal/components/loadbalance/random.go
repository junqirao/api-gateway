package loadbalance

import (
	"math/rand"
	"sort"
	"time"

	"github.com/gogf/gf/v2/net/ghttp"
)

const (
	intSize = 32 << (^uint(0) >> 63) // cf. strconv.IntSize
	maxInt  = 1<<(intSize-1) - 1
)

// random weighted random selection ref: https://github.com/mroth/weightedrand
type random struct {
	rand *rand.Rand
}

func newRandom() Balancer {
	return &random{rand: rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func (ra *random) Selector(_ *ghttp.Request, ups []Weighted) (ref int, ok bool) {
	if len(ups) == 0 {
		return
	}

	sort.Slice(ups, func(i, j int) bool {
		return ups[i].Weight() < ups[j].Weight()
	})
	totals := make([]int, len(ups))
	runningTotal := 0
	for i, c := range ups {
		weight := int(c.Weight())
		if (maxInt - runningTotal) <= weight {
			return
		}
		runningTotal += weight
		totals[i] = runningTotal
	}
	if runningTotal < 1 {
		return
	}

	ref = ups[ra.search(totals, ra.rand.Intn(runningTotal)+1)].Ref()
	ok = true
	return
}

func (ra *random) search(a []int, x int) int {
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
