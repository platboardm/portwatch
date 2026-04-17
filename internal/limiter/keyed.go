package limiter

// KeyedGuard wraps Limiter to provide RAII-style acquire/release via a
// returned release function, reducing the chance of forgotten Release calls.
type KeyedGuard struct {
	l *Limiter
}

// NewKeyed returns a KeyedGuard backed by a Limiter with the given max.
func NewKeyed(max int) *KeyedGuard {
	return &KeyedGuard{l: New(max)}
}

// Acquire attempts to acquire a slot for key.
// On success it returns a no-arg release function and a nil error.
// On failure it returns nil and ErrLimitReached.
func (g *KeyedGuard) Acquire(key string) (release func(), err error) {
	if err = g.l.Acquire(key); err != nil {
		return nil, err
	}
	once := false
	return func() {
		if !once {
			once = true
			g.l.Release(key)
		}
	}, nil
}

// Inflight delegates to the underlying Limiter.
func (g *KeyedGuard) Inflight(key string) int {
	return g.l.Inflight(key)
}
