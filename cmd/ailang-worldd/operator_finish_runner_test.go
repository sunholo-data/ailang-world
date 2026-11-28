package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorFinishReadiness(t *testing.T) {
	for _, phase := range []string{"confirmed", "resume", "terminal"} {
		t.Run(phase, func(t *testing.T) {
			var confirm []bool
			if phase == "confirmed" {
				confirm = []bool{true}
			}
			var change func(*store.OperatorFinishReview)
			var want error = operatorGitHubAuth(6)
			wantTrace := "remote-read\n"
			if phase == "terminal" {
				change = func(m *store.OperatorFinishReview) { m.Checks[0].Unsandboxed = true }
				want, wantTrace = errOperatorFinishTerminal, ""
			}
			f, work, now := operatorLocalFinishFixture(t, change, confirm...)
			path := filepath.Dir(work.manifest.Store)
			tools, err := os.MkdirTemp(path, "approved-")
			if err != nil {
				t.Fatal(err)
			}
			trace := filepath.Join(tools, "trace")
			c := runnerConfig{Home: path, Credentials: work.manifest.Store, UID: os.Getuid(), ToolDirs: []string{tools},
				Executables: make(map[string]string), FinishRemotes: []string{work.manifest.Remote.URL}}
			writeFinishTestTool(t, &c, tools, "gh", "case \"$1\" in auth|label) ;; repo) printf 'WRITE\\n';; *) exit 83;; esac\n")
			writeFinishTestTool(t, &c, tools, "git", "printf 'remote-read\\n' >> "+shellQuote(trace)+"\nexit 1\n")
			before := serviceFileHash(work.manifest.Store)
			ctx := activationTestContext(t)
			h, err := work.read.OperatorHeartbeat(ctx)
			if err != nil {
				t.Fatal(err)
			}
			err = (serviceEnv{}).consumeFinish(ctx, c, h, operatorStartSource{clock: func() int64 { return *now }})
			row, readErr := work.read.OperatorAction(ctx, work.prep.row.ID, *now)
			marker, markerErr := os.ReadFile(trace)
			beforeHead := strings.TrimSpace(f.gitOutput("-C", f.wt, "rev-parse", "HEAD"))
			j := operatorFinishJournal(t, work)
			if !errors.Is(err, want) || readErr != nil || !reflect.DeepEqual(work.prep.row, row) ||
				before != serviceFileHash(work.manifest.Store) || string(marker) != wantTrace || markerErr != nil && !errors.Is(markerErr, os.ErrNotExist) ||
				beforeHead != f.before || j.CommitPhase != "" || j.OpenCheck || len(j.Checks) != 0 || f.checks != 0 ||
				!h.Available(h.Build, "start", *now) {
				t.Fatalf("readiness crossed %s admission: %v, row %s, marker=%q (%v)", phase, err, row.State, marker, markerErr)
			}
		})
	}
}
