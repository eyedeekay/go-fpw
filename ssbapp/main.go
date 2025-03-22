package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	fcw "github.com/eyedeekay/go-fpw"
)

func main() {
	// Command line flags
	startURL := flag.String("url", "", "URL for the site-specific browser (required)")
	profileBase := flag.String("profiles", getDefaultProfileDir(), "Base directory for profiles")
	private := flag.Bool("private", false, "Use private browsing mode")

	flag.Parse()

	// Validate URL
	if *startURL == "" {
		fmt.Fprintf(os.Stderr, "Error: -url flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Parse and validate URL
	uri, err := url.Parse(*startURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	// Create profile directory based on hostname
	profileDir := filepath.Join(*profileBase, sanitizeHostname(uri.Hostname()))
	if err := os.MkdirAll(profileDir, 0755); err != nil {
		log.Fatalf("Failed to create profile directory: %v", err)
	}

	// Check for portable Firefox
	if portablePath := fcw.PortablePath(); portablePath != "" {
		log.Printf("Using portable Firefox installation: %s", portablePath)
	}

	// Create and configure Firefox instance
	ui, err := fcw.WebAppFirefox(profileDir, *private, *startURL)
	if err != nil {
		log.Fatalf("Failed to start Firefox: %v", err)
	}
	defer ui.Close()

	// Wait for browser to close
	<-ui.Done()
}

// getDefaultProfileDir returns the default base directory for profiles
func getDefaultProfileDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "profiles")
	}
	return filepath.Join(homeDir, ".sitebrowsers")
}

// sanitizeHostname makes the hostname safe for use as a directory name
func sanitizeHostname(hostname string) string {
	// Replace potentially problematic characters
	replacer := strings.NewReplacer(
		":", "_",
		"/", "_",
		"\\", "_",
		"?", "_",
		"*", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		" ", "_",
	)
	return replacer.Replace(hostname)
}
