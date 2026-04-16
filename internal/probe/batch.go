package probe

import (
	"sync"
)

// Target pairs a logical name with a TCP address.
type Target struct {
	Name string
	Addr string
}

// Batch runs Probe concurrently for each target and returns all results.
func (p *Prober) Batch(targets []Target) []Result {
	results := make([]Result, len(targets))
	var wg sync.WaitGroup
	for i, tgt := range targets {
		wg.Add(1)
		go func(idx int, t Target) {
			defer wg.Done()
			results[idx] = p.Probe(t.Name, t.Addr)
		}(i, tgt)
	}
	wg.Wait()
	return results
}
