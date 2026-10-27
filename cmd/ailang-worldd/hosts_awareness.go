package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

const worldToolMap = `Your only file and shell access is World's MCP tools:
  read/edit/write files -> ailang-read / ailang-edit / ailang-write / ailang-delete (paths relative to the worktree)
  run commands         -> workspace-exec, the exec profile's commands only (list them with the tool)
  AILANG               -> ailang-check / ailang-run / ailang-cli; docs -> builtins-search / examples-search
`

const worldBoundary = "Out of budget or need a new grant? Say so; the operator mints it. Don't try to work around World."

// hostsWorlddPath is the ailang-worldd the generated scripts call: the binary
// running hosts init, symlinks resolved, so they work whether or not it is on
// the agent's PATH (a bare name read "daemon down" when it was not).
var hostsWorlddPath = func() string {
	exe, err := os.Executable()
	if err != nil {
		return "ailang-worldd"
	}
	if real, err := filepath.EvalSymlinks(exe); err == nil {
		return real
	}
	return exe
}

func hostDirectories(dir string) []string {
	out := []string{dir}
	for _, rel := range []string{"bin", "codex", "state", "claude-plugin", "claude-plugin/.claude-plugin", "claude-plugin/hooks", "pi"} {
		out = append(out, filepath.Join(dir, rel))
	}
	return out
}

func addAwarenessFiles(files map[string][]byte, dir, addr string) {
	putJSON := func(name string, v any) {
		b, _ := json.Marshal(v)
		files[name] = append(b, '\n')
	}
	context := filepath.Join(dir, "bin", "world-context")
	status := filepath.Join(dir, "bin", "world-statusline")
	plugin := filepath.Join(dir, "claude-plugin")
	home, _ := os.UserHomeDir()
	worldd := hostsWorlddPath()
	var settings map[string]any
	_ = json.Unmarshal(files["claude-world.settings.json"], &settings)
	settings["statusLine"] = map[string]any{"type": "command", "command": shellQuote(status), "refreshInterval": 60}
	putJSON("claude-world.settings.json", settings)
	putJSON("claude-plugin/.claude-plugin/plugin.json", map[string]string{
		"name": "ailang-world", "description": "Session awareness for confined World agents", "version": "1.0.0",
	})
	putJSON("claude-plugin/hooks/hooks.json", map[string]any{"hooks": map[string]any{
		"SessionStart": []any{map[string]any{"matcher": "", "hooks": []any{
			map[string]any{"type": "command", "command": shellQuote(context), "timeout": 10},
		}}},
	}})
	files["codex/AGENTS.md"] = []byte("You are working inside AILANG World (confined).\n" + worldToolMap + worldBoundary + "\n")
	hooks := map[string]any{}
	codexHook := filepath.Join(home, ".codex", "herdr-agent-state.sh")
	if info, err := os.Stat(codexHook); err == nil && info.Mode().IsRegular() {
		hooks["SessionStart"] = []any{map[string]any{"hooks": []any{map[string]any{
			"command": "bash " + shellQuote(codexHook) + " session", "timeout": 10, "type": "command",
		}}}}
	}
	putJSON("codex/hooks.json", map[string]any{"hooks": hooks})
	// Capture only the user's status command, never write to user settings.
	userCommand := ""
	var user struct {
		StatusLine struct {
			Command string `json:"command"`
		} `json:"statusLine"`
	}
	if b, err := os.ReadFile(filepath.Join(home, ".claude", "settings.json")); err == nil && json.Unmarshal(b, &user) == nil {
		userCommand = user.StatusLine.Command
	}
	files["bin/world-context"] = []byte("#!/bin/sh\n" + "addr=" + shellQuote(addr) + "\nworldd=" + shellQuote(worldd) + "\n" + `data=$("$worldd" --addr "$addr" whoami --json 2>/dev/null)
code=$?
if [ "$code" -ne 0 ]; then
  case $code in
    11) printf '%s\n' 'World context unavailable: session invalid or expired; ask the operator to mint a new session.' ;;
    126|127) printf 'World context unavailable: no ailang-worldd at %s; run hosts init --replace with the ailang-worldd you use.\n' "$worldd" ;;
    *) printf '%s\n' 'World context unavailable: daemon down; run ailang-worldd doctor and start World.' ;;
  esac
  exit 0
fi
if ! printf '%s' "$data" | jq -e '.episode and .secondsLeft and .grants' >/dev/null 2>&1; then
  printf '%s\n' 'World context unavailable: invalid response; run ailang-worldd doctor.'
  exit 0
fi
printf '%s\n' 'You are working inside AILANG World (confined).'
printf '%s' "$data" | jq -r '"Episode \(.episode) · worktree \(.worktree) · session \((.secondsLeft/3600|floor))h\((.secondsLeft/60|floor)%60)m left · daemon \(.addr)"'
cat <<'WORLD_TOOLS'
` + worldToolMap + `WORLD_TOOLS
printf '%s' "$data" | jq -r '"Grants left: " + ([.grants[] | "\(.effect) \(.left)/\(.budget)"] | join(", "))'
printf 'Every call is recorded. See it: %s/surface  ·  explain an entry: ailang-worldd why <entry>\n' "$addr"
printf '%s\n' ` + shellQuote(worldBoundary) + "\n")
	files["bin/world-statusline"] = []byte("#!/bin/sh\numask 077\naddr=" + shellQuote(addr) + "\nstate=" + shellQuote(filepath.Join(dir, "state")) + "\nuser_command=" + shellQuote(userCommand) + "\nworldd=" + shellQuote(worldd) + "\n" + `input=$(cat; printf '.')
input=${input%.}
# A new opaque cache filename per credential; neither token nor its hash is stored.
# cksum is a non-secret cache discriminator, not credential authority.
key=$(printf '%s' "${WORLD_SESSION:-}" | cksum | cut -d ' ' -f 1)
cache=$state/status-$key.json
now=$(date +%s)
data=
if [ -f "$cache" ] && [ ! -L "$cache" ]; then
  data=$(cat "$cache")
  at=$(printf '%s' "$data" | jq -r '.cachedAt // 0' 2>/dev/null)
  case $at in ''|*[!0-9]*) at=0 ;; esac
  [ "$((now-at))" -lt 30 ] && [ "$now" -ge "$at" ] || data=
fi
if [ -z "$data" ]; then
  response=$("$worldd" --addr "$addr" whoami --json 2>/dev/null)
  code=$?
  case $code in 0) health=ok ;; 11) health=expired ;; 126|127) health=no-cli ;; *) health=down ;; esac
  # Whitelist only display fields: never persist an arbitrary CLI response.
  [ -n "$response" ] || response='{}'
  data=$(printf '%s' "$response" | jq -c --arg health "$health" --argjson now "$now" '{episode:(.episode // "?"),secondsLeft:(.secondsLeft // 0),cachedAt:$now,health:$health}' 2>/dev/null)
  [ -n "$data" ] || data=$(jq -nc --arg h "$health" --argjson n "$now" '{episode:"?",secondsLeft:0,cachedAt:$n,health:$h}')
  tmp=$(mktemp "$state/.status.XXXXXX") || tmp=
  if [ -n "$tmp" ]; then
    printf '%s\n' "$data" >"$tmp"
    chmod 600 "$tmp"
    mv -f "$tmp" "$cache"
  fi
fi
printf '%s' "$data" | jq -r --argjson now "$now" '.secondsLeft = ([0, (.secondsLeft - ($now-.cachedAt))] | max) | "World · \(.episode) · \((.secondsLeft/3600|floor))h\((.secondsLeft/60|floor)%60)m · \(if .health == "ok" and .secondsLeft == 0 then "expired" else .health end)"' | tr -d '\n'
if [ -n "$user_command" ]; then
  extra=$(printf '%s' "$input" | sh -c "$user_command" 2>/dev/null)
  code=$?
  if [ "$code" -eq 0 ] && [ -n "$extra" ]; then
    printf ' │ %s' "$extra"
  fi
fi
printf '\n'
`)
	for _, host := range []string{"claude", "codex", "pi"} {
		name := "bin/" + host + "-world"
		script := string(files[name])
		switch host {
		case "claude":
			script = strings.Replace(script, "exec claude ", "exec claude --plugin-dir "+shellQuote(plugin)+" ", 1)
		case "codex":
			// jq's JSON string is also a TOML basic string; escape control characters.
			prefix := "context=$(" + shellQuote(context) + ")\ninstructions=$(printf '%s' \"$context\" | jq -Rs .)\n"
			script = strings.Replace(script, "exec env ", prefix+"exec env ", 1)
			script = strings.Replace(script, " codex \"$@\"", " codex -c \"developer_instructions=$instructions\" \"$@\"", 1)
		case "pi":
			options := " -e " + shellQuote(filepath.Join(dir, "pi", "world-status.ts")) + " --append-system-prompt \"$(" + shellQuote(context) + ")\""
			piHook := filepath.Join(home, ".pi", "agent", "extensions", "herdr-agent-state.ts")
			if info, err := os.Stat(piHook); err == nil && info.Mode().IsRegular() {
				options += " -e " + shellQuote(piHook)
			}
			script = strings.Replace(script, "exec pi ", "exec pi"+options+" ", 1)
		}
		files[name] = []byte(script)
	}
	files["pi/world-status.ts"] = []byte(`import { execFile } from "node:child_process";
export default function (pi: any) {
  pi.on("session_start", async (_event: any, ctx: any) => {
    const data: any = await new Promise(resolve => {
      execFile(` + jsString(worldd) + `, ["--addr", ` + jsString(addr) + `, "whoami", "--json"],
        { timeout: 5000, maxBuffer: 1048576 }, (error, stdout) => {
          try { resolve(error ? null : JSON.parse(stdout)); } catch { resolve(null); }
        });
    });
    if (!data) {
      ctx.ui.setStatus("world", "World · unavailable · run ailang-worldd doctor");
      return;
    }
    const h = Math.floor(data.secondsLeft / 3600);
    const m = Math.floor(data.secondsLeft / 60) % 60;
    ctx.ui.setStatus("world", ` + "`World · ${data.episode} · ${h}h${m}m`" + `);
  });
}
`)
}

func jsString(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
