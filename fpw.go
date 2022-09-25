package fcw

import (
	"log"
	"os"
	"path/filepath"
)

// BasicFirefox sets up a new Firefox instance, and creates the profile directory if
// it does not already exist.
func BasicFirefox(userdir string, private bool, args ...string) (UI, error) {
	os.MkdirAll(userdir, os.ModePerm)
	add := true
	var cleanedArgs []string
	if private {
		cleanedArgs = append(cleanedArgs, "--private-window")
	}
	for _, arg := range args {
		if arg == "--private-window" {
			if private {
				add = false
			} else {
				add = true
			}
		}
		if add {
			cleanedArgs = append(cleanedArgs, arg)
		}
	}
	log.Println("Args", cleanedArgs)
	userdir, err := filepath.Abs(userdir)
	if err != nil {
		return nil, err
	}
	return NewFirefox("", userdir, 800, 600, cleanedArgs...)
}

// Run creates a basic instance of the Firefox manager with a default profile directory and
// launches duckduckgo.com
func Run() error {
	var FIREFOX, ERROR = BasicFirefox("basic", true, "https://duckduckgo.com")
	if ERROR != nil {
		return ERROR
	}
	defer FIREFOX.Close()
	<-FIREFOX.Done()
	return nil
}
