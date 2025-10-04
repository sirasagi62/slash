# Project Plan: `slash` – Go Command-Line Utility

## Overview
`slash` is a Go-based command-line program designed to perform text processing and file system operations using a simple, extensible subcommand architecture (similar to `git`). The tool will support a set of core subcommands and provide a framework for adding future functionality.

## Goals
- Provide a fast, compiled binary with minimal dependencies.
- Offer a clear, user‑friendly CLI with help documentation.
- Implement a modular design to simplify adding new subcommands.
- Ensure cross‑platform compatibility (Linux, macOS, Windows).

## Core Features (Version 1.0)
1. **`slash version`** – Print the program version.
2. **`slash help`** – Show usage information and list available subcommands.
3. **`slash <note> [args...]`** – Display a markdown note from `~/.config/slash` or `./.slash`. If additional arguments are provided, the note is processed as a Go `text/template` with those arguments.

## Architecture
- **`main.go`** – Entry point; parses top‑level flags and dispatches to subcommand handlers.
- **`cmd/` package** – Contains a subpackage for each subcommand (`version`, `help`, `count`, `replace`, `grep`). Each subcommand implements a common interface:
  ```go
  type Command interface {
      Name() string
      Synopsis() string
      Run(args []string) error
  }
  ```
- **`internal/util/`** – Helper functions for file I/O, error handling, and common utilities.
- **`go.mod`** – Module definition, specifying Go version and any third‑party dependencies (e.g., `github.com/spf13/pflag` for flag parsing, `regexp` from the standard library).

## Development Steps
1. **Initialize the module**
   ```bash
   go mod init github.com/yourusername/slash
   ```
2. **Create project skeleton**
   - `main.go`
   - `cmd/` directory with placeholder files for each subcommand.
   - `internal/util/` for shared helpers.
3. **Implement command interface and registration**
   - Build a map of command name → `Command` implementation.
4. **Develop each subcommand**
   - Write unit tests for core logic (e.g., counting, replace, grep).
   - Ensure proper error messages for missing arguments or file errors.
5. **Add help and version output**
   - Use `embed` package to include a static `VERSION` file.
6. **Testing**
   - Run `go test ./...` to verify all packages.
   - Perform manual CLI testing on all supported platforms.
7. **Packaging**
   - Create a simple Makefile or GoReleaser config for building releases.
   - Provide installation instructions (binary download, `go install`).

## Timeline (Estimated)
| Week | Milestone                              |
|------|----------------------------------------|
| 1    | Project initialization, CLI skeleton  |
| 2    | Implement `version` and `help` commands|
| 3    | Implement `count` and unit tests       |
| 4    | Implement `replace` and `grep` commands|
| 5    | Comprehensive testing and CI setup     |
| 6    | Documentation, packaging, release prep |

## Future Enhancements
- Plugin system for third‑party subcommands.
- Configurable output formats (JSON, CSV).
- Parallel processing for large files.
- Integration with `cobra` for richer CLI features.

---  

*Prepared by the development team on `$(date)`.*
