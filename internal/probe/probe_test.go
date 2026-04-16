package probe

import (
	"net"
	"testing"
	"time"
)

func startTCP(t *testing.T) string {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	t.Cleanup(func() { ln.Close() })
	go func() {
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			c.Close()
		}
	}()
	return ln.Addr().String()
}

func TestNew_PanicsOnZeroTimeout(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic")
		}
	}()
	New(0)
}

func TestProbe_Up(t *testing.T) {
	addr := startTCP(t)
	p := New(time.Second)
	r := p.Probe("svc", addr)
	if !r.Up {
		t.Fatal("expected Up=true")
	}
	if r.Latency <= 0 {
		t.Fatal("expected positive latency")
	}
	if r.Target != "svc" {
		t.Fatalf("unexpected target %q", r.Target)
	}
}

func TestProbe_Down(t *testing.T) {
	p := New(200 * time.Millisecond)
	r := p.Probe("svc", "127.0.0.1:1")
	if r.Up {
		t.Fatal("expected Up=false")
	}
}

func TestProbe_SetsAt(t *testing.T) {
	before := time.Now()
	p := New(200 * time.Millisecond)
	r := p.Probe("svc", "127.0.0.1:1")
	if r.At.Before(before) {
		t.Fatal("At should be >= before probe")
	}
}
