// Package heartbeat provides a periodic liveness signal for the portwatch
// monitor loop.
//
// An Emitter is created with a fixed interval and any number of subscribers
// can register via Subscribe. Each subscriber receives a Beat value on every
// tick. Slow subscribers are never blocked; a beat is silently dropped if the
// subscriber's channel is full.
//
// Usage:
//
//	e := heartbeat.New(10 * time.Second)
//	ch := e.Subscribe()
//	go e.Run(ctx)
//	for b := range ch {
//		log.Printf("alive at %s", b.At)
//	}
package heartbeat
