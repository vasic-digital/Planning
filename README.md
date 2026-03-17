# Planning

`digital.vasic.planning` -- AI planning algorithms with Hierarchical Planning (HiPlan), Monte Carlo Tree Search (MCTS), and Tree of Thoughts (ToT).

## Overview

Planning is a Go module that implements three complementary AI planning algorithms for decomposing complex tasks into executable steps, exploring solution spaces through tree search, and reasoning through multi-step problems via structured thought generation.

**HiPlan (Hierarchical Planning)** decomposes high-level goals into milestones with dependency relationships, then breaks each milestone into concrete action steps. It supports parallel milestone execution with dependency-aware topological sorting, contextual hint generation for each step, and adaptive planning that continues past individual step failures.

**MCTS (Monte Carlo Tree Search)** implements the MASTER framework for code generation and general planning. It uses the UCB1 selection policy with an optional UCT-DP (depth-preferred) extension, parallel simulations, configurable rollout policies, and discounted reward backpropagation. The module includes concrete implementations for code action generation and code quality evaluation.

**Tree of Thoughts (ToT)** extends chain-of-thought reasoning into a tree structure where multiple reasoning paths are explored in parallel. It supports three search strategies (BFS, DFS, beam search), thought pruning below quality thresholds, path-level evaluation with depth-weighted scoring, and backtracking on dead ends.

All three algorithms accept pluggable strategy interfaces for generation, evaluation, and execution, making them adaptable to any domain beyond code generation.

## Architecture

```
+----------------------------------------------------------+
|                    Planning Module                         |
+----------------------------------------------------------+
|                                                            |
|  +----------------+  +----------------+  +---------------+ |
|  |     HiPlan     |  |     MCTS       |  |    ToT        | |
|  +-------+--------+  +-------+--------+  +------+--------+ |
|          |                    |                   |          |
|  +-------v--------+  +-------v--------+  +------v--------+ |
|  | Milestone       |  | MCTSAction     |  | Thought       | |
|  | Generator       |  | Generator      |  | Generator     | |
|  | (LLM-backed)    |  | (Code-backed)  |  | (LLM-backed)  | |
|  +----------------+  +----------------+  +---------------+ |
|  | StepExecutor    |  | RewardFunction |  | Thought       | |
|  | (pluggable)     |  | (Code quality) |  | Evaluator     | |
|  +----------------+  | RolloutPolicy  |  | (LLM-backed)  | |
|                       +----------------+  +---------------+ |
+----------------------------------------------------------+

Data flow:
  Goal/Problem -> Algorithm -> Solution/Plan/Actions
```

## Package Structure

| Package | Purpose |
|---------|---------|
| `planning` | Core module: HiPlan, MCTS, Tree of Thoughts -- all algorithms and types |

### Source Files

| File | Description |
|------|-------------|
| `hiplan.go` | `HiPlan` -- hierarchical planning with milestones, steps, hints, parallel/sequential execution |
| `mcts.go` | `MCTS` -- Monte Carlo Tree Search with UCB1/UCT-DP selection, expansion, simulation, backpropagation |
| `tree_of_thoughts.go` | `TreeOfThoughts` -- BFS/DFS/beam search over thought trees with pruning and evaluation |

## API Reference

### HiPlan Types and Interfaces

```go
// HiPlan implements Hierarchical Planning
type HiPlan struct { ... }

func NewHiPlan(config HiPlanConfig, generator MilestoneGenerator, executor StepExecutor, logger *logrus.Logger) *HiPlan
func (h *HiPlan) CreatePlan(ctx context.Context, goal string) (*HierarchicalPlan, error)
func (h *HiPlan) ExecutePlan(ctx context.Context, plan *HierarchicalPlan) (*PlanResult, error)
func (h *HiPlan) ExecuteStep(ctx context.Context, step *PlanStep, hints []string) (*StepResult, error)
func (h *HiPlan) AddToLibrary(milestone *Milestone)
func (h *HiPlan) GetFromLibrary(id string) (*Milestone, bool)
func (h *HiPlan) GetCurrentPlan() *HierarchicalPlan

// MilestoneGenerator generates milestones from a goal
type MilestoneGenerator interface {
    GenerateMilestones(ctx context.Context, goal string) ([]*Milestone, error)
    GenerateSteps(ctx context.Context, milestone *Milestone) ([]*PlanStep, error)
    GenerateHints(ctx context.Context, step *PlanStep, context string) ([]string, error)
}

// StepExecutor executes plan steps
type StepExecutor interface {
    Execute(ctx context.Context, step *PlanStep, hints []string) (*StepResult, error)
    Validate(ctx context.Context, step *PlanStep, result *StepResult) error
}

// LLMMilestoneGenerator -- LLM-backed implementation
func NewLLMMilestoneGenerator(generateFunc func(ctx, prompt) (string, error), logger) *LLMMilestoneGenerator
```

**Milestone states**: `MilestoneStatePending`, `MilestoneStateInProgress`, `MilestoneStateCompleted`, `MilestoneStateFailed`, `MilestoneStateSkipped`

**Step states**: `PlanStepStatePending`, `PlanStepStateInProgress`, `PlanStepStateCompleted`, `PlanStepStateFailed`

### MCTS Types and Interfaces

```go
// MCTS implements Monte Carlo Tree Search
type MCTS struct { ... }

func NewMCTS(config MCTSConfig, actionGen MCTSActionGenerator, rewardFunc MCTSRewardFunction, rolloutPolicy MCTSRolloutPolicy, logger *logrus.Logger) *MCTS
func (m *MCTS) Search(ctx context.Context, initialState interface{}) (*MCTSResult, error)
func (m *MCTS) UCTValue(node *MCTSNode, parentVisits int) float64

// MCTSActionGenerator generates possible actions from a state
type MCTSActionGenerator interface {
    GetActions(ctx context.Context, state interface{}) ([]string, error)
    ApplyAction(ctx context.Context, state interface{}, action string) (interface{}, error)
}

// MCTSRewardFunction evaluates states
type MCTSRewardFunction interface {
    Evaluate(ctx context.Context, state interface{}) (float64, error)
    IsTerminal(ctx context.Context, state interface{}) (bool, error)
}

// MCTSRolloutPolicy performs simulation rollouts
type MCTSRolloutPolicy interface {
    Rollout(ctx context.Context, state interface{}, depth int) (float64, error)
}

// Concrete implementations for code generation
func NewCodeActionGenerator(generateFunc, logger) *CodeActionGenerator
func NewCodeRewardFunction(evaluateFunc, testFunc, logger) *CodeRewardFunction
func NewDefaultRolloutPolicy(actionGen, rewardFunc) *DefaultRolloutPolicy
```

**Node states**: `MCTSNodeStateUnexpanded`, `MCTSNodeStateExpanded`, `MCTSNodeStateTerminal`

### Tree of Thoughts Types and Interfaces

```go
// TreeOfThoughts implements the ToT reasoning framework
type TreeOfThoughts struct { ... }

func NewTreeOfThoughts(config TreeOfThoughtsConfig, generator ThoughtGenerator, evaluator ThoughtEvaluator, logger *logrus.Logger) *TreeOfThoughts
func (t *TreeOfThoughts) Solve(ctx context.Context, problem string) (*ToTResult, error)

// ThoughtGenerator generates new thoughts
type ThoughtGenerator interface {
    GenerateThoughts(ctx context.Context, parent *Thought, count int) ([]*Thought, error)
    GenerateInitialThoughts(ctx context.Context, problem string, count int) ([]*Thought, error)
}

// ThoughtEvaluator evaluates thought quality
type ThoughtEvaluator interface {
    EvaluateThought(ctx context.Context, thought *Thought) (float64, error)
    EvaluatePath(ctx context.Context, path []*Thought) (float64, error)
    IsTerminal(ctx context.Context, thought *Thought) (bool, error)
}

// LLM-backed implementations
func NewLLMThoughtGenerator(generateFunc, temperature, logger) *LLMThoughtGenerator
func NewLLMThoughtEvaluator(evaluateFunc, logger) *LLMThoughtEvaluator
```

**Thought states**: `ThoughtStatePending`, `ThoughtStateActive`, `ThoughtStateEvaluated`, `ThoughtStatePruned`, `ThoughtStateSelected`

**Search strategies**: `"bfs"` (breadth-first), `"dfs"` (depth-first), `"beam"` (beam search, default)

## Usage Examples

### Hierarchical Planning

```go
config := planning.DefaultHiPlanConfig()
generator := planning.NewLLMMilestoneGenerator(llmGenerate, logger)
hp := planning.NewHiPlan(config, generator, executor, logger)

plan, err := hp.CreatePlan(ctx, "Build a REST API with authentication")
// plan.Milestones: [Setup, Auth, Endpoints, Tests, Deploy]

result, err := hp.ExecutePlan(ctx, plan)
fmt.Printf("Completed: %d/%d milestones\n",
    result.CompletedMilestones, result.CompletedMilestones+result.FailedMilestones)
```

### Monte Carlo Tree Search for code generation

```go
mctsConfig := planning.DefaultMCTSConfig()
mctsConfig.MaxIterations = 500
mctsConfig.UseUCTDP = true

actionGen := planning.NewCodeActionGenerator(llmGenerate, logger)
rewardFunc := planning.NewCodeRewardFunction(evaluateCode, runTests, logger)
rollout := planning.NewDefaultRolloutPolicy(actionGen, rewardFunc)

mcts := planning.NewMCTS(mctsConfig, actionGen, rewardFunc, rollout, logger)
result, err := mcts.Search(ctx, initialCodeState)

fmt.Printf("Best actions: %v\n", result.BestActions)
fmt.Printf("Final reward: %.2f (tree size: %d nodes)\n", result.FinalReward, result.TreeSize)
```

### Tree of Thoughts for problem solving

```go
totConfig := planning.DefaultTreeOfThoughtsConfig()
totConfig.SearchStrategy = "beam"
totConfig.BeamWidth = 3
totConfig.MaxDepth = 8

generator := planning.NewLLMThoughtGenerator(llmGenerate, 0.7, logger)
evaluator := planning.NewLLMThoughtEvaluator(llmEvaluate, logger)
tot := planning.NewTreeOfThoughts(totConfig, generator, evaluator, logger)

result, err := tot.Solve(ctx, "Design a distributed cache with consistency guarantees")
fmt.Printf("Solution path (%d steps, score %.2f):\n", len(result.Solution), result.BestScore)
for _, thought := range result.Solution {
    fmt.Printf("  [%.2f] %s\n", thought.Score, thought.Content)
}
```

## Configuration

### HiPlan Configuration

```go
type HiPlanConfig struct {
    MaxMilestones            int           // Max milestones (default: 20)
    MaxStepsPerMilestone     int           // Max steps per milestone (default: 50)
    EnableParallelMilestones bool          // Parallel execution (default: true)
    MaxParallelMilestones    int           // Concurrency limit (default: 3)
    EnableAdaptivePlanning   bool          // Continue past failures (default: true)
    RetryFailedSteps         bool          // Auto-retry (default: true)
    MaxRetries               int           // Retries per step (default: 3)
    Timeout                  time.Duration // Overall timeout (default: 30m)
    StepTimeout              time.Duration // Per-step timeout (default: 5m)
}
```

### MCTS Configuration

```go
type MCTSConfig struct {
    ExplorationConstant  float64       // UCB1 C parameter (default: sqrt(2))
    DepthPreferenceAlpha float64       // UCT-DP depth bonus (default: 0.5)
    MaxDepth             int           // Max tree depth (default: 50)
    MaxIterations        int           // Total iterations (default: 1000)
    RolloutDepth         int           // Simulation depth (default: 10)
    SimulationCount      int           // Simulations per expansion (default: 1)
    DiscountFactor       float64       // Future reward discount (default: 0.99)
    EnableParallel       bool          // Parallel simulations (default: true)
    ParallelWorkers      int           // Worker count (default: 4)
    Timeout              time.Duration // Search timeout (default: 5m)
    UseUCTDP             bool          // Use depth-preferred UCT (default: true)
}
```

### Tree of Thoughts Configuration

```go
type TreeOfThoughtsConfig struct {
    MaxDepth           int           // Max thought tree depth (default: 10)
    MaxBranches        int           // Branches per node (default: 5)
    MinScore           float64       // Min score to consider (default: 0.3)
    PruneThreshold     float64       // Score below which to prune (default: 0.2)
    SearchStrategy     string        // "bfs", "dfs", or "beam" (default: "beam")
    BeamWidth          int           // Beam width for beam search (default: 3)
    Temperature        float64       // Generation diversity (default: 0.7)
    EnableBacktracking bool          // Allow backtracking (default: true)
    MaxIterations      int           // Total iterations (default: 100)
    Timeout            time.Duration // Search timeout (default: 5m)
}
```

## Testing

```bash
go build ./...
go test ./... -count=1 -race
```

## Integration with HelixAgent

Planning connects to HelixAgent through the adapter at `internal/adapters/planning/`:

- **Debate Planning**: HiPlan is used to decompose complex debate topics into structured milestones, where each milestone represents a debate phase (proposal, critique, optimization, convergence).
- **Code Generation**: MCTS drives iterative code refinement in HelixAgent's code generation pipeline, exploring multiple solution paths and selecting the highest-quality result.
- **Complex Reasoning**: Tree of Thoughts extends single-response reasoning into multi-path exploration for difficult analytical questions, improving answer quality through systematic consideration of alternatives.
- **SpecKit Integration**: HiPlan's milestone decomposition maps directly to SpecKit's 7-phase development flow, providing planning capabilities for large-scale feature development.
