package jitter

import (
	"testing"
	"time"
)

func TestNew_ClampsFactor(t *testing.T) {
	j := New(-0.5)
	if j.factor != 0 {
		t.Fatalf("expected factor 0, got %v", j.factor)
	}
	j2 := New(2.0)
	if j2.factor != 1 {
		t.Fatalf("expected factor 1, got %v", j2.factor)
	}
}

func TestApply_ZeroFactor_ReturnsBase(t *testing.T) {
	j := New(0)
	base := 100 * time.Millisecond
	if got := j.Apply(base); got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApply_WithFactor_WithinBounds(t *testing.T) {
	j := New(0.2)
	base := 1 * time.Second
	for i := 0; i < 200; i++ {
		got := j.Apply(base)
		low := time.Duration(float64(base) * 0.8)
		high := time.Duration(float64(base) * 1.2)
		if got < low || got > high {
			t.Fatalf("result %v outside [%v, %v]", got, low, high)
		}
	}
}

func TestApplyPositive_NeverReturnsZero(t *testing.T) {
	j := New(1.0)
	base := time.Nanosecond
	for i := 0; i < 100; i++ {
		got := j.ApplyPositive(base)
		if got < time.Nanosecond {
			t.Fatalf("expected >= 1ns, got %v", got)
		}
	}
}

func TestApply_NegativeBase_ReturnsBase(t *testing.T) {
	j := New(0.5)
	base := -1 * time.Second
	if got := j.Apply(base); got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestReset_DoesNotPanic(t *testing.T) {
	j := New(0.3)
	j.Reset()
	if j.rng == nil {
		t.Fatal("rng should not be nil after Reset")
	}
}
