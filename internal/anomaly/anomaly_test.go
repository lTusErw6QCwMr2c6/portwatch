package anomaly

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(proto string, port int) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Protocol: proto, Number: port},
		Action: alert.ActionOpened,
	}
}

func TestNew_NotNil(t *testing.T) {
	d := New(time.Second, 3)
	if d == nil {
		t.Fatal("expected non-nil detector")
	}
}

func TestEvaluate_BelowThreshold_NoAnomaly(t *testing.T) {
	d := New(time.Second, 5)
	e := makeEvent("tcp", 8080)
	_, ok := d.Evaluate(e)
	if ok {
		t.Error("expected no anomaly below threshold")
	}
}

func TestEvaluate_AtThreshold_DetectsAnomaly(t *testing.T) {
	d := New(time.Second, 3)
	e := makeEvent("tcp", 9090)
	var det Detection
	var ok bool
	for i := 0; i < 3; i++ {
		det, ok = d.Evaluate(e)
	}
	if !ok {
		t.Fatal("expected anomaly at threshold")
	}
	if det.Severity != SeverityMedium {
		t.Errorf("expected medium severity, got %s", det.Severity)
	}
}

func TestEvaluate_DoubleThreshold_HighSeverity(t *testing.T) {
	d := New(time.Second, 3)
	e := makeEvent("tcp", 443)
	var det Detection
	for i := 0; i < 6; i++ {
		det, _ = d.Evaluate(e)
	}
	if det.Severity != SeverityHigh {
		t.Errorf("expected high severity, got %s", det.Severity)
	}
}

func TestReset_ClearsBuckets(t *testing.T) {
	d := New(time.Second, 2)
	e := makeEvent("udp", 53)
	d.Evaluate(e)
	d.Evaluate(e)
	d.Reset()
	_, ok := d.Evaluate(e)
	if ok {
		t.Error("expected no anomaly after reset")
	}
}

func TestDetection_String(t *testing.T) {
	det := Detection{
		Event:    makeEvent("tcp", 22),
		Reason:   "3 events on tcp:22 within 1s",
		Severity: SeverityLow,
	}
	s := det.String()
	if s == "" {
		t.Error("expected non-empty string")
	}
}
