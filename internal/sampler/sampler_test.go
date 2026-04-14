package sampler_test

import (
	"testing"

	"portwatch/internal/sampler"
)

// fixedSource always returns the same value.
type fixedSource struct{ val float64 }

func (f *fixedSource) Float64() float64 { return f.val }

func TestNew_PanicsOnNegativeRate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative rate")
		}
	}()
	sampler.New(-0.1)
}

func TestNew_PanicsOnRateAboveOne(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for rate > 1.0")
		}
	}()
	sampler.New(1.1)
}

func TestNewWithSource_PanicsOnNilSource(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil source")
		}
	}()
	sampler.NewWithSource(0.5, nil)
}

func TestAllow_RateOne_AlwaysTrue(t *testing.T) {
	s := sampler.NewWithSource(1.0, &fixedSource{val: 0.99})
	for i := 0; i < 10; i++ {
		if !s.Allow() {
			t.Fatal("expected Allow to return true for rate=1.0")
		}
	}
}

func TestAllow_RateZero_AlwaysFalse(t *testing.T) {
	s := sampler.NewWithSource(0.0, &fixedSource{val: 0.0})
	for i := 0; i < 10; i++ {
		if s.Allow() {
			t.Fatal("expected Allow to return false for rate=0.0")
		}
	}
}

func TestAllow_BelowRate_Passes(t *testing.T) {
	// source returns 0.3, rate is 0.5 → 0.3 < 0.5 → allowed
	s := sampler.NewWithSource(0.5, &fixedSource{val: 0.3})
	if !s.Allow() {
		t.Fatal("expected Allow to return true when source < rate")
	}
}

func TestAllow_AboveRate_Suppressed(t *testing.T) {
	// source returns 0.7, rate is 0.5 → 0.7 >= 0.5 → suppressed
	s := sampler.NewWithSource(0.5, &fixedSource{val: 0.7})
	if s.Allow() {
		t.Fatal("expected Allow to return false when source >= rate")
	}
}

func TestRate_ReturnsConfigured(t *testing.T) {
	s := sampler.New(0.42)
	if got := s.Rate(); got != 0.42 {
		t.Fatalf("expected rate 0.42, got %v", got)
	}
}
