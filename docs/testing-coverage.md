# Testing Coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package                          | Coverage | Cognitive | Lines |
|--------|----------------------------------|----------|-----------|-------|
| ❌     | titpetric/yamlexpr               | 54.00%   | 359       | 988   |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | 0.00%    | 80        | 248   |
| ✅     | titpetric/yamlexpr/interpolation | 80.47%   | 88        | 290   |
| ❌     | titpetric/yamlexpr/model         | 64.27%   | 24        | 191   |
| ❌     | titpetric/yamlexpr/stack         | 79.10%   | 156       | 582   |

## Functions

| Status | Package                          | Function                         | Coverage | Cognitive |
|--------|----------------------------------|----------------------------------|----------|-----------|
| ✅     | titpetric/yamlexpr               | Expr.Load                        | 69.20%   | 4         |
| ❌     | titpetric/yamlexpr               | Expr.Parse                       | 61.50%   | 9         |
| ❌     | titpetric/yamlexpr               | Expr.handleForWithContext        | 48.20%   | 42        |
| ❌     | titpetric/yamlexpr               | Expr.handleIncludeWithContext    | 0.00%    | 11        |
| ❌     | titpetric/yamlexpr               | Expr.handleMatrixWithContext     | 0.00%    | 56        |
| ❌     | titpetric/yamlexpr               | Expr.loadAndMergeFileWithContext | 0.00%    | 6         |
| ✅     | titpetric/yamlexpr               | Expr.process                     | 100.00%  | 4         |
| ❌     | titpetric/yamlexpr               | Expr.processMapWithContext       | 75.00%   | 15        |
| ❌     | titpetric/yamlexpr               | Expr.processSliceWithContext     | 54.80%   | 40        |
| ✅     | titpetric/yamlexpr               | Expr.processWithContext          | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr               | Expr.processWithStack            | 75.00%   | 1         |
| ✅     | titpetric/yamlexpr               | MapMatchesSpec                   | 100.00%  | 5         |
| ✅     | titpetric/yamlexpr               | New                              | 75.00%   | 1         |
| ✅     | titpetric/yamlexpr               | ValuesEqual                      | 100.00%  | 0         |
| ❌     | titpetric/yamlexpr               | applyExcludes                    | 0.00%    | 7         |
| ❌     | titpetric/yamlexpr               | applyIncludes                    | 0.00%    | 19        |
| ❌     | titpetric/yamlexpr               | evaluateConditionWithPath        | 54.50%   | 23        |
| ❌     | titpetric/yamlexpr               | expandMatrixBase                 | 0.00%    | 12        |
| ✅     | titpetric/yamlexpr               | isQuoted                         | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr               | isTruthy                         | 22.20%   | 1         |
| ✅     | titpetric/yamlexpr               | isValidVarName                   | 75.00%   | 1         |
| ❌     | titpetric/yamlexpr               | mergeRecursive                   | 0.00%    | 31        |
| ✅     | titpetric/yamlexpr               | parseForExpr                     | 95.20%   | 13        |
| ❌     | titpetric/yamlexpr               | parseMatrixDirective             | 0.00%    | 30        |
| ✅     | titpetric/yamlexpr               | parseYAML                        | 75.00%   | 1         |
| ✅     | titpetric/yamlexpr               | quoteUnquotedComparisons         | 84.60%   | 16        |
| ✅     | titpetric/yamlexpr               | valuesEqual                      | 92.90%   | 7         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | deepEqual                        | 0.00%    | 4         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | main                             | 0.00%    | 13        |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr  | testFixtures                     | 0.00%    | 63        |
| ✅     | titpetric/yamlexpr/interpolation | ContainsInterpolation            | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/interpolation | InterpolateString                | 81.50%   | 10        |
| ✅     | titpetric/yamlexpr/interpolation | InterpolateStringPermissive      | 90.00%   | 9         |
| ❌     | titpetric/yamlexpr/interpolation | InterpolateStringWithContext     | 64.10%   | 29        |
| ❌     | titpetric/yamlexpr/interpolation | InterpolateValue                 | 35.30%   | 21        |
| ✅     | titpetric/yamlexpr/interpolation | InterpolateValuePermissive       | 100.00%  | 3         |
| ❌     | titpetric/yamlexpr/interpolation | InterpolateValueWithContext      | 73.30%   | 11        |
| ✅     | titpetric/yamlexpr/interpolation | extractSingleExpression          | 80.00%   | 2         |
| ✅     | titpetric/yamlexpr/interpolation | isSingleInterpolation            | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/model         | Config.ForDirective              | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | Config.IfDirective               | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | Config.IncludeDirective          | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | Config.MatrixDirective           | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.AppendPath               | 81.80%   | 8         |
| ✅     | titpetric/yamlexpr/model         | Context.Count                    | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/model         | Context.FormatIncludeChain       | 0.00%    | 1         |
| ✅     | titpetric/yamlexpr/model         | Context.Path                     | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.Pop                      | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.Push                     | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.Stack                    | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.WithInclude              | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/model         | Context.WithPath                 | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | DefaultConfig                    | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model         | NewContext                       | 75.00%   | 3         |
| ❌     | titpetric/yamlexpr/model         | WithDirectiveHandler             | 0.00%    | 4         |
| ✅     | titpetric/yamlexpr/model         | WithFS                           | 0.00%    | 0         |
| ❌     | titpetric/yamlexpr/model         | WithSyntax                       | 0.00%    | 8         |
| ✅     | titpetric/yamlexpr/stack         | CanDescend                       | 100.00%  | 5         |
| ✅     | titpetric/yamlexpr/stack         | IsSlice                          | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack         | New                              | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/stack         | NewStack                         | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/stack         | NewStackWithData                 | 100.00%  | 1         |
| ❌     | titpetric/yamlexpr/stack         | PopulateStructFields             | 0.00%    | 17        |
| ✅     | titpetric/yamlexpr/stack         | ResolveValue                     | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack         | SliceToAny                       | 100.00%  | 4         |
| ✅     | titpetric/yamlexpr/stack         | Stack.All                        | 85.70%   | 4         |
| ✅     | titpetric/yamlexpr/stack         | Stack.Copy                       | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/stack         | Stack.Count                      | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/stack         | Stack.ForEach                    | 100.00%  | 12        |
| ✅     | titpetric/yamlexpr/stack         | Stack.GetInt                     | 73.30%   | 5         |
| ✅     | titpetric/yamlexpr/stack         | Stack.GetMap                     | 100.00%  | 5         |
| ✅     | titpetric/yamlexpr/stack         | Stack.GetSlice                   | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/stack         | Stack.GetString                  | 72.70%   | 3         |
| ❌     | titpetric/yamlexpr/stack         | Stack.Lookup                     | 71.40%   | 6         |
| ✅     | titpetric/yamlexpr/stack         | Stack.Pop                        | 90.90%   | 6         |
| ✅     | titpetric/yamlexpr/stack         | Stack.Push                       | 66.70%   | 1         |
| ✅     | titpetric/yamlexpr/stack         | Stack.Resolve                    | 92.30%   | 7         |
| ✅     | titpetric/yamlexpr/stack         | Stack.Set                        | 66.70%   | 1         |
| ✅     | titpetric/yamlexpr/stack         | Stack.resolveStep                | 100.00%  | 8         |
| ❌     | titpetric/yamlexpr/stack         | StructToMap                      | 0.00%    | 18        |
| ✅     | titpetric/yamlexpr/stack         | getCachedPath                    | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack         | resolveMap                       | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/stack         | resolveSliceIndex                | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack         | resolveStruct                    | 100.00%  | 6         |
| ✅     | titpetric/yamlexpr/stack         | resolveValueRecursive            | 88.90%   | 4         |
| ✅     | titpetric/yamlexpr/stack         | splitPathImpl                    | 85.40%   | 31        |
