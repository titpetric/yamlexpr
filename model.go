package yamlexpr

import "github.com/titpetric/yamlexpr/model"

// Model type aliases.
type (
	// Config aliases model.Config.
	Config = model.Config
	// ConfigOption aliases model.ConfigOption.
	ConfigOption = model.ConfigOption
	// DirectiveHandler aliases model.DirectiveHandler.
	DirectiveHandler = model.DirectiveHandler
	// Syntax aliases model.SyntaxHandler.
	Syntax = model.Syntax
)

// Model function/value aliases.
var (
	// DefaultConfig aliases model.DefaultConfig
	DefaultConfig = model.DefaultConfig
	// WithFS aliases model.WithFS
	WithFS = model.WithFS
	// WithSyntax aliases model.WithFS
	WithSyntax = model.WithSyntax
)
