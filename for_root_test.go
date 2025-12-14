package yamlexpr_test

import (
	"io/fs"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
)

// TestLoad_ForRootLevel_ContentValidation validates that Load returns
// the correct number of documents with correct content when processing root-level
// for directives.
func TestLoad_ForRootLevel_ContentValidation(t *testing.T) {
	subFS, err := fs.Sub(fixtureFS, "testdata/fixtures")
	require.NoError(t, err)

	expr := yamlexpr.New(yamlexpr.WithFS(subFS))

	// Test 150-for-root-level-simple.yaml
	// Input has: for: item in items with 3 items
	// Expected: 3 documents, one for each item
	docs, err := expr.Load("150-for-root-level-simple.yaml")
	require.NoError(t, err)
	require.Equal(t, 3, len(docs), "should expand to 3 documents")

	// Verify each document has the correct name and items array
	expectedNames := []string{"alice", "bob", "charlie"}
	for i, expectedName := range expectedNames {
		doc, ok := docs[i].(map[string]any)
		require.True(t, ok, "document %d should be a map", i)
		require.Equal(t, expectedName, doc["name"], "document %d name", i)

		// Verify items array is still present (not consumed by for directive)
		items, ok := doc["items"].([]any)
		require.True(t, ok, "document %d should have items array", i)
		require.Equal(t, 3, len(items), "document %d items should have 3 items", i)
	}
}

// TestLoad_ForRootLevel_WithIndex validates that root-level for loops
// correctly bind index variables.
func TestLoad_ForRootLevel_WithIndex(t *testing.T) {
	subFS, err := fs.Sub(fixtureFS, "testdata/fixtures")
	require.NoError(t, err)

	expr := yamlexpr.New(yamlexpr.WithFS(subFS))

	// Test 151-for-root-level-with-index.yaml
	// Input has: for: (idx, item) in versions with 3 versions
	// Expected: 3 documents with correct index and version values
	docs, err := expr.Load("151-for-root-level-with-index.yaml")
	require.NoError(t, err)
	require.Equal(t, 3, len(docs), "should expand to 3 documents")

	// Test data: versions ["1.0", "1.1", "2.0"]
	testData := []struct {
		idx     int
		version string
	}{
		{0, "1.0"},
		{1, "1.1"},
		{2, "2.0"},
	}

	for i, expected := range testData {
		doc, ok := docs[i].(map[string]any)
		require.True(t, ok, "document %d should be a map", i)

		// Index is stored as int
		require.Equal(t, expected.version, doc["version"], "document %d version", i)
		require.Equal(t, expected.idx, doc["index"], "document %d index", i)

		// Verify versions array is still present
		versions, ok := doc["versions"].([]any)
		require.True(t, ok, "document %d should have versions array", i)
		require.Equal(t, 3, len(versions), "document %d versions should have 3 items", i)
	}
}
