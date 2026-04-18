// Package checkpoint persists the last-known scan state to disk so that
// portwatch can resume gracefully after a restart without flooding alerts.
package checkpoint

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// State holds the data written to / read from the checkpoint file.
type State struct {
	Ports     []string  `json:"ports"`
	SavedAt   time.Time `json:"saved_at"`
	CycleCount int64    `json:"cycle_count"`
}

// Store manages reading and writing checkpoint state.
type Store struct {
	mu   sync.Mutex
	path string
}

// New returns a Store that persists state to path.
func New(path string) *Store {
	return &Store{path: path}
}

// Save writes s to disk atomically.
func (s *Store) Save(state State) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	state.SavedAt = time.Now()
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

// Load reads the persisted state. Returns a zero State and no error when the
// file does not exist yet.
func (s *Store) Load() (State, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.path)
	if os.IsNotExist(err) {
		return State{}, nil
	}
	if err != nil {
		return State{}, err
	}
	var st State
	if err := json.Unmarshal(data, &st); err != nil {
		return State{}, err
	}
	return st, nil
}

// Remove deletes the checkpoint file if it exists.
func (s *Store) Remove() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := os.Remove(s.path)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}
