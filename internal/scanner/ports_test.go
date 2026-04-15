package scanner

import (
	"testing"
)

func TestDiff_NewPortOpened(t *testing.T) {
	prev := Snapshot{}
	curr := Snapshot{
		"tcp:8080": {Port: 8080, Protocol: "tcp", State: "open"},
	}

	opened, closed := Diff(prev, curr)

	if len(opened) != 1 {
		t.Fatalf("expected 1 opened port, got %d", len(opened))
	}
	if opened[0].Port != 8080 {
		t.Errorf("expected port 8080, got %d", opened[0].Port)
	}
	if len(closed) != 0 {
		t.Errorf("expected 0 closed ports, got %d", len(closed))
	}
}

func TestDiff_PortClosed(t *testing.T) {
	prev := Snapshot{
		"tcp:443": {Port: 443, Protocol: "tcp", State: "open"},
	}
	curr := Snapshot{}

	opened, closed := Diff(prev, curr)

	if len(closed) != 1 {
		t.Fatalf("expected 1 closed port, got %d", len(closed))
	}
	if closed[0].Port != 443 {
		t.Errorf("expected port 443, got %d", closed[0].Port)
	}
	if len(opened) != 0 {
		t.Errorf("expected 0 opened ports, got %d", len(opened))
	}
}

func TestDiff_NoChange(t *testing.T) {
	snap := Snapshot{
		"tcp:80": {Port: 80, Protocol: "tcp", State: "open"},
	}
	opened, closed := Diff(snap, snap)

	if len(opened) != 0 || len(closed) != 0 {
		t.Errorf("expected no changes, got opened=%d closed=%d", len(opened), len(closed))
	}
}

func TestFormatPort(t *testing.T) {
	ps := PortState{Port: 22, Protocol: "tcp", State: "open"}
	got := FormatPort(ps)
	want := "TCP/22"
	if got != want {
		t.Errorf("FormatPort() = %q, want %q", got, want)
	}
}

func TestScanPorts_InvalidRange(t *testing.T) {
	_, err := ScanPorts(9000, 8000)
	if err == nil {
		t.Error("expected error for invalid port range, got nil")
	}
}
