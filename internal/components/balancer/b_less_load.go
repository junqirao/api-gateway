package balancer

type lessLoadBalancer struct {
}

func newLessLoad() Balancer {
	return &lessLoadBalancer{}
}

func (l lessLoadBalancer) Pick(objects []any, _ ...any) (o any, err error) {
	var minVal int64 = 0
	for _, object := range objects {
		if measure, ok := object.(Measurable); ok {
			if o == nil {
				o = object
				minVal = measure.Load()
				continue
			}
			if cur := measure.Load(); cur < minVal {
				o = object
				minVal = cur
			}
		}
	}
	return
}
