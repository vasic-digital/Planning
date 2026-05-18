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
