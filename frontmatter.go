package yamlexpr

import (
	"fmt"
	"strings"

	yaml "gopkg.in/yaml.v3"
)

// DocumentContent represents a YAML document with optional frontmatter and sections.
type DocumentContent struct {
	// Frontmatter is the first yaml code block at the start of the file.
	Frontmatter map[string]any
	// Sections are the remaining yaml blocks in the file.
	// With fixtures, the first section is the yamlexpr "template",
	// while the remaining sections are document outputs from Parse().
	Sections []string
}

// TrimLeadingDocumentMarker trims the leading --- from a YAML document if present.
func TrimLeadingDocumentMarker(content string) string {
	content = strings.TrimSpace(content)
	if strings.HasPrefix(content, "---") {
		// Remove the --- and the following newline if present
		content = strings.TrimPrefix(content, "---")
		content = strings.TrimLeft(content, "\n\r")
		content = strings.TrimSpace(content)
	}
	return content
}

// ParseDocument parses a YAML document into frontmatter and sections separated by ---.
// The document format is:
//
//	---
//	title: "Example"
//	description: "..."
//	---
//	first section (usually input)
//	---
//	second section (usually output)
//	---
//	additional sections...
//
// If the document starts with ---, it is trimmed before parsing.
func ParseDocument(content string) (*DocumentContent, error) {
	doc := &DocumentContent{
		Frontmatter: make(map[string]any),
		Sections:    make([]string, 0),
	}

	// Trim leading --- from document if present
	content = TrimLeadingDocumentMarker(content)
	content = strings.TrimSpace(content)

	lines := strings.Split(content, "\n")
	separatorCount := 0
	var currentSection []string

	for _, line := range lines {
		if strings.TrimSpace(line) == "---" {
			separatorCount++

			// Save the current section
			sectionContent := strings.TrimSpace(strings.Join(currentSection, "\n"))
			if separatorCount == 1 && sectionContent != "" {
				// First section is frontmatter
				if err := yaml.Unmarshal([]byte(sectionContent), &doc.Frontmatter); err != nil {
					return nil, fmt.Errorf("error parsing frontmatter: %w", err)
				}
			} else if sectionContent != "" {
				// Other sections are added as-is
				doc.Sections = append(doc.Sections, sectionContent)
			}

			currentSection = []string{}
			continue
		}

		currentSection = append(currentSection, line)
	}

	// Handle last section if no trailing separator
	if len(currentSection) > 0 {
		sectionContent := strings.TrimSpace(strings.Join(currentSection, "\n"))
		if sectionContent != "" {
			doc.Sections = append(doc.Sections, sectionContent)
		}
	}

	return doc, nil
}

// GetFrontmatterField extracts a specific field from frontmatter.
func (d *DocumentContent) GetFrontmatterField(field string) (string, bool) {
	val, ok := d.Frontmatter[field].(string)
	return val, ok
}

// GetFrontmatterFieldWithDefault extracts a field or returns a default.
func (d *DocumentContent) GetFrontmatterFieldWithDefault(field, defaultVal string) string {
	val, ok := d.GetFrontmatterField(field)
	if !ok {
		return defaultVal
	}
	return val
}
