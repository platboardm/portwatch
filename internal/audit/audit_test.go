package audit_test

import (
	"bytes"
	"strings"
	"testing"

	"portwatch/internal/audit"
)

func TestNew_PanicsOnNilWriter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	audit.New(nil)
}

func TestLog_WritesLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Log(audit.KindStartup, "portwatch", "started")
	line := buf.String()
	if !strings.Contains(line, "startup") {
		t.Errorf("expected 'startup' in output, got: %s", line)
	}
	if !strings.Contains(line, "portwatch") {
		t.Errorf("expected target in output, got: %s", line)
	}
	if !strings.Contains(line, "started") {
		t.Errorf("expected message in output, got: %s", line)
	}
}

func TestLog_MultipleEvents(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Log(audit.KindStateChange, "db:5432", "up -> down")
	l.Log(audit.KindAlertSent, "db:5432", "notified")
	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
}

func TestEventString_ContainsAllFields(t *testing.T) {
	var buf bytes.Buffer
	l := audit.New(&buf)
	l.Log(audit.KindAlertDropped, "svc:80", "rate limited")
	s := buf.String()
	for _, want := range []string{"alert_dropped", "svc:80", "rate limited"} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in: %s", want, s)
		}
	}
}
