package eventpipeline_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventpipeline"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Port: port, Protocol: "tcp"},
	}
}

func TestNew_NotNil(t *testing.T) {
	p := eventpipeline.New()
	if p == nil {
		t.Fatal("expected non-nil pipeline")
	}
}

func TestRegister_AddsStage(t *testing.T) {
	p := eventpipeline.New()
	err := p.Register(eventpipeline.Stage{
		Name:    "noop",
		Process: func(e []alert.Event) []alert.Event { return e },
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Len() != 1 {
		t.Fatalf("expected 1 stage, got %d", p.Len())
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	p := eventpipeline.New()
	err := p.Register(eventpipeline.Stage{
		Process: func(e []alert.Event) []alert.Event { return e },
	})
	if err == nil {
		t.Fatal("expected error for empty stage name")
	}
}

func TestRegister_DuplicateName_ReturnsError(t *testing.T) {
	p := eventpipeline.New()
	stage := eventpipeline.Stage{
		Name:    "filter",
		Process: func(e []alert.Event) []alert.Event { return e },
	}
	_ = p.Register(stage)
	if err := p.Register(stage); err == nil {
		t.Fatal("expected error for duplicate stage name")
	}
}

func TestRegister_NilProcess_ReturnsError(t *testing.T) {
	p := eventpipeline.New()
	err := p.Register(eventpipeline.Stage{Name: "bad"})
	if err == nil {
		t.Fatal("expected error for nil Process")
	}
}

func TestRun_NoStages_ReturnsAll(t *testing.T) {
	p := eventpipeline.New()
	events := []alert.Event{makeEvent(80), makeEvent(443)}
	out := p.Run(events)
	if len(out) != 2 {
		t.Fatalf("expected 2 events, got %d", len(out))
	}
}

func TestRun_StagesExecuteInOrder(t *testing.T) {
	p := eventpipeline.New()
	var order []string

	_ = p.Register(eventpipeline.Stage{
		Name: "first",
		Process: func(e []alert.Event) []alert.Event {
			order = append(order, "first")
			return e
		},
	})
	_ = p.Register(eventpipeline.Stage{
		Name: "second",
		Process: func(e []alert.Event) []alert.Event {
			order = append(order, "second")
			return e
		},
	})

	p.Run([]alert.Event{makeEvent(22)})

	if len(order) != 2 || order[0] != "first" || order[1] != "second" {
		t.Fatalf("unexpected execution order: %v", order)
	}
}

func TestRun_EarlyExitOnEmptySlice(t *testing.T) {
	p := eventpipeline.New()
	called := false

	_ = p.Register(eventpipeline.Stage{
		Name:    "drain",
		Process: func(e []alert.Event) []alert.Event { return nil },
	})
	_ = p.Register(eventpipeline.Stage{
		Name: "should-not-run",
		Process: func(e []alert.Event) []alert.Event {
			called = true
			return e
		},
	})

	p.Run([]alert.Event{makeEvent(8080)})

	if called {
		t.Fatal("expected second stage to be skipped after empty result")
	}
}

func TestStages_ReturnsNamesInOrder(t *testing.T) {
	p := eventpipeline.New()
	_ = p.Register(eventpipeline.Stage{Name: "a", Process: func(e []alert.Event) []alert.Event { return e }})
	_ = p.Register(eventpipeline.Stage{Name: "b", Process: func(e []alert.Event) []alert.Event { return e }})
	names := p.Stages()
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Fatalf("unexpected stage names: %v", names)
	}
}
