// Package aggregator provides a thread-safe collector that ingests
// [alert.Alert] values from multiple monitoring goroutines and computes
// a rolled-up [Summary] of overall service health.
//
// Typical usage:
//
//	agg := aggregator.New(100) // keep last 100 alerts
//
//	// inside each monitor callback:
//	agg.Record(al)
//
//	// inside a status handler:
//	summary := agg.Summarise()
//	fmt.Printf("%d/%d targets up\n", summary.Up, summary.Total)
//
// The Aggregator is safe for concurrent use by multiple goroutines.
package aggregator
