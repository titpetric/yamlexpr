package yamlexpr

// Syntax defines the directive keywords used in YAML documents.
// Empty fields retain their default values when merged with defaults.
type Syntax struct {
	// If is the directive keyword for conditional blocks (default: "if")
	If string `json:"if" yaml:"if"`
	// For is the directive keyword for iteration blocks (default: "for")
	For string `json:"for" yaml:"for"`
	// Include is the directive keyword for file inclusion (default: "include")
	Include string `json:"include" yaml:"include"`
}

// DefaultSyntax is the default syntax configuration with standard directive names.
var DefaultSyntax = Syntax{
	If:      "if",
	For:     "for",
	Include: "include",
}

// Config holds configuration options for the Expr evaluator.
type Config struct {
	// syntax defines the directive keywords used in YAML documents
	syntax Syntax
}

// DefaultConfig returns the default configuration with standard directive names.
func DefaultConfig() *Config {
	return &Config{
		syntax: DefaultSyntax,
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
			cfg.syntax.If = syntax.If
		}
		if syntax.For != "" {
			cfg.syntax.For = syntax.For
		}
		if syntax.Include != "" {
			cfg.syntax.Include = syntax.Include
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

// IncludeDirective returns the current include directive keyword.
func (c *Config) IncludeDirective() string {
	return c.syntax.Include
}
