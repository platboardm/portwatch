package checker_test

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/portwatch/portwatch/internal/checker"
)

// startListener opens a TCP listener on a random port and returns the port number.
func startListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port, _ := strconv.Atoi(ln.Addr().(*net.TCPAddr).Port.String())
	// Extract port properly
	port = ln.Addr().(*net.TCPAddr).Port
	return port, func() { ln.Close() }
}

func TestCheck_Up(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not start listener: %v", err)
	}
	defer ln.Close()

	port := ln.Addr().(*net.TCPAddr).Port
	c := checker.New(2 * time.Second)
	res := c.Check("127.0.0.1", port)

	if res.Status != checker.StatusUp {
		t.Errorf("expected StatusUp, got %s: %v", res.Status, res.Err)
	}
	if res.Latency <= 0 {
		t.Errorf("expected positive latency, got %v", res.Latency)
	}
}

func TestCheck_Down(t *testing.T) {
	// Port 1 is almost certainly closed/refused on loopback.
	c := checker.New(500 * time.Millisecond)
	res := c.Check("127.0.0.1", 1)

	if res.Status != checker.StatusDown {
		t.Errorf("expected StatusDown, got %s", res.Status)
	}
	if res.Err == nil {
		t.Error("expected a non-nil error for a closed port")
	}
}

func TestStatusString(t *testing.T) {
	if checker.StatusUp.String() != "UP" {
		t.Errorf("unexpected string for StatusUp")
	}
	if checker.StatusDown.String() != "DOWN" {
		t.Errorf("unexpected string for StatusDown")
	}
}
