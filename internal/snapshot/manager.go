package snapshot

import (
	"fmt"
	"log"

	"github.com/user/portwatch/internal/scanner"
)

// Manager coordinates scanning and persisting snapshots between cycles.
type Manager struct {
	store    *Store
	logger   *log.Logger
	portStart int
	portEnd   int
}

// NewManager creates a Manager that scans [portStart, portEnd] and persists
// state to the given store.
func NewManager(store *Store, logger *log.Logger, portStart, portEnd int) *Manager {
	return &Manager{
		store:     store,
		logger:    logger,
		portStart: portStart,
		portEnd:   portEnd,
	}
}

// Cycle performs one scan, computes the diff against the previous snapshot,
// persists the new snapshot, and returns the diff.
func (m *Manager) Cycle() (scanner.DiffResult, error) {
	prev, err := m.store.Load()
	if err != nil {
		return scanner.DiffResult{}, fmt.Errorf("load snapshot: %w", err)
	}

	curr, err := scanner.ScanPorts(m.portStart, m.portEnd)
	if err != nil {
		return scanner.DiffResult{}, fmt.Errorf("scan ports: %w", err)
	}

	diff := scanner.Diff(prev.Ports, curr)

	newSnap := New(curr)
	if err := m.store.Save(newSnap); err != nil {
		m.logger.Printf("[WARN] failed to persist snapshot: %v", err)
	}

	return diff, nil
}

// Reset clears the persisted snapshot so the next Cycle starts fresh.
func (m *Manager) Reset() error {
	empty := New(make(map[scanner.PortKey]bool))
	return m.store.Save(empty)
}
