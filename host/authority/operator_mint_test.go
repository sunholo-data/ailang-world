package authority

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

type mintRig struct {
	operatorStoreRig
	in  store.OperatorMintInput
	row store.OperatorActionRow
	ro  *store.Credentials
	ctx context.Context
}

func operatorMintRig(t *testing.T, passkey bool, edit func(*approval.OperatorAction, *store.OperatorMintInput)) mintRig {
	t.Helper()
	dir := t.TempDir()
	in := store.OperatorMintInput{Attempt: store.OperatorMintAttempt{Catalogue: 1}, Grants: []approval.Grant{{Effect: "exec", Scope: "repo", Budget: 4}},
		Profile: approval.ProfileRef{Project: "world", Digest: hashref.SumSHA256([]byte("profile fixture bytes")).String()}, CredentialID: strings.Repeat("2", 64)}
	r := newOperatorStoreRigAction(t, passkey, func(a *approval.OperatorAction) {
		in.TokenPath = filepath.Join(dir, a.Episode+".session")
		in.Start = store.EpisodeStart{Schema: "world/episode-start/v1", Episode: a.Episode,
			Preset: store.StartPreset{Name: a.Options.Preset, SHA256: "preset"}, Repo: dir,
			Base: store.StartBase{Branch: "main", SHA: "base"}, Branch: "new", Agent: "codex", TTLS: 90,
			Task: &store.StartTask{Path: filepath.Join(dir, a.Episode+".task.md"), SHA256: hashref.SumSHA256([]byte(a.Task)).String()}}
		var err error
		a.Resolution, err = in.Resolution()
		if err != nil {
			t.Fatal(err)
		}
		if edit != nil {
			edit(a, &in)
		}
	})
	ctx := boundedTestContext(t)
	r.confirm(t)
	h := store.OperatorHeartbeat{Owner: store.OperatorOwner{Nonce: "fresh-runner"}, Protocol: 1, Build: "fixture",
		SeenAt: 1010, Readiness: "start", CatalogueGeneration: in.Attempt.Catalogue}
	var err error
	if h.Owner.Generation, err = r.PutOperatorHeartbeat(ctx, h); err != nil {
		t.Fatal(err)
	}
	if err := r.ClaimOperatorAction(ctx, r.cas, h.Owner, 1011); err != nil {
		t.Fatal(err)
	}
	jgExec(t, r.path, `UPDATE operator_actions SET state = 'running'`)
	jgExec(t, r.path, fmt.Sprintf(`INSERT INTO operator_action_steps
		(kind, identity, action_id, phase, version, owner, created_at)
		VALUES ('mint', '%s', '%s', 'reserved', 7, 'fresh-runner', 1012)`, r.cas.ID, r.cas.ID))
	row, err := r.OperatorAction(ctx, r.cas.ID, 1012)
	if err != nil {
		t.Fatal(err)
	}
	in.Attempt = store.OperatorMintAttempt{CAS: store.OperatorCAS{ID: row.ID, Digest: [32]byte(row.Digest), Generation: row.Generation},
		Runner: h.Owner, Catalogue: h.CatalogueGeneration, StepVersion: 7}
	return mintRig{r, in, row, operatorStoreOpen(t, r.path, store.OpenCredentialsReadOnly), ctx}
}

func (r mintRig) seal(t *testing.T) *approval.OperatorMintProof {
	t.Helper()
	p, err := RecheckClaimedOperatorAction(r.ctx, r.ro, r.action, [32]byte(r.row.Challenge), r.row.ExpiresAt, 1013, r.row)
	if err != nil {
		t.Fatal(err)
	}
	return p.OperatorConfirmation().MintProof()
}

type mintSnapshot struct {
	db               [32]byte
	action           store.OperatorActionRow
	sessions         []store.SessionRow
	start            *store.EpisodeStart
	token, task, mcp bool
}

func (r mintRig) snapshot(t *testing.T) mintSnapshot {
	t.Helper()
	b, err := os.ReadFile(r.path)
	if err != nil {
		t.Fatal(err)
	}
	row, err := r.ro.OperatorAction(r.ctx, r.row.ID, 1015)
	if err != nil {
		t.Fatal(err)
	}
	sessions, err := r.ro.ListSessions(r.ctx, "", 20)
	if err != nil {
		t.Fatal(err)
	}
	start, err := r.ro.EpisodeStart(r.ctx, r.row.Episode)
	if err != nil {
		t.Fatal(err)
	}
	exists := func(p string) bool {
		_, err := os.Lstat(p)
		if err != nil && !os.IsNotExist(err) {
			t.Fatal(err)
		}
		return err == nil
	}
	var token string
	for _, fact := range r.action.Resolution {
		if fact.Name == "mint_token_path" {
			token = fact.Value
		}
	}
	dir := filepath.Dir(token)
	return mintSnapshot{sha256.Sum256(b), row, sessions, start, exists(token), exists(filepath.Join(dir, r.action.Episode+".task.md")),
		exists(filepath.Join(dir, r.action.Episode+".mcp.json"))}
}

func mintAt(now int64) func() int64 { return func() int64 { return now } }

func (r mintRig) unchanged(t *testing.T, before mintSnapshot, err, want error) {
	t.Helper()
	if !errors.Is(err, want) || !reflect.DeepEqual(before, r.snapshot(t)) {
		t.Fatalf("refusal %v, want %v with exact state unchanged", err, want)
	}
}

func TestOperatorMintBinding(t *testing.T) {
	for _, tc := range []struct {
		name, sql string
		key       bool
		change    func(*store.OperatorMintInput)
		now       int64
		class     error
		field     string
	}{
		{name: "ttl", change: func(in *store.OperatorMintInput) { in.Start.TTLS++ }, field: "mint_start"},
		{name: "base", change: func(in *store.OperatorMintInput) { in.Start.Base.SHA += "0" }, field: "mint_start"},
		{name: "grant_effect", change: func(in *store.OperatorMintInput) { in.Grants[0].Effect += "0" }, field: "mint_grants"},
		{name: "grant_scope", change: func(in *store.OperatorMintInput) { in.Grants[0].Scope += "0" }, field: "mint_grants"},
		{name: "grant_budget", change: func(in *store.OperatorMintInput) { in.Grants[0].Budget++ }, field: "mint_grants"},
		{name: "project", change: func(in *store.OperatorMintInput) { in.Profile.Project += "0" }, field: "mint_project"},
		{name: "profile", change: func(in *store.OperatorMintInput) {
			in.Profile.Digest = hashref.SumSHA256([]byte("modified profile bytes")).String()
		}, field: "mint_profile"},
		{name: "action_id", change: func(in *store.OperatorMintInput) { in.Attempt.CAS.ID = "unrelated" }, class: store.ErrOperatorMissing},
		{name: "generation", change: func(in *store.OperatorMintInput) { in.Attempt.CAS.Generation++ }, class: store.ErrOperatorConflict},
		{name: "digest", change: func(in *store.OperatorMintInput) { in.Attempt.CAS.Digest[0] ^= 1 }, class: store.ErrOperatorConflict},
		{name: "runner_generation", change: func(in *store.OperatorMintInput) { in.Attempt.Runner.Generation++ }, class: store.ErrOperatorRunner},
		{name: "runner_owner", sql: "UPDATE operator_runner_state SET owner = 'another'", class: store.ErrOperatorRunner},
		{name: "catalogue_generation", sql: "UPDATE operator_runner_state SET catalogue_generation = 2", class: store.ErrOperatorRunner},
		{name: "catalogue_review", sql: "UPDATE operator_runner_state SET catalogue_generation = 2", change: func(in *store.OperatorMintInput) { in.Attempt.Catalogue++ }, field: "mint_catalogue"},
		{name: "action_owner", sql: "UPDATE operator_actions SET owner = 'different'; UPDATE operator_action_steps SET owner = 'different' WHERE kind='mint'", class: store.ErrOperatorConflict},
		{name: "action_phase", sql: "UPDATE operator_actions SET state = 'claimed'", class: store.ErrOperatorConflict},
		{name: "approver_identity", sql: "UPDATE operator_actions SET approved_by = 'ssh:different'", class: store.ErrOperatorConflict},
		{name: "different_challenge", sql: "UPDATE operator_actions SET challenge=zeroblob(32)", field: "challenge"},
		{name: "step_phase", sql: "UPDATE operator_action_steps SET phase='observed', receipt=X'41' WHERE kind='mint'", class: store.ErrOperatorConflict},
		{name: "step_version", change: func(in *store.OperatorMintInput) { in.Attempt.StepVersion++ }, class: store.ErrOperatorConflict},
		{name: "step_owner", sql: "UPDATE operator_action_steps SET owner='another' WHERE kind='mint'", class: store.ErrOperatorConflict},
		{name: "step_identity", sql: "UPDATE operator_action_steps SET identity='another' WHERE kind='mint'", class: store.ErrOperatorConflict},
		{name: "episode_row", sql: "UPDATE operator_actions SET episode_id='another'", class: store.ErrOperatorConflict},
		{name: "record_episode", change: func(in *store.OperatorMintInput) { in.Start.Episode = "another" }, field: "mint_start"},
		{name: "artifact", sql: "UPDATE operator_actions SET artifact = X'41'", class: approval.ErrNotVerified},
		{name: "revoked", sql: "UPDATE approvers SET revoked_at = 1014", class: store.ErrOperatorConflict},
		{name: "future_enrolment", sql: "UPDATE approvers SET enrolled_at = 1015, last_used_at = NULL WHERE format = 'ssh'", class: store.ErrOperatorConflict},
		{name: "approver_name", sql: "UPDATE approvers SET name='renamed' WHERE format='ssh'", class: store.ErrOperatorConflict},
		{name: "approver_format", sql: "UPDATE approvers SET format='webauthn' WHERE format='ssh'", class: store.ErrOperatorConflict},
		{name: "replaced_key", sql: "UPDATE approvers SET public_key = X'41' WHERE format = 'ssh'", class: store.ErrOperatorConflict},
		{name: "passkey_credential", key: true, sql: "UPDATE webauthn_credentials SET credential_id = X'41'", class: store.ErrOperatorConflict},
		{name: "passkey_counter", key: true, sql: "UPDATE webauthn_credentials SET sign_count = 1", class: store.ErrOperatorConflict},
		{name: "passkey_rp", key: true, sql: "UPDATE webauthn_credentials SET rp_id = 'different'", class: store.ErrOperatorConflict},
		{name: "passkey_origin", key: true, sql: "UPDATE webauthn_credentials SET origin = 'https://different'", class: store.ErrOperatorConflict},
		{name: "passkey_backup", key: true, sql: "UPDATE webauthn_credentials SET backup_eligible = 1 - backup_eligible", class: store.ErrOperatorConflict},
		{name: "expiry_boundary", now: 1600, class: store.ErrOperatorExpired},
		{name: "clock_behind", now: 999, class: store.ErrOperatorConflict},
		{name: "overflow", change: func(in *store.OperatorMintInput) { in.Start.TTLS = math.MaxInt64 }, class: store.ErrOperatorInput},
		{name: "zero_ttl", change: func(in *store.OperatorMintInput) { in.Start.TTLS = 0 }, class: store.ErrOperatorInput},
		{name: "no_grants", change: func(in *store.OperatorMintInput) { in.Grants = nil }, class: store.ErrOperatorInput},
		{name: "profile_blank", change: func(in *store.OperatorMintInput) { in.Profile.Project = "" }, class: store.ErrOperatorInput},
		{name: "credential_upper", change: func(in *store.OperatorMintInput) { in.CredentialID = strings.Repeat("A", 64) }, class: store.ErrOperatorInput},
		{name: "profile_hash", change: func(in *store.OperatorMintInput) { in.Profile.Digest = strings.Repeat("1", 64) }, class: store.ErrOperatorInput},
		{name: "credential_hash", change: func(in *store.OperatorMintInput) { in.CredentialID = "sha" }, class: store.ErrOperatorInput},
		{name: "existing_branch", change: func(in *store.OperatorMintInput) { in.Start.BranchExisted = true }, class: store.ErrOperatorConflict},
		{name: "resume_history", change: func(in *store.OperatorMintInput) { in.Start.Resumes = []store.EpisodeResume{{}} }, class: store.ErrOperatorConflict},
		{name: "launch_history", change: func(in *store.OperatorMintInput) { in.Start.Launches = []store.EpisodeLaunch{{}} }, class: store.ErrOperatorConflict},
		{name: "preset_name", change: func(in *store.OperatorMintInput) { in.Start.Preset.Name += "0" }, class: store.ErrOperatorConflict},
		{name: "missing_task", change: func(in *store.OperatorMintInput) { in.Start.Task = nil }, class: store.ErrOperatorConflict},
		{name: "task_digest", change: func(in *store.OperatorMintInput) { in.Start.Task.SHA256 += "0" }, class: store.ErrOperatorConflict},
		{name: "relative_token", change: func(in *store.OperatorMintInput) {
			in.TokenPath = filepath.Base(in.TokenPath)
			in.Start.Task.Path = filepath.Base(in.Start.Task.Path)
		}, class: store.ErrOperatorConflict},
		{name: "unclean_token", change: func(in *store.OperatorMintInput) {
			in.TokenPath = filepath.Dir(in.TokenPath) + "/./" + filepath.Base(in.TokenPath)
		}, class: store.ErrOperatorConflict},
		{name: "token_fact", field: "mint_token_path"},
		{name: "wrong_token", change: func(in *store.OperatorMintInput) { in.TokenPath += "0" }, class: store.ErrOperatorConflict},
		{name: "wrong_task", change: func(in *store.OperatorMintInput) { in.Start.Task.Path += "0" }, class: store.ErrOperatorConflict},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := operatorMintRig(t, tc.key, func(a *approval.OperatorAction, _ *store.OperatorMintInput) {
				if tc.name == "token_fact" {
					a.Resolution[len(a.Resolution)-1].Value += ".different"
				}
			})
			seal := r.seal(t)
			if tc.sql != "" {
				jgExec(t, r.path, tc.sql)
			}
			before := r.snapshot(t)
			if tc.change != nil {
				tc.change(&r.in)
			}
			now := tc.now
			if now == 0 {
				now = 1014
			}
			err := r.CommitOperatorMint(r.ctx, r.in, seal, mintAt(now))
			if tc.field != "" {
				var mismatch *approval.BindingError
				if !errors.As(err, &mismatch) || mismatch.Field != tc.field || !reflect.DeepEqual(before, r.snapshot(t)) {
					t.Fatalf("expected binding %s and unchanged state, got %v", tc.field, err)
				}
			} else {
				r.unchanged(t, before, err, tc.class)
			}
		})
	}
	t.Run("proof", func(t *testing.T) {
		r := operatorMintRig(t, false, nil)
		before := r.snapshot(t)
		for _, signed := range []*approval.SignedOperatorAction{nil, {}, r.proof.OperatorConfirmation()} {
			if signed.MintProof() != nil {
				t.Fatal("non-runner proof acquired a mint seal")
			}
		}
		for _, proof := range []*approval.OperatorMintProof{nil, {}} {
			r.unchanged(t, before, r.CommitOperatorMint(r.ctx, r.in, proof, mintAt(1014)), approval.ErrNotVerified)
		}
	})
	t.Run("start_options", func(t *testing.T) {
		r := operatorMintRig(t, false, func(a *approval.OperatorAction, in *store.OperatorMintInput) {
			a.Options.Headless = false
			a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
		})
		r.unchanged(t, r.snapshot(t), r.CommitOperatorMint(r.ctx, r.in, r.seal(t), mintAt(1014)), store.ErrOperatorConflict)
	})
	t.Run("read_only", func(t *testing.T) {
		r := operatorMintRig(t, false, nil)
		r.unchanged(t, r.snapshot(t), r.ro.CommitOperatorMint(r.ctx, r.in, r.seal(t), mintAt(1014)), store.ErrOperatorConflict)
	})
	t.Run("thin_challenge", func(t *testing.T) {
		r := operatorMintRig(t, false, nil)
		before := r.snapshot(t)
		r.row.Challenge = nil
		raw, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1014))
		if raw != "" {
			t.Fatal("short challenge returned a token")
		}
		r.unchanged(t, before, err, approval.ErrNotVerified)
	})
	t.Run("changed_review", func(t *testing.T) {
		r := operatorMintRig(t, false, nil)
		before := r.snapshot(t)
		other := r.action
		other.Review += " changed"
		raw, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, other, r.row, mintAt(1014))
		var judge *JudgeError
		if raw != "" || !errors.As(err, &judge) || !reflect.DeepEqual(before, r.snapshot(t)) {
			t.Fatal("changed effective action reached mint")
		}
	})
	t.Run("fresh_clock", func(t *testing.T) {
		r := operatorMintRig(t, false, nil)
		before, calls := r.snapshot(t), 0
		raw, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, func() int64 {
			calls++
			if calls == 1 {
				return 1013
			}
			return r.row.ExpiresAt
		})
		if raw != "" || calls != 2 {
			t.Fatal("expected distinct precheck and transactional clock readings, with no token")
		}
		r.unchanged(t, before, err, store.ErrOperatorExpired)
	})
}

func TestOperatorMintCollision(t *testing.T) {
	for _, kind := range []string{"credential", "start"} {
		t.Run(kind, func(t *testing.T) {
			r := operatorMintRig(t, false, nil)
			if kind == "credential" {
				jgExec(t, r.path, fmt.Sprintf("INSERT INTO session_credentials VALUES ('old', '%s', '[]', 1002, 1001, NULL, NULL)", r.action.Episode))
			} else {
				jgExec(t, r.path, fmt.Sprintf("INSERT INTO episode_starts VALUES ('%s', '{}')", r.action.Episode))
			}
			before := r.snapshot(t)
			token, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1014))
			if token != "" || err == nil || !reflect.DeepEqual(before, r.snapshot(t)) ||
				kind == "start" && !strings.Contains(err.Error(), "UNIQUE constraint") ||
				kind == "credential" && !errors.Is(err, store.ErrOperatorConflict) {
				t.Fatalf("collision %s accepted or changed state: %v", kind, err)
			}
		})
	}
}

func TestOperatorMintRollback(t *testing.T) {
	for _, tc := range []struct{ name, sql, want string }{
		{"session_insert", "CREATE TRIGGER fail BEFORE INSERT ON session_credentials BEGIN SELECT RAISE(ABORT, 'injected session failure'); END", "injected session failure"},
		{"start_insert", "CREATE TRIGGER fail BEFORE INSERT ON episode_starts BEGIN SELECT RAISE(ABORT, 'injected start failure'); END", "injected start failure"},
		{"receipt_update", "CREATE TRIGGER fail AFTER UPDATE ON operator_action_steps BEGIN SELECT RAISE(ABORT, 'injected receipt failure'); END", "injected receipt failure"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := operatorMintRig(t, false, nil)
			jgExec(t, r.path, tc.sql)
			before := r.snapshot(t)
			raw, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1014))
			if raw != "" || err == nil || !strings.Contains(err.Error(), tc.want) || !reflect.DeepEqual(before, r.snapshot(t)) {
				t.Fatalf("fault %s did not roll back all writes: %v", tc.name, err)
			}
			jgExec(t, r.path, "DROP TRIGGER fail")
			if _, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1015)); err != nil {
				t.Fatalf("verified rollback did not permit bounded retry: %v", err)
			}
		})
	}
}

func TestOperatorMintReceiptWindow(t *testing.T) {
	t.Run("transaction_separation", TestOperatorMintRollback)
	r := operatorMintRig(t, true, nil)
	before, _, err := r.ro.Approver(r.ctx, strings.TrimPrefix(r.row.ApprovedBy, "webauthn:"))
	if err != nil {
		t.Fatal(err)
	}
	pk, _, err := r.ro.Passkey(r.ctx, before.Name)
	if err != nil {
		t.Fatal(err)
	}
	token, err := MintOperatorStart(r.ctx, r.Credentials, r.ro, r.in, r.action, r.row, mintAt(1014))
	if err != nil || len(token) != 64 {
		t.Fatalf("mint: %v", err)
	}
	// Model a lost handoff by retaining only its hash, before any file writer can run.
	credential := fmt.Sprintf("%x", sha256.Sum256([]byte(token)))
	if err := r.Close(); err != nil {
		t.Fatal(err)
	}
	db := operatorStoreOpen(t, r.path, store.OpenCredentials)
	receipt, err := db.ReadOperatorMint(r.ctx, r.row.ID)
	if err != nil || receipt.CredentialID != credential || receipt.StepVersion != 8 || receipt.Owner != r.row.Owner || receipt.ExpiresAt != 1104 {
		t.Fatalf("commit did not preserve exact receipt: %+v %v", receipt, err)
	}
	got := r.snapshot(t)
	if got.token || got.task || got.mcp || len(got.sessions) != 1 || got.sessions[0].CredentialID != credential ||
		got.sessions[0].CreatedAt != 1014 || got.sessions[0].ExpiresAt != 1104 || got.sessions[0].Project != r.in.Profile.Project ||
		got.sessions[0].ProfileDigest != r.in.Profile.Digest || got.sessions[0].GrantsJSON != `[{"Effect":"exec","Scope":"repo","ExpiresAt":1104,"Budget":4}]` ||
		got.start == nil || got.start.StartedAt != "1970-01-01T00:16:54Z" || got.start.TTLS != 90 || !reflect.DeepEqual(got.action, r.row) {
		t.Fatalf("commit did not preserve exact start, owner and credential: %+v", got)
	}
	after, _, err := r.ro.Approver(r.ctx, before.Name)
	if err != nil || !reflect.DeepEqual(before, after) {
		t.Fatal("mint advanced approver use")
	}
	pkAfter, _, err := r.ro.Passkey(r.ctx, before.Name)
	if err != nil || !reflect.DeepEqual(pk, pkAfter) {
		t.Fatal("mint advanced the passkey counter")
	}
	token, err = MintOperatorStart(r.ctx, db, r.ro, r.in, r.action, r.row, mintAt(1015))
	if token != "" {
		t.Fatal("committed orphan returned a second token")
	}
	r.unchanged(t, got, err, store.ErrOperatorConflict)
}

// TestOperatorMintScope pins the transaction's own scope guards. A signed
// action mints only in the store it names, and only as a start.
func TestOperatorMintScope(t *testing.T) {
	t.Run("store_copy", func(t *testing.T) {
		r := operatorMintRig(t, false, nil)
		seal := r.seal(t)
		other := filepath.Join(t.TempDir(), "copy.db")
		jgExec(t, r.path, fmt.Sprintf("VACUUM INTO '%s'", other))
		copied := operatorStoreOpen(t, other, store.OpenCredentials)
		before := r.snapshot(t)
		err := copied.CommitOperatorMint(r.ctx, r.in, seal, mintAt(1014))
		sessions, listErr := copied.ListSessions(r.ctx, "", 20)
		start, startErr := copied.EpisodeStart(r.ctx, r.action.Episode)
		if listErr != nil || startErr != nil || len(sessions) != 0 || start != nil {
			t.Fatalf("a copied store minted: %v %v %d %v", listErr, startErr, len(sessions), start)
		}
		r.unchanged(t, before, err, store.ErrOperatorConflict)
	})
	t.Run("finish_kind", func(t *testing.T) {
		// No runner row exists, so only the kind refusal can answer conflict.
		s := newOperatorStoreRigAction(t, false, func(a *approval.OperatorAction) {
			a.Kind, a.Task, a.Options.Headless = "finish", "", false
			a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
		})
		ctx := boundedTestContext(t)
		s.confirm(t)
		jgExec(t, s.path, `UPDATE operator_actions SET state = 'running', owner = 'fresh-runner'`,
			fmt.Sprintf(`INSERT INTO operator_action_steps (kind, identity, action_id, phase, version, owner, created_at)
			VALUES ('mint', '%s', '%s', 'reserved', 7, 'fresh-runner', 1012)`, s.cas.ID, s.cas.ID))
		row, err := s.OperatorAction(ctx, s.cas.ID, 1012)
		if err != nil {
			t.Fatal(err)
		}
		r := mintRig{s, store.OperatorMintInput{Attempt: store.OperatorMintAttempt{CAS: s.cas, StepVersion: 7,
			Runner: store.OperatorOwner{Nonce: "fresh-runner"}, Catalogue: 1}}, row, operatorStoreOpen(t, s.path, store.OpenCredentialsReadOnly), ctx}
		before := r.snapshot(t)
		r.unchanged(t, before, r.CommitOperatorMint(ctx, r.in, nil, mintAt(1014)), store.ErrOperatorConflict)
	})
}
