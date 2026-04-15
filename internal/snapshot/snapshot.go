package snapshot

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Snapshot holds a point-in-time capture of open ports.
type Snapshot struct {
	Timestamp time.Time              `json:"timestamp"`
	Ports     map[scanner.PortKey]bool `json:"ports"`
}

// Store manages reading and writing snapshots to disk.
type Store struct {
	mu   sync.RWMutex
	path string
}

// NewStore creates a new Store backed by the given file path.
func NewStore(path string) *Store {
	return &Store{path: path}
}

// Save writes the current snapshot to disk as JSON.
func (s *Store) Save(snap Snapshot) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	f, err := os.Create(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(snap)
}

// Load reads the last snapshot from disk. Returns an empty snapshot if the
// file does not exist yet.
func (s *Store) Load() (Snapshot, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var snap Snapshot
	f, err := os.Open(s.path)
	if os.IsNotExist(err) {
		return Snapshot{Ports: make(map[scanner.PortKey]bool)}, nil
	}
	if err != nil {
		return snap, err
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&snap); err != nil {
		return snap, err
	}
	if snap.Ports == nil {
		snap.Ports = make(map[scanner.PortKey]bool)
	}
	return snap, nil
}

// New creates a fresh snapshot from a set of port keys.
func New(ports map[scanner.PortKey]bool) Snapshot {
	return Snapshot{
		Timestamp: time.Now(),
		Ports:     ports,
	}
}
