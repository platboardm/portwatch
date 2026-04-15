package window_test

import (
	"testing"
	"time"

	"portwatch/internal/window"
)

func TestNew_PanicsOnZeroDuration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero duration")
		}
	}()
	window.New(0)
}

func TestNew_PanicsOnNegativeDuration(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative duration")
		}
	}()
	window.New(-time.Second)
}

func TestCount_EmptyOnStart(t *testing.T) {
	c := window.New(time.Second)
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestAdd_IncrementsCount(t *testing.T) {
	c := window.New(time.Second)
	c.Add()
	c.Add()
	c.Add()
	if got := c.Count(); got != 3 {
		t.Fatalf("expected 3, got %d", got)
	}
}

func TestCount_EvictsOldEvents(t *testing.T) {
	// Use a very short window so we can observe eviction without sleeping long.
	c := window.New(50 * time.Millisecond)
	c.Add()
	c.Add()
	time.Sleep(80 * time.Millisecond)
	c.Add() // this one is within the window
	if got := c.Count(); got != 1 {
		t.Fatalf("expected 1 after eviction, got %d", got)
	}
}

func TestReset_ClearsAll(t *testing.T) {
	c := window.New(time.Second)
	c.Add()
	c.Add()
	c.Reset()
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestCount_AllEvictedReturnsZero(t *testing.T) {
	c := window.New(30 * time.Millisecond)
	c.Add()
	time.Sleep(60 * time.Millisecond)
	if got := c.Count(); got != 0 {
		t.Fatalf("expected 0 after full eviction, got %d", got)
	}
}
