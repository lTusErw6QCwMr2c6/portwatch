// Package watchdog detects when the scan cycle stalls or takes too long.
package watchdog

import (
	"fmt"
	"sync"
	"time"
)

// Logger is a minimal logging interface.
type Logger interface {
	Warn(msg string)
}

// Watchdog monitors heartbeats and fires an alert if a deadline is missed.
type Watchdog struct {
	mu       sync.Mutex
	timeout  time.Duration
	timer    *time.Timer
	log      Logger
	stalled  bool
	stopCh   chan struct{}
}

// New creates a Watchdog that fires a warning after timeout with no heartbeat.
func New(timeout time.Duration, log Logger) *Watchdog {
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	w := &Watchdog{
		timeout: timeout,
		log:     log,
		stopCh:  make(chan struct{}),
	}
	w.timer = time.AfterFunc(timeout, w.onStall)
	return w
}

// Heartbeat resets the watchdog timer, signalling the cycle is alive.
func (w *Watchdog) Heartbeat() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.stalled = false
	w.timer.Reset(w.timeout)
}

// Stalled reports whether the watchdog has fired without a subsequent heartbeat.
func (w *Watchdog) Stalled() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.stalled
}

// Stop disables the watchdog.
func (w *Watchdog) Stop() {
	w.timer.Stop()
}

func (w *Watchdog) onStall() {
	w.mu.Lock()
	w.stalled = true
	w.mu.Unlock()
	w.log.Warn(fmt.Sprintf("watchdog: no heartbeat received within %s — scan cycle may be stalled", w.timeout))
}
