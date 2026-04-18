package sampling

import "sync"

// Strategy controls how events are sampled before dispatch.
type Strategy int

const (
	StrategyNone  Strategy = iota // pass all events through
	StrategyEvery                 // pass every Nth event
	StrategyRandom                // pass with probability Rate
)

// Config holds sampler configuration.
type Config struct {
	Strategy Strategy
	Rate     float64 // 0.0–1.0 for Random; N for Every
}

// Sampler decides whether an event should be forwarded.
type Sampler struct {
	mu      sync.Mutex
	cfg     Config
	counter map[string]int
	randf   func() float64
}

// New returns a Sampler with the given config.
func New(cfg Config, randf func() float64) *Sampler {
	if randf == nil {
		randf = defaultRand
	}
	return &Sampler{cfg: cfg, counter: make(map[string]int), randf: randf}
}

// Allow returns true if the event identified by key should pass through.
func (s *Sampler) Allow(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	switch s.cfg.Strategy {
	case StrategyEvery:
		n := int(s.cfg.Rate)
		if n <= 1 {
			return true
		}
		s.counter[key]++
		if s.counter[key] >= n {
			s.counter[key] = 0
			return true
		}
		return false
	case StrategyRandom:
		return s.randf() < s.cfg.Rate
	default:
		return true
	}
}

// Reset clears all counters.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter = make(map[string]int)
}

var defaultRand func() float64

func init() {
	import_math_rand()
}
