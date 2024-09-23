package balancer

type (
	// Strategy of balancer
	Strategy string

	// Balancer interface
	Balancer interface {
		// Pick an object
		// when objects size=1, the strategy will be ignored
		// when objects size=0, return ErrNoObject
		// args only used in StrategyHash, accept value type of string
		Pick(objects []any, args ...any) (o any, err error)
	}
)

// New Balancer
func New(strategy Strategy) (b Balancer) {
	switch strategy {
	case StrategyRoundRobin:
		b = newRoundRobin()
	case StrategyRandom:
		b = newRandom()
	case StrategyWeightedRoundRobin:
		b = newWeightedRoundRobin()
	case StrategyWeightedRandom:
		b = newWeightedRandom()
	case StrategyLessLoad:
		b = newLessLoad()
	case StrategyHash:
		b = newHash()
	default:
		// default strategy round-robin
		b = newRoundRobin()
	}

	return &wrapper{b: b, backup: newRoundRobin(), s: strategy}
}
