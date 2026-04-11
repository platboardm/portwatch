package alert_test

import (
	"strings"
	"testing"
	"time"

	"portwatch/internal/alert"
)

func TestSeverityString(t *testing.T) {
	tests := []struct {
		sev  alert.Severity
		want string
	}{
		{alert.SeverityInfo, "INFO"},
		{alert.SeverityCritical, "CRITICAL"},
		{alert.Severity(99), "UNKNOWN"},
	}
	for _, tc := range tests {
		if got := tc.sev.String(); got != tc.want {
			t.Errorf("Severity(%d).String() = %q, want %q", tc.sev, got, tc.want)
		}
	}
}

func TestNew_SetsFields(t *testing.T) {
	before := time.Now().UTC()
	a := alert.New("localhost:8080", alert.SeverityCritical, "port unreachable")
	after := time.Now().UTC()

	if a.Target != "localhost:8080" {
		t.Errorf("Target = %q, want %q", a.Target, "localhost:8080")
	}
	if a.Severity != alert.SeverityCritical {
		t.Errorf("Severity = %v, want SeverityCritical", a.Severity)
	}
	if a.Message != "port unreachable" {
		t.Errorf("Message = %q, want %q", a.Message, "port unreachable")
	}
	if a.OccurredAt.Before(before) || a.OccurredAt.After(after) {
		t.Errorf("OccurredAt %v outside expected range [%v, %v]", a.OccurredAt, before, after)
	}
}

func TestAlert_String_ContainsFields(t *testing.T) {
	a := alert.New("10.0.0.1:443", alert.SeverityInfo, "service recovered")
	s := a.String()

	for _, want := range []string{"INFO", "10.0.0.1:443", "service recovered"} {
		if !strings.Contains(s, want) {
			t.Errorf("Alert.String() = %q, missing %q", s, want)
		}
	}
}

func TestAlert_String_CriticalLabel(t *testing.T) {
	a := alert.New("db:5432", alert.SeverityCritical, "down")
	if !strings.Contains(a.String(), "CRITICAL") {
		t.Errorf("expected CRITICAL in string output, got: %s", a.String())
	}
}
