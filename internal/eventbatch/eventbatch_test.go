package eventbatch

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
)

func makeEvent(port int) alert.Event {
	return alert.Event{Port: port, Proto: "tcp", Type: alert.Opened}
}

func TestNew_InvalidMaxSize_ReturnsError(t *testing.T) {
	_, err := New(0, time.Second, func([]alert.Event) {})
	if err == nil {
		t.Fatal("expected error for maxSize=0")
	}
}

func TestNew_InvalidInterval_ReturnsError(t *testing.T) {
	_, err := New(10, 0, func([]alert.Event) {})
	if err == nil {
		t.Fatal("expected error for interval=0")
	}
}

func TestNew_NilHandler_ReturnsError(t *testing.T) {
	_, err := New(10, time.Second, nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestAdd_FlushesWhenMaxSizeReached(t *testing.T) {
	var mu sync.Mutex
	var received []alert.Event

	b, err := New(3, 10*time.Second, func(events []alert.Event) {
		mu.Lock()
		received = append(received, events...)
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer b.Stop()

	b.Add(makeEvent(80))
	b.Add(makeEvent(443))
	b.Add(makeEvent(8080))

	time.Sleep(50 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 3 {
		t.Fatalf("expected 3 events, got %d", len(received))
	}
}

func TestAdd_FlushesOnInterval(t *testing.T) {
	var mu sync.Mutex
	var received []alert.Event

	b, err := New(100, 50*time.Millisecond, func(events []alert.Event) {
		mu.Lock()
		received = append(received, events...)
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer b.Stop()

	b.Add(makeEvent(22))
	b.Add(makeEvent(25))

	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(received) < 2 {
		t.Fatalf("expected at least 2 events via interval flush, got %d", len(received))
	}
}

func TestStop_FlushesRemainingEvents(t *testing.T) {
	var mu sync.Mutex
	var received []alert.Event

	b, err := New(100, 10*time.Second, func(events []alert.Event) {
		mu.Lock()
		received = append(received, events...)
		mu.Unlock()
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	b.Add(makeEvent(3000))
	b.Stop()

	mu.Lock()
	defer mu.Unlock()
	if len(received) != 1 {
		t.Fatalf("expected 1 event after Stop, got %d", len(received))
	}
}
