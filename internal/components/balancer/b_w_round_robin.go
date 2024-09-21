package balancer

type weightedRoundRobin struct {
}

func newWeightedRoundRobin() Balancer {
	return &weightedRoundRobin{}
}

func (wrr *weightedRoundRobin) Pick(objects []any, _ ...any) (o any, err error) {
	var (
		total int64 = 0
		best  Weighable
		idx   int
	)

	for i := 0; i < len(objects); i++ {
		w, ok := objects[i].(Weighable)
		if w == nil || !ok {
			continue
		}
		w.Current().Add(w.Effective().Load())
		total += w.Effective().Load()
		if w.Effective().Load() < w.Weight() {
			w.Effective().Add(1)
		}
		if best == nil || w.Current().Load() > best.Current().Load() {
			best = w
			idx = i
		}
	}
	if best == nil {
		err = ErrUnWeighable
		return
	}
	best.Current().Add(-total)
	o = objects[idx]
	return
}
