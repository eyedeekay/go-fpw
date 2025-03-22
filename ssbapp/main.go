package main

import (
	"flag"
	"os"
	"path/filepath"

	ssb "github.com/eyedeekay/go-fpw/ssbapp/lib"
)

func main() {
	// Command line flags
	startURL := flag.String("url", "", "URL for the site-specific browser (required)")
	profileBase := flag.String("profiles", getDefaultProfileDir(), "Base directory for profiles")
	private := flag.Bool("private", false, "Use private browsing mode")
	offline := flag.Bool("offline", false, "Use offline mode")

	flag.Parse()

	// Validate URL
	ssb.WebAppFunction(*startURL, *profileBase, *private, *offline)
}

// getDefaultProfileDir returns the default base directory for profiles
func getDefaultProfileDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "profiles")
	}
	return filepath.Join(homeDir, ".sitebrowsers")
}
