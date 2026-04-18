package dedupe

import (
	"testing"
	"time"
)

func TestAllow_FirstEventAllowed(t *testing.T) {
	f := New(5 * time.Second)
	if !f.Allow("port:tcp:8080:opened") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_DuplicateWithinWindowSuppressed(t *testing.T) {
	f := New(5 * time.Second)
	f.Allow("port:tcp:8080:opened")
	if f.Allow("port:tcp:8080:opened") {
		t.Fatal("expected duplicate within window to be suppressed")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	f := New(5 * time.Second)
	f.Allow("port:tcp:8080:opened")
	if !f.Allow("port:tcp:9090:opened") {
		t.Fatal("expected different key to be allowed")
	}
}

func TestAllow_AllowedAfterWindowExpires(t *testing.T) {
	now := time.Now()
	f := New(1 * time.Second)
	f.now = func() time.Time { return now }
	f.Allow("port:tcp:8080:opened")

	// Advance time past the window.
	f.now = func() time.Time { return now.Add(2 * time.Second) }
	if !f.Allow("port:tcp:8080:opened") {
		t.Fatal("expected event to be allowed after window expires")
	}
}

func TestReset_ClearsEntries(t *testing.T) {
	f := New(5 * time.Second)
	f.Allow("port:tcp:8080:opened")
	f.Reset()
	if !f.Allow("port:tcp:8080:opened") {
		t.Fatal("expected event to be allowed after reset")
	}
}

func TestLen_ReflectsActiveEntries(t *testing.T) {
	now := time.Now()
	f := New(1 * time.Second)
	f.now = func() time.Time { return now }
	f.Allow("a")
	f.Allow("b")
	if f.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", f.Len())
	}
	f.now = func() time.Time { return now.Add(2 * time.Second) }
	if f.Len() != 0 {
		t.Fatalf("expected 0 entries after expiry, got %d", f.Len())
	}
}
