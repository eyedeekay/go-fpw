# fcw - Firefox Controller Wrapper

A Go library for managing Firefox instances with custom profiles and settings. Useful for creating site-specific browsers, kiosk applications, or configuring Firefox for overlay networks like Tor or I2P.

## Features

- Launch Firefox with custom profiles and settings
- Create WebApp-style Firefox instances 
- Manage Firefox processes and profiles
- Certificate management support
- Private browsing mode
- Support for portable Firefox installations
- Cross-platform (Windows, macOS, Linux)

## Installation

```bash
go get github.com/eyedeekay/go-fpw
```

## Core Functions

### Basic Usage

```go
// Create and launch Firefox with basic profile
ui, err := fcw.BasicFirefox("profile-dir", false, "https://example.com")

// Create a WebApp-style Firefox instance
ui, err := fcw.WebAppFirefox("webapp-profile", false, false, "https://example.com")

// Manage certificates
cm, err := ui.CertManager()
err = cm.AddCertificate("cert.pem", "nickname")
```

### WebApp Mode Features

When using `WebAppFirefox()`, the following customizations are applied:

- Minimal UI with hidden URL bar
- Custom userChrome.css for app-like appearance
- Disabled telemetry and first-run procedures
- Optional offline mode support
- Copy URL to clipboard extension
- User profile customizations enabled

## Site-Specific Browser Application

The package includes `ssbapp`, a command-line utility for creating isolated Firefox instances for specific websites. See [ssbapp documentation](ssbapp/README.md) for details.

Example usage:
```bash
ssbapp -url "https://example.com" -private -profiles "./my-profiles"
```

## License

MIT License - See [LICENSE](LICENSE) file

## Contributing

Contributions welcome! Please feel free to submit issues and pull requests.