package jitter

import (
	"math/rand"
	"time"
)

// Jitter adds randomised variance to durations to prevent thundering-herd
// effects when multiple goroutines wake at the same scheduled interval.
type Jitter struct {
	factor float64
	rng    *rand.Rand
}

// New returns a Jitter with the given factor (0.0–1.0).
// A factor of 0.2 means ±20 % variance around the base duration.
func New(factor float64) *Jitter {
	if factor < 0 {
		factor = 0
	}
	if factor > 1 {
		factor = 1
	}
	return &Jitter{
		factor: factor,
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Apply returns base ± (factor * base * random[0,1)).
func (j *Jitter) Apply(base time.Duration) time.Duration {
	if j.factor == 0 || base <= 0 {
		return base
	}
	delta := float64(base) * j.factor * j.rng.Float64()
	if j.rng.Intn(2) == 0 {
		return base + time.Duration(delta)
	}
	return base - time.Duration(delta)
}

// ApplyPositive is like Apply but guarantees the result is >= 1ns.
func (j *Jitter) ApplyPositive(base time.Duration) time.Duration {
	d := j.Apply(base)
	if d < time.Nanosecond {
		return time.Nanosecond
	}
	return d
}

// Reset re-seeds the internal RNG.
func (j *Jitter) Reset() {
	j.rng = rand.New(rand.NewSource(time.Now().UnixNano()))
}
