package metrics_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"portwatch/internal/metrics"
)

func TestNewExporter_PanicsOnNilRegistry(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil registry")
		}
	}()
	metrics.NewExporter(nil)
}

func TestWriteText_ContainsMetricNames(t *testing.T) {
	reg := metrics.New()
	reg.Counter("checks_total").Add(10)
	reg.Gauge("targets_up").Set(3)

	ex := metrics.NewExporter(reg)
	var buf bytes.Buffer
	if err := ex.WriteText(&buf); err != nil {
		t.Fatalf("WriteText error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"checks_total", "10", "targets_up", "3"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q:\n%s", want, out)
		}
	}
}

func TestWriteJSON_ValidJSON(t *testing.T) {
	reg := metrics.New()
	reg.Counter("alerts_fired").Add(2)

	ex := metrics.NewExporter(reg)
	var buf bytes.Buffer
	if err := ex.WriteJSON(&buf); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}
	var got map[string]int64
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got["alerts_fired"] != 2 {
		t.Errorf("expected alerts_fired=2, got %d", got["alerts_fired"])
	}
}

func TestWriteText_SortedOutput(t *testing.T) {
	reg := metrics.New()
	reg.Counter("zebra").Inc()
	reg.Counter("apple").Inc()
	reg.Counter("mango").Inc()

	ex := metrics.NewExporter(reg)
	var buf bytes.Buffer
	_ = ex.WriteText(&buf)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// lines[0] is header; lines[1..] are data rows
	if len(lines) < 4 {
		t.Fatalf("expected at least 4 lines, got %d", len(lines))
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[1]), "apple") {
		t.Errorf("expected first data row to be 'apple', got: %s", lines[1])
	}
}
