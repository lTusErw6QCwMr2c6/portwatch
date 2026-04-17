package summary

import (
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/logger"
)

// Emitter periodically logs a summary report.
type Emitter struct {
	builder  *Builder
	interval time.Duration
	log      *logger.Logger
	stop     chan struct{}
}

// NewEmitter creates an Emitter that logs summaries at the given interval.
func NewEmitter(b *Builder, interval time.Duration, log *logger.Logger) *Emitter {
	return &Emitter{
		builder:  b,
		interval: interval,
		log:      log,
		stop:     make(chan struct{}),
	}
}

// Start begins emitting summaries in a background goroutine.
func (e *Emitter) Start(events func() []alert.Event) {
	go func() {
		ticker := time.NewTicker(e.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r := e.builder.Build(events())
				e.log.Info(r.Format())
			case <-e.stop:
				return
			}
		}
	}()
}

// Stop shuts down the emitter.
func (e *Emitter) Stop() {
	close(e.stop)
}
