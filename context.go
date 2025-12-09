package yamlexpr

import (
	"github.com/titpetric/yamlexpr/model"
)

// ExprContext is a backwards-compatible alias for model.Context.
type ExprContext = model.Context

// ExprContextOptions is a backwards-compatible alias for model.ContextOptions.
type ExprContextOptions = model.ContextOptions

// NewExprContext is a backwards-compatible alias for model.NewContext
func NewExprContext(options *ExprContextOptions) *ExprContext {
	return model.NewContext((*model.ContextOptions)(options))
}
