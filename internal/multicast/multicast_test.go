package multicast_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/multicast"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent() alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Port: 8080, Protocol: "tcp"},
	}
}

func TestNew_NotNil(t *testing.T) {
	b := multicast.New(4)
	if b == nil {
		t.Fatal("expected non-nil bus")
	}
}

func TestSubscribe_ReturnsChannel(t *testing.T) {
	b := multicast.New(4)
	ch, err := b.Subscribe("a")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ch == nil {
		t.Fatal("expected non-nil channel")
	}
}

func TestSubscribe_DuplicateName_ReturnsError(t *testing.T) {
	b := multicast.New(4)
	b.Subscribe("a")
	_, err := b.Subscribe("a")
	if err == nil {
		t.Fatal("expected error for duplicate subscriber")
	}
}

func TestPublish_DeliverstToAllSubscribers(t *testing.T) {
	b := multicast.New(4)
	ch1, _ := b.Subscribe("s1")
	ch2, _ := b.Subscribe("s2")
	ev := makeEvent()
	n := b.Publish(ev)
	if n != 2 {
		t.Fatalf("expected 2 deliveries, got %d", n)
	}
	for _, ch := range []<-chan alert.Event{ch1, ch2} {
		select {
		case got := <-ch:
			if got.Port.Port != ev.Port.Port {
				t.Errorf("wrong event received")
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("timed out waiting for event")
		}
	}
}

func TestUnsubscribe_RemovesSubscriber(t *testing.T) {
	b := multicast.New(4)
	b.Subscribe("x")
	if b.Len() != 1 {
		t.Fatal("expected 1 subscriber")
	}
	b.Unsubscribe("x")
	if b.Len() != 0 {
		t.Fatal("expected 0 subscribers after unsubscribe")
	}
}

func TestPublish_FullBuffer_DropsEvent(t *testing.T) {
	b := multicast.New(1)
	b.Subscribe("slow")
	b.Publish(makeEvent()) // fills buffer
	n := b.Publish(makeEvent()) // should drop
	if n != 0 {
		t.Errorf("expected 0 deliveries to full buffer, got %d", n)
	}
}

func TestUnsubscribe_NonExistent_NoError(t *testing.T) {
	b := multicast.New(4)
	// Unsubscribing a name that was never registered should not panic or affect state.
	b.Unsubscribe("ghost")
	if b.Len() != 0 {
		t.Fatalf("expected 0 subscribers, got %d", b.Len())
	}
}
