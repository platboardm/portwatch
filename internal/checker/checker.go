package checker

import (
	"fmt"
	"net"
	"time"
)

// Status represents the availability state of a port.
type Status int

const (
	StatusUp   Status = iota
	StatusDown Status = iota
)

func (s Status) String() string {
	if s == StatusUp {
		return "UP"
	}
	return "DOWN"
}

// Result holds the outcome of a single port check.
type Result struct {
	Host      string
	Port      int
	Status    Status
	Latency   time.Duration
	CheckedAt time.Time
	Err       error
}

// Checker probes TCP connectivity for a given host and port.
type Checker struct {
	Timeout time.Duration
}

// New creates a Checker with the provided dial timeout.
func New(timeout time.Duration) *Checker {
	return &Checker{Timeout: timeout}
}

// Check attempts a TCP connection to host:port and returns a Result.
func (c *Checker) Check(host string, port int) Result {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()

	conn, err := net.DialTimeout("tcp", addr, c.Timeout)
	latency := time.Since(start)

	result := Result{
		Host:      host,
		Port:      port,
		Latency:   latency,
		CheckedAt: start,
	}

	if err != nil {
		result.Status = StatusDown
		result.Err = err
		return result
	}

	_ = conn.Close()
	result.Status = StatusUp
	return result
}
