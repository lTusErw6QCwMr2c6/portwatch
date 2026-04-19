package drain

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestAdd_BuffersEvent(t *testing.T) {
	d := New(4, time.Second, func([]alert.Event) error { return nil })
	_ = d.Add(makeEvent(80))
	if d.Len() != 1 {
		t.Fatalf("expected 1 buffered event, got %d", d.Len())
	}
}

func TestFlush_CallsHandler(t *testing.T) {
	var received []alert.Event
	d := New(4, time.Second, func(evs []alert.Event) error {
		received = evs
		return nil
	})
	_ = d.Add(makeEvent(443))
	_ = d.Add(makeEvent(8080))
	if err := d.Flush(); err != nil {
		t.Fatal(err)
	}
	if len(received) != 2 {
		t.Fatalf("expected 2 events, got %d", len(received))
	}
	if d.Len() != 0 {
		t.Fatal("buffer should be empty after flush")
	}
}

func TestAdd_AutoFlushWhenFull(t *testing.T) {
	flushed := 0
	d := New(2, time.Second, func(evs []alert.Event) error {
		flushed += len(evs)
		return nil
	})
	_ = d.Add(makeEvent(1))
	_ = d.Add(makeEvent(2)) // triggers flush
	if flushed != 2 {
		t.Fatalf("expected 2 flushed, got %d", flushed)
	}
}

func TestStop_FlushesRemaining(t *testing.T) {
	var got []alert.Event
	d := New(10, time.Second, func(evs []alert.Event) error {
		got = append(got, evs...)
		return nil
	})
	_ = d.Add(makeEvent(22))
	if err := d.Stop(); err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 {
		t.Fatalf("expected 1 event after stop, got %d", len(got))
	}
}

func TestFlush_HandlerError_Propagated(t *testing.T) {
	expected := errors.New("handler error")
	d := New(4, time.Second, func([]alert.Event) error { return expected })
	_ = d.Add(makeEvent(9090))
	if err := d.Flush(); !errors.Is(err, expected) {
		t.Fatalf("expected handler error, got %v", err)
	}
}

func TestFlush_Empty_NoError(t *testing.T) {
	d := New(4, time.Second, func([]alert.Event) error {
		t.Fatal("handler should not be called on empty flush")
		return nil
	})
	if err := d.Flush(); err != nil {
		t.Fatal(err)
	}
}
