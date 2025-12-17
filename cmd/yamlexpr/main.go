package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v3"

	"github.com/titpetric/yamlexpr"
)

func main() {
	testFlag := flag.Bool("test-fixtures", false, "Run fixture tests instead of processing input")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [file]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Parse and evaluate a YAML file using yamlexpr\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -test-fixtures  Run all fixture tests\n")
	}

	flag.Parse()

	if *testFlag {
		testFixtures()
		return
	}

	// Read input file or stdin
	var input []byte
	var err error

	if flag.NArg() > 0 {
		input, err = os.ReadFile(flag.Arg(0))
		if err != nil {
			log.Fatalf("error reading file: %v", err)
		}
	} else {
		input, err = io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("error reading stdin: %v", err)
		}
	}

	// Parse YAML to Document
	var docMap map[string]any
	err = yaml.Unmarshal(input, &docMap)
	if err != nil {
		log.Fatalf("error parsing YAML: %v", err)
	}

	// Create evaluator
	expr := yamlexpr.New(nil)

	// Process with yamlexpr
	docs, err := expr.Parse(yamlexpr.Document(docMap))
	if err != nil {
		log.Fatalf("error evaluating: %v", err)
	}

	// Output all documents as YAML
	for i, doc := range docs {
		output, err := yaml.Marshal(map[string]any(doc))
		if err != nil {
			log.Fatalf("error marshaling result: %v", err)
		}
		if i > 0 {
			fmt.Print("---\n")
		}
		fmt.Print(string(output))
	}
}

// testFixtures runs all fixture tests and reports results
func testFixtures() {
	fixturesDir := "testdata/fixtures"
	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		log.Fatalf("error reading fixtures directory: %v", err)
	}

	// Create filesystem rooted at fixtures directory for include support
	fixturesFS := os.DirFS(fixturesDir)
	expr := yamlexpr.New(fixturesFS)
	passed := 0
	failed := 0
	skipped := 0

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		// Skip fixtures with leading underscore silently (helper files, not actual tests)
		if strings.HasPrefix(name, "_") {
			continue
		}

		// Skip non-yaml files
		if filepath.Ext(name) != ".yaml" && filepath.Ext(name) != ".yml" {
			continue
		}

		// Load the fixture file
		data, err := os.ReadFile(filepath.Join(fixturesDir, name))
		if err != nil {
			fmt.Printf("FAIL %s: error reading file: %v\n", name, err)
			failed++
			continue
		}

		// Split by --- separator
		parts := strings.Split(string(data), "\n---\n")
		if len(parts) < 2 {
			fmt.Printf("SKIP %s: no expected output (no --- separator)\n", name)
			skipped++
			continue
		}

		// Parse input (first section) - can be map or array
		var inputDoc any
		err = yaml.Unmarshal([]byte(parts[0]), &inputDoc)
		if err != nil {
			fmt.Printf("FAIL %s: error parsing input: %v\n", name, err)
			failed++
			continue
		}

		// Convert array input to a wrapper map
		var input map[string]any
		if arr, ok := inputDoc.([]any); ok {
			// Wrap array in a temporary key for processing
			input = map[string]any{"_items": arr}
		} else if m, ok := inputDoc.(map[string]any); ok {
			input = m
		} else {
			fmt.Printf("FAIL %s: input must be map or array, got %T\n", name, inputDoc)
			failed++
			continue
		}

		// Parse expected (all remaining sections as separate documents)
		var expectedDocs []any
		for i := 1; i < len(parts); i++ {
			var doc any
			err = yaml.Unmarshal([]byte(parts[i]), &doc)
			if err != nil {
				fmt.Printf("FAIL %s: error parsing expected section %d: %v\n", name, i, err)
				failed++
				continue
			}
			expectedDocs = append(expectedDocs, doc)
		}

		// Process the input
		docs, err := expr.Parse(yamlexpr.Document(input))
		if err != nil {
			fmt.Printf("FAIL %s: error processing: %v\n", name, err)
			failed++
			continue
		}

		// If input was an array, unwrap the _items key from the result
		_, isArrayInput := inputDoc.([]any)
		var actualResults []yamlexpr.Document
		if isArrayInput && len(docs) == 1 {
			// Check if the result has _items key
			if items, ok := docs[0]["_items"]; ok {
				if itemsSlice, ok := items.([]any); ok {
					// Convert items back to Document slice
					for _, item := range itemsSlice {
						if itemMap, ok := item.(map[string]any); ok {
							actualResults = append(actualResults, yamlexpr.Document(itemMap))
						}
					}
				}
			}
		} else {
			// Regular map results
			actualResults = docs
		}

		// Compare results based on number of expected documents
		var expected any
		var result any

		if len(expectedDocs) == 1 {
			// Single expected document
			expected = expectedDocs[0]
			if len(actualResults) == 1 {
				// Single result document
				result = map[string]any(actualResults[0])
			} else {
				// Multiple result documents - convert to slice
				resultSlice := make([]any, len(actualResults))
				for i, doc := range actualResults {
					resultSlice[i] = map[string]any(doc)
				}
				result = resultSlice
			}
		} else {
			// Multiple expected documents
			expectedAny := make([]any, len(expectedDocs))
			copy(expectedAny, expectedDocs)
			expected = expectedAny

			resultSlice := make([]any, len(actualResults))
			for i, doc := range actualResults {
				resultSlice[i] = map[string]any(doc)
			}
			result = resultSlice
		}

		// Compare with expected
		if !deepEqual(expected, result) {
			fmt.Printf("FAIL %s: output mismatch\n", name)
			expYAML, _ := yaml.Marshal(expected)
			resYAML, _ := yaml.Marshal(result)
			fmt.Printf("  Expected YAML:\n%s", string(expYAML))
			fmt.Printf("  Got YAML:\n%s", string(resYAML))
			failed++
			continue
		}

		fmt.Printf("PASS %s\n", name)
		passed++
	}

	fmt.Printf("\n%d passed, %d failed, %d skipped\n", passed, failed, skipped)
	if failed > 0 {
		os.Exit(1)
	}
}

// deepEqual compares two values for equality by re-marshaling and unmarshaling
// to normalize their representation (especially map ordering)
func deepEqual(a, b any) bool {
	// Marshal both to YAML and unmarshal back to normalize map ordering
	aYAML, errA := yaml.Marshal(a)
	bYAML, errB := yaml.Marshal(b)

	if errA != nil || errB != nil {
		return false
	}

	// Unmarshal both back to generic format
	var aNorm, bNorm any
	if err := yaml.Unmarshal(aYAML, &aNorm); err != nil {
		return false
	}
	if err := yaml.Unmarshal(bYAML, &bNorm); err != nil {
		return false
	}

	// Re-marshal with sorted keys to ensure consistent comparison
	aNormYAML, _ := yaml.Marshal(aNorm)
	bNormYAML, _ := yaml.Marshal(bNorm)

	return string(aNormYAML) == string(bNormYAML)
}
