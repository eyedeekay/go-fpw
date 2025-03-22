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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// FirefoxExecutable returns a string which points to the preferred Firefox
// executable file as calculated by the LocateFirefox variable
var FirefoxExecutable = LocateFirefox

// PortablePath determines if there is a "Portable" Firefox in a sub-directory
// of the runtime directory
func PortablePath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Println("An error was encountered detecting the portable path", err)
	}
	listing, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Println(err)
	}
	for _, appdir := range listing {
		if appdir.IsDir() {
			for _, exe := range portableFiles() {
				path := filepath.Join(dir, appdir.Name(), exe)
				if _, err := os.Stat(path); os.IsNotExist(err) {
					continue
				}
				log.Println(path)
				return path
			}
		}
	}
	return "false"
}

func portableFiles() []string {
	var paths []string
	switch runtime.GOOS {
	case "windows":
		paths = []string{
			"firefox.exe",
			"librewolf.exe",
			"waterfox.exe",
			"icecat.exe",
		}
	default:
		paths = []string{
			"firefox-esr",
			"firefox",
			"librewolf",
			"waterfox",
			"icecat",
			"purebrowser",
		}
	}
	return paths
}

func defaultPaths() []string {
	var paths []string
	switch runtime.GOOS {
	case "windows":
		dirs := windowsDefaultPaths()
		exes := portableFiles()
		for _, dir := range dirs {
			for _, exe := range exes {
				paths = append(paths, filepath.Join(dir, exe))
			}
		}
	case "darwin":
		dirs := darwinDefaultPaths()
		exes := portableFiles()
		for _, dir := range dirs {
			for _, exe := range exes {
				paths = append(paths, filepath.Join(dir, exe))
			}
		}
	default:
		dirs := linuxDefaultPaths()
		exes := portableFiles()
		for _, dir := range dirs {
			for _, exe := range exes {
				paths = append(paths, filepath.Join(dir, exe))
			}
		}
	}
	return paths
}

func optBins() []string {
	fi, err := ioutil.ReadDir("/opt/")
	if err != nil {
		log.Println(err.Error())
		return []string{""}
	}
	var optbins []string
	for _, f := range fi {
		if f.IsDir() {
			if strings.HasSuffix(f.Name(), "bin") {
				optbins = append(optbins, f.Name())
			}
		}
	}
	return optbins
}

func windowsDefaultPaths() []string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	programFiles := os.Getenv("ProgramFiles")
	// String localAppData = System.getenv("LOCALAPPDATA");
	// Is there some way Mozilla does adminless installs to LocalAppData? Don't
	// know for sure.
	programFiles86 := os.Getenv("ProgramFiles(x86)")

	tbPath := []string{
		filepath.Join(userHome, "/OneDrive/Desktop/Tor Browser/Browser/"),
		filepath.Join(userHome, "/Desktop/Tor Browser/Browser/"),
	}

	paths := []string{
		tbPath[0],
		tbPath[1],
		filepath.Join(programFiles, "Mozilla Firefox/"),
		filepath.Join(programFiles86, "Mozilla Firefox/"),
		filepath.Join(programFiles, "Waterfox/"),
		filepath.Join(programFiles86, "Waterfox/"),
		filepath.Join(programFiles, "Librewolf/"),
		filepath.Join(programFiles86, "Librewolf/"),
	}

	return paths
}

func linuxDefaultPaths() []string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	pathvar := os.Getenv("PATH")
	elements := strings.Split(pathvar, ":")
	additionalelements := []string{"/opt/bin", filepath.Join(userHome, "bin")}
	optbins := optBins()
	elements = append(elements, additionalelements...)
	elements = append(elements, optbins...)
	return elements
}

func darwinDefaultPaths() []string {
	userHome, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	elements := []string{
		"/Applications/Firefox.app/Contents/MacOS/",
		"/Applications/Waterfox.app/Contents/MacOS/",
		"/Applications/Librewolf.app/Contents/MacOS/",
	}
	pathvar := os.Getenv("PATH")
	pathelements := strings.Split(pathvar, ":")
	additionalelements := []string{"/opt/bin", filepath.Join(userHome, "bin")}
	optbins := optBins()
	elements = append(elements, pathelements...)
	elements = append(elements, additionalelements...)
	elements = append(elements, optbins...)
	return elements
}

// LocateFirefox returns a path to the Firefox binary, or an empty string if
// Firefox installation is not found.
func LocateFirefox() string {
	// If env variable "FIREFOX_BIN" specified and it exists
	if path, ok := os.LookupEnv("FIREFOX_BIN"); ok {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	paths := defaultPaths()

	for _, path := range paths {
		// for _, exe := range exes {

		if info, err := os.Stat(path); os.IsNotExist(err) {
			// err != nil {
			// log.Println(exepath, err)
			continue
		} else {
			if !info.IsDir() {
				log.Println(path)
				return path
			}
		}
		//}
	}
	return ""
}
