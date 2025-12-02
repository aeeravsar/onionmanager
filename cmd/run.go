package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"onionmanager/pkg/service"
	"onionmanager/pkg/tor"
)

// Run starts an existing onion service
func Run(path string) error {
	// Create service instance
	svc, err := service.New(path)
	if err != nil {
		return err
	}

	// Validate that service exists
	if _, err := os.Stat(svc.Path); os.IsNotExist(err) {
		return fmt.Errorf("service does not exist: %s", svc.Path)
	}

	// Check if manager.conf exists
	if _, err := os.Stat(svc.ManagerConfPath()); os.IsNotExist(err) {
		return fmt.Errorf("manager.conf not found: %s", svc.ManagerConfPath())
	}

	// Check if data directory exists
	if _, err := os.Stat(svc.DataDir()); os.IsNotExist(err) {
		return fmt.Errorf("data directory not found: %s", svc.DataDir())
	}

	// Regenerate torrc with current absolute paths
	if err := svc.GenerateTorrc(); err != nil {
		return fmt.Errorf("failed to generate torrc: %w", err)
	}

	// Find Tor binary
	torBinary, err := tor.FindBinary()
	if err != nil {
		return fmt.Errorf("tor binary not found: %w", err)
	}

	// Read onion address
	onionAddr, err := svc.ReadOnionAddress()
	if err != nil {
		return fmt.Errorf("failed to read onion address: %w", err)
	}

	// Print startup info
	fmt.Printf("Starting Onion service from: %s\n", svc.Path)
	fmt.Printf("Onion address: %s\n", onionAddr)
	fmt.Println()

	// Start Tor
	cmd := exec.Command(torBinary, "-f", svc.TorrcPath())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start tor: %w", err)
	}

	// Setup signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Wait for either the process to exit or a signal
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-sigChan:
		fmt.Println("\nShutting down...")
		cmd.Process.Signal(syscall.SIGTERM)
		cmd.Wait()
	case err := <-done:
		if err != nil {
			fmt.Printf("Tor exited with error: %v\n", err)
		}
	}

	return nil
}
