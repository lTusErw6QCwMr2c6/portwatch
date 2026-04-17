package summary_test

import (
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/summary"
)

func makeMetrics(cycles, opened, closed, alerts int) *metrics.Metrics {
	m := metrics.New()
	for i := 0; i < cycles; i++ {
		m.RecordCycle()
	}
	for i := 0; i < opened; i++ {
		m.RecordOpened()
	}
	for i := 0; i < closed; i++ {
		m.RecordClosed()
	}
	for i := 0; i < alerts; i++ {
		m.RecordAlert()
	}
	return m
}

func TestBuild_PopulatesReport(t *testing.T) {
	m := makeMetrics(3, 5, 2, 1)
	b := summary.New(m, 5)
	r := b.Build(nil)
	if r.Cycles != 3 {
		t.Errorf("expected 3 cycles, got %d", r.Cycles)
	}
	if r.Opened != 5 {
		t.Errorf("expected 5 opened, got %d", r.Opened)
	}
	if r.Closed != 2 {
		t.Errorf("expected 2 closed, got %d", r.Closed)
	}
	if r.Alerts != 1 {
		t.Errorf("expected 1 alert, got %d", r.Alerts)
	}
}

func TestBuild_CapsTopEvents(t *testing.T) {
	m := makeMetrics(1, 0, 0, 0)
	events := []alert.Event{
		{Port: scanner.Port{Number: 80}, Kind: alert.Opened},
		{Port: scanner.Port{Number: 443}, Kind: alert.Opened},
		{Port: scanner.Port{Number: 8080}, Kind: alert.Closed},
		{Port: scanner.Port{Number: 3000}, Kind: alert.Opened},
	}
	b := summary.New(m, 2)
	r := b.Build(events)
	if len(r.TopEvents) != 2 {
		t.Errorf("expected 2 top events, got %d", len(r.TopEvents))
	}
}

func TestFormat_ContainsKeyFields(t *testing.T) {
	m := makeMetrics(2, 4, 1, 0)
	b := summary.New(m, 5)
	r := b.Build(nil)
	out := r.Format()
	for _, needle := range []string{"Summary", "Cycles", "Opened", "Closed"} {
		if !strings.Contains(out, needle) {
			t.Errorf("expected %q in output", needle)
		}
	}
}

func TestNew_DefaultMaxTop(t *testing.T) {
	m := makeMetrics(0, 0, 0, 0)
	b := summary.New(m, 0)
	if b == nil {
		t.Fatal("expected non-nil builder")
	}
}
