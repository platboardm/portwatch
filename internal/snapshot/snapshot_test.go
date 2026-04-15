package snapshot_test

import (
	"testing"
	"time"

	"portwatch/internal/snapshot"
)

func makeStatus(target string, up bool) snapshot.Status {
	return snapshot.Status{
		Target: target,
		Host:   "localhost",
		Port:   8080,
		Up:     up,
		Since:  time.Now(),
		Checks: 1,
	}
}

func TestStore_RecordAndLen(t *testing.T) {
	s := snapshot.New()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
	s.Record("svc-a", makeStatus("svc-a", true))
	if s.Len() != 1 {
		t.Fatalf("expected 1, got %d", s.Len())
	}
}

func TestStore_OverwritesSameKey(t *testing.T) {
	s := snapshot.New()
	s.Record("svc-a", makeStatus("svc-a", true))
	s.Record("svc-a", makeStatus("svc-a", false))
	if s.Len() != 1 {
		t.Fatalf("expected 1 entry after overwrite, got %d", s.Len())
	}
	snap := s.Capture()
	if snap.Statuses[0].Up {
		t.Error("expected overwritten status to be down")
	}
}

func TestCapture_ReturnsAllEntries(t *testing.T) {
	s := snapshot.New()
	s.Record("a", makeStatus("a", true))
	s.Record("b", makeStatus("b", false))
	s.Record("c", makeStatus("c", true))

	snap := s.Capture()
	if len(snap.Statuses) != 3 {
		t.Fatalf("expected 3 statuses, got %d", len(snap.Statuses))
	}
	if snap.CapturedAt.IsZero() {
		t.Error("CapturedAt should not be zero")
	}
}

func TestCapture_IsImmutable(t *testing.T) {
	s := snapshot.New()
	s.Record("a", makeStatus("a", true))

	snap1 := s.Capture()
	s.Record("b", makeStatus("b", false))
	snap2 := s.Capture()

	if len(snap1.Statuses) != 1 {
		t.Errorf("snap1 should have 1 entry, got %d", len(snap1.Statuses))
	}
	if len(snap2.Statuses) != 2 {
		t.Errorf("snap2 should have 2 entries, got %d", len(snap2.Statuses))
	}
}
