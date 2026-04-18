package rollup_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/rollup"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(kind alert.Kind, port int) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestAdd_BuffersEvent(t *testing.T) {
	r := rollup.New(200*time.Millisecond, func(_ rollup.Group) {})
	r.Add(makeEvent(alert.Opened, 8080))
	if r.Len() != 1 {
		t.Fatalf("expected 1 buffered event, got %d", r.Len())
	}
}

func TestFlush_GroupsOpenedAndClosed(t *testing.T) {
	results := make(chan rollup.Group, 1)
	r := rollup.New(200*time.Millisecond, func(g rollup.Group) {
		results <- g
	})

	r.Add(makeEvent(alert.Opened, 8080))
	r.Add(makeEvent(alert.Closed, 9090))
	r.Add(makeEvent(alert.Opened, 443))
	r.Flush()

	select {
	case g := <-results:
		if len(g.Opened) != 2 {
			t.Errorf("expected 2 opened, got %d", len(g.Opened))
		}
		if len(g.Closed) != 1 {
			t.Errorf("expected 1 closed, got %d", len(g.Closed))
		}
	case <-time.After(time.Second):
		t.Fatal("flush callback not called")
	}
}

func TestFlush_ClearsBuffer(t *testing.T) {
	r := rollup.New(200*time.Millisecond, func(_ rollup.Group) {})
	r.Add(makeEvent(alert.Opened, 80))
	r.Flush()
	if r.Len() != 0 {
		t.Errorf("expected empty buffer after flush, got %d", r.Len())
	}
}

func TestAdd_AutoFlushAfterWindow(t *testing.T) {
	results := make(chan rollup.Group, 1)
	r := rollup.New(50*time.Millisecond, func(g rollup.Group) {
		results <- g
	})
	r.Add(makeEvent(alert.Opened, 3000))

	select {
	case g := <-results:
		if len(g.Opened) != 1 {
			t.Errorf("expected 1 opened event, got %d", len(g.Opened))
		}
	case <-time.After(time.Second):
		t.Fatal("auto-flush did not fire")
	}
}
