package snapshot_test

import (
	"strings"
	"testing"
	"time"

	"portwatch/internal/snapshot"
)

func TestNewWriter_PanicsOnNilWriter(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for nil writer")
		}
	}()
	snapshot.NewWriter(nil, "")
}

func TestWrite_DefaultTitle(t *testing.T) {
	var buf strings.Builder
	w := snapshot.NewWriter(&buf, "")
	w.Write(snapshot.Snapshot{CapturedAt: time.Now()})
	if !strings.Contains(buf.String(), "Port Status Snapshot") {
		t.Errorf("expected default title in output, got: %s", buf.String())
	}
}

func TestWrite_Empty(t *testing.T) {
	var buf strings.Builder
	w := snapshot.NewWriter(&buf, "Test")
	w.Write(snapshot.Snapshot{CapturedAt: time.Now()})
	if !strings.Contains(buf.String(), "no targets") {
		t.Errorf("expected 'no targets' message, got: %s", buf.String())
	}
}

func TestWrite_ContainsTargetName(t *testing.T) {
	var buf strings.Builder
	w := snapshot.NewWriter(&buf, "Test")
	snap := snapshot.Snapshot{
		CapturedAt: time.Now(),
		Statuses: []snapshot.Status{
			{Target: "api-gateway", Up: true, Since: time.Now()},
		},
	}
	w.Write(snap)
	if !strings.Contains(buf.String(), "api-gateway") {
		t.Errorf("expected target name in output, got: %s", buf.String())
	}
}

func TestWrite_StatusLabels(t *testing.T) {
	var buf strings.Builder
	w := snapshot.NewWriter(&buf, "Test")
	snap := snapshot.Snapshot{
		CapturedAt: time.Now(),
		Statuses: []snapshot.Status{
			{Target: "svc-up", Up: true, Since: time.Now()},
			{Target: "svc-down", Up: false, Since: time.Now()},
		},
	}
	w.Write(snap)
	out := buf.String()
	if !strings.Contains(out, "UP") {
		t.Error("expected UP label in output")
	}
	if !strings.Contains(out, "DOWN") {
		t.Error("expected DOWN label in output")
	}
}
