// Package eventstream provides a fan-out streaming layer that broadcasts
// alert.Event values to multiple named consumers with optional backpressure.
package eventstream

import (
	"errors"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// ErrDuplicateConsumer is returned when a consumer with the same name is registered twice.
var ErrDuplicateConsumer = errors.New("eventstream: duplicate consumer name")

// ErrUnknownConsumer is returned when an operation targets a consumer that does not exist.
var ErrUnknownConsumer = errors.New("eventstream: unknown consumer")

// Stream fans out events to registered consumers.
type Stream struct {
	mu        sync.RWMutex
	consumers map[string]chan alert.Event
	buf       int
}

// New creates a Stream where each consumer channel is buffered to buf events.
func New(buf int) *Stream {
	if buf < 1 {
		buf = 1
	}
	return &Stream{
		consumers: make(map[string]chan alert.Event),
		buf:       buf,
	}
}

// Subscribe registers a new consumer and returns its receive channel.
func (s *Stream) Subscribe(name string) (<-chan alert.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.consumers[name]; ok {
		return nil, ErrDuplicateConsumer
	}
	ch := make(chan alert.Event, s.buf)
	s.consumers[name] = ch
	return ch, nil
}

// Unsubscribe removes a consumer and closes its channel.
func (s *Stream) Unsubscribe(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	ch, ok := s.consumers[name]
	if !ok {
		return ErrUnknownConsumer
	}
	close(ch)
	delete(s.consumers, name)
	return nil
}

// Publish sends an event to all registered consumers.
// If a consumer's buffer is full the event is dropped for that consumer.
func (s *Stream) Publish(ev alert.Event) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, ch := range s.consumers {
		select {
		case ch <- ev:
		default:
		}
	}
}

// Len returns the number of active consumers.
func (s *Stream) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.consumers)
}

// Close unsubscribes all consumers.
func (s *Stream) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for name, ch := range s.consumers {
		close(ch)
		delete(s.consumers, name)
	}
}
