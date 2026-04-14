// Package timeout provides per-target check timeout enforcement,
// wrapping a deadline around each port-check attempt and returning
// a structured error when the deadline is exceeded.
package timeout

import (
	"context"
	"fmt"
	"net"
	"time"
)

// Checker performs a TCP dial with a configurable deadline.
type Checker struct {
	timeout time.Duration
	dialer  func(ctx context.Context, network, addr string) error
}

// New returns a Checker that enforces the given timeout on every check.
// It panics if timeout is zero or negative.
func New(timeout time.Duration) *Checker {
	if timeout <= 0 {
		panic("timeout: timeout must be positive")
	}
	return &Checker{
		timeout: timeout,
		dialer:  defaultDial,
	}
}

// Check attempts to open a TCP connection to host:port within the
// configured timeout. It returns nil on success and a descriptive
// error on failure or timeout.
func (c *Checker) Check(ctx context.Context, host string, port int) error {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	addr := fmt.Sprintf("%s:%d", host, port)
	if err := c.dialer(ctx, "tcp", addr); err != nil {
		if ctx.Err() != nil {
			return fmt.Errorf("timeout: check timed out after %s for %s", c.timeout, addr)
		}
		return fmt.Errorf("timeout: dial %s: %w", addr, err)
	}
	return nil
}

// Timeout returns the configured deadline duration.
func (c *Checker) Timeout() time.Duration {
	return c.timeout
}

func defaultDial(ctx context.Context, network, addr string) error {
	var d net.Dialer
	conn, err := d.DialContext(ctx, network, addr)
	if err != nil {
		return err
	}
	return conn.Close()
}
