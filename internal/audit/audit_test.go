package audit_test

import (
	"bufio"
	"encoding/json"
	"os"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/audit"
	"github.com/user/portwatch/internal/scanner"
)

func tempPath(t *testing.T) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "audit-*.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func makeEvent(kind alert.Kind) alert.Event {
	return alert.Event{Kind: kind, Port: scanner.Port{Number: 8080, Protocol: "tcp"}}
}

func TestNew_CreatesFile(t *testing.T) {
	path := tempPath(t)
	w, err := audit.New(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer w.Close()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("file not created: %v", err)
	}
}

func TestWrite_AppendsJSONLine(t *testing.T) {
	path := tempPath(t)
	w, err := audit.New(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := w.Write(makeEvent(alert.Opened)); err != nil {
		t.Fatalf("write error: %v", err)
	}
	w.Close()

	f, _ := os.Open(path)
	defer f.Close()
	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		t.Fatal("expected at least one line")
	}
	var entry audit.Entry
	if err := json.Unmarshal(scanner.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.Event.Port.Number != 8080 {
		t.Errorf("expected port 8080, got %d", entry.Event.Port.Number)
	}
}

func TestWrite_MultipleEntries(t *testing.T) {
	path := tempPath(t)
	w, _ := audit.New(path)
	w.Write(makeEvent(alert.Opened))
	w.Write(makeEvent(alert.Closed))
	w.Close()

	f, _ := os.Open(path)
	defer f.Close()
	lines := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		lines++
	}
	if lines != 2 {
		t.Errorf("expected 2 lines, got %d", lines)
	}
}
