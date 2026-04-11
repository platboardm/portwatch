package ratelimit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/ratelimit"
)

func TestNew_PanicsOnZeroCooldown(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero cooldown")
		}
	}()
	ratelimit.New(0)
}

func TestNew_PanicsOnNegativeCooldown(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative cooldown")
		}
	}()
	ratelimit.New(-time.Second)
}

func TestAllow_FirstCallAlwaysTrue(t *testing.T) {
	l := ratelimit.New(time.Minute)
	if !l.Allow("host:8080") {
		t.Fatal("first Allow should return true")
	}
}

func TestAllow_SuppressedWithinCooldown(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("host:8080") // prime
	if l.Allow("host:8080") {
		t.Fatal("second Allow within cooldown should return false")
	}
}

func TestAllow_PassesAfterCooldown(t *testing.T) {
	now := time.Now()
	l := ratelimit.New(time.Minute)

	// Manually inject clock via unexported field would require refactor;
	// instead use a short cooldown and sleep.
	fast := ratelimit.New(10 * time.Millisecond)
	fast.Allow("k")
	time.Sleep(20 * time.Millisecond)
	_ = now
	if !fast.Allow("k") {
		t.Fatal("Allow should return true after cooldown elapsed")
	}
}

func TestAllow_IndependentKeys(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("a:80")
	if !l.Allow("b:80") {
		t.Fatal("different keys should be independent")
	}
}

func TestReset_ClearsKey(t *testing.T) {
	l := ratelimit.New(time.Minute)
	l.Allow("host:9090")
	l.Reset("host:9090")
	if !l.Allow("host:9090") {
		t.Fatal("Allow after Reset should return true")
	}
}

func TestCooldown_ReturnsConfigured(t *testing.T) {
	d := 5 * time.Minute
	l := ratelimit.New(d)
	if l.Cooldown() != d {
		t.Fatalf("expected cooldown %v, got %v", d, l.Cooldown())
	}
}
