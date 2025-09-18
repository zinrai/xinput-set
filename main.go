package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	var configPath string
	var dryRun bool

	flag.StringVar(&configPath, "config", "config.yaml", "Path to configuration file")
	flag.BoolVar(&dryRun, "dry-run", false, "Show commands without executing")
	flag.Parse()

	if err := run(configPath, dryRun); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(configPath string, dryRun bool) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	applier := NewApplier(dryRun)
	return applier.ApplyConfig(config)
}
