# AGENTS.md - Planning Module

## Overview

Planning module provides AI planning algorithms including hierarchical planning,
Monte Carlo Tree Search, and Tree of Thoughts for use in AI agent workflows.

## Key Files

- `planning/hiplan.go` — Hierarchical planning (HiPlan, MilestoneGenerator, StepExecutor)
- `planning/mcts.go` — Monte Carlo Tree Search (MCTS, MCTSActionGenerator, MCTSRewardFunction)
- `planning/tree_of_thoughts.go` — Tree of Thoughts (TreeOfThoughts, ThoughtGenerator, ThoughtEvaluator)

## Exported Types Summary

### hiplan.go
- `MilestoneState`, `PlanStepState` — State enums
- `Milestone`, `PlanStep`, `HierarchicalPlan`, `PlanResult`, `MilestoneResult`, `StepResult` — Data types
- `HiPlanConfig`, `DefaultHiPlanConfig()` — Configuration
- `MilestoneGenerator`, `StepExecutor` — Interfaces
- `HiPlan`, `NewHiPlan()` — Core planner
- `LLMMilestoneGenerator`, `NewLLMMilestoneGenerator()` — LLM-backed generator

### mcts.go
- `MCTSNodeState` — State enum
- `MCTSNode`, `MCTSResult` — Tree types
- `MCTSConfig`, `DefaultMCTSConfig()` — Configuration
- `MCTSActionGenerator`, `MCTSRewardFunction`, `MCTSRolloutPolicy` — Interfaces
- `MCTS`, `NewMCTS()` — Core planner
- `CodeActionGenerator`, `CodeRewardFunction`, `DefaultRolloutPolicy` — Concrete implementations

### tree_of_thoughts.go
- `ThoughtState` — State enum
- `Thought`, `ThoughtNode`, `ToTResult` — Tree types
- `TreeOfThoughtsConfig`, `DefaultTreeOfThoughtsConfig()` — Configuration
- `ThoughtGenerator`, `ThoughtEvaluator` — Interfaces
- `TreeOfThoughts`, `NewTreeOfThoughts()` — Core planner
- `LLMThoughtGenerator`, `LLMThoughtEvaluator` — LLM-backed implementations

## Integration with HelixAgent

The adapter at `internal/adapters/planning/adapter.go` bridges the internal
`dev.helix.agent/internal/planning` package to this extracted module.
Use `planningadapter.New(logger)` to obtain an adapter instance.

## Development Standards

- All code must compile and pass `go vet ./...`
- Tests must use table-driven style with `testify`
- No mocks outside unit tests
- Run challenges before submitting: `./challenges/scripts/planning_challenge.sh`

<!-- BEGIN host-power-management addendum (CONST-033) -->

## Host Power Management — Hard Ban (CONST-033)

**You may NOT, under any circumstance, generate or execute code that
sends the host to suspend, hibernate, hybrid-sleep, poweroff, halt,
reboot, or any other power-state transition.** This rule applies to:

- Every shell command you run via the Bash tool.
- Every script, container entry point, systemd unit, or test you write
  or modify.
- Every CLI suggestion, snippet, or example you emit.

**Forbidden invocations** (non-exhaustive — see CONST-033 in
`CONSTITUTION.md` for the full list):

- `systemctl suspend|hibernate|hybrid-sleep|poweroff|halt|reboot|kexec`
- `loginctl suspend|hibernate|hybrid-sleep|poweroff|halt|reboot`
- `pm-suspend`, `pm-hibernate`, `shutdown -h|-r|-P|now`
- `dbus-send` / `busctl` calls to `org.freedesktop.login1.Manager.Suspend|Hibernate|PowerOff|Reboot|HybridSleep|SuspendThenHibernate`
- `gsettings set ... sleep-inactive-{ac,battery}-type` to anything but `'nothing'` or `'blank'`

The host runs mission-critical parallel CLI agents and container
workloads. Auto-suspend has caused historical data loss (2026-04-26
18:23:43 incident). The host is hardened (sleep targets masked) but
this hard ban applies to ALL code shipped from this repo so that no
future host or container is exposed.

**Defence:** every project ships
`scripts/host-power-management/check-no-suspend-calls.sh` (static
scanner) and
`challenges/scripts/no_suspend_calls_challenge.sh` (challenge wrapper).
Both MUST be wired into the project's CI / `run_all_challenges.sh`.

**Full background:** `docs/HOST_POWER_MANAGEMENT.md` and `CONSTITUTION.md` (CONST-033).

<!-- END host-power-management addendum (CONST-033) -->

