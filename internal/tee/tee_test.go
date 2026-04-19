package tee_test

import (
	"errors"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/tee"
)

type recordSink struct {
	events []alert.Event
	err    error
}

func (r *recordSink) Receive(e alert.Event) error {
	r.events = append(r.events, e)
	return r.err
}

func makeEvent() alert.Event {
	return alert.Event{Type: alert.Opened, Port: 8080, Protocol: "tcp"}
}

func TestSend_DeliveriesToAllSinks(t *testing.T) {
	a, b := &recordSink{}, &recordSink{}
	tee := tee.New(false, a, b)

	if err := tee.Send(makeEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(a.events) != 1 || len(b.events) != 1 {
		t.Fatalf("expected 1 event per sink, got %d / %d", len(a.events), len(b.events))
	}
}

func TestSend_ContinuesOnError_WhenStopOnErrFalse(t *testing.T) {
	a := &recordSink{err: errors.New("boom")}
	b := &recordSink{}
	tee := tee.New(false, a, b)

	err := tee.Send(makeEvent())
	if err == nil {
		t.Fatal("expected error")
	}
	if len(b.events) != 1 {
		t.Fatal("second sink should still receive event")
	}
}

func TestSend_StopsOnError_WhenStopOnErrTrue(t *testing.T) {
	a := &recordSink{err: errors.New("boom")}
	b := &recordSink{}
	tee := tee.New(true, a, b)

	if err := tee.Send(makeEvent()); err == nil {
		t.Fatal("expected error")
	}
	if len(b.events) != 0 {
		t.Fatal("second sink should not receive event after stop-on-err")
	}
}

func TestAdd_IncreasesSinkCount(t *testing.T) {
	tee := tee.New(false)
	if tee.Len() != 0 {
		t.Fatalf("expected 0, got %d", tee.Len())
	}
	tee.Add(&recordSink{})
	if tee.Len() != 1 {
		t.Fatalf("expected 1, got %d", tee.Len())
	}
}

func TestSend_NoSinks_ReturnsNil(t *testing.T) {
	tee := tee.New(false)
	if err := tee.Send(makeEvent()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
