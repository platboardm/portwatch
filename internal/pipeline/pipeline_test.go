package pipeline_test

import (
	"bytes"
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user/portwatch/internal/checker"
	"github.com/user/portwatch/internal/pipeline"
)

func newPipeline(t *testing.T, threshold int, cooldown time.Duration) (*pipeline.Pipeline, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	p := pipeline.New("svc", pipeline.Config{
		DebounceThreshold: threshold,
		AlertCooldown:     cooldown,
	}, &buf)
	return p, &buf
}

func TestProcess_NoAlertBeforeDebounceThreshold(t *testing.T) {
	p, buf := newPipeline(t, 3, time.Hour)
	ctx := context.Background()

	p.Process(ctx, "svc", checker.StatusDown)
	p.Process(ctx, "svc", checker.StatusDown)

	if buf.Len() != 0 {
		t.Fatalf("expected no output before threshold, got: %s", buf.String())
	}
}

func TestProcess_AlertAfterDebounceThreshold(t *testing.T) {
	p, buf := newPipeline(t, 2, time.Millisecond)
	ctx := context.Background()

	p.Process(ctx, "svc", checker.StatusDown)
	p.Process(ctx, "svc", checker.StatusDown)

	if !strings.Contains(buf.String(), "svc") {
		t.Fatalf("expected alert containing target name, got: %s", buf.String())
	}
}

func TestProcess_NoSpuriousAlertOnNoChange(t *testing.T) {
	p, buf := newPipeline(t, 1, time.Millisecond)
	ctx := context.Background()

	p.Process(ctx, "svc", checker.StatusDown)
	first := buf.String()

	p.Process(ctx, "svc", checker.StatusDown)
	if buf.String() != first {
		t.Fatal("expected no second alert for unchanged state")
	}
}

func TestProcess_RateLimitSuppressesRepeated(t *testing.T) {
	p, buf := newPipeline(t, 1, time.Hour)
	ctx := context.Background()

	p.Process(ctx, "svc", checker.StatusDown)
	first := buf.String()

	// Force a state flip then back so debounce and state both pass again.
	p.Process(ctx, "svc", checker.StatusUp)
	p.Process(ctx, "svc", checker.StatusDown)

	if buf.String() != first {
		t.Fatal("expected rate limiter to suppress second alert within cooldown")
	}
}

func TestNew_DefaultDebounceThreshold(t *testing.T) {
	// threshold 0 should be clamped to 1 — no panic
	var buf bytes.Buffer
	p := pipeline.New("svc", pipeline.Config{DebounceThreshold: 0, AlertCooldown: time.Millisecond}, &buf)
	p.Process(context.Background(), "svc", checker.StatusDown)
	if buf.Len() == 0 {
		t.Fatal("expected alert with threshold clamped to 1")
	}
}
