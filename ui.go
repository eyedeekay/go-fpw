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
	"runtime"
	"strings"

	"sync"
)

type UI interface {
	Done() <-chan struct{}
	Close() error
}

// FirefoxExecutable returns a string which points to the preferred Firefox
// executable file.
var FirefoxExecutable = LocateFirefox

// LocateFirefox returns a path to the Firefox binary, or an empty string if
// Firefox installation is not found.
func LocateFirefox() string {

	// If env variable "LORCACHROME" specified and it exists
	if path, ok := os.LookupEnv("FIREFOX_BIN"); ok {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	var paths []string
	switch runtime.GOOS {
	case "darwin":
		paths = []string{
			"/Applications/Moxilla Firefox.app/Contents/MacOS/Mozilla Firefox",
			"/Applications/Firefox.app/Contents/MacOS/Mozilla Firefox",
			"/usr/bin/firefox-esr",
			"/usr/bin/firefox",
			"/usr/bin/icecat",
		}
	case "windows":
		paths = []string{
			os.Getenv("LocalAppData") + "/Mozilla Firefox/firefox.exe",
			os.Getenv("ProgramFiles") + "/Mozilla Firefox/firefox.exe",
			os.Getenv("ProgramFiles(x86)") + "/Mozilla Firefox/firefox.exe",
			os.Getenv("LocalAppData") + "/GNU Icecat/icecat.exe",
			os.Getenv("ProgramFiles") + "/GNU Icecat/icecat.exe",
			os.Getenv("ProgramFiles(x86)") + "/GNU Icecat/icecat.exe",
		}
	default:
		paths = []string{
			"/usr/bin/firefox-esr",
			"/usr/bin/firefox",
			"/usr/bin/waterfox",
			"/usr/bin/icecat",
			"/usr/bin/purebrowser",
		}
	}

	for _, path := range paths {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			continue
		}
		return path
	}
	return ""
}

// PromptDownload asks user if he wants to download and install Firefox, and
// opens a download web page if the user agrees.
func PromptDownload() {
	title := "Firefox not found"
	text := "No Firefox installation was found. Would you like to download and install it now?"

	// Ask user for confirmation
	if !messageBox(title, text) {
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

func NewFirefox(url, dir string, width, height int, customArgs ...string) (UI, error) {
	if url == "" {
		url = "about:blank"
	}
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
	c := &firefox{
		id: 2,
	}

	// Start firefox process
	c.cmd = exec.Command(firefoxBinary, args...)
	//pipe, err := c.cmd.StderrPipe()
	//if err != nil {
	//return nil, err
	//}
	if err := c.cmd.Start(); err != nil {
		return nil, err
	}

	// Wait for websocket address to be printed to stderr
	/*re := regexp.MustCompile(`^DevTools listening on (ws://.*?)\r?\n$`)
	m, err := readUntilMatch(pipe, re)
	if err != nil {
		c.kill()
		return nil, err
	}
	wsURL := m[1]

	// Open a websocket
	c.ws, err = websocket.Dial(wsURL, "", "http://127.0.0.1")
	if err != nil {
		c.kill()
		return nil, err
	}

	// Find target and initialize session
	c.target, err = c.findTarget()
	if err != nil {
		c.kill()
		return nil, err
	}

	c.session, err = c.startSession(c.target)
	if err != nil {
		c.kill()
		return nil, err
	}
	go c.readLoop()
	for method, args := range map[string]h{
		"Page.enable":          nil,
		"Target.setAutoAttach": {"autoAttach": true, "waitForDebuggerOnStart": false},
		"Network.enable":       nil,
		"Runtime.enable":       nil,
		"Security.enable":      nil,
		"Performance.enable":   nil,
		"Log.enable":           nil,
	} {
		if _, err := c.send(method, args); err != nil {
			c.kill()
			c.cmd.Wait()
			return nil, err
		}
	}

	if !contains(args, "--headless") {
		win, err := c.getWindowForTarget(c.target)
		if err != nil {
			c.kill()
			return nil, err
		}
		c.window = win.WindowID
	}*/

	return c, nil
}
