package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/portwatch/internal/checkpoint"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestLoad_ReturnsEmptyWhenAbsent(t *testing.T) {
	s := checkpoint.New(tempPath(t))
	st, err := s.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(st.Ports) != 0 || st.CycleCount != 0 {
		t.Fatalf("expected zero state, got %+v", st)
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	s := checkpoint.New(tempPath(t))
	want := checkpoint.State{
		Ports:      []string{"tcp:8080", "tcp:443"},
		CycleCount: 42,
	}
	if err := s.Save(want); err != nil {
		t.Fatalf("save: %v", err)
	}
	got, err := s.Load()
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	if len(got.Ports) != len(want.Ports) {
		t.Fatalf("ports mismatch: want %v got %v", want.Ports, got.Ports)
	}
	if got.CycleCount != want.CycleCount {
		t.Fatalf("cycle count: want %d got %d", want.CycleCount, got.CycleCount)
	}
	if got.SavedAt.IsZero() {
		t.Fatal("SavedAt should be set")
	}
}

func TestSave_SetsTimestamp(t *testing.T) {
	s := checkpoint.New(tempPath(t))
	before := time.Now()
	_ = s.Save(checkpoint.State{})
	st, _ := s.Load()
	if st.SavedAt.Before(before) {
		t.Fatal("SavedAt should be >= before")
	}
}

func TestRemove_DeletesFile(t *testing.T) {
	p := tempPath(t)
	s := checkpoint.New(p)
	_ = s.Save(checkpoint.State{Ports: []string{"tcp:22"}})
	if err := s.Remove(); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		t.Fatal("file should be deleted")
	}
}

func TestRemove_NoErrorWhenAbsent(t *testing.T) {
	s := checkpoint.New(tempPath(t))
	if err := s.Remove(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
