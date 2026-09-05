#!/bin/bash
# Offline regression checks for Anthropic→Codex mission routing and decision state.
set -uo pipefail

ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
DERIVE="$ROOT/tools/launchd/derive-planner-lane.sh"
FIXTURE="$ROOT/tools/launchd/testdata/planner-lane/n-no-backtick-bullet.md"
PASS=0; FAIL=0
ok() { PASS=$((PASS+1)); echo "  PASS: $1"; }
bad() { FAIL=$((FAIL+1)); echo "  FAIL: $1 (got: $2)"; }
want() { [ "$2" = "$3" ] && ok "$1" || bad "$1" "$2"; }

out=$(MISSION_PLANNER_MODEL=codex:gpt-5.6-sol MISSION_ANTHROPIC_AVAILABLE=1 "$DERIVE" "$FIXTURE" 2>/dev/null)
want "planner remains fail-closed to Opus while Anthropic is available" "$out" "opus fail-closed:unparsable-path-entry"

out=$(MISSION_PLANNER_MODEL=codex:gpt-5.6-sol MISSION_ANTHROPIC_AVAILABLE=0 "$DERIVE" "$FIXTURE" 2>/dev/null)
want "planner falls back to Codex Sol when Anthropic is unavailable" "$out" "codex:gpt-5.6-sol anthropic-fallback:fail-closed:unparsable-path-entry"

out=$(MISSION_PLANNER_MODEL=codex:gpt-5.6-sol MISSION_ANTHROPIC_AVAILABLE=0 \
  MISSION_PLANNER_ANTHROPIC_FALLBACK=codex:test-sol "$DERIVE" "$FIXTURE" 2>/dev/null)
want "planner Anthropic fallback is configurable" "$out" "codex:test-sol anthropic-fallback:fail-closed:unparsable-path-entry"

driver="$ROOT/tools/launchd/mission-control.sh"
mission_doc="${MISSION_DOC_PATH:-$ROOT/design_docs/world-mission.md}"
grep -q 'MISSION_EXECUTOR_MODEL:-codex:gpt-5.6-sol' "$driver" \
  && ok "executor remains Codex Sol primary" || bad "executor remains Codex Sol primary" "missing default"
# Chain now leads with the FLAT-RATE ollama rung and keeps the metered OpenRouter twin
# behind it — same weights, so exhaustion degrades the ROUTE, not the model (matches the
# shared ailang driver, 2026-09-05).
grep -q 'MISSION_EXECUTOR_FALLBACK:-pi:ollama/deepseek-v4-flash:0731-cloud,pi:openrouter/deepseek/deepseek-v4-flash-0731:floor' "$driver" \
  && ok "executor chain is flat-rate ollama -> metered twin" || bad "executor chain is flat-rate ollama -> metered twin" "missing fallback"

# ─── QUOTA-DROUGHT SURVIVAL (2026-09-05, Mark attended) ───────────────────────────────
# World is a separate repo and had missed every routing fix in sunholo-data/ailang. Before
# this, its pre-flight probed the EXECUTOR and nothing else, so on a dry Anthropic bucket
# the designer (fable), planner (opus) and evaluator (sonnet) all died mid-iteration.
grep -q '_an_probed' "$driver" && grep -q '_mc_probe "\$an_model"' "$driver" \
  && ok "anthropic lanes get a role pre-flight" || bad "anthropic lanes get a role pre-flight" "no anthropic role probe"
[ "$(grep -c 'for role in DESIGNER PLANNER EXECUTOR EVALUATOR; do' "$driver")" = "3" ] \
  && ok "all three provider loops cover all four roles" \
  || bad "all three provider loops cover all four roles" "a provider loop is missing or role-limited"
# The planner's fallback used to be `opus` — the very model it was pinned to.
grep -q 'MISSION_PLANNER_FALLBACK:-opus}' "$driver" \
  && bad "planner fallback is not itself" "planner still falls back to opus, the model it is pinned to" \
  || ok "planner fallback is not itself"
grep -q 'MISSION_PLANNER_FALLBACK:-codex:gpt-5.6-sol' "$driver" \
  && ok "planner falls to codex first" || bad "planner falls to codex first" "planner chain missing"
grep -q 'MISSION_DESIGNER_FALLBACK:-codex:gpt-6-astra' "$driver" \
  && ok "designer has a codex rung" || bad "designer has a codex rung" "designer chain missing"
grep -q 'MISSION_EVALUATOR_FALLBACK:-pi:ollama/minimax-m3:cloud' "$driver" \
  && ok "evaluator has a chain at all" || bad "evaluator has a chain at all" "evaluator was sonnet-or-nothing"
# codex must stay the evaluator's LAST rung: the executor is codex:gpt-5.6-sol, so a codex
# judge is the vendor-level generator==judge collision, acceptable only as the last resort.
grep -q 'MISSION_EVALUATOR_FALLBACK:-codex:' "$driver" \
  && bad "codex is the evaluator's LAST rung" "a codex judge was promoted to the head of the chain" \
  || ok "codex is the evaluator's LAST rung"
# The ledger must be EMITTED, not just accumulated. World's five silently-demoted
# iterations (18/19/21/22) are why the shared driver grew one; world never got the emit side.
grep -q 'LANES DEGRADED THIS FIRE' "$driver" \
  && ok "lane degradation is reported before the iteration starts" \
  || bad "lane degradation is reported before the iteration starts" "ledger accumulates with no emit site"
# Astra sits BETWEEN the Anthropic rungs, and the selector must dispatch on provider or the
# codex entry would be handed to the claude CLI.
grep -q 'PREFS="\${MISSION_MODEL_PREFS:-claude-opus-5,codex:gpt-6-astra' "$driver" \
  && ok "astra sits between the Anthropic controller rungs" \
  || bad "astra sits between the Anthropic controller rungs" "controller ladder not updated"
grep -q 'PROVIDER-DISPATCHED since 2026-09-05' "$driver" \
  && ok "controller ladder dispatches on provider" \
  || bad "controller ladder dispatches on provider" "a codex PREFS entry would go to the claude CLI"
grep -q 'MISSION_CONTROLLER_FALLBACK:-codex:gpt-5.6-sol' "$driver" \
  && ok "controller has Codex Sol fallback" || bad "controller has Codex Sol fallback" "missing fallback"

"$ROOT/scripts/mission_decisions.sh" --check --file "$mission_doc" >/dev/null \
  && ok "decision ledger validates" || bad "decision ledger validates" "invalid"
open=$("$ROOT/scripts/mission_decisions.sh" --open --file "$mission_doc")
case "$open" in *$'D-WORLD-ROUTE-1\t'*) bad "resolved routing decision is not re-asked" "D-WORLD-ROUTE-1 appeared OPEN" ;; *) ok "resolved routing decision is not re-asked" ;; esac

# Exercise the real controller selector with probes stubbed: all Anthropic
# candidates fail, then the configured Codex subscription fallback succeeds.
lab=$(mktemp -d "${TMPDIR:-/tmp}/mission-routing.XXXXXX") || exit 1
awk '/^_mc_probe_codex\(\) \{/,/^\}/' "$driver" > "$lab/select.sh"
awk '/^_mc_set_controller\(\) \{/,/^\}/' "$driver" >> "$lab/select.sh"
awk '/^select_model\(\) \{/,/^\}/' "$driver" >> "$lab/select.sh"
out=$(/bin/bash -c '
  set -uo pipefail
  . "$1"
  _mc_probe() { return 1; }
  _mc_probe_codex() { [ "$1" = gpt-5.6-sol ]; }
  log() { :; }
  PREFS="claude-opus-5,claude-fable-5"
  CONTROLLER_FALLBACK="codex:gpt-5.6-sol"
  OVERRIDE_FILE="$2/no-override"
  select_model || exit 1
  printf "%s|%s|%s" "$CONTROLLER_ID" "$MISSION_ANTHROPIC_AVAILABLE" "$MODEL_WHY"
' _ "$lab/select.sh" "$lab")
want "controller selector traverses Anthropic to Codex" "$out" "codex:gpt-5.6-sol|0|Anthropic unavailable; subscription fallback"
rm -rf "$lab"

echo ""
echo "==== $PASS passed, $FAIL failed ===="
[ "$FAIL" -eq 0 ]
