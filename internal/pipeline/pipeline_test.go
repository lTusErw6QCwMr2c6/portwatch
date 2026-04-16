package pipeline_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/pipeline"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/scanner"
)

func diff(opened, closed []scanner.Port) scanner.Diff {
	return scanner.Diff{Opened: opened, Closed: closed}
}

func port(n int) scanner.Port {
	return scanner.Port{Number: n, Protocol: "tcp"}
}

func TestRun_NoFilter_ReturnsEvents(t *testing.T) {
	p := pipeline.New(pipeline.Config{})
	events := p.Run(diff([]scanner.Port{port(8080)}, nil))
	if len(events) == 0 {
		t.Fatal("expected events")
	}
}

func TestRun_BelowThreshold_ReturnsNil(t *testing.T) {
	th := alert.Threshold{MinOpened: 5, MinClosed: 5}
	p := pipeline.New(pipeline.Config{Threshold: th})
	events := p.Run(diff([]scanner.Port{port(8080)}, nil))
	if events != nil {
		t.Fatalf("expected nil, got %v", events)
	}
}

func TestRun_WithFilter_ExcludesPort(t *testing.T) {
	f := filter.New([]filter.Rule{{Port: 8080, Protocol: "tcp", Exclude: true}})
	p := pipeline.New(pipeline.Config{Filter: f})
	events := p.Run(diff([]scanner.Port{port(8080)}, nil))
	if len(events) != 0 {
		t.Fatalf("expected 0 events after filter, got %d", len(events))
	}
}

func TestRun_WithRateLimit_SuppressesDuplicate(t *testing.T) {
	rl := ratelimit.New(10 * time.Second)
	p := pipeline.New(pipeline.Config{RateLimit: rl})
	d := diff([]scanner.Port{port(9090)}, nil)
	p.Run(d)
	events := p.Run(d)
	if len(events) != 0 {
		t.Fatalf("expected suppressed duplicate, got %d events", len(events))
	}
}
