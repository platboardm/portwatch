package jitter_test

import (
	"testing"
	"time"

	"portwatch/internal/jitter"
)

// fixedSource always returns the same value, making tests deterministic.
type fixedSource struct{ val int64 }

func (f *fixedSource) Int63n(_ int64) int64 { return f.val }

func TestNew_PanicsOnZeroFactor(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for factor=0")
		}
	}()
	jitter.New(0)
}

func TestNew_PanicsOnFactorAboveOne(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for factor>1")
		}
	}()
	jitter.New(1.1)
}

func TestNewWithSource_PanicsOnNilSource(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil source")
		}
	}()
	jitter.NewWithSource(0.1, nil)
}

func TestApply_ZeroOffset(t *testing.T) {
	// fixedSource returns 0, so Apply should return exactly base.
	j := jitter.NewWithSource(0.5, &fixedSource{val: 0})
	base := 10 * time.Second
	got := j.Apply(base)
	if got != base {
		t.Fatalf("expected %v, got %v", base, got)
	}
}

func TestApply_AddsOffset(t *testing.T) {
	const offset = 1_000_000_000 // 1 second in nanoseconds
	j := jitter.NewWithSource(0.5, &fixedSource{val: offset})
	base := 10 * time.Second
	got := j.Apply(base)
	want := base + time.Duration(offset)
	if got != want {
		t.Fatalf("expected %v, got %v", want, got)
	}
}

func TestApply_NonPositiveBase(t *testing.T) {
	j := jitter.NewWithSource(0.25, &fixedSource{val: 999})
	for _, base := range []time.Duration{0, -1 * time.Second} {
		if got := j.Apply(base); got != base {
			t.Fatalf("Apply(%v): expected %v unchanged, got %v", base, base, got)
		}
	}
}

func TestApply_WithinBounds(t *testing.T) {
	// With a real random source the result must be in [base, base*(1+factor)).
	j := jitter.New(0.25)
	base := 100 * time.Millisecond
	for i := 0; i < 1000; i++ {
		got := j.Apply(base)
		if got < base || got >= time.Duration(float64(base)*1.25)+1 {
			t.Fatalf("Apply result %v out of expected range", got)
		}
	}
}
