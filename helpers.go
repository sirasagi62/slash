package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

type Config struct {
	Editor string `json:"editor"`
}

var ErrPromptNotFound = errors.New("prompt not found")

// fatal prints an error message to stderr and exits with status 1.
func fatal(err error, msg string) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", msg, err)
	os.Exit(1)
}

// parseArgs extracts the command (prompt) and its arguments.
func parseArgs() (string, []string, error) {
	if len(os.Args) < 2 {
		return "", nil, fmt.Errorf("Usage: slash <prompt> [args...]")
	}
	return os.Args[1], os.Args[2:], nil
}

// runHelp prints the help message.
func runHelp() {
	printHelp()
}

// runVersion prints the program version.
func runVersion() {
	fmt.Println(version)
}

// runEdit handles the "edit" sub‑command.
func runEdit(name string, extra []string) error {
	// Guard reserved names
	if name == "help" || name == "version" || name == "edit" {
		return fmt.Errorf("cannot create prompt with reserved name %q", name)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("getting home directory: %w", err)
	}
	path := filepath.Join(homeDir, ".config", "slash", name+".md")

	// If "-p" flag is present, write the content directly.
	if len(extra) >= 2 && extra[0] == "-p" {
		content := strings.Join(extra[1:], " ")
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Println("Prompt replaced.")
		return nil
	}

	// Ensure directory exists.
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	// Create empty file if missing.
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	f.Close()

	// Determine editor: config → $EDITOR → vi
	editor := ""
	cfgPath := filepath.Join(homeDir, ".config", "slash", "config.json")
	if cfgData, err := os.ReadFile(cfgPath); err == nil {
		var cfg Config
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
		return fmt.Errorf("running editor %q: %w", editor, err)
	}
	return nil
}

// loadPrompt reads the prompt file from the possible locations.
func loadPrompt(name string) ([]byte, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("getting home directory: %w", err)
	}
	paths := []string{
		filepath.Join(homeDir, ".config", "slash", name+".md"),
		filepath.Join(".", ".slash", name+".md"),
	}
	for _, p := range paths {
		if data, err := os.ReadFile(p); err == nil {
			return data, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrPromptNotFound, name)
}

// renderTemplate parses and executes a Go text/template with ArgN variables.
func renderTemplate(src []byte, args []string, out io.Writer) error {
	tmpl, err := template.New("prompt").Parse(string(src))
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}
	data := map[string]string{}
	for i, v := range args {
		data[fmt.Sprintf("Arg%d", i+1)] = v
	}
	if err := tmpl.Execute(out, data); err != nil {
		return fmt.Errorf("template execution error: %w", err)
	}
	return nil
}
