package main

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorTerminalRevokeExclusion(t *testing.T) {
	fixture := newEpisodeRun(t)
	for _, scenario := range []string{"holder", "durable_owner", "stale_lookup"} {
		t.Run(scenario, func(t *testing.T) {
			r := fixture
			r.f.db = dbPathOf(t)
			cred, err := credentialsFor(r.f.db, "")
			if err != nil {
				t.Fatal(err)
			}
			mint := func(ep string) string {
				t.Helper()
				var out, stderr bytes.Buffer
				if code := runSessionMint([]string{"--db", r.f.db, "--episode", ep, "--grant", "fs.read=repo:1"}, &out, &stderr,
					testMintEnv("y\n", 1000)); code != exitOK {
					t.Fatalf("mint fixture: %d %s", code, &stderr)
				}
				return hex64TokenRegex.FindString(stderr.String())
			}
			id, other := mint("ep"), mint("other")
			var held *os.File
			refusal := "episode lock held"
			if scenario == "durable_owner" {
				claimFixtureEpisode(t, cred, "ep").Close()
				refusal = "active action " + fixtureActionID
			} else if scenario == "holder" {
				held, err = terminalEpisodeLock(cred, "ep", nil)
				if err != nil {
					t.Fatal(err)
				}
				defer held.Close()
			}

			snapshot := func() any {
				return struct {
					db                 string
					repo, ws, sessions map[string]string
				}{serviceFileHash(cred), provisionFiles(t, filepath.Join(r.f.repo, ".git")),
					provisionFiles(t, r.f.ws), provisionFiles(t, r.sessions)}
			}
			before := snapshot()
			check := func(code int, out, stderr *bytes.Buffer) {
				t.Helper()
				if code != exitUsage || !strings.Contains(stderr.String(), refusal) || out.Len() != 0 ||
					!reflect.DeepEqual(before, snapshot()) {
					t.Fatalf("revoke did not refuse without effects: %d %s", code, stderr)
				}
			}
			switch scenario {
			case "stale_lookup":
				refusal = "session changed during revoke"
				acquire := func(path, ep string, inherited *os.File) (*os.File, error) {
					t.Helper()
					if ep != "ep" || path != cred {
						t.Fatalf("wrong lock target: %q %q", path, ep)
					}
					rw, err := store.OpenCredentials(cred)
					if err != nil {
						t.Fatalf("credentials writer acquired before exclusion: %v", err)
					}
					rw.Close()
					credSQL(t, cred, "UPDATE session_credentials SET episode_id='moved' WHERE episode_id='ep'")
					before = snapshot()
					held, err = terminalEpisodeLock(path, ep, inherited)
					return held, err
				}
				var out, stderr bytes.Buffer
				code := sessionRevokeWithLock(boundedTestContext(t), []string{"--db", r.f.db, id}, &out, &stderr, acquire)
				check(code, &out, &stderr)
				if _, err := held.Stat(); !errors.Is(err, os.ErrClosed) {
					t.Fatalf("revoke retained the descriptor: %v", err)
				}
			default:
				for _, args := range [][]string{{id}, {"--episode", "ep"}} {
					var out, stderr bytes.Buffer
					code := runSessionRevoke(append([]string{"--db", r.f.db}, args...), &out, &stderr)
					check(code, &out, &stderr)
				}
			}
			var out, stderr bytes.Buffer
			if code := runSessionRevoke([]string{"--db", r.f.db, other}, &out, &stderr); code != exitOK {
				t.Fatalf("unrelated revoke returned %d: %s", code, &stderr)
			}
			if rows := r.f.rows(t); len(rows) != 1 || rows[0].CredentialID != id {
				t.Fatalf("wrong rows after unrelated revoke: %+v", rows)
			}
		})
	}
}

func TestOperatorTerminalRevokeLegacy(t *testing.T) {
	for _, route := range []string{"id", "episode", "opaque"} {
		t.Run(route, func(t *testing.T) {
			world := dbPathOf(t)
			episode, lockKey := "legacy", "legacy"
			if route == "opaque" {
				episode = "legacy/team"
				lockKey = fmt.Sprintf("-%x", sha256.Sum256([]byte(episode)))
			}
			id := preSplitRow(t, world, episode)
			cred, err := credentialsFor(world, "")
			if err != nil {
				t.Fatal(err)
			}
			args := []string{"--db", world, id}
			if route != "id" {
				args = []string{"--db", world, "--episode", episode}
			}
			held, err := store.OpenOperatorEpisodeLock(cred, lockKey)
			if err != nil {
				t.Fatal(err)
			}
			defer held.Close()
			before := serviceFileHash(world)
			var out, stderr bytes.Buffer
			code := runSessionRevoke(args, &out, &stderr)
			if code != exitUsage || !strings.Contains(stderr.String(), "episode lock held") || out.Len() != 0 ||
				serviceFileHash(world) != before {
				t.Fatalf("migrated before exclusion: %d %s", code, &stderr)
			}
			if _, err := os.Stat(cred); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("split store exists after refusal: %v", err)
			}
			held.Close()
			if code := runSessionRevoke(args, &out, &stderr); code != exitOK || !strings.Contains(out.String(), id) {
				t.Fatalf("legacy revoke: %d %s", code, &stderr)
			}
			ro, err := store.OpenCredentialsReadOnly(cred)
			if err != nil {
				t.Fatal(err)
			}
			defer ro.Close()
			if _, ok, err := ro.ResolveSession(boundedTestContext(t), id); err != nil || ok {
				t.Fatalf("legacy row still exists: %v %v", ok, err)
			}
		})
	}
}

func TestOperatorRunnerPathUID(t *testing.T) {
	home, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, "record")
	if err := os.WriteFile(path, []byte("fixture"), 0600); err != nil {
		t.Fatal(err)
	}
	before := serviceFileHash(path)
	if err := runnerPath(home, path, os.Getuid(), false, 0600); err != nil {
		t.Fatal(err)
	}
	if err := runnerPath(home, path, os.Getuid()+1, false, 0600); err == nil ||
		!strings.Contains(err.Error(), "path ownership mismatch") || serviceFileHash(path) != before {
		t.Fatalf("foreign UID accepted or evidence changed: %v", err)
	}
}
