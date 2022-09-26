// Package fcw wraps a Firefox process from start-to-finish, and allows
// the user to pass a profile directory to that process where it will find
// pre-configured settings. It's useful for using Firefox as an interface to
// applications, and for applying specific settings for use when browsing
// overlay networks like Tor or I2P.
package fcw

/**
 // MIT License

 // Copyright (c) 2018 Serge Zaitsev

 // Permission is hereby granted, free of charge, to any person obtaining a copy
 // of this software and associated documentation files (the "Software"), to deal
 // in the Software without restriction, including without limitation the rights
 // to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 // copies of the Software, and to permit persons to whom the Software is
 // furnished to do so, subject to the following conditions:

 // The above copyright notice and this permission notice shall be included in all
 // copies or substantial portions of the Software.

 // THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 // IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 // FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 // AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 // LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 // OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 // SOFTWARE.
**/

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
