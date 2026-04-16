package supervisor

import (
	"context"
	"time"
)

// Status represents the health state of a supervised component.
type Status int

const (
	StatusOK Status = iota
	StatusDegraded
	StatusFailed
)

// Component is anything that can be started and stopped.
type Component interface {
	Name() string
	Start(ctx context.Context) error
	Stop() error
}

// Entry tracks a supervised component and its last known status.
type Entry struct {
	Component Component
	Status    Status
	LastError error
	Restarts  int
	LastSeen  time.Time
}

// Supervisor manages lifecycle and restart logic for components.
type Supervisor struct {
	entries    map[string]*Entry
	maxRetries int
	backoff    time.Duration
}

// New creates a Supervisor with the given retry limit and backoff duration.
func New(maxRetries int, backoff time.Duration) *Supervisor {
	return &Supervisor{
		entries:    make(map[string]*Entry),
		maxRetries: maxRetries,
		backoff:    backoff,
	}
}

// Register adds a component to the supervisor.
func (s *Supervisor) Register(c Component) {
	s.entries[c.Name()] = &Entry{
		Component: c,
		Status:    StatusOK,
		LastSeen:  time.Now(),
	}
}

// Run starts all registered components and supervises them.
func (s *Supervisor) Run(ctx context.Context) {
	for _, e := range s.entries {
		go s.supervise(ctx, e)
	}
}

func (s *Supervisor) supervise(ctx context.Context, e *Entry) {
	for {
		err := e.Component.Start(ctx)
		e.LastSeen = time.Now()
		if err == nil || ctx.Err() != nil {
			e.Status = StatusOK
			return
		}
		e.LastError = err
		e.Restarts++
		if e.Restarts > s.maxRetries {
			e.Status = StatusFailed
			return
		}
		e.Status = StatusDegraded
		select {
		case <-ctx.Done():
			return
		case <-time.After(s.backoff):
		}
	}
}

// Entries returns a snapshot of all tracked entries.
func (s *Supervisor) Entries() map[string]*Entry {
	out := make(map[string]*Entry, len(s.entries))
	for k, v := range s.entries {
		copy := *v
		out[k] = &copy
	}
	return out
}
