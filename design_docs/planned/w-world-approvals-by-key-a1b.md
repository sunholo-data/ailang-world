# w-world-approvals-by-key-a1b: enrolled SSH approvers, and approving a mint request by signature

- **Status**: Planned. Designed in iteration 375, **r0**. Queue row 201, `w-world-approvals-by-key`.
  **Sprint boundary: sprint 1 builds M1–M3 (enrolment). Sprint 2 builds M4–M6 (approving by key),
  from this same doc without re-design (§0, §8).**
  This doc covers the **A1b half of parent milestone A1**. Parent:
  [`w-world-approvals-by-key.md`](w-world-approvals-by-key.md). Sibling (built):
  [`w-world-approvals-by-key-a0a1.md`](w-world-approvals-by-key-a0a1.md) (A0 and A1a; PR #533 and
  PR #535). Their direction and rulings are settled and not reopened here.
- **Authority**: **D-WORLD-89** (Mark, attended 2026-10-21, "all defaults"). R1: a verified
  signature from an enrolled key may approve a mint, as well as the typed `y`. R2: the formats are
  WebAuthn and OpenSSH SSHSIG. R3: enrolment and revocation stay typed-only. R4: publishing stays
  typed-only. **D-WORLD-109 = A** (Mark, attended 2026-11-06): the loop builds and judges row 201
  against **fakes, temp stores and test keys only**. It never enrols a real key, never touches
  `~/.ssh`, never opens a real credentials store, never mints in one and never publishes. Mark's
  own enrolment and his first approval with his own key are the **attended exit run** (§8). The
  loop does not run them.
- **Item**: three things.
  - `ailang-worldd approver enroll|list|revoke` and an `approvers` table at credentials schema 5.
  - `approve <id> --signature <file>`, `--sign-with <key>` and `--write-message <file>`. These mint
    from a stored request with **no terminal**, once an enrolled, unrevoked key's SSHSIG over the
    request's exact message verifies.
  - The approver recorded as `approved_by = ssh:<name>`, shown by `approve --list --decided`.
  It changes no `/v1` route and no `.ail` code.
- **Designer lane**: `claude:claude-opus-5-5`.
- **Estimated**: 6 milestones in two sprints, about 3,750 lines in all (§8). **Sprint 1 = M1–M3**,
  about 1,900 lines: the pure enrolment law, store schema 5, and the `approver` verbs.
  **Sprint 2 = M4–M6**, about 1,850 lines: the pure signature law and the store consume, the
  `approve` signature paths, and `--sign-with` with its real-tool tests and CI step. Measured
  calibration: A0/A1a's two sprints overran their estimates by 1.6× and 2.2× (V24). §8 therefore
  names a fallback cut for each sprint.
- **Measurement base**: `3a96cf75` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter375-row201-a1b`, branch
  `sprint/world-iter375-row201-a1b`. sha256 prefixes (V1): `approve.go` `a215b26ad7ed060a`,
  `session_new.go` `2e9baa1c8c349750`, `mint_request.go` `d12b6929cd05fe8d`,
  `host/store/mint_requests.go` `bfe6bdc86f22b957`, `host/store/credentials.go`
  `b1d0e10d183e8d5c`, `host/authority/mint.go` `49d1aed46251752d`, `host/approval/sshsig.go`
  `a40850200782b3e4`, `host/approval/admit.go` `e139596fcb089057`.
- **Instruments**: `go1.26.6 darwin/arm64`; `/usr/bin/ssh-keygen`, `OpenSSH_10.3p1, LibreSSL
  3.3.6`; pinned `AILANG v0.41.0` at `$HOME/.pinned-ailang/ailang`. Every `ssh-keygen` probe ran in
  one `mktemp -d` directory, with `HOME` set to it, `SSH_AUTH_SOCK` unset, `-N ''`/`-f <tempdir
  path>` and `< /dev/null`. The directory and the one temporary `ssh-agent` were removed (V32).
  The designer's shell has no controlling terminal (V33), so no probe could prompt.

---

## 0. Scope, and the sprint boundary

**In A1b:**
- **Enrolment, the root of trust (parent R3).** `approver enroll --name <name> --format ssh
  --pubkey <file>` asks `y` at `/dev/tty` and stores the key. `approver list` is read-only.
  `approver revoke <name>` asks `y` at `/dev/tty` and marks the row revoked. It never deletes the
  row (§2.1). Only the `ssh` format exists in A1b. `--format webauthn` is refused with a message
  naming A2 (§2.10).
- **Approving by signature.** `approve <id> --signature <file>` runs the existing admission
  (`Admit`) and then verifies the signature against the enrolled, unrevoked key that the signature
  names. It then runs the existing re-derivation and the same `provisionEpisode`, with **no
  terminal**, and records `approved_by = ssh:<name>`.
  - `--sign-with <key>` runs `ssh-keygen -Y sign -n world-approval@sunholo -f <key>` over the
    exact message and feeds the result into the `--signature` path.
  - `--write-message <file>` writes the 90-byte message so the operator can sign it elsewhere.
- **The record.** `mint_requests.approved_by` already exists (schema 4). A1b writes
  `ssh:<name>` there, in the same transaction as a guarded `last_used_at` update of the approver
  row. `approve --list --decided` is the CLI decision list. `approver revoke` lists the requests
  the revoked key approved.
- **The threats a0a1 §4 left to A1b**: a revoked key, enrolment without the typed `y`, and the R6
  residual (SSHSIG proves possession of a key, not that a human is present) (§4, §9).
- **Carried forward from iteration 374** (§6): the untested `CR-TOOMANY` branch gets a test in M1.
  The mutation table has one mutant per new branch. Every struct that a security check claims to
  cover completely gets a reflective completeness test. Exit codes use the measured constants.

**Out of A1b** (only the seams are named, §2.10): A2 WebAuthn, A3 the surface and the Tauri
*needs you* flow (including a graphical decision list), A4 `Human.Approve` packets. `sk-*`
security-key signatures stay refused (Q1, the a0a1 default).

**Sprint boundary.** Sprint 1 (M1–M3) is **enrolment**. It contains the pure enrolment law, the
schema 5 `approvers` table with its store methods, and the three `approver` verbs. It does not
change how a session is minted. Sprint 2 (M4–M6) is **approving by key**. It contains the pure
signature law, the consume's approver binding, the three `approve` flags and the CI step. Sprint 2
depends on sprint 1. Nothing in sprint 1 depends on sprint 2. The cut follows a0a1 §0's own
reasoning: enrolment is the root of trust and gets its own review.

## 1. What the code does today (measured at `3a96cf75`)

1. **Exit codes.** `exitOK = 0`, `exitUsage = 1`, `exitFatal = 2` (`cmd/ailang-worldd/main.go:72-74`,
   V2). **Correction to a0a1:** a0a1 §2.5, §2.6a and CL2/CL7/CL13 said refusals exit 2. As built,
   every refusal before the write returns `exitUsage` (1). That covers a decline, an unanswered
   prompt, an expired request and a refused admission (`session_new.go:410`, `:443`, `:460`, `:468`;
   `approve.go:92`). Only a failure at or after the credentials write returns `exitFatal` (2)
   (`session_new.go:477`, `:492`, `:559`). **This doc uses 1 for every refusal and 2 for a failure
   at or after the write.**
2. **`approve` today** (`cmd/ailang-worldd/approve.go`, 222 lines, V3):
   - `runApprove` (`:88`) parses `--db`, `--credentials-db` and `--list`, and checks the id
     against `store.MintRequestIDPattern` (`:117`).
   - It opens the credentials store read-only under a 3 s root (`:130`, `approveReadBudget`,
     `mint_request.go:24`). It reads the row (`:152`) and closes the handle (`:153`).
   - `approveAdmission` (`:54`) decodes the row and runs `approval.Admit`. It maps each verdict to
     a refusal (`:64-82`).
   - It requires `--db`/`--credentials-db` to name the request's stores (`:169-177`). It
     re-derives the request through `prepareMint(fromRequest(mr))` (`:180`), fetches the daemon's
     health and runs `checkDaemon` (`:191-197`). It refuses unless
     `Canonical(fresh) == row.Request` (`:204-207`).
   - It sets `req.consume = &authority.RequestMint{…, Consume: store.RequestConsume{ID, Digest,
     ApprovedBy: "tty", Clock: env.now}}` (`:213-214`) and calls `provisionEpisode` (`:215`).
     **`"tty"` is the only `approved_by` any production code writes** (V26).
3. **`provisionEpisode` always opens the terminal first** (`session_new.go:403` func, `:406`
   `env.openTerminal()`), then prints the y/N.
   - With `answerBy != 0` the read is bounded (`:454-465`). A non-`y` answer refuses
     (`:466-469`).
   - It then takes the writer lock under a 75 s root (`:472`) and mints. With `r.consume != nil`
     it calls `authority.MintForRequest` (`:506-507`).
   - `mintRequest.consume` and `.answerBy` are at `:228` and `:230` (V4).
   - **There is no path through `provisionEpisode` that skips the terminal.** A1b adds the first
     one, and §5 argues it is reachable only with a verified signature.
4. **`sessionEnv`** has four seams: `openTerminal`, `now`, `random` and `after`
   (`session.go:25-38`). The production value opens `/dev/tty` `O_RDWR` (`:49`) (V5). It has no
   signing seam.
5. **The credentials store is schema 4** (`host/store/credentials.go:42`). Its parts (V6):
   - The exact table set per version is `credentialsTables` (`:71-80`). Schema 4 holds
     `episode_starts`, `mint_requests` and `session_credentials`.
   - `enforceSchema` accepts a known version only with exactly that table set (`:318`).
   - `upgradeCredentials` adds each later version's statements in one transaction. Its last step
     is `if version < 4 { mintRequestsSQL }` (`:340-367`).
   - **Only `serve` upgrades a store** (`OpenCredentials` "never upgraded here", `:188-199`;
     `ensureCredentialsStore`, `host/daemon/credentials.go:59`; log line `:76`) (V29).
   - **A bound, started or resumed mint refuses on any store below the binary's version**
     (`credentials.go:432`). So raising the version to 5 makes a v4 store refuse those mints until
     `serve` restarts with the new build. `episode start` already forces that restart through
     `checkDaemon`'s same-build check, so this adds no operator step (R3 of a0a1, and R2 here).
6. **The mint-request store** (`host/store/mint_requests.go`, V7):
   - `mintRequestsSchemaVersion = 4` (`:36`), `MaxPendingMintRequests = 16` (`:39`),
     `ErrTooManyPendingRequests` (`:49`), `ErrIncompleteConsume` (`:56`).
   - `OlderCredentialsStoreError.Error()` hard-codes "mint requests need schema %d" with
     `mintRequestsSchemaVersion` (`:65-69`). A test pins "need schema 4"
     (`mint_requests_test.go:202`) (V28).
   - `PutMintRequest` (`:88`), `MintRequest` (`:142`). `PendingMintRequests` (`:162`) filters
     `consumed_at IS NULL AND expires_at > ?` (`:171-172`). **No method lists consumed rows.**
   - `RequestConsume` is `{ID, Digest, ApprovedBy string; Clock func() int64; Check func(int64)
     error}` (`:234-239`).
   - `MintSessionForRequest` (`:243`) runs a guarded `UPDATE … WHERE request_id = ? AND digest = ?
     AND consumed_at IS NULL AND created_at <= ? AND expires_at > ?` (`:265`). **That WHERE does not
     bind the challenge column.** `classifyConsume` is at `:287`.
7. **The mint core** (`host/authority/mint.go`, V8):
   - `RequestMint{Request approval.MintRequest; RequestCreatedAt int64; Consume
     store.RequestConsume}` (`:60-64`) and `MintForRequest` (`:70`).
   - `mintRecorded`'s consume branch sets `Consume.Check` to `approval.MintMatches` and calls
     `st.MintSessionForRequest` (`:114-124`). It passes the caller's `ApprovedBy` through unread.
8. **The pure approval core** (`host/approval`, V9):
   - `Admit` (`admit.go:44`), `MintMatches` (`:92`), `SignedMessage` (`request.go:211`).
   - `Namespace = "world-approval@sunholo"` (`sshsig.go:19`).
   - `SSHPublicKey{Type string; Blob []byte}` (`:46-49`).
   - `parseKeyBlob` accepts only `ssh-ed25519` and `ecdsa-sha2-nistp256`. It refuses `sk-*` with a
     message naming Q1 (`:103-149`).
   - `Fingerprint` (`:151`). `ParseAuthorizedKey` (`:157`) **keeps no comment**: it returns only
     `Type` and `Blob`.
   - `ParseSSHSIG` (`:182`), `VerifySSHSIG` (`:262`), `VerifyApproval` (`:353`).
   - `imports_test.go` forbids `os`, `os/exec`, `crypto/rand`, `net` and others in the package's
     non-test files (V31).
9. **The request path's `CR-TOOMANY` branch is untested at the CLI** (V10, V11):
   - `createMintRequest` (`mint_request.go:101`) maps `errors.As(err, &older) || errors.Is(err,
     store.ErrTooManyPendingRequests)` to `exitUsage` (`:142-144`).
   - I applied the mutant `false && errors.Is(…)` to a copy. `go test ./cmd/ailang-worldd/ -run
     'Request|Approve|EpisodeStart' -count=1` printed `ok` both before and after the mutant. The
     file was restored and `git status` was clean.
   - No `cmd/ailang-worldd` test names `ErrTooManyPendingRequests` or "already pending" (0 matches).
     The store test does (`mint_requests_test.go:257`).
10. **Fixtures and the tracked-binary gate.** `host/approval/testdata/sshsig/` holds
    `message.hex`, not `message.bin` as a0a1 §2.8 planned (V16). The reason: `generate.sh` says the
    tracked-binary gate rejects any committed blob that holds a NUL byte. That gate is
    `scripts/tracked_binary_gate.sh`, run from `scripts/verify_go.sh:228-240`. The 90-byte message
    holds a NUL after its prefix, so **a fixture that holds the binary message, or any other NUL,
    cannot be committed**. `.gitattributes` holds `host/approval/testdata/sshsig/* -text`.
11. **CI has the real tool.** At `3a96cf75`, CI run `42753484556` (conclusion `success`) found
    `/usr/bin/ssh-keygen` (`OpenSSH_9.6p1 Ubuntu-3ubuntu13.20`) on `ubuntu-24.04`. It printed
    `--- PASS` for `TestSSHSIGCrossCheckWithSSHKeygen` and both subtests (V15). **a0a1's R1
    ("UNVERIFIED") is resolved.**
12. **Nothing of A1b exists yet** (V23). No production Go file has an `approvers` table, an
    `approver` verb, `sign-with`, `write-message` or an `ssh:` approver, and none execs
    `ssh-keygen`. The controls in the same command find `case "approve"` and the cross-check test.
13. **"The decision list" does not exist in code.** The parent's "You · Minted · 12 h · iPhone
    passkey" is a surface (A3) concept. Today the only reader of `approved_by` is the replay
    refusal text (`approve.go:71`) (V26).
14. **The verb pins.**
    - `help_test.go:14-15` lists `{"approve"}` in `everyVerb`. `:44` wants `"approve
      <request-id>", "approve --list"` in the usage.
    - `main.go:124-125` has the usage lines, `:442-449` the `case "approve"` that refuses the
      global `--addr`, and `:664` `ownHelpVerbs`.
    - `keyguard.go:77` has `case "approve"` in `ownVerbHelp` (V14).
15. **Context roots** (`host/store/context_roots_test.go`). `:19-20` pin
    `mint_request.go|createMintRequest` and `approve.go|runApprove`, one each. `:28-36` pin the
    session verbs (V13).
16. **Sprint sizes measured** (V24). Code, tests and fixtures, excluding plan docs: sprint 1
    (PR #533) **2,291** lines against an estimate of about 1,450. Sprint 2 (PR #535) **3,473**
    changed lines against about 1,600.

## 2. Design

### 2.1 The `approvers` table (credentials schema 5)

```sql
-- Approvers (version 5, row 201 A1b). A key enrolled at the daemon's terminal (a typed y,
-- parent R3) may approve a mint request by signature. Revoking MARKS the row and never deletes
-- it, so a name in mint_requests.approved_by ('ssh:<name>') names one key forever, and neither
-- a name nor a key is ever enrolled twice. format is checked in code (host/approval), not here,
-- so A2 can add 'webauthn' (and sign_count by ADD COLUMN) without rebuilding the table.
CREATE TABLE IF NOT EXISTS approvers (
    name         TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
    format       TEXT NOT NULL CHECK (format <> ''),
    public_key   BLOB NOT NULL UNIQUE CHECK (length(public_key) > 0),
    enrolled_at  INTEGER NOT NULL CHECK (enrolled_at > 0),
    last_used_at INTEGER CHECK (last_used_at IS NULL OR last_used_at >= enrolled_at),
    revoked_at   INTEGER CHECK (revoked_at IS NULL OR revoked_at >= enrolled_at)
);
```

- **`public_key` is the SSH wire blob** (`SSHPublicKey.Blob`). The key type and fingerprint are
  derived on read (`parseKeyBlob`, `Fingerprint`). There is no stored copy that could drift from
  the blob.
- **No comment column.** `ParseAuthorizedKey` keeps no comment (§1.8). A comment would also be
  operator-supplied text printed on a y/N, which is a terminal-injection surface (§4). The name is
  the label.
- **Revoke marks, never deletes** (Q3). Four reasons:
  - The audit row stays.
  - `ssh:<name>` stays unambiguous: `name` is the primary key and is never reused.
  - A revoked key can never be re-enrolled: `public_key` is `UNIQUE` across every row, revoked
    ones included. A revoked key is presumed compromised.
  - The guarded consume (§2.5) can tell "revoked" from "unknown".
- **The upgrade.** `credentialsSchemaVersion` goes 4 → 5. `credentialsTables(5)` is
  `["approvers", "episode_starts", "mint_requests", "session_credentials"]`, sorted by name as
  `applicationTables` returns them. Its query is `type = 'table'` excluding `sqlite_%` (V30), so the
  `UNIQUE` autoindex is not a table. `upgradeCredentials` gains `if version < 5 {
  approversSQL }`, a `CREATE TABLE` byte-identical (normalised) to the schema file's, under the
  existing canonical-DDL test pattern (`credentials_test.go:53-65`).
- **`OlderCredentialsStoreError` is generalised.** It gains `Need int` and `Feature string`, and
  `Error()` prints `"%s need schema %d"`. The four existing sites (`mint_requests.go:93`, `:147`,
  `:167`, `:254`) set `Need: 4, Feature: "mint requests"`, so their text is byte-identical. The
  `"need schema 4"` pin at `mint_requests_test.go:202` stays green unedited. The approver sites set
  `Need: 5, Feature: "approvers"`.

### 2.2 The pure enrolment law (`host/approval/approver.go`)

```go
const (
	FormatSSH       = "ssh"
	FormatWebAuthn  = "webauthn" // A2: named, refused in A1b
	MaxKeyFileBytes = 8 << 10
)

var ApproverNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,47}$`)

type Approver struct {
	Name       string
	Format     string       // FormatSSH
	Key        SSHPublicKey // Type + wire Blob
	EnrolledAt int64
	LastUsedAt int64 // 0: never used
	RevokedAt  int64 // 0: active
}

var (
	ErrApproverName        = errors.New("approval: approver name must match ^[a-z0-9][a-z0-9._-]{0,47}$")
	ErrFormatNotYet        = errors.New("approval: webauthn approvers arrive with row 201 A2; A1b enrols ssh keys only")
	ErrApproverFormat      = errors.New("approval: approver format must be ssh")
	ErrApproverKeyFile     = errors.New("approval: the public key file must hold exactly one key line of at most 8 KiB")
	ErrPrivateKeyFile      = errors.New("approval: that is a PRIVATE key file; pass the .pub file (its contents were not shown or kept)")
	ErrApproverRevoked     = errors.New("approval: the approver is revoked")
	ErrApproverClockBehind = errors.New("approval: the clock is behind the approver's enrolment time")
)

// NewApprover is the enrolment law: it parses one public key line into an Approver, or refuses.
func NewApprover(name, format string, keyFile []byte, now int64) (Approver, error) {
	if !ApproverNamePattern.MatchString(name) { return Approver{}, ErrApproverName }            // E1
	if format == FormatWebAuthn { return Approver{}, ErrFormatNotYet }                           // E2
	if format != FormatSSH { return Approver{}, ErrApproverFormat }                              // E3
	if len(keyFile) == 0 || len(keyFile) > MaxKeyFileBytes { return Approver{}, ErrApproverKeyFile } // E4
	if bytes.HasPrefix(bytes.TrimLeft(keyFile, " \t\r\n"), []byte("-----BEGIN")) {
		return Approver{}, ErrPrivateKeyFile                                                     // E5
	}
	line := bytes.TrimSuffix(keyFile, []byte("\n"))
	if bytes.ContainsAny(line, "\n") { return Approver{}, ErrApproverKeyFile }                  // E6
	key, err := ParseAuthorizedKey(line)
	if err != nil { return Approver{}, err }                                                     // E7 (ErrKeyType for sk-*, rsa, p384; ErrAuthorizedKey; ErrBase64)
	return Approver{Name: name, Format: FormatSSH, Key: key, EnrolledAt: now}, nil
}

// CheckRevoke is the revocation law, run before the y/N (the store re-checks it in SQL).
func CheckRevoke(a Approver, now int64) error {
	if a.RevokedAt != 0 { return ErrApproverRevoked }                                            // R1
	if now < a.EnrolledAt { return ErrApproverClockBehind }                                     // R2
	return nil
}

// RenderEnrolment is the body of the enrolment y/N: what will be stored, and what it proves.
func RenderEnrolment(a Approver) string
```

`RenderEnrolment` prints the name, the format, the key type, the `SHA256:` fingerprint, and the
`approved_by` value the key will record (`ssh:<name>`). It also prints the R6 text verbatim:

> This key will approve mint requests with no terminal. A signature proves that something holding
> the private key signed. It does not prove that you were there. Enrol only a key that cannot sign
> without you: one behind a passphrase, in an agent that asks before each use (`ssh-add -c`), or in
> a Secure Enclave or hardware holder (for example Secretive). Compare the fingerprint above with
> `ssh-keygen -l -f <your key>.pub` or `ssh-add -l` before you type y.

The text is pinned by a golden test (AP5). `RenderEnrolment` has no branches.

### 2.3 The pure signature law (`host/approval/signed.go`)

```go
// SignedApproval is a verified signature approval. Its fields are unexported, so only
// ApproveBySignature, in this package, can return one with ok set; the zero value binds nothing.
type SignedApproval struct {
	ok        bool
	digest    [32]byte
	challenge [32]byte
	expiry    int64
	name      string
	format    string
	key       []byte
}

type SignatureInput struct {
	Admit     AdmitInput // the existing admission input (row times, consumed, both digests)
	Challenge [32]byte   // the request row's challenge
	Signature []byte     // the armored SSHSIG
	Approver  Approver   // the enrolled row the signature's key names
}

type NotAdmittedError struct{ Verdict Verdict }

var (
	ErrNotVerified = errors.New("approval: not a verified signature approval")
)

type BindingError struct{ Field string } // "digest" | "challenge" | "expiry"

func ApproveBySignature(in SignatureInput) (*SignedApproval, error) {
	if v := Admit(in.Admit); v != Admitted { return nil, &NotAdmittedError{Verdict: v} }        // P1
	if !ApproverNamePattern.MatchString(in.Approver.Name) { return nil, ErrApproverName }          // P2
	if in.Approver.Format != FormatSSH { return nil, ErrApproverFormat }                           // P3
	if in.Approver.RevokedAt != 0 { return nil, ErrApproverRevoked }                               // P4
	if in.Admit.Now < in.Approver.EnrolledAt { return nil, ErrApproverClockBehind }                // P5
	if err := VerifyApproval(in.Signature, in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt,
		in.Approver.Key); err != nil { return nil, err }                                           // P6
	return &SignedApproval{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge,
		expiry: in.Admit.ExpiresAt, name: in.Approver.Name, format: in.Approver.Format,
		key: bytes.Clone(in.Approver.Key.Blob)}, nil
}

func (s *SignedApproval) Binds(digest, challenge [32]byte, expiry int64) error {
	if s == nil || !s.ok { return ErrNotVerified }                                                 // B1
	if s.digest != digest { return &BindingError{Field: "digest"} }                                // B2
	if s.challenge != challenge { return &BindingError{Field: "challenge"} }                       // B3
	if s.expiry != expiry { return &BindingError{Field: "expiry"} }                                // B4
	return nil
}

// ApprovedBy is the approved_by value: format ":" name ("ssh:mark-mac"); "" unless verified.
func (s *SignedApproval) ApprovedBy() string {
	if s == nil || !s.ok { return "" }                                                             // B5
	return s.format + ":" + s.name
}

// Use is what the consume binds in SQL: the approver's name, format and key (a copy).
func (s *SignedApproval) Use() (name, format string, key []byte, ok bool) {
	if s == nil || !s.ok { return "", "", nil, false }                                             // B6
	return s.name, s.format, bytes.Clone(s.key), true
}
```

- **Order.** Admission comes first, so a replayed request says "already used", not "bad
  signature". The approver's state comes before the cryptography, so a revoked key says "revoked"
  even when its signature is good.
- **Which approver.** `approve` finds the row by the signature's own key blob
  (`ApproverByKey`, §2.4). P6's `VerifySSHSIG` then requires that blob to be **byte-equal** to the
  row's key (`ErrKeyMismatch`, `sshsig.go:262`). A lookup bug cannot make a signature verify
  against a key that did not make it.
- **The signed values are the stored ones.** The digest that is verified is `StoredDigest`, which
  P1 has just shown equals the SHA-256 of the stored bytes. The expiry is the row's `expires_at`,
  and the challenge is the row's. These three are exactly what `SignedMessage` encodes.

### 2.4 Store methods (`host/store/approvers.go`)

```go
const approversSchemaVersion = 5

type ApproverRow struct {
	Name       string
	Format     string
	PublicKey  []byte
	EnrolledAt int64
	LastUsedAt int64 // 0 = NULL
	RevokedAt  int64 // 0 = NULL
}

var (
	ErrApproverNameTaken   = errors.New("store: an approver with that name exists (names are never reused, even after revoke)")
	ErrApproverKeyEnrolled = errors.New("store: that key is already enrolled (a key is never enrolled twice, even after revoke)")
	ErrApproverUnknown     = errors.New("store: no such approver")
	ErrApproverRevoked     = errors.New("store: the approver is revoked")
	ErrApproverKeyChanged  = errors.New("store: the approver's key is not the key that signed")
	ErrApproverClockBehind = errors.New("store: the clock is behind the approver's enrolment time")
	ErrBadApprover         = errors.New("store: invalid approver row")
)
```

| Method | What it does |
|---|---|
| `EnrolApprover(ctx, row)` | requires a deadline. Refuses below v5 (`*OlderCredentialsStoreError{Need: 5, Feature: "approvers"}`). Validates the row, each check its own branch: an empty name, a format other than `ssh`, an empty key, `EnrolledAt <= 0`, a non-zero `LastUsedAt` or `RevokedAt` → `ErrBadApprover`. In one transaction it refuses a taken name (`ErrApproverNameTaken`) and an enrolled key (`ErrApproverKeyEnrolled`, wrapped with the holder's name), then inserts. The `UNIQUE` constraints are the backstop. Driver error text is never parsed |
| `Approver(ctx, name)` | point lookup, `(row, ok, err)`. Below v5 → `Older` |
| `ApproverByKey(ctx, blob)` | point lookup by `public_key`, `(row, ok, err)`. Below v5 → `Older` |
| `ListApprovers(ctx)` | every row, ordered by name, capped at 256. Below v5 → `Older` |
| `RevokeApprover(ctx, name, now)` | below v5 → `Older`. Runs `UPDATE approvers SET revoked_at = ? WHERE name = ? AND revoked_at IS NULL AND enrolled_at <= ?`. If that does not change exactly one row, it reads the row in the same transaction and returns `ErrApproverUnknown`, `ErrApproverRevoked` or `ErrApproverClockBehind` |
| `DecidedMintRequests(ctx, limit)` *(sprint 2, `mint_requests.go`)* | consumed rows, newest `consumed_at` first, capped at 64. Works at v4 and above |

The store does **not** import `host/approval` (0 imports today, V23 control). It holds its own
minimal checks. The name grammar and key parsing are the pure law's, run by the CLI before the
store is called. The store's `length(name) BETWEEN 1 AND 48` CHECK is the backstop.

### 2.5 The consume binds the challenge and the approver (sprint 2)

`RequestConsume` gains two fields:

```go
type RequestConsume struct {
	ID, Digest, ApprovedBy string
	Challenge              string       // NEW: 64 lowercase hex, the row's challenge
	Approver               *ApproverUse // NEW: nil = the typed y
	Clock                  func() int64
	Check                  func(consumedAt int64) error
}
type ApproverUse struct {
	Name, Format string
	PublicKey    []byte
}
var (
	ErrInconsistentApproval    = errors.New("store: approved_by does not match the approver the consume binds")
	ErrRequestChallengeChanged = errors.New("store: the mint request's challenge is not the one signed")
)
```

`MintSessionForRequest` adds these, in order:
- **C1** The existing completeness `if` also requires `isLowerHex64(rc.Challenge)`, giving
  `ErrIncompleteConsume`.
- **C2** `rc.Approver != nil` with an empty `Name`, `Format` or `PublicKey` gives
  `ErrIncompleteConsume`.
- **C3** `rc.Approver == nil && rc.ApprovedBy != "tty"` gives `ErrInconsistentApproval`.
- **C4** `rc.Approver != nil && rc.ApprovedBy != rc.Approver.Format + ":" + rc.Approver.Name`
  gives `ErrInconsistentApproval`.
- **C5** `rc.Approver != nil && c.version < approversSchemaVersion` gives
  `*OlderCredentialsStoreError{Need: 5, Feature: "approvers"}`.
- **C6** The guarded request `UPDATE` adds `AND challenge = ?`.
- **C7** `classifyConsume` gains `case challenge != rc.Challenge: ErrRequestChallengeChanged`,
  after the digest case.
- **C8** When `rc.Approver != nil`, it runs `UPDATE approvers SET last_used_at = ? WHERE name = ?
  AND format = ? AND public_key = ? AND revoked_at IS NULL AND enrolled_at <= ?` with the same
  in-transaction reading `at`. This comes after the request update and before the session insert.
  If that does not change exactly one row, `classifyApproverUse` reads the row in the transaction:
  - **C9** absent → `ErrApproverUnknown`;
  - **C10** `revoked_at` set → `ErrApproverRevoked`;
  - **C11** format or key differs → `ErrApproverKeyChanged`;
  - **C12** `enrolled_at > at` → `ErrApproverClockBehind`.
  
  Any of these rolls back the request update. The request stays pending, and no session row is
  written.

**C8 is the revocation check that cannot go stale.** `approve` verifies the signature on a
read-only handle with no lock (§2.7). The mint takes the writer lock later, inside
`provisionEpisode`. A `revoke` that commits in that window turns C8's update into zero rows, so
the mint is refused. This mirrors a0a1's in-transaction expiry clock.

**The store's law over `approved_by`** (C3, C4, C8): every `ssh:<name>` written to
`mint_requests.approved_by` is written in the same transaction as a one-row update of an active
`approvers` row with that name, format and key. `tty` is written only with no approver.

### 2.6 `authority.MintForRequest` derives `approved_by` (sprint 2)

`RequestMint` gains three fields: `Challenge [32]byte`, `ExpiresAt int64` (the row's) and
`Signed *approval.SignedApproval`. In `mintRecorded`'s consume branch (`mint.go:114`), before the
store call:

```go
if rm.Consume.ApprovedBy != "" { return …, ErrApprovedBySetByCaller }        // A1: the caller never names the approver
digest, err := hex.DecodeString(rm.Consume.Digest)
if err != nil || len(digest) != 32 { return …, ErrBadConsumeDigest }        // A2
consume.Challenge = hex.EncodeToString(rm.Challenge[:])
if rm.Signed != nil {                                                        // A3
	if err := rm.Signed.Binds([32]byte(digest), rm.Challenge, rm.ExpiresAt); err != nil { return …, err } // A4
	name, format, key, _ := rm.Signed.Use()
	consume.ApprovedBy = rm.Signed.ApprovedBy()
	consume.Approver = &store.ApproverUse{Name: name, Format: format, PublicKey: key}
} else {
	consume.ApprovedBy = "tty"                                               // A5
}
```

So `cmd/ailang-worldd` stops writing `ApprovedBy` (`approve.go:214` loses `ApprovedBy: "tty"`).
The value comes from the same field that decides whether a terminal was used (§2.8). The two
authority tests that pass `ApprovedBy: "tty"` (`host/authority/mint_request_test.go:56`, `:78`) are
edited to pass none and still assert `approved_by = tty` on the row (V34).

### 2.7 `approve` with a signature (sprint 2)

```
approve <request-id> [--db <path>] [--credentials-db <path>]
        [--signature <file> | --sign-with <key> | --write-message <file>]
approve --list [--decided] [--db <path>] [--credentials-db <path>]
```

**Flag rules** (each a branch, exit 1 with the help text):
- **F1** More than one of `--signature`, `--sign-with` and `--write-message`.
- **F2** Any of them together with `--list`.
- **F3** `--decided` without `--list`.

**Steps 1–5 are today's, with the approver lookup moved inside the read-only handle:**

1. Read the row, run `approveAdmission` (unchanged), and decode the row's challenge hex (**F4**:
   a challenge that is not 64 hex → refused).
2. **`--write-message <file>`** (**F5**/**F6**): writes `SignedMessage(digest, challenge,
   expires_at)`, 90 bytes, with `O_CREATE|O_EXCL`, mode `0644`. The message is not secret. An
   existing file is refused (F5) and never overwritten. It then prints the exact out-of-band
   command and the next step:
   `ssh-keygen -Y sign -n world-approval@sunholo -f <your key> < <file> > <file>.sig`, then
   `ailang-worldd approve <id> --signature <file>.sig --db … --credentials-db …`. It exits 0. It
   opens no terminal, writes no store row, and the request stays pending (F6).
3. **`--sign-with <key>`**:
   - **F7** `env.sign == nil` → `exitFatal`, an internal error, like `random`.
   - Otherwise it calls `env.sign(ctx, key, msg, stderr)`. `ctx` is a new root bounded by the
     request's remaining lifetime, `expires_at - env.now()` (§2.9). **F8** a signing error →
     refused, exit 1, naming `ssh-keygen`'s failure. The returned bytes enter step 4 as if read
     from a file. **World never reads the private key.** Only `ssh-keygen` does.
4. **`--signature <file>`**:
   - **F9** `Lstat` must show a regular file of at most 16 KiB, so a FIFO or `/dev/zero` is
     refused before any read.
   - **F10** A read error → refused.
   - **F11** `ParseSSHSIG` fails → refused, naming the typed error.
   - It looks up the approver by the signature's key blob on the **same read-only handle**, before
     it closes (`approve.go:153` moves after this lookup).
   - **F12** Below v5 → refused, naming the `serve` restart.
   - **F13** No such key → refused: "the signature's key `SHA256:…` is not an enrolled approver;
     see `ailang-worldd approver list`".
   - It then calls `ApproveBySignature`. **F14** an error → refused, naming the reason (revoked
     at T, bad signature, wrong namespace, and so on). All of this happens **before any network
     call**.
5. Unchanged: the same-store check, `prepareMint(fromRequest(mr))`, `fetchDaemonHealth`,
   `checkDaemon`, and canonical byte equality.
6. **The mint.** `req.consume = &authority.RequestMint{Request: mr, RequestCreatedAt:
   row.CreatedAt, Challenge: ch, ExpiresAt: row.ExpiresAt, Signed: signed, Consume:
   store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Clock: env.now}}`.
   - On the signature path `answerBy` stays 0, because there is no prompt to bound. Expiry is
     still enforced in the consume transaction by the in-transaction clock (a0a1 §2.4 [r3]).
   - On the typed path `answerBy = row.ExpiresAt` as today, and `Signed` is nil.

**`approve --list --decided`** (**F15**/**F16**) prints the consumed rows newest first: the request
id, episode, `approved_by`, the consumed time (RFC 3339) and the first 16 hex characters of the
credential id. With none it prints `no decided mint requests` (F15). This is the CLI decision list
the parent's §2.4 step 4 names. The surface's graphical one is A3, reading the same rows.

### 2.8 `provisionEpisode`'s no-terminal branch (sprint 2)

```go
signed := r.consume != nil && r.consume.Signed != nil
var term io.ReadWriteCloser
if !signed {                                   // V1: the typed path, unchanged
	term, err = env.openTerminal()
	… // unchanged refusal, defer term.Close()
}
…
if signed {                                    // V2: no prompt; the request and its approver go to stdout
	fmt.Fprint(stdout, r.presetPrompt)
	fmt.Fprintf(stdout, "%s and mint a session for episode %s with %d grant(s) [%s]%s, expiry +%ds: approved by %s\n",
		action, r.episode, len(grants), grantList(grants), bound, r.ttl, r.consume.Signed.ApprovedBy())
} else {
	… // the existing answerBy check, y/N, bounded or unbounded read, isYes: unchanged
}
```

- **`signed` cannot be true on any path but `approve --signature`/`--sign-with`.** `Signed` is
  set at one site (§2.7 step 6), from `ApproveBySignature`'s return. A `SignedApproval` with
  `ok == true` cannot be built outside `host/approval`, because its fields are unexported. Even a
  forged or zero `*SignedApproval` reaching `MintForRequest` is refused by A4's `Binds` (B1). The
  zero value's `ApprovedBy()` is `""`, and the store's C4 refuses that.
- The comment at `session_new.go:24-28` gains: "`approve --signature` runs steps 4–7 without
  step 2–3's terminal, only with a verified SSHSIG approval (row 201 A1b)".
- `mintRequestFieldMap["consume"]` (`mint_request.go:223`) is reworded ("set by approve from the
  stored row and, on the signature path, the verified approval"). No new `mintRequest` field is
  added, so CL14's reflection is unaffected (V12).

### 2.9 `--sign-with`: the signing seam (sprint 2)

```go
// sessionEnv gains:
sign func(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error)

// realSessionEnv: signWithSSHKeygen
func signWithSSHKeygen(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "ssh-keygen", "-Y", "sign", "-n", approval.Namespace, "-f", key)
	cmd.Stdin = bytes.NewReader(msg)
	var out limitedBuffer // 16 KiB cap
	cmd.Stdout, cmd.Stderr = &out, stderr
	if err := cmd.Run(); err != nil { return nil, fmt.Errorf("ssh-keygen -Y sign: %w", err) } // G1
	if out.overflow { return nil, errSignatureTooLarge }                                      // G2
	if out.Len() == 0 { return nil, errEmptySignature }                                       // G3
	return out.Bytes(), nil
}
```

Measured on this machine, with temp keys only (V17–V21):
- **Signing from stdin writes the signature to stdout** (318 bytes for Ed25519) and prints
  `Signing data on standard input` to stderr. Nothing is written beside the key. This is a0a1
  V12's stdin rule, which avoids the `Overwrite (y/n)?` trap.
- **`-f <x>.pub` with no agent still signs**, using the private file `x` beside it. With the
  private file absent and no agent it fails: `No private key found for public key …`, rc 255.
  With a temporary `ssh-agent` holding the key, `-f only.pub` signs, and the signature is
  **byte-equal** to the file-key signature (Ed25519 is deterministic). So `--sign-with <x>.pub`
  does **not** guarantee that the agent signed. The help text says so (R6).
- **A passphrase key with no terminal fails fast.** It printed `Enter passphrase …: Load key …:
  incorrect passphrase supplied`, rc 255, in 0 s. It did not hang. At a real terminal `ssh-keygen`
  prompts on `/dev/tty` itself. That is `ssh-keygen`'s prompt, not World's, and it is bounded by
  the context deadline.
- `ssh-keygen` is found through `PATH`. **PATH is not a trust boundary here.** Whatever runs, its
  output is verified by P6 against the enrolled key over World's own message. A substituted binary
  can only refuse or produce something that fails to verify. A binary that can reach the
  operator's agent could already sign anything (R6).
- **The context root** is a new `context.Background()` root in `approve.go` (`signForApproval`),
  with timeout `expires_at - now` seconds, added to `TestProductionContextRoots`.

### 2.10 The `approver` verbs (sprint 1)

`case "approver"` in `main.go` refuses the global `--addr` (**M1**), like `session` and `approve`.
It dispatches as follows:
- **D1** no subcommand → help, exit 1;
- **D2** an unknown subcommand → exit 1;
- **D3** `--help` → exit 0.

Each subcommand takes `--db` (default `~/.ailang/world/world.db`) and `--credentials-db`, like
`session`. It requires an existing world store (`requireExistingStore`).

**`approver enroll --name <name> --format ssh --pubkey <file>`**
- **N1** A flag error or a positional argument → help, exit 1. **N2** an empty `--name` or
  `--pubkey` → exit 1. **N3** `requireExistingStore` fails → exit 1.
- **N4** `Lstat(--pubkey)` must show a regular file of at most `MaxKeyFileBytes`. The file is read
  **once**, and those bytes are what is parsed, shown and stored.
- **N5** `approval.NewApprover` fails → exit 1, naming the reason. `ErrPrivateKeyFile` never
  echoes the bytes.
- **Pre-checks on a read-only handle under a 3 s root, before any terminal:**
  - **N6** below v5 → exit 1, naming the `serve` restart;
  - **N7** the name is taken (revoked names included) → exit 1;
  - **N8** the key is enrolled, revoked or not → exit 1, naming the holder.
  
  So the operator is never asked `y` for something that will be refused.
- **The fence:**
  - **N9** `env.openTerminal()` fails → `sessionTTYFix`, exit 1.
  - It prints `RenderEnrolment(a)` and `Enrol approver <name>? [y/N] `.
  - **N10** a non-`y` answer → aborted, exit 1.
  - The read is unbounded, like `session new`'s (R7 of a0a1, R5 here).
- **The write:**
  - **N11** `openCredentialsForSession` under a new 3 s root fails → exit 2, with the existing
    writer-busy text.
  - `EnrolApprover`: **N12** `ErrApproverNameTaken` or `ErrApproverKeyEnrolled` (a race with
    another enrol) → exit 1; **N13** any other error → exit 2.
- **N14** It prints `enrolled approver <name>: <key type> <SHA256:…>; approvals by it record
  approved_by=ssh:<name>`.

**`approver list`** (read-only, 3 s root)
- **L1** below v5 → exit 1, naming the restart.
- **L2** none → `no approvers`.
- **L3** otherwise one line per approver: name, format, key type, fingerprint, `enrolled <RFC
  3339>`, `last used <RFC 3339|never>`, and `active` or `revoked <RFC 3339>`.

**`approver revoke <name>`**
- **K1** not exactly one positional argument → exit 1.
- **K2** read-only lookup: unknown → exit 1. **K3** `CheckRevoke` fails (already revoked, or the
  clock is behind) → exit 1.
- **K4** the fence: no terminal → exit 1. It prints `Revoke approver <name> (<type> <fp>)? It will
  never approve again, and its name and key can never be enrolled again. [y/N] `. **K5** a non-`y`
  answer → aborted, exit 1.
- **K6** `RevokeApprover` returns `ErrApproverRevoked` or `ErrApproverUnknown` (a race) → exit 1;
  any other error → exit 2.
- **K7** It prints `revoked approver <name>`.
- **K8** *(sprint 2, with `DecidedMintRequests`)*: it then lists the decided requests whose
  `approved_by` is `ssh:<name>`, with the credential id prefix and the line `session revoke --db …
  <credential_id-hash>`. **K9** if there are none it prints `no session was approved by it`. It
  revokes no session itself (Q4).

**Why revoke is fenced.** Revoking only ever *removes* authority, and parent R3 makes it typed all
the same. The cost is that an emergency revoke needs the daemon's terminal. That is accepted (R3)
and recorded as R4.

**Seams named, not built (A2):** `--format webauthn` is the E2 refusal. A2 adds `FormatWebAuthn` to
the pure law, `sign_count` by `ALTER TABLE ADD COLUMN` at schema 6, `approved_by =
webauthn:<name>`, and the browser half of enrolment (a one-time code typed at this same y/N). The
store's C4 law (`format:name`) already covers `webauthn:` without a change. **A3** reads
`PendingMintRequests` (*needs you*) and `DecidedMintRequests` (the decision list).

## 3. Pure core, effects at the boundary (coding standards S1–S3)

| Piece | Where | Pure? | Tested by |
|---|---|---|---|
| `NewApprover`, `CheckRevoke`, `RenderEnrolment`, `ApproverNamePattern` | `host/approval/approver.go` | yes | AP1–AP6 |
| `ApproveBySignature`, `SignedApproval.{Binds, ApprovedBy, Use}` | `host/approval/signed.go` | yes | SA1–SA5 |
| the test signer (`approvaltest.Signer`) | `host/approval/approvaltest/` (non-production, imported only by `_test.go` files) | yes | TS1, TS3, XS2 |
| the store: `approvers`, the consume binding | `host/store/approvers.go`, `mint_requests.go` | effect | AS1–AS9, CS1–CS6 |
| `approved_by` derivation | `host/authority/mint.go` | no I/O beyond the store call | AU1–AU3 |
| `/dev/tty`, the clock, the timer, **the signer** | `sessionEnv` (`openTerminal`, `now`, `after`, **`sign`**) | seam | CLI tests inject them |
| `ssh-keygen` | `signWithSSHKeygen` in `cmd/` only; tests; the fixture generator | effect | XS1, XS2 |

- **Why not `.ail` (S1).** The cryptography cannot be written in AILANG. The pinned v0.41.0
  `std/crypto` has `sha256bytes`, `sha256hex`, `hmacsha256`, `rsa_verify_pkcs1v15` and
  `constanttimeequal`, and **0** builtins matching `ed25519|ecdsa|sha512` (V22). The enrolment law
  is a name regex plus a call into the Go SSH parser, and the signature law ends in `VerifySSHSIG`.
  Neither has an invariant Z3 could state without the cryptography. a0a1's Q2 (move `Admit` into
  the kernel at A4) is unchanged.
- **Why not a package (S3).** The consumers are `cmd/ailang-worldd`, `host/authority` and (A3) the
  daemon, all host Go. The approval core already lives in `host/approval`, and this extends it.
- **S2 kept.** `host/approval` stays standard-library only with no I/O. `imports_test.go` forbids
  `os`, `os/exec`, `crypto/rand` and `net` (V31). `regexp` and `bytes` are already allowed.
  `approvaltest` uses `crypto/ed25519` with a caller-supplied seed, so it needs no randomness. A
  new guard (TS3) asserts that no non-test `.go` file in the repo imports it.

## 4. Threats (A1b)

| Threat | What A1b does | Pinned by |
|---|---|---|
| **An agent enrols a key** | Enrolment is only `approver enroll`, behind `env.openTerminal()` and `isYes` (N9, N10). Agents have no controlling terminal (parent §1). `EnrolApprover` has one production caller, after `isYes` | CA3, CA4 |
| **The pubkey file is swapped between the operator's look and the enrol** | Read once (N4). The y/N shows the fingerprint of **those** bytes, and those bytes are stored. The prompt tells the operator to compare with `ssh-keygen -l`/`ssh-add -l` | AP5, CA2 |
| **A private key passed as `--pubkey`** | `ErrPrivateKeyFile` (E5). The contents are never echoed or kept | AP1, CA5 |
| **Terminal injection through the prompt** | The name grammar is `[a-z0-9._-]`. No key comment is stored or shown. The fingerprint is base64 | AP1, AP5 |
| **A revoked key approves** | P4 refuses before the cryptography. C8's guarded `UPDATE … AND revoked_at IS NULL` refuses inside the mint transaction | SA1, CS3, CR6 |
| **Revoke races approve** | `approve` verifies on a read-only handle. A revoke that commits before the mint's transaction makes C8 change zero rows, so `ErrApproverRevoked` rolls back and the request stays pending | CS3, CR6 |
| **Re-enrolling a revoked key, or reusing a name** | `UNIQUE(public_key)` and `PRIMARY KEY(name)` across all rows, with the pre-checks N7/N8 | AS2, CA6 |
| **`approved_by` claims a key that did not sign** | `approve` never sets it (A1). It comes from the verified `SignedApproval` (A3). The store writes `ssh:<name>` only with C8's one-row update of that name, format and key (C4, C8) | AU1, CS2, CS4 |
| **The no-terminal path reached without a signature** | `signed` is keyed on `consume.Signed`, settable only from `ApproveBySignature`. A zero or foreign `SignedApproval` fails `Binds` (B1) and has an empty `ApprovedBy` (B5), which C4 refuses | SA3, AU2, CR4 |
| **A signature for request R1 used for R2** | Each request's challenge is in the message (existing). C6 now also binds `challenge` in the consume's `WHERE` | SA2, CS5 |
| **A signature in another namespace (`git`, `file`)** | `VerifySSHSIG` namespace equality (existing, `ErrNamespace`) | SA1 (P6 arm), CR5 |
| **`--sign-with` hangs at a passphrase or touch prompt** | The context deadline is the request's remaining lifetime. A passphrase key with no terminal fails fast (V19) | CR9 |
| **A substituted `ssh-keygen` on `PATH`** | Its output is verified by P6 against the enrolled key over World's message. It can only refuse | CR8 (a fake signer returning another key's signature is refused) |
| **`--signature` from a FIFO, a device or a huge file** | `Lstat` regular, at most 16 KiB (F9) | CR7 |
| **`--write-message` overwrites a file** | `O_EXCL` (F5) | CR10 |
| **An agent process with access to an unlocked `ssh-agent` or an unencrypted key file** | **Not prevented. This is R6.** SSHSIG proves possession of the private key, not that a human is present. With an enrolled key in an agent that signs without confirmation, any process with `SSH_AUTH_SOCK` can run `approve --sign-with`. Mitigation is operator policy, stated at the enrolment y/N (§2.2): a passphrase, `ssh-add -c`, a Secure Enclave or hardware holder. World's exec sandbox denies agents reads of the operator's home (a0a1 V22). `sk-*` keys, whose signatures carry a touch bit, stay refused until an attended fixture exists (Q1) | stated (AP5 pins the text) |
| **A process that can write `credentials.db`** | Unchanged and out of scope. It can already insert a session row, and now an approver row | existing |
| **Clock skew** | One clock, the approving host's. P5 and C12 refuse a clock behind the enrolment. The request's expiry is enforced in the transaction (existing). A signer's clock is never read | SA1, CS3 |

## 5. By-construction argument

1. **No key is enrolled without a typed `y`.** `EnrolApprover` has exactly one production caller,
   `runApproverEnroll`. It runs after `env.openTerminal()` succeeded and `isYes` returned true
   (N9, N10). CA3's spy and CA4's `n` arm fail if either is bypassed.
2. **The terminal is skipped only with a verified approval for this exact request.**
   - `provisionEpisode` skips the terminal iff `r.consume.Signed != nil` (V1).
   - `approve` sets `Signed` only from `ApproveBySignature`. That function returns a non-nil
     value only after admission, an active approver and `VerifyApproval` over `(stored digest,
     row challenge, row expiry)` all pass (P1–P6).
   - `MintForRequest` re-checks `Binds` against the consume's own digest, challenge and expiry
     (A4).
3. **`approved_by = ssh:<name>` means `<name>`'s active key signed this request.** The value
   comes only from the `SignedApproval` (A1, A3). The store writes it only alongside a one-row
   guarded update of that approver's row, `revoked_at IS NULL` (C4, C8), in the transaction that
   consumes the request and inserts the session.
4. **`approved_by = tty` means a `y` was typed.** A5 sets `tty` exactly when `Signed == nil`, and
   then V1 opens the terminal and reads `y` (unchanged a0a1 §5.1).
5. **A revocation that commits before the mint's transaction wins** (C8, §2.5).
6. **Nothing minted differs from the request.** Unchanged from a0a1 §5.2: the same re-derivation,
   `MintMatches` and consume. The signature path runs steps 1–5 of §2.7 exactly as the typed path
   does. It only adds work before them.
7. **World never reads private key material.** `--sign-with` passes a path to `ssh-keygen`.
   `--pubkey` refuses a private key file without echoing it (E5).

## 6. Test plan

**Rules the planner must carry (iteration 374 and 373 retros):**
- **One mutant per new production branch.** The mutation table below has one row per new
  production `if`, `return` or `switch` arm in this design's own code sketches. Those are
  E1–E7, R1–R2, P1–P6, B1–B6, C1–C12, A1–A5, F1–F16, G1–G3, V1–V2, N1–N14, L1–L3, K1–K9, D1–D3,
  M1, and the `EnrolApprover`/`RevokeApprover` validation branches. It does not stop at mutants
  for the headline properties. The planner copies this list into the sprint plan. The evaluator
  applies each mutant and records which assertion fires (S6, "load-bearing criteria require
  mutations"). A branch the planner adds that is not listed here gets its own row in the plan.
- **Reflective completeness tests** (RC1–RC5), for every struct that a security check claims to
  cover completely, **including element types**. Each asserts its own field list, so a new field
  is red until someone classifies it.
- **Exit codes** are asserted by constant (`exitUsage`, `exitFatal`), never by the literal, so a
  wrong doc number cannot leak into a test.
- Where the deliverable is a refusal, the arm checks the typed error with `errors.Is`/`errors.As`
  **and** that nothing was written.

### Sprint 1

**M1a — `CR-TOOMANY` (carried from iteration 374), `cmd/ailang-worldd/approve_test.go`**

| Id | Test | Asserts | Mutant killed |
|---|---|---|---|
| CT1 | `TestEpisodeStartRequestRefusesWhenSixteenArePending` | seed 16 live requests with `PutMintRequest` (the rig's store, `approveT0`). `episode start … --request` → `exitUsage`; stderr contains `16 mint requests are already pending`; still 16 rows; no terminal opened | `false && errors.Is(err, store.ErrTooManyPendingRequests)` at `mint_request.go:142` (the request would return `exitFatal`); measured surviving today (V10) |

**M1b — the pure enrolment law, `host/approval/approver_test.go`**

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| AP1 | `TestNewApproverRefusals` | one arm per branch: bad names (`""`, `Mark`, `a b`, `-x`, 49 chars, `x\x1b[2J`) → `ErrApproverName` (E1); `webauthn` → `ErrFormatNotYet` (E2); `gpg` → `ErrApproverFormat` (E3); empty and 8 KiB+1 → `ErrApproverKeyFile` (E4); a derived private-key file (the armor line `-----BEGIN OPENSSH PRIVATE KEY-----` and base64 noise, made in test code) → `ErrPrivateKeyFile`, and the error text does not contain the noise (E5); two key lines → `ErrApproverKeyFile` (E6); real `rsa3072.pub`, `p384.pub` and a derived `sk-ssh-ed25519@openssh.com` retype → `ErrKeyType` (E7) | each of E1–E7 deleted; E2 placed after E3 (webauthn says "must be ssh" instead of naming A2) |
| AP2 | `TestNewApproverAcceptsTheRealFixtures` | the checked-in `testdata/sshsig/ed25519.pub` and `p256.pub` → `Format "ssh"`, the right `Key.Type`, `Key.Blob` equal to `ParseAuthorizedKey`'s, `EnrolledAt == now`, `LastUsedAt == RevokedAt == 0`; with and without a trailing newline | the name or `now` not copied; a trailing-newline file refused |
| AP3 | `TestCheckRevoke` | active → nil; `RevokedAt != 0` → `ErrApproverRevoked` (R1); `now == EnrolledAt` → nil; `now == EnrolledAt-1` → `ErrApproverClockBehind` (R2) | R1, R2 deleted; `<=` for `<` |
| AP4 | `TestApproverNamePatternBoundaries` | 1 and 48 characters accepted, 49 refused; `.`, `_`, `-` accepted after the first character | `{0,48}` |
| AP5 | `TestRenderEnrolmentGolden` | for the `ed25519.pub` fixture: golden text containing the name, `ssh`, `ssh-ed25519`, the `SHA256:` fingerprint from `SOURCE.txt`, `approved_by=ssh:<name>`, and the R6 paragraph verbatim | the fingerprint or the R6 text dropped |
| RC4 | `TestEnrolmentPromptShowsEveryPersistedField` | reflects over `Approver` and `SSHPublicKey` (the element type). Each field is either rendered by `RenderEnrolment` (checked by changing it and seeing the text change) or listed as excluded with a reason (`EnrolledAt`, `LastUsedAt`, `RevokedAt`: "timestamps, set by the store"; `Key.Blob`: "shown as its fingerprint"). The field count equals the dispositions | a new `Approver` field that the prompt never shows |
| — | `TestApprovalImportsOnlyTheStandardLibrary` (existing) | still green with `approver.go` | an `os` import |
| TS1 | `TestTestSignerVerifies` (`approvaltest`, sprint 1 because M3's CLI tests need keys) | `Signer` from a fixed seed: `AuthorizedKey()` parses with `ParseAuthorizedKey`, and `Sign(Namespace, msg)` verifies with `VerifySSHSIG` | — (its only job is to be correct; XS2 in M6 pins it to the real tool) |
| TS3 | `TestApprovalTestIsTestOnly` | walks the repo; no non-`_test.go` file imports `host/approval/approvaltest`. Instrument-health control: at least one `_test.go` file does | — (an import guard, labelled as a control) |

The **new real fixtures** for AP1 are `rsa3072.pub` and `p384.pub`. They live in a new
`host/approval/testdata/approvers/` and are made by `generate.sh` there:
`ssh-keygen -q -t rsa -b 3072 -N '' -f "$tmp/rsa3072" < /dev/null` and `ssh-keygen -q -t ecdsa
-b 384 -N '' -f "$tmp/p384" < /dev/null` in a `mktemp -d` that it deletes. Only the `.pub` files
and a `SOURCE.txt` (`ssh -V`, the time, `ssh-keygen -l`) are kept. These are one-line text files
with **no NUL**, so `scripts/tracked_binary_gate.sh` (`verify_go.sh:228-240`) passes. Their
SHA-256 values are pinned in `approver_test.go` (AP0, `TestApproverFixturesAreTheGeneratedBytes`),
and `.gitattributes` gains `host/approval/testdata/approvers/* -text`. Existing fixtures are not
regenerated, so SG1's pins do not move.

**M2 — store schema 5, `host/store/approvers_test.go` and edited credentials tests**

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| AS1 | `TestCredentialsSchemaMatchesCanonicalDDL` (edited) | four tables; `approvers` DDL equals the canonical copy and `approversSQL`; `user_version` 5 | `approversSQL` drifting from the schema file |
| AS2 | `TestEnrolApproverRefusals` | v4 handle → `*OlderCredentialsStoreError` with `Need == 5` and text `approvers need schema 5`; each `ErrBadApprover` arm (empty name, format `gpg`, empty key, `EnrolledAt 0`, `LastUsedAt != 0`, `RevokedAt != 0`); a taken name → `ErrApproverNameTaken`; the same key under a new name → `ErrApproverKeyEnrolled` naming the holder; both again after the first is **revoked** → still refused; 0 extra rows each time | each validation `if` deleted; the version check deleted; the pre-checks deleted (the `UNIQUE` backstop error is not `errors.Is` the typed error, so the arm still reds) |
| AS3 | `TestApproverLookups` | `Approver` and `ApproverByKey` find the row (all six fields equal), miss with `ok == false`, and return `Older` below v5; `ListApprovers` is ordered by name and capped at 256 | lookup by a prefix; missing order |
| AS4 | `TestRevokeApprover` | active → `revoked_at` set, the row kept; again → `ErrApproverRevoked`; unknown → `ErrApproverUnknown`; `now < enrolled_at` → `ErrApproverClockBehind`, unchanged; below v5 → `Older` | the `WHERE revoked_at IS NULL` dropped; a `DELETE` instead of a mark; the classify arms swapped |
| AS5 | `TestCredentialsUpgradeToV5` | stores at v1, v3 and v4 → v5, every existing row unchanged (sessions, starts, requests), idempotent; a session verb's `OpenCredentials` leaves v4 alone | the `version < 5` statement missing; a CLI open upgrades |
| AS6 | `TestOlderCredentialsStoreErrorText` | `Need 4` → the exact pre-A1b text (byte-equal to a pinned string); `Need 5` → `approvers need schema 5` and the restart hint | the generalised `Error()` changing the v4 text |
| RC2 | `TestApproverRowMatchesTheTable` | reflects over `ApproverRow` against `PRAGMA table_info(approvers)` through a declared field↔column map; enrols a row with a distinct non-zero value in every field (via raw SQL for `last_used_at`/`revoked_at`) and requires `Approver()` to return each one | a column dropped from the `SELECT` or the scan (that field comes back zero) |
| edited | `TestCredentialsMethodsRequireADeadline` | covers the 5 new methods | — |
| edited | `TestCredentialsUpgradeToV4`, `credentials_binding_test.go:106` | expectations move to `credentialsSchemaVersion` (5) where they pin "current" | — |

**M3 — the `approver` verbs, `cmd/ailang-worldd/approver_test.go`** (with the existing
`spyEnv`/`approveEnv`/`noTerminalEnv` patterns and a temp store; keys from
`approvaltest.Signer` or the checked-in `.pub` fixtures, never `~/.ssh`)

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| CA1 | `TestApproverDispatchAndHelp` | no subcommand → `exitUsage` + help (D1); `frob` → `exitUsage` (D2); `--help` → `exitOK` (D3); the global `--addr` → `exitUsage` (M1) | each branch deleted |
| CA2 | `TestApproverEnrollStoresWhatItShowed` | `y` → one row whose blob equals the file's parsed blob; the prompt (captured from the spy terminal) contains that blob's fingerprint and the R6 paragraph; stdout matches N14; the file is rewritten **after** the prompt is printed, and the stored blob is still the first one | reading the file twice (once to show, once to store) |
| CA3 | `TestApproverEnrollWithoutATerminalEnrolsNothing` | `noTerminalEnv` → `sessionTTYFix`, `exitUsage`, 0 rows (N9) | the fence bypassed |
| CA4 | `TestApproverEnrollDeclinedEnrolsNothing` | `n`, empty answer and EOF → `exitUsage`, 0 rows (N10) | `isYes` skipped |
| CA5 | `TestApproverEnrollRefusals` | flag misuse (N1), an empty `--name`/`--pubkey` (N2), a missing world store (N3), a FIFO or 8 KiB+1 `--pubkey` (N4), a private key file (N5: stderr lacks the file's bytes), `--format webauthn` (N5 via E2: stderr names A2), a v4 store (N6: names the restart), a taken name (N7), an enrolled key (N8: names the holder). Each → `exitUsage`, **no terminal opened** (spy), 0 rows | each pre-check deleted (the arm then reaches the spy terminal) |
| CA6 | `TestApproverEnrollRaceAndFailure` | an `EnrolApprover` seam returning `ErrApproverNameTaken` → `exitUsage` (N12); returning a plain error → `exitFatal` (N13); the writer lock held by another handle → `exitFatal` with the busy text (N11) | N12 mapped to `exitFatal`; N11's error swallowed |
| CA7 | `TestApproverList` | v4 → refused naming the restart (L1); empty → `no approvers` (L2); one active and one revoked → both lines with fingerprint, `never`/time, `active`/`revoked <time>` (L3); opens no terminal and holds no writer lock | the state column always `active` |
| CA8 | `TestApproverRevoke` | `y` → `revoked_at` set, `exitOK`, K7's line; unknown (K2), already revoked (K3), no terminal (K4), `n` (K5) → `exitUsage` and unchanged; arity (K1); `RevokeApprover` returning `ErrApproverRevoked` → `exitUsage`, a plain error → `exitFatal` (K6) | each branch deleted; revoke without the fence |
| edited | `TestEveryVerbHasHelp` (`help_test.go:14`, `:44`) | `{"approver"}`, `{"approver","enroll"}`, `{"approver","list"}`, `{"approver","revoke"}`; usage contains `approver enroll`, `approver list`, `approver revoke` | — |
| edited | `TestHelpIsNeverBlockedByTheRegistryKeyGuard` (`keyguard_test.go:32`) | `approver --help` is answered with the key set (`keyguard.go` `ownVerbHelp` gains `case "approver"`) | — |
| edited | `TestProductionContextRoots` | +3 entries, 5 roots in all: `approver.go\|runApproverEnroll\|Background: 2` (pre-check, write), `approver.go\|runApproverList\|Background: 1`, `approver.go\|runApproverRevoke\|Background: 2` | — |

### Sprint 2

**M4 — the signature law, the consume binding and authority**

`host/approval/signed_test.go`. The keys come from the checked-in fixtures, where the
`fixtureMessageHex` digest, challenge and expiry are the message; elsewhere they come from
`approvaltest.Signer`.

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| SA1 | `TestApproveBySignatureOrder` | one arm per branch: consumed + a good signature → `*NotAdmittedError{RefuseConsumed}` (P1); expired → `{RefuseExpired}`; a bad name (P2); format `webauthn` (P3); revoked + good signature → `ErrApproverRevoked` (P4); `Now < EnrolledAt` (P5); the real `git-ns-ed25519.sig` → `ErrNamespace`, `p256.sig` against the ed25519 approver → `ErrKeyMismatch`, a flipped byte → `ErrBadSignature` (P6). The clean arm returns non-nil | each branch deleted; P4 moved after P6 (a revoked key with a bad signature then says "bad signature"; the arm uses a good one) |
| SA2 | `TestSignedApprovalBinds` | the clean arm → nil; a zero `SignedApproval{}` and a nil pointer → `ErrNotVerified` (B1); digest, challenge or expiry changed by one → `*BindingError` naming that field (B2–B4) | each comparison deleted |
| SA3 | `TestSignedApprovalAccessorsFailClosed` | zero and nil → `ApprovedBy() == ""` (B5), `Use()` `ok == false` (B6); verified → `"ssh:<name>"`, and `Use()`'s key is a copy (mutating it does not change a second call) | B5/B6 deleted; the key returned uncloned |
| RC1 | `TestSignedApprovalFieldsAreAllBound` | reflects over `SignedApproval`. Every field has a disposition (`ok`: B1; `digest`, `challenge`, `expiry`: `Binds`; `name`, `format`: `ApprovedBy` and `Use`; `key`: `Use`), and each is exercised by perturbing that field in an otherwise-valid value. The element type of `key` (`[]byte`) is perturbed by one byte. The field count equals the dispositions | a field added that `Binds` does not compare |

`host/store` (`approvers_test.go`, `mint_requests_test.go` edited):

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| CS1 | `TestMintSessionForRequestNeedsTheChallenge` | an empty or 63-hex `Challenge` → `ErrIncompleteConsume` (C1); an approver missing `Name`, `Format` or `PublicKey` → `ErrIncompleteConsume` (C2) | the challenge clause dropped; C2 deleted |
| CS2 | `TestApprovedByMatchesTheApprover` | `Approver == nil` with `ssh:x` → `ErrInconsistentApproval` (C3); an approver `x` with `ApprovedBy "ssh:y"` or `"tty"` → `ErrInconsistentApproval` (C4); a v4 store with an approver → `Older{Need 5}` (C5). Each leaves 0 session rows and the request pending | C3, C4, C5 deleted |
| CS3 | `TestSignedConsumeChecksTheApproverInTheTransaction` | approver active → ok: the request is `approved_by = ssh:x`, the approver's `last_used_at == at`, 1 session row. **The clock seam revokes `x` (raw SQL on a second connection) just before `Clock()` returns** → `ErrApproverRevoked` (C10), request pending, 0 session rows, `last_used_at` unchanged; unknown name → `ErrApproverUnknown` (C9); another key or format → `ErrApproverKeyChanged` (C11); `enrolled_at > at` → `ErrApproverClockBehind` (C12) | the approver `UPDATE` deleted (C8: the revoked arm mints); `RowsAffected` unchecked; the update placed after the session insert and not rolled back; the classify arms swapped |
| CS4 | `TestTypedConsumeTouchesNoApprover` | `Approver == nil`, `ApprovedBy "tty"` → ok on a v5 store **with** an approver row, whose `last_used_at` stays NULL | the approver update run unconditionally |
| CS5 | `TestConsumeBindsTheChallenge` | the row's challenge rewritten by raw SQL after the consume is built → `ErrRequestChallengeChanged` (C6, C7), 0 session rows | `AND challenge = ?` dropped (it mints); C7's case deleted (it classifies as the generic "changed 0 rows") |
| CS6 | `TestDecidedMintRequests` | consumed rows only, newest `consumed_at` first, cap 64, `approved_by` and `credential_id` returned; works on v4 | the `consumed_at IS NOT NULL` filter dropped |
| RC3 | `TestRequestConsumeCompleteness` | reflects over `RequestConsume` **and its element type `ApproverUse`**. Each field has a declared disposition: required (zeroing it → `ErrIncompleteConsume`: `ID`, `Digest`, `Challenge`, `Clock`, `Check`, `ApproverUse.Name`, `.Format`, `.PublicKey`), derived-and-checked (`ApprovedBy`: C3/C4), or optional with a meaning (`Approver`: nil = tty). Each required field is zeroed in turn on an otherwise valid consume. The field count equals the dispositions | a new `RequestConsume` or `ApproverUse` field the completeness check ignores |

`host/authority/mint_request_test.go` (edited):

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| AU1 | `TestMintForRequestDerivesApprovedBy` | a caller-set `ApprovedBy` (`"tty"` or `"ssh:mallory"`) → `ErrApprovedBySetByCaller`, 0 rows (A1); a bad digest hex (A2); `Signed == nil` → `approved_by = tty`, no approver touched (A5); a verified `Signed` (from the test signer) → `approved_by = ssh:<name>`, `last_used_at` set (A3) | A1 deleted (A5 then silently overwrites the caller's value and the arm mints, so the arm's `errors.Is(err, ErrApprovedBySetByCaller)` fires); A2 deleted (the bad hex reaches `Binds` or the store instead of `ErrBadConsumeDigest`); A5's `"tty"` changed |
| AU2 | `TestMintForRequestRefusesAnUnboundApproval` | a `Signed` made for another request's digest, challenge or expiry → the `*BindingError`, request pending (A4); a zero `&approval.SignedApproval{}` → `ErrNotVerified` | A4 deleted |
| RC5 | `TestRequestMintFieldsAreSourced` | reflects over `authority.RequestMint` and its new fields. Each is sourced from the stored row (`Challenge`, `ExpiresAt`, `RequestCreatedAt`), the decoded request (`Request`), the verified approval (`Signed`) or the consume (`Consume`). A new field is red until classified | — (completeness) |

**M5 — the `approve` signature paths, `cmd/ailang-worldd/approve_sig_test.go`**

These use the existing `approveRig`. The approver is enrolled into the rig's temp store by the
store method directly, so the test needs no terminal. Keys come from `approvaltest.Signer`. The
signature is made over `SignedMessage(row digest, row challenge, row expires_at)`.

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| CR1 | `TestApproveFlagRules` | F1 (two of the three), F2 (with `--list`), F3 (`--decided` alone) → `exitUsage` + help, 0 changes | each rule deleted |
| CR2 | `TestApproveBySignatureMintsWithNoTerminal` | the spy's `openTerminal` **fails the test if called**; `exitOK`; the request row has `approved_by = ssh:<name>`; the approver's `last_used_at` is set; stdout has the request block and `approved by ssh:<name>`; and **the full a0a1 effect check `approveEffectProblems` (CL16) passes** on worktree, branch, session file, mode, token hash, MCP config, agent, model and task | V1 deleted (the spy fires); V2's line missing |
| CR3 | `TestApproveTypedPathUnchanged` | with no signature flag every a0a1 CL test stays green **unedited**, and this arm checks `approved_by = tty` with an approver enrolled | `Signed` set on the typed path |
| CR4 | `TestApproveBySignatureRefusalsBeforeAnyEffect` | each → `exitUsage`, no terminal, no network (the `stubDaemon` counts 0 requests), request pending, 0 session rows: a FIFO and 16 KiB+1 file (F9); unreadable (F10); garbage (F11); a v4 store (F12: names the restart); an unenrolled key (F13: names the fingerprint and `approver list`); a revoked approver (F14: "revoked at"); a signature over another request's message (F14: bad signature); consumed and expired requests → admission's existing refusals, before the signature is read | each branch deleted; the lookup moved after the network call (the stub counts 1) |
| CR5 | `TestApproveBySignatureRefusesTheWrongNamespace` | a test-signer signature under `git` → `exitUsage`, stderr names the namespace | — (P6 at the CLI) |
| CR6 | `TestApproveRevokeRaceIsRefusedAtTheWrite` | the `now` seam revokes the approver (raw SQL) on its call inside the consume → `exitFatal` (a failure at the write), stderr `approver is revoked`; the worktree this run created is removed (the existing step 7); request pending | C8 bypassed at the CLI |
| CR7 | `TestApproveDecidedList` | after one typed and one signed approve: `--list --decided` prints both with `approved_by` `tty` and `ssh:<name>`, newest first (F16); an empty store → `no decided mint requests` (F15) | F15/F16 swapped; pending rows included |
| CR8 | `TestApproveSignWithFakeSigner` | the `sign` seam: `nil` → `exitFatal` internal (F7); returns an error → `exitUsage`, stderr names it (F8); returns **another enrolled approver's valid signature over a different message** → refused (F14), nothing minted; returns the right signature → mints as `ssh:<name>`, and the seam received exactly `SignedMessage(...)` with `approval.Namespace` passed through to it | F7/F8 deleted; the seam's output trusted without verification |
| CR9 | `TestApproveSignWithDeadline` | the seam records `ctx.Deadline()`; with `now = expires_at - 37` the deadline is 37 s ahead | the bound taken from `RequestLifetime` |
| CR10 | `TestApproveWriteMessage` | writes exactly the 90 bytes of `SignedMessage` (F6), `exitOK`, prints the `ssh-keygen -Y sign -n world-approval@sunholo` line and the `--signature` next step, no terminal, request pending; an existing file → `exitUsage`, its contents unchanged (F5); a consumed request → admission refuses and no file is written | `O_EXCL` dropped; writing before admission |
| CR11 | `TestApproveChallengeHex` | the row's challenge rewritten by raw SQL to 64 non-hex characters → `exitUsage` (F4) | F4 deleted (a zero challenge flows on) |
| CR12 | `TestApproverRevokeListsWhatItApproved` | after a signed approve, `approver revoke` prints the request id, the credential prefix and a `session revoke` line (K8); for an approver that approved nothing, `no session was approved by it` (K9); no session is revoked | K8 filtering by the wrong `approved_by` |
| CR13 | `TestProvisionEpisodeSkipsTheTerminalOnlyWhenSigned` | a direct `provisionEpisode` call with `consume.Signed == nil` opens the terminal (spy); with a verified `Signed` it does not; with a zero `&SignedApproval{}` it skips the terminal and is **refused at A4** (`ErrNotVerified`, `exitFatal`, nothing minted, the created worktree removed) | V1 keyed on `r.consume != nil` alone |
| edited | `TestEveryVerbHasHelp` | usage contains `approve <request-id> --signature`, `--sign-with`, `--write-message`, `approve --list --decided` | — |
| edited | `TestProductionContextRoots` | `approve.go\|signForApproval\|Background: 1` | — |

**M6 — the real tool, `host/approval/approvaltest/sign_crosscheck_test.go` and
`cmd/ailang-worldd/approve_signwith_test.go`**

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| XS1 | `TestSignWithSSHKeygenEndToEnd` (subtests `ed25519`, `ecdsa-p256`) | `ssh-keygen -q -t … -N '' -f $tmp/k < /dev/null` in `t.TempDir()`, with `HOME=$tmp` and `SSH_AUTH_SOCK=` in the child's environment; enrol `k.pub` through the store; `approve <id> --sign-with $tmp/k` with the **production** `signWithSSHKeygen` → `exitOK`, `approved_by = ssh:<name>`; `--sign-with $tmp/k.pub` with `k` moved away and no agent → `exitUsage` naming `No private key found` (G1); a fake `ssh-keygen` script first on `PATH` that exits 0 and prints nothing → `exitUsage` (G3) | G1 deleted (the missing-key arm reaches the parser with no bytes and says "bad armor", which the arm distinguishes); G3 deleted |
| XS2 | `TestTestSignerAgreesWithSSHKeygen` | the `approvaltest.Signer` output for the fixture message is accepted by `ssh-keygen -Y check-novalidate -n world-approval@sunholo` (rc 0, `Good "world-approval@sunholo" signature`); `ssh-keygen -l -f` of its `AuthorizedKey()` prints the same fingerprint as `Fingerprint()`; a flipped byte → rc 255 | a test signer that is self-consistent but not SSHSIG (for example, hashing with SHA-256 while claiming `sha512`) |
| — | G2 (stdout over 16 KiB) | a fake `ssh-keygen` printing 16 KiB+1 → `exitUsage` | G2 deleted |

Both tests use the a0a1 skip pattern. With `WORLD_REQUIRE_SSH_KEYGEN=1` a missing `ssh-keygen` is
`t.Fatal`. Otherwise it is `t.Skip("SKIP COUNTED: …")`. The **CI step** at `ci.yml:141-155` is
extended in place:
- its `-run` becomes `'^(TestSSHSIGCrossCheckWithSSHKeygen|TestTestSignerAgreesWithSSHKeygen|TestSignWithSSHKeygenEndToEnd)$'`
  over `./host/approval/... ./cmd/ailang-worldd/`;
- its PASS-grep list gains `TestTestSignerAgreesWithSSHKeygen`,
  `TestSignWithSSHKeygenEndToEnd/ed25519` and `TestSignWithSSHKeygenEndToEnd/ecdsa-p256`.

The runner's `ssh-keygen` is measured present (V15), so the step is expected green.

**Must stay green, unedited** (the typed path): every a0a1 `TestApprove*`, `TestEpisodeStartRequest*`,
`TestMintRequestFieldMapIsComplete`, `TestApproveRebuildsTheSameMintRequest`, `TestSessionMint_*`,
`TestSessionNew*`, `TestEpisodeStart*`, `TestExistingCallersReadUnbounded`, `TestReadTermLineWithin`,
`TestSSHSIG*`, `TestDaemonDependencyAllowlist` and the quickstart tests. The exceptions are the
named edits: `approve.go:214`'s `ApprovedBy` removal, the two authority tests (V34) and the schema
version pins.

## 7. Conflict surface

**Sprint 1 (M1–M3):**

| File | Change |
|---|---|
| `host/approval/approver.go` (+ `approver_test.go`) | **new**, pure |
| `host/approval/testdata/approvers/{rsa3072.pub,p384.pub,SOURCE.txt,generate.sh}` | **new**, generated, text, no NUL |
| `.gitattributes` | +1 line (`host/approval/testdata/approvers/* -text`) |
| `host/approval/approvaltest/signer.go` (+ `signer_test.go`) | **new**, test support (moves to sprint 1 because M3's CLI tests need keys) |
| `host/store/credentials_schema.sql` | `approvers` table |
| `host/store/credentials.go` | version 5; `credentialsTables(5)`; `upgradeCredentials` `if version < 5` |
| `host/store/approvers.go` (+ `approvers_test.go`) | **new**: 5 methods, typed errors |
| `host/store/mint_requests.go` | `OlderCredentialsStoreError{Need, Feature}`; the 4 sites set them |
| `host/store/credentials_test.go`, `credentials_binding_test.go`, `mint_requests_test.go`, `context_roots_test.go` | table list, version pins, the deadline test, +5 roots |
| `host/daemon/credentials.go` | none needed: the log line already names the version generically (`:76`, V29) |
| `cmd/ailang-worldd/approver.go` (+ `approver_test.go`) | **new**: `runApprover`, enrol, list, revoke |
| `cmd/ailang-worldd/approve_test.go` | +CT1 |
| `cmd/ailang-worldd/main.go`, `keyguard.go`, `help_test.go` | `case "approver"`, usage lines, `ownHelpVerbs`, `ownVerbHelp`, `everyVerb` |

**Sprint 2 (M4–M6):**

| File | Change |
|---|---|
| `host/approval/signed.go` (+ `signed_test.go`) | **new**, pure |
| `host/store/mint_requests.go` | `RequestConsume.{Challenge, Approver}`, C1–C12, `DecidedMintRequests` |
| `host/authority/mint.go` (+ `mint_request_test.go`) | `RequestMint.{Challenge, ExpiresAt, Signed}`, A1–A5 |
| `cmd/ailang-worldd/session.go` | `sessionEnv.sign`, `signWithSSHKeygen`, `realSessionEnv` |
| `cmd/ailang-worldd/session_new.go` | V1/V2 in `provisionEpisode`; the caller comment at `:24-28` |
| `cmd/ailang-worldd/approve.go` (+ `approve_sig_test.go`, `approve_signwith_test.go`) | the three flags, `--decided`, the approver lookup inside the read-only handle, `ApprovedBy` removed from `:214` |
| `cmd/ailang-worldd/approver.go` | K8/K9 |
| `cmd/ailang-worldd/mint_request.go` | the `consume` map entry's wording |
| `.github/workflows/ci.yml` | the existing row 201 step widened (§6 M6) |
| `docs/QUICKSTART.md` | **new §15 "Approve with an enrolled SSH key"** *(attended — pending first verbatim run)*. It runs on the §0 scratch store `/tmp/se-world/world.db` with a **scratch key** made by `ssh-keygen -t ed25519 -N '' -f /tmp/se-world/demo-key < /dev/null`, never `~/.ssh`: `approver enroll`, `episode start --request`, `approve --sign-with`, `approve --list --decided`, `approver revoke`, then a refused `approve` with the revoked key |

**Not touched**: `world/*.ail`, `packages/`, `scripts/verify_ail.sh`, `tools/launchd/*`,
`host/daemon/daemon.go` (no route), `cmd/world-publish/*` (R4), `session revoke`, launch presets,
`--resume`. No open PRs at measurement time (V25).

## 8. Milestones

Every milestone ends with `go build ./... && go vet ./<pkg>/` plus its own tests. `go build` does
not compile test files (MEMORY: go-build-is-not-a-compile-fence), so `go vet` or `go test -run
'^$'` is the compile fence. Each sprint ends with the full gates:
- `./scripts/verify_ail.sh` (no `.ail` change, so exactly green);
- `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` (it includes the tracked-binary
  gate);
- `go build ./... && go test ./...`;
- `go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1`.

**Sprint 1 = M1–M3, "enrolment", about 1,900 lines.**
- **M1: `CR-TOOMANY` and the pure enrolment law.** About 120 production lines and 400 test lines,
  plus fixtures. CT1, AP0–AP5, RC4, and `approvaltest` with TS1/TS3.
  `go test ./cmd/ailang-worldd/ -run '^TestEpisodeStartRequestRefusesWhenSixteenArePending$' -count=1`;
  `go test ./host/approval/... -run 'Approver|Enrolment|CheckRevoke|TestSigner|ApprovalTest|Imports' -count=1`
- **M2: store schema 5.** About 230 + 420 lines. AS1–AS6, RC2, and the edited credentials tests.
  `go test ./host/store/ -run 'Credentials|Approver|OlderCredentials|MintRequest|TestProductionContextRoots' -count=1`
- **M3: `approver enroll|list|revoke`.** About 270 + 460 lines. CA1–CA8, help, keyguard and roots.
  `go test ./cmd/ailang-worldd/ -run '^(TestApprover|TestEveryVerbHasHelp|TestHelpIsNeverBlocked)' -count=1`
- **Sprint 1 fallback cut:** if the sprint plan's own estimate exceeds 2,000 lines, M3 moves to the
  head of sprint 2, and sprint 2's M6 becomes a third sprint.

**Sprint 2 = M4–M6, "approving by key", about 1,850 lines.**
- **M4: the signature law, the consume binding and authority.** About 170 + 530 lines. SA1–SA3,
  RC1, CS1–CS6, RC3, AU1–AU2, RC5.
  `go test ./host/approval/ ./host/store/ ./host/authority/ -run 'Signed|ApproveBySignature|Consume|ApprovedBy|DecidedMintRequests|MintForRequest|RequestMint' -count=1`
- **M5: the `approve` paths.** About 230 + 600 lines. CR1–CR13 and the edited help and roots
  tests.
  `go test ./cmd/ailang-worldd/ -run '^(TestApprove|TestProvisionEpisodeSkips|TestApproverRevokeLists|TestEveryVerbHasHelp)' -count=1`
  and `go test ./host/store/ -run '^TestProductionContextRoots$' -count=1`
- **M6: `--sign-with` with the real tool, the CI step, and QUICKSTART §15.** About 70 + 250 lines,
  20 YAML lines and 40 documentation lines. XS1, XS2, G2.
  `WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/... ./cmd/ailang-worldd/ -run '^(TestTestSignerAgreesWithSSHKeygen|TestSignWithSSHKeygenEndToEnd)$' -count=1 -v`
  must print `--- PASS` for every named test and subtest.
- **Sprint 2 fallback cut:** M6 becomes its own sprint.

**A1b exit criterion (attended, Mark; NOT run by the loop, D-WORLD-109).** Mark runs these on the
scratch store and his own machine:
1. `approver enroll --name <his-name> --format ssh --pubkey <his key>.pub`, with a key in the
   holder he chooses (Q5), typing `y` after comparing the fingerprint;
2. `episode start … --request`, then `approve <id> --sign-with <his key>.pub`, which mints with no
   World prompt;
3. `approve --list --decided`, which shows `ssh:<his-name>`;
4. `approver revoke <his-name>` (typed);
5. a second request, `approve --sign-with` with the same key, refused as revoked.

QUICKSTART §15 is executed verbatim in that run and its *attended — pending* tag removed. The
loop's own judgement uses only M1–M6's tests.

## 9. Residuals and risks

- **R1** Size. A0/A1a overran by 1.6× and 2.2× (V24). The fallback cuts in §8 are pre-named.
- **R2** Schema 5 makes a v4 store refuse bound, started and resumed mints from this build's CLI
  until `serve` restarts (`credentials.go:432`). `episode start` already forces that restart
  through its same-build check, so no new step is added. An older daemon refuses a v5 store, as
  before.
- **R3** Both sprints raise or use the schema, so sprint 1 alone already costs the one restart.
- **R4** Revocation is typed (parent R3). An emergency revoke needs the daemon's terminal. A key
  that cannot sign without its owner (Q5) narrows the window this matters in.
- **R5** The enrolment and revocation y/N reads are unbounded, like `session new`'s (a0a1 R7).
  There is no request expiry to bound them by.
- **R6** **SSHSIG proves possession of a private key, not the presence of a human.** It is not
  prevented, only stated at enrolment (§2.2) and in the help for `--sign-with`. `--sign-with
  <x>.pub` signs with the private file `x` beside it when no agent holds the key (V18), so naming a
  `.pub` does not mean an agent signed. `sk-*` touch-bit signatures remain refused (Q1).
- **R7** `approver list` and `approve --list --decided` are CLI only. The surface's decision list
  is A3.
- **R8** The approvers table has no cap. Rows only arrive through a typed `y`, so growth is human
  paced.

## 10. Open questions for Mark (each with a recommended default)

- **Q1** (carried from a0a1) `sk-*` security-key signatures: refuse them until an attended fixture
  from a real security key exists (a public key and a signature, no secret)? **Default: refuse**
  (unchanged; `sshsig.go:103-149` already does).
- **Q2** Approver names follow a strict grammar, `^[a-z0-9][a-z0-9._-]{0,47}$` (`mark-mac`,
  `mark-iphone`), not free text like the parent's "Mark's iPhone". They are printed on prompts and
  stored as `ssh:<name>`. **Default: strict.** A display label can come with A3.
- **Q3** Revoking marks the row and never deletes it. Names and keys are never reused, so a
  revoked key cannot be re-enrolled and `ssh:<name>` means one key forever. **Default: yes.**
- **Q4** Revoking an approver does not revoke the sessions it approved. It lists them with the
  `session revoke` line. **Default: list, do not auto-revoke.** Auto-revoke would turn one typed
  act into many silent ones.
- **Q5** For the exit run, which holder for Mark's key? **Default: an agent that confirms each use
  (Secretive on the Secure Enclave, or `ssh-add -c`), not an unencrypted file.** This is the only
  mitigation of R6 until `sk-*` or WebAuthn (A2).
- **Q6** Two sprints, enrolment first (§0), with the §8 fallback cuts. **Default: yes.**

## 11. What this is NOT doing

- No WebAuthn (A2), no surface, notification or Tauri change, no graphical decision list and no
  `/v1` route (A3), no `Human.Approve` packets (A4).
- No `sk-*` keys (Q1), and no RSA, P-384, P-521, DSA or certificate keys.
- No change to the typed `y` path's behaviour, `session new`/`mint`/`revoke`, launch presets,
  `--launch` or `--resume`. No change to publishing (R4 of the parent).
- No auto-revocation of sessions on approver revoke (Q4).
- No `.ail`, kernel, package or verify-gate change, and no new Go module.
- No real key, no `~/.ssh`, no real store, no mint in one and no publish, by the loop or by its
  tests (D-WORLD-109). Every test key comes from `t.TempDir()`, a fixed test seed, or a
  generator's `mktemp -d`.

## 12. Verification Log

Commands ran in the worktree at `3a96cf75` unless stated otherwise. Each empty result is paired
with a positive control in the same command.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and instruments | `git log --oneline -1`; `shasum -a 256 cmd/ailang-worldd/{approve,session_new,mint_request}.go host/store/{mint_requests,credentials}.go host/authority/mint.go host/approval/{sshsig,admit}.go \| awk '{print substr($1,1,16), $2}'`; `go version`; `ssh -V`; `$HOME/.pinned-ailang/ailang --version` | `3a96cf75 Merge pull request #536 …`; `a215b26ad7ed060a`, `2e9baa1c8c349750`, `d12b6929cd05fe8d`, `bfe6bdc86f22b957`, `b1d0e10d183e8d5c`, `49d1aed46251752d`, `a40850200782b3e4`, `e139596fcb089057`; `go1.26.6 darwin/arm64`; `OpenSSH_10.3p1, LibreSSL 3.3.6`; `AILANG v0.41.0` |
| V2 | Exit codes; refusals are `exitUsage` | `grep -n "exitOK\s*=\|exitUsage\s*=\|exitFatal\s*=" cmd/ailang-worldd/main.go`; `sed -n 403,560p cmd/ailang-worldd/session_new.go`; `sed -n 85,93p cmd/ailang-worldd/approve.go` | `72: exitOK = 0`, `73: exitUsage = 1`, `74: exitFatal = 2`; `awk` over `:403-560` for `return exit…`: `410`, `443`, `460`, `468` `return exitUsage` (no terminal, expired before the prompt, no answer within, aborted); `477`, `492`, `559` `return exitFatal`; `approve.go:92 return exitUsage` in `refuse` |
| V3 | `approve`'s flow | `cat -n cmd/ailang-worldd/approve.go` | `:54 func approveAdmission`; `:88 func runApprove`; `:117 store.MintRequestIDPattern`; `:130 context.WithTimeout(context.Background(), approveReadBudget)`; `:152 ro.MintRequest`; `:153 _ = ro.Close()`; `:169-177` the same-store refusal; `:180 prepareMint(fromRequest(mr))`; `:191 fetchDaemonHealth`; `:195 checkDaemon`; `:204-207 bytes.Equal(fb, row.Request)`; `:213-214 … ApprovedBy: "tty", Clock: env.now}}`; `:215 provisionEpisode` |
| V4 | `provisionEpisode` always opens the terminal | `grep -n "consume\|answerBy\|func provisionEpisode\|openTerminal\|readTermLine\|isYes\|context.Background\|MintForRequest" cmd/ailang-worldd/session_new.go` | `228: consume *authority.RequestMint`; `230: answerBy int64`; `403: func provisionEpisode`; `406: term, err := env.openTerminal()`; `441: if r.answerBy != 0 && r.answerBy-now <= 0`; `454-456` the bounded read; `464 answer = readTermLine(term)`; `466 if !isYes(answer)`; `472 context.WithTimeout(context.Background(), sessionNewBudget)`; `506-507 if r.consume != nil { … authority.MintForRequest(…)` |
| V5 | `sessionEnv` seams | `cat -n cmd/ailang-worldd/session.go \| sed -n 1,80p` | `:25 type sessionEnv struct {`, fields `openTerminal`, `now`, `random`, `after`; `:49 return os.OpenFile("/dev/tty", os.O_RDWR, 0)`; `:53 after: time.After`; no `sign` field |
| V6 | Credentials schema 4, the table sets, the upgrade, the bound-mint refusal | `grep -n "credentialsSchemaVersion\|case known\|func " host/store/credentials.go`; `sed -n 36,100p` and `sed -n 296,392p` and `sed -n 428,436p` of the same file; `cat -n host/store/credentials_schema.sql` | `:42 const credentialsSchemaVersion = 4`; `:71 func credentialsTables` with `case version >= 4: return []string{"episode_starts", mintRequestsTable, credentialsTable}`; `:318 case known && slices.Equal(tables, credentialsTables(version))`; `:340 func upgradeCredentials` with `if version < 4 { statements = append(statements, mintRequestsSQL) }`; `:432 if (row.Bound() \|\| start != nil \|\| resume != nil) && c.version < credentialsSchemaVersion {`; schema `:44-56 CREATE TABLE IF NOT EXISTS mint_requests (…)` |
| V7 | The mint-request store | `grep -n "^func \|^type RequestConsume\|^const mintRequestsSchemaVersion\|^const MaxPending\|ErrTooManyPendingRequests =\|ErrIncompleteConsume \|need schema %d" host/store/mint_requests.go`; `sed -n 140,312p` same file | `:36 mintRequestsSchemaVersion = 4`; `:39 MaxPendingMintRequests = 16`; `:49 ErrTooManyPendingRequests = errors.New("store: 16 mint requests are already pending; …")`; `:56 ErrIncompleteConsume … "needs an id, a digest, approved_by, a clock and a check"`; `:66 … "mint requests need schema %d: "+`; `:88 PutMintRequest`; `:142 MintRequest`; `:162 PendingMintRequests`; `:172 … WHERE consumed_at IS NULL AND expires_at > ? …`; `:234 type RequestConsume struct` = `ID, Digest, ApprovedBy string; Clock func() int64; Check func(consumedAt int64) error`; `:243 MintSessionForRequest`; `:265 UPDATE mint_requests SET consumed_at = ?, approved_by = ?, credential_id = ? WHERE request_id = ? AND digest = ? AND consumed_at IS NULL AND created_at <= ? AND expires_at > ?` (no `challenge`); `:287 classifyConsume` |
| V8 | The mint core | `grep -n "RequestMint\|func MintForRequest\|func mintRecorded" host/authority/*.go \| grep -v _test`; `sed -n 74,135p host/authority/mint.go` | `mint.go:60 type RequestMint struct` = `Request approval.MintRequest; RequestCreatedAt int64; Consume store.RequestConsume`; `:70 func MintForRequest`; `:74 func mintRecorded`; `:114 case rm != nil:` … `consume.Check = func(consumedAt int64) error { return approval.MintMatches(…) }`; `err = st.MintSessionForRequest(ctx, row, consume)`; `ApprovedBy` not read |
| V9 | The pure approval API | `grep -n "^func \|^type \|^const\|^var" host/approval/{admit,request,sshsig}.go`; `sed -n 95,200p host/approval/sshsig.go` | `admit.go:44 func Admit`, `:92 func MintMatches`; `request.go:211 func SignedMessage`; `sshsig.go:19 const Namespace = "world-approval@sunholo"`; `:46 type SSHPublicKey struct { Type string; Blob []byte }`; `:103 func parseKeyBlob` with `case "sk-ssh-ed25519@openssh.com", "sk-ecdsa-sha2-nistp256@openssh.com": … refused until an attended hardware fixture exists (row 201 Q1)` and `default: … is not ssh-ed25519 or ecdsa-sha2-nistp256`; `:151 Fingerprint`; `:157 ParseAuthorizedKey` returns `SSHPublicKey{Type: km.typ, Blob: blob}` (no comment); `:182 ParseSSHSIG`; `:262 VerifySSHSIG` (`ErrKeyMismatch` when `!bytes.Equal(s.PublicKey.Blob, key.Blob)`); `:353 VerifyApproval` |
| V10 | `CR-TOOMANY` survives today | `sed -n 101,153p cmd/ailang-worldd/mint_request.go`; then with `AILANG_BIN=$HOME/.pinned-ailang/ailang`: `cp mint_request.go /tmp/mr_backup_iter375.go`; `go test ./cmd/ailang-worldd/ -run 'Request\|Approve\|EpisodeStart' -count=1`; `sed -i '' 's/… \|\| errors.Is(err, store.ErrTooManyPendingRequests)/… \|\| false \&\& errors.Is(…)/'`; `grep -n "false && errors.Is"`; the same `go test`; restore; `git status --short`; `git diff --stat` | `:101 func createMintRequest`; `:142 if errors.As(err, &older) \|\| errors.Is(err, store.ErrTooManyPendingRequests) {` `:144 return exitUsage`; before: `ok … cmd/ailang-worldd 3.281s`; `142: … \|\| false && errors.Is(err, store.ErrTooManyPendingRequests) {`; after: `ok … 3.353s` (**the mutant survives**); after the restore both `git status` and `git diff --stat` are empty |
| V11 | No CLI test pins it (with controls) | `grep -rn "ErrTooManyPendingRequests\|already pending" cmd/ailang-worldd/*_test.go \| wc -l`; control `grep -rln "OlderCredentialsStoreError\|is schema 3\|serve. with this build" cmd/ailang-worldd/*_test.go`; `grep -rn "ErrTooManyPendingRequests" host/store/*_test.go` | `0`; control: `cmd/ailang-worldd/presplit_unix_test.go`; `host/store/mint_requests_test.go:257 … !errors.Is(err, ErrTooManyPendingRequests)` |
| V12 | The field map's `consume` entry; `mintRequest`'s fields | `sed -n 180,230p cmd/ailang-worldd/mint_request.go`; `sed -n 205,232p cmd/ailang-worldd/session_new.go` | `:223 "consume": "excluded: set by approve from the stored row (id, digest, approved_by tty, env.now)"`; `mintRequest` ends `consume *authority.RequestMint` / `answerBy int64` |
| V13 | Context-roots inventory | `grep -n "approve\|mint_request\|session_new.go\|session.go" host/store/context_roots_test.go` | `:19 "cmd/ailang-worldd/mint_request.go\|createMintRequest\|Background": 1`; `:20 "cmd/ailang-worldd/approve.go\|runApprove\|Background": 1`; `:28-36` the session verbs; no `approver.go` entry |
| V14 | Verb and help pins | `grep -n "approve\|\"session\"" cmd/ailang-worldd/help_test.go`; `grep -n "approve\|approver" cmd/ailang-worldd/main.go`; `sed -n 655,666p cmd/ailang-worldd/main.go`; `grep -n "case \"approve\":" cmd/ailang-worldd/keyguard.go` | `help_test.go:14 {"session"}, …`, `:15 {"approve"},`, `:44 "approve <request-id>", "approve --list",`; `main.go:124-125` the two approve usage lines; `:442 case "approve":`; `:664 … "approve": true, …` in `ownHelpVerbs`; `keyguard.go:77 case "approve":`; no `approver` anywhere |
| V15 | CI has `ssh-keygen`; the row 201 step passed | `sed -n 137,155p .github/workflows/ci.yml`; `gh run list --branch dev --workflow ci.yml --limit 1 --json databaseId,headSha,conclusion`; `gh run view 42753484556 --log \| grep "SSHSIG cross-check" \| grep -E "OpenSSH_\|PASS: TestSSHSIG\|/usr/bin/ssh-keygen"` | `:141 - name: Row 201 SSHSIG cross-check with real ssh-keygen, verbose` … `:153 for t in TestSSHSIGCrossCheckWithSSHKeygen …/ed25519 …/ecdsa-p256`; `[{"conclusion":"success","databaseId":42753484556,"headSha":"3a96cf755f3e119de72d68e098194e786da04e54"}]`; `/usr/bin/ssh-keygen`; `OpenSSH_9.6p1 Ubuntu-3ubuntu13.20, OpenSSL 3.0.13 30 Jan 2024`; `--- PASS: TestSSHSIGCrossCheckWithSSHKeygen (0.10s)`, `--- PASS: …/ed25519 (0.06s)`, `--- PASS: …/ecdsa-p256 (0.05s)` |
| V16 | Fixtures, the NUL rule, the gate | `ls host/approval/testdata/sshsig`; `cat host/approval/testdata/sshsig/generate.sh`; `cat .gitattributes`; `grep -n -i "binary\|NUL\|tracked" scripts/verify_go.sh` | `SOURCE.txt ed25519.pub ed25519.sig generate.sh git-ns-ed25519.sig message.hex p256.pub p256.sig` (no `message.bin`); the generator: "The 90-byte message is kept as message.hex (text), never as a binary file: the tracked-binary gate (scripts/tracked_binary_gate.sh) rejects any committed blob that holds a NUL byte"; `host/approval/testdata/sshsig/* -text`; `verify_go.sh:228 echo "── tracked-binary hygiene gate"`, `:230 … whole-blob NUL detector`, `:239 bash scripts/test_tracked_binary_gate.sh`, `:240 bash scripts/tracked_binary_gate.sh --allow-png …` |
| V17 | `ssh-keygen` keys, fingerprints, stdin signing, `check-novalidate` (temp keys only) | in `T=$(mktemp -d /tmp/world-iter375-designer.XXXXXX)`, `HOME=$T`, `SSH_AUTH_SOCK` unset: `ssh-keygen -q -t {ed25519, ecdsa -b 256, rsa -b 3072, ecdsa -b 384} -N "" -C probe-* -f $T/{ed,p256,rsa,p384} < /dev/null`; `ssh-keygen -l -f` each; a 90-byte `msg.bin`; `ssh-keygen -Y sign -n world-approval@sunholo -f $T/ed < msg.bin > ed.sig`; `ssh-keygen -Y check-novalidate -n world-approval@sunholo -s ed.sig < msg.bin`, then with `-n git`, then on a tampered message; `ls -la $T/.ssh` | `rc=0` ×4; `256 SHA256:mnLou9rs… probe-ed (ED25519)`, `256 … (ECDSA)`, `3072 … (RSA)`, `384 … (ECDSA)`; `fp len=50`; `90`; `Signing data on standard input`, `sign stdin rc=0 bytes=318`, first line `-----BEGIN SSH SIGNATURE-----`; `Good "world-approval@sunholo" signature with ED25519 key SHA256:mnLou9rs…` rc 0; `Couldn't verify signature: namespace does not match` rc 255; `Signature verification failed: incorrect signature` rc 255; `$T/.ssh: No such file or directory` |
| V18 | `-f <x>.pub` falls back to the private file; an agent signs; Ed25519 is deterministic | same `$T`: `ssh-keygen -Y sign … -f $T/ed.pub < msg.bin > x.sig` (private `ed` beside it, no agent); `-f $T/pubonly/only.pub` (no private file, no agent); `ssh-agent -a $T/agent.sock`, `ssh-add $T/ed`, `ssh-add -l`, `-f only.pub` again; `cmp w.sig ed.sig`; `kill $SSH_AGENT_PID` | `rc=0 bytes=318`, and `x.sig` verifies `Good … ED25519 key SHA256:mnLou9rs…`; `No private key found for public key "…/pubonly/only.pub"` `rc=255 bytes=0`; `Identity added: …/ed (probe-ed)`; `rc=0 bytes=318`, `Good …` rc 0; `equal`; `agent gone` |
| V19 | A passphrase key with no terminal fails fast | same `$T`: `ssh-keygen -q -t ed25519 -N "probe-pass-123" -f $T/pw < /dev/null`; `ssh-keygen -Y sign … -f $T/pw < msg.bin` in the background, polled 10 s with `date +%s` | `rc=255 bytes=0 wall=0s`; `Enter passphrase for "…/pw": Load key "…/pw": incorrect passphrase supplied to decrypt private key` |
| V20 | `sk-*` keys need hardware | same `$T`: `ssh-keygen -q -t ed25519-sk -N "" -f $T/sk < /dev/null` | `No FIDO SecurityKeyProvider specified`, `Key enrollment failed: invalid format`, rc 255 |
| V21 | `-Y verify` with an allowed-signers file | same `$T`: `allowed` = `approver-a namespaces="world-approval@sunholo" <ed.pub>`; `ssh-keygen -Y verify -f allowed -I approver-a …`; then `-I approver-b` | `Good "world-approval@sunholo" signature for approver-a with ED25519 key SHA256:mnLou9rs…` rc 0; `Could not verify signature.` rc 255 |
| V22 | The pinned AILANG has no Ed25519, ECDSA or SHA-512 | `$HOME/.pinned-ailang/ailang builtins list \| grep -iE "_crypto_"`; `… \| grep -ciE 'ed25519\|ecdsa\|sha512'`; control `… \| grep -ciE 'string'` | `_crypto_constanttimeequal`, `_crypto_hmacsha256`, `_crypto_rsa_verify_pkcs1v15`, `_crypto_sha256bytes`, `_crypto_sha256hex` (all `[pure] std/crypto`); `0`; control `38` |
| V23 | Nothing of A1b exists; the store does not import approval (with controls) | `grep -rnE "approvers\|\"approver\"\|sign-with\|write-message\|ssh:" --include="*.go" host cmd \| grep -v _test.go`; control `grep -rn '"approve"' cmd/ailang-worldd/main.go`; `grep -rln '"ssh-keygen"' --include="*.go" host cmd \| grep -v _test`; control `grep -rln '"ssh-keygen"' --include="*_test.go" host cmd`; `grep -rln 'host/approval"' --include='*.go' host/store \| wc -l`; control the same over `host/authority` | only `episode_finish_dispatch.go:141` (an `ssh://git@github.com/` regexp, unrelated); control `main.go:442 case "approve":`, `:664`; no production file; control `host/approval/sshsig_crosscheck_test.go`; `0`; control `2` |
| V24 | The A0/A1a sprint sizes | `git diff --numstat dba0e8b2^1 dba0e8b2 \| awk` (code vs `design_docs\|sprint_`); the same for `ca8a4c17` (insertions + deletions) | `sprint1 code+tests+fixtures=2291 plan/docs=2549`; `sprint2 code changed=3473 plan/docs=2263` |
| V25 | No open PRs | `gh pr list --state open --json number,headRefName,title --limit 20` | `[]` |
| V26 | `approved_by` is written as `tty` only, and read only by the replay refusal (with control) | `grep -rn "approved_by\|ApprovedBy" --include="*.go" host cmd \| grep -v _test`; `grep -rln "decision list\|decisionList\|DecisionList" --include="*.go" host cmd`; control `grep -n "decision list" design_docs/planned/w-world-approvals-by-key.md` | writers: `approve.go:214 … ApprovedBy: "tty"` (the only literal); readers: `approve.go:71` (the "already used … by %s" text), `mint_request.go:223` (the map comment), the store's own SQL; decision list in Go: no match; control: `w-world-approvals-by-key.md:58 … records **which key approved it** in the episode's record and the decision list` |
| V27 | `session revoke`'s usage (for K8's printed line) | `grep -n "session revoke" cmd/ailang-worldd/main.go` | `:141 ailang-worldd session revoke --db <path> [--credentials-db <path>] <credential_id-hash> \| --episode <ep>` |
| V28 | The tests pinned to schema 4 (the edit sites) | `grep -rn "user_version = 4\|want 4\|schema 4" --include="*_test.go" host/store cmd \| head`; `sed -n 53,65p host/store/credentials_test.go` | `mint_requests_test.go:202 … "need schema 4"`, `:125 (c) The result is schema 4 …`; `credentials_binding_test.go:106 … want %d", i, v, credentialsSchemaVersion`; `credentials_test.go:56 requireExactTableNames(…, []string{"episode_starts", "mint_requests", "session_credentials"})` |
| V29 | Only serve upgrades; the CLI writer never upgrades; the log text | `sed -n 188,216p host/store/credentials.go`; `sed -n 150,168p cmd/ailang-worldd/session.go`; `sed -n 55,85p host/daemon/credentials.go` | `OpenCredentials … An existing version-1 store is opened as it is, never upgraded here`; `openCredentialsForSession` calls `store.OpenCredentials(creds)`; `:59 func ensureCredentialsStore`; `:76 "ailang-worldd: upgraded the credentials store %s to schema %d (every existing credential unchanged)\n"` |
| V30 | `applicationTables` lists tables only | `sed -n "/^func applicationTables/,/^}/p" host/store/credentials.go \| grep -n "sqlite_master"` | `SELECT name FROM sqlite_master WHERE type = 'table' AND name NOT LIKE 'sqlite\_%' ESCAPE '\' ORDER BY name` |
| V31 | `host/approval`'s import guard | `sed -n 1,78p host/approval/imports_test.go` | `approvalForbidden` = `os`, `os/exec`, `os/signal`, `net`, `net/http`, `syscall`, `io/ioutil`, `crypto/rand`, `math/rand`, `math/rand/v2`, `unsafe`; a non-stdlib path (a dot in the first segment) is also refused; `TestApprovalImportsOnlyTheStandardLibrary` |
| V32 | Probes cleaned up; no git writes | `rm -rf "$T"`; `ls -d /tmp/world-iter375-designer*`; `pgrep -fl "ssh-agent -a /tmp/world-iter375" \|\| echo …`; `git status --short` | `No such file or directory`; `no probe agent running`; only this doc listed |
| V33 | The designer's shell has no controlling terminal | `( : < /dev/tty ) 2>&1 && echo openable \|\| echo "not openable"` | `/dev/tty: Device not configured`, `tty: not openable` |
| V34 | The authority tests that pass `ApprovedBy: "tty"` | `grep -n "ApprovedBy\|tty" host/authority/mint_request_test.go \| head -8` | `:56 return store.RequestConsume{ID: row.RequestID, Digest: row.Digest, ApprovedBy: "tty", Clock: …}`; `:78 … got.ApprovedBy != "tty" …` |

## 13. Quorum verification log

*(Left empty for the controller.)*
