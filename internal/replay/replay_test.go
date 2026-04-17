package replay

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
)

func writeEntries(t *testing.T, entries []Entry) string {
	t.Helper()
	f, err := os.CreateTemp("", "replay-*.log")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	for _, e := range entries {
		if err := enc.Encode(e); err != nil {
			t.Fatal(err)
		}
	}
	return f.Name()
}

func makeEntry(ts time.Time) Entry {
	return Entry{
		Timestamp: ts,
		Event: alert.Event{
			Kind: alert.Opened,
			Port: scanner.Port{Number: 8080, Protocol: "tcp"},
		},
	}
}

func TestRead_ReturnsAllEntries(t *testing.T) {
	now := time.Now()
	path := writeEntries(t, []Entry{makeEntry(now), makeEntry(now.Add(time.Second))})
	defer os.Remove(path)

	entries, err := New(path).Read(Options{})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestRead_SinceFilter(t *testing.T) {
	now := time.Now()
	path := writeEntries(t, []Entry{makeEntry(now), makeEntry(now.Add(2 * time.Second))})
	defer os.Remove(path)

	entries, err := New(path).Read(Options{Since: now.Add(time.Second)})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
}

func TestRead_LimitRespected(t *testing.T) {
	now := time.Now()
	path := writeEntries(t, []Entry{makeEntry(now), makeEntry(now), makeEntry(now)})
	defer os.Remove(path)

	entries, err := New(path).Read(Options{Limit: 2})
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestRead_MissingFile(t *testing.T) {
	_, err := New("/nonexistent/path.log").Read(Options{})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}
