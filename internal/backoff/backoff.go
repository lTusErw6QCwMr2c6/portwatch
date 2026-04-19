package backoff

import (
	"math"
	"time"
)

// Strategy defines how delays are calculated between retries.
type Strategy int

const (
	Fixed Strategy = iota
	Exponential
)

// Config holds backoff configuration.
type Config struct {
	Strategy   Strategy
	BaseDelay  time.Duration
	MaxDelay   time.Duration
	Multiplier float64
	MaxRetries int
}

// DefaultConfig returns a sensible exponential backoff config.
func DefaultConfig() Config {
	return Config{
		Strategy:   Exponential,
		BaseDelay:  500 * time.Millisecond,
		MaxDelay:   30 * time.Second,
		Multiplier: 2.0,
		MaxRetries: 5,
	}
}

// Backoff tracks retry state for a single operation.
type Backoff struct {
	cfg     Config
	attempt int
}

// New creates a new Backoff using the given config.
func New(cfg Config) *Backoff {
	return &Backoff{cfg: cfg}
}

// Next returns the delay for the current attempt and advances the counter.
// Returns false when max retries have been exhausted.
func (b *Backoff) Next() (time.Duration, bool) {
	if b.cfg.MaxRetries > 0 && b.attempt >= b.cfg.MaxRetries {
		return 0, false
	}
	var delay time.Duration
	switch b.cfg.Strategy {
	case Exponential:
		f := math.Pow(b.cfg.Multiplier, float64(b.attempt))
		delay = time.Duration(float64(b.cfg.BaseDelay) * f)
	default:
		delay = b.cfg.BaseDelay
	}
	if delay > b.cfg.MaxDelay {
		delay = b.cfg.MaxDelay
	}
	b.attempt++
	return delay, true
}

// Reset resets the attempt counter.
func (b *Backoff) Reset() {
	b.attempt = 0
}

// Attempt returns the current attempt number.
func (b *Backoff) Attempt() int {
	return b.attempt
}
