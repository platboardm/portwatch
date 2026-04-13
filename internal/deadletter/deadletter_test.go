package deadletter_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/deadletter"
)

func makeAlert(target string) alert.Alert {
	return alert.New(target, alert.SeverityCritical, "port unreachable")
}

func TestNew_PanicsOnZeroCap(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero capacity")
		}
	}()
	deadletter.New(0)
}

func TestPush_AndLen(t *testing.T) {
	q := deadletter.New(10)
	q.Push(makeAlert("svc-a"), "timeout", 3)
	q.Push(makeAlert("svc-b"), "connection refused", 1)
	if got := q.Len(); got != 2 {
		t.Fatalf("Len() = %d, want 2", got)
	}
}

func TestDrain_ReturnsAndClears(t *testing.T) {
	q := deadletter.New(10)
	q.Push(makeAlert("svc-a"), "timeout", 3)
	entries := q.Drain()
	if len(entries) != 1 {
		t.Fatalf("Drain() returned %d entries, want 1", len(entries))
	}
	if q.Len() != 0 {
		t.Fatal("queue should be empty after Drain()")
	}
}

func TestDrain_Fields(t *testing.T) {
	q := deadletter.New(10)
	before := time.Now()
	q.Push(makeAlert("svc-x"), "dial error", 5)
	after := time.Now()

	entries := q.Drain()
	e := entries[0]

	if e.Alert.Target != "svc-x" {
		t.Errorf("Target = %q, want %q", e.Alert.Target, "svc-x")
	}
	if e.Reason != "dial error" {
		t.Errorf("Reason = %q, want %q", e.Reason, "dial error")
	}
	if e.Attempts != 5 {
		t.Errorf("Attempts = %d, want 5", e.Attempts)
	}
	if e.FailedAt.Before(before) || e.FailedAt.After(after) {
		t.Errorf("FailedAt %v outside expected range", e.FailedAt)
	}
}

func TestPush_EvictsOldestWhenFull(t *testing.T) {
	q := deadletter.New(3)
	q.Push(makeAlert("first"), "err", 1)
	q.Push(makeAlert("second"), "err", 1)
	q.Push(makeAlert("third"), "err", 1)
	q.Push(makeAlert("fourth"), "err", 1) // should evict "first"

	if q.Len() != 3 {
		t.Fatalf("Len() = %d, want 3", q.Len())
	}
	entries := q.Drain()
	if entries[0].Alert.Target != "second" {
		t.Errorf("oldest entry = %q, want %q", entries[0].Alert.Target, "second")
	}
	if entries[2].Alert.Target != "fourth" {
		t.Errorf("newest entry = %q, want %q", entries[2].Alert.Target, "fourth")
	}
}
