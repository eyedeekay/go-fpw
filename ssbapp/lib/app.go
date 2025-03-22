package ssb

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

func WebAppFunction(startURL, profileBase string, private, offline bool) {
	if startURL == "" {
		fmt.Fprintf(os.Stderr, "Error: -url flag is required\n")
		flag.Usage()
		os.Exit(1)
	}

	// Parse and validate URL
	uri, err := url.Parse(startURL)
	if err != nil {
		log.Fatalf("Invalid URL: %v", err)
	}

	// Create profile directory based on hostname
	profileDir := filepath.Join(profileBase, sanitizeHostname(uri.Hostname()))
	if err := os.MkdirAll(profileDir, 0o755); err != nil {
		log.Fatalf("Failed to create profile directory: %v", err)
	}

	// Check for portable Firefox
	if portablePath := fcw.PortablePath(); portablePath != "" {
		log.Printf("Using portable Firefox installation: %s", portablePath)
	}

	// Create and configure Firefox instance
	ui, err := fcw.WebAppFirefox(profileDir, private, offline, startURL)
	if err != nil {
		log.Fatalf("Failed to start Firefox: %v", err)
	}
	defer ui.Close()

	// Wait for browser to close
	<-ui.Done()
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
