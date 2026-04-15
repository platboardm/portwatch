package retry_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/retry"
)

var errTemp = errors.New("temporary error")

func TestNew_PanicsOnZeroAttempts(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	retry.New(0)
}

func TestDo_SucceedsOnFirstAttempt(t *testing.T) {
	p := retry.New(3)
	p.BaseDelay = time.Millisecond
	calls := 0
	err := p.Do(context.Background(), func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

func TestDo_RetriesUntilSuccess(t *testing.T) {
	p := retry.New(5)
	p.BaseDelay = time.Millisecond
	var count int32
	err := p.Do(context.Background(), func() error {
		if atomic.AddInt32(&count, 1) < 3 {
			return errTemp
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 calls, got %d", count)
	}
}

func TestDo_ExhaustsAttempts(t *testing.T) {
	p := retry.New(3)
	p.BaseDelay = time.Millisecond
	var count int32
	err := p.Do(context.Background(), func() error {
		atomic.AddInt32(&count, 1)
		return errTemp
	})
	if !errors.Is(err, retry.ErrExhausted) {
		t.Fatalf("expected ErrExhausted, got %v", err)
	}
	if count != 3 {
		t.Fatalf("expected 3 calls, got %d", count)
	}
}

func TestDo_RespectsContextCancel(t *testing.T) {
	p := retry.New(10)
	p.BaseDelay = 50 * time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	var count int32
	go func() {
		time.Sleep(10 * time.Millisecond)
		cancel()
	}()
	err := p.Do(ctx, func() error {
		atomic.AddInt32(&count, 1)
		return errTemp
	})
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
	if count > 3 {
		t.Fatalf("too many calls before cancel: %d", count)
	}
}

func TestDo_DelayGrows(t *testing.T) {
	p := retry.New(4)
	p.BaseDelay = 10 * time.Millisecond
	p.MaxDelay = 100 * time.Millisecond
	p.Multiplier = 3.0
	start := time.Now()
	_ = p.Do(context.Background(), func() error { return errTemp })
	elapsed := time.Since(start)
	// 10 + 30 + 90 = 130 ms minimum between 4 attempts
	if elapsed < 100*time.Millisecond {
		t.Fatalf("delays too short: %v", elapsed)
	}
}
