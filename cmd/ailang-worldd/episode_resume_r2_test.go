package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/worktree"
)

func resumeR2(t *testing.T, f presetFixture, ep string, now int64) *store.EpisodeStart {
	t.Helper()
	term := &promptTerm{in: strings.NewReader("y\n")}
	env := presetEnv(term)
	env.now = func() int64 { return now }
	var out, errOut bytes.Buffer
	if code := runEpisodeStart(f.args("--resume", ep), &out, &errOut, env); code != 0 {
		t.Fatal(code, errOut.String(), term.out.String())
	}
	return readPresetStart(t, f, ep)
}

func TestEpisodeResumeRefusesAfterRevokeByID(t *testing.T) {
	f, ep, _ := resumeFixture(t)
	start := resumeR2(t, f, ep, presetClock()+43201)
	id := start.Resumes[0].CredentialID
	if id == "" {
		t.Fatal("resume did not record its credential")
	}
	var out, errOut bytes.Buffer
	if code := runSessionRevoke([]string{"--db", f.r.f.db, id}, &out, &errOut); code != 0 {
		t.Fatal(code, errOut.String())
	}
	before := f.r.f.rows(t)
	var opened bool
	env := spyEnv(&opened)
	env.now = func() int64 { return presetClock() + 86402 }
	errOut.Reset()
	if code := runEpisodeStart(f.args("--resume", ep), &out, &errOut, env); code == 0 || opened || !strings.Contains(errOut.String(), "no credential to resume (revoked or finished)") || !strings.Contains(errOut.String(), "revoke ends an episode") {
		t.Fatal(code, opened, errOut.String())
	}
	if len(before) != 1 || !reflect.DeepEqual(before, f.r.f.rows(t)) || len(readPresetStart(t, f, ep).Resumes) != 1 {
		t.Fatal("revoked episode minted a third credential")
	}
}

func TestEpisodeResumeChainUsesCurrentCredential(t *testing.T) {
	f, ep, _ := resumeFixture(t)
	first := resumeR2(t, f, ep, presetClock()+43201)
	second := resumeR2(t, f, ep, presetClock()+86402)
	if len(second.Resumes) != 2 || second.Resumes[1].ResumedFrom != first.Resumes[0].CredentialID || second.Resumes[1].CredentialID == "" {
		t.Fatal(second.Resumes)
	}
}

func TestEpisodeResumeExpiryBoundary(t *testing.T) {
	for _, kind := range []string{"at-expiry", "after-expiry"} {
		t.Run(kind, func(t *testing.T) {
			f, ep, _ := resumeFixture(t)
			expiry := f.r.f.rows(t)[0].ExpiresAt
			if kind == "after-expiry" {
				resumeR2(t, f, ep, expiry+1)
				return
			}
			before := launchFileSnapshot(t, f.r.sessions)
			var opened bool
			env := spyEnv(&opened)
			env.now = func() int64 { return expiry }
			var out, errOut bytes.Buffer
			if code := runEpisodeStart(f.args("--resume", ep), &out, &errOut, env); code == 0 || opened || !strings.Contains(errOut.String(), "(F12)") {
				t.Fatal(code, opened, errOut.String())
			}
			if !reflect.DeepEqual(before, launchFileSnapshot(t, f.r.sessions)) || len(f.r.f.rows(t)) != 1 {
				t.Fatal("boundary minted or wrote")
			}
		})
	}
}

func TestEpisodeResumeFailureAfterMintRestores(t *testing.T) {
	f, ep, _ := resumeFixture(t)
	originalID := f.r.f.rows(t)[0].CredentialID
	session := filepath.Join(f.r.sessions, ep+".session")
	mcp := filepath.Join(f.r.sessions, ep+".mcp.json")
	old, err := os.ReadFile(session)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(mcp, []byte("original MCP"), 0o600); err != nil {
		t.Fatal(err)
	}
	before := launchFileSnapshot(t, f.r.sessions)
	oldHook := resumeAfterRename
	t.Cleanup(func() { resumeAfterRename = oldHook })
	resumeAfterRename = func() error {
		if _, err := os.Stat(session); !os.IsNotExist(err) {
			t.Fatal("hook ran before rename", err)
		}
		return errors.New("post-mint fixture failure")
	}
	term := &promptTerm{in: strings.NewReader("y\n")}
	env := presetEnv(term)
	env.now = func() int64 { return presetClock() + 43201 }
	var out, errOut bytes.Buffer
	if code := runEpisodeStart(f.args("--resume", ep), &out, &errOut, env); code == 0 || !strings.Contains(errOut.String(), "post-mint fixture failure") {
		t.Fatal(code, errOut.String())
	}
	if !reflect.DeepEqual(before, launchFileSnapshot(t, f.r.sessions)) {
		t.Fatal("old files or modes not restored")
	}
	start := readPresetStart(t, f, ep)
	rows := f.r.f.rows(t)
	if len(rows) != 1 || rows[0].CredentialID != originalID || len(start.Resumes) != 1 || !start.Resumes[0].Failed || start.Resumes[0].CredentialID == "" {
		t.Fatal(rows, start.Resumes)
	}
	got, _ := os.ReadFile(session)
	if string(got) != string(old) {
		t.Fatal("old token changed")
	}
	resumeAfterRename = oldHook
	start = resumeR2(t, f, ep, env.now()+1)
	if len(start.Resumes) != 2 || start.Resumes[1].ResumedFrom != originalID {
		t.Fatal(start.Resumes)
	}
	if _, err := os.Stat(session + ".expired-" + strconv.FormatInt(env.now(), 10)); !os.IsNotExist(err) {
		t.Fatal("failed attempt archive remains")
	}
}

func TestEpisodeResumeGuardsR2(t *testing.T) {
	for _, kind := range []string{"preset-name", "daemon-down", "nonregular-session", "nonregular-mcp", "wrong-branch", "foreign-worktree"} {
		t.Run(kind, func(t *testing.T) {
			f, ep, start := resumeFixture(t)
			args := f.args("--resume", ep)
			dir := filepath.Join(f.r.f.ws, ep)
			want := ""
			switch kind {
			case "preset-name":
				var p launchPreset
				if err := json.Unmarshal(f.data, &p); err != nil {
					t.Fatal(err)
				}
				// finishPreset accepts the alias, but the pinned name must not.
				p.Name = "other"
				data, err := json.Marshal(p)
				if err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(filepath.Join(f.dir, "other.json"), data, 0o600); err != nil {
					t.Fatal(err)
				}
				args[1] = "other"
				want = "does not match start preset"
			case "daemon-down":
				h := matchingHealth(t, f.r.f, f.r.profiles, f.profile)
				h.DBPath = f.r.f.db
				h.ToolRef = ""
				stubDaemon(t, h, nil)
				want = "tool binary"
			case "nonregular-session", "nonregular-mcp":
				ext := ".session"
				if kind == "nonregular-mcp" {
					ext = ".mcp.json"
				}
				path := filepath.Join(f.r.sessions, ep+ext)
				_ = os.Remove(path)
				if err := os.Mkdir(path, 0o700); err != nil {
					t.Fatal(err)
				}
				want = "regular file"
			case "wrong-branch":
				if _, err := worktree.Output(presetTestContext(t), f.r.f.git, []string{"-C", dir, "checkout", "-q", "-b", "wrong"}); err != nil {
					t.Fatal(err)
				}
				want = "not on start branch"
			case "foreign-worktree":
				if err := worktree.Remove(presetTestContext(t), f.r.f.git, f.r.f.repo, dir); err != nil {
					t.Fatal(err)
				}
				foreign := newEpisodeFixture(t)
				if err := worktree.Add(presetTestContext(t), foreign.git, foreign.repo, dir, start.Branch, "HEAD"); err != nil {
					t.Fatal(err)
				}
				want = "is not a worktree of"
			}
			before := launchFileSnapshot(t, f.r.sessions)
			rows := f.r.f.rows(t)
			var opened bool
			env := spyEnv(&opened)
			env.now = func() int64 { return presetClock() + 43201 }
			var out, errOut bytes.Buffer
			if code := runEpisodeStart(args, &out, &errOut, env); code == 0 || opened || !strings.Contains(strings.ToLower(errOut.String()), want) {
				t.Fatal(code, opened, errOut.String())
			}
			if !reflect.DeepEqual(before, launchFileSnapshot(t, f.r.sessions)) || !reflect.DeepEqual(rows, f.r.f.rows(t)) || !reflect.DeepEqual(start, readPresetStart(t, f, ep)) {
				t.Fatal("guard wrote state")
			}
		})
	}
}
