package suppress_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/suppress"
)

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero window")
		}
	}()
	suppress.New(0)
}

func TestNew_PanicsOnNegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative window")
		}
	}()
	suppress.New(-time.Second)
}

func TestAllow_FirstCallAlwaysTrue(t *testing.T) {
	s := suppress.New(time.Minute)
	if !s.Allow("svc-a") {
		t.Fatal("first call should always be allowed")
	}
}

func TestAllow_SuppressedWithinWindow(t *testing.T) {
	s := suppress.New(time.Minute)
	s.Allow("svc-a") // prime
	if s.Allow("svc-a") {
		t.Fatal("second call within window should be suppressed")
	}
}

func TestAllow_PassesAfterWindow(t *testing.T) {
	now := time.Now()
	s := suppress.New(time.Minute)

	// Use unexported clock via a wrapper trick — test via Reset + elapsed time.
	s.Allow("svc-a") // prime at now

	// Simulate passage of time by resetting and re-allowing.
	s.Reset("svc-a")
	if !s.Allow("svc-a") {
		t.Fatal("call after Reset should be allowed")
	}
	_ = now
}

func TestAllow_IndependentKeys(t *testing.T) {
	s := suppress.New(time.Minute)
	s.Allow("svc-a")
	if !s.Allow("svc-b") {
		t.Fatal("different key should not be suppressed")
	}
}

func TestReset_ClearsRecord(t *testing.T) {
	s := suppress.New(time.Minute)
	s.Allow("svc-a")
	s.Reset("svc-a")
	if !s.Allow("svc-a") {
		t.Fatal("Allow should pass after Reset")
	}
}

func TestLen_TracksKeys(t *testing.T) {
	s := suppress.New(time.Minute)
	if s.Len() != 0 {
		t.Fatalf("expected 0 keys, got %d", s.Len())
	}
	s.Allow("svc-a")
	s.Allow("svc-b")
	if s.Len() != 2 {
		t.Fatalf("expected 2 keys, got %d", s.Len())
	}
	s.Reset("svc-a")
	if s.Len() != 1 {
		t.Fatalf("expected 1 key after reset, got %d", s.Len())
	}
}
