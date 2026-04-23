package eventpause

import (
	"testing"
)

func TestNew_GateIsOpen(t *testing.T) {
	g := New()
	if g.IsPaused() {
		t.Fatal("expected gate to be open after New")
	}
}

func TestAllow_OpenGate_ReturnsNil(t *testing.T) {
	g := New()
	if err := g.Allow(); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestAllow_PausedGate_ReturnsErrPaused(t *testing.T) {
	g := New()
	g.Pause()
	if err := g.Allow(); err != ErrPaused {
		t.Fatalf("expected ErrPaused, got %v", err)
	}
}

func TestResume_ReopensGate(t *testing.T) {
	g := New()
	g.Pause()
	g.Resume()
	if g.IsPaused() {
		t.Fatal("expected gate to be open after Resume")
	}
	if err := g.Allow(); err != nil {
		t.Fatalf("expected nil after resume, got %v", err)
	}
}

func TestStats_TracksAllowedAndDropped(t *testing.T) {
	g := New()

	_ = g.Allow()
	_ = g.Allow()

	g.Pause()
	_ = g.Allow()

	allowed, dropped := g.Stats()
	if allowed != 2 {
		t.Fatalf("expected 2 allowed, got %d", allowed)
	}
	if dropped != 1 {
		t.Fatalf("expected 1 dropped, got %d", dropped)
	}
}

func TestReset_ClearsCountersAndOpensGate(t *testing.T) {
	g := New()
	g.Pause()
	_ = g.Allow()
	g.Reset()

	allowed, dropped := g.Stats()
	if allowed != 0 || dropped != 0 {
		t.Fatalf("expected zero stats after Reset, got allowed=%d dropped=%d", allowed, dropped)
	}
	if g.IsPaused() {
		t.Fatal("expected gate to be open after Reset")
	}
}

func TestPause_IsIdempotent(t *testing.T) {
	g := New()
	g.Pause()
	g.Pause()
	if !g.IsPaused() {
		t.Fatal("expected gate to remain paused")
	}
}
