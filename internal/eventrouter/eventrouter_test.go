package eventrouter_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventrouter"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto, kind string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: port, Protocol: proto},
		Kind: kind,
	}
}

func TestNew_NotNil(t *testing.T) {
	r := eventrouter.New()
	if r == nil {
		t.Fatal("expected non-nil router")
	}
}

func TestRegister_AddsRoute(t *testing.T) {
	r := eventrouter.New()
	err := r.Register(eventrouter.Route{
		Name:   "test",
		Match:  func(e alert.Event) bool { return true },
		Handle: func(e alert.Event) error { return nil },
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Len() != 1 {
		t.Errorf("expected 1 route, got %d", r.Len())
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	r := eventrouter.New()
	err := r.Register(eventrouter.Route{
		Name:   "",
		Match:  func(e alert.Event) bool { return true },
		Handle: func(e alert.Event) error { return nil },
	})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestDispatch_MatchingRoute_CallsHandler(t *testing.T) {
	r := eventrouter.New()
	called := false
	_ = r.Register(eventrouter.Route{
		Name:   "catch-all",
		Match:  func(e alert.Event) bool { return true },
		Handle: func(e alert.Event) error { called = true; return nil },
	})
	if err := r.Dispatch(makeEvent(80, "tcp", "opened")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Error("expected handler to be called")
	}
}

func TestDispatch_NoMatch_ReturnsErrNoMatch(t *testing.T) {
	r := eventrouter.New()
	_ = r.Register(eventrouter.Route{
		Name:   "never",
		Match:  func(e alert.Event) bool { return false },
		Handle: func(e alert.Event) error { return nil },
	})
	err := r.Dispatch(makeEvent(443, "tcp", "opened"))
	if !errors.Is(err, eventrouter.ErrNoMatch) {
		t.Errorf("expected ErrNoMatch, got %v", err)
	}
}

func TestDispatch_FirstMatchWins(t *testing.T) {
	r := eventrouter.New()
	hits := []string{}
	_ = r.Register(eventrouter.Route{
		Name:  "first",
		Match: func(e alert.Event) bool { return true },
		Handle: func(e alert.Event) error {
			hits = append(hits, "first")
			return nil
		},
	})
	_ = r.Register(eventrouter.Route{
		Name:  "second",
		Match: func(e alert.Event) bool { return true },
		Handle: func(e alert.Event) error {
			hits = append(hits, "second")
			return nil
		},
	})
	_ = r.Dispatch(makeEvent(22, "tcp", "opened"))
	if len(hits) != 1 || hits[0] != "first" {
		t.Errorf("expected only first route to fire, got %v", hits)
	}
}
