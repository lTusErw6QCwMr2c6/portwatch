package eventchain_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventchain"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Number: port, Protocol: proto},
	}
}

func TestNew_NotNil(t *testing.T) {
	c := eventchain.New()
	if c == nil {
		t.Fatal("expected non-nil Chain")
	}
}

func TestLink_And_Parent(t *testing.T) {
	c := eventchain.New()
	parent := makeEvent(80, "tcp")
	child := makeEvent(8080, "tcp")

	c.Link(parent, child)

	pID, ok := c.Parent(child)
	if !ok {
		t.Fatal("expected parent to be found")
	}
	if pID != parent.Port.String() {
		t.Errorf("expected parent %q, got %q", parent.Port.String(), pID)
	}
}

func TestParent_UnknownChild_ReturnsFalse(t *testing.T) {
	c := eventchain.New()
	_, ok := c.Parent(makeEvent(9999, "tcp"))
	if ok {
		t.Fatal("expected no parent for unknown child")
	}
}

func TestChildren_ReturnsAll(t *testing.T) {
	c := eventchain.New()
	parent := makeEvent(443, "tcp")
	child1 := makeEvent(8443, "tcp")
	child2 := makeEvent(9443, "tcp")

	c.Link(parent, child1)
	c.Link(parent, child2)

	kids := c.Children(parent)
	if len(kids) != 2 {
		t.Fatalf("expected 2 children, got %d", len(kids))
	}
}

func TestChildren_NoChildren_ReturnsEmpty(t *testing.T) {
	c := eventchain.New()
	kids := c.Children(makeEvent(22, "tcp"))
	if len(kids) != 0 {
		t.Fatalf("expected empty children, got %d", len(kids))
	}
}

func TestReset_ClearsAll(t *testing.T) {
	c := eventchain.New()
	parent := makeEvent(80, "tcp")
	child := makeEvent(8080, "tcp")
	c.Link(parent, child)

	c.Reset()

	_, ok := c.Parent(child)
	if ok {
		t.Fatal("expected parent to be cleared after reset")
	}
	if len(c.Children(parent)) != 0 {
		t.Fatal("expected children to be cleared after reset")
	}
}
