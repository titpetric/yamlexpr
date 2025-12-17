package testdata

import (
	"bufio"
	"os"
	"strings"
	"testing"

	yaml "gopkg.in/yaml.v3"
)

// FixtureCase represents a single test case from a fixture file
type FixtureCase struct {
	// Metadata from frontmatter (optional)
	Title       string
	Description string
	Category    string
	Tags        []string

	// Test data
	Input  interface{}
	Output interface{}

	// Source file info
	SourceFile string
	LineNumber int
}

// FixtureLoader loads fixtures from YAML files with optional frontmatter
type FixtureLoader struct {
	file *os.File
}

// NewFixtureLoader creates a new fixture loader for a file
func NewFixtureLoader(path string) (*FixtureLoader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return &FixtureLoader{file: file}, nil
}

// Close closes the underlying file
func (fl *FixtureLoader) Close() error {
	if fl.file != nil {
		return fl.file.Close()
	}
	return nil
}

// Load parses the fixture and returns test cases
func (fl *FixtureLoader) Load(sourceFile string) ([]FixtureCase, error) {
	defer fl.file.Close()

	scanner := bufio.NewScanner(fl.file)
	var cases []FixtureCase

	// Check for frontmatter
	var frontmatter map[string]interface{}
	var hasFrontmatter bool

	if scanner.Scan() && strings.TrimSpace(scanner.Text()) == "---" {
		hasFrontmatter = true
		// Read frontmatter
		var fmLines []string
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) == "---" {
				break
			}
			fmLines = append(fmLines, line)
		}

		if len(fmLines) > 0 {
			fmYAML := strings.Join(fmLines, "\n")
			yaml.Unmarshal([]byte(fmYAML), &frontmatter)
		}
	}

	// Read input section
	var inputLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "---" {
			break
		}
		inputLines = append(inputLines, line)
	}

	// Read output section
	var outputLines []string
	for scanner.Scan() {
		line := scanner.Text()
		outputLines = append(outputLines, line)
	}

	// Parse input and output YAML
	var input, output interface{}

	if len(inputLines) > 0 {
		inputYAML := strings.TrimSpace(strings.Join(inputLines, "\n"))
		yaml.Unmarshal([]byte(inputYAML), &input)
	}

	if len(outputLines) > 0 {
		outputYAML := strings.TrimSpace(strings.Join(outputLines, "\n"))
		yaml.Unmarshal([]byte(outputYAML), &output)
	}

	// Build fixture case
	fc := FixtureCase{
		Input:      input,
		Output:     output,
		SourceFile: sourceFile,
	}

	// Extract metadata from frontmatter
	if frontmatter != nil {
		if title, ok := frontmatter["title"].(string); ok {
			fc.Title = title
		}
		if desc, ok := frontmatter["description"].(string); ok {
			fc.Description = desc
		}
		if cat, ok := frontmatter["category"].(string); ok {
			fc.Category = cat
		}
		if tagsRaw, ok := frontmatter["tags"].([]interface{}); ok {
			for _, tag := range tagsRaw {
				if tagStr, ok := tag.(string); ok {
					fc.Tags = append(fc.Tags, tagStr)
				}
			}
		}
	}

	cases = append(cases, fc)

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cases, nil
}

// LoadAllFixtures loads all fixture files from a directory
func LoadAllFixtures(t *testing.T, patterns ...string) []FixtureCase {
	var allCases []FixtureCase

	// Default patterns
	if len(patterns) == 0 {
		patterns = []string{
			"testdata/fixtures/*.yaml",
			"testdata/fixtures-by-feature/**/*.yaml",
		}
	}

	// TODO: Implement directory walking
	// For now, this is a placeholder for the implementation

	return allCases
}
