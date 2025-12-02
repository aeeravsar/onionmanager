package tor

import (
	"fmt"
	"os"
	"os/exec"
)

// FindBinary locates the Tor binary on the system
func FindBinary() (string, error) {
	// Check TOR_BINARY environment variable first
	if torBinary := os.Getenv("TOR_BINARY"); torBinary != "" {
		if _, err := os.Stat(torBinary); err == nil {
			// Check if executable
			if err := exec.Command(torBinary, "--version").Run(); err == nil {
				return torBinary, nil
			}
		}
		return "", fmt.Errorf("TOR_BINARY is set to %s but binary is not executable or not found", torBinary)
	}

	// Default to /usr/bin/tor
	defaultPath := "/usr/bin/tor"
	if _, err := os.Stat(defaultPath); err == nil {
		// Check if executable
		if err := exec.Command(defaultPath, "--version").Run(); err == nil {
			return defaultPath, nil
		}
	}

	return "", fmt.Errorf("tor binary not found at %s", defaultPath)
}
