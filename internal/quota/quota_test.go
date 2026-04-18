package quota

import (
	"testing"
	"time"
)

func TestAllow_UnderLimit(t *testing.T) {
	q := New(Config{Limit: 3, Window: time.Second})
	for i := 0; i < 3; i++ {
		if !q.Allow("k") {
			t.Fatalf("expected allow on call %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	q := New(Config{Limit: 2, Window: time.Second})
	q.Allow("k")
	q.Allow("k")
	if q.Allow("k") {
		t.Fatal("expected deny after limit reached")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	q := New(Config{Limit: 1, Window: time.Second})
	q.Allow("a")
	if !q.Allow("b") {
		t.Fatal("different key should be allowed")
	}
}

func TestAllow_ResetsAfterWindow(t *testing.T) {
	q := New(Config{Limit: 1, Window: 20 * time.Millisecond})
	q.Allow("k")
	if q.Allow("k") {
		t.Fatal("expected deny within window")
	}
	time.Sleep(30 * time.Millisecond)
	if !q.Allow("k") {
		t.Fatal("expected allow after window expired")
	}
}

func TestRemaining_Full(t *testing.T) {
	q := New(Config{Limit: 5, Window: time.Second})
	if r := q.Remaining("k"); r != 5 {
		t.Fatalf("expected 5, got %d", r)
	}
}

func TestRemaining_DecreasesOnAllow(t *testing.T) {
	q := New(Config{Limit: 5, Window: time.Second})
	q.Allow("k")
	q.Allow("k")
	if r := q.Remaining("k"); r != 3 {
		t.Fatalf("expected 3, got %d", r)
	}
}

func TestReset_ClearsState(t *testing.T) {
	q := New(Config{Limit: 1, Window: time.Second})
	q.Allow("k")
	q.Reset()
	if !q.Allow("k") {
		t.Fatal("expected allow after reset")
	}
}

func TestDefaultConfig_SaneValues(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Limit <= 0 {
		t.Fatal("limit must be positive")
	}
	if cfg.Window <= 0 {
		t.Fatal("window must be positive")
	}
}
