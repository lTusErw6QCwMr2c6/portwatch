package fence_test

import (
	"testing"

	"github.com/user/portwatch/internal/fence"
)

func TestNew_ValidRange(t *testing.T) {
	f, err := fence.New([]fence.Range{{Start: 1, End: 1024, Protocol: "tcp"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Len() != 1 {
		t.Fatalf("expected 1 range, got %d", f.Len())
	}
}

func TestNew_InvalidRange_ReturnsError(t *testing.T) {
	_, err := fence.New([]fence.Range{{Start: 9000, End: 80}})
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestAllow_PortInRange(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Start: 80, End: 443, Protocol: "tcp"}})
	if !f.Allow(80, "tcp") {
		t.Error("expected port 80/tcp to be allowed")
	}
	if !f.Allow(443, "tcp") {
		t.Error("expected port 443/tcp to be allowed")
	}
}

func TestAllow_PortOutsideRange(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Start: 80, End: 443, Protocol: "tcp"}})
	if f.Allow(8080, "tcp") {
		t.Error("expected port 8080 to be denied")
	}
}

func TestAllow_ProtocolMismatch_Denied(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Start: 53, End: 53, Protocol: "udp"}})
	if f.Allow(53, "tcp") {
		t.Error("expected tcp on udp-only range to be denied")
	}
}

func TestAllow_EmptyProtocol_MatchesBoth(t *testing.T) {
	f, _ := fence.New([]fence.Range{{Start: 8080, End: 8080}})
	if !f.Allow(8080, "tcp") {
		t.Error("expected tcp to be allowed on protocol-agnostic range")
	}
	if !f.Allow(8080, "udp") {
		t.Error("expected udp to be allowed on protocol-agnostic range")
	}
}

func TestRanges_ReturnsCopy(t *testing.T) {
	orig := []fence.Range{{Start: 1000, End: 2000}}
	f, _ := fence.New(orig)
	r := f.Ranges()
	r[0].Start = 9999
	if f.Ranges()[0].Start == 9999 {
		t.Error("Ranges should return a copy, not a reference")
	}
}
