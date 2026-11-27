package store

import (
	"bytes"
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

// OperatorMintAttempt identifies an already reserved mint step.
type OperatorMintAttempt struct {
	CAS         OperatorCAS
	Runner      OperatorOwner
	Catalogue   int64
	StepVersion int64
}

// OperatorMintInput is the effective mint reconstructed by the caller. Clock
// is read after the transaction starts. No token bytes enter the store.
type OperatorMintInput struct {
	Attempt      OperatorMintAttempt
	Start        EpisodeStart
	Grants       []approval.Grant
	Profile      approval.ProfileRef
	CredentialID string
	TokenPath    string
}

// Resolution supplies the mint fields of the reviewed action and is repeated
// on the effective inputs at commit. StartedAt comes from the commit clock.
func (in OperatorMintInput) Resolution() ([]approval.OperatorResolved, error) {
	start := in.Start
	start.StartedAt = ""
	body, err := json.Marshal(start)
	if err != nil {
		return nil, err
	}
	grants, _ := json.Marshal(in.Grants) // approval.Grant contains only strings and an integer.
	return []approval.OperatorResolved{
		{Name: "mint_catalogue", Value: strconv.FormatInt(in.Attempt.Catalogue, 10)},
		{Name: "mint_grants", Value: string(grants)},
		{Name: "mint_profile", Value: in.Profile.Digest},
		{Name: "mint_project", Value: in.Profile.Project},
		{Name: "mint_start", Value: fmt.Sprintf("%x", sha256.Sum256(body))},
		{Name: "mint_token_path", Value: in.TokenPath},
	}, nil
}

func mintFact(a approval.OperatorAction, name string) string {
	for _, fact := range a.Resolution {
		if fact.Name == name {
			return fact.Value
		}
	}
	return ""
}

func operatorMintApprover(ctx context.Context, tx *sql.Tx, proof *approval.OperatorMintProof, b approval.OperatorApprovalBinding, now int64) error {
	var current bool
	if err := tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM approvers WHERE
		name = ? AND format = ? AND public_key = ? AND revoked_at IS NULL AND enrolled_at <= ?)`,
		b.Name, b.Format, []byte(b.PublicKey), now).Scan(&current); err != nil {
		return err
	}
	if !current {
		return ErrOperatorConflict
	}
	if id, count, ok := proof.PasskeyUse(); ok {
		if err := tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM webauthn_credentials WHERE name = ?
			AND credential_id = ? AND rp_id = ? AND origin = ? AND backup_eligible = ? AND sign_count >= ?)`,
			b.Name, id, b.RPID, b.Origin, b.BackupEligible, count).Scan(&current); err != nil {
			return err
		}
		if !current {
			return ErrOperatorConflict
		}
	}
	return nil
}

// CommitOperatorMint atomically records the session, start, and observed mint
// receipt. It is an inert transaction seam pending paired-process admission.
func (c *Credentials) CommitOperatorMint(ctx context.Context, in OperatorMintInput, proof *approval.OperatorMintProof, clock func() int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	if c.lock == nil {
		return ErrOperatorConflict
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	now := clock()
	try := in.Attempt
	row, err := operatorCASRow(ctx, tx, try.CAS)
	if err != nil {
		return err
	}
	if now >= row.ExpiresAt {
		return ErrOperatorExpired
	}
	if now < row.CreatedAt || row.State != "running" || row.Owner != try.Runner.Nonce || row.Kind != "start" {
		return ErrOperatorConflict
	}
	var current bool
	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM operator_runner_state WHERE
		generation = ? AND owner = ? AND catalogue_generation = ?)`,
		try.Runner.Generation, try.Runner.Nonce, try.Catalogue).Scan(&current)
	if err != nil {
		return err
	}
	if !current {
		return ErrOperatorRunner
	}
	a, err := decodeOperatorMintAction(row)
	if err != nil {
		return err
	}
	if a.Store != c.lock.dbPath || !a.Options.Headless {
		return ErrOperatorConflict
	}
	b, err := proof.Bind(a, [32]byte(row.Challenge), row.ExpiresAt, row.Artifact)
	if err != nil {
		return err
	}
	if row.ApprovedBy != b.Format+":"+b.Name {
		return ErrOperatorConflict
	}
	if err := operatorMintApprover(ctx, tx, proof, b, now); err != nil {
		return err
	}
	session, err := operatorMintSession(in, a, now)
	if err != nil {
		return err
	}
	err = tx.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM session_credentials WHERE episode_id = ?)`, a.Episode).Scan(&current)
	if err != nil {
		return err
	}
	if current {
		return ErrOperatorConflict
	}
	if err := insertSessionRow(ctx, tx, session); err != nil {
		return err
	}
	start := in.Start
	start.StartedAt = time.Unix(now, 0).UTC().Format(time.RFC3339)
	startJSON, err := insertEpisodeStart(ctx, tx, a.Episode, start)
	if err != nil {
		return err
	}
	receipt := in.CredentialID + fmt.Sprintf(":%x", sha256.Sum256(startJSON))
	res, err := tx.ExecContext(ctx, `UPDATE operator_action_steps SET phase = 'observed', version = version + 1, receipt = ?
		WHERE kind = 'mint' AND identity = ? AND action_id = ? AND phase = 'reserved' AND version = ? AND owner = ?`,
		[]byte(receipt), a.ID, a.ID, try.StepVersion, row.Owner)
	if err := operatorChangedOne(res, err); err != nil {
		return err
	}
	return tx.Commit()
}

func decodeOperatorMintAction(row OperatorActionRow) (approval.OperatorAction, error) {
	if sha256.Sum256(row.Action) != [32]byte(row.Digest) {
		return approval.OperatorAction{}, ErrOperatorConflict
	}
	a, err := approval.DecodeOperatorAction(row.Action)
	if err != nil {
		return approval.OperatorAction{}, err
	}
	if a.ID != row.ID || a.Episode != row.Episode || a.Kind != row.Kind {
		return approval.OperatorAction{}, ErrOperatorConflict
	}
	return a, nil
}

func operatorMintSession(in OperatorMintInput, a approval.OperatorAction, now int64) (SessionRow, error) {
	start := in.Start
	if start.BranchExisted || len(start.Resumes) != 0 || len(start.Launches) != 0 ||
		start.Preset.Name != a.Options.Preset || start.Task == nil ||
		start.Task.SHA256 != hashref.SumSHA256([]byte(a.Task)).String() {
		return SessionRow{}, ErrOperatorConflict
	}
	if !filepath.IsAbs(in.TokenPath) || filepath.Clean(in.TokenPath) != in.TokenPath ||
		filepath.Base(in.TokenPath) != a.Episode+".session" ||
		start.Task.Path != filepath.Join(filepath.Dir(in.TokenPath), a.Episode+".task.md") {
		return SessionRow{}, ErrOperatorConflict
	}
	if _, profileErr := hashref.Parse(in.Profile.Digest); start.TTLS <= 0 || now > math.MaxInt64-start.TTLS || len(in.Grants) == 0 ||
		in.Profile.Project == "" || profileErr != nil || !operatorDigestString(in.CredentialID) {
		return SessionRow{}, ErrOperatorInput
	}
	facts, err := in.Resolution()
	if err != nil {
		return SessionRow{}, err
	}
	for _, fact := range facts {
		if value := mintFact(a, fact.Name); value != fact.Value {
			return SessionRow{}, &approval.BindingError{Field: fact.Name}
		}
	}
	expiry := now + start.TTLS
	grants := make([]approval.FactGrant, len(in.Grants))
	for i, grant := range in.Grants {
		grants[i] = approval.FactGrant{Effect: grant.Effect, Scope: grant.Scope, Budget: grant.Budget, ExpiresAt: expiry}
	}
	wire, _ := json.Marshal(grants) // approval.FactGrant contains only strings and integers.
	return SessionRow{CredentialID: in.CredentialID, EpisodeID: a.Episode, CreatedAt: now, ExpiresAt: expiry,
		GrantsJSON: string(wire), Project: in.Profile.Project, ProfileDigest: in.Profile.Digest}, nil
}

func operatorDigestString(s string) bool {
	b, err := hex.DecodeString(s)
	return err == nil && len(b) == sha256.Size && hex.EncodeToString(b) == s
}

// OperatorMintReceipt is redacted evidence for attended recovery. CAS and
// step version pin the exact action and receipt read before the terminal prompt.
type OperatorMintReceipt struct {
	CAS                              OperatorCAS
	Owner, Episode, Store, TokenPath string
	StepVersion, ExpiresAt           int64
	CredentialID, StartDigest        string
}

type operatorMintReader interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func (c *Credentials) ReadOperatorMint(ctx context.Context, id string) (OperatorMintReceipt, error) {
	if err := c.operatorAccess(ctx); err != nil {
		return OperatorMintReceipt{}, err
	}
	return readOperatorMint(ctx, c.db, id, false, false)
}

func readOperatorMint(ctx context.Context, db operatorMintReader, id string, fileStage, launchStage bool) (OperatorMintReceipt, error) {
	row, err := scanOperatorAction(db.QueryRowContext(ctx, `SELECT `+operatorColumns+` FROM operator_actions WHERE id = ?`, id))
	if err != nil {
		return OperatorMintReceipt{}, err
	}
	if row.Kind != "start" || row.State != "running" && row.State != "uncertain" {
		return OperatorMintReceipt{}, ErrOperatorConflict
	}
	a, err := decodeOperatorMintAction(row)
	if err != nil {
		return OperatorMintReceipt{}, err
	}
	r := OperatorMintReceipt{CAS: OperatorCAS{ID: id, Digest: [32]byte(row.Digest), Generation: row.Generation},
		Owner: row.Owner, Episode: a.Episode, Store: a.Store, TokenPath: mintFact(a, "mint_token_path")}
	var receipt []byte
	err = db.QueryRowContext(ctx, `SELECT version, receipt FROM operator_action_steps WHERE kind = 'mint'
		AND identity = ? AND action_id = ? AND phase = 'observed' AND owner = ?
		AND NOT EXISTS (SELECT 1 FROM operator_action_steps WHERE action_id = ? AND kind NOT IN ('operation','worktree','mint') AND NOT (? AND kind = 'files')
		AND NOT (? AND kind = 'launch' AND phase='reserved'))`,
		id, id, row.Owner, id, fileStage, launchStage).Scan(&r.StepVersion, &receipt)
	if err != nil {
		return OperatorMintReceipt{}, fmt.Errorf("mint evidence unavailable or later effect reserved: %w", err)
	}
	cred, startDigest, _ := bytes.Cut(receipt, []byte(":"))
	r.CredentialID, r.StartDigest = string(cred), string(startDigest)
	if !operatorDigestString(r.CredentialID) || !operatorDigestString(r.StartDigest) {
		return OperatorMintReceipt{}, ErrOperatorConflict
	}
	var start []byte
	err = db.QueryRowContext(ctx, `SELECT s.start_json, c.expires_at FROM episode_starts s JOIN session_credentials c
		ON c.episode_id = s.episode WHERE s.episode = ? AND c.credential_id = ?`, a.Episode, r.CredentialID).Scan(&start, &r.ExpiresAt)
	if err != nil {
		return OperatorMintReceipt{}, fmt.Errorf("mint start or credential unavailable: %w", err)
	}
	if fmt.Sprintf("%x", sha256.Sum256(start)) != r.StartDigest {
		return OperatorMintReceipt{}, errors.New("mint start record differs from receipt")
	}
	return r, nil
}

// AbortOperatorMint revokes only the receipt's credential and releases the
// exact reviewed owner after the attended caller proves no token was installed.
func (c *Credentials) AbortOperatorMint(ctx context.Context, reviewed OperatorMintReceipt, now int64) error {
	if err := c.operatorAccess(ctx); err != nil {
		return err
	}
	if c.lock == nil || reviewed.Store != c.lock.dbPath {
		return ErrOperatorConflict
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	current, err := readOperatorMint(ctx, tx, reviewed.CAS.ID, false, false)
	if err != nil {
		return err
	}
	if current != reviewed {
		return ErrOperatorConflict
	}
	res, err := tx.ExecContext(ctx, `DELETE FROM session_credentials WHERE credential_id = ? AND episode_id = ?`,
		current.CredentialID, current.Episode)
	if err := operatorChangedOne(res, err); err != nil {
		return err
	}
	res, err = tx.ExecContext(ctx, `UPDATE operator_actions SET state = 'failed', owner = NULL, generation = generation + 1, terminal_at = ?
		WHERE id = ? AND digest = ? AND generation = ? AND owner = ? AND state IN ('running','uncertain')
		AND EXISTS (SELECT 1 FROM operator_action_steps WHERE kind = 'mint' AND identity = ?
			AND action_id = ? AND phase = 'observed' AND version = ? AND owner = ?)`,
		now, current.CAS.ID, current.CAS.Digest[:], current.CAS.Generation, current.Owner,
		current.CAS.ID, current.CAS.ID, current.StepVersion, current.Owner)
	if err := operatorChangedOne(res, err); err != nil {
		return err
	}
	return tx.Commit()
}
