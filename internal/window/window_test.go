package window

import (
	"testing"
	"time"
)

func TestCount_EmptyKey_ReturnsZero(t *testing.T) {
	c := New(time.Second)
	if got := c.Count("k"); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_IncrementsCount(t *testing.T) {
	c := New(time.Second)
	c.Add("k")
	c.Add("k")
	if got := c.Count("k"); got != 2 {
		t.Fatalf("expected 2, got %d", got)
	}
}

func TestCount_EvictsExpiredEntries(t *testing.T) {
	c := New(50 * time.Millisecond)
	c.Add("k")
	time.Sleep(80 * time.Millisecond)
	if got := c.Count("k"); got != 0 {
		t.Fatalf("expected 0 after expiry, got %d", got)
	}
}

func TestCount_DifferentKeysAreIndependent(t *testing.T) {
	c := New(time.Second)
	c.Add("a")
	c.Add("a")
	c.Add("b")
	if c.Count("a") != 2 {
		t.Fatalf("expected 2 for a")
	}
	if c.Count("b") != 1 {
		t.Fatalf("expected 1 for b")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	c := New(time.Second)
	c.Add("x")
	c.Add("y")
	c.Reset()
	if c.Count("x") != 0 || c.Count("y") != 0 {
		t.Fatal("expected all counts cleared after reset")
	}
}

func TestAdd_OnlyRecentCountedAfterMix(t *testing.T) {
	c := New(60 * time.Millisecond)
	c.Add("k")
	time.Sleep(70 * time.Millisecond)
	c.Add("k")
	if got := c.Count("k"); got != 1 {
		t.Fatalf("expected 1 (only recent), got %d", got)
	}
}
