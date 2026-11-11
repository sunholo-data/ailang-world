package store

// Row 201 A2 (design_docs/planned/w-world-approvals-by-key-a2.md §2.10): passkeys (credentials
// schema 6). A webauthn approver is an approvers row (format webauthn, public_key the canonical
// COSE key) and one webauthn_credentials row, written together by EnrolPasskey after the typed
// enrolment code and y. Like approvers, a passkey row is never deleted; its sign_count only moves
// forward, inside the consume that uses it. The store does not import host/approval: these checks
// are the store's own backstop for the pure law the CLI runs first.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ---- M4a ----

const passkeysTable = "webauthn_credentials"

// passkeysSchemaVersion is the first credentials schema with webauthn_credentials.
const passkeysSchemaVersion = 6

// passkeysSQL is the version 5 -> 6 upgrade: credentials_schema.sql's webauthn_credentials table,
// column for column (TestCredentialsSchemaMatchesCanonicalDDL compares the two).
const passkeysSQL = `CREATE TABLE webauthn_credentials (
    name            TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
    credential_id   BLOB NOT NULL UNIQUE CHECK (length(credential_id) BETWEEN 1 AND 1023),
    rp_id           TEXT NOT NULL CHECK (rp_id <> ''),
    origin          TEXT NOT NULL CHECK (origin <> ''),
    alg             INTEGER NOT NULL CHECK (alg IN (-7, -8)),
    aaguid          BLOB NOT NULL CHECK (length(aaguid) = 16),
    backup_eligible INTEGER NOT NULL CHECK (backup_eligible IN (0, 1)),
    sign_count      INTEGER NOT NULL CHECK (sign_count BETWEEN 0 AND 4294967295)
)`

// ---- M4b ----

// maxSignCount is the largest WebAuthn sign counter (32 bits); maxActivePasskeys caps ActivePasskeys.
const (
	maxSignCount      = 1<<32 - 1
	maxActivePasskeys = 64
)

// PasskeyRow is one webauthn_credentials row.
type PasskeyRow struct {
	Name, RPID, Origin   string
	CredentialID, AAGUID []byte
	Alg                  int64 // -7 (ES256) or -8 (EdDSA)
	BackupEligible       bool
	SignCount            int64 // 0 .. 2^32-1
}

var (
	ErrPasskeyCredentialEnrolled = errors.New("store: that passkey credential is already enrolled (a credential is never enrolled twice, even after revoke)")
	ErrSignCountRegression       = errors.New("store: the passkey's sign counter did not advance (a cloned authenticator, or an older assertion)")
	ErrPasskeyChanged            = errors.New("store: the approver's passkey is not the one that signed")
)

func (c *Credentials) olderForPasskeys() error {
	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: passkeysSchemaVersion, Feature: "passkeys"}
}

// checkPasskeyRows is EnrolPasskey's validation (S1-S10): each refusal is ErrBadApprover.
func checkPasskeyRows(a ApproverRow, p PasskeyRow) error {
	switch {
	case a.Format != "webauthn": // S1
		return fmt.Errorf("%w: format %q is not webauthn", ErrBadApprover, a.Format)
	case a.Name != p.Name: // S2
		return fmt.Errorf("%w: the approver %q and its passkey %q differ", ErrBadApprover, a.Name, p.Name)
	case a.Name == "" || len(a.PublicKey) == 0: // S3
		return fmt.Errorf("%w: empty name or public key", ErrBadApprover)
	case a.EnrolledAt <= 0: // S4
		return fmt.Errorf("%w: enrolled_at %d is not positive", ErrBadApprover, a.EnrolledAt)
	case a.LastUsedAt != 0 || a.RevokedAt != 0: // S5
		return fmt.Errorf("%w: a new approver is neither used nor revoked", ErrBadApprover)
	case len(p.CredentialID) < 1 || len(p.CredentialID) > 1023: // S6
		return fmt.Errorf("%w: a credential ID is 1 to 1023 bytes, not %d", ErrBadApprover, len(p.CredentialID))
	case p.Alg != -7 && p.Alg != -8: // S7
		return fmt.Errorf("%w: alg %d is not -7 or -8", ErrBadApprover, p.Alg)
	case len(p.AAGUID) != 16: // S8
		return fmt.Errorf("%w: an AAGUID is 16 bytes, not %d", ErrBadApprover, len(p.AAGUID))
	case p.RPID == "" || p.Origin == "": // S9
		return fmt.Errorf("%w: empty relying party ID or origin", ErrBadApprover)
	case p.SignCount < 0 || p.SignCount > maxSignCount: // S10
		return fmt.Errorf("%w: sign count %d is not 0 to 2^32-1", ErrBadApprover, p.SignCount)
	}
	return nil
}

// EnrolPasskey stores one new, active webauthn approver and its passkey in one transaction. It
// refuses a name, a key or a credential ID that any row holds, revoked rows included, naming the
// holder; the UNIQUE constraints are the backstop and their error text is never parsed.
func (c *Credentials) EnrolPasskey(ctx context.Context, a ApproverRow, p PasskeyRow) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.version < passkeysSchemaVersion {
		return c.olderForPasskeys()
	}
	if err := checkPasskeyRows(a, p); err != nil {
		return err
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: enrol passkey: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	for _, pre := range []struct {
		query string
		arg   any
		err   error
	}{
		{`SELECT name FROM approvers WHERE name = ?`, a.Name, ErrApproverNameTaken},
		{`SELECT name FROM approvers WHERE public_key = ?`, a.PublicKey, ErrApproverKeyEnrolled},
		{`SELECT name FROM webauthn_credentials WHERE credential_id = ?`, p.CredentialID, ErrPasskeyCredentialEnrolled},
	} {
		var holder string
		err := tx.QueryRowContext(ctx, pre.query, pre.arg).Scan(&holder)
		if err == nil {
			return fmt.Errorf("%w: approver %q holds it", pre.err, holder)
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("store: enrol passkey: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES (?, ?, ?, ?)`,
		a.Name, a.Format, a.PublicKey, a.EnrolledAt); err != nil {
		return fmt.Errorf("store: enrol passkey: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO webauthn_credentials (name, credential_id, rp_id, origin, alg, aaguid, backup_eligible, sign_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		p.Name, p.CredentialID, p.RPID, p.Origin, p.Alg, p.AAGUID, p.BackupEligible, p.SignCount); err != nil {
		return fmt.Errorf("store: enrol passkey: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: enrol passkey: commit: %w", err)
	}
	return nil
}
