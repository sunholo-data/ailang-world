#!/usr/bin/env bash
# mission-control.sh — continuous outer-loop iterations for ONE mission (default: V1).
# PORTABLE (M-MISSION-PORTABILITY M1, 2026-07-21, attended): MISSION_PROFILE=<name>
# sources ~/.config/ailang/mission-<name>.env (MISSION_NAME/REPO/DOC/...). MISSION_NAME=v1
# (the default) keeps the LEGACY state paths + log name EXACTLY as before — bit-for-bit,
# no migration; any other name gets fully namespaced state so two missions never collide.
#
# Fires a headless controller session that runs the mission-control skill:
# observe mission state → pick top backlog item → route through the inner-loop
# skills (design-doc → sprint-plan → execute → evaluate) → record → retro.
# See design_docs/v1-mission.md for the charter and guardrails.
#
# Scheduled via launchd StartInterval every 2h (see the plist); the overlap
# guard below makes this effectively "back-to-back iterations, ≤2h idle gap".
# Iterations are cloud-model work: they NEVER take rig.lock (GPU mutex only —
# GPU-touching sprint steps take it per-step inside the session).
#
# MODEL SELECTION (fleet Phase A, 2026-07-14): ordered preference probing.
# MISSION_MODEL_PREFS (default "claude-opus-5,claude-opus-4-8,claude-fable-5"
# — Opus 5 first since 2026-07-27 (Mark), 4.8 kept as probe fallback — OPUS-FIRST
# since 2026-07-16, Mark: Fable is reserved for high-cognition ROLES — design
# synthesis + evaluation, both bounded pinned sub-agents — never the long
# orchestration session, which burned the weekly Fable bucket at 2h cadence)
# is walked each iteration with a 1-token probe; first usable model wins. A
# quota-limited probe falls through to the next candidate; transient errors
# retry once. Fable last = emergency fallback only (a controller on Fable
# beats no controller). Semantics of the ordered list follow
# internal/ai/routing.go AIRoutingPolicy.Order (the third-vocabulary rule in
# m-mission-adaptive-multiprovider-routing); it lives in bash because the
# driver must select BEFORE any Go/claude process exists.
# Manual pins still win: MISSION_MODEL env (absolute) or
# ~/.ailang/state/mission-model ("<model> [expiry-epoch]", auto-expires).
#
# Transient Anthropic errors (Overloaded/dropped socket) are retried with backoff
# (TRANSIENT-RETRY block); deliberate watchdog kills are not.
# Kill switch: touch ~/.ailang/state/mission-control.disabled
# Portable to macOS bash 3.2. No GNU timeout on this rig → bash watchdog below.
set -uo pipefail

REPO="${MISSION_WORKDIR:-$(cd "$(dirname "$0")/../.." && pwd)}"
cd "$REPO" || exit 1

# launchd PATH is restricted; claude lives in ~/.local/bin, go tools in ~/go/bin.
export PATH="$HOME/go/bin:$HOME/.local/bin:$PATH"

# Dead-slot fix (Mark, attended 2026-08-10). The harness terminates background
# tasks at 600s. A controller that spawns its executor as a background Agent and
# ends its turn to wait is killed there — the driver then logs
# `iteration complete (rc=0)`, NO watchdog fires (HARD TIMEOUT/STALL: both
# absent), and the slot ends with a plausible transcript, zero charter rows and
# zero commits. A clean rc=0 is what a dead iteration looks like. Measured: the
# only 2 hits of "Background tasks still running after 600s" in 67 iterations
# are exactly the 2 orphaned slots (iter-64 stranded 525 lines of SM.C).
# 0 = no ceiling. The driver's own HARD_TIMEOUT/STALL watchdogs remain the bound.
export CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0

# --- MISSION PROFILE + STATE NAMESPACE (M1, 2026-07-21) ----------------------
[ -n "${MISSION_PROFILE:-}" ] && [ -f "$HOME/.config/ailang/mission-${MISSION_PROFILE}.env" ] \
  && . "$HOME/.config/ailang/mission-${MISSION_PROFILE}.env"
MISSION_NAME="${MISSION_NAME:-v1}"
MISSION_REPO="${MISSION_REPO:-sunholo-data/ailang}"
MISSION_DOC="${MISSION_DOC:-design_docs/v1-mission.md}"
export MISSION_NAME MISSION_REPO MISSION_DOC
STATE_DIR="$HOME/.ailang/state"
if [ "$MISSION_NAME" = "v1" ]; then
  # LEGACY paths — bit-for-bit compat with the live V1 loop (no migration).
  LOG=/tmp/ailang-mission-control.log
  KILL_SWITCH="$STATE_DIR/mission-control.disabled"
  PIDFILE="$STATE_DIR/mission-control.pid"
  OVERRIDE_FILE="$STATE_DIR/mission-model"
  LAST_MODEL_FILE="$STATE_DIR/mission-model-last"
  EXEC_ONCE_FILE="$STATE_DIR/mission-executor-model-once"
  GH_ISSUE_FILE="$STATE_DIR/mission-gh-issue"
  MSG_FROM="mission-control"
else
  LOG="/tmp/ailang-mission-${MISSION_NAME}.log"
  KILL_SWITCH="$STATE_DIR/mission-${MISSION_NAME}.disabled"
  PIDFILE="$STATE_DIR/mission-${MISSION_NAME}.pid"
  OVERRIDE_FILE="$STATE_DIR/mission-${MISSION_NAME}-model"
  LAST_MODEL_FILE="$STATE_DIR/mission-${MISSION_NAME}-model-last"
  EXEC_ONCE_FILE="$STATE_DIR/mission-${MISSION_NAME}-executor-model-once"
  GH_ISSUE_FILE="$STATE_DIR/mission-${MISSION_NAME}-gh-issue"
  MSG_FROM="mission-${MISSION_NAME}"
fi
# -----------------------------------------------------------------------------
[ -f "$HOME/.config/ailang/secrets.env" ] && . "$HOME/.config/ailang/secrets.env"

# BILLING GUARD (2026-07-10): the mission MUST bill the Claude subscription,
# never API credits. secrets.env exports ANTHROPIC_API_KEY for other tools —
# strip it so claude's only auth paths are subscription ones (keychain OAuth,
# or CLAUDE_CODE_OAUTH_TOKEN if set). Subscription-or-nothing by construction.
# 2026-07-27 extension (same construction, codex lane): secrets.env also exports the
# METERED OPENAI_API_KEY — strip it so codex's only auth is the ChatGPT-subscription
# OAuth in ~/.codex/auth.json (auth_mode=chatgpt). Metered OpenAI runs happen outside
# mission iterations, deliberately.
unset ANTHROPIC_API_KEY ANTHROPIC_AUTH_TOKEN OPENAI_API_KEY

log() { echo "[$(date '+%F %H:%M:%S')] $*" | tee -a "$LOG"; }

# --- stall detection (see the stall watchdog below) -------------------------
# _mc_descendants PID → echoes PID and every descendant PID (one per line).
_mc_descendants() {
  local pid="$1"; echo "$pid"
  local kids k; kids=$(pgrep -P "$pid" 2>/dev/null)
  for k in $kids; do _mc_descendants "$k"; done
}
# _mc_etime_secs "[[DD-]HH:]MM:SS" → seconds (macOS ps has no `etimes`).
_mc_etime_secs() {
  local t="${1// /}" dd=0 hh=0 mm=0 ss=0 rest nf
  [ -n "$t" ] || { echo 0; return; }
  case "$t" in *-*) dd=${t%%-*}; rest=${t#*-} ;; *) rest="$t" ;; esac
  nf=$(( $(printf '%s' "$rest" | tr -cd ':' | wc -c) + 1 ))
  if [ "$nf" -ge 3 ]; then hh=${rest%%:*}; rest=${rest#*:}; fi
  mm=${rest%%:*}; ss=${rest##*:}
  echo $(( 10#${dd:-0}*86400 + 10#${hh:-0}*3600 + 10#${mm:-0}*60 + 10#${ss:-0} ))
}
# _mc_progress_bytes → a byte counter that GROWS while the controller works,
# or non-zero rc when no such instrument exists for this provider.
#
# WHY THIS REPLACED THE CPU TEST (measured 2026-09-02). The old idleness arm was
# instantaneous `ps %cpu` summed over the tree, and the comment above it claimed
# "a session doing real work reads non-idle and is NOT flagged ... we miss late
# stalls, never kill live work". That claim is false, and it was the load-bearing
# one. Sampled against a LIVE v1 controller whose transcript was growing 15-45KB
# per 30s, the same expression read 0.10 / 0.30 / 0.80 / 1.40 — under the 2%
# floor on most samples, because an agent spends its wall-clock BLOCKED ON THE
# MODEL API, not on CPU. Cost of the false premise in one day: 4 V1 and 3 world
# iterations killed; the 21:13 kill landed on a session at Gate 5 that was
# committing its own iteration-321 record, which is why 321 has no record.
#
# The session transcript is the one signal the controller cannot forget to emit —
# the harness appends to it on every assistant message and every tool result. The
# heartbeat cannot be the primary arm: it is stamped by the AGENT at gate
# boundaries, and that same v1 slot reached Gate 5 having stamped only `gate-0`,
# so a heartbeat-only test reads "dead" on a session writing its own record.
_mc_progress_bytes() {
  local dir newest total=0 got=0 sz
  case "${CONTROLLER_PROVIDER:-claude}" in
    claude)
      # `claude -p` appends to ~/.claude/projects/<cwd, / and . mapped to ->/<uuid>.jsonl
      dir="$HOME/.claude/projects/$(printf '%s' "$PWD" | tr '/.' '--')"
      if [ -d "$dir" ]; then
        newest=$(ls -t "$dir"/*.jsonl 2>/dev/null | head -1)
        if [ -n "$newest" ] && [ -f "$newest" ]; then
          sz=$(wc -c < "$newest" 2>/dev/null | tr -d ' ')
          total=$((total + ${sz:-0})); got=1
        fi
      fi
      ;;
  esac
  # The driver log is a second arm, and the ONLY one for codex/pi — both stream
  # into it, while `claude -p` buffers to the end, which is exactly why claude
  # needs the transcript above.
  if [ -n "${LOG:-}" ] && [ -f "$LOG" ]; then
    sz=$(wc -c < "$LOG" 2>/dev/null | tr -d ' ')
    total=$((total + ${sz:-0})); got=1
  fi
  [ "$got" -eq 1 ] || return 1
  echo "$total"
}
# _mc_stalled PID → true only when FOUR arms agree, sample over sample, that the
# tree has made no progress: (1) a descendant alive ≥ STALL_CHILD_AGE (the wedged
# tool call of iteration 13's `until COND; do sleep 30; done`), (2) the progress
# counter above is unchanged since the previous sample, (3) the gate heartbeat is
# unchanged since the previous sample, and (4) the tree is under STALL_CPU_PCT.
# Any one arm showing movement resets the caller's hit counter, so live work is
# never killed — and unlike the CPU-only predecessor, that is now a property the
# suite can kill a mutant on rather than a claim in a comment.
#
# Fails OPEN: with no progress instrument we cannot tell a wedge from live work,
# so we refuse to guess and say so in the log; HARD_TIMEOUT still bounds the slot.
_mc_stalled() {
  local root="$1" pids p secs cpu long=0 prog hb
  pids=$(_mc_descendants "$root")
  for p in $pids; do
    [ "$p" = "$root" ] && continue
    secs=$(_mc_etime_secs "$(ps -o etime= -p "$p" 2>/dev/null)")
    [ "${secs:-0}" -ge "${STALL_CHILD_AGE:-2400}" ] && { long=1; break; }
  done
  [ "$long" -eq 1 ] || { _MC_STALL_WHY="no-long-child"; return 1; }

  if ! prog=$(_mc_progress_bytes); then
    if [ "${_MC_STALL_NOPROG_LOGGED:-0}" != "1" ]; then
      log "WARNING: stall watchdog has NO progress instrument (provider=${CONTROLLER_PROVIDER:-claude} pwd=$PWD) — early kill DISABLED for this attempt; HARD_TIMEOUT still applies"
      _MC_STALL_NOPROG_LOGGED=1
    fi
    _MC_STALL_WHY="no-progress-instrument"
    return 1
  fi
  hb=$(wc -c < "${_mc_heartbeat:-${AILANG_STATE_DIR:-$HOME/.ailang/state}/mission-${MISSION_NAME:-none}-heartbeat}" 2>/dev/null | tr -d ' '); hb="${hb:-0}"

  # A first sample can prove nothing — seed the baseline and report live.
  if [ -z "${_MC_PROG_PREV:-}" ]; then
    _MC_PROG_PREV="$prog"; _MC_HB_PREV="$hb"; _MC_STALL_WHY="seeding"; return 1
  fi
  if [ "$prog" != "$_MC_PROG_PREV" ] || [ "$hb" != "${_MC_HB_PREV:-}" ]; then
    _MC_STALL_WHY="progress prog=${_MC_PROG_PREV}->${prog} hb=${_MC_HB_PREV:-}->${hb}"
    _MC_PROG_PREV="$prog"; _MC_HB_PREV="$hb"
    return 1
  fi
  cpu=$(ps -o %cpu= -p "$(echo $pids | tr ' ' ',')" 2>/dev/null | awk '{s+=$1} END{printf "%d", s+0}')
  [ "${cpu:-0}" -lt "${STALL_CPU_PCT:-2}" ] || { _MC_STALL_WHY="cpu=$cpu"; return 1; }
  _MC_STALL_WHY="flat prog=$prog hb=$hb cpu=$cpu"
  return 0
}

# ----------------------------------------------------------------------------

# --- model selection (fleet Phase A) -----------------------------------------
# ASTRA RUNG (2026-09-05, Mark attended). Opus still leads; astra sits BETWEEN the
# Anthropic rungs so a dry bucket reaches a working controller after one failed probe
# instead of three. Matches the shared driver's ladder, which world had never received.
PREFS="${MISSION_MODEL_PREFS:-claude-opus-5,codex:gpt-6-astra,claude-opus-4-8,claude-fable-5-1}"
CONTROLLER_FALLBACK="${MISSION_CONTROLLER_FALLBACK:-codex:gpt-5.6-sol,codex:gpt-6-astra,pi:ollama/glm-5.3:cloud}"
QUOTA_SIG="usage limit|rate.?limit|quota|exceeded|too many requests|weekly limit"
PROBE_TIMEOUT="${MISSION_PROBE_TIMEOUT:-120}"   # per-probe wall-clock cap, seconds

# _mc_bounded SECONDS CMD... — run CMD with a hard wall-clock cap.
# rc = CMD's rc, or 124 on expiry (mirrors GNU `timeout`, which this rig does not have).
# Combined stdout+stderr lands in $MC_BOUNDED_OUT.
#
# Why (2026-07-27): a model probe is a network call to a third party and CAN hang. Observed that
# day: `codex exec --model <unknown-model>` ran past 180s with no output. Both probes below used
# to be unbounded command substitutions, so one hung probe would burn the whole 6h fire before the
# driver's HARD_TIMEOUT reclaimed it — the exact failure class as mission-control Standing rule 6
# ("every wait is bounded"), which the loop enforces on itself but the driver did not.
_mc_bounded() {
  local secs="$1"; shift
  local out_f rc deadline pid
  out_f=$(mktemp -t mc_bounded) || { MC_BOUNDED_OUT=""; return 125; }
  ( exec "$@" ) >"$out_f" 2>&1 &
  pid=$!
  deadline=$(( $(date +%s) + secs ))
  while kill -0 "$pid" 2>/dev/null; do
    if [ "$(date +%s)" -ge "$deadline" ]; then
      kill "$pid" 2>/dev/null; sleep 2; kill -9 "$pid" 2>/dev/null
      MC_BOUNDED_OUT="$(cat "$out_f" 2>/dev/null)"; rm -f "$out_f"
      return 124
    fi
    sleep 2
  done
  wait "$pid"; rc=$?
  MC_BOUNDED_OUT="$(cat "$out_f" 2>/dev/null)"; rm -f "$out_f"
  return "$rc"
}

# _mc_probe MODEL → 0 usable | 1 quota-limited | 2 unusable (auth/transient×2/timeout×2)
_mc_probe() {
  local m="$1" out rc
  _mc_bounded "$PROBE_TIMEOUT" claude -p 'reply with exactly: ok' --model "$m"; rc=$?
  out="$MC_BOUNDED_OUT"
  [ "$rc" -eq 0 ] && return 0
  [ "$rc" -eq 124 ] && log "model $m probe timed out after ${PROBE_TIMEOUT}s"
  if printf '%s' "$out" | grep -qiE "$QUOTA_SIG"; then return 1; fi
  # transient? retry once
  sleep 5
  _mc_bounded "$PROBE_TIMEOUT" claude -p 'reply with exactly: ok' --model "$m"; rc=$?
  out="$MC_BOUNDED_OUT"
  [ "$rc" -eq 0 ] && return 0
  [ "$rc" -eq 124 ] && log "model $m probe timed out after ${PROBE_TIMEOUT}s (retry)"
  printf '%s' "$out" | grep -qiE "$QUOTA_SIG" && return 1
  return 2
}

# _mc_probe_codex MODEL → 0 usable | non-zero unusable. The OpenAI API key is
# stripped above, so a pass proves the ChatGPT-subscription OAuth lane works.
_mc_probe_codex() {
  local m="$1" rc
  _mc_bounded "$PROBE_TIMEOUT" codex exec --skip-git-repo-check --model "$m" 'reply with exactly: ok'
  rc=$?
  [ "$rc" -eq 124 ] && log "controller fallback codex:$m probe timed out after ${PROBE_TIMEOUT}s"
  [ "$rc" -ne 0 ] && log "controller fallback codex:$m probe failed (rc=$rc): $(printf '%s' "$MC_BOUNDED_OUT" | tail -3 | tr '\n' ' ')"
  return "$rc"
}

_mc_uptime_secs() {
  local b now
  b=$(sysctl -n kern.boottime 2>/dev/null | awk -F'sec *= *' 'NF>1 {printf "%d", $2+0; exit}')
  [ -n "$b" ] && [ "$b" -gt 0 ] 2>/dev/null || return 1
  now=$(date +%s)
  echo $(( now - b ))
}

# _mc_boot_offset NAME — seconds this mission waits out of a boot stampede.
# Spacing is 7 minutes, which is longer than a controller's probe+startup
# preamble (the v1 slot that burned 240s on opus probes is the worst measured),
# so each mission's spawn burst has finished before the next one begins. v1 is 0
# because it has the shortest interval (90m) and the deepest ladder — it is the
# loop we least want to delay. Unknown missions get 0: a new mission must be
# added here deliberately, and defaulting it into someone else's slot would be
# worse than leaving it at boot.
_mc_boot_offset() {
  case "${1:-}" in
    v1)     echo 0    ;;
    world)  echo 420  ;;
    docs)   echo 840  ;;
    motoko) echo 1260 ;;
    *)      echo 0    ;;
  esac
}

# _mc_mem_snapshot — echo "AVAIL_MB COMPRESSED_MB", rc=1 if vm_stat cannot answer.
#
# AVAIL = free + inactive + speculative + purgeable. NOT `free` alone: at the
# 09-05 09:23 event free was 4030 pages (66 MB) while inactive held 506169
# (7.7 GB) that the pager could still reclaim, so a free-only threshold would
# have to sit absurdly low to avoid firing constantly. Measured separation with
# this expression: ~7.8 GB at each of the three OOM events, ~104 GB on a healthy
# idle box — two orders of magnitude apart, so the threshold is not delicate.
#
# The compressor arm is the second signal, for the case where `inactive` still
# looks healthy but the machine is already paging hard: 66 GB compressed (holding
# 131 GB) at every event, 0 on a fresh boot.
#
# NOT memoryPressure: the kernel's own flag read `false` throughout the 09-03
# panic. .claude/rules/local-models.md carries that trap.
_mc_mem_snapshot() {
  command -v vm_stat >/dev/null 2>&1 || return 1
  vm_stat 2>/dev/null | awk '
    /page size of/ { for (i=1; i<=NF; i++) if ($i == "of") { ps=$(i+1); break } }
    /^Pages free:/                    { free=$3 }
    /^Pages inactive:/                { inact=$3 }
    /^Pages speculative:/             { spec=$3 }
    /^Pages purgeable:/               { purge=$3 }
    /^Pages occupied by compressor:/  { comp=$5 }
    END {
      if (ps == "" || free == "") exit 1
      gsub(/\./, "", free); gsub(/\./, "", inact); gsub(/\./, "", spec)
      gsub(/\./, "", purge); gsub(/\./, "", comp)
      printf "%d %d\n", (free+inact+spec+purge) * ps / 1048576, comp * ps / 1048576
    }'
}

# _mc_mem_ok AVAIL_MB COMPRESSED_MB — 0 = there is room to start an iteration.
# Thresholds are STARTING VALUES, not measured ones: nobody has profiled an
# iteration's peak footprint. They are chosen to sit far from both observed
# states (refuse at 7.8 GB avail / 66 GB compressed, pass at 104 GB / 0) and are
# logged with the live numbers on every fire so the log tells us the real values.
_mc_mem_ok() {
  [ "${1:-0}" -ge "${MEM_MIN_AVAIL_MB:-16384}" ] || return 1
  [ "${2:-0}" -le "${MEM_MAX_COMP_MB:-49152}" ]  || return 1
  return 0
}

_mc_set_controller() {
  local requested="$1"
  MODEL_WHY="$2"
  case "$requested" in
    codex:*) CONTROLLER_PROVIDER=codex; MODEL="${requested#codex:}"; MISSION_ANTHROPIC_AVAILABLE=0 ;;
    pi:*) CONTROLLER_PROVIDER=pi; MODEL="${requested#pi:}"; MISSION_ANTHROPIC_AVAILABLE=0 ;;
    claude:*) CONTROLLER_PROVIDER=claude; MODEL="${requested#claude:}"; MISSION_ANTHROPIC_AVAILABLE=1 ;;
    *) CONTROLLER_PROVIDER=claude; MODEL="$requested"; MISSION_ANTHROPIC_AVAILABLE=1 ;;
  esac
  CONTROLLER_ID="${CONTROLLER_PROVIDER}:${MODEL}"
  export CONTROLLER_PROVIDER CONTROLLER_ID MODEL MODEL_WHY MISSION_ANTHROPIC_AVAILABLE
}

select_model() {
  # 1. absolute pin
  if [ -n "${MISSION_MODEL:-}" ]; then _mc_set_controller "$MISSION_MODEL" "env pin"; return 0; fi
  # 2. override file pin (optional expiry epoch)
  if [ -f "$OVERRIDE_FILE" ]; then
    local ov_model ov_until now
    read -r ov_model ov_until < "$OVERRIDE_FILE" 2>/dev/null || true
    now=$(date +%s)
    if [ -n "${ov_until:-}" ] && [ "$now" -ge "${ov_until:-0}" ]; then
      rm -f "$OVERRIDE_FILE"
      log "model override expired — resuming preference probing"
    elif [ -n "${ov_model:-}" ]; then
      _mc_set_controller "$ov_model" "override file"; return 0
    fi
  fi
  # 3. ordered preference probing.
  #
  # PROVIDER-DISPATCHED since 2026-09-05 (Mark, attended: "put astra ahead of each
  # fable instance, that falls back to fable"). This list used to be Anthropic-only —
  # every entry went to `_mc_probe`, the claude CLI probe — so a non-Anthropic model
  # could ONLY be expressed in CONTROLLER_FALLBACK, which is reached after EVERY
  # Anthropic candidate has failed. There was therefore no way to say
  # "opus, then astra, then fable": a codex entry could sit before opus (never) or
  # after fable (too late), but not BETWEEN them. That ordering is the whole ask.
  #
  # Bare and `claude:`-prefixed entries keep the exact 0/1/2 quota-vs-unusable
  # semantics they had; only the dispatch is new. _mc_set_controller already parses
  # every prefix, so a matched entry needs no special-casing beyond its probe.
  local m why rcode
  for m in $(printf '%s' "$PREFS" | tr ',' ' '); do
    case "$m" in
      codex:*)
        if _mc_probe_codex "${m#codex:}"; then
          _mc_set_controller "$m" "probe ok"; return 0
        fi
        log "controller preference $m unusable — falling through"
        ;;
      pi:*)
        _mc_bounded "$PROBE_TIMEOUT" pi --mode json --no-session --no-tools --model "${m#pi:}" -p 'reply with exactly: ok'
        rcode=$?
        if [ "$rcode" -eq 0 ]; then _mc_set_controller "$m" "probe ok"; return 0; fi
        log "controller preference $m probe failed (rc=$rcode within ${PROBE_TIMEOUT}s) — falling through"
        ;;
      *)
        _mc_probe "$m"; rcode=$?
        case "$rcode" in
          0) _mc_set_controller "$m" "probe ok"; return 0 ;;
          1) log "model $m quota-limited — falling through" ;;
          2) log "model $m unusable (auth/transient) — falling through" ;;
        esac
        ;;
    esac
  done
  # 4. cross-provider fallback CHAIN, walked in order (Mark 2026-08-31 — see the
  # CONTROLLER_FALLBACK comment above). Every rung is probe-gated; an unsupported
  # entry is skipped loudly rather than aborting the walk, so one typo cannot
  # disable the rungs behind it.
  log "all Anthropic controller candidates unavailable — walking fallback chain ($CONTROLLER_FALLBACK)"
  local fb
  for fb in $(printf '%s' "$CONTROLLER_FALLBACK" | tr ',' ' '); do
    case "$fb" in
      codex:*)
        m="${fb#codex:}"
        if _mc_probe_codex "$m"; then
          _mc_set_controller "$fb" "Anthropic unavailable; subscription fallback"
          return 0
        fi
        ;;
      pi:*)
        m="${fb#pi:}"
        # Same probe shape as the role-lane pi loop: --no-tools keeps it ~1 reply
        # token, --no-session avoids polluting ~/.pi/sessions; rc is the verdict.
        # rc captured explicitly: after `if cmd; then...fi` falls through, $? is the
        # IF's status (0), not cmd's — logging it would report every failure as rc=0.
        _mc_bounded "$PROBE_TIMEOUT" pi --mode json --no-session --no-tools --model "$m" -p 'reply with exactly: ok'
        rcode=$?
        if [ "$rcode" -eq 0 ]; then
          _mc_set_controller "$fb" "Anthropic+codex unavailable; pi fallback rung"
          return 0
        fi
        log "pi controller rung '$m' probe failed (rc=$rcode within ${PROBE_TIMEOUT}s) — falling through"
        ;;
      *) log "unsupported CONTROLLER_FALLBACK entry '$fb' (expected codex:<model> or pi:<model>) — skipping" ;;
    esac
  done
  return 1
}
# ----------------------------------------------------------------------------

HARD_TIMEOUT="${MISSION_TIMEOUT:-21600}"   # 6h wall-clock kill per iteration
# Stall watchdog (2026-07-12): a wedged unbounded poll loop (iteration 13's
# `until COND; do sleep 30; done`) otherwise burns the whole 6h slot before
# HARD_TIMEOUT. Kill early once the tree has made NO PROGRESS — see _mc_stalled
# for the four arms, and for why the original CPU-only test killed live work. The
# grace and the child-age gate sit past the skill's 30-min bounded-wait cap so a
# COMPLIANT wait can never trip it. All env-overridable.
STALL_GRACE="${MISSION_STALL_GRACE:-2400}"       # 40m before the first check
STALL_CHILD_AGE="${MISSION_STALL_CHILD_AGE:-2400}" # a descendant alive ≥40m = wedged
STALL_INTERVAL="${MISSION_STALL_INTERVAL:-120}"  # 2m between samples
# 5 × 2m = 10 minutes of PROVEN no-progress before a kill. Was 3, on an arm that
# could not see progress at all; with real arms the window is worth widening,
# because what a wrong kill destroys is a whole iteration.
STALL_SAMPLES="${MISSION_STALL_SAMPLES:-5}"      # consecutive no-progress hits → kill
STALL_CPU_PCT="${MISSION_STALL_CPU_PCT:-2}"      # tree %cpu floor — the weakest of the four arms
export STALL_CHILD_AGE

# TRANSIENT-RETRY (2026-07-14): Anthropic capacity is flaky some evenings —
# `claude -p` does its own internal retries then exits rc=1 on a persistent
# "API Error: Overloaded" / dropped socket, losing the whole iteration (2 lost
# 2026-07-14). Retry the run on a TIGHTLY-ANCHORED transient signature (claude's
# own "API Error:" emissions + socket-closed), with backoff. NEVER retried:
# watchdog kills (rc 143/137 = deliberate stall/timeout), quota/429 (that's
# Phase A's start-probe fall-through job, not a same-model retry), or any other
# genuine rc. Signature is anchored so an unrelated "503" in a test's output
# (e.g. the httpbin fixture) cannot trigger a false retry.
TRANSIENT_RETRIES="${MISSION_TRANSIENT_RETRIES:-3}"   # total attempts incl. the first
TRANSIENT_BACKOFF="${MISSION_TRANSIENT_BACKOFF:-45}"  # base seconds, ×attempt (45s,90s)
TRANSIENT_SIG="API Error: Overloaded|socket connection was closed|overloaded_error|API Error: 5[0-9][0-9]|API Error: Internal|API Error: Connection|API Error: Request timed out"

# PER-ROLE MODEL ROUTING (2026-07-15, m-mission-agentic-provider-routing M1): the charter's routing
# table was never enforced — every inner role ran on the controller's single session --model, so with
# the driver on Fable 100% of each iteration billed Fable (memory:
# project-mission-routing-table-never-enforced). Fix: the controller session keeps $MODEL; the HEAVY
# roles are spawned by mission-control Gate 3 as model-PINNED sub-agents that read these env vars.
# Defaults track the charter routing table; M3 will A/B the planner down-tier — keep it at the proven
# Opus until there's evidence. Cross-provider AGENT executors (codex/motoko) ride the same env once
# fleet Phase C wires them into the spawn (a value like "codex:gpt-5.6" is resolved by the skill).
# 2026-07-16 (Mark): Fable = high-cognition ROLES only. The controller session is opus-first (see
# PREFS above); Fable bills exactly two BOUNDED pinned sub-agents per iteration: the designer
# (deep spec synthesis, fired only when a new doc is needed) and the evaluator (adversarial judge,
# ≠ the opus executor → generator≠judge holds).
# NB: these are in-session Agent/Task-tool model ALIASES (opus|fable|sonnet|haiku) — NOT the full
# IDs (claude-opus-4-8) the driver's own `claude -p --model` flag takes. Two different interfaces:
# the controller session is launched with a full ID; the sub-agents it spawns are pinned by alias.
# A "provider:model" value (e.g. codex:gpt-5.6-sol) instead signals cross-provider agent routing via
# provider_executor (fleet Phase C), which the skill resolves — not the Agent tool.
# MISSION_EXECUTOR_MODEL specifically accepts EITHER form: an Agent alias (opus) OR a
# provider:model value (codex:gpt-5.6-sol — the DEFAULT since 2026-07-27) which the mission-control
# Gate-3 recipe routes to a bounded `codex exec` run in the sprint worktree (M1b). Default flipped
# to codex 2026-07-27 (Mark, quota relief): codex now authenticates via ChatGPT SUBSCRIPTION
# (auth.json auth_mode=chatgpt; metered API-key backup at ~/.codex/auth.json.apikey.bak), so it is
# a quota lane, not metered $ — the old "never bills metered $ by accident" rationale is moot.
# The pre-flight probe below falls back to opus for the fire when the codex bucket is spent.
# Generator≠judge holds: evaluator stays sonnet regardless of executor lane.
# Weekly rolling bookkeeping issue (2026-07-16, Mark): the issue number lives in a state file so
# the skill's Monday-07:00 rotation (aligned to the quota reset) moves threads without a driver
# edit. Precedence: env pin > state file > 329 (the original thread). Exported so the skill's
# gh snippets see the same number the driver reports to.
MISSION_GH_ISSUE="${MISSION_GH_ISSUE:-$(head -1 "$GH_ISSUE_FILE" 2>/dev/null)}"
[ -z "${MISSION_GH_ISSUE:-}" ] && [ "$MISSION_NAME" = "v1" ] && MISSION_GH_ISSUE=329
export MISSION_GH_ISSUE

# designer default is the claude-CLI lane (claude:<full-id>), NOT the bare "fable" alias: the
# Agent tool pins only sonnet|opus|haiku (F1, iteration 31), so under an opus-first controller a
# bare "fable" would silently fall back to opus. claude:claude-fable-5 = a REAL bounded Fable run.
export MISSION_DESIGNER_MODEL="${MISSION_DESIGNER_MODEL:-claude:claude-fable-5}"
# Designer chain (2026-09-05) — world had none. Same rungs and order as the skill's
# fable -> astra -> deepseek rotation, so a PINNED designer degrades the way a rotating one does.
export MISSION_DESIGNER_FALLBACK="${MISSION_DESIGNER_FALLBACK:-codex:gpt-6-astra,pi:ollama/deepseek-v4-flash:0731-cloud,pi:openrouter/deepseek/deepseek-v4-flash-0731}"
# Per-iteration METERED-spend ceiling (2026-07-18, Mark: "make sure costs don't go crazy"):
# the sum of all metered-API spend (codex $ + gemini $) within ONE iteration must stay under
# this. Enforced by the skill's Gate-3 metered ledger; quota-bucket (subscription) spend is
# NOT counted — this caps dollars, not tokens.
export MISSION_METERED_BUDGET_USD="${MISSION_METERED_BUDGET_USD:-5}"
export MISSION_PLANNER_MODEL="${MISSION_PLANNER_MODEL:-opus}"
export MISSION_EXECUTOR_MODEL="${MISSION_EXECUTOR_MODEL:-codex:gpt-5.6-sol}"
# EXECUTOR FALLBACK CHAIN — ailang#611 (2026-08-11). Ratified: codex default,
# deepseek the replacement when codex is dry, opus the last resort. Kept identical
# to the V1 driver so both missions route the same way.
export MISSION_EXECUTOR_FALLBACK="${MISSION_EXECUTOR_FALLBACK:-pi:ollama/deepseek-v4-flash:0731-cloud,pi:openrouter/deepseek/deepseek-v4-flash-0731:floor}"
# The planner's fallback used to be `opus` — the model it was already pinned to, so on a
# dry Anthropic bucket it fell back onto the failure. codex first, then the flat-rate pi rung.
export MISSION_PLANNER_FALLBACK="${MISSION_PLANNER_FALLBACK:-codex:gpt-5.6-sol,pi:ollama/kimi-k3:cloud,pi:openrouter/moonshotai/kimi-k3}"
export MISSION_PLANNER_ANTHROPIC_FALLBACK="${MISSION_PLANNER_ANTHROPIC_FALLBACK:-codex:gpt-5.6-sol}"
# evaluator default = sonnet (2026-07-16, Mark directive on #399: "default can be gemini (if able
# to git clone the codebase etc)? otherwise sonnet-5"). gemini managed_agents is NOT viable as the
# evaluator today — VERIFIED iteration 38: (1) architecturally the request body carries only
# Directive+SystemPrompt over a server-side CapRemoteSandbox (managed_agents.go:164), so it cannot
# see the sprint's UNCOMMITTED worktree changes nor re-run local tests — at most it could clone the
# public origin/dev, which lacks the changes; (2) the backend live-timed-out (http2 timeout, same
# class as iters 36-37). So the ladder resolves to sonnet-5: pinnable via the Agent tool (fable is
# not — F1), distinct from the opus executor (generator≠judge holds), cheap, behavioral (re-runs
# tests locally). This also RETIRES the per-iteration fable→sonnet re-route (iters 31/36) into a
# standing default. gemini-as-evaluator is a queued follow-up (diff-bridge + backend reliability).
export MISSION_EVALUATOR_MODEL="${MISSION_EVALUATOR_MODEL:-sonnet}"
# Evaluator chain (2026-09-05) — world had none at all: sonnet or nothing, which wedged
# Gate 4 on every Anthropic outage. codex is LAST on purpose: the executor is
# codex:gpt-5.6-sol, so an astra judge is the vendor-level generator==judge collision.
# minimax (independent vendor) takes the first rung; astra stands only between "judged by
# the executor's own vendor" and no judge at all.
export MISSION_EVALUATOR_FALLBACK="${MISSION_EVALUATOR_FALLBACK:-pi:ollama/minimax-m3:cloud,pi:openrouter/minimax/minimax-m3,codex:gpt-6-astra}"

# ─── ROLE-GENERIC LANE PRE-FLIGHT — PORTED FROM THE SHARED AILANG DRIVER 2026-09-05 ───
# World is a SEPARATE GitHub repo carrying its own copy of this driver and has no
# lib/pin-root.sh, so it never re-execs from the driver pin and has silently missed every
# routing fix made in sunholo-data/ailang. What stood here until now probed the EXECUTOR
# and nothing else: planner (opus), designer (fable) and evaluator (sonnet) were never
# probed by anything, so a dry Anthropic bucket killed three of five roles mid-iteration.
#
# Worse, the lane-degradation LEDGER below was written BECAUSE of this mission — the World
# mission spent five iterations (18/19/21/22) silently demoted from codex to opus, each
# mis-attributed to a spent quota — and world is the one mission that never received it.
#
# This block is a verbatim port. Keep it verbatim: it is the same code the shared driver
# runs, and the whole point of de-forking world later is that this stops being a copy.
# Codex-lane pre-flight, ROLE-GENERIC (m-planner-codex-lane): probe once per DISTINCT
# codex model, fall back per-role on ANY non-zero rc (#486: probe MUST carry --model;
# an unusable pin is exactly as fatal as spent quota). Export AFTER fallback so the
# EXPORTED env — what the routing-evidence row reports — stays honest.
# BASH 3.2 (L19): ':'-delimited string sets, NOT associative arrays; no ${var,,}.
#
# The probe MUST carry --model (#486, 2026-07-27): without it codex exercises its DEFAULT model,
# so a pinned-but-unreachable model false-greens the lane. Live evidence that day: codex-cli
# 0.137.0 answered the model-less probe on gpt-5.5 (rc=0) while `--model gpt-5.6-sol` returned a
# 400 "requires a newer version of Codex" — the driver exported the codex pin as healthy and the
# failure only surfaced inside the skill's Gate-3 recipe, one silent fallback later.
#
# Fall back on ANY non-zero rc, not just quota signatures: an unusable model pin is exactly as
# fatal to the lane as a spent quota, and the old quota-only gate is what let #486 through. The
# skill's Gate-3 recipe re-probes and would fall back anyway; doing it here keeps the EXPORTED
# env honest, which is what the routing-evidence row reports.
_cx_probed=":"   # models probed this fire (dedupe: planner+executor share the default model)
_cx_failed=":"   # models whose probe failed
# LANE-DEGRADATION LEDGER (motoko mission iteration 0, 2026-08-12; Mark ratified the fix).
# Until now a lane demotion was `log`ged here and NOWHERE ELSE — none of this driver's four
# `gh issue comment` sites covers it, so the human channel saw nothing. That is exactly how the
# World mission spent FIVE iterations (18/19/21/22) silently demoted from codex to opus, each
# mis-attributed to a spent quota, before iter-23 found the real cause. A fallback visible only in
# a routing-evidence row written AFTER the fact is still a silent fallback (Critical Principle 2):
# by then the iteration has already run on the wrong lane.
# ROLE FALLBACK CHAINS (2026-08-26). MISSION_<ROLE>_FALLBACK may now be a
# COMMA-SEPARATED chain, walked left to right, with opus as the implicit tail:
#
#   codex -> pi:ollama/<m>:cloud -> pi:openrouter/<twin> -> opus
#            flat-rate             metered                 Anthropic
#
# The Ollama Cloud quota is a subscription with an UNPUBLISHED denominator
# (/api/usage reports consumption but no limit), so we cannot predict exhaustion
# — only survive it. The OpenRouter rung is the same weights on a metered route,
# so exhaustion degrades the ROUTE and not the model.
# bash 3.2 (L19/L21): plain string splitting, no arrays or ${x//}.
_chain_head() { printf '%s' "${1%%,*}"; }
_chain_tail() { case "$1" in *,*) printf '%s' "${1#*,}" ;; *) printf '' ;; esac; }

# Accumulate here; emit ONCE below, AFTER every early exit and BEFORE the iteration starts.
# bash 3.2 (L19/L21): no associative arrays — ';'-delimited "model=rc", newline-delimited ledger.
_lane_degraded=""   # newline-delimited markdown bullets, one per degraded role
_cx_rcmap=""        # "model=rc;" so the emit site names the probe's exit code, not just the lane
_pi_rcmap=""
# ─── PHASE ORDER CORRECTED + BOOT STAGGER / MEMORY GATE PORTED 2026-09-05 ───
# World ran its lane pre-flights BEFORE the kill switch and the overlap guard, so a
# DISABLED mission and a mission yielding to its own previous iteration both still
# paid a full round of probe tokens. Upstream orders these deliberately: kill switch,
# overlap, stagger, memory gate, THEN probes — a fire that is not going to run should
# spend nothing. Moving world's guards up is that same ordering.
#
# The stagger and the gate are ported verbatim. World is now on a KeepAlive schedule,
# which makes both load-bearing rather than nice-to-have: KeepAlive restarts the job
# the moment it exits, so without the gate a box already out of memory would get a new
# iteration added to it immediately instead of at the next interval.

# 1. Kill switch — the intended "off" state, exit silently.
if [ -f "$KILL_SWITCH" ]; then
  log "kill switch present ($KILL_SWITCH) — skip"; exit 0
fi

# 1b. ONE iteration at a time (2026-07-10, continuous mode): two concurrent
#     controllers would stomp the charter/log in the main tree and could pick
#     the same queue item. If one is still running, yield this slot.
#     PIDFILE-based (2026-07-16): the old `pgrep -f "claude -p Run one mission"`
#     matched ANY process whose cmdline contained the phrase — including a
#     human's monitoring shell (`pgrep -f "claude -p Run one mission"` itself!),
#     which made a kickstarted fire yield against its own observer. A pidfile
#     + liveness check cannot false-positive.
if [ -f "$PIDFILE" ]; then
  oldpid=$(head -1 "$PIDFILE" 2>/dev/null)
  if [ -n "$oldpid" ] && kill -0 "$oldpid" 2>/dev/null; then
    log "previous iteration still running (pid $oldpid) — yield (next interval retries)"; exit 0
  fi
  rm -f "$PIDFILE"   # stale pidfile from a crashed/killed run — proceed
fi

# 3b. BOOT STAGGER (2026-09-05). See _mc_boot_offset for the measurement. Placed
#     AFTER the kill switch, the overlap yield and the dry run — a disabled
#     mission, a yielding one and a wiring check must all still be instant — and
#     BEFORE the probes, so a staggered fire spends zero tokens while it waits.
#     Holding the job for the offset is safe: launchd will not start a second
#     copy of a StartInterval job while the first is still running, so the worst
#     case is one skipped slot on a mission whose interval is 90m or longer.
BOOT_WINDOW="${MISSION_BOOT_WINDOW:-900}"
_up=$(_mc_uptime_secs || echo "")
_off=$(_mc_boot_offset "$MISSION_NAME")
if [ -z "$_up" ]; then
  log "boot stagger: kern.boottime unreadable — stagger SKIPPED this fire (not silent: gate disabled, not passed)"
elif [ "$_up" -lt "$BOOT_WINDOW" ] && [ "$_off" -gt 0 ]; then
  log "boot stagger: up ${_up}s (< ${BOOT_WINDOW}s window) — waiting ${_off}s so $MISSION_NAME does not start alongside the other missions"
  sleep "$_off"
fi

# 3c. MEMORY GATE (2026-09-05). Refuses to ADD an iteration to a box that is
#     already out of memory. Waits rather than skipping outright, because a
#     skipped slot costs motoko 13h — a transient spike should delay a fire, not
#     cancel it. On expiry it yields exactly like the overlap guard above:
#     exit 0, no notification, mission-recovery and the next interval retry.
MEM_MIN_AVAIL_MB=$(( ${MISSION_MIN_AVAIL_GB:-16} * 1024 ))
MEM_MAX_COMP_MB=$(( ${MISSION_MAX_COMPRESSED_GB:-48} * 1024 ))
MEM_WAIT="${MISSION_MEM_WAIT:-600}"
MEM_POLL="${MISSION_MEM_POLL:-60}"
_mem_deadline=$(( $(date +%s) + MEM_WAIT ))
while :; do
  _snap=$(_mc_mem_snapshot || echo "")
  if [ -z "$_snap" ]; then
    # Fail OPEN, loudly. vm_stat is macOS-only; refusing on a box that cannot
    # answer would wedge every mission rather than protect anything, and the
    # gate is an admission control, not a correctness guarantee.
    log "memory gate: vm_stat unavailable — gate DISABLED for this fire"
    break
  fi
  _avail=${_snap%% *}; _comp=${_snap##* }
  if _mc_mem_ok "$_avail" "$_comp"; then
    log "memory gate: ok (avail=${_avail}MB >= ${MEM_MIN_AVAIL_MB}MB, compressed=${_comp}MB <= ${MEM_MAX_COMP_MB}MB)"
    break
  fi
  if [ "$(date +%s)" -ge "$_mem_deadline" ]; then
    log "memory gate: STILL SHORT after ${MEM_WAIT}s (avail=${_avail}MB, compressed=${_comp}MB) — yield (next interval retries)"
    exit 0
  fi
  log "memory gate: low memory (avail=${_avail}MB, compressed=${_comp}MB) — waiting ${MEM_POLL}s"
  sleep "$MEM_POLL"
done

# ANTHROPIC-LANE PRE-FLIGHT, ROLE-GENERIC (2026-09-05, Mark attended: "give me rotation
# and fallbacks to codex so we can keep going on the missions ... before we usually ran out
# of quota on friday"). THIS is the rung that was missing, and its absence is why a drought
# did not pause the fleet so much as half-run it.
#
# The two loops below only ever look at `codex:*` and `pi:*` values, so a role pinned to an
# ANTHROPIC model — the `sonnet` evaluator, a `claude:*` designer — was never probed by
# anything. On a dry bucket the controller would fall to its own codex rung and the
# iteration would START, then die at the first Anthropic-only gate. Continues the
# NO-SINGLE-PROVIDER-ROLE directive (Mark 2026-08-26) from declaring chains to actually
# walking them: MISSION_<ROLE>_FALLBACK existed for the evaluator since that day, and
# nothing on any code path read it (the skill greps zero MISSION_*_FALLBACK, verified).
#
# Placed BEFORE the codex loop on purpose: a role handed from anthropic to `codex:gpt-6-astra`
# is then probed by that loop, and on to pi by the next — one chain across three providers,
# not three disconnected pre-flights.
#
# Cost is one `claude -p` per DISTINCT anthropic model per fire (deduped, same as the codex
# and pi loops), and `_mc_probe` already distinguishes quota-limited (rc=1) from unusable
# (rc=2) — a distinction the degradation ledger reports, because "Friday" and "broken pin"
# have very different resume conditions.
# BASH 3.2 (L19/L21): ':'-delimited string sets, no associative arrays, no ${var,,}.
_an_probed=":"   # anthropic models probed this fire
_an_failed=":"   # anthropic models whose probe failed
_an_rcmap=""     # "model=rc;" so the emit site names the probe's exit code
for role in DESIGNER PLANNER EXECUTOR EVALUATOR; do
  var="MISSION_${role}_MODEL"; val="${!var:-}"
  # Only anthropic-shaped values: a bare Agent alias (sonnet/opus) or an explicit claude: pin.
  case "$val" in
    ""|codex:*|pi:*) continue ;;
  esac
  an_model="${val#claude:}"
  case "$_an_probed" in *":${an_model}:"*) : ;; *)   # not yet probed
    _an_probed="${_an_probed}${an_model}:"
    _mc_probe "$an_model"; an_rc=$?
    if [ "$an_rc" -ne 0 ]; then
      _an_failed="${_an_failed}${an_model}:"
      if [ "$an_rc" -eq 1 ]; then an_why="quota-limited"; else an_why="unusable (rc=$an_rc)"; fi
      log "anthropic model '$an_model' $an_why"
      _an_rcmap="${_an_rcmap}${an_model}=${an_rc};"
    fi
  ;; esac
  case "$_an_failed" in *":${an_model}:"*)
    role_lc=$(printf '%s' "$role" | tr 'A-Z' 'a-z')   # ${role,,} is bash-4.0-only (L21)
    fbvar="MISSION_${role}_FALLBACK"; _chain="${!fbvar:-}"
    _an_rc_for=$(printf '%s' "$_an_rcmap" | tr ';' '\n' | grep "^${an_model}=" | head -1 | cut -d= -f2)
    [ -n "$_an_rc_for" ] || _an_rc_for="unknown"
    # NO implicit opus tail here, unlike the codex loop. Opus IS anthropic, so on the exact
    # failure this loop exists for it is the one destination guaranteed to be dry too —
    # handing to it would launder a drought into a second failure one gate later. A role
    # with no chain keeps its pin and is reported, which the skill can see and act on.
    if [ -z "$_chain" ]; then
      log "anthropic ${role_lc} lane '$an_model' $an_why and NO fallback chain configured — pin kept, reported to the ledger"
      _lane_degraded="${_lane_degraded}
- \`${role_lc}\`: **anthropic** lane \`${an_model}\` unusable (probe rc=\`${_an_rc_for}\`) → NO fallback chain, pin kept"
      continue
    fi
    fb=$(_chain_head "$_chain")
    remvar="MISSION_${role}_CHAIN_REMAINING"
    printf -v "$remvar" '%s' "$(_chain_tail "$_chain")"; export "$remvar"
    log "anthropic ${role_lc} lane -> falling back to '$fb' for this fire (model '$an_model', $an_why)"
    _lane_degraded="${_lane_degraded}
- \`${role_lc}\`: **anthropic** lane \`${an_model}\` unusable (probe rc=\`${_an_rc_for}\` — ${an_why}) → handed to \`${fb}\`"
    printf -v "$var" '%s' "$fb"; export "$var"
  ;; esac
done

for role in DESIGNER PLANNER EXECUTOR EVALUATOR; do
  var="MISSION_${role}_MODEL"; val="${!var}"
  case "$val" in codex:*)
    cx_model="${val#codex:}"
    case "$_cx_probed" in *":${cx_model}:"*) : ;; *)   # not yet probed
      _cx_probed="${_cx_probed}${cx_model}:"
      _mc_bounded "$PROBE_TIMEOUT" codex exec --skip-git-repo-check --model "$cx_model" 'reply with exactly: ok'
      cx_rc=$?; cx_out="$MC_BOUNDED_OUT"
      if [ "$cx_rc" -ne 0 ]; then
        _cx_failed="${_cx_failed}${cx_model}:"
        # why-classification happens ONCE, at probe time (timeout / quota-sig / other)
        if [ "$cx_rc" -eq 124 ]; then cx_why="probe timed out after ${PROBE_TIMEOUT}s"
        elif printf '%s' "$cx_out" | grep -qiE "$QUOTA_SIG"; then cx_why="quota-limited"
        else cx_why="probe failed (rc=$cx_rc)"; fi
        log "codex model '$cx_model' unusable: $cx_why"
        log "codex probe output: $(printf '%s' "$cx_out" | tail -3 | tr '\n' ' ')"
        _cx_rcmap="${_cx_rcmap}${cx_model}=${cx_rc};"
      fi
    ;; esac
    case "$_cx_failed" in *":${cx_model}:"*)
      role_lc=$(printf '%s' "$role" | tr 'A-Z' 'a-z')   # ${role,,} is bash-4.0-only (L21)
      # Hand off to the NEXT link, not straight to opus (#611). A `pi:*` value here
      # is probed by the pi loop below, which degrades to opus on its own failure —
      # that is what makes codex -> deepseek -> opus a real chain. `%s` rather than
      # a bare format string: the value is data, and a stray % would be a directive.
      # Continue an in-flight chain if the anthropic pre-flight already started one,
      # otherwise start this role's chain from the top. Without this an
      # anthropic->codex handoff whose codex rung then fails would RESTART at the
      # head of _FALLBACK — i.e. hand back to the codex rung that just failed.
      remvar="MISSION_${role}_CHAIN_REMAINING"; _chain="${!remvar:-}"
      [ -n "$_chain" ] || { fbvar="MISSION_${role}_FALLBACK"; _chain="${!fbvar:-opus}"; }
      fb=$(_chain_head "$_chain")
      # Remember what is left so the pi loop can advance instead of jumping to opus.
      remvar="MISSION_${role}_CHAIN_REMAINING"
      printf -v "$remvar" '%s' "$(_chain_tail "$_chain")"; export "$remvar"
      log "codex ${role_lc} lane -> falling back to '$fb' for this fire (model '$cx_model')"
      _cx_rc_for=$(printf '%s' "$_cx_rcmap" | tr ';' '\n' | grep "^${cx_model}=" | head -1 | cut -d= -f2)
      [ -n "$_cx_rc_for" ] || _cx_rc_for="unknown"
      _lane_degraded="${_lane_degraded}
- \`${role_lc}\`: **codex** lane \`${cx_model}\` unusable (probe rc=\`${_cx_rc_for}\`$([ "$_cx_rc_for" = "124" ] && printf ' — TIMEOUT after %ss' "$PROBE_TIMEOUT")) → handed to \`${fb}\`"
      printf -v "$var" '%s' "$fb"; export "$var"
    ;; esac
  ;; esac
done
# pi-lane pre-flight, ROLE-GENERIC (mirrors the codex loop above; added 2026-08-06,
# Mark: DeepSeek executor lane — trial record in models.yml pi-or-deepseek-v4-flash).
# Probe once per DISTINCT pi model, fall back per-role on ANY non-zero rc — an
# unusable pin is exactly as fatal as a spent bucket (#486). The OpenRouter key
# rides ~/.pi/agent/models.json (custom provider), not env, so this probe is
# headless-safe. --no-tools keeps it ~1 reply-token; --no-session avoids polluting
# ~/.pi/sessions. BASH 3.2 (L19): ':'-delimited string sets, NOT associative arrays.
_pi_probed=":"   # models probed this fire (dedupe: planner+executor could share one)
_pi_failed=":"   # models whose probe failed
for role in DESIGNER PLANNER EXECUTOR EVALUATOR; do
  var="MISSION_${role}_MODEL"
  # while-loop so a chain advance re-enters the probe for the NEW value; a
  # non-pi value (or a settled pi value) breaks out at the bottom.
  while :; do
  val="${!var}"
  case "$val" in pi:*)
    pi_model="${val#pi:}"
    case "$_pi_probed" in *":${pi_model}:"*) : ;; *)   # not yet probed
      _pi_probed="${_pi_probed}${pi_model}:"
      _mc_bounded "$PROBE_TIMEOUT" pi --mode json --no-session --no-tools --model "$pi_model" -p 'reply with exactly: ok'
      pi_rc=$?; pi_out="$MC_BOUNDED_OUT"
      if [ "$pi_rc" -ne 0 ]; then
        _pi_failed="${_pi_failed}${pi_model}:"
        if [ "$pi_rc" -eq 124 ]; then pi_why="probe timed out after ${PROBE_TIMEOUT}s"
        else pi_why="probe failed (rc=$pi_rc)"; fi
        log "pi model '$pi_model' unusable: $pi_why"
        log "pi probe output: $(printf '%s' "$pi_out" | tail -3 | tr '\n' ' ')"
        _pi_rcmap="${_pi_rcmap}${pi_model}=${pi_rc};"
      fi
    ;; esac
    case "$_pi_failed" in *":${pi_model}:"*)
      role_lc=$(printf '%s' "$role" | tr 'A-Z' 'a-z')   # ${role,,} is bash-4.0-only (L21)
      _pi_rc_for=$(printf '%s' "$_pi_rcmap" | tr ';' '\n' | grep "^${pi_model}=" | head -1 | cut -d= -f2)
      [ -n "$_pi_rc_for" ] || _pi_rc_for="unknown"
      # Advance along the chain rather than jumping to opus. The Ollama Cloud rung
      # can be exhausted by a quota whose denominator is unpublished, so the
      # OpenRouter twin — same weights, metered route — is the rung that keeps the
      # loop on the SAME model instead of degrading capability.
      remvar="MISSION_${role}_CHAIN_REMAINING"; _rem="${!remvar:-}"
      if [ -n "$_rem" ]; then
        _next=$(_chain_head "$_rem")
        printf -v "$remvar" '%s' "$(_chain_tail "$_rem")"; export "$remvar"
        log "pi ${role_lc} lane '$pi_model' unusable -> advancing to '$_next'"
        _lane_degraded="${_lane_degraded}
- \`${role_lc}\`: **pi** lane \`${pi_model}\` unusable (probe rc=\`${_pi_rc_for}\`$([ "$_pi_rc_for" = "124" ] && printf ' — TIMEOUT after %ss' "$PROBE_TIMEOUT")) → advanced to \`${_next}\`"
        printf -v "$var" '%s' "$_next"; export "$var"
        continue
      fi
      log "pi ${role_lc} lane -> falling back to opus for this fire (model '$pi_model')"
      _lane_degraded="${_lane_degraded}
- \`${role_lc}\`: **pi** lane \`${pi_model}\` unusable (probe rc=\`${_pi_rc_for}\`$([ "$_pi_rc_for" = "124" ] && printf ' — TIMEOUT after %ss' "$PROBE_TIMEOUT")) → handed to \`opus\` (end of chain)"
      printf -v "$var" 'opus'; export "$var"
    ;; esac
  ;; esac
  break
  done
done


# 3. Dry run — verify wiring without spending tokens (no probes fired).
if [ "${MISSION_DRY_RUN:-0}" = "1" ]; then
  # lanes= mirrors the shared driver, so a drought can be SIMULATED here rather than waited for:
  #   MISSION_PROFILE=world MISSION_DRY_RUN=1 MISSION_EVALUATOR_MODEL=claude-drought-sim
  if [ -n "$_lane_degraded" ]; then
    _dry_lanes="DEGRADED($(printf '%s' "$_lane_degraded" | grep -c '^- '))$(printf '%s' "$_lane_degraded" | tr '\n' ' ')"
  else
    _dry_lanes="ok"
  fi
  log "DRY RUN ok: mission=$MISSION_NAME repo-slug=$MISSION_REPO doc=$MISSION_DOC workdir=$REPO pidfile=$PIDFILE prefs=$PREFS timeout=${HARD_TIMEOUT}s | roles: designer=$MISSION_DESIGNER_MODEL planner=$MISSION_PLANNER_MODEL executor=$MISSION_EXECUTOR_MODEL evaluator=$MISSION_EVALUATOR_MODEL | lanes=$_dry_lanes"; exit 0
fi

# 4. Select the model (probe doubles as the subscription-auth check: API keys
#    are stripped above, so a passing probe proves keychain/token auth too).
if ! select_model; then
  log "NO usable controller in Anthropic prefs ($PREFS) or fallback ($CONTROLLER_FALLBACK). Refusing."
  ailang messages send controlplane \
    "mission-control refused to start: no usable controller in Anthropic prefs ($PREFS) or fallback ($CONTROLLER_FALLBACK). Per-model reasons are in the driver log. Zero tokens spent beyond probes." \
    --title "Mission iteration blocked: no usable model" --from "$MSG_FROM" 2>/dev/null
  [ -n "${MISSION_GH_ISSUE:-}" ] && gh issue comment "$MISSION_GH_ISSUE" --repo "$MISSION_REPO" \
    --body "⚠️ Mission iteration did not start: **no usable controller** in Anthropic preferences (\`$PREFS\`) or fallback (\`$CONTROLLER_FALLBACK\`). Per-model detail is in the driver log. Will retry next interval." 2>/dev/null
  exit 1
fi

# Announce model CHANGES on #329 (not every iteration — only transitions).
PREV_MODEL=$(cat "$LAST_MODEL_FILE" 2>/dev/null || true)
if [ -n "$CONTROLLER_ID" ] && [ "$CONTROLLER_ID" != "${PREV_MODEL:-}" ]; then
  printf '%s\n' "$CONTROLLER_ID" > "$LAST_MODEL_FILE"
  if [ -n "${PREV_MODEL:-}" ]; then
    log "controller model change: ${PREV_MODEL} → ${CONTROLLER_ID} (${MODEL_WHY})"
    [ -n "${MISSION_GH_ISSUE:-}" ] && gh issue comment "$MISSION_GH_ISSUE" --repo "$MISSION_REPO" \
      --body "🔁 Controller model: **${PREV_MODEL} → ${CONTROLLER_ID}** (${MODEL_WHY}) at $(date '+%F %H:%M %Z'). Automatic — Anthropic preference order \`$PREFS\`, then \`$CONTROLLER_FALLBACK\`; reverts when a higher-preference probe succeeds again." 2>/dev/null || true
  fi
fi

# ONE-SHOT executor override (2026-07-16, Mark: fleet live-fire tests). If armed, the file's value
# overrides MISSION_EXECUTOR_MODEL for exactly THIS iteration and is deleted on consumption.
# Placed after every early-exit (kill switch, overlap yield, no-model refusal) so a fire that does
# not actually run can never burn the shot. Arm with e.g.:
#   echo "codex:gpt-5.6-sol" > ~/.ailang/state/mission-executor-model-once
if [ -f "$EXEC_ONCE_FILE" ]; then
  once=$(head -1 "$EXEC_ONCE_FILE" 2>/dev/null)
  rm -f "$EXEC_ONCE_FILE"
  if [ -n "$once" ]; then
    export MISSION_EXECUTOR_MODEL="$once"
    log "one-shot executor override consumed: executor=$once (this iteration only)"
  fi
fi

# LANE DEGRADATION REPORT. World's five silently-demoted iterations (18/19/21/22, 2026-07)
# are why the shared driver grew a ledger; world never got the emit side. This is the
# minimum honest version: one line per degraded role, in the launchd log, BEFORE the
# iteration starts — never only in a routing-evidence row written after the fact.
if [ -n "$_lane_degraded" ]; then
  log "LANES DEGRADED THIS FIRE ($(printf '%s' "$_lane_degraded" | grep -c '^- ') role(s)):$(printf '%s' "$_lane_degraded" | tr '\n' ' ')"
fi
log "=== mission iteration starting (controller=$CONTROLLER_ID via ${MODEL_WHY}, timeout=${HARD_TIMEOUT}s | roles: designer=$MISSION_DESIGNER_MODEL planner=$MISSION_PLANNER_MODEL executor=$MISSION_EXECUTOR_MODEL evaluator=$MISSION_EVALUATOR_MODEL) ==="

PROMPT="Run one mission-control iteration: invoke the mission-control skill for \
${MISSION_DOC} and follow its gates. You are a scheduled run; \
there is no human present — park anything needing human input and report via \
ailang messages and the GitHub bookkeeping issue, per the skill. \
The authoritative runtime instructions are /Users/voightkampff/dev/sunholo-data/ailang/.claude/skills/mission-control/SKILL.md; \
read and follow that file even when the controller provider is Codex."

# _mc_run_once → runs the selected provider with BOTH watchdogs, waits, sets global RC.
# Watchdogs are per-attempt (fresh PIDs each retry).
_mc_run_once() {
  if [ "$CONTROLLER_PROVIDER" = "codex" ]; then
    codex exec --skip-git-repo-check \
      --dangerously-bypass-approvals-and-sandbox \
      --model "$MODEL" -C "$REPO" "$PROMPT" >>"$LOG" 2>&1 &
  else
    claude -p "$PROMPT" \
      --model "$MODEL" \
      --permission-mode bypassPermissions \
      >>"$LOG" 2>&1 &
  fi
  CONTROLLER_PID=$!
  printf '%s\n' "$CONTROLLER_PID" > "$PIDFILE"   # overlap guard reads this (per-attempt: retries refresh it)

  # Watchdog: TERM at the wall limit, KILL 60s later. (No GNU timeout on macOS.)
  (
    sleep "$HARD_TIMEOUT"
    if kill -0 "$CONTROLLER_PID" 2>/dev/null; then
      echo "[$(date '+%F %H:%M:%S')] HARD TIMEOUT ${HARD_TIMEOUT}s — killing $CONTROLLER_PID" >>"$LOG"
      kill -TERM "$CONTROLLER_PID" 2>/dev/null; sleep 60; kill -KILL "$CONTROLLER_PID" 2>/dev/null
    fi
  ) &
  WATCHDOG_PID=$!

  # Stall watchdog: after the grace window, sample for the wedged-tool fingerprint
  # (no progress + a descendant alive ≥ STALL_CHILD_AGE). STALL_SAMPLES consecutive
  # hits → kill early so the slot recycles instead of idling to HARD_TIMEOUT. hits
  # resets on ANY arm showing movement, so live work is never killed — and the
  # progress arms make that a property the suite can kill a mutant on (2026-09-02).
  (
    sleep "$STALL_GRACE"
    hits=0
    while kill -0 "$CONTROLLER_PID" 2>/dev/null; do
      if _mc_stalled "$CONTROLLER_PID"; then hits=$((hits + 1)); else hits=0; fi
      if [ "$hits" -ge "$STALL_SAMPLES" ]; then
        echo "[$(date '+%F %H:%M:%S')] STALL: $CONTROLLER_PROVIDER $CONTROLLER_PID made NO PROGRESS across $STALL_SAMPLES samples ($((STALL_SAMPLES * STALL_INTERVAL))s) with a descendant alive ≥${STALL_CHILD_AGE}s — killing early [${_MC_STALL_WHY:-unknown}]" >>"$LOG"
        kill -TERM "$CONTROLLER_PID" 2>/dev/null; sleep 30; kill -KILL "$CONTROLLER_PID" 2>/dev/null
        break
      fi
      sleep "$STALL_INTERVAL"
    done
  ) &
  STALL_PID=$!

  wait "$CONTROLLER_PID"; RC=$?
  kill "$WATCHDOG_PID" "$STALL_PID" 2>/dev/null
  return "$RC"
}

# Run with transient-retry. On a non-zero exit that is NOT a deliberate watchdog
# kill (143/137) AND whose THIS-attempt output carries a transient signature,
# back off and re-run — up to TRANSIENT_RETRIES total attempts.
# DURATION CLOCK (2026-09-05). World wrote no durable record of how long an iteration
# takes — the only trace was the /tmp driver log, which does not survive a reboot, so
# after the 09-05 restart world was the ONE mission with zero duration history and its
# idle time could not be estimated at all. Same file and format as the shared driver's
# slot-verdict log, so the same tooling reads all four missions.
_mc_iter_start=$(date +%s)
attempt=1
while : ; do
  logpos=$(wc -l < "$LOG" 2>/dev/null || echo 0)
  _mc_run_once; RC=$?
  [ "$RC" -eq 0 ] && break
  case "$RC" in 143|137) break ;; esac   # watchdog kill — never retry
  if [ "$attempt" -lt "$TRANSIENT_RETRIES" ] \
     && tail -n +$((logpos + 1)) "$LOG" 2>/dev/null | grep -qiE "$TRANSIENT_SIG"; then
    backoff=$(( TRANSIENT_BACKOFF * attempt ))
    log "transient API error (rc=$RC) attempt $attempt/$TRANSIENT_RETRIES — retrying in ${backoff}s (Anthropic capacity)"
    sleep "$backoff"
    attempt=$((attempt + 1))
    continue
  fi
  break
done

rm -f "$PIDFILE"   # this instance owns the run; yield paths above never reach here

# Durable one-line record per completed slot: ~/.ailang/state survives the reboots that
# wipe /tmp. Written for EVERY outcome, not just rc=0 — a crashed slot's elapsed time is
# exactly what tells a timeout apart from an instant refusal.
_mc_iter_elapsed=$(( $(date +%s) - _mc_iter_start ))
case "$RC" in 0) _mc_iter_verdict=COMPLETED ;; 143|137) _mc_iter_verdict=KILLED ;; *) _mc_iter_verdict=CRASHED ;; esac
printf '%s verdict=%s rc=%s attempt=%s/%s elapsed_s=%s controller=%s\n' \
  "$(date -u '+%Y-%m-%dT%H:%M:%SZ')" "$_mc_iter_verdict" "$RC" "$attempt" "$TRANSIENT_RETRIES" \
  "$_mc_iter_elapsed" "${CONTROLLER_ID:-unknown}" \
  >> "$STATE_DIR/mission-${MISSION_NAME}-slot-verdicts.log" 2>/dev/null || true
log "slot-verdict: $_mc_iter_verdict rc=$RC attempt=$attempt/$TRANSIENT_RETRIES elapsed_s=$_mc_iter_elapsed mission=$MISSION_NAME"

if [ "$RC" -ne 0 ]; then
  log "iteration exited rc=$RC"
  ailang messages send controlplane \
    "mission-control iteration exited rc=$RC (timeout or crash). Log: $LOG" \
    --title "Mission iteration FAILED (rc=$RC)" --from "$MSG_FROM" 2>/dev/null
  [ -n "${MISSION_GH_ISSUE:-}" ] && gh issue comment "$MISSION_GH_ISSUE" --repo "$MISSION_REPO" \
    --body "⚠️ Mission iteration **FAILED to complete** (rc=$RC — timeout or crash) at $(date '+%F %H:%M %Z'). Log on the rig: \`$LOG\`. The queue is untouched; the next interval will retry." 2>/dev/null
else
  log "iteration complete (rc=0)"
  # The skill itself sends the substantive report (Gate 5, both channels).
fi
exit "$RC"
