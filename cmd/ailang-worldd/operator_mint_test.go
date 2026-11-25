package main

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

type recoveryFixture struct {
	*operatorServiceFixture
	db                    *sql.DB
	path, token, worktree string
	receipt               store.OperatorMintReceipt
}

const reconcileMintExample = "ailang-worldd service runner reconcile-mint oa-0123456789abcdef0123456789abcdef"

func newRecoveryFixture(t *testing.T) recoveryFixture {
	t.Helper()
	f := operatorService(t)
	f.mustInstall(t)
	f.env.terminal.now = func() int64 { return 1020 }
	f.writes = nil
	ep := "minted"
	sessions := filepath.Join(f.home, "mint-sessions")
	if err := os.Mkdir(sessions, 0700); err != nil {
		t.Fatal(err)
	}
	in := store.OperatorMintInput{
		TokenPath: filepath.Join(sessions, ep+".session"),
		Profile:   approval.ProfileRef{Project: "world", Digest: serviceHash([]byte("profile fixture bytes"))},
		Grants:    []approval.Grant{{Effect: "fs.read", Scope: "repo", Budget: 2}},
	}
	in.Attempt.Catalogue = 1
	var a approval.OperatorAction
	fd := claimFixtureEpisodeAction(t, f.config.Credentials, ep, func(action *approval.OperatorAction) {
		action.Options.Preset = "sol"
		in.Start = store.EpisodeStart{Schema: "world/episode-start/v1", Episode: ep, Preset: store.StartPreset{Name: "sol"}, TTLS: 20,
			Task: &store.StartTask{Path: filepath.Join(sessions, ep+".task.md"), SHA256: serviceHash([]byte(action.Task))}}
		var err error
		action.Resolution, err = in.Resolution()
		if err != nil {
			t.Fatal(err)
		}
		a = *action
		a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
	})
	defer fd.Close()
	db, err := sql.Open("sqlite", runnerDB(f.config.Credentials))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	r := recoveryFixture{operatorServiceFixture: f, db: db, path: in.TokenPath, worktree: filepath.Join(f.home, "workspaces", ep)}
	r.exec(t, "UPDATE operator_actions SET state='running' WHERE id=?", a.ID)
	r.exec(t, `INSERT INTO operator_action_steps (kind,identity,action_id,phase,version,owner,created_at)
		VALUES ('mint',?,?,'reserved',7,'test-runner',1012)`, a.ID, a.ID)
	ctx := boundedTestContext(t)
	err = runnerWrite(f.config.Credentials, func(c *store.Credentials) error {
		row, err := c.OperatorAction(ctx, a.ID, 1012)
		if err != nil {
			return err
		}
		h, err := c.OperatorHeartbeat(ctx)
		if err != nil {
			return err
		}
		in.Attempt.CAS = store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation}
		in.Attempt.Runner, in.Attempt.StepVersion = h.Owner, 7
		r.token, err = authority.MintOperatorStart(ctx, c, c, in, a, row, func() int64 { return 1012 })
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	r.receipt, err = readRunnerMint(f.config.Credentials, a.ID)
	if err != nil {
		t.Fatal(err)
	}
	r.exec(t, `INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at)
		VALUES ('unrelated', ?, '[]', 1100, 1013)`, ep)
	if err := os.Mkdir(r.worktree, 0700); err != nil {
		t.Fatal(err)
	}
	for _, p := range []string{filepath.Join(r.worktree, "keep"), in.Start.Task.Path, filepath.Join(sessions, ep+".mcp.json")} {
		if err := os.WriteFile(p, []byte("preserve these bytes"), 0600); err != nil {
			t.Fatal(err)
		}
	}
	return r
}

func (r recoveryFixture) exec(t *testing.T, query string, args ...any) {
	t.Helper()
	if _, err := r.db.ExecContext(boundedTestContext(t), query, args...); err != nil {
		t.Fatal(err)
	}
}

// rewriteAction replaces the stored action and its digest together, so only
// the resolver's own store and token-path identity checks can refuse it.
func (r recoveryFixture) rewriteAction(t *testing.T, edit func(*approval.OperatorAction)) {
	t.Helper()
	var raw []byte
	if err := r.db.QueryRowContext(boundedTestContext(t), "SELECT action FROM operator_actions").Scan(&raw); err != nil {
		t.Fatal(err)
	}
	a, err := approval.DecodeOperatorAction(raw)
	if err != nil {
		t.Fatal(err)
	}
	edit(&a)
	if raw, err = approval.CanonicalOperatorAction(a); err != nil {
		t.Fatal(err)
	}
	digest := sha256.Sum256(raw)
	r.exec(t, "UPDATE operator_actions SET action=?, digest=?", raw, digest[:])
}

func (r recoveryFixture) run(t *testing.T, command string) (int, string) {
	t.Helper()
	args := strings.Fields(command)
	if len(args) < 2 || args[0] != "ailang-worldd" || args[1] != "service" {
		t.Fatalf("unexpected remedy %q", command)
	}
	return r.operatorServiceFixture.run(args[2:]...)
}

func TestOperatorMintRecovery(t *testing.T) {
	for _, tc := range []struct{ name, want string }{
		{"revocation", "Orphaned mint revoked"},
		{"unsafe_parent", "cannot establish private token directory"},
		{"missing_lock", "episode recovery lock held or unavailable"},
		{"wrong_receipt_phase", "mint evidence unavailable"},
		{"unowned", "operator action conflicts"},
		{"busy", "episode recovery lock held or unavailable"},
		{"stale_owner", "operator action conflicts"},
		{"stale_step", "operator action conflicts"},
		{"appeared", "token appeared after review"},
		{"adopt", "Matching private token retained"},
		{"expired", "matching token expired"},
		{"later_effect", "later effect reserved"},
		{"wrong_step_owner", "mint evidence unavailable"},
		{"start_changed", "mint start record differs"},
		{"action_changed", "operator action conflicts"},
		{"receipt_changed", "operator action conflicts"},
		{"missing_credential", "mint start or credential unavailable"},
		{"abort_rollback", "injected abort failure"},
		{"store_identity", "mint token path or store identity differs"},
		{"token_name", "mint token path or store identity differs"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newRecoveryFixture(t)
			id := r.receipt.CAS.ID
			writeToken := func() {
				t.Helper()
				if err := writeTokenFile(r.path, r.token); err != nil {
					t.Fatal(err)
				}
			}
			switch tc.name {
			case "revocation":
				r.exec(t, "UPDATE operator_actions SET state='uncertain'")
			case "action_changed":
				r.exec(t, `UPDATE operator_actions SET action=CAST(json_set(CAST(action AS TEXT),'$.review','changed') AS BLOB)`)
			case "unsafe_parent":
				if err := os.Chmod(filepath.Dir(r.path), 0750); err != nil {
					t.Fatal(err)
				}
			case "missing_lock":
				if err := os.Remove(operatorEpisodeLockPath(r.config.Credentials, r.receipt.Episode)); err != nil {
					t.Fatal(err)
				}
			case "unowned":
				r.exec(t, "UPDATE operator_actions SET state='failed', owner=NULL, terminal_at=1019")
			case "wrong_receipt_phase":
				r.exec(t, "UPDATE operator_action_steps SET phase='reserved', receipt=NULL WHERE kind='mint'")
			case "busy":
				held, err := openRunnerLock(operatorEpisodeLockPath(r.config.Credentials, r.receipt.Episode), 0)
				if err != nil {
					t.Fatal(err)
				}
				defer held.Close()
			case "adopt", "expired":
				writeToken()
				if tc.name == "expired" {
					r.env.terminal.now = func() int64 { return r.receipt.ExpiresAt }
				}
			case "later_effect":
				r.exec(t, `INSERT INTO operator_action_steps (kind,identity,action_id,phase,version,owner,created_at)
					VALUES ('launch',?,?,'reserved',1,'test-runner',1013)`, id, id)
			case "wrong_step_owner":
				r.exec(t, "UPDATE operator_action_steps SET owner='other' WHERE kind='mint'")
			case "start_changed":
				r.exec(t, "UPDATE episode_starts SET start_json='{}'")
			case "receipt_changed":
				r.exec(t, "UPDATE operator_action_steps SET receipt=X'41' WHERE kind='mint'")
			case "missing_credential":
				r.exec(t, "DELETE FROM session_credentials WHERE credential_id=?", r.receipt.CredentialID)
			case "abort_rollback":
				r.exec(t, "CREATE TRIGGER failed_abort AFTER UPDATE ON operator_actions BEGIN SELECT RAISE(ABORT,'injected abort failure'); END")
			case "store_identity", "token_name":
				r.rewriteAction(t, func(a *approval.OperatorAction) {
					if tc.name == "store_identity" {
						a.Store = filepath.Join(filepath.Dir(a.Store), "other.db")
						return
					}
					for i := range a.Resolution {
						if a.Resolution[i].Name == "mint_token_path" {
							a.Resolution[i].Value = filepath.Join(filepath.Dir(r.path), "other.session")
						}
					}
				})
			}
			before := serviceFileHash(r.config.Credentials)
			files := provisionFiles(t, filepath.Dir(r.path))
			worktree := provisionFiles(t, r.worktree)
			if tc.name == "revocation" {
				// The returned token is deliberately never published. Invoke the
				// command printed by the refusal, byte for byte, after a no-tty run.
				r.env.terminal.openTerminal = noTerminalEnv(1020).openTerminal
				code, out := r.run(t, reconcileMintExample)
				_, command, ok := strings.Cut(out, "Retry at a terminal: ")
				if code != exitUsage || !ok || strings.TrimSpace(command) != reconcileMintExample || !strings.Contains(out, "no controlling terminal") ||
					serviceFileHash(r.config.Credentials) != before {
					t.Fatalf("no tty or changed remedy: %d %s", code, out)
				}
				r.env.terminal.openTerminal = testMintEnv("n\n", 1020).openTerminal
				if code, out := r.run(t, strings.TrimSpace(command)); code != exitUsage || !strings.Contains(out, "aborted (not confirmed)") ||
					serviceFileHash(r.config.Credentials) != before {
					t.Fatalf("declined recovery: %d %s", code, out)
				}
				r.env.terminal.openTerminal = testMintEnv("y\n", 1020).openTerminal
			}
			if tc.name == "stale_owner" || tc.name == "stale_step" || tc.name == "appeared" {
				r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
					switch tc.name {
					case "stale_owner":
						r.exec(t, "UPDATE operator_actions SET owner='new-owner', generation=generation+1 WHERE id=?", id)
						r.exec(t, "UPDATE operator_action_steps SET owner='new-owner' WHERE kind='mint'")
						before = serviceFileHash(r.config.Credentials)
					case "stale_step":
						r.exec(t, "UPDATE operator_action_steps SET version=version+1 WHERE kind='mint'")
						before = serviceFileHash(r.config.Credentials)
					case "appeared":
						writeToken()
						files = provisionFiles(t, filepath.Dir(r.path))
					}
					return testMintEnv("y\n", 1020).openTerminal()
				}
			}
			code, out := r.run(t, reconcileMintExample)
			wantCode := exitUsage
			if tc.name == "revocation" || tc.name == "adopt" {
				wantCode = exitOK
			}
			if code != wantCode || !strings.Contains(out, tc.want) || strings.Contains(out, r.token) ||
				!reflect.DeepEqual(files, provisionFiles(t, filepath.Dir(r.path))) || !reflect.DeepEqual(worktree, provisionFiles(t, r.worktree)) ||
				len(r.writes) != 0 || r.ctl.boots != 0 || r.ctl.stops != 0 {
				t.Fatalf("recovery %s: %d %s", tc.name, code, out)
			}
			if tc.name != "revocation" {
				if before != serviceFileHash(r.config.Credentials) {
					t.Fatal("recovery refusal or matched-token read modified the store")
				}
				return
			}
			var state, owner string
			var version, generation, minted, other, starts int64
			var receipt []byte
			for _, q := range []struct {
				sql   string
				slots []any
			}{
				{"SELECT state, coalesce(owner,''), generation FROM operator_actions", []any{&state, &owner, &generation}},
				{"SELECT version, receipt FROM operator_action_steps WHERE kind='mint'", []any{&version, &receipt}},
				{"SELECT count(*) FROM session_credentials WHERE credential_id != 'unrelated'", []any{&minted}},
				{"SELECT count(*) FROM session_credentials WHERE credential_id = 'unrelated'", []any{&other}},
				{"SELECT count(*) FROM episode_starts", []any{&starts}},
			} {
				if err := r.db.QueryRowContext(boundedTestContext(t), q.sql).Scan(q.slots...); err != nil {
					t.Fatal(err)
				}
			}
			if state != "failed" || owner != "" || generation != r.receipt.CAS.Generation+1 || version != r.receipt.StepVersion ||
				string(receipt) != r.receipt.CredentialID+":"+r.receipt.StartDigest || minted != 0 || other != 1 || starts != 1 {
				t.Fatalf("orphan disposition changed evidence or unrelated authority: %s %s %d %d %d %d", state, owner, version, minted, other, starts)
			}
			requireProvisionLockReleased(t, r.config.Credentials, r.receipt.Episode)
		})
	}
}

func TestOperatorMintTokenReadback(t *testing.T) {
	for _, name := range []string{"uid", "mode", "link", "symlink", "symlink_valid", "directory", "fifo", "bytes", "length", "short"} {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			p := filepath.Join(dir, "token")
			token := strings.Repeat("a", 64)
			uid, mode := os.Getuid(), os.FileMode(0600)
			var err error
			switch name {
			case "symlink", "symlink_valid":
				target := filepath.Join(dir, "target")
				if name == "symlink_valid" {
					if err := os.WriteFile(target, []byte(token), 0600); err != nil {
						t.Fatal(err)
					}
				}
				err = os.Symlink(target, p)
			case "directory":
				err = os.Mkdir(p, 0700)
			case "fifo":
				err = syscall.Mkfifo(p, 0600)
			default:
				data := token
				switch name {
				case "uid":
					uid++
				case "mode":
					mode = 0644
				case "bytes":
					data = strings.Repeat("b", 64)
				case "length":
					data += "a"
				case "short":
					token = token[:63]
					data = token
				}
				err = os.WriteFile(p, []byte(data), mode)
				if err == nil && name == "link" {
					err = os.Link(p, filepath.Join(dir, "other"))
				}
			}
			if err != nil {
				t.Fatal(err)
			}
			before := provisionFiles(t, dir)
			present, err := readMintToken(p, uid, fmt.Sprintf("%x", sha256.Sum256([]byte(token))))
			if present || err == nil || !reflect.DeepEqual(before, provisionFiles(t, dir)) {
				t.Fatalf("unsafe file evidence changed or accepted: present=%v err=%v", present, err)
			}
			// The file type is refused as a type, before any mode or owner check.
			if (name == "directory" || name == "fifo") && err.Error() != "mint token is not a regular file" {
				t.Fatalf("%s refused for the wrong reason: %v", name, err)
			}
		})
	}
}

func TestOperatorMintRecoveryArgs(t *testing.T) {
	e := serviceEnv{goos: "darwin", label: worldServiceLabel}
	for _, args := range [][]string{{"reconcile-mint"}, {"reconcile-mint", "wrong"}, {"reconcile-mint", fixtureActionID, "extra"}} {
		var msg bytes.Buffer
		if code := runServiceRunner(args, &msg, &msg, e); code != exitUsage || !strings.Contains(msg.String(), "expected reconcile-mint ACTION-ID") {
			t.Fatalf("recovery syntax: %v: %d %s", args, code, &msg)
		}
	}
}

func TestQuickstartRunnerMintRecovery(t *testing.T) {
	b, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	const sidecar = "Terminal episode operations use private mode-0600 sidecars named\n`<creds>.operator-episode-<ep>.lock` beside the canonical credentials database.\nKeep these sidecars in place across stop and uninstall."
	if !strings.Contains(string(b), sidecar) || !strings.Contains(string(b), "\n"+reconcileMintExample+"\n") {
		t.Fatal("missing the canonical sidecar description or exact attended command")
	}
}
