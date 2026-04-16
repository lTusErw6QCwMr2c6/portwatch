package supervisor

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockComponent struct {
	name    string
	err     error
	calls   int
	maxFail int
}

func (m *mockComponent) Name() string { return m.name }

func (m *mockComponent) Start(_ context.Context) error {
	m.calls++
	if m.calls <= m.maxFail {
		return m.err
	}
	return nil
}

func (m *mockComponent) Stop() error { return nil }

func TestNew_NotNil(t *testing.T) {
	s := New(3, 10*time.Millisecond)
	if s == nil {
		t.Fatal("expected non-nil supervisor")
	}
}

func TestRegister_AddsEntry(t *testing.T) {
	s := New(3, 10*time.Millisecond)
	c := &mockComponent{name: "scanner"}
	s.Register(c)
	if _, ok := s.entries["scanner"]; !ok {
		t.Fatal("expected entry for scanner")
	}
}

func TestSupervise_RestartsOnError(t *testing.T) {
	s := New(5, 5*time.Millisecond)
	c := &mockComponent{name: "watcher", err: errors.New("boom"), maxFail: 2}
	s.Register(c)
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	s.Run(ctx)
	<-ctx.Done()
	e := s.entries["watcher"]
	if e.Restarts == 0 {
		t.Error("expected at least one restart")
	}
}

func TestSupervise_FailsAfterMaxRetries(t *testing.T) {
	s := New(2, 5*time.Millisecond)
	c := &mockComponent{name: "probe", err: errors.New("fail"), maxFail: 99}
	s.Register(c)
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()
	s.Run(ctx)
	time.Sleep(200 * time.Millisecond)
	e := s.entries["probe"]
	if e.Status != StatusFailed {
		t.Errorf("expected StatusFailed, got %v", e.Status)
	}
	cancel()
}

func TestEntries_ReturnsCopy(t *testing.T) {
	s := New(1, 5*time.Millisecond)
	s.Register(&mockComponent{name: "x"})
	out := s.Entries()
	out["x"].Restarts = 99
	if s.entries["x"].Restarts == 99 {
		t.Error("Entries should return a copy, not a reference")
	}
}
