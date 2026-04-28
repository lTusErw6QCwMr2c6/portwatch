package eventfunnel

import (
	"testing"
	"time"

	"github.com/example/portwatch/internal/alert"
	"github.com/example/portwatch/internal/scanner"
)

func TestDrain_CallsHandlerForEachEvent(t *testing.T) {
	f := New(8)
	ch := make(chan alert.Event, 4)
	_ = f.Add(Source{Name: "drain-src", Ch: ch})

	ch <- alert.Event{Type: alert.Opened, Port: scanner.Port{Number: 8080, Protocol: "tcp"}}
	ch <- alert.Event{Type: alert.Closed, Port: scanner.Port{Number: 9090, Protocol: "tcp"}}

	var count int
	done := make(chan struct{})
	go func() {
		defer close(done)
		Drain(f, func(_ alert.Event) { count++ })
	}()

	timeout := time.After(300 * time.Millisecond)
	for count < 2 {
		select {
		case <-timeout:
			t.Fatal("timed out waiting for drain")
		case <-time.After(10 * time.Millisecond):
		}
	}
	f.Stop()
	<-done
}

func TestDrainAsync_DoneClosedAfterStop(t *testing.T) {
	f := New(4)
	ch := make(chan alert.Event)
	_ = f.Add(Source{Name: "x", Ch: ch})

	done := DrainAsync(f, func(_ alert.Event) {})
	f.Stop()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("done channel not closed after Stop")
	}
}
