package model

import (
	"io/fs"
)

// Syntax defines the directive keywords used in YAML documents.
// Empty fields retain their default values when merged with defaults.
type Syntax struct {
	// If is the directive keyword for conditional blocks (default: "if").
	If string `json:"if" yaml:"if"`
	// For is the directive keyword for iteration blocks (default: "for").
	For string `json:"for" yaml:"for"`
	// Include is the directive keyword for file inclusion/composition (default: "include").
	Include string `json:"include" yaml:"include"`
	// Matrix is the directive keyword for matrix iteration (default: "matrix").
	Matrix string `json:"matrix" yaml:"matrix"`
}

// DefaultSyntax is the default syntax configuration with standard directive names.
var DefaultSyntax = Syntax{
	If:      "if",
	For:     "for",
	Include: "include",
	Matrix:  "matrix",
}

// Config holds configuration options for the Expr evaluator.
type Config struct {
	// syntax defines the directive keywords used in YAML documents
	Syntax Syntax
	// handlers maps directive names to their handler functions
	Handlers map[string]DirectiveHandler
	// handlerOrder tracks the order handlers were registered (for deterministic evaluation)
	HandlerOrder []string
	// filesystem is the FS used for loading resources (can be nil)
	FS fs.FS
}

// DefaultConfig returns the default configuration with standard directive names.
func DefaultConfig() *Config {
	return &Config{
		Syntax:   DefaultSyntax,
		Handlers: make(map[string]DirectiveHandler),
	}
}

// ConfigOption is a functional option for configuring an Expr instance.
type ConfigOption func(*Config)

// WithSyntax sets custom directive syntax, preserving defaults for empty fields.
// Empty string values in the Syntax struct will use the default keywords.
//
// Example:
//
//	e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
//		If:      "v-if",
//		For:     "v-for",
//		Include: "v-include",
//	}))
//
// Or partially customize (empty fields keep defaults):
//
//	e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
//		If:  "v-if",
//		For: "v-for",
//		// Include remains "include"
//	}))
func WithSyntax(syntax Syntax) ConfigOption {
	return func(cfg *Config) {
		// Apply non-empty fields, keeping defaults for empty ones
		if syntax.If != "" {
			cfg.Syntax.If = syntax.If
		}
		if syntax.For != "" {
			cfg.Syntax.For = syntax.For
		}
		if syntax.Include != "" {
			cfg.Syntax.Include = syntax.Include
		}
		if syntax.Matrix != "" {
			cfg.Syntax.Matrix = syntax.Matrix
		}
	}
}

// IfDirective returns the current if directive keyword.
func (c *Config) IfDirective() string {
	return c.Syntax.If
}

// ForDirective returns the current for directive keyword.
func (c *Config) ForDirective() string {
	return c.Syntax.For
}

// IncludeDirective returns the current include directive keyword.
func (c *Config) IncludeDirective() string {
	return c.Syntax.Include
}

// MatrixDirective returns the current matrix directive keyword.
func (c *Config) MatrixDirective() string {
	return c.Syntax.Matrix
}

// WithDirectiveHandler registers a custom handler for a directive name.
// The handler will be called for any block containing the specified directive.
//
// Example:
//
//	e := yamlexpr.New(fs,
//		yamlexpr.WithDirectiveHandler("matrix", myMatrixHandler),
//		yamlexpr.WithDirectiveHandler("repeat", myRepeatHandler),
//	)
//
// If a handler is registered for a built-in directive (if, for, include),
// it overrides the default implementation for that directive.
func WithDirectiveHandler(directive string, handler DirectiveHandler) ConfigOption {
	return func(cfg *Config) {
		if cfg.Handlers == nil {
			cfg.Handlers = make(map[string]DirectiveHandler)
		}
		// Track order of registration if this is a new handler
		if _, exists := cfg.Handlers[directive]; !exists {
			cfg.HandlerOrder = append(cfg.HandlerOrder, directive)
		}
		cfg.Handlers[directive] = handler
	}
}

// WithFS sets the filesystem for resource loading (include directive).
// If not set, only in-memory processing is available.
//
// Example:
//
//	e := yamlexpr.New(yamlexpr.WithFS(myFS))
func WithFS(filesystem fs.FS) ConfigOption {
	return func(cfg *Config) {
		cfg.FS = filesystem
	}
}
