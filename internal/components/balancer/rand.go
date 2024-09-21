package balancer

import (
	math "math/rand"
	"sync"
	"time"
)

// rand is a multi goroutine safe random number generator
type rand struct {
	r  *math.Rand
	mu sync.Mutex
}

func newRand() *rand {
	return &rand{r: math.New(math.NewSource(time.Now().UnixNano()))}
}

func (r *rand) IntN(n int) int {
	r.mu.Lock()
	res := r.r.Intn(n)
	r.mu.Unlock()
	return res
}
