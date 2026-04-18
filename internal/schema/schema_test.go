package schema_test

import (
	"testing"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/schema"
)

func makeEvent(proto string, opened bool) alert.Event {
	return alert.Event{
		Port:   scanner.Port{Port: 8080, Protocol: proto},
		Opened: opened,
	}
}

func TestValidate_ValidOpenedTCP(t *testing.T) {
	v := schema.New(schema.DefaultRules())
	if err := v.Validate(makeEvent("tcp", true)); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_ValidClosedUDP(t *testing.T) {
	v := schema.New(schema.DefaultRules())
	if err := v.Validate(makeEvent("udp", false)); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidate_InvalidProtocol(t *testing.T) {
	v := schema.New(schema.DefaultRules())
	err := v.Validate(makeEvent("sctp", true))
	if err == nil {
		t.Fatal("expected error for invalid protocol")
	}
}

func TestValidate_MissingProtocol(t *testing.T) {
	v := schema.New(schema.DefaultRules())
	err := v.Validate(makeEvent("", true))
	if err == nil {
		t.Fatal("expected error for missing protocol")
	}
}

func TestValidateAll_ReturnsErrorMap(t *testing.T) {
	v := schema.New(schema.DefaultRules())
	events := []alert.Event{
		makeEvent("tcp", true),
		makeEvent("bad", false),
		makeEvent("udp", true),
	}
	errs := v.ValidateAll(events)
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if _, ok := errs[1]; !ok {
		t.Error("expected error at index 1")
	}
}

func TestValidateAll_AllValid_EmptyMap(t *testing.T) {
	v := schema.New(schema.DefaultRules())
	events := []alert.Event{
		makeEvent("tcp", true),
		makeEvent("udp", false),
	}
	errs := v.ValidateAll(events)
	if len(errs) != 0 {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestNew_CustomRule(t *testing.T) {
	rules := []schema.Rule{
		{Field: "protocol", Required: true, AllowedValues: []string{"tcp"}},
	}
	v := schema.New(rules)
	if err := v.Validate(makeEvent("udp", true)); err == nil {
		t.Fatal("expected udp to fail custom rule allowing only tcp")
	}
}
