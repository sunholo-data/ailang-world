# Sprint plan — w-world-approvals-by-key-a0a1, sprint 2 (iteration 374, row 201)

- **Design**: [w-world-approvals-by-key-a0a1.md](w-world-approvals-by-key-a0a1.md), **r3**. This
  sprint is the design's **M4–M6**: credentials store schema 4 (`mint_requests`), `episode start
  --request` plus `approve --list`, and `ailang-worldd approve <id>` by the typed `y`. It is
  preceded by **M3b**, a small test-only milestone in `host/approval` that closes the sprint-1
  judge's three test gaps. The design is **not** reopened. Where this plan and the design differ,
  this plan was measured against the tree and wins. The differences are in §1.3, and the design
  premises found false are in §1.4.
- **Format model**: sprint 1's plan
  ([w-world-approvals-by-key-a0a1-sprint-plan.md](w-world-approvals-by-key-a0a1-sprint-plan.md)).
  Sprint 1 landed at merge `dba0e8b2`.
- **Planner base**: `ab4ae3aa` (origin/dev), worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter374-row201s2`, branch
  `sprint/world-iter374-row201-s2`.
- **Executor**: `claude-haiku-5-5`, run one or two milestones at a time. Every name, signature,
  constant, SQL statement, refusal string and assertion below is fixed. **Copy, do not paraphrase.**
- **Estimate**: about **2.25 days** of executor time over **ten** milestones, about **2,550 lines**
  (≈795 Go production, ≈1,710 Go test of which M3b is ≈120, ≈30 SQL and help text, ≈45 QUICKSTART;
  the per-milestone figures in §3 sum to this). The design says
  ~1,600 for M4–M6. The difference is M3b (≈120), the F3 transaction callback, and tests the
  design's rows imply but do not name: ST4b, ST8, BM1 and the authority-level AU1 (§1.3).
- **Touched**: `host/approval/*_test.go` (M3b only), `host/store/**`, `host/authority/mint.go`
  (+ one new test file), `host/daemon/credentials.go` (one log line), `cmd/ailang-worldd/**`,
  `docs/QUICKSTART.md` (§14 appended). **Not touched**: `host/approval/*.go` production files,
  `host/approval/testdata/**`, `.gitattributes`, `.github/workflows/ci.yml`, `host/broker/**`,
  `host/daemon/daemon.go` (no route), `cmd/world-publish/**`, `world/*.ail`, `packages/**`,
  `scripts/**`, `tools/launchd/*`, any skill, `~/.ailang/state/**`.

---

## 0. Binding constraints (each is restated in the milestone it governs)

1. **Fakes, temp stores and test keys only (D-WORLD-109 = A).** Every store is under
   `t.TempDir()`. Every daemon is the `stubDaemon` seam (`episode_test.go:71`). Every terminal is
   a fake `io.ReadWriteCloser`. Every key is made inside a test in `t.TempDir()` by `ssh-keygen`
   with `SSH_AUTH_SOCK=` emptied (the existing `runSSHKeygen`). Nothing opens `~/.ailang`,
   `~/.ssh`, a real credentials store, a real daemon or `/dev/tty`. No test mints into, or
   publishes to, anything that outlives the test. **Never run a bare `ssh-keygen`**: with no
   arguments it starts generating `~/.ssh/id_ed25519` interactively (the planner hit its prompt
   while probing; no file was written, checked with `ls -la ~/.ssh`).
2. **No hand-typed signature bytes, no committed binary blob.** M3b derives every new arm from a
   real fixture or a real in-test `ssh-keygen` signature by a **named transform**. No new file goes
   under `host/approval/testdata/`, so the SG1 pinned-hash table and the tracked-binary gate
   (`scripts/tracked_binary_gate.sh`, run by `verify_go.sh` with its `--allow-font` pins) are
   untouched.
3. **`provisionEpisode` stays the only fence-and-mint path.** `approve` reaches a mint only through
   it. `episode start --request` never calls it, never opens a terminal and never calls
   `authority.*`.
4. **Existing callers are unchanged.** `session new`, `session mint`, `episode start`, the launch
   preset path and resume keep `answerBy == 0` and `consume == nil`, so their y/N stays unbounded
   and they mint exactly as at HEAD. Their tests (§6 "must stay green") are **not edited**.
5. **No `t.Parallel()`** in any new test: `cmd/ailang-worldd`, `host/store`, `host/authority` and
   `host/approval` tests write files and fork (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
6. **Every test row names the mutant it kills** (§4: 83 mutants). The controller and the judge run
   them. A refusal is tested at the branch: deleting the branch turns exactly its arm red.
7. **Bisectable.** Each milestone ends compiling (`go build ./...` **and** `go vet <pkg>`, because
   `go build` does not compile `_test.go` files) with its own tests green and every earlier
   milestone's tests still green.
8. **Context roots.** Exactly two new production `context.Background()` roots:
   `cmd/ailang-worldd/mint_request.go|createMintRequest` (M5b) and
   `cmd/ailang-worldd/approve.go|runApprove` (M6b). No `TODO()`, no `WithoutCancel`.
9. **Standard library only** in the new code, plus this module's own packages.
   `TestDaemonDependencyAllowlist` must stay green.

---

## 1. Pre-flight

### 1.1 Baseline at `ab4ae3aa` (measured by the planner, 2026-11-10)

`AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`, commit `24ee108`), `go1.26.6
darwin/arm64`, `OpenSSH_10.3p1`.

| Command | Result |
|---|---|
| `git log --oneline -1`; `git status --short` | `ab4ae3aa Merge pull request #534 …`; clean |
| `go build ./...` | rc 0 |
| `go vet ./host/store/ ./cmd/ailang-worldd/ ./host/authority/ ./host/approval/` | rc 0 |
| `go test -run '^$' ./...` (compile fence for every `_test.go`) | no failure line |
| `go test ./host/approval/ -count=1` | `ok … 0.399s` |
| `WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -run '^TestSSHSIGCrossCheckWithSSHKeygen$' -count=1 -v` | `--- PASS: TestSSHSIGCrossCheckWithSSHKeygen`, both subtests run |
| `go test ./host/store/ -run 'Credentials\|TestProductionContextRoots\|TestEpisodeStarts\|TestMigrateCredentials\|TestPreSplitStoreOpenModes' -count=1` | `ok … 0.432s` |
| `go test -race ./host/store/ -run '^TestCredentials' -count=1` | `ok … 2.666s` (the race detector works on `host/store` here) |
| `go test ./host/authority/ -count=1` | `ok … 1.044s` |
| `go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1` | `ok … 0.709s` |
| `go test ./host/store/ -run '^TestProductionContextRoots$' -count=1 -v` | `--- PASS: TestProductionContextRoots` |
| `go test ./cmd/ailang-worldd/ -run '^(TestSessionMint_\|TestSessionNew\|TestEpisodeStart\|TestEveryVerbHasHelp\|TestHelpIsNeverBlockedByTheRegistryKeyGuard\|TestSessionMintGrantsExpireWithTheSession)' -count=1 -v` | 32 `--- PASS`, 0 `--- SKIP`, 0 `--- FAIL` |
| `AILANG_BIN=… ./scripts/verify_ail.sh` | rc 0: `verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `VERIFY_GO_LEGS=plain AILANG_BIN=… ./scripts/verify_go.sh` | **rc 0**: tracked-binary gate 24/24 self-test cases ok, `3 pinned font(s) validated`, `0 other binary blobs among 6578 tracked files`; 64 evidence tests passed exactly once; `go test ./... -count=1` **36 packages ok, 0 FAIL** (`cmd/ailang-worldd` 187 s, `host/daemon` 355 s, `host/broker` 302 s); `✓ go gate PASSED (plain)` |
| `gh pr list --state open` | `[]` (no open PR to conflict with) |

### 1.2 Codebase claims checked (the design's line numbers are still exact)

| # | Claim | Evidence (command → output) |
|---|---|---|
| C1 | The five files the design measured are byte-identical at `ab4ae3aa`, so its line numbers hold | `shasum -a 256 cmd/ailang-worldd/{session,session_new,episode}.go host/store/credentials.go host/authority/mint.go \| cut -c1-16` → `3566d5377c9fc01a`, `ef5163a5dc6e7aa5`, `5588c8659f51669b`, `3b4d630f89676737`, `1dfdf1080000889f` (= design V1) |
| C2 | Re-derived lines in `session_new.go` | `grep -n` → `:26` "no other code creates a credential"; `:181 type mintOptions`; `:203 type mintRequest`; `:228 func prepareMint`; `:245 r := &mintRequest{mintOptions: o, credPath: credPath}` (the flags are kept as typed); `:359 checkBindingFits`, `:372 if ro.Version() < store.CredentialsSchemaVersion`; `:396 func provisionEpisode`; `:399 openTerminal`; `:408 now := env.now()`; `:430 if r.presetPrompt != ""`; `:438 if !isYes(readTermLine(term))`; `:444 context.WithTimeout(context.Background(), sessionNewBudget)`; `:479/:481/:483` the three mint calls; `:522 if created {` |
| C3 | `sessionEnv` has two fields and no randomness or timer | `sed -n 24,47p cmd/ailang-worldd/session.go` → `openTerminal func() (io.ReadWriteCloser, error)`, `now func() int64`; `realSessionEnv` at `:40`; `readTermLine` at `:469`, `isYes` at `:475` |
| C4 | `episode start` order | `grep -n` → `:147-148` seams `episodeSelfBuild`, `fetchDaemonHealth`; `:205 launch := fs.Bool`; `:227 fail := func` (returns `exitUsage`); `:232 resumeSet := false`, `:234 if resumeSet {`; `:245 if *preset != "" && *preset != "se-tools" {`; `:305 prepareMint(…)`; `:335 checkDaemon(health, req, hostPort, self)`; `:338-340 req.afterToken = … writeMCPConfig(mcpFile, hostPort, token)`; `:343 provisionEpisode`; `:348 printLaunch` |
| C5 | Exit codes | `sed -n 70,75p cmd/ailang-worldd/main.go` → `exitOK = 0`, `exitUsage = 1`, `exitFatal = 2`. The decline at `session_new.go:438-440` returns `exitUsage` |
| C6 | `provisionEpisode` already has four callers | `grep -n "provisionEpisode(" cmd/ailang-worldd/*.go \| grep -v _test` → `episode.go:343`, `episode_resume.go:210`, `episode_preset.go:301`, `session_new.go:606` |
| C7 | Credentials store at version 3, its table check and upgrade | `grep -n` → `credentials.go:40 const credentialsSchemaVersion = 3`; `:304 case known && ((version < 3 && len(tables) == 1 …) \|\| (version == 3 && len(tables) == 2 && tables[0] == "episode_starts" …))`; `:326 func upgradeCredentialsV3`; `:332 statements := []string{episodeStartsSQL}`; `:402 func (c *Credentials) mintSession`; `:412` the `c.version < credentialsSchemaVersion` gate for bound/start/resume rows; `:415 BeginTx`; `:428` the unbound `INSERT`. `episode_starts.go:10 const episodeStartsSQL = \`CREATE TABLE episode_starts (` |
| C8 | Every other gate that compares against the **current** version | `grep -rn "CredentialsSchemaVersion" --include="*.go" cmd host \| grep -v _test` → `session_new.go:372`, `episode_preset.go:235`, `host/daemon/credentials.go:61` (plus `credentials.go:412` above) |
| C9 | Tests that pin version 3 or its table list | `grep -rn '"episode_starts", "session_credentials"' --include="*_test.go" host cmd` → only `host/store/credentials_test.go:39`. `episode_starts_test.go:30` asserts `Version() != 2` for a v2 store (unaffected). No cmd/daemon test builds a v3 credentials store (`grep -rln "user_version" --include="*_test.go" cmd host/daemon host/authority` → only the pre-split **world**-store tests, `PRAGMA user_version = 4` of the world schema) |
| C10 | SQLite drops `IF NOT EXISTS` from `sqlite_master.sql` | scratch `sqlite3`: `CREATE TABLE IF NOT EXISTS u (a INT); SELECT sql …` → `CREATE TABLE u (a INT)`. So the schema file and a `CREATE TABLE mint_requests (` constant normalise to the same DDL (as `episode_starts` does today, `credentials_test.go:43`) |
| C11 | `requireDeadline` and the methods-need-a-deadline test | `store.go:250 func requireDeadline`; `credentials_test.go:220 TestCredentialsMethodsRequireADeadline` (a hand-written map of 4 methods) |
| C12 | The writer lock conflicts inside one process | `writer_lock_unix.go:41 syscall.Flock(…, LOCK_EX\|LOCK_NB)` on a fresh `os.OpenFile` each time, so two `OpenCredentials` in one process contend (ST6 relies on it) |
| C13 | Mint core | `host/authority/mint.go:57 func mintRecorded`; `:70-76` token and `hashHexToken`; `:81-95` the row; `:96-100` the store call; `resolver.go:207 func hashHexToken` |
| C14 | `broker.Capability` and `daemon.BuildIdentity` | `host/broker/decide.go:15-20` `{Effect string; Scope string; ExpiresAt int64; Budget int64}`, no tags; `host/daemon/build.go:27-38` `{VCSRevision, VCSModified, ExecutableSHA256, GoVersion}` with a `String()` at `:44` |
| C15 | Help plumbing | `help_test.go:10-14 everyVerb`, `:37-45` the usage want-list; `keyguard_test.go:17` **also** iterates `everyVerb`; `keyguard.go:38 helpOnly`, `:71 ownVerbHelp`; `main.go:650 ownHelpVerbs`, `:659 genericVerbHelp`, `:701 usageLinesFor`; `case "session"` refuses the global `--addr` at `:429-438` |
| C16 | Context-root inventory | `host/store/context_roots_test.go:17-37`; it already has `cmd/world-publish/main.go\|runApprove\|Background` (another package, another file key) |
| C17 | Test seams that exist | `episode_test.go:71 stubDaemon`, `:97 spyEnv`, `:51 matchingHealth` (calls `pinnedToolRef` and **`t.Skipf`s** on a platform with no tool pin), `:80 newEpisodeRun`; `session_test.go:29 testMintEnv`, `:38 noTerminalEnv`; `session_new_test.go:25 episodeFixture`, `:462 promptTerm` |
| C18 | No test name collides | `grep -rn "func TestApprove\|func TestMintRequest\|func TestEpisodeStartRequest\|func TestReadTermLine\|func TestMintForRequest\|func TestPutMintRequest\|func TestMintSessionForRequest\|func TestBuildMintRequest\|func TestConcurrentApprovals\|func TestExistingCallers\|func TestSSHSIGStrict" host cmd` → nothing; control `grep -rn "func TestEpisodeStart\|func TestStartLaunch\|func TestEpisodeResume" cmd/ailang-worldd/*_test.go \| wc -l` → `21` |
| C19 | No `t.Parallel()` in the touched test packages | `grep -rn "t.Parallel()" cmd/ailang-worldd/*_test.go host/store/*_test.go host/authority/*_test.go \| wc -l` → `0` |
| C20 | QUICKSTART tests cut only by named headings | `grep -rn QUICKSTART --include="*_test.go"` → cuts at `### 10.`, `### 11.`, §9, "Watch the world live", the decision-packet section; none counts `## N.` headings. The file ends at `## 13.` (`:1687`), 1771 lines |
| C21 | `ssh-keygen -Y sign -O hashalg=sha256` works and `-Y verify` accepts it | scratch run, keys in a `mktemp -d` deleted on exit (re-checked absent): Ed25519 and P-256 both `sign rc=0`, the blob carries `sha256`, `Good "world-approval@sunholo" signature …`, `verify rc=0` |
| C22 | The non-zero-pad-bits transform is real | scratch Go program on `p256.pub` (104-byte blob, `%3 == 2`): `z5w=` → `z5x=`; `base64.StdEncoding` decodes it to the **same** 104 bytes; `StdEncoding.Strict()` → `illegal base64 data at input byte 139`. `p256.sig` blob is 263 bytes (`%3 == 2`); `ed25519.sig` 192 and `ed25519.pub` 51 (`%3 == 0`, no pad bits) |

### 1.2a The sprint-1 judge gaps, measured

| Gap | Command → output |
|---|---|
| (a) the SHA-256 hash path of SSHSIG is never exercised | `grep -c hashalg host/approval/{sshsig,sshsig_crosscheck,request,admit}_test.go` → `0 0 0 0`; control `grep -c Namespace host/approval/sshsig_crosscheck_test.go` → `5`. `signedData`'s `case "sha256"` (`sshsig.go:326-328`) has no caller with `HashAlg == "sha256"` |
| (b) three strictness branches have no killing test | `grep -n "Strict\|PadBits\|padding" host/approval/*_test.go` → only `request_test.go:96 TestDecodeIsStrict` (JSON, not base64); control `grep -n "Strict()" host/approval/sshsig.go` → `:166`, `:197`. The only `ErrTrailing` arm is `sshsig_test.go:598` `withTrailingByte`, which appends **after** the signature string (the outer `r.end()` at `:225`), not inside it (`:248`) or after the mpint pair (`:298`) |
| (c) parse-level `ErrHashAlg` is not asserted directly | `grep -n ErrHashAlg host/approval/*_test.go` → only `sshsig_test.go:594`, inside SG3, which falls through to `VerifySSHSIG` when the parse succeeds, so `signedData`'s own `ErrHashAlg` (`:330`) answers for a deleted parse check (`:232-234`) |

### 1.3 Planner decisions where the design leaves a choice or cannot be built literally

- **D1 — M3b fixtures.** The SHA-256 path is tested inside XC1 with keys in `t.TempDir()` (three
  new arms per key type, signed with `-O hashalg=sha256`). This changes no checked-in fixture, so
  `fixtureSHA256`, `SOURCE.txt` and the D7 lottery are untouched, and the existing CI step (which
  requires XC1's two subtests to `--- PASS`) makes the arms unskippable on `ubuntu-24.04` with no
  `ci.yml` edit. The strictness and parse-level arms use the checked-in fixtures (new test SG8).
- **D2 — F3 resolution: a check callback inside the transaction** (§1.4 F3). `RequestConsume`
  gains `Check func(consumedAt int64) error`. `MintSessionForRequest` reads `Clock()` once after
  `BeginTx`, calls `Check(thatReading)`, then runs the guarded `UPDATE`, then the `INSERT`.
  `authority.MintForRequest` sets `Check` to `approval.MintMatches(request, facts{…, ConsumedAt:
  consumedAt})`. So `MintMatches` still runs **before the write** of the session row, and its
  `created_at` upper bound is the in-transaction reading (sprint-1 risk R5, design §2.4 [r3]).
  Both `Clock` and `Check` are required: nil is `ErrIncompleteConsume` (fail closed).
- **D3 — `insertSessionRow`.** `mintSession`'s `INSERT` (`credentials.go:426-438`) moves verbatim
  into `func insertSessionRow(ctx context.Context, tx *sql.Tx, row SessionRow) error`, used by both
  `mintSession` and `MintSessionForRequest`. A pure extraction: no behaviour change.
- **D4 — Exit codes by class, not by the design's numbers** (§1.4 F1). Refusals before the
  credentials lock → `exitUsage` (as `episode start`'s `fail` and the decline do). Failures at or
  after the write (including `ErrRequestExpired` from the store) → `exitFatal` (as
  `provisionEpisode` does for a mint error). Tests assert the **named constants**.
- **D5 — The request-id grammar lives in `host/store`** as `MintRequestIDPattern`
  (`^mr-[0-9a-f]{16}$`), enforced by `PutMintRequest` (§1.4 F2) and used by `approve` for its
  grammar check.
- **D6 — All four store methods refuse a store below version 4** with
  `*OlderCredentialsStoreError` (the design names it for `PutMintRequest`; on a v3 store the other
  three would otherwise fail with a raw `no such table`).
- **D7 — `createMintRequest` opens the store with `openCredentialsForSession`**, the mint path's
  own open (`session.go:148`), so the pre-split move behaves exactly as for a mint.
- **D8 — One clock read before the prompt.** The `answerBy` bound reuses `provisionEpisode`'s
  `now := env.now()` (`:408`): nothing blocks between it and the prompt. The `env.now()` calls on
  `approve`'s path are therefore exactly, in order: **[1]** `runApprove`'s admission, **[2]**
  `provisionEpisode`'s `:408`, **[3]** `Clock()` inside the transaction. Tests rely on this order.
- **D9 — One renderer.** `renderRequest` prints the **full** digest for both `--request`'s output
  and `approve`'s pre-prompt block (the design asks for the first 16 characters at create; the full
  digest contains them, and one renderer cannot drift from itself).
- **D10 — `approve` has its own help** (S7 (a)): `approveHelp`, registered in `ownHelpVerbs`
  (`main.go:650`) and `ownVerbHelp` (`keyguard.go:71`), because `keyguard_test.go:17` iterates
  `everyVerb` and compares the help printed with and without the registry key byte for byte.
- **D11 — `approve` refuses an `env.after == nil`** up front (`exitFatal`), as a request path
  refuses `env.random == nil` (design §2.2). A test that forgets a seam fails loudly.
- **D12 — `approve` also re-checks that the MCP config does not exist yet** (`episode start` checks
  both files at `episode.go:293-300`; `prepareMint` checks only the session file).
- **D13 — Planner-added tests**: **SG8** (M3b strictness and parse-level arms), **ST4b** (a failed
  insert or a refusing `Check` leaves the request pending), **ST8** (bad request rows; the weak SQL
  `CHECK` as its instrument control), **ST2b** (a v4 store missing `mint_requests` is refused),
  **BM1** (relative inputs are made absolute), **AU1** (`MintForRequest` checks inside the
  transaction), **CL1b** (a nil `random` refuses).

### 1.4 Premises in the design found false (each routed around minimally)

- **F1 — Exit codes.** The design's CL2 ("exit 2"), CL7 ("`n` → exit 2"), CL12 ("bad grammar →
  exit 2; unknown id → exit 1") and §2.6a ("refuses with exit 2") assume 2 = usage. At HEAD
  `exitUsage = 1` and `exitFatal = 2` (`main.go:71-75`, C5), the decline returns `exitUsage`
  (`session_new.go:438-440`), and CL17 forbids changing it. Routed by D4. CL12 tells the grammar
  refusal from the unknown-id refusal by message, not by code.
- **F2 — The DDL does not enforce the id grammar.** `request_id GLOB 'mr-[0-9a-f]*'` matches one
  hex character and then anything. Scratch `sqlite3` with the design's exact `CHECK`: `INSERT …
  VALUES ('mr-0zzzzzzzzzzzzzzz')` succeeds and reads back. The DDL stays verbatim (it is the
  design's canonical text); `PutMintRequest` enforces `MintRequestIDPattern` (D5), and ST8 uses the
  weak `CHECK` as its instrument control.
- **F3 — "`MintMatches` before the write" vs "`created_at` ≤ the in-transaction reading"**
  (sprint 1 F3/R5). The reading exists only inside `MintSessionForRequest` after `BeginTx`. Routed
  by D2.
- **F4 — CL15's "equal on every non-func field" does not hold for ordinary inputs.**
  `prepareMint` keeps the flag values as typed (`session_new.go:245`), while the request stores
  canonical or absolute paths (`WorldDB`, `CredentialsDB` canonical; `WorkspaceRoot`, `Repo`,
  `ProfilesDir` absolute). A defaulted `--credentials-db` is `""` on `episode start`'s side and a
  path on `approve`'s; `t.TempDir()` on darwin is under the `/var` → `/private/var` symlink. And
  `verb` differs (`"episode start"` vs `"approve"`). Routed: CL15 feeds canonical, absolute inputs
  (an `EvalSymlinks`'d base and an explicit canonical `--credentials-db`) and zeroes `verb` (an
  "excluded" field); BM1 pins the normalisation of relative inputs separately.
- **F5 — The field map is two entries short.** §2.1a maps the 29 fields of `mintOptions` +
  `mintRequest` at the design base. This sprint adds `consume` and `answerBy` to `mintRequest`, and
  CL14 requires every field to be mapped, so the map gains two `excluded:` entries (M6a).
- **F6 — "`approve` is a third caller."** It is the **fifth** (C6). The `:24-26` comment lists
  all five.
- **F7 — "Add `approve` to `everyVerb`" is not a one-line edit.** `everyVerb` also drives
  `TestHelpIsNeverBlockedByTheRegistryKeyGuard` (C15), so `approve`'s help must come out of
  `helpOnly` byte-identically. Routed by D10.
- **F8 — Not false, but unstated: the bump moves three existing gates.** `session_new.go:372`,
  `episode_preset.go:235` and `credentials.go:412` compare with the **current** version (C8). At
  version 4 a v3 store refuses bound mints, preset starts and resumes until `serve` of this build
  restarts, exactly as v2 stores did at the v3 bump. The design's "no extra operator step" holds
  for `episode start` (it already requires a same-build daemon, which upgraded the store at its
  start). The comparisons are **kept** (that precedent); this is risk R2, not a code change.

---

## 2. Shared definitions (the executor copies these exactly)

### 2.1 `host/store/credentials_schema.sql` (append, after `episode_starts`)

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

### 2.2 `host/store/credentials.go` edits (M4a)

- `const credentialsSchemaVersion = 4`; its doc comment gains: "version 4 adds mint requests (row
  201)".
- Add `"slices"` to the imports, and:

```go
// credentialsTables is the exact, sorted application table set of a credentials store at version.
func credentialsTables(version int) []string {
	switch {
	case version >= 4:
		return []string{"episode_starts", mintRequestsTable, credentialsTable}
	case version == 3:
		return []string{"episode_starts", credentialsTable}
	default:
		return []string{credentialsTable}
	}
}
```

- `enforceSchema`'s second case (`:304`) becomes
  `case known && slices.Equal(tables, credentialsTables(version)):` (body unchanged except the call
  below). The `known && len(tables) == 0` case and every other case are unchanged.
- Rename `upgradeCredentialsV3` → `upgradeCredentials` (doc: "adds what each version after
  `version` added, in one transaction"), and replace its statement list (`:332-335`) with:

```go
	var statements []string
	if version < 2 {
		statements = append(statements, credentialsUpgradeV2...)
	}
	if version < 3 {
		statements = append(statements, episodeStartsSQL)
	}
	if version < 4 {
		statements = append(statements, mintRequestsSQL)
	}
```

- D3: extract `insertSessionRow` (M4c, §2.4).
- `host/daemon/credentials.go:76`: the log line becomes
  `"ailang-worldd: upgraded the credentials store %s to schema %d (every existing credential unchanged)\n"`.

### 2.3 `host/store/mint_requests.go` (new; M4a holds only the first block, M4b the rest up to `PendingMintRequests`, M4c `MintSessionForRequest`)

```go
package store

// Row 201 A0 (design_docs/planned/w-world-approvals-by-key-a0a1.md §2.4): stored, single-use
// mint requests. `episode start --request` puts one; `approve` consumes it in the same
// transaction that inserts its session row.

// ---- M4a ----
const mintRequestsTable = "mint_requests"

// mintRequestsSQL is the version 3 -> 4 upgrade: credentials_schema.sql's mint_requests table,
// column for column (TestCredentialsUpgradeToV4 compares the two).
const mintRequestsSQL = `CREATE TABLE mint_requests (
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
)`

// ---- M4b ----
// mintRequestsSchemaVersion is the first credentials schema with mint_requests.
const mintRequestsSchemaVersion = 4

// MaxPendingMintRequests caps the unconsumed, unexpired requests one store holds (design §4, R4).
const MaxPendingMintRequests = 16

// mintRequestPruneAfter: an unconsumed request is deleted once it has been expired this long.
const mintRequestPruneAfter = 86400

// MintRequestIDPattern is a request id: "mr-" and 16 lowercase hex (8 random bytes). The table's
// CHECK is weaker (plan §1.4 F2), so PutMintRequest enforces this.
var MintRequestIDPattern = regexp.MustCompile(`^mr-[0-9a-f]{16}$`)

var (
	ErrTooManyPendingRequests = errors.New("store: 16 mint requests are already pending; approve them or let them expire (ailang-worldd approve --list)")
	ErrRequestUnknown         = errors.New("store: no such mint request")
	ErrRequestConsumed        = errors.New("store: the mint request was already used")
	ErrRequestExpired         = errors.New("store: the mint request expired before it was consumed")
	ErrRequestClockBehind     = errors.New("store: the clock is behind the mint request's creation time")
	ErrRequestDigestChanged   = errors.New("store: the mint request's digest is not the one approved")
	ErrBadMintRequest         = errors.New("store: invalid mint request row")
	ErrIncompleteConsume      = errors.New("store: a mint request consume needs an id, a digest, approved_by, a clock and a check")
)

// OlderCredentialsStoreError refuses a mint-request operation on a store below version 4.
type OlderCredentialsStoreError struct {
	Path    string
	Version int
}

func (e *OlderCredentialsStoreError) Error() string {
	return fmt.Sprintf("store: credentials store %q is schema %d; mint requests need schema %d: "+
		"restart `ailang-worldd serve` with this build once (it upgrades the store at startup), then run this again",
		e.Path, e.Version, mintRequestsSchemaVersion)
}

// MintRequestRow is one mint_requests row.
type MintRequestRow struct {
	RequestID    string
	Request      []byte // the canonical request bytes (approval.Canonical)
	Digest       string // 64 lowercase hex: SHA-256 of Request when it was made
	Challenge    string // 64 lowercase hex
	EpisodeID    string
	CreatedAt    int64
	ExpiresAt    int64
	Consumed     bool   // consumed_at IS NOT NULL
	ConsumedAt   int64  // 0 unless Consumed
	ApprovedBy   string // "" unless Consumed
	CredentialID string // "" unless Consumed
}

func (c *Credentials) PutMintRequest(ctx context.Context, row MintRequestRow, now int64) error
func (c *Credentials) MintRequest(ctx context.Context, id string) (MintRequestRow, bool, error)
func (c *Credentials) PendingMintRequests(ctx context.Context, now int64, limit int) ([]MintRequestRow, error)

func isLowerHex64(s string) bool // len 64, every byte in [0-9a-f]
const mintRequestColumns = `request_id, request, digest, challenge, episode_id, created_at, expires_at, consumed_at, approved_by, credential_id`
func scanMintRequest(scan func(...any) error) (MintRequestRow, error) // sql.NullInt64 / sql.NullString for the last three

// ---- M4c ----
// RequestConsume is what MintSessionForRequest consumes. Clock is read exactly once, inside the
// transaction, after BeginTx (design §2.4 [r3]); that reading is consumed_at and the expiry
// bound. Check runs with that reading before the guarded UPDATE; an error rolls everything back
// (plan §1.3 D2: authority.MintForRequest puts approval.MintMatches here).
type RequestConsume struct {
	ID, Digest, ApprovedBy string
	Clock                  func() int64
	Check                  func(consumedAt int64) error
}

func (c *Credentials) MintSessionForRequest(ctx context.Context, session SessionRow, rc RequestConsume) error
func insertSessionRow(ctx context.Context, tx *sql.Tx, row SessionRow) error // D3
func classifyConsume(ctx context.Context, tx *sql.Tx, rc RequestConsume, at int64, changed int64) error
```

**`PutMintRequest`**, in this order, each its own branch:
1. `requireDeadline(ctx)`.
2. `c.version < mintRequestsSchemaVersion` → `&OlderCredentialsStoreError{Path: c.display, Version: c.version}`.
3. Row checks, each `fmt.Errorf("%w: <what>", ErrBadMintRequest)`:
   `!MintRequestIDPattern.MatchString(row.RequestID)`; `len(row.Request) == 0`;
   `!isLowerHex64(row.Digest)`; `!isLowerHex64(row.Challenge)`; `row.EpisodeID == ""`;
   `row.ExpiresAt <= row.CreatedAt`; `row.Consumed` ("a new request is unconsumed").
4. `tx := c.db.BeginTx(ctx, nil)`; `defer tx.Rollback()`.
5. `DELETE FROM mint_requests WHERE consumed_at IS NULL AND expires_at <= ?` with
   `now - mintRequestPruneAfter`.
6. `SELECT count(*) FROM mint_requests WHERE consumed_at IS NULL AND expires_at > ?` with `now`;
   `n >= MaxPendingMintRequests` → `ErrTooManyPendingRequests`.
7. `INSERT INTO mint_requests (request_id, request, digest, challenge, episode_id, created_at,
   expires_at) VALUES (?, ?, ?, ?, ?, ?, ?)`; then `tx.Commit()`.

**`MintRequest`**: `requireDeadline`; version check (as step 2); `SELECT ` + `mintRequestColumns` +
` FROM mint_requests WHERE request_id = ?`; `sql.ErrNoRows` → `(MintRequestRow{}, false, nil)`.

**`PendingMintRequests`**: `requireDeadline`; version check; `limit <= 0 || limit >
MaxPendingMintRequests` → `MaxPendingMintRequests`; `SELECT ` + columns + ` FROM mint_requests
WHERE consumed_at IS NULL AND expires_at > ? ORDER BY created_at DESC, request_id DESC LIMIT ?`.

**`MintSessionForRequest`**, in this order:
1. `requireDeadline(ctx)`.
2. `rc.ID == "" || rc.Digest == "" || rc.ApprovedBy == "" || rc.Clock == nil || rc.Check == nil`
   → `ErrIncompleteConsume`.
3. `session.CredentialID == "" || session.EpisodeID == ""` → `fmt.Errorf("store: mint session:
   empty credential_id or episode_id")`; a half binding → the same message `mintSession` uses.
4. `c.version < mintRequestsSchemaVersion` → `*OlderCredentialsStoreError`.
5. `tx, err := c.db.BeginTx(ctx, nil)` (error → `fmt.Errorf("store: consume mint request: begin:
   %w", err)`); `defer tx.Rollback()`.
6. `at := rc.Clock()` — the **only** call.
7. `if err := rc.Check(at); err != nil { return err }`.
8. `res, err := tx.ExecContext(ctx, `UPDATE mint_requests SET consumed_at = ?, approved_by = ?,
   credential_id = ? WHERE request_id = ? AND digest = ? AND consumed_at IS NULL AND created_at <=
   ? AND expires_at > ?`, at, rc.ApprovedBy, session.CredentialID, rc.ID, rc.Digest, at, at)`;
   `n, err := res.RowsAffected()`.
9. `n != 1` → `return classifyConsume(ctx, tx, rc, at, n)`.
10. `insertSessionRow(ctx, tx, session)` → `fmt.Errorf("store: mint session: %w", err)`.
11. `tx.Commit()` → `fmt.Errorf("store: consume mint request: commit: %w", err)`.

**`classifyConsume`**: `SELECT digest, created_at, expires_at, consumed_at IS NOT NULL FROM
mint_requests WHERE request_id = ?` in `tx`. In this order: no row → `ErrRequestUnknown`;
consumed → `ErrRequestConsumed`; `at < created_at` → `ErrRequestClockBehind`; `at >= expires_at` →
`ErrRequestExpired`; `digest != rc.Digest` → `ErrRequestDigestChanged`; otherwise
`fmt.Errorf("store: mint request %s: the guarded update changed %d rows", rc.ID, changed)`. Every
sentinel is wrapped as `fmt.Errorf("%w: %s", ErrX, rc.ID)`.

### 2.4 `host/authority/mint.go` (M6a)

```go
import "github.com/sunholo-data/ailang-world/host/approval"

// RequestMint is what a mint from a stored request needs beyond mintRecorded's arguments
// (row 201 §2.6): the decoded request, its row's created_at, and the consume.
type RequestMint struct {
	Request          approval.MintRequest
	RequestCreatedAt int64
	Consume          store.RequestConsume // ID, Digest, ApprovedBy, Clock; MintForRequest sets Check
}

// MintForRequest mints exactly the session the request describes and consumes the request in the
// same transaction. approval.MintMatches checks all seven columns of the row it built inside that
// transaction, with the in-transaction clock reading as created_at's upper bound, before the row
// is written.
func MintForRequest(ctx context.Context, st *store.Credentials, episodeID string, binding ProfileBinding, grants []broker.Capability, ttl, now int64, out io.Writer, rm RequestMint) (rawToken, credentialID string, row store.SessionRow, err error) {
	return mintRecorded(ctx, st, episodeID, binding, grants, ttl, now, out, nil, nil, &rm)
}
```

`mintRecorded` gains a last parameter `rm *RequestMint`; `MintStarted` and `MintResumed` pass
`nil`. Its store call (`:96-100`) becomes:

```go
	switch {
	case rm != nil:
		consume := rm.Consume
		consume.Check = func(consumedAt int64) error {
			return approval.MintMatches(rm.Request, approval.MintFacts{
				CredentialID: row.CredentialID, EpisodeID: row.EpisodeID, GrantsJSON: row.GrantsJSON,
				ExpiresAt: row.ExpiresAt, CreatedAt: row.CreatedAt, Project: row.Project,
				ProfileDigest: row.ProfileDigest, TokenHash: hashHexToken(rawToken),
				RequestCreatedAt: rm.RequestCreatedAt, ConsumedAt: consumedAt,
			})
		}
		err = st.MintSessionForRequest(ctx, row, consume)
	case resume != nil:
		err = st.MintSessionWithResume(ctx, row, *resume)
	default:
		err = st.MintSessionWithStart(ctx, row, start)
	}
```

### 2.5 `cmd/ailang-worldd/session.go` (M5b: `random`; M6a: `after`, `readTermLineWithin`)

```go
	// random fills b with cryptographically random bytes (crypto/rand in production). A request
	// path that finds it nil refuses rather than fall back (row 201 §2.2).
	random func(b []byte) error
	// after is time.After in production: the timer that bounds approve's y/N (row 201 §2.6a).
	after func(time.Duration) <-chan time.Time
```

`realSessionEnv` adds `random: func(b []byte) error { _, err := rand.Read(b); return err }`
(`crypto/rand`) and `after: time.After`.

```go
// readTermLineWithin is readTermLine bounded by d: the read runs in a goroutine that sends to a
// buffered channel of size 1, and the first of the line and after(d) wins; ok is false on timeout.
// It is only safe for one-shot CLI entrypoints that exit the process on return: a timed-out read
// leaves its goroutine blocked on the terminal until the process exits (row 201 §2.6a [r3]).
func readTermLineWithin(r io.Reader, d time.Duration, after func(time.Duration) <-chan time.Time) (string, bool) {
	ch := make(chan string, 1)
	go func() { ch <- readTermLine(r) }()
	select {
	case line := <-ch:
		return line, true
	case <-after(d):
		return "", false
	}
}
```

### 2.6 `cmd/ailang-worldd/session_new.go` (M6a)

`mintRequest` gains, after `presetPrompt`:

```go
	// consume is set only by approve (row 201): the stored request this mint consumes, in the
	// mint's own transaction.
	consume *authority.RequestMint
	// answerBy bounds the y/N to this Unix time; 0 (every caller but approve) is unbounded.
	answerBy int64
```

In `provisionEpisode`, **immediately before** `if r.presetPrompt != "" {` (`:430`):

```go
	requestID := ""
	if r.consume != nil {
		requestID = r.consume.Consume.ID
	}
	if r.answerBy != 0 && r.answerBy-now <= 0 {
		fmt.Fprintf(stderr, "%srequest %s expired before the prompt; make a new request with episode start --request\n", prefix, requestID)
		return exitUsage
	}
```

and the read at `:438` becomes:

```go
	var answer string
	if r.answerBy != 0 {
		remaining := time.Duration(r.answerBy-now) * time.Second
		line, ok := readTermLineWithin(term, remaining, env.after)
		if !ok {
			fmt.Fprintf(stderr, "%sno answer within the request's remaining lifetime (%s); nothing was minted and request %s is still pending until %s\n",
				prefix, remaining, requestID, time.Unix(r.answerBy, 0).UTC().Format(time.RFC3339))
			return exitUsage
		}
		answer = line
	} else {
		answer = readTermLine(term)
	}
	if !isYes(answer) {
```

(the decline body is unchanged). The mint branch at `:478` gains a **first** arm:

```go
	if r.consume != nil {
		rawToken, credentialID, row, err = authority.MintForRequest(mintCtx, st, r.episode, binding, grants, r.ttl, now, out, *r.consume)
	} else if r.resume != nil {
```

The `:24-26` comment becomes: "Steps 2-7 are provisionEpisode, the ONE fence-and-mint path:
`episode start` (episode.go, row 189 M3), its launch-preset and resume paths, and `approve`
(approve.go, row 201, which consumes a stored mint request in the mint's own transaction) run it
too, after their own validation, so the verbs cannot drift apart and no other code creates a
credential."

### 2.7 `cmd/ailang-worldd/mint_request.go` (new; M5a, plus `createMintRequest` in M5b and two map entries in M6a)

```go
package main

// Row 201 A0 (design_docs/planned/w-world-approvals-by-key-a0a1.md §2.1a, §2.5): the mint
// request `episode start --request` stores and `approve` mints from.

// approveReadBudget bounds the store reads of `episode start --request` and `approve` (3 s).
const approveReadBudget = 3 * time.Second

// buildMintRequest is design §2.5 step 3: the MintRequest of a validated mintRequest. Every value
// approve acts on is in it; paths are canonical (stores) or absolute (the rest).
func buildMintRequest(r *mintRequest, build daemon.BuildIdentity, addr, agent, model, task, mcpFile string) (approval.MintRequest, error)

// fromRequest is the §2.1a map read backwards: the mintOptions approve validates, built from the
// decoded request only. No flag or default reaches it.
func fromRequest(mr approval.MintRequest) mintOptions

func fromCapabilities(cs []broker.Capability) []approval.Grant // {Effect, Scope, Budget}
func toCapabilities(gs []approval.Grant) []broker.Capability   // {Effect, Scope, Budget}, ExpiresAt 0

// mcpAfterToken is the afterToken that writes the episode's MCP config; nil iff mcpFile is "".
func mcpAfterToken(mcpFile, addr string) func(token string) error

// renderRequest is the request block --request prints and approve shows above its y/N.
func renderRequest(id string, mr approval.MintRequest, digestHex string, expiresAt, now int64) string

var mintRequestFieldMap = map[string]string{ /* exactly the table below */ }

// requestOnlyLeaves are the MintRequest leaves no mintOptions/mintRequest field maps to: approve
// uses them directly.
var requestOnlyLeaves = []string{"Schema", "MCPConfig", "Agent", "Model", "Task",
	"Daemon.Addr", "Daemon.VCSRevision", "Daemon.VCSModified", "Daemon.ExecutableSHA256", "Daemon.GoVersion"}
```

**`buildMintRequest`**, in order: any `r.grants[i].ExpiresAt != 0` → `fmt.Errorf("grant %d carries
an expiry; a request's grants expire at the mint", i)`; `worldDB, err := store.CanonicalPath(r.db)`;
`wsRoot, err := filepath.Abs(r.wsRoot)`; `repo := ""` and, if `r.repo != ""`, `repo, err =
filepath.Abs(r.repo)`; `profilesDir, err := filepath.Abs(r.profilesDir)`; `prof :=
approval.ProfileRef{}` and, if `r.prof != nil`, `{Project: r.prof.Project, Digest: r.prof.Digest}`.
Then:

```go
	mr := approval.MintRequest{Schema: approval.RequestSchema, Episode: r.episode, WorldDB: worldDB,
		WorkspaceRoot: wsRoot, Worktree: r.dir, Reuse: r.reuse, Degraded: r.degraded, Repo: repo, Git: r.git,
		Branch: r.branch, Grants: fromCapabilities(r.grants), TTLSeconds: r.ttl, Profile: prof,
		ProfilesDir: profilesDir, CredentialsDB: r.credPath, SessionFile: r.outPath, MCPConfig: mcpFile,
		Agent: agent, Model: model, Task: task,
		Daemon: approval.DaemonRef{Addr: addr, VCSRevision: build.VCSRevision, VCSModified: build.VCSModified,
			ExecutableSHA256: build.ExecutableSHA256, GoVersion: build.GoVersion}}
	return mr, mr.Validate()
```

**`fromRequest`**:

```go
	return mintOptions{verb: "approve", episode: mr.Episode, db: mr.WorldDB, credsDB: mr.CredentialsDB,
		wsRoot: mr.WorkspaceRoot, repo: mr.Repo, branch: mr.Branch, allowNonGit: false,
		grants: toCapabilities(mr.Grants), ttl: mr.TTLSeconds, outPath: mr.SessionFile, outDirCreate: true,
		profile: mr.Profile.Project, profilesDir: mr.ProfilesDir}
```

**`mcpAfterToken`**: `if mcpFile == "" { return nil }`; otherwise `return func(token string) error {
return writeMCPConfig(mcpFile, addr, token) }`.

**`renderRequest`** (exact lines; `grantList` is `session_new.go:669`):

```go
	var b strings.Builder
	fmt.Fprintf(&b, "mint request %s\n", id)
	fmt.Fprintf(&b, "  digest    %s\n", digestHex)
	fmt.Fprintf(&b, "  expires   %s (in %s)\n", time.Unix(expiresAt, 0).UTC().Format(time.RFC3339), time.Duration(expiresAt-now)*time.Second)
	fmt.Fprintf(&b, "  episode   %s\n", mr.Episode)
	where := mr.Worktree + " (new, from " + mr.Repo + " on branch " + mr.Branch + ")"
	if mr.Reuse {
		where = mr.Worktree + " (existing worktree)"
	}
	fmt.Fprintf(&b, "  worktree  %s\n", where)
	fmt.Fprintf(&b, "  grants    %s\n", grantList(toCapabilities(mr.Grants)))
	fmt.Fprintf(&b, "  session   %ds; token file %s\n", mr.TTLSeconds, mr.SessionFile)
	profile := "unbound"
	if mr.Profile.Project != "" {
		profile = mr.Profile.Project + " " + mr.Profile.Digest + " (" + mr.ProfilesDir + ")"
	}
	fmt.Fprintf(&b, "  profile   %s\n", profile)
	mcp, task := "none", "none"
	if mr.MCPConfig != "" {
		mcp = mr.MCPConfig
	}
	if mr.Task != "" {
		task = mr.Task
	}
	fmt.Fprintf(&b, "  agent     %s %s; MCP config %s\n", mr.Agent, mr.Model, mcp)
	fmt.Fprintf(&b, "  task      %s\n", task)
	fmt.Fprintf(&b, "  daemon    %s build %s\n", mr.Daemon.Addr, daemon.BuildIdentity{VCSRevision: mr.Daemon.VCSRevision,
		VCSModified: mr.Daemon.VCSModified, ExecutableSHA256: mr.Daemon.ExecutableSHA256, GoVersion: mr.Daemon.GoVersion})
	fmt.Fprintf(&b, "  stores    %s, %s\n", mr.WorldDB, mr.CredentialsDB)
	return b.String()
```

**`mintRequestFieldMap`** (M5a: 33 entries; M6a adds the last two). Keys are Go field names of
`mintOptions` and `mintRequest` (the embedded struct flattened) plus `grants[].<Capability
field>`. Values are `MintRequest.<Go path>` (`[]` marks a slice element), or start with
`excluded: `, `fixed: ` or `nil: `.

| Key | Value |
|---|---|
| `verb` | `excluded: display only, the message prefix; approve sets "approve"` |
| `episode` | `MintRequest.Episode` |
| `db` | `MintRequest.WorldDB` |
| `credsDB` | `MintRequest.CredentialsDB` |
| `wsRoot` | `MintRequest.WorkspaceRoot` |
| `repo` | `MintRequest.Repo` |
| `branch` | `MintRequest.Branch` |
| `allowNonGit` | `fixed: false (episode start never sets it)` |
| `grants` | `MintRequest.Grants` |
| `grants[].Effect` | `MintRequest.Grants[].Effect` |
| `grants[].Scope` | `MintRequest.Grants[].Scope` |
| `grants[].Budget` | `MintRequest.Grants[].Budget` |
| `grants[].ExpiresAt` | `fixed: 0, set at the mint to now + ttl` |
| `ttl` | `MintRequest.TTLSeconds` |
| `outPath` | `MintRequest.SessionFile` |
| `outDirCreate` | `fixed: true (episode start's value)` |
| `profile` | `MintRequest.Profile.Project` |
| `profilesDir` | `MintRequest.ProfilesDir` |
| `credPath` | `MintRequest.CredentialsDB` |
| `dir` | `MintRequest.Worktree` |
| `reuse` | `MintRequest.Reuse` |
| `degraded` | `MintRequest.Degraded` |
| `git` | `MintRequest.Git` |
| `prof` | `MintRequest.Profile.Digest` |
| `profFile` | `excluded: derived as ProfilesDir/<Profile.Project>.json` |
| `afterToken` | `excluded: set by approve from MintRequest.MCPConfig and MintRequest.Daemon.Addr (nil iff no MCP config)` |
| `resume` | `nil: set only by episode_resume.go; --request refuses --resume` |
| `beforeToken` | `nil: set only by episode_resume.go` |
| `resumeRollback` | `nil: set only by episode_resume.go` |
| `start` | `nil: set only by episode_preset.go; --request refuses launch presets` |
| `startPoint` | `fixed: "" (set only by episode_preset.go)` |
| `existingBranch` | `fixed: false (set only by episode_preset.go)` |
| `presetPrompt` | `excluded: display only; approve sets renderRequest(request)` |
| `consume` **(M6a)** | `excluded: set by approve from the stored row (id, digest, approved_by tty, env.now)` |
| `answerBy` **(M6a)** | `excluded: set by approve to the stored row's expires_at` |

**`createMintRequest`** (M5b):

```go
// createMintRequest is `episode start --request` steps 3-5 (design §2.5): it stores the request
// and prints it. It opens no terminal and mints nothing.
func createMintRequest(r *mintRequest, build daemon.BuildIdentity, addr, agent, model, task, mcpFile string, stdout, stderr io.Writer, env sessionEnv) int
```

In order: `env.random == nil` → stderr `ailang-worldd episode start: internal error: no randomness
source (sessionEnv.random)`, `exitFatal`. `mr, err := buildMintRequest(…)` → `exitUsage` with the
error. `var idb [8]byte; var ch [32]byte`; `env.random(idb[:])` then `env.random(ch[:])` (that
order; an error → `exitFatal`). `id := "mr-" + hex.EncodeToString(idb[:])`. `canon, err :=
approval.Canonical(mr)`; `sum := sha256.Sum256(canon)`. `now := env.now()`. `row :=
store.MintRequestRow{RequestID: id, Request: canon, Digest: hex.EncodeToString(sum[:]), Challenge:
hex.EncodeToString(ch[:]), EpisodeID: mr.Episode, CreatedAt: now, ExpiresAt: now +
int64(approval.RequestLifetime/time.Second)}`. Then the **one** root:
`ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)`; `defer cancel()`;
`st, err := openCredentialsForSession(ctx, r.db, r.credPath)` (error → `exitFatal`); `defer
st.Close()`; `st.PutMintRequest(ctx, row, now)`: `*store.OlderCredentialsStoreError` or
`ErrTooManyPendingRequests` → stderr `ailang-worldd episode start: <err>`, `exitUsage`; any other
error → `exitFatal`. On success print to stdout `renderRequest(id, mr, row.Digest, row.ExpiresAt,
now)`, then `next: ailang-worldd approve <id> --db <mr.WorldDB> --credentials-db
<mr.CredentialsDB>` and `nothing is minted until a person approves it at a terminal; it expires in
10m0s`, each on its own line. Return `exitOK`.

### 2.8 `cmd/ailang-worldd/episode.go` (M5b)

- Flag, after `headless`: `request := fs.Bool("request", false, "store a mint request instead of minting (no terminal; approve it with ailang-worldd approve <id>)")`.
- Insert, between `fs.Visit(… resumeSet …)` (`:233`) and `if resumeSet {` (`:234`):

```go
	if *request {
		why := "a request carries no launch callbacks or start records (row 201 §2.5)"
		switch {
		case resumeSet:
			return fail("--request cannot be combined with --resume: %s; fix: drop --request", why)
		case *launch:
			return fail("--request cannot be combined with --launch: %s; fix: drop --launch, then run the launch line approve prints", why)
		case *headless:
			return fail("--request cannot be combined with --headless: %s; fix: drop --headless", why)
		case *preset != "" && *preset != "se-tools":
			return fail("--request cannot be combined with a launch preset (--preset %s): %s; fix: use --grant or --preset se-tools", *preset, why)
		}
	}
```

- After `checkDaemon` succeeds (`:335-337`) and **before** `if mcpFile != "" {` (`:338`):

```go
	if *request {
		return createMintRequest(req, health.Build, hostPort, *agent, *model, taskPath, mcpFile, stdout, stderr, env)
	}
```

- `episodeHelp`: append `[--request]` to the `usage: ailang-worldd episode start <episode> …`
  block's last flag line, and add under the flag list:
  ```
    --request             store a mint request instead of minting: no terminal, no worktree,
                          no token. It prints the request id; 'ailang-worldd approve <id>' mints
                          it after a y/N. It expires in 10 minutes. Not with --launch, --headless,
                          --resume or a launch preset.
  ```
- `main.go` `usage` (`:119-123`): append ` [--request]` to the `episode start <episode> --repo`
  line's last continuation line.

### 2.9 `cmd/ailang-worldd/approve.go` (new; M6b)

```go
package main

// Row 201 A0 (design §2.6): `ailang-worldd approve <request-id>` mints the session a stored mint
// request describes after the typed y, through provisionEpisode; `approve --list` shows the
// pending requests. The request is consumed in the mint's own transaction.

const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
       ailang-worldd approve --list [--db <path>] [--credentials-db <path>]

ATTENDED. Mints the session a stored mint request describes (row 201), after a
y/N typed at /dev/tty. 'ailang-worldd episode start <episode> ... --request'
makes the request without a terminal; nothing is minted until approve.

approve <request-id> (mr- and 16 hex):
  1. reads the request and refuses, before any terminal opens, one that is
     unknown, already used, expired, or whose stored bytes do not match their
     digest;
  2. re-derives the request from the world as it is now (the exec profile, the
     serving daemon's build, the worktree, the session and MCP files) and
     refuses, naming the first field, if anything differs;
  3. runs episode start's own path: the /dev/tty fence, a y/N showing the whole
     request, then the worktree and the mint. The y/N waits at most until the
     request expires. The request is consumed in the mint's own transaction, so
     it mints at most once;
  4. prints the launch line, as episode start does. The token is never printed.
--list prints the pending requests, read-only. A request expires 10 minutes
after it is made; at most 16 are pending per credentials store.

  --db <path>               the world store the request was made for
                            (default ~/.ailang/world/world.db)
  --credentials-db <path>   default credentials.db beside --db
The global --addr is refused: the daemon's address is part of the request.
`

func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int
func approveAdmission(row store.MintRequestRow, now int64) (approval.MintRequest, string) // ("", refusal) or (decoded, "")
```

`runApprove`, in order (prefix `ailang-worldd approve: `; every refusal before step 11 is
`exitUsage` unless stated):
1. `flag.NewFlagSet("ailang-worldd approve", flag.ContinueOnError)`, output `stderr`, `Usage =
   func() {}`; flags `--db` (default `defaultWorldPaths().db`), `--credentials-db` (default `""`,
   `credentialsDBUsage`), `--list` (bool). `parseInterleaved`. `flag.ErrHelp` → print
   `approveHelp` to **stdout**, `exitOK`. Other parse error → `approveHelp` to stderr, `exitUsage`.
2. `--list` with a positional, or no `--list` and `len(positional) != 1` → `approveHelp` to
   stderr, `exitUsage`.
3. Not `--list`: `!store.MintRequestIDPattern.MatchString(id)` → `%q is not a mint request id (mr-
   and 16 lowercase hex)`. Then `env.after == nil` → `internal error: no prompt timer
   (sessionEnv.after)`, `exitFatal` (D11).
4. `credPath, err := credentialsFor(*db, *credsDB)`.
5. The **one** root: `ctx, cancel := context.WithTimeout(context.Background(),
   approveReadBudget)`. `ro, err := store.OpenCredentialsReadOnly(credPath)` (error → refusal
   naming `credPath` and the error). `--list`: `rows, err := ro.PendingMintRequests(ctx,
   env.now(), 0)`; `cancel()`; `ro.Close()`; print `no pending mint requests` when empty, else one
   line each: `fmt.Fprintf(stdout, "%s  episode %s  digest %s  expires %s\n", r.RequestID,
   r.EpisodeID, r.Digest[:16], time.Unix(r.ExpiresAt, 0).UTC().Format(time.RFC3339))`; `exitOK`.
   Otherwise `row, ok, err := ro.MintRequest(ctx, id)`; `cancel()`; `ro.Close()`. Store error →
   refusal with the error. `!ok` → `no mint request %s in %s; fix: ailang-worldd approve --list
   shows the pending ones, or make one with episode start --request`.
6. `now := env.now()` — call **[1]** (D8). `mr, refusal := approveAdmission(row, now)`; a refusal
   → print it, `exitUsage`. `approveAdmission`: `approval.Decode(row.Request)` error → `request %s
   is not a canonical mint request (%v); fix: make a new request with episode start --request`;
   then `stored, err := hex.DecodeString(row.Digest)` (bad hex or `len(stored) != 32` → the digest-mismatch text, never a panic);
   `approval.Admit(approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt:
   row.ExpiresAt, Consumed: row.Consumed, StoredDigest: [32]byte(stored), BytesDigest:
   sha256.Sum256(row.Request)})`:
   - `RefuseConsumed` → `request %s was already used at %s by %s (credential %s); fix: make a new request with episode start --request`
   - `RefuseClockBehind` → `the clock is behind request %s's creation time (%s); refusing`
   - `RefuseExpired` → `request %s expired at %s; fix: make a new request with episode start --request`
   - `RefuseDigestMismatch` → `request %s's stored bytes do not match its digest; refusing; fix: make a new request with episode start --request`
   - any other non-`Admitted` verdict → `request %s: unexpected admission verdict %s`.
   (Times print as `time.Unix(…, 0).UTC().Format(time.RFC3339)`.)
7. Same store (design §2.1a `db` row, [r3]): `w, err := store.CanonicalPath(*db)`; `err != nil
   || w != mr.WorldDB || credPath != mr.CredentialsDB` → `--db/--credentials-db name %s and %s, but
   request %s was made for %s and %s; fix: pass --db %s --credentials-db %s`.
8. `req, err := prepareMint(fromRequest(mr))` → `the world changed since request %s was made:
   %v; it stays pending`. Then `mr.MCPConfig != ""` and `os.Lstat(mr.MCPConfig)` succeeds →
   `the world changed since request %s was made: %s already exists; it stays pending` (D12).
9. `health, err := fetchDaemonHealth(mr.Daemon.Addr)` → `refusing: %v` (the
   `daemonUnreachableError` text included). `checkDaemon(health, req, mr.Daemon.Addr,
   episodeSelfBuild())` → `refusing: %v`.
10. `fresh, err := buildMintRequest(req, health.Build, mr.Daemon.Addr, mr.Agent, mr.Model,
    mr.Task, mr.MCPConfig)` (error → the "world changed" text); `fb, _ := approval.Canonical(fresh)`;
    `!bytes.Equal(fb, row.Request)` → `the world changed since request %s was made: %s differs;
    it stays pending` with `approval.Diff(mr, fresh)`.
11. `req.presetPrompt = renderRequest(row.RequestID, mr, row.Digest, row.ExpiresAt, now)`;
    `req.afterToken = mcpAfterToken(mr.MCPConfig, mr.Daemon.Addr)`; `req.answerBy =
    row.ExpiresAt`; `req.consume = &authority.RequestMint{Request: mr, RequestCreatedAt:
    row.CreatedAt, Consume: store.RequestConsume{ID: row.RequestID, Digest: row.Digest,
    ApprovedBy: "tty", Clock: env.now}}`.
12. `if code := provisionEpisode(req, stdout, stderr, env); code != exitOK { return code }`.
13. `printLaunch(stdout, launchInfo{episode: mr.Episode, agent: mr.Agent, model: mr.Model, addr:
    mr.Daemon.Addr, worktree: req.dir, branch: mr.Branch, reused: req.reuse, sessionFile:
    mr.SessionFile, mcpFile: mr.MCPConfig, task: mr.Task, grants: toCapabilities(mr.Grants)})`;
    `exitOK`.

`main.go`: in `run`'s switch, beside `case "session"`:

```go
	case "approve":
		if addrGiven {
			fmt.Fprintln(stderr, "ailang-worldd: --addr is a client flag and is not valid for 'approve'; "+
				"the daemon's address is part of the mint request")
			return exitUsage
		}
		return runApprove(rest[1:], stdout, stderr, realSessionEnv())
```

`usage`: after the `episode start <episode> --repo` block, two lines:

```
  ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
  ailang-worldd approve --list [--db <path>] [--credentials-db <path>]
```

`ownHelpVerbs` gains `"approve": true`; `ownVerbHelp` gains `case "approve": return approveHelp,
true`.

### 2.10 Test helpers (exact; `cmd/ailang-worldd/approve_test.go` unless stated)

```go
// approveT0 is the fixed creation time every approve test uses (2027-01-15T08:00:00Z).
const approveT0 int64 = 1800000000

// seqRandom fills each call's buffer with the next bytes of 1, 2, 3, ... (mod 256), so the first
// request's id is mr-0102030405060708 and its challenge is 09..28.
func seqRandom() func([]byte) error

// requestEnv is episode start --request's env: openTerminal records into *opened and fails,
// now returns now, random is seqRandom().
func requestEnv(opened *bool, now int64) sessionEnv

// clockSeq returns before for the first n calls and after from then on (D8 call order).
func clockSeq(n int, before, after int64) func() int64

// approveEnv: openTerminal returns term (and records into *opened when non-nil), now is clock,
// after is time.After.
func approveEnv(term io.ReadWriteCloser, opened *bool, clock func() int64) sessionEnv

type approveRig struct {
	r        episodeRun
	ep, id   string
	credPath string
	profile  *broker.ExecProfile
	task     string
}

// newApproveRig writes profile "proj", a task file, stubs a matching daemon, and runs
// `episode start <ep> --profile proj --task <task> --model claude-opus-5-5 --request` at createdAt
// with requestEnv. It fails the test unless that exits 0 with no terminal opened, and returns the
// rig with the request id read back from PendingMintRequests.
func newApproveRig(t *testing.T, ep string, createdAt int64, extra ...string) approveRig // = newApproveRigIn(t, newEpisodeRun(t), …)
// newApproveRigIn is newApproveRig in an existing run, so several requests share one store
// (writeTestProfile rewrites the same bytes, so the digest is unchanged; stubDaemon is re-applied).
func newApproveRigIn(t *testing.T, r episodeRun, ep string, createdAt int64, extra ...string) approveRig
func (a approveRig) args(extra ...string) []string           // id, --db, db, --credentials-db, credPath, extra...
func (a approveRig) row(t *testing.T) store.MintRequestRow    // read-only MintRequest; t.Fatal if absent
func (a approveRig) sessionRows(t *testing.T) []store.SessionRow
func (a approveRig) worktree() string                          // <ws>/<ep>
```

`host/store` test helper (in `mint_requests_test.go`):

```go
// makeV3Credentials writes a credentials store at user_version 3: canonicalCredentialsDDL,
// episodeStartsSQL, the given rows.
func makeV3Credentials(t *testing.T, path string, rows []SessionRow)
// testMintRequestRow is a valid unconsumed row: id "mr-" + fmt.Sprintf("%016x", n), Request
// []byte(`{"n":<n>}`), Digest strings.Repeat("ab", 32), Challenge strings.Repeat("cd", 32),
// EpisodeID "ep", CreatedAt created, ExpiresAt created+600.
func testMintRequestRow(n int, created int64) MintRequestRow
// okConsume returns a RequestConsume for row with ApprovedBy "tty", a Clock returning at and
// counting calls into *clocks, and a Check returning nil and counting into *checks.
func okConsume(row MintRequestRow, at int64, clocks, checks *int) RequestConsume
```

---

## 3. Milestones

Every milestone ends with `go build ./...` **and** `go vet` of each package it touched, plus
`test -z "$(gofmt -l <touched dirs>)"`. No `t.Parallel()`. Wrap any loop or `$VAR` in `bash -c`
(the controller's shell is zsh). For any `go test -v` run of `cmd/ailang-worldd`, also run the
**skip guard**: `… -v 2>&1 | grep -E '^\s*--- SKIP' && exit 1 || true` (C17: `matchingHealth`
skips on a platform with no tool pin; on this rig and on CI it must not).

### M3b — close the sprint-1 judge's test gaps (`host/approval`, tests only) · ~0 prod + ~120 test LOC

**Files**: `host/approval/sshsig_test.go` (one transform helper set + SG8),
`host/approval/sshsig_crosscheck_test.go` (three arms). **No production or testdata change.**

**New transforms** (`sshsig_test.go`, beside the existing ones; each takes a real fixture):

```go
// withSigStringTrailingByte appends one 0x00 INSIDE the signature string, after the inner blob
// (parseSSHSIGBlob's sr.end(), sshsig.go:248).
func withSigStringTrailingByte(t *testing.T, a []byte) []byte  // p.signature = append(bytes.Clone(p.signature), 0)
// withMpintPairTrailingByte appends one 0x00 after the s mpint INSIDE a p256 signature blob
// (VerifySSHSIG's sr.end(), sshsig.go:298).
func withMpintPairTrailingByte(t *testing.T, a []byte) []byte // f := readStringsTest(p.signature, 2); packStrings(f[0], append(bytes.Clone(f[1]), 0))
// withNonZeroPadBits sets the lowest bit of the last base64 character before the '=' padding:
// StdEncoding still decodes it to the same bytes, Strict() refuses it. t.Fatal when s has no
// padding or the character did not change.
func withNonZeroPadBits(t *testing.T, s string) string
// armorNonStrict armors blob like armorTest but with withNonZeroPadBits applied to the base64.
func armorNonStrict(t *testing.T, blob []byte) []byte
```

`withNonZeroPadBits`: `const alphabet =
"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"`; `body :=
strings.TrimRight(s, "=")`; `len(body) == len(s)` → `t.Fatal("no padding: no slack bits")`; `i :=
len(body)-1`; `c := alphabet[strings.IndexByte(alphabet, body[i])|1]`; `c == body[i]` →
`t.Fatal`; return `body[:i] + string(c) + s[len(body):]`.

**SG8 `TestSSHSIGStrictnessArms`** (fixtures only). `ed`, `p`, `edK`, `pK`, `msg` as in SG3. Each
arm asserts `errors.Is` and, where given, the reason substring:

| Arm | Input | Assertion | Kills |
|---|---|---|---|
| sig-string-trailing | `ParseSSHSIG(withSigStringTrailingByte(t, ed))` | `s == nil`, `errors.Is(err, ErrTrailing)` | MUT-SG-SIGTRAIL |
| mpint-pair-trailing | `s, err := ParseSSHSIG(withMpintPairTrailingByte(t, p))` → `err == nil`; then `VerifySSHSIG(s, Namespace, msg, pK)` | `errors.Is(err, ErrBadSignature)`, contains `trailing bytes` | MUT-SG-MPINTTRAIL |
| nonstrict-sig | `ParseSSHSIG(armorNonStrict(t, unarmorTest(t, p)))` | `s == nil`, `errors.Is(err, ErrBase64)`. **Control**: `base64.StdEncoding.DecodeString` of the mutated body equals `unarmorTest(t, p)` byte for byte | MUT-SG-STRICT |
| nonstrict-pub | `f := strings.Fields(string(fixtureBytes(t, "p256.pub")))`; `ParseAuthorizedKey([]byte(f[0] + " " + withNonZeroPadBits(t, f[1]) + " " + f[2]))` | `errors.Is(err, ErrBase64)`. **Control**: the unmutated three-field line parses | MUT-AK-STRICT |
| parse-sha1 | `ParseSSHSIG(withHashAlg(t, ed, "sha1"))` | `s == nil`, `errors.Is(err, ErrHashAlg)` (the parse itself, never `VerifySSHSIG`) | MUT-SG-HASHALG |
| verify-sha1 | `s, _ := ParseSSHSIG(ed)` (nil err asserted); `s.HashAlg = "sha1"`; `VerifySSHSIG(s, Namespace, msg, edK)` | `errors.Is(err, ErrHashAlg)` | MUT-SD-DEFAULT |
| relabel-sha256 | `s, err := ParseSSHSIG(withHashAlg(t, ed, "sha256"))` → `err == nil`, `s.HashAlg == "sha256"`; `VerifySSHSIG(s, Namespace, msg, edK)` | `errors.Is(err, ErrBadSignature)` (signed over SHA-512, so a relabel must not verify) | (sanity; with MUT-SD-SHA256 it still refuses) |

**XC1 additions** (`crossCheckSSHKeygen`, step 2): `good256, rc256 := runSSHKeygen(t, msg, "-Y",
"sign", "-n", Namespace, "-O", "hashalg=sha256", "-f", key)`, added to the map as `"good-sha256"`
(so it gets the rc, non-empty and key-type checks). Then assert `ParseSSHSIG(good).HashAlg ==
"sha512"` and `ParseSSHSIG(good256).HashAlg == "sha256"` (`t.Fatalf` naming both). Step 4 arms
gain, after `flipped-signature`:

```go
		{"good-sha256", good256, msg, true},
		{"tampered-message-sha256", good256, flipByte(msg, 0), false},
		{"flipped-signature-sha256", withLastSigByteFlipped(t, good256), msg, false},
```

Subtest names are unchanged, so the CI step's three `--- PASS` greps still apply and now cover
the SHA-256 arms on `ubuntu-24.04` (the arms are loop iterations, not subtests). Kills
MUT-SD-SHA256 and MUT-SG-SHA512ONLY.

**Acceptance**
```
test -z "$(gofmt -l host/approval)"
go build ./... && go vet ./host/approval/
go test ./host/approval/ -run '^TestSSHSIGStrictnessArms$' -count=1 -v
WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -run '^TestSSHSIGCrossCheckWithSSHKeygen$' -count=1 -v   # parent + ed25519 + ecdsa-p256 PASS
go test ./host/approval/ -count=1 && go test -race ./host/approval/ -count=1
test -z "$(git status --short host/approval/testdata .gitattributes .github)"                              # no fixture, attribute or CI change
```

### M4a — schema 4: the table, the version, the upgrade, the table-set check · ~45 prod + ~110 test LOC

**Files**: `host/store/credentials_schema.sql` (§2.1), `host/store/credentials.go` (§2.2),
new `host/store/mint_requests.go` (only the `---- M4a ----` block of §2.3),
`host/daemon/credentials.go:76` (§2.2), `host/store/credentials_test.go` (ST1 edit, ST2b), new
`host/store/mint_requests_test.go` (`makeV3Credentials`, ST2).

**Tests**
- **ST1 `TestCredentialsSchemaMatchesCanonicalDDL`** (edit): the table list becomes
  `[]string{"episode_starts", "mint_requests", "session_credentials"}`; add a hardcoded
  `canonicalMintRequestsDDL` (byte-identical to `mintRequestsSQL`'s text, review-visible like
  `canonicalCredentialsDDL`) and require `normalizeDDL(ddl["mint_requests"]) ==
  normalizeDDL(canonicalMintRequestsDDL)` and `normalizeDDL(mintRequestsSQL) ==
  normalizeDDL(canonicalMintRequestsDDL)`. The `user_version` assertion stays symbolic. Kills
  MUT-ST-DDL-SCHEMA, MUT-ST-DDL-CONST.
- **ST2 `TestCredentialsUpgradeToV4`**: for `v` in `{1, 2, 3}` (subtests `v1`, `v2`, `v3`), a store
  at version `v` (`makeV1Credentials(…, rigRows())`; v2 = `canonicalCredentialsDDL` +
  `setPositiveFixtureVersion(2)` + one bound row; `makeV3Credentials` with one bound row plus one
  `episode_starts` row inserted raw). Then:
  (a) `OpenCredentials(path)` → `Version() == v`; close; `userVersion == v`; `schemaState` tables
  unchanged (a session verb never upgrades).
  (b) `UpgradeCredentials` → `(true, nil)`; again → `(false, nil)`.
  (c) `userVersion == 4`; tables exactly `[episode_starts mint_requests session_credentials]`;
  `normalizeDDL(mint_requests) == normalizeDDL(canonicalMintRequestsDDL)`; every original session
  row `ResolveSession`s `reflect.DeepEqual` to what was inserted; for v3 the `episode_starts` row's
  `start_json` is unchanged. Kills MUT-ST-UPG-NOMR, MUT-ST-CLIUPG.
- **ST2b `TestCredentialsV4WithoutMintRequestsIsRefused`**: a fresh store; raw `DROP TABLE
  mint_requests`; `OpenCredentials` and `OpenCredentialsReadOnly` both fail with
  `*CredentialsSchemaError` (`errors.As`); `schemaState` unchanged by the refused opens. Kills
  MUT-ST-TABLESET.

**Acceptance**
```
test -z "$(gofmt -l host/store host/daemon)"
go build ./... && go vet ./host/store/ ./host/daemon/
go test ./host/store/ -run '^(TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsUpgradeToV4|TestCredentialsV4WithoutMintRequestsIsRefused)$' -count=1 -v
go test ./host/store/ -run 'Credentials|TestEpisodeStarts|TestMigrateCredentials|TestPreSplitStoreOpenModes|TestDroppedCredentialsTableIsRefused' -count=1
go test ./host/daemon/ -run 'Credentials|Split|Binding' -count=1
```

### M4b — put, get and list requests · ~130 prod + ~150 test LOC

**Files**: `host/store/mint_requests.go` (the `---- M4b ----` block of §2.3),
`host/store/mint_requests_test.go`, `host/store/credentials_test.go`
(`TestCredentialsMethodsRequireADeadline` edit).

**Tests**
- **ST3 `TestPutMintRequestNeedsVersion4`**: `makeV3Credentials`; `OpenCredentials` (stays 3).
  `PutMintRequest`, `MintRequest` and `PendingMintRequests` each give `errors.As(&older)` with
  `older.Version == 3` and text containing ``restart `ailang-worldd serve` `` and `need schema 4`;
  `userVersion == 3` afterwards. Kills MUT-ST-OLDER, MUT-ST-OLDER-GET, MUT-ST-OLDER-LIST.
- **ST7 `TestPutMintRequestCapsAndPrunes`** (one store):
  1. At `now = 1000`: put A = `testMintRequestRow(1, 1000)` (exp 1600), B = `(2, 999)` (exp 1599),
     C = `(3, 1001)` (exp 1601). Consume B with `MintSessionForRequest`… **not available until
     M4c**, so here mark B consumed with a raw `UPDATE mint_requests SET consumed_at = 1100,
     approved_by = 'tty', credential_id = 'x' WHERE request_id = ?` (M4c's ST4 covers the real
     consume).
  2. At `now = 88000` (= 1600 + 86400): put D = `(4, 88000)`. Then `MintRequest(A)` → not found
     (pruned at the boundary `expires_at == now - 86400`); `MintRequest(B)` → found, `Consumed`;
     `MintRequest(C)` → found (1601 is inside the day); D found.
  3. Put E1..E15 = `(10..24, 88000)` → all ok (16 live pending: D + E; C is expired and must not
     count). Put F = `(30, 88000)` → `errors.Is(ErrTooManyPendingRequests)`; `MintRequest(F)` not
     found.
  4. `PendingMintRequests(ctx, 88000, 0)` → exactly 16 rows, none of A/B/C, ordered by
     `created_at DESC, request_id DESC` (all `created_at` 88000, so `request_id` descending: the
     first is E15's id).
  Kills MUT-ST-CAP, MUT-ST-CAPGE, MUT-ST-CAPLIVE, MUT-ST-PRUNE, MUT-ST-PRUNEEQ, MUT-ST-PRUNECONS,
  MUT-ST-PENDFILTER, MUT-ST-PENDEXP.
- **ST8 `TestPutMintRequestRefusesABadRow`**: from `testMintRequestRow(1, 1000)`, each arm alone →
  `errors.Is(ErrBadMintRequest)` and no row stored: id `mr-0zzzzzzzzzzzzzzz`; id `mr-0123`;
  `Request` nil; `Digest` upper-case hex; `Challenge` 63 chars; `EpisodeID ""`; `ExpiresAt =
  CreatedAt`; `Consumed true`. **Instrument control (F2)**: a raw `INSERT INTO mint_requests
  (request_id, request, digest, challenge, episode_id, created_at, expires_at) VALUES
  ('mr-0zzzzzzzzzzzzzzz', x'00', <64 a>, <64 a>, 'ep', 1, 2)` on the same file **succeeds**, so
  the Go check is load-bearing. The clean row is accepted. Kills MUT-ST-ID.
- **edit `TestCredentialsMethodsRequireADeadline`**: add `PutMintRequest`, `MintRequest`,
  `PendingMintRequests` with `context.Background()` → each `errors.Is(ErrNoDeadline)`. Kills
  MUT-ST-DEADLINE.

**Acceptance**
```
test -z "$(gofmt -l host/store)"
go build ./... && go vet ./host/store/
go test ./host/store/ -run '^(TestPutMintRequestNeedsVersion4|TestPutMintRequestCapsAndPrunes|TestPutMintRequestRefusesABadRow|TestCredentialsMethodsRequireADeadline)$' -count=1 -v
go test ./host/store/ -run 'Credentials|MintRequest|TestEpisodeStarts|TestMigrateCredentials' -count=1
```

### M4c — consume in the mint's transaction · ~95 prod + ~190 test LOC

**Files**: `host/store/mint_requests.go` (`---- M4c ----` of §2.3), `host/store/credentials.go`
(D3 `insertSessionRow` extraction), `host/store/mint_requests_test.go`,
`host/store/credentials_test.go` (deadline edit, fourth method).

**Tests** (`session := SessionRow{CredentialID: fmt.Sprintf("%064x", k), EpisodeID: "ep",
GrantsJSON: "[]", ExpiresAt: 2000, CreatedAt: 1010}` unless stated)
- **ST4 `TestMintSessionForRequestConsumesOnce`**: put `R = testMintRequestRow(1, 1000)`; consume
  with `okConsume(R, 1010, &clocks, &checks)` → nil. `MintRequest(R)` → `Consumed`, `ConsumedAt ==
  1010`, `ApprovedBy == "tty"`, `CredentialID == session.CredentialID`; `ResolveSession` finds the
  session row `DeepEqual`; `clocks == 1`, `checks == 1`, and the `Check` saw `1010`. A second
  consume with another credential id → `errors.Is(ErrRequestConsumed)`; `ListSessions` → exactly
  1 row. Kills MUT-ST-ROWS, MUT-ST-CONSUMEDW, MUT-ST-CLOCK2.
- **ST4b `TestMintSessionForRequestRollsBackOnFailure`**: (a) insert a session row with id X via
  `MintSession` first; consume `R` with `session.CredentialID = X` → an error (PK collision);
  `MintRequest(R)` still unconsumed. (b) a `Check` that returns `errSentinel` → `errors.Is(err,
  errSentinel)`; `R` unconsumed; no session row with that id. Kills MUT-ST-NOTX, MUT-ST-NOCHECK.
- **ST5 `TestMintSessionForRequestRefusals`** (fresh store per arm, `R` created 1000, exp 1600;
  each refusal also asserts **0** session rows and `R` unconsumed):

  | Arm | Change | Want | Kills |
  |---|---|---|---|
  | at-expiry | Clock 1600 | `ErrRequestExpired` | MUT-ST-EXPGE |
  | expiry-minus-1 | Clock 1599 | nil (consumed) | MUT-ST-EXPGE (boundary) |
  | at-created | Clock 1000 | nil | MUT-ST-CREATEDLE |
  | clock-behind | Clock 999 | `ErrRequestClockBehind` | MUT-ST-CLASSIFY |
  | other-digest | `rc.Digest = strings.Repeat("ef", 32)` | `ErrRequestDigestChanged` | MUT-ST-DIGESTW |
  | unknown | `rc.ID = "mr-ffffffffffffffff"` | `ErrRequestUnknown` | MUT-ST-CLASSIFY |
  | nil-clock | `rc.Clock = nil` | `ErrIncompleteConsume` | MUT-ST-CLOCKNIL |
  | nil-check | `rc.Check = nil` | `ErrIncompleteConsume`, and `clocks == 0` | MUT-ST-CHECKNIL |
  | expired-ctx | `ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(-time.Second))` | non-nil error, and `clocks == 0` (the Clock is read only after `BeginTx`) | MUT-ST-CLOCKPRE |
  | v3-store | `makeV3Credentials` | `*OlderCredentialsStoreError` | MUT-ST-OLDER-CONSUME |
- **ST6 `TestConcurrentApprovalsMintOnce`** (run under `-race`): one file with `R` put; 8
  goroutines released together by a closed channel; each does `OpenCredentials(path)` → on
  `IsWriterAlreadyActive(err)` record "locked"; otherwise `MintSessionForRequest` with its own
  credential id `fmt.Sprintf("%064x", 100+i)` and a Clock of 1010, record the error, `Close`. Then:
  exactly 1 nil; every other result is "locked" or `errors.Is(ErrRequestConsumed)`; `ListSessions`
  → 1 row; `R.CredentialID` is the winner's. (Kills nothing alone beyond MUT-ST-ROWS; it is the
  design's concurrency pin and the `-race` leg.)
- **edit `TestCredentialsMethodsRequireADeadline`**: add `MintSessionForRequest(Background, row,
  okConsume(…))` → `ErrNoDeadline`.

**Acceptance**
```
test -z "$(gofmt -l host/store)"
go build ./... && go vet ./host/store/
go test ./host/store/ -run '^(TestMintSessionForRequestConsumesOnce|TestMintSessionForRequestRollsBackOnFailure|TestMintSessionForRequestRefusals|TestConcurrentApprovalsMintOnce|TestCredentialsMethodsRequireADeadline)$' -count=1 -v
go test -race ./host/store/ -run '^TestConcurrentApprovalsMintOnce$' -count=5
go test ./host/store/ -count=1
```

### M5a — the request from a validated mint, and back (`cmd/ailang-worldd`) · ~150 prod + ~190 test LOC

**Files**: new `cmd/ailang-worldd/mint_request.go` (§2.7 without `createMintRequest` and without
the two M6a map entries), new `cmd/ailang-worldd/mint_request_test.go`.

**Tests**
- **CL14 `TestMintRequestFieldMapIsComplete`**. Write a test-local `fieldMapProblems(m
  map[string]string) []string` and assert it returns nothing for `mintRequestFieldMap`:
  1. `fields`: the Go names of `reflect.TypeOf(mintRequest{})`'s fields, recursing into the
     anonymous `mintOptions` (its own name is not a key). Each must be a key ("unmapped field X").
  2. Each field of `reflect.TypeOf(broker.Capability{})` must be a key `grants[].<Name>`.
  3. Each key must be a field from 1 or a `grants[].<Capability field>` ("stale key X").
  4. Each value must start with `MintRequest.`, `excluded: `, `fixed: ` or `nil: `.
  5. Each `MintRequest.<path>` must resolve on `reflect.TypeOf(approval.MintRequest{})`: split on
     `.`; a segment ending in `[]` names a slice field and continues in its element type.
  6. `leaves`: every leaf of `approval.MintRequest` (structs recursed as `Parent.Child`, slices of
     structs as `Name[].Child`, everything else a leaf). Each leaf must equal some value's path or
     be in `requestOnlyLeaves`; no `requestOnlyLeaves` entry may also be a target.
  7. Non-vacuity: `len(fields) >= 29` (31 after M6a), 4 Capability fields, `len(leaves) == 28`.
  Positive controls: `fieldMapProblems` of a copy without `"git"` reports exactly `unmapped field
  git`; with an extra key `"nope"` it reports `stale key nope`. Kills MUT-CL14-DROP,
  MUT-CL14-NEWFIELD, MUT-CL14-CAPFIELD, MUT-CL14-REQLEAF.
- **CL15 `TestApproveRebuildsTheSameMintRequest`**, two arms (`claude` bound to profile `proj`
  with a task; `codex` unbound, no MCP). `base, _ := filepath.EvalSymlinks(t.TempDir())`; build an
  `episodeFixture`-like world under `base` (store via `store.Open`, `ws`, a one-commit git `repo`,
  `profiles` with `writeTestProfile`), `t.Setenv("HOME", base)`. Side A = `prepareMint(mintOptions{
  verb: "episode start", episode: "ep1", db: <canonical db>, credsDB: <credentialsFor(db, "")>,
  wsRoot: <ws>, repo: <repo>, branch: "world/ep1", grants: episodeStartGrants("", nil, 50, bound),
  ttl: 3600, outPath: <base>/sessions/ep1.session, outDirCreate: true, profile: "proj" or "",
  profilesDir: <profiles>})` (episode.go:305-307's literal) with `A.afterToken = func(token string)
  error { return writeMCPConfig(mcpFile, "127.0.0.1:7644", token) }` when claude. `mr :=
  buildMintRequest(A, testBuild, "127.0.0.1:7644", agent, "m1", task, mcpFile)` (nil error); `dec
  := approval.Decode(approval.Canonical(mr))`; side B = `prepareMint(fromRequest(dec))` with
  `B.afterToken = mcpAfterToken(dec.MCPConfig, dec.Daemon.Addr)`. Then: the nil-ness of
  `afterToken`, `beforeToken`, `resumeRollback` is equal; copy both, zero `verb` and the three func
  fields, `reflect.DeepEqual(a, b)` (on failure print `%+v` of both). For claude: call
  `A.afterToken("tok")`, read `mcpFile`, remove it, call `B.afterToken("tok")`, read again → equal
  bytes. Kills MUT-FR-PROFILESDIR, MUT-FR-CREDSDB, MUT-FR-OUTDIRCREATE.
- **BM1 `TestBuildMintRequestPathsAndRefusals`**: `t.Chdir(base)`; `prepareMint` with **relative**
  `db`, `wsRoot`, `repo`, `profilesDir` → `buildMintRequest` → nil error; `WorldDB ==
  store.CanonicalPath(db)`, `WorkspaceRoot`, `Repo`, `ProfilesDir` absolute and equal to
  `filepath.Abs` of the inputs; `Daemon` fields equal `testBuild`'s and `Addr`. Arm: `r.grants[0].
  ExpiresAt = 5` → error containing `carries an expiry`. Arm: `renderRequest(…)` contains the id,
  the full digest, `grantList`, `SessionFile`, `MCPConfig`, `Profile.Digest`, `Daemon.Addr`, and
  `(in 10m0s)` for `expiresAt = now + 600`. Kills MUT-BM-REPOABS, MUT-BM-WSABS, MUT-BM-EXPIRES.

**Acceptance**
```
test -z "$(gofmt -l cmd/ailang-worldd)"
go build ./... && go vet ./cmd/ailang-worldd/
go test ./cmd/ailang-worldd/ -run '^(TestMintRequestFieldMapIsComplete|TestApproveRebuildsTheSameMintRequest|TestBuildMintRequestPathsAndRefusals)$' -count=1 -v   # + skip guard
go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1
```

### M5b — `episode start --request` (create) · ~110 prod + ~170 test LOC

**Files**: `cmd/ailang-worldd/session.go` (`random` only, §2.5), `cmd/ailang-worldd/mint_request.go`
(`createMintRequest`, §2.7), `cmd/ailang-worldd/episode.go` (§2.8), `cmd/ailang-worldd/main.go`
(usage `[--request]`), new `cmd/ailang-worldd/approve_test.go` (§2.10 helpers + CL1–CL3),
`host/store/context_roots_test.go` (+ `"cmd/ailang-worldd/mint_request.go|createMintRequest|Background": 1, // row 201 M5: --request's store open and insert (3 s)`).

**Tests**
- **CL1 `TestEpisodeStartRequestStoresAndOpensNoTerminal`**: `newEpisodeRun`; profile `proj`;
  `stubDaemon(matchingHealth(…, p))`; a task file; `runEpisodeStart(r.args("ep1", "--profile",
  "proj", "--task", task, "--model", "claude-opus-5-5", "--request"), …, requestEnv(&opened,
  approveT0))` → `exitOK`, `opened == false`. Read-only `PendingMintRequests(approveT0)` → exactly
  1 row: `RequestID == "mr-0102030405060708"`; `Challenge == hex(09..28)`; `CreatedAt ==
  approveT0`; `ExpiresAt == approveT0+600`; `Digest == hex(sha256(Request))`. `Decode(Request)`:
  `Episode "ep1"`, `Grants == fromCapabilities(episodeStartGrants("", nil, 50, true))`, `TTLSeconds
  3600`, `Profile == {proj, p.Digest}`, `Agent "claude"`, `Model "claude-opus-5-5"`, `Task` the
  absolute task path, `SessionFile` `<sessions>/ep1.session`, `MCPConfig` `<sessions>/ep1.mcp.json`,
  `Daemon.Addr == defaultEpisodeAddr`, `Daemon.VCSRevision == testBuild.VCSRevision`. Stdout
  contains the id, `ailang-worldd approve mr-0102030405060708`, and `renderRequest`'s `grants`
  line. `<ws>/ep1` does not exist; `r.assertNoSessionFiles(t)`; 0 session rows. Kills
  MUT-CR-PROVISION, MUT-CR-EXPIRY, MUT-CR-CHALLENGE.
- **CL1b `TestEpisodeStartRequestRefusesWithoutRandomness`**: `requestEnv` with `random = nil` →
  `exitFatal`, stderr contains `no randomness source`, no `mint_requests` row. Kills
  MUT-CR-RANDOMNIL.
- **CL2 `TestEpisodeStartRequestRefusesUnsupportedModes`**: arms (each with `requestEnv`, `opened`
  stays false, `exitUsage`, stderr contains the quoted text, and no request row):
  `r.args("ep1", "--request", "--launch")` → `--request cannot be combined with --launch`;
  `r.args("ep1", "--request", "--headless")` → `… with --headless`; `[]string{"--request",
  "--resume", "ep1", "--preset", "p", "--db", …}` → `… with --resume`; `r.args("ep1", "--request",
  "--preset", "mypreset")` → `… with a launch preset`. **Control**: `r.args("ep1", "--request",
  "--preset", "se-tools")` → `exitOK`, 1 row. Kills MUT-CR-LAUNCH, MUT-CR-HEADLESS, MUT-CR-RESUME,
  MUT-CR-PRESET.
- **CL3 `TestEpisodeStartRequestRefusesAnOlderStore`**: before running, write a v3 credentials
  store at `credentialsFor(r.f.db, "")` with raw SQL (the session_credentials and episode_starts
  DDL, `PRAGMA user_version = 3`). Unbound (`--grant Workspace.Read=worktree:1`, no `--profile`, so
  `checkBindingFits` does not answer first). → `exitUsage`, stderr contains `mint requests need
  schema 4` and ``restart `ailang-worldd serve` ``; `user_version` still 3; no `mint_requests`
  table.
- **edit `TestProductionContextRoots`** (the pin above).

**Acceptance**
```
test -z "$(gofmt -l cmd/ailang-worldd host/store)"
go build ./... && go vet ./cmd/ailang-worldd/ ./host/store/
go test ./cmd/ailang-worldd/ -run '^(TestEpisodeStartRequest)' -count=1 -v                     # + skip guard
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
go test ./cmd/ailang-worldd/ -run '^(TestEpisodeStart|TestSessionNew|TestSessionMint_|TestStartLaunch|TestEpisodeResume|TestMintRequestFieldMap|TestApproveRebuilds|TestBuildMintRequest)' -count=1
```

### M6a — the bounded y/N, the consume branch, `MintForRequest` · ~95 prod + ~170 test LOC

**Files**: `cmd/ailang-worldd/session.go` (`after`, `readTermLineWithin`, §2.5),
`cmd/ailang-worldd/session_new.go` (§2.6), `cmd/ailang-worldd/mint_request.go` (the two map
entries), `host/authority/mint.go` (§2.4), new `host/authority/mint_request_test.go` (AU1),
`cmd/ailang-worldd/approve_test.go` (CL17, CL18), `cmd/ailang-worldd/mint_request_test.go` (CL14
non-vacuity bound `>= 31`).

**Tests**
- **AU1 `TestMintForRequestChecksTheRowInsideTheTransaction`** (fresh temp store per arm). The
  request: `approval.MintRequest` valid, `Episode "ep"`, one grant `{Workspace.Read, worktree,
  5}`, `TTLSeconds 600`, `Profile {proj, sha256:<64 hex>}`, the rest absolute dummy paths; its row
  `mr-0123456789abcdef`, `Request = Canonical`, `Digest = hex(sha256)`, created 1000, exp 1600.
  Good call: `MintForRequest(ctx, st, "ep", ProfileBinding{proj, digest}, []broker.Capability{{…,
  Budget: 5}}, 600, 1010, nil, RequestMint{Request: mr, RequestCreatedAt: 1000, Consume:
  store.RequestConsume{ID, Digest, ApprovedBy: "tty", Clock: func() int64 { return 1020 }}})` → nil;
  the row is consumed by `tty` with the returned credential id; the session row has `CreatedAt
  1010`, `ExpiresAt 1610`, its grant `ExpiresAt 1610`; `credentialID == hex(sha256(rawToken))`.
  Mismatch arms, each `errors.As(&mm)` with `*approval.MintMismatchError` and the exact `Column`,
  request unconsumed, 0 session rows: budget 6 → `grants_json`; episode `"other"` →
  `episode_id`; zero binding → `project`; ttl 601 → `expires_at`; `now = 1021` (Clock 1020) →
  `created_at`; `now = 999` → `created_at`. Arm: Clock 1600 → `errors.Is(store.ErrRequestExpired)`.
  Kills MUT-AU-NOCHECK, MUT-AU-PRECONSUME.
- **CL17 `TestExistingCallersReadUnbounded`**: `after := func(time.Duration) <-chan time.Time {
  t.Error("after called on an unbounded caller"); return nil }`. (a) `runSessionNew(f.args("ep1",
  "--grant", "Workspace.Read=worktree:1"), …, env)` with a `fakeTerm` answering `"y\n"`, `now
  1000`, that `after` → `exitOK`, 1 session row. (b) plain `runEpisodeStart(r.args("ep2"), …)`
  with `stubDaemon(matchingHealth(…))`, the same terminal and `after` → `exitOK`. Kills
  MUT-PE-BOUNDALL.
- **CL18 `TestReadTermLineWithin`**: (a) `strings.NewReader("y\n")` with `after` returning a
  never-firing `make(chan time.Time)` → `("y", true)`. (b) `pr, pw := io.Pipe()`;
  `t.Cleanup(func() { _ = pw.Close() })`; `after` returns a closed channel; run it in a goroutine
  and require it returns `("", false)` within a 5 s guard (`select` on `time.After(5 *
  time.Second)` → `t.Fatal`). Kills MUT-RTW-READONLY.
- CL14 stays green with the two new entries (MUT-FM-CONSUME: delete the `consume` entry → red).

**Acceptance**
```
test -z "$(gofmt -l cmd/ailang-worldd host/authority)"
go build ./... && go vet ./cmd/ailang-worldd/ ./host/authority/
go test ./host/authority/ -count=1 -v -run '^TestMintForRequestChecksTheRowInsideTheTransaction$'
go test ./host/authority/ -count=1
go test ./cmd/ailang-worldd/ -run '^(TestExistingCallersReadUnbounded|TestReadTermLineWithin|TestMintRequestFieldMapIsComplete)$' -count=1 -v   # + skip guard
go test ./cmd/ailang-worldd/ -run '^(TestSessionMint_|TestSessionNew|TestEpisodeStart|TestStartLaunch|TestEpisodeResume|TestSessionMintGrantsExpireWithTheSession)' -count=1
```

### M6b — the `approve` verb · ~170 prod + ~190 test LOC

**Files**: new `cmd/ailang-worldd/approve.go` (§2.9), `cmd/ailang-worldd/main.go` (`case
"approve"`, the two usage lines, `ownHelpVerbs`), `cmd/ailang-worldd/keyguard.go`
(`ownVerbHelp`), `cmd/ailang-worldd/help_test.go` (`everyVerb` + `{"approve"}`; want-list +
`"approve <request-id>"`, `"approve --list"`), `host/store/context_roots_test.go` (+
`"cmd/ailang-worldd/approve.go|runApprove|Background": 1, // row 201 M6: approve's read-only request lookup or --list (3 s)`),
`cmd/ailang-worldd/approve_test.go`.

**Tests** (rig: `a := newApproveRig(t, "ep1", approveT0)`)
- **CL4 `TestApproveListShowsPending`**: one `r := newEpisodeRun(t)`; `newApproveRigIn(t, r, …)`
  for A `ep1` at `approveT0-700`, B `ep2` at `approveT0`, C `ep3` at `approveT0`; approve C with
  `y` and clock `approveT0+1`. `runApprove([]string{
  "--list", "--db", …, "--credentials-db", …}, …, approveEnv(nil, nil, func() int64 { return
  approveT0 + 2 }))` → `exitOK`; stdout has B's id and not A's (expired) or C's (consumed); B's
  line contains `episode ep2` and `B.row.Digest[:16]`. With nothing pending (list at
  `approveT0+1000`) → `no pending mint requests`. (Store mutants MUT-ST-PENDFILTER/PENDEXP also
  turn this red.)
- **CL6 `TestApproveWithoutATerminalMintsNothing`**: `noTerminalEnv(approveT0+1)` with `after:
  time.After` added → `exitUsage`, stderr contains `sessionTTYFix`; the row is unconsumed; no
  worktree; 0 session rows.
- **CL7 `TestApproveDeclinedLeavesTheRequestPending`**: `promptTerm{in: "n\n"}` → `exitUsage`,
  stderr contains `aborted (not confirmed)`; unconsumed; no worktree; 0 session rows; `term.out`
  contains the `mint request <id>` block.
- **CL8 `TestApproveRefusesReplay`**: first `approve` with `y` at `approveT0+1` → `exitOK`, the row
  `Consumed` by `tty`; second `approve` with a spy (`opened`) → `exitUsage`, stderr contains
  `already used` and `by tty`, `opened == false`; 1 session row. Kills MUT-AP-ADMIT,
  MUT-PE-NOCONSUMEBRANCH (the first approve would leave the row unconsumed).
- **CL12 `TestApproveRefusesBadIDsAndMisuse`**: (a) `mr-XYZ` → `exitUsage`, `is not a mint request
  id`; (b) `mr-0000000000000000` → `exitUsage`, `no mint request`; (c) `run([]string{"--addr",
  "127.0.0.1:1", "approve", a.id})` → `exitUsage`, `not valid for 'approve'`; (d) env without
  `after` → `exitFatal`, `no prompt timer`, terminal not opened; (e) another world store: a fresh `dir2/world.db` (`store.Open`) and
  `approve <id> --db dir2/world.db --credentials-db <a.credPath>` (the row is found; the world
  store is not the request's) → `exitUsage`, stderr contains `was made for`, terminal not opened; (f) `approve --help` and `help approve` → `exitOK`, stdout == `approveHelp`. Kills
  MUT-AP-GRAMMAR, MUT-AP-ADDR, MUT-AP-AFTERNIL, MUT-AP-SAMESTORE.
- **edits**: `TestEveryVerbHasHelp` (`everyVerb` + want-list), `TestProductionContextRoots`
  (the pin). `TestHelpIsNeverBlockedByTheRegistryKeyGuard` must pass unedited.

**Acceptance**
```
test -z "$(gofmt -l cmd/ailang-worldd host/store)"
go build ./... && go vet ./cmd/ailang-worldd/ ./host/store/
go test ./cmd/ailang-worldd/ -run '^(TestApproveListShowsPending|TestApproveWithoutATerminalMintsNothing|TestApproveDeclinedLeavesTheRequestPending|TestApproveRefusesReplay|TestApproveRefusesBadIDsAndMisuse|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard)$' -count=1 -v   # + skip guard
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1
```

### M6c — expiry, stale world, tampering, the prompt bound (tests only) · ~0 prod + ~230 test LOC

**Files**: `cmd/ailang-worldd/approve_test.go`.

- **CL9 `TestApproveRefusesExpired`**: arm 1: clock `approveT0+600` with a spy → `exitUsage`,
  `expired at`, terminal not opened. Arm 2: `clockSeq(2, approveT0+1, approveT0+601)` (D8: calls
  [1] and [2] see `+1`, the in-transaction [3] sees `+601`), `promptTerm` `y` → `exitFatal`, stderr
  contains `expired before it was consumed`; unconsumed; 0 session rows; `<ws>/ep1` removed and
  stderr contains `removed the worktree`. Kills MUT-AP-ADMITEXP, MUT-AP-CLOCKREUSE.
- **CL10 `TestApproveRefusesAStaleWorld`** (a fresh rig per arm; each arm: `exitUsage`, terminal
  not opened (spy), the row unconsumed, 0 session rows): (a) rewrite profile `proj` with extra
  `, "x": 1` and re-`stubDaemon` with the **new** digest → stderr contains `profile.digest`;
  (b) `stubDaemon` with `Build = otherBuild` (`VCSRevision` 40 `f`s) and then `episodeSelfBuild =
  func() daemon.BuildIdentity { return otherBuild }` → `daemon.vcs_revision`; (c) write
  `<sessions>/ep1.session` → stderr contains `already exists`; (d) write `<sessions>/ep1.mcp.json`
  → `already exists`; (e) `git worktree add` `<ws>/ep1` from the repo → `reuse`. Kills
  MUT-AP-NODERIVE, MUT-AP-MCPEXISTS.
- **CL11 `TestApproveRefusesATamperedRow`**: open the credentials file with `sql.Open("sqlite",
  …)`. Arm 1: replace the first `"budget":50` in `request` with `"budget":999`, keep `digest` →
  `exitUsage`, `do not match its digest`, terminal not opened. Arm 2 (fresh rig): the same edit and
  `digest = hex(sha256(new bytes))` → with `y`: `exitOK`; `term.out` contains
  `Workspace.Read=worktree:999`; the session row's `grants_json` decodes to a `Workspace.Read`
  grant with `Budget 999`. Kills MUT-AP-DIGEST.
- **CL13 `TestApproveUnansweredPromptTimesOut`**: arm 1: `pr, pw := io.Pipe()`;
  `t.Cleanup(func() { _ = pw.Close() })`; terminal `promptTerm{in: pr}`; clock always
  `exp-37` (`exp = approveT0+600`); `after` records its argument and returns a closed channel.
  `runApprove` in a goroutine with a 5 s guard → `exitUsage`; stderr contains `no answer within
  the request's remaining lifetime (37s)` and `still pending`; the recorded duration is exactly
  `37 * time.Second`; unconsumed; no worktree; 0 session rows. Arm 2: `clockSeq(1, exp-1, exp)` →
  `exitUsage`, `expired before the prompt`, `term.out` is **empty**. Kills MUT-PE-PLAINREAD,
  MUT-PE-BOUNDFROMLIFETIME, MUT-PE-PROMPTCHECK.
- **CL19 `TestApproveRefusesExpiryDuringWorktreePreparation`**: `wt := a.worktree()`; clock
  `func() int64 { if _, err := os.Stat(wt); err == nil { return exp }; return approveT0 + 1 }`;
  `promptTerm` `y` → `exitFatal`, stderr contains `expired before it was consumed`; the row is
  unconsumed; 0 session rows; `wt` is gone. The worktree exists only after `worktree.Add`, so only
  a reading taken **after** preparation sees `exp`. Kills MUT-AP-CLOCKEARLY (and
  MUT-AP-CLOCKREUSE).

**Acceptance**
```
test -z "$(gofmt -l cmd/ailang-worldd)"
go build ./... && go vet ./cmd/ailang-worldd/
go test ./cmd/ailang-worldd/ -run '^(TestApproveRefusesExpired|TestApproveRefusesAStaleWorld|TestApproveRefusesATamperedRow|TestApproveUnansweredPromptTimesOut|TestApproveRefusesExpiryDuringWorktreePreparation)$' -count=1 -v   # + skip guard
go test -race ./cmd/ailang-worldd/ -run '^(TestApproveUnansweredPromptTimesOut|TestReadTermLineWithin)$' -count=1
```

### M6d — the full effect check, and QUICKSTART §14 · ~0 prod + ~190 test + ~45 doc LOC

**Files**: `cmd/ailang-worldd/approve_test.go`, `docs/QUICKSTART.md`.

```go
// approveFacts is everything an approve did outside the session row (CL5, CL16).
type approveFacts struct {
	WorktreeDir         string // mr.Worktree if it exists, else ""
	WorktreeIsGit       bool   // isGitWorktree(WorktreeDir)
	Branch              string // git branch --show-current in WorktreeDir
	SessionFile         string // mr.SessionFile if it exists, else ""
	SessionMode         os.FileMode
	TokenHash           string // hex(sha256(the session file's bytes))
	MCPFile             string // mr.MCPConfig if it exists, else ""
	MCPAddr, MCPToken   string // from the MCP config's world server URL and Authorization header
	Launch              string // approve's stdout
	ApprovedBy          string // the request row
	RowCredentialID     string // the request row
	SessionCredentialID string // the one session row
}
func collectApproveFacts(t *testing.T, a approveRig, mr approval.MintRequest, stdout string) approveFacts
// approveEffectProblems returns one entry per mismatch, each starting with its name:
// worktree, worktree-git, branch, session_file, session_mode, token_hash, mcp, mcp_for_codex,
// agent, model, task, approved_by, credential_id.
func approveEffectProblems(mr approval.MintRequest, f approveFacts) []string
```

`approveEffectProblems` checks, in order: `WorktreeDir != mr.Worktree`; `!WorktreeIsGit`; `Branch
!= mr.Branch`; `SessionFile != mr.SessionFile`; `SessionMode != 0o600`; `TokenHash !=
SessionCredentialID`; claude: `MCPFile != mr.MCPConfig || MCPAddr != mr.Daemon.Addr ||
hex(sha256(MCPToken)) != SessionCredentialID` → `mcp`; codex: `MCPFile != ""` → `mcp_for_codex`;
`Launch` lacks `"  claude --tools"` (claude) or `" codex exec "` (codex) → `agent`; `mr.Model !=
""` and `Launch` lacks `"--model "+shellQuote(mr.Model)` (claude) or `"-m "+shellQuote(mr.Model)`
(codex) → `model`; `mr.Task != ""` and `Launch` lacks `"task: "+mr.Task` → `task`; `ApprovedBy !=
"tty"`; `RowCredentialID != SessionCredentialID`.

- **CL5 `TestApproveMintsExactlyTheRequest`**: rig; `y` at `approveT0+1` → `exitOK`. `mr :=
  Decode(a.row(t).Request)`. (1) The stored session row passes `approval.MintMatches(mr,
  MintFacts{…from the row…, TokenHash: hex(sha256(session file)), RequestCreatedAt: approveT0,
  ConsumedAt: a.row(t).ConsumedAt})`. (2) `approveEffectProblems(mr, collectApproveFacts(…))` is
  empty. Kills MUT-AP-LAUNCHMODEL, MUT-AP-TASK (and, through the re-derivation refusal,
  MUT-AP-GRANTSDEFAULT).
- **CL16 `TestApproveEffectsCheckCatchesEachMismatch`** (pure, no rig): a clean `mr` (claude) and
  clean `facts` → no problem. Twelve arms, each altering one fact → the problems include an entry
  starting with its name: `WorktreeDir`, `Branch`, `SessionFile`, `SessionMode = 0o644`,
  `TokenHash`, `MCPAddr`, codex `mr` with `MCPFile` set → `mcp_for_codex`, `Launch` without the
  agent, without the model, without the task, `ApprovedBy = "ssh:x"`, `RowCredentialID`. Kills
  MUT-EFF-SKIP-BRANCH, MUT-EFF-SKIP-MODE, MUT-EFF-SKIP-APPROVEDBY.
- **QUICKSTART §14** — append at the end of `docs/QUICKSTART.md`:
  `## 14. Mint by request *(attended — pending first verbatim run)*`, then: one paragraph (what
  `--request` stores, that it needs no terminal and mints nothing, that `approve` shows the whole
  request and needs the typed `y`, the 10-minute expiry, at most 16 pending, and that this guide
  runs it against the §0 scratch store `/tmp/se-world/world.db` with its daemon serving, never a
  real store); a bash block:
  ```bash
  ailang-worldd episode start req-demo --repo "$PWD" --db /tmp/se-world/world.db --request
  ailang-worldd approve --list --db /tmp/se-world/world.db
  ailang-worldd approve mr-<the 16 hex it printed> --db /tmp/se-world/world.db
  ```
  then a short list of what `approve` refuses before any terminal opens (unknown, already used,
  expired, bytes not matching the digest, a changed world named by its first field) and that a
  `y` not typed before the expiry mints nothing. No other part of the file changes.

**Acceptance**
```
test -z "$(gofmt -l cmd/ailang-worldd)"
go build ./... && go vet ./cmd/ailang-worldd/
go test ./cmd/ailang-worldd/ -run '^(TestApproveMintsExactlyTheRequest|TestApproveEffectsCheckCatchesEachMismatch)$' -count=1 -v   # + skip guard
go test ./cmd/ailang-worldd/ -run '^(TestEpisodeStartRequest|TestApprove|TestMintRequestFieldMap|TestBuildMintRequest|TestExistingCallersReadUnbounded|TestReadTermLineWithin|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard)' -count=1 -v   # + skip guard
go test ./host/broker/ ./host/daemon/ -run 'Quickstart|QUICKSTART|Guide' -count=1
test "$(tail -c 1 docs/QUICKSTART.md | xxd -p)" = "0a"
```

---

## 4. Mutation table (83 mutants; the controller and judge run them)

Procedure: apply **one** mutant (a single edit), run the named test(s), expect **FAIL** (red, a
panic, or a compile error where marked), restore the file and `cmp` it with the saved copy.
**Never commit a mutant.** Line numbers are at `ab4ae3aa` (M3b) or refer to §2's definitions.

| Id | Mutant (single edit) | Must fail |
|---|---|---|
| MUT-SG-SIGTRAIL | `sshsig.go:248-250` (`sr.end()` after the inner blob) deleted | SG8 sig-string-trailing |
| MUT-SG-MPINTTRAIL | `sshsig.go:298-300` (P-256 `sr.end()`) deleted | SG8 mpint-pair-trailing |
| MUT-SG-STRICT | `sshsig.go:197` `StdEncoding.Strict()` → `StdEncoding` | SG8 nonstrict-sig |
| MUT-AK-STRICT | `sshsig.go:166` the same | SG8 nonstrict-pub |
| MUT-SG-HASHALG | `sshsig.go:232-234` (parse-level hash check) deleted | SG8 parse-sha1 |
| MUT-SD-DEFAULT | `signedData`'s `default: return nil, ErrHashAlg` deleted | SG8 verify-sha1 |
| MUT-SD-SHA256 | `signedData`'s `sha256` case hashes with `sha512.Sum512` | XC1 good-sha256 (both subtests) |
| MUT-SG-SHA512ONLY | `sshsig.go:232` `h != "sha256" && h != "sha512"` → `h != "sha512"` | XC1 good-sha256 |
| MUT-ST-UPG-NOMR | `upgradeCredentials` omits the `mintRequestsSQL` step | ST2 |
| MUT-ST-CLIUPG | `OpenCredentials` → `openCredentialsWritable(path, true)` | ST2 (a) |
| MUT-ST-TABLESET | `enforceSchema` also accepts `credentialsTables(3)` at version 4 | ST2b |
| MUT-ST-DDL-SCHEMA | drop `CHECK (length(digest) = 64)` in `credentials_schema.sql` only | ST1 |
| MUT-ST-DDL-CONST | the same in `mintRequestsSQL` only | ST1, ST2 |
| MUT-ST-OLDER | `PutMintRequest`'s version check deleted | ST3 |
| MUT-ST-OLDER-GET | `MintRequest`'s version check deleted | ST3 |
| MUT-ST-OLDER-LIST | `PendingMintRequests`'s version check deleted | ST3 |
| MUT-ST-CAP | the pending-count refusal deleted | ST7 step 3 |
| MUT-ST-CAPGE | `n >= MaxPendingMintRequests` → `>` | ST7 step 3 |
| MUT-ST-CAPLIVE | the count drops `AND expires_at > ?` | ST7 step 3 |
| MUT-ST-PRUNE | the `DELETE` deleted | ST7 step 2 |
| MUT-ST-PRUNEEQ | `expires_at <= ?` → `<` in the `DELETE` | ST7 step 2 (A) |
| MUT-ST-PRUNECONS | the `DELETE` drops `consumed_at IS NULL` | ST7 step 2 (B) |
| MUT-ST-ID | `PutMintRequest`'s `MintRequestIDPattern` check deleted | ST8 |
| MUT-ST-PENDFILTER | `PendingMintRequests` drops `consumed_at IS NULL` | ST7 step 4, CL4 |
| MUT-ST-PENDEXP | `PendingMintRequests` drops `expires_at > ?` | ST7 step 4, CL4 |
| MUT-ST-DEADLINE | `PutMintRequest`'s `requireDeadline` deleted | deadline test |
| MUT-ST-ROWS | the `n != 1` branch deleted (always insert) | ST4 |
| MUT-ST-NOTX | the `UPDATE` runs on `c.db` instead of `tx` | ST4b (a) (or a deadline error: red either way) |
| MUT-ST-CLOCK2 | `Clock()` called again for the `WHERE` arguments | ST4 (`clocks == 1`) |
| MUT-ST-CLOCKPRE | `at := rc.Clock()` moved above `BeginTx` | ST5 expired-ctx |
| MUT-ST-NOCHECK | the `rc.Check(at)` call deleted | ST4b (b) |
| MUT-ST-CHECKNIL | `rc.Check == nil` dropped from the guard | ST5 nil-check (panic) |
| MUT-ST-CLOCKNIL | `rc.Clock == nil` dropped from the guard | ST5 nil-clock (panic) |
| MUT-ST-EXPGE | SQL `expires_at > ?` → `>= ?` | ST5 expiry-minus-1 |
| MUT-ST-CREATEDLE | SQL `created_at <= ?` → `< ?` | ST5 at-created |
| MUT-ST-DIGESTW | SQL drops `AND digest = ?` | ST5 other-digest |
| MUT-ST-CONSUMEDW | SQL drops `AND consumed_at IS NULL` | ST4 second consume |
| MUT-ST-CLASSIFY | `classifyConsume` returns `ErrRequestConsumed` for every case | ST5 |
| MUT-ST-OLDER-CONSUME | `MintSessionForRequest`'s version check deleted | ST5 v3-store |
| MUT-CL14-DROP | delete the `"git"` map entry | CL14 |
| MUT-CL14-NEWFIELD | add `extra string` to `mintRequest` | CL14 |
| MUT-CL14-CAPFIELD | add `Note string` to `broker.Capability` | CL14 (or compile error) |
| MUT-CL14-REQLEAF | add `Extra string` to `approval.MintRequest` | CL14 (and RQ1) |
| MUT-FR-PROFILESDIR | `fromRequest` omits `profilesDir` | CL15 |
| MUT-FR-CREDSDB | `fromRequest` omits `credsDB` | CL15 |
| MUT-FR-OUTDIRCREATE | `fromRequest` omits `outDirCreate` | CL15 |
| MUT-BM-REPOABS | `buildMintRequest` stores `r.repo` as typed | BM1 |
| MUT-BM-WSABS | `buildMintRequest` stores `r.wsRoot` as typed | BM1 |
| MUT-BM-EXPIRES | the grant-expiry refusal deleted | BM1 |
| MUT-CR-PROVISION | the `if *request { return createMintRequest(…) }` branch deleted | CL1 |
| MUT-CR-LAUNCH | the `--launch` refusal arm deleted | CL2 |
| MUT-CR-HEADLESS | the `--headless` refusal arm deleted | CL2 |
| MUT-CR-RESUME | the `--resume` refusal arm deleted | CL2 |
| MUT-CR-PRESET | the launch-preset refusal arm deleted | CL2 |
| MUT-CR-RANDOMNIL | the `env.random == nil` refusal deleted | CL1b (panic) |
| MUT-CR-EXPIRY | `ExpiresAt: now + r.ttl` | CL1 |
| MUT-CR-CHALLENGE | the challenge drawn before the id | CL1 |
| MUT-PE-BOUNDALL | `if r.answerBy != 0` (the read) → `if true` | CL17 |
| MUT-RTW-READONLY | `readTermLineWithin` → `return <-ch, true` | CL18 (guard) |
| MUT-AU-NOCHECK | `consume.Check` → `func(int64) error { return nil }` | AU1 |
| MUT-AU-PRECONSUME | `ConsumedAt: consumedAt` → `ConsumedAt: now` | AU1 `now = 1021` |
| MUT-FM-CONSUME | delete the `"consume"` map entry | CL14 |
| MUT-AP-ADMIT | `runApprove` skips `approval.Admit` | CL8 |
| MUT-AP-SAMESTORE | the same-store refusal deleted | CL12 (e) |
| MUT-AP-ADDR | `case "approve"`'s `addrGiven` refusal deleted | CL12 (c) |
| MUT-AP-GRAMMAR | the id-grammar refusal deleted | CL12 (a) |
| MUT-AP-AFTERNIL | the `env.after == nil` refusal deleted | CL12 (d) (panic) |
| MUT-PE-NOCONSUMEBRANCH | `provisionEpisode`'s `r.consume != nil` arm deleted | CL8 |
| MUT-AP-ADMITEXP | `RefuseExpired` treated as `Admitted` | CL9 arm 1 |
| MUT-AP-CLOCKREUSE | `Clock: env.now` → `Clock: func() int64 { return now }` (the admission reading) | CL9 arm 2, CL19 |
| MUT-AP-CLOCKEARLY | in `provisionEpisode` before step 5: `t0 := r.consume.Consume.Clock(); r.consume.Consume.Clock = func() int64 { return t0 }` | CL19 |
| MUT-AP-NODERIVE | the `!bytes.Equal(fb, row.Request)` refusal deleted | CL10 (a), (b), (e) |
| MUT-AP-MCPEXISTS | the MCP-config-exists refusal deleted | CL10 (d) |
| MUT-AP-DIGEST | `BytesDigest: sha256.Sum256(row.Request)` → `BytesDigest: [32]byte(stored)` | CL11 arm 1 |
| MUT-PE-PLAINREAD | the bounded branch → `answer = readTermLine(term)` always | CL13 arm 1 (guard) |
| MUT-PE-BOUNDFROMLIFETIME | `remaining := approval.RequestLifetime` | CL13 arm 1 (≠ 37 s) |
| MUT-PE-PROMPTCHECK | the "expired before the prompt" refusal deleted | CL13 arm 2 |
| MUT-AP-GRANTSDEFAULT | `fromRequest`'s grants → `sePresetGrants` | CL5 (refused as a stale world) |
| MUT-AP-LAUNCHMODEL | `printLaunch(…model: "")` in `runApprove` | CL5 |
| MUT-AP-TASK | `printLaunch(…task: "")` in `runApprove` | CL5 |
| MUT-EFF-SKIP-BRANCH | `approveEffectProblems` skips the branch check | CL16 |
| MUT-EFF-SKIP-MODE | it skips the session-mode check | CL16 |
| MUT-EFF-SKIP-APPROVEDBY | it skips the `approved_by` check | CL16 |

Known equivalent: making `readTermLineWithin`'s channel unbuffered only leaks the goroutine on a
timeout; nothing in a test can observe it (it is why the design asks for size 1). It is not
counted.

---

## 5. Files (final `git status --short` after M6d)

| File | Milestone | Change |
|---|---|---|
| `host/approval/sshsig_test.go` | M3b | transforms + SG8 |
| `host/approval/sshsig_crosscheck_test.go` | M3b | three SHA-256 arms |
| `host/store/credentials_schema.sql` | M4a | `mint_requests` |
| `host/store/credentials.go` | M4a, M4c | version 4, `credentialsTables`, `upgradeCredentials`, `insertSessionRow` |
| `host/store/mint_requests.go` | M4a–M4c | new |
| `host/store/mint_requests_test.go` | M4a–M4c | new: ST2, ST3, ST4, ST4b, ST5, ST6, ST7, ST8 |
| `host/store/credentials_test.go` | M4a–M4c | ST1 edit, ST2b, deadline edit |
| `host/store/context_roots_test.go` | M5b, M6b | +2 pins |
| `host/daemon/credentials.go` | M4a | log line |
| `host/authority/mint.go` | M6a | `RequestMint`, `MintForRequest`, `mintRecorded(…, rm)` |
| `host/authority/mint_request_test.go` | M6a | new: AU1 |
| `cmd/ailang-worldd/mint_request.go` | M5a, M5b, M6a | new |
| `cmd/ailang-worldd/mint_request_test.go` | M5a, M6a | new: CL14, CL15, BM1 |
| `cmd/ailang-worldd/session.go` | M5b, M6a | `random`, `after`, `readTermLineWithin` |
| `cmd/ailang-worldd/session_new.go` | M6a | `consume`, `answerBy`, bound, mint arm, comment |
| `cmd/ailang-worldd/episode.go` | M5b | `--request`, refusals, branch, help |
| `cmd/ailang-worldd/approve.go` | M6b | new |
| `cmd/ailang-worldd/approve_test.go` | M5b–M6d | new: CL1–CL13, CL16–CL19 |
| `cmd/ailang-worldd/main.go` | M5b, M6b | usage, `case "approve"`, `ownHelpVerbs` |
| `cmd/ailang-worldd/keyguard.go` | M6b | `ownVerbHelp` |
| `cmd/ailang-worldd/help_test.go` | M6b | `everyVerb`, want-list |
| `docs/QUICKSTART.md` | M6d | §14 appended |
| `design_docs/planned/w-world-approvals-by-key-a0a1-sprint2-plan.md`, `design_docs/planned/sprint_w-world-approvals-by-key-a0a1-s2.json` | planner | this plan |

---

## 6. Sprint gates (controller, after M6d; every one was green at the base, §1.1)

```
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh            # no .ail change: must stay exactly green
test -z "$(gofmt -l host/approval host/store host/authority host/daemon cmd/ailang-worldd)"
go build ./... && go vet ./host/approval/ ./host/store/ ./host/authority/ ./host/daemon/ ./cmd/ailang-worldd/
go test -run '^$' ./...                                                   # compile fence for every _test.go
go test ./host/approval/ -count=1 && go test -race ./host/approval/ -count=1
WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -run '^TestSSHSIGCrossCheckWithSSHKeygen$' -count=1 -v   # parent + 2 subtests PASS
go test ./host/store/ ./host/authority/ -count=1
go test -race ./host/store/ -run '^(TestConcurrentApprovalsMintOnce|TestMintSessionForRequest|TestPutMintRequest|TestCredentials)' -count=1   # ST6's -race leg
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1
go test ./host/verifygate/ -run '^TestNoParallelWriteForkPackagesOutsideVerifygate$' -count=1
go test ./cmd/ailang-worldd/ -count=1 -v -run '^(TestEpisodeStartRequest|TestApprove|TestMintRequestFieldMap|TestBuildMintRequest|TestExistingCallersReadUnbounded|TestReadTermLineWithin|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard|TestSessionMint_|TestSessionNew|TestEpisodeStart|TestStartLaunch|TestEpisodeResume|TestSessionMintGrantsExpireWithTheSession)' 2>&1 | tee /tmp/cl.log | tail -3
bash -c '! grep -E "^\s*--- (SKIP|FAIL)" /tmp/cl.log'                     # no silent skip (C17)
# AFTER the controller commits (the binary gate reads HEAD's committed blobs):
VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh   # go build + go test ./... + scripts/tracked_binary_gate.sh with its --allow-font pins
git status --short                                                         # only the §5 files
```

**Must stay green, unedited** (design §6): `TestSessionMint_*`, `TestSessionNew*`,
`TestSessionNewWithoutAControllingTerminalRefusesEndToEnd`,
`TestSessionNewReadsTheAnswerFromThePtyWhenStdinIsDevNull`, `TestEpisodeStart*` (other than the
new `TestEpisodeStartRequest*`), `TestStartLaunch*`, `TestEpisodeResume*`,
`TestSessionMintGrantsExpireWithTheSession`, every `TestCredentials*` except the three edited ones,
`TestMigrateCredentials*`, `TestPreSplitStoreOpenModes`, `TestEpisodeStartsUpgradeV2AndAtomicMint`,
`TestDroppedCredentialsTableIsRefused`, `TestDaemonDependencyAllowlist`,
`TestHelpIsNeverBlockedByTheRegistryKeyGuard`, and the QUICKSTART tests.

---

## 7. Risks

- **R1 — Size.** ~2,550 lines against the design's ~1,600 for M4–M6, over the design's
  2,000-line sprint ceiling by the M3b and the planner-added tests (§1.3 D13). The ten milestones
  keep each executor run under ~250 lines. If the controller must cut, M6d's CL16 and QUICKSTART
  §14 are the separable tail (CL5 then needs its own minimal effect asserts).
- **R2 — The bump moves three version gates** (§1.4 F8). Until `serve` of this build restarts, a
  v3 store refuses `session new --profile`, preset starts and resumes, naming the restart.
  `episode start` already forces that restart. Nothing in the suite builds a v3 store outside
  `host/store` (C9); the full `go test ./...` in `verify_go.sh` is the check.
- **R3 — Silent skips.** Every approve test rides `matchingHealth`, which `t.Skipf`s on a
  platform with no tool pin (C17). The skip guard in §3/§6 makes a skip red here and on CI.
- **R4 — The blocked reader in CL13/CL18.** The goroutine stays blocked on the pipe until
  `t.Cleanup` closes the writer. Forgetting the cleanup leaks a goroutine per run and could trip a
  goroutine-leak check; none exists in this package today, but keep the cleanup.
- **R5 — Real git in many tests.** CL5–CL19 add ~15 worktree adds to a package that already takes
  ~187 s; budget ~20 s more. No new CI timeout risk measured.
- **R6 — The older-daemon refusal** (design R3) is unchanged: a daemon of an older build refuses a
  v4 store, as it refused v2 and v3.
- **R7 — QUICKSTART §14 is attended and unverified** (S7). It is marked *pending first verbatim
  run*; Mark's first run is the check (A1b's attended episode, D-WORLD-109).
- **R8 — Haiku and reflection.** CL14's walker is the most intricate code in the sprint. §3 M5a
  gives its exact algorithm and two positive controls; the controller should read it, not only run
  it.
- **R9 — `approve` reaches the daemon over loopback** (`fetchDaemonHealth`). Tests always stub it;
  an unstubbed test would hit `127.0.0.1:7644` and fail closed (unreachable), never mint.

---

## 8. Acceptance criteria (the judge scores against these)

1. M3b: the SHA-256 hash path is exercised by XC1 with real `ssh-keygen -O hashalg=sha256`
   signatures on both key types and the PASS-grepped subtests; SG8 kills the signature-string and
   mpint-pair trailing-byte branches, both `Strict()` decodes and the parse-level `ErrHashAlg`. No
   fixture, `.gitattributes` or `ci.yml` change; no hand-typed signature bytes.
2. The credentials store is schema 4 with exactly three tables; `mint_requests` matches §2.1 in a
   fresh store and in stores upgraded from 1, 2 and 3, every existing row unchanged; a session
   verb never upgrades; a v4 store without `mint_requests` is refused unmodified.
3. `PutMintRequest`, `MintRequest`, `PendingMintRequests` and `MintSessionForRequest` exist with
   §2.3's signatures, each requires a deadline and refuses a store below 4 with
   `*OlderCredentialsStoreError`; `PutMintRequest` enforces the id grammar the `CHECK` does not, the
   16-pending cap and the one-day prune (consumed rows never pruned).
4. `MintSessionForRequest` reads `Clock()` exactly once after `BeginTx`, runs `Check` with it before
   the guarded `UPDATE`, inserts the session row only when exactly one request row changed, and
   classifies every refusal with a typed error, leaving the request pending and no session row.
5. `authority.MintForRequest` runs `approval.MintMatches` inside that transaction with the
   in-transaction reading as `created_at`'s upper bound (sprint-1 R5 / §1.4 F3).
6. `episode start --request` opens no terminal, creates no worktree, session or MCP file, mints
   nothing, stores one request (id `mr-`+16 hex, 32-byte challenge, `expires_at = created_at +
   600`), prints it and `ailang-worldd approve <id>`, and refuses `--launch`, `--headless`,
   `--resume` and launch presets.
7. `approve` refuses, before any terminal opens: a bad id, an unknown id, a used, expired or
   tampered request, another store, a changed world (naming the field), the global `--addr`. It
   mints only through `provisionEpisode`; `--list` shows only pending requests.
8. The y/N of `approve` is bounded by the request's remaining lifetime and leaves the request
   pending on timeout; every existing caller's y/N is unchanged (`after` never called).
9. An approve whose request expires during worktree preparation is refused at the write, the
   request stays pending and the new worktree is removed (CL19).
10. `mintRequestFieldMap` covers every field of `mintOptions`, `mintRequest` and
    `broker.Capability`, and every `MintRequest` leaf; `fromRequest` rebuilds `episode start`'s
    `mintRequest` exactly (CL14, CL15).
11. CL5 shows the minted row matches the request on all seven columns and every other effect
    (worktree, branch, token file and mode, MCP config, launch line, `approved_by`, credential id);
    CL16 shows the effect check is not vacuous for any of its twelve fields.
12. Exactly two new production context roots (§0.8); `TestProductionContextRoots`,
    `TestDaemonDependencyAllowlist`, `TestEveryVerbHasHelp` and
    `TestHelpIsNeverBlockedByTheRegistryKeyGuard` pass; `approve --help` prints `approveHelp`.
13. QUICKSTART §14 exists at the end of the file, marked attended and pending its first verbatim
    run; nothing else in the file changed.
14. Each §4 mutant, applied alone, makes its named test(s) fail.
15. Every §6 gate passes, including `verify_go.sh` (and its tracked-binary gate) after the commit;
    no test skips; `git status --short` lists only the §5 files.
16. No real key, real store, mint or publish anywhere in the sprint (D-WORLD-109).
