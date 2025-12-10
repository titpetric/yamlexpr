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

//go:embed testdata/fixtures
var fixtureFS embed.FS

// fixtureTest executes the standard fixture assertions against all supplied
// filenames using the provided filesystem.
func fixtureTest(tb testing.TB, fixtureFS fs.FS, filenames ...string) {
	tb.Helper()

	expr := yamlexpr.New(yamlexpr.WithFS(fixtureFS))

	for _, fn := range filenames {
		result, err := expr.Load(fn)
		assert.NoError(tb, err)
		assert.NotNil(tb, result)
		assert.IsType(tb, map[string]any{}, result)
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

		if filepath.Ext(name) != ".yaml" {
			continue
		}

		files = append(files, name)
	}

	assert.NotEmpty(t, files, "no .yaml fixtures found")

	fixtureTest(t, subFS, files...)
}
