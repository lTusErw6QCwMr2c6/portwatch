package fingerprint_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/fingerprint"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto, process, action string) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Port: port, Proto: proto, Process: process},
		Action: action,
	}
}

func TestCompute_ReturnsDeterministicFingerprint(t *testing.T) {
	g := fingerprint.New(nil)
	e := makeEvent(8080, "tcp", "nginx", "opened")
	f1 := g.Compute(e)
	f2 := g.Compute(e)
	if f1 != f2 {
		t.Fatalf("expected same fingerprint, got %s and %s", f1, f2)
	}
}

func TestCompute_DifferentPorts_DifferentFingerprints(t *testing.T) {
	g := fingerprint.New(nil)
	e1 := makeEvent(80, "tcp", "", "opened")
	e2 := makeEvent(443, "tcp", "", "opened")
	if g.Compute(e1) == g.Compute(e2) {
		t.Fatal("expected different fingerprints for different ports")
	}
}

func TestCompute_CustomFields_UsesOnlySpecified(t *testing.T) {
	g := fingerprint.New([]string{"port", "proto"})
	e1 := makeEvent(8080, "tcp", "nginx", "opened")
	e2 := makeEvent(8080, "tcp", "caddy", "closed")
	if g.Compute(e1) != g.Compute(e2) {
		t.Fatal("expected same fingerprint when process and action excluded")
	}
}

func TestCompute_NonEmptyHex(t *testing.T) {
	g := fingerprint.New(nil)
	e := makeEvent(22, "tcp", "sshd", "opened")
	f := g.Compute(e)
	if len(f) == 0 {
		t.Fatal("expected non-empty fingerprint")
	}
}

func TestComputeAll_ReturnsCorrectCount(t *testing.T) {
	g := fingerprint.New(nil)
	events := []alert.Event{
		makeEvent(80, "tcp", "", "opened"),
		makeEvent(443, "tcp", "", "opened"),
		makeEvent(22, "tcp", "", "closed"),
	}
	fingerprints := g.ComputeAll(events)
	if len(fingerprints) != 3 {
		t.Fatalf("expected 3 fingerprints, got %d", len(fingerprints))
	}
}
