// Package probe provides a configurable TCP probe that records latency
// alongside the up/down result for a target.
package probe

import (
	"net"
	"time"
)

// Result holds the outcome of a single probe attempt.
type Result struct {
	Target  string
	Addr    string
	Up      bool
	Latency time.Duration
	At      time.Time
}

// Prober performs TCP probes against a target address.
type Prober struct {
	timeout time.Duration
	dial    func(network, addr string, timeout time.Duration) (net.Conn, error)
}

// New returns a Prober with the given dial timeout.
// Panics if timeout is zero or negative.
func New(timeout time.Duration) *Prober {
	if timeout <= 0 {
		panic("probe: timeout must be positive")
	}
	return &Prober{timeout: timeout, dial: net.DialTimeout}
}

// Probe dials addr and returns a Result. A successful connection means Up=true.
func (p *Prober) Probe(target, addr string) Result {
	start := time.Now()
	conn, err := p.dial("tcp", addr, p.timeout)
	latency := time.Since(start)
	if err != nil {
		return Result{Target: target, Addr: addr, Up: false, Latency: latency, At: start}
	}
	_ = conn.Close()
	return Result{Target: target, Addr: addr, Up: true, Latency: latency, At: start}
}
