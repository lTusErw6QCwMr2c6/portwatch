package sampling

import (
	"testing"
)

func TestAllow_None_AlwaysPasses(t *testing.T) {
	s := New(Config{Strategy: StrategyNone}, nil)
	for i := 0; i < 10; i++ {
		if !s.Allow("k") {
			t.Fatal("expected all events to pass with StrategyNone")
		}
	}
}

func TestAllow_Every_PassesEveryN(t *testing.T) {
	s := New(Config{Strategy: StrategyEvery, Rate: 3}, nil)
	results := make([]bool, 6)
	for i := 0; i < 6; i++ {
		results[i] = s.Allow("k")
	}
	// indices 2 and 5 should be true (every 3rd)
	for i, want := range []bool{false, false, true, false, false, true} {
		if results[i] != want {
			t.Errorf("index %d: got %v want %v", i, results[i], want)
		}
	}
}

func TestAllow_Every_DifferentKeysAreIndependent(t *testing.T) {
	s := New(Config{Strategy: StrategyEvery, Rate: 2}, nil)
	if s.Allow("a") {
		t.Error("first call for 'a' should not pass")
	}
	if s.Allow("b") {
		t.Error("first call for 'b' should not pass")
	}
	if !s.Allow("a") {
		t.Error("second call for 'a' should pass")
	}
}

func TestAllow_Random_RespectsRate(t *testing.T) {
	calls := 0
	s := New(Config{Strategy: StrategyRandom, Rate: 0.5}, func() float64 {
		calls++
		if calls%2 == 0 {
			return 0.3 // < 0.5 → allow
		}
		return 0.7 // >= 0.5 → deny
	})
	if !s.Allow("x") {
		t.Error("expected allow on even call")
	}
	if s.Allow("x") {
		t.Error("expected deny on odd call")
	}
}

func TestReset_ClearsCounters(t *testing.T) {
	s := New(Config{Strategy: StrategyEvery, Rate: 2}, nil)
	s.Allow("k") // counter = 1
	s.Reset()
	if s.Allow("k") {
		t.Error("after reset counter should restart; first call should not pass")
	}
}
