# Testing Coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package       | Coverage | Cognitive | Lines |
|--------|---------------|----------|-----------|-------|
| ❌     | .             | 61.65%   | 356       | 924   |
| ❌     | cmd/yamlexpr  | 0.00%    | 67        | 338   |
| ✅     | frontmatter   | 100.00%  | 17        | 63    |
| ❌     | interpolation | 72.39%   | 88        | 268   |
| ❌     | model         | 46.30%   | 24        | 134   |
| ❌     | stack         | 74.20%   | 156       | 534   |

## Functions

| Status | Package       | Function                                       | Coverage | Cognitive |
|--------|---------------|------------------------------------------------|----------|-----------|
| ✅     | .             | Expr.Load                                      | 69.23%   | 4         |
| ❌     |               | Expr.Parse                                     | 69.23%   | 9         |
| ❌     |               | Expr.handleForWithContext                      | 48.21%   | 42        |
| ❌     |               | Expr.handleIncludeWithContext                  | 22.22%   | 11        |
| ❌     |               | Expr.handleMatrixWithContext                   | 0.00%    | 56        |
| ❌     |               | Expr.loadAndMergeFileWithContext               | 11.76%   | 7         |
| ✅     |               | Expr.process                                   | 100.00%  | 4         |
| ✅     |               | Expr.processMapWithContext                     | 87.50%   | 15        |
| ❌     |               | Expr.processSliceWithContext                   | 67.74%   | 40        |
| ✅     |               | Expr.processWithContext                        | 100.00%  | 1         |
| ✅     |               | Expr.processWithStack                          | 75.00%   | 1         |
| ✅     |               | MapMatchesSpec                                 | 100.00%  | 5         |
| ✅     |               | New                                            | 75.00%   | 1         |
| ✅     |               | ValuesEqual                                    | 100.00%  | 0         |
| ✅     |               | applyExcludes                                  | 100.00%  | 7         |
| ❌     |               | applyIncludes                                  | 63.16%   | 19        |
| ❌     |               | evaluateConditionWithPath                      | 54.55%   | 23        |
| ✅     |               | expandMatrixBase                               | 100.00%  | 15        |
| ✅     |               | isQuoted                                       | 100.00%  | 3         |
| ✅     |               | isTruthy                                       | 22.22%   | 1         |
| ✅     |               | isValidVarName                                 | 75.00%   | 1         |
| ✅     |               | mergeRecursive                                 | 90.91%   | 24        |
| ✅     |               | parseForExpr                                   | 95.24%   | 13        |
| ✅     |               | parseMatrixDirective                           | 90.91%   | 30        |
| ✅     |               | parseYAML                                      | 75.00%   | 1         |
| ✅     |               | quoteUnquotedComparisons                       | 84.62%   | 16        |
| ✅     |               | valuesEqual                                    | 92.86%   | 7         |
| ✅     | cmd/yamlexpr  | GenCommand.Help                                | 0.00%    | 0         |
| ❌     |               | GenCommand.Run                                 | 0.00%    | 10        |
| ✅     |               | NewGenCommand                                  | 0.00%    | 0         |
| ✅     |               | NewProcessCommand                              | 0.00%    | 0         |
| ✅     |               | NewTestCommand                                 | 0.00%    | 0         |
| ✅     |               | ProcessCommand.Help                            | 0.00%    | 0         |
| ❌     |               | ProcessCommand.Run                             | 0.00%    | 8         |
| ❌     |               | ProcessCommand.processStdin                    | 0.00%    | 3         |
| ✅     |               | TestCommand.Help                               | 0.00%    | 0         |
| ❌     |               | TestCommand.Run                                | 0.00%    | 12        |
| ❌     |               | main                                           | 0.00%    | 1         |
| ❌     |               | printDocuments                                 | 0.00%    | 5         |
| ✅     |               | printUsage                                     | 0.00%    | 0         |
| ❌     |               | renderFixture                                  | 0.00%    | 7         |
| ❌     |               | runFixture                                     | 0.00%    | 12        |
| ❌     |               | start                                          | 0.00%    | 9         |
| ✅     | frontmatter   | DocumentContent.GetFrontmatterField            | 100.00%  | 0         |
| ✅     |               | DocumentContent.GetFrontmatterFieldWithDefault | 100.00%  | 1         |
| ✅     |               | ParseDocument                                  | 100.00%  | 16        |
| ✅     | interpolation | ContainsInterpolation                          | 100.00%  | 1         |
| ✅     |               | InterpolateString                              | 81.48%   | 10        |
| ✅     |               | InterpolateStringPermissive                    | 90.00%   | 9         |
| ❌     |               | InterpolateStringWithContext                   | 64.10%   | 29        |
| ❌     |               | InterpolateValue                               | 35.29%   | 21        |
| ✅     |               | InterpolateValuePermissive                     | 100.00%  | 3         |
| ❌     |               | InterpolateValueWithContext                    | 73.33%   | 11        |
| ✅     |               | extractSingleExpression                        | 80.00%   | 2         |
| ✅     |               | isSingleInterpolation                          | 100.00%  | 2         |
| ✅     | model         | Config.ForDirective                            | 100.00%  | 0         |
| ✅     |               | Config.IfDirective                             | 100.00%  | 0         |
| ✅     |               | Config.IncludeDirective                        | 100.00%  | 0         |
| ✅     |               | Config.MatrixDirective                         | 100.00%  | 0         |
| ✅     |               | Context.AppendPath                             | 81.82%   | 8         |
| ✅     |               | Context.Count                                  | 0.00%    | 0         |
| ❌     |               | Context.FormatIncludeChain                     | 0.00%    | 1         |
| ✅     |               | Context.Path                                   | 100.00%  | 0         |
| ✅     |               | Context.Pop                                    | 100.00%  | 0         |
| ✅     |               | Context.Push                                   | 100.00%  | 0         |
| ✅     |               | Context.Stack                                  | 100.00%  | 0         |
| ✅     |               | Context.WithInclude                            | 0.00%    | 0         |
| ✅     |               | Context.WithPath                               | 100.00%  | 0         |
| ✅     |               | DefaultConfig                                  | 100.00%  | 0         |
| ✅     |               | NewContext                                     | 75.00%   | 3         |
| ❌     |               | WithDirectiveHandler                           | 0.00%    | 4         |
| ✅     |               | WithFS                                         | 0.00%    | 0         |
| ❌     |               | WithSyntax                                     | 0.00%    | 8         |
| ✅     | stack         | CanDescend                                     | 100.00%  | 5         |
| ✅     |               | IsSlice                                        | 100.00%  | 2         |
| ✅     |               | New                                            | 100.00%  | 0         |
| ✅     |               | NewStack                                       | 100.00%  | 0         |
| ✅     |               | NewStackWithData                               | 100.00%  | 1         |
| ❌     |               | PopulateStructFields                           | 0.00%    | 17        |
| ✅     |               | ResolveValue                                   | 100.00%  | 2         |
| ✅     |               | SliceToAny                                     | 100.00%  | 4         |
| ✅     |               | Stack.All                                      | 85.71%   | 4         |
| ✅     |               | Stack.Copy                                     | 0.00%    | 0         |
| ✅     |               | Stack.Count                                    | 0.00%    | 0         |
| ✅     |               | Stack.ForEach                                  | 100.00%  | 12        |
| ✅     |               | Stack.GetInt                                   | 73.33%   | 5         |
| ✅     |               | Stack.GetMap                                   | 100.00%  | 5         |
| ✅     |               | Stack.GetSlice                                 | 100.00%  | 3         |
| ✅     |               | Stack.GetString                                | 72.73%   | 3         |
| ❌     |               | Stack.Lookup                                   | 71.43%   | 6         |
| ✅     |               | Stack.Pop                                      | 90.91%   | 6         |
| ✅     |               | Stack.Push                                     | 66.67%   | 1         |
| ✅     |               | Stack.Resolve                                  | 92.31%   | 7         |
| ✅     |               | Stack.Set                                      | 66.67%   | 1         |
| ✅     |               | Stack.resolveStep                              | 100.00%  | 8         |
| ❌     |               | StructToMap                                    | 0.00%    | 18        |
| ✅     |               | getCachedPath                                  | 100.00%  | 2         |
| ✅     |               | resolveMap                                     | 100.00%  | 1         |
| ✅     |               | resolveSliceIndex                              | 100.00%  | 2         |
| ✅     |               | resolveStruct                                  | 100.00%  | 6         |
| ✅     |               | resolveValueRecursive                          | 88.89%   | 4         |
| ✅     |               | splitPathImpl                                  | 85.37%   | 31        |
