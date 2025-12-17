package testdata

// This file contains code examples from README.md to verify they compile.
// It's not executed, just verified to compile by the Go toolchain.

import (
	"fmt"
	"log"
	"os"

	"github.com/titpetric/yamlexpr"
)

// Example1_LoadAndEvaluateYAMLFromFile demonstrates loading and evaluating a YAML file.
func Example1_LoadAndEvaluateYAMLFromFile() {
	// Create an Expr evaluator with the current directory as the filesystem
	expr := yamlexpr.New(os.DirFS("."))

	// Load and evaluate a YAML file
	// Returns a slice of Documents (may be multiple if using root-level for: or matrix:)
	docs, err := expr.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Process each resulting document
	for i, doc := range docs {
		fmt.Printf("Document %d: %+v\n", i, doc)
	}
}

// Example2_ParseAndEvaluateYAMLData demonstrates parsing and evaluating YAML data.
func Example2_ParseAndEvaluateYAMLData() {
	// Create an Expr evaluator
	expr := yamlexpr.New(os.DirFS("."))

	// Prepare your YAML data as Document (map[string]any)
	data := yamlexpr.Document{
		"name": "production",
		"env": map[string]any{
			"debug": false,
			"port":  8080,
		},
		"services": []any{
			map[string]any{"name": "api", "active": true},
			map[string]any{"name": "worker", "active": false},
		},
	}

	// Parse with variable interpolation, conditionals, and composition
	docs, err := expr.Parse(data)
	if err != nil {
		log.Fatal(err)
	}

	// Process each resulting document
	for i, doc := range docs {
		fmt.Printf("Document %d: %+v\n", i, doc)
	}
}

// Example3_UseCustomDirectiveSyntax demonstrates using custom directive syntax.
func Example3_UseCustomDirectiveSyntax() {
	// Use Vue.js-style directives
	expr := yamlexpr.New(os.DirFS("."), yamlexpr.WithSyntax(yamlexpr.Syntax{
		If:      "v-if",
		For:     "v-for",
		Include: "v-include",
		Matrix:  "v-matrix",
	}))

	// Now your YAML can use v-if, v-for, v-include, and v-matrix
	docs, err := expr.Load("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// Process resulting documents
	for _, doc := range docs {
		fmt.Printf("%+v\n", doc)
	}
}
