package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/sunholo-data/ailang-world/host/verbs"
)

const runHelp = `usage: ailang-worldd run <preset> [<task-text>...] [flags]

Start and launch a confined episode in one command using a standing session
allowance (row 224/225). By default, draws from the preset's active standing
allowance (--from-allowance) and launches in headless mode (--launch --headless).

Flags:
  --task <file>           read task prompt from file (alternative to inline text)
  --episode <id>          explicit episode id (default: auto-derived)
  --interactive           launch interactive mode instead of headless
  --prompt                ask for explicit TTY y/N instead of using a standing allowance
  --model <name>          override preset model
  --branch <name>         override episode branch
  --ttl <seconds>         override session TTL
  --presets-dir <dir>     launch presets directory
  --exec-profiles-dir <d> exec profiles directory
  --hosts-dir <dir>       hosts init output directory

Examples:
  ailang-worldd run world-agy "add unit tests for allowance revocation"
  ailang-worldd run docparse --task /tmp/task.md
  ailang-worldd run daneel --interactive
`

// runEpisodeStartFunc is a test seam for runWrapper; production calls runEpisodeStart.
var runEpisodeStartFunc = runEpisodeStart

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
	interactive := fs.Bool("interactive", false, "launch interactive mode instead of headless")
	prompt := fs.Bool("prompt", false, "require explicit TTY y/N instead of standing allowance")
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
	if len(pos) == 0 {
		return fail("<preset> is required\n\n%s", runHelp)
	}
	presetName := pos[0]
	inlineTask := strings.TrimSpace(strings.Join(pos[1:], " "))
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

	if !*interactive && resolvedTask == "" {
		return fail("headless run requires a task (pass inline task text, --task <file>, or --interactive)")
	}

	startArgs := []string{"--preset", presetName, "--launch"}
	if !*interactive {
		startArgs = append(startArgs, "--headless")
	}
	if !*prompt {
		startArgs = append(startArgs, "--from-allowance")
	}
	if resolvedTask != "" {
		startArgs = append(startArgs, "--task", resolvedTask)
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
