package eventmerge_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventmerge"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Port: port, Protocol: "tcp"},
	}
}

func TestNew_NotNil(t *testing.T) {
	m := eventmerge.New(16)
	if m == nil {
		t.Fatal("expected non-nil Merger")
	}
}

func TestAdd_RegistersSource(t *testing.T) {
	m := eventmerge.New(16)
	ch := make(chan alert.Event)
	if err := m.Add("src1", ch); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestAdd_EmptyName_ReturnsError(t *testing.T) {
	m := eventmerge.New(16)
	ch := make(chan alert.Event)
	if err := m.Add("", ch); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAdd_DuplicateName_ReturnsError(t *testing.T) {
	m := eventmerge.New(16)
	ch := make(chan alert.Event)
	_ = m.Add("dup", ch)
	if err := m.Add("dup", ch); err == nil {
		t.Fatal("expected error for duplicate name")
	}
}

func TestStart_MergesEvents(t *testing.T) {
	m := eventmerge.New(32)

	ch1 := make(chan alert.Event, 4)
	ch2 := make(chan alert.Event, 4)

	_ = m.Add("a", ch1)
	_ = m.Add("b", ch2)

	ch1 <- makeEvent(80)
	ch2 <- makeEvent(443)
	close(ch1)
	close(ch2)

	m.Start()

	got := make(map[int]bool)
	timeout := time.After(2 * time.Second)
	for i := 0; i < 2; i++ {
		select {
		case ev := <-m.Out():
			got[ev.Port.Port] = true
		case <-timeout:
			t.Fatal("timed out waiting for merged events")
		}
	}

	if !got[80] || !got[443] {
		t.Errorf("missing expected events; got %v", got)
	}
}

func TestStop_ClosesOutput(t *testing.T) {
	m := eventmerge.New(8)
	ch := make(chan alert.Event)
	_ = m.Add("s", ch)
	m.Start()
	m.Stop()

	select {
	case _, ok := <-m.Out():
		if ok {
			t.Fatal("expected output channel to be closed")
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for channel close")
	}
}
