package balancer

import (
	"hash/maphash"
	"sync"
)

type hashBalancer struct {
	pool *sync.Pool
}

func newHash() Balancer {
	return &hashBalancer{
		pool: &sync.Pool{
			New: func() any {
				return new(maphash.Hash)
			},
		},
	}
}

func (b *hashBalancer) Pick(objects []any, args ...any) (o any, err error) {
	if len(args) == 0 {
		err = ErrMissingArgs
		return
	}
	s := args[0].(string)
	if s == "" {
		err = ErrMissingArgs
		return
	}
	o = objects[b.calc(s)%len(objects)]
	return
}

func (b *hashBalancer) calc(s string) int {
	hash := b.pool.Get().(*maphash.Hash)
	// renew seed make sure result is distributed
	// hash.SetSeed(maphash.MakeSeed())
	_, _ = hash.WriteString(s)
	res := int(hash.Sum64() % maxInt)
	hash.Reset()
	b.pool.Put(hash)
	return res
}
