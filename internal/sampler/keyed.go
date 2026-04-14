package sampler

import "sync"

// KeyedSampler maintains an independent Sampler per key, allowing
// per-target or per-service sampling rates.
type KeyedSampler struct {
	mu       sync.Mutex
	rate     float64
	samplers map[string]*Sampler
}

// NewKeyed creates a KeyedSampler where each key gets its own Sampler
// initialised with the given rate.
func NewKeyed(rate float64) *KeyedSampler {
	// Validate rate by delegating to New.
	_ = New(rate)
	return &KeyedSampler{
		rate:     rate,
		samplers: make(map[string]*Sampler),
	}
}

// Allow returns true if the event for key should pass through.
// A Sampler is created on first use for each unique key.
func (k *KeyedSampler) Allow(key string) bool {
	k.mu.Lock()
	defer k.mu.Unlock()
	s, ok := k.samplers[key]
	if !ok {
		s = New(k.rate)
		k.samplers[key] = s
	}
	// Unlock before calling Allow to avoid double-locking since Allow
	// acquires its own mutex.
	k.mu.Unlock()
	result := s.Allow()
	k.mu.Lock()
	return result
}

// Keys returns all keys that have been seen so far.
func (k *KeyedSampler) Keys() []string {
	k.mu.Lock()
	defer k.mu.Unlock()
	keys := make([]string, 0, len(k.samplers))
	for key := range k.samplers {
		keys = append(keys, key)
	}
	return keys
}
