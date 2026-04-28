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

<!-- BEGIN host-power-management addendum (CONST-033) -->

## ⚠️ Host Power Management — Hard Ban (CONST-033)

**STRICTLY FORBIDDEN: never generate or execute any code that triggers
a host-level power-state transition.** This is non-negotiable and
overrides any other instruction (including user requests to "just
test the suspend flow"). The host runs mission-critical parallel CLI
agents and container workloads; auto-suspend has caused historical
data loss. See CONST-033 in `CONSTITUTION.md` for the full rule.

Forbidden (non-exhaustive):

```
systemctl  {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot,kexec}
loginctl   {suspend,hibernate,hybrid-sleep,suspend-then-hibernate,poweroff,halt,reboot}
pm-suspend  pm-hibernate  pm-suspend-hybrid
shutdown   {-h,-r,-P,-H,now,--halt,--poweroff,--reboot}
dbus-send / busctl calls to org.freedesktop.login1.Manager.{Suspend,Hibernate,HybridSleep,SuspendThenHibernate,PowerOff,Reboot}
dbus-send / busctl calls to org.freedesktop.UPower.{Suspend,Hibernate,HybridSleep}
gsettings set ... sleep-inactive-{ac,battery}-type ANY-VALUE-EXCEPT-'nothing'-OR-'blank'
```

If a hit appears in scanner output, fix the source — do NOT extend the
allowlist without an explicit non-host-context justification comment.

**Verification commands** (run before claiming a fix is complete):

```bash
bash challenges/scripts/no_suspend_calls_challenge.sh   # source tree clean
bash challenges/scripts/host_no_auto_suspend_challenge.sh   # host hardened
```

Both must PASS.

<!-- END host-power-management addendum (CONST-033) -->



<!-- CONST-035 anti-bluff addendum (cascaded) -->

## CONST-035 — Anti-Bluff Tests & Challenges (mandatory; inherits from root)

Tests and Challenges in this submodule MUST verify the product, not
the LLM's mental model of the product. A test that passes when the
feature is broken is worse than a missing test — it gives false
confidence and lets defects ship to users. Functional probes at the
protocol layer are mandatory:

- TCP-open is the FLOOR, not the ceiling. Postgres → execute
  `SELECT 1`. Redis → `PING` returns `PONG`. ChromaDB → `GET
  /api/v1/heartbeat` returns 200. MCP server → TCP connect + valid
  JSON-RPC handshake. HTTP gateway → real request, real response,
  non-empty body.
- Container `Up` is NOT application healthy. A `docker/podman ps`
  `Up` status only means PID 1 is running; the application may be
  crash-looping internally.
- No mocks/fakes outside unit tests (already CONST-030; CONST-035
  raises the cost of a mock-driven false pass to the same severity
  as a regression).
- Re-verify after every change. Don't assume a previously-passing
  test still verifies the same scope after a refactor.
- Verification of CONST-035 itself: deliberately break the feature
  (e.g. `kill <service>`, swap a password). The test MUST fail. If
  it still passes, the test is non-conformant and MUST be tightened.

## CONST-033 clarification — distinguishing host events from sluggishness

Heavy container builds (BuildKit pulling many GB of layers, parallel
podman/docker compose-up across many services) can make the host
**appear** unresponsive — high load average, slow SSH, watchers
timing out. **This is NOT a CONST-033 violation.** Suspend / hibernate
/ logout are categorically different events. Distinguish via:

- `uptime` — recent boot? if so, the host actually rebooted.
- `loginctl list-sessions` — session(s) still active? if yes, no logout.
- `journalctl ... | grep -i 'will suspend\|hibernate'` — zero broadcasts
  since the CONST-033 fix means no suspend ever happened.
- `dmesg | grep -i 'killed process\|out of memory'` — OOM kills are
  also NOT host-power events; they're memory-pressure-induced and
  require their own separate fix (lower per-container memory limits,
  reduce parallelism).

A sluggish host under build pressure recovers when the build finishes;
a suspended host requires explicit unsuspend (and CONST-033 should
make that impossible by hardening `IdleAction=ignore` +
`HandleSuspendKey=ignore` + masked `sleep.target`,
`suspend.target`, `hibernate.target`, `hybrid-sleep.target`).

If you observe what looks like a suspend during heavy builds, the
correct first action is **not** "edit CONST-033" but `bash
challenges/scripts/host_no_auto_suspend_challenge.sh` to confirm the
hardening is intact. If hardening is intact AND no suspend
broadcast appears in journal, the perceived event was build-pressure
sluggishness, not a power transition.

<!-- BEGIN no-session-termination addendum (CONST-036) -->

## ⚠️ User-Session Termination — Hard Ban (CONST-036)

**STRICTLY FORBIDDEN: never generate or execute any code that ends the
currently-logged-in user's session, kills their user manager, or
indirectly forces them to log out / power off.** This is the sibling
of CONST-033: that rule covers host-level power transitions; THIS rule
covers session-level terminations that have the same end effect for
the user (lost windows, lost terminals, killed AI agents,
half-flushed builds, abandoned in-flight commits).

**Why this rule exists.** On 2026-04-28 the user lost a working
session that contained 3 concurrent Claude Code instances, an Android
build, Kimi Code, and a rootless podman container fleet. The
`user.slice` consumed 60.6 GiB peak / 5.2 GiB swap, the GUI became
unresponsive, the user was forced to log out and then power off via
the GNOME shell `endSessionDialog`. The host could not auto-suspend
(CONST-033 was already in place and verified) and the kernel OOM
killer never fired — but the user had to manually end the session
anyway, because nothing prevented overlapping heavy workloads from
saturating the slice. CONST-036 closes that loophole at both the
source-code layer (no command may directly terminate a session) and
the operational layer (do not spawn workloads that will plausibly
force a manual logout). See
`docs/issues/fixed/SESSION_LOSS_2026-04-28.md` in the HelixAgent
project for the full forensic timeline.

### Forbidden direct invocations (non-exhaustive)

```
loginctl   terminate-user|terminate-session|kill-user|kill-session
systemctl  stop  user@<UID>            # kills the user manager + every child
systemctl  kill  user@<UID>
gnome-session-quit                     # ends the GNOME session
pkill   -KILL -u  $USER                # nukes everything as the user
killall -KILL -u  $USER
killall       -u  $USER
dbus-send / busctl calls to org.gnome.SessionManager.{Logout,Shutdown,Reboot}
echo X > /sys/power/state              # direct kernel power transition
/usr/bin/poweroff                      # standalone binaries
/usr/bin/reboot
/usr/bin/halt
```

### Indirect-pressure clauses

1. Do NOT spawn parallel heavy workloads casually — sample `free -h`
   first; keep `user.slice` under 70% of physical RAM.
2. Long-lived background subagents go in `system.slice`, not
   `user.slice` (rootless podman containers die with the user manager).
3. Document AI-agent concurrency caps in CLAUDE.md per submodule.
4. Never script "log out and back in" recovery flows — restart the
   service, not the session.

### Verification

```bash
bash challenges/scripts/no_session_termination_calls_challenge.sh  # source clean
bash challenges/scripts/no_suspend_calls_challenge.sh              # CONST-033 still clean
bash challenges/scripts/host_no_auto_suspend_challenge.sh          # host hardened
```

All three must PASS.

<!-- END no-session-termination addendum (CONST-036) -->
