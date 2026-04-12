// Package pipeline wires together the checker, debounce, ratelimit, and
// notifier stages into a single reusable processing pipeline for a target.
package pipeline

import (
	"context"
	"io"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/checker"
	"github.com/user/portwatch/internal/debounce"
	"github.com/user/portwatch/internal/notifier"
	"github.com/user/portwatch/internal/ratelimit"
	"github.com/user/portwatch/internal/state"
)

// Config holds the tunable parameters for a Pipeline.
type Config struct {
	// DebounceThreshold is the number of consecutive identical results required
	// before a state change is accepted.
	DebounceThreshold int

	// AlertCooldown is the minimum duration between repeated alerts for the
	// same target.
	AlertCooldown time.Duration
}

// Pipeline processes a single check result through debounce → state → ratelimit
// → notify stages.
type Pipeline struct {
	cfg      Config
	debounce *debounce.Debouncer
	store    *state.Store
	rl       *ratelimit.Limiter
	n        *notifier.Notifier
}

// New constructs a Pipeline for the given target name, writing alerts to w.
// Panics if cfg.DebounceThreshold < 1 or cfg.AlertCooldown <= 0.
func New(target string, cfg Config, w io.Writer) *Pipeline {
	if cfg.DebounceThreshold < 1 {
		cfg.DebounceThreshold = 1
	}
	return &Pipeline{
		cfg:      cfg,
		debounce: debounce.New(cfg.DebounceThreshold),
		store:    state.New(),
		rl:       ratelimit.New(cfg.AlertCooldown),
		n:        notifier.New(w),
	}
}

// Process takes a raw checker.Status for the named target, runs it through
// the pipeline stages and emits an alert when appropriate.
func (p *Pipeline) Process(ctx context.Context, target string, status checker.Status) {
	if !p.debounce.Confirm(bool(status == checker.StatusUp)) {
		return
	}

	changed, prev := p.store.Set(target, state.Status(status))
	if !changed {
		return
	}

	sev := alert.SeverityWarning
	if status == checker.StatusDown {
		sev = alert.SeverityCritical
	}

	a := alert.New(target, string(prev), string(state.Status(status)), sev)

	if !p.rl.Allow(target) {
		return
	}

	p.n.Notify(ctx, a)
}
