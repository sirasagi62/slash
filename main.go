package main

import (
	"fmt"
	"os"
)

const version = "0.1.0"

func printHelp() {
	fmt.Println(`Usage: slash <prompt> [args...]
Commands:
  edit <name>    Create a new prompt
  help          Show this help message
  version       Print the program version

If <prompt> is a prompt name, its content is printed. If additional arguments are provided,
the prompt is treated as a Go text/template with Arg1, Arg2, ... variables.`)
}

func main() {
	// Parse command‑line arguments
	cmd, args, err := parseArgs()
	if err != nil {
		fatal(err, "invalid arguments")
	}

	switch cmd {
	case "help", "-h", "--help":
		runHelp()
	case "version":
		runVersion()
	case "edit":
		if len(args) == 0 {
			fatal(fmt.Errorf("missing prompt name"), "Usage: slash edit <name>")
		}
		if err := runEdit(args[0], args[1:]); err != nil {
			fatal(err, "edit failed")
		}
	default:
		// Treat cmd as a prompt name
		content, err := loadPrompt(cmd)
		if err != nil {
			fatal(err, fmt.Sprintf("Note not found: %s", cmd))
		}
		if len(args) > 0 {
			if err := renderTemplate(content, args, os.Stdout); err != nil {
				fatal(err, "template execution error")
			}
		} else {
			fmt.Print(string(content))
		}
	}
}
