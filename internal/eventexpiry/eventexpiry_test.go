package eventexpiry

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Port: port, Protocol: "tcp"},
	}
}

func TestNew_NotNil(t *testing.T) {
	r := New(time.Second, nil)
	if r == nil {
		t.Fatal("expected non-nil registry")
	}
}

func TestTrack_IncreasesLen(t *testing.T) {
	r := New(time.Second, nil)
	r.Track("a", makeEvent(80))
	r.Track("b", makeEvent(443))
	if r.Len() != 2 {
		t.Fatalf("expected 2, got %d", r.Len())
	}
}

func TestCancel_RemovesEntry(t *testing.T) {
	r := New(time.Second, nil)
	r.Track("k", makeEvent(22))
	if !r.Cancel("k") {
		t.Fatal("expected Cancel to return true")
	}
	if r.Len() != 0 {
		t.Fatal("expected empty registry after cancel")
	}
}

func TestCancel_UnknownKey_ReturnsFalse(t *testing.T) {
	r := New(time.Second, nil)
	if r.Cancel("missing") {
		t.Fatal("expected false for unknown key")
	}
}

func TestExpiry_CallbackFired(t *testing.T) {
	var mu sync.Mutex
	var fired []alert.Event

	r := New(20*time.Millisecond, func(e alert.Event) {
		mu.Lock()
		fired = append(fired, e)
		mu.Unlock()
	})

	r.Track("x", makeEvent(8080))

	time.Sleep(60 * time.Millisecond)

	mu.Lock()
	n := len(fired)
	mu.Unlock()

	if n != 1 {
		t.Fatalf("expected 1 expiry callback, got %d", n)
	}
	if r.Len() != 0 {
		t.Fatal("expected registry to be empty after expiry")
	}
}

func TestTrack_ResetTimer_OnDuplicate(t *testing.T) {
	called := make(chan struct{}, 1)
	r := New(40*time.Millisecond, func(_ alert.Event) {
		called <- struct{}{}
	})

	r.Track("dup", makeEvent(9000))
	time.Sleep(20 * time.Millisecond)
	// re-track resets the timer
	r.Track("dup", makeEvent(9000))

	select {
	case <-called:
		// good
	case <-time.After(120 * time.Millisecond):
		t.Fatal("expiry callback never fired after reset")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	r := New(time.Second, nil)
	r.Track("a", makeEvent(1))
	r.Track("b", makeEvent(2))
	r.Reset()
	if r.Len() != 0 {
		t.Fatalf("expected 0 after Reset, got %d", r.Len())
	}
}
