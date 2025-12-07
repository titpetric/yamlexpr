package yamlexpr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
	"github.com/titpetric/yamlexpr/stack"
)

func TestInterpolateString(t *testing.T) {
	st := stack.New(map[string]any{
		"name": "John",
		"age":  30,
		"user": map[string]any{
			"email": "john@example.com",
		},
	})

	// Simple interpolation
	result := yamlexpr.InterpolateString("Hello ${name}", st)
	require.Equal(t, "Hello John", result)

	// Multiple interpolations
	result = yamlexpr.InterpolateString("${name} is ${age} years old", st)
	require.Equal(t, "John is 30 years old", result)

	// Nested path
	result = yamlexpr.InterpolateString("Email: ${user.email}", st)
	require.Equal(t, "Email: john@example.com", result)

	// Missing variable
	result = yamlexpr.InterpolateString("Hello ${missing}", st)
	require.Equal(t, "Hello ${missing}", result)

	// No interpolation
	result = yamlexpr.InterpolateString("No variables here", st)
	require.Equal(t, "No variables here", result)
}

func TestContainsInterpolation(t *testing.T) {
	require.True(t, yamlexpr.ContainsInterpolation("Hello ${name}"))
	require.False(t, yamlexpr.ContainsInterpolation("Hello name"))
	require.False(t, yamlexpr.ContainsInterpolation("${unclosed"))
	require.False(t, yamlexpr.ContainsInterpolation("${"))
}
