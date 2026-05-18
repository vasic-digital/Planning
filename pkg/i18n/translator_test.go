// Copyright 2026 vasic-digital. All rights reserved.
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"context"
	"testing"
)

// TestNoopTranslator_ReturnsMsgIDVerbatim asserts the standalone-default
// translator behavior: on T(ctx, "planning_foo"), it MUST return
// "planning_foo" exactly. Real production translators wire by the
// consuming project will resolve the ID to localized text; the noop
// guarantees the call site never crashes when no translator is wired.
// This is the contract that lets the submodule stay decoupled per
// CONST-051(B) while still routing user-facing strings through the
// Translator seam per CONST-046.
func TestNoopTranslator_ReturnsMsgIDVerbatim(t *testing.T) {
	tr := Default()
	got := tr.T(context.Background(), "planning_milestone_prompt_intro")
	if got != "planning_milestone_prompt_intro" {
		t.Fatalf("NoopTranslator.T returned %q; expected msgID verbatim", got)
	}
}

// TestNoopTranslator_IgnoresArgs confirms positional substitution args
// are accepted but unused by the noop. Real translators substitute;
// the noop is a pure pass-through so unit tests can assert "translator
// was called with this msgID" without coupling to a bundle.
func TestNoopTranslator_IgnoresArgs(t *testing.T) {
	tr := NoopTranslator{}
	got := tr.T(context.Background(), "planning_step_prompt_intro", "goal-x", 42)
	if got != "planning_step_prompt_intro" {
		t.Fatalf("NoopTranslator.T with args returned %q; expected msgID verbatim", got)
	}
}
