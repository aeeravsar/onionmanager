package main

import (
	"fmt"
	"os"

	"onionmanager/cmd"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "new":
		if len(os.Args) != 3 {
			fmt.Println("Usage: onionmanager new <path>")
			os.Exit(1)
		}
		if err := cmd.New(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "run":
		if len(os.Args) != 3 {
			fmt.Println("Usage: onionmanager run <path>")
			os.Exit(1)
		}
		if err := cmd.Run(os.Args[2]); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "help", "-h", "--help":
		printUsage()

	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("onionmanager - Portable Tor Onion Service Manager")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  onionmanager new <path>   Create a new onion service")
	fmt.Println("  onionmanager run <path>   Run an existing onion service")
	fmt.Println("  onionmanager help         Show this help message")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  TOR_BINARY                Path to Tor binary (default: /usr/bin/tor)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  onionmanager new ~/my-service")
	fmt.Println("  onionmanager run ~/my-service")
}
