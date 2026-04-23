package eventtransform

import (
	"errors"
	"fmt"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// TransformFunc is a function that transforms an event in place.
type TransformFunc func(e *alert.Event) (*alert.Event, error)

// Transformer applies a named chain of transform functions to events.
type Transformer struct {
	mu     sync.RWMutex
	stages []stage
}

type stage struct {
	name string
	fn   TransformFunc
}

// New returns an empty Transformer.
func New() *Transformer {
	return &Transformer{}
}

// Register adds a named transform stage. Names must be unique.
func (t *Transformer) Register(name string, fn TransformFunc) error {
	if name == "" {
		return errors.New("eventtransform: stage name must not be empty")
	}
	if fn == nil {
		return errors.New("eventtransform: transform func must not be nil")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	for _, s := range t.stages {
		if s.name == name {
			return fmt.Errorf("eventtransform: stage %q already registered", name)
		}
	}
	t.stages = append(t.stages, stage{name: name, fn: fn})
	return nil
}

// Apply runs the event through all registered stages in order.
// If any stage returns an error the pipeline stops and the error is returned.
func (t *Transformer) Apply(e alert.Event) (*alert.Event, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	cur := &e
	var err error
	for _, s := range t.stages {
		cur, err = s.fn(cur)
		if err != nil {
			return nil, fmt.Errorf("eventtransform: stage %q: %w", s.name, err)
		}
		if cur == nil {
			return nil, fmt.Errorf("eventtransform: stage %q returned nil event", s.name)
		}
	}
	return cur, nil
}

// Len returns the number of registered stages.
func (t *Transformer) Len() int {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.stages)
}
