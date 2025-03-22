# Site-Specific Browser

A command line utility that creates isolated Firefox instances for specific websites, using the fcw (Firefox Controller Wrapper) library.

## Features

- Creates isolated Firefox profiles for each website
- Runs Firefox in "App Mode" for a more focused browsing experience
- Supports private browsing mode
- Automatically manages profile directories
- Portable Firefox support
- Custom window dimensions
- Clean URL-based profile naming

## Installation

```bash
go install github.com/eyedeekay/go-fpw/ssbapp@latest
```

## Usage

Basic usage:

```bash
ssbapp -url "https://example.com"
```

Advanced options:

```bash
ssbapp -url "https://example.com" \
       -profiles "/custom/profile/path" \
       -private
```

### Command Line Flags

- `-url`: The website URL to create a dedicated browser for (required)
- `-profiles`: Base directory for storing profiles (default: ~/.sitebrowsers)
- `-private`: Enable private browsing mode (default: false)
- `-offline`: Enable offline, localhost-only mode (default:false)

## Profile Management

Each website gets its own isolated profile directory, named after the site's hostname (sanitized for filesystem compatibility). Profiles are stored by default in:

- Linux/macOS: `~/.sitebrowsers/`
- Windows: `%USERPROFILE%\.sitebrowsers\`

## Examples

Launch GitHub in its own browser instance:
```bash
ssbapp -url "https://github.com"
```

Create a private browsing instance for DuckDuckGo:
```bash
ssbapp -url "https://duckduckgo.com" -private
```

Use custom profile location:
```bash
ssbapp -url "https://example.com" -profiles "./my-browsers"
```

Use offline mode:
```bash
ssbapp -url "http://localhost:7657" -offline
```

## Dependencies

- Firefox browser installed on the system
- Go 1.16 or later

## License

MIT License - See LICENSE file

## Contributing

Feel free to submit issues and pull requests.