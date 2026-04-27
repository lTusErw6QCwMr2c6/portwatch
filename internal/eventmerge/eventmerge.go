// Package eventmerge provides a mechanism to merge multiple event streams
// into a single ordered output channel.
package eventmerge

import (
	"errors"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// Merger fans-in multiple named event sources into one stream.
type Merger struct {
	mu      sync.Mutex
	sources map[string]<-chan alert.Event
	out     chan alert.Event
	stop    chan struct{}
	wg      sync.WaitGroup
}

// New creates a new Merger with a buffered output channel of the given size.
func New(bufSize int) *Merger {
	if bufSize < 1 {
		bufSize = 64
	}
	return &Merger{
		sources: make(map[string]<-chan alert.Event),
		out:     make(chan alert.Event, bufSize),
		stop:    make(chan struct{}),
	}
}

// Add registers a named source channel. Returns an error if the name is
// already registered or the merger has already been started.
func (m *Merger) Add(name string, ch <-chan alert.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if name == "" {
		return errors.New("eventmerge: source name must not be empty")
	}
	if _, exists := m.sources[name]; exists {
		return errors.New("eventmerge: source already registered: " + name)
	}
	m.sources[name] = ch
	return nil
}

// Start begins draining all registered sources into the output channel.
// Each source is consumed in its own goroutine.
func (m *Merger) Start() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, ch := range m.sources {
		src := ch
		m.wg.Add(1)
		go func() {
			defer m.wg.Done()
			for {
				select {
				case ev, ok := <-src:
					if !ok {
						return
					}
					select {
					case m.out <- ev:
					case <-m.stop:
						return
					}
				case <-m.stop:
					return
				}
			}
		}()
	}
}

// Out returns the merged output channel.
func (m *Merger) Out() <-chan alert.Event { return m.out }

// Stop signals all goroutines to exit and waits for them to finish.
func (m *Merger) Stop() {
	close(m.stop)
	m.wg.Wait()
	close(m.out)
}
