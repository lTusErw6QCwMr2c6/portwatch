package debounce_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/portwatch/internal/debounce"
)

func TestTrigger_FiresAfterWindow(t *testing.T) {
	var mu sync.Mutex
	fired := []string{}

	d := debounce.New(50*time.Millisecond, func(key string) {
		mu.Lock()
		fired = append(fired, key)
		mu.Unlock()
	})

	d.Trigger("port:8080")
	time.Sleep(100 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if len(fired) != 1 || fired[0] != "port:8080" {
		t.Fatalf("expected callback for port:8080, got %v", fired)
	}
}

func TestTrigger_ResetsOnRepeatedCall(t *testing.T) {
	var mu sync.Mutex
	count := 0

	d := debounce.New(60*time.Millisecond, func(_ string) {
		mu.Lock()
		count++
		mu.Unlock()
	})

	// Trigger three times in quick succession — only one callback expected.
	d.Trigger("port:9090")
	time.Sleep(20 * time.Millisecond)
	d.Trigger("port:9090")
	time.Sleep(20 * time.Millisecond)
	d.Trigger("port:9090")
	time.Sleep(120 * time.Millisecond)

	mu.Lock()
	defer mu.Unlock()
	if count != 1 {
		t.Fatalf("expected 1 callback, got %d", count)
	}
}

func TestCancel_PreventsCallback(t *testing.T) {
	fired := false

	d := debounce.New(50*time.Millisecond, func(_ string) {
		fired = true
	})

	d.Trigger("port:3000")
	d.Cancel("port:3000")
	time.Sleep(100 * time.Millisecond)

	if fired {
		t.Fatal("callback should not have fired after Cancel")
	}
}

func TestPending_ReflectsActiveTimers(t *testing.T) {
	d := debounce.New(200*time.Millisecond, func(_ string) {})

	if d.Pending() != 0 {
		t.Fatalf("expected 0 pending, got %d", d.Pending())
	}

	d.Trigger("a")
	d.Trigger("b")

	if got := d.Pending(); got != 2 {
		t.Fatalf("expected 2 pending, got %d", got)
	}

	d.Cancel("a")
	if got := d.Pending(); got != 1 {
		t.Fatalf("expected 1 pending after cancel, got %d", got)
	}
}
