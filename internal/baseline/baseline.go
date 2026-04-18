// Package baseline records a trusted set of ports and flags deviations.
package baseline

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/user/portwatch/internal/scanner"
)

// Baseline holds a snapshot of ports considered "normal".
type Baseline struct {
	mu      sync.RWMutex
	ports   map[string]scanner.Port
	SavedAt time.Time `json:"saved_at"`
}

// New returns an empty Baseline.
func New() *Baseline {
	return &Baseline{ports: make(map[string]scanner.Port)}
}

// Set replaces the current baseline with the given ports.
func (b *Baseline) Set(ports []scanner.Port) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.ports = make(map[string]scanner.Port, len(ports))
	for _, p := range ports {
		b.ports[key(p)] = p
	}
	b.SavedAt = time.Now()
}

// Deviations returns ports in current that are absent from the baseline.
func (b *Baseline) Deviations(current []scanner.Port) []scanner.Port {
	b.mu.RLock()
	defer b.mu.RUnlock()
	var out []scanner.Port
	for _, p := range current {
		if _, ok := b.ports[key(p)]; !ok {
			out = append(out, p)
		}
	}
	return out
}

// Save persists the baseline to a JSON file.
func (b *Baseline) Save(path string) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	type file struct {
		SavedAt time.Time      `json:"saved_at"`
		Ports   []scanner.Port `json:"ports"`
	}
	ports := make([]scanner.Port, 0, len(b.ports))
	for _, p := range b.ports {
		ports = append(ports, p)
	}
	data, err := json.MarshalIndent(file{SavedAt: b.SavedAt, Ports: ports}, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Load reads a baseline from a JSON file.
func Load(path string) (*Baseline, error) {
	type file struct {
		SavedAt time.Time      `json:"saved_at"`
		Ports   []scanner.Port `json:"ports"`
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var f file
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	b := New()
	b.Set(f.Ports)
	b.SavedAt = f.SavedAt
	return b, nil
}

func key(p scanner.Port) string {
	return p.Protocol + ":" + p.Address
}
