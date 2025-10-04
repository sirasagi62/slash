# slash

`slash` is a simple Go command‑line utility for managing markdown prompts.
It lets you create, view, and template prompts directly from the terminal.

## Features

- **`slash help`** – Show usage information.
- **`slash edit <name>`** – Create or edit a markdown prompt in `~/.config/slash` (or `./.slash`).  
  Use `slash edit <name> -p "<content>"` to replace the prompt directly from the CLI.
- **`slash <prompt> [args...]`** – Display a prompt. If extra arguments are supplied,
  the prompt is treated as a Go `text/template` with variables `{{.Arg1}}`, `{{.Arg2}}`, ….
- **`slash version`** – Print the program version.

## Installation

The recommended way to install `slash` is via `go install`:
> **Note:** The project requires **Go 1.25.1** (or any newer 1.x release).

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
slash edit todo

# (or create/replace it in one shot)
slash edit todo -p "# TODO List\n- {{.Arg1}}\n- {{.Arg2}}"

# Display the prompt (raw markdown)
slash todo

# Use templating – the prompt can contain {{.Arg1}}, {{.Arg2}}, …
slash todo "Buy milk" "Call mom"
```

## Configuration

`slash` looks for a JSON config file at:

```
$HOME/.config/slash/config.json
```

```json
{
  "editor": "code --wait"
}
```

If the file is missing, `slash` falls back to the `$EDITOR` environment variable,
and finally to `vi`. The config file is optional; you only need it to override the default editor.

## Requirements

- Go **1.25.1** or newer
- A POSIX‑compatible shell (Linux/macOS) or PowerShell on Windows

## Contributing

See [`agents/plan.md`](agents/plan.md) for the project roadmap and design details. Feel free to open
issues or submit pull requests.

## License

This project is licensed under the terms of the `LICENSE` file in the repository.
