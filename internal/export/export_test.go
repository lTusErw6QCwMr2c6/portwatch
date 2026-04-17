package export_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/export"
	"github.com/user/portwatch/internal/process"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(t alert.EventType, port int) alert.Event {
	return alert.Event{
		Timestamp: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Type:      t,
		Port: scanner.Port{
			Port:     port,
			Protocol: "tcp",
			Process:  process.Info{Name: "nginx", PID: 42},
		},
	}
}

func TestWrite_JSON_ValidOutput(t *testing.T) {
	var buf bytes.Buffer
	ex := export.New(&buf, export.FormatJSON)
	events := []alert.Event{makeEvent(alert.EventOpened, 80)}
	if err := ex.Write(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var records []export.Record
	if err := json.Unmarshal(buf.Bytes(), &records); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	if records[0].Port != 80 {
		t.Errorf("expected port 80, got %d", records[0].Port)
	}
	if records[0].Process != "nginx" {
		t.Errorf("expected process nginx, got %s", records[0].Process)
	}
}

func TestWrite_CSV_ContainsHeader(t *testing.T) {
	var buf bytes.Buffer
	ex := export.New(&buf, export.FormatCSV)
	events := []alert.Event{makeEvent(alert.EventClosed, 443)}
	if err := ex.Write(events); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(lines[0], "timestamp") {
		t.Errorf("expected header row, got: %s", lines[0])
	}
	if !strings.Contains(lines[1], "443") {
		t.Errorf("expected port 443 in data row, got: %s", lines[1])
	}
}

func TestWrite_UnknownFormat_ReturnsError(t *testing.T) {
	var buf bytes.Buffer
	ex := export.New(&buf, export.Format("xml"))
	if err := ex.Write(nil); err == nil {
		t.Error("expected error for unknown format")
	}
}

func TestWrite_Empty_NoError(t *testing.T) {
	var buf bytes.Buffer
	ex := export.New(&buf, export.FormatJSON)
	if err := ex.Write([]alert.Event{}); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
