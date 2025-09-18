package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the entire configuration file
type Config struct {
	Profiles map[string]Profile `yaml:"profiles"`
}

// Profile represents a device configuration profile
type Profile struct {
	Description string    `yaml:"description"`
	Detection   Detection `yaml:"detection"`
	Actions     []Action  `yaml:"actions"`
}

// Detection contains device detection and validation rules
type Detection struct {
	Filter     Filter     `yaml:"filter"`
	Validation Validation `yaml:"validation"`
}

// Filter defines how to find the device
type Filter struct {
	NameContains string `yaml:"name_contains"`
}

// Validation defines how to validate the device
type Validation struct {
	HasProperties []string `yaml:"has_properties"`
}

// Action represents a configuration action to apply
type Action struct {
	Command string `yaml:"command"`
	Args    string `yaml:"args"`
}

// LoadConfig loads and parses the YAML configuration file
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig performs basic validation on the configuration
func validateConfig(config *Config) error {
	if len(config.Profiles) == 0 {
		return fmt.Errorf("no profiles defined")
	}

	for name, profile := range config.Profiles {
		if err := validateProfile(name, profile); err != nil {
			return err
		}
	}

	return nil
}

func validateProfile(name string, profile Profile) error {
	if profile.Detection.Filter.NameContains == "" {
		return fmt.Errorf("profile %s: filter.name_contains is required", name)
	}

	for i, action := range profile.Actions {
		if err := validateAction(name, i+1, action); err != nil {
			return err
		}
	}

	return nil
}

func validateAction(profileName string, actionIndex int, action Action) error {
	if action.Command == "" {
		return fmt.Errorf("profile %s: action %d has no command", profileName, actionIndex)
	}

	if action.Args == "" {
		return fmt.Errorf("profile %s: action %d has no args", profileName, actionIndex)
	}

	return nil
}
