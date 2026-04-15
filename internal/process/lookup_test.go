package process

import (
	"testing"
)

func TestInfo_String_WithProcess(t *testing.T) {
	info := Info{PID: 1234, Name: "nginx"}
	got := info.String()
	want := "nginx(1234)"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestInfo_String_Unknown(t *testing.T) {
	info := Info{}
	got := info.String()
	want := "unknown"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestParseOutput_ValidOutput(t *testing.T) {
	input := "p4242\ncnginx\nf10\n"
	info, err := parseOutput(input)
	if err != nil {
		t.Fatalf("parseOutput() unexpected error: %v", err)
	}
	if info.PID != 4242 {
		t.Errorf("PID = %d, want 4242", info.PID)
	}
	if info.Name != "nginx" {
		t.Errorf("Name = %q, want %q", info.Name, "nginx")
	}
}

func TestParseOutput_EmptyOutput(t *testing.T) {
	info, err := parseOutput("")
	if err != nil {
		t.Fatalf("parseOutput() unexpected error: %v", err)
	}
	if info.PID != 0 || info.Name != "" {
		t.Errorf("expected empty Info, got %+v", info)
	}
}

func TestParseOutput_InvalidPID(t *testing.T) {
	_, err := parseOutput("pabc\ncnginx\n")
	if err == nil {
		t.Error("parseOutput() expected error for invalid PID, got nil")
	}
}

func TestLookup_InvalidPort(t *testing.T) {
	// Port 0 should return empty info without error (lsof finds nothing).
	info, err := Lookup(0)
	if err != nil {
		t.Fatalf("Lookup(0) unexpected error: %v", err)
	}
	// We cannot assert a specific process on port 0; just verify the call succeeds.
	_ = info
}
