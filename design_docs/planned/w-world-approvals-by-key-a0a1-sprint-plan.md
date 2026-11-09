# Sprint plan — w-world-approvals-by-key-a0a1, sprint 1 (iteration 373, row 201)

- **Design**: [w-world-approvals-by-key-a0a1.md](w-world-approvals-by-key-a0a1.md), **r3**. This
  sprint is the design's **M1–M3 only**: the pure request core, the SSHSIG verifier with real
  `ssh-keygen` fixtures, and the `ssh-keygen` cross-check with its CI step. **M4–M6 (store schema 4,
  `episode start --request`, `approve`) are the next sprint. Do not start them.** Sprint 1 ships no
  behaviour change: no verb, route, schema or `.ail` change.
- **Planner base**: `63766f2ffee52ec8d39a8e3f6e62801d875f0a57` (= `693f704e` + the r3 design doc),
  worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter373-row201`, branch
  `sprint/world-iter373-row201`.
- **Executor**: `claude-haiku-5-5`. Every name, signature, constant, string and assertion below is
  fixed. Where this plan and the design differ, this plan was checked against the tree and wins. The
  differences are listed in §1.3, and the two design premises that turned out false are in §1.4.
- **Estimate**: about **1.25 days** of executor time over five milestones, about **1,770 lines**
  (≈580 Go production, ≈1,150 Go test, ≈20 YAML, ≈22 script/attributes, plus generated fixtures).
  The design says ~1,450. The extra lines come from AP1, SG7 and the named fixture transforms (§1.3).
- **Touched**: only `host/approval/**` (new), `.gitattributes` (new, one line) and
  `.github/workflows/ci.yml` (one step). **Not touched**: `cmd/**`, `host/store/**`,
  `host/authority/**`, `host/daemon/**`, `host/broker/**`, `docs/QUICKSTART.md`, `scripts/**`,
  `world/*.ail`, `packages/**`, `tools/launchd/*`, any skill, `~/.ailang/state/**`.

---

## 0. Binding constraints (each is restated in the milestone it governs)

1. **Test keys only (D-WORLD-109 = A).** Every key is made inside a test (`t.TempDir()`) or by
   `generate.sh` in a `mktemp -d` directory that it deletes on exit. No private key is ever checked
   in or written under the repo. Nothing in this sprint opens `~/.ailang/state`, a real credentials
   store, a real `ssh-agent` (XC1 sets `SSH_AUTH_SOCK=` empty) or `~/.ssh`. SG1 asserts the fixture
   directory holds **exactly** the eight pinned names, so a stray key file is red.
2. **Real, pinned, regenerable fixtures.** The fixtures are made by real `ssh-keygen -Y sign`
   through `generate.sh`, **signing from stdin**. They are checked in exactly as produced. The three
   `-Y sign` command lines are quoted in a comment directly above the SG1 pinned-hash table
   (`fixtureSHA256`). SG1 pins the SHA-256 of every file in the directory. **Whoever runs
   `generate.sh` must update `fixtureSHA256` and `SOURCE.txt` in the same commit.** The parser fails
   loudly with a typed error and a nil result on any non-empty input it cannot fully parse. Every
   refusal arm is derived from a real fixture by a **named transform function** (§2.4). None is
   typed by hand.
3. **Standard library only in `host/approval`.** Production files may import only standard-library
   packages, and none of `os`, `os/exec`, `os/signal`, `net`, `net/http`, `syscall`, `io/ioutil`,
   `crypto/rand`, `math/rand`, `math/rand/v2` or `unsafe` (no I/O, no randomness: coding standard
   S2). Test AP1 enforces this. `ecdsa.ParseUncompressedPublicKey` **exists** in go1.26.6
   (`go doc crypto/ecdsa.ParseUncompressedPublicKey`, §1.2 C1), so the planned API is used as is.
4. **Every test row names the mutant it kills.** For refusal branches the unit of mutation is the
   branch. §4 lists 68 mutants. The controller and the judge run them.
5. **Acceptance commands** per milestone (§3) and sprint gates (§6).
6. **CI step**: the design's §2.8 step text goes into `.github/workflows/ci.yml`, job `go-verify`
   (`runs-on: ubuntu-24.04`, `:121`), **after the `go version` step (`:132-135`) and before
   `Install PINNED v0.41.0 ailang` (`:137`)**. Whether `ubuntu-24.04` ships `ssh-keygen` is
   **UNVERIFIED**. The step fails loudly with the `apt-get install openssh-client` hint if it is
   missing.
7. **Baseline** was measured before planning (§1.1). Every gate in §6 is green at the base.

---

## 1. Pre-flight

### 1.1 Baseline at `63766f2f` (measured by the planner, 2026-11-09)

| Command | Result |
|---|---|
| `git status --short` | empty (clean) |
| `go test ./host/approval/ 2>&1 \| head` | `# ./host/approval` / `stat …/host/approval: directory not found` / `FAIL ./host/approval [setup failed]` (no package, as expected) |
| `go build ./...` | rc 0 |
| `go test -run '^$' ./...` (compile fence for test files) | every package `ok` or `no test files`, no failure |
| `go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1` | `ok … 0.392s` |
| `go test ./host/store/ -run '^TestProductionContextRoots$' -count=1` | `ok … 0.268s` |
| `go test ./host/verifygate/ -run '^TestNoParallelWriteForkPackagesOutsideVerifygate$' -count=1` | `ok … 0.544s` |
| `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` (`AILANG v0.41.0`, commit `24ee108`) | rc 0: `verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./...` | rc 0, 35 `ok` packages, 6 m 03 s wall |
| `bash scripts/tracked_binary_gate.sh <verify_go.sh's allowance flags>` | rc 0: `0 other binary blobs among 6552 tracked files` |
| `actionlint .github/workflows/ci.yml` (1.7.12, with shellcheck) | rc 0, no output |
| `ls .gitattributes`; `git config --get core.autocrlf` | absent; unset (rc 1) |
| `go version`; `ssh -V`; `command -v ssh-keygen xxd` | `go1.26.6 darwin/arm64`; `OpenSSH_10.3p1, LibreSSL 3.3.6`; `/usr/bin/ssh-keygen`, `/usr/bin/xxd` |

### 1.2 Codebase and toolchain claims checked

| # | Claim | Evidence |
|---|---|---|
| C1 | `ecdsa.ParseUncompressedPublicKey(curve elliptic.Curve, data []byte) (*PublicKey, error)` exists | `go doc crypto/ecdsa.ParseUncompressedPublicKey` (go1.26.6). `ecdsa.Verify(pub, hash, r, s *big.Int) bool` also exists and is not deprecated |
| C2 | `go.mod` is `go 1.26.6`; module `github.com/sunholo-data/ailang-world` | `head -5 go.mod` |
| C3 | The SSHSIG layout and the stdlib verification recipe are right | Planner scratch run (keys in a `mktemp -d`, deleted): magic `SSHSIG`, version 1, reserved empty, hash `sha512`, no trailing bytes; Ed25519 sig blob 64 bytes and `ed25519.Verify(pk, signedData, sig)` true; P-256 curve `nistp256`, Q 65 bytes, parsed by `ParseUncompressedPublicKey`, sig blob = two mpints (r had a leading `0x00`, 33 bytes), `ecdsa.Verify(pub, sha256(signedData), r, s)` true; armor lines at most 70 chars, one trailing `\n`; `SHA256:`+`base64.RawStdEncoding(sha256(blob))` equals `ssh-keygen -l` for both keys |
| C4 | `ssh-keygen -Y verify -f allowed -I <principal> -n world-approval@sunholo -s <sig> < msg` with an allowed-signers line `<principal> namespaces="world-approval@sunholo" <type> <base64>` gives rc 0 (good) and rc 255 (tampered) | planner scratch run |
| C5 | The plan's `generate.sh` (§2.5) works from any cwd, writes exactly the 8 names, leaves no private key and no NUL byte | dry run in a scratch tree, then deleted: `gen_rc=0`; `SOURCE.txt ed25519.pub ed25519.sig generate.sh git-ns-ed25519.sig message.hex p256.pub p256.sig`; `grep -rl 'PRIVATE KEY'` empty; NUL scan empty; `ssh-keygen -Y check-novalidate` of the p256 sig rc 0 |
| C6 | The golden canonical JSON, digest, signed message and message hash in §2.4 | computed by an independent scratch program with the design's exact struct (§2.1); printed `CANONLEN 1005`, `MSGLEN 90`, `MSGHEXLEN 180`, and `json.Marshal` (HTML escaping on) differs from the canonical bytes, so the `<`/`>`/`&` in `Model` makes RQ1 see the escaping mutant |
| C7 | `scripts/tracked_binary_gate.sh` (run by `verify_go.sh:239`, i.e. CI's `go build + test gate` step) rejects any **committed** blob containing a NUL byte, with no allowance outside PNG/font dirs | header comment `:20-23`; reproduced in a throwaway git repo: a committed `host/approval/testdata/sshsig/message.bin` holding `world-approval/v1\0…` → `FATAL: binary blob(s) rejected: … message.bin  not-allowlisted`, rc 1; a text-only `message.hex` → rc 0 |
| C8 | `TestNoParallelWriteForkPackagesOutsideVerifygate` (`host/verifygate/forklocked_write_test.go:717`) names any `t.Parallel()` in a package whose tests both write files and fork | read `:690-735`. XC1 writes in `t.TempDir()` and execs `ssh-keygen`, so **no `t.Parallel()` anywhere in `host/approval`** |
| C9 | `host/pins` forbids `v0.N.N` release literals in production Go | `host/pins/scan_test.go:22` (`v0\.[0-9]+\.[0-9]+`), `:72-75` (testdata and `_test.go` excluded). The planned code has none (`world/mint-request/v1` and `world-approval/v1` do not match) |
| C10 | `ci.yml` job count must stay 3 | `host/verifygate/mission_config_gate_test.go:329-330` (regex `^  [a-z0-9-]+:$`). The new step is indented 6 spaces and adds no job |
| C11 | `TestDaemonDependencyAllowlist` covers only `./host/daemon/...` and `./cmd/ailang-worldd/...` | `host/daemon/daemon_test.go:793`, `:899`. `host/approval` is imported by neither in sprint 1, so the test is a guard only |
| C12 | `TestProductionContextRoots` walks every non-test `.go` under `host/` and `cmd/` | `host/store/context_roots_test.go:114-119`. `host/approval` must add **no** `context` root (it imports no `context` at all) |
| C13 | The personal-email gate does not scan `host/` | `scripts/check_no_personal_email.sh:11-18`. The fixture key comments are `world-approval-fixture-*` anyway |
| C14 | `broker.Capability` is `{Effect string; Scope string; ExpiresAt int64; Budget int64}` with no json tags; `store.SessionRow` has the 7 columns, `GrantsJSON` a `string` | `host/broker/decide.go:15-20`; `host/store/credentials.go:76-84`; `host/authority/mint.go:83` `json.Marshal(grants)` |
| C15 | Insertion point in `ci.yml` | `:116 go-verify:`, `:121 runs-on: ubuntu-24.04`, `:127-130 Set up Go`, `:132-135 go version`, `:136` blank, `:137 Install PINNED v0.41.0 ailang`. The full step (§3 M3) passed `actionlint` (with shellcheck) in a scratch copy of `ci.yml` |
| C16 | Shell for the controller is zsh | memory note. Every loop or `$VAR` command below is wrapped in `bash -c` |

### 1.3 Planner decisions where the design leaves a choice or cannot be built literally

- **D1 — `message.hex`, not `message.bin` (forced; §1.4 F1).** The signed message is checked in as
  `testdata/sshsig/message.hex`: the 180 lowercase hex characters plus one `\n`. `generate.sh` writes
  the binary message only into its `mktemp -d` directory and signs that. Tests decode the hex.
  Everything the design says of `message.bin` (SG1's equality with `fixtureMessageHex` and with
  `SignedMessage(fixture)`) holds for `message.hex` after hex-decoding. The fixture set is still
  eight files.
- **D2 — The `Verdict` constant is `Admitted`, not `Admit` (forced; §1.4 F2).** A function `Admit`
  and a constant `Admit` cannot share a package. `Admitted = iota + 1`, so the zero `Verdict` is
  not "admit". That makes it fail closed.
- **D3 — `MintFacts.GrantsJSON` is the raw column string.** `MintMatches` decodes it (strictly,
  unknown keys refused) into `[]FactGrant`. This follows the §2.3 table's "decodes" arm, and the
  type mirrors `store.SessionRow.GrantsJSON` (C14). `approval` still imports neither `store` nor
  `broker`.
- **D4 — `MintMatches` returns `*MintMismatchError{Column, Want, Got}`.** RQ7 checks
  `errors.As` and the exact `Column`, so deleting one comparison cannot hide behind another one.
- **D5 — Sentinels the design does not name** but needs: request `ErrSchema`, `ErrString`,
  `ErrEpisode`, `ErrPath`, `ErrGitReuse`, `ErrNoGrants`, `ErrGrantEmpty`, `ErrNegative`,
  `ErrHalfProfile`, `ErrAgent`, `ErrMCPConfig`, `ErrDecode` and `ErrNegativeExpiry`; SSHSIG
  `ErrAuthorizedKey` (a malformed `.pub` line).
- **D6 — Planner-added tests**: **AP1** `TestApprovalImportsOnlyTheStandardLibrary` (constraint 3)
  and **SG7** `TestParseAuthorizedKeyRefusals`. Without SG7, `ParseAuthorizedKey`'s refusal
  branches would have no test that kills them.
- **D7 — The mpint refusals wrap `ErrBadSignature`** with the reasons `negative mpint` and
  `non-minimal mpint`. Their SG3 arms are derived so that **with the check deleted the signature
  would verify**: strip the real leading `0x00` from a high-bit r or s for the negative arm, and
  prepend `0x00` to s for the non-minimal arm. So deleting either check makes its arm accept, and
  the arm goes red. The negative arm needs a high-bit r or s in `p256.sig`, which happens about 75%
  of the time. SG1 asserts it, and if it is missing the executor re-runs `generate.sh` (M2a step 4).
- **D8 — `VerifySSHSIG` steps 3–5 use the expected `key`**, not `s.PublicKey`, and build the signed
  data with the **argument** `namespace`. Then deleting step 1 or step 2 changes the error, and the
  arm sees it.
- **D9 — SG1 also asserts the exact directory listing** (the eight names). That enforces
  constraint 1 mechanically.
- **D10 — XC1 hygiene**: no `t.Parallel()` (C8). Each `ssh-keygen` runs under
  `exec.CommandContext` with a 30 s timeout and `Env = append(os.Environ(), "SSH_AUTH_SOCK=")`, so
  no real agent is consulted (constraint 1). Arms run as loop iterations inside the two subtests,
  not as nested subtests, so the CI grep reads exactly three PASS lines.
- **D11 — `Validate` does not cross-check `Worktree == WorkspaceRoot/Episode`.** §2.1 lists no such
  branch. Sprint 2's re-derivation compares it.
- **D12 — `RequestLifetime = 10 * time.Minute` and `MessagePrefix = "world-approval/v1"`** are
  exported constants in `request.go`. The request-id grammar (`mr-` + 16 hex) is left to sprint 2,
  which first uses it.

### 1.4 Premises in the design found false

- **F1 — `message.bin` cannot be committed.** Design §2.8 and SG1 check in a 90-byte binary
  `message.bin`. By construction it has a NUL at offset 17, and the fixture challenge or digest can
  add more. The fleet-owned `scripts/tracked_binary_gate.sh`, run by `verify_go.sh`, which is CI's
  `go build + test gate` step, rejects any committed blob that contains a NUL (C7, reproduced). The
  only allowances are PNG and pinned fonts, and widening the gate is a reviewed fleet edit
  (D-WORLD-DRIVER-1). The design as written would have turned CI red **after** the merge-candidate
  commit. Fixed by D1. Round 1–3 quorum did not catch it. V18 ran the generator outside git.
- **F2 — `Verdict: Admit | …` next to `func Admit` does not compile** (`Admit redeclared in this
  block`, reproduced). Fixed by D2.
- **F3 (not a sprint-1 defect; logged for sprint 2).** §2.3 says `authority.MintForRequest` runs
  `MintMatches` on the built row **before** the write. But the r3 `created_at` arm bounds
  `created_at` by "the in-transaction consume reading", and per §2.4 that reading exists only
  inside `MintSessionForRequest` after `BeginTx`. As written, sprint 2 cannot do both. It must run
  the check inside the store transaction (for example a check callback), or apply the upper bound
  there. The pure function in this sprint takes `ConsumedAt` as an input and is unaffected.

---

## 2. Shared definitions (the executor copies these exactly)

### 2.1 `host/approval/request.go`

```go
// Package approval is the pure core of row 201's approvals by key
// (design_docs/planned/w-world-approvals-by-key-a0a1.md): the mint request,
// its canonical bytes and digest, the signed message, the admission law, the
// mint-row check and the SSHSIG verifier. Standard library only; no I/O, no
// clock, no randomness (coding standard S2): every effect is the caller's.
package approval

const RequestSchema = "world/mint-request/v1"

// MessagePrefix starts every signed approval message (design §2.2).
const MessagePrefix = "world-approval/v1"

// RequestLifetime is a request's fixed lifetime: expires_at = created_at + 600.
const RequestLifetime = 10 * time.Minute
```

The four types `MintRequest`, `Grant`, `ProfileRef`, `DaemonRef` are **exactly** the design's §2.1
block: same field order, same json tags, same comments. The `MintRequest` field order is `Schema,
Episode, WorldDB, WorkspaceRoot, Worktree, Reuse, Degraded, Repo, Git, Branch, Grants, TTLSeconds,
Profile, ProfilesDir, CredentialsDB, SessionFile, MCPConfig, Agent, Model, Task, Daemon`. RQ1's
golden pins it.

Sentinels (one `var ( … )` block, each `errors.New("approval: …")`):

```go
ErrSchema, ErrString, ErrEpisode, ErrPath, ErrGitReuse, ErrNoGrants, ErrGrantEmpty,
ErrNegative, ErrHalfProfile, ErrAgent, ErrMCPConfig, ErrDecode, ErrNotCanonical, ErrNegativeExpiry
```

Functions:

```go
func (r MintRequest) Validate() error
func Canonical(r MintRequest) ([]byte, error)
func Decode(b []byte) (MintRequest, error)
func Digest(r MintRequest) ([32]byte, error)                         // sha256.Sum256(Canonical(r))
func Diff(a, b MintRequest) string                                    // "" when equal
func SignedMessage(digest, challenge [32]byte, expiry int64) ([]byte, error)
func MessageHash(digest, challenge [32]byte, expiry int64) ([32]byte, error) // sha256.Sum256(SignedMessage(...))

// unexported helpers
var episodePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
func jsonName(f reflect.StructField) string // strings.Split(f.Tag.Get("json"), ",")[0]
func stringLeaves(v reflect.Value, path string, fn func(path, s string) error) error
func diffValue(a, b reflect.Value, path string) string
func absClean(p string) bool // filepath.IsAbs(p) && filepath.Clean(p) == p
```

**`Validate` checks these, in this order.** Each one is its own `if` branch and returns
`fmt.Errorf("%w: <json path> …", Err…)`, naming the JSON path:

1. `r.Schema != RequestSchema` → `ErrSchema`.
2. `stringLeaves(reflect.ValueOf(r), "", …)` over **every** string leaf (structs and slices are
   walked, and the path uses json names such as `grants[0].scope` or `daemon.addr`). A leaf that is
   `!utf8.ValidString(s)` or `strings.ContainsRune(s, 0)` → `ErrString`.
3. `!episodePattern.MatchString(r.Episode)` → `ErrEpisode`.
4. Required paths, in this order: `world_db`, `workspace_root`, `worktree`, `profiles_dir`,
   `credentials_db`, `session_file`. Any one that is not `absClean` → `ErrPath`. Then the optional
   paths, `repo`, `git`, `mcp_config`, `task`: any one that is **non-empty** and not `absClean` →
   `ErrPath`. Both are a `for _, p := range []struct{name, val string}{…}` loop.
5. `(r.Git != "") == r.Reuse` → `ErrGitReuse` (Git must be non-empty iff `!Reuse`).
6. `len(r.Grants) == 0` → `ErrNoGrants`.
7. For each grant `i`: `Effect == "" || Scope == ""` → `ErrGrantEmpty`, then `Budget < 0` →
   `ErrNegative` (`grants[i].budget`).
8. `r.TTLSeconds < 0` → `ErrNegative` (`ttl_s`).
9. `(r.Profile.Project == "") != (r.Profile.Digest == "")` → `ErrHalfProfile`.
10. `r.Agent != "claude" && r.Agent != "codex"` → `ErrAgent`.
11. `(r.Agent == "claude" && r.MCPConfig == "") || (r.Agent == "codex" && r.MCPConfig != "")` →
    `ErrMCPConfig`.

**`Canonical`**: `var buf bytes.Buffer; enc := json.NewEncoder(&buf); enc.SetEscapeHTML(false);
enc.Encode(r)`, then `bytes.TrimSuffix(buf.Bytes(), []byte("\n"))`. It does not call `Validate`.

**`Decode`**: `dec := json.NewDecoder(bytes.NewReader(b)); dec.DisallowUnknownFields()`. A decode
error becomes `fmt.Errorf("%w: %v", ErrDecode, err)`. Then `c, err := Canonical(r)`, and
`!bytes.Equal(c, b)` → `ErrNotCanonical` with a zero `MintRequest`. This one comparison is what
refuses trailing bytes, a second concatenated value, reformatting, reordering and case-folded keys.
`Decode` does not call `Validate`. Callers do.

**`Diff`**: `diffValue(reflect.ValueOf(a), reflect.ValueOf(b), "")`. A struct walks its fields in
order with path `parent.jsonName`. A slice of a different length returns its own path (for example
`grants`), and otherwise walks its elements with paths like `grants[1]`. A leaf compares with
`a.Interface() != b.Interface()` and returns the path. The first difference wins.

**`SignedMessage`**: `expiry < 0` → `nil, ErrNegativeExpiry`. Otherwise the message is
`MessagePrefix ‖ 0x00 ‖ digest[:] ‖ challenge[:] ‖ binary.BigEndian.AppendUint64(uint64(expiry))`,
always 90 bytes.

### 2.2 `host/approval/admit.go`

```go
type Verdict int

const (
	Admitted Verdict = iota + 1 // the zero Verdict is not an admission
	RefuseConsumed
	RefuseClockBehind
	RefuseExpired
	RefuseDigestMismatch
)

// String: "admit", "consumed", "clock-behind", "expired", "digest-mismatch"; otherwise fmt.Sprintf("verdict(%d)", int(v)).
func (v Verdict) String() string

type AdmitInput struct {
	Now, CreatedAt, ExpiresAt int64
	Consumed                  bool
	StoredDigest, BytesDigest [32]byte // the digest column; SHA-256 of the stored bytes
}

// Admit is the admission law of design §2.3, in its fixed order.
func Admit(in AdmitInput) Verdict {
	if in.Consumed { return RefuseConsumed }
	if in.Now < in.CreatedAt { return RefuseClockBehind }
	if in.Now >= in.ExpiresAt { return RefuseExpired }
	if in.StoredDigest != in.BytesDigest { return RefuseDigestMismatch }
	return Admitted
}

// FactGrant is one element of session_credentials.grants_json. Its field names
// and order match broker.Capability, which has no json tags (design V34).
type FactGrant struct {
	Effect    string
	Scope     string
	ExpiresAt int64
	Budget    int64
}

// MintFacts is the session_credentials row the mint built (its seven columns),
// plus the token's hash and the two times that bound the mint (design §2.3).
type MintFacts struct {
	CredentialID     string // credential_id
	EpisodeID        string // episode_id
	GrantsJSON       string // grants_json, raw; MintMatches decodes it
	ExpiresAt        int64  // expires_at
	CreatedAt        int64  // created_at
	Project          string // project
	ProfileDigest    string // profile_digest
	TokenHash        string // hashHexToken(the token handed to the writer)
	RequestCreatedAt int64  // the request row's created_at: lower bound for CreatedAt
	ConsumedAt       int64  // the in-transaction consume reading: upper bound for CreatedAt
}

type MintMismatchError struct{ Column, Want, Got string }

// Error: fmt.Sprintf("approval: the minted row's %s is %s, the request needs %s", e.Column, e.Got, e.Want)
func (e *MintMismatchError) Error() string

func MintMatches(r MintRequest, f MintFacts) error
```

**`MintMatches` checks these in this order.** Each check is one branch that returns
`&MintMismatchError{Column: …}`:

1. `f.CredentialID != f.TokenHash` → `credential_id`.
2. `f.EpisodeID != r.Episode` → `episode_id`.
3. Decode `f.GrantsJSON` with `json.NewDecoder` plus `DisallowUnknownFields` into `[]FactGrant`. On
   a decode error, or `len != len(r.Grants)`, or any element whose `Effect`, `Scope` or `Budget`
   differs from `r.Grants[i]` (compared in order) → `grants_json`.
4. Any element with `ExpiresAt != f.ExpiresAt` → `grants_json.ExpiresAt`.
5. `f.ExpiresAt != f.CreatedAt + r.TTLSeconds` → `expires_at`.
6. `f.CreatedAt < f.RequestCreatedAt` → `created_at`.
7. `f.CreatedAt > f.ConsumedAt` → `created_at`.
8. `f.Project != r.Profile.Project` → `project`.
9. `f.ProfileDigest != r.Profile.Digest` → `profile_digest`.
10. `return nil`.

### 2.3 `host/approval/sshsig.go`

```go
const Namespace = "world-approval@sunholo"

const (
	armorBegin = "-----BEGIN SSH SIGNATURE-----"
	armorEnd   = "-----END SSH SIGNATURE-----"
	keyEd25519 = "ssh-ed25519"
	keyP256    = "ecdsa-sha2-nistp256"
)

var (
	ErrArmor, ErrBase64, ErrMagic, ErrVersion, ErrTruncated, ErrTrailing, ErrReserved, ErrHashAlg,
	ErrKeyType, ErrNamespace, ErrKeyMismatch, ErrSigType, ErrBadSignature, ErrAuthorizedKey error // errors.New("approval: sshsig: …")
)

type SSHPublicKey struct{ Type string; Blob []byte }   // Blob = the SSH wire encoding
type SSHSig struct{ PublicKey SSHPublicKey; Namespace, HashAlg string; SigType string; Sig []byte }

func ParseAuthorizedKey(line []byte) (SSHPublicKey, error)
func (k SSHPublicKey) Fingerprint() string // "SHA256:" + base64.RawStdEncoding.EncodeToString(sha256(k.Blob))
func ParseSSHSIG(armored []byte) (*SSHSig, error)
func VerifySSHSIG(s *SSHSig, namespace string, message []byte, key SSHPublicKey) error
func VerifyApproval(armored []byte, digest, challenge [32]byte, expiry int64, key SSHPublicKey) error

// unexported
type sshReader struct{ b []byte }
func (r *sshReader) uint32() (uint32, error)     // len < 4 → ErrTruncated
func (r *sshReader) readString() ([]byte, error) // uint32 length; uint64(n) > uint64(len(r.b)) → ErrTruncated; returns r.b[:n:n]
func (r *sshReader) end() error                  // len(r.b) != 0 → ErrTrailing
func appendString(dst, s []byte) []byte          // uint32 BE length ‖ s
type keyMaterial struct{ typ string; ed ed25519.PublicKey; ec *ecdsa.PublicKey }
func parseKeyBlob(blob []byte) (keyMaterial, error)
func parseSSHSIGBlob(b []byte) (*SSHSig, error)
func parseMpint(b []byte) (*big.Int, error)
func signedData(namespace, hashAlg string, message []byte) ([]byte, error)
```

Allowed imports (constraint 3): `bytes`, `crypto/ecdsa`, `crypto/ed25519`, `crypto/elliptic`,
`crypto/sha256`, `crypto/sha512`, `encoding/base64`, `encoding/binary`, `errors`, `fmt`,
`math/big`, `strings`.

**`parseKeyBlob`.** Read the type with `readString`, then switch on it:
- `keyEd25519`: read `pk`, then `end()`. `len(pk) != ed25519.PublicKeySize` →
  `fmt.Errorf("%w: ssh-ed25519 key is %d bytes, want 32", ErrKeyType, len(pk))`.
- `keyP256`: read `curve` and `q`, then `end()`. `string(curve) != "nistp256"` → `ErrKeyType`.
  `ecdsa.ParseUncompressedPublicKey(elliptic.P256(), q)` erroring →
  `fmt.Errorf("%w: ecdsa-sha2-nistp256 point: %v", ErrKeyType, err)`.
- `"sk-ssh-ed25519@openssh.com"` or `"sk-ecdsa-sha2-nistp256@openssh.com"`:
  `fmt.Errorf("%w: %s security-key signatures are refused until an attended hardware fixture exists (row 201 Q1)", ErrKeyType, t)`.
- `default`: `fmt.Errorf("%w: %q is not ssh-ed25519 or ecdsa-sha2-nistp256", ErrKeyType, t)`.

Every error path returns a zero `keyMaterial`.

**`ParseSSHSIG`.** It **never returns a non-nil `*SSHSig` together with an error**, and never a nil
error unless it parsed everything.
1. `len(armored) == 0` → `nil, ErrArmor`.
2. `s := strings.TrimSuffix(string(armored), "\n")` (exactly one optional trailing newline), then
   `lines := strings.Split(s, "\n")`. `len(lines) < 3`, `lines[0] != armorBegin`,
   `lines[len-1] != armorEnd` or any empty body line → `ErrArmor`. A CRLF file fails here, because
   `lines[0]` ends in `\r`.
3. `blob, err := base64.StdEncoding.Strict().DecodeString(strings.Join(body, ""))` → on error
   `fmt.Errorf("%w: %v", ErrBase64, err)`.
4. `return parseSSHSIGBlob(blob)`.

**`parseSSHSIGBlob`.**
1. `len(b) < 6` → `ErrTruncated`. `!bytes.Equal(b[:6], []byte("SSHSIG"))` → `ErrMagic`.
2. `r := sshReader{b[6:]}`. Read the version: `!= 1` → `fmt.Errorf("%w: %d", ErrVersion, v)`.
3. Read `pk`, `ns`, `reserved`, `hashAlg`, `sig` (5 × `readString`), then `r.end()`.
4. `len(reserved) != 0` → `ErrReserved`. `hashAlg` not `sha256` or `sha512` → `ErrHashAlg`.
5. `km, err := parseKeyBlob(pk)`. Return the error if there is one.
6. `sr := sshReader{sig}`: read `sigType` and `sigBlob`, then `sr.end()`.
7. Return `&SSHSig{PublicKey: SSHPublicKey{Type: km.typ, Blob: bytes.Clone(pk)}, Namespace:
   string(ns), HashAlg: string(hashAlg), SigType: string(sigType), Sig: bytes.Clone(sigBlob)}`.

**`VerifySSHSIG`** runs a fixed sequence (design §2.7, D8):
0. `s == nil` → `ErrBadSignature`.
1. `s.Namespace != namespace` → `ErrNamespace`.
2. `!bytes.Equal(s.PublicKey.Blob, key.Blob)` → `ErrKeyMismatch`.
3. `km, err := parseKeyBlob(key.Blob)` → return the error (`ErrKeyType` family).
4. `s.SigType != km.typ` → `ErrSigType`.
5. `sd, err := signedData(namespace, s.HashAlg, message)`. Then:
   - `keyEd25519`: `len(s.Sig) != ed25519.SignatureSize || !ed25519.Verify(km.ed, sd, s.Sig)` →
     `ErrBadSignature`.
   - `keyP256`: `sr := sshReader{s.Sig}`; read `rb`, `sb`, then `end()`. Any read error →
     `fmt.Errorf("%w: %v", ErrBadSignature, err)`. Then `parseMpint(rb)` and `parseMpint(sb)`, then
     `h := sha256.Sum256(sd)`; `!ecdsa.Verify(km.ec, h[:], r, s)` → `ErrBadSignature`.
   - `default` → `ErrKeyType`.
6. `return nil`.

**`signedData`** is `"SSHSIG" ‖ appendString(namespace) ‖ appendString(empty reserved) ‖
appendString(hashAlg) ‖ appendString(H(message))`, where `H` is `sha512.Sum512` for `sha512` and
`sha256.Sum256` for `sha256`. Any other `hashAlg` → `ErrHashAlg`.

**`parseMpint`**:
```go
if len(b) > 0 && b[0]&0x80 != 0 {
	return nil, fmt.Errorf("%w: negative mpint", ErrBadSignature)
}
if len(b) > 0 && b[0] == 0 && (len(b) == 1 || b[1]&0x80 == 0) {
	return nil, fmt.Errorf("%w: non-minimal mpint", ErrBadSignature)
}
return new(big.Int).SetBytes(b), nil
```

**`VerifyApproval`**: `ParseSSHSIG(armored)`, then `SignedMessage(digest, challenge, expiry)`, then
`VerifySSHSIG(s, Namespace, msg, key)`. Return the first error.

**`ParseAuthorizedKey`**:
1. `s := strings.TrimSuffix(string(line), "\n")`. `s == ""` or `strings.ContainsAny(s, "\r\n")` →
   `ErrAuthorizedKey`.
2. `f := strings.Fields(s)`. `len(f) < 2` → `ErrAuthorizedKey`.
3. `blob, err := base64.StdEncoding.Strict().DecodeString(f[1])`. On error →
   `fmt.Errorf("%w: %v", ErrBase64, err)`.
4. `km, err := parseKeyBlob(blob)`. Return the error if there is one.
5. `f[0] != km.typ` → `fmt.Errorf("%w: the line says %q, the key is %q", ErrKeyType, f[0], km.typ)`.
6. `return SSHPublicKey{Type: km.typ, Blob: blob}, nil`.

**The phrase `PRIVATE KEY` must not appear anywhere under `host/approval/`.** The acceptance grep in
M2a relies on that.

### 2.4 Test fixtures and constants (exact)

In `host/approval/request_test.go`:

```go
func fixtureRequest() MintRequest {
	return MintRequest{
		Schema: RequestSchema, Episode: "fixture-ep",
		WorldDB: "/ws/state/world.db", WorkspaceRoot: "/ws/workspace", Worktree: "/ws/workspace/fixture-ep",
		Repo: "/src/ailang-world", Git: "/usr/bin/git", Branch: "world/fixture-ep",
		Grants:     []Grant{{Effect: "Workspace.Read", Scope: "fixture-ep", Budget: 0}, {Effect: "Workspace.Write", Scope: "fixture-ep", Budget: 25}},
		TTLSeconds: 3600, Profile: ProfileRef{Project: "ailang-world", Digest: strings.Repeat("ab", 32)},
		ProfilesDir: "/ws/profiles", CredentialsDB: "/ws/state/credentials.db",
		SessionFile: "/ws/sessions/fixture-ep/session.token", MCPConfig: "/ws/sessions/fixture-ep/mcp.json",
		Agent: "claude", Model: "claude-opus<5>&x", Task: "/ws/tasks/fixture-ep.md",
		Daemon: DaemonRef{Addr: "127.0.0.1:7377", VCSRevision: "0123456789abcdef0123456789abcdef01234567",
			VCSModified: false, ExecutableSHA256: strings.Repeat("cd", 32), GoVersion: "go1.26.6"},
	}
}

const fixtureCanonicalJSON = `{"schema":"world/mint-request/v1","episode":"fixture-ep","world_db":"/ws/state/world.db","workspace_root":"/ws/workspace","worktree":"/ws/workspace/fixture-ep","reuse":false,"degraded":false,"repo":"/src/ailang-world","git":"/usr/bin/git","branch":"world/fixture-ep","grants":[{"effect":"Workspace.Read","scope":"fixture-ep","budget":0},{"effect":"Workspace.Write","scope":"fixture-ep","budget":25}],"ttl_s":3600,"profile":{"project":"ailang-world","digest":"abababababababababababababababababababababababababababababababab"},"profiles_dir":"/ws/profiles","credentials_db":"/ws/state/credentials.db","session_file":"/ws/sessions/fixture-ep/session.token","mcp_config":"/ws/sessions/fixture-ep/mcp.json","agent":"claude","model":"claude-opus<5>&x","task":"/ws/tasks/fixture-ep.md","daemon":{"addr":"127.0.0.1:7377","vcs_revision":"0123456789abcdef0123456789abcdef01234567","vcs_modified":false,"executable_sha256":"cdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcdcd","go_version":"go1.26.6"}}`

const fixtureDigestHex = "0ace927ef952f541b1e454b066f6d6b1941fb5f3d710c784c22dd80a53d771eb"

const fixtureMessageHashHex = "29181d224d982f659d1bd4703238d3f8489e26968ffe2a22bd786fa3bf15fef5"

func fixtureDigest(t *testing.T) [32]byte // hex-decodes fixtureDigestHex into [32]byte; t.Fatal on error
```

These goldens were computed by an independent program (§1.2 C6). **If the code disagrees, fix the
code. Never re-paste a golden from the code's own output.**

In `host/approval/sshsig_test.go`, created in **M1a** with only this block and grown in M2. The
`fixtureMessageHex` line **must stay a standalone top-level `const` on one line at column 0**
(not inside a `const ( … )` block), because `generate.sh` extracts it with
`sed -n 's/^const fixtureMessageHex = "\([0-9a-f]*\)"$/\1/p'`:

```go
// fixtureMessageHex is SignedMessage(fixtureDigest, fixtureChallenge(), fixtureExpiry):
// the 90 bytes every SSHSIG fixture signs. testdata/sshsig/generate.sh reads this line.
const fixtureMessageHex = "776f726c642d617070726f76616c2f7631000ace927ef952f541b1e454b066f6d6b1941fb5f3d710c784c22dd80a53d771eba0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf000000006955b900"

// fixtureExpiry is 2026-01-01T00:00:00Z.
const fixtureExpiry int64 = 1767225600

func fixtureChallenge() [32]byte {
	var c [32]byte
	for i := range c {
		c[i] = byte(0xa0 + i)
	}
	return c
}
```

**Test-side SSHSIG helpers** (`sshsig_test.go`, M2a). They are independent of the production
parser:

```go
func fixtureBytes(t *testing.T, name string) []byte // os.ReadFile(filepath.Join("testdata", "sshsig", name)); t.Fatal on error
func fixtureMessage(t *testing.T) []byte             // hex.DecodeString(fixtureMessageHex)
func fixtureKey(t *testing.T, name string) SSHPublicKey // ParseAuthorizedKey(fixtureBytes(t, name)); t.Fatal on error
func unarmorTest(t *testing.T, armored []byte) []byte  // own split on "\n", drop BEGIN/END, base64.StdEncoding decode
func armorTest(blob []byte) []byte                      // BEGIN line, base64 wrapped at 70 columns, END line, one trailing "\n"
func sshStringTest(b []byte) []byte                     // uint32 BE length ‖ b
func readStringsTest(t *testing.T, b []byte, n int) [][]byte // exactly n length-prefixed strings, t.Fatal on short or trailing

type sigParts struct {
	magic                                             []byte
	version                                           uint32
	publicKey, namespace, reserved, hashAlg, signature []byte
}
func splitSig(t *testing.T, armored []byte) sigParts // unarmorTest, then magic (6), version (4), readStringsTest(…, 5)
func (p sigParts) armored() []byte                    // re-encode in order, then armorTest
```

**Named transforms.** Each one takes a real fixture's armored bytes and returns new armored bytes
built from `splitSig`, one field changed, then `.armored()`. Each name says what it changes:

```go
func withVersion(t *testing.T, a []byte, v uint32) []byte
func withMagic(t *testing.T, a []byte, m string) []byte          // m must be 6 bytes
func withReserved(t *testing.T, a []byte, r string) []byte
func withHashAlg(t *testing.T, a []byte, alg string) []byte
func withTrailingByte(t *testing.T, a []byte) []byte              // one 0x00 after the signature string, inside the blob
func withKeyType(t *testing.T, a []byte, typ string) []byte       // rewrites the FIRST string of the publickey blob, keeps the rest
func withCurveName(t *testing.T, a []byte, curve string) []byte   // p256 only: rewrites the second string of the publickey blob
func withPointByteFlipped(t *testing.T, a []byte) []byte          // p256 only: Q[33] ^= 0x01 (point leaves the curve)
func withShortEd25519Key(t *testing.T, a []byte) []byte           // ed25519 only: pk[:31]
func withSigType(t *testing.T, a []byte, typ string) []byte       // rewrites sig-format-id inside the signature string
func withLastSigByteFlipped(t *testing.T, a []byte) []byte        // last byte of the inner sig blob ^= 0x01
func withNegativeMpint(t *testing.T, a []byte) []byte             // p256: strip the leading 0x00 of the first of r, s that has one (a high-bit value); t.Fatal if neither does
func withNonMinimalMpint(t *testing.T, a []byte) []byte           // p256: prepend 0x00 to s
func flipByte(b []byte, i int) []byte                             // copy with b[i] ^= 0x01
```

### 2.5 `host/approval/testdata/sshsig/generate.sh` (exact; mode 0755; run by hand, never in CI)

This is the design's §2.8 generator. Only the three message lines differ, because of D1 (the binary
message lives only in `$tmp`), and the closing reminder is new. It was dry-run by the planner
(§1.2 C5).

```bash
#!/usr/bin/env bash
# Regenerates the SSHSIG fixtures with REAL ssh-keygen (row 201 A1a). Test keys only:
# they are made in a temp dir and deleted on exit; only public keys and signatures are kept.
# The 90-byte message is kept as message.hex (text), never as a binary file: the tracked-binary
# gate (scripts/tracked_binary_gate.sh) rejects any committed blob that holds a NUL byte.
# After a run, update fixtureSHA256 in ../../sshsig_test.go IN THE SAME COMMIT (test SG1).
set -euo pipefail
here=$(cd "$(dirname "$0")" && pwd)
hex=$(sed -n 's/^const fixtureMessageHex = "\([0-9a-f]*\)"$/\1/p' "$here/../../sshsig_test.go")
[ ${#hex} -eq 180 ] || { echo "fixtureMessageHex not found" >&2; exit 1; }
tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
printf '%s\n' "$hex" > "$here/message.hex"
printf '%s' "$hex" | xxd -r -p > "$tmp/message.bin"
ssh-keygen -q -t ed25519 -N '' -C world-approval-fixture-ed25519 -f "$tmp/ed25519"
ssh-keygen -q -t ecdsa -b 256 -N '' -C world-approval-fixture-p256 -f "$tmp/p256"
ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/ed25519" < "$tmp/message.bin" > "$here/ed25519.sig"
ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/p256" < "$tmp/message.bin" > "$here/p256.sig"
ssh-keygen -Y sign -n git -f "$tmp/ed25519" < "$tmp/message.bin" > "$here/git-ns-ed25519.sig"
cp "$tmp/ed25519.pub" "$tmp/p256.pub" "$here/"
{ ssh -V 2>&1; date -u +%FT%TZ; ssh-keygen -l -f "$here/ed25519.pub"; ssh-keygen -l -f "$here/p256.pub"; } > "$here/SOURCE.txt"
echo "regenerated; now update fixtureSHA256 in sshsig_test.go: shasum -a 256 $here/*" >&2
```

### 2.6 `.gitattributes` (new, repo root, exactly one line)

```
host/approval/testdata/sshsig/* -text
```

---

## 3. Milestones

Every milestone ends with `go build ./... && go vet ./host/approval/` (`go build` does not compile
test files) plus its own tests. **No `t.Parallel()` anywhere in `host/approval`** (C8). Run every
command from the worktree root. Wrap loops in `bash -c` (the controller's shell is zsh).

### M1a — request core: types, Validate, Canonical, Decode, Digest, Diff, SignedMessage (design M1, part 1) · ~190 prod + ~330 test LOC

**Files**
- New `host/approval/request.go` (§2.1).
- New `host/approval/request_test.go` (fixtures from §2.4; RQ1–RQ5, RQ8).
- New `host/approval/sshsig_test.go` holding **only** the §2.4 `fixtureMessageHex`,
  `fixtureExpiry` and `fixtureChallenge` block.
- New `host/approval/imports_test.go` (AP1).

**Tests** (each kills the mutant named; mutant ids refer to §4)
- **RQ1 `TestCanonicalRequestGolden`**: `Canonical(fixtureRequest())` equals
  `[]byte(fixtureCanonicalJSON)` byte for byte. It has no trailing `\n` and contains the literal
  `claude-opus<5>&x`. `hex(Digest(fixtureRequest()))` equals `fixtureDigestHex`. *Kills*
  MUT-RQ-ORDER, MUT-RQ-HTML, MUT-RQ-NEWLINE.
- **RQ2 `TestDecodeIsStrict`**: `Decode([]byte(fixtureCanonicalJSON))` `reflect.DeepEqual`s
  `fixtureRequest()`, and re-encoding it gives the same bytes. Each arm gives a zero result and an
  error that `errors.Is(ErrNotCanonical) || errors.Is(ErrDecode)`. The arms are:
  (a) `{"extra":1,` spliced in after the leading `{`;
  (b) the canonical bytes + `"\n"`;
  (c) the `json.Indent(…, "", "  ")` output;
  (d) the canonical bytes with `"reuse":false,"degraded":false` replaced by
  `"degraded":false,"reuse":false` (assert the replacement changed the bytes);
  (e) the canonical bytes twice;
  (f) the canonical bytes with `"schema":` replaced by `"Schema":`;
  (g) empty input.
  Arms (b), (c), (d), (e) and (f) must be exactly `ErrNotCanonical`. *Kills* MUT-RQ-NOREENCODE.
- **RQ3 `TestDigestCoversEveryField`**: a test-side reflection walk over every leaf of
  `fixtureRequest()`. Structs are recursed, and each slice element is recursed with the path
  `grants[i]`. For each leaf, deep-copy the request (copy the `Grants` slice), change only that
  leaf (string `+= "x"`, bool `= !bool`, int64 `+= 1`), and require that `Digest` differs from the
  base digest, naming the path on failure. Assert the leaf count is `>= 20` (it is 31) and log it.
  *Kills* MUT-RQ-HIDDEN.
- **RQ4 `TestValidateRefusals`**: these must give `nil`: `fixtureRequest()`; a codex variant
  (`Agent "codex"`, `MCPConfig ""`); a reuse variant (`Reuse true`, `Git ""`). Then a table of
  `{name, mutate func(*MintRequest), want error, pathSubstring string}`. Each arm requires
  `errors.Is(err, want)` and `strings.Contains(err.Error(), pathSubstring)`:

  | Arm | Mutation | want | path | Kills |
  |---|---|---|---|---|
  | schema | `Schema = "world/mint-request/v2"` | ErrSchema | `schema` | MUT-V-SCHEMA |
  | utf8 | `Model = "\xff"` | ErrString | `model` | MUT-V-STRING |
  | nul-branch | `Branch = "a\x00b"` | ErrString | `branch` | MUT-V-STRING |
  | nul-grant | `Grants[0].Scope = "a\x00"` | ErrString | `grants[0].scope` | MUT-V-STRING |
  | nul-daemon | `Daemon.Addr = "1\x00"` | ErrString | `daemon.addr` | MUT-V-STRING |
  | nul-path | `WorldDB = "/ws/\x00"` | ErrString | `world_db` | MUT-V-STRING (fires before ErrPath) |
  | episode-case | `Episode = "Bad_Ep"` | ErrEpisode | `episode` | MUT-V-EPISODE |
  | episode-empty | `Episode = ""` | ErrEpisode | `episode` | MUT-V-EPISODE |
  | episode-long | `Episode = strings.Repeat("a", 65)` | ErrEpisode | `episode` | MUT-V-EPISODE |
  | world-db-rel | `WorldDB = "ws/state/world.db"` | ErrPath | `world_db` | MUT-V-PATH |
  | world-db-empty | `WorldDB = ""` | ErrPath | `world_db` | MUT-V-PATH |
  | ws-root-unclean | `WorkspaceRoot = "/ws/workspace/"` | ErrPath | `workspace_root` | MUT-V-PATH |
  | worktree-dotdot | `Worktree = "/ws/workspace/../x"` | ErrPath | `worktree` | MUT-V-PATH |
  | profiles-dir-rel | `ProfilesDir = "profiles"` | ErrPath | `profiles_dir` | MUT-V-PATH, MUT-V-PATHFIELD |
  | creds-rel | `CredentialsDB = "credentials.db"` | ErrPath | `credentials_db` | MUT-V-PATH |
  | session-rel | `SessionFile = "s.token"` | ErrPath | `session_file` | MUT-V-PATH |
  | repo-rel | `Repo = "src/ailang-world"` | ErrPath | `repo` | MUT-V-PATH |
  | git-rel | `Git = "git"` | ErrPath | `git` | MUT-V-PATH |
  | mcp-rel | `MCPConfig = "mcp.json"` | ErrPath | `mcp_config` | MUT-V-PATH |
  | task-rel | `Task = "task.md"` | ErrPath | `task` | MUT-V-PATH |
  | git-missing | `Git = ""` (Reuse false) | ErrGitReuse | `git` | MUT-V-GITREUSE |
  | git-on-reuse | `Reuse = true` (Git kept) | ErrGitReuse | `git` | MUT-V-GITREUSE |
  | no-grants-nil | `Grants = nil` | ErrNoGrants | `grants` | MUT-V-NOGRANTS |
  | no-grants-empty | `Grants = []Grant{}` | ErrNoGrants | `grants` | MUT-V-NOGRANTS |
  | effect-empty | `Grants[1].Effect = ""` | ErrGrantEmpty | `grants[1]` | MUT-V-GRANTEMPTY |
  | scope-empty | `Grants[0].Scope = ""` | ErrGrantEmpty | `grants[0]` | MUT-V-GRANTEMPTY |
  | budget-neg | `Grants[1].Budget = -1` | ErrNegative | `grants[1].budget` | MUT-V-BUDGET |
  | ttl-neg | `TTLSeconds = -1` | ErrNegative | `ttl_s` | MUT-V-TTL |
  | half-profile-a | `Profile.Project = ""` | ErrHalfProfile | `profile` | MUT-V-HALFPROFILE |
  | half-profile-b | `Profile.Digest = ""` | ErrHalfProfile | `profile` | MUT-V-HALFPROFILE |
  | agent | `Agent = "gemini"` | ErrAgent | `agent` | MUT-V-AGENT |
  | mcp-for-codex | `Agent = "codex"` (MCPConfig kept) | ErrMCPConfig | `mcp_config` | MUT-V-MCP |
  | mcp-missing | `MCPConfig = ""` (claude) | ErrMCPConfig | `mcp_config` | MUT-V-MCP |
- **RQ5 `TestSignedMessageGolden`**: `m, err := SignedMessage(fixtureDigest(t), fixtureChallenge(), fixtureExpiry)`.
  Then: `hex(m) == fixtureMessageHex`; `len(m) == 90`; `string(m[:17]) == MessagePrefix`;
  `m[17] == 0`; `m[18:50]` is the digest; `m[50:82]` is the challenge;
  `binary.BigEndian.Uint64(m[82:90]) == uint64(fixtureExpiry)`. `hex(MessageHash(…))` equals
  `fixtureMessageHashHex`, which equals `sha256(m)`. `SignedMessage(…, -1)` gives `nil` and
  `errors.Is(err, ErrNegativeExpiry)`. Expiry `0` is accepted. *Kills* MUT-SM-NUL, MUT-SM-LE,
  MUT-SM-NEGEXP.
- **RQ8 `TestDiffNamesTheField`**: `Diff(f, f) == ""`. Then: `Profile.Digest` changed gives
  `profile.digest`; `Grants[1].Budget++` gives `grants[1].budget`; `Daemon.ExecutableSHA256`
  changed gives `daemon.executable_sha256`; a third grant appended gives `grants`; `Episode` and
  `Model` both changed give `episode` (the first in order). *Kills* MUT-DIFF-EMPTY.
- **AP1 `TestApprovalImportsOnlyTheStandardLibrary`** (`imports_test.go`): `go/parser.ParseDir` is
  deprecated in go1.26.6 (`go doc go/parser.ParseDir`), so read the directory with `os.ReadDir(".")` and parse each non-`_test.go` `.go` file
  with `parser.ParseFile(fset, name, nil, parser.ImportsOnly)`. Every import path (unquoted) must
  have no `.` in its first `/`-segment, and must not be in the forbidden set `{os, os/exec,
  os/signal, net, net/http, syscall, io/ioutil, crypto/rand, math/rand, math/rand/v2, unsafe}`.
  Non-vacuity: at least 1 production file and at least 5 imports seen (after M2 it is ≥ 3 files;
  the assertion stays `>= 1` so M1a can pass). Positive control: the same check run on the
  in-memory source `package x; import _ "github.com/sunholo-data/ailang-world/host/store"`
  (`parser.ParseFile(fset, "x.go", src, parser.ImportsOnly)`) must report exactly one violation.
  *Kills* MUT-AP-IMPORT.

**Acceptance**
```
test -z "$(gofmt -l host/approval)"
go build ./... && go vet ./host/approval/
go test ./host/approval/ -run '^(TestCanonicalRequestGolden|TestDecodeIsStrict|TestDigestCoversEveryField|TestValidateRefusals|TestSignedMessageGolden|TestDiffNamesTheField|TestApprovalImportsOnlyTheStandardLibrary)$' -count=1 -v
go test -run '^$' ./host/approval/
```

### M1b — admission law and mint-row check (design M1, part 2) · ~110 prod + ~190 test LOC

**Files**: new `host/approval/admit.go` (§2.2), new `host/approval/admit_test.go`.

**Tests**
- **RQ6 `TestAdmitLaw`**: base `AdmitInput{Now: 1000, CreatedAt: 900, ExpiresAt: 1500}` with equal
  digests (`[32]byte{1}`). Table:

  | Arm | Change from base | Want |
  |---|---|---|
  | clean | — | Admitted |
  | at-expiry | `Now = 1500` | RefuseExpired |
  | expiry-minus-1 | `Now = 1499` | Admitted |
  | at-created | `Now = 900` | Admitted |
  | clock-behind | `Now = 899` | RefuseClockBehind |
  | consumed | `Consumed = true` | RefuseConsumed |
  | consumed+expired | `Consumed = true, Now = 2000` | RefuseConsumed |
  | consumed+behind | `Consumed = true, Now = 1` | RefuseConsumed |
  | digest | `BytesDigest = [32]byte{2}` | RefuseDigestMismatch |
  | digest+expired | `BytesDigest = [32]byte{2}, Now = 1500` | RefuseExpired |
  | digest+behind | `BytesDigest = [32]byte{2}, Now = 899` | RefuseClockBehind |

  It also asserts the five `String()` names (§2.2), `Verdict(0).String() == "verdict(0)"`, and
  `Admitted != Verdict(0)`. *Kills* MUT-AD-CONSUMED, MUT-AD-CLOCK, MUT-AD-EXPIRED, MUT-AD-DIGEST,
  MUT-AD-GE, MUT-AD-ORDER, MUT-AD-ZERO.
- **RQ7 `TestMintMatchesEveryColumn`**: `r := fixtureRequest()`. Make the token hash test-side as
  `th := hex(sha256([]byte(strings.Repeat("aa", 32))))`. The test-side
  `type capJSON struct{ Effect, Scope string; ExpiresAt, Budget int64 }` (this field order) is
  marshalled into `GrantsJSON`. Use a helper `grantsJSON(t, exp int64, gs ...Grant) string`.
  Clean facts:
  `{CredentialID: th, EpisodeID: "fixture-ep", GrantsJSON: grantsJSON(t, 4600, r.Grants...),
  ExpiresAt: 4600, CreatedAt: 1000, Project: "ailang-world", ProfileDigest: strings.Repeat("ab", 32),
  TokenHash: th, RequestCreatedAt: 900, ConsumedAt: 1010}`.
  The clean facts give `nil`, and so do the boundary arms `CreatedAt = 900` (with `ExpiresAt` and
  the grants' `ExpiresAt` set to 4500) and `CreatedAt = 1010` (4610). Each refusal arm requires
  `errors.As(err, &mm)` and an exact `mm.Column`:

  | Arm | Change | Column | Kills |
  |---|---|---|---|
  | cred-id | `CredentialID = hex(sha256("other"))` | `credential_id` | MUT-MM-CREDID |
  | episode | `EpisodeID = "other-ep"` | `episode_id` | MUT-MM-EPISODE |
  | grants-undecodable | `GrantsJSON = "not json"` | `grants_json` | (decode arm) |
  | grants-unknown-key | `GrantsJSON = `[{"Effect":"Workspace.Read","Scope":"fixture-ep","ExpiresAt":4600,"Budget":0,"Extra":1},…second…]`` | `grants_json` | (strict decode) |
  | grants-extra | three grants (the third a copy of the second) | `grants_json` | (length check; MUT-MM-GRANTSCOUNT survives this arm by design, the next three kill it) |
  | grants-reordered | `grantsJSON(t, 4600, r.Grants[1], r.Grants[0])` | `grants_json` | MUT-MM-GRANTSCOUNT |
  | grants-budget | second grant `Budget: 26` | `grants_json` | MUT-MM-GRANTSCOUNT |
  | grants-scope | first grant `Scope: "other"` | `grants_json` | MUT-MM-GRANTSCOUNT |
  | grant-exp-zero | grants marshalled with `ExpiresAt 0` | `grants_json.ExpiresAt` | MUT-MM-GRANTEXP |
  | grant-exp-plus1 | grants marshalled with `ExpiresAt 4601` | `grants_json.ExpiresAt` | MUT-MM-GRANTEXP |
  | expires-off-by-one | `ExpiresAt = 4601` and grants with 4601 | `expires_at` | MUT-MM-EXPIRES |
  | created-before | `CreatedAt = 899`, `ExpiresAt = 4499`, grants with 4499 | `created_at` | MUT-MM-CREATEDLOW |
  | created-after | `CreatedAt = 1011`, `ExpiresAt = 4611`, grants with 4611 | `created_at` | MUT-MM-CREATEDHIGH |
  | other-digest | `ProfileDigest = strings.Repeat("ef", 32)` | `profile_digest` | MUT-MM-DIGEST |
  | bound-req-unbound-row | `Project = "", ProfileDigest = ""` | `project` | MUT-MM-PROJECT |
  | unbound-req-bound-row | `r.Profile = ProfileRef{}` (facts kept) | `project` | MUT-MM-PROJECT |

**Acceptance**
```
test -z "$(gofmt -l host/approval)"
go build ./... && go vet ./host/approval/
go test ./host/approval/ -run '^(TestAdmitLaw|TestMintMatchesEveryColumn)$' -count=1 -v
# design §8's M1 gate, verbatim:
go test ./host/approval/ -run '^(TestCanonical|TestDecode|TestDigest|TestValidate|TestSignedMessage|TestAdmit|TestMintMatches|TestDiff)' -count=1
```

### M2a — fixtures, generator, `.gitattributes`, the parser half of `sshsig.go` (design M2, part 1) · ~170 prod + ~22 script/attr + ~230 test LOC

**Order matters. Do the steps in this order.**
1. Write `host/approval/testdata/sshsig/generate.sh` **exactly** as §2.5, then
   `chmod 755 host/approval/testdata/sshsig/generate.sh`.
2. Write `.gitattributes` exactly as §2.6.
3. Write the parser half of `host/approval/sshsig.go` (§2.3): the constants, errors, types,
   `sshReader`, `appendString`, `parseKeyBlob`, `ParseSSHSIG`, `parseSSHSIGBlob`,
   `ParseAuthorizedKey`, `Fingerprint`. Leave `VerifySSHSIG`, `parseMpint`, `signedData` and
   `VerifyApproval` for M2b.
4. Run `bash host/approval/testdata/sshsig/generate.sh` once. It uses real `ssh-keygen` and test
   keys only (constraint 1). Then check that `p256.sig` has a high-bit r or s. In a scratch `_test`
   run or in SG1 itself (below), `splitSig` then read r and s: **one of them must start with `0x00`
   followed by a byte `>= 0x80`**. If neither does, re-run `generate.sh` (D7). Never edit a
   fixture by hand.
5. Run `shasum -a 256 host/approval/testdata/sshsig/*` and paste the eight hashes into
   `fixtureSHA256` (below). **From now on, any re-run of `generate.sh` must update `fixtureSHA256`
   and `SOURCE.txt` in the same commit.**

**Comment and table in `sshsig_test.go`** (directly above the table, constraint 2):

```go
// The SSHSIG fixtures under testdata/sshsig/ are produced ONLY by
// testdata/sshsig/generate.sh, with real ssh-keygen, signing from stdin with
// throwaway keys made in a mktemp -d directory that it deletes:
//
//	ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/ed25519" < "$tmp/message.bin" > "$here/ed25519.sig"
//	ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/p256" < "$tmp/message.bin" > "$here/p256.sig"
//	ssh-keygen -Y sign -n git -f "$tmp/ed25519" < "$tmp/message.bin" > "$here/git-ns-ed25519.sig"
//
// SOURCE.txt records the tool version. Whoever runs generate.sh updates this
// table and SOURCE.txt in the same commit (design §2.8, test SG1).
var fixtureSHA256 = map[string]string{
	"SOURCE.txt":          "<64 hex>",
	"ed25519.pub":         "<64 hex>",
	"ed25519.sig":         "<64 hex>",
	"generate.sh":         "<64 hex>",
	"git-ns-ed25519.sig":  "<64 hex>",
	"message.hex":         "<64 hex>",
	"p256.pub":            "<64 hex>",
	"p256.sig":            "<64 hex>",
}
```

**Tests** (in `sshsig_test.go`, plus the §2.4 helpers and transforms)
- **SG1 `TestSSHSIGFixturesAreTheGeneratedBytes`**:
  (a) The sorted `os.ReadDir("testdata/sshsig")` names equal the sorted keys of `fixtureSHA256`.
  There are exactly 8, and no other file is allowed (D9, constraint 1).
  (b) The SHA-256 of each file equals its pinned value, and a failure names the file.
  (c) `message.hex` is exactly `fixtureMessageHex + "\n"`, and its hex-decoding equals
  `SignedMessage(fixtureDigest(t), fixtureChallenge(), fixtureExpiry)`.
  (d) Each `.sig` is non-empty and starts with `"-----BEGIN SSH SIGNATURE-----\n"`.
  (e) No fixture byte is `0x00` or `\r`.
  (f) The D7 precondition: in `splitSig(p256.sig)`, r or s has a leading `0x00` before a byte
  `>= 0x80`.
  *Kills* MUT-FX-EDIT, MUT-FX-CRLF, MUT-FX-EXTRA, and a message-encoding drift (MUT-SM-LE also
  turns this red).
- **SG4 `TestParseSSHSIGFailsLoudly`**: for every fixture `.sig` (`ed25519.sig`, `p256.sig`,
  `git-ns-ed25519.sig`) and every prefix length `0 <= n < len(blob)` of its decoded blob,
  `ParseSSHSIG(armorTest(blob[:n]))` returns `nil` and a non-nil error, and does not panic (wrap
  each call in a helper that `recover`s and fails). Then non-empty garbage gives `nil` and an
  error:
  - 64 seeded random byte slices from `rand.New(rand.NewPCG(201, 373))` (`math/rand/v2`), lengths
    1–512, both raw and wrapped with `armorTest`;
  - `"-----BEGIN PGP SIGNATURE-----\n\niQEz\n-----END PGP SIGNATURE-----\n"`;
  - `{"sig":"x"}`;
  - `armorTest([]byte("hello"))`;
  - armor fault arms (assert `errors.Is(err, ErrArmor)`): the fixture with every `\n` turned into
    `\r\n`; the fixture plus an extra `"\n"` (two trailing newlines); `"\n"` + the fixture (a
    leading blank line); `"x"` + the fixture; the fixture twice; the fixture with one empty body
    line inserted.
  Empty input gives `errors.Is(err, ErrArmor)`. Positive control: every fixture `.sig` parses with
  a nil error. *Kills* MUT-SG-ZEROVALUE, MUT-SG-BOUNDS (as a panic), MUT-SG-ARMOR.
- **SG6 `TestFingerprintMatchesSourceTxt`**: read `SOURCE.txt`. Line 1 starts with `OpenSSH_`. The
  line containing `(ED25519)` has `strings.Fields(...)[1] == fixtureKey(t, "ed25519.pub").Fingerprint()`,
  and the line containing `(ECDSA)` matches `p256.pub` the same way. Both fingerprints start with
  `SHA256:` and contain no `=`. *Kills* MUT-FP-PAD, MUT-FP-TEXT.
- **SG7 `TestParseAuthorizedKeyRefusals`**: `ed25519.pub` gives `Type == "ssh-ed25519"`, and
  `p256.pub` gives `Type == "ecdsa-sha2-nistp256"`. A two-field line (the comment dropped) parses.
  Each of these arms gives an error matched with `errors.Is`:
  - `""` → ErrAuthorizedKey;
  - `"ssh-ed25519"` → ErrAuthorizedKey;
  - the `ed25519.pub` line twice → ErrAuthorizedKey;
  - `"ssh-ed25519 !!!notbase64"` → ErrBase64;
  - `namespaces="x"` prepended to `ed25519.pub` → ErrBase64;
  - outer type `ssh-ed25519` with `p256.pub`'s base64 → ErrKeyType;
  - outer and inner type `sk-ssh-ed25519@openssh.com` (the blob made by rewriting the first
    string of `ed25519.pub`'s decoded blob, then re-encoded) → ErrKeyType.
  *Kills* MUT-AK-OUTER, MUT-SG-SKTYPE.

**Acceptance**
```
test -z "$(gofmt -l host/approval)"
go build ./... && go vet ./host/approval/
test "$(LC_ALL=C ls -A host/approval/testdata/sshsig | tr '\n' ' ')" = "SOURCE.txt ed25519.pub ed25519.sig generate.sh git-ns-ed25519.sig message.hex p256.pub p256.sig "
test -x host/approval/testdata/sshsig/generate.sh
test -z "$(grep -rl 'PRIVATE KEY' host/approval .gitattributes)"
bash -c 'for f in host/approval/testdata/sshsig/*; do LC_ALL=C tr -d "\000" < "$f" | cmp -s - "$f" || { echo "NUL byte in $f"; exit 1; }; done'
git check-attr text -- host/approval/testdata/sshsig/ed25519.sig   # must print: host/approval/testdata/sshsig/ed25519.sig: text: unset
bash -c 'xxd -r -p host/approval/testdata/sshsig/message.hex | ssh-keygen -Y check-novalidate -n world-approval@sunholo -s host/approval/testdata/sshsig/ed25519.sig'
bash -c 'xxd -r -p host/approval/testdata/sshsig/message.hex | ssh-keygen -Y check-novalidate -n world-approval@sunholo -s host/approval/testdata/sshsig/p256.sig'
go test ./host/approval/ -run '^(TestSSHSIGFixturesAreTheGeneratedBytes|TestParseSSHSIGFailsLoudly|TestFingerprintMatchesSourceTxt|TestParseAuthorizedKeyRefusals)$' -count=1 -v
go test ./host/approval/ -count=1
```

### M2b — the verifier half of `sshsig.go` (design M2, part 2) · ~110 prod + ~230 test LOC

**Files**: `host/approval/sshsig.go` (add `VerifySSHSIG`, `signedData`, `parseMpint` and
`VerifyApproval` as in §2.3) and `host/approval/sshsig_test.go`.

**Tests**
- **SG2 `TestSSHSIGFixturesVerify`**: for `(ed25519.sig, ed25519.pub, "ssh-ed25519")` and
  `(p256.sig, p256.pub, "ecdsa-sha2-nistp256")`: `s, err := ParseSSHSIG(sig)` gives a nil error and
  non-nil `s`, with `s.PublicKey.Type`, `s.SigType` equal to the type, `s.HashAlg == "sha512"` and
  `s.Namespace == Namespace`. `bytes.Equal(s.PublicKey.Blob, fixtureKey(t, pub).Blob)` holds.
  `VerifySSHSIG(s, Namespace, fixtureMessage(t), key)` is nil. *Kills* MUT-SG-ECDSA512,
  MUT-SG-REFUSEALL, and a parser that returns a zero value (the type assertions).
- **SG3 `TestSSHSIGRefusals`**: a table of `{name, armored, message, key, want}`. Each arm runs
  `s, err := ParseSSHSIG(armored)`, and if `err == nil` then `err = VerifySSHSIG(s, Namespace,
  message, key)`. It requires `errors.Is(err, want)`, plus a reason substring where one is given.
  Here `msg` = `fixtureMessage(t)`, `ed` = `ed25519.sig`, `p` = `p256.sig`, and `edK`/`pK` are the
  `.pub` keys.

  | Arm | armored | message | key | want (reason) | Kills |
  |---|---|---|---|---|---|
  | git-namespace | `git-ns-ed25519.sig` | msg | edK | ErrNamespace | MUT-SG-NS |
  | wrong-key | ed | msg | pK | ErrKeyMismatch | MUT-SG-KEYEQ |
  | msg-flip-ed / -p256 | ed / p | `flipByte(msg, 0)` | edK / pK | ErrBadSignature | MUT-SG-ACCEPTALL |
  | challenge-flip-ed / -p256 | ed / p | `flipByte(msg, 50)` | edK / pK | ErrBadSignature | MUT-SG-ACCEPTALL |
  | expiry-flip-ed / -p256 | ed / p | `flipByte(msg, 89)` | edK / pK | ErrBadSignature | MUT-SG-ACCEPTALL |
  | sig-flip-ed / -p256 | `withLastSigByteFlipped(ed / p)` | msg | edK / pK | ErrBadSignature | MUT-SG-ACCEPTALL |
  | sk-type | `withKeyType(ed, "sk-ssh-ed25519@openssh.com")` | msg | edK | ErrKeyType (`security-key`) | MUT-SG-SKTYPE |
  | rsa-type | `withKeyType(ed, "rsa-sha2-512")` | msg | edK | ErrKeyType | (default arm) |
  | p384-type | `withKeyType(p, "ecdsa-sha2-nistp384")` | msg | pK | ErrKeyType | (default arm) |
  | curve | `withCurveName(p, "nistp384")` | msg | pK | ErrKeyType | MUT-SG-CURVE |
  | point | `withPointByteFlipped(p)` | msg | pK | ErrKeyType (`point`) | (stdlib point check) |
  | short-ed-key | `withShortEd25519Key(ed)` | msg | edK | ErrKeyType | MUT-SG-EDLEN |
  | hash-sha1 | `withHashAlg(ed, "sha1")` | msg | edK | ErrHashAlg | MUT-SG-HASHALG |
  | version-2 | `withVersion(ed, 2)` | msg | edK | ErrVersion | MUT-SG-VERSION |
  | magic | `withMagic(ed, "SSHSIH")` | msg | edK | ErrMagic | MUT-SG-MAGIC |
  | reserved | `withReserved(ed, "x")` | msg | edK | ErrReserved | MUT-SG-RESERVED |
  | trailing | `withTrailingByte(ed)` | msg | edK | ErrTrailing | MUT-SG-TRAILING |
  | sigtype-ed | `withSigType(ed, "ecdsa-sha2-nistp256")` | msg | edK | ErrSigType | MUT-SG-SIGTYPE |
  | sigtype-p256 | `withSigType(p, "ssh-ed25519")` | msg | pK | ErrSigType | MUT-SG-SIGTYPE |
  | mpint-negative | `withNegativeMpint(t, p)` | msg | pK | ErrBadSignature (`negative mpint`) | MUT-SG-NEGMPINT |
  | mpint-non-minimal | `withNonMinimalMpint(t, p)` | msg | pK | ErrBadSignature (`non-minimal mpint`) | MUT-SG-NONMIN |

  It also includes `VerifySSHSIG(nil, Namespace, msg, edK)` → ErrBadSignature.
- **SG5 `TestVerifyApprovalBindsTheRequest`**: for both `(ed25519.sig, edK)` and `(p256.sig, pK)`,
  `VerifyApproval(sig, fixtureDigest(t), fixtureChallenge(), fixtureExpiry, key)` is nil. With the
  digest's byte 0 flipped, the challenge's byte 0 flipped, or `fixtureExpiry+1`, each gives
  `errors.Is(err, ErrBadSignature)`. A negative expiry gives `errors.Is(err, ErrNegativeExpiry)`.
  `git-ns-ed25519.sig` gives `ErrNamespace`. *Kills* MUT-VA-NS, MUT-VA-NOVERIFY.

**Acceptance**
```
test -z "$(gofmt -l host/approval)"
go build ./... && go vet ./host/approval/
go test ./host/approval/ -run '^(TestSSHSIGFixturesVerify|TestSSHSIGRefusals|TestVerifyApprovalBindsTheRequest)$' -count=1 -v
# design §8's M2 gate (plus SG7):
go test ./host/approval/ -run '^TestSSHSIG|^TestParseSSHSIG|^TestVerifyApproval|^TestFingerprint|^TestParseAuthorizedKey' -count=1
go test -race ./host/approval/ -count=1
```

### M3 — the `ssh-keygen` cross-check and its CI step (design M3) · ~170 test + ~20 YAML LOC

**Files**: new `host/approval/sshsig_crosscheck_test.go`; `.github/workflows/ci.yml` (one step).

**XC1 `TestSSHSIGCrossCheckWithSSHKeygen`**, with subtests `ed25519` and `ecdsa-p256`:
- At the top of the parent, before any subtest: if `exec.LookPath("ssh-keygen")` fails, then with
  `os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1"` call
  `t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")`; otherwise call
  `t.Skip("SKIP COUNTED: ssh-keygen not found; the SSHSIG cross-check did NOT run")`.
- Helper `runSSHKeygen(t *testing.T, stdin []byte, args ...string) (stdout []byte, rc int)`:
  `exec.CommandContext` with a 30 s `context.WithTimeout`;
  `cmd.Env = append(os.Environ(), "SSH_AUTH_SOCK=")`; `cmd.Stdin = bytes.NewReader(stdin)` (or nil);
  stdout and stderr captured, and stderr logged with `t.Logf` when `rc != 0`. A timeout is
  `t.Fatal`. `rc` comes from `*exec.ExitError`, and any other error is `t.Fatal`.
- Per subtest (`args` = `-t ed25519` or `-t ecdsa -b 256`, `wantType` = `ssh-ed25519` or
  `ecdsa-sha2-nistp256`), with `dir := t.TempDir()` (test keys only, constraint 1):
  1. `ssh-keygen -q <args> -N "" -C xc1-<name> -f dir/key` and the same for `dir/other`. Each must
     give rc 0.
  2. `msg := fixtureMessage(t)`. Then
     `good := runSSHKeygen(t, msg, "-Y", "sign", "-n", Namespace, "-f", dir/key)`,
     `gitSig := … "-n", "git" …` and `otherSig := … "-f", dir/other`. All sign **from stdin**.
     Each must give rc 0, a **non-empty** stdout, and `ParseSSHSIG` nil error with
     `PublicKey.Type == wantType` (a silent no-write fails).
  3. `pub := ParseAuthorizedKey(os.ReadFile(dir/key.pub))`. Write `dir/allowed` as
     `"xc1 namespaces=\"world-approval@sunholo\" " + fields[0] + " " + fields[1] + "\n"` from
     `key.pub`.
  4. Arms, as `{name, sig, message}`: `good` (good, msg); `tampered-message` (good,
     `flipByte(msg, 0)`); `git-namespace` (gitSig, msg); `another-key` (otherSig, msg);
     `flipped-signature` (`withLastSigByteFlipped(t, good)`, msg). For each arm, write the sig to
     `dir/<name>.sig` and run
     `ssh-keygen -Y verify -f dir/allowed -I xc1 -n world-approval@sunholo -s dir/<name>.sig`
     with the message on stdin. `ours := ParseSSHSIG` then `VerifySSHSIG(s, Namespace, message,
     pub)`. Require `(rc == 0) == (ours == nil)`, naming the arm. Also require that `good` has
     `rc == 0 && ours == nil` (non-vacuous), and that each of the other four has `rc != 0`
     (instrument sanity).
  5. `ssh-keygen -l -f dir/key.pub` stdout contains `pub.Fingerprint()`.
- No `t.Parallel()`. No nested subtests: the arms are loop iterations.
- *Kills*: MUT-SG-REFUSEALL (the `good` arm), MUT-SG-ACCEPTALL (the tamper arms), MUT-FP-PAD
  (step 5), MUT-XC-SKIP and MUT-XC-REQUIRE (the CI step and the negative control below).

**CI step.** In `.github/workflows/ci.yml`, job `go-verify`, insert this block **after line 136**
(the blank line after the `go version` step, `:132-135`) and **before** `:137`
`- name: Install PINNED v0.41.0 ailang (AILANG_BIN — one pin, current AILANG)`. Keep one blank line
after the block. The step body is the design's §2.8 text, byte for byte. Only the four comment
lines are added. The planner ran `actionlint` (with shellcheck) on exactly this insertion: rc 0.

```yaml
      # Row 201 (w-world-approvals-by-key-a0a1 §2.8, test XC1): the SSHSIG verifier is cross-checked
      # against REAL ssh-keygen. WORLD_REQUIRE_SSH_KEYGEN=1 turns a missing ssh-keygen into a red,
      # and the PASS grep turns a SKIP into a red. Whether ubuntu-24.04 ships ssh-keygen is
      # UNVERIFIED; if it does not, this step fails naming the apt-get line.
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

**Step simulation** (runs the step body exactly as committed, on this rig; the controller's shell
is zsh, so it goes through `bash`):

```bash
bash -c 'set -euo pipefail
export RUNNER_TEMP="$(mktemp -d)"; trap "rm -rf \"$RUNNER_TEMP\"" EXIT
awk '"'"'/- name: Row 201 SSHSIG cross-check/{f=1;next} f&&/^      - name:/{exit} f'"'"' .github/workflows/ci.yml \
  | sed -n '"'"'/^        run: |$/,$p'"'"' | sed '"'"'1d; s/^          //'"'"' > "$RUNNER_TEMP/step.sh"
test -s "$RUNNER_TEMP/step.sh"
bash "$RUNNER_TEMP/step.sh"'
```

**Acceptance**
```
test -z "$(gofmt -l host/approval)"
go build ./... && go vet ./host/approval/
WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -run '^TestSSHSIGCrossCheckWithSSHKeygen$' -count=1 -v   # prints --- PASS for the parent and both subtests
<the step simulation above>                                                                                 # rc 0
# negative control 1 (ssh-keygen absent + required → red). The go binary's own dir has no ssh-keygen on this rig (planner-checked):
bash -c 'PATH="$(dirname "$(command -v go)")"; export PATH; ! command -v ssh-keygen && ! WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -run "^TestSSHSIGCrossCheckWithSSHKeygen$" -count=1'
# negative control 2 (absent, not required → a COUNTED skip):
bash -c 'PATH="$(dirname "$(command -v go)")"; export PATH; go test ./host/approval/ -run "^TestSSHSIGCrossCheckWithSSHKeygen$" -count=1 -v 2>&1 | grep -q "SKIP COUNTED: ssh-keygen not found"'
actionlint .github/workflows/ci.yml
go test ./host/approval/ -count=1
```

---

## 4. Mutation table (68 mutants; the controller and judge run them)

Procedure: apply **one** mutant (a single edit), run the named test(s), expect **FAIL** (a red
test, or a panic counted as a failure), then restore the file and `cmp` it against the saved copy.
Fixture mutants (MUT-FX-*) are applied to a copy and restored the same way. **Never commit a
mutant.**

| Id | Mutant (single edit) | Must fail |
|---|---|---|
| MUT-RQ-ORDER | swap the declarations of `Reuse` and `Degraded` in `MintRequest` | RQ1 |
| MUT-RQ-HTML | delete `enc.SetEscapeHTML(false)` in `Canonical` | RQ1 |
| MUT-RQ-NEWLINE | `Canonical` returns `buf.Bytes()` untrimmed | RQ1 |
| MUT-RQ-NOREENCODE | `Decode`: delete the `!bytes.Equal(c, b)` refusal | RQ2 |
| MUT-RQ-HIDDEN | `Model`'s tag → `json:"-"` | RQ3 (and RQ1) |
| MUT-V-SCHEMA | `Validate` step 1 deleted | RQ4 schema |
| MUT-V-STRING | `Validate` step 2 deleted | RQ4 utf8, nul-* |
| MUT-V-EPISODE | `Validate` step 3 deleted | RQ4 episode-* |
| MUT-V-PATH | `Validate` step 4 (both loops) deleted | RQ4 path arms |
| MUT-V-PATHFIELD | `profiles_dir` removed from the required-path list | RQ4 profiles-dir-rel |
| MUT-V-GITREUSE | `Validate` step 5 deleted | RQ4 git-missing, git-on-reuse |
| MUT-V-NOGRANTS | `Validate` step 6 deleted | RQ4 no-grants-* |
| MUT-V-GRANTEMPTY | the `Effect == "" \|\| Scope == ""` branch deleted | RQ4 effect-empty, scope-empty |
| MUT-V-BUDGET | the `Budget < 0` branch deleted | RQ4 budget-neg |
| MUT-V-TTL | `Validate` step 8 deleted | RQ4 ttl-neg |
| MUT-V-HALFPROFILE | `Validate` step 9 deleted | RQ4 half-profile-* |
| MUT-V-AGENT | `Validate` step 10 deleted | RQ4 agent |
| MUT-V-MCP | `Validate` step 11 deleted | RQ4 mcp-* |
| MUT-SM-NUL | `SignedMessage` omits the `0x00` | RQ5, SG1 |
| MUT-SM-LE | `SignedMessage` uses `binary.LittleEndian` | RQ5, SG1 |
| MUT-SM-NEGEXP | the `expiry < 0` refusal deleted | RQ5 |
| MUT-DIFF-EMPTY | `Diff` body → `return ""` | RQ8 |
| MUT-AP-IMPORT | add `_ "os"` to `request.go`'s imports | AP1 |
| MUT-AD-CONSUMED | `Admit` step 1 deleted | RQ6 consumed* |
| MUT-AD-CLOCK | `Admit` step 2 deleted | RQ6 clock-behind |
| MUT-AD-EXPIRED | `Admit` step 3 deleted | RQ6 at-expiry |
| MUT-AD-DIGEST | `Admit` step 4 deleted | RQ6 digest |
| MUT-AD-GE | `in.Now >= in.ExpiresAt` → `in.Now > in.ExpiresAt` | RQ6 at-expiry |
| MUT-AD-ORDER | move the `Consumed` check after the expiry check | RQ6 consumed+expired |
| MUT-AD-ZERO | `Admitted Verdict = iota + 1` → `Admitted Verdict = iota` | RQ6 zero-verdict assertion |
| MUT-MM-CREDID | `MintMatches` step 1 deleted | RQ7 cred-id |
| MUT-MM-EPISODE | step 2 deleted | RQ7 episode |
| MUT-MM-GRANTSCOUNT | step 3's per-element comparison deleted (count only) | RQ7 grants-reordered, -budget, -scope |
| MUT-MM-GRANTEXP | step 4 deleted | RQ7 grant-exp-* |
| MUT-MM-EXPIRES | step 5 deleted | RQ7 expires-off-by-one |
| MUT-MM-CREATEDLOW | step 6 deleted | RQ7 created-before |
| MUT-MM-CREATEDHIGH | step 7 deleted | RQ7 created-after |
| MUT-MM-PROJECT | step 8 deleted | RQ7 *-unbound-*, unbound-req-* |
| MUT-MM-DIGEST | step 9 deleted | RQ7 other-digest |
| MUT-FX-EDIT | change one base64 character in the middle of `ed25519.sig` | SG1 |
| MUT-FX-CRLF | convert `p256.pub` to CRLF | SG1 |
| MUT-FX-EXTRA | add a file `testdata/sshsig/extra.pub` (a copy of `p256.pub`) | SG1 |
| MUT-SG-ZEROVALUE | `ParseSSHSIG`: on a `parseSSHSIGBlob` error, `return &SSHSig{}, nil` | SG4 |
| MUT-SG-BOUNDS | `readString`: delete the `uint64(n) > uint64(len(r.b))` check | SG4 (panic) |
| MUT-SG-ARMOR | `strings.TrimSuffix(s, "\n")` → `strings.TrimSpace(s)` | SG4 armor arms |
| MUT-FP-PAD | `Fingerprint`: `RawStdEncoding` → `StdEncoding` | SG6, XC1 |
| MUT-FP-TEXT | `Fingerprint` hashes `[]byte(base64.StdEncoding.EncodeToString(k.Blob))` | SG6 |
| MUT-AK-OUTER | `ParseAuthorizedKey` step 5 deleted | SG7 outer-type arm |
| MUT-SG-SKTYPE | `parseKeyBlob`: `case keyEd25519:` → `case keyEd25519, "sk-ssh-ed25519@openssh.com":` | SG3 sk-type, SG7 |
| MUT-SG-NS | `VerifySSHSIG` step 1 deleted | SG3 git-namespace |
| MUT-SG-KEYEQ | `VerifySSHSIG` step 2 deleted | SG3 wrong-key |
| MUT-SG-CURVE | `parseKeyBlob`: the `nistp256` curve check deleted | SG3 curve |
| MUT-SG-EDLEN | `parseKeyBlob`: the `len(pk) != ed25519.PublicKeySize` check deleted | SG3 short-ed-key |
| MUT-SG-SIGTYPE | `VerifySSHSIG` step 4 deleted | SG3 sigtype-* |
| MUT-SG-ECDSA512 | the P-256 branch hashes `sd` with `sha512.Sum512` (first 32 bytes) | SG2, XC1 |
| MUT-SG-ACCEPTALL | `VerifySSHSIG` step 5 → `return nil` | SG3 flip arms, SG5, XC1 |
| MUT-SG-REFUSEALL | `VerifySSHSIG` step 5 → `return ErrBadSignature` | SG2, SG5, XC1 good |
| MUT-SG-NEGMPINT | `parseMpint`'s negative check deleted | SG3 mpint-negative |
| MUT-SG-NONMIN | `parseMpint`'s non-minimal check deleted | SG3 mpint-non-minimal |
| MUT-SG-RESERVED | the `len(reserved) != 0` check deleted | SG3 reserved |
| MUT-SG-HASHALG | the hash-algorithm check in `parseSSHSIGBlob` deleted | SG3 hash-sha1 |
| MUT-SG-VERSION | the version check deleted | SG3 version-2 |
| MUT-SG-MAGIC | the magic check deleted | SG3 magic |
| MUT-SG-TRAILING | `parseSSHSIGBlob` step 3's `r.end()` deleted | SG3 trailing |
| MUT-VA-NS | `VerifyApproval` passes `"git"` instead of `Namespace` | SG5 |
| MUT-VA-NOVERIFY | `VerifyApproval` returns `nil` after a successful parse | SG5 |
| MUT-XC-SKIP | the first line of XC1 → `t.Skip("mutant")` | the M3 step simulation (the PASS grep reds) |
| MUT-XC-REQUIRE | XC1 ignores `WORLD_REQUIRE_SSH_KEYGEN` (always skips when absent) | M3 negative control 1 |

---

## 5. Files (final `git status --short` after M3)

| File | Milestone | Change |
|---|---|---|
| `host/approval/request.go` | M1a | new |
| `host/approval/request_test.go` | M1a | new: RQ1–RQ5, RQ8 |
| `host/approval/imports_test.go` | M1a | new: AP1 |
| `host/approval/admit.go` | M1b | new |
| `host/approval/admit_test.go` | M1b | new: RQ6, RQ7 |
| `host/approval/sshsig.go` | M2a, M2b | new |
| `host/approval/sshsig_test.go` | M1a (fixture constants), M2a, M2b | new: SG1–SG7, helpers, transforms |
| `host/approval/testdata/sshsig/{SOURCE.txt,ed25519.pub,ed25519.sig,generate.sh,git-ns-ed25519.sig,message.hex,p256.pub,p256.sig}` | M2a | new, generated (`generate.sh` hand-written, mode 0755) |
| `.gitattributes` | M2a | new, one line |
| `host/approval/sshsig_crosscheck_test.go` | M3 | new: XC1 |
| `.github/workflows/ci.yml` | M3 | one step, inserted at `:137` |
| `design_docs/planned/w-world-approvals-by-key-a0a1-sprint-plan.md`, `design_docs/planned/sprint_w-world-approvals-by-key-a0a1.json` | planner | this plan |

---

## 6. Sprint gates (controller, after M3; all green at base, §1.1)

```
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh
go build ./... && go vet ./host/approval/
go test ./host/approval/ -count=1
go test -race ./host/approval/ -count=1
go test -run '^$' ./...                                     # compile fence for every _test.go
go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
go test ./host/verifygate/ -run '^TestNoParallelWriteForkPackagesOutsideVerifygate$' -count=1
go test ./host/pins/ -count=1
WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -run '^TestSSHSIGCrossCheckWithSSHKeygen$' -count=1 -v
actionlint .github/workflows/ci.yml
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./...       # full suite (≈6 min at base)
# AFTER the controller commits (the gate reads HEAD's committed blobs, not the working tree):
VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh   # includes scripts/tracked_binary_gate.sh (§1.4 F1)
git status --short                                          # only the §5 files
```

---

## 7. Risks

- **R1 (design R1)**: whether the `ubuntu-24.04` image has `ssh-keygen` is **UNVERIFIED**. If it
  does not, the step fails with `::error::ssh-keygen missing (apt-get install openssh-client)`.
  The fix is one `sudo apt-get install -y openssh-client` line in the same step.
- **R2**: the CI `go build + test gate` step (`verify_go.sh`) and the `go-race-verify` job also run
  XC1 without `WORLD_REQUIRE_SSH_KEYGEN`. There it can at most skip-and-count. Only the new step
  is load-bearing for XC1.
- **R3 (D7)**: there is about a 25% chance per `generate.sh` run that `p256.sig` has no high-bit r
  or s. SG1(f) catches it, and the fix is to re-run the generator. Each re-run means new hashes
  in the same commit.
- **R4**: the fixtures are not reproducible byte for byte (fresh keys, and ECDSA is randomised). SG1
  pins what was produced, and `SOURCE.txt` records the tool and the time. That is by design.
- **R5 (§1.4 F3)**: sprint 2 must reconcile "`MintMatches` before the write" with the
  `created_at ≤ in-transaction consume reading` bound.
- **R6**: OpenSSH output formats (armor width 70, `ssh-keygen -l` layout) were measured on
  `OpenSSH_10.3p1`. The parser does not depend on line width. XC1 on Linux uses whatever
  `ssh-keygen` the runner has. That is the point of the check.

---

## 8. Acceptance criteria (the judge scores against these)

1. `host/approval` exists with exactly the §2.1–§2.3 exported names and signatures. The verdict
   constant is `Admitted` (`iota + 1`).
2. Production files import only the standard library and none of the §0.3 forbidden packages. AP1
   passes and its positive control fires.
3. `Canonical(fixtureRequest())` equals `fixtureCanonicalJSON`, and the digest equals
   `fixtureDigestHex`. Neither golden was re-derived from the implementation.
4. `Decode` refuses unknown fields and every non-canonical byte string (`ErrNotCanonical`).
5. `Validate` has the 11 ordered branches of §2.1, each with its own sentinel, and RQ4 has one or
   more arms per branch.
6. `SignedMessage` is 90 bytes, `MessagePrefix ‖ 0x00 ‖ digest ‖ challenge ‖ BE64(expiry)`, and
   refuses a negative expiry.
7. `Admit` follows the §2.3 order. `MintMatches` checks all seven columns plus the grants'
   `ExpiresAt`, and returns `*MintMismatchError` naming the column.
8. `ParseSSHSIG` returns `(nil, typed error)` on every input it cannot fully parse, including every
   prefix of every fixture, and never panics (SG4).
9. `VerifySSHSIG` runs the fixed sequence of §2.3. P-256 hashes the signed data with SHA-256. The
   mpint refusals are in place.
10. The fixtures are the 8 files of §5, generated by `generate.sh` with real `ssh-keygen` signing
    from stdin. No private key and no NUL byte are present. The three `-Y sign` lines are quoted
    above `fixtureSHA256`, and SG1 pins every file's SHA-256 plus the exact directory listing.
11. `.gitattributes` is exactly `host/approval/testdata/sshsig/* -text`.
12. Every refusal arm in SG3, SG4 and SG7 is built by a named transform from a real fixture. No
    signature or key bytes are typed by hand.
13. XC1 passes locally with `WORLD_REQUIRE_SSH_KEYGEN=1` (both subtests PASS). It fails when
    `ssh-keygen` is absent and required, and prints `SKIP COUNTED` when absent and not required.
14. The CI step sits in job `go-verify` between `go version` and `Install PINNED v0.41.0 ailang`,
    and its `run:` body is byte-identical to design §2.8. `actionlint` is clean, and the step
    simulation passes.
15. Each §4 mutant, applied alone, makes its named test(s) fail.
16. Every §6 gate passes, including `verify_go.sh` (the tracked-binary gate) after the commit.
    `git status --short` lists only the §5 files. No `cmd/`, `host/store`, `host/authority`,
    `host/daemon`, `.ail`, `scripts/`, `tools/launchd/*`, QUICKSTART or skill changes.
