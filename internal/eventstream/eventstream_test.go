package eventstream_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventstream"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Number: port, Protocol: "tcp"},
	}
}

func TestNew_NotNil(t *testing.T) {
	s := eventstream.New(4)
	if s == nil {
		t.Fatal("expected non-nil stream")
	}
}

func TestSubscribe_ReturnsChannel(t *testing.T) {
	s := eventstream.New(4)
	ch, err := s.Subscribe("a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
	if s.Len() != 1 {
		t.Fatalf("expected 1 consumer, got %d", s.Len())
	}
}

func TestSubscribe_Duplicate_ReturnsError(t *testing.T) {
	s := eventstream.New(4)
	_, _ = s.Subscribe("a")
	_, err := s.Subscribe("a")
	if err != eventstream.ErrDuplicateConsumer {
		t.Fatalf("expected ErrDuplicateConsumer, got %v", err)
	}
}

func TestPublish_DeliverstToAllSubscribers(t *testing.T) {
	s := eventstream.New(4)
	ch1, _ := s.Subscribe("one")
	ch2, _ := s.Subscribe("two")

	ev := makeEvent(8080)
	s.Publish(ev)

	for _, ch := range []<-chan alert.Event{ch1, ch2} {
		select {
		case got := <-ch:
			if got.Port.Number != 8080 {
				t.Errorf("unexpected port %d", got.Port.Number)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("timed out waiting for event")
		}
	}
}

func TestUnsubscribe_ClosesChannel(t *testing.T) {
	s := eventstream.New(4)
	ch, _ := s.Subscribe("x")
	if err := s.Unsubscribe("x"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 0 {
		t.Fatalf("expected 0 consumers, got %d", s.Len())
	}
	// channel should be closed
	select {
	case _, ok := <-ch:
		if ok {
			t.Error("expected channel to be closed")
		}
	default:
		t.Error("expected closed channel to be readable")
	}
}

func TestUnsubscribe_UnknownConsumer_ReturnsError(t *testing.T) {
	s := eventstream.New(4)
	err := s.Unsubscribe("missing")
	if err != eventstream.ErrUnknownConsumer {
		t.Fatalf("expected ErrUnknownConsumer, got %v", err)
	}
}

func TestClose_RemovesAllConsumers(t *testing.T) {
	s := eventstream.New(4)
	_, _ = s.Subscribe("a")
	_, _ = s.Subscribe("b")
	s.Close()
	if s.Len() != 0 {
		t.Fatalf("expected 0 consumers after Close, got %d", s.Len())
	}
}

func TestPublish_FullBuffer_DropsEvent(t *testing.T) {
	s := eventstream.New(1)
	ch, _ := s.Subscribe("slow")

	// fill the buffer
	s.Publish(makeEvent(1))
	// this should not block
	s.Publish(makeEvent(2))

	got := <-ch
	if got.Port.Number != 1 {
		t.Errorf("expected first event port 1, got %d", got.Port.Number)
	}
	// second event was dropped; channel should be empty
	select {
	case extra := <-ch:
		t.Errorf("expected empty channel, got port %d", extra.Port.Number)
	default:
	}
}
