package yamlexpr

// ForLoopExpr represents a parsed for loop expression.
// It holds the variable names to bind and the source collection to iterate over.
//
// Variables can include "_" to omit a specific binding position (e.g., ignoring index).
//
// Example: for the expression "item in items", Variables is ["item"] and Source is "items".
// For the expression "(idx, item) in items", Variables is ["idx", "item"] and Source is "items".
type ForLoopExpr struct {
	// Variables is a list of variable names to bind. Can include "_" to omit.
	Variables []string

	// Source is the name of the variable to iterate over.
	Source string
}
