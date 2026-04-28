package eventfunnel

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/scanner"
)

func makeEvent(port uint16, proto string) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Number: port, Protocol: proto},
	}
}

func TestNew_NotNil(t *testing.T) {
	f := New(0)
	if f == nil {
		t.Fatal("expected non-nil funnel")
	}
}

func TestAdd_RegistersSource(t *testing.T) {
	f := New(8)
	ch := make(chan alert.Event)
	err := f.Add(Source{Name: "a", Ch: ch})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Len() != 1 {
		t.Fatalf("expected 1 source, got %d", f.Len())
	}
	f.Stop()
}

func TestAdd_EmptyName_ReturnsError(t *testing.T) {
	f := New(8)
	ch := make(chan alert.Event)
	err := f.Add(Source{Name: "", Ch: ch})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
	f.Stop()
}

func TestAdd_DuplicateName_ReturnsError(t *testing.T) {
	f := New(8)
	ch := make(chan alert.Event)
	_ = f.Add(Source{Name: "dup", Ch: ch})
	err := f.Add(Source{Name: "dup", Ch: ch})
	if err == nil {
		t.Fatal("expected error for duplicate name")
	}
	f.Stop()
}

func TestOut_DeliversEventsFromAllSources(t *testing.T) {
	f := New(16)
	ch1 := make(chan alert.Event, 4)
	ch2 := make(chan alert.Event, 4)
	_ = f.Add(Source{Name: "s1", Ch: ch1})
	_ = f.Add(Source{Name: "s2", Ch: ch2})

	ch1 <- makeEvent(80, "tcp")
	ch2 <- makeEvent(443, "tcp")

	collected := map[uint16]bool{}
	timeout := time.After(500 * time.Millisecond)
	for len(collected) < 2 {
		select {
		case ev := <-f.Out():
			collected[ev.Port.Number] = true
		case <-timeout:
			t.Fatal("timed out waiting for events")
		}
	}
	f.Stop()
}

func TestDrainAsync_ReceivesAllEvents(t *testing.T) {
	f := New(16)
	ch := make(chan alert.Event, 4)
	_ = f.Add(Source{Name: "src", Ch: ch})

	var got []alert.Event
	results := make(chan alert.Event, 4)
	done := DrainAsync(f, func(ev alert.Event) { results <- ev })

	ch <- makeEvent(22, "tcp")
	ch <- makeEvent(53, "udp")

	timeout := time.After(500 * time.Millisecond)
	for len(got) < 2 {
		select {
		case ev := <-results:
			got = append(got, ev)
		case <-timeout:
			t.Fatal("timed out")
		}
	}
	f.Stop()
	<-done
}
