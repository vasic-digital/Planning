// Copyright 2026 vasic-digital. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package planning

import (
	"context"
	"fmt"

	"digital.vasic.planning/pkg/i18n"
)

// English fallback templates used when no real i18n.Translator is wired
// (NoopTranslator path). Per CONST-046 these are NOT the canonical
// strings — the canonical source is pkg/i18n/bundles/active.en.yaml.
// They are duplicated here so that the standalone NoopTranslator path
// produces a sensible default without forcing the consuming project to
// ship a translator. A real translator returns localized text instead
// of the msgID, and these fallbacks are skipped.
const (
	fallbackMilestonePromptIntro = `Create a hierarchical plan for the following goal:
"%s"

Generate 3-5 high-level milestones that need to be achieved.
For each milestone, provide:
1. A clear name
2. A description
3. Any dependencies on other milestones (by number)

Format as numbered list with details.`

	fallbackStepPromptIntro = `For the milestone: "%s"
Description: %s

Generate 3-7 specific action steps to complete this milestone.
Each step should be:
1. Concrete and actionable
2. Independently executable
3. Verifiable for completion

Format as numbered list.`

	fallbackStepExecutionFailed = "execution failed"

	// --- Round 2 (CONST-046 Phase 4 round 96) ---
	// 7 additional user-facing strings: hints prompt, MCTS code-action
	// prompts (2), Tree-of-Thoughts prompts (3), evaluator prompt, the
	// HiPlan context label, plus a locale-overridable terminal-keyword
	// CSV. All routed through resolveOrFallback below.

	fallbackHintsPromptIntro = `Given this context:
%s

Generate 2-3 specific hints or tips to help execute this step successfully.
Focus on:
- Edge cases to consider
- Best practices
- Common pitfalls to avoid

Format as bullet points.`

	fallbackMCTSActionsPromptIntro = `Given the current code state:
%s

Generate 3-5 possible next coding actions or modifications.
Each action should be distinct and meaningful.
Format: one action per line.`

	fallbackMCTSApplyActionPromptIntro = `Current code state:
%s

Apply this action: %s

Return the updated code state after applying the action.`

	fallbackToTThoughtsPromptIntro = `Given the current reasoning step:
"%s"

Generate %d distinct next steps or approaches to continue solving this problem.
Each step should be different and explore a unique angle.
Format each step on a new line starting with a number.`

	fallbackToTInitialThoughtsPromptIntro = `Given the problem:
"%s"

Generate %d distinct initial approaches or strategies to solve this problem.
Each approach should be different and explore a unique angle.
Format each approach on a new line starting with a number.`

	fallbackToTEvaluatePromptIntro = `Evaluate the following reasoning step on a scale of 0.0 to 1.0:
"%s"

Consider:
- Logical validity
- Progress toward solution
- Feasibility
- Clarity

Respond with only a number between 0.0 and 1.0.`

	fallbackHiPlanContextLabel = `Milestone: %s
Description: %s
Step: %s
Action: %s`

	// fallbackToTTerminalKeywords is the English-locale CSV of tokens
	// the IsTerminal detector matches against. Consuming projects MUST
	// override via the i18n_tot_terminal_keywords msgID for non-English
	// locales, otherwise the detector silently fails to recognise
	// solution states the LLM emits in the user's language.
	fallbackToTTerminalKeywords = "solution,answer,result,conclusion,final"
)

// resolveOrFallback routes a user-facing string through tr.T. When the
// translator is the noop (msgID-verbatim path), the call site receives
// the msgID back and we substitute the bundled English fallback. When
// a real translator is wired, its result is used directly.
//
// This is the single seam every CONST-046 migration in the Planning
// module passes through: msgID stays grep-able, args stay positional,
// and the bundled fallback is the source of truth for English.
func resolveOrFallback(ctx context.Context, tr i18n.Translator, msgID, fallback string, args ...any) string {
	if tr == nil {
		tr = i18n.Default()
	}
	got := tr.T(ctx, msgID, args...)
	if got == msgID {
		// NoopTranslator path — translator did not resolve. Use the
		// English fallback bundled with the migration.
		return fmt.Sprintf(fallback, args...)
	}
	return got
}
