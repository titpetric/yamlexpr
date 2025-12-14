package yamlexpr_test

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/titpetric/yamlexpr"
)

//go:embed all:testdata/fixtures
var fixtureFS embed.FS

// fixtureTest executes the standard fixture assertions against all supplied
// filenames using the provided filesystem.
func fixtureTest(tb testing.TB, fixtureFS fs.FS, filenames ...string) {
	tb.Helper()

	expr := yamlexpr.New(yamlexpr.WithFS(fixtureFS))

	for _, fn := range filenames {
		docs, err := expr.Load(fn)
		assert.NoError(tb, err)
		assert.NotNil(tb, docs)
		assert.NotEmpty(tb, docs)
		// For regular (non-expansion) fixtures, should return single document
		if len(docs) == 1 {
			_, ok := docs[0].(map[string]any)
			assert.True(tb, ok)
		}
	}
}

// TestExpr_Load scans the embedded fixtures directory for all *.yml files,
// skipping any whose base name begins with '_', and runs standard fixture
// validation against each of them.
func TestExpr_Load(t *testing.T) {
	subFS, err := fs.Sub(fixtureFS, "testdata/fixtures")
	assert.NoError(t, err)

	entries, err := fs.ReadDir(subFS, ".")
	assert.NoError(t, err)

	var files []string

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		if strings.HasPrefix(name, "_") {
			continue
		}

		// Skip root-level for/matrix fixtures (150-199: root-level for, 200+: root-level matrix)
		// These are tested in TestExpr_LoadMulti
		if strings.HasPrefix(name, "1") && name >= "150" || strings.HasPrefix(name, "2") {
			continue
		}

		if filepath.Ext(name) != ".yaml" {
			continue
		}

		files = append(files, name)
	}

	assert.NotEmpty(t, files, "no .yaml fixtures found")

	fixtureTest(t, subFS, files...)
}

// TestExpr_Load_RootLevelExpansion tests root-level for directive that expands into multiple documents.
// Tests fixtures 150-199 (root-level for directive).
func TestExpr_Load_RootLevelExpansion(t *testing.T) {
	subFS, err := fs.Sub(fixtureFS, "testdata/fixtures")
	assert.NoError(t, err)

	entries, err := fs.ReadDir(subFS, ".")
	assert.NoError(t, err)

	expr := yamlexpr.New(yamlexpr.WithFS(subFS))

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		if strings.HasPrefix(name, "_") {
			continue
		}

		if filepath.Ext(name) != ".yaml" {
			continue
		}

		// Only process root-level for fixtures (150-199)
		// Note: matrix fixtures (200+) have list-level matrix, not root-level
		if !(strings.HasPrefix(name, "1") && name >= "150" && name < "200") {
			continue
		}

		t.Run(name, func(t *testing.T) {
			docs, err := expr.Load(name)
			assert.NoError(t, err)
			assert.NotNil(t, docs)
			assert.NotEmpty(t, docs, "Load should return at least one document")

			// Verify all returned items are maps
			for i, doc := range docs {
				assert.NotNil(t, doc, "document %d should not be nil", i)
				_, ok := doc.(map[string]any)
				assert.True(t, ok, "document %d should be a map, got %T", i, doc)
			}
		})
	}
}
