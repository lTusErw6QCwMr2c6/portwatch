package eventbatch

import (
	"errors"
	"sync"
	"time"

	"github.com/user/portwatch/internal/alert"
)

// Batch accumulates events and flushes them either when the batch reaches
// a maximum size or when a flush interval elapses.
type Batch struct {
	mu       sync.Mutex
	events   []alert.Event
	maxSize  int
	interval time.Duration
	handler  func([]alert.Event)
	stop     chan struct{}
	wg       sync.WaitGroup
}

// New creates a Batch that flushes to handler when maxSize events are
// buffered or when interval elapses, whichever comes first.
func New(maxSize int, interval time.Duration, handler func([]alert.Event)) (*Batch, error) {
	if maxSize <= 0 {
		return nil, errors.New("eventbatch: maxSize must be positive")
	}
	if interval <= 0 {
		return nil, errors.New("eventbatch: interval must be positive")
	}
	if handler == nil {
		return nil, errors.New("eventbatch: handler must not be nil")
	}
	b := &Batch{
		events:  make([]alert.Event, 0, maxSize),
		maxSize: maxSize,
		interval: interval,
		handler: handler,
		stop:    make(chan struct{}),
	}
	b.wg.Add(1)
	go b.run()
	return b, nil
}

// Add appends an event to the batch, flushing immediately if maxSize is reached.
func (b *Batch) Add(e alert.Event) {
	b.mu.Lock()
	b.events = append(b.events, e)
	if len(b.events) >= b.maxSize {
		b.flush()
	}
	b.mu.Unlock()
}

// Stop flushes any remaining events and stops the background ticker.
func (b *Batch) Stop() {
	close(b.stop)
	b.wg.Wait()
	b.mu.Lock()
	if len(b.events) > 0 {
		b.flush()
	}
	b.mu.Unlock()
}

// flush dispatches buffered events to the handler and resets the buffer.
// Caller must hold b.mu.
func (b *Batch) flush() {
	if len(b.events) == 0 {
		return
	}
	batch := make([]alert.Event, len(b.events))
	copy(batch, b.events)
	b.events = b.events[:0]
	go b.handler(batch)
}

func (b *Batch) run() {
	defer b.wg.Done()
	ticker := time.NewTicker(b.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			b.mu.Lock()
			b.flush()
			b.mu.Unlock()
		case <-b.stop:
			return
		}
	}
}
