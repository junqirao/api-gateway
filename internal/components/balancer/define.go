package balancer

import (
	"errors"
)

// Strategy define
const (
	StrategyRoundRobin         Strategy = "round-robin"
	StrategyRandom             Strategy = "random"
	StrategyWeightedRoundRobin Strategy = "weighted-round-robin"
	StrategyWeightedRandom     Strategy = "weighted-random"
	StrategyLessLoad           Strategy = "less-load"
	StrategyHash               Strategy = "hash"
)

// error define
var (
	ErrNoObject = errors.New("no object")

	ErrUnWeighable = errors.New("not weighable object, check implementation")
	ErrMissingArgs = errors.New("missing args")
)

const (
	intSize = 32 << (^uint(0) >> 63) // cf. strconv.IntSize
	maxInt  = 1<<(intSize-1) - 1
)
