package handlers_test

import (
	"testing"
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
	// This test documents the expected behavior
	// See testdata/fixtures/ for actual test cases
	t.Log("Matrix handler tests are in fixture tests")
	t.Log("Pattern: testdata/fixtures/2XX-matrix-*.yaml")
	t.Log("See MATRIX_SPEC_COMPARISON.md for test specifications")
}
