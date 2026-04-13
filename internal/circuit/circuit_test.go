package circuit_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/circuit"
)

func TestStateString(t *testing.T) {
	cases := []struct {
		s    circuit.State
		want string
	}{
		{circuit.StateClosed, "closed"},
		{circuit.StateOpen, "open"},
		{circuit.StateHalfOpen, "half-open"},
		{circuit.State(99), "unknown(99)"},
	}
	for _, tc := range cases {
		if got := tc.s.String(); got != tc.want {
			t.Errorf("State(%d).String() = %q; want %q", int(tc.s), got, tc.want)
		}
	}
}

func TestNew_PanicsOnZeroThreshold(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	circuit.New(0, time.Second)
}

func TestNew_PanicsOnZeroWindow(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic")
		}
	}()
	circuit.New(1, 0)
}

func TestAllow_ClosedByDefault(t *testing.T) {
	b := circuit.New(3, time.Second)
	if !b.Allow() {
		t.Error("expected Allow() == true for fresh breaker")
	}
}

func TestBreaker_OpensAtThreshold(t *testing.T) {
	b := circuit.New(3, time.Second)
	for i := 0; i < 3; i++ {
		b.RecordFailure()
	}
	if b.State() != circuit.StateOpen {
		t.Errorf("expected StateOpen after threshold failures, got %s", b.State())
	}
	if b.Allow() {
		t.Error("expected Allow() == false when open")
	}
}

func TestBreaker_SuccessResetsClosed(t *testing.T) {
	b := circuit.New(2, time.Second)
	b.RecordFailure()
	b.RecordFailure()
	b.RecordSuccess()
	if b.State() != circuit.StateClosed {
		t.Errorf("expected StateClosed after success, got %s", b.State())
	}
}

func TestBreaker_HalfOpenAfterWindow(t *testing.T) {
	b := circuit.New(1, 20*time.Millisecond)
	b.RecordFailure()
	time.Sleep(30 * time.Millisecond)
	if !b.Allow() {
		t.Error("expected Allow() == true in half-open after reset window")
	}
	if b.State() != circuit.StateHalfOpen {
		t.Errorf("expected StateHalfOpen, got %s", b.State())
	}
}
