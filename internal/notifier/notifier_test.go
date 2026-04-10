package notifier_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"portwatch/internal/notifier"
)

func makeEvent(up bool) notifier.Event {
	return notifier.Event{
		Target:    "web",
		Host:      "localhost",
		Port:      8080,
		Up:        up,
		Timestamp: time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestNotify_Down(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf)

	if err := n.Notify(makeEvent(false)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "DOWN") {
		t.Errorf("expected DOWN in output, got: %q", got)
	}
	if !strings.Contains(got, "web") {
		t.Errorf("expected target name in output, got: %q", got)
	}
	if !strings.Contains(got, "localhost:8080") {
		t.Errorf("expected host:port in output, got: %q", got)
	}
}

func TestNotify_Up(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf)

	if err := n.Notify(makeEvent(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := buf.String()
	if !strings.Contains(got, "UP") {
		t.Errorf("expected UP in output, got: %q", got)
	}
}

func TestNew_NilWriter(t *testing.T) {
	// Should not panic when nil is passed; defaults to stdout.
	n := notifier.New(nil)
	if n == nil {
		t.Fatal("expected non-nil notifier")
	}
}

func TestNotify_TimestampFormat(t *testing.T) {
	var buf bytes.Buffer
	n := notifier.New(&buf)

	if err := n.Notify(makeEvent(true)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// RFC3339 timestamp should appear in output.
	if !strings.Contains(buf.String(), "2024-01-15T10:00:00Z") {
		t.Errorf("expected RFC3339 timestamp in output, got: %q", buf.String())
	}
}
