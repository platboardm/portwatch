package envelope_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/envelope"
)

func makeAlert() alert.Alert {
	return alert.New("db", "localhost", 5432, alert.SeverityDown)
}

func TestNew_SetsFields(t *testing.T) {
	a := makeAlert()
	before := time.Now()
	e := envelope.New(a, envelope.PriorityHigh, "ops")
	after := time.Now()

	if e.Alert != a {
		t.Errorf("expected alert to be preserved")
	}
	if e.Priority != envelope.PriorityHigh {
		t.Errorf("expected priority high, got %v", e.Priority)
	}
	if e.Channel != "ops" {
		t.Errorf("expected channel ops, got %q", e.Channel)
	}
	if e.CreatedAt.Before(before) || e.CreatedAt.After(after) {
		t.Errorf("unexpected timestamp %v", e.CreatedAt)
	}
}

func TestNew_TraceIDNotEmpty(t *testing.T) {
	e := envelope.New(makeAlert(), envelope.PriorityNormal, "")
	if e.TraceID == "" {
		t.Error("expected non-empty trace ID")
	}
}

func TestNew_UniqueTraceIDs(t *testing.T) {
	e1 := envelope.New(makeAlert(), envelope.PriorityNormal, "")
	e2 := envelope.New(makeAlert(), envelope.PriorityNormal, "")
	if e1.TraceID == e2.TraceID {
		t.Error("expected unique trace IDs")
	}
}

func TestPriorityString(t *testing.T) {
	cases := []struct {
		p    envelope.Priority
		want string
	}{
		{envelope.PriorityLow, "low"},
		{envelope.PriorityNormal, "normal"},
		{envelope.PriorityHigh, "high"},
		{envelope.PriorityCritical, "critical"},
		{envelope.Priority(99), "priority(99)"},
	}
	for _, tc := range cases {
		if got := tc.p.String(); got != tc.want {
			t.Errorf("Priority(%d).String() = %q, want %q", int(tc.p), got, tc.want)
		}
	}
}

func TestNew_EmptyChannel(t *testing.T) {
	e := envelope.New(makeAlert(), envelope.PriorityLow, "")
	if e.Channel != "" {
		t.Errorf("expected empty channel, got %q", e.Channel)
	}
}

func TestTraceID_IsHex(t *testing.T) {
	e := envelope.New(makeAlert(), envelope.PriorityNormal, "")
	const hexChars = "0123456789abcdef"
	for _, c := range e.TraceID {
		if !strings.ContainsRune(hexChars, c) {
			t.Errorf("trace ID %q contains non-hex char %q", e.TraceID, c)
		}
	}
}
