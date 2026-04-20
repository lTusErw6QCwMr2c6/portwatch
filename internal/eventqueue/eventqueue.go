package eventqueue

import (
	"errors"
	"sync"

	"github.com/user/portwatch/internal/alert"
)

// ErrQueueFull is returned when the queue has reached its maximum capacity.
var ErrQueueFull = errors.New("event queue is full")

// Queue is a bounded, thread-safe FIFO queue for alert events.
type Queue struct {
	mu      sync.Mutex
	items   []alert.Event
	cap     int
	dropped int
}

// New creates a new Queue with the given capacity.
func New(capacity int) *Queue {
	if capacity <= 0 {
		capacity = 64
	}
	return &Queue{
		items: make([]alert.Event, 0, capacity),
		cap:   capacity,
	}
}

// Push adds an event to the back of the queue.
// Returns ErrQueueFull if the queue is at capacity.
func (q *Queue) Push(e alert.Event) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.cap {
		q.dropped++
		return ErrQueueFull
	}
	q.items = append(q.items, e)
	return nil
}

// Pop removes and returns the front event.
// Returns false if the queue is empty.
func (q *Queue) Pop() (alert.Event, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return alert.Event{}, false
	}
	e := q.items[0]
	q.items = q.items[1:]
	return e, true
}

// Len returns the current number of events in the queue.
func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// Dropped returns the total number of events dropped due to a full queue.
func (q *Queue) Dropped() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.dropped
}

// Drain removes and returns all events currently in the queue.
func (q *Queue) Drain() []alert.Event {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]alert.Event, len(q.items))
	copy(out, q.items)
	q.items = q.items[:0]
	return out
}

// Stats returns a snapshot of the queue's current length, capacity, and
// total number of dropped events. Useful for metrics and diagnostics.
func (q *Queue) Stats() (length, capacity, dropped int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items), q.cap, q.dropped
}
