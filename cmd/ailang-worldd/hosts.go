package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

const hostsHelp = `usage: ailang-worldd hosts init [--out-dir <dir>] [--addr <url>] [--replace] [--pi-mcp-adapter <dir>]

Write private World launch settings (directories 0700, configs 0600, launchers 0700).
Files: claude-world.settings.json, claude-world.mcp.json, codex/config.toml,
codex/auth.json (symlink only, when ~/.codex/auth.json exists), pi-world-mcp.json,
bin/claude-world, bin/codex-world, bin/pi-world, hosts.json;
bin/world-context, bin/world-statusline, claude-plugin/, codex/AGENTS.md,
codex/hooks.json, pi/world-status.ts. Runtime status caches live under state/.
Changed files refuse unless --replace; unchanged bytes are a no-op.
With an installed service, scripts call its stable installed binary path.
Defaults: --out-dir ~/.ailang-world/hosts; --addr http://127.0.0.1:7644 or service.json's address.
--pi-mcp-adapter defaults to ~/.pi/agent/npm/node_modules/pi-mcp-adapter;
its package.json must name pi.extensions. Missing adapters produce a refusing launcher.
`

type hostsManifest struct {
	Schema             string            `json:"schema"`
	Addr               string            `json:"addr"`
	Files              map[string]string `json:"files"`
	PiMCPAdapter       string            `json:"pi_mcp_adapter"`
	PiMCPAdapterSHA256 string            `json:"pi_mcp_adapter_sha256"`
}

func defaultHostsDir() string {
	h, _ := os.UserHomeDir()
	return filepath.Join(h, ".ailang-world", "hosts")
}

func hostFiles(dir, addr string, adapters ...string) map[string][]byte {
	adapter := defaultPiMCPAdapter()
	if len(adapters) > 0 {
		adapter = adapters[0]
	}
	mcp := strings.TrimRight(addr, "/") + "/mcp/"
	files := map[string][]byte{}
	putJSON := func(name string, v any) { b, _ := json.Marshal(v); files[name] = append(b, '\n') }
	putJSON("claude-world.settings.json", map[string]any{"permissions": map[string]any{
		"deny":                         []string{"Bash", "Edit", "Write", "Read", "Glob", "Grep", "NotebookEdit", "WebFetch", "WebSearch"},
		"disableBypassPermissionsMode": "disable"}})
	putJSON("claude-world.mcp.json", map[string]any{"mcpServers": map[string]any{"world": map[string]any{
		"type": "http", "url": mcp, "headers": map[string]string{"Authorization": "Bearer ${WORLD_SESSION}"}}}})
	files["codex/config.toml"] = []byte(fmt.Sprintf("sandbox_mode = \"read-only\"\napproval_policy = \"never\"\nweb_search = \"disabled\"\n[features]\nhooks = true\nshell_tool = false\nunified_exec = false\n[mcp_servers.world]\nurl = %q\nbearer_token_env_var = \"WORLD_SESSION\"\nrequired = true\ndefault_tools_approval_mode = \"approve\"\n", mcp))
	putJSON("pi-world-mcp.json", map[string]any{"mcpServers": map[string]any{"world": map[string]any{
		"url": mcp, "auth": "bearer", "bearerTokenEnv": "WORLD_SESSION", "lifecycle": "eager"}}})
	commands := map[string]string{
		"claude": "claude --settings " + shellQuote(filepath.Join(dir, "claude-world.settings.json")) + " --strict-mcp-config --mcp-config " + shellQuote(filepath.Join(dir, "claude-world.mcp.json")),
		"codex":  "env CODEX_HOME=" + shellQuote(filepath.Join(dir, "codex")) + " codex",
		"pi":     "pi --no-builtin-tools --no-extensions -e " + shellQuote(adapter) + " --mcp-config " + shellQuote(filepath.Join(dir, "pi-world-mcp.json")),
	}
	for host, command := range commands {
		files["bin/"+host+"-world"] = []byte("#!/bin/sh\nif [ -z \"${WORLD_SESSION:-}\" ]; then\n  cat >&2 <<'WORLD_ERROR'\n" + host + "-world: no World session (WORLD_SESSION is unset).\n  why: World agents need a session you approved.\n  next: start one with `ailang-worldd episode start …`, then export WORLD_SESSION=\"$(cat <session file>)\"\nWORLD_ERROR\n  exit 3\nfi\nexec " + command + " \"$@\"\n")
	}
	guard := ""
	if _, err := piMCPAdapterDigest(adapter); err != nil {
		guard = "cat >&2 <<'WORLD_ADAPTER_ERROR'\n" + piAdapterWhat(adapter) + "\n  why: " + piAdapterWhy + "\n  next: " + piAdapterFix + "\nWORLD_ADAPTER_ERROR\nexit 3\n"
	} else {
		guard = "if [ ! -f " + shellQuote(filepath.Join(adapter, "package.json")) + " ]; then\ncat >&2 <<'WORLD_ADAPTER_ERROR'\n" + piAdapterWhat(adapter) + "\n  why: " + piAdapterWhy + "\n  next: " + piAdapterFix + "\nWORLD_ADAPTER_ERROR\nexit 3\nfi\n"
	}
	files["bin/pi-world"] = []byte(strings.Replace(string(files["bin/pi-world"]), "exec pi ", guard+"exec pi ", 1))
	addAwarenessFiles(files, dir, addr)
	return files
}

func runHosts(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		fmt.Fprint(stderr, hostsHelp)
		return exitUsage
	}
	if args[0] != "init" {
		if args[0] == "--help" || args[0] == "help" {
			fmt.Fprint(stdout, hostsHelp)
			return exitOK
		}
		fmt.Fprint(stderr, hostsHelp)
		return exitUsage
	}
	fs := flag.NewFlagSet("hosts init", flag.ContinueOnError)
	fs.SetOutput(stderr)
	out := fs.String("out-dir", defaultHostsDir(), "output directory")
	addr := fs.String("addr", defaultServiceAddr(), "World URL")
	adapterFlag := fs.String("pi-mcp-adapter", defaultPiMCPAdapter(), "pi MCP extension directory")
	replace := fs.Bool("replace", false, "replace changed files")
	if code, done := parseVerbFlags(fs, args[1:], hostsHelp, stdout, stderr); done {
		return code
	}
	if fs.NArg() != 0 {
		fmt.Fprint(stderr, hostsHelp)
		return exitUsage
	}
	dir, err := filepath.Abs(*out)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	adapter, err := filepath.Abs(*adapterFlag)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return exitUsage
	}
	digest, _ := piMCPAdapterDigest(adapter)
	files := hostFiles(dir, *addr, adapter)
	manifest := hostsManifest{Schema: "world/hosts/v1", Addr: *addr, Files: map[string]string{}, PiMCPAdapter: adapter, PiMCPAdapterSHA256: digest}
	for name, b := range files {
		manifest.Files[name] = fmt.Sprintf("sha256:%x", sha256.Sum256(b))
	}
	b, _ := json.Marshal(manifest)
	files["hosts.json"] = append(b, '\n')
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	fail := func(err error) int { fmt.Fprintf(stderr, "ailang-worldd hosts init: %v\n", err); return exitUsage }
	// Refuse redirected output directories before inspecting any destination.
	for _, d := range hostDirectories(dir) {
		info, e := os.Lstat(d)
		if e == nil && (!info.IsDir() || info.Mode()&os.ModeSymlink != 0) {
			return fail(fmt.Errorf("%s is not a real directory", d))
		}
		if e != nil && !os.IsNotExist(e) {
			return fail(e)
		}
	}
	// Preflight every destination before writing any file. Never follow a destination symlink.
	for _, name := range names {
		p := filepath.Join(dir, name)
		info, e := os.Lstat(p)
		if os.IsNotExist(e) {
			continue
		}
		if e != nil {
			return fail(e)
		}
		if !info.Mode().IsRegular() {
			return fail(fmt.Errorf("%s is not a regular file", p))
		}
		old, e := os.ReadFile(p)
		if e != nil {
			return fail(e)
		}
		if !bytes.Equal(old, files[name]) && !*replace {
			return fail(fmt.Errorf("%s differs; use --replace", p))
		}
	}
	home, _ := os.UserHomeDir()
	auth := filepath.Join(home, ".codex", "auth.json")
	link := filepath.Join(dir, "codex", "auth.json")
	_, authErr := os.Stat(auth) // Metadata only: never read or copy credentials.
	if authErr == nil {
		if info, e := os.Lstat(link); e == nil {
			target, e := os.Readlink(link)
			if info.Mode()&os.ModeSymlink == 0 || e != nil || target != auth {
				if !*replace {
					return fail(fmt.Errorf("%s differs; use --replace", link))
				}
			}
		}
	}
	for _, d := range hostDirectories(dir) {
		if e := os.MkdirAll(d, 0700); e != nil {
			return fail(e)
		}
		if e := os.Chmod(d, 0700); e != nil {
			return fail(e)
		}
	}
	for _, name := range names {
		p := filepath.Join(dir, name)
		mode := os.FileMode(0600)
		if strings.HasPrefix(name, "bin/") {
			mode = 0700
		}
		old, e := os.ReadFile(p)
		if e == nil && bytes.Equal(old, files[name]) {
			continue
		}
		f, e := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
		if e != nil {
			return fail(e)
		}
		_, e = f.Write(files[name])
		closeErr := f.Close()
		if e != nil {
			return fail(e)
		}
		if closeErr != nil {
			return fail(closeErr)
		}
		if e = os.Chmod(p, mode); e != nil {
			return fail(e)
		}
	}
	if authErr == nil {
		target, e := os.Readlink(link)
		if e != nil || target != auth {
			if _, e = os.Lstat(link); e == nil {
				if e = os.Remove(link); e != nil {
					return fail(e)
				}
			}
			if e = os.Symlink(auth, link); e != nil {
				return fail(e)
			}
		}
	} else if os.IsNotExist(authErr) {
		fmt.Fprintln(stdout, "! Codex will ask you to log in: ~/.codex/auth.json does not exist.")
	} else {
		return fail(authErr)
	}
	if r, ok := homeServiceRecord(); ok {
		path, _ := os.Executable()
		if serviceFileHash(path) != serviceFileHash(r.Binary) {
			fmt.Fprintf(stdout, "! the scripts call the installed World binary %s, not this one (%s)\n", r.Binary, path)
		}
	}
	fmt.Fprintf(stdout, "World hosts ready in %s\n", dir)
	return exitOK
}
