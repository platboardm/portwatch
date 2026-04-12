// Package pipeline provides a composable event-processing pipeline that
// connects the individual portwatch stages — debounce, state tracking,
// rate-limiting, and notification — into a single, easy-to-use unit.
//
// # Overview
//
// A Pipeline is created per monitored target.  Each time the scheduler fires
// a check, the raw [checker.Status] result is handed to [Pipeline.Process].
// Internally the result passes through four sequential gates:
//
//  1. Debounce  – transient flaps are ignored until N consecutive identical
//     results are observed.
//  2. State     – the new status is compared with the last known status; if
//     nothing changed the pipeline short-circuits.
//  3. Rate-limit – even genuine transitions are suppressed if an alert was
//     already sent within the configured cooldown window.
//  4. Notify    – a formatted alert is written to the configured [io.Writer].
//
// # Usage
//
//	p := pipeline.New("api:8080", pipeline.Config{
//	    DebounceThreshold: 2,
//	    AlertCooldown:     5 * time.Minute,
//	}, os.Stdout)
//
//	p.Process(ctx, "api:8080", checker.StatusDown)
package pipeline
