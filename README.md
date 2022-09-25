# fcw
--
    import "github.com/eyedeekay/go-fpw"


## Usage

```go
var FirefoxExecutable = LocateFirefox
```
FirefoxExecutable returns a string which points to the preferred Firefox
executable file as calculated by the LocateFirefox variable

#### func  LocateFirefox

```go
func LocateFirefox() string
```
LocateFirefox returns a path to the Firefox binary, or an empty string if
Firefox installation is not found.

#### func  MessageBox

```go
func MessageBox(title, text string) bool
```
MessageBox creates a dialog box which prompts the user to download and install
Firefox if they have not already.

#### func  PortablePath

```go
func PortablePath() string
```
PortablePath determines if there is a "Portable" Firefox in a sub-directory of
the runtime directory

#### func  PromptDownload

```go
func PromptDownload()
```
PromptDownload asks user if they want to download and install Firefox, and opens
a download web page if the user agrees.

#### func  Run

```go
func Run() error
```
Run creates a basic instance of the Firefox manager with a default profile
directory and launches duckduckgo.com

#### type UI

```go
type UI interface {
	Done() <-chan struct{}
	Close() error
	Log() string
}
```

UI is a wrapper/manager for a Firefox external process.

#### func  BasicFirefox

```go
func BasicFirefox(userdir string, private bool, args ...string) (UI, error)
```
BasicFirefox sets up a new Firefox instance, and creates the profile directory
if it does not already exist.

#### func  NewFirefox

```go
func NewFirefox(url, dir string, width, height int, customArgs ...string) (UI, error)
```
NewFirefox creates a new instance of the Firefox manager.
