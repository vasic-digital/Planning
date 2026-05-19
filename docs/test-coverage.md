# Test-Coverage Ledger — round-270

This ledger maps every exported symbol of `digital.vasic.planning`
to the test or Challenge that exercises it with captured runtime
evidence. Per CONST-035, CONST-050(B), and the 2026-05-19 operator
mandate quoted below, no symbol may PASS without a corresponding
runtime-evidence exercise.

> Verbatim 2026-05-19 operator mandate: "all existing tests and
> Challenges do work in anti-bluff manner - they MUST confirm that
> all tested codebase really works as expected! We had been in
> position that all tests do execute with success and all
> Challenges as well, but in reality the most of the features does
> not work and can't be used! This MUST NOT be the case and
> execution of tests and Challenges MUST guarantee the quality, the
> completition and full usability by end users of the product!"

Operative rule (Article XI §11.9): **The bar for shipping is not
"tests pass" but "users can use the feature."** Every PASS in the
table below carries either a unit test, an integration test, an
e2e test, or a challenge-runner section that produces positive
runtime evidence — no metadata-only / grep-only PASS counts.

## Module surface

`digital.vasic.planning` ships TWO Go packages:

- **`planning`** — root domain package: `hiplan.go` + `mcts.go` +
  `tree_of_thoughts.go` + `i18n_defaults.go` (i18n string-ID
  fallbacks for CONST-046).
- **`pkg/i18n`** — minimal, dependency-free Translator interface
  + `NoopTranslator` standalone default; per CONST-051(B) this
  submodule depends on no consuming-project package.

## Symbol → exerciser map

### `planning/hiplan.go`

| Symbol | Kind | Exercised by |
|--------|------|--------------|
| `MilestoneState` | type | runner Section 1 (state-field round-trip per locale) + `hiplan_test.go` |
| `MilestoneStatePending` | const | runner Section 1 (default-state assertion) + `hiplan_test.go` |
| `MilestoneStateInProgress` | const | `hiplan_test.go` (mid-execution transition) |
| `MilestoneStateCompleted` | const | `hiplan_test.go` |
| `MilestoneStateFailed` | const | `hiplan_test.go` |
| `MilestoneStateSkipped` | const | `hiplan_test.go` |
| `Milestone` | struct | runner Section 1 (per-locale Name+Description byte-exact + Library round-trip) + `hiplan_test.go` |
| `PlanStepState` | type | runner Section 1 (default-state on generated steps) + `hiplan_test.go` |
| `PlanStepStatePending` | const | runner Section 1 + `hiplan_test.go` |
| `PlanStepStateInProgress` | const | `hiplan_test.go` |
| `PlanStepStateCompleted` | const | `hiplan_test.go` |
| `PlanStepStateFailed` | const | `hiplan_test.go` |
| `PlanStep` | struct | runner Section 1 (per-locale Description+Action+Hints byte-exact) + `hiplan_test.go` |
| `HiPlanConfig` | struct | runner Section 1 (custom config built per Section) + `hiplan_test.go` |
| `DefaultHiPlanConfig` | func | runner Section 1 (used as base) + `hiplan_test.go` (full-field assertion) |
| `MilestoneGenerator` | interface | runner Section 1 (`echoingMilestoneGen` satisfies all 3 methods) |
| `StepExecutor` | interface | runner Section 1 (`echoingExecutor` satisfies Execute+Validate) |
| `StepResult` | struct | runner Section 1 (Outputs map round-trip with locale bytes) |
| `HierarchicalPlan` | struct | runner Section 1 (per-locale Goal byte-exact, Milestones len, State, Progress) |
| `HiPlan` | struct | runner Section 1 (real `NewHiPlan` per locale) |
| `NewHiPlan` | func | runner Section 1 (5 locales, custom + default configs) + `hiplan_test.go` |
| `HiPlan.SetTranslator` | method | runner Section 4 (nil-input branch — restores Default) + `i18n_callsites_test.go` |
| `HiPlan.CreatePlan` | method | runner Section 1 (per-locale Goal byte-exact assertion) + `hiplan_test.go` |
| `HiPlan.ExecutePlan` | method | runner Section 1 (per-locale success=true, completed=1) + `hiplan_test.go` |
| `HiPlan.ExecuteStep` | method | runner Section 1 (per-locale Outputs.description byte-exact) + `hiplan_test.go` |
| `HiPlan.AddToLibrary` | method | runner Section 1 (per-locale add) + `hiplan_test.go` |
| `HiPlan.GetFromLibrary` | method | runner Section 1 (per-locale retrieve + Description byte-exact) + `hiplan_test.go` |
| `HiPlan.GetCurrentPlan` | method | runner Section 1 (ID equality vs just-created plan) + `hiplan_test.go` |
| `MilestoneResult` | struct | runner Section 1 (indirect via ExecutePlan) |
| `PlanResult` | struct | runner Section 1 (per-locale success+completed+failed+duration) + `hiplan_test.go` |
| `PlanResult.MarshalJSON` | method | `hiplan_test.go` (JSON round-trip) |
| `LLMMilestoneGenerator` | struct | runner Section 3 (NewLLMMilestoneGenerator constructor) + `hiplan_test.go` |
| `NewLLMMilestoneGenerator` | func | runner Section 3 (constructor) + `hiplan_test.go` |
| `LLMMilestoneGenerator.SetTranslator` | method | `i18n_callsites_test.go` |
| `LLMMilestoneGenerator.GenerateMilestones` | method | `hiplan_test.go` |
| `LLMMilestoneGenerator.GenerateSteps` | method | `hiplan_test.go` |
| `LLMMilestoneGenerator.GenerateHints` | method | `hiplan_test.go` |

### `planning/mcts.go`

| Symbol | Kind | Exercised by |
|--------|------|--------------|
| `MCTSNodeState` | type | `mcts_test.go` |
| `MCTSNodeStateUnexpanded` | const | `mcts_test.go` (default state) |
| `MCTSNodeStateExpanded` | const | `mcts_test.go` |
| `MCTSNodeStateTerminal` | const | `mcts_test.go` |
| `MCTSNode` | struct | runner Section 2 (UCTValue + AverageReward + AddReward) + `mcts_test.go` |
| `MCTSNode.AverageReward` | method | runner Section 2 (per-locale avg=0.5 after AddReward(0.5)) |
| `MCTSNode.AddReward` | method | runner Section 2 (per-locale visits=6 after add) |
| `MCTSConfig` | struct | runner Section 2 (custom budgets per Section) + `mcts_test.go` |
| `DefaultMCTSConfig` | func | runner Section 2 + `mcts_test.go` |
| `MCTSActionGenerator` | interface | runner Section 2 (`deterministicActionGen` satisfies GetActions+ApplyAction) |
| `MCTSRewardFunction` | interface | runner Section 2 (`deterministicReward` satisfies Evaluate+IsTerminal) |
| `MCTSRolloutPolicy` | interface | runner Section 2 (NewDefaultRolloutPolicy satisfies Rollout) |
| `MCTS` | struct | runner Section 2 (real `NewMCTS` per locale) |
| `NewMCTS` | func | runner Section 2 (5 locales) + `mcts_test.go` |
| `MCTS.Search` | method | runner Section 2 (per-locale TreeSize>0, FinalReward in [0,1], best-action carries seed bytes) + `mcts_test.go` |
| `MCTS.UCTValue` | method | runner Section 2 (per-locale visits=5 parentVisits=100 → 1.857) + `mcts_test.go` |
| `MCTSResult` | struct | runner Section 2 (per-locale full field walk) + `mcts_test.go` |
| `MCTSResult.MarshalJSON` | method | `mcts_test.go` |
| `CodeActionGenerator` | struct | runner Section 2 (NewCodeActionGenerator constructor) + `mcts_test.go` |
| `NewCodeActionGenerator` | func | runner Section 2 (constructor) + `mcts_test.go` |
| `CodeActionGenerator.SetTranslator` | method | `i18n_callsites_test.go` |
| `CodeActionGenerator.GetActions` | method | `mcts_test.go` |
| `CodeActionGenerator.ApplyAction` | method | `mcts_test.go` |
| `CodeRewardFunction` | struct | runner Section 2 (NewCodeRewardFunction constructor) + `mcts_test.go` |
| `NewCodeRewardFunction` | func | runner Section 2 (constructor with stub funcs) + `mcts_test.go` |
| `CodeRewardFunction.Evaluate` | method | `mcts_test.go` |
| `CodeRewardFunction.IsTerminal` | method | `mcts_test.go` |
| `DefaultRolloutPolicy` | struct | runner Section 2 (NewDefaultRolloutPolicy) + `mcts_test.go` |
| `NewDefaultRolloutPolicy` | func | runner Section 2 (constructor with our deterministic strategies) + `mcts_test.go` |
| `DefaultRolloutPolicy.Rollout` | method | runner Section 2 (indirect via Search rollouts) + `mcts_test.go` |

### `planning/tree_of_thoughts.go`

| Symbol | Kind | Exercised by |
|--------|------|--------------|
| `ThoughtState` | type | `tree_of_thoughts_test.go` |
| `ThoughtStatePending` | const | runner Section 3 (default state on generated thoughts) + `tree_of_thoughts_test.go` |
| `ThoughtStateActive` | const | runner Section 3 (root state) + `tree_of_thoughts_test.go` |
| `ThoughtStateEvaluated` | const | `tree_of_thoughts_test.go` |
| `ThoughtStatePruned` | const | `tree_of_thoughts_test.go` |
| `ThoughtStateSelected` | const | `tree_of_thoughts_test.go` |
| `Thought` | struct | runner Section 3 (per-locale Content carries ToTThought bytes) + `tree_of_thoughts_test.go` |
| `ThoughtNode` | struct | runner Section 3 (indirect via Solve) + `tree_of_thoughts_test.go` |
| `TreeOfThoughtsConfig` | struct | runner Section 3 (custom MaxDepth/MaxBranches/BeamWidth per strategy) + `tree_of_thoughts_test.go` |
| `DefaultTreeOfThoughtsConfig` | func | runner Section 3 (base config) + `tree_of_thoughts_test.go` |
| `ThoughtGenerator` | interface | runner Section 3 (`echoingThoughtGen` satisfies both methods) |
| `ThoughtEvaluator` | interface | runner Section 3 (`echoingThoughtEval` satisfies all 3 methods) |
| `TreeOfThoughts` | struct | runner Section 3 (real `NewTreeOfThoughts` per locale × per strategy = 15 instances) |
| `NewTreeOfThoughts` | func | runner Section 3 (15 instances) + `tree_of_thoughts_test.go` |
| `TreeOfThoughts.Solve` | method | runner Section 3 (5 locales × 3 strategies = 15 invocations, NodesExplored>0, BestScore in [0,1] or sentinel -1, locale ToTThought bytes preserved in selected path when non-empty) + `tree_of_thoughts_test.go` |
| `ToTResult` | struct | runner Section 3 (per-locale per-strategy NodesExplored / Iterations / BestScore / Solution / TreeDepth / Strategy / Problem) + `tree_of_thoughts_test.go` |
| `ToTResult.GetSolutionContent` | method | runner Section 3 (per-strategy locale-bytes check on selected path) + `tree_of_thoughts_test.go` |
| `ToTResult.MarshalJSON` | method | `tree_of_thoughts_test.go` |
| `LLMThoughtGenerator` | struct | runner Section 3 (NewLLMThoughtGenerator constructor) + `tree_of_thoughts_test.go` |
| `NewLLMThoughtGenerator` | func | runner Section 3 (constructor) + `tree_of_thoughts_test.go` |
| `LLMThoughtGenerator.SetTranslator` | method | `i18n_callsites_test.go` |
| `LLMThoughtGenerator.GenerateThoughts` | method | `tree_of_thoughts_test.go` |
| `LLMThoughtGenerator.GenerateInitialThoughts` | method | `tree_of_thoughts_test.go` |
| `LLMThoughtEvaluator` | struct | runner Section 3 (NewLLMThoughtEvaluator constructor) + `tree_of_thoughts_test.go` |
| `NewLLMThoughtEvaluator` | func | runner Section 3 (constructor) + `tree_of_thoughts_test.go` |
| `LLMThoughtEvaluator.SetTranslator` | method | `i18n_callsites_test.go` |
| `LLMThoughtEvaluator.EvaluateThought` | method | `tree_of_thoughts_test.go` |
| `LLMThoughtEvaluator.EvaluatePath` | method | `tree_of_thoughts_test.go` |
| `LLMThoughtEvaluator.IsTerminal` | method | `tree_of_thoughts_test.go` |

### `pkg/i18n/translator.go`

| Symbol | Kind | Exercised by |
|--------|------|--------------|
| `Translator` | interface | runner Section 4 (Default returns conforming impl) + `translator_test.go` |
| `NoopTranslator` | struct | runner Section 4 (T msg-id fallthrough) + `translator_test.go` |
| `NoopTranslator.T` | method | runner Section 4 (T returns msgID verbatim) + `translator_test.go` |
| `Default` | func | runner Section 4 (returns non-nil Translator) + `translator_test.go` |

## Test runs (round-270 evidence captured)

### `go test -count=1 -p 1 ./...`

```
ok  	digital.vasic.planning/pkg/i18n     (~0.001s)
ok  	digital.vasic.planning/planning     (~0.009s)
ok  	digital.vasic.planning/tests/benchmark   [no tests to run]
ok  	digital.vasic.planning/tests/e2e    (~0.003s)
ok  	digital.vasic.planning/tests/integration (~0.001s)
ok  	digital.vasic.planning/tests/security    (~0.002s)
ok  	digital.vasic.planning/tests/stress      (~0.001s)
```

The full suite passes without `-race`. There is a single pre-existing
race condition in `TestHiPlan_CreatePlan_ConcurrentAccess`
(`planning/hiplan_test.go:1321`) dating back to the initial commit
(`27a207e`); it is tracked separately and is NOT introduced by
round-270. Round-270 deliverables — runner, fixture, ledger, challenge
wrapper — pass under both `-race` and non-race runs because they only
touch the public API of `HiPlan`/`MCTS`/`TreeOfThoughts` and do not
exercise the concurrent-CreatePlan code path under the race detector.

### `challenges/runner/main.go -fixtures tests/fixtures/planning/payloads.json`

```
=== Round-270 Planning Challenge Runner ===
... PASS lines across 4 sections, 5 locales × 3 ToT strategies ...
=== Summary: 64 PASS, 0 FAIL ===
```

Per-locale runtime evidence captured:

- Section 1 (HiPlan): 5 × (CreatePlan + GetCurrentPlan + ExecutePlan
  + ExecuteStep + Library round-trip) = 25 PASS. Every PASS carries
  rune count of the Goal, milestone/step counts, and a byte-equality
  check for `Goal`, `Milestone.Name`, `Step.Description`,
  `StepResult.Outputs["description"]`, and `Milestone.Description`
  via the library retrieval.
- Section 2 (MCTS): 5 × (Search + UCTValue + AverageReward) = 15
  PASS, plus 3 constructor PASS for
  `NewCodeActionGenerator` + `NewCodeRewardFunction` +
  `NewDefaultRolloutPolicy`. Every Search PASS carries TreeSize,
  best-action count, FinalReward, and rune count of the MCTS state.
- Section 3 (Tree-of-Thoughts): 5 locales × 3 strategies (bfs, dfs,
  beam) = 15 Solve PASS, plus 3 constructor PASS for
  `NewLLMThoughtGenerator` + `NewLLMThoughtEvaluator` +
  `NewLLMMilestoneGenerator`. Every Solve PASS carries
  NodesExplored, Iterations, selected-path length, BestScore (or
  sentinel -1 for non-converging strategies on shallow budgets), and
  rune count of the problem string.
- Section 4 (i18n fallback): 3 PASS — Default returns non-nil
  Translator; NoopTranslator.T returns msgID verbatim;
  HiPlan.SetTranslator(nil) survives the nil branch and restores
  Default.

### `bash challenges/scripts/planning_describe_challenge.sh`

Clean mode exit 0; `--anti-bluff-mutate` exit 99 (paired mutation
correctly detected — ledger-vs-source drift caught when
`CreatePlan` is renamed to `CreatePlan_MUTATED` in a tmp copy of
the ledger).

## Anti-bluff invariants

This round addresses every taxonomy entry in CLAUDE.md §"Bluff
taxonomy":

- **Wrapper bluff** — the describe-challenge wrapper uses
  PASS/FAIL counters with a separate guard, never inline arithmetic
  on a command that prints + exits non-zero.
- **Contract bluff** — every public method, type, constant, and
  constructor listed above is exercised by a runtime test or
  challenge section. The ledger surface is closed and audited. For
  ToT, all three advertised search strategies (bfs, dfs, beam) are
  independently invoked per locale.
- **Structural bluff** — no `check_file_exists` PASS without a
  paired functional assertion. Every PASS carries either a rune
  count, a byte-equality check, a boolean expected/got, or a
  numeric-field comparison.
- **Comment bluff** — the README's `## Anti-bluff guarantees`
  section is enforced by `planning_describe_challenge.sh` Section 5.
- **Skip bluff** — the runner has no dead branches and never calls
  `t.Skip`.

## Cross-reference to constitutional anchors

| Anchor | Layer | How honoured |
|--------|-------|--------------|
| CONST-035 / Article XI §11.9 | end-user-usability | every PASS line carries runtime evidence (locale, rune count, boolean, byte equality) |
| CONST-046 | no-hardcoded-content | runner Section 4 exercises the i18n.Translator seam; ledger entries for `SetTranslator` per planner type |
| CONST-050(A) | no-fakes-beyond-unit-tests | runner uses only public types — no test helpers from `planning` package's internal layer; the only fakes are the `*_test.go` files for the unit suite, plus the runner's strategy-injection seams which are the *advertised public extension point*, not internal mocks |
| CONST-050(B) | 100%-test-type coverage | unit tests + challenge runner + paired-mutation gate together cover unit + integration-style + meta-test layers; plus pre-existing `tests/{e2e,integration,security,stress,benchmark}` directories |
| CONST-051(B) | decoupling | `pkg/i18n` declares its own Translator interface — no import of any consuming-project package |
| CONST-053 | .gitignore | `.gitignore` covers `/bin/`, `*.test`, `coverage.out`, IDE state, build artefacts, `.env`, secrets, language caches |

The 2026-05-19 operator mandate is preserved verbatim above and in
the runner's package doc comment.
