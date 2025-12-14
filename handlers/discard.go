package handlers

import (
	"fmt"

	"github.com/titpetric/yamlexpr/model"
)

// NewDiscardHandler returns a directive handler for the "discard" directive.
// The discard directive omits a block if set to true, similar to if: false.
func DiscardHandlerBuiltin() DirectiveHandler {
	return func(ctx *model.Context, block map[string]any, value any) (any, bool, error) {
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

// NewDiscardHandler is deprecated - use DiscardHandlerBuiltin instead.
// This function is kept for backwards compatibility.
func NewDiscardHandler() any {
	return DiscardHandlerBuiltin()
}
