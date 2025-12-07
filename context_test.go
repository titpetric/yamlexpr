package yamlexpr_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
	"github.com/titpetric/yamlexpr/stack"
)

func TestNewExprContext(t *testing.T) {
	st := stack.New(map[string]any{"name": "test"})
	ctx := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{
		Stack: st,
		Path:  "root",
	})
	require.NotNil(t, ctx)
	require.Equal(t, "root", ctx.Path())
	require.Equal(t, st, ctx.Stack())
}

func TestNewExprContext_Defaults(t *testing.T) {
	ctx := yamlexpr.NewExprContext(nil)
	require.NotNil(t, ctx)
	require.Equal(t, "", ctx.Path())
	require.NotNil(t, ctx.Stack())
}

func TestExprContext_WithPath(t *testing.T) {
	st := stack.New(nil)
	ctx := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{
		Stack: st,
		Path:  "root",
	})
	newCtx := ctx.WithPath("root.child")
	require.Equal(t, "root", ctx.Path())          // Original unchanged
	require.Equal(t, "root.child", newCtx.Path()) // New context has new path
	require.Equal(t, st, newCtx.Stack())          // Stack is shared
}

func TestExprContext_AppendPath_Key(t *testing.T) {
	ctx := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{Path: "root"})
	newCtx := ctx.AppendPath("child")
	require.Equal(t, "root.child", newCtx.Path())
}

func TestExprContext_AppendPath_NestedKey(t *testing.T) {
	ctx := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{Path: "root.database"})
	newCtx := ctx.AppendPath("config")
	require.Equal(t, "root.database.config", newCtx.Path())
}

func TestExprContext_AppendPath_ArrayIndex(t *testing.T) {
	ctx := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{Path: "root"})
	newCtx := ctx.AppendPath("[0]")
	require.Equal(t, "root[0]", newCtx.Path())
}

func TestExprContext_AppendPath_ArrayIndexNested(t *testing.T) {
	ctx := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{Path: "root[0]"})
	newCtx := ctx.AppendPath("items")
	require.Equal(t, "root[0].items", newCtx.Path())
}

func TestExprContext_AppendPath_EmptyPath(t *testing.T) {
	ctx := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{Path: ""})
	newCtx := ctx.AppendPath("root")
	require.Equal(t, "root", newCtx.Path())
}

func TestExprContext_AppendPath_EmptySegment(t *testing.T) {
	ctx := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{Path: "root"})
	newCtx := ctx.AppendPath("")
	require.Equal(t, "root", newCtx.Path())
}

func TestExprContext_WithInclude(t *testing.T) {
	ctx := yamlexpr.NewExprContext(nil)
	ctx1 := ctx.WithInclude("base.yaml")
	ctx2 := ctx1.WithInclude("config.yaml")

	require.Equal(t, "base.yaml", ctx1.FormatIncludeChain())
	require.Equal(t, "base.yaml -> config.yaml", ctx2.FormatIncludeChain())
	require.Equal(t, "", ctx.FormatIncludeChain()) // Original unchanged
}

func TestExprContext_FormatIncludeChain_Empty(t *testing.T) {
	ctx := yamlexpr.NewExprContext(nil)
	require.Equal(t, "", ctx.FormatIncludeChain())
}

func TestExprContext_StackScope(t *testing.T) {
	st := stack.New(map[string]any{"x": 1})
	ctx := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{Stack: st})

	// Push new scope
	ctx.PushStackScope(map[string]any{"y": 2})
	val, ok := ctx.Stack().Resolve("y")
	require.True(t, ok)
	require.Equal(t, 2, val)

	// Pop scope
	ctx.PopStackScope()
	val2, ok2 := ctx.Stack().Resolve("y")
	require.False(t, ok2)
	require.Nil(t, val2)
}

func TestExprContext_SharedStack(t *testing.T) {
	st := stack.New(map[string]any{"x": 1})
	ctx1 := yamlexpr.NewExprContext(&yamlexpr.ExprContextOptions{
		Stack: st,
		Path:  "path1",
	})
	ctx2 := ctx1.WithPath("path2")

	// Different paths
	require.Equal(t, "path1", ctx1.Path())
	require.Equal(t, "path2", ctx2.Path())

	// Same stack (shared)
	require.Equal(t, ctx1.Stack(), ctx2.Stack())
}
