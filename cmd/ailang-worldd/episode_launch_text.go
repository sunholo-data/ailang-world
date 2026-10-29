package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode/utf16"

	"github.com/sunholo-data/ailang-world/host/broker"
)

const launchChannelNote = "---\nChannel note (trial operator): work only through the World MCP tools (`mcp__world__*`); there is no shell and no native file access, and paths are relative to the episode worktree. `workspace-exec` runs the operator's profiled commands by id: %s. Each call is sandboxed (no network, short timeout). When you are done, finish by summarising the change you made.\n"

const launchHostSandbox = "Your host's own file and shell tools are switched off or read-only on purpose; World's tools above are how you read, write and run checks, and they are not governed by that sandbox."

func launchPrompt(task string, commands []string) string {
	sep := "\n"
	if strings.HasSuffix(task, "\n") {
		sep = ""
	}
	cmds := "(none profiled)"
	if len(commands) > 0 {
		quoted := make([]string, len(commands))
		for i, c := range commands {
			quoted[i] = "`" + c + "`"
		}
		cmds = strings.Join(quoted, ", ")
	}
	return task + sep + "\n" + fmt.Sprintf(launchChannelNote, cmds)
}

// launchContext mirrors trial.world_context, including its absent-field defaults.
func launchContext(body []byte, addr, token string) (string, error) {
	base := addr
	if !strings.HasPrefix(base, "http") {
		base = "http://" + base
	}
	var me map[string]any
	if len(body) > 0 {
		decoder := json.NewDecoder(bytes.NewReader(body))
		decoder.UseNumber()
		if err := decoder.Decode(&me); err != nil {
			return "", err
		}
	}
	str := func(m map[string]any, key, fallback string) string {
		if v, ok := m[key]; ok {
			return fmt.Sprint(v)
		}
		return fallback
	}
	number := func(m map[string]any, key string) int64 {
		n, _ := m[key].(json.Number)
		value, err := n.Int64()
		if err != nil {
			decimal, _ := n.Float64()
			return int64(decimal)
		}
		return value
	}
	out := "You are working inside AILANG World (confined)."
	if ep, _ := me["episode"].(string); ep != "" {
		left := number(me, "secondsLeft")
		out += fmt.Sprintf("\nEpisode %s · worktree %s · session %dh%dm left · daemon %s", ep, str(me, "worktree", "?"), left/3600, left/60%60, str(me, "addr", base))
		if git, ok := me["git"].(map[string]any); ok {
			if e, _ := git["error"].(string); e != "" {
				out += "\nBase unknown: " + e
			} else {
				sha := []rune(str(git, "base", ""))
				out += fmt.Sprintf("\nBase %s (%s, \"%s\") · dev is %d commits ahead of it", string(sha[:min(7, len(sha))]), str(git, "base_date", ""), str(git, "base_subject", ""), number(git, "behind"))
				if dirty, _ := git["dirty"].(bool); dirty {
					out += " · worktree has uncommitted changes"
				}
			}
		}
	}
	out += "\n" + worldToolMap
	if grants, _ := me["grants"].([]any); len(grants) > 0 {
		parts := make([]string, 0, len(grants))
		for _, g := range grants {
			m, _ := g.(map[string]any)
			parts = append(parts, fmt.Sprintf("%s %s/%s", str(m, "effect", "None"), str(m, "left", "None"), str(m, "budget", "None")))
		}
		out += "Grants left: " + strings.Join(parts, ", ") + "\n"
	}
	out += "Every call is recorded. See it: " + strings.TrimRight(base, "/") + "/surface  ·  explain an entry: ailang-worldd why <entry>\n"
	out += launchHostSandbox + "\n" + worldBoundary
	if token != "" && strings.Contains(out, token) {
		return "", fmt.Errorf("world_context: the context would carry the session token")
	}
	return out, nil
}

// pythonJSONString matches json.dumps(s, ensure_ascii=True), without HTML escaping.
func pythonJSONString(s string) string {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	_ = e.Encode(s)
	var out strings.Builder
	for _, r := range strings.TrimSuffix(b.String(), "\n") {
		if r < 127 {
			out.WriteRune(r)
		} else if r <= 0xffff {
			fmt.Fprintf(&out, "\\u%04x", r)
		} else {
			hi, lo := utf16.EncodeRune(r)
			fmt.Fprintf(&out, "\\u%04x\\u%04x", hi, lo)
		}
	}
	return out.String()
}

const claudeLaunchIsolation = `{"enabledPlugins":{"cc-plugin-agents-md@builtin":false,"cc-plugin-telemetry@builtin":false,"cc-plugin-plugin-authoring@builtin":false},"autoMemoryEnabled":false}`

func launchArgv(agent, model, prompt, dir, mcp, addr, context, budget string, grants []broker.Capability) []string {
	if agent == "claude" {
		a := []string{"claude", "-p", prompt}
		if model != "" {
			a = append(a, "--model", model)
		}
		a = append(a, "--output-format", "stream-json", "--verbose", "--no-session-persistence", "--disable-slash-commands",
			"--tools", "", "--strict-mcp-config", "--mcp-config", mcp, "--allowedTools", strings.Join(allowedTools(grants), ","),
			"--setting-sources", "", "--settings", claudeLaunchIsolation)
		if budget != "" {
			a = append(a, "--max-budget-usd", budget)
		}
		if context != "" {
			a = append(a, "--append-system-prompt", context)
		}
		return a
	}
	a := []string{"codex", "exec", "--json", "--ephemeral", "--ignore-user-config", "--skip-git-repo-check"}
	if model != "" {
		a = append(a, "-m", model)
	}
	a = append(a, "-s", "read-only", "-C", dir)
	for _, feature := range []string{"plugins", "apps", "browser_use", "computer_use", "image_generation", "multi_agent", "shell_tool", "unified_exec"} {
		a = append(a, "--disable", feature)
	}
	a = append(a, "-c", "web_search=disabled", "-c", "mcp_servers.world.url=http://"+addr+"/mcp/",
		"-c", "mcp_servers.world.bearer_token_env_var=WORLD_SESSION", "-c", "mcp_servers.world.required=true",
		"-c", "mcp_servers.world.default_tools_approval_mode=approve")
	if context != "" {
		a = append(a, "-c", "developer_instructions="+pythonJSONString(context))
	}
	return append(a, prompt)
}

func redactLaunchArgv(argv []string, prompt, token string) []string {
	out := make([]string, len(argv))
	for i, a := range argv {
		switch {
		case a == prompt:
			out[i] = "<prompt " + serviceHash([]byte(prompt)) + ">"
		case token != "" && strings.Contains(a, token), strings.Contains(strings.ToLower(a), "bearer "):
			out[i] = "<redacted>"
		default:
			out[i] = a
		}
	}
	return out
}
