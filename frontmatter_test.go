package yamlexpr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
)

var ParseDocument = yamlexpr.ParseDocument

// TestParseDocument_SimpleContent tests parsing content without frontmatter.
func TestParseDocument_SimpleContent(t *testing.T) {
	content := "key: value"
	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Empty(t, doc.Frontmatter)
	require.Len(t, doc.Sections, 1)
	require.Contains(t, doc.Sections[0], "key: value")
}

// TestParseDocument_WithFrontmatter tests parsing content with frontmatter.
func TestParseDocument_WithFrontmatter(t *testing.T) {
	content := `---
title: "Test Document"
description: "A test"
---
key: value`

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Equal(t, "Test Document", doc.Frontmatter["title"])
	require.Equal(t, "A test", doc.Frontmatter["description"])
	require.Len(t, doc.Sections, 1)
	require.Contains(t, doc.Sections[0], "key: value")
}

// TestParseDocument_MultipleSections tests parsing with multiple sections.
func TestParseDocument_MultipleSections(t *testing.T) {
	content := `---
title: "Test"
---
input:
  - item1
  - item2
---
output:
  - item1
  - item2`

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Equal(t, "Test", doc.Frontmatter["title"])
	require.Len(t, doc.Sections, 2)
	require.Contains(t, doc.Sections[0], "input:")
	require.Contains(t, doc.Sections[1], "output:")
}

// TestParseDocument_LeadingMarkerWithFrontmatter tests with leading --- before frontmatter.
func TestParseDocument_LeadingMarkerWithFrontmatter(t *testing.T) {
	content := `---
title: "Example"
category: "test"
tags: ["a", "b"]
---
servers:
  - api
  - worker
---
servers:
  - api
  - worker`

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Equal(t, "Example", doc.Frontmatter["title"])
	require.Equal(t, "test", doc.Frontmatter["category"])
	require.NotNil(t, doc.Frontmatter["tags"])
	require.Len(t, doc.Sections, 2)
}

// TestParseDocument_EmptyFrontmatter tests with empty frontmatter section.
func TestParseDocument_EmptyFrontmatter(t *testing.T) {
	content := `---
---
key: value`

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Empty(t, doc.Frontmatter)
	require.Len(t, doc.Sections, 1)
	require.Contains(t, doc.Sections[0], "key: value")
}

// TestParseDocument_NoFrontmatterNoMarkers tests plain YAML.
func TestParseDocument_NoFrontmatterNoMarkers(t *testing.T) {
	content := `key: value
nested:
  field: data`

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Empty(t, doc.Frontmatter)
	require.Len(t, doc.Sections, 1)
	require.Contains(t, doc.Sections[0], "key: value")
}

// TestParseDocument_InvalidFrontmatterYAML tests invalid frontmatter YAML.
func TestParseDocument_InvalidFrontmatterYAML(t *testing.T) {
	content := `---
title: [unclosed array
---
key: value`

	_, err := ParseDocument(content)

	require.Error(t, err)
	require.Contains(t, err.Error(), "parsing frontmatter")
}

// TestParseDocument_GetFrontmatterField tests retrieving individual fields.
func TestParseDocument_GetFrontmatterField(t *testing.T) {
	content := `---
title: "My Document"
category: "example"
---
content`

	doc, err := ParseDocument(content)
	require.NoError(t, err)

	title, ok := doc.GetFrontmatterField("title")
	require.True(t, ok)
	require.Equal(t, "My Document", title)

	category, ok := doc.GetFrontmatterField("category")
	require.True(t, ok)
	require.Equal(t, "example", category)

	missing, ok := doc.GetFrontmatterField("missing")
	require.False(t, ok)
	require.Empty(t, missing)
}

// TestParseDocument_GetFrontmatterFieldWithDefault tests field retrieval with defaults.
func TestParseDocument_GetFrontmatterFieldWithDefault(t *testing.T) {
	content := `---
title: "Test"
---
content`

	doc, err := ParseDocument(content)
	require.NoError(t, err)

	title := doc.GetFrontmatterFieldWithDefault("title", "default")
	require.Equal(t, "Test", title)

	missing := doc.GetFrontmatterFieldWithDefault("missing", "default")
	require.Equal(t, "default", missing)
}

// TestParseDocument_ComplexFixtureFormat tests realistic fixture format.
func TestParseDocument_ComplexFixtureFormat(t *testing.T) {
	content := `---
title: "Simple Value Iteration"
description: "Iterate over a list of values with a single loop variable"
category: "basics"
tags: ["for", "iteration", "simple"]
---
servers:
  - for: server in server_list
    name: "${server}"
server_list:
  - "api"
  - "worker"
  - "cache"
---
servers:
  - name: "api"
  - name: "worker"
  - name: "cache"
server_list:
  - "api"
  - "worker"
  - "cache"`

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Equal(t, "Simple Value Iteration", doc.Frontmatter["title"])
	require.Equal(t, "Iterate over a list of values with a single loop variable", doc.Frontmatter["description"])
	require.Equal(t, "basics", doc.Frontmatter["category"])
	require.Len(t, doc.Sections, 2)
	require.Contains(t, doc.Sections[0], "for: server in server_list")
	require.Contains(t, doc.Sections[1], "- name: \"api\"")
}

// TestParseDocument_PreservesIndentation tests that indentation is preserved.
func TestParseDocument_PreservesIndentation(t *testing.T) {
	content := `---
title: "Test"
---
items:
  - name: first
    value: 1
  - name: second
    value: 2`

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Len(t, doc.Sections, 1)
	// Check that the section still has proper indentation
	require.Contains(t, doc.Sections[0], "  - name: first")
}

// TestParseDocument_MultilineFrontmatterField tests multiline fields in frontmatter.
func TestParseDocument_MultilineFrontmatterField(t *testing.T) {
	content := `---
title: "Test"
description: |
  This is a
  multiline description
---
content: value`

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Equal(t, "Test", doc.Frontmatter["title"])
	require.NotEmpty(t, doc.Frontmatter["description"])
	require.Len(t, doc.Sections, 1)
}

// TestParseDocument_TagsArray tests array fields in frontmatter.
func TestParseDocument_TagsArray(t *testing.T) {
	content := `---
title: "Test"
tags: ["tag1", "tag2", "tag3"]
---
content`

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Equal(t, "Test", doc.Frontmatter["title"])
	tags, ok := doc.Frontmatter["tags"].([]interface{})
	require.True(t, ok)
	require.Len(t, tags, 3)
}

// TestParseDocument_EmptyDocument tests parsing empty content.
func TestParseDocument_EmptyDocument(t *testing.T) {
	content := ""

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Empty(t, doc.Frontmatter)
	require.Empty(t, doc.Sections)
}

// TestParseDocument_OnlyWhitespace tests whitespace-only content.
func TestParseDocument_OnlyWhitespace(t *testing.T) {
	content := "   \n\n   \n"

	doc, err := ParseDocument(content)

	require.NoError(t, err)
	require.Empty(t, doc.Frontmatter)
	require.Empty(t, doc.Sections)
}
