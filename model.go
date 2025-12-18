package yamlexpr

import (
	"github.com/titpetric/yamlexpr/frontmatter"
	"github.com/titpetric/yamlexpr/model"
)

// Model type aliases.
type (
	// Config aliases model.Config.
	Config = model.Config
	// ConfigOption aliases model.ConfigOption.
	ConfigOption = model.ConfigOption
	// Context aliases model.Context.
	Context = model.Context
	// ContextOptions aliases model.ContextOptions.
	ContextOptions = model.ContextOptions
	// DirectiveHandler aliases model.DirectiveHandler.
	DirectiveHandler = model.DirectiveHandler
	// Syntax aliases model.SyntaxHandler.
	Syntax = model.Syntax
	// DocumentContent aliases frontmatter.DocumentContent.
	DocumentContent = frontmatter.DocumentContent
)

// Model function/value aliases.
var (
	// DefaultConfig aliases model.DefaultConfig.
	DefaultConfig = model.DefaultConfig
	// NewContext aliases model.NewContext.
	NewContext = model.NewContext
	// WithFS aliases model.WithFS.
	WithFS = model.WithFS
	// WithSyntax aliases model.WithFS.
	WithSyntax = model.WithSyntax
	// ParseDocument aliases frontmatter.ParseDocument.
	ParseDocument = frontmatter.ParseDocument
)
