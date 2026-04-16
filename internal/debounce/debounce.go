package debounce

import (
	"sync"
	"time"
)

// Debouncer suppresses repeated triggers for the same key within a given window.
// Unlike a rate limiter, it resets the timer on each new trigger, firing only
// after the activity has settled.
type Debouncer struct {
	mu      sync.Mutex
	window  time.Duration
	timers  map[string]*time.Timer
	callback func(key string)
}

// New creates a Debouncer that calls fn(key) after no further triggers for
// that key have occurred within window.
func New(window time.Duration, fn func(key string)) *Debouncer {
	return &Debouncer{
		window:   window,
		timers:   make(map[string]*time.Timer),
		callback: fn,
	}
}

// Trigger schedules fn(key) to be called after the debounce window. If Trigger
// is called again for the same key before the window elapses, the timer resets.
func (d *Debouncer) Trigger(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
	}

	d.timers[key] = time.AfterFunc(d.window, func() {
		d.mu.Lock()
		delete(d.timers, key)
		d.mu.Unlock()
		d.callback(key)
	})
}

// Cancel stops any pending timer for key without invoking the callback.
func (d *Debouncer) Cancel(key string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if t, ok := d.timers[key]; ok {
		t.Stop()
		delete(d.timers, key)
	}
}

// Pending returns the number of keys currently waiting to fire.
func (d *Debouncer) Pending() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return len(d.timers)
}
