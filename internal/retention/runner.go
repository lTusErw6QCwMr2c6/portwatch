package retention

import (
	"time"
)

// Logger is a minimal logging interface.
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Runner periodically applies a retention policy.
type Runner struct {
	policy   *Policy
	interval time.Duration
	log      Logger
	stop     chan struct{}
}

// NewRunner creates a Runner that applies the policy on the given interval.
func NewRunner(policy *Policy, interval time.Duration, log Logger) *Runner {
	return &Runner{
		policy:   policy,
		interval: interval,
		log:      log,
		stop:     make(chan struct{}),
	}
}

// Start begins the periodic retention loop in a goroutine.
func (r *Runner) Start() {
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				removed, err := r.policy.Apply()
				if err != nil {
					r.log.Errorf("retention: apply failed: %v", err)
					continue
				}
				if len(removed) > 0 {
					r.log.Infof("retention: removed %d file(s)", len(removed))
				}
			case <-r.stop:
				return
			}
		}
	}()
}

// Stop halts the runner.
func (r *Runner) Stop() {
	close(r.stop)
}
