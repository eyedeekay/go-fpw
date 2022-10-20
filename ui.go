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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"sync"
)

// UI is a wrapper/manager for a Firefox external process.
type UI interface {
	Done() <-chan struct{}
	Close() error
	Log() string
}

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
	elements := []string{"/Applications/Firefox.app/Contents/MacOS/",
		"/Applications/Waterfox.app/Contents/MacOS/",
		"/Applications/Librewolf.app/Contents/MacOS/"}
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
		//for _, exe := range exes {

		if info, err := os.Stat(path); os.IsNotExist(err) {
			//err != nil {
			//log.Println(exepath, err)
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

// PromptDownload asks user if they want to download and install Firefox, and
// opens a download web page if the user agrees.
func PromptDownload() {
	title := "Firefox not found"
	text := "No Firefox installation was found. Would you like to download and install it now?"

	// Ask user for confirmation
	if !MessageBox(title, text) {
		return
	}

	// Open download page
	url := "https://www.mozilla.org/firefox/"
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Run()
	case "darwin":
		exec.Command("open", url).Run()
	case "windows":
		r := strings.NewReplacer("&", "^&")
		exec.Command("cmd", "/c", "start", r.Replace(url)).Run()
	}
}

type firefox struct {
	sync.Mutex
	cmd *exec.Cmd
	//ws       *websocket.Conn
	id      int32
	target  string
	session string
	window  int
	//pending  map[int]chan result
}

type ui struct {
	firefox *firefox
	done    chan struct{}
	tmpDir  string
}

func (u *ui) Log() string {
	out, err := u.firefox.cmd.Output()
	if err != nil {
		return err.Error()
	}
	return string(out)
}

func (u *ui) Done() <-chan struct{} {
	return u.done
}

func (u *ui) Close() error {
	// ignore err, as the firefox process might be already dead, when user close the window.
	u.firefox.kill()
	<-u.done
	if u.tmpDir != "" {
		if err := os.RemoveAll(u.tmpDir); err != nil {
			return err
		}
	}
	return nil
}

var firefoxArgs = []string{
	"--no-remote",
	"--new-instance",
}

// NewFirefox creates a new instance of the Firefox manager.
func NewFirefox(url, dir string, width, height int, customArgs ...string) (UI, error) {
	tmpDir := ""
	if dir == "" {
		name, err := ioutil.TempDir("", "ffox")
		if err != nil {
			return nil, err
		}
		dir, tmpDir = name, name
	}
	args := append(firefoxArgs, "--profile")
	args = append(args, dir)
	args = append(args, "--window-size")
	args = append(args, fmt.Sprintf("%d,%d", width, height))
	args = append(args, customArgs...)
	args = append(args, url)
	//args = append(args, "--remote-debugging-port=0")
	log.Println(FirefoxExecutable(), args)

	firefox, err := newFirefoxWithArgs(FirefoxExecutable(), args...)
	done := make(chan struct{})
	if err != nil {
		return nil, err
	}

	go func() {
		firefox.cmd.Wait()
		close(done)
	}()
	return &ui{firefox: firefox, done: done, tmpDir: tmpDir}, nil
}

func (c *firefox) kill() error {
	if state := c.cmd.ProcessState; state == nil || !state.Exited() {
		return c.cmd.Process.Kill()
	}
	return nil
}

func newFirefoxWithArgs(firefoxBinary string, args ...string) (*firefox, error) {
	// The first two IDs are used internally during the initialization
	if firefoxBinary == "" {
		PromptDownload()
		return nil, fmt.Errorf("Firefox not found.")
	}
	c := &firefox{
		id: 2,
	}

	// Start firefox process
	c.cmd = exec.Command(firefoxBinary, args...)
	if err := c.cmd.Start(); err != nil {
		return nil, err
	}

	return c, nil
}
