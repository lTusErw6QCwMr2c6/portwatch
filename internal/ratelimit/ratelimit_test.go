package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestAllow_FirstEventAllowed(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected first event to be allowed")
	}
}

func TestAllow_DuplicateWithinCooldownSuppressed(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("tcp:8080")
	if l.Allow("tcp:8080") {
		t.Fatal("expected duplicate event within cooldown to be suppressed")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("tcp:8080")
	if !l.Allow("tcp:9090") {
		t.Fatal("expected different key to be allowed independently")
	}
}

func TestAllow_AllowedAfterCooldownExpires(t *testing.T) {
	l := ratelimit.New(20 * time.Millisecond)
	l.Allow("tcp:8080")
	time.Sleep(30 * time.Millisecond)
	if !l.Allow("tcp:8080") {
		t.Fatal("expected event to be allowed after cooldown expired")
	}
}

func TestReset_ClearsAllKeys(t *testing.T) {
	l := ratelimit.New(1 * time.Second)
	l.Allow("tcp:8080")
	l.Allow("udp:5353")
	l.Reset()
	if !l.Allow("tcp:8080") {
		t.Fatal("expected event to be allowed after reset")
	}
}

func TestPrune_RemovesExpiredEntries(t *testing.T) {
	l := ratelimit.New(20 * time.Millisecond)
	l.Allow("tcp:8080")
	time.Sleep(30 * time.Millisecond)
	l.Prune()
	// After pruning the expired entry, the key should be allowed again
	if !l.Allow("tcp:8080") {
		t.Fatal("expected pruned key to be allowed through")
	}
}
