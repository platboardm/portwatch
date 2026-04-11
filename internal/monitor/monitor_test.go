package monitor_test

import (
	"bytes"
	"context"
	"net"
	"testing"
	"time"

	"github.com/user/portwatch/internal/checker"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notifier"
)

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("freePort: %v", err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	l.Close()
	return port
}

func TestMonitor_DetectsDown(t *testing.T) {
	port := freePort(t)
	cfg := &config.Config{
		Interval: 1,
		Targets: []config.Target{
			{Name: "test", Host: "127.0.0.1", Port: port},
		},
	}

	var buf bytes.Buffer
	c := checker.New(500 * time.Millisecond)
	n := notifier.New(&buf)
	m := monitor.New(cfg, c, n)

	ctx, cancel := context.WithTimeout(context.Background(), 2500*time.Millisecond)
	defer cancel()

	m.Run(ctx)

	if buf.Len() == 0 {
		t.Error("expected at least one notification, got none")
	}
	if got := buf.String(); !containsString(got, "DOWN") {
		t.Errorf("expected DOWN in output, got: %s", got)
	}
}

func TestMonitor_NoSpuriousAlerts(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port

	cfg := &config.Config{
		Interval: 1,
		Targets: []config.Target{
			{Name: "stable", Host: "127.0.0.1", Port: port},
		},
	}

	var buf bytes.Buffer
	c := checker.New(500 * time.Millisecond)
	n := notifier.New(&buf)
	m := monitor.New(cfg, c, n)

	ctx, cancel := context.WithTimeout(context.Background(), 3500*time.Millisecond)
	defer cancel()
	m.Run(ctx)

	// Service stayed UP the whole time — only one notification expected.
	count := bytes.Count(buf.Bytes(), []byte("\n"))
	if count > 1 {
		t.Errorf("expected 1 notification, got %d: %s", count, buf.String())
	}
}

func containsString(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsRune(s, sub))
}

func containsRune(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
