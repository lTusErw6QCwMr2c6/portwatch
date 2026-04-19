package eventid

import (
	"strings"
	"testing"
)

func TestNew_NotEmpty(t *testing.T) {
	id := New()
	if id.String() == "" {
		t.Fatal("expected non-empty ID string")
	}
}

func TestNew_SequenceIncreases(t *testing.T) {
	Reset()
	a := New()
	b := New()
	if b.Sequence <= a.Sequence {
		t.Fatalf("expected b.Sequence > a.Sequence, got %d <= %d", b.Sequence, a.Sequence)
	}
}

func TestString_ContainsThreeParts(t *testing.T) {
	id := New()
	parts := strings.Split(id.String(), "-")
	if len(parts) != 3 {
		t.Fatalf("expected 3 parts in ID string, got %d: %s", len(parts), id.String())
	}
}

func TestBefore_EarlierIsBefore(t *testing.T) {
	Reset()
	a := New()
	b := New()
	if !a.Before(b) {
		t.Fatal("expected a to be before b")
	}
	if b.Before(a) {
		t.Fatal("expected b not to be before a")
	}
}

func TestNew_UniqueRandom(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 50; i++ {
		id := New()
		if seen[id.Random] {
			t.Fatalf("duplicate random segment: %s", id.Random)
		}
		seen[id.Random] = true
	}
}

func TestReset_ResetsSequence(t *testing.T) {
	_ = New()
	_ = New()
	Reset()
	id := New()
	if id.Sequence != 1 {
		t.Fatalf("expected sequence 1 after reset, got %d", id.Sequence)
	}
}
