// Package main is the entry-point for the portwatch CLI daemon.
//
// Usage:
//
//	portwatch [flags]
//
// Flags:
//
//	-config string
//		Path to the YAML configuration file (default "config.yaml").
//	-version
//		Print the build version and exit.
//
// portwatch loads the supplied configuration, wires together the checker,
// notifier, and monitor subsystems, then runs until it receives SIGINT or
// SIGTERM.  All alerts are written to stdout in a structured, human-readable
// format.
package main
