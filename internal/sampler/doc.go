// Package sampler provides probabilistic sampling for alert events in portwatch.
//
// # Overview
//
// When alert volume is high, it may be desirable to forward only a
// representative fraction of events to downstream sinks. The sampler
// package provides two types for this purpose:
//
//   - [Sampler] — a single sampler with a fixed rate applied globally.
//   - [KeyedSampler] — maintains independent samplers per key (e.g. per
//     target name), so each service is sampled independently.
//
// # Usage
//
//	s := sampler.New(0.25) // allow ~25% of events
//	if s.Allow() {
//	    sink.Send(alert)
//	}
//
// Rates are expressed as a float64 in the range [0.0, 1.0].
// A rate of 1.0 passes all events; 0.0 suppresses all events.
package sampler
