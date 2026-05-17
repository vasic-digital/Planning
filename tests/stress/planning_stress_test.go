// ACKNOWLEDGED: This stress suite previously used mocks in violation of
// CONST-050(A) (no-fakes-beyond-unit-tests). Per round-26 §11.4 audit
// (2026-05-17), every mock type that previously satisfied
// planning.MilestoneGenerator / StepExecutor / MCTSActionGenerator /
// MCTSRewardFunction / ThoughtGenerator / ThoughtEvaluator has been removed
// from this file. Tests now LOUDLY SKIP (per CONST-035 bluff taxonomy: loud
// skip > silent mock pass) until a real LLM (Ollama / Bedrock / etc.) and
// real executor are wired through helix-deps.yaml + the planning module's
// LLMMilestoneGenerator / CodeActionGenerator implementations.
//
// Tracking ticket: SKIP-OK: #PLANNING-INT-REAL — stress coverage owed
// against a real planning topology (real LLM + real executor + real reward
// function). See vasic-digital/Planning issue tracker for re-wiring plan.
package stress

import (
	"testing"
)

const skipReason = "SKIP-OK: #PLANNING-INT-REAL — requires real LLM + executor; stress coverage owed per CONST-050(A) round-26 §11.4 audit"

func TestHiPlan_ConcurrentPlanCreation_Stress(t *testing.T) {
	t.Skip(skipReason)
}

func TestHiPlan_ConcurrentPlanExecution_Stress(t *testing.T) {
	t.Skip(skipReason)
}

func TestHiPlan_ConcurrentLibraryAccess_Stress(t *testing.T) {
	t.Skip(skipReason)
}

func TestMCTS_ConcurrentSearches_Stress(t *testing.T) {
	t.Skip(skipReason)
}

func TestMCTSNode_ConcurrentRewardUpdates_Stress(t *testing.T) {
	t.Skip(skipReason)
}

func TestTreeOfThoughts_ConcurrentSolves_Stress(t *testing.T) {
	t.Skip(skipReason)
}

func TestHiPlan_RapidCreateAndExecute_Stress(t *testing.T) {
	t.Skip(skipReason)
}

func TestMixedAlgorithms_ConcurrentExecution_Stress(t *testing.T) {
	t.Skip(skipReason)
}
