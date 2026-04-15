package cooldown

import (
	"testing"
	"time"
)

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on zero window")
		}
	}()
	New(0)
}

func TestNew_PanicsOnNegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic on negative window")
		}
	}()
	New(-time.Second)
}

func TestAllow_FirstCallAlwaysTrue(t *testing.T) {
	tr := New(time.Minute)
	if !tr.Allow("svc-a") {
		t.Fatal("first call should always be allowed")
	}
}

func TestAllow_SuppressedWithinWindow(t *testing.T) {
	now := time.Now()
	tr := New(time.Minute)
	tr.nowFunc = func() time.Time { return now }

	tr.Allow("svc-a") // record trigger
	if tr.Allow("svc-a") {
		t.Fatal("second call within window should be suppressed")
	}
}

func TestAllow_PassesAfterWindow(t *testing.T) {
	now := time.Now()
	tr := New(time.Minute)
	tr.nowFunc = func() time.Time { return now }
	tr.Allow("svc-a")

	tr.nowFunc = func() time.Time { return now.Add(time.Minute) }
	if !tr.Allow("svc-a") {
		t.Fatal("call after window should be allowed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	now := time.Now()
	tr := New(time.Minute)
	tr.nowFunc = func() time.Time { return now }

	tr.Allow("svc-a")
	if !tr.Allow("svc-b") {
		t.Fatal("different key should not be affected")
	}
}

func TestReset_AllowsImmediately(t *testing.T) {
	now := time.Now()
	tr := New(time.Minute)
	tr.nowFunc = func() time.Time { return now }

	tr.Allow("svc-a")
	tr.Reset("svc-a")
	if !tr.Allow("svc-a") {
		t.Fatal("after reset, Allow should return true")
	}
}

func TestRemaining_ZeroWhenNotTracked(t *testing.T) {
	tr := New(time.Minute)
	if r := tr.Remaining("unknown"); r != 0 {
		t.Fatalf("expected 0, got %v", r)
	}
}

func TestRemaining_PositiveWithinWindow(t *testing.T) {
	now := time.Now()
	tr := New(time.Minute)
	tr.nowFunc = func() time.Time { return now }
	tr.Allow("svc-a")

	tr.nowFunc = func() time.Time { return now.Add(20 * time.Second) }
	r := tr.Remaining("svc-a")
	if r <= 0 || r > time.Minute {
		t.Fatalf("unexpected remaining: %v", r)
	}
}

func TestRemaining_ZeroAfterWindow(t *testing.T) {
	now := time.Now()
	tr := New(time.Minute)
	tr.nowFunc = func() time.Time { return now }
	tr.Allow("svc-a")

	tr.nowFunc = func() time.Time { return now.Add(2 * time.Minute) }
	if r := tr.Remaining("svc-a"); r != 0 {
		t.Fatalf("expected 0 after window, got %v", r)
	}
}
