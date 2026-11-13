package store

// Row 201 A4 sprint 2 (design §2.7): credentials schema 8's key_decisions, the evidence of a key
// decision. The authority is the world store's decision, never this row; decision_ref links the two
// once it landed (NULL until then). Format-agnostic: sprint 3's passkey arm adds no column (S2-D3);
// its counter is webauthn_credentials.sign_count, moved by usePasskey in RecordKeyDecision (K4).

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

const keyDecisionsTable = "key_decisions"

// keyDecisionsSchemaVersion is the first credentials schema with key_decisions.
const keyDecisionsSchemaVersion = 8

// keyDecisionsSQL is the version 7 -> 8 upgrade: credentials_schema.sql's key_decisions table,
// column for column (TestCredentialsSchema8KeyDecisions compares the two), the design's DDL verbatim.
const keyDecisionsSQL = `CREATE TABLE key_decisions (
    request_ref   TEXT PRIMARY KEY CHECK (length(request_ref) = 71 AND request_ref GLOB 'sha256:*'),
    packet_ref    TEXT NOT NULL CHECK (length(packet_ref) = 71 AND packet_ref GLOB 'sha256:*'),
    verb          TEXT NOT NULL CHECK (verb IN ('approve', 'deny')),
    format        TEXT NOT NULL CHECK (format IN ('ssh', 'webauthn')),
    approver      TEXT NOT NULL CHECK (approver <> ''),
    artifact      BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
    expiry        INTEGER NOT NULL,
    decided_at    INTEGER NOT NULL,
    decision_ref  TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*'))
)`

var (
	// ErrBadKeyDecision refuses a row or a ref that breaks the table's grammar (K2).
	ErrBadKeyDecision = errors.New("store: invalid key decision row")
	// ErrKeyDecisionLanded refuses recording a request whose key decision already landed (K0).
	ErrKeyDecisionLanded = errors.New("store: the request's key decision already landed")
	// ErrKeyDecisionNotPending refuses landing a request with no unlanded key decision.
	ErrKeyDecisionNotPending = errors.New("store: no unlanded key decision for the request")
)

// KeyDecisionRow is one key_decisions row.
type KeyDecisionRow struct {
	RequestRef  string // sha256: and 64 lowercase hex, the request decided
	PacketRef   string // the one packet carrying it
	Verb        string // "approve" or "deny"
	Format      string // "ssh" (sprint 2) or "webauthn" (sprint 3)
	Approver    string // "<format>:<name>", the decision's decidedBy
	Artifact    []byte // the armored SSHSIG (or the assertion), 1..16384 bytes
	Expiry      int64  // the signed expiry
	DecidedAt   int64  // set by RecordKeyDecision to its now
	DecisionRef string // "" until LandKeyDecision
}

const keyDecisionColumns = `request_ref, packet_ref, verb, format, approver, artifact, expiry, decided_at, coalesce(decision_ref, '')`

func (c *Credentials) olderForKeyDecisions() error {
	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: keyDecisionsSchemaVersion, Feature: "key decisions"}
}

// isDecisionRef is the table's ref grammar, held in code to lower-case hex (the CHECKs hold length and prefix).
func isDecisionRef(s string) bool { return len(s) == 71 && s[:7] == "sha256:" && isLowerHex64(s[7:]) }

// checkKeyDecisionRow is RecordKeyDecision's grammar (K2): each refusal is ErrBadKeyDecision.
func checkKeyDecisionRow(row KeyDecisionRow, use ApproverUse) error {
	switch {
	case !isDecisionRef(row.RequestRef) || !isDecisionRef(row.PacketRef):
		return fmt.Errorf("%w: request %q, packet %q", ErrBadKeyDecision, row.RequestRef, row.PacketRef)
	case row.Verb != "approve" && row.Verb != "deny":
		return fmt.Errorf("%w: verb %q", ErrBadKeyDecision, row.Verb)
	case row.Format != use.Format || row.Approver != use.Format+":"+use.Name:
		return fmt.Errorf("%w: %s by %q is not the approver %s:%s", ErrBadKeyDecision, row.Format, row.Approver, use.Format, use.Name)
	case (use.Format == "webauthn") != (use.Passkey != nil): // the counter binding is exactly a passkey's (A2 C13)
		return fmt.Errorf("%w: a %s decision with passkey use %v", ErrBadKeyDecision, use.Format, use.Passkey != nil)
	case len(row.Artifact) < 1 || len(row.Artifact) > MaxSurfaceApprovalArtifactBytes:
		return fmt.Errorf("%w: an artifact is 1 to %d bytes, not %d", ErrBadKeyDecision, MaxSurfaceApprovalArtifactBytes, len(row.Artifact))
	}
	return nil
}

func scanKeyDecision(scan func(...any) error) (KeyDecisionRow, error) {
	var r KeyDecisionRow
	err := scan(&r.RequestRef, &r.PacketRef, &r.Verb, &r.Format, &r.Approver, &r.Artifact, &r.Expiry, &r.DecidedAt, &r.DecisionRef)
	return r, err
}

// RecordKeyDecision records row for the approver use at now, in one transaction (design §2.7): K1
// deadline, schema 8; K2 grammar; K0 a landed row refuses, an unlanded one equal in artifact and
// signed fields is kept (the exact-match retry); K3 useApprover on both paths; K4, for a passkey and
// only off the exact match, usePasskey: the assertion law's W15 against the stored counter, in this
// transaction (design §2.7 [r2]; the decision judge defers it), so a retry of the exact artifact of
// an unlanded row is not refused by its own first use; K5 insert or replace the unlanded row.
func (c *Credentials) RecordKeyDecision(ctx context.Context, row KeyDecisionRow, use ApproverUse, now int64) error {
	if err := requireDeadline(ctx); err != nil { // K1
		return err
	}
	if c.version < keyDecisionsSchemaVersion { // K1
		return c.olderForKeyDecisions()
	}
	if err := checkKeyDecisionRow(row, use); err != nil { // K2
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: record key decision: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	old, err := scanKeyDecision(tx.QueryRowContext(ctx, `SELECT `+keyDecisionColumns+` FROM key_decisions WHERE request_ref = ?`, row.RequestRef).Scan)
	found := err == nil
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("store: record key decision: %w", err)
	}
	if found && old.DecisionRef != "" { // K0: landed
		return fmt.Errorf("%w: %s by %s", ErrKeyDecisionLanded, row.RequestRef, old.Approver)
	}
	exact := found && bytes.Equal(old.Artifact, row.Artifact) && old.PacketRef == row.PacketRef && old.Verb == row.Verb &&
		old.Format == row.Format && old.Approver == row.Approver && old.Expiry == row.Expiry // K0: the exact match
	if err := useApprover(ctx, tx, use, now); err != nil { // K3, on both paths
		return err
	}
	if !exact && use.Passkey != nil { // K4: W15, forward only (both zero: unchanged), off the exact match
		if err := usePasskey(ctx, tx, use.Name, *use.Passkey); err != nil {
			return err
		}
	}
	if !exact { // K5
		if _, err := tx.ExecContext(ctx, `INSERT INTO key_decisions (request_ref, packet_ref, verb, format, approver, artifact, expiry, decided_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(request_ref) DO UPDATE SET packet_ref = excluded.packet_ref, verb = excluded.verb, format = excluded.format, approver = excluded.approver,
artifact = excluded.artifact, expiry = excluded.expiry, decided_at = excluded.decided_at WHERE key_decisions.decision_ref IS NULL`,
			row.RequestRef, row.PacketRef, row.Verb, row.Format, row.Approver, row.Artifact, row.Expiry, now); err != nil {
			return fmt.Errorf("store: record key decision: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: record key decision: commit: %w", err)
	}
	return nil
}

// LandKeyDecision links request's unlanded key decision to the world decision decisionRef.
func (c *Credentials) LandKeyDecision(ctx context.Context, request, decisionRef string) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.version < keyDecisionsSchemaVersion {
		return c.olderForKeyDecisions()
	}
	if !isDecisionRef(decisionRef) {
		return fmt.Errorf("%w: decision %q", ErrBadKeyDecision, decisionRef)
	}
	res, err := c.db.ExecContext(ctx, `UPDATE key_decisions SET decision_ref = ? WHERE request_ref = ? AND decision_ref IS NULL`, decisionRef, request)
	if err != nil {
		return fmt.Errorf("store: land key decision: %w", err)
	}
	if n, err := res.RowsAffected(); err != nil || n != 1 {
		return fmt.Errorf("%w: %s (%v)", ErrKeyDecisionNotPending, request, err)
	}
	return nil
}

// KeyDecision reads request's key decision; ok is false when there is none.
func (c *Credentials) KeyDecision(ctx context.Context, request string) (KeyDecisionRow, bool, error) {
	if err := requireDeadline(ctx); err != nil {
		return KeyDecisionRow{}, false, err
	}
	if c.version < keyDecisionsSchemaVersion {
		return KeyDecisionRow{}, false, c.olderForKeyDecisions()
	}
	row, err := scanKeyDecision(c.db.QueryRowContext(ctx, `SELECT `+keyDecisionColumns+` FROM key_decisions WHERE request_ref = ?`, request).Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return KeyDecisionRow{}, false, nil
	}
	if err != nil {
		return KeyDecisionRow{}, false, fmt.Errorf("store: read key decision: %w", err)
	}
	return row, true, nil
}
