package testdata

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/require"

	"github.com/titpetric/yamlexpr"
)

// TestReadmeExamples verifies that all README code examples actually run correctly.
// These tests use Load() with MapFS YAML files to test the real usage pattern.
func TestReadmeExamples(t *testing.T) {
	t.Run("InterpolationExample", func(t *testing.T) {
		// Create YAML file with interpolation
		fs := fstest.MapFS{
			"config.yaml": {Data: []byte(`
user_name: Alice
registry: docker.io
app: myapp
version: "1.0.0"
greeting: "Hello, ${user_name}!"
image: "${registry}/${app}:${version}"
`)},
		}

		expr := yamlexpr.New(fs)

		// Load and evaluate YAML file
		docs, err := expr.Load("config.yaml")
		require.NoError(t, err)
		require.Len(t, docs, 1)

		doc := docs[0]
		require.Equal(t, "Hello, Alice!", doc["greeting"])
		require.Equal(t, "docker.io/myapp:1.0.0", doc["image"])
	})

	t.Run("ConditionalExample", func(t *testing.T) {
		// Create YAML file with conditionals
		fs := fstest.MapFS{
			"config.yaml": {Data: []byte(`
name: production
debug:
  if: "true"
  enabled: true
  level: verbose
other: always_present
`)},
		}

		expr := yamlexpr.New(fs)

		// Load and evaluate YAML file
		docs, err := expr.Load("config.yaml")
		require.NoError(t, err)
		require.Len(t, docs, 1)

		doc := docs[0]
		require.Contains(t, doc, "debug")
		require.Contains(t, doc, "other")
		require.Equal(t, "production", doc["name"])
	})

	t.Run("ForLoopExample", func(t *testing.T) {
		// Create YAML file with for loop
		fs := fstest.MapFS{
			"config.yaml": {Data: []byte(`
server_list:
  - api
  - worker
  - cache
servers:
  - for: "server in server_list"
    name: "${server}"
`)},
		}

		expr := yamlexpr.New(fs)

		// Load and evaluate YAML file
		docs, err := expr.Load("config.yaml")
		require.NoError(t, err)
		require.Len(t, docs, 1)

		doc := docs[0]
		servers := doc["servers"].([]any)
		require.Len(t, servers, 3)
		require.Equal(t, "api", servers[0].(map[string]any)["name"])
		require.Equal(t, "worker", servers[1].(map[string]any)["name"])
		require.Equal(t, "cache", servers[2].(map[string]any)["name"])
	})

	t.Run("DocumentExpansionExample", func(t *testing.T) {
		// Create YAML file with root-level for loop (document expansion)
		fs := fstest.MapFS{
			"config.yaml": {Data: []byte(`
for: "item in items"
name: "${item}"
items:
  - alice
  - bob
  - charlie
`)},
		}

		expr := yamlexpr.New(fs)

		// Load and evaluate YAML file - should expand to multiple documents
		docs, err := expr.Load("config.yaml")
		require.NoError(t, err)
		require.Len(t, docs, 3) // Should expand to 3 documents

		require.Equal(t, "alice", docs[0]["name"])
		require.Equal(t, "bob", docs[1]["name"])
		require.Equal(t, "charlie", docs[2]["name"])
	})

	t.Run("MatrixExample", func(t *testing.T) {
		// Create YAML file with matrix (should expand to multiple documents)
		fs := fstest.MapFS{
			"config.yaml": {Data: []byte(`
matrix:
  os:
    - linux
    - windows
  version:
    - 12
    - 14
job_name: "Test ${os} v${version}"
os: "${os}"
version: "${version}"
`)},
		}

		expr := yamlexpr.New(fs)

		// Load and evaluate YAML file
		docs, err := expr.Load("config.yaml")
		require.NoError(t, err)
		require.Len(t, docs, 4) // 2 OS × 2 versions = 4 combinations

		// Verify all combinations exist
		results := make(map[string]bool)
		for _, doc := range docs {
			results[doc["job_name"].(string)] = true
		}

		require.True(t, results["Test linux v12"])
		require.True(t, results["Test linux v14"])
		require.True(t, results["Test windows v12"])
		require.True(t, results["Test windows v14"])
	})

	t.Run("IncludeCompositionExample", func(t *testing.T) {
		// Create YAML files with include directive
		fs := fstest.MapFS{
			"config.yaml": {Data: []byte(`
name: myapp
database:
  include: "_db-defaults.yaml"
  pool_size: 20
`)},
			"_db-defaults.yaml": {Data: []byte(`
host: localhost
port: 5432
timeout: 30
`)},
		}

		expr := yamlexpr.New(fs)

		// Load and evaluate YAML file
		docs, err := expr.Load("config.yaml")
		require.NoError(t, err)
		require.Len(t, docs, 1)

		doc := docs[0]
		require.Equal(t, "myapp", doc["name"])
		dbConfig := doc["database"].(map[string]any)
		require.Equal(t, "localhost", dbConfig["host"])
		require.Equal(t, 5432, dbConfig["port"])
		require.Equal(t, 20, dbConfig["pool_size"])
	})

	t.Run("CombinedFeaturesExample", func(t *testing.T) {
		// Create YAML files combining multiple features
		fs := fstest.MapFS{
			"config.yaml": {Data: []byte(`
services:
  - for: "svc in service_list"
    if: "${svc.enabled}"
    name: "${svc.name}"
    port: ${svc.port}
    include: "_service-defaults.yaml"
service_list:
  - name: api
    enabled: true
    port: 8080
  - name: disabled-worker
    enabled: false
    port: 9000
  - name: cache
    enabled: true
    port: 6379
`)},
			"_service-defaults.yaml": {Data: []byte(`
timeout: 30
retries: 3
`)},
		}

		expr := yamlexpr.New(fs)

		// Load and evaluate YAML file
		docs, err := expr.Load("config.yaml")
		require.NoError(t, err)
		require.Len(t, docs, 1)

		doc := docs[0]
		services := doc["services"].([]any)
		require.Len(t, services, 2) // Only 2 enabled services

		// Check first service
		svc1 := services[0].(map[string]any)
		require.Equal(t, "api", svc1["name"])
		require.Equal(t, 8080, svc1["port"])
		require.Equal(t, 30, svc1["timeout"])

		// Check second service
		svc2 := services[1].(map[string]any)
		require.Equal(t, "cache", svc2["name"])
		require.Equal(t, 6379, svc2["port"])
		require.Equal(t, 30, svc2["timeout"])
	})
}
