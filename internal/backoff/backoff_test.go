package backoff_test

import (
	"testing"
	"time"

	"portwatch/internal/backoff"
)

func TestNew_DefaultMax(t *testing.T) {
	s := backoff.New(time.Second, 0)
	if s.Max != 30*time.Second {
		t.Fatalf("expected max 30s, got %v", s.Max)
	}
}

func TestNew_PanicsOnZeroBase(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero base")
		}
	}()
	backoff.New(0, time.Minute)
}

func TestDelay_AttemptZero(t *testing.T) {
	s := backoff.New(500*time.Millisecond, 0)
	if got := s.Delay(0); got != 500*time.Millisecond {
		t.Fatalf("attempt 0: want 500ms, got %v", got)
	}
}

func TestDelay_GrowsExponentially(t *testing.T) {
	s := backoff.New(time.Second, time.Minute)
	prev := s.Delay(0)
	for i := 1; i <= 4; i++ {
		curr := s.Delay(i)
		if curr <= prev {
			t.Fatalf("attempt %d: delay %v should be greater than %v", i, curr, prev)
		}
		prev = curr
	}
}

func TestDelay_CapsAtMax(t *testing.T) {
	max := 5 * time.Second
	s := backoff.New(time.Second, max)
	for i := 0; i <= 20; i++ {
		if d := s.Delay(i); d > max {
			t.Fatalf("attempt %d: delay %v exceeds max %v", i, d, max)
		}
	}
}

func TestDelay_NegativeAttemptTreatedAsZero(t *testing.T) {
	s := backoff.New(time.Second, 0)
	if got, want := s.Delay(-5), time.Second; got != want {
		t.Fatalf("negative attempt: want %v, got %v", want, got)
	}
}

func TestDelay_CustomMultiplier(t *testing.T) {
	s := backoff.New(time.Second, time.Minute)
	s.Multiplier = 3.0
	// attempt 1 should be base * 3^1 = 3s
	if got, want := s.Delay(1), 3*time.Second; got != want {
		t.Fatalf("want %v, got %v", want, got)
	}
}
