# Testing Coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package                  | Coverage | Cognitive | Lines |
|--------|--------------------------|----------|-----------|-------|
| ✅     | titpetric/yamlexpr       | 83.92%   | 239       | 829   |
| ❌     | titpetric/yamlexpr/stack | 78.07%   | 92        | 329   |

## Functions

| Status | Package                  | Function                         | Coverage | Cognitive |
|--------|--------------------------|----------------------------------|----------|-----------|
| ✅     | titpetric/yamlexpr       | Config.ForDirective              | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | Config.IfDirective               | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | Config.IncludeDirective          | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | DefaultConfig                    | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | Expr.Load                        | 69.20%   | 4         |
| ✅     | titpetric/yamlexpr       | Expr.Process                     | 100.00%  | 4         |
| ✅     | titpetric/yamlexpr       | Expr.ProcessWithStack            | 75.00%   | 1         |
| ❌     | titpetric/yamlexpr       | Expr.handleForWithContext        | 60.70%   | 42        |
| ❌     | titpetric/yamlexpr       | Expr.handleIncludeWithContext    | 22.20%   | 11        |
| ✅     | titpetric/yamlexpr       | Expr.loadAndMergeFileWithContext | 75.00%   | 3         |
| ✅     | titpetric/yamlexpr       | Expr.processMapWithContext       | 86.40%   | 14        |
| ❌     | titpetric/yamlexpr       | Expr.processSliceWithContext     | 66.70%   | 29        |
| ✅     | titpetric/yamlexpr       | Expr.processWithContext          | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr       | ExprContext.AppendPath           | 100.00%  | 8         |
| ✅     | titpetric/yamlexpr       | ExprContext.FormatIncludeChain   | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr       | ExprContext.Path                 | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | ExprContext.PopStackScope        | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | ExprContext.PushStackScope       | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | ExprContext.Stack                | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | ExprContext.WithInclude          | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | ExprContext.WithPath             | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr       | New                              | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr       | NewExprContext                   | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr       | WithSyntax                       | 100.00%  | 6         |
| ✅     | titpetric/yamlexpr       | containsInterpolation            | 100.00%  | 1         |
| ❌     | titpetric/yamlexpr       | evaluateConditionWithPath        | 54.50%   | 23        |
| ✅     | titpetric/yamlexpr       | interpolateString                | 88.90%   | 5         |
| ✅     | titpetric/yamlexpr       | interpolateStringHelper          | 66.70%   | 1         |
| ❌     | titpetric/yamlexpr       | interpolateStringWithContext     | 65.40%   | 15        |
| ✅     | titpetric/yamlexpr       | isQuoted                         | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr       | isTruthy                         | 22.20%   | 1         |
| ✅     | titpetric/yamlexpr       | isValidVarName                   | 75.00%   | 1         |
| ❌     | titpetric/yamlexpr       | mergeRecursive                   | 38.50%   | 31        |
| ✅     | titpetric/yamlexpr       | parseForExpr                     | 95.20%   | 13        |
| ✅     | titpetric/yamlexpr       | parseYAML                        | 75.00%   | 1         |
| ✅     | titpetric/yamlexpr       | quoteUnquotedComparisons         | 84.60%   | 16        |
| ✅     | titpetric/yamlexpr/stack | New                              | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/stack | Stack.All                        | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/stack | Stack.Copy                       | 100.00%  | 0         |
| ❌     | titpetric/yamlexpr/stack | Stack.ForEach                    | 46.70%   | 12        |
| ✅     | titpetric/yamlexpr/stack | Stack.GetInt                     | 46.70%   | 5         |
| ✅     | titpetric/yamlexpr/stack | Stack.GetMap                     | 40.00%   | 5         |
| ✅     | titpetric/yamlexpr/stack | Stack.GetSlice                   | 80.00%   | 5         |
| ✅     | titpetric/yamlexpr/stack | Stack.GetString                  | 72.70%   | 3         |
| ✅     | titpetric/yamlexpr/stack | Stack.Lookup                     | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/stack | Stack.Pop                        | 81.80%   | 6         |
| ✅     | titpetric/yamlexpr/stack | Stack.Push                       | 66.70%   | 1         |
| ✅     | titpetric/yamlexpr/stack | Stack.Resolve                    | 84.60%   | 7         |
| ✅     | titpetric/yamlexpr/stack | Stack.Set                        | 66.70%   | 1         |
| ❌     | titpetric/yamlexpr/stack | Stack.resolveStep                | 77.80%   | 7         |
| ✅     | titpetric/yamlexpr/stack | getCachedPath                    | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack | splitPathImpl                    | 85.40%   | 31        |
