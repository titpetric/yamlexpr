package yamlexpr

import (
	"io/fs"

	"github.com/titpetric/yamlexpr/model"
)

// DirectiveHandler is a backwards-compatible alias for model.DirectiveHandler.
type DirectiveHandler = model.DirectiveHandler

// Syntax defines the directive keywords used in YAML documents.
// Empty fields retain their default values when merged with defaults.
type Syntax struct {
	// If is the directive keyword for conditional blocks (default: "if").
	If string `json:"if" yaml:"if"`
	// For is the directive keyword for iteration blocks (default: "for").
	For string `json:"for" yaml:"for"`
	// Embed is the directive keyword for file embedding/inclusion (default: "embed").
	Embed string `json:"embed" yaml:"embed"`
}

// DefaultSyntax is the default syntax configuration with standard directive names.
var DefaultSyntax = Syntax{
	If:    "if",
	For:   "for",
	Embed: "embed",
}

// Config holds configuration options for the Expr evaluator.
type Config struct {
	// syntax defines the directive keywords used in YAML documents
	syntax Syntax
	// handlers maps directive names to their handler functions
	handlers map[string]DirectiveHandler
	// handlerOrder tracks the order handlers were registered (for deterministic evaluation)
	handlerOrder []string
	// registerStandard indicates if standard handlers should be registered
	registerStandard bool
	// filesystem is the FS used for loading resources (can be nil)
	filesystem fs.FS
}

// DefaultConfig returns the default configuration with standard directive names.
func DefaultConfig() *Config {
	return &Config{
		syntax:   DefaultSyntax,
		handlers: make(map[string]DirectiveHandler),
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
//		If:    "v-if",
//		For:   "v-for",
//		Embed: "v-embed",
//	}))
//
// Or partially customize (empty fields keep defaults):
//
//	e := yamlexpr.New(fs, yamlexpr.WithSyntax(yamlexpr.Syntax{
//		If:  "v-if",
//		For: "v-for",
//		// Embed remains "embed"
//	}))
func WithSyntax(syntax Syntax) ConfigOption {
	return func(cfg *Config) {
		// Apply non-empty fields, keeping defaults for empty ones
		if syntax.If != "" {
			cfg.syntax.If = syntax.If
		}
		if syntax.For != "" {
			cfg.syntax.For = syntax.For
		}
		if syntax.Embed != "" {
			cfg.syntax.Embed = syntax.Embed
		}
	}
}

// IfDirective returns the current if directive keyword.
func (c *Config) IfDirective() string {
	return c.syntax.If
}

// ForDirective returns the current for directive keyword.
func (c *Config) ForDirective() string {
	return c.syntax.For
}

// EmbedDirective returns the current embed directive keyword.
func (c *Config) EmbedDirective() string {
	return c.syntax.Embed
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
// If a handler is registered for a built-in directive (if, for, embed),
// it overrides the default implementation for that directive.
func WithDirectiveHandler(directive string, handler DirectiveHandler) ConfigOption {
	return func(cfg *Config) {
		if cfg.handlers == nil {
			cfg.handlers = make(map[string]DirectiveHandler)
		}
		// Track order of registration if this is a new handler
		if _, exists := cfg.handlers[directive]; !exists {
			cfg.handlerOrder = append(cfg.handlerOrder, directive)
		}
		cfg.handlers[directive] = handler
	}
}

// WithFS sets the filesystem for resource loading (embed directive).
// If not set, only in-memory processing is available.
//
// Example:
//
//	e := yamlexpr.New(yamlexpr.WithFS(myFS), yamlexpr.WithStandardHandlers())
func WithFS(filesystem fs.FS) ConfigOption {
	return func(cfg *Config) {
		cfg.filesystem = filesystem
	}
}

// WithStandardHandlers registers the standard handlers (for, if, embed).
// This is a convenience option to enable the built-in directives.
//
// Example:
//
//	e := yamlexpr.New(yamlexpr.WithFS(fs), yamlexpr.WithStandardHandlers())
//
// This is equivalent to manually registering each handler:
//
//	e := yamlexpr.New(yamlexpr.WithFS(fs),
//		yamlexpr.WithDirectiveHandler("for", handlers.ForHandlerBuiltin(e, "for")),
//		yamlexpr.WithDirectiveHandler("if", handlers.IfHandlerBuiltin("if")),
//		yamlexpr.WithDirectiveHandler("embed", handlers.EmbedHandlerBuiltin(e, "embed")),
//	)
func WithStandardHandlers() ConfigOption {
	return func(cfg *Config) {
		cfg.registerStandard = true
	}
}

// GetHandler returns the handler for a directive, or nil if not registered.
func (c *Config) GetHandler(directive string) DirectiveHandler {
	if c.handlers == nil {
		return nil
	}
	return c.handlers[directive]
}
