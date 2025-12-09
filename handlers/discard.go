package handlers

import (
	"fmt"
)

// NewDiscardHandler returns a directive handler for the "discard" directive.
// The discard directive omits a block if set to true, similar to if: false.
//
// Signature matches yamlexpr.DirectiveHandler:
//
//	func(ctx *ExprContext, block map[string]any, value any) (result any, consumed bool, err error)
//
// Usage in YAML:
//
//	steps:
//	  - name: "build"
//	    run: "npm run build"
//	  - name: "test"
//	    run: "npm test"
//	    discard: false
//	  - name: "publish"
//	    run: "npm publish"
//	    discard: true
//
// The "publish" step will be omitted in the output.
//
// Priority: 10 (runs after if/for but before matrix)
//
// Example usage:
//
//	e := yamlexpr.New(fs,
//	    yamlexpr.WithDirectiveHandler("discard", handlers.NewDiscardHandler(), 10),
//	)
func NewDiscardHandler() any {
	return func(ctx any, block map[string]any, value any) (any, bool, error) {
		// Handle various input types
		switch v := value.(type) {
		case bool:
			if v {
				// discard: true → omit this block
				return nil, true, nil
			}
			// discard: false → continue with normal processing
			return nil, false, nil

		case string:
			// Handle string representations of boolean
			switch v {
			case "true", "1", "yes":
				return nil, true, nil
			case "false", "0", "no", "":
				return nil, false, nil
			default:
				return nil, false, fmt.Errorf("discard value must be boolean or 'true'/'false', got string '%s'", v)
			}

		case int:
			if v != 0 {
				return nil, true, nil
			}
			return nil, false, nil

		case nil:
			// nil is falsy
			return nil, false, nil

		default:
			return nil, false, fmt.Errorf("discard must be boolean, got %T", value)
		}
	}
}
