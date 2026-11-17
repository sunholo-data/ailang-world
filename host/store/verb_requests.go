package store

// Row 243 V3 sprint 3 M7 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.4): credentials
// schema 9's verb_requests, one per `episode finish --confirm-on-surface`, inert until a verified key
// confirms it; the authority is the world store's finish approval (decision_ref), never this row.

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
)

const verbRequestsTable = "verb_requests"

// verbRequestsSchemaVersion is the first credentials schema with verb_requests.
const verbRequestsSchemaVersion = 9

// verbRequestsSQL is the version 8 -> 9 upgrade: credentials_schema.sql's verb_requests table,
// column for column (TestCredentialsSchema9VerbRequests compares a fresh store with an upgraded one).
const verbRequestsSQL = `CREATE TABLE verb_requests (
    request_id   TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'vr-[0-9a-f]*'),
    request      BLOB NOT NULL CHECK (length(request) BETWEEN 1 AND 131072),
    digest       TEXT NOT NULL CHECK (length(digest) = 64),
    challenge    TEXT NOT NULL CHECK (length(challenge) = 64),
    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
    created_at   INTEGER NOT NULL,
    expires_at   INTEGER NOT NULL CHECK (expires_at > created_at),
    consumed_at  INTEGER,
    outcome      TEXT CHECK (outcome IS NULL OR outcome IN ('approved', 'superseded')),
    approved_by  TEXT CHECK (approved_by IS NULL OR approved_by <> ''),
    artifact     BLOB CHECK (artifact IS NULL OR length(artifact) BETWEEN 1 AND 16384),
    decision_ref TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*')),
    CHECK ((consumed_at IS NULL) = (outcome IS NULL)),
    CHECK ((outcome IS 'approved') = (approved_by IS NOT NULL AND artifact IS NOT NULL AND decision_ref IS NOT NULL))
)`

// MaxPendingVerbRequests caps the unconsumed, unexpired verb requests one store holds (design T14).
const MaxPendingVerbRequests = 16

// MaxVerbRequestBytes is the table's request size bound (approval.MaxVerbRequestBytes, 128 KiB).
const MaxVerbRequestBytes = 128 << 10

// verbRequestPruneAfter: an unconsumed verb request is deleted once it has been expired this long.
const verbRequestPruneAfter = 86400

// The outcomes a consumed verb request carries.
const (
	VerbOutcomeApproved   = "approved"
	VerbOutcomeSuperseded = "superseded"
)

// VerbRequestIDPattern is a verb request id: "vr-" and 16 lowercase hex. The table's CHECK is
// weaker, so PutVerbRequest enforces this.
var VerbRequestIDPattern = regexp.MustCompile(`^vr-[0-9a-f]{16}$`)

var (
	ErrBadVerbRequest             = errors.New("store: invalid verb request row")
	ErrTooManyPendingVerbRequests = errors.New("store: 16 verb requests are already pending; confirm them or let them expire")
	ErrVerbRequestUnknown         = errors.New("store: no such verb request")
	ErrVerbRequestConsumed        = errors.New("store: the verb request was already confirmed")
	ErrVerbRequestSuperseded      = errors.New("store: the verb request was superseded by a newer finish of its episode")
	ErrVerbRequestExpired         = errors.New("store: the verb request expired before it was confirmed")
	ErrVerbRequestClockBehind     = errors.New("store: the clock is behind the verb request's creation time")
	ErrVerbRequestDigestChanged   = errors.New("store: the verb request's digest is not the one confirmed")
)

// VerbRequestRow is one verb_requests row.
type VerbRequestRow struct {
	RequestID   string
	Request     []byte // the canonical world/verb-request/v1 bytes (approval.CanonicalVerb)
	Digest      string // 64 lowercase hex: SHA-256 of Request
	Challenge   string // 64 lowercase hex, random, single use
	EpisodeID   string
	CreatedAt   int64
	ExpiresAt   int64
	ConsumedAt  int64  // 0 while pending
	Outcome     string // "" while pending, else VerbOutcomeApproved or VerbOutcomeSuperseded
	ApprovedBy  string // "<format>:<name>" of an approved row
	Artifact    []byte // the verified artifact of an approved row
	DecisionRef string // the world store's finish approval decision of an approved row
}

func (c *Credentials) olderForVerbRequests() error {
	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: verbRequestsSchemaVersion, Feature: "verb requests"}
}

// sha256Hex is the lowercase hex SHA-256 of b.
func sha256Hex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// checkVerbRequestRow is PutVerbRequest's grammar (S2): each refusal is ErrBadVerbRequest.
func checkVerbRequestRow(row VerbRequestRow) error {
	switch {
	case !VerbRequestIDPattern.MatchString(row.RequestID):
		return fmt.Errorf("%w: request id %q", ErrBadVerbRequest, row.RequestID)
	case len(row.Request) < 1 || len(row.Request) > MaxVerbRequestBytes:
		return fmt.Errorf("%w: a request is 1 to %d bytes, not %d", ErrBadVerbRequest, MaxVerbRequestBytes, len(row.Request))
	case row.Digest != sha256Hex(row.Request): // so the digest is 64 lowercase hex
		return fmt.Errorf("%w: the digest is not the request's SHA-256", ErrBadVerbRequest)
	case !isLowerHex64(row.Challenge):
		return fmt.Errorf("%w: the challenge is not 64 lowercase hex", ErrBadVerbRequest)
	case row.EpisodeID == "":
		return fmt.Errorf("%w: empty episode id", ErrBadVerbRequest)
	case row.ExpiresAt <= row.CreatedAt:
		return fmt.Errorf("%w: expires_at %d is not after created_at %d", ErrBadVerbRequest, row.ExpiresAt, row.CreatedAt)
	case row.ConsumedAt != 0 || row.Outcome != "" || row.ApprovedBy != "" || row.Artifact != nil || row.DecisionRef != "":
		return fmt.Errorf("%w: a new request is pending", ErrBadVerbRequest)
	}
	return nil
}

// PutVerbRequest stores one new, pending request in one transaction (design §2.4): S1 deadline and
// schema 9; S2 grammar; S3 the episode's pending request is superseded; S4 requests expired a day
// ago are pruned; S5 at MaxPendingVerbRequests live ones it refuses; S6 insert.
func (c *Credentials) PutVerbRequest(ctx context.Context, row VerbRequestRow, now int64) error {
	if err := requireDeadline(ctx); err != nil { // S1
		return err
	}
	if c.version < verbRequestsSchemaVersion { // S1
		return c.olderForVerbRequests()
	}
	if err := checkVerbRequestRow(row); err != nil { // S2
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: put verb request: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	if _, err := tx.ExecContext(ctx, `UPDATE verb_requests SET consumed_at = ?, outcome = 'superseded' WHERE episode_id = ? AND consumed_at IS NULL`, now, row.EpisodeID); err != nil { // S3
		return fmt.Errorf("store: put verb request: supersede: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM verb_requests WHERE consumed_at IS NULL AND expires_at <= ?`, now-verbRequestPruneAfter); err != nil { // S4
		return fmt.Errorf("store: put verb request: prune: %w", err)
	}
	var live int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM verb_requests WHERE consumed_at IS NULL AND expires_at > ?`, now).Scan(&live); err != nil {
		return fmt.Errorf("store: put verb request: count: %w", err)
	}
	if live >= MaxPendingVerbRequests { // S5
		return ErrTooManyPendingVerbRequests
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO verb_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		row.RequestID, row.Request, row.Digest, row.Challenge, row.EpisodeID, row.CreatedAt, row.ExpiresAt); err != nil { // S6
		return fmt.Errorf("store: put verb request: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: put verb request: commit: %w", err)
	}
	return nil
}

const verbRequestColumns = `request_id, request, digest, challenge, episode_id, created_at, expires_at, coalesce(consumed_at, 0), coalesce(outcome, ''), coalesce(approved_by, ''), artifact, coalesce(decision_ref, '')`

func scanVerbRequest(scan func(...any) error) (VerbRequestRow, error) {
	var r VerbRequestRow
	err := scan(&r.RequestID, &r.Request, &r.Digest, &r.Challenge, &r.EpisodeID, &r.CreatedAt, &r.ExpiresAt, &r.ConsumedAt, &r.Outcome, &r.ApprovedBy, &r.Artifact, &r.DecisionRef)
	return r, err
}

// VerbRequest reads one request by id; ok is false when there is none.
func (c *Credentials) VerbRequest(ctx context.Context, id string) (VerbRequestRow, bool, error) {
	if err := requireDeadline(ctx); err != nil {
		return VerbRequestRow{}, false, err
	}
	if c.version < verbRequestsSchemaVersion {
		return VerbRequestRow{}, false, c.olderForVerbRequests()
	}
	row, err := scanVerbRequest(c.db.QueryRowContext(ctx, `SELECT `+verbRequestColumns+` FROM verb_requests WHERE request_id = ?`, id).Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return VerbRequestRow{}, false, nil
	}
	if err != nil {
		return VerbRequestRow{}, false, fmt.Errorf("store: read verb request: %w", err)
	}
	return row, true, nil
}

// PendingVerbRequests lists the live, pending requests, newest first. limit <= 0 or above
// MaxPendingVerbRequests means MaxPendingVerbRequests.
func (c *Credentials) PendingVerbRequests(ctx context.Context, now int64, limit int) ([]VerbRequestRow, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if c.version < verbRequestsSchemaVersion {
		return nil, c.olderForVerbRequests()
	}
	if limit <= 0 || limit > MaxPendingVerbRequests {
		limit = MaxPendingVerbRequests
	}
	rows, err := c.db.QueryContext(ctx, `SELECT `+verbRequestColumns+` FROM verb_requests WHERE consumed_at IS NULL AND expires_at > ? ORDER BY created_at DESC, request_id DESC LIMIT ?`, now, limit)
	if err != nil {
		return nil, fmt.Errorf("store: list verb requests: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []VerbRequestRow
	for rows.Next() {
		row, err := scanVerbRequest(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("store: list verb requests: %w", err)
		}
		out = append(out, row)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: list verb requests: %w", err)
	}
	return out, nil
}

// checkVerbConsume is ConsumeVerbRequest's grammar: each refusal is ErrBadVerbRequest.
func checkVerbConsume(id, digest string, use ApproverUse, approvedBy string, artifact []byte, decisionRef string) error {
	switch {
	case !VerbRequestIDPattern.MatchString(id) || !isLowerHex64(digest):
		return fmt.Errorf("%w: request %q, digest %q", ErrBadVerbRequest, id, digest)
	case approvedBy != use.Format+":"+use.Name:
		return fmt.Errorf("%w: %q is not the approver %s:%s", ErrBadVerbRequest, approvedBy, use.Format, use.Name)
	case (use.Format == "webauthn") != (use.Passkey != nil): // the counter binding is exactly a passkey's (A2 C13)
		return fmt.Errorf("%w: a %s confirmation with passkey use %v", ErrBadVerbRequest, use.Format, use.Passkey != nil)
	case len(artifact) < 1 || len(artifact) > MaxSurfaceApprovalArtifactBytes:
		return fmt.Errorf("%w: an artifact is 1 to %d bytes, not %d", ErrBadVerbRequest, MaxSurfaceApprovalArtifactBytes, len(artifact))
	case !isDecisionRef(decisionRef):
		return fmt.Errorf("%w: decision %q", ErrBadVerbRequest, decisionRef)
	}
	return nil
}

// ConsumeVerbRequest confirms request id at now in one transaction, after the world store's finish
// approval decisionRef landed (design §2.4, §2.6 C15): the row is re-read and refused in Admit's
// order (consumed or superseded, clock behind, expired, digest); the approver is re-checked and
// marked used (useApprover) and, for a passkey, its counter moved forward (usePasskey); then one
// UPDATE sets consumed_at, outcome, approved_by, artifact and decision_ref. Any refusal writes
// nothing.
func (c *Credentials) ConsumeVerbRequest(ctx context.Context, id, digest string, use ApproverUse, approvedBy string, artifact []byte, decisionRef string, now int64) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.version < verbRequestsSchemaVersion {
		return c.olderForVerbRequests()
	}
	if err := checkVerbConsume(id, digest, use, approvedBy, artifact, decisionRef); err != nil {
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: consume verb request: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	row, err := scanVerbRequest(tx.QueryRowContext(ctx, `SELECT `+verbRequestColumns+` FROM verb_requests WHERE request_id = ?`, id).Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: %s", ErrVerbRequestUnknown, id)
	}
	if err != nil {
		return fmt.Errorf("store: consume verb request: %w", err)
	}
	switch {
	case row.Outcome == VerbOutcomeSuperseded:
		return fmt.Errorf("%w: %s", ErrVerbRequestSuperseded, id)
	case row.Outcome != "":
		return fmt.Errorf("%w: %s by %s", ErrVerbRequestConsumed, id, row.ApprovedBy)
	case now < row.CreatedAt:
		return fmt.Errorf("%w: %s", ErrVerbRequestClockBehind, id)
	case now >= row.ExpiresAt:
		return fmt.Errorf("%w: %s", ErrVerbRequestExpired, id)
	case digest != row.Digest || sha256Hex(row.Request) != row.Digest:
		return fmt.Errorf("%w: %s", ErrVerbRequestDigestChanged, id)
	}
	if err := useApprover(ctx, tx, use, now); err != nil {
		return err
	}
	if use.Passkey != nil {
		if err := usePasskey(ctx, tx, use.Name, *use.Passkey); err != nil {
			return err
		}
	}
	res, err := tx.ExecContext(ctx, `UPDATE verb_requests SET consumed_at = ?, outcome = 'approved', approved_by = ?, artifact = ?, decision_ref = ? WHERE request_id = ? AND consumed_at IS NULL`,
		now, approvedBy, artifact, decisionRef, id)
	if err != nil {
		return fmt.Errorf("store: consume verb request: %w", err)
	}
	if n, err := res.RowsAffected(); err != nil || n != 1 {
		return fmt.Errorf("%w: %s (%d rows, %v)", ErrVerbRequestConsumed, id, n, err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: consume verb request: commit: %w", err)
	}
	return nil
}
