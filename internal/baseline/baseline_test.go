package baseline

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/scanner"
)

func ports(addrs ...string) []scanner.Port {
	var ps []scanner.Port
	for _, a := range addrs {
		ps = append(ps, scanner.Port{Protocol: "tcp", Address: a})
	}
	return ps
}

func TestNew_Empty(t *testing.T) {
	b := New()
	if b == nil {
		t.Fatal("expected non-nil baseline")
	}
	if got := b.Deviations(ports(":80")); len(got) != 1 {
		t.Fatalf("expected 1 deviation from empty baseline, got %d", len(got))
	}
}

func TestSet_And_Deviations(t *testing.T) {
	b := New()
	b.Set(ports(":80", ":443"))

	current := ports(":80", ":443", ":8080")
	dev := b.Deviations(current)
	if len(dev) != 1 {
		t.Fatalf("expected 1 deviation, got %d", len(dev))
	}
	if dev[0].Address != ":8080" {
		t.Errorf("unexpected deviation: %v", dev[0].Address)
	}
}

func TestDeviations_NoChange(t *testing.T) {
	b := New()
	b.Set(ports(":80", ":443"))
	dev := b.Deviations(ports(":80", ":443"))
	if len(dev) != 0 {
		t.Fatalf("expected no deviations, got %d", len(dev))
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "baseline.json")

	b := New()
	b.Set(ports(":22", ":80"))
	if err := b.Save(path); err != nil {
		t.Fatalf("save: %v", err)
	}

	b2, err := Load(path)
	if err != nil {
		t.Fatalf("load: %v", err)
	}
	dev := b2.Deviations(ports(":22", ":80"))
	if len(dev) != 0 {
		t.Fatalf("expected no deviations after round-trip, got %d", len(dev))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "missing.json"))
	if !os.IsNotExist(err) {
		t.Fatalf("expected not-exist error, got %v", err)
	}
}
