package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
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
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: slash <prompt> [args...]")
		os.Exit(1)
	}

	prompt := os.Args[1]
	args := os.Args[2:]

	// Help command
	if prompt == "help" || prompt == "-h" || prompt == "--help" {
		printHelp()
		os.Exit(0)
	}

	// Version command
	if prompt == "version" {
		fmt.Println(version)
		os.Exit(0)
	}

	// Handle editing of a prompt
	if prompt == "edit" {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: slash edit <name>")
			os.Exit(1)
		}
		name := args[0]

		// Prevent creating prompts with reserved names
		if name == "help" || name == "version" || name == "new" {
			fmt.Fprintf(os.Stderr, "Error: cannot create prompt with reserved name \"%s\"\n", name)
			os.Exit(1)
		}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		path := filepath.Join(homeDir, ".config", "slash", name+".md")

		// If -p option is provided, replace the prompt content without opening editor
		if len(args) >= 2 && args[1] == "-p" {
			content := strings.Join(args[2:], " ")
			// Ensure the target directory exists
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
				os.Exit(1)
			}
			// Write the provided content to the file (overwrite or create)
			if err := os.WriteFile(path, []byte(content), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Prompt replaced.")
			os.Exit(0)
		}

		// Ensure the target directory exists
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating directory: %v\n", err)
			os.Exit(1)
		}

		// Create the file if it does not exist
		f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file: %v\n", err)
			os.Exit(1)
		}
		f.Close()

		// Determine the editor to use:
		// 1. Try reading from ~/.config/slash/config.json ("editor" field)
		// 2. Fallback to $EDITOR environment variable
		// 3. Default to "vi"
		editor := ""
		configPath := filepath.Join(homeDir, ".config", "slash", "config.json")
		if cfgData, err := os.ReadFile(configPath); err == nil {
			var cfg struct {
				Editor string `json:"editor"`
			}
			if jsonErr := json.Unmarshal(cfgData, &cfg); jsonErr == nil && cfg.Editor != "" {
				editor = cfg.Editor
			}
		}
		if editor == "" {
			editor = os.Getenv("EDITOR")
		}
		if editor == "" {
			editor = "vi"
		}
		cmd := exec.Command(editor, path)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error running editor: %v\n", err)
			os.Exit(1)
		}
		// Exit after creating and editing the prompt
		os.Exit(0)
	}

	// Determine possible prompt locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	paths := []string{
		filepath.Join(homeDir, ".config", "slash", prompt+".md"),
		filepath.Join(".", ".slash", prompt+".md"),
	}

	var content []byte
	found := false
	for _, p := range paths {
		if data, err := os.ReadFile(p); err == nil {
			content = data
			found = true
			break
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "Note not found: %s\n", prompt)
		os.Exit(1)
	}

	// If there are extra arguments, treat the prompt as a Go text/template
	if len(args) > 0 {
		tmpl, err := template.New("prompt").Parse(string(content))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Template parse error: %v\n", err)
			os.Exit(1)
		}
		// Build a simple map of Arg1, Arg2, ... for template execution
		data := map[string]string{}
		for i, v := range args {
			key := fmt.Sprintf("Arg%d", i+1)
			data[key] = v
		}
		if err := tmpl.Execute(os.Stdout, data); err != nil {
			fmt.Fprintf(os.Stderr, "Template execution error: %v\n", err)
			os.Exit(1)
		}
	} else {
		// No template processing; just print the raw content
		fmt.Print(string(content))
	}
}
