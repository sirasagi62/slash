package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

func printHelp() {
	fmt.Println(`Usage: slash <note> [args...]
Commands:
  new <name>    Create a new note
  help          Show this help message

If <note> is a note name, its content is printed. If additional arguments are provided,
the note is treated as a Go text/template with Arg1, Arg2, ... variables.`)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: slash <note> [args...]")
		os.Exit(1)
	}

	note := os.Args[1]
	args := os.Args[2:]

	// Help command
	if note == "help" || note == "-h" || note == "--help" {
		printHelp()
		os.Exit(0)
	}

	// Handle creation of a new note
	if note == "new" {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "Usage: slash new <name>")
			os.Exit(1)
		}
		name := args[0]

		homeDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
			os.Exit(1)
		}
		path := filepath.Join(homeDir, ".config", "slash", name+".md")

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

		// Open the file in the user's editor
		editor := os.Getenv("EDITOR")
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
		// Exit after creating and editing the note
		os.Exit(0)
	}

	// Determine possible note locations
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	paths := []string{
		filepath.Join(homeDir, ".config", "slash", note+".md"),
		filepath.Join(".", ".slash", note+".md"),
	}

	var content []byte
	found := false
	for _, p := range paths {
		if data, err := ioutil.ReadFile(p); err == nil {
			content = data
			found = true
			break
		}
	}
	if !found {
		fmt.Fprintf(os.Stderr, "Note not found: %s\n", note)
		os.Exit(1)
	}

	// If there are extra arguments, treat the note as a Go text/template
	if len(args) > 0 {
		tmpl, err := template.New("note").Parse(string(content))
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
