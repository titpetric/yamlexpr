package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/titpetric/yamlexpr"
)

// GenCommand renders a feature's fixtures as its documentation page: each
// fixture is a section with the input and the documents it produces, so the
// docs and the tests are one set of files.
type GenCommand struct{}

// NewGenCommand constructs the gen command.
func NewGenCommand() *GenCommand {
	return &GenCommand{}
}

// Help returns the usage text for the gen command.
func (c *GenCommand) Help() string {
	return `yamlexpr gen -feature NAME [-dir DIR]

Renders the fixtures of one feature as markdown on stdout: the title and
description from each fixture's frontmatter, the input, and the documents it
produces. The fixtures live under testdata/fixtures-by-feature/NAME; a
feature with none renders nothing.`
}

// Run renders the named feature.
func (c *GenCommand) Run(args []string) int {
	fset := flag.NewFlagSet("gen", flag.ContinueOnError)
	feature := fset.String("feature", "", "feature to render")
	dir := fset.String("dir", "testdata/fixtures-by-feature", "fixture directory")
	if err := fset.Parse(args); err != nil {
		return 1
	}
	if *feature == "" {
		fmt.Fprintln(os.Stderr, "error: -feature names the fixtures to render")
		return 1
	}

	entries, err := os.ReadDir(filepath.Join(*dir, *feature))
	if err != nil {
		// A feature with no fixtures renders nothing, so a page can be
		// wired before its fixtures are written.
		fmt.Fprintf(os.Stderr, "warning: no fixtures for %s: %v\n", *feature, err)
		return 0
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".yaml" || strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	for i, name := range names {
		if err := renderFixture(filepath.Join(*dir, *feature, name), i > 0); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s: %v\n", name, err)
			return 1
		}
	}

	return 0
}

// renderFixture writes one fixture as a section: the heading, the
// description, the input, and the output documents.
func renderFixture(path string, separate bool) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	fixture, err := yamlexpr.ParseDocument(string(data))
	if err != nil {
		return err
	}
	if len(fixture.Sections) < 2 {
		return fmt.Errorf("the fixture holds no expected output")
	}

	if separate {
		fmt.Println()
	}

	title := fixture.GetFrontmatterFieldWithDefault("title", strings.TrimSuffix(filepath.Base(path), ".yaml"))
	fmt.Printf("## %s\n", title)

	if description, ok := fixture.GetFrontmatterField("description"); ok {
		fmt.Printf("\n%s\n", description)
	}

	fmt.Printf("\n**Input:**\n\n```yaml\n%s\n```\n", fixture.Sections[0])
	fmt.Printf("\n**Output:**\n")

	outputs := fixture.Sections[1:]
	if len(outputs) > 1 {
		fmt.Printf("\nRendering produces **%d** documents:\n", len(outputs))
	}
	for _, section := range outputs {
		fmt.Printf("\n```yaml\n%s\n```\n", section)
	}

	return nil
}
