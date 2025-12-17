package yamlexpr_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
)

// TestExpr_Load tests the Load function.
func TestExpr_Load(t *testing.T) {
	expr := yamlexpr.New(os.DirFS("testdata/fixtures"))

	docs, err := expr.Load("001-simple-pass-through.yaml")
	require.NoError(t, err)
	require.NotNil(t, docs)
	require.Len(t, docs, 1)
	require.IsType(t, yamlexpr.Document{}, docs[0])
}
