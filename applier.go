package main

import (
	"fmt"
	"os/exec"
)

// Applier applies xinput configurations
type Applier struct {
	dryRun bool
}

// NewApplier creates a new Applier
func NewApplier(dryRun bool) *Applier {
	return &Applier{
		dryRun: dryRun,
	}
}

// ApplyConfig applies all profiles in the configuration
func (a *Applier) ApplyConfig(config *Config) error {
	for profileName, profile := range config.Profiles {
		if err := a.applyProfile(profileName, profile); err != nil {
			return err
		}
	}
	return nil
}

func (a *Applier) applyProfile(name string, profile Profile) error {
	fmt.Printf("Processing profile: %s\n", name)
	if profile.Description != "" {
		fmt.Printf("Description: %s\n", profile.Description)
	}

	// Find and validate device
	device, err := FindValidDevice(profile.Detection.Filter, profile.Detection.Validation)
	if err != nil {
		return fmt.Errorf("finding valid device: %w", err)
	}
	if device == nil {
		return fmt.Errorf("no valid device found matching \"%s\"",
			profile.Detection.Filter.NameContains)
	}

	fmt.Printf("Selected device: id=%s\n", device.ID)

	// Apply actions
	for _, action := range profile.Actions {
		if err := a.executeAction(device.ID, action); err != nil {
			return fmt.Errorf("applying profile: %w", err)
		}
	}

	if !a.dryRun {
		fmt.Println("Configuration applied successfully.")
	}
	return nil
}

func (a *Applier) executeAction(deviceID string, action Action) error {
	// Parse args string into array
	args := ParseArgs(action.Command, action.Args)

	// Build command: xinput [command] [device_id] [args...]
	cmd := append([]string{"xinput", action.Command, deviceID}, args...)

	// Format command for display
	cmdStr := FormatCommand(cmd)

	if a.dryRun {
		fmt.Printf("[DRY-RUN] Would execute: %s\n", cmdStr)
		return nil
	}

	fmt.Printf("Executing: %s\n", cmdStr)
	if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
		return fmt.Errorf("failed to execute: %s: %w", cmdStr, err)
	}
	return nil
}
