package eventprojector_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventprojector"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(proto, addr string, kind alert.Kind) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{
			Protocol: proto,
			Address:  addr,
		},
	}
}

func TestNew_NotNil(t *testing.T) {
	p := eventprojector.New()
	if p == nil {
		t.Fatal("expected non-nil projector")
	}
}

func TestApply_And_Get(t *testing.T) {
	p := eventprojector.New()
	e := makeEvent("tcp", "0.0.0.0:8080", alert.Opened)
	p.Apply(e)

	got, ok := p.Get("tcp", "0.0.0.0:8080")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if got.Port.Address != e.Port.Address {
		t.Errorf("address mismatch: got %s want %s", got.Port.Address, e.Port.Address)
	}
}

func TestGet_Missing_ReturnsFalse(t *testing.T) {
	p := eventprojector.New()
	_, ok := p.Get("tcp", "0.0.0.0:9999")
	if ok {
		t.Fatal("expected no entry")
	}
}

func TestApply_Overwrites_PreviousEntry(t *testing.T) {
	p := eventprojector.New()
	p.Apply(makeEvent("tcp", "0.0.0.0:443", alert.Opened))
	p.Apply(makeEvent("tcp", "0.0.0.0:443", alert.Closed))

	got, ok := p.Get("tcp", "0.0.0.0:443")
	if !ok {
		t.Fatal("expected entry")
	}
	if got.Kind != alert.Closed {
		t.Errorf("expected Closed, got %v", got.Kind)
	}
}

func TestSnapshot_ReturnsCopy(t *testing.T) {
	p := eventprojector.New()
	p.Apply(makeEvent("tcp", "0.0.0.0:80", alert.Opened))
	p.Apply(makeEvent("udp", "0.0.0.0:53", alert.Opened))

	snap := p.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	p := eventprojector.New()
	p.Apply(makeEvent("tcp", "0.0.0.0:22", alert.Opened))
	p.Remove("tcp", "0.0.0.0:22")

	_, ok := p.Get("tcp", "0.0.0.0:22")
	if ok {
		t.Fatal("expected entry to be removed")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	p := eventprojector.New()
	p.Apply(makeEvent("tcp", "0.0.0.0:80", alert.Opened))
	p.Apply(makeEvent("tcp", "0.0.0.0:443", alert.Opened))
	p.Reset()

	if p.Len() != 0 {
		t.Errorf("expected 0 entries after reset, got %d", p.Len())
	}
}

func TestLen_ReflectsCount(t *testing.T) {
	p := eventprojector.New()
	if p.Len() != 0 {
		t.Fatalf("expected 0, got %d", p.Len())
	}
	p.Apply(makeEvent("tcp", "0.0.0.0:8080", alert.Opened))
	if p.Len() != 1 {
		t.Fatalf("expected 1, got %d", p.Len())
	}
}
