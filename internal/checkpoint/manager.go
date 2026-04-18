package checkpoint

import (
	"time"
)

// Manager wraps Store and provides a higher-level API used by the watcher loop.
type Manager struct {
	store    *Store
	interval time.Duration
	last     time.Time
}

// NewManager returns a Manager that auto-saves at most every interval.
func NewManager(store *Store, interval time.Duration) *Manager {
	return &Manager{store: store, interval: interval}
}

// MaybeSave persists state only when the flush interval has elapsed.
func (m *Manager) MaybeSave(state State) error {
	if time.Since(m.last) < m.interval {
		return nil
	}
	if err := m.store.Save(state); err != nil {
		return err
	}
	m.last = time.Now()
	return nil
}

// ForceSave always persists state regardless of the interval.
func (m *Manager) ForceSave(state State) error {
	if err := m.store.Save(state); err != nil {
		return err
	}
	m.last = time.Now()
	return nil
}

// Load delegates to the underlying Store.
func (m *Manager) Load() (State, error) {
	return m.store.Load()
}
