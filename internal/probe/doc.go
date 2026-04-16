// Package probe implements concurrent TCP probing with latency measurement.
//
// Basic usage:
//
//	p := probe.New(2 * time.Second)
//
//	// Single probe
//	r := p.Probe("api", "localhost:8080")
//	fmt.Println(r.Up, r.Latency)
//
//	// Batch probe (concurrent)
//	results := p.Batch([]probe.Target{
//		{Name: "api",  Addr: "localhost:8080"},
//		{Name: "db",   Addr: "localhost:5432"},
//	})
//
// Each Result records whether the target was reachable, the round-trip
// latency, and the timestamp of the attempt.
package probe
