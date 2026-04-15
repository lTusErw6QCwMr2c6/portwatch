package filter_test

import (
	"testing"

	"github.com/user/portwatch/internal/filter"
	"github.com/user/portwatch/internal/scanner"
)

func ports(nums ...uint16) []scanner.Port {
	out := make([]scanner.Port, len(nums))
	for i, n := range nums {
		out[i] = scanner.Port{Number: n, Protocol: "tcp"}
	}
	return out
}

func TestApply_NoRules_ReturnsAll(t *testing.T) {
	f := filter.New(nil)
	input := ports(80, 443, 8080)
	got := f.Apply(input)
	if len(got) != len(input) {
		t.Fatalf("expected %d ports, got %d", len(input), len(got))
	}
}

func TestApply_ExcludeRule_RemovesPorts(t *testing.T) {
	f := filter.New([]filter.Rule{
		{PortStart: 80, PortEnd: 80, Protocols: []string{"tcp"}, Exclude: true},
	})
	got := f.Apply(ports(80, 443))
	if len(got) != 1 || got[0].Number != 443 {
		t.Fatalf("expected only port 443, got %v", got)
	}
}

func TestApply_IncludeRule_AllowsOnlyMatching(t *testing.T) {
	f := filter.New([]filter.Rule{
		{PortStart: 443, PortEnd: 443, Protocols: []string{"tcp"}, Exclude: false},
		{PortStart: 1, PortEnd: 65535, Exclude: true},
	})
	got := f.Apply(ports(80, 443, 8080))
	if len(got) != 1 || got[0].Number != 443 {
		t.Fatalf("expected only port 443, got %v", got)
	}
}

func TestApply_ProtocolMismatch_NotFiltered(t *testing.T) {
	f := filter.New([]filter.Rule{
		{PortStart: 80, PortEnd: 80, Protocols: []string{"udp"}, Exclude: true},
	})
	// Port 80 is tcp — the exclude rule targets udp, so it should pass through.
	got := f.Apply(ports(80))
	if len(got) != 1 {
		t.Fatalf("expected port 80 to pass through, got %v", got)
	}
}

func TestApply_RangeExclusion(t *testing.T) {
	f := filter.New([]filter.Rule{
		{PortStart: 1024, PortEnd: 49151, Exclude: true},
	})
	input := ports(80, 8080, 443, 3000)
	got := f.Apply(input)
	// Only 80 and 443 are below 1024.
	if len(got) != 2 {
		t.Fatalf("expected 2 ports, got %d: %v", len(got), got)
	}
}
