// Package eventfunnel provides a fan-in aggregator that merges multiple
// named event streams into a single output channel with optional priority.
package eventfunnel

import (
	"errors"
	"sync"

	"github.com/example/portwatch/internal/alert"
)

// Source represents a named input channel.
type Source struct {
	Name     string
	Priority int
	Ch       <-chan alert.Event
}

// Funnel merges multiple event sources into one output channel.
type Funnel struct {
	mu      sync.RWMutex
	sources []Source
	out     chan alert.Event
	stop    chan struct{}
	wg      sync.WaitGroup
}

// New creates a Funnel with the given output buffer size.
func New(bufSize int) *Funnel {
	if bufSize <= 0 {
		bufSize = 64
	}
	return &Funnel{
		out:  make(chan alert.Event, bufSize),
		stop: make(chan struct{}),
	}
}

// Add registers a named source channel. Returns an error on duplicate name.
func (f *Funnel) Add(s Source) error {
	if s.Name == "" {
		return errors.New("eventfunnel: source name must not be empty")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, existing := range f.sources {
		if existing.Name == s.Name {
			return errors.New("eventfunnel: duplicate source name: " + s.Name)
		}
	}
	f.sources = append(f.sources, s)
	f.wg.Add(1)
	go f.forward(s)
	return nil
}

// Out returns the merged output channel.
func (f *Funnel) Out() <-chan alert.Event { return f.out }

// Stop signals all forwarders to exit and closes the output channel.
func (f *Funnel) Stop() {
	close(f.stop)
	f.wg.Wait()
	close(f.out)
}

// Len returns the number of registered sources.
func (f *Funnel) Len() int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	return len(f.sources)
}

func (f *Funnel) forward(s Source) {
	defer f.wg.Done()
	for {
		select {
		case ev, ok := <-s.Ch:
			if !ok {
				return
			}
			select {
			case f.out <- ev:
			case <-f.stop:
				return
			}
		case <-f.stop:
			return
		}
	}
}
