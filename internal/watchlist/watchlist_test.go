package watchlist_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/watchlist"
)

func makeEvent(proto string, port int) alert.Event {
	return alert.Event{
		Port: scanner.Port{Protocol: proto, Port: port},
		Type: alert.Opened,
	}
}

func TestAdd_And_Match(t *testing.T) {
	w := watchlist.New()
	w.Add(watchlist.Entry{Port: 22, Protocol: "tcp", Label: "ssh"})

	ev := makeEvent("tcp", 22)
	e, ok := w.Match(ev)
	if !ok {
		t.Fatal("expected match for port 22/tcp")
	}
	if e.Label != "ssh" {
		t.Errorf("expected label ssh, got %s", e.Label)
	}
}

func TestMatch_NoEntry_ReturnsFalse(t *testing.T) {
	w := watchlist.New()
	_, ok := w.Match(makeEvent("tcp", 9999))
	if ok {
		t.Error("expected no match for unregistered port")
	}
}

func TestRemove_DeletesEntry(t *testing.T) {
	w := watchlist.New()
	w.Add(watchlist.Entry{Port: 80, Protocol: "tcp", Label: "http"})
	w.Remove("tcp", 80)
	_, ok := w.Match(makeEvent("tcp", 80))
	if ok {
		t.Error("expected entry to be removed")
	}
}

func TestMatch_ProtocolMismatch_ReturnsFalse(t *testing.T) {
	w := watchlist.New()
	w.Add(watchlist.Entry{Port: 53, Protocol: "udp", Label: "dns"})
	_, ok := w.Match(makeEvent("tcp", 53))
	if ok {
		t.Error("expected no match for wrong protocol")
	}
}

func TestAll_ReturnsSnapshot(t *testing.T) {
	w := watchlist.New()
	w.Add(watchlist.Entry{Port: 443, Protocol: "tcp", Label: "https"})
	w.Add(watchlist.Entry{Port: 8080, Protocol: "tcp", Label: "alt-http"})

	all := w.All()
	if len(all) != 2 {
		t.Errorf("expected 2 entries, got %d", len(all))
	}
}
