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

	// Parse YAML to any
	var data any
	err = yaml.Unmarshal(input, &data)
	if err != nil {
		log.Fatalf("error parsing YAML: %v", err)
	}

	// Create evaluator with standard handlers
	expr := yamlexpr.New(nil)

	// Process with yamlexpr
	result, err := expr.Process(data, nil)
	if err != nil {
		log.Fatalf("error evaluating: %v", err)
	}

	// Output as YAML
	output, err := yaml.Marshal(result)
	if err != nil {
		log.Fatalf("error marshaling result: %v", err)
	}

	fmt.Print(string(output))
}

// testFixtures runs all fixture tests and reports results
func testFixtures() {
	fixturesDir := "testdata/fixtures"
	entries, err := os.ReadDir(fixturesDir)
	if err != nil {
		log.Fatalf("error reading fixtures directory: %v", err)
	}

	expr := yamlexpr.New(nil)
	passed := 0
	failed := 0
	skipped := 0

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()

		// Skip non-yaml files and fixtures with leading underscore
		if strings.HasPrefix(name, "_") {
			skipped++
			continue
		}
		if strings.HasPrefix(name, "2") {
			// Skip matrix fixtures (200+ range)
			skipped++
			continue
		}
		if filepath.Ext(name) != ".yaml" {
			skipped++
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

		// Parse input and expected
		var input any
		err = yaml.Unmarshal([]byte(parts[0]), &input)
		if err != nil {
			fmt.Printf("FAIL %s: error parsing input: %v\n", name, err)
			failed++
			continue
		}

		var expected any
		err = yaml.Unmarshal([]byte(parts[1]), &expected)
		if err != nil {
			fmt.Printf("FAIL %s: error parsing expected: %v\n", name, err)
			failed++
			continue
		}

		// Process the input
		result, err := expr.Process(input, nil)
		if err != nil {
			fmt.Printf("FAIL %s: error processing: %v\n", name, err)
			failed++
			continue
		}

		// Compare with expected
		if !deepEqual(expected, result) {
			fmt.Printf("FAIL %s: output mismatch\n", name)
			fmt.Printf("  Expected: %v\n", expected)
			fmt.Printf("  Got: %v\n", result)
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

// deepEqual compares two values for equality recursively
func deepEqual(a, b any) bool {
	// Use YAML marshaling for comparison to handle floating point precision
	aYAML, errA := yaml.Marshal(a)
	bYAML, errB := yaml.Marshal(b)

	if errA != nil || errB != nil {
		return false
	}

	return string(aYAML) == string(bYAML)
}
