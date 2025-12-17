package yamlexpr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
)

func TestNew(t *testing.T) {
	e := yamlexpr.New(nil)
	require.NotNil(t, e)
}

func TestExpr_Parse(t *testing.T) {
	e := yamlexpr.New(nil)

	doc := yamlexpr.Document{
		"name":  "${user.name}",
		"items": []any{"a", "b"},
		"user": map[string]any{
			"name": "John",
		},
	}

	want := yamlexpr.Document{
		"name":  "John",
		"items": []any{"a", "b"},
		"user": map[string]any{
			"name": "John",
		},
	}

	docs, err := e.Parse(doc)

	require.NoError(t, err)
	require.Len(t, docs, 1)
	require.Equal(t, want, docs[0])
}
