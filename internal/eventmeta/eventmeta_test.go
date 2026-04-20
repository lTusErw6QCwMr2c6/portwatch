package eventmeta_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventmeta"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(proto string, port int) alert.Event {
	return alert.Event{
		Type: alert.Opened,
		Port: scanner.Port{Protocol: proto, Number: port},
	}
}

func TestNew_NotNil(t *testing.T) {
	s := eventmeta.New()
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestSet_And_Get(t *testing.T) {
	s := eventmeta.New()
	e := makeEvent("tcp", 8080)
	s.Set(e, eventmeta.Meta{"env": "prod", "team": "platform"})

	m, ok := s.Get(e)
	if !ok {
		t.Fatal("expected metadata to be present")
	}
	if m["env"] != "prod" {
		t.Errorf("expected env=prod, got %s", m["env"])
	}
	if m["team"] != "platform" {
		t.Errorf("expected team=platform, got %s", m["team"])
	}
}

func TestGet_Missing_ReturnsFalse(t *testing.T) {
	s := eventmeta.New()
	e := makeEvent("udp", 53)
	_, ok := s.Get(e)
	if ok {
		t.Fatal("expected false for missing event")
	}
}

func TestSet_MergesKeys(t *testing.T) {
	s := eventmeta.New()
	e := makeEvent("tcp", 443)
	s.Set(e, eventmeta.Meta{"a": "1"})
	s.Set(e, eventmeta.Meta{"b": "2"})

	m, _ := s.Get(e)
	if m["a"] != "1" || m["b"] != "2" {
		t.Errorf("expected merged meta, got %v", m)
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	s := eventmeta.New()
	e := makeEvent("tcp", 22)
	s.Set(e, eventmeta.Meta{"owner": "ops"})
	s.Delete(e)

	_, ok := s.Get(e)
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestReset_ClearsAll(t *testing.T) {
	s := eventmeta.New()
	s.Set(makeEvent("tcp", 80), eventmeta.Meta{"x": "1"})
	s.Set(makeEvent("udp", 161), eventmeta.Meta{"y": "2"})
	s.Reset()

	if s.Len() != 0 {
		t.Errorf("expected Len 0 after Reset, got %d", s.Len())
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	s := eventmeta.New()
	e := makeEvent("tcp", 9000)
	s.Set(e, eventmeta.Meta{"key": "original"})

	m, _ := s.Get(e)
	m["key"] = "mutated"

	m2, _ := s.Get(e)
	if m2["key"] != "original" {
		t.Errorf("Get should return a copy; store was mutated")
	}
}
