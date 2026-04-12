package reporter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/user/portwatch/internal/aggregator"
	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/reporter"
)

func makeAgg(t *testing.T) *aggregator.Aggregator {
	t.Helper()
	return aggregator.New(64)
}

func TestNew_PanicsOnNilAggregator(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil aggregator")
		}
	}()
	reporter.New(nil, &bytes.Buffer{}, "test")
}

func TestNew_PanicsOnNilWriter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil writer")
		}
	}()
	reporter.New(makeAgg(t), nil, "test")
}

func TestWrite_ContainsTitle(t *testing.T) {
	agg := makeAgg(t)
	var buf bytes.Buffer
	r := reporter.New(agg, &buf, "my service")
	if err := r.Write(); err != nil {
		t.Fatalf("Write() error: %v", err)
	}
	if !strings.Contains(buf.String(), "my service") {
		t.Errorf("output missing title: %q", buf.String())
	}
}

func TestWrite_DefaultTitle(t *testing.T) {
	agg := makeAgg(t)
	var buf bytes.Buffer
	r := reporter.New(agg, &buf, "")
	_ = r.Write()
	if !strings.Contains(buf.String(), "portwatch status") {
		t.Errorf("expected default title in output: %q", buf.String())
	}
}

func TestWrite_ShowsUpDown(t *testing.T) {
	agg := makeAgg(t)
	agg.Record(alert.New("api:8080", alert.SeverityInfo))
	agg.Record(alert.New("db:5432", alert.SeverityCritical))

	var buf bytes.Buffer
	r := reporter.New(agg, &buf, "test")
	if err := r.Write(); err != nil {
		t.Fatalf("Write() error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "api:8080") {
		t.Errorf("expected api:8080 in output: %q", out)
	}
	if !strings.Contains(out, "db:5432") {
		t.Errorf("expected db:5432 in output: %q", out)
	}
	if !strings.Contains(out, "total=") {
		t.Errorf("expected summary line in output: %q", out)
	}
}
