# Testing Coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package                     | Coverage | Cognitive | Lines |
|--------|-----------------------------|----------|-----------|-------|
| ✅     | titpetric/yamlexpr          | 81.96%   | 107       | 476   |
| ❌     | titpetric/yamlexpr/handlers | 57.40%   | 308       | 990   |
| ✅     | titpetric/yamlexpr/model    | 100.00%  | 12        | 94    |
| ❌     | titpetric/yamlexpr/stack    | 78.07%   | 92        | 329   |

## Functions

| Status | Package                     | Function                         | Coverage | Cognitive |
|--------|-----------------------------|----------------------------------|----------|-----------|
| ✅     | titpetric/yamlexpr          | Config.EmbedDirective            | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr          | Config.ForDirective              | 100.00%  | 0         |
| ❌     | titpetric/yamlexpr          | Config.GetHandler                | 0.00%    | 1         |
| ✅     | titpetric/yamlexpr          | Config.IfDirective               | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr          | DefaultConfig                    | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr          | Expr.Load                        | 69.20%   | 4         |
| ✅     | titpetric/yamlexpr          | Expr.LoadAndMergeFileWithContext | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr          | Expr.Process                     | 100.00%  | 4         |
| ✅     | titpetric/yamlexpr          | Expr.ProcessMapWithContext       | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr          | Expr.ProcessWithContext          | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr          | Expr.ProcessWithStack            | 75.00%   | 1         |
| ✅     | titpetric/yamlexpr          | Expr.RegisterHandler             | 80.00%   | 2         |
| ✅     | titpetric/yamlexpr          | Expr.loadAndMergeFileWithContext | 83.30%   | 3         |
| ✅     | titpetric/yamlexpr          | Expr.processMapWithContext       | 92.60%   | 28        |
| ✅     | titpetric/yamlexpr          | Expr.processSliceWithContext     | 88.90%   | 5         |
| ✅     | titpetric/yamlexpr          | Expr.processWithContext          | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr          | New                              | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr          | NewExprContext                   | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr          | NewExtended                      | 100.00%  | 0         |
| ❌     | titpetric/yamlexpr          | WithDirectiveHandler             | 0.00%    | 4         |
| ✅     | titpetric/yamlexpr          | WithFS                           | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr          | WithStandardHandlers             | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr          | WithSyntax                       | 100.00%  | 6         |
| ✅     | titpetric/yamlexpr          | isValidVarName                   | 100.00%  | 1         |
| ❌     | titpetric/yamlexpr          | mergeRecursive                   | 53.80%   | 31        |
| ✅     | titpetric/yamlexpr          | parseForExpr                     | 95.20%   | 13        |
| ✅     | titpetric/yamlexpr          | parseYAML                        | 75.00%   | 1         |
| ✅     | titpetric/yamlexpr/handlers | BuildScope                       | 90.90%   | 9         |
| ✅     | titpetric/yamlexpr/handlers | ContainsInterpolation            | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/handlers | EmbedHandlerBuiltin              | 91.70%   | 19        |
| ❌     | titpetric/yamlexpr/handlers | EvaluateConditionWithPath        | 76.50%   | 23        |
| ❌     | titpetric/yamlexpr/handlers | ForHandlerBuiltin                | 56.10%   | 60        |
| ✅     | titpetric/yamlexpr/handlers | IfHandlerBuiltin                 | 85.70%   | 4         |
| ❌     | titpetric/yamlexpr/handlers | IfHandlerImpl.Handle             | 0.00%    | 2         |
| ✅     | titpetric/yamlexpr/handlers | InterpolateString                | 81.50%   | 10        |
| ✅     | titpetric/yamlexpr/handlers | InterpolateStringWithContext     | 80.80%   | 15        |
| ✅     | titpetric/yamlexpr/handlers | InterpolateValue                 | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/handlers | InterpolationHandlerBuiltin      | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/handlers | IsQuoted                         | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/handlers | IsTruthy                         | 77.80%   | 1         |
| ✅     | titpetric/yamlexpr/handlers | NewDiscardHandler                | 100.00%  | 11        |
| ✅     | titpetric/yamlexpr/handlers | NewForHandler                    | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/handlers | NewIfHandler                     | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/handlers | NewInterpolationHandler          | 100.00%  | 0         |
| ❌     | titpetric/yamlexpr/handlers | NewMatrixHandler                 | 0.00%    | 6         |
| ✅     | titpetric/yamlexpr/handlers | ParseForExpr                     | 90.50%   | 34        |
| ✅     | titpetric/yamlexpr/handlers | QuoteUnquotedComparisons         | 92.30%   | 16        |
| ❌     | titpetric/yamlexpr/handlers | applyEmbeds                      | 0.00%    | 19        |
| ❌     | titpetric/yamlexpr/handlers | applyExcludes                    | 0.00%    | 7         |
| ❌     | titpetric/yamlexpr/handlers | expandMatrixBase                 | 0.00%    | 12        |
| ✅     | titpetric/yamlexpr/handlers | isNumeric                        | 83.30%   | 7         |
| ❌     | titpetric/yamlexpr/handlers | jobMatchesSpec                   | 0.00%    | 5         |
| ❌     | titpetric/yamlexpr/handlers | parseMatrixDirective             | 0.00%    | 30        |
| ✅     | titpetric/yamlexpr/handlers | trimSpace                        | 100.00%  | 4         |
| ❌     | titpetric/yamlexpr/handlers | valuesEqual                      | 0.00%    | 7         |
| ✅     | titpetric/yamlexpr/model    | Context.AppendPath               | 100.00%  | 8         |
| ✅     | titpetric/yamlexpr/model    | Context.FormatIncludeChain       | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/model    | Context.Path                     | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model    | Context.PopStackScope            | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model    | Context.PushStackScope           | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model    | Context.Stack                    | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model    | Context.WithInclude              | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model    | Context.WithPath                 | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model    | NewContext                       | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/stack    | New                              | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/stack    | Stack.All                        | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/stack    | Stack.Copy                       | 100.00%  | 0         |
| ❌     | titpetric/yamlexpr/stack    | Stack.ForEach                    | 46.70%   | 12        |
| ✅     | titpetric/yamlexpr/stack    | Stack.GetInt                     | 46.70%   | 5         |
| ✅     | titpetric/yamlexpr/stack    | Stack.GetMap                     | 40.00%   | 5         |
| ✅     | titpetric/yamlexpr/stack    | Stack.GetSlice                   | 80.00%   | 5         |
| ✅     | titpetric/yamlexpr/stack    | Stack.GetString                  | 72.70%   | 3         |
| ✅     | titpetric/yamlexpr/stack    | Stack.Lookup                     | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/stack    | Stack.Pop                        | 81.80%   | 6         |
| ✅     | titpetric/yamlexpr/stack    | Stack.Push                       | 66.70%   | 1         |
| ✅     | titpetric/yamlexpr/stack    | Stack.Resolve                    | 84.60%   | 7         |
| ✅     | titpetric/yamlexpr/stack    | Stack.Set                        | 66.70%   | 1         |
| ❌     | titpetric/yamlexpr/stack    | Stack.resolveStep                | 77.80%   | 7         |
| ✅     | titpetric/yamlexpr/stack    | getCachedPath                    | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack    | splitPathImpl                    | 85.40%   | 31        |
