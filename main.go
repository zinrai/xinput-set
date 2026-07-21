package main

import (
	"flag"
	"fmt"
	"os"
)

// Injected at build time by goreleaser ldflags.
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	var configPath string
	var dryRun bool
	var showVersion bool

	flag.StringVar(&configPath, "config", "config.yaml", "Path to configuration file")
	flag.BoolVar(&dryRun, "dry-run", false, "Show commands without executing")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.Parse()

	if showVersion {
		fmt.Printf("xinput-set %s (commit %s, built %s)\n", version, commit, date)
		return
	}

	if err := run(configPath, dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath string, dryRun bool) error {
	if err := CheckXinput(); err != nil {
		return err
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	applier := NewApplier(dryRun)
	return applier.ApplyConfig(config)
}
