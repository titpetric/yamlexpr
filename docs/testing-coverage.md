# Testing Coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package                          | Coverage | Cognitive | Lines |
|--------|----------------------------------|----------|-----------|-------|
| ❌     | titpetric/yamlexpr               | 79.29%   | 356       | 994   |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | 5.22%    | 161       | 655   |
| ✅     | titpetric/yamlexpr/frontmatter   | 0.00%    | 0         | 0     |
| ✅     | titpetric/yamlexpr/interpolation | 0.00%    | 0         | 0     |
| ✅     | titpetric/yamlexpr/model         | 0.00%    | 0         | 0     |
| ✅     | titpetric/yamlexpr/stack         | 0.00%    | 0         | 0     |

## Functions

| Status | Package                          | Function                                       | Coverage | Cognitive |
|--------|----------------------------------|------------------------------------------------|----------|-----------|
| ✅     | titpetric/yamlexpr               | Expr.Load                                      | 69.20%   | 4         |
| ✅     | titpetric/yamlexpr               | Expr.Parse                                     | 92.30%   | 9         |
| ❌     | titpetric/yamlexpr               | Expr.handleForWithContext                      | 66.10%   | 42        |
| ❌     | titpetric/yamlexpr               | Expr.handleIncludeWithContext                  | 77.80%   | 11        |
| ❌     | titpetric/yamlexpr               | Expr.handleMatrixWithContext                   | 0.00%    | 56        |
| ✅     | titpetric/yamlexpr               | Expr.loadAndMergeFileWithContext               | 82.40%   | 7         |
| ✅     | titpetric/yamlexpr               | Expr.process                                   | 100.00%  | 4         |
| ✅     | titpetric/yamlexpr               | Expr.processMapWithContext                     | 95.80%   | 15        |
| ✅     | titpetric/yamlexpr               | Expr.processSliceWithContext                   | 87.10%   | 40        |
| ✅     | titpetric/yamlexpr               | Expr.processWithContext                        | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr               | Expr.processWithStack                          | 75.00%   | 1         |
| ❌     | titpetric/yamlexpr               | MapMatchesSpec                                 | 0.00%    | 5         |
| ✅     | titpetric/yamlexpr               | New                                            | 75.00%   | 1         |
| ✅     | titpetric/yamlexpr               | ValuesEqual                                    | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr               | applyExcludes                                  | 0.00%    | 7         |
| ❌     | titpetric/yamlexpr               | applyIncludes                                  | 0.00%    | 19        |
| ❌     | titpetric/yamlexpr               | evaluateConditionWithPath                      | 54.50%   | 23        |
| ❌     | titpetric/yamlexpr               | expandMatrixBase                               | 0.00%    | 15        |
| ✅     | titpetric/yamlexpr               | isQuoted                                       | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr               | isTruthy                                       | 22.20%   | 1         |
| ❌     | titpetric/yamlexpr               | isValidVarName                                 | 0.00%    | 1         |
| ✅     | titpetric/yamlexpr               | mergeRecursive                                 | 90.90%   | 24        |
| ❌     | titpetric/yamlexpr               | parseForExpr                                   | 0.00%    | 13        |
| ❌     | titpetric/yamlexpr               | parseMatrixDirective                           | 0.00%    | 30        |
| ✅     | titpetric/yamlexpr               | parseYAML                                      | 75.00%   | 1         |
| ✅     | titpetric/yamlexpr               | quoteUnquotedComparisons                       | 84.60%   | 16        |
| ❌     | titpetric/yamlexpr               | valuesEqual                                    | 0.00%    | 7         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | ExtractFrontmatterField                        | 0.00%    | 2         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | ExtractInput                                   | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | ExtractOutput                                  | 0.00%    | 2         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | GenCommand.Help                                | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | GenCommand.Name                                | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | GenCommand.Run                                 | 0.00%    | 1         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | NewGenCommand                                  | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | NewProcessCommand                              | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | NewTestCommand                                 | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | ProcessCommand.Help                            | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | ProcessCommand.Name                            | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | ProcessCommand.Run                             | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | ProcessCommand.run                             | 0.00%    | 13        |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | TestCommand.Help                               | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | TestCommand.Name                               | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | TestCommand.Run                                | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | TestCommand.run                                | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | countFixtures                                  | 0.00%    | 12        |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | generateDocs                                   | 0.00%    | 19        |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | generateDocsForFeature                         | 0.00%    | 1         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | generateJSONDocs                               | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | generateMarkdownDocs                           | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | generateMarkdownDocsInternal                   | 0.00%    | 27        |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | main                                           | 0.00%    | 1         |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | parseFixture                                   | 90.00%   | 3         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | parseFixtureDoc                                | 66.70%   | 11        |
| ✅     | titpetric/yamlexpr/cmd/yamlexpr  | printUsage                                     | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | runFixtureTests                                | 0.00%    | 16        |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | start                                          | 0.00%    | 9         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | testFixture                                    | 0.00%    | 39        |
| ✅     | titpetric/yamlexpr/frontmatter   | DocumentContent.GetFrontmatterField            | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/frontmatter   | DocumentContent.GetFrontmatterFieldWithDefault | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/frontmatter   | ParseDocument                                  | 0.00%    | 16        |
| ❌     | titpetric/yamlexpr/interpolation | ContainsInterpolation                          | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/interpolation | InterpolateString                              | 0.00%    | 10        |
| ❌     | titpetric/yamlexpr/interpolation | InterpolateStringPermissive                    | 0.00%    | 9         |
| ❌     | titpetric/yamlexpr/interpolation | InterpolateStringWithContext                   | 0.00%    | 29        |
| ❌     | titpetric/yamlexpr/interpolation | InterpolateValue                               | 0.00%    | 21        |
| ❌     | titpetric/yamlexpr/interpolation | InterpolateValuePermissive                     | 0.00%    | 3         |
| ❌     | titpetric/yamlexpr/interpolation | InterpolateValueWithContext                    | 0.00%    | 11        |
| ❌     | titpetric/yamlexpr/interpolation | extractSingleExpression                        | 0.00%    | 2         |
| ❌     | titpetric/yamlexpr/interpolation | isSingleInterpolation                          | 0.00%    | 2         |
| ✅     | titpetric/yamlexpr/model         | Config.ForDirective                            | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | Config.IfDirective                             | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | Config.IncludeDirective                        | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | Config.MatrixDirective                         | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/model         | Context.AppendPath                             | 0.00%    | 8         |
| ✅     | titpetric/yamlexpr/model         | Context.Count                                  | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/model         | Context.FormatIncludeChain                     | 0.00%    | 1         |
| ✅     | titpetric/yamlexpr/model         | Context.Path                                   | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.Pop                                    | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.Push                                   | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.Stack                                  | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.WithInclude                            | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.WithPath                               | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | DefaultConfig                                  | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/model         | NewContext                                     | 0.00%    | 3         |
| ❌     | titpetric/yamlexpr/model         | WithDirectiveHandler                           | 0.00%    | 4         |
| ✅     | titpetric/yamlexpr/model         | WithFS                                         | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/model         | WithSyntax                                     | 0.00%    | 8         |
| ❌     | titpetric/yamlexpr/stack         | CanDescend                                     | 0.00%    | 5         |
| ❌     | titpetric/yamlexpr/stack         | IsSlice                                        | 0.00%    | 2         |
| ✅     | titpetric/yamlexpr/stack         | New                                            | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/stack         | NewStack                                       | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/stack         | NewStackWithData                               | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/stack         | PopulateStructFields                           | 0.00%    | 17        |
| ❌     | titpetric/yamlexpr/stack         | ResolveValue                                   | 0.00%    | 2         |
| ❌     | titpetric/yamlexpr/stack         | SliceToAny                                     | 0.00%    | 4         |
| ❌     | titpetric/yamlexpr/stack         | Stack.All                                      | 0.00%    | 4         |
| ✅     | titpetric/yamlexpr/stack         | Stack.Copy                                     | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/stack         | Stack.Count                                    | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/stack         | Stack.ForEach                                  | 0.00%    | 12        |
| ❌     | titpetric/yamlexpr/stack         | Stack.GetInt                                   | 0.00%    | 5         |
| ❌     | titpetric/yamlexpr/stack         | Stack.GetMap                                   | 0.00%    | 5         |
| ❌     | titpetric/yamlexpr/stack         | Stack.GetSlice                                 | 0.00%    | 3         |
| ❌     | titpetric/yamlexpr/stack         | Stack.GetString                                | 0.00%    | 3         |
| ❌     | titpetric/yamlexpr/stack         | Stack.Lookup                                   | 0.00%    | 6         |
| ❌     | titpetric/yamlexpr/stack         | Stack.Pop                                      | 0.00%    | 6         |
| ❌     | titpetric/yamlexpr/stack         | Stack.Push                                     | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/stack         | Stack.Resolve                                  | 0.00%    | 7         |
| ❌     | titpetric/yamlexpr/stack         | Stack.Set                                      | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/stack         | Stack.resolveStep                              | 0.00%    | 8         |
| ❌     | titpetric/yamlexpr/stack         | StructToMap                                    | 0.00%    | 18        |
| ❌     | titpetric/yamlexpr/stack         | getCachedPath                                  | 0.00%    | 2         |
| ❌     | titpetric/yamlexpr/stack         | resolveMap                                     | 0.00%    | 1         |
| ❌     | titpetric/yamlexpr/stack         | resolveSliceIndex                              | 0.00%    | 2         |
| ❌     | titpetric/yamlexpr/stack         | resolveStruct                                  | 0.00%    | 6         |
| ❌     | titpetric/yamlexpr/stack         | resolveValueRecursive                          | 0.00%    | 4         |
| ❌     | titpetric/yamlexpr/stack         | splitPathImpl                                  | 0.00%    | 31        |
