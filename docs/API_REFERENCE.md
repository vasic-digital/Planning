# Planning - API Reference

**Module:** `digital.vasic.planning`
**Package:** `planning`

## Constructor Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `NewHiPlan` | `NewHiPlan(config *HiPlanConfig, gen MilestoneGenerator, exec StepExecutor, logger *logrus.Logger) *HiPlan` | Creates a hierarchical planner. |
| `DefaultHiPlanConfig` | `DefaultHiPlanConfig() *HiPlanConfig` | Returns default HiPlan configuration. |
| `NewMCTS` | `NewMCTS(config *MCTSConfig, gen MCTSActionGenerator, reward MCTSRewardFunction, rollout MCTSRolloutPolicy, logger *logrus.Logger) *MCTS` | Creates an MCTS planner. |
| `DefaultMCTSConfig` | `DefaultMCTSConfig() *MCTSConfig` | Returns default MCTS configuration. |
| `NewTreeOfThoughts` | `NewTreeOfThoughts(config *TreeOfThoughtsConfig, gen ThoughtGenerator, eval ThoughtEvaluator, logger *logrus.Logger) *TreeOfThoughts` | Creates a Tree of Thoughts planner. |
| `DefaultTreeOfThoughtsConfig` | `DefaultTreeOfThoughtsConfig() *TreeOfThoughtsConfig` | Returns default ToT configuration. |

## HiPlan Types

### HiPlan

Main hierarchical planner struct.

**Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `Plan` | `(h *HiPlan) Plan(ctx context.Context, req PlanRequest) (*PlanResult, error)` | Generates a hierarchical execution plan. |

### HiPlanConfig

```go
type HiPlanConfig struct {
    MaxDepth             int
    MaxStepsPerMilestone int
    Timeout              time.Duration
    EnableParallel       bool
}
```

### Strategy Interfaces

```go
type MilestoneGenerator interface {
    GenerateMilestones(ctx context.Context, goal string, constraints []string) ([]*Milestone, error)
}

type StepExecutor interface {
    ExecuteStep(ctx context.Context, step *PlanStep) (*StepResult, error)
}
```

### Data Types

```go
type PlanRequest struct {
    Goal        string
    Constraints []string
    Context     map[string]interface{}
}

type HierarchicalPlan struct {
    ID         string
    Goal       string
    Milestones []*Milestone
    Status     string
    CreatedAt  time.Time
}

type Milestone struct {
    ID           string
    Name         string
    Description  string
    Steps        []*PlanStep
    Dependencies []string
    Status       string
}

type PlanStep struct {
    ID          string
    Name        string
    Description string
    Action      string
    Parameters  map[string]interface{}
}

type PlanResult struct {
    Plan     *HierarchicalPlan
    Duration time.Duration
    Error    error
}
```

Concrete implementation: `LLMMilestoneGenerator` (uses an LLM provider).

## MCTS Types

### MCTS

Main Monte Carlo Tree Search planner struct.

**Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `Search` | `(m *MCTS) Search(ctx context.Context, problem MCTSProblem) (*MCTSResult, error)` | Explores the solution space and returns the best path. |

### MCTSConfig

```go
type MCTSConfig struct {
    MaxIterations       int
    ExplorationConstant float64    // UCB1 constant (default: sqrt(2))
    MaxDepth            int
    SimulationDepth     int
    Timeout             time.Duration
}
```

### Strategy Interfaces

```go
type MCTSActionGenerator interface {
    GenerateActions(ctx context.Context, state string) ([]*MCTSAction, error)
}

type MCTSRewardFunction interface {
    Evaluate(ctx context.Context, state string, actions []*MCTSAction) (float64, error)
}

type MCTSRolloutPolicy interface {
    Rollout(ctx context.Context, node *MCTSNode) (float64, error)
}
```

### Data Types

```go
type MCTSProblem struct {
    InitialState string
    Goal         string
    MaxDepth     int
}

type MCTSNode struct {
    State      string
    Action     *MCTSAction
    Parent     *MCTSNode
    Children   []*MCTSNode
    Visits     int
    TotalValue float64
    Depth      int
}

type MCTSAction struct {
    Name   string
    Reward float64
}

type MCTSResult struct {
    BestPath   []*MCTSAction
    BestReward float64
    Iterations int
    NodesExplored int
    Duration   time.Duration
}
```

Concrete implementations: `CodeActionGenerator`, `CodeRewardFunction`,
`DefaultRolloutPolicy`.

## Tree of Thoughts Types

### TreeOfThoughts

Main ToT planner struct.

**Methods:**

| Method | Signature | Description |
|--------|-----------|-------------|
| `Solve` | `(t *TreeOfThoughts) Solve(ctx context.Context, problem ToTProblem) (*ToTResult, error)` | Explores reasoning paths and returns the best solution. |

### TreeOfThoughtsConfig

```go
type TreeOfThoughtsConfig struct {
    MaxBranches    int
    MaxDepth       int
    PruneThreshold float64
    SearchStrategy SearchStrategy  // BreadthFirst or DepthFirst
    MaxThoughts    int
}
```

### Strategy Interfaces

```go
type ThoughtGenerator interface {
    Generate(ctx context.Context, problem string, parent *Thought) ([]*Thought, error)
}

type ThoughtEvaluator interface {
    Evaluate(ctx context.Context, thought *Thought) (float64, error)
}
```

### Data Types

```go
type ToTProblem struct {
    Question    string
    MaxThoughts int
}

type Thought struct {
    ID       string
    Content  string
    Score    float64
    Depth    int
    ParentID string
}

type ThoughtNode struct {
    Thought  *Thought
    Children []*ThoughtNode
    Parent   *ThoughtNode
}

type ToTResult struct {
    BestPath    []*Thought
    BestScore   float64
    TotalNodes  int
    PrunedNodes int
    Duration    time.Duration
}
```

Concrete implementations: `LLMThoughtGenerator`, `LLMThoughtEvaluator`.

## Search Strategy Enum

| Constant | Value | Description |
|----------|-------|-------------|
| `SearchBreadthFirst` | `"bfs"` | Explore all siblings before going deeper |
| `SearchDepthFirst` | `"dfs"` | Explore depth-first, prune and backtrack |
