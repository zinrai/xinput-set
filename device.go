package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Device represents an input device
type Device struct {
	ID   string
	Name string
}

// FindValidDevice searches for a device matching the filter and validation
func FindValidDevice(filter Filter, validation Validation) (*Device, error) {
	// Get all matching devices
	devices, err := FindMatchingDevices(filter)
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, nil
	}

	fmt.Printf("Found %d device(s) matching \"%s\"\n", len(devices), filter.NameContains)

	// Try to find a valid device among the matches
	for _, device := range devices {
		if ValidateDevice(device.ID, validation) {
			fmt.Printf("  Device id=%s: validation passed\n", device.ID)
			return &device, nil
		}

		// Get missing properties for error message
		missingProps := GetMissingProperties(device.ID, validation.HasProperties)
		if len(missingProps) > 0 {
			fmt.Printf("  Device id=%s: missing required properties %v\n", device.ID, missingProps)
		} else {
			fmt.Printf("  Device id=%s: validation failed\n", device.ID)
		}
	}

	return nil, nil
}

// FindMatchingDevices returns all devices matching the filter
func FindMatchingDevices(filter Filter) ([]Device, error) {
	devices, err := ListDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	var matches []Device
	for _, device := range devices {
		if strings.Contains(device.Name, filter.NameContains) {
			matches = append(matches, device)
		}
	}

	return matches, nil
}

// ListDevices returns all input devices from xinput
func ListDevices() ([]Device, error) {
	output, err := RunCommand("xinput", "list")
	if err != nil {
		return nil, fmt.Errorf("failed to run xinput list: %w", err)
	}

	var devices []Device
	lines := strings.Split(string(output), "\n")
	idRegex := regexp.MustCompile(`id=(\d+)`)

	for _, line := range lines {
		device := parseDeviceLine(line, idRegex)
		if device == nil {
			continue
		}
		devices = append(devices, *device)
	}

	return devices, nil
}

func parseDeviceLine(line string, idRegex *regexp.Regexp) *Device {
	matches := idRegex.FindStringSubmatch(line)
	if len(matches) <= 1 {
		return nil
	}

	deviceID := matches[1]
	deviceName := extractDeviceName(line)
	if deviceName == "" {
		return nil
	}

	return &Device{
		ID:   deviceID,
		Name: deviceName,
	}
}

// extractDeviceName extracts the device name from an xinput list line
func extractDeviceName(line string) string {
	// Remove leading/trailing whitespace
	line = strings.TrimSpace(line)

	// Check if line contains device info
	if !containsDeviceInfo(line) {
		return ""
	}

	// Clean up the line
	line = cleanDeviceLine(line)

	// Extract name (everything before "id=")
	idx := strings.Index(line, "id=")
	if idx <= 0 {
		return ""
	}

	name := strings.TrimSpace(line[:idx])
	// Remove any remaining tab characters and collapse spaces
	name = strings.ReplaceAll(name, "\t", " ")
	name = regexp.MustCompile(`\s+`).ReplaceAllString(name, " ")
	return name
}

func containsDeviceInfo(line string) bool {
	// Skip lines that don't have device information
	hasTreeChar := strings.HasPrefix(line, "↳") || strings.HasPrefix(line, "⎜")
	hasID := strings.Contains(line, "id=")
	return hasTreeChar || hasID
}

func cleanDeviceLine(line string) string {
	// Remove the tree characters
	line = strings.TrimPrefix(line, "↳")
	line = strings.TrimPrefix(line, "⎜")
	return strings.TrimSpace(line)
}

// ValidateDevice checks if the device has all required properties
func ValidateDevice(deviceID string, validation Validation) bool {
	if len(validation.HasProperties) == 0 {
		return true
	}

	properties, err := GetDeviceProperties(deviceID)
	if err != nil {
		return false
	}

	for _, requiredProp := range validation.HasProperties {
		if !hasProperty(properties, requiredProp) {
			return false
		}
	}

	return true
}

func hasProperty(properties []string, requiredProp string) bool {
	for _, prop := range properties {
		if strings.Contains(prop, requiredProp) {
			return true
		}
	}
	return false
}

// GetDeviceProperties returns the list of properties for a device
func GetDeviceProperties(deviceID string) ([]string, error) {
	output, err := RunCommand("xinput", "list-props", deviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to list device properties: %w", err)
	}

	return parseProperties(string(output)), nil
}

func parseProperties(output string) []string {
	lines := strings.Split(output, "\n")
	var properties []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Device") {
			continue
		}
		properties = append(properties, line)
	}

	return properties
}

// GetMissingProperties returns a list of properties that are missing from the device
func GetMissingProperties(deviceID string, requiredProps []string) []string {
	properties, err := GetDeviceProperties(deviceID)
	if err != nil {
		return requiredProps
	}

	var missing []string
	for _, requiredProp := range requiredProps {
		if hasProperty(properties, requiredProp) {
			continue
		}
		missing = append(missing, requiredProp)
	}

	return missing
}
