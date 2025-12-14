package handlers

import (
	"fmt"

	"github.com/titpetric/yamlexpr/model"
)

// IncludeHandler creates a handler for the "include" directive.
// Requires a Processor to load and merge files.
//
// Include is used to compose YAML files. It loads external YAML files and merges
// their content into the current structure.
func IncludeHandler(proc Processor, includeDirective string) DirectiveHandler {
	return func(ctx *model.Context, block map[string]any, value any) (any, bool, error) {
		// Result map to merge embedded content into
		result := make(map[string]any)

		// Handle single file
		if filename, ok := value.(string); ok {
			if err := proc.LoadAndMergeFileWithContext(ctx, filename, result); err != nil {
				return nil, false, err
			}
		} else if files, ok := value.([]any); ok {
			// Handle list of files
			for _, f := range files {
				if filename, ok := f.(string); ok {
					if err := proc.LoadAndMergeFileWithContext(ctx, filename, result); err != nil {
						return nil, false, err
					}
				}
			}
		} else {
			return nil, false, fmt.Errorf("include must be a string or list of strings, got %T", value)
		}

		// Return the merged content but don't consume all processing
		// This allows the current map's other keys to still be processed
		return result, false, nil
	}
}
