package metrics

import (
	"testing"
	"time"
)

func TestNew_ZeroValues(t *testing.T) {
	c := New()
	s := c.Read()
	if s.Cycles != 0 || s.Opened != 0 || s.Closed != 0 || s.Alerts != 0 {
		t.Fatalf("expected zero snapshot, got %+v", s)
	}
}

func TestRecordCycle_IncrementsAndTimestamps(t *testing.T) {
	c := New()
	before := time.Now()
	c.RecordCycle()
	c.RecordCycle()
	s := c.Read()
	if s.Cycles != 2 {
		t.Fatalf("expected 2 cycles, got %d", s.Cycles)
	}
	if s.LastCycleAt.Before(before) {
		t.Fatalf("LastCycleAt not updated")
	}
}

func TestRecordOpened_Accumulates(t *testing.T) {
	c := New()
	c.RecordOpened(3)
	c.RecordOpened(2)
	if got := c.Read().Opened; got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestRecordClosed_Accumulates(t *testing.T) {
	c := New()
	c.RecordClosed(1)
	c.RecordClosed(4)
	if got := c.Read().Closed; got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestRecordAlert_Increments(t *testing.T) {
	c := New()
	c.RecordAlert()
	c.RecordAlert()
	c.RecordAlert()
	if got := c.Read().Alerts; got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	c := New()
	c.RecordCycle()
	c.RecordOpened(10)
	c.RecordClosed(5)
	c.RecordAlert()
	c.Reset()
	s := c.Read()
	if s.Cycles != 0 || s.Opened != 0 || s.Closed != 0 || s.Alerts != 0 {
		t.Fatalf("expected zero after reset, got %+v", s)
	}
}
