package yamlexpr

import (
	"github.com/titpetric/yamlexpr/handlers"
	"github.com/titpetric/yamlexpr/model"
)

// ExprContext is an alias for model.Context.
type ExprContext = model.Context

// ExprContextOptions is an alias for model.ContextOptions.
type ExprContextOptions = model.ContextOptions

// Matrix is an alias for handlers.MatrixDirective.
type Matrix = handlers.MatrixDirective

// NewExprContext creates a new evaluation context with the given options.
func NewExprContext(options *ExprContextOptions) *ExprContext {
	return model.NewContext((*model.ContextOptions)(options))
}
