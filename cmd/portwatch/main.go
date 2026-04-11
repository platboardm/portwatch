package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/notifier"
)

var version = "dev"

func main() {
	cfgPath := flag.String("config", "config.yaml", "path to configuration file")
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("portwatch %s\n", version)
		os.Exit(0)
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error loading config: %v\n", err)
		os.Exit(1)
	}

	n := notifier.New(os.Stdout)
	m := monitor.New(cfg, n)

	ctx, stop := signal.NotifyContext(
		signal.NotifyContext(nil, syscall.SIGINT, syscall.SIGTERM),
	)
	defer stop()

	fmt.Printf("portwatch %s starting — watching %d target(s)\n", version, len(cfg.Targets))

	if err := m.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "monitor exited with error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("portwatch stopped")
}
