package eventtag

import (
	"bytes"
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/logger"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{Port: port, Protocol: proto, Type: alert.Opened}
}

func TestNew_NotNil(t *testing.T) {
	if New() == nil {
		t.Fatal("expected non-nil Tagger")
	}
}

func TestRegister_AddsRule(t *testing.T) {
	tr := New()
	if err := tr.Register("http", 80, "tcp", []string{"web"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 1 {
		t.Fatalf("expected 1 rule, got %d", tr.Len())
	}
}

func TestRegister_DuplicateName_ReturnsError(t *testing.T) {
	tr := New()
	_ = tr.Register("http", 80, "tcp", []string{"web"})
	if err := tr.Register("http", 443, "tcp", []string{"tls"}); err == nil {
		t.Fatal("expected error for duplicate rule name")
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	tr := New()
	if err := tr.Register("", 80, "tcp", []string{"web"}); err == nil {
		t.Fatal("expected error for empty rule name")
	}
}

func TestApply_MatchingRule_ReturnsTags(t *testing.T) {
	tr := New()
	_ = tr.Register("dns", 53, "udp", []string{"infra", "dns"})
	tags := tr.Apply(makeEvent(53, "udp"))
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestApply_ProtocolMismatch_ReturnsNil(t *testing.T) {
	tr := New()
	_ = tr.Register("dns", 53, "udp", []string{"infra"})
	if tags := tr.Apply(makeEvent(53, "tcp")); tags != nil {
		t.Fatalf("expected nil, got %v", tags)
	}
}

func TestApply_EmptyProtocol_MatchesBoth(t *testing.T) {
	tr := New()
	_ = tr.Register("any53", 53, "", []string{"dns"})
	if tags := tr.Apply(makeEvent(53, "tcp")); len(tags) == 0 {
		t.Fatal("expected match for empty protocol rule")
	}
	if tags := tr.Apply(makeEvent(53, "udp")); len(tags) == 0 {
		t.Fatal("expected match for empty protocol rule")
	}
}

func TestRemove_DeletesRule(t *testing.T) {
	tr := New()
	_ = tr.Register("http", 80, "tcp", []string{"web"})
	if !tr.Remove("http") {
		t.Fatal("expected Remove to return true")
	}
	if tr.Len() != 0 {
		t.Fatalf("expected 0 rules after remove, got %d", tr.Len())
	}
}

func TestHandler_CallsNext(t *testing.T) {
	tr := New()
	_ = tr.Register("ssh", 22, "tcp", []string{"admin"})
	var buf bytes.Buffer
	log := logger.New(&buf)
	var got []string
	h := NewHandler(tr, log, func(_ alert.Event, tags []string) { got = tags })
	h.Handle(makeEvent(22, "tcp"))
	if len(got) == 0 {
		t.Fatal("expected tags to be forwarded to next")
	}
}
