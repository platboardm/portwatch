// Package metrics provides lightweight in-process counters and gauges
// for tracking portwatch runtime statistics such as checks performed,
// alerts fired, and current target states.
package metrics

import (
	"sync"
	"sync/atomic"
)

// Counter is a monotonically increasing uint64 counter.
type Counter struct {
	v uint64
}

// Inc increments the counter by 1.
func (c *Counter) Inc() { atomic.AddUint64(&c.v, 1) }

// Add increments the counter by n.
func (c *Counter) Add(n uint64) { atomic.AddUint64(&c.v, n) }

// Value returns the current counter value.
func (c *Counter) Value() uint64 { return atomic.LoadUint64(&c.v) }

// Gauge holds an int64 value that can go up or down.
type Gauge struct {
	mu sync.RWMutex
	v  int64
}

// Set replaces the gauge value.
func (g *Gauge) Set(v int64) {
	g.mu.Lock()
	g.v = v
	g.mu.Unlock()
}

// Inc increments the gauge by 1.
func (g *Gauge) Inc() {
	g.mu.Lock()
	g.v++
	g.mu.Unlock()
}

// Dec decrements the gauge by 1.
func (g *Gauge) Dec() {
	g.mu.Lock()
	g.v--
	g.mu.Unlock()
}

// Value returns the current gauge value.
func (g *Gauge) Value() int64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.v
}

// Registry holds a named set of counters and gauges.
type Registry struct {
	mu       sync.RWMutex
	counters map[string]*Counter
	gauges   map[string]*Gauge
}

// New returns an initialised Registry.
func New() *Registry {
	return &Registry{
		counters: make(map[string]*Counter),
		gauges:   make(map[string]*Gauge),
	}
}

// Counter returns the named Counter, creating it if necessary.
func (r *Registry) Counter(name string) *Counter {
	r.mu.Lock()
	defer r.mu.Unlock()
	if c, ok := r.counters[name]; ok {
		return c
	}
	c := &Counter{}
	r.counters[name] = c
	return c
}

// Gauge returns the named Gauge, creating it if necessary.
func (r *Registry) Gauge(name string) *Gauge {
	r.mu.Lock()
	defer r.mu.Unlock()
	if g, ok := r.gauges[name]; ok {
		return g
	}
	g := &Gauge{}
	r.gauges[name] = g
	return g
}

// Snapshot returns a point-in-time copy of all counter and gauge values.
func (r *Registry) Snapshot() map[string]int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make(map[string]int64, len(r.counters)+len(r.gauges))
	for k, c := range r.counters {
		out[k] = int64(c.Value())
	}
	for k, g := range r.gauges {
		out[k] = g.Value()
	}
	return out
}
