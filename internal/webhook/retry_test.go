package webhook_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/portwatch/internal/webhook"
)

func TestRetry_SucceedsOnSecondAttempt(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if calls.Add(1) < 2 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := webhook.New(srv.URL, 0)
	rn := webhook.NewRetry(n, 3, 10*time.Millisecond)

	if err := rn.Notify(context.Background(), makeAlert(t)); err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}
	if got := calls.Load(); got != 2 {
		t.Errorf("expected 2 calls, got %d", got)
	}
}

func TestRetry_ExhaustsAttempts(t *testing.T) {
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls.Add(1)
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	n := webhook.New(srv.URL, 0)
	rn := webhook.NewRetry(n, 3, 10*time.Millisecond)

	if err := rn.Notify(context.Background(), makeAlert(t)); err == nil {
		t.Error("expected error after exhausted retries, got nil")
	}
	if got := calls.Load(); got != 3 {
		t.Errorf("expected 3 calls, got %d", got)
	}
}

func TestRetry_RespectsContextCancel(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	ctx, cancel := context.WithCancel(context.Background())
	n := webhook.New(srv.URL, 0)
	rn := webhook.NewRetry(n, 5, 200*time.Millisecond)

	// Cancel after a short delay so the first attempt fires but back-off is interrupted.
	go func() { time.Sleep(50 * time.Millisecond); cancel() }()

	if err := rn.Notify(ctx, makeAlert(t)); err == nil {
		t.Error("expected error on context cancel, got nil")
	}
}

func TestNewRetry_PanicsOnZeroAttempts(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic for maxAttempts=0")
		}
	}()
	webhook.NewRetry(webhook.New("http://example.com", 0), 0, 0)
}
