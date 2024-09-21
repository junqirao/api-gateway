package balancer

import (
	"errors"
	"log"
)

type (
	wrapper struct {
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
	o, err = w.b.Pick(objects, args...)
	switch {
	case err == nil:
		return
	case w.backup != nil,
		errors.Is(err, ErrUnWeighable),
		errors.Is(err, ErrMissingArgs):
		// retry backup balancer
		log.Default().Printf("[balancer] use backup balancer: %s", err)
		o, err = w.backup.Pick(objects, args...)
	}
	return
}
