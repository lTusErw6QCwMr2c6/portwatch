package eventbridge_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventbridge"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Port: port, Proto: "tcp"},
	}
}

func TestNew_NotNil(t *testing.T) {
	b := eventbridge.New()
	if b == nil {
		t.Fatal("expected non-nil bridge")
	}
}

func TestRegisterSource_And_Sources(t *testing.T) {
	b := eventbridge.New()
	if err := b.RegisterSource("scanner"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srcs := b.Sources(); len(srcs) != 1 || srcs[0] != "scanner" {
		t.Fatalf("unexpected sources: %v", srcs)
	}
}

func TestRegisterSource_Duplicate_ReturnsError(t *testing.T) {
	b := eventbridge.New()
	_ = b.RegisterSource("scanner")
	if err := b.RegisterSource("scanner"); err == nil {
		t.Fatal("expected error for duplicate source")
	}
}

func TestSubscribe_UnknownSource_ReturnsError(t *testing.T) {
	b := eventbridge.New()
	if err := b.Subscribe("ghost", func(alert.Event) {}); err == nil {
		t.Fatal("expected error for unknown source")
	}
}

func TestEmit_DeliverstToSubscribers(t *testing.T) {
	b := eventbridge.New()
	_ = b.RegisterSource("scanner")

	var got []alert.Event
	_ = b.Subscribe("scanner", func(e alert.Event) { got = append(got, e) })
	_ = b.Subscribe("scanner", func(e alert.Event) { got = append(got, e) })

	ev := makeEvent(8080)
	if err := b.Emit("scanner", ev); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 deliveries, got %d", len(got))
	}
}

func TestEmit_UnknownSource_ReturnsError(t *testing.T) {
	b := eventbridge.New()
	if err := b.Emit("missing", makeEvent(22)); err == nil {
		t.Fatal("expected error emitting to unknown source")
	}
}

func TestEmit_NoSubscribers_NoError(t *testing.T) {
	b := eventbridge.New()
	_ = b.RegisterSource("scanner")
	if err := b.Emit("scanner", makeEvent(443)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
