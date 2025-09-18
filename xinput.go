package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// RunCommand executes a command and returns its output
func RunCommand(command string, args ...string) ([]byte, error) {
	// Always show what command is being executed
	fmt.Printf("Executing: %s %s\n", command, strings.Join(args, " "))

	cmd := exec.Command(command, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("command failed: %w\nstderr: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

// CheckXinput verifies that xinput is available
func CheckXinput() error {
	_, err := exec.LookPath("xinput")
	if err != nil {
		return fmt.Errorf("xinput command not found. Please install xinput and try again")
	}
	return nil
}

// init checks for xinput availability when the program starts
func init() {
	if err := CheckXinput(); err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
}
