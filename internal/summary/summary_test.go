package summary_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/summary"
)

func baseTime() time.Time {
	return time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC)
}

func TestNew_DefaultTitle(t *testing.T) {
	var buf strings.Builder
	sw := summary.New(&buf, "")
	snap := summary.Snapshot{At: baseTime(), Targets: nil}
	if err := sw.Write(snap); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "Port Status Summary") {
		t.Errorf("expected default title in output, got: %s", buf.String())
	}
}

func TestNew_PanicsOnNilWriter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil writer")
		}
	}()
	summary.New(nil, "")
}

func TestWrite_Empty(t *testing.T) {
	var buf strings.Builder
	sw := summary.New(&buf, "Test")
	snap := summary.Snapshot{At: baseTime()}
	_ = sw.Write(snap)
	if !strings.Contains(buf.String(), "(no targets)") {
		t.Errorf("expected '(no targets)' in output, got: %s", buf.String())
	}
}

func TestWrite_ContainsTargetName(t *testing.T) {
	var buf strings.Builder
	sw := summary.New(&buf, "Test")
	snap := summary.Snapshot{
		At: baseTime(),
		Targets: []summary.TargetStatus{
			{Name: "api", Addr: "localhost:8080", Up: true, Since: baseTime().Add(-5 * time.Minute)},
			{Name: "db", Addr: "localhost:5432", Up: false, Since: baseTime().Add(-2 * time.Minute)},
		},
	}
	_ = sw.Write(snap)
	out := buf.String()
	for _, want := range []string{"api", "db", "UP", "DOWN", "localhost:8080", "localhost:5432"} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestWrite_CustomTitle(t *testing.T) {
	var buf strings.Builder
	sw := summary.New(&buf, "My Custom Title")
	_ = sw.Write(summary.Snapshot{At: baseTime()})
	if !strings.Contains(buf.String(), "My Custom Title") {
		t.Errorf("expected custom title in output, got: %s", buf.String())
	}
}

func TestWrite_IncludesDuration(t *testing.T) {
	var buf strings.Builder
	sw := summary.New(&buf, "Test")
	snap := summary.Snapshot{
		At: baseTime(),
		Targets: []summary.TargetStatus{
			{Name: "svc", Addr: "host:9000", Up: true, Since: baseTime().Add(-90 * time.Second)},
		},
	}
	_ = sw.Write(snap)
	if !strings.Contains(buf.String(), "1m30s") {
		t.Errorf("expected duration '1m30s' in output, got: %s", buf.String())
	}
}
