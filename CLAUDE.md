# CLAUDE.md - Planning Module


## Definition of Done

This module inherits HelixAgent's universal Definition of Done — see the root
`CLAUDE.md` and `docs/development/definition-of-done.md`. In one line: **no
task is done without pasted output from a real run of the real system in the
same session as the change.** Coverage and green suites are not evidence.

### Acceptance demo for this module

```bash
# HiPlan + MCTS + Tree-of-Thoughts planning algorithms
cd Planning && GOMAXPROCS=2 nice -n 19 go test -count=1 -race -v \
  -run 'TestFullPlanningWorkflow_HiPlan_E2E|TestFullPlanningWorkflow_MCTS_E2E|TestFullPlanningWorkflow_TreeOfThoughts_E2E' \
  ./tests/e2e/...
```
Expect: three E2E PASS — each algorithm produces a valid plan for its test task.


## Overview

`digital.vasic.planning` is a Go module for AI planning algorithms: HiPlan (hierarchical planning),
Monte Carlo Tree Search (MCTS), and Tree of Thoughts (ToT).

**Module**: `digital.vasic.planning` (Go 1.24+)

## Build & Test

```bash
go build ./...
go test ./... -count=1 -race
```

## Package Structure

| Package | Purpose |
|---------|---------|
| `planning` | Core planning algorithms: HiPlan, MCTS, ToT |

## Key Types

### HiPlan (Hierarchical Planning)

- `HiPlan` — Main hierarchical planner struct
- `HiPlanConfig` / `DefaultHiPlanConfig()` — Configuration with defaults
- `MilestoneGenerator` — Interface for generating plan milestones
- `StepExecutor` — Interface for executing individual steps
- `HierarchicalPlan` / `Milestone` / `PlanStep` / `PlanResult` — Plan data types
- `LLMMilestoneGenerator` — LLM-backed milestone generator implementation

### MCTS (Monte Carlo Tree Search)

- `MCTS` — Main MCTS planner struct
- `MCTSConfig` / `DefaultMCTSConfig()` — Configuration with defaults
- `MCTSActionGenerator` / `MCTSRewardFunction` / `MCTSRolloutPolicy` — Strategy interfaces
- `MCTSNode` / `MCTSResult` — Tree node and result types
- `CodeActionGenerator` / `CodeRewardFunction` / `DefaultRolloutPolicy` — Concrete implementations

### Tree of Thoughts

- `TreeOfThoughts` — Main ToT planner struct
- `TreeOfThoughtsConfig` / `DefaultTreeOfThoughtsConfig()` — Configuration with defaults
- `ThoughtGenerator` / `ThoughtEvaluator` — Strategy interfaces
- `Thought` / `ThoughtNode` / `ToTResult` — Thought tree data types
- `LLMThoughtGenerator` / `LLMThoughtEvaluator` — LLM-backed implementations

## Mandatory Development Standards

- 100% test coverage across unit, integration, and benchmark tests
- No mocks outside unit tests — all other tests use real implementations
- Challenges must validate real-life use cases, not just return codes
- Follow Conventional Commits: `feat(planning): ...`, `fix(planning): ...`

## Integration Seams

| Direction | Sibling modules |
|-----------|-----------------|
| Upstream (this module imports) | none |
| Downstream (these import this module) | HelixLLM |

*Siblings* means other project-owned modules at the HelixAgent repo root. The root HelixAgent app and external systems are not listed here — the list above is intentionally scoped to module-to-module seams, because drift *between* sibling modules is where the "tests pass, product broken" class of bug most often lives. See root `CLAUDE.md` for the rules that keep these seams contract-tested.
