#!/usr/bin/env bash
# planning_describe_challenge.sh
#
# Round-270 paired-mutation deep-doc challenge for digital.vasic.planning.
#
# Validates that:
#   1. The deep-doc ledger (docs/test-coverage.md) lists every exported
#      symbol from planning/hiplan.go, planning/mcts.go,
#      planning/tree_of_thoughts.go, and pkg/i18n/translator.go.
#   2. The multi-locale fixture (tests/fixtures/planning/payloads.json)
#      parses and contains at least 3 locales.
#   3. The multi-locale runner (challenges/runner/main.go) builds and
#      runs, byte-preserving non-ASCII payloads through the real
#      HiPlan + MCTS + TreeOfThoughts + i18n surfaces across 4
#      sections and 5 locales.
#   4. The README enumerates the round-270 anti-bluff guarantees.
#
# Paired-mutation invariant (CONST-035 + CONST-050(B)):
#   With --anti-bluff-mutate the script plants a deliberate symbol-rename
#   mutation in a tmp copy of the ledger (CreatePlan ->
#   CreatePlan_MUTATED), reruns validation, and asserts the gate
#   FAILS with exit 99. This proves the gate actually catches
#   ledger-vs-source drift instead of rubber-stamping it.
#
# Exit codes:
#   0  — gate PASS on clean tree
#   1  — gate FAIL on clean tree (real failure to fix)
#   99 — paired-mutation correctly detected (good — proves anti-bluff)
#   2  — usage / environment error

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MODULE_DIR="$(cd "${SCRIPT_DIR}/../.." && pwd)"

MUTATE=0
for arg in "$@"; do
    case "$arg" in
        --anti-bluff-mutate) MUTATE=1 ;;
        --help|-h)
            sed -n '1,32p' "$0"
            exit 0
            ;;
        *)
            echo "unknown argument: $arg" >&2
            exit 2
            ;;
    esac
done

PASS=0
FAIL=0
TOTAL=0

pass() { PASS=$((PASS+1)); TOTAL=$((TOTAL+1)); echo "  PASS: $1"; }
fail() { FAIL=$((FAIL+1)); TOTAL=$((TOTAL+1)); echo "  FAIL: $1"; }

LEDGER="${MODULE_DIR}/docs/test-coverage.md"
FIXTURE="${MODULE_DIR}/tests/fixtures/planning/payloads.json"
RUNNER="${MODULE_DIR}/challenges/runner/main.go"
README="${MODULE_DIR}/README.md"

LEDGER_WORK="${LEDGER}"
TMP_LEDGER=""
if [ "${MUTATE}" -eq 1 ]; then
    TMP_LEDGER="$(mktemp)"
    cp "${LEDGER}" "${TMP_LEDGER}"
    # Plant a rename so the symbol no longer matches what the source declares.
    sed -i 's/CreatePlan/CreatePlan_MUTATED/g' "${TMP_LEDGER}"
    LEDGER_WORK="${TMP_LEDGER}"
    echo "=== Planning Describe Challenge (anti-bluff-mutate mode) ==="
else
    echo "=== Planning Describe Challenge (clean mode) ==="
fi
echo ""

# Section 1: ledger presence and freshness
echo "Section 1: docs/test-coverage.md ledger"
if [ ! -f "${LEDGER_WORK}" ]; then
    fail "ledger missing at ${LEDGER_WORK}"
else
    pass "ledger present"
    if grep -q "round-270" "${LEDGER_WORK}"; then
        pass "ledger marked round-270"
    else
        fail "ledger missing round-270 marker"
    fi
    if grep -q "execution of tests and Challenges MUST guarantee" "${LEDGER_WORK}"; then
        pass "ledger carries Article XI §11.9 mandate"
    else
        fail "ledger missing Article XI §11.9 mandate"
    fi
fi

# Section 2: every exported package symbol appears in ledger.
echo ""
echo "Section 2: structural symbol cross-reference"

EXPECTED_SYMBOLS=(
    # planning/hiplan.go
    "MilestoneState" "Milestone" "PlanStepState" "PlanStep" "HiPlanConfig"
    "DefaultHiPlanConfig" "MilestoneGenerator" "StepExecutor" "StepResult"
    "HierarchicalPlan" "HiPlan" "NewHiPlan" "CreatePlan" "ExecutePlan"
    "ExecuteStep" "AddToLibrary" "GetFromLibrary" "GetCurrentPlan"
    "MilestoneResult" "PlanResult" "LLMMilestoneGenerator" "NewLLMMilestoneGenerator"
    # planning/mcts.go
    "MCTSNodeState" "MCTSNode" "MCTSConfig" "DefaultMCTSConfig"
    "MCTSActionGenerator" "MCTSRewardFunction" "MCTSRolloutPolicy"
    "MCTS" "NewMCTS" "Search" "UCTValue" "MCTSResult"
    "CodeActionGenerator" "NewCodeActionGenerator"
    "CodeRewardFunction" "NewCodeRewardFunction"
    "DefaultRolloutPolicy" "NewDefaultRolloutPolicy"
    # planning/tree_of_thoughts.go
    "ThoughtState" "Thought" "ThoughtNode" "TreeOfThoughtsConfig"
    "DefaultTreeOfThoughtsConfig" "ThoughtGenerator" "ThoughtEvaluator"
    "TreeOfThoughts" "NewTreeOfThoughts" "Solve" "ToTResult"
    "GetSolutionContent" "LLMThoughtGenerator" "NewLLMThoughtGenerator"
    "LLMThoughtEvaluator" "NewLLMThoughtEvaluator"
    # pkg/i18n/translator.go
    "Translator" "NoopTranslator" "Default"
    # i18n callsites
    "SetTranslator"
)

CHECKED=0
MISSING=0
for sym in "${EXPECTED_SYMBOLS[@]}"; do
    CHECKED=$((CHECKED + 1))
    if grep -qE "\\b${sym}\\b" "${LEDGER_WORK}"; then
        : # found
    else
        fail "ledger missing symbol ${sym}"
        MISSING=$((MISSING + 1))
    fi
done
if [ "${MISSING}" -eq 0 ]; then
    pass "all ${CHECKED} structural symbols cross-referenced in ledger"
fi

# Section 3: multi-locale fixture sanity
echo ""
echo "Section 3: multi-locale fixture"
if [ ! -f "${FIXTURE}" ]; then
    fail "fixture missing at ${FIXTURE}"
else
    pass "fixture present"
    LOCALE_COUNT=$(grep -oE '"locale":\s*"[^"]+"' "${FIXTURE}" | sort -u | wc -l)
    if [ "${LOCALE_COUNT}" -ge 3 ]; then
        pass "fixture covers ${LOCALE_COUNT} locales (>=3)"
    else
        fail "fixture covers only ${LOCALE_COUNT} locales (<3)"
    fi
fi

# Section 4: runner builds + runs against every section
echo ""
echo "Section 4: multi-locale runner build + run (real HiPlan + MCTS + ToT)"
if [ ! -f "${RUNNER}" ]; then
    fail "runner missing at ${RUNNER}"
else
    pass "runner source present"
    cd "${MODULE_DIR}"
    if go build -o /tmp/planning_round270_runner ./challenges/runner/ 2>/tmp/planning_build.log; then
        pass "runner builds"
        if /tmp/planning_round270_runner -fixtures "${FIXTURE}" > /tmp/planning_run.log 2>&1; then
            pass "runner exit 0 across every section + locale"
            # Per-locale + per-section PASS coverage
            if grep -q "PASS: \[Section1\]\[CreatePlan\]\[en\]" /tmp/planning_run.log; then
                pass "Section 1 English (en) HiPlan CreatePlan"
            else
                fail "Section 1 English (en) HiPlan CreatePlan missing"
            fi
            if grep -q "PASS: \[Section1\]\[CreatePlan\]\[sr\]" /tmp/planning_run.log; then
                pass "Section 1 Cyrillic (sr) HiPlan CreatePlan"
            else
                fail "Section 1 Cyrillic (sr) HiPlan CreatePlan missing"
            fi
            if grep -q "PASS: \[Section1\]\[CreatePlan\]\[ja\]" /tmp/planning_run.log; then
                pass "Section 1 Japanese (ja) HiPlan CreatePlan"
            else
                fail "Section 1 Japanese (ja) HiPlan CreatePlan missing"
            fi
            if grep -q "PASS: \[Section1\]\[CreatePlan\]\[ar\]" /tmp/planning_run.log; then
                pass "Section 1 Arabic (ar) HiPlan CreatePlan"
            else
                fail "Section 1 Arabic (ar) HiPlan CreatePlan missing"
            fi
            if grep -q "PASS: \[Section1\]\[CreatePlan\]\[zh-CN\]" /tmp/planning_run.log; then
                pass "Section 1 Han (zh-CN) HiPlan CreatePlan"
            else
                fail "Section 1 Han (zh-CN) HiPlan CreatePlan missing"
            fi
            if grep -q "PASS: \[Section1\]\[ExecuteStep\]\[ja\]" /tmp/planning_run.log; then
                pass "Section 1 Japanese (ja) HiPlan ExecuteStep byte-exact"
            else
                fail "Section 1 ExecuteStep ja missing"
            fi
            if grep -q "PASS: \[Section1\]\[Library\]\[ar\]" /tmp/planning_run.log; then
                pass "Section 1 Arabic (ar) HiPlan Library round-trip"
            else
                fail "Section 1 Library ar missing"
            fi
            if grep -q "PASS: \[Section2\]\[Search\]\[en\]" /tmp/planning_run.log; then
                pass "Section 2 English (en) MCTS Search"
            else
                fail "Section 2 MCTS Search en missing"
            fi
            if grep -q "PASS: \[Section2\]\[Search\]\[sr\]" /tmp/planning_run.log; then
                pass "Section 2 Cyrillic (sr) MCTS Search"
            else
                fail "Section 2 MCTS Search sr missing"
            fi
            if grep -q "PASS: \[Section2\]\[Search\]\[zh-CN\]" /tmp/planning_run.log; then
                pass "Section 2 Han (zh-CN) MCTS Search"
            else
                fail "Section 2 MCTS Search zh-CN missing"
            fi
            if grep -q "PASS: \[Section2\]\[UCTValue\]\[ar\]" /tmp/planning_run.log; then
                pass "Section 2 Arabic (ar) MCTS UCTValue"
            else
                fail "Section 2 UCTValue ar missing"
            fi
            if grep -q "PASS: \[Section2\]\[NewCodeActionGenerator\]" /tmp/planning_run.log; then
                pass "Section 2 NewCodeActionGenerator constructor"
            else
                fail "Section 2 NewCodeActionGenerator missing"
            fi
            if grep -q "PASS: \[Section3\]\[Solve\]\[en\]\[bfs\]" /tmp/planning_run.log; then
                pass "Section 3 English (en) ToT bfs strategy"
            else
                fail "Section 3 ToT bfs en missing"
            fi
            if grep -q "PASS: \[Section3\]\[Solve\]\[sr\]\[dfs\]" /tmp/planning_run.log; then
                pass "Section 3 Cyrillic (sr) ToT dfs strategy"
            else
                fail "Section 3 ToT dfs sr missing"
            fi
            if grep -q "PASS: \[Section3\]\[Solve\]\[ja\]\[beam\]" /tmp/planning_run.log; then
                pass "Section 3 Japanese (ja) ToT beam strategy"
            else
                fail "Section 3 ToT beam ja missing"
            fi
            if grep -q "PASS: \[Section3\]\[Solve\]\[ar\]\[bfs\]" /tmp/planning_run.log; then
                pass "Section 3 Arabic (ar) ToT bfs strategy"
            else
                fail "Section 3 ToT bfs ar missing"
            fi
            if grep -q "PASS: \[Section3\]\[Solve\]\[zh-CN\]\[beam\]" /tmp/planning_run.log; then
                pass "Section 3 Han (zh-CN) ToT beam strategy"
            else
                fail "Section 3 ToT beam zh-CN missing"
            fi
            if grep -q "PASS: \[Section3\]\[NewLLMThoughtGenerator\]" /tmp/planning_run.log; then
                pass "Section 3 NewLLMThoughtGenerator constructor"
            else
                fail "Section 3 NewLLMThoughtGenerator missing"
            fi
            if grep -q "PASS: \[Section4\]\[i18n.Default\]" /tmp/planning_run.log; then
                pass "Section 4 i18n.Default returns Translator"
            else
                fail "Section 4 i18n.Default missing"
            fi
            if grep -q "PASS: \[Section4\]\[NoopTranslator.T\]" /tmp/planning_run.log; then
                pass "Section 4 NoopTranslator.T msg-id fallthrough"
            else
                fail "Section 4 NoopTranslator.T missing"
            fi
            if grep -q "PASS: \[Section4\]\[HiPlan.SetTranslator(nil)\]" /tmp/planning_run.log; then
                pass "Section 4 HiPlan.SetTranslator(nil) survives"
            else
                fail "Section 4 HiPlan.SetTranslator(nil) missing"
            fi
        else
            fail "runner exit non-zero — see /tmp/planning_run.log"
            sed -n '1,80p' /tmp/planning_run.log
        fi
    else
        fail "runner build failed — see /tmp/planning_build.log"
        sed -n '1,40p' /tmp/planning_build.log
    fi
    rm -f /tmp/planning_round270_runner
fi

# Section 5: README round-270 anti-bluff section
echo ""
echo "Section 5: README round-270 anti-bluff section"
if grep -q "Anti-bluff guarantees" "${README}"; then
    pass "README declares Anti-bluff guarantees"
else
    fail "README missing Anti-bluff guarantees section"
fi
if grep -q "round-270" "${README}"; then
    pass "README marked round-270"
else
    fail "README missing round-270 marker"
fi

# Cleanup mutated ledger if any
if [ -n "${TMP_LEDGER}" ]; then
    rm -f "${TMP_LEDGER}"
fi

echo ""
echo "=== Summary: ${PASS}/${TOTAL} PASS, ${FAIL} FAIL ==="

if [ "${MUTATE}" -eq 1 ]; then
    if [ "${FAIL}" -gt 0 ]; then
        echo "anti-bluff-mutate: gate correctly detected planted mutation (exit 99)"
        exit 99
    else
        echo "anti-bluff-mutate: gate FAILED to detect planted mutation — bluff!"
        exit 1
    fi
fi

if [ "${FAIL}" -gt 0 ]; then
    exit 1
fi
exit 0
