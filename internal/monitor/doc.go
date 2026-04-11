// Package monitor orchestrates the port-watching lifecycle.
//
// It spawns one goroutine per target defined in the configuration and
// polls each target on the configured interval. State transitions
// (UP → DOWN, DOWN → UP) are forwarded to the notifier so that
// duplicate alerts are suppressed while the status remains unchanged.
//
// Basic usage:
//
//	cfg, _ := config.Load("portwatch.yaml")
//	c := checker.New(2 * time.Second)
//	n := notifier.New(os.Stdout)
//	m := monitor.New(cfg, c, n)
//	m.Run(ctx)
package monitor
