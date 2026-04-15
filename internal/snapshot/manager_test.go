package snapshot_test

import (
	"log"
	"os"
	"testing"

	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/snapshot"
)

func silentLogger() *log.Logger {
	return log.New(os.Stderr, "", 0)
}

func TestManager_Reset_ClearsSnapshot(t *testing.T) {
	path := tempFile(t)
	defer os.Remove(path)

	store := snapshot.NewStore(path)
	ports := map[scanner.PortKey]bool{{Port: 9090, Proto: "tcp"}: true}
	if err := store.Save(snapshot.New(ports)); err != nil {
		t.Fatalf("setup Save failed: %v", err)
	}

	m := snapshot.NewManager(store, silentLogger(), 1, 65535)
	if err := m.Reset(); err != nil {
		t.Fatalf("Reset failed: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load after Reset failed: %v", err)
	}
	if len(loaded.Ports) != 0 {
		t.Errorf("expected empty snapshot after Reset, got %d ports", len(loaded.Ports))
	}
}

func TestNewManager_InvalidRange_CycleErrors(t *testing.T) {
	path := tempFile(t)
	defer os.Remove(path)

	store := snapshot.NewStore(path)
	m := snapshot.NewManager(store, silentLogger(), 9999, 1) // invalid: start > end

	_, err := m.Cycle()
	if err == nil {
		t.Error("expected error for invalid port range, got nil")
	}
}

func TestNewManager_NotNil(t *testing.T) {
	store := snapshot.NewStore("/tmp/test-manager.json")
	m := snapshot.NewManager(store, silentLogger(), 1, 1024)
	if m == nil {
		t.Error("expected non-nil Manager")
	}
}
