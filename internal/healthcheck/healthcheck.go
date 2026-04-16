package healthcheck

import (
	"fmt"
	"sync"
	"time"
)

// Status represents the health state of a component.
type Status string

const (
	StatusOK      Status = "ok"
	StatusDegraded Status = "degraded"
	StatusDown    Status = "down"
)

// Check holds the result of a single health check.
type Check struct {
	Name      string
	Status    Status
	Message   string
	CheckedAt time.Time
}

func (c Check) String() string {
	return fmt.Sprintf("[%s] %s: %s", c.Status, c.Name, c.Message)
}

// Monitor tracks health checks for named components.
type Monitor struct {
	mu     sync.RWMutex
	checks map[string]Check
}

// New returns an initialised Monitor.
func New() *Monitor {
	return &Monitor{checks: make(map[string]Check)}
}

// Record stores a health check result for the given component name.
func (m *Monitor) Record(name string, status Status, message string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.checks[name] = Check{
		Name:      name,
		Status:    status,
		Message:   message,
		CheckedAt: time.Now(),
	}
}

// All returns a copy of all recorded checks.
func (m *Monitor) All() []Check {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]Check, 0, len(m.checks))
	for _, c := range m.checks {
		out = append(out, c)
	}
	return out
}

// Overall returns StatusOK if all checks are OK, StatusDegraded if any are
// degraded, and StatusDown if any are down.
func (m *Monitor) Overall() Status {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := StatusOK
	for _, c := range m.checks {
		if c.Status == StatusDown {
			return StatusDown
		}
		if c.Status == StatusDegraded {
			result = StatusDegraded
		}
	}
	return result
}
