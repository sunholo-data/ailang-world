package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sunholo-data/ailang-world/host/verbs"
)

const runHelp = `usage: ailang-worldd run [<preset>] [<task-text>...] [flags]

Start and launch a confined episode in one command using a standing session
allowance (row 224/225).

If <preset> is omitted (or the first word is not a known preset name), the
preset matching the current working directory's repository is auto-detected.
With no task and no --headless flag, launches in interactive mode (--interactive);
with inline task text or --task <file>, launches in headless mode (--headless).
Draws from the preset's active standing allowance by default (--from-allowance).

Flags:
  --task <file>           read task prompt from file (alternative to inline text)
  --episode <id>          explicit episode id (default: auto-derived)
  --interactive           launch interactive mode (overrides preset launch mode)
  --headless              launch headless mode (requires task text or --task)
  --prompt                ask for explicit TTY y/N instead of using a standing allowance
  --agent <name>          override agent harness (agy|claude|codex)
  --model <name>          override preset model
  --branch <name>         override episode branch
  --ttl <seconds>         override session TTL
  --presets-dir <dir>     launch presets directory
  --exec-profiles-dir <d> exec profiles directory
  --hosts-dir <dir>       hosts init output directory

Examples:
  ailang-worldd run                                  # interactive in current repo
  ailang-worldd run "add unit tests for foo"         # headless in current repo
  ailang-worldd run docparse "fix PDF parser"        # explicit preset from anywhere
  ailang-worldd run daneel --interactive             # explicit interactive preset
  ailang-worldd run --agent claude                   # interactive claude in current repo
`

// runEpisodeStartFunc and runGetwdFunc are test seams for runWrapper.
var (
	runEpisodeStartFunc = runEpisodeStart
	runGetwdFunc        = os.Getwd
)

type presetMatchCandidate struct {
	name    string
	repo    string
	agent   string
	matched int // repo path length (prefer most specific match)
}

func canonicalDirPath(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return filepath.Clean(dir)
	}
	if real, err := filepath.EvalSymlinks(abs); err == nil {
		return filepath.Clean(real)
	}
	return filepath.Clean(abs)
}

func presetExists(presetsDir, name string) bool {
	if name == "" || strings.ContainsAny(name, `/\`) {
		return false
	}
	info, err := os.Stat(filepath.Join(presetsDir, name+".json"))
	return err == nil && info.Mode().IsRegular()
}

// detectPresetForDir finds the best preset in presetsDir whose repo contains cwd.
// When multiple presets match the same repo, prefers preferredAgent (or "agy" by
// default), then shorter preset name, then alphabetical order.
func detectPresetForDir(cwd, presetsDir, preferredAgent string) (string, error) {
	realCwd := canonicalDirPath(cwd)
	entries, err := os.ReadDir(presetsDir)
	if err != nil {
		return "", fmt.Errorf("read presets dir %s: %w", presetsDir, err)
	}
	var matches []presetMatchCandidate
	for _, ent := range entries {
		if ent.IsDir() || !strings.HasSuffix(ent.Name(), ".json") {
			continue
		}
		raw, err := os.ReadFile(filepath.Join(presetsDir, ent.Name()))
		if err != nil {
			continue
		}
		var p struct {
			Name  string `json:"name"`
			Repo  string `json:"repo"`
			Agent string `json:"agent"`
		}
		if json.Unmarshal(raw, &p) != nil || p.Repo == "" {
			continue
		}
		name := strings.TrimSuffix(ent.Name(), ".json")
		if p.Name != "" {
			name = p.Name
		}
		realRepo := canonicalDirPath(p.Repo)
		if realCwd == realRepo || strings.HasPrefix(realCwd, realRepo+string(filepath.Separator)) {
			matches = append(matches, presetMatchCandidate{
				name:    name,
				repo:    realRepo,
				agent:   p.Agent,
				matched: len(realRepo),
			})
		}
	}
	if len(matches) == 0 {
		return "", fmt.Errorf("no launch preset in %s matches current directory %s", presetsDir, cwd)
	}
	wantAgent := preferredAgent
	if wantAgent == "" {
		wantAgent = "agy"
	}
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].matched != matches[j].matched {
			return matches[i].matched > matches[j].matched
		}
		mi := matches[i].agent == wantAgent
		mj := matches[j].agent == wantAgent
		if mi != mj {
			return mi
		}
		if len(matches[i].name) != len(matches[j].name) {
			return len(matches[i].name) < len(matches[j].name)
		}
		return matches[i].name < matches[j].name
	})
	return matches[0].name, nil
}

func runWrapper(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	for _, a := range args {
		if a == "--help" || a == "-h" || a == "-help" {
			fmt.Fprint(stdout, runHelp)
			return exitOK
		}
		if a == "--" {
			break
		}
	}

	fs := flag.NewFlagSet("ailang-worldd run", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	taskFile := fs.String("task", "", "task prompt file")
	episodeID := fs.String("episode", "", "explicit episode id")
	interactiveFlag := fs.Bool("interactive", false, "launch interactive mode (overrides preset launch mode)")
	headlessFlag := fs.Bool("headless", false, "launch headless mode")
	prompt := fs.Bool("prompt", false, "require explicit TTY y/N instead of standing allowance")
	agent := fs.String("agent", "", "override preset agent harness (agy|claude|codex)")
	model := fs.String("model", "", "override preset model")
	branch := fs.String("branch", "", "override episode branch")
	ttl := fs.String("ttl", "", "override session TTL")
	presetsDir := fs.String("presets-dir", "", "launch presets directory")
	profilesDir := fs.String("exec-profiles-dir", "", "exec profiles directory")
	hostsDir := fs.String("hosts-dir", "", "hosts init output directory")

	pos, err := verbs.ParseInterleaved(fs, args)
	if err != nil {
		if errors.Is(err, flag.ErrHelp) {
			fmt.Fprint(stdout, runHelp)
			return exitOK
		}
		fmt.Fprint(stderr, runHelp)
		return exitUsage
	}
	fail := func(format string, a ...any) int {
		fmt.Fprintf(stderr, "ailang-worldd run: "+format+"\n", a...)
		return exitUsage
	}
	if *interactiveFlag && *headlessFlag {
		return fail("--interactive and --headless are mutually exclusive")
	}

	effectivePresetsDir := *presetsDir
	if effectivePresetsDir == "" {
		effectivePresetsDir = defaultPresetsDir()
	}

	var presetName, inlineTask string
	if len(pos) > 0 && presetExists(effectivePresetsDir, pos[0]) {
		presetName = pos[0]
		inlineTask = strings.TrimSpace(strings.Join(pos[1:], " "))
	} else {
		cwd, cwdErr := runGetwdFunc()
		if cwdErr != nil {
			return fail("resolve working directory: %v", cwdErr)
		}
		detected, detErr := detectPresetForDir(cwd, effectivePresetsDir, *agent)
		if detErr != nil {
			if len(pos) == 0 {
				return fail("%v; pass a <preset> name (see 'ailang-worldd preset list')", detErr)
			}
			return fail("%q is not a preset in %s, and %v", pos[0], effectivePresetsDir, detErr)
		}
		presetName = detected
		inlineTask = strings.TrimSpace(strings.Join(pos, " "))
	}

	if *taskFile != "" && inlineTask != "" {
		return fail("cannot pass both --task <file> and inline task text")
	}

	resolvedTask := *taskFile
	if inlineTask != "" {
		tmp, err := os.CreateTemp("", "world-task-*.md")
		if err != nil {
			return fail("create temp task file: %v", err)
		}
		tmpPath := tmp.Name()
		defer func() { _ = os.Remove(tmpPath) }()
		if err := tmp.Chmod(0o600); err != nil {
			_ = tmp.Close()
			return fail("chmod temp task file: %v", err)
		}
		if _, err := io.WriteString(tmp, inlineTask+"\n"); err != nil {
			_ = tmp.Close()
			return fail("write temp task file: %v", err)
		}
		if err := tmp.Close(); err != nil {
			return fail("close temp task file: %v", err)
		}
		resolvedTask = tmpPath
	}

	// Default mode: interactive when no task is provided (unless --headless was explicitly asked),
	// headless when a task is provided (unless --interactive was explicitly asked).
	useInteractive := *interactiveFlag || (!*headlessFlag && resolvedTask == "")
	if !useInteractive && resolvedTask == "" {
		return fail("headless run requires a task (pass inline task text, --task <file>, or drop --headless)")
	}

	startArgs := []string{"--preset", presetName, "--launch"}
	if useInteractive {
		startArgs = append(startArgs, "--interactive")
	} else {
		startArgs = append(startArgs, "--headless")
	}
	if !*prompt {
		startArgs = append(startArgs, "--from-allowance")
	}
	if resolvedTask != "" {
		startArgs = append(startArgs, "--task", resolvedTask)
	}
	if *agent != "" {
		startArgs = append(startArgs, "--agent", *agent)
	}
	if *model != "" {
		startArgs = append(startArgs, "--model", *model)
	}
	if *branch != "" {
		startArgs = append(startArgs, "--branch", *branch)
	}
	if *ttl != "" {
		startArgs = append(startArgs, "--ttl", *ttl)
	}
	if *presetsDir != "" {
		startArgs = append(startArgs, "--presets-dir", *presetsDir)
	}
	if *profilesDir != "" {
		startArgs = append(startArgs, "--exec-profiles-dir", *profilesDir)
	}
	if *hostsDir != "" {
		startArgs = append(startArgs, "--hosts-dir", *hostsDir)
	}
	if *episodeID != "" {
		startArgs = append(startArgs, *episodeID)
	}

	return runEpisodeStartFunc(startArgs, stdout, stderr, env)
}
