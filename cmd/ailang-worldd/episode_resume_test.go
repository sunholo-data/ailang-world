package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func resumeFixture(t *testing.T) (presetFixture, string, *store.EpisodeStart) {
	t.Helper()
	f := newPresetFixture(t)
	var out, errOut bytes.Buffer
	term := &promptTerm{in: strings.NewReader("y\n")}
	if code := runEpisodeStart(f.args("--task", f.task), &out, &errOut, presetEnv(term)); code != 0 {
		t.Fatal(code, errOut.String())
	}
	ep := "p-2610281352"
	return f, ep, readPresetStart(t, f, ep)
}

func TestEpisodeResumeAC11(t *testing.T) {
	f, ep, original := resumeFixture(t)
	oldRows := f.r.f.rows(t)
	tokenPath := filepath.Join(f.r.sessions, ep+".session")
	oldToken, err := os.ReadFile(tokenPath)
	if err != nil {
		t.Fatal(err)
	}
	mcpPath := filepath.Join(f.r.sessions, ep+".mcp.json")
	if err := os.WriteFile(mcpPath, oldToken, 0o600); err != nil {
		t.Fatal(err)
	}
	trees := presetGit(t, f.r.f, "worktree", "list", "--porcelain")
	term := &promptTerm{in: strings.NewReader("y\n")}
	env := presetEnv(term)
	env.now = func() int64 { return presetClock() + 43201 }
	var out, errOut bytes.Buffer
	if code := runEpisodeStart(f.args("--resume", ep), &out, &errOut, env); code != 0 {
		t.Fatal(code, errOut.String(), term.out.String())
	}
	archived := tokenPath + ".expired-" + strconv.FormatInt(env.now(), 10)
	bytes, err := os.ReadFile(archived)
	if err != nil || string(bytes) != string(oldToken) {
		t.Fatal("old token not retained", err)
	}
	oldMCP := mcpPath + ".expired-" + strconv.FormatInt(env.now(), 10)
	mcpBytes, err := os.ReadFile(oldMCP)
	if err != nil || string(mcpBytes) != string(oldToken) {
		t.Fatal("old MCP config not retained", err)
	}
	for _, path := range []string{archived, tokenPath, oldMCP} {
		info, err := os.Stat(path)
		if err != nil || info.Mode().Perm() != 0o600 {
			t.Fatal(path, info, err)
		}
	}
	newToken, err := os.ReadFile(tokenPath)
	if err != nil || string(newToken) == string(oldToken) {
		t.Fatal("no new token", err)
	}
	rows := f.r.f.rows(t)
	if len(rows) != len(oldRows)+1 || rows[len(rows)-1].CredentialID != serviceHash(newToken)[len("sha256:"):] {
		t.Fatal(rows)
	}
	if after := presetGit(t, f.r.f, "worktree", "list", "--porcelain"); after != trees {
		t.Fatal("worktree list changed")
	}
	if !strings.Contains(term.out.String(), "Use the existing worktree") || !strings.Contains(term.out.String(), "unchanged") {
		t.Fatal(term.out.String())
	}
	start := readPresetStart(t, f, ep)
	if len(start.Resumes) != 1 || start.Resumes[0].ResumedFrom != oldRows[len(oldRows)-1].CredentialID {
		t.Fatal(start.Resumes)
	}
	t.Log("resume prompt:\n" + term.out.String())
	t.Logf("resume entry: %+v", start.Resumes[0])
	start.Resumes = nil
	if !reflect.DeepEqual(original, start) {
		t.Fatal("start fields changed")
	}
}

func TestEpisodeResumeDeclineAndRefusals(t *testing.T) {
	for _, kind := range []string{"decline", "live", "no-rows", "profile", "worktree", "lock", "task", "grant", "positional"} {
		t.Run(kind, func(t *testing.T) {
			f, ep, _ := resumeFixture(t)
			now := presetClock() + 43201
			extra := []string{"--resume", ep}
			want := ""
			switch kind {
			case "live":
				now = presetClock()
				want = "has a live session"
			case "no-rows":
				path, _ := credentialsFor(f.r.f.db, "")
				c, err := store.OpenCredentials(path)
				if err != nil {
					t.Fatal(err)
				}
				for _, row := range f.r.f.rows(t) {
					if err := c.RevokeSession(presetTestContext(t), row.CredentialID); err != nil {
						t.Fatal(err)
					}
				}
				_ = c.Close()
				want = "no credential to resume"
			case "profile":
				path := filepath.Join(f.r.profiles, "proj.json")
				b, _ := os.ReadFile(path)
				_ = os.WriteFile(path, append(b, '\n'), 0o600)
				want = "changed since episode"
			case "worktree":
				if err := worktree.Remove(presetTestContext(t), f.r.f.git, f.r.f.repo, filepath.Join(f.r.f.ws, ep)); err != nil {
					t.Fatal(err)
				}
				want = "worktree is gone"
			case "lock":
				_ = os.WriteFile(filepath.Join(f.r.sessions, ep+".launch.lock"), []byte(fmt.Sprint(os.Getpid())), 0o600)
				want = "launch lock"
			case "task":
				extra = append(extra, "--task", f.task)
				want = "--resume takes only"
			case "grant":
				extra = append(extra, "--grant", "Workspace.Read=worktree:1")
				want = "--resume takes only"
			case "positional":
				extra = append(extra, "another")
				want = "--resume takes only"
			}
			before := launchFileSnapshot(t, f.r.sessions)
			rows := f.r.f.rows(t)
			start := readPresetStart(t, f, ep)
			var opened bool
			env := spyEnv(&opened)
			env.now = func() int64 { return now }
			term := &promptTerm{in: strings.NewReader("n\n")}
			if kind == "decline" {
				env = presetEnv(term)
				env.now = func() int64 { return now }
				want = "aborted"
			}
			var out, errOut bytes.Buffer
			if code := runEpisodeStart(f.args(extra...), &out, &errOut, env); code == 0 || !strings.Contains(errOut.String(), want) || opened {
				t.Fatal(code, opened, errOut.String())
			}
			if !reflect.DeepEqual(before, launchFileSnapshot(t, f.r.sessions)) || !reflect.DeepEqual(rows, f.r.f.rows(t)) || !reflect.DeepEqual(start, readPresetStart(t, f, ep)) {
				t.Fatal("refusal changed state")
			}
			if kind == "live" && !strings.Contains(errOut.String(), "episode launch") {
				t.Fatal(errOut.String())
			}
		})
	}
}

func TestEpisodeResumeNoStartAndUnknownFlag(t *testing.T) {
	f := newPresetFixture(t)
	path, _ := credentialsFor(f.r.f.db, "")
	c, err := store.OpenCredentials(path)
	if err != nil {
		t.Fatal(err)
	}
	_ = c.Close()
	for _, extra := range [][]string{{"--resume", "ep"}, {"--resume", "ep", "--yes"}} {
		var opened bool
		var out, errOut bytes.Buffer
		want := "has no start record"
		if len(extra) > 2 {
			want = "--resume takes only"
		}
		if code := runEpisodeStart(f.args(extra...), &out, &errOut, spyEnv(&opened)); code == 0 || opened || !strings.Contains(errOut.String(), want) {
			t.Fatal(code, opened, errOut.String())
		}
	}
}

func TestEpisodeResumeChangedPresetKeepsGrants(t *testing.T) {
	f, ep, before := resumeFixture(t)
	oldRows := f.r.f.rows(t)
	_ = os.WriteFile(f.file, append(f.data, '\n'), 0o600)
	term := &promptTerm{in: strings.NewReader("y\n")}
	env := presetEnv(term)
	env.now = func() int64 { return presetClock() + 43201 }
	var out, errOut bytes.Buffer
	if code := runEpisodeStart(f.args("--resume", ep), &out, &errOut, env); code != 0 || !strings.Contains(term.out.String(), "changed since start") {
		t.Fatal(code, errOut.String(), term.out.String())
	}
	after := readPresetStart(t, f, ep)
	if after.Preset != before.Preset || len(after.Resumes) != 1 || after.Resumes[0].PresetFileSHA256 == before.Preset.SHA256 {
		t.Fatal(after)
	}
	rows := f.r.f.rows(t)
	// provisionEpisode only re-dates grant expiry.
	if strings.ReplaceAll(rows[len(rows)-1].GrantsJSON, strconv.FormatInt(env.now()+43200, 10), "EXP") != strings.ReplaceAll(oldRows[len(oldRows)-1].GrantsJSON, strconv.FormatInt(presetClock()+43200, 10), "EXP") {
		t.Fatal("grants changed")
	}
}
