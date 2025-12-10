package handlers_test

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/titpetric/yamlexpr"
)

// Note: Full matrix handler tests require integration with Expr.processMapWithContext
// which is a private method. These tests verify the pure logic functions indirectly
// through YAML fixture tests in expr_fixtures_test.go

// Example fixture test pattern (in expr_fixtures_test.go):
//
// testdata/fixtures/210-matrix-simple.yaml:
//   Input:
//     jobs:
//       build:
//         matrix:
//           os: [linux, windows]
//           version: [12, 14]
//         name: "Build on $os v$version"
//   Expected:
//     jobs:
//       build:
//         - name: "Build on linux v12"
//           os: linux
//           version: 12
//         - name: "Build on linux v14"
//           os: linux
//           version: 14
//         - name: "Build on windows v12"
//           os: windows
//           version: 12
//         - name: "Build on windows v14"
//           os: windows
//           version: 14

func TestMatrixHandler_Documentation(t *testing.T) {
	fixtureTest(t, os.DirFS("../testdata/fixtures"), "200-matrix-simple.yaml")
}

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
