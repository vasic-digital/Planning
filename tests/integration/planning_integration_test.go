// ACKNOWLEDGED: This integration suite previously used mocks in violation of
// CONST-050(A) (no-fakes-beyond-unit-tests). Per round-26 §11.4 audit
// (2026-05-17), every mock type that previously satisfied
// planning.MilestoneGenerator / StepExecutor / MCTSActionGenerator /
// MCTSRewardFunction / ThoughtGenerator / ThoughtEvaluator has been removed
// from this file. Tests now LOUDLY SKIP (per CONST-035 bluff taxonomy: loud
// skip > silent mock pass) until a real LLM (Ollama / Bedrock / etc.) and
// real executor are wired through helix-deps.yaml + the planning module's
// LLMMilestoneGenerator / CodeActionGenerator implementations.
//
// Tracking ticket: SKIP-OK: #PLANNING-INT-REAL — integration coverage owed
// against a real planning topology (real LLM + real executor + real reward
// function). See vasic-digital/Planning issue tracker for re-wiring plan.
package integration

import (
	"testing"
)

const skipReason = "SKIP-OK: #PLANNING-INT-REAL — requires real LLM + executor; integration coverage owed per CONST-050(A) round-26 §11.4 audit"

func TestHiPlan_CreateAndExecutePlan_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestHiPlan_SequentialExecution_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestHiPlan_FailedSteps_AdaptivePlanning_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestMCTS_SearchAndExplore_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestMCTS_WithRolloutPolicy_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestTreeOfThoughts_BeamSearch_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestTreeOfThoughts_BFSSearch_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestTreeOfThoughts_DFSSearch_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestHiPlan_LibraryOperations_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestHiPlan_WithDependencies_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestMCTS_UCTValue_Integration(t *testing.T) {
	t.Skip(skipReason)
}

func TestAllThreeAlgorithms_Integration(t *testing.T) {
	t.Skip(skipReason)
}
