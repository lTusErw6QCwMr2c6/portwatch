package eventqueue

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Port: port, Protocol: "tcp"},
	}
}

func TestNew_DefaultCapacity(t *testing.T) {
	q := New(0)
	if q.cap != 64 {
		t.Fatalf("expected default cap 64, got %d", q.cap)
	}
}

func TestPush_And_Pop(t *testing.T) {
	q := New(4)
	e := makeEvent(8080)
	if err := q.Push(e); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, ok := q.Pop()
	if !ok {
		t.Fatal("expected event, got none")
	}
	if got.Port.Port != 8080 {
		t.Fatalf("expected port 8080, got %d", got.Port.Port)
	}
}

func TestPop_EmptyQueue_ReturnsFalse(t *testing.T) {
	q := New(4)
	_, ok := q.Pop()
	if ok {
		t.Fatal("expected false from empty queue")
	}
}

func TestPush_ExceedsCapacity_ReturnsError(t *testing.T) {
	q := New(2)
	_ = q.Push(makeEvent(1))
	_ = q.Push(makeEvent(2))
	err := q.Push(makeEvent(3))
	if err != ErrQueueFull {
		t.Fatalf("expected ErrQueueFull, got %v", err)
	}
	if q.Dropped() != 1 {
		t.Fatalf("expected 1 dropped, got %d", q.Dropped())
	}
}

func TestLen_ReflectsContents(t *testing.T) {
	q := New(8)
	if q.Len() != 0 {
		t.Fatal("expected empty queue")
	}
	_ = q.Push(makeEvent(80))
	_ = q.Push(makeEvent(443))
	if q.Len() != 2 {
		t.Fatalf("expected len 2, got %d", q.Len())
	}
}

func TestDrain_ReturnsAllAndClears(t *testing.T) {
	q := New(8)
	_ = q.Push(makeEvent(22))
	_ = q.Push(makeEvent(80))
	out := q.Drain()
	if len(out) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out))
	}
	if q.Len() != 0 {
		t.Fatal("expected queue empty after drain")
	}
}

func TestFIFO_Order(t *testing.T) {
	q := New(4)
	ports := []int{1, 2, 3}
	for _, p := range ports {
		_ = q.Push(makeEvent(p))
	}
	for _, want := range ports {
		e, _ := q.Pop()
		if e.Port.Port != want {
			t.Fatalf("expected port %d, got %d", want, e.Port.Port)
		}
	}
}
