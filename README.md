# slash

`slash` is a simple Go command‑line utility for managing markdown prompts.
It lets you create, view, and template prompts directly from the terminal.

## Features

- **`slash help`** – Show usage information.
- **`slash edit <name>`** – Create a new markdown prompt in `~/.config/slash`.
- **`slash <prompt> [args...]`** – Display a note. If additional arguments are supplied,
  the prompt is processed as a Go `text/template` with variables `{{.Arg1}}`, `{{.Arg2}}`, ….
- **`slash version`** – Print the program version.

## Installation

The recommended way to install `slash` is via `go install`:

```bash
go install github.com/sirasagi62/slash@latest
```

The command compiles the program and places the binary in your `$GOBIN` directory
(usually `$HOME/go/bin`). Make sure that directory is on your `PATH`:

```bash
export PATH=$PATH:$HOME/go/bin   # add to your shell profile (~/.bashrc, ~/.zshrc, …)
```

## Usage

```bash
# Show help
slash help

# Show version
slash version

# Create a new prompt called "todo"
slash new todo

# Display the prompt
slash todo

# Use templating – the prompt can contain {{.Arg1}}, {{.Arg2}}, …
slash todo "Buy milk" "Call mom"
```

## Requirements

- Go 1.20 or newer (the project uses Go 1.25.1)
- A POSIX‑compatible shell (Linux/macOS) or PowerShell on Windows

## Contributing

See `plan.md` for the project roadmap and design details. Feel free to open
issues or submit pull requests.

## License

This project is licensed under the terms of the `LICENSE` file in the repository.

---

*Generated on $(date).*
