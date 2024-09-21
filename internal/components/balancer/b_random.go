package balancer

type randomBalancer struct {
	r *rand
}

func newRandom() Balancer {
	return &randomBalancer{
		r: newRand(),
	}
}

func (r *randomBalancer) Pick(objects []any, args ...any) (o any, err error) {
	if len(objects) == 0 {
		err = ErrNoObject
		return
	}
	return objects[r.r.IntN(len(objects))], nil
}
