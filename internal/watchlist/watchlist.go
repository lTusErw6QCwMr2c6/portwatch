// Package watchlist provides port-specific watch rules that trigger
// elevated alerting when a monitored port changes state.
package watchlist

import (
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Entry describes a single watched port rule.
type Entry struct {
	Port     int
	Protocol string
	Label    string
}

// Watcher holds a set of watched port entries and matches them against events.
type Watcher struct {
	mu      sync.RWMutex
	entries map[string]Entry // key: "proto:port"
}

// New returns an empty Watcher.
func New() *Watcher {
	return &Watcher{entries: make(map[string]Entry)}
}

// Add registers a port entry in the watchlist.
func (w *Watcher) Add(e Entry) {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.entries[key(e.Protocol, e.Port)] = e
}

// Remove deletes a port entry from the watchlist.
func (w *Watcher) Remove(protocol string, port int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.entries, key(protocol, port))
}

// Match returns the Entry and true if the event's port is on the watchlist.
func (w *Watcher) Match(ev alert.Event) (Entry, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	e, ok := w.entries[key(ev.Port.Protocol, ev.Port.Port)]
	return e, ok
}

// All returns a snapshot of all current watchlist entries.
func (w *Watcher) All() []Entry {
	w.mu.RLock()
	defer w.mu.RUnlock()
	out := make([]Entry, 0, len(w.entries))
	for _, e := range w.entries {
		out = append(out, e)
	}
	return out
}

func key(protocol string, port int) string {
	return fmt.Sprintf("%s:%d", protocol, port)
}
