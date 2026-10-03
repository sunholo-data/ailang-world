#!/bin/bash
# Attended helper for row 134's M6: serve the eight se-tools transitions from a
# scratch World store and smoke them with curl, pi and Claude Code, every agent
# limited to World's MCP tools. Runbook: docs/QUICKSTART.md §9 — where they
# disagree, the runbook wins.
#
# WHAT THIS IS NOT: an automation of the two attended steps. The irreversible
# transition publish is NEVER run or built here (AC30,
# TestNoCIStepOrScriptReachesThePublishEntrypoint): `publish` only checks the
# daemon is stopped and points you at the Publish block of QUICKSTART §9, which
# you paste yourself from the repo root. `mint` is TTY-fenced in ailang-worldd;
# this script runs it (as `session new`, row 138) with stdin connected to
# /dev/tty (the fence's own requirement, honestly met) and you still type the
# confirmation. Every other
# step needs no human.
#
# WHAT IT ABSORBS (each measured on the first M6 run, 2026-10-03):
#   - AILANG_REGISTRY_API_KEY in the environment makes the daemon refuse to start
#     (broker guard) — it is unset for every child here.
#   - `fence=tty reason=stdin-is-not-the-controlling-terminal` in embedded
#     terminals — publish/mint read stdin from /dev/tty.
#   - POST /v1/commit is session-gated, so the genesis commit needs the minted
#     session and therefore runs at `serve`, after `mint`, not at `prepare`.
#   - /mcp/ answers as SSE (`event: message` / `data: …`) — parsed here.
#   - a backgrounded daemon is stopped by its pidfile, never by job control.
#   - pi-mcp-adapter registers direct tools from its cache; `pi` warms it once.
#   - pi and claude read a non-TTY stdin as extra prompt input and wait for EOF
#     (measured: a backgrounded `pi -p` hung) — both get </dev/null.
#
# Usage (from anywhere; paste ONE line at a time — interactive zsh does not
# treat `#` as a comment):
#   tools/attended/se_smoke.sh prepare   build CLIs, store, worktree, bootstrap
#   tools/attended/se_smoke.sh publish   ATTENDED: you paste QUICKSTART §9's Publish block
#   tools/attended/se_smoke.sh mint      ATTENDED: session new, the 8 se-tools grants
#   tools/attended/se_smoke.sh serve     daemon with tools + genesis commit
#   tools/attended/se_smoke.sh check     curl tools/list (expects 8) + one read
#   tools/attended/se_smoke.sh pi        pi smoke, World tools only
#   tools/attended/se_smoke.sh claude    Claude Code smoke, World tools only
#   tools/attended/se_smoke.sh log       log entries and effect journal counts
#   tools/attended/se_smoke.sh stop      stop the daemon
#   tools/attended/se_smoke.sh clean     stop + delete the scratch state
#
# Paths (override by env): SE_BASE (default ~/.ailang/se-smoke), SE_EPISODE
# (ep1), SE_PORT (7644), PIN (~/.pinned-ailang/ailang, v0.41.0), TOOL
# (~/.pinned-ailang-tools/v0.52.1/ailang), SE_EXAMPLES (~/.ailang/examples).
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
SE_BASE="${SE_BASE:-$HOME/.ailang/se-smoke}"
EP="${SE_EPISODE:-ep1}"
PORT="${SE_PORT:-7644}"
PIN="${PIN:-$HOME/.pinned-ailang/ailang}"
TOOL="${TOOL:-$HOME/.pinned-ailang-tools/v0.52.1/ailang}"
EXAMPLES="${SE_EXAMPLES:-$HOME/.ailang/examples}"

BIN="$SE_BASE/bin"
STORE_DIR="$SE_BASE/world"
DB="$STORE_DIR/world.db"
WS="$SE_BASE/ws"
PROJ="$SE_BASE/proj"
SESSION="$SE_BASE/session"
PIDFILE="$SE_BASE/daemon.pid"
LOG="$SE_BASE/daemon.log"
URL="http://127.0.0.1:$PORT"
TOOLS_CSV="ailang-read,ailang-write,ailang-edit,ailang-check,ailang-run,builtins-search,examples-search,ailang-cli"

unset AILANG_REGISTRY_API_KEY

die() { printf '\n  ✗ %s\n\n' "$*" >&2; exit 1; }
say() { printf '  %s\n' "$*"; }
ok()  { printf '  ✓ %s\n' "$*"; }

need_pin_and_tool() {
  [ -x "$PIN" ] || die "PIN not executable: $PIN"
  [ -x "$TOOL" ] || die "TOOL not executable: $TOOL"
  local pv tv
  pv="$("$PIN" --version)"; tv="$("$TOOL" --version)"
  case "$pv" in *'AILANG v0.41.0'*) ;; *) die "PIN is not AILANG v0.41.0: $PIN" ;; esac
  case "$tv" in *'AILANG v0.52.1'*) ;; *) die "TOOL is not AILANG v0.52.1: $TOOL" ;; esac
  ok "PIN v0.41.0, TOOL v0.52.1"
}

daemon_running() { [ -f "$PIDFILE" ] && kill -0 "$(cat "$PIDFILE")" 2>/dev/null; }

wait_health() {
  local i
  for i in $(seq 1 50); do
    curl -sf "$URL/v1/health" >/dev/null 2>&1 && return 0
    daemon_running || { tail -5 "$LOG" >&2; die "daemon exited during startup"; }
    sleep 0.2
  done
  die "daemon did not answer /v1/health on $URL"
}

start_daemon() { # args: extra serve flags
  daemon_running && die "a daemon is already running (pid $(cat "$PIDFILE")); run: $0 stop"
  "$BIN/ailang-worldd" serve --db "$DB" --bind "127.0.0.1:$PORT" --ailang-bin "$PIN" "$@" >"$LOG" 2>&1 &
  echo $! >"$PIDFILE"
  wait_health
  ok "daemon pid $(cat "$PIDFILE") on $URL"
}

stop_daemon() {
  if daemon_running; then
    kill "$(cat "$PIDFILE")"; sleep 1
    daemon_running && { kill -9 "$(cat "$PIDFILE")" 2>/dev/null || true; }
    ok "daemon stopped"
  else
    say "no daemon running"
  fi
  rm -f "$PIDFILE"
}

need_session() { [ -s "$SESSION" ] || die "no session at $SESSION — run: $0 mint"; }

mcp() { # args: JSON-RPC body; prints the `data:` JSON
  curl -s -H "Authorization: Bearer $(cat "$SESSION")" -H 'Content-Type: application/json' \
    -H 'Accept: application/json, text/event-stream' -d "$1" "$URL/mcp/" |
    python3 -c 'import sys
for line in sys.stdin:
    if line.startswith("data: "):
        print(line[6:].strip()); break'
}

cmd_prepare() {
  need_pin_and_tool
  mkdir -p "$BIN" "$STORE_DIR" "$WS"
  (cd "$REPO_ROOT" && go build -o "$BIN/ailang-worldd" ./cmd/ailang-worldd)
  ok "built ailang-worldd into $BIN"
  if [ ! -d "$WS/$EP" ]; then
    [ -d "$PROJ/.git" ] || { git init -q "$PROJ" && git -C "$PROJ" -c user.email=se-smoke@example.com -c user.name=se-smoke commit -q --allow-empty -m init; }
    git -C "$PROJ" worktree add -q --detach "$WS/$EP"
    printf 'module hello\n\nimport std/io (println)\n\nexport func main() -> () ! {IO} {\n  println("hello from ailang-run")\n}\n' >"$WS/$EP/hello.ail"
  fi
  ok "worktree $WS/$EP (hello.ail)"
  start_daemon   # first start bootstraps the epoch registry for PIN
  stop_daemon
  say "next (attended): $0 publish"
}

cmd_publish() {
  daemon_running && die "stop the daemon first (publish needs single-writer authority): $0 stop"
  say "attended: from $REPO_ROOT, paste the Publish block of docs/QUICKSTART.md §9 with"
  say "  --store $DB --ailang-bin $PIN  (append </dev/tty in an embedded terminal)"
  say "then: $0 mint"
}

cmd_mint() {
  daemon_running && die "stop the daemon first: $0 stop"
  # session new reuses the worktree prepare made and mints the eight
  # se-tools grants (--preset se-tools = the manifest's declaredEffects,
  # budget 50 each); --out is written once at 0600 and must not exist.
  rm -f "$SESSION"
  "$BIN/ailang-worldd" session new "$EP" --db "$DB" --workspace-root "$WS" --repo "$PROJ" \
    --preset se-tools --ttl 14400 --out "$SESSION" </dev/tty
  say "next: $0 serve"
}

genesis_json() {
  python3 - "$PIN" <<'EOF'
import json, hashlib, base64, sys
def sha(b): return "sha256:" + hashlib.sha256(b).hexdigest()
payload = json.dumps({"goal": "se-tools smoke"}).encode()
interp = open(sys.argv[1], "rb").read()
eh = sha(b"se-genesis-entry")
print(json.dumps({"observedHead": "",
  "objects": [{"hash": sha(payload), "interfaceHash": sha(b"iface-v1"),
               "semanticId": "world/demo/genesis-goal", "provenance": "se-smoke",
               "payload": base64.b64encode(payload).decode()}],
  "nextWorld": {"ref": sha(b"se-world-1"), "revision": 0, "stateRoot": sha(b"se-state-1"), "logHead": eh},
  "entry": {"header": {"entryIndex": 0, "semanticsEpoch": 1, "transitionFn": sha(payload),
                       "interpreter": sha(interp), "prevEntryHash": sha(b"genesis"), "writtenBy": "se-smoke"},
            "entryHash": eh, "transitionRef": sha(payload)}}))
EOF
}

cmd_serve() {
  need_pin_and_tool; need_session
  local ex=()
  [ -d "$EXAMPLES" ] && ex=(--examples-dir "$EXAMPLES")
  start_daemon --workspace-root "$WS" --tool-ailang-bin "$TOOL" ${ex[@]+"${ex[@]}"}
  local head
  head="$(curl -s "$URL/v1/head" || true)"   # plain text: the selected head ref, or empty
  case "$head" in
    sha256:*) ok "world head present (${head:0:23}…)" ;;
    *) genesis_json >"$SE_BASE/genesis.json"
       "$BIN/ailang-worldd" --addr "$URL" commit --file "$SE_BASE/genesis.json" --session "$(cat "$SESSION")" >/dev/null
       ok "genesis world committed" ;;
  esac
  say "next: $0 check"
}

cmd_check() {
  need_session; daemon_running || die "daemon not running: $0 serve"
  local list call
  list="$(mcp '{"jsonrpc":"2.0","id":1,"method":"tools/list"}')"
  echo "$list" | python3 -c 'import json,sys
names=sorted(t["name"] for t in json.load(sys.stdin)["result"]["tools"])
print("  tools/list:", ", ".join(names))
sys.exit(0 if len(names)==8 else 1)' || die "expected 8 tools"
  call="$(mcp '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"ailang-read","arguments":{"path":"hello.ail"}}}')"
  echo "$call" | python3 -c 'import json,sys
sc=json.load(sys.stdin)["result"]["structuredContent"]
assert sc["ok"], sc
print("  ailang-read ok; effect", sc["world"]["effects"][0]["status"], sc["world"]["effects"][0]["record"][:23]+"…")'
  ok "MCP surface live"
}

cmd_pi() {
  need_session; daemon_running || die "daemon not running: $0 serve"
  local cfg="$SE_BASE/pi-mcp.json" names
  printf '{"mcpServers":{"world":{"url":"%s/mcp/","auth":"bearer","bearerTokenEnv":"WORLD_SESSION","directTools":true,"toolPrefix":"server"}}}\n' "$URL" >"$cfg"
  names="$(echo "$TOOLS_CSV" | sed 's/\([^,]*\)/world_\1/g')"
  export WORLD_SESSION; WORLD_SESSION="$(cat "$SESSION")"
  cd "$WS/$EP"
  say "warming pi-mcp-adapter cache…"
  pi --mcp-config "$cfg" --no-builtin-tools -p "Reply with the single word ok." </dev/null >/dev/null
  pi --mcp-config "$cfg" --no-builtin-tools --tools "$names" \
    -p "${SE_PI_PROMPT:-Read hello.ail, change its message to 'hello from pi', check it, then run it and report its stdout.}" </dev/null
}

cmd_claude() {
  need_session; daemon_running || die "daemon not running: $0 serve"
  local d allowed
  d="$(mktemp -d)"; trap 'rm -rf "$d"' RETURN   # .mcp.json holds the raw bearer token
  (cd "$d" && claude mcp add --transport http world "$URL/mcp/" --header "Authorization: Bearer $(cat "$SESSION")" -s project >/dev/null)
  allowed="$(echo "$TOOLS_CSV" | sed 's/\([^,]*\)/mcp__world__\1/g')"
  cd "$WS/$EP"
  claude -p "${SE_CLAUDE_PROMPT:-Read hello.ail, change its message to 'hello from claude', check it, then run it and report its stdout.}" \
    --tools "" --strict-mcp-config --mcp-config "$d/.mcp.json" --allowedTools "$allowed" </dev/null
}

cmd_log() {
  need_session; daemon_running || die "daemon not running: $0 serve"
  local i=0 out
  while out="$(curl -sf -H "Authorization: Bearer $(cat "$SESSION")" "$URL/v1/log/$i")"; do
    echo "$out" | python3 -c 'import json,sys
d=json.load(sys.stdin); e=d.get("entry",d); h=e.get("header",e)
print("  %3d  %-18s %s" % (h.get("entryIndex",-1), h.get("writtenBy",""), (h.get("transitionFn") or "")[:23]))'
    i=$((i+1))
  done
  say "$i log entries; effect intents/outcomes: $(sqlite3 -readonly "file:$DB?mode=ro" "select count(*) from journal where kind='intent' and invocation_id like 'effect:$EP:%'")/$(sqlite3 -readonly "file:$DB?mode=ro" "select count(*) from journal where kind='outcome' and invocation_id like 'effect:$EP:%'")"
}

cmd_clean() {
  stop_daemon
  [ -d "$PROJ/.git" ] && git -C "$PROJ" worktree remove --force "$WS/$EP" 2>/dev/null || true
  rm -rf "$SE_BASE"
  ok "removed $SE_BASE"
}

case "${1:-}" in
  prepare) cmd_prepare ;; publish) cmd_publish ;; mint) cmd_mint ;;
  serve) cmd_serve ;; check) cmd_check ;; pi) cmd_pi ;; claude) cmd_claude ;;
  log) cmd_log ;; stop) stop_daemon ;; clean) cmd_clean ;;
  *) sed -n '2,/^set -euo/p' "$0" | sed 's/^# \{0,1\}//' | sed '$d'; exit 2 ;;
esac
