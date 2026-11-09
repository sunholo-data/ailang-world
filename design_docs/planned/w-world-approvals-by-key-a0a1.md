# w-world-approvals-by-key-a0a1 — the mint request (A0) and the SSHSIG verifier (A1a)

- **Status**: Planned (designed iteration 373, **r3** = controller carve-out applying the round-2 reviewers' fixes, §13; **r2**: r1 was BLOCKED 2 of 3 at quorum round 1. r2
  fixes the two objections, marked **[r2]**: the field map and widened mint check in §2.1a and
  §2.3, and the bounded prompt in §2.6a. See §13). Queue row 201, `w-world-approvals-by-key`.
  **[r2] Sprint boundary: this sprint builds M1–M3. M4–M6 are the next sprint, from this same doc
  (§0, §8).**
  This doc covers parent milestones **A0** and the **verifier half of A1**. Parent:
  [`w-world-approvals-by-key.md`](w-world-approvals-by-key.md); its direction and rulings are
  settled and not reopened here.
- **Authority**: **D-WORLD-89** (Mark, attended 2026-10-21, "all defaults": R1 a verified signature
  from an enrolled key may approve a mint, as well as the typed `y`; R2 the formats are WebAuthn
  and OpenSSH SSHSIG; R3 enrolment and revocation stay typed-only; R4 publishing stays typed-only;
  R5 ordering). **D-WORLD-109 = A** (Mark, attended 2026-11-06): the loop builds and judges row 201
  against **fakes and test keys only**. It never enrols a real key, never mints in a real store and
  never publishes. Every key in this sprint is generated inside a test (`t.TempDir()`) or by the
  checked-in fixture generator, whose private keys live only in a `mktemp -d` directory that it
  deletes (§2.8). Every store is a test store under `t.TempDir()`. Mark's own enrolment and his
  first key-approval run stay attended (A1b).
- **Item**: the first two steps of row 201. It adds a stored, single-use **mint request** that the
  typed `y` approves, and a pure, standard-library **SSHSIG verifier** tested on real
  `ssh-keygen` signatures. It changes no `/v1` route and no `.ail` code. It raises the credentials
  store's schema from 3 to 4 (one new table).
- **Designer lane**: `claude:claude-opus-5-5`.
- **Estimated**: 6 milestones in two sprints, about 3,000 lines in all once r2's additions are
  counted (§8). **[r2]** **Sprint 1 = M1–M3**, about 1,450 lines: the pure core, the SSHSIG
  verifier, and the cross-check with its CI step. **Sprint 2 = M4–M6**, about 1,600 lines: store
  schema 4, `episode start --request`, `approve`. Each sprint is under the 2,000-line ceiling.
- **Measurement base**: `693f704e` (= origin/dev), worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter373-row201`, branch
  `sprint/world-iter373-row201`. sha256 prefixes: `session.go` `3566d5377c9fc01a`,
  `session_new.go` `ef5163a5dc6e7aa5`, `episode.go` `5588c8659f51669b`, `host/store/credentials.go`
  `3b4d630f89676737`, `host/authority/mint.go` `1dfdf1080000889f` (V1).
- **Instruments**: `go1.26.6 darwin/arm64`; `/usr/bin/ssh-keygen` from `OpenSSH_10.3p1, LibreSSL
  3.3.6`; pinned `AILANG v0.41.0` at `$HOME/.pinned-ailang/ailang`. Scratch work went to
  `/tmp/world-iter373/designer-*` with throwaway keys only, and those directories were deleted at
  the end (V27). No Linux host was used, so every statement about the CI runner is **UNVERIFIED**
  and the design makes the CI step fail loudly rather than assume (§2.8, V14).

---

## 0. Scope, and what moved out

The brief asked for A0 and A1 in one sprint of 1,200–2,000 lines, or else a cut to A0 plus the
SSHSIG verifier. A0 alone is a new pure package, a schema bump and two CLI paths, and A1 adds a
parser for another tool's format, a second table and three more verbs. My estimate for both is about
3,500 lines with tests (this cut is about 2,800, §8; A1b's enrolment, its table and the signature
paths are the other ~700), so **this doc makes the cut**.

**In this sprint:**
- **A0**: the mint request (canonical bytes, SHA-256 digest, 32-byte one-time challenge, 10-minute
  expiry, single consumption, storage), `episode start --request` to create one, and
  `ailang-worldd approve <request-id>` to mint from it, approved by the typed `y` at `/dev/tty`.
- **A1a**: a pure SSHSIG parser and verifier (Ed25519 and ECDSA P-256), a function that binds a
  signature to a stored request, real `ssh-keygen` fixtures, a cross-check against `ssh-keygen -Y
  verify`, and a CI step that cannot pass while skipping.

**Moved out to A1b** (the next row-201 sprint), and why:
- `ailang-worldd approver enroll|list|revoke` and the `approvers` table. Enrolment is the root of
  trust (R3). It deserves its own review, not the tail of a sprint that is already over its line
  budget.
- `approve <id> --signature <file>` and `approve <id> --sign-with <key>`: the signature path to a
  mint. It needs an enrolled key to verify against, so it lands with enrolment.
- Recording which key approved a session in the episode record and the decision list. In A0
  `approved_by` is always `tty` (§2.4).
- Revoked-key and enrolment threats (§4 marks them A1b).

**Out of row 201's A0/A1 entirely** (later parent milestones; only their seams are named, §2.9):
A2 WebAuthn, A3 the surface and Tauri flow, A4 `Human.Approve` decision packets.

**[r2] Sprint boundary.** At about 3,000 lines the cut is still too big for one sprint. **This
sprint builds M1–M3**: the pure request core, the SSHSIG verifier with its real fixtures, and the
`ssh-keygen` cross-check with its CI step. They touch only `host/approval/`, its testdata,
`.gitattributes` and `ci.yml`, and change no behaviour. **M4–M6** (store schema 4,
`episode start --request`, `approve`) are **the next sprint**, built from this same doc without
re-design. Nothing in M1–M3 depends on M4–M6.

A0 on its own does not improve the operator's experience. It is still one typed `y` at a terminal,
now split into two commands. It is the base the signature paths need: a request whose exact bytes
can be signed, stored and used once.

## 1. What the code does today (measured)

1. **One fence-and-mint path.** `session new` and `episode start` both mint through
   `provisionEpisode` (`cmd/ailang-worldd/session_new.go:396`). Its own comment says that no other
   code creates a credential (`:24-26`). It opens the terminal first (`:399`, via
   `env.openTerminal`), reads the clock (`:408`), prints a y/N naming the worktree, every grant, the
   profile binding and the TTL (`:436-437`), takes the credentials writer lock under a 75 s root
   (`:444`), adds the worktree, mints (`:475-485`) and rolls back on failure (`:503-518`) (V2).
2. **The seam is `sessionEnv`** (`cmd/ailang-worldd/session.go:24-32`): `openTerminal` and `now`.
   The production value opens `/dev/tty` `O_RDWR` (`:40-43`). It has **no randomness seam** (V3).
3. **Plain `episode start`** (`cmd/ailang-worldd/episode.go:185`) validates everything with
   `prepareMint` (`:305`), checks the serving daemon's health and **build identity**
   (`checkDaemon`, `:431`, which requires `daemon.SameBuild`), then calls `provisionEpisode`
   (`:343`). It does not set `r.start`, so it mints through `sessionMint` = `authority.MintBound`.
   Only the launch-preset path sets `req.start` (`episode_preset.go:240`) (V4).
4. **The credentials store** is schema version 3 (`host/store/credentials.go:40`). The open check
   accepts exactly the known table sets per version (`:304`). Only `serve` upgrades a store
   (`ensureCredentialsStore`, `host/daemon/credentials.go:59`), and the CLI verbs refuse a bound
   mint into an older store (`session_new.go:372`, `episode_preset.go:235`). Every store method
   requires a context deadline (`requireDeadline`, `credentials.go:403`) (V5, V6).
5. **The mint core** is `authority.mintRecorded` (`host/authority/mint.go:57`). It makes the token,
   clamps the grant expiry to the session and calls the store's one transaction (V7).
6. **Nothing in the repo does approvals by key yet.** There is no `sshsig`, `world-approval`,
   `approver` or mint-request code. `golang.org/x/crypto` is not a dependency and no Go file
   imports it (V8). The daemon's dependency allowlist test (`host/daemon/daemon_test.go:899`) holds
   `host/daemon/...` and `cmd/ailang-worldd/...` to the standard library plus a closed module list
   (`:775`, `:793`) (V9).
7. **Context roots are inventoried.** `TestProductionContextRoots` pins every production
   `context.Background()` root by file and function (`host/store/context_roots_test.go:26-37`) (V10).
8. **CI** runs `go test` on `ubuntu-24.04` and `ubuntu-latest`. No workflow mentions `ssh-keygen`.
   The house pattern for a test that can skip is a dedicated verbose step that fails unless every
   named test prints `--- PASS:` (`.github/workflows/ci.yml:238-258`) (V14).
9. **Revocation is not fenced today.** `session revoke` says "No terminal fence is required"
   (`session.go:340`). Parent R3 makes *approver* revocation typed; that is A1b and does not
   change `session revoke` (V2).

## 2. Design

### 2.1 The mint request

A new package `host/approval` holds the pure core. It imports only the standard library: it does
not import `host/broker`, `host/daemon` or `host/store`, so A3's daemon can import it later
without an import cycle.

```go
const RequestSchema = "world/mint-request/v1"

type MintRequest struct {
	Schema        string     `json:"schema"`         // RequestSchema
	Episode       string     `json:"episode"`        // ^[a-z0-9][a-z0-9-]{0,63}$
	WorldDB       string     `json:"world_db"`       // [r2] canonical --db path
	WorkspaceRoot string     `json:"workspace_root"` // [r2] absolute --workspace-root
	Worktree      string     `json:"worktree"`       // absolute, clean: <workspace-root>/<episode>
	Reuse         bool       `json:"reuse"`          // the worktree exists already
	Degraded      bool       `json:"degraded"`       // non-git (always false on this path, §2.1a)
	Repo          string     `json:"repo"`           // [r2] absolute (filepath.Abs at create), or ""
	Git           string     `json:"git"`            // [r2] the resolved git binary; "" when Reuse
	Branch        string     `json:"branch"`
	Grants        []Grant    `json:"grants"`         // in order, at least one
	TTLSeconds    int64      `json:"ttl_s"`          // ≥ 0
	Profile       ProfileRef `json:"profile"`        // both "" when unbound
	ProfilesDir   string     `json:"profiles_dir"`   // [r2] absolute --exec-profiles-dir
	CredentialsDB string     `json:"credentials_db"` // canonical path the mint writes
	SessionFile   string     `json:"session_file"`   // absolute; the token's only file
	MCPConfig     string     `json:"mcp_config"`     // absolute for claude, "" for codex
	Agent         string     `json:"agent"`          // claude | codex
	Model         string     `json:"model"`
	Task          string     `json:"task"`           // absolute or ""
	Daemon        DaemonRef  `json:"daemon"`
}
type Grant struct {
	Effect string `json:"effect"`
	Scope  string `json:"scope"`
	Budget int64  `json:"budget"`
}
type ProfileRef struct {
	Project string `json:"project"`
	Digest  string `json:"digest"`
}
type DaemonRef struct { // the requesting daemon's build identity (parent §2.2) and address
	Addr             string `json:"addr"`
	VCSRevision      string `json:"vcs_revision"`
	VCSModified      bool   `json:"vcs_modified"`
	ExecutableSHA256 string `json:"executable_sha256"`
	GoVersion        string `json:"go_version"`
}
```

- **Everything the approval executes is inside the request.** This covers the parent's list
  (episode, worktree, every grant with scope and budget, TTL, exec profile digest, daemon build
  identity). It also covers every value `approve` acts on: where the token is written, where the
  MCP config goes, which repo and branch the worktree comes from, and which credentials store is
  written. If the token path were outside the digest, a substituted request could send the token
  somewhere an agent can read. **Rule:** `approve` builds its mint only from the decoded request
  bytes. No flag or default can change what is minted (§5). **[r2]** §2.1a maps every field of
  the mint path to this struct or to a named exclusion, and two tests keep that map complete.
- Grants carry no `ExpiresAt`. The mint sets it to `now + TTLSeconds`, as `provisionEpisode` does
  today (`:410-412`), so the digest does not depend on the time.
- `Validate()` refuses each of these, as its own named branch: the wrong schema; a bad episode
  grammar; a path that is not absolute and clean (**[r2]** including `WorldDB`, `WorkspaceRoot`,
  `ProfilesDir`, a non-empty `Repo`, and `Git`, which is non-empty iff `!Reuse`); no grants; an empty effect or scope; a negative
  budget or TTL; a half profile; `MCPConfig` set for codex or empty for claude; an agent outside
  `{claude, codex}`; a string that is not valid UTF-8 or holds NUL.

**Canonical bytes.** `Canonical(r)` is the compact JSON of `MintRequest` in struct field order,
produced by an encoder with `SetEscapeHTML(false)` and with no trailing newline. `Decode(b)` is
strict. It refuses unknown fields and trailing data, and it re-encodes the result: if
`Canonical(decoded) != b` it refuses with `ErrNotCanonical`. So exactly one byte string stands for
each request, and a reformatted, reordered or padded copy is refused.

**Digest.** `Digest(r) = SHA-256(Canonical(r))`. It is stored and shown as 64 lowercase hex.

**Why JSON and not a binary encoding.** The bytes are stored, shown in refusals and read back by
A3's surface. JSON is readable there. The re-encode check makes it canonical. Go's encoder writes
struct fields in declaration order with deterministic escaping, and test RQ1 pins the bytes against
a golden copy.

#### 2.1a [r2] Every field of the mint path, mapped

Round 1 found a real gap. "Everything the approval executes is inside the request" was asserted,
never checked against the structs, and the §2.6 byte-equality check cannot see a field that never
entered `MintRequest`. Both structs at `693f704e` are dumped in V28 and their setters in V29. The
table maps **every** field to a `MintRequest` field or to a named exclusion. `approve` rebuilds its
`mintRequest` only through this map (`fromRequest`, §2.6 step 3).

| `mintOptions` / `mintRequest` field | Disposition |
|---|---|
| `verb` | **excluded, display only**: the message prefix. `approve` sets `"approve"` |
| `episode` | `Episode` |
| `db` | `WorldDB` (canonical). `approve`'s `--db`/`--credentials-db` only *locate* the request. **[r3, kimi verbatim]** approve refuses unless canonical(--db) resolves to the request's WorldDB and the resolved credentials path (--credentials-db or its derivation) equals CredentialsDB; it passes WorldDB on. `db` reaches the mint only through `MigrateCredentials`, a no-op on a schema-4 store, and through the same-file check |
| `credsDB` | `CredentialsDB` (`approve` passes the request's canonical path) |
| `wsRoot` | `WorkspaceRoot` |
| `repo` | `Repo`, made absolute at create, so `approve` in another directory gets the same repository |
| `branch` | `Branch` |
| `allowNonGit` | **fixed `false`**: `episode start` never sets it (`episode.go:305-307`, V29). `approve` passes `false` |
| `grants` | `Grants` (effect, scope, budget, in order). `ExpiresAt` is **fixed 0** in the request and set at the mint to `now + ttl` (`session_new.go:409-412`). **[r3]** `broker.Capability` is exactly `{Effect, Scope, ExpiresAt, Budget}` with no json tags, so `session_credentials.grants_json` carries the keys `Effect`, `Scope`, `ExpiresAt`, `Budget` (`json.Marshal(grants)`, `host/authority/mint.go:83`) (V34). Element dispositions, enforced by CL14: `grants[].Effect -> MintRequest.grants[].effect`, `grants[].Scope -> MintRequest.grants[].scope`, `grants[].Budget -> MintRequest.grants[].budget`, `grants[].ExpiresAt -> fixed: 0, set at mint` |
| `ttl` | `TTLSeconds` |
| `outPath` | `SessionFile` |
| `outDirCreate` | **fixed `true`**, `episode start`'s value (`episode.go:307`) |
| `profile` | `Profile.Project` |
| `profilesDir` | `ProfilesDir` |
| `credPath` | derived from `WorldDB` and `CredentialsDB`; re-derived and compared through `CredentialsDB` |
| `dir` | `Worktree` (derived; re-derived and compared) |
| `reuse`, `degraded` | `Reuse`, `Degraded` (derived; re-derived and compared) |
| `git` | `Git`: the `exec.LookPath("git")` result, which depends on `PATH`. Re-derived and compared, so an `approve` whose `PATH` finds another `git` is refused |
| `prof` | `Profile.Digest` (+ `Project`). The digest pins the profile's bytes, and `checkDaemon` compares it with the daemon's loaded copy |
| `profFile` | derived as `ProfilesDir/Project.json`; covered by `ProfilesDir` + `Project` |
| `afterToken` | set by `approve` exactly as `episode.go:339` sets it, from `MCPConfig` and `Daemon.Addr`; `nil` iff `MCPConfig == ""` |
| `resume`, `beforeToken`, `resumeRollback` | **must be nil**: set only by `episode_resume.go:190-201` (V29). `--request` refuses `--resume` (CL2) |
| `start` | **must be nil**: set only by `episode_preset.go:240`. `--request` refuses launch presets (CL2) |
| `startPoint`, `existingBranch` | **must be `""` / `false`**: set only by `episode_preset.go:239`. So the plain path never reaches the `existingBranch` action text (`session_new.go:433-434`) or `AddExisting` (`:457-458`) |
| `presetPrompt` | **excluded, display only**: `approve` sets it to `Render(request)`, which is derived from the request |
| *(outside `mintRequest`)* `agent`, `model`, `task`, `mcpFile`, `hostPort`, the daemon's `Build` | `Agent`, `Model`, `Task`, `MCPConfig`, `Daemon.Addr`, `Daemon.*` |

The map is code, not prose. `cmd/ailang-worldd/mint_request.go` declares
`var mintRequestFieldMap = map[string]string{...}`, keyed by the Go field name, with the value
either `"MintRequest.<path>"` or `"excluded: <reason>"` / `"fixed: <value>"` / `"nil: <reason>"`.
Two tests enforce it:
- **CL14 `TestMintRequestFieldMapIsComplete`** reflects over `mintOptions` and `mintRequest`,
  flattening the embedded struct. It fails if any field is missing from the map, if the map names a
  field that does not exist, or if a `MintRequest.<path>` target does not exist (reflecting over
  `MintRequest`). It also runs the reverse check: every `MintRequest` leaf is a target, or is in
  the request-only list (`Agent`, `Model`, `Task`, `MCPConfig`, `Daemon.*`). **A future field on
  any of the three structs is red until someone decides its disposition.** **[r3, kimi verbatim]**
  `TestMintRequestFieldMapIsComplete` additionally reflects over broker.Capability and fails on any
  Capability field lacking a disposition in mintRequestFieldMap (with the map gaining entries like
  `grants[].effect -> MintRequest.grants[].effect`, `grants[].expires_at -> fixed: 0, set at mint`;
  the Go field names are `Effect`, `Scope`, `Budget`, `ExpiresAt`, V34). So the element type of
  `[]broker.Capability` is inside the completeness test, not invisible to it.
- **CL15 `TestApproveRebuildsTheSameMintRequest`** is behavioural. It builds a `mintRequest`
  through `episode start`'s path (`prepareMint` plus the `episode.go` setters) and one through
  `approve`'s `fromRequest(Decode(Canonical(...)))`, both for a bound claude request, and requires
  them equal on every non-func field (`reflect.DeepEqual` after zeroing the three func fields).
  It requires equal nil-ness of the func fields, and requires that `afterToken` from each writes
  byte-identical MCP configs into two temp files. **This is the check that sees a dropped
  field**: a field `fromRequest` forgets comes back as its zero value and is not equal.

### 2.2 Challenge, expiry and the signed message

- **Request id**: `mr-` followed by 16 lowercase hex characters (8 random bytes), grammar
  `^mr-[0-9a-f]{16}$`. It is only a handle. It is not secret and it is not signed.
- **Challenge**: 32 random bytes per request, stored as hex. It is not secret either. It makes
  every request's signed message unique, so a signature over one request cannot approve another
  request with the same digest.
- **Expiry**: `expires_at = created_at + 600` (`RequestLifetime = 10 * time.Minute`, the parent's
  default). It is a constant, not a flag.
- **Randomness** comes from a new `sessionEnv.random func([]byte) error` seam, which is
  `crypto/rand.Read` in `realSessionEnv`. A request path that finds `random == nil` refuses, so a
  test that forgets the seam fails loudly and does not quietly fall back to real randomness.

**Signed message** (`SignedMessage(digest, challenge [32]byte, expiry int64) ([]byte, error)`):

```
"world-approval/v1" ‖ 0x00 ‖ digest (32 bytes) ‖ challenge (32 bytes) ‖ expiry (uint64 big-endian Unix seconds, 8 bytes)
```

That is 90 bytes, with a fixed prefix and fixed widths, so there is no length ambiguity. A negative
expiry is refused. A1b's SSHSIG path signs exactly these bytes under the namespace
`world-approval@sunholo`. A2's WebAuthn path uses `SHA-256(message)` as the WebAuthn challenge
(parent §2.2); this sprint only exposes `MessageHash`.

### 2.3 The admission law (pure)

```go
type AdmitInput struct {
	Now, CreatedAt, ExpiresAt int64
	Consumed                  bool
	StoredDigest, BytesDigest [32]byte // the digest column; SHA-256 of the stored bytes
}
func Admit(in AdmitInput) Verdict // Verdict: Admit | RefuseConsumed | RefuseClockBehind | RefuseExpired | RefuseDigestMismatch
```

The order is fixed, and each step is a branch with its own test row (RQ6):
1. `Consumed` gives **RefuseConsumed**. A replay always says "already used", even if it has also
   expired.
2. `Now < CreatedAt` gives **RefuseClockBehind**. The clock went back, which would otherwise make
   a request live longer than 10 minutes.
3. `Now >= ExpiresAt` gives **RefuseExpired**. At `ExpiresAt` itself the request has expired.
4. `StoredDigest != BytesDigest` gives **RefuseDigestMismatch**. The row's bytes and digest
   disagree.
5. Otherwise, **Admit**.

A1b appends signature verification as step 6. The store's consume statement (§2.4) enforces
steps 1–4 a second time inside the mint transaction, so a check that passed before the y/N cannot
go stale by the time of the write.

**[r2] `MintMatches(r MintRequest, f MintFacts) error`** replaces r1's narrow `SessionMatches`.
It compares **every column the mint writes** to `session_credentials` against the request. That
table has exactly seven columns (V30): `credential_id`, `episode_id`, `grants_json`,
`expires_at`, `created_at`, `project` and `profile_digest`. It has no `source` or `created_by`
column. `MintFacts` is a plain struct (`approval` does not import `store`) holding the seven
columns, with `grants_json` decoded, plus the token's hash and the two times that bound the mint.
Each arm returns an error naming the column:

| Arm | Requires | Mutant it kills |
|---|---|---|
| `credential_id` | equals the resolver's hash of the token handed to the writer (`hashHexToken`) | a row whose id is not the written token's hash |
| `episode_id` | `== r.Episode` | episode taken from a flag or default |
| `grants_json` | decodes; same length and order; each `effect`/`scope`/`budget` equal | grants rebuilt from the preset, or reordered |
| grant `ExpiresAt` | each `== expires_at` (the clamp in `mint.go`) | a grant left at 0, or outliving the session |
| `expires_at` | `== created_at + r.TTLSeconds` | TTL from the default 3600 |
| `created_at` | within `[request created_at, the in-transaction consume reading]` **[r3]** | a stale or forged mint time |
| `project`, `profile_digest` | `== r.Profile` (both empty iff unbound) | an unbound mint of a bound request, or another digest |

`authority.MintForRequest` runs `MintMatches` on the row it built **before** the write, and refuses
if it fails (§2.6). Effects outside the row (worktree, branch, files, the printed launch line)
happen after the mint. The tests check them with a second, test-side check (CL5 and CL16, §6),
because production cannot inspect its own later writes before making them.

`Diff(a, b MintRequest) string` names the first field that differs (for example
`profile.digest`). `approve`'s stale-world refusal uses it.

### 2.4 Storage: credentials store schema 4

The request lives in the **credentials store** because consuming it and inserting the session row
must be one SQLite transaction (§5). That store is the only file `provisionEpisode` writes under
its own lock (V5). One new table:

```sql
-- Mint requests (version 4, row 201 A0). A request is inert until approved; approving consumes
-- it in the same transaction that inserts its session row. Consumed rows stay as the audit of
-- who approved which session; unconsumed rows are pruned a day after they expire.
CREATE TABLE IF NOT EXISTS mint_requests (
    request_id    TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'mr-[0-9a-f]*'),
    request       BLOB NOT NULL,
    digest        TEXT NOT NULL CHECK (length(digest) = 64),
    challenge     TEXT NOT NULL CHECK (length(challenge) = 64),
    episode_id    TEXT NOT NULL CHECK (episode_id <> ''),
    created_at    INTEGER NOT NULL,
    expires_at    INTEGER NOT NULL CHECK (expires_at > created_at),
    consumed_at   INTEGER,
    approved_by   TEXT CHECK (approved_by IS NULL OR approved_by <> ''),
    credential_id TEXT,
    CHECK ((consumed_at IS NULL) = (approved_by IS NULL) AND (consumed_at IS NULL) = (credential_id IS NULL))
);
```

- `credentialsSchemaVersion` goes from 3 to 4. The upgrade adds this table, plus the earlier steps
  for a store at version 1 or 2, in one transaction. The open check (`credentials.go:304`) accepts
  version 4 with exactly three tables. `ensureCredentialsStore`'s log line names the upgrade
  generically, not "exec profile binding".
- **The upgrade costs no extra operator step.** `episode start` already refuses a daemon that is
  not this build (V4), so the operator already restarts `serve` with the new binary, and that
  restart upgrades the store (V6). A daemon of an older build refuses a version-4 store, as it
  already refuses any newer store.
- `approved_by` is `tty` in A0. A1b writes `ssh:<approver name>`. No `approvers` table is added
  here (§0). A1b adds it at version 5. Its upgrade costs the same single restart, for the reason in
  the previous bullet.

New store methods, in `host/store/mint_requests.go`. Each one requires a context deadline, like
every other method:

| Method | What it does |
|---|---|
| `PutMintRequest(ctx, row, now)` | refuses a store below version 4 (`*OlderCredentialsStoreError`, which names the serve restart); in one transaction it deletes unconsumed rows whose `expires_at <= now - 86400`, refuses if 16 or more unconsumed, unexpired rows remain (`ErrTooManyPendingRequests`), then inserts |
| `MintRequest(ctx, id)` | point lookup, `(row, ok, err)` |
| `PendingMintRequests(ctx, now, limit)` | unconsumed rows with `expires_at > now`, newest first |
| `MintSessionForRequest(ctx, session SessionRow, c RequestConsume)` | one transaction. `UPDATE mint_requests SET consumed_at=?, approved_by=?, credential_id=? WHERE request_id=? AND digest=? AND consumed_at IS NULL AND created_at <= ? AND expires_at > ?`. If the update did not change exactly one row, it reads the row back, classifies the reason as `ErrRequestUnknown`, `ErrRequestConsumed`, `ErrRequestExpired`, `ErrRequestClockBehind` or `ErrRequestDigestChanged`, and rolls back. Otherwise it inserts the `session_credentials` row and commits |

`RequestConsume{ID, Digest, ApprovedBy string; Clock func() int64}`. **[r3, gpt verbatim]** Consumption obtains a fresh clock reading inside the mint transaction, immediately before the guarded UPDATE and after worktree preparation; that reading is used for expiry admission, `consumed_at`, and `MintMatches`' consume-time bound. The guarantee is defined precisely as **admission at transactional consumption**, not commit completion. `MintSessionForRequest` calls `c.Clock()` exactly once, after `BeginTx`, and refuses a nil `Clock`. `Digest` is the digest the approver
saw. If anything rewrote the row between the prompt and the write, the `WHERE digest = ?` matches
nothing and the mint is refused.

### 2.5 Creating a request: `episode start --request`

`ailang-worldd episode start <episode> [the same flags] --request`:

1. Runs every check plain `episode start` runs, unchanged: `prepareMint`, the daemon health and
   build check, and the session and MCP files not existing yet.
2. Refuses `--request` together with `--launch`, `--headless`, `--resume` or a launch preset
   (exit 2, naming why). Those paths carry callbacks and start records that this sprint does not
   serialise. They are later work (§11).
3. Builds the `MintRequest` from the validated `mintRequest`, the daemon's `HealthResponse.Build`
   and its address, then runs `Validate()`.
4. **Does not open `/dev/tty`.** A request grants nothing, and the approval is the gate. Under a 3 s
   root (`createMintRequest`, a new context root, §7), it opens the credentials store writable and
   calls `PutMintRequest`.
5. Prints the request as a block: the id, the first 16 hex characters of the digest, the expiry
   (RFC 3339 and "in 10m"), the same lines the y/N would name, and the next step,
   `ailang-worldd approve <id>`. It prints no token, because none exists yet.

`ailang-worldd approve --list` prints the pending requests read-only, using a read-only handle and
a 3 s root, in the same `runApprove` root.

### 2.6 Approving with the typed `y`: `ailang-worldd approve <request-id>`

`approve <request-id> [--db <path>] [--credentials-db <path>]`. The global `--addr` is refused, as
it is for `session`. The daemon address comes from the request.

1. Checks the grammar of `<request-id>`, then reads the row read-only under a 3 s root. An unknown
   id is refused with exit 1.
2. Runs `Decode` on the stored bytes, then `Admit` with `now`. Each refusal names its reason and
   the fix ("create a new request with `episode start --request`"). **No terminal is opened
   before this step passes** (CL9 pins this with the spy).
3. **Re-derives the request from the world as it is now**: it rebuilds `mintOptions` from the
   decoded request through **[r2]** `fromRequest`, which follows the §2.1a map field by field (CL14,
   CL15). It refuses if `--db`/`--credentials-db` do not resolve to the request's `CredentialsDB`.
   It then runs `prepareMint`, fetches the daemon's health at `Daemon.Addr`, runs
   `checkDaemon`, and builds a fresh `MintRequest` the same way §2.5 step 3 does. If
   `Canonical(fresh) != stored bytes`, it refuses with "the world changed since the request was
   made: `<Diff>`". This covers a changed profile file, a restarted daemon of another build, a
   worktree that now exists and a session file that appeared. The request stays pending.
4. Calls `provisionEpisode` with three new `mintRequest` fields:
   - `presetPrompt` (already printed before the y/N, `:430-431`) is set to the request block,
     with the id and the full digest.
   - `consume = &store.RequestConsume{ID, Digest, ApprovedBy: "tty", Clock: env.now}`.
     **[r3, gpt verbatim]** consumption obtains a fresh clock reading inside the mint transaction,
     immediately before the guarded UPDATE and after worktree preparation (r2 sampled it once after
     the y/N, so work between that read and the transaction could expire the request while SQL still
     compared against the earlier timestamp). A human who takes longer than the expiry at the prompt,
     or a worktree preparation that runs past it, is refused at the write, and the worktree that this
     run created is rolled back by the existing step 7.
   - **[r2]** `answerBy = expires_at`: the y/N read is bounded by the request's remaining lifetime
     (§2.6a).
5. `provisionEpisode`'s mint branch gains `else if r.consume != nil` →
   `authority.MintForRequest(...)`. That builds the row exactly as `mintRecorded` does, checks
   **[r2]** `approval.MintMatches(request, facts)` on all seven columns (§2.3, step 2 of §5), then calls
   `st.MintSessionForRequest`. The token file and MCP config are written by the existing code
   after the mint.
6. Prints the launch block, as `episode start` does (`printLaunch`), from the request's fields.

`provisionEpisode` stays the only fence-and-mint path. `approve` is a third caller with a request
attached, and the comment at `session_new.go:24-26` gains that caller.

#### 2.6a [r2] The y/N is bounded by the request's lifetime

**The defect.** `provisionEpisode` opens the terminal (`session_new.go:399`), prints the y/N and
reads the answer with `readTermLine(term)` (`:438`) with no deadline. Its 75 s context
(`sessionNewBudget`) starts only after the read (`:444`, V2, V33). So a prompt nobody answers
waits forever. That is **pre-existing** for `session new` and `episode start`. `approve` is a new
caller, and the request's expiry does not bound the wait, because r1 checked expiry only after the
read returned.

**The mechanism.**
- `mintRequest` gains `answerBy int64` (Unix seconds). It is 0 for every existing caller, which
  keeps the existing unbounded read. `approve` sets it to the request's `expires_at`.
- `sessionEnv` gains `after func(time.Duration) <-chan time.Time`; `realSessionEnv` sets it to
  `time.After`.
- In step 3 of `provisionEpisode`, when `answerBy != 0`, it reads `remaining := answerBy -
  env.now()` just before printing the prompt.
  - If `remaining <= 0`, it refuses without prompting: `request <id> expired before the prompt;
    make a new request`.
  - Otherwise it calls `readTermLineWithin(term, remaining*time.Second, env.after)` (**[r3, gemini]** its doc comment states it is only safe for one-shot CLI entrypoints that exit the process on return, because a timed-out read leaves its goroutine blocked on the fd). That runs
    `readTermLine(term)` in a goroutine that sends to a buffered channel of size 1, and `select`s
    on that channel against `after(d)`.
  - When the timer wins, it refuses with exit 2: `ailang-worldd approve: no answer within the
    request's remaining lifetime (<d>); nothing was minted and request <id> is still pending
    until <expires_at>`.
- The refusal comes before the credentials lock, the worktree and the mint, so it writes nothing,
  and the request stays pending.

**Why it works on a tty on darwin and Linux.** It never needs the blocked read to return. It does
not use `SetReadDeadline`, and it does not rely on `Close` to unblock the read. I measured both on
darwin with `/dev/tty` as the real controlling terminal (a child of `pty.fork`, V32):
- `SetReadDeadline` returns `file type does not support deadline`.
- After the timer fired, `Close()` returned nil, but **the read was still blocked 500 ms later**.
- The `select` on the timer returned at 500 ms, `main` returned, and the process exited with
  status 0 at 1.02 s wall time. The blocked reader died with the process.

On an ordinary pty slave, `/dev/ttysNNN`, a deadline *did* work (`i/o timeout` at 200 ms, V31).
So the behaviour depends on the device, which is why the design relies on neither. `approve` is a
short-lived process: after the refusal `runApprove` returns, `main` exits, and the abandoned
goroutine goes with it. Nothing after the refusal waits on the terminal. The deferred `term.Close()`
still runs and returns at once (V32). Linux was **not measured**, but the mechanism uses only a
timer, a channel and `select`, which behave the same on every OS Go supports.

**Unchanged for the existing callers.** `session new`, `session mint`, `episode start`, presets
and resume keep `answerBy == 0`, so their y/N is unbounded exactly as at HEAD. That pre-existing
wait is **out of scope and logged** as residual R7 (§9), not fixed here. CL17 shows it is
unchanged: with `answerBy == 0`, `env.after` is never called.

### 2.7 The SSHSIG verifier (A1a)

`host/approval/sshsig.go`, standard library only: `crypto/ed25519`, `crypto/ecdsa`,
`crypto/elliptic`, `crypto/sha256`, `crypto/sha512`, `encoding/base64`, `encoding/binary`,
`math/big`, `bytes`, `errors`. **Correction to the parent §2.1 list:** `crypto/sha512` is
required, because `ssh-keygen -Y sign` hashes with SHA-512 by default (V12, V13). Both are in the
standard library, so `TestDaemonDependencyAllowlist` is unaffected.

```go
const Namespace = "world-approval@sunholo"

type SSHPublicKey struct{ Type string; Blob []byte }   // Blob = the SSH wire encoding
func ParseAuthorizedKey(line []byte) (SSHPublicKey, error) // "ssh-ed25519 AAAA… comment"
func (k SSHPublicKey) Fingerprint() string                 // "SHA256:" + unpadded base64, as ssh-keygen -l prints

type SSHSig struct{ PublicKey SSHPublicKey; Namespace, HashAlg string; SigType string; Sig []byte }
func ParseSSHSIG(armored []byte) (*SSHSig, error)
func VerifySSHSIG(s *SSHSig, namespace string, message []byte, key SSHPublicKey) error
func VerifyApproval(armored []byte, digest, challenge [32]byte, expiry int64, key SSHPublicKey) error
```

**Parse** (PROTOCOL.sshsig). The input must be exactly one `-----BEGIN SSH SIGNATURE-----` …
`-----END SSH SIGNATURE-----` block. Body lines are standard base64, and nothing other than one
optional trailing newline may surround the block. The decoded blob is `"SSHSIG"`, `uint32 version
= 1`, then the strings `publickey`, `namespace`, `reserved`, `hash_algorithm` and `signature`,
with no trailing bytes. Every length is bounds-checked before slicing.

**Every failure is a typed error, and the parser never returns a value with a nil error unless it
parsed completely.** The named errors: `ErrArmor` (including empty input), `ErrBase64`,
`ErrMagic`, `ErrVersion`, `ErrTruncated`, `ErrTrailing`, `ErrReserved` (non-empty; ssh-keygen
writes it empty, V12), `ErrHashAlg` (anything but `sha256` or `sha512`) and `ErrKeyType`.

**Verify** is a fixed sequence of named refusals:
1. The namespace must equal the argument (`ErrNamespace`).
2. The signature's public key blob must be **byte-equal** to the expected key's blob
   (`ErrKeyMismatch`).
3. The key type must be `ssh-ed25519` or `ecdsa-sha2-nistp256` (`ErrKeyType`). The curve name
   must be `nistp256` and the point must parse with `ecdsa.ParseUncompressedPublicKey`
   (go1.26.6, V13).
4. The signature's format id must equal the key type (`ErrSigType`).
5. The signed data is built as `"SSHSIG" ‖ string(namespace) ‖ string(reserved) ‖
   string(hash_alg) ‖ string(H(message))`.
   - **Ed25519**: the signature is 64 bytes, checked with `ed25519.Verify` over the signed data.
   - **ECDSA P-256**: the signature blob is two `mpint`s (`r`, `s`). A negative or non-minimal
     `mpint` is refused. The check is `ecdsa.Verify` over `SHA-256(signed data)`.

   Any failure is `ErrBadSignature`.

`VerifyApproval` builds `SignedMessage(digest, challenge, expiry)` and calls `VerifySSHSIG` with
`Namespace`. This is the single function A1b's `approve --signature` will call after `Admit`.

**Key types refused in this sprint:**
- **`sk-ssh-ed25519@openssh.com`, `sk-ecdsa-sha2-nistp256@openssh.com`.** A security-key
  signature carries a flags byte and a counter, and its user-presence bit is the thing that would
  prove a touch. No fixture can be made without hardware: `ssh-keygen -t ed25519-sk` on this rig
  fails with "No FIDO SecurityKeyProvider specified" (V15). So this sprint **refuses `sk-*`** with
  `ErrKeyType` and a message naming why, and records it as open question Q1.
- **RSA, ECDSA P-384/P-521, DSA and certificates.** They are outside R2's list. Each is refused
  with `ErrKeyType`.

A prototype of exactly this verifier, about 120 lines using the standard library only, was run
against real signatures. It accepted the Ed25519 and P-256 signatures and refused a tampered
message, a `git`-namespace signature, another key's signature and a P-256 signature presented with
the Ed25519 key. `ssh-keygen -Y verify` gave the same answer on every arm (V16, V17). The prototype
is not shipped. It only shows that the design is buildable without `x/crypto`.

### 2.8 Fixtures: real, pinned, regenerable

The directory `host/approval/testdata/sshsig/` is checked in and holds:

| File | Made by |
|---|---|
| `message.bin` | the 90-byte `SignedMessage` of the fixture request, from `fixtureMessageHex` in `sshsig_test.go` |
| `ed25519.pub`, `ed25519.sig` | `ssh-keygen -t ed25519` / `-Y sign -n world-approval@sunholo` |
| `p256.pub`, `p256.sig` | `ssh-keygen -t ecdsa -b 256` / `-Y sign -n world-approval@sunholo` |
| `git-ns-ed25519.sig` | the same Ed25519 key, `-Y sign -n git`: a real signature for another namespace |
| `SOURCE.txt` | `ssh -V`, the UTC time, and `ssh-keygen -l` of both public keys |
| `generate.sh` | the generator below |

That is three separate `ssh-keygen -Y sign` invocations and two key types. **No private key is
checked in.** The generator makes the keys in a `mktemp -d` directory and deletes it on exit
(V18 ran it: zero private keys were left in the fixture directory).

`generate.sh` (exact; run by hand on a rig, never in CI):

```bash
#!/usr/bin/env bash
# Regenerates the SSHSIG fixtures with REAL ssh-keygen (row 201 A1a). Test keys only:
# they are made in a temp dir and deleted on exit; only public keys and signatures are kept.
set -euo pipefail
here=$(cd "$(dirname "$0")" && pwd)
hex=$(sed -n 's/^const fixtureMessageHex = "\([0-9a-f]*\)"$/\1/p' "$here/../../sshsig_test.go")
[ ${#hex} -eq 180 ] || { echo "fixtureMessageHex not found" >&2; exit 1; }
tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
printf '%s' "$hex" | xxd -r -p > "$here/message.bin"
ssh-keygen -q -t ed25519 -N '' -C world-approval-fixture-ed25519 -f "$tmp/ed25519"
ssh-keygen -q -t ecdsa -b 256 -N '' -C world-approval-fixture-p256 -f "$tmp/p256"
ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/ed25519" < "$here/message.bin" > "$here/ed25519.sig"
ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/p256" < "$here/message.bin" > "$here/p256.sig"
ssh-keygen -Y sign -n git -f "$tmp/ed25519" < "$here/message.bin" > "$here/git-ns-ed25519.sig"
cp "$tmp/ed25519.pub" "$tmp/p256.pub" "$here/"
{ ssh -V 2>&1; date -u +%FT%TZ; ssh-keygen -l -f "$here/ed25519.pub"; ssh-keygen -l -f "$here/p256.pub"; } > "$here/SOURCE.txt"
```

- **Sign from stdin, never from a file argument.** `ssh-keygen -Y sign -f k msg.bin` writes
  `msg.bin.sig`. If that file exists it asks "Overwrite (y/n)?", and with no answer it **exits 0
  without writing** (V12). A generator that signed in file mode could keep a stale signature and
  still look green. From stdin, the signature goes to stdout.
- **The comment requirement.** `sshsig_test.go` carries, directly above the fixture table, a
  comment naming `testdata/sshsig/generate.sh` as the only source of these bytes and quoting its
  three `-Y sign` lines. `SOURCE.txt` records the tool version.
- **Pinned against VCS transformations.** A new root `.gitattributes` (none exists today, and
  `core.autocrlf` is unset on this rig, V19) holds `host/approval/testdata/sshsig/* -text`. Test
  SG1 also pins the SHA-256 of every fixture file (**[r3, kimi verbatim]** regeneration changes keys and therefore the SG1 pinned-hash table and SOURCE.txt must be updated in the same commit by whoever runs generate.sh), so a CRLF conversion, a re-wrap or an edit is
  red at once.
- **Derived refusal arms.** Arms such as a flipped signature byte, a retyped key (`sk-…`, `rsa-…`)
  or a truncated blob are **made in test code from the real fixtures by named transforms**
  (decode, change one field, re-armor). They are never typed by hand. Each transform is a small
  function whose name says what it changes.

**The cross-check** (`host/approval/sshsig_crosscheck_test.go`, test XC1) generates fresh keys in
`t.TempDir()` with `ssh-keygen` for each of Ed25519 and P-256, then signs `message.bin` from
stdin. On each arm (good, tampered message, `git` namespace, another key, flipped signature byte)
it requires `ssh-keygen -Y verify` to exit 0 exactly when our `VerifySSHSIG` returns nil. The
`good` arm must be accepted by both, so the test is not vacuous. It also asserts that every
signature `ssh-keygen` produced is non-empty and parses, so a silent no-write fails.

**If `ssh-keygen` is absent:** when `WORLD_REQUIRE_SSH_KEYGEN=1`, the test calls `t.Fatal`.
Otherwise it calls `t.Skip("SKIP COUNTED: ssh-keygen not found; the SSHSIG cross-check did NOT
run")`. The core verifier never depends on it: SG1–SG6 use only the checked-in fixtures. A new CI
step, placed after Go is set up in the `ubuntu-24.04` job, makes the skip impossible to miss:

```yaml
      - name: Row 201 SSHSIG cross-check with real ssh-keygen, verbose
        timeout-minutes: 3
        run: |
          set -euo pipefail
          command -v ssh-keygen || { echo "::error::ssh-keygen missing (apt-get install openssh-client)"; exit 1; }
          ssh -V
          log="$RUNNER_TEMP/sshsig-crosscheck.log"
          rc=0
          WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -count=1 -v \
            -run '^TestSSHSIGCrossCheckWithSSHKeygen$' >"$log" 2>&1 || rc=$?
          cat "$log"
          [ "$rc" -eq 0 ] || exit "$rc"
          for t in TestSSHSIGCrossCheckWithSSHKeygen TestSSHSIGCrossCheckWithSSHKeygen/ed25519 TestSSHSIGCrossCheckWithSSHKeygen/ecdsa-p256; do
            grep -qE -- "^ *--- PASS: $t " "$log" || { echo "::error::$t did not PASS (SKIP?)"; exit 1; }
          done
```

Whether the `ubuntu-24.04` image ships `ssh-keygen` is **UNVERIFIED** here (no network, V14). If
it does not, the step fails with the install line, and the fix is one `apt-get` line in the same
step.

### 2.9 Seams left for later milestones (named, not built)

- **A1b**: `approve <id> --signature <file>` = `Admit`, then `VerifyApproval` against the
  `approvers` row, then the same `provisionEpisode` with `consume.ApprovedBy = "ssh:<name>"` and no
  terminal. `approve <id> --sign-with <key>` runs `ssh-keygen -Y sign -n world-approval@sunholo
  -f <key>` from stdin (it accepts a public key whose private half is in `ssh-agent`, V13). Also
  `approve <id> --write-message <file>`, `approver enroll|list|revoke` (typed `y`), and the
  `approvers` table at version 5.
- **A2**: `MessageHash` as the WebAuthn challenge; `approved_by = "webauthn:<name>"`.
- **A3**: the daemon would have to write the credentials store to create and consume requests from
  the surface. Today it never holds that store's writer lock (V6). That is a deliberate boundary,
  so A3 must design it. `PendingMintRequests` is the *needs you* source.
- **A4**: `Human.Approve` (`host/broker/approve.go:18`, V26) packets get their own canonical form and
  digest. Whether `Admit` moves into the `.ail` kernel then is Q2.

## 3. Pure core, effects at the boundary (coding-standards S1–S3)

| Piece | Where | Pure? | Tested by |
|---|---|---|---|
| `MintRequest.Validate`, `Canonical`, `Decode`, `Digest`, `Diff` | `host/approval/request.go` | yes | RQ1–RQ4, RQ8 |
| `SignedMessage`, `MessageHash` | `host/approval/request.go` | yes | RQ5 |
| `Admit`, `MintMatches` **[r2]** | `host/approval/admit.go` | yes | RQ6, RQ7 |
| `ParseAuthorizedKey`, `Fingerprint`, `ParseSSHSIG`, `VerifySSHSIG`, `VerifyApproval` | `host/approval/sshsig.go` | yes | SG1–SG6, XC1 |
| clock, randomness, `/dev/tty`, **[r2]** the prompt timer | `sessionEnv` (`now`, `random`, `openTerminal`, `after`) | seam | CL1–CL18 inject them |
| the store | `host/store/mint_requests.go` | effect | ST1–ST7 |
| daemon health | `getDaemonHealth` (existing) | effect | the existing `stubDaemon` rig |
| `ssh-keygen` | tests and the generator only; no production code runs it in this sprint | effect | XC1 |

**Why not `.ail` with Z3 contracts (S1), and why not a package (S3).**
- **The cryptography cannot be written in AILANG.** The pinned `AILANG v0.41.0` `std/crypto`
  offers `sha256bytes`, `sha256hex`, `hmacsha256`, `rsa_verify_pkcs1v15` and
  `constanttimeequal`. It has **no Ed25519, ECDSA or SHA-512** (V20). So neither the verifier nor
  the digest of an SSHSIG signed-data block can live in `.ail`, and no registry package could hold
  them either.
- **The admission law could be written in AILANG, and I recommend not doing it now.** `Admit` is
  four comparisons over integers and booleans, which Z3 can encode. But:
  - Its inputs come from a store row, the clock and a hash, which are all host effects, and Go
    executes it.
  - A `.ail` copy would be a transcription under a drift test (the broker precedent), and it would
    have to join `world/`, the frozen kernel. That raises the verified floor and triggers S8's
    six-file coupling (`REQUIRED_VERIFIED`, `EXACT_TOTAL_VERIFIED=16`, the package golden and the
    runbook digests; V21).

  That cost is out of proportion to a five-row table that RQ6 tests exhaustively. Q2 asks whether
  to revisit at A4, where `Human.Approve`, a kernel effect, uses the same law.
- **Why `host/approval` and not a package.** Its consumers are the mint path in
  `cmd/ailang-worldd` and `host/authority`, which are host Go, and the cryptography is host-only.
  S2 is kept: every effect (tty, clock, randomness, store, network, subprocess) is injected or sits
  in `cmd/` or `host/store`, and `host/approval` makes no I/O call.

## 4. Threats (A0 and A1a)

| Threat | What this sprint does | Pinned by |
|---|---|---|
| **Replay**: approve the same request twice | `Admit` step 1. The consume `UPDATE … WHERE consumed_at IS NULL` sits in the mint transaction and must change exactly one row. For signatures, each request has its own challenge, so a signature over request R1 does not verify for R2, even with the same digest | RQ6, ST4, CL8, SG3 (challenge-flip arm), SG5 |
| **Expiry** | 10 minutes. Checked before the terminal opens, and again in SQL at the write, using a clock read **after** the y/N | RQ6, ST5, CL9 |
| **Digest substitution**: approve a small request, mint a big one | The mint is built only from the decoded bytes. The digest is recomputed from those bytes. The consume requires `digest = ?`, the digest the approver saw. **[r2]** `MintMatches` checks all seven columns of the built row before the write. `fromRequest` follows a field map that tests keep complete. A signature (A1b) covers the digest | RQ3, RQ7, ST5, CL5, CL11, CL14, CL15, SG5 |
| **[r2] A field dropped between request and mint**: one that never entered `MintRequest`, so byte equality cannot see it | The §2.1a map covers every field of `mintOptions`/`mintRequest`. CL14 fails on an unmapped field, and CL15 fails if `approve`'s rebuild differs from `episode start`'s in any field | CL14, CL15 |
| **[r2] An unanswered prompt holds `approve` forever** | The y/N read is bounded by `expires_at - now` with a goroutine, a timer and `select`. On timeout it refuses and leaves the request pending. This does not depend on tty read deadlines (V31, V32) | CL13, CL18 |
| **Stale world**: the profile, daemon or worktree changed after the request | `approve` re-derives the request and requires byte equality; otherwise it refuses naming the field | CL10 |
| **Namespace confusion**: a signature made for `git` or `file` | Namespace equality, plus the namespace inside the signed data, so the cryptography fails too. The fixture is a real `-n git` signature | SG3, XC1 |
| **Wrong key** | The signature's key blob must equal the expected blob, byte for byte | SG3, XC1 |
| **Revoked key** | **A1b.** There is no key store in this sprint | — |
| **Enrolment without the typed `y`** | **A1b.** No enrolment exists in this sprint. The only approval here is the typed `y` at `/dev/tty`, through the unchanged fence | CL6 |
| **An agent process holding a key FILE** | **A1 proves possession of the private key, nothing more.** SSHSIG for Ed25519 and P-256 carries no user-presence bit. If a key is an unencrypted file that an agent process can read, the agent can sign. World's exec sandbox denies reads of the operator's home (`exec_settings.go:113`, V22), but an agent harness running as the operator *outside* that sandbox can read whatever the operator can. A1 cannot tell a passphrase-, agent-confirm- or Secure-Enclave-held key from a bare file. That is operator policy (parent §4 row 2), and A1b's enrolment text must say so. `sk-*` keys, whose signatures do carry a touch bit, are refused in this sprint (Q1) | stated, not tested |
| **An agent with write access to `credentials.db`** | Unchanged and out of scope: such a process can already insert a session row directly. `prepareMint` refuses a workspace root that contains the store (`session_new.go:257-264`) | existing |
| **Clock skew** | One clock: the approving host's, which is also the store's host in A0/A1. A clock that went back is refused (`RefuseClockBehind`). One that jumped forward expires the request early, which is fail-closed. A signer's clock is never read | RQ6, ST5 |
| **Concurrent consumption**: two approvals race | The credentials writer lock does not wait, so the second `approve` fails fast. Even without the lock, the guarded `UPDATE` and the `INSERT` are one transaction, so exactly one mints. With 8 racers: one success, and one session row | ST6 (`-race`) |
| **Request flooding** | At most 16 pending requests per store. Expired ones are pruned after a day | ST7 |
| **Token sent to a readable path** | The session file and MCP paths are in the digest and shown on the prompt. `approve` re-runs `checkOutDir`, which requires a private directory | CL5, CL10 |
| **A parser that accepts what it cannot read** | Every non-empty input it cannot parse is an error with a nil result. The fixtures assert the key type each one parsed to | SG2, SG4 |

## 5. By-construction argument

1. **Nothing is minted without a typed `y`.** `approve` reaches the mint only through
   `provisionEpisode`. Its first statement is the unchanged `env.openTerminal()`, then `isYes` on
   the terminal's line (`session_new.go:399`, `:438`). There is no other mint branch in `approve`.
2. **Nothing is minted that was not in the request.** `approve` builds `mintOptions` only from
   `Decode(stored bytes)`. `Decode` accepts only canonical bytes. The re-derived request must be
   byte-equal to them, and **[r2]** `fromRequest` reproduces every field `episode start` sets (CL15) under a complete field map (CL14). `MintForRequest` refuses a row that `MintMatches` rejects on any of the seven columns. So the
   session's episode, grants, TTL and profile are the request's.
3. **The bytes are the bytes the human saw.** The prompt is rendered from the decoded request, and
   the consume binds `digest = SHA-256(those bytes)`. A row rewritten after the prompt no longer
   matches.
4. **At most one session per request.** Consumption and insertion share one transaction, and the
   consume must change exactly one row whose `consumed_at IS NULL`. A rolled-back insert leaves the
   request unconsumed, and a rolled-back consume leaves no session row.
5. **No mint after expiry.** The SQL `expires_at > ?` uses a clock read after the y/N.
5a. **[r2] `approve` cannot wait forever at the prompt.** With `answerBy != 0` the only read of
   the answer is `readTermLineWithin`, which returns by `after(expires_at - now)` at the latest
   (§2.6a).
6. **Creating a request mints nothing.** `createMintRequest` calls `PutMintRequest` only. It never
   calls `provisionEpisode` or `authority.*`, and it opens no terminal (CL1's spy).
7. **The verifier accepts only the expected key, namespace and message.** Each is an equality on
   the parsed data before any signature check, and the namespace and message hash are inside the
   bytes that are verified.

## 6. Test plan

Each row names the mutant it kills. Where the deliverable is a refusal, the unit is the branch:
deleting that branch must turn exactly that arm red, or turn it into a different typed error,
which the arm also checks with `errors.Is`.

**M1 `host/approval/request_test.go`, `admit_test.go`**

| Id | Test | Asserts | Mutant killed |
|---|---|---|---|
| RQ1 | `TestCanonicalRequestGolden` | a fixed fixture request gives golden bytes and a golden digest | field order changed; HTML escaping on; trailing newline |
| RQ2 | `TestDecodeIsStrict` | the round trip is equal; arms: unknown field, trailing bytes, reformatted JSON, reordered fields, two concatenated values → `ErrNotCanonical` or a decode error | `Decode` without the re-encode comparison (the reformatted arm is accepted) |
| RQ3 | `TestDigestCoversEveryField` | reflect over every leaf field (each grant member, profile, daemon); changing any one changes the digest; the leaf count is asserted ≥ 20, so a reflection walk that finds nothing is red | a field tagged `json:"-"` or dropped from the encoding |
| RQ4 | `TestValidateRefusals` | one arm per `Validate` branch (§2.1) | each branch deleted |
| RQ5 | `TestSignedMessageGolden` | golden 90 bytes (`fixtureMessageHex`); a negative expiry is refused | no NUL after the prefix; little-endian expiry |
| RQ6 | `TestAdmitLaw` | every branch, plus the boundaries `Now == ExpiresAt` (expired), `ExpiresAt-1` (admit), `Now == CreatedAt` (admit), and consumed+expired → consumed | `>` for `>=`; each branch deleted; order swapped |
| RQ7 | **[r2]** `TestMintMatchesEveryColumn` | one arm per row of the §2.3 table: `credential_id` not the token's hash; another `episode_id`; an extra, reordered or budget+1 grant; a grant `ExpiresAt` of 0 or `expires_at+1`; `expires_at` off by one from `created_at+ttl`; `created_at` before the request or after the consume; another `profile_digest`; a bound request with an unbound row. The clean arm passes | each arm's comparison deleted; grants compared by count only |
| RQ8 | `TestDiffNamesTheField` | names `profile.digest`, `grants[1].budget`, `daemon.executable_sha256` | `Diff` returns "" on a difference |

**M2 `host/approval/sshsig_test.go`** (fixtures only, no `ssh-keygen`)

| Id | Test | Asserts | Mutant killed |
|---|---|---|---|
| SG1 | `TestSSHSIGFixturesAreTheGeneratedBytes` | SHA-256 of every fixture file equals a pinned table; `message.bin` equals `hex(fixtureMessageHex)` and `SignedMessage(fixture)`; each `.sig` is non-empty and starts with the armor line | a CRLF-converted or edited fixture; a message encoding drift |
| SG2 | `TestSSHSIGFixturesVerify` | `ed25519.sig` and `p256.sig` verify against their `.pub`; the parsed key types are `ssh-ed25519` and `ecdsa-sha2-nistp256`, the hash is `sha512`, the namespace is `Namespace` | ECDSA hashed with SHA-512; a parser that returns a zero value with nil |
| SG3 | `TestSSHSIGRefusals` | one arm per named error (§2.7): real `git-ns-ed25519.sig` → `ErrNamespace`; `ed25519.sig` with `p256.pub` → `ErrKeyMismatch`; message, challenge or expiry byte flipped → `ErrBadSignature`; last signature byte flipped → `ErrBadSignature`; key retyped `sk-ssh-ed25519@openssh.com` or `rsa-sha2-512` → `ErrKeyType`; hash `sha1` → `ErrHashAlg`; version 2 → `ErrVersion`; magic → `ErrMagic`; reserved non-empty → `ErrReserved`; trailing byte → `ErrTrailing`; signature type ≠ key type → `ErrSigType`; negative `mpint` → `ErrBadSignature` | each branch deleted |
| SG4 | `TestParseSSHSIGFailsLoudly` | **every** prefix of each fixture blob, re-armored, → error and nil result, no panic; non-empty garbage (seeded random bytes, a PGP armored block, JSON, the right armor around a wrong body) → error and nil; empty input → `ErrArmor` | a parser that returns `(&SSHSig{}, nil)` on what it cannot read; a missing bounds check (panic) |
| SG5 | `TestVerifyApprovalBindsTheRequest` | the fixture's (digest, challenge, expiry) → ok; each of the three changed → refused | `VerifyApproval` omitting one of the three from the message |
| SG6 | `TestFingerprintMatchesSourceTxt` | `Fingerprint()` of both `.pub` files equals the `ssh-keygen -l` lines recorded in `SOURCE.txt` | padded base64; hashing the text line, not the blob |

**M3 `host/approval/sshsig_crosscheck_test.go`**

| Id | Test | Asserts | Mutant killed |
|---|---|---|---|
| XC1 | `TestSSHSIGCrossCheckWithSSHKeygen` (subtests `ed25519`, `ecdsa-p256`) | keys in `t.TempDir()`; signing from stdin; on 5 arms `ssh-keygen -Y verify` exits 0 ⇔ ours returns nil; `good` is accepted by both; fingerprints agree | a verifier that refuses everything (the `good` arm); one that accepts everything (the tamper arms) |

**M4 `host/store/mint_requests_test.go`, edited credentials tests**

| Id | Test | Asserts | Mutant killed |
|---|---|---|---|
| ST1 | `TestCredentialsSchemaMatchesCanonicalDDL` (edited) | three tables, the canonical `mint_requests` DDL, `user_version` 4 | — |
| ST2 | `TestCredentialsUpgradeToV4` | stores at v1 and v3 → v4, every existing row unchanged, idempotent; a session verb leaves v3 alone (the existing v1 test's pattern) | the upgrade drops a statement; a CLI open upgrades |
| ST3 | `TestPutMintRequestNeedsVersion4` | a v3 handle → `*OlderCredentialsStoreError`, whose text names the serve restart | the version check deleted (a raw "no such table" instead) |
| ST4 | `TestMintSessionForRequestConsumesOnce` | first call ok, the row carries `approved_by=tty` and the credential id; the second → `ErrRequestConsumed`; exactly 1 session row | `RowsAffected` not checked |
| ST5 | `TestMintSessionForRequestRefusals` | `now == expires_at`; `now < created_at`; another digest; unknown id → each typed error, and **0** session rows | the insert placed before the guarded update and kept on refusal |
| ST6 | `TestConcurrentApprovalsMintOnce` (`-race`) | 8 goroutines, each with its own `OpenCredentials` on one file: exactly 1 succeeds, the rest get `WriterAlreadyActive` or `ErrRequestConsumed`, and there is 1 session row | the consume moved outside the transaction |
| ST7 | `TestPutMintRequestCapsAndPrunes` | the 16th pending insert is ok and the 17th is refused; a row expired more than a day ago is pruned; a consumed row is never pruned | the cap or the consumed-row guard deleted |
| edited | `TestCredentialsMethodsRequireADeadline` | covers the 4 new methods | — |
| edited | `host/daemon` `ensureCredentialsStore` test of the upgrade log line, if one pins the text (none found, V23) | — | — |

**M5 and M6 `cmd/ailang-worldd/approve_test.go`** (the existing `stubDaemon`, `spyEnv`,
`testMintEnv` and `noTerminalEnv` helpers, plus a `random` seam)

| Id | Test | Asserts | Mutant killed |
|---|---|---|---|
| CL1 | `TestEpisodeStartRequestStoresAndOpensNoTerminal` | the spy shows no terminal opened; 1 pending row; decoded episode, grants, TTL and profile equal what the y/N path would mint; stdout has the id and `ailang-worldd approve <id>`; no worktree, session or MCP file | the request path calls `provisionEpisode` |
| CL2 | `TestEpisodeStartRequestRefusesUnsupportedModes` | `--request` with `--launch`, `--headless`, `--resume` or a launch preset → exit 2, 0 rows | a refusal branch deleted |
| CL3 | `TestEpisodeStartRequestRefusesAnOlderStore` | a v3 store → refused naming the restart, 0 rows | — (ST3 at the CLI) |
| CL4 | `TestApproveListShowsPending` | lists the pending row, not the consumed or expired ones | the pending filter deleted |
| CL5 | `TestApproveMintsExactlyTheRequest` | `y` → **[r2]** `MintMatches` on the stored row (all seven columns), **and** the test-side effect check `requireApproveEffects` (CL16) on everything else the mint did: the worktree exists at `Worktree` and is a git worktree whose `git branch --show-current` is `Branch`; the session file is at `SessionFile`, mode 0600, and holds the token whose hash is `credential_id`; for claude the MCP config at `MCPConfig` names `Daemon.Addr` and that token, and for codex there is none; the printed launch block names `Agent`, `Model` and `Task`; the request row has `approved_by=tty` and `credential_id` equal to the session's | grants rebuilt from defaults; the worktree on another branch; the token file elsewhere; agent, model or task from flag defaults |
| CL6 | `TestApproveWithoutATerminalMintsNothing` | `noTerminalEnv` → refused with `sessionTTYFix`; the request is still pending; no worktree | the fence bypassed on the request path |
| CL7 | `TestApproveDeclinedLeavesTheRequestPending` | `n` → exit 2, pending, nothing created | consume before the y/N |
| CL8 | `TestApproveRefusesReplay` | a second approve → "already used … by tty", with **no terminal opened**; 1 session row | `Admit` step 1 skipped |
| CL9 | `TestApproveRefusesExpired` | arm 1: the clock at `+600` → refused before the terminal opens (spy); arm 2: `now` returns T and then T+601 → refused at the write, and the worktree this run created is removed | arm 2: the consume reuses the pre-prompt `now` |
| CL10 | `TestApproveRefusesAStaleWorld` | profile bytes changed → refusal names `profile.digest`; daemon build changed → `daemon.…`; session file now present → refused; the request stays pending in each | the re-derivation equality skipped |
| CL11 | `TestApproveRefusesATamperedRow` | SQL edits the stored bytes (budget 1→999) and keeps the digest → `RefuseDigestMismatch`; edit bytes and digest together → the prompt shows `999`, and the minted grant is the one shown | the digest recomputation skipped |
| CL12 | `TestApproveRefusesBadIDs` | a bad grammar → exit 2; an unknown id → exit 1; the global `--addr` → refused | — |
| CL13 [r2] | `TestApproveUnansweredPromptTimesOut` | **[r3, gemini]** `t.Cleanup` closes the mock reader pipe to reap the background goroutine before test completion; a spy terminal whose `Read` blocks until test cleanup; `now` returns `expires_at-37`; `after` records its argument and returns a closed channel. `runApprove` returns within a 5 s guard, exit 2, and stderr names "no answer within"; the recorded duration is exactly 37 s; the request is still pending (`consumed_at` NULL); no worktree and no session row. Arm 2: `now == expires_at` at the prompt → refused, and the prompt is never written | `approve` using the plain `readTermLine` (the guard fires); a bound taken from `RequestLifetime` or `TTLSeconds` rather than what remains (≠ 37 s); arm 2's pre-prompt check deleted |
| CL14 [r2] | `TestMintRequestFieldMapIsComplete` | reflects over `mintOptions` and `mintRequest` (flattened) and `MintRequest`: every field is mapped or excluded exactly once, every target exists, and every `MintRequest` leaf is a target or request-only | a new field added to `mintRequest` (or `MintRequest`) without a map entry; a map entry deleted |
| CL15 [r2] | `TestApproveRebuildsTheSameMintRequest` | `episode start`'s `mintRequest` and `approve`'s `fromRequest(Decode(Canonical(…)))` are equal on every non-func field; func fields have the same nil-ness; the two `afterToken`s write byte-identical MCP configs | `fromRequest` forgetting any one field (for example `profilesDir`, `git` or `wsRoot`), which byte equality alone cannot see |
| CL16 [r2] | `TestApproveEffectsCheckCatchesEachMismatch` | the test-side checker `requireApproveEffects` against facts from a good approve with **one** fact altered per arm (worktree path, branch, session file path, its mode, token hash, MCP addr, MCP present for codex, agent, model, task, `approved_by`, the request row's `credential_id`): each arm red, the clean facts green | the checker silently skipping a field, which would make CL5 vacuous for it |
| CL17 [r2] | `TestExistingCallersReadUnbounded` | `session new` and `episode start` with an `after` seam that fails the test if called, and a terminal answering `y`: both mint as before | the bound applied to every caller (a behaviour change outside scope) |
| CL18 [r2] | `TestReadTermLineWithin` | a reader that returns `"y\n"` before the timer → `("y", true)`; a blocked reader with a fired timer → `("", false)` without waiting for the reader | `select` on the read only; a timer that is never selected |
| CL19 [r3] | `TestApproveRefusesExpiryDuringWorktreePreparation` | **[r3, gpt verbatim]** answers y before expiry, advances the injected clock during worktree preparation to `expires_at`, and requires `ErrRequestExpired`, zero session rows, an unconsumed request, and rollback of the newly created worktree | the consume reading taken before worktree preparation (r2's shape); `MintSessionForRequest` using a caller-supplied timestamp instead of `Clock()` inside the transaction |
| edited | `TestEveryVerbHasHelp`, `everyVerb` (`help_test.go:10`, V11) | `approve` added; usage contains `approve <request-id>` | — |
| edited | `TestProductionContextRoots` | +2 roots (§7) | — |

**Must stay green, unedited.** These are the typed-`y` path's existing tests (V24):
`TestSessionMint_*` (`session_test.go`), `TestSessionNew*` (`session_new_test.go`),
`TestSessionNewWithoutAControllingTerminalRefusesEndToEnd` and
`TestSessionNewReadsTheAnswerFromThePtyWhenStdinIsDevNull` (`session_pty_test.go`),
`TestEpisodeStart*` (`episode_test.go`), `TestStartLaunch*`, `TestEpisodeResume*`,
`TestSessionMintGrantsExpireWithTheSession`, every `TestCredentials*`, `TestMigrateCredentials*`
and `TestPreSplitStoreOpenModes`, `TestDaemonDependencyAllowlist`, and the quickstart tests.

## 7. Conflict surface

**[r2]** Sprint 1 (M1–M3) touches only the first four rows plus `.github/workflows/ci.yml`. Every
other row belongs to sprint 2 (M4–M6).

| File | Change |
|---|---|
| `host/approval/request.go`, `admit.go`, `sshsig.go` | **new** (pure) |
| `host/approval/request_test.go`, `admit_test.go`, `sshsig_test.go`, `sshsig_crosscheck_test.go` | **new** |
| `host/approval/testdata/sshsig/*` (8 files, §2.8) | **new**, generated |
| `.gitattributes` | **new**: one line |
| `host/store/credentials_schema.sql` | `mint_requests` table |
| `host/store/credentials.go` | version 4; the table-set check at `:304`; the upgrade function generalised to reach 4 |
| `host/store/mint_requests.go` (+ test) | **new**: 4 methods, typed errors, `OlderCredentialsStoreError` |
| `host/store/credentials_test.go`, `credentials_binding_test.go` | canonical DDL, table names, version expectations |
| `host/store/context_roots_test.go` | +`"cmd/ailang-worldd/mint_request.go\|createMintRequest\|Background": 1` and `"cmd/ailang-worldd/approve.go\|runApprove\|Background": 1`. **These are the only new production roots.** `approve`'s mint runs under `provisionEpisode`'s existing root (`:444`). No `TODO()` or `WithoutCancel` is added |
| `host/daemon/credentials.go` | the upgrade log line text only |
| `host/authority/mint.go` | `MintForRequest`; `mintRecorded` takes an optional consume. The existing signatures are unchanged |
| `cmd/ailang-worldd/session.go` | `sessionEnv.random` and **[r2]** `sessionEnv.after`; `realSessionEnv` sets them (`crypto/rand.Read`, `time.After`); **[r2]** `readTermLineWithin` beside `readTermLine` |
| `cmd/ailang-worldd/session_new.go` | `mintRequest.consume` and **[r2]** `mintRequest.answerBy`; one mint branch; a second `env.now()` for the consume; **[r2]** the bounded read when `answerBy != 0` (step 3 only); the comment's caller list |
| `cmd/ailang-worldd/episode.go` | `--request` flag and help; the branch to `createMintRequest` after `checkDaemon` |
| `cmd/ailang-worldd/mint_request.go` | **new**: `buildMintRequest`, `createMintRequest`, **[r2]** `fromRequest` and `mintRequestFieldMap` (§2.1a), conversions to and from `broker.Capability` |
| `cmd/ailang-worldd/approve.go` (+ `approve_test.go`, **[r2]** `mint_request_test.go` for CL14/CL15) | **new**: `runApprove`, `--list` |
| `cmd/ailang-worldd/main.go` | `case "approve"`, usage line, global `--addr` refused |
| `cmd/ailang-worldd/help_test.go` | `everyVerb` and the usage want-list |
| `.github/workflows/ci.yml` | one new step (§2.8) |
| `docs/QUICKSTART.md` | **new §14 "Mint by request"** at the end of the file, *attended — pending first verbatim run* (S7). It shows `episode start --request` then `approve <id>`, against a test store. #515 (row 256) also edits this file but not its end (V25) |

**Not touched**: `world/*.ail`, `packages/`, `scripts/verify_ail.sh`, `tools/launchd/*`,
`host/daemon/daemon.go` (no route), `cmd/world-publish/*` (publishing stays typed-only, R4), and
`session revoke`.

## 8. Milestones

Every milestone ends with `go build ./... && go vet ./<pkg>/` (`go build` does not compile test
files) plus its own tests. Each sprint ends with the full gates.

**[r2] Sprint 1 = M1–M3 (this sprint, about 1,450 lines).**

- **M1: request core** (`host/approval` request, admit; about 250 production + 400 test lines,
  including r2's seven-column `MintMatches`).
  RQ1–RQ8. `go test ./host/approval/ -run '^(TestCanonical|TestDecode|TestDigest|TestValidate|TestSignedMessage|TestAdmit|TestMintMatches|TestDiff)' -count=1`
- **M2: SSHSIG verifier and fixtures** (about 280 + 380 lines, plus fixtures, `generate.sh`,
  `.gitattributes`). SG1–SG6. `go test ./host/approval/ -run '^TestSSHSIG|^TestParseSSHSIG|^TestVerifyApproval|^TestFingerprint' -count=1`
- **M3: cross-check and CI step** (about 150 test lines, 20 YAML lines). XC1.
  `WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -run '^TestSSHSIGCrossCheckWithSSHKeygen$' -count=1 -v`
  must print `--- PASS` for both subtests.
- **Sprint 1 gates**: `./scripts/verify_ail.sh`, `go build ./... && go test ./...`, and
  `go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1` (`host/approval` is
  not yet imported by the daemon core; this runs as a guard for sprint 2).

**[r2] Sprint 2 = M4–M6 (the next sprint, same doc, about 1,600 lines).**

- **M4: store schema 4** (about 200 + 330 lines). ST1–ST7.
  `go test ./host/store/ -run 'Credentials|MintRequest|MintSessionForRequest|ConcurrentApprovals|PutMintRequest|TestProductionContextRoots' -count=1`
  and `go test -race ./host/store/ -run '^TestConcurrentApprovalsMintOnce$' -count=1`
- **M5: create** (`episode start --request`, `approve --list`; about 150 + 220 lines). CL1–CL4.
  `go test ./cmd/ailang-worldd/ -run '^(TestEpisodeStartRequest|TestApproveList)' -count=1`
- **M6: approve by typed `y`** (about 230 + 520 lines, plus QUICKSTART §14). CL5–CL18 and the
  edited help and context-root tests.
  `go test ./cmd/ailang-worldd/ ./host/authority/ -run '^(TestApprove|TestMintRequestFieldMap|TestExistingCallersReadUnbounded|TestReadTermLineWithin|TestEveryVerbHasHelp|TestSessionMint_|TestSessionNew|TestEpisodeStart)' -count=1`
  and `go test ./host/store/ -run '^TestProductionContextRoots$' -count=1`
- **Sprint 2 gates**: `./scripts/verify_ail.sh` (no `.ail` change, so it must stay exactly
  green), `go build ./... && go test ./...`, and
  `go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1`.

Nothing in M1–M3 depends on M4–M6. Sprint 1 ships no behaviour change: no verb, route or schema
changes until sprint 2.

## 9. Residuals and risks

- **R1** The CI image's `ssh-keygen` is **UNVERIFIED** (§2.8). The step fails loudly with the fix.
- **R2** A0 adds a step for the operator (two commands where there was one) and gives no UX gain
  until A1b or A3. That is intended (§0).
- **R3** Schema 4 makes an older daemon refuse the store. The same was true of versions 2 and 3,
  and `episode start` already forces a same-build restart (V4, V6).
- **R4** The 16-pending cap and the one-day prune are judgements. No measurement chose them.
- **R5** `--request` does not cover launch presets, `--launch` or `--resume` (CL2). Mark's
  `MINT_BLOCK_N.sh` episodes use presets (parent §1), so the UX row 201 is for needs those paths
  in A1b or A3.
- **R6** In A1b, SSHSIG approval proves possession of a key, not the presence of a human (§4).
- **R7 [r2]** **Pre-existing, logged, not fixed:** the y/N of `session new`, `session mint`,
  `episode start`, launch presets and resume has no time bound (`session_new.go:438` reads before
  the 75 s root at `:444`; `session.go:275` likewise; V33). Only `approve` gains the bound
  (§2.6a). Bounding the others is a separate row, because it changes attended behaviour Mark
  relies on today.

## 10. Open questions for Mark (each with a recommended default)

- **Q1** `sk-*` (security-key) signatures: refuse them in A1 until an attended fixture from a real
  security key exists (a public key and signature, no secret), or accept them unverified by a
  fixture? **Default: refuse.** A1b can add them with Mark's attended fixture.
- **Q2** Move the admission law (`Admit`) into the `.ail` kernel with Z3 contracts when A4 brings
  `Human.Approve` packets under the same law? **Default: no for A0/A1 (§3). Revisit at A4.**
- **Q3** Creating a mint request needs no terminal, because it grants nothing and is capped at 16
  pending; only approving needs the `y`. **Default: yes.** This matches parent §2.4/§3, where an
  agent may later *propose* a request.

## 11. What this is NOT doing

- No signature-based mint, enrolment, `approvers` table, or key attribution in the record (A1b).
- No WebAuthn (A2), no surface, notification or Tauri change, no `/v1` route (A3), no
  `Human.Approve` packets (A4).
- No change to `session new`, `session mint`, `session revoke`, launch presets, `--launch` or
  `--resume`. No change to publishing (R4).
- No `.ail`, kernel, package or verify-gate change. No new Go module.
- No real key, no real store and no publish, by the loop or by its tests (D-WORLD-109).

## 12. Verification Log

Commands ran in the worktree at `693f704e` unless they say otherwise. An empty result is paired
with a positive control in the same command.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and instruments | `git log --oneline -1`; `shasum -a 256 cmd/ailang-worldd/{session,session_new,episode}.go host/store/credentials.go host/authority/mint.go \| cut -c1-16`; `go version`; `ssh -V`; `$HOME/.pinned-ailang/ailang --version` | `693f704e Merge pull request #530 …`; `3566d5377c9fc01a`, `ef5163a5dc6e7aa5`, `5588c8659f51669b`, `3b4d630f89676737`, `1dfdf1080000889f`; `go1.26.6 darwin/arm64`; `OpenSSH_10.3p1, LibreSSL 3.3.6`; `AILANG v0.41.0` |
| V2 | One fence-and-mint path; its steps; revoke is unfenced | `cat -n cmd/ailang-worldd/session_new.go \| sed -n 1,520p`; `cat -n cmd/ailang-worldd/session.go`; `grep -n "No terminal fence is required" cmd/ailang-worldd/session.go` | `:24-26` "Steps 2-7 are provisionEpisode, the ONE fence-and-mint path … no other code creates a credential"; `:396 func provisionEpisode`; `:399 term, err := env.openTerminal()`; `:408 now := env.now()`; `:430-431 if r.presetPrompt != "" { fmt.Fprint(term, r.presetPrompt)`; `:436-438` the y/N and `isYes(readTermLine(term))`; `:444 ctx, cancel := context.WithTimeout(context.Background(), sessionNewBudget)`; `:481 MintStarted`, `:483 sessionMint(mintCtx, …)`; `:503-518` revoke on failure; `session.go:340 // credential. No terminal fence is required — revocation is an ordinary` |
| V3 | The seam has no randomness | `sed -n 20,47p cmd/ailang-worldd/session.go` | `:24 type sessionEnv struct {`, fields `openTerminal func() (io.ReadWriteCloser, error)` and `now func() int64` only; `:43 return os.OpenFile("/dev/tty", os.O_RDWR, 0)` |
| V4 | Plain `episode start` validates, checks the daemon build, mints unrecorded | `sed -n 185,360p cmd/ailang-worldd/episode.go`; `sed -n 419,485p` same file; `grep -n "req.start = " cmd/ailang-worldd/episode_preset.go` | `:305 req, err := prepareMint(mintOptions{verb: "episode start", …` (no `start`); `:343 if code := provisionEpisode(req, …`; `:431 func checkDaemon(…)`, `:432 if !daemon.SameBuild(self, h.Build) {`; `episode_preset.go:240 req.start = &store.EpisodeStart{…` |
| V5 | Credentials schema 3, its table check, deadlines | `grep -n "const credentialsSchemaVersion\|case known && ((version\|func (c \*Credentials) mintSession\|requireDeadline(ctx); err" host/store/credentials.go`; `cat host/store/credentials_schema.sql` | `:40 const credentialsSchemaVersion = 3`; `:304 case known && ((version < 3 && len(tables) == 1 …) \|\| (version == 3 && len(tables) == 2 && tables[0] == "episode_starts" …))`; `:402 func (c *Credentials) mintSession`; `:403 if err := requireDeadline(ctx); err != nil {`; schema has `session_credentials` and `episode_starts` only |
| V6 | Only serve upgrades; the daemon never keeps the writer lock; CLIs refuse older stores | `sed -n 40,85p host/daemon/credentials.go`; `grep -n "Version() < store.CredentialsSchemaVersion\|version < store.CredentialsSchemaVersion" cmd/ailang-worldd/*.go` | `:59 func ensureCredentialsStore(…)`; comment "so the daemon never keeps the credentials writer lock. Only serve upgrades"; log "upgraded the credentials store %s to schema %d (exec profile binding; …)"; `session_new.go:372`, `episode_preset.go:235` |
| V7 | Mint core | `sed -n 20,130p host/authority/mint.go` | `:43 func MintBound`, `:48 func MintStarted`, `:57 func mintRecorded(…)`; `rand.Read(raw)`; `clampGrantExpiry`; `st.MintSessionWithStart(ctx, row, start)` |
| V8 | No approval code and no `x/crypto` (with controls) | `grep -rniE "sshsig\|world-approval\|approver\|mint_request\|MintRequest" --include="*.go" host cmd \| head`; control `grep -rn "func provisionEpisode" --include="*.go" cmd`; `grep -rln "golang.org/x/crypto" --include="*.go" host cmd \| wc -l`; control `grep -rln '"modernc.org/sqlite"' --include="*.go" host cmd \| head -2`; `cat go.mod`; `ls host/approval cmd/ailang-worldd/approve*.go` | only `host/coordinator/decisionpacket_producer_test.go` matches (`approveReq`, `approveRigWith`: a `Human.Approve` test rig, not approvals by key); control: `session_new.go:396`; `0`; control: `host/broker/allowlist_test.go`, `host/daemon/integrity_test.go`; `go.mod` requires `ailang v0.47.2`, `x/text`, `modernc.org/sqlite` only; both `No such file or directory` |
| V9 | The dependency allowlist | `sed -n 757,800p host/daemon/daemon_test.go`; `grep -n "func TestDaemonDependencyAllowlist" host/daemon/daemon_test.go` | `:775 var allowedDepModules = []string{` (this repo, sqlite and its indirect deps, `ailang/serveapi/protocol`); `:793 var daemonCorePatterns = []string{"./host/daemon/...", "./cmd/ailang-worldd/..."}`; `:899 func TestDaemonDependencyAllowlist` |
| V10 | Context roots inventory | `sed -n 15,40p host/store/context_roots_test.go` | `:26 "cmd/ailang-worldd/session.go\|runSessionMint\|Background": 1`, `:29 "cmd/ailang-worldd/session_new.go\|provisionEpisode\|Background": 1`; no `approve.go` or `mint_request.go` entry |
| V11 | Verb list and help pin | `grep -n "case \"" cmd/ailang-worldd/main.go`; `sed -n 1,60p cmd/ailang-worldd/help_test.go` | no `case "approve"` (control: `case "session"` at `:429`); `:10 var everyVerb = [][]string{`; `TestEveryVerbHasHelp` checks `--help` and a usage want-list |
| V12 | Real `ssh-keygen` signatures; their shape; the file-mode trap | in `/tmp/world-iter373/designer-scratch`: `ssh-keygen -q -t ed25519 -N "" -f k_ed25519`; `ssh-keygen -q -t ecdsa -b 256 -N "" -f k_p256`; `ssh-keygen -Y sign -n world-approval@sunholo -f k_ed25519 msg.bin`, then the same with `k_p256` (file mode); then from stdin `… < msg.bin > ed25519.sig` and `> p256.sig`; `awk '{print length}' p256.sig` | file mode, second call: `msg.bin.sig already exists.` `Overwrite (y/n)?` then `sign p256 rc=0`, and the existing `.sig` still decodes as `ssh-ed25519` (not written); stdin: `Signing data on standard input`, rc 0 both; armor lines at most 70 characters; the decoded Ed25519 blob shows reserved `AAAAAA` (empty) and hash `c2hhNTEy` = `sha512` |
| V13 | `-Y sign` semantics; the stdlib API | `man ssh-keygen \| col -b \| sed -n 430,450p`; `… sed -n 258,290p`; `go doc crypto/ecdsa.ParseUncompressedPublicKey` | "Signatures are written to … standard output if the message to be signed was read from standard input"; "-f … may refer to either a private key, or a public key with the private half available via ssh-agent(1)"; "hashalg … Valid algorithms are "sha256" and "sha512." The default is "sha512.""; `func ParseUncompressedPublicKey(curve elliptic.Curve, data []byte) (*PublicKey, error)` |
| V14 | CI runners; no `ssh-keygen` in CI (with control); the PASS-grep pattern | `grep -c ssh-keygen .github/workflows/ci.yml`; control `grep -n runs-on .github/workflows/ci.yml`; `sed -n 225,320p .github/workflows/ci.yml` | `0`; control: `:20 ubuntu-latest`, `:121 ubuntu-24.04`, `:423 ubuntu-24.04`; `:238` the "Row 140 … verbose" step, whose `for t in $tests; do grep -q -- "^--- PASS: $t " "$log" \|\| { echo "::error::$t did not PASS on linux (SKIP?)"; exit 1; }`. Whether the image ships `ssh-keygen`: not checkable offline (**UNVERIFIED**) |
| V15 | `sk-*` keys cannot be made without hardware | `ssh-keygen -q -t ed25519-sk -N "" -f k_sk </dev/null`; `ls k_sk*` | `No FIDO SecurityKeyProvider specified`, `Key enrollment failed: invalid format`, rc 255; `No such file or directory` |
| V16 | `ssh-keygen -Y verify` arms | `allowed` holds both keys with `namespaces="world-approval@sunholo"`; `ssh-keygen -Y verify -f allowed -I principal -n world-approval@sunholo -s {ed25519,p256}.sig < msg.bin`; then a tampered message, `-n git`, a `-n git` signature and a third key's signature | `Good "world-approval@sunholo" signature … ED25519` rc 0; `… ECDSA` rc 0; `incorrect signature` rc 255; `namespace does not match` rc 255 (twice); `Could not verify signature.` rc 255 |
| V17 | A standard-library-only verifier agrees on every arm | `GO111MODULE=off go run /tmp/world-iter373/designer-proto/main.go /tmp/world-iter373/designer-scratch` (throwaway prototype, about 120 lines; imports only `crypto/ecdsa`, `crypto/ed25519`, `crypto/elliptic`, `crypto/sha256`, `crypto/sha512`, `encoding/base64`, `encoding/binary`, `math/big` and basics) | `ed25519 good -> <nil>`; `p256 good -> <nil>`; `ed25519 tampered msg -> ed25519 bad`; `git-namespace sig -> namespace "git"`; `unenrolled key -> not the enrolled key`; `p256 sig vs ed key -> not the enrolled key` |
| V18 | The generator in §2.8 works and leaves no private key | the generator copied to `/tmp/world-iter373/designer-gen-src/generate.sh`; run in a fresh `mktemp -d` tree holding a stub `host/approval/sshsig_test.go` with a 90-byte `fixtureMessageHex`; then `ls`; `wc -c < message.bin`; `ssh-keygen -Y verify` on each `.sig` | `gen rc=0`; `SOURCE.txt ed25519.pub ed25519.sig generate.sh git-ns-ed25519.sig message.bin p256.pub p256.sig`; `90`; `SOURCE.txt` = `OpenSSH_10.3p1 …`, the time, `256 SHA256:nI4P… world-approval-fixture-ed25519 (ED25519)`, `256 SHA256:X+Nl… (ECDSA)`; `ed25519 verify rc=0`, `p256 verify rc=0`, `git-ns refused rc=255`; `private keys left in fixture dir: 0` |
| V19 | No `.gitattributes`; autocrlf unset (with control) | `ls -la .gitattributes`; control `ls -la go.mod`; `git config --get core.autocrlf; echo rc=$?` | `No such file or directory`; control lists `go.mod`; `autocrlf rc=1` (unset) |
| V20 | Pinned AILANG has no Ed25519, ECDSA or SHA-512 | `$HOME/.pinned-ailang/ailang builtins list \| grep -E "std/crypto"`; `… \| grep -iE "ed25519\|ecdsa\|p256" \| wc -l`; control `… \| grep -ciE "string\|list"` | `_crypto_constanttimeequal`, `_crypto_hmacsha256`, `_crypto_rsa_verify_pkcs1v15`, `_crypto_sha256bytes`, `_crypto_sha256hex`; `0`; control `66` |
| V21 | Kernel modules and the floor constants (S8 cost) | `ls world/`; `grep -n "^REQUIRED_VERIFIED\|^EXACT_TOTAL_VERIFIED" scripts/verify_ail.sh` | `contracts.ail logepoch.ail transitions.ail types.ail`; `:323 REQUIRED_VERIFIED = {`, `:414 EXACT_TOTAL_VERIFIED=16` |
| V22 | The exec sandbox denies reads of the operator's home | `sed -n 95,125p host/broker/exec_settings.go` | `:113 DenyRead:   []string{spec.OperatorHome, spec.StateDir, spec.WorkspaceRoot},` |
| V23 | No test pins the upgrade log text or schema 3's table list outside the store (with control) | `grep -rln "upgraded the credentials store" --include="*_test.go" host cmd`; `grep -rn '"episode_starts", "session_credentials"' --include="*_test.go" host cmd` | first: no output; second: only `host/store/credentials_test.go:39` (the control, and an edit site) |
| V24 | The typed-`y` tests that must stay green exist | `grep -n "^func Test" cmd/ailang-worldd/{session_new,session_pty,episode}_test.go`; `grep -n "^func Test.*Mint" cmd/ailang-worldd/*_test.go` | `TestSessionNewWithoutATerminalProvisionsNothing`, `TestSessionNewDeclinedProvisionsNothing`, `TestSessionNewMintFailureRollsBackTheWorktree`, `TestSessionNewWithoutAControllingTerminalRefusesEndToEnd`, `TestEpisodeStartWithoutATerminalProvisionsNothing`, `TestEpisodeStartRefusesAMismatchedDaemonBeforeTheTerminal`, `TestSessionMint_RequiresControllingTerminal`, `TestSessionMint_ConfirmationRejected`, `TestSessionMintGrantsExpireWithTheSession`, … |
| V25 | In-flight overlap: #515's branch (fetched locally) | `git diff --name-only origin/dev...origin/sprint/world-iter362-row256 -- host cmd docs` | touches `docs/QUICKSTART.md`, `doctor*.go`, `exec_probe.go`, `exec_startup.go`, `daemon.go`, `workspace.go` (+ tests); of this design's files, only `QUICKSTART.md` (§14 is appended at the end). Open PRs on GitHub were **not** checked (no network); the controller checks before the build |
| V26 | `Human.Approve` exists as the A4 seam | `grep -rn "EffectHumanApprove\s*=" host/broker/*.go` | `host/broker/approve.go:18 EffectHumanApprove = "Human.Approve"` |
| V27 | Scratch cleaned, tree clean (r1; re-checked at r2) | `rm -rf` of `/tmp/world-iter373/designer-{scratch,proto,gen-src,gen.*}` (r1) and `designer-r2` (r2), then `ls /tmp/world-iter373/`; `git status --short` | only `designer-directive.md` (and the controller's own files) left; only this doc listed |
| V28 [r2] | The two mint structs at `693f704e` (controller-verified too) | `sed -n 180,224p cmd/ailang-worldd/session_new.go \| grep -vE '^\s*//'` | `type mintOptions struct { verb string; episode string; db string; credsDB string; wsRoot string; repo string; branch string; allowNonGit bool; grants []broker.Capability; ttl int64; outPath string; outDirCreate bool; profile string; profilesDir string }`; `type mintRequest struct { mintOptions; credPath, dir string; reuse, degraded bool; git string; prof *broker.ExecProfile; profFile string; afterToken func(token string) error; resume *store.EpisodeResume; beforeToken func() error; resumeRollback func() error; start *store.EpisodeStart; startPoint string; existingBranch bool; presetPrompt string }`: 14 + 15 fields, every one in the §2.1a table |
| V29 [r2] | Who sets the path-specific fields | `grep -n "\.startPoint\b\|\.existingBranch\b\|presetPrompt\s*=\|\.afterToken =\|\.beforeToken =\|\.resumeRollback =\|\.resume =\|req\.start =\|outDirCreate:" cmd/ailang-worldd/*.go \| grep -v _test` | `episode.go:307 … outDirCreate: true …` (no `allowNonGit` in that call, `:305-307`); `episode.go:339 req.afterToken = func(token string) error { return writeMCPConfig(mcpFile, hostPort, token) }`; `episode_preset.go:217` (`outDirCreate: true`), `:239 req.startPoint, req.existingBranch = p.BaseSHA, existing`, `:240 req.start = …`, `:268 req.presetPrompt = …`, `:276 req.afterToken = …`; `episode_resume.go:188 req.presetPrompt`, `:190 req.resume`, `:192 req.resumeRollback`, `:201 req.beforeToken`; readers `session_new.go:421`, `:433`, `:457`, `:460` |
| V30 [r2] | `session_credentials` has seven columns and no `source`/`created_by` (with control) | `sed -n "/CREATE TABLE IF NOT EXISTS session_credentials/,/);/p" host/store/credentials_schema.sql`; `grep -c "source\|created_by" host/store/credentials_schema.sql`; control `grep -c expires_at` same file | `credential_id`, `episode_id`, `grants_json`, `expires_at`, `created_at`, `project`, `profile_digest`; `0`; control `1` |
| V31 [r2] | On a pty slave, a read deadline works on darwin | throwaway `/tmp/world-iter373/designer-r2/ttyprobe.go` (opens the path `O_RDWR`, `SetReadDeadline(+200ms)`, reads a line in a goroutine, `select` with a 500 ms timer, then `Close`), built with `go build`; driver `python3 -I drive.py` passing `os.ttyname(slave)` of `pty.openpty()` | `slave: /dev/ttys081 rc: 0 wall: 0.45s`; `SetReadDeadline: <nil>`; `read returned after 200ms: "" read /dev/ttys081: i/o timeout` |
| V32 [r2] | On `/dev/tty` as the controlling terminal, darwin supports no deadline; `Close` does not unblock the read; timer + `select` + exit does | the same probe with argument `/dev/tty`, run as the child of `pty.fork()` (so the pty is its controlling terminal) by `python3 -I drive_ctty.py` | `SetReadDeadline: file type does not support deadline`; `timer fired after 500ms; read still blocked`; `Close: <nil> took 0s`; `read still blocked 500ms after Close`; `main returns (process exits)`; `child status 0 wall 1.02s`. Linux: not measured |
| V33 [r2] | The existing y/N reads have no deadline; the roots start after them (controller-verified too) | `sed -n 396,445p cmd/ailang-worldd/session_new.go`; `grep -n "readTermLine\|context.Background" cmd/ailang-worldd/session.go` | `:399 term, err := env.openTerminal()`; `:436-437` the prompt; `:438 if !isYes(readTermLine(term)) {`; `:444 ctx, cancel := context.WithTimeout(context.Background(), sessionNewBudget)`; `session.go:275 if !isYes(readTermLine(term)) {`, `:282 mintCtx, cancelMint := context.WithTimeout(context.Background(), 3*time.Second)` |
| V34 [r3] | `broker.Capability` fields and the `grants_json` keys (controller-measured) | `grep -n -A6 'type Capability struct' host/broker/*.go`; `grep -n 'json.Marshal' host/authority/mint.go`; `sed -n 426,434p host/store/credentials.go` | `decide.go:15 type Capability struct {` / `Effect string` / `Scope string` / `ExpiresAt int64` / `Budget int64` / `}` — no json tags, so the keys are `Effect`, `Scope`, `ExpiresAt`, `Budget`; `mint.go:83 grantsJSON, err := json.Marshal(grants)`; `credentials.go:428 INSERT INTO session_credentials (credential_id, episode_id, grants_json, expires_at, created_at)` (bound form adds `project, profile_digest`) |

## 13. Quorum verification log

- **Round 1** (iteration 373, r1; seats `oc-kimi-k3`, `gpt6-1-sol`, `gemini-4-argon`; $0.355):
  **BLOCKED 2/3.** `gemini-4-argon` PASS, noting the ~2,800-line size. `oc-kimi-k3` and
  `gpt6-1-sol` REJECT.
  - **kimi**: "everything the approval executes is inside the request" was asserted but never
    verified against `mintOptions`/`mintRequest`. Byte-equality re-derivation cannot see a dropped
    field, and `SessionMatches` was too narrow. **Answered by** §2.1a (field map and four new
    `MintRequest` fields: `WorldDB`, `WorkspaceRoot`, `Git`, `ProfilesDir`; `Repo` made
    absolute), §2.3 (`MintMatches` over all seven columns), §6 CL5 (widened), CL14, CL15, CL16,
    RQ7 (widened), and V28–V30. One premise in the objection does not hold: `session_credentials`
    has **no** `source` or `created_by` column (V30). The seven real columns are all checked.
  - **gpt**: `approve` reused `provisionEpisode`, whose y/N read has no deadline, and expiry was
    checked only after the read. **Answered by** §2.6a (`answerBy`, `sessionEnv.after`,
    `readTermLineWithin`), CL13, CL17, CL18, R7, and V31–V33. The mechanism uses no tty read
    deadline, because `/dev/tty` on darwin supports none (V32).
  - **gemini (size)**: **answered by** the sprint boundary in §0 and §8: this sprint is M1–M3,
    and the next is M4–M6.
- **r2** = the designer's revision of the same file, applying exactly these fixes.
- **Round 2** (iteration 373, r2; same seats; $0.429): **BLOCKED 2/3.** `gemini-4-argon` PASS
  (non-blocking: a timed-out `readTermLineWithin` leaves its goroutine blocked; reap it in tests).
  `oc-kimi-k3` REJECT: the field map's `grants` leaf was asserted, `broker.Capability` never dumped,
  and CL14 could not see the element type; also the `db` row conflated `WorldDB`/`CredentialsDB`
  and SG1's pinned hashes vs regeneration was unstated. `gpt6-1-sol` REJECT: `consume.Now` was
  sampled before worktree preparation, so the SQL expiry check used a stale time.
  Surfaces per round (rule (a)): r1 = field completeness (cmd) + prompt wait (cmd); r2 = Capability
  element completeness (cmd/broker) + consume-time clock (store). Neither r2 objection disputes the
  direction; both carry a concrete `proposed_fix`.
- **r3** = **controller narrow-refinement carve-out** (iteration 373): the reviewers' r2
  `proposed_fix` text applied verbatim, nothing controller-invented beyond V34's measured output.
  kimi → §2.1a `db` and `grants` rows, CL14 bullet, §2.8/SG1 regeneration sentence, V34. gpt →
  §2.4 `RequestConsume.Clock` (one reading inside the transaction, after `BeginTx`), §2.6 step 4,
  §2.3 `created_at` bound, CL19. gemini → the `readTermLineWithin` comment constraint and CL13's
  `t.Cleanup`. Routed to the planner without a third quorum, per the carve-out.
