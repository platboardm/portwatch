package throttle_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/throttle"
)

func TestNew_PanicsOnZeroLimit(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero limit")
		}
	}()
	throttle.New(0, time.Second)
}

func TestNew_PanicsOnNegativeWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative window")
		}
	}()
	throttle.New(1, -time.Second)
}

func TestAllow_FirstCallAlwaysTrue(t *testing.T) {
	th := throttle.New(3, time.Minute)
	if !th.Allow("svc") {
		t.Fatal("first call should be allowed")
	}
}

func TestAllow_RespectsLimit(t *testing.T) {
	th := throttle.New(2, time.Minute)
	if !th.Allow("svc") {
		t.Fatal("call 1 should be allowed")
	}
	if !th.Allow("svc") {
		t.Fatal("call 2 should be allowed")
	}
	if th.Allow("svc") {
		t.Fatal("call 3 should be suppressed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	th := throttle.New(1, time.Minute)
	if !th.Allow("a") {
		t.Fatal("key a call 1 should be allowed")
	}
	if th.Allow("a") {
		t.Fatal("key a call 2 should be suppressed")
	}
	if !th.Allow("b") {
		t.Fatal("key b should be independent and allowed")
	}
}

func TestAllow_PassesAfterWindowExpires(t *testing.T) {
	var fakeNow time.Time
	th := throttle.New(1, 50*time.Millisecond)
	// inject controllable clock
	fakeNow = time.Now()
	th2 := throttle.New(1, 50*time.Millisecond)
	_ = th2

	// Use real timing for simplicity.
	if !th.Allow("svc") {
		t.Fatal("first call should be allowed")
	}
	if th.Allow("svc") {
		t.Fatal("second call within window should be suppressed")
	}
	time.Sleep(60 * time.Millisecond)
	if !th.Allow("svc") {
		t.Fatalf("call after window expiry should be allowed (fakeNow=%v)", fakeNow)
	}
}

func TestReset_ClearsKey(t *testing.T) {
	th := throttle.New(1, time.Minute)
	th.Allow("svc")
	if th.Allow("svc") {
		t.Fatal("should be suppressed before reset")
	}
	th.Reset("svc")
	if !th.Allow("svc") {
		t.Fatal("should be allowed after reset")
	}
}
