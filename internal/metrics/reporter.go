package metrics

import (
	"fmt"
	"io"
	"time"
)

// Reporter periodically writes a metrics summary to a writer.
type Reporter struct {
	collector *Collector
	out       io.Writer
	interval  time.Duration
	stop      chan struct{}
}

// NewReporter creates a Reporter that prints to out every interval.
func NewReporter(c *Collector, out io.Writer, interval time.Duration) *Reporter {
	return &Reporter{
		collector: c,
		out:       out,
		interval:  interval,
		stop:      make(chan struct{}),
	}
}

// Start begins the reporting loop in a background goroutine.
func (r *Reporter) Start() {
	go func() {
		ticker := time.NewTicker(r.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				r.write()
			case <-r.stop:
				return
			}
		}
	}()
}

// Stop halts the reporting loop.
func (r *Reporter) Stop() {
	close(r.stop)
}

func (r *Reporter) write() {
	s := r.collector.Read()
	fmt.Fprintf(r.out, "[metrics] cycles=%d opened=%d closed=%d alerts=%d last_cycle=%s\n",
		s.Cycles, s.Opened, s.Closed, s.Alerts, formatTime(s.LastCycleAt))
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "never"
	}
	return t.Format(time.RFC3339)
}
