# Planning - Getting Started

**Module:** `digital.vasic.planning`

## Installation

```go
import "digital.vasic.planning/planning"
```

## Quick Start: Using Each Planner

### HiPlan (Hierarchical Planning)

Decomposes high-level goals into milestones and executable steps:

```go
package main

import (
    "context"
    "fmt"

    "digital.vasic.planning/planning"
    "github.com/sirupsen/logrus"
)

func main() {
    logger := logrus.New()
    config := planning.DefaultHiPlanConfig()

    // Provide implementations for milestone generation and step execution
    milestoneGen := &myMilestoneGenerator{}
    stepExec := &myStepExecutor{}

    planner := planning.NewHiPlan(config, milestoneGen, stepExec, logger)

    result, err := planner.Plan(context.Background(), planning.PlanRequest{
        Goal:        "Build a REST API for user management",
        Constraints: []string{"Use Go", "Include authentication"},
        Context:     map[string]interface{}{"language": "go"},
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Plan: %s\n", result.Plan.ID)
    for _, m := range result.Plan.Milestones {
        fmt.Printf("  Milestone: %s (%d steps)\n", m.Name, len(m.Steps))
    }
}
```

### MCTS (Monte Carlo Tree Search)

Explores solution spaces to find optimal action sequences:

```go
    mConfig := planning.DefaultMCTSConfig()
    mConfig.MaxIterations = 500
    mConfig.ExplorationConstant = 1.41

    actionGen := &planning.CodeActionGenerator{}
    rewardFn := &planning.CodeRewardFunction{}
    rollout := &planning.DefaultRolloutPolicy{}

    mcts := planning.NewMCTS(mConfig, actionGen, rewardFn, rollout, logger)

    result, err := mcts.Search(context.Background(), planning.MCTSProblem{
        InitialState: "function signature with empty body",
        Goal:         "implement sorting algorithm",
        MaxDepth:     10,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Best action sequence (%d steps):\n", len(result.BestPath))
    for _, action := range result.BestPath {
        fmt.Printf("  %s (reward: %.2f)\n", action.Name, action.Reward)
    }
```

### Tree of Thoughts

Explores multiple reasoning paths with evaluation and pruning:

```go
    tConfig := planning.DefaultTreeOfThoughtsConfig()
    tConfig.MaxBranches = 3
    tConfig.MaxDepth = 4
    tConfig.SearchStrategy = planning.SearchBreadthFirst

    thoughtGen := &planning.LLMThoughtGenerator{Provider: llmProvider}
    thoughtEval := &planning.LLMThoughtEvaluator{Provider: llmProvider}

    tot := planning.NewTreeOfThoughts(tConfig, thoughtGen, thoughtEval, logger)

    result, err := tot.Solve(context.Background(), planning.ToTProblem{
        Question:    "Design a caching strategy for a read-heavy microservice",
        MaxThoughts: 50,
    })
    if err != nil {
        panic(err)
    }

    fmt.Printf("Best thought path (score: %.2f):\n", result.BestScore)
    for _, thought := range result.BestPath {
        fmt.Printf("  [%.2f] %s\n", thought.Score, thought.Content)
    }
```

## Configuring Search Parameters

### HiPlan Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `MaxDepth` | `5` | Maximum decomposition depth |
| `MaxStepsPerMilestone` | `10` | Maximum steps in one milestone |
| `Timeout` | `10m` | Overall planning timeout |
| `EnableParallel` | `true` | Parallelize independent milestones |

### MCTS Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `MaxIterations` | `1000` | Number of search iterations |
| `ExplorationConstant` | `1.414` | UCB1 exploration factor (sqrt(2)) |
| `MaxDepth` | `20` | Maximum tree depth |
| `SimulationDepth` | `10` | Rollout simulation depth |
| `Timeout` | `5m` | Search timeout |

### Tree of Thoughts Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `MaxBranches` | `3` | Thoughts generated per node |
| `MaxDepth` | `5` | Maximum thought depth |
| `PruneThreshold` | `0.3` | Minimum score to keep a branch |
| `SearchStrategy` | `BreadthFirst` | `BreadthFirst` or `DepthFirst` |
| `MaxThoughts` | `100` | Total thought budget |

## Strategy Interfaces

Each algorithm relies on injected strategies for domain-specific logic:

| Algorithm | Interface | Methods |
|-----------|-----------|---------|
| HiPlan | `MilestoneGenerator` | `GenerateMilestones(ctx, goal, constraints) ([]*Milestone, error)` |
| HiPlan | `StepExecutor` | `ExecuteStep(ctx, step) (*StepResult, error)` |
| MCTS | `MCTSActionGenerator` | `GenerateActions(ctx, state) ([]*MCTSAction, error)` |
| MCTS | `MCTSRewardFunction` | `Evaluate(ctx, state, actions) (float64, error)` |
| MCTS | `MCTSRolloutPolicy` | `Rollout(ctx, node) (float64, error)` |
| ToT | `ThoughtGenerator` | `Generate(ctx, problem, parent) ([]*Thought, error)` |
| ToT | `ThoughtEvaluator` | `Evaluate(ctx, thought) (float64, error)` |

The module ships concrete implementations prefixed with `Code*` (for code
tasks) and `LLM*` (for LLM-backed reasoning).

## Integration with HelixAgent

The Planning module integrates through:

- **Adapter** at `internal/adapters/planning/adapter.go`
- Used by the debate orchestrator for planning debate strategies
- Used by the agentic workflow system for task decomposition

## Next Steps

- See [ARCHITECTURE.md](ARCHITECTURE.md) for algorithm design details
- See [API_REFERENCE.md](API_REFERENCE.md) for the full type catalog
