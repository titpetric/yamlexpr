package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"

	"github.com/titpetric/yamlexpr"
)

// ProcessCommand evaluates YAML files and writes the documents they produce.
type ProcessCommand struct{}

// NewProcessCommand constructs the process command.
func NewProcessCommand() *ProcessCommand {
	return &ProcessCommand{}
}

// Help returns the usage text for the process command.
func (c *ProcessCommand) Help() string {
	return `yamlexpr process [file...]

Evaluates each YAML file and prints the documents it produces, separated by
---. Includes are resolved relative to the file's own directory. With no
file, stdin is read and evaluated with includes resolved relative to the
working directory.`
}

// Run evaluates each file in turn, or stdin when none is named.
func (c *ProcessCommand) Run(args []string) int {
	if len(args) == 0 {
		if err := c.processStdin(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
		return 0
	}

	for _, filename := range args {
		expr := yamlexpr.New(os.DirFS(filepath.Dir(filename)))

		docs, err := expr.Load(filepath.Base(filename))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}

		if err := printDocuments(docs); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			return 1
		}
	}

	return 0
}

// processStdin evaluates one document read from stdin.
func (c *ProcessCommand) processStdin() error {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}

	// Decoded as a plain map: yaml fills the nested mappings of a named map
	// type with that same type, and the processor matches map[string]any.
	var doc map[string]any
	if err := yaml.Unmarshal(input, &doc); err != nil {
		return fmt.Errorf("error parsing input: %w", err)
	}

	docs, err := yamlexpr.New(os.DirFS(".")).Parse(yamlexpr.Document(doc))
	if err != nil {
		return err
	}

	return printDocuments(docs)
}

// printDocuments prints the documents as YAML, separated the way a multi
// document file writes them.
func printDocuments(docs []yamlexpr.Document) error {
	for i, doc := range docs {
		if i > 0 {
			fmt.Println("---")
		}
		out, err := yaml.Marshal(doc)
		if err != nil {
			return err
		}
		fmt.Print(string(out))
	}
	return nil
}
