package geoip

import (
	"strings"
	"testing"
)

func TestLookup_LoopbackIsPrivate(t *testing.T) {
	info := Lookup("127.0.0.1")
	if !info.Private {
		t.Errorf("expected private, got %+v", info)
	}
	if info.Country != "private" {
		t.Errorf("expected country=private, got %s", info.Country)
	}
}

func TestLookup_RFC1918IsPrivate(t *testing.T) {
	for _, ip := range []string{"10.0.0.1", "192.168.1.1", "172.16.0.1"} {
		info := Lookup(ip)
		if !info.Private {
			t.Errorf("expected %s to be private", ip)
		}
	}
}

func TestLookup_PublicIP_NotPrivate(t *testing.T) {
	info := Lookup("8.8.8.8")
	if info.Private {
		t.Errorf("expected public IP to not be private")
	}
	if info.Country != "unknown" {
		t.Errorf("expected unknown country for stub, got %s", info.Country)
	}
}

func TestLookup_InvalidIP_ReturnsUnknown(t *testing.T) {
	info := Lookup("not-an-ip")
	if info.Country != "unknown" {
		t.Errorf("expected unknown, got %s", info.Country)
	}
	if info.Private {
		t.Error("invalid IP should not be marked private")
	}
}

func TestInfo_String_Private(t *testing.T) {
	info := Info{IP: "192.168.0.1", Country: "private", Private: true}
	if !strings.Contains(info.String(), "private") {
		t.Errorf("expected 'private' in string, got %s", info.String())
	}
}

func TestInfo_String_Public(t *testing.T) {
	info := Info{IP: "8.8.8.8", Country: "US"}
	s := info.String()
	if !strings.Contains(s, "8.8.8.8") || !strings.Contains(s, "US") {
		t.Errorf("unexpected string: %s", s)
	}
}
