package store

// Row 201 A2 (design_docs/planned/w-world-approvals-by-key-a2.md §2.10): passkeys (credentials
// schema 6). A webauthn approver is an approvers row (format webauthn, public_key the canonical
// COSE key) and one webauthn_credentials row, written together by EnrolPasskey after the typed
// enrolment code and y. Like approvers, a passkey row is never deleted; its sign_count only moves
// forward, inside the consume that uses it. The store does not import host/approval: these checks
// are the store's own backstop for the pure law the CLI runs first.

import (
	"bytes"
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

// ---- M4c ----

// EnrolledPasskey is an approvers row with its passkey row.
type EnrolledPasskey struct {
	Approver ApproverRow
	Passkey  PasskeyRow
}

const passkeyColumns = `w.name, w.credential_id, w.rp_id, w.origin, w.alg, w.aaguid, w.backup_eligible, w.sign_count`

func scanPasskey(scan func(...any) error, dst ...any) (PasskeyRow, error) {
	var p PasskeyRow
	err := scan(append(dst, &p.Name, &p.CredentialID, &p.RPID, &p.Origin, &p.Alg, &p.AAGUID, &p.BackupEligible, &p.SignCount)...)
	return p, err
}

// enrolledPasskeys reads approvers joined with their passkeys, where says which.
func (c *Credentials) enrolledPasskeys(ctx context.Context, where string, args ...any) ([]EnrolledPasskey, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if c.version < passkeysSchemaVersion {
		return nil, c.olderForPasskeys()
	}
	rows, err := c.db.QueryContext(ctx, `SELECT a.name, a.format, a.public_key, a.enrolled_at, a.last_used_at, a.revoked_at, `+passkeyColumns+
		` FROM approvers a JOIN webauthn_credentials w ON w.name = a.name WHERE `+where, args...)
	if err != nil {
		return nil, fmt.Errorf("store: read passkeys: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []EnrolledPasskey
	for rows.Next() {
		var e EnrolledPasskey
		var lastUsed, revoked sql.NullInt64
		a := &e.Approver
		if e.Passkey, err = scanPasskey(rows.Scan, &a.Name, &a.Format, &a.PublicKey, &a.EnrolledAt, &lastUsed, &revoked); err != nil {
			return nil, fmt.Errorf("store: read passkeys: %w", err)
		}
		a.LastUsedAt, a.RevokedAt = lastUsed.Int64, revoked.Int64
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: read passkeys: %w", err)
	}
	return out, nil
}

// PasskeyByCredentialID reads the approver whose passkey has credential ID id, byte for byte, and
// that passkey. ok is false when none.
func (c *Credentials) PasskeyByCredentialID(ctx context.Context, id []byte) (ApproverRow, PasskeyRow, bool, error) {
	got, err := c.enrolledPasskeys(ctx, `w.credential_id = ?`, id)
	if err != nil || len(got) == 0 {
		return ApproverRow{}, PasskeyRow{}, false, err
	}
	return got[0].Approver, got[0].Passkey, true, nil
}

// Passkey reads approver name's passkey, revoked approvers included. ok is false when none.
func (c *Credentials) Passkey(ctx context.Context, name string) (PasskeyRow, bool, error) {
	got, err := c.enrolledPasskeys(ctx, `w.name = ?`, name)
	if err != nil || len(got) == 0 {
		return PasskeyRow{}, false, err
	}
	return got[0].Passkey, true, nil
}

// ActivePasskeys lists the active webauthn approvers with their passkeys, ordered by name, at most 64.
func (c *Credentials) ActivePasskeys(ctx context.Context) ([]EnrolledPasskey, error) {
	return c.enrolledPasskeys(ctx, `a.revoked_at IS NULL AND a.format = 'webauthn' ORDER BY a.name LIMIT ?`, maxActivePasskeys)
}

// ---- M5a ----

// PasskeyUse is what a passkey consume binds: the credential that signed and its new counter.
type PasskeyUse struct {
	CredentialID []byte
	SignCount    int64 // the assertion's counter, 0 .. 2^32-1
}

// usePasskey is the consume's counter binding (design §2.10): inside the consume's transaction,
// after useApprover, it moves name's sign_count to p.SignCount only for this credential and only
// forward (both zero: unchanged). Otherwise it names why (C16-C18), and the caller's rollback undoes
// the request update and last_used_at.
func usePasskey(ctx context.Context, tx *sql.Tx, name string, p PasskeyUse) error {
	res, err := tx.ExecContext(ctx, `UPDATE webauthn_credentials SET sign_count = ? WHERE name = ? AND credential_id = ? AND (sign_count < ? OR (sign_count = 0 AND ? = 0))`,
		p.SignCount, name, p.CredentialID, p.SignCount, p.SignCount)
	if err != nil {
		return fmt.Errorf("store: use passkey: %w", err)
	}
	if n, err := res.RowsAffected(); err != nil || n == 1 {
		return err
	}
	var id []byte
	var count int64
	err = tx.QueryRowContext(ctx, `SELECT credential_id, sign_count FROM webauthn_credentials WHERE name = ?`, name).Scan(&id, &count)
	switch {
	case errors.Is(err, sql.ErrNoRows): // C16
		return fmt.Errorf("%w: %q has no passkey", ErrApproverUnknown, name)
	case err != nil:
		return fmt.Errorf("store: use passkey: %w", err)
	case !bytes.Equal(id, p.CredentialID): // C17
		return fmt.Errorf("%w: %q", ErrPasskeyChanged, name)
	default: // C18
		return fmt.Errorf("%w: %q holds %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
			ErrSignCountRegression, name, count, p.SignCount, name)
	}
}
