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
	"embed"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

//go:embed copy_tab_url_to_clipboard-1.0.xpi
var extraExtension []byte

//go:embed actually_work_offline-2.0-fx.xpi
var offlineExtension []byte

//go:embed userChrome.css
var userChrome []byte

//go:embed user-overrides.js
var userOverrides []byte

var dummy embed.FS

// BasicFirefox sets up a new Firefox instance, and creates the profile directory if
// it does not already exist.
func BasicFirefox(userdir string, private bool, args ...string) (UI, error) {
	userdir = directory(userdir)
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
		if arg != "" {
			if add {
				cleanedArgs = append(cleanedArgs, arg)
			}
		}
	}
	log.Println("Args", cleanedArgs)
	userdir, err := filepath.Abs(userdir)
	if err != nil {
		return nil, err
	}
	return NewFirefox("", userdir, 800, 600, cleanedArgs...)
}

func increment(i *int) int {
	*i++
	return *i
}

// WebAppFirefox sets up a new Firefox instance, and creates the profile directory if
// it does not already exist. It turns Firefox into a WebApp-Viewer with the provided
// profile
func WebAppFirefox(userdir string, private bool, args ...string) (UI, error) {
	// i := 0
	// log.Println(increment(&i))
	userdir = directory(userdir)
	// log.Println(increment(&i))
	add := true
	var cleanedArgs []string
	if private {
		cleanedArgs = append(cleanedArgs, "--private-window")
	}
	// log.Println(increment(&i))
	for _, arg := range args {
		if arg == "--private-window" {
			if private {
				add = false
			} else {
				add = true
			}
		}
		if arg != "" {
			if add {
				cleanedArgs = append(cleanedArgs, arg)
			}
		}
	}
	log.Println("Args", cleanedArgs)
	userdir, err := filepath.Abs(userdir)
	if err != nil {
		return nil, err
	}
	log.Println("Unpacking App" + userdir)
	userdir, err = UnpackApp(userdir)
	if err != nil {
		return nil, err
	}
	log.Println("Unpacked App" + userdir)
	ui, err := NewFirefox("", userdir, 800, 600, cleanedArgs...)
	if err != nil {
		return nil, err
	}
	return ui, nil
}

// UnpackApp unpacks a "App" mode profile into the "profileDir" and returns the
// path to the profile and possibly, an error if something goes wrong. If everything
// works, the error will be nil
func UnpackApp(profileDir string) (string, error) {
	if err := os.MkdirAll(filepath.Join(profileDir, "chrome"), 0o755); err != nil {
		return filepath.Join(profileDir), err
	}
	if err := forceUserChromeCSS(filepath.Join(profileDir, "chrome", "userChrome.css")); err != nil {
		return filepath.Join(profileDir), err
	}
	if err := appifyUserJS(filepath.Join(profileDir, "user-overrides.js")); err != nil {
		return filepath.Join(profileDir), err
	}
	if err := appifyUserJS(filepath.Join(profileDir, "user.js")); err != nil {
		return filepath.Join(profileDir), err
	}
	if err := appifyUserJS(filepath.Join(profileDir, "prefs.js")); err != nil {
		return filepath.Join(profileDir), err
	}
	return profileDir, nil
}

/*


 */

func appifyUserJS(profile string) error {
	prefset := "user_pref(\"extensions.autoDisableScopes\", 0);\n"
	prefset += "user_pref(\"extensions.enabledScopes\", 1);\n"
	prefset += "user_pref(\"toolkit.legacyUserProfileCustomizations.stylesheets\", true);\n"
	extDir := filepath.Join(filepath.Dir(profile), "extensions")
	if err := os.MkdirAll(extDir, 0o755); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(extDir, "{786c38ae-eac8-41df-ad3b-3c737603bead}.xpi"), extraExtension, 0o644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(filepath.Join(extDir, "awo@eyedeekay.github.io.xpi"), offlineExtension, 0o644); err != nil {
		return err
	}
	if _, err := os.Stat(profile); err != nil {
		if strings.Contains(profile, "prefs.js") {
			if err := ioutil.WriteFile(profile, []byte(prefset), 0o644); err != nil {
				return err
			}
		}
	} else {
		content, err := ioutil.ReadFile(profile)
		if err != nil {
			return err
		}
		finished := false
		lines := strings.Split(string(content), "\n")
		for i, line := range lines {
			if strings.Contains(line, "toolkit.legacyUserProfileCustomizations.stylesheets\"") {
				if strings.Contains(line, "true") {
					return nil
				} else {
					line = strings.Replace(line, "false", "true", 1)
					finished = true
				}
			}
			lines[i] = line
		}
		out := strings.Join(lines, "\n")
		if err := ioutil.WriteFile(profile, []byte(out), 0o644); err != nil {
			return err
		}
		if !finished {
			f, err := os.OpenFile(profile,
				os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return err
			}
			defer f.Close()
			if _, err := f.WriteString(prefset); err != nil {
				return err
			}
		}
	}
	return nil
}

func DeAppifyUserJS(profile string) error {
	if _, err := os.Stat(profile); err != nil {
		return nil
	}
	extDir := filepath.Join(profile, "extensions")
	if err := os.Remove(filepath.Join(extDir, "{786c38ae-eac8-41df-ad3b-3c737603bead}.xpi")); err != nil {
		log.Println(err)
	} else {
		log.Println("Removed {786c38ae-eac8-41df-ad3b-3c737603bead}.xpi")
	}
	if err := os.Remove(filepath.Join(extDir, "awo@eyedeekay.github.io.xpi")); err != nil {
		log.Println(err)
	} else {
		log.Println("Removed awp@eyedeekay.github.io.xpi")
	}
	content, err := ioutil.ReadFile(profile)
	if err != nil {
		return err
	}
	lines := strings.Split(string(content), "\n")
	for i, line := range lines {
		if strings.Contains(line, "toolkit.legacyUserProfileCustomizations.stylesheets\"") {
			if strings.Contains(line, "false") {
				return nil
			} else {
				line = strings.Replace(line, "true", "false", 1)
			}
		}
		lines[i] = line
	}
	out := strings.Join(lines, "\n")
	if err := ioutil.WriteFile(profile, []byte(out), 0o644); err != nil {
		return err
	}
	return nil
}

func forceUserChromeCSS(profile string) error {
	if err := ioutil.WriteFile(profile, userChrome, 0o644); err != nil {
		return err
	}
	return nil
}

// Run creates a basic instance of the Firefox manager with a default profile directory and
// launches duckduckgo.com
func Run() error {
	FIREFOX, ERROR := WebAppFirefox("basic", false, "")
	if ERROR != nil {
		return ERROR
	}
	defer FIREFOX.Close()
	<-FIREFOX.Done()
	return nil
}
