package timeout_test

import (
	"context"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"portwatch/internal/timeout"
)

func startTCPServer(t *testing.T) (port int, stop func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()
	return ln.Addr().(*net.TCPAddr).Port, func() { ln.Close() }
}

func TestNew_PanicsOnZeroTimeout(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero timeout")
		}
	}()
	timeout.New(0)
}

func TestNew_PanicsOnNegativeTimeout(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative timeout")
		}
	}()
	timeout.New(-time.Second)
}

func TestTimeout_ReturnsConfiguredDuration(t *testing.T) {
	c := timeout.New(5 * time.Second)
	if got := c.Timeout(); got != 5*time.Second {
		t.Errorf("Timeout() = %v, want 5s", got)
	}
}

func TestCheck_SucceedsOnOpenPort(t *testing.T) {
	port, stop := startTCPServer(t)
	defer stop()

	c := timeout.New(2 * time.Second)
	if err := c.Check(context.Background(), "127.0.0.1", port); err != nil {
		t.Errorf("expected nil error, got: %v", err)
	}
}

func TestCheck_FailsOnClosedPort(t *testing.T) {
	c := timeout.New(500 * time.Millisecond)
	// port 1 is almost certainly closed/refused
	err := c.Check(context.Background(), "127.0.0.1", 1)
	if err == nil {
		t.Fatal("expected error for closed port, got nil")
	}
}

func TestCheck_TimesOutOnStalledDial(t *testing.T) {
	// Use a non-routable address to force a timeout rather than a refusal.
	c := timeout.New(100 * time.Millisecond)
	err := c.Check(context.Background(), "192.0.2.1", 9999)
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "timeout") {
		t.Errorf("error should mention timeout, got: %v", err)
	}
}

func TestCheck_RespectsParentContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	c := timeout.New(5 * time.Second)
	err := c.Check(ctx, "127.0.0.1", 9999)
	if err == nil {
		t.Fatal("expected error on cancelled context")
	}
	_ = fmt.Sprintf("%v", err) // ensure error is printable
}
