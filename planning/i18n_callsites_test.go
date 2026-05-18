// Copyright 2026 vasic-digital. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package planning

import (
	"context"
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
