package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

const filesExample = "ailang-worldd service runner reconcile-files oa-0123456789abcdef0123456789abcdef"

func mintOperatorFileFixture(t *testing.T, r *worktreeStartRig) string {
	t.Helper()
	if err := r.start(t, func(string) error { return nil }); err != nil {
		t.Fatal(err)
	}
	r.state(t, "running", "observed", "created", 2, true)
	ctx := boundedTestContext(t)
	ro, err := store.OpenCredentialsReadOnly(runnerDB(r.creds))
	if err != nil {
		t.Fatal(err)
	}
	defer ro.Close()
	row, err := ro.OperatorAction(ctx, r.action.ID, 1012)
	if err != nil {
		t.Fatal(err)
	}
	in := r.in
	in.Attempt = r.at
	in.Attempt.CAS.Generation = row.Generation
	in.Attempt.StepVersion = 1
	var token string
	err = runnerWrite(r.creds, func(db *store.Credentials) error {
		key := store.OperatorStepKey{Action: in.Attempt.CAS, Owner: row.Owner, Kind: "mint", Identity: row.ID, Version: 1}
		if err := db.ReserveOperatorStep(ctx, key, 1012); err != nil {
			return err
		}
		token, err = authority.MintOperatorStart(ctx, db, ro, in, r.action, row, r.src.clock)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(r.f.r.sessions, 0700); err != nil {
		t.Fatal(err)
	}
	return token
}

func publishOperatorFilesFixture(t *testing.T, ctx context.Context, r *worktreeStartRig, token string, clock func() int64, checkpoint func(string) error) error {
	t.Helper()
	held := r.borrow(t)
	defer held.Close()
	return publishOperatorStartFiles(ctx, r.config, held, r.action.ID, token, clock, checkpoint)
}

func TestOperatorFilesReconcile(t *testing.T) {
	t.Run("metadata", TestOperatorFileReadbackMetadata)
	doc, err := os.ReadFile("../../docs/QUICKSTART.md")
	_, example, found := strings.Cut(string(doc), "\nailang-worldd service runner reconcile-files ")
	line, _, _ := strings.Cut(example, "\n")
	example = "ailang-worldd service runner reconcile-files " + line
	if err != nil || !found || example != filesExample {
		t.Fatalf("missing executable QUICKSTART example: %v", err)
	}
	for _, tc := range []struct {
		name, at string
		count    int
	}{
		{"no_mcp", "", 2}, {"all_files", "", 3},
		{"before_reservation", "before reservation", 0}, {"before_first", "before ep.task.md", 0},
		{"between_task_token", "before ep.session", 1}, {"between_token_mcp", "before ep.mcp.json", 2},
		{"before_readback", "before readback", 3}, {"before_observe", "before observe", 3},
		{"competing_file", "", 2}, {"changed_file", "", 3}, {"replaced_after_readback", "", 3}, {"replaced_before_readback", "", 3},
		{"changed_during_review", "before observe", 3}, {"parent_changed_during_review", "before observe", 3},
		{"expired_reopen", "before ep.task.md", 0},
		{"owner_changed", "before observe", 3}, {"version_changed", "before observe", 3},
		{"later_step", "before observe", 3}, {"post_step_version", "before observe", 3}, {"observe_error", "", 3},
		{"publish_owner_drift", "", 3}, {"publish_mint_drift", "", 3}, {"late_key_revoked", "", 3}, {"late_signature_drift", "", 3},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := newWorktreeStartRig(t, tc.name != "no_mcp")
			token := mintOperatorFileFixture(t, r)
			mcpRef := filepath.Join(t.TempDir(), "reference")
			if err := writeMCPConfig(mcpRef, r.config.Address, token); err != nil {
				t.Fatal(err)
			}
			names := [3]string{"ep.task.md", "ep.session", "ep.mcp.json"}
			data := [3]string{r.action.Task, token, string(readStartFixtureFile(t, mcpRef))}
			ctx := presetTestContext(t)
			initial := observeProvision(t, r.f.r, r.creds)
			if tc.name == "no_mcp" {
				err := publishOperatorStartFiles(ctx, r.config, nil, r.action.ID, token, r.src.clock, func(string) error {
					t.Fatal("unheld publication reached an effect boundary")
					return nil
				})
				if !errors.Is(err, store.ErrOperatorConflict) || !reflect.DeepEqual(initial, observeProvision(t, r.f.r, r.creds)) {
					t.Fatalf("missing borrow changed files, Git or credentials: %v", err)
				}
			}
			var baseline, want store.OperatorFileState
			baseline, _, err := readRunnerFiles(ctx, r.creds, r.action.ID)
			if err != nil {
				t.Fatal(err)
			}
			step := func(phase, state, receipt string, version int64) {
				t.Helper()
				var gotPhase, gotState, gotReceipt, owner string
				var gotVersion, generation, starts, sessions, terminal int64
				q := `SELECT coalesce(f.phase,''), a.state, coalesce(f.receipt,''), coalesce(f.version,0),
					coalesce(a.owner,''), a.generation, coalesce(a.terminal_at,0), (SELECT count(*) FROM episode_starts),
					(SELECT count(*) FROM session_credentials) FROM operator_actions a
					LEFT JOIN operator_action_steps f ON f.kind='files' AND f.action_id=a.id`
				err := r.db.QueryRow(q).Scan(&gotPhase, &gotState, &gotReceipt, &gotVersion, &owner, &generation, &terminal, &starts, &sessions)
				wantGeneration, wantOwner, wantTerminal := baseline.Mint.CAS.Generation, baseline.Mint.Owner, int64(0)
				if state == "uncertain" {
					wantGeneration++
				}
				if state == "failed" {
					wantOwner, wantTerminal, wantGeneration = "", r.env.terminal.now(), wantGeneration+version-1
				}
				if err != nil || gotPhase != phase || gotState != state || gotReceipt != receipt || gotVersion != version ||
					owner != wantOwner || terminal != wantTerminal || generation != wantGeneration || starts != 1 || sessions != 1 {
					t.Fatalf("unexpected persisted state %s/%s/%s/%d ownership or mint counts: %v", gotPhase, gotState, gotReceipt, gotVersion, err)
				}
			}
			checkFiles := func() {
				t.Helper()
				for i, name := range names {
					p := filepath.Join(r.f.r.sessions, name)
					st, err := os.Lstat(p)
					if i >= tc.count {
						if !errors.Is(err, os.ErrNotExist) {
							t.Fatalf("unexpected file %s: %v", name, err)
						}
						continue
					}
					if err != nil || st.Mode() != 0600 || uint64(st.Sys().(*syscall.Stat_t).Nlink) != 1 ||
						st.Sys().(*syscall.Stat_t).Uid != uint32(os.Geteuid()) || string(readStartFixtureFile(t, p)) != data[i] {
						t.Fatalf("wrong bytes or metadata for %s: %v", name, err)
					}
				}
				snapshot := observeProvision(t, r.f.r, r.creds)
				if snapshot.worktrees != initial.worktrees || snapshot.branches != initial.branches ||
					!reflect.DeepEqual(snapshot.workspace, initial.workspace) || len(r.writes) != 0 || r.ctl.boots != 0 || r.ctl.stops != 0 {
					t.Fatal("publication or disposition changed Git or the service")
				}
			}
			if tc.name == "all_files" {
				mask := syscall.Umask(0777)
				defer syscall.Umask(mask)
			}
			crash := errors.New("injected file interruption")
			err = publishOperatorFilesFixture(t, ctx, r, token, r.src.clock, func(point string) error {
				if point == "before ep.task.md" {
					step("reserved", "running", "", 1)
					want, _, err = readRunnerFiles(ctx, r.creds, r.action.ID)
					if err != nil || strings.Contains(want.Step.Identity, token) {
						t.Fatal("reservation missing or stores raw credentials")
					}
					var target operatorFileSet
					if json.Unmarshal([]byte(want.Step.Identity), &target) != nil ||
						target.Hashes[0] != serviceHash([]byte(data[0])) || target.Hashes[1] != serviceHash([]byte(data[1])) ||
						(tc.name != "no_mcp" && target.Hashes[2] != serviceHash([]byte(data[2]))) {
						t.Fatal("reservation does not bind actual file bytes")
					}
					lock, err := openRunnerLock(operatorEpisodeLockPath(r.creds, r.action.Episode), 0)
					if err == nil {
						lock.Close()
						t.Fatal("episode lock was released during publication")
					}
				}
				switch {
				case tc.name == "competing_file" && point == "before ep.session":
					data[1] = "foreign"
					seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, names[1]), []byte(data[1]))
				case tc.name == "changed_file" && point == "before readback":
					data[0] = strings.Repeat("x", len(r.action.Task))
					seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, names[0]), []byte(data[0]))
				case tc.name == "replaced_after_readback" && point == "before observe",
					tc.name == "replaced_before_readback" && point == "before readback":
					path := filepath.Join(r.f.r.sessions, names[1])
					old := filepath.Join(r.f.r.sessions, "retained-original")
					if err := os.Rename(path, old); err != nil {
						t.Fatal(err)
					}
					seedStartFixtureFile(t, path, readStartFixtureFile(t, old))
				case point == "before observe":
					var drift string
					switch tc.name {
					case "publish_owner_drift":
						drift = "UPDATE operator_actions SET owner='other'"
					case "publish_mint_drift":
						drift = "UPDATE session_credentials SET expires_at=expires_at+1"
					case "late_key_revoked":
						r.exec(t, "UPDATE approvers SET revoked_at=1012 WHERE name='fixture'")
					case "late_signature_drift":
						r.exec(t, "UPDATE operator_actions SET artifact=X'01'")
					case "observe_error":
						drift = "SELECT RAISE(ABORT,'file observation refused')"
					}
					if drift != "" {
						r.exec(t, "CREATE TRIGGER fail_file_observe AFTER UPDATE OF phase ON operator_action_steps WHEN old.kind='files' BEGIN "+drift+"; END")
					}
					if tc.at == point {
						return crash
					}
				case tc.at == point:
					return crash
				}
				return nil
			})
			checkFiles()
			if tc.name == "no_mcp" || tc.name == "all_files" {
				if err != nil {
					t.Fatal(err)
				}
				step("observed", "running", "published", 2)
				exact := observeProvision(t, r.f.r, r.creds)
				if err := publishOperatorFilesFixture(t, ctx, r, token, r.src.clock, func(string) error { return errors.New("unexpected replay") }); !errors.Is(err, store.ErrOperatorConflict) ||
					!reflect.DeepEqual(exact, observeProvision(t, r.f.r, r.creds)) {
					t.Fatal("replay modified a published start")
				}
				r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { t.Error("prompt for observed row"); return nil, crash }
				if code, msg := r.run(t, example); code != exitUsage || !strings.Contains(msg, store.ErrOperatorMissing.Error()) ||
					!reflect.DeepEqual(exact, observeProvision(t, r.f.r, r.creds)) {
					t.Fatalf("observed row reached disposition: %d %s", code, msg)
				}
				return
			}
			if tc.name == "before_reservation" {
				step("", "running", "", 0)
				before := observeProvision(t, r.f.r, r.creds)
				r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) { t.Error("prompt without reservation"); return nil, crash }
				if c, msg := r.run(t, example); !errors.Is(err, crash) || c != exitUsage || !strings.Contains(msg, store.ErrOperatorMissing.Error()) ||
					!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
					t.Fatalf("wrong refusal before reservation: %v %d %s", err, c, msg)
				}
				return
			}
			version := int64(1)
			if tc.name == "competing_file" || tc.name == "changed_file" || tc.name == "replaced_after_readback" || tc.name == "replaced_before_readback" || strings.HasPrefix(tc.name, "late_") {
				version++
				step("reserved", "uncertain", "", version)
				if err == nil || !strings.Contains(err.Error(), "file publication uncertain") {
					t.Fatalf("wrong mismatch outcome: %v", err)
				}
			} else {
				step("reserved", "running", "", version)
				if tc.name == "observe_error" {
					if err == nil || !strings.Contains(err.Error(), "file observation refused") {
						t.Fatalf("wrong database failure: %v", err)
					}
					r.exec(t, "DROP TRIGGER fail_file_observe")
				} else if tc.name == "publish_owner_drift" || tc.name == "publish_mint_drift" {
					if !errors.Is(err, store.ErrOperatorConflict) {
						t.Fatalf("wrong post-step CAS refusal: %v", err)
					}
					r.exec(t, "DROP TRIGGER fail_file_observe")
				} else if !errors.Is(err, crash) {
					t.Fatalf("lost interruption class: %v", err)
				}
			}
			_, remedy, ok := strings.Cut(err.Error(), "; ")
			if !ok || remedy != filesExample {
				t.Fatalf("wrong file remedy: %v", err)
			}
			before := observeProvision(t, r.f.r, r.creds)
			if tc.name == "between_task_token" || tc.name == "between_token_mcp" || tc.name == "competing_file" || tc.name == "changed_file" {
				r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
					t.Error("unattributed or partial bundle reached attended stop")
					return nil, crash
				}
				code, msg := r.run(t, example)
				if code != exitUsage || !strings.Contains(msg, "file evidence remains uncertain") ||
					!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
					t.Fatalf("partial or foreign evidence was disposed: %d %s", code, msg)
				}
				return
			}
			for _, answer := range []string{"no tty", "n\n", "y\n"} {
				r.env.terminal = testMintEnv(answer, 1012)
				if answer == "no tty" {
					r.env.terminal = noTerminalEnv(1012)
				}
				if tc.name == "expired_reopen" {
					r.env.terminal.now = func() int64 { return baseline.Session.ExpiresAt }
				}
				confirm := r.env.terminal.openTerminal
				if answer == "y\n" {
					r.env.terminal.openTerminal = func() (io.ReadWriteCloser, error) {
						switch tc.name {
						case "changed_during_review":
							data[0] += "changed"
							seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, names[0]), []byte(data[0]))
						case "parent_changed_during_review":
							old := r.f.r.sessions + "-retained"
							if err := os.Rename(r.f.r.sessions, old); err != nil {
								t.Fatal(err)
							}
							if err := os.Mkdir(r.f.r.sessions, 0700); err != nil {
								t.Fatal(err)
							}
							// Same signed bytes and modes in a new parent do not make it the reserved parent.
							for _, name := range names {
								seedStartFixtureFile(t, filepath.Join(r.f.r.sessions, name), readStartFixtureFile(t, filepath.Join(old, name)))
							}
							files := launchFileSnapshot(t, old)
							t.Cleanup(func() {
								if !reflect.DeepEqual(files, launchFileSnapshot(t, old)) {
									t.Error("disposition changed the original directory")
								}
							})
						case "owner_changed":
							r.exec(t, "CREATE TRIGGER owner_after_file AFTER UPDATE OF phase ON operator_action_steps WHEN old.kind='files' BEGIN UPDATE operator_actions SET owner='other'; END")
						case "version_changed":
							r.exec(t, "UPDATE operator_action_steps SET version=version+1 WHERE kind='files'")
						case "later_step":
							r.exec(t, "CREATE TRIGGER later_after_file AFTER UPDATE OF phase ON operator_action_steps WHEN old.kind='files' BEGIN INSERT INTO operator_action_steps(kind,identity,action_id,phase,version,owner,created_at) SELECT 'launch',id,id,'reserved',1,owner,1012 FROM operator_actions; END")
						case "post_step_version":
							r.exec(t, "CREATE TRIGGER version_after_file AFTER UPDATE OF phase ON operator_action_steps WHEN old.kind='files' BEGIN UPDATE operator_action_steps SET version=version+1 WHERE kind='files'; END")
						}
						before = observeProvision(t, r.f.r, r.creds)
						return confirm()
					}
				}
				code, msg := r.run(t, example)
				switch {
				case answer != "y\n":
					if code != exitUsage || !strings.Contains(msg, "Retry at a terminal: "+remedy+"\n") ||
						!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
						t.Fatalf("unattended disposition changed state: %d %s", code, msg)
					}
				case tc.name == "changed_during_review" || tc.name == "parent_changed_during_review" || tc.name == "owner_changed" || tc.name == "version_changed" || tc.name == "later_step" || tc.name == "post_step_version":
					if code != exitUsage || !strings.Contains(msg, "file recovery refused") ||
						!reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
						t.Fatalf("stale disposition changed evidence: %d %s", code, msg)
					}
				default:
					if code != exitOK || !strings.Contains(msg, "File reservation disposed for "+r.action.ID) {
						t.Fatalf("attended disposition: %d %s", code, msg)
					}
					receipt := "stopped: retained"
					if tc.count == 0 {
						receipt = "stopped: absent"
					}
					step("observed", "failed", receipt, version+1)
				}
			}
			checkFiles()
		})
	}
}

func TestOperatorFileReadbackMetadata(t *testing.T) {
	for _, name := range []string{"exact", "digest", "owner", "mode", "symlink", "hardlink", "fifo", "at_limit", "oversize"} {
		t.Run(name, func(t *testing.T) {
			base, err := filepath.EvalSymlinks(t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			if err := os.Chmod(base, 0700); err != nil {
				t.Fatal(err)
			}
			path, alternate := filepath.Join(base, "file"), filepath.Join(base, "alternate")
			set := operatorFileSet{Dir: base, UID: uint64(os.Geteuid())}
			f, _ := openOperatorFilesDir(set)
			if f == nil {
				t.Fatal("fixture directory unavailable")
			}
			defer f.Close()
			body := []byte("from publisher")
			if name == "at_limit" || name == "oversize" {
				body = make([]byte, approval.MaxOperatorTaskBytes+4096)
			}
			installed := body
			if name == "oversize" {
				installed = append(body, 1)
			}
			if _, err := installOperatorFile(f, "alternate", installed); err != nil {
				t.Fatal(err)
			}
			uid := set.UID
			want := serviceHash(body)
			switch name {
			case "owner":
				uid++
				err = os.Rename(alternate, path)
			case "mode":
				err = os.Rename(alternate, path)
				if err == nil {
					err = os.Chmod(path, 0640)
				}
			case "symlink":
				err = os.Symlink(alternate, path)
			case "hardlink":
				err = os.Link(alternate, path)
			case "fifo":
				want, err = serviceHash(nil), syscall.Mkfifo(path, 0600)
			default:
				err = os.Rename(alternate, path)
				if name == "digest" {
					want = serviceHash([]byte("other contents"))
				}
			}
			if err != nil {
				t.Fatal(err)
			}
			before := launchFileSnapshot(t, base)
			got, _ := readOperatorFile(f, "file", want, uid)
			if (got.Status == "exact") != (name == "exact" || name == "at_limit") {
				t.Fatalf("wrong readback class: %s", got.Status)
			}
			if _, err := installOperatorFile(f, "file", []byte("replacement")); !errors.Is(err, os.ErrExist) ||
				!reflect.DeepEqual(before, launchFileSnapshot(t, base)) {
				t.Fatalf("install accepted existing pathname or changed evidence: %v", err)
			}
		})
	}
}

func TestOperatorFilePublicationAdmission(t *testing.T) {
	for _, name := range []string{"missing_mint", "wrong_worktree_receipt", "wrong_worktree_owner", "wrong_worktree_phase", "wrong_mint_owner", "expired_credential",
		"expired_action", "expiry_after_reservation", "credential_expiry_after_reservation", "wrong_grants", "wrong_ttl", "revoked_approver", "revoked_credential",
		"wrong_token", "wrong_approved_address", "reserve_owner_drift", "reserve_mint_drift", "future_mint", "existing_exact_file", "linked_destination", "wrong_destination_mode", "worktree_moved", "lock_held"} {
		t.Run(name, func(t *testing.T) {
			ttl := int64(43200)
			if name == "credential_expiry_after_reservation" {
				ttl = 60
			}
			r := newWorktreeStartRig(t, true, ttl)
			token := mintOperatorFileFixture(t, r)
			held := r.borrow(t)
			defer held.Close()
			borrow := held
			clock := r.src.clock
			ctx := presetTestContext(t)
			want := store.ErrOperatorConflict
			switch name {
			case "missing_mint":
				r.exec(t, "DELETE FROM operator_action_steps WHERE kind='mint'")
				want = nil
			case "wrong_worktree_receipt":
				r.exec(t, "UPDATE operator_action_steps SET receipt=? WHERE kind='worktree'", []byte("stopped: retained"))
				want = nil
			case "wrong_worktree_owner":
				r.exec(t, "UPDATE operator_action_steps SET owner='other' WHERE kind='worktree'")
				want = sql.ErrNoRows
			case "wrong_worktree_phase":
				r.exec(t, "UPDATE operator_action_steps SET phase='reserved',receipt=NULL WHERE kind='worktree'")
				want = sql.ErrNoRows
			case "wrong_mint_owner":
				r.exec(t, "UPDATE operator_action_steps SET owner='other' WHERE kind='mint'")
				want = nil
			case "expired_credential":
				r.exec(t, "UPDATE session_credentials SET expires_at=1012")
				want = store.ErrOperatorExpired
			case "expired_action":
				clock = func() int64 { return 1600 }
				want = store.ErrOperatorExpired
			case "expiry_after_reservation", "credential_expiry_after_reservation":
				now := int64(1012)
				clock = func() int64 { return now }
				err := publishOperatorStartFiles(ctx, r.config, borrow, r.action.ID, token, clock, func(point string) error {
					if point == "before ep.task.md" {
						now = 1600
						if name == "credential_expiry_after_reservation" {
							now = 1012 + ttl
						}
					}
					return nil
				})
				if err == nil || err.Error() != "file publication uncertain; "+filesExample {
					t.Fatalf("expired reserved publication: %v", err)
				}
				st, _, err := readRunnerFiles(ctx, r.creds, r.action.ID)
				files, listErr := os.ReadDir(r.f.r.sessions)
				if err != nil || st.Phase != "reserved" || st.Step.Version != 2 || listErr != nil || len(files) != 0 {
					t.Fatalf("expired reservation changed its files or lost the receipt: %v %v", err, listErr)
				}
				return
			case "wrong_grants":
				r.exec(t, "UPDATE session_credentials SET grants_json='[]'")
			case "wrong_ttl":
				r.exec(t, "UPDATE session_credentials SET expires_at=expires_at+1")
			case "revoked_approver":
				r.exec(t, "UPDATE approvers SET revoked_at=1012 WHERE name='fixture'")
				want = nil
			case "revoked_credential":
				r.exec(t, "DELETE FROM session_credentials")
				want = nil
			case "wrong_token":
				token += "other"
			case "wrong_approved_address":
				r.config.Address = "127.0.0.1:1"
			case "reserve_owner_drift":
				r.exec(t, "CREATE TRIGGER owner_after_file_reserve AFTER INSERT ON operator_action_steps WHEN new.kind='files' BEGIN UPDATE operator_actions SET owner='other'; END")
			case "reserve_mint_drift":
				r.exec(t, "CREATE TRIGGER mint_after_file_reserve AFTER INSERT ON operator_action_steps WHEN new.kind='files' BEGIN UPDATE session_credentials SET expires_at=expires_at+1; END")
			case "future_mint":
				clock = func() int64 { return 1011 }
			case "existing_exact_file":
				seedStartFixtureFile(t, r.in.TokenPath, []byte(token))
			case "linked_destination":
				if err := os.Remove(r.f.r.sessions); err != nil {
					t.Fatal(err)
				}
				decoy := t.TempDir()
				if err := os.Chmod(decoy, 0700); err != nil {
					t.Fatal(err)
				}
				if err := os.Symlink(decoy, r.f.r.sessions); err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() {
					if listed, err := os.ReadDir(decoy); err != nil || len(listed) != 0 {
						t.Errorf("linked directory changed: %v", err)
					}
				})
			case "wrong_destination_mode":
				if err := os.Chmod(r.f.r.sessions, 0750); err != nil {
					t.Fatal(err)
				}
			case "worktree_moved":
				presetGit(t, r.f.r.f, "worktree", "move", r.path(), filepath.Join(t.TempDir(), "other"))
			case "lock_held":
				// The canonical lock stays held while the caller lends a different episode's lock.
				var err error
				borrow, err = store.OpenOperatorEpisodeLock(r.creds, "other")
				if err != nil {
					t.Fatal(err)
				}
				defer borrow.Close()
			}
			before := observeProvision(t, r.f.r, r.creds)
			err := publishOperatorStartFiles(ctx, r.config, borrow, r.action.ID, token, clock, func(point string) error {
				if point != "before reservation" ||
					(name != "expired_credential" && name != "expired_action" && name != "wrong_grants" &&
						name != "wrong_ttl" && name != "revoked_approver" && !strings.HasPrefix(name, "reserve_") && name != "future_mint") {
					t.Errorf("unexpected effect boundary: %s", point)
				}
				return nil
			})
			if err == nil || want != nil && !errors.Is(err, want) || !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
				t.Fatalf("expected refusal with unchanged store, files and worktree: %v", err)
			}
			if name == "lock_held" && !strings.Contains(err.Error(), "wrong episode lock") {
				t.Fatalf("expected wrong episode lock refusal: %v", err)
			}
			var files int
			if err := r.db.QueryRow("SELECT count(*) FROM operator_action_steps WHERE kind='files'").Scan(&files); err != nil || files != 0 {
				t.Fatalf("file step created by refused start: %d %v", files, err)
			}
		})
	}
}

func TestOperatorFileAllowanceOrigin(t *testing.T) {
	r := newWorktreeStartRig(t, true)
	if err := r.start(t, func(string) error { return nil }); err != nil {
		t.Fatal(err)
	}
	start := r.in.Start
	if err := os.Mkdir(r.f.r.sessions, 0700); err != nil {
		t.Fatal(err)
	}
	ctx := boundedTestContext(t)
	now := r.src.clock()
	start.StartedAt = time.Unix(now, 0).UTC().Format(time.RFC3339)
	al := approval.StandingAllowance{Schema: approval.AllowanceSchema, PresetName: start.Preset.Name,
		PresetSHA256: strings.TrimPrefix(start.Preset.SHA256, "sha256:"), Profile: r.in.Profile,
		Repos: []string{start.Repo}, Grants: r.in.Grants, TTLSeconds: start.TTLS, MaxPerDay: 1, ValidDays: 1,
		WorldDB: r.config.World, CredentialsDB: r.creds, WorkspaceRoot: r.f.r.f.ws, ExcludeMissionLoop: true}
	payload, err := approval.CanonicalAllowance(al)
	if err != nil {
		t.Fatal(err)
	}
	var token string
	err = runnerWrite(r.creds, func(db *store.Credentials) error {
		if err := db.GrantAllowance(ctx, store.StandingAllowanceRow{AllowanceID: "sa-0123456789abcdef",
			PresetName: al.PresetName, PresetSHA256: al.PresetSHA256, ProfileProject: al.Profile.Project,
			ProfileDigest: al.Profile.Digest, Payload: payload, Digest: strings.TrimPrefix(serviceHash(payload), "sha256:"),
			Challenge: strings.Repeat("a", 64), MaxPerDay: 1, TTLS: start.TTLS, ApprovedAt: now - 1, ExpiresAt: now + 3600, ApprovedBy: "tty"}); err != nil {
			return err
		}
		var grants []broker.Capability
		for _, g := range r.in.Grants {
			grants = append(grants, broker.Capability{Effect: g.Effect, Scope: g.Scope, Budget: g.Budget})
		}
		token, _, _, _, _, err = authority.MintFromAllowance(ctx, db, start.Episode,
			authority.ProfileBinding{Project: al.Profile.Project, Digest: al.Profile.Digest}, grants, start.TTLS, now, nil,
			authority.AllowanceMintRequest{PresetName: al.PresetName, PresetSHA256: al.PresetSHA256, Repo: start.Repo, Start: start})
		return err
	})
	if err != nil || token == "" {
		t.Fatalf("real allowance draw: %v", err)
	}
	var minted store.EpisodeStart
	var startJSON string
	if err := r.db.QueryRow("SELECT start_json FROM episode_starts WHERE episode=?", start.Episode).Scan(&startJSON); err != nil ||
		json.Unmarshal([]byte(startJSON), &minted) != nil || minted.Allowance == nil {
		t.Fatalf("allowance provenance missing: %v", err)
	}
	before := observeProvision(t, r.f.r, r.creds)
	err = publishOperatorFilesFixture(t, ctx, r, token, r.src.clock, func(point string) error {
		t.Errorf("allowance entered publication at %s", point)
		return nil
	})
	if !errors.Is(err, sql.ErrNoRows) || !reflect.DeepEqual(before, observeProvision(t, r.f.r, r.creds)) {
		t.Fatalf("allowance draw was adopted: %v", err)
	}
	var draws, steps int
	if err := r.db.QueryRow("SELECT (SELECT count(*) FROM allowance_draws), (SELECT count(*) FROM operator_action_steps WHERE kind IN ('mint','files'))").Scan(&draws, &steps); err != nil ||
		draws != 1 || steps != 0 {
		t.Fatalf("allowance or stage evidence changed: %d/%d %v", draws, steps, err)
	}
}
