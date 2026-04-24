package mute_test

import (
	"testing"
	"time"

	"portwatch/internal/mute"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func TestIsMuted_FalseByDefault(t *testing.T) {
	s := mute.New(fixedNow(base))
	if s.IsMuted("svc-a") {
		t.Fatal("expected not muted before any window is registered")
	}
}

func TestIsMuted_TrueWithinWindow(t *testing.T) {
	s := mute.New(fixedNow(base))
	s.Mute("svc-a", base.Add(10*time.Minute))
	if !s.IsMuted("svc-a") {
		t.Fatal("expected muted within active window")
	}
}

func TestIsMuted_FalseAfterWindowExpires(t *testing.T) {
	s := mute.New(fixedNow(base.Add(20 * time.Minute)))
	s.Mute("svc-a", base.Add(10*time.Minute))
	if s.IsMuted("svc-a") {
		t.Fatal("expected not muted after window has expired")
	}
}

func TestUnmute_RemovesWindow(t *testing.T) {
	s := mute.New(fixedNow(base))
	s.Mute("svc-b", base.Add(1*time.Hour))
	s.Unmute("svc-b")
	if s.IsMuted("svc-b") {
		t.Fatal("expected not muted after explicit unmute")
	}
}

func TestMute_ReplacesExistingWindow(t *testing.T) {
	s := mute.New(fixedNow(base))
	s.Mute("svc-c", base.Add(5*time.Minute))
	// Replace with a window already in the past.
	s.Mute("svc-c", base.Add(-1*time.Minute))
	if s.IsMuted("svc-c") {
		t.Fatal("expected not muted after window replaced with expired one")
	}
}

func TestActive_ReturnsOnlyActiveWindows(t *testing.T) {
	s := mute.New(fixedNow(base))
	s.Mute("active", base.Add(1*time.Hour))
	s.Mute("expired", base.Add(-1*time.Minute))

	active := s.Active()
	if len(active) != 1 {
		t.Fatalf("expected 1 active window, got %d", len(active))
	}
	if active[0].Target != "active" {
		t.Errorf("expected target 'active', got %q", active[0].Target)
	}
}

func TestActive_EmptyWhenNoneMuted(t *testing.T) {
	s := mute.New(fixedNow(base))
	if got := s.Active(); len(got) != 0 {
		t.Fatalf("expected empty active list, got %d entries", len(got))
	}
}
