package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
)

type OperatorFileState struct {
	Mint      OperatorMintReceipt
	Session   SessionRow
	Worktree  string
	StartJSON string
	Step      OperatorStepKey
	Phase     string
}

func readOperatorFiles(ctx context.Context, tx *sql.Tx, id string) (OperatorFileState, approval.OperatorAction, error) {
	var state OperatorFileState
	var action approval.OperatorAction
	mint, err := readOperatorMint(ctx, tx, id, true)
	if err != nil {
		return state, action, err
	}
	state.Mint = mint
	row, err := operatorCASRow(ctx, tx, mint.CAS)
	if err != nil {
		return state, action, err
	}
	action, err = decodeOperatorMintAction(row)
	if err != nil || !action.Options.Headless {
		return state, action, ErrOperatorConflict
	}
	err = tx.QueryRowContext(ctx, `SELECT w.identity, s.start_json, c.credential_id, c.episode_id, c.grants_json,
		c.expires_at, c.created_at, coalesce(c.project,''), coalesce(c.profile_digest,'')
		FROM operator_action_steps w JOIN episode_starts s ON s.episode = ?
		JOIN session_credentials c ON c.episode_id=s.episode AND c.credential_id=?
		WHERE w.action_id = ? AND w.kind = 'worktree' AND w.owner = ? AND w.phase = 'observed' AND w.receipt = ?
		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'worktree') = 1
		AND (SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'mint') = 1`,
		mint.Episode, mint.CredentialID, id, mint.Owner, []byte("created"), id, id).Scan(
		&state.Worktree, &state.StartJSON, &state.Session.CredentialID, &state.Session.EpisodeID, &state.Session.GrantsJSON,
		&state.Session.ExpiresAt, &state.Session.CreatedAt, &state.Session.Project, &state.Session.ProfileDigest)
	if err != nil {
		return state, action, err
	}
	state.Step = OperatorStepKey{Action: mint.CAS, Owner: mint.Owner, Kind: "files"}
	err = tx.QueryRowContext(ctx, `SELECT identity, version, phase FROM operator_action_steps
		WHERE action_id = ? AND kind = 'files' AND owner = ?`, id, mint.Owner).
		Scan(&state.Step.Identity, &state.Step.Version, &state.Phase)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return state, action, err
	}
	var count int
	err = tx.QueryRowContext(ctx, "SELECT count(*) FROM operator_action_steps WHERE action_id = ? AND kind = 'files'", id).Scan(&count)
	if err != nil || (count != 0 && (count != 1 || state.Step.Identity == "")) {
		return state, action, ErrOperatorConflict
	}
	return state, action, nil
}

func (c *Credentials) ReadOperatorFiles(ctx context.Context, id string) (OperatorFileState, approval.OperatorAction, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorFileState{}, approval.OperatorAction{}, err
	}
	tx, err := c.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return OperatorFileState{}, approval.OperatorAction{}, err
	}
	defer tx.Rollback()
	return readOperatorFiles(ctx, tx, id)
}

const (
	OperatorFilesUncertain       = ""
	OperatorFilesPublished       = "published"
	OperatorFilesStoppedAbsent   = "stopped: absent"
	OperatorFilesStoppedRetained = "stopped: retained"
)

// RecordOperatorFiles reserves an empty phase, or observes or disposes the exact reservation.
func (c *Credentials) RecordOperatorFiles(ctx context.Context, reviewed OperatorFileState, value string, proof *approval.OperatorMintProof, clock func() int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	if c.lock == nil || c.lock.dbPath != reviewed.Mint.Store {
		return ErrOperatorConflict
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	current, action, err := readOperatorFiles(ctx, tx, reviewed.Mint.CAS.ID)
	if err != nil {
		return err
	}
	if current != reviewed {
		return ErrOperatorConflict
	}
	key, now := current.Step, clock()
	switch current.Phase {
	case "":
		if value == "" || proof == nil {
			return ErrOperatorInput
		}
		row, err := operatorCASRow(ctx, tx, key.Action)
		if err != nil {
			return err
		}
		if err := admitOperatorFiles(ctx, tx, row, current, action, proof, now); err != nil {
			return err
		}
		key.Identity, key.Version = value, 1
		if err := reserveOperatorStep(ctx, tx, key, now); err != nil {
			return err
		}
	case "reserved":
		switch value {
		case OperatorFilesPublished, OperatorFilesStoppedAbsent, OperatorFilesStoppedRetained:
			if err := observeOperatorStep(ctx, tx, key, []byte(value)); err != nil {
				return err
			}
		case OperatorFilesUncertain:
			res, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET version=version+1 WHERE action_id=?
				AND kind='files' AND identity=? AND owner=? AND version=? AND phase='reserved'`,
				key.Action.ID, key.Identity, key.Owner, key.Version)
			if err := operatorChangedOne(res, err); err != nil {
				return err
			}
		default:
			return ErrOperatorInput
		}
	default:
		return ErrOperatorConflict
	}
	// The step update can succeed before a trigger or other drift invalidates
	// the action. Enforce the independent owner CAS before committing either.
	phase, version, receipt := "reserved", key.Version, []byte(nil)
	if current.Phase == "reserved" {
		version++
		if value != OperatorFilesUncertain {
			phase, receipt = "observed", []byte(value)
		}
	}
	if current.Phase == "" || value == OperatorFilesPublished {
		after, action, err := readOperatorFiles(ctx, tx, key.Action.ID)
		expected := current
		expected.Step = key
		expected.Phase, expected.Step.Version = phase, version
		if err != nil || after != expected {
			return ErrOperatorConflict
		}
		row, err := operatorCASRow(ctx, tx, key.Action)
		if err != nil {
			return err
		}
		if err := admitOperatorFiles(ctx, tx, row, after, action, proof, clock()); err != nil {
			return err
		}
	}
	stopped := current.Phase == "reserved" && (value == OperatorFilesStoppedAbsent || value == OperatorFilesStoppedRetained)
	uncertain := current.Phase == "reserved" && value == OperatorFilesUncertain
	res, err := tx.ExecContext(ctx, `UPDATE operator_actions
		SET state=CASE WHEN ? THEN 'failed' WHEN ? THEN 'uncertain' ELSE state END,
		    owner=CASE WHEN ? THEN NULL ELSE owner END,
		    terminal_at=CASE WHEN ? THEN ? ELSE terminal_at END,
		    generation=generation+CASE WHEN ? THEN 1 ELSE 0 END
		WHERE id=? AND digest=? AND generation=? AND owner=? AND kind='start' AND state IN ('running','uncertain')
		AND EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id=? AND kind='files'
			AND identity=? AND owner=? AND phase=? AND version=? AND receipt IS ?)
		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id=? AND kind NOT IN ('operation','worktree','mint','files'))`,
		stopped, uncertain, stopped, stopped, now, stopped || uncertain,
		key.Action.ID, key.Action.Digest[:], key.Action.Generation, key.Owner,
		key.Action.ID, key.Identity, key.Owner, phase, version, receipt, key.Action.ID)
	if err := operatorChangedOne(res, err); err != nil {
		return err
	}
	return tx.Commit()
}

// admitOperatorFiles validates the persisted mint against the same law that
// made it, including the present approver, session and original expiry.
func admitOperatorFiles(ctx context.Context, tx *sql.Tx, row OperatorActionRow, current OperatorFileState,
	action approval.OperatorAction, proof *approval.OperatorMintProof, now int64) error {
	if now >= row.ExpiresAt || now >= current.Session.ExpiresAt {
		return ErrOperatorExpired
	}
	if row.State != "running" || now < row.CreatedAt || now < current.Session.CreatedAt || len(row.Challenge) != 32 {
		return ErrOperatorConflict
	}
	if proof == nil {
		return ErrOperatorInput
	}
	b, err := proof.Bind(action, [32]byte(row.Challenge), row.ExpiresAt, row.Artifact)
	if err != nil {
		return err
	}
	if row.ApprovedBy != b.Format+":"+b.Name {
		return ErrOperatorConflict
	}
	if err := operatorMintApprover(ctx, tx, proof, b, now); err != nil {
		return err
	}
	catalogue, err := strconv.ParseInt(mintFact(action, "mint_catalogue"), 10, 64)
	if err != nil {
		return ErrOperatorConflict
	}
	in := OperatorMintInput{
		Attempt:      OperatorMintAttempt{Catalogue: catalogue},
		Profile:      approval.ProfileRef{Project: mintFact(action, "mint_project"), Digest: mintFact(action, "mint_profile")},
		CredentialID: current.Mint.CredentialID, TokenPath: current.Mint.TokenPath,
	}
	if json.Unmarshal([]byte(current.StartJSON), &in.Start) != nil || in.Start.Allowance != nil ||
		json.Unmarshal([]byte(mintFact(action, "mint_grants")), &in.Grants) != nil ||
		in.Start.StartedAt != time.Unix(current.Session.CreatedAt, 0).UTC().Format(time.RFC3339) {
		return ErrOperatorConflict
	}
	want, err := operatorMintSession(in, action, current.Session.CreatedAt)
	if err != nil {
		return err
	}
	if want != current.Session {
		return ErrOperatorConflict
	}
	return nil
}
