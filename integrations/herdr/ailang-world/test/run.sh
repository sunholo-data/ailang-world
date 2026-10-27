#!/bin/sh
# test/run.sh: the ailang-world herdr plugin's tests. POSIX sh and jq; case 9
# also needs python3 >= 3.11 (tomllib). Fake herdr, fake ailang-worldd: no
# network, no real herdr, nothing written outside a temporary directory.
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
  printf 'run.sh: not executable:%s\n  why: herdr runs each command directly.\n  next: chmod 755 %s/bin/*\n' "$notexec" "$PLUGIN" >&2
  exit 2
fi

TMP=$(mktemp -d "${TMPDIR:-/tmp}/world-herdr-test.XXXXXX") || exit 2
trap 'rm -rf "$TMP"' EXIT
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
case ${1-} in
  doctor)
    case " $* " in
      *" --session "*) printf '%s\n' "${FAKE_SESS_JSON:-$ok_sess}"; exit "${FAKE_SESS_CODE:-0}" ;;
      *) printf '%s\n' "${FAKE_DAEMON_JSON:-$ok_daemon}"; exit "${FAKE_DAEMON_CODE:-0}" ;;
    esac
    ;;
  episode)
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
  mkdir -p "$T/conf" "$T/state" "$T/repo" "$T/sessions" "$T/hosts/bin"
  export HERDR_BIN_PATH="$TMP/fake/herdr" HERDR_PLUGIN_ROOT="$PLUGIN" HERDR_SOCKET_PATH="$T/herdr.sock"
  export HERDR_PLUGIN_STATE_DIR="$T/state" HERDR_PLUGIN_CONFIG_DIR="$T/conf"
  export FAKE_HERDR_LOG="$T/herdr.log" FAKE_WORLDD_LOG="$T/worldd.log" FAKE_TOKEN="tok-SECRET-$1-0123456789abcdef"
  unset FAKE_DAEMON_CODE FAKE_DAEMON_JSON FAKE_SESS_CODE FAKE_SESS_JSON FAKE_MINT_STDIN FAKE_MINT_CODE
  unset HERDR_PLUGIN_EVENT_JSON WORLD_WATCH_ONCE
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
    world-config-init) continue ;;
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
tables = {'actions': {'id', 'title', 'contexts', 'command'}, 'events': {'on', 'command'},
          'startup': {'command'}, 'panes': {'id', 'placement', 'command'}}
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
        if set(e) != keys:
            errs.append('%s entry keys %s' % (table, sorted(e)))
        cmd = e.get('command')
        if not (isinstance(cmd, list) and cmd and all(isinstance(a, str) and a for a in cmd)):
            errs.append('%s command %r' % (table, cmd))
        else:
            print(cmd[0])
panes = {p.get('id') for p in m.get('panes', [])}
if any(p.get('placement') != 'popup' for p in m.get('panes', [])):
    errs.append('a pane is not placement = "popup"')
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
mint_line_plain() {
  [ "$(grep -cE '^[[:space:]]*"\$WORLDD"' "$BIN/world-agent-popup")" -eq 1 ] &&
    grep -qxE '[[:space:]]*"\$WORLDD" "\$@"' "$BIN/world-agent-popup"
}
lib_runs_only_doctor() {
  [ -z "$(grep -F '"$WORLDD"' "$PLUGIN/lib/world.sh" | grep -v doctor | grep -vF -- '-x "$WORLDD"')" ]
}
check "10 the popup runs ailang-worldd once, with no redirection" mint_line_plain
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
run world-open "$T/answers" world-agent-popup
printf 'pane current\nworkspace list\nplugin pane open world-agent-popup\n' >"$T/want"
same "12 world-open asks herdr where, then opens the popup" "$T/want" "$T/herdr.log"
check "12 the hand-off names the pane and workspace" [ "$(cat "$T/state/invoker")" = "$(printf 'p-cur\nmyrepo')" ]
: >"$T/herdr.log"
printf 'claude\nep12\n' >"$T/answers"
run world-agent-popup "$T/answers"
check "12 the popup splits the invoking pane" grep -qx "pane split p-cur --direction right --cwd $T/state/empty/ep12 --no-focus" "$T/herdr.log"
refute "12 the popup did not ask herdr again" has "$T/herdr.log" "pane current"
check "12 the hand-off is used once" [ ! -e "$T/state/invoker" ]

printf '\n%d passed, %d failed\n' "$PASS" "$FAIL"
[ "$FAIL" -eq 0 ]
