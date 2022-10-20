# fcw

Package fcw wraps a Firefox process from start-to-finish, and allows
the user to pass a profile directory to that process where it will find
pre-configured settings. It's useful for using Firefox as an interface to
applications, and for applying specific settings for use when browsing
overlay networks like Tor or I2P.

## Variables

FirefoxExecutable returns a string which points to the preferred Firefox
executable file as calculated by the LocateFirefox variable

```golang
var FirefoxExecutable = LocateFirefox
```

## Functions

### func [LocateFirefox](/ui.go#L214)

`func LocateFirefox() string`

LocateFirefox returns a path to the Firefox binary, or an empty string if
Firefox installation is not found.

### func [MessageBox](/messagebox.go#L40)

`func MessageBox(title, text string) bool`

MessageBox creates a dialog box which prompts the user to download and install Firefox if they
have not already.

### func [PortablePath](/ui.go#L53)

`func PortablePath() string`

PortablePath determines if there is a "Portable" Firefox in a sub-directory
of the runtime directory

### func [PromptDownload](/ui.go#L243)

`func PromptDownload()`

PromptDownload asks user if they want to download and install Firefox, and
opens a download web page if the user agrees.

### func [Run](/fpw.go#L270)

`func Run() error`

Run creates a basic instance of the Firefox manager with a default profile directory and
launches duckduckgo.com

### func [UnpackApp](/fpw.go#L111)

`func UnpackApp(profileDir string) (string, error)`

UnpackApp unpacks a "App" mode profile into the "profileDir" and returns the
path to the profile and possibly, an error if something goes wrong. If everything
works, the error will be nil

## Types

### type [UI](/ui.go#L41)

`type UI interface { ... }`

UI is a wrapper/manager for a Firefox external process.

#### func [BasicFirefox](/fpw.go#L42)

`func BasicFirefox(userdir string, private bool, args ...string) (UI, error)`

BasicFirefox sets up a new Firefox instance, and creates the profile directory if
it does not already exist.

#### func [NewFirefox](/ui.go#L312)

`func NewFirefox(url, dir string, width, height int, customArgs ...string) (UI, error)`

NewFirefox creates a new instance of the Firefox manager.

#### func [WebAppFirefox](/fpw.go#L74)

`func WebAppFirefox(userdir string, private bool, args ...string) (UI, error)`

WebAppFirefox sets up a new Firefox instance, and creates the profile directory if
it does not already exist. It turns Firefox into a WebApp-Viewer with the provided
profile

## Sub Packages

* [check](./check)

---
Readme created from Go doc with [goreadme](https://github.com/posener/goreadme)
