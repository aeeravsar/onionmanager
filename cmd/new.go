package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"onionmanager/pkg/service"
	"onionmanager/pkg/tor"
)

// New creates a new onion service
func New(path string) error {
	// Create service instance
	svc, err := service.New(path)
	if err != nil {
		return err
	}

	// Check if path already exists
	if _, err := os.Stat(svc.Path); err == nil {
		// Path exists, check if it's empty
		entries, err := os.ReadDir(svc.Path)
		if err != nil {
			return fmt.Errorf("failed to read directory: %w", err)
		}
		if len(entries) > 0 {
			return fmt.Errorf("directory already exists and is not empty: %s", svc.Path)
		}
	}

	// Create directory structure
	if err := os.MkdirAll(svc.Path, 0755); err != nil {
		return fmt.Errorf("failed to create service directory: %w", err)
	}

	if err := os.MkdirAll(svc.DataDir(), 0700); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Prompt for port mappings and create manager.conf
	managerConf, err := promptForManagerConf()
	if err != nil {
		return err
	}

	// Save manager.conf
	if err := os.WriteFile(svc.ManagerConfPath(), []byte(managerConf), 0644); err != nil {
		return fmt.Errorf("failed to save manager.conf: %w", err)
	}

	// Generate torrc
	if err := svc.GenerateTorrc(); err != nil {
		return fmt.Errorf("failed to generate torrc: %w", err)
	}

	// Find Tor binary
	torBinary, err := tor.FindBinary()
	if err != nil {
		return fmt.Errorf("tor binary not found: %w", err)
	}

	// Run Tor to generate keys
	fmt.Println("Generating onion service keys...")
	cmd := exec.Command(torBinary, "-f", svc.TorrcPath())
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start tor: %w", err)
	}

	// Wait for hostname file to be created
	hostnamePath := svc.DataDir() + "/hidden_service/hostname"
	for i := 0; i < 30; i++ {
		if _, err := os.Stat(hostnamePath); err == nil {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}

	// Stop Tor
	if err := cmd.Process.Kill(); err != nil {
		// Process might have already exited
		_ = err
	}
	cmd.Wait()

	// Read the onion address
	onionAddr, err := svc.ReadOnionAddress()
	if err != nil {
		return fmt.Errorf("failed to read onion address: %w", err)
	}

	// Print success message
	fmt.Printf("\nCreated new Onion service at: %s\n", svc.Path)
	fmt.Printf("Onion address: %s\n\n", onionAddr)
	fmt.Printf("Configuration saved to: %s\n", svc.ManagerConfPath())
	fmt.Printf("Edit manager.conf to customize Tor settings.\n\n")
	fmt.Printf("Service ready. Run with: onionmanager run %s\n", path)

	return nil
}

func promptForManagerConf() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var builder strings.Builder

	// Add header comment
	builder.WriteString("# Edit torrc options except HiddenServiceDir and DataDirectory here\n\n")

	portCount := 0
	for {
		// Prompt for virtual port
		fmt.Print("Enter virtual port: ")
		virtualStr, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		virtualStr = strings.TrimSpace(virtualStr)
		virtual, err := strconv.Atoi(virtualStr)
		if err != nil || virtual < 1 || virtual > 65535 {
			fmt.Println("Invalid port number. Must be between 1 and 65535.")
			continue
		}

		// Prompt for local port
		fmt.Print("Enter local port: ")
		localStr, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		localStr = strings.TrimSpace(localStr)
		local, err := strconv.Atoi(localStr)
		if err != nil || local < 1 || local > 65535 {
			fmt.Println("Invalid port number. Must be between 1 and 65535.")
			continue
		}

		// Add to manager.conf
		builder.WriteString(fmt.Sprintf("HiddenServicePort %d 127.0.0.1:%d\n", virtual, local))
		portCount++

		// Ask if user wants to add another
		fmt.Print("Add another port mapping? (y/n): ")
		answer, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "y" && answer != "yes" {
			break
		}
	}

	if portCount == 0 {
		return "", fmt.Errorf("at least one port mapping is required")
	}

	return builder.String(), nil
}
