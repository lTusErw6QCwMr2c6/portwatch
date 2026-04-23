package eventtransform_test

import (
	"errors"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/eventtransform"
	"github.com/user/portwatch/internal/scanner"
)

func makeEvent(port int, proto string) alert.Event {
	return alert.Event{
		Port:     scanner.Port{Number: port, Protocol: proto},
		Protocol: proto,
	}
}

func TestNew_NotNil(t *testing.T) {
	tr := eventtransform.New()
	if tr == nil {
		t.Fatal("expected non-nil Transformer")
	}
}

func TestRegister_AddsStage(t *testing.T) {
	tr := eventtransform.New()
	if err := tr.Register("upper", eventtransform.UppercaseProtocol()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 1 {
		t.Fatalf("expected 1 stage, got %d", tr.Len())
	}
}

func TestRegister_EmptyName_ReturnsError(t *testing.T) {
	tr := eventtransform.New()
	if err := tr.Register("", eventtransform.UppercaseProtocol()); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestRegister_DuplicateName_ReturnsError(t *testing.T) {
	tr := eventtransform.New()
	_ = tr.Register("dup", eventtransform.UppercaseProtocol())
	if err := tr.Register("dup", eventtransform.UppercaseProtocol()); err == nil {
		t.Fatal("expected error for duplicate stage name")
	}
}

func TestApply_UppercasesProtocol(t *testing.T) {
	tr := eventtransform.New()
	_ = tr.Register("upper", eventtransform.UppercaseProtocol())
	e := makeEvent(80, "tcp")
	out, err := tr.Apply(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Protocol != "TCP" {
		t.Fatalf("expected TCP, got %s", out.Protocol)
	}
}

func TestApply_StampNow_SetsTimestamp(t *testing.T) {
	tr := eventtransform.New()
	_ = tr.Register("stamp", eventtransform.StampNow())
	e := makeEvent(443, "tcp")
	before := time.Now().UTC()
	out, err := tr.Apply(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Timestamp.Before(before) {
		t.Fatal("expected timestamp to be set to now")
	}
}

func TestApply_SetLabel_AddsLabel(t *testing.T) {
	tr := eventtransform.New()
	_ = tr.Register("label", eventtransform.SetLabel("env", "prod"))
	e := makeEvent(8080, "udp")
	out, err := tr.Apply(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Labels["env"] != "prod" {
		t.Fatalf("expected label env=prod, got %v", out.Labels)
	}
}

func TestApply_StageError_Propagates(t *testing.T) {
	tr := eventtransform.New()
	_ = tr.Register("fail", func(e *alert.Event) (*alert.Event, error) {
		return nil, errors.New("boom")
	})
	e := makeEvent(22, "tcp")
	_, err := tr.Apply(e)
	if err == nil {
		t.Fatal("expected error from failing stage")
	}
}

func TestApply_NoStages_ReturnsOriginal(t *testing.T) {
	tr := eventtransform.New()
	e := makeEvent(53, "udp")
	out, err := tr.Apply(e)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Port.Number != 53 {
		t.Fatalf("expected port 53, got %d", out.Port.Number)
	}
}
