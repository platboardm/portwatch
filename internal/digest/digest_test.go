package digest_test

import (
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/digest"
)

func makeEntry(target, status string, sev alert.Severity) digest.Entry {
	return digest.Entry{
		Target:    target,
		Status:    status,
		Severity:  sev,
		ChangedAt: time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
}

func TestNew_DefaultTitle(t *testing.T) {
	d := digest.New("")
	var sb strings.Builder
	_ = d.Write(&sb)
	if !strings.Contains(sb.String(), "portwatch digest") {
		t.Errorf("expected default title in output, got: %s", sb.String())
	}
}

func TestNew_CustomTitle(t *testing.T) {
	d := digest.New("my custom title")
	var sb strings.Builder
	_ = d.Write(&sb)
	if !strings.Contains(sb.String(), "my custom title") {
		t.Errorf("expected custom title in output, got: %s", sb.String())
	}
}

func TestWrite_Empty(t *testing.T) {
	d := digest.New("test")
	var sb strings.Builder
	if err := d.Write(&sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sb.String(), "no events") {
		t.Errorf("expected 'no events' message, got: %s", sb.String())
	}
}

func TestWrite_ContainsEntries(t *testing.T) {
	d := digest.New("test")
	d.Add(makeEntry("db:5432", "down", alert.SeverityCritical))
	d.Add(makeEntry("api:8080", "up", alert.SeverityInfo))

	var sb strings.Builder
	_ = d.Write(&sb)
	out := sb.String()

	if !strings.Contains(out, "db:5432") {
		t.Errorf("expected target db:5432 in output")
	}
	if !strings.Contains(out, "api:8080") {
		t.Errorf("expected target api:8080 in output")
	}
	if !strings.Contains(out, "total events: 2") {
		t.Errorf("expected total count in output, got: %s", out)
	}
}

func TestLen_And_Reset(t *testing.T) {
	d := digest.New("test")
	d.Add(makeEntry("svc:80", "down", alert.SeverityCritical))
	d.Add(makeEntry("svc:443", "down", alert.SeverityCritical))

	if d.Len() != 2 {
		t.Fatalf("expected len 2, got %d", d.Len())
	}
	d.Reset()
	if d.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", d.Len())
	}
}

func TestWrite_SeverityLabel(t *testing.T) {
	d := digest.New("test")
	d.Add(makeEntry("svc:9000", "down", alert.SeverityCritical))
	var sb strings.Builder
	_ = d.Write(&sb)
	if !strings.Contains(sb.String(), "critical") {
		t.Errorf("expected severity label 'critical' in output, got: %s", sb.String())
	}
}
