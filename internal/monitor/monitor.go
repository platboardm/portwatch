// Package monitor orchestrates periodic port checking and dispatches
// notifications when a target's status changes.
package monitor

import (
	"context"
	"sync"
	"time"

	"github.com/user/portwatch/internal/checker"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/notifier"
)

// Monitor runs a check loop for every configured target.
type Monitor struct {
	cfg      *config.Config
	checker  *checker.Checker
	notifier *notifier.Notifier
	states   map[string]checker.Status
	mu       sync.Mutex
}

// New creates a Monitor from the provided configuration.
func New(cfg *config.Config, c *checker.Checker, n *notifier.Notifier) *Monitor {
	return &Monitor{
		cfg:      cfg,
		checker:  c,
		notifier: n,
		states:   make(map[string]checker.Status),
	}
}

// Run starts one goroutine per target and blocks until ctx is cancelled.
func (m *Monitor) Run(ctx context.Context) {
	var wg sync.WaitGroup
	for _, t := range m.cfg.Targets {
		wg.Add(1)
		go func(t config.Target) {
			defer wg.Done()
			m.loop(ctx, t)
		}(t)
	}
	wg.Wait()
}

func (m *Monitor) loop(ctx context.Context, t config.Target) {
	interval := time.Duration(m.cfg.Interval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.check(t)
		}
	}
}
func (m *Monitor) check(t config.Target) {
	event := m.checker.Check(t)\tm.mu.Lock()
	prev, seen := m.states[eventAddr]
	m.states[event.Addr] = event.Status
	m.mu.Unlock()

	if !seen || prev != event.Status {
		m.notifier.Notify(event)
	}
}
