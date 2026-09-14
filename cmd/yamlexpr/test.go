package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"

	"github.com/titpetric/yamlexpr"
)

// TestCommand runs the fixture files and reports each as a pass or a fail.
type TestCommand struct{}

// NewTestCommand constructs the test command.
func NewTestCommand() *TestCommand {
	return &TestCommand{}
}

// Help returns the usage text for the test command.
func (c *TestCommand) Help() string {
	return `yamlexpr test [-dir DIR]

Runs every fixture below the directory, testdata/fixtures-by-feature by
default. A fixture is frontmatter, an input document, and the documents it
should produce, separated by ---; a file opening on an underscore is data
for other fixtures and is skipped. Includes are resolved relative to the
fixture's own directory. Exits 1 when any fixture fails.`
}

// Run walks the fixtures and evaluates each against what it expects.
func (c *TestCommand) Run(args []string) int {
	fset := flag.NewFlagSet("test", flag.ContinueOnError)
	dir := fset.String("dir", "testdata/fixtures-by-feature", "fixture directory")
	if err := fset.Parse(args); err != nil {
		return 1
	}

	var passed, failed, skipped int

	err := filepath.WalkDir(*dir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Ext(path) != ".yaml" || strings.HasPrefix(entry.Name(), "_") {
			return nil
		}

		name, err := filepath.Rel(*dir, path)
		if err != nil {
			name = path
		}

		switch reason, err := runFixture(path); {
		case reason != "":
			fmt.Printf("SKIP %s: %s\n", name, reason)
			skipped++
		case err != nil:
			fmt.Printf("FAIL %s: %v\n", name, err)
			failed++
		default:
			fmt.Printf("PASS %s\n", name)
			passed++
		}
		return nil
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	fmt.Printf("\n%d passed, %d failed, %d skipped\n", passed, failed, skipped)
	if failed > 0 {
		return 1
	}
	return 0
}

// runFixture evaluates one fixture's input and compares the documents it
// produced against the ones the fixture expects. A fixture whose frontmatter
// carries skip describes behaviour that is not implemented yet: it stays in
// the generated docs, and the reason comes back instead of a verdict.
func runFixture(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	fixture, err := yamlexpr.ParseDocument(string(data))
	if err != nil {
		return "", err
	}
	if reason, ok := fixture.GetFrontmatterField("skip"); ok {
		return reason, nil
	}
	if len(fixture.Sections) < 2 {
		return "", fmt.Errorf("the fixture holds no expected output")
	}

	// Decoded as a plain map: yaml fills the nested mappings of a named map
	// type with that same type, and the processor matches map[string]any.
	var input map[string]any
	if err := yaml.Unmarshal([]byte(fixture.Sections[0]), &input); err != nil {
		return "", fmt.Errorf("error parsing input: %w", err)
	}

	docs, err := yamlexpr.New(os.DirFS(filepath.Dir(path))).Parse(yamlexpr.Document(input))
	if err != nil {
		return "", err
	}

	expected := fixture.Sections[1:]
	if len(docs) != len(expected) {
		return "", fmt.Errorf("produced %d documents, want %d", len(docs), len(expected))
	}

	for i, section := range expected {
		var want map[string]any
		if err := yaml.Unmarshal([]byte(section), &want); err != nil {
			return "", fmt.Errorf("error parsing expected document %d: %w", i+1, err)
		}
		if !yamlexpr.ValuesEqual(want, map[string]any(docs[i])) {
			got, _ := yaml.Marshal(docs[i])
			return "", fmt.Errorf("document %d does not match:\nwant:\n%s\ngot:\n%s", i+1, section, got)
		}
	}

	return "", nil
}
