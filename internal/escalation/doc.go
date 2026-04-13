// Package escalation tracks continuous outage durations and signals when a
// target has been down long enough to warrant an escalated response.
//
// # Overview
//
// An [Escalator] is created with a threshold duration. Each time an [alert.Alert]
// is evaluated, the escalator records when the target first became critical. If
// the target remains critical beyond the threshold, Evaluate returns true,
// indicating that the situation should be escalated (e.g. paging on-call staff
// or posting to a high-priority webhook channel).
//
// Recovery events (non-critical alerts) automatically clear the tracked downtime
// for the affected target.
//
// # Example
//
//	esk := escalation.New(30 * time.Minute)
//
//	if esk.Evaluate(a) {
//		// target has been down for > 30 minutes — escalate
//	}
package escalation
