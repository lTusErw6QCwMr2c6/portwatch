package snapshot_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func tempFile(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp("", "snapshot-*.json")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	f.Close()
	os.Remove(f.Name()) // start with no file so Load returns empty
	return f.Name()
}

func TestLoad_ReturnsEmptyWhenFileAbsent(t *testing.T) {
	store := snapshot.NewStore("/tmp/portwatch-nonexistent-test.json")
	snap, err := store.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(snap.Ports) != 0 {
		t.Errorf("expected empty ports, got %d", len(snap.Ports))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tempFile(t)
	defer os.Remove(path)

	store := snapshot.NewStore(path)
	ports := map[scanner.PortKey]bool{
		{Port: 80, Proto: "tcp"}:   true,
		{Port: 443, Proto: "tcp"}: true,
	}
	orig := snapshot.Snapshot{Timestamp: time.Now().Truncate(time.Second), Ports: ports}

	if err := store.Save(orig); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if len(loaded.Ports) != len(orig.Ports) {
		t.Errorf("expected %d ports, got %d", len(orig.Ports), len(loaded.Ports))
	}
	for k := range orig.Ports {
		if !loaded.Ports[k] {
			t.Errorf("missing port %v in loaded snapshot", k)
		}
	}
}

func TestNew_SetsTimestamp(t *testing.T) {
	before := time.Now()
	snap := snapshot.New(map[scanner.PortKey]bool{})
	after := time.Now()

	if snap.Timestamp.Before(before) || snap.Timestamp.After(after) {
		t.Errorf("timestamp %v not within expected range", snap.Timestamp)
	}
}

func TestSave_OverwritesPreviousSnapshot(t *testing.T) {
	path := tempFile(t)
	defer os.Remove(path)

	store := snapshot.NewStore(path)

	first := snapshot.New(map[scanner.PortKey]bool{{Port: 22, Proto: "tcp"}: true})
	if err := store.Save(first); err != nil {
		t.Fatalf("first Save failed: %v", err)
	}

	second := snapshot.New(map[scanner.PortKey]bool{{Port: 8080, Proto: "tcp"}: true})
	if err := store.Save(second); err != nil {
		t.Fatalf("second Save failed: %v", err)
	}

	loaded, _ := store.Load()
	if loaded.Ports[scanner.PortKey{Port: 22, Proto: "tcp"}] {
		t.Error("old port should not be present after overwrite")
	}
	if !loaded.Ports[scanner.PortKey{Port: 8080, Proto: "tcp"}] {
		t.Error("new port should be present after overwrite")
	}
}
