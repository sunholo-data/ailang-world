# lib/world.sh: helpers shared by the ailang-world herdr plugin's commands.
# Sourced by bin/*, never run. The caller sets ROOT (the plugin directory).
#
# herdr is called only through $HERDR_BIN_PATH. World is called only through
# $WORLDD: `doctor --quick --json` here, `episode start` in the popup.
# Session tokens never pass through this file: the pane reads its session file.

# world_error WHAT WHY NEXT: the three-part message, on stderr.
world_error() {
  printf 'ailang-world: %s\n  why: %s\n  next: %s\n' "$1" "$2" "$3" >&2
}

# world_die WHAT WHY NEXT: the three-part message, then exit 2.
world_die() {
  world_error "$1" "$2" "$3"
  exit 2
}

# world_need_config: source $HERDR_PLUGIN_CONFIG_DIR/world.env, or refuse.
world_need_config() {
  if [ -z "${HERDR_PLUGIN_CONFIG_DIR:-}" ]; then
    world_die "HERDR_PLUGIN_CONFIG_DIR is not set, so there is no world.env to read." \
      "herdr sets it for plugin commands; this command ran outside herdr." \
      "run it from herdr, or set HERDR_PLUGIN_CONFIG_DIR to the plugin's config directory"
  fi
  WORLD_ENV=$HERDR_PLUGIN_CONFIG_DIR/world.env
  if [ ! -f "$WORLD_ENV" ]; then
    world_die "no plugin config at $WORLD_ENV." \
      "the plugin needs World's paths, your repos and your grants before it can start an agent." \
      "run $ROOT/bin/world-config-init (or herdr's \"World: write config\" action), then edit $WORLD_ENV"
  fi
  # shellcheck disable=SC1090
  . "$WORLD_ENV"
  for _v in WORLDD WORLD_ADDR WORLD_DB WORLD_WS WORLD_SESSIONS WORLD_HOSTS WORLD_TTL WORLD_REPOS; do
    eval "_val=\${$_v:-}"
    if [ -z "$_val" ]; then
      world_die "$WORLD_ENV does not set $_v." \
        "every World pane needs it (world-config-init's comments say what each one is)." \
        "set $_v in $WORLD_ENV"
    fi
  done
  : "${WORLD_PROFILE:=}" "${WORLD_GRANTS:=}" "${WORLD_MODELS:=}"
}

# world_init: the config first, then the state directory, herdr and jq.
world_init() {
  world_need_config
  if [ -z "${HERDR_PLUGIN_STATE_DIR:-}" ]; then
    world_die "HERDR_PLUGIN_STATE_DIR is not set." \
      "the plugin records its panes there; herdr sets it for plugin commands." \
      "run this from herdr (the plugin's actions, events and startup do)"
  fi
  mkdir -p "$HERDR_PLUGIN_STATE_DIR" || world_die "cannot create $HERDR_PLUGIN_STATE_DIR." \
    "the plugin records its panes there." "check that directory's parent and permissions"
  if [ -z "${HERDR_BIN_PATH:-}" ] || [ ! -x "$HERDR_BIN_PATH" ]; then
    world_die "HERDR_BIN_PATH does not name the herdr binary." \
      "the plugin calls herdr only through it, and herdr sets it for plugin commands." \
      "run this from herdr (the plugin's actions, events and startup do)"
  fi
  if ! command -v jq >/dev/null 2>&1; then
    world_die "jq is not on PATH." \
      "the plugin reads herdr's and doctor's JSON with jq." \
      "install jq (brew install jq), then run the action again"
  fi
}

# world_herdr ARGS...: herdr, never reading the caller's stdin.
world_herdr() {
  "$HERDR_BIN_PATH" "$@" </dev/null
}

# world_notify TITLE BODY
world_notify() {
  world_herdr notification show "$1" --body "$2" >/dev/null ||
    printf 'ailang-world: herdr could not show a notification: %s\n' "$2" >&2
}

world_now() {
  date '+%Y-%m-%dT%H:%M:%S%z'
}

# world_squote S: S as one single-quoted shell word.
world_squote() {
  printf "'%s'" "$(printf '%s' "$1" | sed "s/'/'\\\\''/g")"
}

# world_safe S: S with anything but [A-Za-z0-9._-] replaced, for file names.
world_safe() {
  printf '%s' "$1" | tr -c 'A-Za-z0-9._-' '_'
}

world_valid_episode() {
  printf '%s' "$1" | grep -Eq '^[a-z0-9][a-z0-9-]{0,63}$'
}

# world_exit_pause CODE: in a popup, keep the last message on screen.
world_exit_pause() {
  if [ -t 0 ] && [ -t 1 ]; then
    printf '\nPress Enter to close.'
    IFS= read -r _ || :
  fi
  return "$1"
}

# world_current: INV_PANE and INV_WORKSPACE from herdr's current pane, where
# still unset. Returns 1 when herdr does not say.
world_current() {
  _cur=$(world_herdr pane current) || return 1
  if [ -z "$INV_PANE" ]; then
    INV_PANE=$(printf '%s' "$_cur" | jq -r 'first(.. | objects | select(has("pane_id")) | .pane_id | tostring) // empty' 2>/dev/null)
  fi
  if [ -z "$INV_WORKSPACE" ]; then
    _wsid=$(printf '%s' "$_cur" | jq -r 'first(.. | objects | select(has("workspace_id")) | .workspace_id | tostring) // empty' 2>/dev/null)
    [ -n "$_wsid" ] || return 1
    _list=$(world_herdr workspace list) || return 1
    INV_WORKSPACE=$(printf '%s' "$_list" | jq -r --arg id "$_wsid" \
      'first(.. | objects | select(has("workspace_id") and (.workspace_id | tostring) == $id) | (.label // .name // empty)) // empty' 2>/dev/null)
  fi
  [ -n "$INV_PANE" ] && [ -n "$INV_WORKSPACE" ]
}

# world_invoker PANE WORKSPACE: the pane the action was invoked from and its
# workspace label. Arguments win; then the hand-off bin/world-open wrote; then herdr.
world_invoker() {
  INV_PANE=$1
  INV_WORKSPACE=$2
  _hand=$HERDR_PLUGIN_STATE_DIR/invoker
  if [ -z "$INV_PANE" ] && [ -z "$INV_WORKSPACE" ] && [ -f "$_hand" ]; then
    INV_PANE=$(sed -n 1p "$_hand")
    INV_WORKSPACE=$(sed -n 2p "$_hand")
    rm -f "$_hand"
  fi
  if [ -z "$INV_PANE" ] || [ -z "$INV_WORKSPACE" ]; then
    world_current || world_die "herdr did not say which pane and workspace this is." \
      "$HERDR_BIN_PATH pane current and workspace list gave no pane id or workspace label." \
      "run the action from a pane, or pass --workspace <label> --pane <id>"
  fi
}

# world_repo LABEL: REPO_DIR and REPO_PROFILE from WORLD_REPOS, or refuse.
# A line is <herdr workspace label>=<repo path>[:<exec profile>].
world_repo() {
  REPO_DIR=
  REPO_PROFILE=
  _found=
  _ifs=$IFS
  IFS='
'
  set -f
  for _line in $WORLD_REPOS; do
    _line=$(printf '%s' "$_line" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    case $_line in '' | '#'*) continue ;; esac
    [ "${_line%%=*}" = "$1" ] || continue
    _val=${_line#*=}
    case $_val in
      *:*)
        _p=${_val##*:}
        case $_p in
          */*) REPO_DIR=$_val ;;
          '') REPO_DIR=${_val%:} ;;
          *) REPO_DIR=${_val%:*} REPO_PROFILE=$_p ;;
        esac
        ;;
      *) REPO_DIR=$_val ;;
    esac
    _found=1
    break
  done
  set +f
  IFS=$_ifs
  if [ -z "$_found" ]; then
    world_die "herdr workspace \"$1\" has no entry in WORLD_REPOS." \
      "the plugin starts each episode on the repo WORLD_REPOS names for the workspace, and it names none for this one." \
      "add a line  $1=/path/to/repo[:<exec profile>]  to WORLD_REPOS in $WORLD_ENV"
  fi
  case $REPO_DIR in "~/"*) REPO_DIR=$HOME/${REPO_DIR#"~/"} ;; esac
  if [ ! -d "$REPO_DIR" ]; then
    world_die "WORLD_REPOS maps \"$1\" to $REPO_DIR, which is not a directory." \
      "an episode's worktree is added from that repository." \
      "fix the \"$1=\" line of WORLD_REPOS in $WORLD_ENV"
  fi
  [ -n "$REPO_PROFILE" ] || REPO_PROFILE=$WORLD_PROFILE
}

# world_model HOST: MODEL from WORLD_MODELS (empty when unset).
world_model() {
  MODEL=
  set -f
  for _w in $WORLD_MODELS; do
    case $_w in "$1="*) MODEL=${_w#*=} ;; esac
  done
  set +f
}

# world_ask_host: HOST, read from stdin.
world_ask_host() {
  printf 'Host [claude/codex/pi]: '
  IFS= read -r HOST || HOST=
  case $HOST in
    claude | codex | pi) ;;
    *) world_die "\"$HOST\" is not a host this plugin starts." \
      "World has launchers for claude, codex and pi (ailang-worldd hosts init)." \
      "run the action again and answer claude, codex or pi" ;;
  esac
}

# world_doctor ARGS...: doctor's quick JSON pass. Sets DOC_OUT and DOC_CODE.
world_doctor() {
  if [ ! -x "$WORLDD" ]; then
    DOC_OUT=
    DOC_CODE=127
    return 0
  fi
  DOC_OUT=$("$WORLDD" --addr "$WORLD_ADDR" doctor --quick --json --db "$WORLD_DB" "$@" </dev/null 2>&1)
  DOC_CODE=$?
}

# world_doctor_problem: PROB_CHECK, PROB_DETAIL and PROB_FIX from the last
# doctor pass. A failed check wins; a host warning counts too, because that
# launch would be observed, not confined. Returns 1 when nothing is wrong.
world_doctor_problem() {
  PROB_CHECK=
  PROB_DETAIL=
  PROB_FIX=
  if [ "$DOC_CODE" -eq 127 ] && [ -z "$DOC_OUT" ]; then
    PROB_CHECK=daemon
    PROB_DETAIL="$WORLDD is not an executable ailang-worldd"
    PROB_FIX="set WORLDD in $WORLD_ENV to the path of ailang-worldd"
    return 0
  fi
  _p=$(printf '%s\n' "$DOC_OUT" | grep '^{' | jq -sc \
    '(map(select(.mark? == "fail")) | first) // (map(select(.check? == "host" and .mark? == "warn")) | first) // empty' 2>/dev/null)
  if [ -n "$_p" ]; then
    PROB_CHECK=$(printf '%s' "$_p" | jq -r '.check // "doctor"')
    PROB_DETAIL=$(printf '%s' "$_p" | jq -r '.detail // ""')
    PROB_FIX=$(printf '%s' "$_p" | jq -r '.fix // ""')
    return 0
  fi
  [ "$DOC_CODE" -eq 0 ] && return 1
  case $DOC_CODE in
    10) PROB_CHECK=daemon ;;
    11) PROB_CHECK=session ;;
    12) PROB_CHECK=host ;;
    13) PROB_CHECK="exec tools" ;;
    *) PROB_CHECK=doctor ;;
  esac
  _first=$(printf '%s\n' "$DOC_OUT" | grep -v '^{' | sed -n 1p)
  PROB_DETAIL="doctor exited $DOC_CODE${_first:+: $_first}"
  PROB_FIX="run $WORLDD --addr $WORLD_ADDR doctor --db $WORLD_DB --quick for the full report"
  return 0
}

# world_session_left: LEFT ("<h>h<m>m left") and LEFT_MIN from the last
# doctor pass's session line. Returns 1 when it has none.
world_session_left() {
  LEFT=
  LEFT_MIN=
  _d=$(printf '%s\n' "$DOC_OUT" | grep '^{' | jq -r 'select(.check? == "session" and .mark? == "ok") | .detail' 2>/dev/null | sed -n 1p)
  _hm=$(printf '%s' "$_d" | sed -n 's/^[^:]*: \([0-9][0-9]*\) h \([0-9][0-9]*\) m left.*$/\1 \2/p')
  [ -n "$_hm" ] || return 1
  set -- $_hm
  LEFT="${1}h${2}m left"
  LEFT_MIN=$(($1 * 60 + $2))
}

# world_split PANE DIR: NEW_PANE, a pane split to the right of PANE in DIR.
world_split() {
  _out=$(world_herdr pane split "$1" --direction right --cwd "$2" --no-focus) ||
    world_die "herdr could not split pane $1." \
      "$HERDR_BIN_PATH pane split exited non-zero (its message is above)." \
      "check that pane $1 is still open, then run the action again"
  NEW_PANE=$(printf '%s' "$_out" | jq -r '.result.pane.pane_id // empty | tostring' 2>/dev/null)
  if [ -z "$NEW_PANE" ]; then
    world_die "herdr split pane $1 but named no new pane." \
      "its answer had no result.pane.pane_id." \
      "check herdr's version (0.9.3 or later), then run the action again"
  fi
}

# world_open_plain PANE REPO HOST MODE PAIR_OF: a plain, unconfined pane running
# HOST in REPO. MODE is plain (the comparison arm) or fallback (World failed).
world_open_plain() {
  _host=$3
  _mode=$4
  _pair=$5
  world_split "$1" "$2"
  world_herdr pane run "$NEW_PANE" "$_host" >/dev/null
  if [ "$_mode" = fallback ]; then
    world_herdr pane rename "$NEW_PANE" "fallback · $_host" >/dev/null
    world_herdr pane report-metadata "$NEW_PANE" --source ailang-world --token arm=plain \
      --token world=fallback --state-label "idle=world: fallback" >/dev/null
  else
    if [ -n "$_pair" ]; then
      world_herdr pane rename "$NEW_PANE" "plain · $_host · $_pair" >/dev/null
      world_herdr pane report-metadata "$NEW_PANE" --source ailang-world --token arm=plain \
        --token "pair_of=$_pair" >/dev/null
    else
      world_herdr pane rename "$NEW_PANE" "plain · $_host" >/dev/null
      world_herdr pane report-metadata "$NEW_PANE" --source ailang-world --token arm=plain >/dev/null
    fi
  fi
  jq -nc --arg pane "$NEW_PANE" --arg host "$_host" --arg repo "$2" --arg started "$(world_now)" \
    --arg world "$_mode" --arg pair "$_pair" \
    '{pane: $pane, arm: "plain", world: $world, host: $host, repo: $repo, started: $started}
     + (if $pair == "" then {} else {pair_of: $pair} end)' >>"$HERDR_PLUGIN_STATE_DIR/arms.jsonl"
}

# world_mark_ended FILE PANE: add "ended" to PANE's open rows in FILE.
world_mark_ended() {
  [ -f "$1" ] || return 0
  _tmp=$1.tmp.$$
  if jq -c --arg p "$2" --arg t "$(world_now)" \
    'if (.pane | tostring) == $p and .ended == null then . + {ended: $t} else . end' "$1" >"$_tmp"; then
    mv "$_tmp" "$1"
  else
    rm -f "$_tmp"
    world_error "could not mark pane $2 ended in $1." "jq could not read that file." "check $1 holds one JSON object per line"
  fi
}
