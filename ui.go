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
	"strconv"
	"strings"
	"sync"
)

// UI is a wrapper/manager for a Firefox external process.
type UI interface {
	Done() <-chan struct{}
	Close() error
	Log() string
	CertManager() (*CertManager, error)
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
	// ws       *websocket.Conn
	id          int32
	target      string
	session     string
	window      int
	certManager *CertManager
	profileDir  string
	// pending  map[int]chan result
}

type ui struct {
	*firefox
	done   chan struct{}
	tmpDir string
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
	defer DeAppifyUserJS(u.firefox.profileDir)
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

func trimBlankArgs(args []string) (trimmedArgs []string) {
	for _, v := range args {
		if v != "" {
			trimmedArgs = append(trimmedArgs, v)
		}
	}
	return
}

func randir() string {
	return "1"
}

func directory(dir string) string {
	i := 0
	if file, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0o755)
	} else {
		if !file.IsDir() {
			for {
				dir = dir + "-" + strconv.Itoa(i)
				i++
				if _, err := os.Stat(dir); os.IsNotExist(err) {
					os.MkdirAll(dir, 0o755)
					return dir
				}
			}
		}
	}
	return dir
}

func (f *firefox) CertManager() (*CertManager, error) {
	if f.certManager == nil {
		cm, err := NewCertManager(f.target)
		if err != nil {
			return nil, err
		}
		f.certManager = cm
	}
	return f.certManager, nil
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
	} else {
		dir = directory(dir)
	}
	args := append(firefoxArgs, "--profile")
	args = append(args, dir)
	args = append(args, "--window-size")
	args = append(args, fmt.Sprintf("%d,%d", width, height))
	args = append(args, customArgs...)
	args = append(args, url)
	args = trimBlankArgs(args)
	// args = append(args, "--remote-debugging-port=0")
	log.Println(FirefoxExecutable(), args)

	firefox, err := newFirefoxWithArgs(FirefoxExecutable(), args...)
	done := make(chan struct{})
	if err != nil {
		return nil, err
	}
	firefox.profileDir = dir

	go func() {
		firefox.cmd.Wait()
		close(done)
	}()
	return &ui{
		firefox: firefox,
		done:    done,
		tmpDir:  tmpDir,
	}, nil
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
