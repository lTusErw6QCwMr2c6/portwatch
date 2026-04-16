package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestAllow_UnderLimit(t *testing.T) {
	th := throttle.New(time.Second, 3)
	for i := 0; i < 3; i++ {
		if !th.Allow("key") {
			t.Fatalf("expected allow on hit %d", i+1)
		}
	}
}

func TestAllow_ExceedsLimit(t *testing.T) {
	th := throttle.New(time.Second, 2)
	th.Allow("key")
	th.Allow("key")
	if th.Allow("key") {
		t.Fatal("expected deny after limit exceeded")
	}
}

func TestAllow_DifferentKeysAreIndependent(t *testing.T) {
	th := throttle.New(time.Second, 1)
	if !th.Allow("a") {
		t.Fatal("expected allow for key a")
	}
	if !th.Allow("b") {
		t.Fatal("expected allow for key b")
	}
	if th.Allow("a") {
		t.Fatal("expected deny for key a on second hit")
	}
}

func TestAllow_AllowedAfterWindowExpires(t *testing.T) {
	th := throttle.New(50*time.Millisecond, 1)
	th.Allow("key")
	time.Sleep(60 * time.Millisecond)
	if !th.Allow("key") {
		t.Fatal("expected allow after window expired")
	}
}

func TestReset_ClearsHits(t *testing.T) {
	th := throttle.New(time.Second, 1)
	th.Allow("key")
	th.Reset()
	if !th.Allow("key") {
		t.Fatal("expected allow after reset")
	}
}

func TestCount_ReturnsRecentHits(t *testing.T) {
	th := throttle.New(time.Second, 5)
	th.Allow("key")
	th.Allow("key")
	if c := th.Count("key"); c != 2 {
		t.Fatalf("expected count 2, got %d", c)
	}
}
