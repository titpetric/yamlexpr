package handlers

import (
	"fmt"

	"github.com/titpetric/yamlexpr/model"
)

// NewForHandler returns a handler for the "for" directive.
// Implements loop expansion with support for:
// - Simple iteration: "item in items"
// - Index and item: "(idx, item) in items"
// - Nested paths: "item in item.subitems"
// - Omitted variables: "(_, item) in items" or "(idx, _) in items"
//
// The for directive expands a template for each item in the collection,
// making loop variables available for interpolation.
//
// Priority: 100 (runs after include, before custom handlers)
func NewForHandler() *ForHandlerImpl {
	return &ForHandlerImpl{}
}

// ForHandlerImpl implements the for directive handler
type ForHandlerImpl struct{}

// ForLoopExpr represents a parsed for loop expression
type ForLoopExpr struct {
	// Variables is a list of variable names to bind. Can include "_" to omit.
	Variables []string
	// Source is the name of the variable to iterate over
	Source string
}

// ParseForExpr parses a for loop expression string.
// Supports:
//   - "item in items" - iterates over items, binding each to 'item'
//   - "item in item.subitems" - iterate over nested path
//   - "(idx, item) in items" - iterates over items, binding index to 'idx' and item to 'item'
//   - "(key, value) in items" - for map iteration
//   - "_" can be used to omit a variable
func ParseForExpr(expr string) (*ForLoopExpr, error) {
	// Find the " in " keyword (with spaces around it)
	inIdx := -1
	parenDepth := 0
	for i := 0; i <= len(expr)-4; i++ {
		if expr[i] == '(' {
			parenDepth++
		} else if expr[i] == ')' {
			parenDepth--
		} else if parenDepth == 0 && expr[i:i+4] == " in " {
			inIdx = i
			break
		}
	}

	if inIdx == -1 {
		return nil, fmt.Errorf("no ' in ' found in for expression")
	}

	varStr := expr[:inIdx]
	sourceStr := expr[inIdx+4:]

	// Parse variables
	var vars []string
	if len(varStr) > 0 && varStr[0] == '(' {
		// Parenthesized variables: "(idx, item)"
		if len(varStr) == 0 || varStr[len(varStr)-1] != ')' {
			return nil, fmt.Errorf("missing closing parenthesis")
		}
		varStr = varStr[1 : len(varStr)-1]

		// Check for trailing comma
		if len(varStr) > 0 && varStr[len(varStr)-1] == ',' {
			return nil, fmt.Errorf("trailing comma in variable list")
		}

		// Split by comma and trim whitespace
		start := 0
		for i := 0; i <= len(varStr); i++ {
			if i == len(varStr) || varStr[i] == ',' {
				part := varStr[start:i]
				// Trim spaces
				trimmed := trimSpace(part)
				if trimmed == "" {
					return nil, fmt.Errorf("empty variable name in list")
				}
				// Variable names cannot start with a digit
				if len(trimmed) > 0 && trimmed[0] >= '0' && trimmed[0] <= '9' {
					return nil, fmt.Errorf("invalid variable name '%s': must not start with a digit", trimmed)
				}
				vars = append(vars, trimmed)
				start = i + 1
			}
		}
	} else {
		// Single variable: "item"
		trimmed := trimSpace(varStr)
		if trimmed == "" {
			return nil, fmt.Errorf("no variable name found")
		}
		// Check for trailing comma in single var (e.g., "item,")
		if trimmed[len(trimmed)-1] == ',' {
			return nil, fmt.Errorf("trailing comma in variable list")
		}
		vars = append(vars, trimmed)
	}

	source := trimSpace(sourceStr)
	if source == "" {
		return nil, fmt.Errorf("no source expression found")
	}

	return &ForLoopExpr{
		Variables: vars,
		Source:    source,
	}, nil
}

// trimSpace is a simple string trimming helper
func trimSpace(s string) string {
	start := 0
	for start < len(s) && s[start] == ' ' {
		start++
	}
	end := len(s)
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}

// BuildScope creates a variable scope map for a loop iteration
// varNames: the variable names to bind (e.g., ["idx", "item"])
// idx: the current iteration index
// item: the current item value
// Returns a map with bound variables (skips "_" variables)
func BuildScope(varNames []string, idx int, item any) map[string]any {
	scope := make(map[string]any)
	for i, varName := range varNames {
		if varName == "_" {
			// Skip underscore variables (intentional omission)
			continue
		}

		// Bind the appropriate value based on position
		switch i {
		case 0:
			// First variable is usually the item (or index if 2 variables)
			if len(varNames) == 2 {
				scope[varName] = idx
			} else {
				scope[varName] = item
			}
		case 1:
			// Second variable is the item
			scope[varName] = item
		default:
			// Additional variables (for potential future use)
			scope[varName] = item
		}
	}
	return scope
}

// ForHandler creates a handler for the "for" directive.
// Requires a Processor to recursively expand templates.
func ForHandler(proc Processor, forDirective string) DirectiveHandler {
	return func(ctx *model.Context, block map[string]any, value any) (any, bool, error) {
		// Get the collection to iterate over and parse the for expression
		var items []any
		var loopVars *ForLoopExpr

		switch v := value.(type) {
		case []any:
			// Direct array literal: for: [1, 2, 3]
			items = v
			// Default to single "item" variable for direct arrays
			loopVars = &ForLoopExpr{
				Variables: []string{"item"},
				Source:    "",
			}
		case string:
			// Parse as new for expression (e.g., "item in items" or "(idx, item) in items")
			var err error
			loopVars, err = ParseForExpr(v)
			if err != nil {
				pathCtx := ""
				if ctx.Path() != "" {
					pathCtx = fmt.Sprintf(" at %s.for", ctx.Path())
				}
				return nil, false, fmt.Errorf("invalid for expression '%s'%s: %w", v, pathCtx, err)
			}

			// Resolve the source variable from the stack
			sourceVal, ok := ctx.Stack().Resolve(loopVars.Source)
			if !ok {
				pathCtx := ""
				if ctx.Path() != "" {
					pathCtx = fmt.Sprintf(" at %s.for", ctx.Path())
				}
				return nil, false, fmt.Errorf("undefined variable '%s'%s", loopVars.Source, pathCtx)
			}

			// Convert source to slice
			if slice, ok := sourceVal.([]any); ok {
				items = slice
			} else {
				pathCtx := ""
				if ctx.Path() != "" {
					pathCtx = fmt.Sprintf(" at %s.for", ctx.Path())
				}
				return nil, false, fmt.Errorf("for: variable '%s' must be an array, got %T%s", loopVars.Source, sourceVal, pathCtx)
			}
		default:
			pathCtx := ""
			if ctx.Path() != "" {
				pathCtx = fmt.Sprintf(" at %s.for", ctx.Path())
			}
			return nil, false, fmt.Errorf("for: expected array or string expression, got %T%s", value, pathCtx)
		}

		// If empty collection, return empty slice
		if len(items) == 0 {
			return []any{}, true, nil
		}

		// Iterate over items and expand template
		result := make([]any, 0, len(items))
		for idx, item := range items {
			// Build the scope map with only non-underscore variables
			scope := make(map[string]any)
			for i, varName := range loopVars.Variables {
				if varName == "_" {
					// Skip underscore variables (intentional omission)
					continue
				}

				// Bind the appropriate value based on position
				switch i {
				case 0:
					// First variable is usually the item (or index if 2 variables)
					if len(loopVars.Variables) == 2 {
						scope[varName] = idx
					} else {
						scope[varName] = item
					}
				case 1:
					// Second variable is the item
					scope[varName] = item
				default:
					// Additional variables (for potential future use)
					scope[varName] = item
				}
			}

			// Create new stack scope with loop variables
			ctx.PushStackScope(scope)

			// Create a fresh copy of the template for each iteration (all keys except for directive)
			template := make(map[string]any)
			for k, v := range block {
				if k != forDirective {
					template[k] = v
				}
			}

			// Create context for this iteration
			itemCtx := ctx.AppendPath(fmt.Sprintf("[%d]", idx))

			// Process template with current item in scope
			expanded, err := proc.ProcessMapWithContext(itemCtx, template)
			if err != nil {
				ctx.PopStackScope()
				return nil, false, err
			}
			if expanded != nil {
				result = append(result, expanded)
			}

			// Pop the scope for this iteration
			ctx.PopStackScope()
		}

		return result, true, nil
	}
}

// ExpandForAtRoot expands a root-level for directive into multiple document configurations.
// Returns a slice of processed maps, one for each item in the iteration.
// This is used by Expr.LoadMulti to handle root-level for loop expansion.
func ExpandForAtRoot(ctx *model.Context, template map[string]any, processor Processor, forDirective string) ([]any, error) {
	forValue, ok := template[forDirective]
	if !ok {
		return nil, fmt.Errorf("for key not found")
	}

	// Parse the for expression
	loopVars, err := ParseForExpr(forValue.(string))
	if err != nil {
		return nil, fmt.Errorf("error parsing for expression: %w", err)
	}

	// Resolve the source variable from the stack
	sourceVal, ok := ctx.Stack().Resolve(loopVars.Source)
	if !ok {
		return nil, fmt.Errorf("undefined variable '%s'", loopVars.Source)
	}

	// Convert source to slice
	items, ok := sourceVal.([]any)
	if !ok {
		return nil, fmt.Errorf("for: variable '%s' must be an array, got %T", loopVars.Source, sourceVal)
	}

	// Collect all keys from template for null-filling (except for directive)
	allKeys := make(map[string]bool)
	for k := range template {
		if k != forDirective {
			allKeys[k] = true
		}
	}

	// Process each item
	result := make([]any, 0, len(items))
	for idx, item := range items {
		// Build scope with loop variables
		scope := make(map[string]any)
		for i, varName := range loopVars.Variables {
			if varName == "_" {
				continue
			}

			switch i {
			case 0:
				// First variable is usually the item (or index if 2 variables)
				if len(loopVars.Variables) == 2 {
					scope[varName] = idx
				} else {
					scope[varName] = item
				}
			case 1:
				// Second variable is the item
				scope[varName] = item
			default:
				// Additional variables (for potential future use)
				scope[varName] = item
			}
		}

		// Ensure all keys are present (fill missing with null)
		for k := range allKeys {
			if _, exists := scope[k]; !exists {
				scope[k] = nil
			}
		}

		// Create template copy without for directive
		itemTemplate := make(map[string]any)
		for k, v := range template {
			if k != forDirective {
				itemTemplate[k] = v
			}
		}

		// Create context with item variables and process template
		st := ctx.Stack()
		st.Push(scope)

		itemCtx := ctx.WithPath(ctx.Path())

		processed, err := processor.ProcessMapWithContext(itemCtx, itemTemplate)
		st.Pop()

		if err != nil {
			return nil, err
		}

		if processed != nil {
			result = append(result, processed)
		}
	}

	return result, nil
}

// Note: The old for loop handling was done in expr.go's handleForWithContext
// because it needs access to the template processing and recursive evaluation.
// This handler extracts the parsing and scope building logic for reuse.
