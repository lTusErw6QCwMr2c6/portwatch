package cooldown_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/cooldown"
)

func TestAllow_FirstCallAllowed(t *testing.T) {
	c := cooldown.New(100 * time.Millisecond)
	if !c.Allow("key1") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallWithinCooldownBlocked(t *testing.T) {
	c := cooldown.New(200 * time.Millisecond)
	c.Allow("key1")
	if c.Allow("key1") {
		t.Fatal("expected second call within cooldown to be blocked")
	}
}

func TestAllow_AllowedAfterCooldownExpires(t *testing.T) {
	c := cooldown.New(30 * time.Millisecond)
	c.Allow("key1")
	time.Sleep(50 * time.Millisecond)
	if !c.Allow("key1") {
		t.Fatal("expected call after cooldown expiry to be allowed")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	c := cooldown.New(200 * time.Millisecond)
	c.Allow("a")
	if !c.Allow("b") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestRemaining_ZeroWhenNotSeen(t *testing.T) {
	c := cooldown.New(100 * time.Millisecond)
	if r := c.Remaining("missing"); r != 0 {
		t.Fatalf("expected 0, got %v", r)
	}
}

func TestRemaining_PositiveAfterAllow(t *testing.T) {
	c := cooldown.New(500 * time.Millisecond)
	c.Allow("key")
	if r := c.Remaining("key"); r <= 0 {
		t.Fatalf("expected positive remaining, got %v", r)
	}
}

func TestReset_ClearsKeys(t *testing.T) {
	c := cooldown.New(500 * time.Millisecond)
	c.Allow("key")
	c.Reset()
	if !c.Allow("key") {
		t.Fatal("expected key to be allowed after reset")
	}
}
