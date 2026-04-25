package eventmute_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventmute"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{
		Port: scanner.Port{Number: port, Protocol: proto},
		Type: alert.Opened,
	}
}

func TestNew_NotNil(t *testing.T) {
	m := eventmute.New()
	if m == nil {
		t.Fatal("expected non-nil Muter")
	}
}

func TestMute_EmptyName_ReturnsError(t *testing.T) {
	m := eventmute.New()
	err := m.Mute(eventmute.Rule{Port: 80, Protocol: "tcp", Duration: time.Second})
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestMute_ZeroDuration_ReturnsError(t *testing.T) {
	m := eventmute.New()
	err := m.Mute(eventmute.Rule{Name: "r", Port: 80, Protocol: "tcp", Duration: 0})
	if err == nil {
		t.Fatal("expected error for zero duration")
	}
}

func TestAllow_NoRules_ReturnsTrue(t *testing.T) {
	m := eventmute.New()
	if !m.Allow(makeEvent(443, "tcp")) {
		t.Fatal("expected Allow=true with no rules")
	}
}

func TestAllow_MatchingRule_ReturnsFalse(t *testing.T) {
	m := eventmute.New()
	_ = m.Mute(eventmute.Rule{Name: "block443", Port: 443, Protocol: "tcp", Duration: time.Hour})
	if m.Allow(makeEvent(443, "tcp")) {
		t.Fatal("expected Allow=false for muted port")
	}
}

func TestAllow_ProtocolMismatch_ReturnsTrue(t *testing.T) {
	m := eventmute.New()
	_ = m.Mute(eventmute.Rule{Name: "block443tcp", Port: 443, Protocol: "tcp", Duration: time.Hour})
	if !m.Allow(makeEvent(443, "udp")) {
		t.Fatal("expected Allow=true when protocol does not match")
	}
}

func TestAllow_EmptyProtocolRule_MatchesAny(t *testing.T) {
	m := eventmute.New()
	_ = m.Mute(eventmute.Rule{Name: "block443", Port: 443, Protocol: "", Duration: time.Hour})
	if m.Allow(makeEvent(443, "udp")) {
		t.Fatal("expected Allow=false when rule protocol is empty (matches any)")
	}
}

func TestAllow_ExpiredRule_ReturnsTrue(t *testing.T) {
	m := eventmute.New()
	_ = m.Mute(eventmute.Rule{Name: "short", Port: 80, Protocol: "tcp", Duration: time.Millisecond})
	time.Sleep(5 * time.Millisecond)
	if !m.Allow(makeEvent(80, "tcp")) {
		t.Fatal("expected Allow=true after mute expires")
	}
}

func TestUnmute_LiftsMute(t *testing.T) {
	m := eventmute.New()
	_ = m.Mute(eventmute.Rule{Name: "block22", Port: 22, Protocol: "tcp", Duration: time.Hour})
	m.Unmute("block22")
	if !m.Allow(makeEvent(22, "tcp")) {
		t.Fatal("expected Allow=true after Unmute")
	}
}

func TestPurge_RemovesExpiredEntries(t *testing.T) {
	m := eventmute.New()
	_ = m.Mute(eventmute.Rule{Name: "short", Port: 9000, Protocol: "tcp", Duration: time.Millisecond})
	_ = m.Mute(eventmute.Rule{Name: "long", Port: 9001, Protocol: "tcp", Duration: time.Hour})
	time.Sleep(5 * time.Millisecond)
	m.Purge()
	if m.Len() != 1 {
		t.Fatalf("expected 1 active rule after purge, got %d", m.Len())
	}
}
