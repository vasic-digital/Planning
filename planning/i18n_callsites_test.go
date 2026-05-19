// Copyright 2026 vasic-digital. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package planning

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"digital.vasic.planning/pkg/i18n"
)

// fakeTranslator is a unit-test-only stub. Per CONST-050(A), mocks are
// permitted inside *_test.go (this file). It returns a sentinel
// "<TRANSLATED:msg_id>" so call-site tests assert against that exact
// marker — NEVER the original English literal. If a regression
// reintroduces the hardcoded literal at the call site, the assertion
// against the sentinel fails — which is the §11.4 anti-bluff guarantee.
type fakeTranslator struct {
	seen []string // msgIDs observed
}

func (f *fakeTranslator) T(_ context.Context, msgID string, _ ...any) string {
	f.seen = append(f.seen, msgID)
	return "<TRANSLATED:" + msgID + ">"
}

// TestLLMMilestoneGenerator_GenerateMilestones_RoutesPromptThroughTranslator
// asserts the milestone-prompt call site goes through the Translator.
// The fakeTranslator returns a sentinel; the test asserts the sentinel
// is what reaches the LLM generateFunc — NOT the original English
// "Create a hierarchical plan..." literal. If a future change inlines
// the literal again, this test fails (anti-bluff guarantee).
func TestLLMMilestoneGenerator_GenerateMilestones_RoutesPromptThroughTranslator(t *testing.T) {
	tr := &fakeTranslator{}
	var capturedPrompt string
	g := NewLLMMilestoneGenerator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "1. m1\n2. m2", nil
		},
		nil,
	)
	g.SetTranslator(tr)

	if _, err := g.GenerateMilestones(context.Background(), "build a CLI tool"); err != nil {
		t.Fatalf("GenerateMilestones returned error: %v", err)
	}

	want := "<TRANSLATED:planning_milestone_prompt_intro>"
	if capturedPrompt != want {
		t.Fatalf("LLM received %q; expected sentinel %q (call site bypassed Translator → CONST-046 violation)", capturedPrompt, want)
	}
	if len(tr.seen) != 1 || tr.seen[0] != "planning_milestone_prompt_intro" {
		t.Fatalf("translator observed msgIDs %v; expected [planning_milestone_prompt_intro]", tr.seen)
	}
}

// TestLLMMilestoneGenerator_GenerateSteps_RoutesPromptThroughTranslator
// is the sibling assertion for the step-prompt call site.
func TestLLMMilestoneGenerator_GenerateSteps_RoutesPromptThroughTranslator(t *testing.T) {
	tr := &fakeTranslator{}
	var capturedPrompt string
	g := NewLLMMilestoneGenerator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "1. s1\n2. s2", nil
		},
		nil,
	)
	g.SetTranslator(tr)

	m := &Milestone{Name: "milestone-x", Description: "do X"}
	if _, err := g.GenerateSteps(context.Background(), m); err != nil {
		t.Fatalf("GenerateSteps returned error: %v", err)
	}

	want := "<TRANSLATED:planning_step_prompt_intro>"
	if capturedPrompt != want {
		t.Fatalf("LLM received %q; expected sentinel %q (call site bypassed Translator → CONST-046 violation)", capturedPrompt, want)
	}
	if len(tr.seen) != 1 || tr.seen[0] != "planning_step_prompt_intro" {
		t.Fatalf("translator observed msgIDs %v; expected [planning_step_prompt_intro]", tr.seen)
	}
}

// TestLLMMilestoneGenerator_NoTranslator_UsesEnglishFallback documents the
// standalone path: when no real translator is wired, the bundled English
// fallback drives the prompt — keeping the module standalone-buildable
// per CONST-051(B). The fallback MUST contain the literal English text
// from active.en.yaml so a manual operator inspection at the LLM
// boundary sees coherent prose.
func TestLLMMilestoneGenerator_NoTranslator_UsesEnglishFallback(t *testing.T) {
	var capturedPrompt string
	g := NewLLMMilestoneGenerator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "1. m1", nil
		},
		nil,
	)
	// NoopTranslator default — no SetTranslator call.

	if _, err := g.GenerateMilestones(context.Background(), "demo goal"); err != nil {
		t.Fatalf("GenerateMilestones returned error: %v", err)
	}
	if !strings.Contains(capturedPrompt, "Create a hierarchical plan") {
		t.Fatalf("standalone fallback prompt missing English header; got %q", capturedPrompt)
	}
	if !strings.Contains(capturedPrompt, "demo goal") {
		t.Fatalf("standalone fallback prompt missing substituted goal; got %q", capturedPrompt)
	}
}

// TestHiPlan_SetTranslator_NilResetsToDefault is the defensive-default
// assertion. SetTranslator(nil) MUST reset to NoopTranslator, never
// leave the translator field nil (which would panic at call sites).
func TestHiPlan_SetTranslator_NilResetsToDefault(t *testing.T) {
	hp := NewHiPlan(DefaultHiPlanConfig(), nil, nil, nil)
	hp.SetTranslator(nil)
	if hp.translator == nil {
		t.Fatalf("SetTranslator(nil) left translator field nil; expected NoopTranslator default")
	}
	if got := hp.translator.T(context.Background(), "planning_step_execution_failed"); got != "planning_step_execution_failed" {
		t.Fatalf("translator after SetTranslator(nil) returned %q; expected NoopTranslator msgID-verbatim", got)
	}
	_ = i18n.NoopTranslator{} // import-keep
}

// ============================================================================
// Round 2 (CONST-046 Phase 4 round 96) — 8 call-site assertions + 2 standalone
// fallback assertions. Each paired-mutation test plants a sentinel through a
// fakeTranslator and asserts that exact sentinel reaches the LLM (or the
// IsTerminal detector). If a future refactor re-inlines any literal at the
// migrated call site, the corresponding assertion fails — the §11.4 anti-bluff
// guarantee for this round.
// ============================================================================

// TestLLMMilestoneGenerator_GenerateHints_RoutesPromptThroughTranslator
// asserts the hint-prompt call site goes through the Translator.
func TestLLMMilestoneGenerator_GenerateHints_RoutesPromptThroughTranslator(t *testing.T) {
	tr := &fakeTranslator{}
	var capturedPrompt string
	g := NewLLMMilestoneGenerator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "- hint1\n- hint2", nil
		},
		nil,
	)
	g.SetTranslator(tr)

	step := &PlanStep{ID: "s1", Action: "do thing"}
	if _, err := g.GenerateHints(context.Background(), step, "some context"); err != nil {
		t.Fatalf("GenerateHints returned error: %v", err)
	}
	want := "<TRANSLATED:planning_hints_prompt_intro>"
	if capturedPrompt != want {
		t.Fatalf("LLM received %q; expected sentinel %q (CONST-046 violation)", capturedPrompt, want)
	}
	if len(tr.seen) != 1 || tr.seen[0] != "planning_hints_prompt_intro" {
		t.Fatalf("translator observed msgIDs %v; expected [planning_hints_prompt_intro]", tr.seen)
	}
}

// TestLLMMilestoneGenerator_GenerateHints_NoTranslator_UsesEnglishFallback
// documents the standalone path keeps emitting coherent English.
func TestLLMMilestoneGenerator_GenerateHints_NoTranslator_UsesEnglishFallback(t *testing.T) {
	var capturedPrompt string
	g := NewLLMMilestoneGenerator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "- h", nil
		},
		nil,
	)
	step := &PlanStep{ID: "s1", Action: "do"}
	if _, err := g.GenerateHints(context.Background(), step, "CTX-XYZ"); err != nil {
		t.Fatalf("GenerateHints returned error: %v", err)
	}
	if !strings.Contains(capturedPrompt, "Given this context") {
		t.Fatalf("standalone fallback hint prompt missing English header; got %q", capturedPrompt)
	}
	if !strings.Contains(capturedPrompt, "CTX-XYZ") {
		t.Fatalf("standalone fallback prompt missing substituted context; got %q", capturedPrompt)
	}
}

// TestHiPlan_BuildContext_RoutesThroughTranslator asserts the context label
// composer goes through the Translator. Paired-mutation: with a fake
// translator returning a sentinel, the assembled context MUST be the
// sentinel — NOT the hardcoded "Milestone:..." literal.
func TestHiPlan_BuildContext_RoutesThroughTranslator(t *testing.T) {
	tr := &fakeTranslator{}
	hp := NewHiPlan(DefaultHiPlanConfig(), nil, nil, nil)
	hp.SetTranslator(tr)

	m := &Milestone{Name: "M1", Description: "Do M1"}
	s := &PlanStep{ID: "S1", Action: "act"}
	got := hp.buildContext(m, s)
	want := "<TRANSLATED:planning_hiplan_context_label>"
	if got != want {
		t.Fatalf("buildContext returned %q; expected sentinel %q (CONST-046 violation)", got, want)
	}
	if len(tr.seen) != 1 || tr.seen[0] != "planning_hiplan_context_label" {
		t.Fatalf("translator observed msgIDs %v; expected [planning_hiplan_context_label]", tr.seen)
	}
}

// TestHiPlan_BuildContext_NoTranslator_UsesEnglishFallback asserts the
// standalone path still emits the structured English label.
func TestHiPlan_BuildContext_NoTranslator_UsesEnglishFallback(t *testing.T) {
	hp := NewHiPlan(DefaultHiPlanConfig(), nil, nil, nil)
	m := &Milestone{Name: "M1", Description: "Do M1"}
	s := &PlanStep{ID: "S1", Action: "act"}
	got := hp.buildContext(m, s)
	for _, want := range []string{"Milestone: M1", "Description: Do M1", "Step: S1", "Action: act"} {
		if !strings.Contains(got, want) {
			t.Fatalf("fallback context label missing %q; got %q", want, got)
		}
	}
}

// TestCodeActionGenerator_GetActions_RoutesPromptThroughTranslator
// asserts the MCTS action-listing prompt site goes through the Translator.
func TestCodeActionGenerator_GetActions_RoutesPromptThroughTranslator(t *testing.T) {
	tr := &fakeTranslator{}
	var capturedPrompt string
	g := NewCodeActionGenerator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "1. a\n2. b", nil
		},
		nil,
	)
	g.SetTranslator(tr)

	if _, err := g.GetActions(context.Background(), "code-state-x"); err != nil {
		t.Fatalf("GetActions returned error: %v", err)
	}
	want := "<TRANSLATED:planning_mcts_actions_prompt_intro>"
	if capturedPrompt != want {
		t.Fatalf("LLM received %q; expected sentinel %q (CONST-046 violation)", capturedPrompt, want)
	}
}

// TestCodeActionGenerator_ApplyAction_RoutesPromptThroughTranslator
// asserts the MCTS action-application prompt site goes through the Translator.
func TestCodeActionGenerator_ApplyAction_RoutesPromptThroughTranslator(t *testing.T) {
	tr := &fakeTranslator{}
	var capturedPrompt string
	g := NewCodeActionGenerator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "new-state", nil
		},
		nil,
	)
	g.SetTranslator(tr)

	if _, err := g.ApplyAction(context.Background(), "state", "act"); err != nil {
		t.Fatalf("ApplyAction returned error: %v", err)
	}
	want := "<TRANSLATED:planning_mcts_apply_action_prompt_intro>"
	if capturedPrompt != want {
		t.Fatalf("LLM received %q; expected sentinel %q (CONST-046 violation)", capturedPrompt, want)
	}
}

// TestLLMThoughtGenerator_GenerateThoughts_RoutesPromptThroughTranslator
// asserts the ToT child-thoughts prompt site goes through the Translator.
func TestLLMThoughtGenerator_GenerateThoughts_RoutesPromptThroughTranslator(t *testing.T) {
	tr := &fakeTranslator{}
	var capturedPrompt string
	g := NewLLMThoughtGenerator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "1. t1", nil
		},
		0.5, nil,
	)
	g.SetTranslator(tr)

	parent := &Thought{ID: "p1", Content: "think about X"}
	if _, err := g.GenerateThoughts(context.Background(), parent, 3); err != nil {
		t.Fatalf("GenerateThoughts returned error: %v", err)
	}
	want := "<TRANSLATED:planning_tot_thoughts_prompt_intro>"
	if capturedPrompt != want {
		t.Fatalf("LLM received %q; expected sentinel %q (CONST-046 violation)", capturedPrompt, want)
	}
}

// TestLLMThoughtGenerator_GenerateInitialThoughts_RoutesPromptThroughTranslator
// asserts the ToT initial-thoughts prompt site goes through the Translator.
func TestLLMThoughtGenerator_GenerateInitialThoughts_RoutesPromptThroughTranslator(t *testing.T) {
	tr := &fakeTranslator{}
	var capturedPrompt string
	g := NewLLMThoughtGenerator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "1. approach", nil
		},
		0.5, nil,
	)
	g.SetTranslator(tr)

	if _, err := g.GenerateInitialThoughts(context.Background(), "problem-X", 2); err != nil {
		t.Fatalf("GenerateInitialThoughts returned error: %v", err)
	}
	want := "<TRANSLATED:planning_tot_initial_thoughts_prompt_intro>"
	if capturedPrompt != want {
		t.Fatalf("LLM received %q; expected sentinel %q (CONST-046 violation)", capturedPrompt, want)
	}
}

// TestLLMThoughtEvaluator_EvaluateThought_RoutesPromptThroughTranslator
// asserts the ToT evaluator prompt site goes through the Translator.
func TestLLMThoughtEvaluator_EvaluateThought_RoutesPromptThroughTranslator(t *testing.T) {
	tr := &fakeTranslator{}
	var capturedPrompt string
	e := NewLLMThoughtEvaluator(
		func(_ context.Context, prompt string) (string, error) {
			capturedPrompt = prompt
			return "0.5", nil
		},
		nil,
	)
	e.SetTranslator(tr)

	if _, err := e.EvaluateThought(context.Background(), &Thought{Content: "step-X"}); err != nil {
		t.Fatalf("EvaluateThought returned error: %v", err)
	}
	want := "<TRANSLATED:planning_tot_evaluate_prompt_intro>"
	if capturedPrompt != want {
		t.Fatalf("LLM received %q; expected sentinel %q (CONST-046 violation)", capturedPrompt, want)
	}
}

// TestLLMThoughtEvaluator_TerminalKeywords_LocaleOverride asserts that a
// real Translator can override the English-default terminal-keyword list.
// This is the CONST-046 anti-bluff guarantee for the IsTerminal detector:
// the English defaults silently fail to recognise solution states in
// non-English LLM outputs unless the consuming project overrides them.
// A Serbian-locale translator returning "rešenje,odgovor,rezultat,zaključak,konačno"
// MUST cause IsTerminal to match a Serbian content line.
func TestLLMThoughtEvaluator_TerminalKeywords_LocaleOverride(t *testing.T) {
	srTr := &fixedTranslator{
		responses: map[string]string{
			"planning_tot_terminal_keywords": "rešenje,odgovor,rezultat,zaključak,konačno",
		},
	}
	e := NewLLMThoughtEvaluator(
		func(_ context.Context, _ string) (string, error) { return "0", nil },
		nil,
	)
	e.SetTranslator(srTr)

	// Confirm English defaults are GONE — IsTerminal MUST NOT match
	// the English "solution" any more after a Serbian override.
	thoughtEN := &Thought{Content: "this is the solution"}
	gotEN, err := e.IsTerminal(context.Background(), thoughtEN)
	if err != nil {
		t.Fatalf("IsTerminal(EN) returned error: %v", err)
	}
	if gotEN {
		t.Fatalf("after Serbian override, English 'solution' still matched IsTerminal — locale override silently failed (CONST-046 violation)")
	}

	// And the Serbian keyword DOES match.
	thoughtSR := &Thought{Content: "ovo je naše konačno rešenje"}
	gotSR, err := e.IsTerminal(context.Background(), thoughtSR)
	if err != nil {
		t.Fatalf("IsTerminal(SR) returned error: %v", err)
	}
	if !gotSR {
		t.Fatalf("Serbian terminal-keyword 'rešenje' / 'konačno' did NOT match IsTerminal under Serbian translator — locale override silently failed (CONST-046 violation)")
	}
}

// fixedTranslator is a unit-test-only stub returning a pre-configured
// per-msgID response. Permitted in *_test.go per CONST-050(A). Unlike
// fakeTranslator (which returns a sentinel), fixedTranslator returns
// the real translation so downstream parsers (parseTerminalKeywords,
// fmt.Sprintf substitution, etc.) can execute against realistic input.
type fixedTranslator struct {
	responses map[string]string
}

func (f *fixedTranslator) T(_ context.Context, msgID string, args ...any) string {
	if v, ok := f.responses[msgID]; ok {
		if len(args) > 0 {
			return fmt.Sprintf(v, args...)
		}
		return v
	}
	return msgID
}
