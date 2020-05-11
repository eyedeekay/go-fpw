//+build !windows

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
	"os/exec"
	"runtime"
	"strings"
	"syscall"
)

func MessageBox(title, text string) bool {
	if runtime.GOOS == "linux" {
		err := exec.Command("zenity", "--question", "--title", title, "--text", text).Run()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				return exitError.Sys().(syscall.WaitStatus).ExitStatus() == 0
			}
		}
	} else if runtime.GOOS == "darwin" {
		script := `set T to button returned of ` +
			`(display dialog "%s" with title "%s" buttons {"No", "Yes"} default button "Yes")`
		out, err := exec.Command("osascript", "-e", fmt.Sprintf(script, text, title)).Output()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				return exitError.Sys().(syscall.WaitStatus).ExitStatus() == 0
			}
		}
		return strings.TrimSpace(string(out)) == "Yes"
	}
	return false
}
