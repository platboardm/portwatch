// Package config provides loading and validation of portwatch configuration
// files written in YAML format.
//
// A configuration file lists one or more targets, each describing a host/port
// pair to monitor along with polling interval and connection timeout.
//
// Example usage:
//
//	cfg, err := config.Load("portwatch.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, t := range cfg.Targets {
//		fmt.Printf("monitoring %s on %s:%d\n", t.Name, t.Host, t.Port)
//	}
package config
