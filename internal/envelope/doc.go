// Package envelope provides a thin wrapper around [alert.Alert] that attaches
// routing metadata — a priority level, a destination channel name, and a
// unique hex trace ID — so that downstream sinks and middleware can make
// delivery decisions without inspecting alert internals.
//
// # Basic usage
//
//	e := envelope.New(a, envelope.PriorityHigh, "ops")
//	fmt.Println(e.TraceID, e.Channel)
//
// # Routing
//
// A [Router] maps priority thresholds to named channels. Rules are evaluated
// in order from highest to lowest MinPriority; the first matching rule wins.
//
//	rules := []envelope.Rule{
//		{MinPriority: envelope.PriorityCritical, Channel: "pagerduty"},
//		{MinPriority: envelope.PriorityHigh,     Channel: "slack"},
//	}
//	router := envelope.NewRouter(rules, "email")
//	routed := router.Route(e)
package envelope
