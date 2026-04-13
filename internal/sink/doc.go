// Package sink implements a concurrent fan-out dispatcher that forwards
// alert.Alert values to one or more named Sender destinations.
//
// # Overview
//
// A [Sink] holds a map of named [Sender] implementations and dispatches
// each alert to all of them in parallel. Any errors are collected and
// returned as a slice of [SendError] values so the caller can decide how
// to handle partial failures.
//
// [Multi] wraps a Sink for fire-and-forget use: it logs failures via a
// standard *log.Logger and never propagates errors to the caller, making
// it convenient inside pipelines where alerting is best-effort.
//
// # Example
//
//	senders := map[string]sink.Sender{
//		"webhook": webhookClient,
//		"stdout":  notifierClient,
//	}
//	s := sink.New(senders)
//	if errs := s.Dispatch(ctx, a); len(errs) > 0 {
//		for _, e := range errs {
//			log.Println(e)
//		}
//	}
package sink
