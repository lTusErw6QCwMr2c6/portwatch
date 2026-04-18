package circuitbreaker_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuitbreaker"
)

func TestNew_DefaultsApplied(t *testing.T) {
	b := circuitbreaker.New(0, 0)
	if b == nil {
		t.Fatal("expected non-nil breaker")
	}
	if b.State() != circuitbreaker.StateClosed {
		t.Errorf("expected closed state, got %v", b.State())
	}
}

func TestAllow_ClosedAllowsCalls(t *testing.T) {
	b := circuitbreaker.New(3, time.Second)
	if err := b.Allow(); err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestRecordFailure_OpensAfterThreshold(t *testing.T) {
	b := circuitbreaker.New(3, time.Minute)
	b.RecordFailure()
	b.RecordFailure()
	if b.State() != circuitbreaker.StateClosed {
		t.Error("expected still closed after 2 failures")
	}
	b.RecordFailure()
	if b.State() != circuitbreaker.StateOpen {
		t.Errorf("expected open after threshold, got %v", b.State())
	}
}

func TestAllow_RejectsWhenOpen(t *testing.T) {
	b := circuitbreaker.New(1, time.Minute)
	b.RecordFailure()
	if err := b.Allow(); err != circuitbreaker.ErrOpen {
		t.Errorf("expected ErrOpen, got %v", err)
	}
}

func TestAllow_HalfOpenAfterCooldown(t *testing.T) {
	b := circuitbreaker.New(1, 10*time.Millisecond)
	b.RecordFailure()
	time.Sleep(20 * time.Millisecond)
	if err := b.Allow(); err != nil {
		t.Errorf("expected nil after cooldown, got %v", err)
	}
	if b.State() != circuitbreaker.StateHalfOpen {
		t.Errorf("expected half-open, got %v", b.State())
	}
}

func TestRecordSuccess_ClosesBreakerAndClearsFailures(t *testing.T) {
	b := circuitbreaker.New(1, time.Minute)
	b.RecordFailure()
	b.RecordSuccess()
	if b.State() != circuitbreaker.StateClosed {
		t.Errorf("expected closed after success, got %v", b.State())
	}
	if err := b.Allow(); err != nil {
		t.Errorf("expected allow after reset, got %v", err)
	}
}

func TestReset_ForcesClosed(t *testing.T) {
	b := circuitbreaker.New(1, time.Minute)
	b.RecordFailure()
	b.Reset()
	if b.State() != circuitbreaker.StateClosed {
		t.Errorf("expected closed after Reset, got %v", b.State())
	}
}
