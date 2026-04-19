package dispatch_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/dispatch"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestNew_NotNil(t *testing.T) {
	r := dispatch.New()
	if r == nil {
		t.Fatal("expected non-nil router")
	}
}

func TestRegister_IncreasesLen(t *testing.T) {
	r := dispatch.New()
	r.Register("a", func(alert.Event) error { return nil })
	if r.Len() != 1 {
		t.Fatalf("expected 1 handler, got %d", r.Len())
	}
}

func TestDeregister_RemovesHandler(t *testing.T) {
	r := dispatch.New()
	r.Register("a", func(alert.Event) error { return nil })
	r.Deregister("a")
	if r.Len() != 0 {
		t.Fatalf("expected 0 handlers, got %d", r.Len())
	}
}

func TestDispatch_CallsAllHandlers(t *testing.T) {
	r := dispatch.New()
	called := map[string]bool{}
	r.Register("x", func(alert.Event) error { called["x"] = true; return nil })
	r.Register("y", func(alert.Event) error { called["y"] = true; return nil })

	errs := r.Dispatch(makeEvent(8080))
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if !called["x"] || !called["y"] {
		t.Fatal("expected both handlers to be called")
	}
}

func TestDispatch_CollectsErrors(t *testing.T) {
	r := dispatch.New()
	r.Register("bad", func(alert.Event) error { return errors.New("fail") })

	errs := r.Dispatch(makeEvent(443))
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestDispatch_NoHandlers_ReturnsNil(t *testing.T) {
	r := dispatch.New()
	errs := r.Dispatch(makeEvent(22))
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}
