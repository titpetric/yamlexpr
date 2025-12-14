# Testing Coverage

Testing criteria for a passing coverage requirement:

- Line coverage of 80%
- Cognitive complexity of 0
- Have cognitive complexity < 5, but have any coverage

Low cognitive complexity means there are few conditional branches to cover. Tests with cognitive complexity 0 would be covered by invocation.

## Packages

| Status | Package                         | Coverage | Cognitive | Lines |
|--------|---------------------------------|----------|-----------|-------|
| ✅     | titpetric/yamlexpr              | 83.52%   | 121       | 488   |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr | 0.00%    | 34        | 162   |
| ❌     | titpetric/yamlexpr/handlers     | 57.35%   | 427       | 1312  |
| ✅     | titpetric/yamlexpr/merge        | 92.88%   | 65        | 285   |
| ✅     | titpetric/yamlexpr/model        | 100.00%  | 12        | 94    |
| ✅     | titpetric/yamlexpr/stack        | 81.93%   | 156       | 578   |

## Functions

| Status | Package                         | Function                         | Coverage | Cognitive |
|--------|---------------------------------|----------------------------------|----------|-----------|
| ✅     | titpetric/yamlexpr              | Config.ForDirective              | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr              | Config.IfDirective               | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr              | Config.IncludeDirective          | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr              | DefaultConfig                    | 100.00%  | 0         |
| ❌     | titpetric/yamlexpr              | Expr.Load                        | 66.70%   | 9         |
| ✅     | titpetric/yamlexpr              | Expr.LoadAndMergeFileWithContext | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr              | Expr.Process                     | 100.00%  | 4         |
| ✅     | titpetric/yamlexpr              | Expr.ProcessMapWithContext       | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr              | Expr.ProcessWithContext          | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr              | Expr.ProcessWithStack            | 75.00%   | 1         |
| ✅     | titpetric/yamlexpr              | Expr.RegisterHandler             | 80.00%   | 2         |
| ✅     | titpetric/yamlexpr              | Expr.expandRootLevel             | 87.50%   | 2         |
| ✅     | titpetric/yamlexpr              | Expr.loadAndMergeFileWithContext | 83.30%   | 3         |
| ✅     | titpetric/yamlexpr              | Expr.processMapWithContext       | 92.60%   | 28        |
| ✅     | titpetric/yamlexpr              | Expr.processSliceWithContext     | 90.90%   | 10        |
| ✅     | titpetric/yamlexpr              | Expr.processWithContext          | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr              | New                              | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr              | NewExprContext                   | 100.00%  | 0         |
| ❌     | titpetric/yamlexpr              | WithDirectiveHandler             | 0.00%    | 4         |
| ✅     | titpetric/yamlexpr              | WithFS                           | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr              | WithSyntax                       | 100.00%  | 6         |
| ✅     | titpetric/yamlexpr              | isMapWithHandler                 | 71.40%   | 4         |
| ✅     | titpetric/yamlexpr              | isValidVarName                   | 100.00%  | 1         |
| ❌     | titpetric/yamlexpr              | mergeRecursive                   | 53.80%   | 31        |
| ✅     | titpetric/yamlexpr              | parseForExpr                     | 95.20%   | 13        |
| ✅     | titpetric/yamlexpr              | parseYAML                        | 75.00%   | 1         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr | deepEqual                        | 0.00%    | 2         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr | main                             | 0.00%    | 9         |
| ❌     | titpetric/yamlexpr/cmd/yamlexpr | testFixtures                     | 0.00%    | 23        |
| ✅     | titpetric/yamlexpr/handlers     | BuildScope                       | 90.90%   | 9         |
| ✅     | titpetric/yamlexpr/handlers     | ContainsInterpolation            | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/handlers     | DiscardHandlerBuiltin            | 100.00%  | 11        |
| ❌     | titpetric/yamlexpr/handlers     | EvaluateConditionWithPath        | 79.40%   | 23        |
| ✅     | titpetric/yamlexpr/handlers     | ExpandForAtRoot                  | 84.40%   | 35        |
| ❌     | titpetric/yamlexpr/handlers     | ExpandMatrixAtRoot               | 0.00%    | 22        |
| ❌     | titpetric/yamlexpr/handlers     | ForHandler                       | 66.70%   | 60        |
| ✅     | titpetric/yamlexpr/handlers     | IfHandler                        | 85.70%   | 4         |
| ❌     | titpetric/yamlexpr/handlers     | IfHandlerImpl.Handle             | 0.00%    | 2         |
| ✅     | titpetric/yamlexpr/handlers     | IncludeHandler                   | 91.70%   | 19        |
| ✅     | titpetric/yamlexpr/handlers     | InterpolateString                | 81.50%   | 10        |
| ✅     | titpetric/yamlexpr/handlers     | InterpolateStringPermissive      | 90.00%   | 9         |
| ❌     | titpetric/yamlexpr/handlers     | InterpolateStringWithContext     | 66.70%   | 29        |
| ❌     | titpetric/yamlexpr/handlers     | InterpolateValue                 | 35.30%   | 21        |
| ✅     | titpetric/yamlexpr/handlers     | InterpolateValuePermissive       | 100.00%  | 3         |
| ❌     | titpetric/yamlexpr/handlers     | InterpolateValueWithContext      | 73.30%   | 11        |
| ✅     | titpetric/yamlexpr/handlers     | InterpolationHandlerBuiltin      | 50.00%   | 0         |
| ✅     | titpetric/yamlexpr/handlers     | IsQuoted                         | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/handlers     | IsTruthy                         | 77.80%   | 1         |
| ✅     | titpetric/yamlexpr/handlers     | NewDiscardHandler                | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/handlers     | NewForHandler                    | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/handlers     | NewIfHandler                     | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/handlers     | NewInterpolationHandler          | 100.00%  | 0         |
| ❌     | titpetric/yamlexpr/handlers     | NewMatrixHandler                 | 0.00%    | 6         |
| ✅     | titpetric/yamlexpr/handlers     | ParseForExpr                     | 92.90%   | 34        |
| ✅     | titpetric/yamlexpr/handlers     | QuoteUnquotedComparisons         | 92.30%   | 16        |
| ❌     | titpetric/yamlexpr/handlers     | applyEmbeds                      | 0.00%    | 19        |
| ❌     | titpetric/yamlexpr/handlers     | applyExcludes                    | 0.00%    | 7         |
| ❌     | titpetric/yamlexpr/handlers     | expandMatrixBase                 | 0.00%    | 12        |
| ✅     | titpetric/yamlexpr/handlers     | extractSingleExpression          | 80.00%   | 2         |
| ❌     | titpetric/yamlexpr/handlers     | isComplexExpression              | 0.00%    | 3         |
| ✅     | titpetric/yamlexpr/handlers     | isNumeric                        | 83.30%   | 7         |
| ✅     | titpetric/yamlexpr/handlers     | isSingleInterpolation            | 100.00%  | 2         |
| ❌     | titpetric/yamlexpr/handlers     | jobMatchesSpec                   | 0.00%    | 5         |
| ❌     | titpetric/yamlexpr/handlers     | parseMatrixDirective             | 0.00%    | 30        |
| ✅     | titpetric/yamlexpr/handlers     | trimSpace                        | 100.00%  | 4         |
| ❌     | titpetric/yamlexpr/handlers     | valuesEqual                      | 0.00%    | 7         |
| ✅     | titpetric/yamlexpr/merge        | Deduplicate                      | 100.00%  | 1         |
| ❌     | titpetric/yamlexpr/merge        | DeduplicateWithStats             | 0.00%    | 1         |
| ✅     | titpetric/yamlexpr/merge        | MergeMap.Data                    | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/merge        | MergeMap.Distinct                | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/merge        | MergeMap.Merge                   | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/merge        | MergeMap.Stats                   | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/merge        | MergeMap.buildPath               | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/merge        | MergeMap.deepCopyMap             | 100.00%  | 5         |
| ✅     | titpetric/yamlexpr/merge        | MergeMap.mergeRecursive          | 100.00%  | 13        |
| ✅     | titpetric/yamlexpr/merge        | MergeMap.recordKey               | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/merge        | MergeMap.recordKeysRecursive     | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/merge        | NewMergeMap                      | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/merge        | appendIfDistinct                 | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/merge        | areEqual                         | 100.00%  | 5         |
| ✅     | titpetric/yamlexpr/merge        | buildPath                        | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/merge        | deduplicateMapNoStats            | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/merge        | deduplicateSliceWithStats        | 100.00%  | 12        |
| ✅     | titpetric/yamlexpr/merge        | deduplicateValueNoStats          | 60.00%   | 1         |
| ✅     | titpetric/yamlexpr/merge        | deduplicateValueWithStats        | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/merge        | deduplicateWithStats             | 100.00%  | 9         |
| ✅     | titpetric/yamlexpr/merge        | formatMapValue                   | 83.30%   | 4         |
| ✅     | titpetric/yamlexpr/merge        | generateComparisonKey            | 100.00%  | 4         |
| ✅     | titpetric/yamlexpr/model        | Context.AppendPath               | 100.00%  | 8         |
| ✅     | titpetric/yamlexpr/model        | Context.FormatIncludeChain       | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/model        | Context.Path                     | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model        | Context.PopStackScope            | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model        | Context.PushStackScope           | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model        | Context.Stack                    | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model        | Context.WithInclude              | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model        | Context.WithPath                 | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/model        | NewContext                       | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/stack        | CanDescend                       | 100.00%  | 5         |
| ✅     | titpetric/yamlexpr/stack        | IsSlice                          | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack        | New                              | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/stack        | NewStack                         | 100.00%  | 0         |
| ✅     | titpetric/yamlexpr/stack        | NewStackWithData                 | 100.00%  | 1         |
| ❌     | titpetric/yamlexpr/stack        | PopulateStructFields             | 0.00%    | 17        |
| ✅     | titpetric/yamlexpr/stack        | ResolveValue                     | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack        | SliceToAny                       | 100.00%  | 4         |
| ✅     | titpetric/yamlexpr/stack        | Stack.All                        | 85.70%   | 4         |
| ✅     | titpetric/yamlexpr/stack        | Stack.Copy                       | 0.00%    | 0         |
| ✅     | titpetric/yamlexpr/stack        | Stack.ForEach                    | 100.00%  | 12        |
| ✅     | titpetric/yamlexpr/stack        | Stack.GetInt                     | 73.30%   | 5         |
| ✅     | titpetric/yamlexpr/stack        | Stack.GetMap                     | 100.00%  | 5         |
| ✅     | titpetric/yamlexpr/stack        | Stack.GetSlice                   | 100.00%  | 3         |
| ✅     | titpetric/yamlexpr/stack        | Stack.GetString                  | 72.70%   | 3         |
| ❌     | titpetric/yamlexpr/stack        | Stack.Lookup                     | 71.40%   | 6         |
| ✅     | titpetric/yamlexpr/stack        | Stack.Pop                        | 90.90%   | 6         |
| ✅     | titpetric/yamlexpr/stack        | Stack.Push                       | 66.70%   | 1         |
| ✅     | titpetric/yamlexpr/stack        | Stack.Resolve                    | 92.30%   | 7         |
| ✅     | titpetric/yamlexpr/stack        | Stack.Set                        | 66.70%   | 1         |
| ✅     | titpetric/yamlexpr/stack        | Stack.resolveStep                | 100.00%  | 8         |
| ❌     | titpetric/yamlexpr/stack        | StructToMap                      | 0.00%    | 18        |
| ✅     | titpetric/yamlexpr/stack        | getCachedPath                    | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack        | resolveMap                       | 100.00%  | 1         |
| ✅     | titpetric/yamlexpr/stack        | resolveSliceIndex                | 100.00%  | 2         |
| ✅     | titpetric/yamlexpr/stack        | resolveStruct                    | 100.00%  | 6         |
| ✅     | titpetric/yamlexpr/stack        | resolveValueRecursive            | 88.90%   | 4         |
| ✅     | titpetric/yamlexpr/stack        | splitPathImpl                    | 85.40%   | 31        |
