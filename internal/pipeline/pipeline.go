package pipeline

import (
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/scanner"
)

// Pipeline processes a port diff through filter, rate-limit, and alert stages.
type Pipeline struct {
	filter    *filter.Filter
	ratelimit *ratelimit.RateLimiter
	threshold alert.Threshold
}

// Config holds pipeline stage configuration.
type Config struct {
	Filter    *filter.Filter
	RateLimit *ratelimit.RateLimiter
	Threshold alert.Threshold
}

// New creates a Pipeline from the given Config.
func New(cfg Config) *Pipeline {
	if cfg.Threshold == (alert.Threshold{}) {
		cfg.Threshold = alert.DefaultThreshold()
	}
	return &Pipeline{
		filter:    cfg.Filter,
		ratelimit: cfg.RateLimit,
		threshold: cfg.Threshold,
	}
}

// Run applies all pipeline stages to a diff and returns alert events.
func (p *Pipeline) Run(diff scanner.Diff) []alert.Event {
	if p.filter != nil {
		diff = p.filter.Apply(diff)
	}

	if !alert.Evaluate(diff, p.threshold) {
		return nil
	}

	events := alert.BuildEvents(diff)

	if p.ratelimit == nil {
		return events
	}

	allowed := events[:0]
	for _, e := range events {
		if p.ratelimit.Allow(e.Key()) {
			allowed = append(allowed, e)
		}
	}
	return allowed
}
