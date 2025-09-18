package main

import (
	"fmt"
	"strings"
)

// ParseArgs parses the argument string based on the command type
func ParseArgs(command, args string) []string {
	if command == "set-button-map" {
		return strings.Fields(args)
	}

	if command != "set-prop" {
		return strings.Fields(args)
	}

	// Handle set-prop command
	return parseSetPropArgs(args)
}

func parseSetPropArgs(args string) []string {
	parts := strings.Fields(args)
	if len(parts) <= 1 {
		return []string{args}
	}

	// Find where the numeric values start
	firstValueIndex := findFirstNumericIndex(parts)
	if firstValueIndex == -1 {
		return []string{args}
	}

	// Split into property name and values
	propName := strings.Join(parts[:firstValueIndex], " ")
	values := parts[firstValueIndex:]

	result := []string{propName}
	result = append(result, values...)
	return result
}

func findFirstNumericIndex(parts []string) int {
	for i, part := range parts {
		if isNumeric(part) {
			return i
		}
	}
	return -1
}

func isNumeric(s string) bool {
	if s == "" {
		return false
	}

	// Remove leading minus sign if present
	if s[0] == '-' {
		s = s[1:]
	}

	if s == "" {
		return false
	}

	// Check if remaining string is numeric (with possible decimal point)
	hasDot := false
	for _, c := range s {
		if c == '.' {
			if hasDot {
				return false // Multiple dots
			}
			hasDot = true
			continue
		}
		if c < '0' || c > '9' {
			return false
		}
	}

	return true
}

// FormatCommand formats a command slice for display
func FormatCommand(cmd []string) string {
	formatted := make([]string, len(cmd))

	for i, arg := range cmd {
		formatted[i] = formatArgument(i, arg)
	}

	return strings.Join(formatted, " ")
}

func formatArgument(index int, arg string) string {
	// First 3 arguments: xinput, subcommand, device ID
	if index <= 2 {
		return arg
	}

	// Fourth argument (property name) needs quotes if it contains spaces
	if index == 3 && strings.Contains(arg, " ") {
		return fmt.Sprintf("\"%s\"", arg)
	}

	// Other arguments don't need quotes
	return arg
}
