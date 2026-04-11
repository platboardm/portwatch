package history_test

import (
	"testing"

	"portwatch/internal/checker"
	"portwatch/internal/history"
)

func TestRecorder_Record(t *testing.T) {
	rec := history.NewRecorder(50)
	rec.Record("api", "localhost", 9000, checker.StatusUp)
	rec.Record("api", "localhost", 9000, checker.StatusDown)

	if rec.Len() != 2 {
		t.Fatalf("expected 2 entries, got %d", rec.Len())
	}

	entries := rec.Entries()
	if entries[0].Status != "up" {
		t.Errorf("first entry: want status \"up\", got %q", entries[0].Status)
	}
	if entries[1].Status != "down" {
		t.Errorf("second entry: want status \"down\", got %q", entries[1].Status)
	}
}

func TestRecorder_Fields(t *testing.T) {
	rec := history.NewRecorder(10)
	rec.Record("db", "10.0.0.1", 5432, checker.StatusUp)

	e := rec.Entries()[0]
	if e.Target != "db" {
		t.Errorf("want target \"db\", got %q", e.Target)
	}
	if e.Host != "10.0.0.1" {
		t.Errorf("want host \"10.0.0.1\", got %q", e.Host)
	}
	if e.Port != 5432 {
		t.Errorf("want port 5432, got %d", e.Port)
	}
	if e.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
}

func TestRecorder_Capacity(t *testing.T) {
	rec := history.NewRecorder(3)
	for i := 0; i < 7; i++ {
		rec.Record("svc", "localhost", 80, checker.StatusUp)
	}
	if rec.Len() != 3 {
		t.Fatalf("expected recorder capped at 3, got %d", rec.Len())
	}
}
