package tag

import (
	"testing"
)

func TestLookup_WellKnownPort(t *testing.T) {
	tagger := New(nil)
	tag := tagger.Lookup(22, "tcp")
	if tag.Label != "ssh" {
		t.Errorf("expected ssh, got %s", tag.Label)
	}
}

func TestLookup_UnknownPort(t *testing.T) {
	tagger := New(nil)
	tag := tagger.Lookup(9999, "tcp")
	if tag.Label != "unknown" {
		t.Errorf("expected unknown, got %s", tag.Label)
	}
}

func TestLookup_CustomOverridesWellKnown(t *testing.T) {
	tagger := New(map[string]string{"tcp:80": "my-app"})
	tag := tagger.Lookup(80, "tcp")
	if tag.Label != "my-app" {
		t.Errorf("expected my-app, got %s", tag.Label)
	}
}

func TestLookup_CustomNewEntry(t *testing.T) {
	tagger := New(map[string]string{"tcp:8443": "custom-tls"})
	tag := tagger.Lookup(8443, "tcp")
	if tag.Label != "custom-tls" {
		t.Errorf("expected custom-tls, got %s", tag.Label)
	}
}

func TestLookupAll_ReturnsTags(t *testing.T) {
	tagger := New(nil)
	tags := tagger.LookupAll([]int{22, 80, 9999}, "tcp")
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}
	if tags[0].Label != "ssh" {
		t.Errorf("expected ssh, got %s", tags[0].Label)
	}
	if tags[1].Label != "http" {
		t.Errorf("expected http, got %s", tags[1].Label)
	}
	if tags[2].Label != "unknown" {
		t.Errorf("expected unknown, got %s", tags[2].Label)
	}
}

func TestLookup_UDP(t *testing.T) {
	tagger := New(nil)
	tag := tagger.Lookup(53, "udp")
	if tag.Label != "dns" {
		t.Errorf("expected dns, got %s", tag.Label)
	}
}

func TestTag_Fields(t *testing.T) {
	tagger := New(nil)
	tag := tagger.Lookup(443, "tcp")
	if tag.Port != 443 || tag.Protocol != "tcp" || tag.Label != "https" {
		t.Errorf("unexpected tag fields: %+v", tag)
	}
}
