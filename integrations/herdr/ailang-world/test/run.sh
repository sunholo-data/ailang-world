#!/bin/sh
# test/run.sh: the ailang-world herdr plugin's tests. POSIX sh and jq; cases 9
# and 13-19 also need python3 >= 3.11 (tomllib). Fake herdr, fake
# ailang-worldd: no network, no real herdr, nothing written outside a
# temporary directory. Every case points HERDR_CONFIG_FILE at a temporary
# file, so herdr's real config is never read or written.
set -u
HERE=$(CDPATH='' cd -- "$(dirname -- "$0")" && pwd)
PLUGIN=$(dirname "$HERE")
BIN=$PLUGIN/bin

if ! command -v jq >/dev/null 2>&1; then
  printf 'run.sh: jq is not on PATH.\n  why: the plugin and these tests read JSON with jq.\n  next: install jq, then run again\n' >&2
  exit 2
fi
notexec=
for f in "$BIN"/*; do [ -x "$f" ] || notexec="$notexec $f"; done
if [ -n "$notexec" ]; then
  printf 'run.sh: not executable:%s\n  why: herdr runs each command directly.\n  next: chmod 755%s\n' "$notexec" "$notexec" >&2
  exit 2
fi

TMP=$(mktemp -d "${TMPDIR:-/tmp}/world-herdr-test.XXXXXX") || exit 2
export HOME="$TMP/home" XDG_CONFIG_HOME="$TMP/xdg"
mkdir -p "$HOME" "$XDG_CONFIG_HOME"
# A watcher an action started (case 23) is stopped before the directory goes.
stop_watchers() {
  for _p in "$TMP"/*/state/world-watch.lock/pid; do
    _w=$(cat "$_p" 2>/dev/null) || continue
    [ -z "$_w" ] || [ "$_w" = "$$" ] || kill "$_w" 2>/dev/null
  done
  return 0
}
trap 'stop_watchers; rm -rf "$TMP"' EXIT
PASS=0
FAIL=0
ok() { PASS=$((PASS + 1)); printf 'ok   %s\n' "$1"; }
bad() { FAIL=$((FAIL + 1)); printf 'FAIL %s\n' "$1"; }
check() { _d=$1; shift; if "$@"; then ok "$_d"; else bad "$_d"; fi; }
refute() { _d=$1; shift; if "$@"; then bad "$_d"; else ok "$_d"; fi; }
same() { # DESC WANT GOT
  if cmp -s "$2" "$3"; then ok "$1"; else bad "$1"; diff "$2" "$3" | sed 's/^/      /'; fi
}
has() { grep -qF -- "$2" "$1"; }
count() { grep -cF -- "$2" "$1" || :; }

# ---------------------------------------------------------------- the fakes
mkdir -p "$TMP/fake"
cat >"$TMP/fake/herdr" <<'EOF'
#!/bin/sh
# Fake herdr: one line per call, its argv joined by spaces.
printf '%s\n' "$*" >>"$FAKE_HERDR_LOG"
case "${1-} ${2-}" in
  "pane split") printf '%s\n' '{"result":{"pane":{"pane_id":"p-new"}}}' ;;
  "pane current") printf '%s\n' '{"result":{"pane":{"pane_id":"p-cur","workspace_id":"w-1"}}}' ;;
  "workspace list") printf '%s\n' '{"result":{"workspaces":[{"workspace_id":"w-0","label":"other"},{"workspace_id":"w-1","label":"myrepo"}]}}' ;;
esac
exit 0
EOF
cat >"$TMP/fake/ailang-worldd" <<'EOF'
#!/bin/sh
# Fake ailang-worldd: doctor answers from FAKE_DAEMON_* (no --session) or
# FAKE_SESS_* (--session); episode start writes a session file.
printf '%s\n' "$*" >>"$FAKE_WORLDD_LOG"
[ "${1-}" != --addr ] || shift 2
ok_daemon='{"check":"daemon","mark":"ok","detail":"ailang-worldd test at http://127.0.0.1:47330, store /x/world.db","fix":"","class":"daemon"}
{"failed":0}'
ok_sess='{"check":"daemon","mark":"ok","detail":"ailang-worldd test at http://127.0.0.1:47330, store /x/world.db","fix":"","class":"daemon"}
{"check":"session","mark":"ok","detail":"ep: 2 h 10 m left; grants Workspace.Read=worktree:200","fix":"","class":"session"}
{"check":"host","mark":"ok","detail":"claude at /usr/local/bin/claude: 2.1.0","fix":"","class":"host"}
{"check":"host","mark":"ok","detail":"launch files match hosts.json in /x/hosts","fix":"","class":"host"}
{"failed":0}'
expired_sess='{"check":"daemon","mark":"ok","detail":"d","fix":"","class":"daemon"}
{"check":"session","mark":"fail","detail":"session expired at 15:04 UTC","fix":"start a new session","class":"session"}
{"failed":1}'
# arg FLAG ARGS...: the value after FLAG.
arg() {
  _f=$1
  shift
  _prev=
  for _a in "$@"; do
    [ "$_prev" != "$_f" ] || { printf '%s' "$_a"; return 0; }
    _prev=$_a
  done
}
case ${1-} in
  doctor)
    case " $* " in
      *"/${FAKE_EXPIRED_EP:-none}.session "*) printf '%s\n' "$expired_sess"; exit 11 ;;
      *" --session "*) printf '%s\n' "${FAKE_SESS_JSON:-$ok_sess}"; exit "${FAKE_SESS_CODE:-0}" ;;
      *) printf '%s\n' "${FAKE_DAEMON_JSON:-$ok_daemon}"; exit "${FAKE_DAEMON_CODE:-0}" ;;
    esac
    ;;
  provenance)
    printf 'World-Provenance: store=/x/world.db episode=%s entries=1810-1820\n' "$(arg --episode "$@")"
    ;;
  why)
    printf 'entry %s: fake why chain\n' "${2-}"
    ;;
  session)
    [ "${2-}" = revoke ] || exit 2
    printf 'revoked 1 session of episode %s\n' "$(arg --episode "$@")"
    ;;
  episode)
    if [ "${2-}" = report ]; then
      printf '{"episode":"%s","entries":[1810,1820],"total_calls":42,"calls_by_tool":{"ailang-read":30,"ailang-write":12},"budgets":[{"effect":"Workspace.Read","scope":"worktree","minted":200,"used":30,"left":170},{"effect":"Ailang.Check","scope":"worktree","minted":100,"used":91,"left":9},{"effect":"Workspace.Write","scope":"worktree","minted":200,"used":12,"left":188}],"plan_refusals":{"count":0},"handler_refusals":{"count":0}}\n' "$3"
      exit 0
    fi
    [ "${2-}" = start ] || exit 2
    ep=$3
    sessions=
    prev=
    for a in "$@"; do
      [ "$prev" != --sessions-dir ] || sessions=$a
      prev=$a
    done
    if [ -n "${FAKE_MINT_STDIN:-}" ]; then
      IFS= read -r line || line='<eof>'
      printf '%s\n' "$line" >"$FAKE_MINT_STDIN"
    fi
    [ "${FAKE_MINT_CODE:-0}" = 0 ] || exit "$FAKE_MINT_CODE"
    mkdir -p "$sessions"
    (umask 077 && printf '%s\n' "$FAKE_TOKEN" >"$sessions/$ep.session")
    printf 'fake episode start: minted %s\n' "$ep"
    ;;
  *) exit 2 ;;
esac
EOF
chmod 755 "$TMP/fake/herdr" "$TMP/fake/ailang-worldd"

# setup NAME: a fresh case directory, config, state and logs.
setup() {
  T=$TMP/$1
  mkdir -p "$T/conf" "$T/state" "$T/repo" "$T/sessions" "$T/hosts/bin" "$T/herdr"
  export HERDR_BIN_PATH="$TMP/fake/herdr" HERDR_PLUGIN_ROOT="$PLUGIN" HERDR_SOCKET_PATH="$T/herdr.sock"
  export HERDR_PLUGIN_STATE_DIR="$T/state" HERDR_PLUGIN_CONFIG_DIR="$T/conf"
  # Never herdr's real config: every case has its own (absent until a case writes it).
  export HERDR_CONFIG_FILE="$T/herdr/config.toml"
  export FAKE_HERDR_LOG="$T/herdr.log" FAKE_WORLDD_LOG="$T/worldd.log" FAKE_TOKEN="tok-SECRET-$1-0123456789abcdef"
  unset FAKE_DAEMON_CODE FAKE_DAEMON_JSON FAKE_SESS_CODE FAKE_SESS_JSON FAKE_MINT_STDIN FAKE_MINT_CODE FAKE_EXPIRED_EP
  unset HERDR_PLUGIN_EVENT_JSON WORLD_WATCH_ONCE WORLD_WATCH_INTERVAL WORLD_KEYS
  # The tests may run inside herdr: its invoking context must not leak in.
  unset HERDR_PANE_ID HERDR_WORKSPACE_ID HERDR_TAB_ID HERDR_PLUGIN_CONTEXT_JSON HERDR_PLUGIN_ACTION_ID
  : >"$FAKE_HERDR_LOG"
  : >"$FAKE_WORLDD_LOG"
  cat >"$T/conf/world.env" <<EOF
WORLDD="$TMP/fake/ailang-worldd"
WORLD_ADDR="http://127.0.0.1:47330"
WORLD_DB="$T/world.db"
WORLD_WS="$T/ws"
WORLD_SESSIONS="$T/sessions"
WORLD_HOSTS="$T/hosts"
WORLD_PROFILE=""
WORLD_TTL="3600"
WORLD_GRANTS="--grant Workspace.Read=worktree:200 --grant Workspace.Write=worktree:200"
WORLD_REPOS="
myrepo=$T/repo:myprofile
"
WORLD_MODELS="claude=claude-opus-5-5 codex=gpt-6.1-sol pi="
EOF
}

# run SCRIPT ANSWERS ARGS...: run a command with ANSWERS as stdin.
run() {
  _s=$1
  _a=$2
  shift 2
  "$BIN/$_s" "$@" <"$_a" >"$T/out" 2>"$T/err"
  CODE=$?
}

line_no() { grep -nF -- "$2" "$1" | sed -n '1s/:.*//p'; }

# hold_watch_lock: the watcher looks running (this shell's pid), so an action starts none.
hold_watch_lock() { mkdir -p "$T/state/world-watch.lock" && printf '%s\n' "$$" >"$T/state/world-watch.lock/pid"; }
release_watch_lock() { rm -rf "$T/state/world-watch.lock"; }
# wait_for CMD...: up to 5 s for CMD to succeed.
wait_for() {
  _n=0
  while [ "$_n" -lt 50 ]; do
    "$@" && return 0
    sleep 0.1
    _n=$((_n + 1))
  done
  return 1
}
watch_pid() { cat "$T/state/world-watch.lock/pid" 2>/dev/null; }
watch_alive() { _w=$(watch_pid) && [ -n "$_w" ] && [ "$_w" != "$$" ] && kill -0 "$_w" 2>/dev/null; }
watch_gone() { [ ! -d "$T/state/world-watch.lock" ]; }

DOCTOR_DAEMON="--addr http://127.0.0.1:47330 doctor --quick --json"
DAEMON_DOWN='{"check":"daemon","mark":"fail","detail":"nothing is listening at http://127.0.0.1:47330","fix":"start World: ailang-worldd serve --db /x/world.db …","class":"daemon"}
{"failed":1}'

# ---------------------------------------------------------------- 1. happy path
setup happy
printf 'claude\nep1\nMINT-STDIN-SENTINEL\n' >"$T/answers"
export FAKE_MINT_STDIN="$T/mint-stdin"
run world-agent-popup "$T/answers" --workspace myrepo --pane p-1
check "1 popup exits 0" [ "$CODE" -eq 0 ]
cat >"$T/want" <<EOF
pane split p-1 --direction right --cwd $T/state/empty/ep1 --no-focus
pane run p-new WORLD_SESSION="\$(cat '$T/sessions/ep1.session')" exec '$T/hosts/bin/claude-world'
pane rename p-new ep1 · claude
pane report-metadata p-new --source ailang-world --token episode=ep1 --token level=confined --token world=ok --token session=2h10m left
EOF
same "1 herdr calls: split in an empty dir, run the launcher on the session path, rename, label confined" "$T/want" "$T/herdr.log"
refute "1 no --env on any herdr call" has "$T/herdr.log" "--env"
cat >"$T/want" <<EOF
$DOCTOR_DAEMON --db $T/world.db
episode start ep1 --repo $T/repo --profile myprofile --agent claude --model claude-opus-5-5 --ttl 3600 --db $T/world.db --workspace-root $T/ws --sessions-dir $T/sessions --addr http://127.0.0.1:47330 --grant Workspace.Read=worktree:200 --grant Workspace.Write=worktree:200
$DOCTOR_DAEMON --db $T/world.db --session $T/sessions/ep1.session --host claude --hosts-dir $T/hosts
EOF
same "1 worldd calls: doctor, episode start (no --task), doctor --session --host" "$T/want" "$T/worldd.log"
empty_dir_ok() { [ -d "$T/state/empty/ep1" ] && [ -z "$(ls -A "$T/state/empty/ep1")" ]; }
check "1 the empty directory exists and is empty" empty_dir_ok
check "1 panes.jsonl records the confined pane" jq -e --arg repo "$T/repo" \
  'select(.pane == "p-new") | .episode == "ep1" and .host == "claude" and .repo == $repo and .level == "confined" and (.started | length > 0) and .ended == null' \
  "$T/state/panes.jsonl"
check "1 panes.jsonl has one row" [ "$(wc -l <"$T/state/panes.jsonl" | tr -d ' ')" -eq 1 ]
refute "1 the token is in no argv, output or state" grep -rqF -- "$FAKE_TOKEN" "$T/herdr.log" "$T/worldd.log" "$T/out" "$T/err" "$T/state"
check "1 the popup says what started" has "$T/out" "✓ ep1 · claude started in pane p-new: confined"
# 10 (runtime): the mint read the popup's own stdin, the line after the answers.
check "10 episode start inherits the popup's stdin" [ "$(cat "$T/mint-stdin")" = MINT-STDIN-SENTINEL ]

# ---------------------------------------------------------------- 2. daemon down, then f, then p
setup down-p
export FAKE_DAEMON_CODE=10 FAKE_DAEMON_JSON="$DAEMON_DOWN"
printf 'codex\nep2\nf\np\n' >"$T/answers"
run world-agent-popup "$T/answers" --workspace myrepo --pane p-1
check "2 popup exits 0" [ "$CODE" -eq 0 ]
cat >"$T/want" <<EOF
pane split p-1 --direction right --cwd $T/repo --no-focus
pane run p-new codex
pane rename p-new fallback · codex
pane report-metadata p-new --source ailang-world --token arm=plain --token world=fallback --state-label idle=world: fallback
notification show ailang-world --body World unavailable: nothing is listening at http://127.0.0.1:47330. Opened a plain pane (not confined).
EOF
same "2 a plain pane in the repo dir, labelled world=fallback, and a notification" "$T/want" "$T/herdr.log"
check "2 the failed check is shown" has "$T/out" "what: nothing is listening at http://127.0.0.1:47330"
check "2 f shows the fix again" [ "$(count "$T/out" "start World: ailang-worldd serve")" -ge 2 ]
refute "2 no mint" has "$T/worldd.log" "episode start"
refute "2 never labelled confined" has "$T/herdr.log" "level=confined"
check "2 fallbacks.jsonl holds the check, detail and fix" jq -e \
  '.pane == "p-new" and .phase == "before-mint" and .check == "daemon" and .detail == "nothing is listening at http://127.0.0.1:47330" and .fix == "start World: ailang-worldd serve --db /x/world.db …"' \
  "$T/state/fallbacks.jsonl"
check "2 arms.jsonl records the fallback" jq -e '.arm == "plain" and .world == "fallback" and .host == "codex"' "$T/state/arms.jsonl"
check "2 no World pane recorded" [ ! -e "$T/state/panes.jsonl" ]

# ---------------------------------------------------------------- 3. daemon down, then q (no arguments: herdr says where)
setup down-q
export FAKE_DAEMON_CODE=10 FAKE_DAEMON_JSON="$DAEMON_DOWN"
printf 'claude\n\nq\n' >"$T/answers"
run world-agent-popup "$T/answers"
check "3 popup exits 0" [ "$CODE" -eq 0 ]
printf 'pane current\nworkspace list\n' >"$T/want"
same "3 no pane: herdr was only asked where the action ran" "$T/want" "$T/herdr.log"
refute "3 no mint" has "$T/worldd.log" "episode start"
check "3 quit says nothing started" has "$T/out" "Quit: nothing was minted and no pane was started."
check "3 no state rows" [ -z "$(ls "$T/state" | grep jsonl)" ]

# ---------------------------------------------------------------- 4. session check fails after the mint (pi, default episode id)
setup sess
export FAKE_SESS_CODE=11 FAKE_SESS_JSON='{"check":"daemon","mark":"ok","detail":"ailang-worldd test","fix":"","class":"daemon"}
{"check":"session","mark":"fail","detail":"h0000000000'"'"'s session expired at 15:04 UTC","fix":"start a new session: ailang-worldd session new h0000000000 … (or episode start)","class":"session"}
{"failed":1}'
printf 'pi\n\np\n' >"$T/answers"
run world-agent-popup "$T/answers" --workspace myrepo --pane p-1
check "4 popup exits 0" [ "$CODE" -eq 0 ]
check "4 the mint ran with the default episode id, --agent codex for pi, no --model" \
  grep -qE "^episode start h[0-9]{10} --repo $T/repo --profile myprofile --agent codex --ttl 3600 " "$T/worldd.log"
check "4 the second doctor checks the session and pi" grep -qE -- "--session $T/sessions/h[0-9]{10}\.session --host pi --hosts-dir $T/hosts\$" "$T/worldd.log"
minted=$(line_no "$T/out" "fake episode start: minted")
offered=$(line_no "$T/out" "✗ World check failed: session")
offer_after_mint() { [ -n "$minted" ] && [ -n "$offered" ] && [ "$minted" -lt "$offered" ]; }
check "4 the fallback offer comes after the mint" offer_after_mint
check "4 it says the episode remains" has "$T/out" "is minted; its worktree and session remain"
check "4 plain pane runs pi in the repo" has "$T/herdr.log" "pane run p-new pi"
check "4 fallbacks.jsonl: session, after the mint" jq -e '.check == "session" and .phase == "after-mint" and (.fix | startswith("start a new session"))' "$T/state/fallbacks.jsonl"
refute "4 never labelled confined" has "$T/herdr.log" "level=confined"

# 4b. a host warning (exit 0) is not confined either: the offer appears.
setup host-warn
export FAKE_SESS_JSON='{"check":"daemon","mark":"ok","detail":"ailang-worldd test","fix":"","class":"daemon"}
{"check":"session","mark":"ok","detail":"ep4b: 1 h 0 m left","fix":"","class":"session"}
{"check":"host","mark":"warn","detail":"bin/claude-world: launch config missing or digest differs; launch would be observed, not confined","fix":"ailang-worldd hosts init --out-dir /x --replace","class":"host"}
{"failed":0}'
printf 'claude\nep4b\nq\n' >"$T/answers"
run world-agent-popup "$T/answers" --workspace myrepo --pane p-1
check "4b host drift offers the fallback" has "$T/out" "✗ World check failed: host"
check "4b quit after the mint names the revoke" has "$T/out" "session revoke --db $T/world.db --episode ep4b"
check "4b no pane" [ ! -s "$T/herdr.log" ]

# ---------------------------------------------------------------- 5. no WORLD_REPOS entry
setup norepo
printf 'claude\nep5\n' >"$T/answers"
run world-agent-popup "$T/answers" --workspace elsewhere --pane p-1
check "5 refuses" [ "$CODE" -ne 0 ]
check "5 what" has "$T/err" 'herdr workspace "elsewhere" has no entry in WORLD_REPOS.'
check "5 why" has "$T/err" "  why: the plugin starts each episode on the repo WORLD_REPOS names"
check "5 next" has "$T/err" "  next: add a line  elsewhere=/path/to/repo[:<exec profile>]  to WORLD_REPOS in $T/conf/world.env"
check "5 no mint, no doctor" [ ! -s "$T/worldd.log" ]
check "5 no herdr call" [ ! -s "$T/herdr.log" ]

# ---------------------------------------------------------------- 6. world.env missing
setup noconf
rm -f "$T/conf/world.env"
printf '' >"$T/answers"
export HERDR_PLUGIN_EVENT_JSON='{"pane_id":"p-1"}'
missing_config_message() {
  has "$T/err" "ailang-world: no plugin config at $T/conf/world.env." &&
    has "$T/err" "  why: " && has "$T/err" "  next: run $PLUGIN/bin/world-config-init"
}
no_config_dir_message() {
  has "$T/err" "HERDR_PLUGIN_CONFIG_DIR is not set" && has "$T/err" "  why: " && has "$T/err" "  next: "
}
for f in "$BIN"/*; do
  s=${f##*/}
  case $s in
    # world-keys edits herdr's config, not World's: it needs no world.env (cases 13-19).
    world-config-init | world-keys) continue ;;
    world-open) set -- world-agent-popup ;;
    world-event) set -- pane.closed ;;
    *) set -- ;;
  esac
  run "$s" "$T/answers" "$@"
  check "6 $s refuses" [ "$CODE" -ne 0 ]
  check "6 $s: what, why, next" missing_config_message
done
check "6 nothing reached herdr" [ ! -s "$T/herdr.log" ]
check "6 nothing reached World" [ ! -s "$T/worldd.log" ]
(unset HERDR_PLUGIN_CONFIG_DIR && "$BIN/world-agent-popup" </dev/null >"$T/out" 2>"$T/err")
check "6 no config dir: what, why, next" no_config_dir_message
run world-config-init "$T/answers"
check "6 world-config-init exits 0" [ "$CODE" -eq 0 ]
check "6 world-config-init writes world.env" [ -f "$T/conf/world.env" ]
check "6 world.env is private" [ "$(ls -l "$T/conf/world.env" | cut -c1-10)" = "-rw-------" ]
check "6 world.env sources, with ten grants and a repo line" sh -c \
  '. "$1" && set -- $WORLD_GRANTS && [ $# -eq 20 ] && [ -n "$WORLD_REPOS" ] && [ -n "$WORLD_HOSTS" ]' sh "$T/conf/world.env"
run world-config-init "$T/answers"
check "6 world-config-init refuses a second time" [ "$CODE" -ne 0 ]
check "6 ... saying it never overwrites" has "$T/err" "already exists"

# ---------------------------------------------------------------- 7. the watcher
setup watch
printf '%s\n' \
  '{"pane":"p-7","episode":"ep7","host":"claude","repo":"r","started":"t","level":"confined"}' \
  '{"pane":"p-old","episode":"epold","host":"codex","repo":"r","started":"t","level":"confined","ended":"t2"}' \
  >"$T/state/panes.jsonl"
export WORLD_WATCH_ONCE=1
: >"$T/answers"
last_label() { grep '^pane report-metadata p-7 ' "$T/herdr.log" | tail -n 1; }
notes() { count "$T/herdr.log" "notification show"; }
label() { printf 'pane report-metadata p-7 --source ailang-world --token episode=ep7 --token level=confined --token world=%s --token session=%s' "$1" "$2"; }

run world-watch "$T/answers"
check "7 ok: tokens" [ "$(last_label)" = "$(label ok '2h10m left')" ]
check "7 ok: no notification" [ "$(notes)" -eq 0 ]

export FAKE_SESS_JSON='{"check":"daemon","mark":"ok","detail":"d","fix":"","class":"daemon"}
{"check":"session","mark":"ok","detail":"ep7: 0 h 45 m left; grants Workspace.Read=worktree:200","fix":"","class":"session"}
{"failed":0}'
run world-watch "$T/answers"
check "7 45 min: tokens" [ "$(last_label)" = "$(label ok '0h45m left')" ]
check "7 45 min: one notification" [ "$(notes)" -eq 1 ]
check "7 45 min: it names the time left" has "$T/herdr.log" "Episode ep7's session has 0h45m left. Fix: "

export FAKE_SESS_CODE=11 FAKE_SESS_JSON='{"check":"daemon","mark":"ok","detail":"d","fix":"","class":"daemon"}
{"check":"session","mark":"fail","detail":"ep7'"'"'s session expired at 15:04 UTC","fix":"start a new session: ailang-worldd session new ep7 … (or episode start)","class":"session"}
{"failed":1}'
run world-watch "$T/answers"
check "7 expired: tokens" [ "$(last_label)" = "$(label expired expired)" ]
check "7 expired: one more notification" [ "$(notes)" -eq 2 ]
check "7 expired: it names the detail and the fix" has "$T/herdr.log" \
  "Episode ep7: world expired: ep7's session expired at 15:04 UTC. Fix: start a new session"
run world-watch "$T/answers"
check "7 expired again: no new notification" [ "$(notes)" -eq 2 ]

export FAKE_SESS_CODE=10 FAKE_SESS_JSON="$DAEMON_DOWN"
run world-watch "$T/answers"
check "7 down: tokens" [ "$(last_label)" = "$(label down unknown)" ]
check "7 down: one more notification" [ "$(notes)" -eq 3 ]
refute "7 an ended pane is never checked" has "$T/worldd.log" "epold"
check "7 the lock is released" [ ! -d "$T/state/world-watch.lock" ]
mkdir "$T/state/world-watch.lock" && printf '%s\n' "$$" >"$T/state/world-watch.lock/pid"
before=$(wc -l <"$T/worldd.log" | tr -d ' ')
run world-watch "$T/answers"
check "7 a second instance exits 0" [ "$CODE" -eq 0 ]
check "7 ... saying one is running" has "$T/out" "already running"
check "7 ... without running doctor" [ "$(wc -l <"$T/worldd.log" | tr -d ' ')" -eq "$before" ]
rm -rf "$T/state/world-watch.lock"
unset WORLD_WATCH_ONCE

# ---------------------------------------------------------------- 8. events
setup events
printf '%s\n' '{"pane":"p-new","episode":"ep8","host":"claude","repo":"r","started":"t","level":"confined"}' >"$T/state/panes.jsonl"
printf '%s\n' '{"pane":"p-plain","arm":"plain","world":"plain","host":"codex","repo":"r","started":"t"}' >"$T/state/arms.jsonl"
: >"$T/answers"
export HERDR_PLUGIN_EVENT_JSON='{"event":"pane.closed","data":{"pane_id":"p-new"}}'
run world-event "$T/answers" pane.closed
check "8 pane.closed exits 0" [ "$CODE" -eq 0 ]
check "8 the World pane is marked ended" jq -e 'select(.pane == "p-new") | .ended != null' "$T/state/panes.jsonl"
check "8 a notification says what remains" has "$T/herdr.log" \
  "notification show ailang-world --body Episode ep8's pane closed. Its worktree and session remain; finish it from the CLI (row 218 will make this one step)"
export HERDR_PLUGIN_EVENT_JSON='{"pane_id":"p-new"}'
run world-event "$T/answers" pane.exited
check "8 pane.exited after pane.closed: no second notification" [ "$(count "$T/herdr.log" "notification show")" -eq 1 ]
export HERDR_PLUGIN_EVENT_JSON='{"pane":{"pane_id":"p-plain"}}'
run world-event "$T/answers" pane.closed
check "8 a plain pane is marked ended" jq -e 'select(.pane == "p-plain") | .ended != null' "$T/state/arms.jsonl"
check "8 a plain pane gets no notification" [ "$(count "$T/herdr.log" "notification show")" -eq 1 ]

# ---------------------------------------------------------------- 9. the manifest
py=
if command -v python3 >/dev/null 2>&1 && python3 -I -c 'import tomllib' 2>/dev/null; then py=python3; fi
if [ -z "$py" ]; then
  bad "9 the manifest parses as TOML: needs python3 >= 3.11 (tomllib), not found"
else
  "$py" -I - "$PLUGIN/herdr-plugin.toml" >"$TMP/commands" 2>"$TMP/manifest-errors" <<'PY'
import sys, tomllib
with open(sys.argv[1], 'rb') as f:
    m = tomllib.load(f)
errs = []
top = {'id', 'name', 'version', 'min_herdr_version', 'platforms'}
# herdr 0.9.3 (checked against installed plugins and plugins.mdx): required keys,
# then the optional ones.
tables = {'actions': {'id', 'title', 'contexts', 'command'}, 'events': {'id', 'on', 'command'},
          'startup': {'command'}, 'panes': {'id', 'title', 'command'}}
optional = {'actions': {'description'}, 'events': set(), 'startup': set(),
            'panes': {'placement', 'width', 'height'}}
for k in top:
    if not m.get(k):
        errs.append('missing ' + k)
for k in m:
    if k not in top and k not in tables:
        errs.append('unknown key ' + k)
if m.get('min_herdr_version') != '0.9.3':
    errs.append('min_herdr_version is not 0.9.3')
if m.get('platforms') != ['macos', 'linux']:
    errs.append('platforms')
for table, keys in tables.items():
    for e in m.get(table, []):
        if not keys <= set(e) <= keys | optional[table]:
            errs.append('%s entry keys %s' % (table, sorted(e)))
        cmd = e.get('command')
        if not (isinstance(cmd, list) and cmd and all(isinstance(a, str) and a for a in cmd)):
            errs.append('%s command %r' % (table, cmd))
        else:
            print(cmd[0])
panes = {p.get('id') for p in m.get('panes', [])}
for p in m.get('panes', []):
    if p.get('placement', 'overlay') not in ('overlay', 'popup'):
        errs.append('pane %s placement %r' % (p.get('id'), p.get('placement')))
    if str(p.get('id', '')).endswith('-popup') and p.get('placement') != 'popup':
        errs.append('pane %s is not placement = "popup"' % p.get('id'))
actions = {a.get('id'): a for a in m.get('actions', [])}
for aid, title in (('world-agent', 'New World agent'), ('plain-agent', 'New plain agent')):
    a = actions.get(aid)
    if not a or a.get('title') != title or a.get('contexts') != ['workspace', 'pane']:
        errs.append('action ' + aid)
for a in actions.values():
    c = a.get('command') or []
    if c[:1] == ['bin/world-open'] and (len(c) != 2 or c[1] not in panes):
        errs.append('world-open opens no [[panes]] id: %r' % c)
if actions.get('world-agent', {}).get('command') != ['bin/world-open', 'world-agent-popup']:
    errs.append('world-agent does not open world-agent-popup')
for p in m.get('panes', []):
    if p.get('id') == 'world-agent-popup' and p.get('command') != ['bin/world-agent-popup']:
        errs.append('world-agent-popup does not run bin/world-agent-popup')
if sorted(e.get('on') for e in m.get('events', [])) != ['pane.closed', 'pane.exited']:
    errs.append('events are not pane.closed and pane.exited')
if len(m.get('startup', [])) != 1:
    errs.append('one [[startup]] expected')
for e in errs:
    print('manifest: ' + e, file=sys.stderr)
sys.exit(1 if errs else 0)
PY
  if [ $? -eq 0 ]; then ok "9 the manifest parses as TOML with herdr 0.9.3's fields"; else
    bad "9 the manifest parses as TOML with herdr 0.9.3's fields"
    sed 's/^/      /' "$TMP/manifest-errors"
  fi
  runnable() { [ -f "$PLUGIN/$1" ] && [ -x "$PLUGIN/$1" ]; }
  while IFS= read -r c; do
    check "9 $c exists and is executable" runnable "$c"
  done <"$TMP/commands"
fi

# ---------------------------------------------------------------- 10. the mint is never automated
set -- "$BIN"/* "$PLUGIN"/lib/*.sh
refute "10 no script pipes yes" grep -nw 'yes' "$@"
refute "10 nothing is piped into ailang-worldd" grep -nE '\|[[:space:]]*"?\$\{?WORLDD' "$@"
refute "10 no script, expect or unbuffer" grep -nE '(^|[;&|(`]|\$\()[[:space:]]*(script|expect|unbuffer)([[:space:]]|$)' "$@"
refute "10 the plugin never touches /dev/tty" grep -n '/dev/tty' "$@"
# The popup's one ailang-worldd command line is the mint, with no pipe or redirection.
# The mint inherits stdin and stdout; only stderr may be kept (to last-mint.err), never stdin.
mint_line_plain() {
  [ "$(grep -cE '^[[:space:]]*"\$WORLDD"' "$BIN/world-agent-popup")" -eq 1 ] &&
    grep -qxE '[[:space:]]*"\$WORLDD" "\$@"( 2>"\$mint_err")?' "$BIN/world-agent-popup"
}
lib_runs_only_doctor() {
  [ -z "$(grep -F '"$WORLDD"' "$PLUGIN/lib/world.sh" | grep -v doctor | grep -vF -- '-x "$WORLDD"')" ]
}
check "10 the popup runs ailang-worldd once, stdin and stdout inherited (stderr kept)" mint_line_plain
check "10 lib runs ailang-worldd only for doctor" lib_runs_only_doctor

# ---------------------------------------------------------------- 11. New plain agent
setup plain
printf 'codex\nep1\n' >"$T/answers"
run plain-agent-popup "$T/answers" --workspace myrepo --pane p-1
check "11 plain popup exits 0" [ "$CODE" -eq 0 ]
cat >"$T/want" <<EOF
pane split p-1 --direction right --cwd $T/repo --no-focus
pane run p-new codex
pane rename p-new plain · codex · ep1
pane report-metadata p-new --source ailang-world --token arm=plain --token pair_of=ep1
EOF
same "11 the host runs plainly in the repo, labelled arm=plain, pair_of" "$T/want" "$T/herdr.log"
check "11 arms.jsonl" jq -e --arg repo "$T/repo" '.arm == "plain" and .world == "plain" and .pair_of == "ep1" and .repo == $repo' "$T/state/arms.jsonl"
check "11 World is not called" [ ! -s "$T/worldd.log" ]

# ---------------------------------------------------------------- 12. world-open hands the invoking pane to the popup
setup open
: >"$T/answers"
hold_watch_lock
run world-open "$T/answers" world-agent-popup
release_watch_lock
printf 'pane current\nworkspace list\nplugin pane open --plugin ailang-world --entrypoint world-agent-popup --focus\n' >"$T/want"
same "12 world-open asks herdr where, then opens the popup" "$T/want" "$T/herdr.log"
check "12 the hand-off names the pane and workspace" [ "$(cat "$T/state/invoker")" = "$(printf 'p-cur\nmyrepo')" ]
refute "12 a running watcher is not started again" has "$T/err" "started world-watch"
: >"$T/herdr.log"
printf 'claude\nep12\n' >"$T/answers"
run world-agent-popup "$T/answers"
check "12 the popup splits the invoking pane" grep -qx "pane split p-cur --direction right --cwd $T/state/empty/ep12 --no-focus" "$T/herdr.log"
refute "12 the popup did not ask herdr again" has "$T/herdr.log" "pane current"
check "12 the hand-off is used once" [ ! -e "$T/state/invoker" ]

# ---------------------------------------------------------------- 13-19. key bindings
BEGIN_M='# >>> ailang-world key bindings (bin/world-keys) >>>'
END_M='# <<< ailang-world key bindings <<<'
# real_config FILE: a herdr config shaped like a real one (other plugins' bindings).
real_config() {
  cat >"$1" <<'EOF'
# herdr config
[ui.sidebar.agents]
rows = [["state_icon", "machine", "workspace", "tab"], ["agent"], ["$episode", "$world", "$session"]]

[[keys.command]]
key = "prefix+shift+v"
type = "plugin_action"
command = "herdr-file-viewer.open"
description = "File viewer"

[[keys.command]]
key = "prefix+shift+q"
type = "plugin_action"
command = "herdr-agent-quota.show"
description = "Agent quota"
EOF
}
# keys_setup NAME: a case whose herdr config is a copy of orig.toml, mode 0640.
keys_setup() {
  setup "$1"
  CFG=$HERDR_CONFIG_FILE
  BACKUP=$T/state/herdr-config.original.toml
  real_config "$T/orig.toml"
  cp "$T/orig.toml" "$CFG"
  chmod 640 "$CFG"
}
toml_ok() { [ -n "$py" ] && "$py" -I -c 'import sys,tomllib; tomllib.load(open(sys.argv[1],"rb"))' "$1" 2>/dev/null; }
block_count() { grep -cxF -- "$BEGIN_M" "$1" || :; }
outside_block() { awk -v b="$BEGIN_M" -v e="$END_M" '$0 == b { s = 1; next } s && $0 == e { s = 0; next } !s' "$1"; }
the_block() { awk -v b="$BEGIN_M" -v e="$END_M" '$0 == b { s = 1 } s { print } s && $0 == e { s = 0 }' "$1"; }
mode_of() { ls -l "$1" | cut -c1-10; }
reloads() { count "$T/herdr.log" "server reload-config"; }
only_config_beside() { [ "$(ls -A "$T/herdr")" = config.toml ]; }
three_part() { has "$T/err" "ailang-world: $1" && has "$T/err" "  why: " && has "$T/err" "  next: "; }
# want_block A L O FILE: the block world-keys writes for keys A, L and O.
want_block() {
  {
    printf '%s\n' "$BEGIN_M"
    printf '[[keys.command]]\nkey = "%s"\ntype = "plugin_action"\ncommand = "ailang-world.world-agent"\ndescription = "New World agent"\n\n' "$1"
    printf '[[keys.command]]\nkey = "%s"\ntype = "plugin_action"\ncommand = "ailang-world.plain-agent"\ndescription = "New plain agent"\n\n' "$2"
    printf '[[keys.command]]\nkey = "%s"\ntype = "plugin_action"\ncommand = "ailang-world.world-panel"\ndescription = "World panel"\n' "$3"
    printf '%s\n' "$END_M"
  } >"$4"
}

# 13. --check prints the block and changes nothing.
keys_setup keys-check
: >"$T/answers"
run world-keys "$T/answers" --check
check "13 --check exits 0" [ "$CODE" -eq 0 ]
check "13 it names the config" has "$T/out" "herdr config: $CFG"
want_block prefix+shift+a prefix+shift+l prefix+shift+o "$T/want"
the_block "$T/out" >"$T/got"
same "13 it prints the exact block" "$T/want" "$T/got"
check "13 no conflicts" has "$T/out" "Conflicts: none."
same "13 the config is unchanged" "$T/orig.toml" "$CFG"
check "13 no backup" [ ! -e "$BACKUP" ]
check "13 no herdr call" [ ! -s "$T/herdr.log" ]
check "13 nothing written beside the config" only_config_beside

# 14. --apply, answer n (and the default mode, answer empty): nothing changes.
keys_setup keys-no
printf 'n\n' >"$T/answers"
run world-keys "$T/answers" --apply
check "14 answer n exits 0" [ "$CODE" -eq 0 ]
check "14 it asked" has "$T/out" "Add these key bindings to $CFG? [y/N]"
check "14 it says nothing changed" has "$T/out" "No change: $CFG was not written, and no backup was made."
same "14 the config is byte-identical" "$T/orig.toml" "$CFG"
check "14 no backup" [ ! -e "$BACKUP" ]
check "14 no reload" [ "$(reloads)" -eq 0 ]
: >"$T/answers"
run world-keys "$T/answers"
check "14 the default is --apply, and no answer is no" has "$T/out" "No change: $CFG was not written"
same "14 still byte-identical" "$T/orig.toml" "$CFG"
check "14 nothing written beside the config" only_config_beside

# 15. --apply, answer y.
keys_setup keys-y
mode_before=$(mode_of "$CFG")
printf 'y\n' >"$T/answers"
run world-keys "$T/answers" --apply
check "15 answer y exits 0" [ "$CODE" -eq 0 ]
check "15 one block" [ "$(block_count "$CFG")" -eq 1 ]
want_block prefix+shift+a prefix+shift+l prefix+shift+o "$T/want"
the_block "$CFG" >"$T/got"
same "15 the block is the one shown" "$T/want" "$T/got"
outside_block "$CFG" >"$T/outside"
same "15 every original line is unchanged and in order" "$T/orig.toml" "$T/outside"
head -n "$(wc -l <"$T/orig.toml" | tr -d ' ')" "$CFG" >"$T/head"
same "15 ... with the block after them" "$T/orig.toml" "$T/head"
check "15 it parses as TOML" toml_ok "$CFG"
same "15 the backup equals the original" "$T/orig.toml" "$BACKUP"
check "15 server reload-config was called once" [ "$(reloads)" -eq 1 ]
check "15 the config keeps its mode" [ "$(mode_of "$CFG")" = "$mode_before" ]
check "15 it reports the keys" has "$T/out" "✓ Key bindings in $CFG; herdr reloaded its config."
check "15 ... each one" sh -c 'grep -q "prefix+shift+a  *New World agent" "$1" && grep -q "prefix+shift+l  *New plain agent" "$1" && grep -q "prefix+shift+o  *World panel" "$1"' sh "$T/out"
check "15 no temporary file is left beside the config" only_config_beside

# 16. Idempotent: a second apply leaves one block, and the backup is the first original.
cp "$CFG" "$T/after1"
printf 'y\n' >"$T/answers"
run world-keys "$T/answers" --apply
check "16 the second apply exits 0" [ "$CODE" -eq 0 ]
check "16 still exactly one block" [ "$(block_count "$CFG")" -eq 1 ]
same "16 the file is the same as after the first" "$T/after1" "$CFG"
same "16 the backup is still the first original" "$T/orig.toml" "$BACKUP"
check "16 it says the backup was kept" has "$T/out" "kept from the first time: $BACKUP"
# A block with lines after it is replaced where it is.
printf '\n[theme]\nname = "slate"\n' >>"$CFG"
cp "$CFG" "$T/with-theme"
export WORLD_KEYS="plain=prefix+shift+u"
printf 'y\n' >"$T/answers"
run world-keys "$T/answers" --apply
unset WORLD_KEYS
want_block prefix+shift+a prefix+shift+u prefix+shift+o "$T/want"
the_block "$CFG" >"$T/got"
same "16 a changed key replaces the block in place" "$T/want" "$T/got"
outside_block "$CFG" >"$T/outside"
outside_block "$T/with-theme" >"$T/want"
same "16 ... and the lines after it stay after it" "$T/want" "$T/outside"
check "16 still one block" [ "$(block_count "$CFG")" -eq 1 ]
# No config yet: it is created holding just the block, and the backup is empty.
setup keys-new
CFG=$HERDR_CONFIG_FILE
printf 'y\n' >"$T/answers"
run world-keys "$T/answers" --apply
check "16 a missing config: exits 0" [ "$CODE" -eq 0 ]
want_block prefix+shift+a prefix+shift+l prefix+shift+o "$T/want"
same "16 a missing config is created holding just the block" "$T/want" "$CFG"
backup_empty() { [ -f "$T/state/herdr-config.original.toml" ] && [ ! -s "$T/state/herdr-config.original.toml" ]; }
check "16 ... its backup is empty" backup_empty
check "16 ... it parses as TOML" toml_ok "$CFG"

# 17. Conflicts refuse, three-part, before asking; WORLD_KEYS moves the key.
keys_setup keys-conflict
{
  cat "$T/orig.toml"
  printf '\n[[keys.command]]\nkey = "prefix+shift+a"\ntype = "plugin_action"\ncommand = "other.thing"\n'
} >"$T/conflict.toml"
cp "$T/conflict.toml" "$CFG"
printf 'y\n' >"$T/answers"
run world-keys "$T/answers" --apply
check "17 a key bound elsewhere refuses" [ "$CODE" -ne 0 ]
check "17 what, why, next" three_part "prefix+shift+a (New World agent) is already bound in $CFG, outside the ailang-world block."
check "17 next names WORLD_KEYS in world.env" has "$T/err" "  next: pick another key: set the optional WORLD_KEYS=\"agent=… plain=… panel=…\" in $T/conf/world.env"
refute "17 it never asked" has "$T/out" "[y/N]"
same "17 the file is unchanged" "$T/conflict.toml" "$CFG"
check "17 no backup" [ ! -e "$BACKUP" ]
check "17 no reload" [ "$(reloads)" -eq 0 ]
run world-keys "$T/answers" --check
refused_with() { [ "$CODE" -ne 0 ] && has "$T/err" "$1"; }
check "17 --check reports it too, and exits non-zero" refused_with "prefix+shift+a (New World agent) is already bound"
printf 'WORLD_KEYS="agent=prefix+shift+n"\n' >>"$T/conf/world.env"
run world-keys "$T/answers" --apply
check "17 a herdr default key refuses" [ "$CODE" -ne 0 ]
check "17 ... what, why, next" three_part "prefix+shift+n (New World agent) is one of herdr's default keys."
same "17 ... the file is unchanged" "$T/conflict.toml" "$CFG"
printf 'WORLD_KEYS="agent=prefix+shift+e"\n' >>"$T/conf/world.env"
run world-keys "$T/answers" --apply
check "17 with WORLD_KEYS overriding the agent key, it applies" [ "$CODE" -eq 0 ]
check "17 ... binding prefix+shift+e" grep -qxF 'key = "prefix+shift+e"' "$CFG"
outside_block "$CFG" >"$T/outside"
same "17 ... and the other binding is untouched" "$T/conflict.toml" "$T/outside"

# 18. A result that is not valid TOML is refused; the file is unchanged.
setup keys-invalid
CFG=$HERDR_CONFIG_FILE
printf '# mine\nkeys.command = 1\n' >"$T/orig.toml"
cp "$T/orig.toml" "$CFG"
printf 'y\n' >"$T/answers"
run world-keys "$T/answers" --apply
check "18 refused" [ "$CODE" -ne 0 ]
check "18 what, why, next" three_part "the new $CFG would not be valid TOML: "
check "18 ... it says the config is unchanged" has "$T/err" "  why: herdr could not load it, so world-keys left $CFG unchanged"
same "18 the file is unchanged" "$T/orig.toml" "$CFG"
check "18 no draft is left beside it" only_config_beside
check "18 no reload" [ "$(reloads)" -eq 0 ]
printf 'not = [valid\n' >"$CFG"
cp "$CFG" "$T/broken.toml"
run world-keys "$T/answers" --apply
check "18 a config that is already broken is refused" refused_with "ailang-world: $CFG is not valid TOML now, before any change: "
check "18 ... what, why, next" three_part "$CFG is not valid TOML now, before any change: "
same "18 ... and left as it was" "$T/broken.toml" "$CFG"

# 19. --uninstall after apply restores the original, byte for byte.
keys_setup keys-uninstall
mode_before=$(mode_of "$CFG")
printf 'y\n' >"$T/answers"
run world-keys "$T/answers" --apply
cp "$CFG" "$T/applied"
printf 'n\n' >"$T/answers"
run world-keys "$T/answers" --uninstall
check "19 uninstall asks" has "$T/out" "Remove the ailang-world key bindings from $CFG? [y/N]"
check "19 answer n changes nothing" has "$T/out" "No change: $CFG was not written."
same "19 ... the file is as applied" "$T/applied" "$CFG"
printf 'y\n' >"$T/answers"
run world-keys "$T/answers" --uninstall
check "19 answer y exits 0" [ "$CODE" -eq 0 ]
same "19 the file is byte-identical to the original" "$T/orig.toml" "$CFG"
check "19 it reloaded (apply, then uninstall)" [ "$(reloads)" -eq 2 ]
check "19 the config keeps its mode" [ "$(mode_of "$CFG")" = "$mode_before" ]
check "19 nothing left beside the config" only_config_beside
run world-keys "$T/answers" --uninstall
check "19 with no block it says so" has "$T/out" "No ailang-world key bindings in $CFG: nothing to remove, nothing changed."
same "19 ... and changes nothing" "$T/orig.toml" "$CFG"
check "19 ... and does not reload" [ "$(reloads)" -eq 2 ]
# A block with the user's lines after it: only the block goes.
printf 'y\n' >"$T/answers"
run world-keys "$T/answers" --apply
printf '[theme]\nname = "slate"\n' >>"$CFG"
{ cat "$T/orig.toml"; printf '[theme]\nname = "slate"\n'; } >"$T/want"
run world-keys "$T/answers" --uninstall
same "19 lines after the block stay; only the block is removed" "$T/want" "$CFG"

# ---------------------------------------------------------------- 20. the panel, one frame
setup panel
printf '%s\n' \
  "{\"pane\":\"p-20\",\"episode\":\"ep20\",\"host\":\"claude\",\"repo\":\"$T/repo\",\"started\":\"t\",\"level\":\"confined\"}" \
  "{\"pane\":\"p-21\",\"episode\":\"ep21\",\"host\":\"codex\",\"repo\":\"$T/repo\",\"started\":\"t\",\"level\":\"confined\"}" \
  "{\"pane\":\"p-old\",\"episode\":\"epold\",\"host\":\"codex\",\"repo\":\"$T/repo\",\"started\":\"t\",\"level\":\"confined\",\"ended\":\"t2\"}" \
  "{\"pane\":\"p-x\",\"episode\":\"epx\",\"host\":\"claude\",\"repo\":\"/elsewhere\",\"started\":\"t\",\"level\":\"confined\"}" \
  >"$T/state/panes.jsonl"
printf '%s\n' \
  "{\"pane\":\"p-22\",\"arm\":\"plain\",\"world\":\"plain\",\"host\":\"codex\",\"repo\":\"$T/repo\",\"started\":\"t\",\"pair_of\":\"ep20\"}" \
  "{\"pane\":\"p-23\",\"arm\":\"plain\",\"world\":\"fallback\",\"host\":\"claude\",\"repo\":\"$T/repo\",\"started\":\"t\"}" \
  >"$T/state/arms.jsonl"
(umask 077 && printf '%s\n' "$FAKE_TOKEN" >"$T/sessions/ep20.session" && printf '%s\n' "$FAKE_TOKEN" >"$T/sessions/ep21.session")
export FAKE_EXPIRED_EP=ep21
run world-panel /dev/null --once --workspace myrepo
check "20 --once exits 0" [ "$CODE" -eq 0 ]
cp "$T/out" "$TMP/panel-frame.txt"
check "20 the header: World ok, the address, the workspace and repo, the watcher" \
  has "$T/out" "World ok · http://127.0.0.1:47330 · workspace myrepo → $T/repo · watcher stopped"
check "20 p-20: episode, host, confined, ok, time left, calls, lowest budget, selected" \
  grep -qE '^> p-20 +ep20 +claude +confined +ok +2h10m left +42  9 Ailang\.Check$' "$T/out"
check "20 p-21: expired" grep -qE '^  p-21 +ep21 +codex +confined +expired +expired +42  9 Ailang\.Check$' "$T/out"
check "20 p-22: a plain arm, with the episode it pairs" grep -qE '^  p-22 +ep20 +codex +plain +- +- +- +not in World$' "$T/out"
check "20 p-23: a fallback" grep -qE '^  p-23 +- +claude +fallback +- +- +- +not in World$' "$T/out"
refute "20 no ended pane" has "$T/out" "p-old"
refute "20 no pane of another repo" has "$T/out" "epx"
check "20 the selected episode's last entries" has "$T/out" "ep20: last entries (provenance --episode ep20)"
check "20 ... the provenance line" has "$T/out" "  World-Provenance: store=/x/world.db episode=ep20 entries=1810-1820"
check "20 ... and the last entry's index" has "$T/out" "  last entry: #1820 (w, then 1820, prints its why chain)"
check "20 the surface link" grep -qx 'surface: http://127.0.0.1:47330/surface' "$T/out"
check "20 the footer names the keys" has "$T/out" "j/k select · n new World agent · l new plain agent · o open surface · w why <entry> · r revoke session · q quit"
check "20 doctor checks each World pane's session by path" has "$T/worldd.log" "--session $T/sessions/ep20.session --host claude --hosts-dir $T/hosts"
check "20 the calls come from episode report --json" has "$T/worldd.log" "--addr http://127.0.0.1:47330 episode report ep20 --json"
refute "20 the token is in no output, argv or state" grep -rqF -- "$FAKE_TOKEN" "$T/out" "$T/err" "$T/herdr.log" "$T/worldd.log" "$T/state"
check "20 with --workspace, herdr is not asked" [ ! -s "$T/herdr.log" ]
run world-panel /dev/null --once
check "20 without --workspace, herdr names it" has "$T/herdr.log" "workspace list"
check "20 ... and the frame is the same workspace's" has "$T/out" "workspace myrepo → $T/repo"

# ---------------------------------------------------------------- 21. the panel's keys
panel_keys() { printf "$1" >"$T/answers"; : >"$T/herdr.log"; : >"$T/worldd.log"; run world-panel "$T/answers" --keys --workspace myrepo; }
panel_keys 'n\nl\nq\n'
check "21 q exits 0" [ "$CODE" -eq 0 ]
check "21 n invokes New World agent" grep -qx 'plugin action invoke world-agent --plugin ailang-world' "$T/herdr.log"
check "21 l invokes New plain agent" grep -qx 'plugin action invoke plain-agent --plugin ailang-world' "$T/herdr.log"
panel_keys 'q\n'
one_frame() { [ "$CODE" -eq 0 ] && [ "$(count "$T/out" "j/k select")" -eq 1 ]; }
check "21 q at once: exits 0 after one frame" one_frame
panel_keys 'r\nn\nq\n'
check "21 r asks first, naming the episode" has "$T/out" "Revoke episode ep20's session (pane p-20)?"
refute "21 r then n revokes nothing" has "$T/worldd.log" "session revoke"
check "21 ... and says so" has "$T/out" "Not revoked: ep20's session is unchanged."
panel_keys 'r\ny\nq\n'
check "21 r then y revokes the selected episode's session" grep -qx "session revoke --db $T/world.db --episode ep20" "$T/worldd.log"
check "21 ... once" [ "$(count "$T/worldd.log" "session revoke")" -eq 1 ]
check "21 ... and says so" has "$T/out" "Revoked ep20's session."
panel_keys 'j\nr\ny\n'
revoked_ep21() { [ "$CODE" -eq 0 ] && grep -qx "session revoke --db $T/world.db --episode ep21" "$T/worldd.log"; }
check "21 j selects the next row; end of input exits 0" revoked_ep21
panel_keys 'j\nj\nr\nq\n'
check "21 a plain row has no session to revoke" has "$T/out" "No World pane is selected"
refute "21 ... nothing revoked" has "$T/worldd.log" "session revoke"
panel_keys 'w\n1820\nq\n'
check "21 w prints why for the entry asked" has "$T/out" "entry 1820: fake why chain"
check "21 ... from ailang-worldd why" grep -qx -- "--addr http://127.0.0.1:47330 why 1820" "$T/worldd.log"
panel_keys 'w\nlast\nq\n'
refute "21 w refuses a non-number" has "$T/worldd.log" " why "
panel_keys 'o\nq\n'
check "21 o prints the surface URL" has "$T/out" "surface: http://127.0.0.1:47330/surface (not opened: no terminal)"
refute "21 the token is in no output or argv" grep -rqF -- "$FAKE_TOKEN" "$T/out" "$T/err" "$T/herdr.log" "$T/worldd.log"
unset FAKE_EXPIRED_EP

# ---------------------------------------------------------------- 22. world-open uses herdr's invoking context
setup open-env
export HERDR_PANE_ID=p-env HERDR_WORKSPACE_ID=w-1
: >"$T/answers"
hold_watch_lock
run world-open "$T/answers" world-agent-popup
printf 'workspace list\nplugin pane open --plugin ailang-world --entrypoint world-agent-popup --focus\n' >"$T/want"
same "22 with HERDR_PANE_ID/HERDR_WORKSPACE_ID: only the workspace's label is asked" "$T/want" "$T/herdr.log"
check "22 the hand-off names the invoking pane" [ "$(cat "$T/state/invoker")" = "$(printf 'p-env\nmyrepo')" ]
: >"$T/herdr.log"
rm -f "$T/state/invoker"
printf 'codex\n\n' >"$T/answers"
run plain-agent-popup "$T/answers"
check "22 a popup with no hand-off uses them too" grep -qx "pane split p-env --direction right --cwd $T/repo --no-focus" "$T/herdr.log"
refute "22 no pane current call" has "$T/herdr.log" "pane current"
release_watch_lock
unset HERDR_PANE_ID HERDR_WORKSPACE_ID

# ---------------------------------------------------------------- 23. the actions start the watcher
setup watch-start
export WORLD_WATCH_INTERVAL=1
: >"$T/answers"
run world-open "$T/answers" world-agent-popup
check "23 world-open starts the watcher when it isn't running" has "$T/err" "started world-watch"
check "23 ... which takes its lock" wait_for watch_alive
first=$(watch_pid)
run world-open "$T/answers" world-agent-popup
refute "23 a second world-open starts none" has "$T/err" "started world-watch"
run world-open-panel "$T/answers"
refute "23 world-open-panel starts none either" has "$T/err" "started world-watch"
check "23 world-open-panel opens the panel split, focused" \
  grep -qx 'plugin pane open --plugin ailang-world --entrypoint world-panel --placement split --focus' "$T/herdr.log"
check "23 still one watcher: the first" [ "$(watch_pid)" = "$first" ]
kill "$first" 2>/dev/null
check "23 the watcher stops and releases its lock" wait_for watch_gone
run world-open-panel "$T/answers"
check "23 world-open-panel starts it when it isn't running" has "$T/err" "started world-watch"
check "23 ... which takes its lock" wait_for watch_alive
kill "$(watch_pid)" 2>/dev/null
wait_for watch_gone || :
unset WORLD_WATCH_INTERVAL

# ---------------------------------------------------------------- 24. the manifest's new action and pane
if [ -z "$py" ]; then
  bad "24 the new actions and panes: needs python3 >= 3.11 (tomllib), not found"
else
  "$py" -I - "$PLUGIN/herdr-plugin.toml" 2>"$TMP/manifest-errors" <<'PY'
import sys, tomllib
with open(sys.argv[1], 'rb') as f:
    m = tomllib.load(f)
errs = []
actions = {a.get('id'): a for a in m.get('actions', [])}
panes = {p.get('id'): p for p in m.get('panes', [])}
want_actions = {
    'world-keys': ('World: install key bindings', ['bin/world-open', 'world-keys-popup']),
    'world-panel': ('World panel', ['bin/world-open-panel']),
}
for aid, (title, cmd) in want_actions.items():
    a = actions.get(aid)
    if not a:
        errs.append('no action ' + aid)
        continue
    if not {'id', 'title', 'contexts', 'command'} <= set(a):
        errs.append('action %s keys %s' % (aid, sorted(a)))
    if a.get('title') != title or a.get('contexts') != ['workspace', 'pane'] or a.get('command') != cmd:
        errs.append('action %s: %r' % (aid, a))
want_panes = {
    'world-keys-popup': ('World key bindings', 'popup', ['bin/world-keys']),
    'world-panel': ('World', 'overlay', ['bin/world-panel']),
}
for pid, (title, placement, cmd) in want_panes.items():
    p = panes.get(pid)
    if not p:
        errs.append('no pane ' + pid)
        continue
    if not {'id', 'title', 'command'} <= set(p):
        errs.append('pane %s keys %s' % (pid, sorted(p)))
    if p.get('title') != title or p.get('placement') != placement or p.get('command') != cmd:
        errs.append('pane %s: %r' % (pid, p))
for e in errs:
    print('manifest: ' + e, file=sys.stderr)
sys.exit(1 if errs else 0)
PY
  if [ $? -eq 0 ]; then ok "24 the new actions and panes carry herdr's required keys"; else
    bad "24 the new actions and panes carry herdr's required keys"
    sed 's/^/      /' "$TMP/manifest-errors"
  fi
fi
runnable() { [ -f "$PLUGIN/$1" ] && [ -x "$PLUGIN/$1" ]; }
for s in world-keys world-panel world-open-panel world-open; do
  check "24 bin/$s exists and is executable" runnable "bin/$s"
done
only_keys_name_herdr_config() { [ -z "$(grep -lF '.config/herdr' "$BIN"/* "$PLUGIN"/lib/*.sh | grep -v '/bin/world-keys$')" ]; }
check "24 only bin/world-keys names herdr's config directory" only_keys_name_herdr_config

# ---------------------------------------------------------------- 25. the host menu: 1/2/3, Enter repeats the last
setup menu
printf '2\nep25\ny\n' >"$T/answers"
run world-agent-popup "$T/answers" --workspace myrepo --pane p-1
check "25 answering 2 starts codex" has "$T/herdr.log" "/bin/codex-world'"
check "25 the menu shows the numbers" has "$T/out" "[1] claude  [2] codex  [3] pi"
check "25 the host is remembered" [ "$(cat "$T/state/last-host")" = codex ]
: >"$T/herdr.log"
printf '\nep25b\ny\n' >"$T/answers"
run world-agent-popup "$T/answers" --workspace myrepo --pane p-1
check "25 Enter repeats the last host (codex)" has "$T/herdr.log" "/bin/codex-world'"
: >"$T/herdr.log"
printf '3\nep25c\ny\n' >"$T/answers"
run world-agent-popup "$T/answers" --workspace myrepo --pane p-1
check "25 answering 3 starts pi" has "$T/herdr.log" "/bin/pi-world'"
: >"$T/herdr.log"
printf '9\nep25d\n' >"$T/answers"
run world-agent-popup "$T/answers" --workspace myrepo --pane p-1
check "25 an unknown answer is refused" [ "$CODE" -ne 0 ]
refute "25 an unknown answer starts nothing" has "$T/herdr.log" "pane split"
check "25 the refusal is kept in last-error.txt" has "$T/state/last-error.txt" '"9" is not a host this plugin starts.'

# ---------------------------------------------------------------- 26. home service defaults (never the real HOME)
mkdir -p "$HOME/.ailang-world"
cat >"$HOME/.ailang-world/service.json" <<EOF
{"schema":"world/service/v1","address":"127.0.0.1:7999","binary":"$HOME/.ailang-world/bin/ailang-worldd","store":"$HOME/data/world.db","workspace":"$HOME/workspaces"}
EOF
"$BIN/world-config-init" --config-dir "$TMP/config-service" >"$TMP/config-out" 2>"$TMP/config-err"
check "26 service record is named" has "$TMP/config-out" "$HOME/.ailang-world/service.json"
(
  . "$TMP/config-service/world.env"
  [ "$WORLDD" = "$HOME/.ailang-world/bin/ailang-worldd" ] &&
  [ "$WORLD_ADDR" = http://127.0.0.1:7999 ] &&
  [ "$WORLD_DB" = "$HOME/data/world.db" ] &&
  [ "$WORLD_WS" = "$HOME/workspaces" ] &&
  [ "$WORLD_SESSIONS" = "$HOME/data/sessions" ]
)
check "26 service values and sessions path" [ "$?" -eq 0 ]
rm "$HOME/.ailang-world/service.json"
"$BIN/world-config-init" --config-dir "$TMP/config-default" >"$TMP/config-out" 2>"$TMP/config-err"
check "26 absent record uses 7644" has "$TMP/config-default/world.env" 'WORLD_ADDR="http://127.0.0.1:7644"'
printf '%s\n' '{"schema":"world/service/v1","uninstalled":true,"address":"127.0.0.1:7999"}' >"$HOME/.ailang-world/service.json"
"$BIN/world-config-init" --config-dir "$TMP/config-uninstalled" >"$TMP/config-out" 2>"$TMP/config-err"
check "26 uninstalled record uses 7644" has "$TMP/config-uninstalled/world.env" 'WORLD_ADDR="http://127.0.0.1:7644"'
rm "$HOME/.ailang-world/service.json"

printf '\n%d passed, %d failed\n' "$PASS" "$FAIL"
[ "$FAIL" -eq 0 ]
