package balancer

import (
	"context"
	"errors"

	"github.com/gogf/gf/v2/frame/g"
)

const (
	maxAttempts = 5
)

type (
	Filter  func(o any) bool
	Filters []Filter
	wrapper struct {
		s         Strategy
		b, backup Balancer
	}
)

func (w *wrapper) Pick(objects []any, args ...any) (o any, err error) {
	// nothing to pick
	if len(objects) == 0 {
		err = ErrNoObject
		return
	}
	// pick first object if only one
	if len(objects) == 1 {
		o = objects[0]
		return
	}

	// try pick
	var (
		arg     = ""
		filters Filters
	)

	if len(args) > 0 {
		arg = args[0].(string)
	}
	if len(args) > 1 {
		if fs, ok := args[1].(Filters); ok {
			filters = fs
		}
	}

	// max attempts = min(maxAttempts,len(objects))
	for i := 0; i < min(maxAttempts, len(objects)); i++ {
		o, err = w.b.Pick(objects, arg)
		switch {
		case err == nil:
		case w.backup != nil,
			errors.Is(err, ErrUnWeighable),
			errors.Is(err, ErrMissingArgs):
			// retry backup balancer
			g.Log().Infof(context.TODO(), "[balancer] use backup balancer: %s", err)
			o, err = w.backup.Pick(objects, arg)
		}
		if filters != nil && !filters.Check(o) {
			if w.s == StrategyHash {
				// retry may get the same object
				return
			}
			continue
		}
		return
	}

	return
}

func (f Filters) Check(o any) bool {
	for _, filter := range f {
		if !filter(o) {
			return false
		}
	}
	return true
}
