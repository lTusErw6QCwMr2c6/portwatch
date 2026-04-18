package healthcheck

import (
	"testing"
)

func TestNew_EmptyMonitor(t *testing.T) {
	m := New()
	if len(m.All()) != 0 {
		t.Fatal("expected no checks on new monitor")
	}
}

func TestRecord_StoresCheck(t *testing.T) {
	m := New()
	m.Record("scanner", StatusOK, "running")
	checks := m.All()
	if len(checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(checks))
	}
	if checks[0].Name != "scanner" {
		t.Errorf("unexpected name: %s", checks[0].Name)
	}
}

func TestRecord_OverwritesPrevious(t *testing.T) {
	m := New()
	m.Record("scanner", StatusOK, "running")
	m.Record("scanner", StatusDown, "failed")
	checks := m.All()
	if len(checks) != 1 {
		t.Fatalf("expected 1 check, got %d", len(checks))
	}
	if checks[0].Status != StatusDown {
		t.Errorf("expected StatusDown, got %s", checks[0].Status)
	}
}

func TestOverall_AllOK(t *testing.T) {
	m := New()
	m.Record("a", StatusOK, "")
	m.Record("b", StatusOK, "")
	if m.Overall() != StatusOK {
		t.Errorf("expected ok, got %s", m.Overall())
	}
}

func TestOverall_AnyDegraded(t *testing.T) {
	m := New()
	m.Record("a", StatusOK, "")
	m.Record("b", StatusDegraded, "slow")
	if m.Overall() != StatusDegraded {
		t.Errorf("expected degraded, got %s", m.Overall())
	}
}

func TestOverall_AnyDown(t *testing.T) {
	m := New()
	m.Record("a", StatusDegraded, "")
	m.Record("b", StatusDown, "crashed")
	if m.Overall() != StatusDown {
		t.Errorf("expected down, got %s", m.Overall())
	}
}

func TestOverall_Empty(t *testing.T) {
	m := New()
	if m.Overall() != StatusOK {
		t.Errorf("expected ok for empty monitor, got %s", m.Overall())
	}
}

func TestCheck_String(t *testing.T) {
	c := Check{Name: "watcher", Status: StatusOK, Message: "running"}
	expected := "[ok] watcher: running"
	if c.String() != expected {
		t.Errorf("got %q, want %q", c.String(), expected)
	}
}
