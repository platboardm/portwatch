package heartbeat_test

import (
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/heartbeat"
)

func TestNew_PanicsOnZeroInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for zero interval")
		}
	}()
	heartbeat.New(0)
}

func TestNew_PanicsOnNegativeInterval(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for negative interval")
		}
	}()
	heartbeat.New(-time.Second)
}

func TestEmitter_ReceivesBeat(t *testing.T) {
	e := heartbeat.New(20 * time.Millisecond)
	ch := e.Subscribe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go e.Run(ctx)

	select {
	case b := <-ch:
		if b.At.IsZero() {
			t.Fatal("beat timestamp is zero")
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("timed out waiting for beat")
	}
}

func TestEmitter_MultipleSubscribers(t *testing.T) {
	e := heartbeat.New(20 * time.Millisecond)
	ch1 := e.Subscribe()
	ch2 := e.Subscribe()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go e.Run(ctx)

	for _, ch := range []<-chan heartbeat.Beat{ch1, ch2} {
		select {
		case <-ch:
		case <-time.After(200 * time.Millisecond):
			t.Fatal("timed out waiting for beat on subscriber")
		}
	}
}

func TestEmitter_ClosesOnCancel(t *testing.T) {
	e := heartbeat.New(50 * time.Millisecond)
	ch := e.Subscribe()

	ctx, cancel := context.WithCancel(context.Background())
	go e.Run(ctx)
	cancel()

	select {
	case _, ok := <-ch:
		if ok {
			// drain until closed
		}
	case <-time.After(300 * time.Millisecond):
		t.Fatal("channel not closed after context cancel")
	}
}
