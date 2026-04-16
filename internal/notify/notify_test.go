package notify

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(kind alert.Kind) alert.Event {
	return alert.Event{
		Kind: kind,
		Port: scanner.Port{Number: 8080, Protocol: "tcp"},
	}
}

func TestStdoutChannel_Send_WritesLine(t *testing.T) {
	var buf bytes.Buffer
	ch := &StdoutChannel{Writer: &buf, Prefix: "NOTIFY"}

	if err := ch.Send(makeEvent(alert.Opened)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	line := buf.String()
	if !strings.Contains(line, "NOTIFY") {
		t.Errorf("expected prefix in output, got: %s", line)
	}
	if !strings.Contains(line, "8080") {
		t.Errorf("expected port number in output, got: %s", line)
	}
}

func TestDispatcher_Dispatch_CallsAllChannels(t *testing.T) {
	var buf1, buf2 bytes.Buffer
	ch1 := &StdoutChannel{Writer: &buf1, Prefix: "A"}
	ch2 := &StdoutChannel{Writer: &buf2, Prefix: "B"}

	d := NewDispatcher(ch1, ch2)
	errs := d.Dispatch(makeEvent(alert.Opened))
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if buf1.Len() == 0 || buf2.Len() == 0 {
		t.Error("expected both channels to receive the event")
	}
}

type failChannel struct{}

func (f *failChannel) Send(_ alert.Event) error {
	return errors.New("send failed")
}

func TestDispatcher_Dispatch_CollectsErrors(t *testing.T) {
	d := NewDispatcher(&failChannel{}, &failChannel{})
	errs := d.Dispatch(makeEvent(alert.Closed))
	if len(errs) != 2 {
		t.Errorf("expected 2 errors, got %d", len(errs))
	}
}

func TestDispatcher_DispatchAll_SendsAllEvents(t *testing.T) {
	var buf bytes.Buffer
	ch := &StdoutChannel{Writer: &buf, Prefix: "X"}
	d := NewDispatcher(ch)

	events := []alert.Event{makeEvent(alert.Opened), makeEvent(alert.Closed)}
	errs := d.DispatchAll(events)
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	lines := strings.Count(buf.String(), "\n")
	if lines != 2 {
		t.Errorf("expected 2 lines, got %d", lines)
	}
}
