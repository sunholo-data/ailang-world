# Sprint plan — w-world-approvals-by-key-a1b, sprint 1 (iteration 375, row 201)

> **THE §8 FALLBACK CUT FIRED. This sprint builds M1 and M2 only. M3 (`approver
> enroll|list|revoke`) moves to the head of sprint 2, and the design's M6 becomes a third sprint.**
> The planner's per-file estimate for M1–M3 is **≈2,750 changed lines** (M1 ≈700, M2 ≈840, M3
> ≈1,210), over the design's 2,000-line threshold (§1.5). Sprint 1 is therefore **≈1,540 lines in
> six sub-milestones**: the `CR-TOOMANY` test, the test signer, the pure enrolment law, and
> credentials schema 5 with its five store methods. It adds **no CLI verb and no user-facing
> surface**. M3's research is not lost: Appendix A carries its decisions, two design premises found
> false that only M3 (and sprint 2's F9) touch, its sub-milestones and its mutant list, for the
> sprint-2 planner.

- **Design**: [w-world-approvals-by-key-a1b.md](w-world-approvals-by-key-a1b.md), **r2** (commit
  `5e2bf772`). The design is **not** reopened. Where this plan and the design differ, this plan was
  measured against the tree and wins. The differences are in §1.3 (decisions) and §1.4 (premises
  found false).
- **Format model**: [w-world-approvals-by-key-a0a1-sprint2-plan.md](w-world-approvals-by-key-a0a1-sprint2-plan.md)
  (iteration 374; landed as PR #535, merge `ca8a4c17`).
- **Planner base**: `5e2bf772` (= `3a96cf75` + the three design commits; no code changed since
  `3a96cf75`), worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter375-row201-a1b`,
  branch `sprint/world-iter375-row201-a1b`.
- **Executor**: as routed by the controller (iteration 374: `claude-haiku-5-5`), **two
  sub-milestones per run, three runs**: (M1a, M1b), (M1c, M2a), (M2b, M2c). Every name, signature,
  constant, SQL statement, error string and assertion below is fixed. **Copy, do not paraphrase.**
  §2's production code was compiled and smoke-tested by the planner in a scratch export (§1.2 C14–C18).
- **Estimate**: ≈**1,540 changed lines** (insertions + deletions, the V24 measure): ≈440
  production (Go, SQL, and the non-test `approvaltest` package), ≈1,075 Go test, ≈26 fixture/attribute
  lines. ≈1.25 executor days.
- **Touched**: `host/approval/approver.go` (new), `host/approval/approver_test.go` (new),
  `host/approval/approvaltest/` (new package: `signer.go`, `signer_test.go`),
  `host/approval/testdata/approvers/` (new: `generate.sh`, `rsa3072.pub`, `p384.pub`, `SOURCE.txt`),
  `.gitattributes` (+1 line), `host/store/{credentials_schema.sql,credentials.go,mint_requests.go}`,
  `host/store/approvers.go` (new), `host/store/{approvers_test.go (new),credentials_test.go,mint_requests_test.go}`,
  `cmd/ailang-worldd/approve_test.go` (+CT1 only).
  **Not touched**: any other `cmd/ailang-worldd` file (no verb, no help, no context root),
  `host/store/context_roots_test.go`, `host/authority/**`, `host/daemon/**`, `host/approval/{sshsig,admit,request}.go`,
  `host/approval/testdata/sshsig/**` (SG1's pins do not move), `.github/workflows/ci.yml`,
  `docs/QUICKSTART.md`, `world/*.ail`, `packages/**`, `scripts/**`, `tools/launchd/*`, any skill,
  `~/.ssh`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated in the milestone it governs)

1. **Fakes, temp stores and test keys only (D-WORLD-109 = A).** Every store is under `t.TempDir()`.
   Keys come from the checked-in `.pub` fixtures, from `approvaltest.NewSigner(seed)` (a fixed
   seed, no randomness), or from `host/approval/testdata/approvers/generate.sh` run **once** by the
   executor in M1b. Nothing opens `~/.ailang/world`, `~/.ssh`, a real credentials store, a real
   daemon or `/dev/tty`.
2. **`ssh-keygen` discipline.** The only `ssh-keygen` key generation in this sprint is
   `generate.sh`, which runs in `mktemp -d` with `HOME` set to it, `SSH_AUTH_SOCK` unset, `-N ''`,
   `-f "$tmp/…"` and `< /dev/null`, and deletes the directory on exit. **Never run a bare
   `ssh-keygen`** (with no `-f` it offers to write `~/.ssh/id_*`; iteration 374's planner hit that
   prompt).
3. **No committed binary.** Every new fixture is a one-line or few-line text file with no NUL
   byte (the tracked-binary gate, `scripts/tracked_binary_gate.sh`, run by `verify_go.sh`).
   `.gitattributes` gains `host/approval/testdata/approvers/* -text`, as `sshsig/` has.
4. **`host/approval` stays pure.** `approver.go` imports only `bytes`, `errors`, `fmt`, `regexp`.
   `TestApprovalImportsOnlyTheStandardLibrary` (it reads only `host/approval/*.go`, not
   subdirectories) must stay green. `host/store` does **not** import `host/approval` (production or
   test).
5. **`approvaltest` is test support only.** It lives under `host/approval/approvaltest/`, has no
   `_test.go`-only build tag (Go cannot import those), and **no non-test file in the repo may import
   it** (TS3).
6. **No behaviour change for existing callers.** Mint requests, `approve`, `episode start
   --request`, sessions: unchanged. No existing error text changes: `OlderCredentialsStoreError`
   for mint requests prints byte-identical text (AS6).
7. **No `t.Parallel()`** in any new test (`host/store` and `cmd/ailang-worldd` write files and fork).
8. **Every new production branch has one mutant** (§4, 74 mutants). A refusal arm asserts the typed
   error with `errors.Is`/`errors.As` **and** that nothing was written.
9. **Exit codes by constant** (`exitUsage`, `exitFatal`), never by literal (design §1.1).
10. **Bisectable.** Each sub-milestone ends compiling — `go build ./...` **and** `go vet <pkg>`,
    because `go build` does not compile `_test.go` files — with its own tests and every earlier
    sub-milestone's tests green.
11. **No new context root.** This sprint adds no `context.Background()` in production code
    (`TestProductionContextRoots` stays green unedited). M3's three roots move with M3.
12. **No git writes by the executor** (§7). The controller rebuilds commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `5e2bf772` (measured by the planner, 2026-11-10)

Every command ran in the worktree with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` in the same
shell (`AILANG v0.41.0`, commit `24ee108`), `go1.26.6 darwin/arm64`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git status --short` | 0 | — | clean |
| `./scripts/verify_ail.sh` | **0** | 52 s | `verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `./scripts/verify_go.sh` (default legs = all), **plain leg** | — | ≈8 min | tracked-binary gate `24 passed, 0 failed`, `3 pinned font(s) validated`, `0 other binary blobs among 6593 tracked files`; `all 64 required top-level evidence tests passed exactly once`; `go test ./... -count=1`: **every package `ok`, 0 FAIL** (`cmd/ailang-worldd` 229 s, `host/daemon` 377 s, `host/broker` 344 s, `host/store` 23.6 s, `host/approval` 2.4 s). That leg **is** the `go test ./...` baseline (same command, `-count=1`) |
| `./scripts/verify_go.sh`, race leg | — | — | the first run's race leg was killed by the planner shell's 600 s background limit (an instrument limit); re-run as `VERIFY_GO_LEGS=race ./scripts/verify_go.sh`: **rc 1, three load timeouts — see §1.1a** |
| `go build ./... && go vet ./...` | **0** | 2 s | clean |
| `go test -run '^$' ./...` (compile fence for every `_test.go`) | **0** | 10 s | no failure line |
| `go test ./host/store/ -run '^TestProductionContextRoots$' -count=1` | **0** | <1 s | `ok … 0.249s` |
| `go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1` | **0** | 1 s | `ok … 0.364s` |
| `gh pr list --state open` | — | — | not re-run; design V25 measured `[]` at the same base |

### 1.1a Race leg (re-run)

**The local race leg is RED at the base, from load, not from code.** Two runs:

| Command | rc | Time | Result |
|---|---|---|---|
| `VERIFY_GO_LEGS=race ./scripts/verify_go.sh` | **1** | 511 s | the race-detector known-positive control fires as designed (`Found 2 data race(s)`, exit 66); `go test ./... -count=1 -race -timeout 8m`: 57 packages `ok` (incl. `cmd/ailang-worldd` 338 s, `host/store` 87 s, `host/approval` 3.9 s); **three `-timeout 8m` package timeouts**: `host/broker` 481.9 s (`TestPublishValidationBudgetOnly/validation_expires` running), `host/daemon` 482.6 s (`TestWorkbenchTimelinePaging`), `host/verifygate` 480.5 s |
| `go test -race -count=1 -timeout 8m ./host/broker/ ./host/daemon/ ./host/verifygate/` (alone) | **1** | 482 s | `host/broker` ok 459.9 s; `host/daemon` timed out at 8 m (`TestWorkspacePackageCacheCoversLock`); `host/verifygate` FAIL: `TestMissionConfigDefaultFlowIntegration/valid_ledger: default-flow fixture timed out / instrument failure: context deadline exceeded` |

`uptime` during the second run: `load averages: 12.81 11.81 11.07` on 16 cores (other sessions on
this machine). Every red is a deadline or a package timeout, none is a race report or an
assertion. The same code (`3a96cf75`) is green on CI, including its separate `go-race-verify` job
(run `42753484556`, conclusion `success`, design V15). **Rule for the controller:** compare the
race leg with this baseline. A timeout-only red in `host/broker`, `host/daemon` or `host/verifygate`
is the pre-existing local instrument limit, not a regression; this sprint touches none of those
packages. The sprint's own packages are race-checked directly in §6 (`go test -race
./host/approval/... ./host/store/`), and the authoritative race leg is CI's `go-race-verify` job
(`VERIFY_GO_LEGS: race`, `ci.yml:517-523`).

### 1.2 Codebase claims checked

| # | Claim | Evidence (command → output) |
|---|---|---|
| C1 | The eight files the design measured are byte-identical at the base, so its line numbers hold | `shasum -a 256 cmd/ailang-worldd/{approve,session_new,mint_request}.go host/store/{mint_requests,credentials}.go host/authority/mint.go host/approval/{sshsig,admit}.go \| cut -c1-16` → `a215b26ad7ed060a d12b…` — all eight equal design V1 |
| C2 | Credentials schema 4; table sets; the upgrade's last step | `grep -n` → `credentials.go:42 const credentialsSchemaVersion = 4`; `:71 func credentialsTables`, `:73 case version >= 4:`; `:353 if version < 4 {`; `credentials_schema.sql` is 56 lines, ending with the `mint_requests` table |
| C3 | `OlderCredentialsStoreError` and its four sites | `mint_requests.go:59-69` (`Path`, `Version`; `Error()` hard-codes "mint requests need schema %d" with `mintRequestsSchemaVersion`); sites `:93`, `:147`, `:167`, `:254`; no other constructor in `host`/`cmd` (`grep -rn OlderCredentialsStoreError`) |
| C4 | The tests pinned to "need schema 4" check by `Contains` | `mint_requests_test.go:202` `strings.Contains(err.Error(), "need schema 4")`; `cmd/ailang-worldd/approve_test.go` (TestEpisodeStartRequestRefusesAnOlderStore) wants `"mint requests need schema 4"` |
| C5 | Exactly two existing tests pin the three-table set | `grep -rn '"mint_requests"' --include=*_test.go host cmd` → `credentials_test.go:56` and `mint_requests_test.go:132` only |
| C6 | No cmd/daemon/authority test builds a v4 **credentials** store | `grep -rn "user_version = [0-9]" --include=*_test.go cmd host/daemon host/authority` → `presplit_unix_test.go:38` and `credentials_split_test.go:44` (both a pre-split **world** store with `session_credentials` only) and `approve_test.go:370` (`= 3`) |
| C7 | `applicationTables` and `schemaState` skip `sqlite_*` objects, so `UNIQUE`'s autoindex is not a table | `credentials.go` `applicationTables`: `type = 'table' AND name NOT LIKE 'sqlite\_%'`; `schema_version_test.go:179 schemaState` the same filter over every object |
| C8 | Store test helpers exist | `rawDB` (`schema_version_test.go:169`), `schemaState` (`:179`), `setPositiveFixtureVersion` (`:297`), `tableDDL` (`journal_test.go:968`), `normalizeDDL` (`:854`), `requireExactTableNames` (`:867`), `makeV1Credentials` (`credentials_binding_test.go:16`), `makeV3Credentials` (`mint_requests_test.go:30`), `rigRows` (`credentials_migrate_test.go:14`), `oneBoundRow` (`mint_requests_test.go:55`), `testCtx` (`:175`), `okConsume` (`:352`), `consumeSession` (`:363`), `openTestCredentials` (`credentials_test.go:87`) |
| C9 | `host/approval`'s import guard reads only its own directory | `imports_test.go:28 os.ReadDir(".")`; forbidden: `os`, `os/exec`, `os/signal`, `net`, `net/http`, `syscall`, `io/ioutil`, `crypto/rand`, `math/rand`, `math/rand/v2`, `unsafe` |
| C10 | SG1 pins only `testdata/sshsig/` | `sshsig_test.go:44 fixtureSHA256` (8 files), `:361 len(pinned) != 8` over `testdata/sshsig`; a new `testdata/approvers/` does not move it |
| C11 | Fixture fingerprints | `testdata/sshsig/SOURCE.txt`: ed25519 `SHA256:pHMYJm+CndqyZGrzMHQQBH0tNeNGj1U9sy0JrVsYdSQ`, p256 `SHA256:yMwGgMZdq1aVNYRG/HO2T69s/Fb14yPZW0GBWtwRyEo`; `Fingerprint()` uses `RawStdEncoding` (same form) |
| C12 | No name collisions | `grep -rn "func Test\(NewApprover\|CheckRevoke\|ApproverName\|RenderEnrolment\|EnrolmentPrompt\|ApproverFixtures\|TestSigner\|ApprovalTest\|EnrolApprover\|ApproverLookups\|RevokeApprover\|CredentialsUpgradeToV5\|OlderCredentialsStoreErrorText\|ApproverRow\|CredentialsV5\|EpisodeStartRequestRefusesWhenSixteen\|ListApprovers\)" host cmd` → nothing; control `grep -c "func TestCredentials" host/store/*_test.go` totals 8. No production identifier `Approver`, `FormatSSH`, `ErrApprover*`, `ApproverRow`, `approversSQL` exists |
| C13 | No `t.Parallel()` in the touched packages | `grep -rn "t.Parallel()" host/approval host/store cmd/ailang-worldd --include=*_test.go \| wc -l` → `0` |
| C14 | §2.1 and §2.2 compile and behave | scratch export (`git archive HEAD \| tar -x -C /tmp/w375-plan-scratch`, outside the repo): `gofmt -l` empty, `go build`/`go vet ./host/approval/...` clean; a smoke test: `NewApprover("mark-mac","ssh",ed25519.pub,…)` → fingerprint `SHA256:pHMY…YdSQ`; `RenderEnrolment` printed §2.1's layout; the signer's key parses and its signature verifies; a one-byte-different message → refused |
| C15 | `approvaltest.Signer` is ssh-keygen compatible (planner measurement; XS2 pins it in sprint 2) | in `T=$(mktemp -d)` with `HOME=$T`, `SSH_AUTH_SOCK` unset: the signer's `msg.sig` and `k.pub` → `ssh-keygen -Y check-novalidate -n world-approval@sunholo -s msg.sig < msg` → `Good "world-approval@sunholo" signature with ED25519 key SHA256:z/fSv0Z0RZS+Lccbc6ZoOWwt/fbj1VFJGSBKQbo4icE`, rc 0; `ssh-keygen -l -f k.pub` → the same fingerprint; a tampered message → `incorrect signature`, rc 255; `$T` removed. Signature length 318 bytes, as ssh-keygen's |
| C16 | §2.4–§2.7 compile, and the bump reds exactly the two pinned tests | scratch: `go build ./... && go vet ./host/store/` clean; `go test ./host/store/ ./host/approval/... ./host/authority/ -count=1` → **only** `TestCredentialsSchemaMatchesCanonicalDDL` and `TestCredentialsUpgradeToV4/{v1,v2,v3}` FAIL; `host/approval`, `host/authority` ok. Then, with the bump and CT1 in place, `go test ./cmd/ailang-worldd/ ./host/daemon/ ./host/authority/ ./host/store/ ./host/approval/... -count=1` (302 s) → `cmd/ailang-worldd` ok 144.8 s, `host/daemon` ok 298.4 s, `host/authority` ok, `host/approval` ok, `host/store` the same two reds only. **The schema bump breaks nothing outside the two pinned tests** |
| C17 | The store's refusals behave as §3 asserts | scratch smoke: same name → `ErrApproverNameTaken` (typed); same key → `ErrApproverKeyEnrolled … it is approver "mark-mac"'s key` (typed); a raw `INSERT` of a duplicate key → `constraint failed: UNIQUE constraint failed: approvers.public_key (2067)`, **`errors.Is(…, ErrApproverKeyEnrolled) == false`** (the backstop is untyped, so AS2's pre-check mutants red); `ApproverByKey` of a 2-byte prefix of a 3-byte key → `ok == false`; revoke at 999 (< enrolled 1000) → `ErrApproverClockBehind`; at 1000 → nil; again → `ErrApproverRevoked`; unknown → `ErrApproverUnknown`; v4 handle → `… is schema 4; approvers need schema 5: restart …`; v4 `PutMintRequest` still works; `UpgradeCredentials` v4 → `true`, tables `[approvers episode_starts mint_requests session_credentials]` |
| C18 | CT1 is green on the base and kills the `CR-TOOMANY` mutant | scratch: CT1 (§3 M1a) → `--- PASS … (0.33s)`; with `mint_request.go:142` mutated to `errors.As(err, &older) \|\| false && errors.Is(err, store.ErrTooManyPendingRequests)` → `--- FAIL … rc 2, want exitUsage`; restored (`grep -c "false && errors.Is"` → `0`). Agrees with design V10/V11 |
| C20 | Without `ORDER BY`, SQLite returns `approvers` rows in insertion (rowid) order, so MUT-LS-ORDER is observable | scratch: enrol `b`, `a`, `c`; `SELECT … FROM approvers LIMIT 256` → `[b a c]` |
| C21 | TS3's walk has a real population | `find . \( -name '.*' -a ! -name '.' -o -name node_modules -o -name vendor \) -prune -o -name '*.go' -print \| wc -l` → `660` (TS3's floor is 300) |
| C19 | `cmd/ailang-worldd` already opens files with `O_NONBLOCK` and fstat-checks them | `episode_launch_interactive.go:252 os.OpenFile(source.path, os.O_RDONLY\|syscall.O_NOFOLLOW\|syscall.O_NONBLOCK, 0)`, then `f.Stat()`, `IsRegular()`, `io.LimitReader(f, capBytes+1)` (the precedent Appendix A uses) |

### 1.3 Planner decisions where the design leaves a choice

- **D1 — The cut** (§1.5). M3 moves to sprint 2. The design's §8 text names this exact cut.
- **D2 — Shared store helpers.** `Approver` and `ApproverByKey` share one private
  `approverWhere(ctx, where, arg)`; every below-v5 refusal is built by one
  `(*Credentials).olderForApprovers()`. Each still requires a deadline and refuses below v5 (the
  design's table), with one mutant per shared check (§4).
- **D3 — `RenderEnrolment`'s exact layout** (§2.1) and the R6 text as an exported constant
  `EnrolmentWarning`, so A3 can show the same words. AP5's golden spells the text out literally (it
  does **not** reference the constant), so an edit to the constant reds.
- **D4 — `RenderEnrolment` prints `approved_by=<format>:<name>`**, not a hard-coded `ssh:`, so the
  prompt cannot drift from sprint 2's `SignedApproval.ApprovedBy()` (`format + ":" + name`, design
  §2.3). For every A1b approver the output is `approved_by=ssh:<name>`, as the design asks.
- **D5 — RC4 classifies `Key.Blob` as shown (as its fingerprint)** and checks it by perturbation —
  stronger than the design's "excluded: shown as its fingerprint", same intent. The three
  timestamps are excluded, and RC4 checks that perturbing them does **not** change the text.
- **D6 — `approvaltest` API**: `NewSigner(seed byte) *Signer` (32-byte seed = `seed` repeated),
  `PublicKey() approval.SSHPublicKey`, `AuthorizedKey() []byte` (`"ssh-ed25519 <b64> approvaltest\n"`),
  `Sign(namespace string, msg []byte) []byte` (armored SSHSIG, sha512, 70-column lines). Its test
  file is the **external** package `approvaltest_test`, so TS3's control ("at least one `_test.go`
  imports it") holds in this sprint (§1.4 P6).
- **D7 — TS3's walk** starts at the repo root (the first ancestor holding `go.mod`) and skips every
  directory whose name starts with `.` (`.git`, `.claude`, `.wt-*`) plus `node_modules` and
  `vendor`. It parses imports only (`parser.ImportsOnly`).
- **D8 — AS6 compares the full text at all four mint-request sites** (not a `Contains`), because
  the existing pins (C4) would not notice a site that drops `Feature` ("…; need schema 4" still
  contains "need schema 4"). See §1.4 P4.
- **D9 — Planner-added tests**: **AS7** `TestCredentialsV5WithoutApproversIsRefused` (the v5 table
  set is enforced; the twin of the existing ST2b), **AS3b** `TestListApprovers` (split out of AS3 so
  M2b and M2c each stay under ~350 lines), and extra arms named in §3 (the 8192-byte accept
  boundary, leading-whitespace private key, revoked-and-behind order, the cross-approver revoke).
- **D10 — `ErrApproverKeyChanged` is declared now** (design §2.4's var block) and first used in
  sprint 2 (C11). An unused package-level `var` compiles and `go vet` does not flag it.
- **D11 — `generate.sh` is executable** (`0755`, like `testdata/sshsig/generate.sh`, git mode
  `100755`). Snapshots use `cp -p` so the controller's commit keeps the mode (§7).

### 1.4 Premises in the design found false or incomplete (each routed around minimally)

- **P1 — "Open, then `fstat` the descriptor" does not refuse a FIFO with a plain `os.Open`**
  (design §2.7 F9, §2.10 N4, CA5, CR4; **affects M3 and sprint 2 only**). `open(2)` of a FIFO with
  `O_RDONLY` blocks until a writer appears, **before** any `fstat` can run, and a context root does
  not bound `os.Open` or `Read`. Measured in a throwaway module in `mktemp -d` (deleted): `mkfifo f`;
  `os.Open(f)` → `BLOCKED >2s in open`; `os.OpenFile(f, os.O_RDONLY|syscall.O_NONBLOCK, 0)` →
  `opened; regular=false mode=prw-r--r--`. **Route** (Appendix A, A-D1): open with
  `O_RDONLY|syscall.O_NONBLOCK` (the C19 precedent), then `f.Stat()`, `IsRegular()`, size, and one
  `io.LimitReader(f, max+1)` read. The design's mutant "`f.Stat()` replaced by a path `Lstat` (the
  FIFO arm then blocks)" is **not killed** by a FIFO arm with no concurrent swap (a path `Lstat`
  also refuses a FIFO); Appendix A replaces it with "`O_NONBLOCK` dropped" (the arm blocks; the 5 s
  guard fires) and "`IsRegular` check dropped" (the FIFO reads as 0 bytes and the arm's
  `not a regular file` text check fires).
- **P2 — §2.10a's code sketch contradicts its own [r2] reader contract** (M3 only). The sketch of
  `confirmApproverAct` still calls `readTermLineWithin` and comments "(also EOF and read failure:
  line == "")", which V39 showed false. **Route** (A-D2): follow the [r2] text —
  `readTermLineErrWithin(r, d, after) (line string, err error, ok bool)` returns the **raw,
  untrimmed** line (because `readTermLine`, `session.go:477-480`, trims and drops the error), and
  `confirmApproverAct` confirms only when `ok && err == nil && strings.HasSuffix(line, "\n") &&
  isYes(line)`.
- **P3 — `credentials_binding_test.go:106` needs no edit.** The design's §6 M2 lists it as edited.
  It already compares with the symbolic `credentialsSchemaVersion`; in the scratch bump (C16) it
  stays green unedited. The only edits are `credentials_test.go:56` and `mint_requests_test.go:132`.
- **P4 — AS6 as written would not kill a site that drops `Feature`** (incomplete, not false).
  Routed by D8.
- **P5 — The sprint-1 size.** Design: "about 2,050" for M1–M3. Planner per-file estimate: ≈2,750.
  Routed by D1 (§1.5).
- **P6 — TS3's control in sprint 1.** "At least one `_test.go` file [imports approvaltest]" is true
  only through M3's CLI tests in the design's plan; with M3 cut, none would. Routed by D6 (the
  signer's own test is an external test package that imports it).
- **P7 — `isYes` accepts `yes` as well as `y`** (`session.go`). The design's reader contract says
  "the answer is y". M3 (Appendix A) keeps `isYes`, as every other attended verb does; the
  difference is named, not changed.

### 1.5 The estimate, and why the cut fires

Lines are insertions + deletions (`git diff --numstat`), the measure design V24 used for the
1.6× and 2.2× overruns.

| Sub-milestone | Production | Test | Fixture | Total |
|---|---|---|---|---|
| M1a — CT1, `approvaltest` (TS1, TS3) | 75 (test-support package) | 150 | — | **225** |
| M1b — `NewApprover`, `CheckRevoke`, fixtures (AP0–AP4) | 80 | 230 | 26 | **336** |
| M1c — `RenderEnrolment` (AP5, RC4) | 15 | 125 | — | **140** |
| M2a — schema 5, `OlderCredentialsStoreError` (AS1, AS5, AS6, AS7) | 89 | 209 | — | **298** |
| M2b — `EnrolApprover`, `Approver`, `ApproverByKey` (AS2, AS3, RC2) | 110 | 229 | — | **339** |
| M2c — `ListApprovers`, `RevokeApprover` (AS3b, AS4) | 70 | 133 | — | **203** |
| **Sprint 1 (this plan)** | **439** | **1,076** | **26** | **≈1,541** |
| M3 (Appendix A: four sub-milestones) | ≈395 | ≈815 | — | ≈1,210 |
| **M1–M3** | | | | **≈2,750 > 2,000 → cut** |

The design's own figures were M1 ≈520 + fixtures, M2 ≈650, M3 ≈880. The difference is the tests
the design's rows imply but do not size (AP0's fixture pins, the RC4 walker, AS5's three source
versions, AS6's four sites, AS7, the 257-row cap arm) and M3's 10 CLI tests with their timer,
FIFO and partial-read rigs.

---

## 2. Shared definitions (the executor copies these exactly)

### 2.1 `host/approval/approver.go` (new; M1b writes all but the last two declarations, M1c adds `EnrolmentWarning` and `RenderEnrolment`)

```go
package approval

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.2): the pure enrolment
// law. An approver is a name and one SSH public key, enrolled at the daemon's terminal (a typed
// y); its signature may later approve a mint request (sprint 2). No I/O here (coding standard S2).

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
)

const (
	// FormatSSH is the only approver format A1b enrols.
	FormatSSH = "ssh"
	// FormatWebAuthn is named so it can be refused with a message naming row 201 A2.
	FormatWebAuthn = "webauthn"
	// MaxKeyFileBytes caps a --pubkey file (one authorized_keys line).
	MaxKeyFileBytes = 8 << 10
)

// ApproverNamePattern is an approver's name: lowercase, digits, '.', '_' and '-', 1 to 48
// characters, starting with a letter or digit. It is printed on prompts and stored as ssh:<name>.
var ApproverNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]{0,47}$`)

// Approver is one enrolled approver as the pure law sees it.
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
// The checks run in a fixed order (E1-E7); the first that fails is returned.
func NewApprover(name, format string, keyFile []byte, now int64) (Approver, error) {
	if !ApproverNamePattern.MatchString(name) { // E1
		return Approver{}, ErrApproverName
	}
	if format == FormatWebAuthn { // E2 (before E3, so webauthn names A2)
		return Approver{}, ErrFormatNotYet
	}
	if format != FormatSSH { // E3
		return Approver{}, ErrApproverFormat
	}
	if len(keyFile) == 0 || len(keyFile) > MaxKeyFileBytes { // E4
		return Approver{}, ErrApproverKeyFile
	}
	if bytes.HasPrefix(bytes.TrimLeft(keyFile, " \t\r\n"), []byte("-----BEGIN")) { // E5: never echoed
		return Approver{}, ErrPrivateKeyFile
	}
	line := bytes.TrimSuffix(keyFile, []byte("\n"))
	if bytes.ContainsAny(line, "\n") { // E6
		return Approver{}, ErrApproverKeyFile
	}
	key, err := ParseAuthorizedKey(line)
	if err != nil { // E7: ErrKeyType (sk-*, rsa, p384), ErrAuthorizedKey, ErrBase64
		return Approver{}, err
	}
	return Approver{Name: name, Format: FormatSSH, Key: key, EnrolledAt: now}, nil
}

// CheckRevoke is the revocation law, run before the y/N (the store re-checks it in SQL).
func CheckRevoke(a Approver, now int64) error {
	if a.RevokedAt != 0 { // R1
		return ErrApproverRevoked
	}
	if now < a.EnrolledAt { // R2
		return ErrApproverClockBehind
	}
	return nil
}

// ---- M1c ----

// EnrolmentWarning is the R6 text every enrolment y/N shows, verbatim (design §2.2).
const EnrolmentWarning = "This key will approve mint requests with no terminal. A signature proves that something holding\n" +
	"the private key signed. It does not prove that you were there. Enrol only a key that cannot sign\n" +
	"without you: one behind a passphrase, in an agent that asks before each use (`ssh-add -c`), or in\n" +
	"a Secure Enclave or hardware holder (for example Secretive). Compare the fingerprint above with\n" +
	"`ssh-keygen -l -f <your key>.pub` or `ssh-add -l` before you type y.\n"

// RenderEnrolment is the body of the enrolment y/N: what will be stored, and what it proves. It
// has no branches (AP5 pins it; RC4 checks it shows every persisted field).
func RenderEnrolment(a Approver) string {
	return fmt.Sprintf("approver  %s\nformat    %s\nkey       %s %s\nrecords   approved_by=%s:%s\n\n%s",
		a.Name, a.Format, a.Key.Type, a.Key.Fingerprint(), a.Format, a.Name, EnrolmentWarning)
}
```

In M1b the file has no `fmt` use yet: the M1b import block is `bytes`, `errors`, `regexp`; M1c adds
`"fmt"`. (Leave the `// ---- M1c ----` marker in the file; it is a comment.)

`RenderEnrolment` for the `ed25519.pub` fixture as `mark-mac` (C14, measured) is exactly:

```
approver  mark-mac
format    ssh
key       ssh-ed25519 SHA256:pHMYJm+CndqyZGrzMHQQBH0tNeNGj1U9sy0JrVsYdSQ
records   approved_by=ssh:mark-mac

This key will approve mint requests with no terminal. A signature proves that something holding
the private key signed. It does not prove that you were there. Enrol only a key that cannot sign
without you: one behind a passphrase, in an agent that asks before each use (`ssh-add -c`), or in
a Secure Enclave or hardware holder (for example Secretive). Compare the fingerprint above with
`ssh-keygen -l -f <your key>.pub` or `ssh-add -l` before you type y.
```

(ending with one `\n` after `type y.`).

### 2.2 `host/approval/approvaltest/signer.go` (new; M1a)

```go
// Package approvaltest is TEST SUPPORT for host/approval: an SSHSIG signer over a fixed Ed25519
// seed, so tests make keys and signatures with no randomness, no ssh-keygen and no ~/.ssh. Only
// _test.go files may import it (TestApprovalTestIsTestOnly); production code never does.
package approvaltest

import (
	"bytes"
	"crypto/ed25519"
	"crypto/sha512"
	"encoding/base64"
	"encoding/binary"
	"strings"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// Signer is an Ed25519 key that signs as `ssh-keygen -Y sign` does (SSHSIG, sha512).
type Signer struct{ priv ed25519.PrivateKey }

// NewSigner returns the signer whose 32-byte Ed25519 seed is the byte seed, repeated.
func NewSigner(seed byte) *Signer {
	return &Signer{priv: ed25519.NewKeyFromSeed(bytes.Repeat([]byte{seed}, ed25519.SeedSize))}
}

func sshString(dst, s []byte) []byte {
	dst = binary.BigEndian.AppendUint32(dst, uint32(len(s)))
	return append(dst, s...)
}

func (s *Signer) blob() []byte {
	return sshString(sshString(nil, []byte("ssh-ed25519")), s.priv.Public().(ed25519.PublicKey))
}

// PublicKey is the signer's key as host/approval holds it.
func (s *Signer) PublicKey() approval.SSHPublicKey {
	return approval.SSHPublicKey{Type: "ssh-ed25519", Blob: s.blob()}
}

// AuthorizedKey is the signer's .pub file: "ssh-ed25519 <base64> approvaltest\n".
func (s *Signer) AuthorizedKey() []byte {
	return []byte("ssh-ed25519 " + base64.StdEncoding.EncodeToString(s.blob()) + " approvaltest\n")
}

// Sign returns the armored SSHSIG signature of msg under namespace, hashed with sha512, wrapped
// at 70 columns as ssh-keygen writes it.
func (s *Signer) Sign(namespace string, msg []byte) []byte {
	h := sha512.Sum512(msg)
	signed := []byte("SSHSIG")
	signed = sshString(signed, []byte(namespace))
	signed = sshString(signed, nil)
	signed = sshString(signed, []byte("sha512"))
	signed = sshString(signed, h[:])
	sig := sshString(sshString(nil, []byte("ssh-ed25519")), ed25519.Sign(s.priv, signed))
	blob := binary.BigEndian.AppendUint32([]byte("SSHSIG"), 1)
	blob = sshString(blob, s.blob())
	blob = sshString(blob, []byte(namespace))
	blob = sshString(blob, nil)
	blob = sshString(blob, []byte("sha512"))
	blob = sshString(blob, sig)
	b64 := base64.StdEncoding.EncodeToString(blob)
	var out strings.Builder
	out.WriteString("-----BEGIN SSH SIGNATURE-----\n")
	for len(b64) > 70 {
		out.WriteString(b64[:70] + "\n")
		b64 = b64[70:]
	}
	out.WriteString(b64 + "\n-----END SSH SIGNATURE-----\n")
	return []byte(out.String())
}
```

### 2.3 `host/approval/testdata/approvers/generate.sh` (new, mode 0755; M1b)

```bash
#!/usr/bin/env bash
# Regenerates the approver refusal fixtures (row 201 A1b, test AP1) with REAL ssh-keygen. Test
# keys only: made in a temp dir with HOME pointed at it and no agent, deleted on exit; only the
# two public keys and SOURCE.txt are kept (text files with no NUL byte: the tracked-binary gate).
# After a run, update approverFixtureSHA256 in ../../approver_test.go IN THE SAME COMMIT (AP0).
set -euo pipefail
here=$(cd "$(dirname "$0")" && pwd)
tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
export HOME="$tmp"; unset SSH_AUTH_SOCK
ssh-keygen -q -t rsa -b 3072 -N '' -C world-approver-fixture-rsa3072 -f "$tmp/rsa3072" < /dev/null
ssh-keygen -q -t ecdsa -b 384 -N '' -C world-approver-fixture-p384 -f "$tmp/p384" < /dev/null
cp "$tmp/rsa3072.pub" "$tmp/p384.pub" "$here/"
{ ssh -V 2>&1; date -u +%FT%TZ; ssh-keygen -l -f "$here/rsa3072.pub"; ssh-keygen -l -f "$here/p384.pub"; } > "$here/SOURCE.txt"
echo "regenerated; now update approverFixtureSHA256 in approver_test.go: shasum -a 256 $here/*" >&2
```

The executor runs it **exactly once**, from the worktree, as `bash
host/approval/testdata/approvers/generate.sh`, then `chmod 0755` it, then pins the four SHA-256
values in AP0. `.gitattributes` gains, after the `sshsig` line, exactly:
`host/approval/testdata/approvers/* -text`.

### 2.4 `host/store/credentials_schema.sql` (append after `mint_requests`; M2a)

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

### 2.5 `host/store/credentials.go` edits (M2a)

1. `:38-42` — the comment's last line and the constant become:
   ```go
   // version 4 adds mint requests (row 201); version 5 adds approvers (row 201 A1b).
   const credentialsSchemaVersion = 5
   ```
2. `credentialsTables` (`:71-80`): replace `case version >= 4:` with
   ```go
   	case version >= 5:
   		return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable}
   	case version == 4:
   		return []string{"episode_starts", mintRequestsTable, credentialsTable}
   ```
   (the `version == 3` and `default` arms unchanged).
3. `upgradeCredentials` (`:353-355`): after the `if version < 4 { … mintRequestsSQL }` block add
   ```go
   	if version < 5 {
   		statements = append(statements, approversSQL)
   	}
   ```

Nothing else in the file changes. `credentials.go:432`'s bound-mint gate keeps comparing with the
current version (design R2): a v4 store refuses bound, started and resumed mints from this build
until `serve` of this build restarts, as v3 stores did at the v4 bump.

### 2.6 `host/store/mint_requests.go` edits (M2a)

Replace `:59-69` with:

```go
// OlderCredentialsStoreError refuses an operation on a store below the version its Feature needs:
// mint requests need 4, approvers need 5 (row 201 A1b §2.1).
type OlderCredentialsStoreError struct {
	Path    string
	Version int    // the store's version
	Need    int    // the first version with Feature
	Feature string // "mint requests" or "approvers"
}

func (e *OlderCredentialsStoreError) Error() string {
	return fmt.Sprintf("store: credentials store %q is schema %d; %s need schema %d: "+
		"restart `ailang-worldd serve` with this build once (it upgrades the store at startup), then run this again",
		e.Path, e.Version, e.Feature, e.Need)
}
```

and at each of the four sites (`:93`, `:147`, `:167`, `:254`) replace
`&OlderCredentialsStoreError{Path: c.display, Version: c.version}` with
`&OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: mintRequestsSchemaVersion, Feature: "mint requests"}`.
(`grep -c 'Need: mintRequestsSchemaVersion, Feature: "mint requests"' host/store/mint_requests.go` → `4`.)

### 2.7 `host/store/approvers.go` (new; M2a holds the first block, M2b the second, M2c the third)

```go
package store

// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.1, §2.4): the approvers
// table (credentials schema 5). A row is written only by `approver enroll` after a typed y; it is
// marked revoked, never deleted, so a name and a key are each enrolled at most once, ever. The
// store does not import host/approval: the name grammar and key parsing are the pure law's, run
// by the CLI first; the checks here are the store's own backstop.

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ---- M2a ----
const approversTable = "approvers"

// approversSchemaVersion is the first credentials schema with approvers.
const approversSchemaVersion = 5

// approversSQL is the version 4 -> 5 upgrade: credentials_schema.sql's approvers table, column
// for column (TestCredentialsSchemaMatchesCanonicalDDL compares the two).
const approversSQL = `CREATE TABLE approvers (
    name         TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
    format       TEXT NOT NULL CHECK (format <> ''),
    public_key   BLOB NOT NULL UNIQUE CHECK (length(public_key) > 0),
    enrolled_at  INTEGER NOT NULL CHECK (enrolled_at > 0),
    last_used_at INTEGER CHECK (last_used_at IS NULL OR last_used_at >= enrolled_at),
    revoked_at   INTEGER CHECK (revoked_at IS NULL OR revoked_at >= enrolled_at)
)`

// ApproverRow is one approvers row. LastUsedAt and RevokedAt are 0 for NULL.
type ApproverRow struct {
	Name       string
	Format     string
	PublicKey  []byte // the SSH wire blob
	EnrolledAt int64
	LastUsedAt int64 // 0 = NULL (never used)
	RevokedAt  int64 // 0 = NULL (active)
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

// maxListedApprovers caps ListApprovers.
const maxListedApprovers = 256

// ---- M2b ----

// olderForApprovers is the refusal of every approver method below version 5.
func (c *Credentials) olderForApprovers() error {
	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: approversSchemaVersion, Feature: "approvers"}
}

const approverColumns = `name, format, public_key, enrolled_at, last_used_at, revoked_at`

func scanApprover(scan func(...any) error) (ApproverRow, error) {
	var r ApproverRow
	var lastUsed, revoked sql.NullInt64
	if err := scan(&r.Name, &r.Format, &r.PublicKey, &r.EnrolledAt, &lastUsed, &revoked); err != nil {
		return ApproverRow{}, err
	}
	r.LastUsedAt, r.RevokedAt = lastUsed.Int64, revoked.Int64
	return r, nil
}

// EnrolApprover stores one new, active approver. It refuses a name or a key that any row holds,
// revoked rows included; the UNIQUE constraints are the backstop and their error text is never
// parsed.
func (c *Credentials) EnrolApprover(ctx context.Context, row ApproverRow) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.version < approversSchemaVersion {
		return c.olderForApprovers()
	}
	if row.Name == "" {
		return fmt.Errorf("%w: empty name", ErrBadApprover)
	}
	if row.Format != "ssh" {
		return fmt.Errorf("%w: format %q is not ssh", ErrBadApprover, row.Format)
	}
	if len(row.PublicKey) == 0 {
		return fmt.Errorf("%w: empty public key", ErrBadApprover)
	}
	if row.EnrolledAt <= 0 {
		return fmt.Errorf("%w: enrolled_at %d is not positive", ErrBadApprover, row.EnrolledAt)
	}
	if row.LastUsedAt != 0 {
		return fmt.Errorf("%w: a new approver has never been used", ErrBadApprover)
	}
	if row.RevokedAt != 0 {
		return fmt.Errorf("%w: a new approver is not revoked", ErrBadApprover)
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: enrol approver: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	var n int
	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM approvers WHERE name = ?`, row.Name).Scan(&n); err != nil {
		return fmt.Errorf("store: enrol approver: %w", err)
	}
	if n != 0 {
		return fmt.Errorf("%w: %q", ErrApproverNameTaken, row.Name)
	}
	var holder string
	err = tx.QueryRowContext(ctx, `SELECT name FROM approvers WHERE public_key = ?`, row.PublicKey).Scan(&holder)
	if err == nil {
		return fmt.Errorf("%w: it is approver %q's key", ErrApproverKeyEnrolled, holder)
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("store: enrol approver: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES (?, ?, ?, ?)`,
		row.Name, row.Format, row.PublicKey, row.EnrolledAt); err != nil {
		return fmt.Errorf("store: enrol approver: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: enrol approver: commit: %w", err)
	}
	return nil
}

// Approver reads one approver by name. ok is false when no such row exists.
func (c *Credentials) Approver(ctx context.Context, name string) (ApproverRow, bool, error) {
	return c.approverWhere(ctx, `name = ?`, name)
}

// ApproverByKey reads the approver whose public_key is blob, byte for byte. ok is false when none.
func (c *Credentials) ApproverByKey(ctx context.Context, blob []byte) (ApproverRow, bool, error) {
	return c.approverWhere(ctx, `public_key = ?`, blob)
}

func (c *Credentials) approverWhere(ctx context.Context, where string, arg any) (ApproverRow, bool, error) {
	if err := requireDeadline(ctx); err != nil {
		return ApproverRow{}, false, err
	}
	if c.version < approversSchemaVersion {
		return ApproverRow{}, false, c.olderForApprovers()
	}
	row, err := scanApprover(c.db.QueryRowContext(ctx, `SELECT `+approverColumns+` FROM approvers WHERE `+where, arg).Scan)
	if errors.Is(err, sql.ErrNoRows) {
		return ApproverRow{}, false, nil
	}
	if err != nil {
		return ApproverRow{}, false, fmt.Errorf("store: read approver: %w", err)
	}
	return row, true, nil
}

// ---- M2c ----

// ListApprovers lists every approver, revoked ones included, ordered by name, at most 256.
func (c *Credentials) ListApprovers(ctx context.Context) ([]ApproverRow, error) {
	if err := requireDeadline(ctx); err != nil {
		return nil, err
	}
	if c.version < approversSchemaVersion {
		return nil, c.olderForApprovers()
	}
	rows, err := c.db.QueryContext(ctx, `SELECT `+approverColumns+` FROM approvers ORDER BY name LIMIT ?`, maxListedApprovers)
	if err != nil {
		return nil, fmt.Errorf("store: list approvers: %w", err)
	}
	defer func() { _ = rows.Close() }()
	var out []ApproverRow
	for rows.Next() {
		r, err := scanApprover(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("store: list approvers: %w", err)
		}
		out = append(out, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("store: list approvers: %w", err)
	}
	return out, nil
}

// RevokeApprover marks the active approver name revoked at now. The row is kept. It refuses an
// unknown name, a row already revoked, and a clock behind the enrolment, each with its own error.
func (c *Credentials) RevokeApprover(ctx context.Context, name string, now int64) error {
	if err := requireDeadline(ctx); err != nil {
		return err
	}
	if c.version < approversSchemaVersion {
		return c.olderForApprovers()
	}
	tx, err := c.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("store: revoke approver: begin: %w", err)
	}
	defer func() { _ = tx.Rollback() }()
	res, err := tx.ExecContext(ctx, `UPDATE approvers SET revoked_at = ? WHERE name = ? AND revoked_at IS NULL AND enrolled_at <= ?`, now, name, now)
	if err != nil {
		return fmt.Errorf("store: revoke approver: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("store: revoke approver: %w", err)
	}
	if n != 1 {
		return classifyRevoke(ctx, tx, name, now)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("store: revoke approver: commit: %w", err)
	}
	return nil
}

// classifyRevoke names why a revoke changed no row, reading the row in the same transaction.
// The order matches approval.CheckRevoke: revoked before clock-behind.
func classifyRevoke(ctx context.Context, tx *sql.Tx, name string, now int64) error {
	row, err := scanApprover(tx.QueryRowContext(ctx, `SELECT `+approverColumns+` FROM approvers WHERE name = ?`, name).Scan)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return fmt.Errorf("%w: %q", ErrApproverUnknown, name)
	case err != nil:
		return fmt.Errorf("store: revoke approver: %w", err)
	case row.RevokedAt != 0:
		return fmt.Errorf("%w: %q at %d", ErrApproverRevoked, name, row.RevokedAt)
	case row.EnrolledAt > now:
		return fmt.Errorf("%w: %q enrolled at %d, now %d", ErrApproverClockBehind, name, row.EnrolledAt, now)
	default:
		return fmt.Errorf("store: revoke approver %q: changed no row", name)
	}
}
```

In M2a the file's import block is only what the M2a block uses (`errors`); M2b adds `context`,
`database/sql`, `fmt`. `gofmt` and `go vet` decide; the executor adjusts imports only.

### 2.8 Test helpers (exact)

`host/store/approvers_test.go` (M2a):

```go
// makeV4Credentials writes a schema-4 credentials store at path: the version-4 table set
// (episode_starts, mint_requests, session_credentials), the given bound rows and user_version 4.
func makeV4Credentials(t *testing.T, path string, rows []SessionRow) {
	t.Helper()
	makeV3Credentials(t, path, rows)
	db := rawDB(t, path)
	if _, err := db.Exec(mintRequestsSQL); err != nil {
		t.Fatal(err)
	}
	setPositiveFixtureVersion(t, db, 4)
}

// approverCount counts the approvers rows of c's store.
func approverCount(t *testing.T, c *Credentials) int // SELECT count(*) FROM approvers via c.db; t.Fatal on error

// testApprover is a valid new row: Name name, Format "ssh", PublicKey []byte("key-" + name), EnrolledAt at.
func testApprover(name string, at int64) ApproverRow
```

`host/store/credentials_test.go` (M2a) gains, after `canonicalMintRequestsDDL`:

```go
// canonicalApproversDDL is the current (version 5, row 201 A1b) approvers table, review-visible
// like canonicalMintRequestsDDL.
var canonicalApproversDDL = `CREATE TABLE approvers (
    name         TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
    format       TEXT NOT NULL CHECK (format <> ''),
    public_key   BLOB NOT NULL UNIQUE CHECK (length(public_key) > 0),
    enrolled_at  INTEGER NOT NULL CHECK (enrolled_at > 0),
    last_used_at INTEGER CHECK (last_used_at IS NULL OR last_used_at >= enrolled_at),
    revoked_at   INTEGER CHECK (revoked_at IS NULL OR revoked_at >= enrolled_at)
)`
```

`host/approval/approver_test.go` (M1b), package `approval` (internal, so it may use `sshReader`
and `appendString`; it must **not** import `approvaltest`, which imports `approval`: an import cycle):

```go
// approverFixtureBytes reads one file of testdata/approvers.
func approverFixtureBytes(t *testing.T, name string) []byte // os.ReadFile(filepath.Join("testdata", "approvers", name)); t.Fatal on error

// paddedKeyLine is ed25519.pub's line (no newline) + " " + 'x' padding, exactly n bytes long.
// t.Fatal when n is too small to hold the line and one pad byte.
func paddedKeyLine(t *testing.T, n int) []byte

// skRetype is ed25519.pub with its key type rewritten to sk-ssh-ed25519@openssh.com in the line
// and in the blob (appendString(nil, "sk-ssh-ed25519@openssh.com") + the blob after its type string).
func skRetype(t *testing.T) []byte

// derivedPrivateKeyFile is an OpenSSH private-key armor around derived noise (no real key):
// "-----BEGIN OPENSSH PRIVATE KEY-----\n" + noise + "\n-----END OPENSSH PRIVATE KEY-----\n", where
// s := sha512.Sum512([]byte("approver-test-noise")); noise := base64.StdEncoding.EncodeToString(bytes.Repeat(s[:], 4)).
func derivedPrivateKeyFile() (file []byte, noise string)

const (
	fixtureEd25519Fingerprint = "SHA256:pHMYJm+CndqyZGrzMHQQBH0tNeNGj1U9sy0JrVsYdSQ"
	fixtureP256Fingerprint    = "SHA256:yMwGgMZdq1aVNYRG/HO2T69s/Fb14yPZW0GBWtwRyEo"
)
```

---

## 3. Milestones

Every sub-milestone ends with: `test -z "$(gofmt -l <touched dirs>)"`; `go build ./...` **and**
`go vet` of each touched package (the compile fence for `_test.go`); its own tests with `-count=1`;
and the earlier sub-milestones' tests. No `t.Parallel()`. Run any loop or `$VAR` expansion under
`bash -c` (the controller's shell is zsh). Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in the
**same** shell as any `cmd/ailang-worldd` or `host/daemon` test (unset, `host/daemon` tests red —
an instrument fault, not a regression). After each sub-milestone, snapshot (§7).

### M1a — `CR-TOOMANY`, and the test signer · ~75 test-support + ~150 test LOC

**Files**: `cmd/ailang-worldd/approve_test.go` (+CT1, appended; every import it needs —
`context`, `fmt`, `strings`, `time`, `host/store` — is already there), new
`host/approval/approvaltest/signer.go` (§2.2), new
`host/approval/approvaltest/signer_test.go` (**package `approvaltest_test`**).

**CT1 `TestEpisodeStartRequestRefusesWhenSixteenArePending`** (appended after
`TestEpisodeStartRequestRefusesAnOlderStore`; validated by the planner, C18):
1. `r := newEpisodeRun(t)`; `p := writeTestProfile(t, r.profiles, "proj", "")`;
   `stubDaemon(t, matchingHealth(t, r.f, r.profiles, p), nil)`; a task file `do the thing\n` in
   `t.TempDir()`.
2. `credPath, _ := credentialsFor(r.f.db, "")`; `st, err := store.OpenCredentials(credPath)`; a 5 s
   context; for `i := 1..store.MaxPendingMintRequests`: `st.PutMintRequest(ctx,
   store.MintRequestRow{RequestID: fmt.Sprintf("mr-%016x", 0xf00+i), Request: []byte(fmt.Sprintf(`{"n":%d}`, i)),
   Digest: strings.Repeat("ab", 32), Challenge: strings.Repeat("cd", 32), EpisodeID: "seed",
   CreatedAt: approveT0, ExpiresAt: approveT0 + 600}, approveT0)` (each `t.Fatalf` on error);
   `st.Close()`.
3. `runEpisodeStart(r.args("ep1", "--profile", "proj", "--task", task, "--model", "claude-opus-5-5", "--request"), &stdout, &stderr, requestEnv(&opened, approveT0))`
   → `code == exitUsage` (by constant); stderr contains `16 mint requests are already pending`;
   `opened == false`; `pendingRows(t, r.f.db, approveT0)` has exactly `store.MaxPendingMintRequests`
   rows and none has `RequestID == firstRequestID`.
Kills **MUT-CT-TOOMANY**.

**TS1 `TestTestSignerVerifies`**: `s1, s2 := approvaltest.NewSigner(1), approvaltest.NewSigner(2)`;
`msg, err := approval.SignedMessage([32]byte{1}, [32]byte{2}, 1800000600)` (nil err). Asserts:
(a) `approval.ParseAuthorizedKey(s1.AuthorizedKey())` → `Type == "ssh-ed25519"`, `Blob` equal to
`s1.PublicKey().Blob`; `AuthorizedKey()` ends with `"\n"`; (b) `s1.PublicKey().Blob` ≠
`s2.PublicKey().Blob`; (c) `bytes.Equal(s1.Sign(approval.Namespace, msg), approvaltest.NewSigner(1).Sign(approval.Namespace, msg))`
(deterministic); (d) `p, err := approval.ParseSSHSIG(s1.Sign(approval.Namespace, msg))` → nil err,
`p.HashAlg == "sha512"`, `p.Namespace == approval.Namespace`, `p.PublicKey.Blob` equal;
`approval.VerifySSHSIG(p, approval.Namespace, msg, s1.PublicKey()) == nil`; (e) the same over
`msg` with byte 0 flipped → `errors.Is(err, approval.ErrBadSignature)`; (f) against
`s2.PublicKey()` → `errors.Is(err, approval.ErrKeyMismatch)`; (g)
`approval.VerifyApproval(s1.Sign(approval.Namespace, msg), [32]byte{1}, [32]byte{2}, 1800000600, s1.PublicKey()) == nil`
and with expiry `1800000601` → `errors.Is(err, approval.ErrBadSignature)`.
Its job is to be correct, not to kill a production mutant; XS2 (sprint 2) pins it to the real
tool, and the planner already measured that agreement (C15).

**TS3 `TestApprovalTestIsTestOnly`** (an import guard, labelled a control): find the repo root
(walk up from the working directory to the first directory with `go.mod`); `filepath.WalkDir` it,
skipping directories named `node_modules`, `vendor`, or starting with `.` (other than the root);
for every `.go` file, `parser.ParseFile(fset, path, nil, parser.ImportsOnly)`; record whether it
imports `github.com/sunholo-data/ailang-world/host/approval/approvaltest`. Assert: no file without
the `_test.go` suffix imports it (`t.Errorf` naming each); **instrument controls**: at least 300 `.go`
files were parsed, and at least one `_test.go` file imports it (this file itself, an external test
package — D6).

**Acceptance**
```
test -z "$(gofmt -l host/approval cmd/ailang-worldd)"
go build ./... && go vet ./host/approval/... ./cmd/ailang-worldd/
export AILANG_BIN=$HOME/.pinned-ailang/ailang
go test ./cmd/ailang-worldd/ -run '^TestEpisodeStartRequestRefusesWhenSixteenArePending$' -count=1 -v      # --- PASS, no SKIP
go test ./host/approval/approvaltest/ -run '^(TestTestSignerVerifies|TestApprovalTestIsTestOnly)$' -count=1 -v
go test ./host/approval/... -count=1
go test ./cmd/ailang-worldd/ -run '^(TestEpisodeStartRequest|TestApprove)' -count=1
```

### M1b — the enrolment law and its fixtures · ~80 prod + ~230 test + ~26 fixture LOC

**Files**: new `host/approval/approver.go` (§2.1 down to `CheckRevoke`), new
`host/approval/approver_test.go` (§2.8 helpers, AP0–AP4), new `host/approval/testdata/approvers/`
(`generate.sh` §2.3, then its outputs `rsa3072.pub`, `p384.pub`, `SOURCE.txt`), `.gitattributes`
(+1 line, §2.3).

**Steps**: write `generate.sh`; `chmod 0755`; run it once (`bash host/approval/testdata/approvers/generate.sh`);
check the outputs are NUL-free text (acceptance); write `approver.go`; write the tests; pin AP0's
hashes from `shasum -a 256 host/approval/testdata/approvers/*`.

**AP0 `TestApproverFixturesAreTheGeneratedBytes`**: `var approverFixtureSHA256 = map[string]string{"SOURCE.txt": …, "generate.sh": …, "p384.pub": …, "rsa3072.pub": …}`
(the four measured values). `os.ReadDir("testdata/approvers")` holds exactly those four names;
each file's SHA-256 equals its pin; each file contains no `0x00` byte and ends with `"\n"`;
`rsa3072.pub` starts with `ssh-rsa ` and `p384.pub` with `ecdsa-sha2-nistp384 `.

**AP1 `TestNewApproverRefusals`** — table-driven, `now = 1000`, `good := fixtureBytes(t,
"ed25519.pub")`. Every arm asserts `errors.Is(err, want)` **and** `reflect.DeepEqual(a, Approver{})`.

| Arm | Call | Want | Kills |
|---|---|---|---|
| name-empty, name-upper `Mark`, name-space `a b`, name-dash-first `-x`, name-dot-first `.x`, name-49 `strings.Repeat("a", 49)`, name-escape `"x\x1b[2J"`, name-newline `"mark\n"` | `NewApprover(name, "ssh", good, now)` | `ErrApproverName` | MUT-E1, MUT-PAT-ANCHOR (`mark\n`, `a b`) |
| webauthn | `NewApprover("mark", "webauthn", good, now)` | `ErrFormatNotYet` | MUT-E2, MUT-E2-ORDER |
| webauthn-empty-key | `NewApprover("mark", "webauthn", nil, now)` | `ErrFormatNotYet` | MUT-E2-ORDER (format checks precede the file checks) |
| format-gpg, format-empty, format-SSH | `NewApprover("mark", f, good, now)` | `ErrApproverFormat` | MUT-E3 |
| key-nil, key-empty | `NewApprover("mark", "ssh", k, now)` | `ErrApproverKeyFile` | MUT-E4, MUT-E4-EMPTY |
| key-8193 | `paddedKeyLine(t, 8193)` | `ErrApproverKeyFile` | MUT-E4 |
| private | `derivedPrivateKeyFile()` → `file` | `ErrPrivateKeyFile`; `err.Error()` does **not** contain `noise[:24]` | MUT-E5 |
| private-leading-space | `append([]byte("\n  \t"), file...)` | `ErrPrivateKeyFile` | MUT-E5-TRIM |
| two-lines | `append(bytes.Clone(good), good...)` | `ErrApproverKeyFile` | MUT-E6 |
| blank-second-line | `append(bytes.Clone(good), '\n')` | `ErrApproverKeyFile` | MUT-E6 |
| rsa | `approverFixtureBytes(t, "rsa3072.pub")` | `ErrKeyType` | MUT-E7 |
| p384 | `approverFixtureBytes(t, "p384.pub")` | `ErrKeyType` | MUT-E7 |
| sk | `skRetype(t)` | `ErrKeyType`, and the text contains `row 201 Q1` | MUT-E7 |
| one-field | `[]byte("ssh-ed25519\n")` | `ErrAuthorizedKey` | MUT-E7 |
| bad-base64 | `[]byte("ssh-ed25519 !!!! c\n")` | `ErrBase64` | MUT-E7 |

**AP2 `TestNewApproverAcceptsTheRealFixtures`**: for `ed25519.pub` (type `ssh-ed25519`,
fingerprint `fixtureEd25519Fingerprint`) and `p256.pub` (`ecdsa-sha2-nistp256`,
`fixtureP256Fingerprint`), each **with** its trailing newline and with it trimmed:
`a, err := NewApprover("mark-mac", "ssh", b, 1234)` → nil err; `a.Name == "mark-mac"`, `a.Format ==
FormatSSH`, `a.Key.Type` as above, `bytes.Equal(a.Key.Blob, want.Blob)` where `want, _ :=
ParseAuthorizedKey(b)`, `a.Key.Fingerprint()` as above, `a.EnrolledAt == 1234`, `a.LastUsedAt == 0`,
`a.RevokedAt == 0`. Boundary arm: `NewApprover("mark-mac", "ssh", paddedKeyLine(t, MaxKeyFileBytes), 1234)`
→ nil err, `len == 8192` asserted first. Kills MUT-NA-NAME, MUT-NA-NOW, MUT-NA-NL, MUT-E4-GE.

**AP3 `TestCheckRevoke`** (`a := Approver{Name: "x", Format: "ssh", EnrolledAt: 1000}`): `(a, 1000)`
→ nil; `(a, 2000)` → nil; `(a, 999)` → `ErrApproverClockBehind`; `a.RevokedAt = 1500`: `(a, 2000)`
→ `ErrApproverRevoked`; and revoked **and** behind, `(a, 999)` → `ErrApproverRevoked`. Kills
MUT-R1, MUT-R2, MUT-R2-LE, MUT-R-ORDER.

**AP4 `TestApproverNamePatternBoundaries`**: accepted `a`, `0`, `strings.Repeat("a", 48)`,
`a.b_c-d`, `mark-mac`, `a-`, `9.`; refused `strings.Repeat("a", 49)`, `""`, `.a`, `_a`, `-a`, `A`,
`a/b`, `a:b`, `a b`, `ä`, `a\n`. Kills MUT-PAT-48, MUT-PAT-FIRST, MUT-PAT-ANCHOR.

`TestApprovalImportsOnlyTheStandardLibrary` (existing) stays green with `approver.go`.

**Acceptance**
```
test -z "$(gofmt -l host/approval)"
go build ./... && go vet ./host/approval/...
bash -c 'for f in host/approval/testdata/approvers/*; do tr -d "\000" < "$f" | cmp -s - "$f" || { echo "NUL in $f"; exit 1; }; done'
test -x host/approval/testdata/approvers/generate.sh
grep -qx 'host/approval/testdata/approvers/\* -text' .gitattributes
go test ./host/approval/ -run '^(TestApproverFixturesAreTheGeneratedBytes|TestNewApproverRefusals|TestNewApproverAcceptsTheRealFixtures|TestCheckRevoke|TestApproverNamePatternBoundaries|TestApprovalImportsOnlyTheStandardLibrary)$' -count=1 -v
go test ./host/approval/... -count=1
test -z "$(git status --short host/approval/testdata/sshsig)"            # SG1's fixtures untouched
```

### M1c — the enrolment prompt body · ~15 prod + ~125 test LOC

**Files**: `host/approval/approver.go` (+`EnrolmentWarning`, `RenderEnrolment`, `"fmt"` import),
`host/approval/approver_test.go` (+AP5, RC4).

**AP5 `TestRenderEnrolmentGolden`**: `a, _ := NewApprover("mark-mac", "ssh", fixtureBytes(t,
"ed25519.pub"), 1800000000)`; `RenderEnrolment(a)` equals a **literal** golden string in the test
(the exact text of §2.1, spelled out, not built from `EnrolmentWarning`), compared with `!=` and a
line diff in the failure message. Kills MUT-RE-FP, MUT-RE-R6, MUT-RE-REC.

**RC4 `TestEnrolmentPromptShowsEveryPersistedField`**:
```go
// enrolmentPromptFields is the disposition of every Approver field (and every field of its
// element type SSHPublicKey): shown on the enrolment y/N, or excluded with the reason.
var enrolmentPromptFields = map[string]string{
	"Name":       "shown",
	"Format":     "shown",
	"Key.Type":   "shown",
	"Key.Blob":   "shown as its SHA256 fingerprint",
	"EnrolledAt": "excluded: a timestamp the store sets",
	"LastUsedAt": "excluded: a timestamp the store sets",
	"RevokedAt":  "excluded: a timestamp the store sets",
}
```
Algorithm: (1) walk `reflect.TypeOf(Approver{})`; a field of struct kind (`Key`) is walked into,
giving paths `Key.Type`, `Key.Blob`; every path must be a key of the map and the path count must
equal `len(enrolmentPromptFields)` (a new field in `Approver` **or** `SSHPublicKey` reds). (2) base
`a` from `ed25519.pub` as `mark-mac`, `EnrolledAt 1000`, `LastUsedAt 1500`, `RevokedAt 1700`; `base
:= RenderEnrolment(a)`. (3) For each path, perturb a copy through reflection — a `string`
`SetString(v.String() + "x")`; a `[]byte` cloned with byte 0 XOR `0xff`; an `int64` `+7` — and
render: a `shown…` path must change the text, an `excluded…` path must not. Kills MUT-RC4-FIELD,
MUT-RC4-KEYFIELD, MUT-RE-TYPE (and, again, MUT-RE-FP via `Key.Blob`).

**Acceptance**
```
test -z "$(gofmt -l host/approval)"
go build ./... && go vet ./host/approval/...
go test ./host/approval/ -run '^(TestRenderEnrolmentGolden|TestEnrolmentPromptShowsEveryPersistedField)$' -count=1 -v
go test ./host/approval/... -count=1 && go test -race ./host/approval/... -count=1
```

### M2a — credentials schema 5 · ~89 prod + ~209 test LOC

**Files**: `host/store/credentials_schema.sql` (§2.4), `host/store/credentials.go` (§2.5),
`host/store/mint_requests.go` (§2.6), new `host/store/approvers.go` (§2.7, `---- M2a ----` block
only), `host/store/credentials_test.go` (AS1 edit, `canonicalApproversDDL`, AS7),
`host/store/mint_requests_test.go` (one-line edit at `:132`), new `host/store/approvers_test.go`
(§2.8 `makeV4Credentials`, AS5, AS6).

**AS1 `TestCredentialsSchemaMatchesCanonicalDDL`** (edit): `:56`'s list becomes
`[]string{"approvers", "episode_starts", "mint_requests", "session_credentials"}`; add
`normalizeDDL(ddl["approvers"]) == normalizeDDL(canonicalApproversDDL)` and
`normalizeDDL(approversSQL) == normalizeDDL(canonicalApproversDDL)` (each its own `t.Fatalf`). The
`user_version` assertion stays symbolic. Kills MUT-DDL-SCHEMA, MUT-DDL-CONST.

**`TestCredentialsUpgradeToV4`** (edit, `mint_requests_test.go:132` only): the list becomes
`[]string{"approvers", "episode_starts", "mint_requests", "session_credentials"}` (an upgrade from
1, 2 or 3 now lands on 5). Nothing else in it changes.

**AS5 `TestCredentialsUpgradeToV5`** (subtests `v1`, `v3`, `v4`): v1 =
`makeV1Credentials(t, path, rigRows())`; v3 = `makeV3Credentials(t, path, []SessionRow{oneBoundRow()})`
plus raw `INSERT INTO episode_starts (episode, start_json) VALUES ('w900', '{"episode":"w900","note":"kept"}')`;
v4 = `makeV4Credentials` with the same row and start, plus one request inserted raw:
`INSERT INTO mint_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at) VALUES ('mr-00000000000000aa', x'7b7d', <64×"ab">, <64×"cd">, 'ep', 1000, 1600)`.
Then: (a) `OpenCredentials(path)` → `Version() == v`; close; `userVersion == v`; `tableNamesAt`
unchanged (a session verb never upgrades). (b) `UpgradeCredentials(path)` → `(true, nil)`; again →
`(false, nil)`. (c) `userVersion == credentialsSchemaVersion` (5); `requireExactTableNames` the four
tables; `normalizeDDL(ddl["approvers"]) == normalizeDDL(canonicalApproversDDL)`; `SELECT count(*)
FROM approvers` → 0; every original session row `ResolveSession`s `reflect.DeepEqual` to what was
inserted; v3/v4: the `episode_starts` row's `start_json` unchanged; v4: the request row's seven
columns read back unchanged and `consumed_at IS NULL`. Kills MUT-UPG-NOV5, MUT-UPG-CLI, MUT-CT-V4CASE.

**AS6 `TestOlderCredentialsStoreErrorText`**: (1) on a v3 store (`makeV3Credentials(t, path, nil)`,
`c, _ := OpenCredentials(path)`): `want := fmt.Sprintf("store: credentials store %q is schema 3; mint
requests need schema 4: restart `+"`ailang-worldd serve`"+` with this build once (it upgrades the store at
startup), then run this again", path)`; each of `PutMintRequest(ctx, testMintRequestRow(1, 1000),
1000)`, `MintRequest(ctx, "mr-0000000000000001")`, `PendingMintRequests(ctx, 1000, 0)`,
`MintSessionForRequest(ctx, consumeSession(1), okConsume(testMintRequestRow(1, 1000), 1010, new(int), new(int)))`
returns an error with `err.Error() == want` exactly, and `errors.As` gives `Need == 4`, `Feature ==
"mint requests"`, `Version == 3`. (2) `(&OlderCredentialsStoreError{Path: "p", Version: 4, Need: 5,
Feature: "approvers"}).Error()` equals `"store: credentials store \"p\" is schema 4; approvers need schema 5: restart `+"`ailang-worldd serve`"+` with this build once (it upgrades the store at startup), then run this again"`.
Kills MUT-OLD-FMT, MUT-OLD-PUT, MUT-OLD-GET, MUT-OLD-LIST, MUT-OLD-CONSUME.

**AS7 `TestCredentialsV5WithoutApproversIsRefused`** (`credentials_test.go`, the twin of the
existing `TestCredentialsV4WithoutMintRequestsIsRefused`): a fresh store; raw `DROP TABLE approvers`;
`OpenCredentials` and `OpenCredentialsReadOnly` both fail with `*CredentialsSchemaError`
(`errors.As`); `schemaState` unchanged by the refused opens. Kills MUT-TABLESET5.

`TestCredentialsV4WithoutMintRequestsIsRefused` (existing, unedited) now drops `mint_requests` from a
v5 store and must stay green.

**Acceptance**
```
test -z "$(gofmt -l host/store)"
go build ./... && go vet ./host/store/
go test ./host/store/ -run '^(TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsUpgradeToV4|TestCredentialsUpgradeToV5|TestOlderCredentialsStoreErrorText|TestCredentialsV5WithoutApproversIsRefused|TestCredentialsV4WithoutMintRequestsIsRefused)$' -count=1 -v
go test ./host/store/ -count=1
export AILANG_BIN=$HOME/.pinned-ailang/ailang
go test ./cmd/ailang-worldd/ -run '^(TestEpisodeStartRequest|TestApprove|TestSessionMint_|TestSessionNew)' -count=1
go test ./host/daemon/ -run 'Credentials|Split|Binding' -count=1 && go test ./host/authority/ -count=1
```

### M2b — enrol and look up · ~110 prod + ~229 test LOC

**Files**: `host/store/approvers.go` (`---- M2b ----` block), `host/store/approvers_test.go` (§2.8
`approverCount`, `testApprover`; AS2, AS3, RC2), `host/store/credentials_test.go`
(`TestCredentialsMethodsRequireADeadline`: +3 entries).

**AS2 `TestEnrolApproverRefusals`**. Each refusal arm asserts its typed error **and**
`approverCount` unchanged.
- **v4**: `makeV4Credentials`; `EnrolApprover(ctx, testApprover("mark-mac", 1000))` →
  `errors.As(err, &older)` with `Version == 4`, `Need == 5`, `Feature == "approvers"`; text contains
  `approvers need schema 5` and ``restart `ailang-worldd serve` ``; `userVersion == 4`; `SELECT
  count(*) FROM sqlite_master WHERE name = 'approvers'` → 0. (MUT-EN-OLD)
- **bad rows** on a fresh store, each `errors.Is(err, ErrBadApprover)`, 0 rows: `Name ""`
  (MUT-EN-NAME), `Format "gpg"` (MUT-EN-FORMAT), `PublicKey nil` (MUT-EN-KEY), `EnrolledAt 0`
  (MUT-EN-ENROLLED, MUT-EN-ENROLLED-LT), `EnrolledAt -1`, `LastUsedAt 1500` (MUT-EN-LASTUSED),
  `RevokedAt 1500` (MUT-EN-REVOKED).
- **good**: `testApprover("mark-mac", 1000)` → nil; 1 row.
- **name taken**: `ApproverRow{Name: "mark-mac", Format: "ssh", PublicKey: []byte("key-other"), EnrolledAt: 1100}`
  → `errors.Is(err, ErrApproverNameTaken)`; 1 row. (MUT-EN-NAMETAKEN)
- **key enrolled**: `ApproverRow{Name: "other", …, PublicKey: []byte("key-mark-mac"), …}` →
  `errors.Is(err, ErrApproverKeyEnrolled)` and `strings.Contains(err.Error(), `"mark-mac"`)`; 1 row.
  (MUT-EN-KEYENROLLED, MUT-EN-HOLDER)
- **after revoke**: raw `UPDATE approvers SET revoked_at = 1200 WHERE name = 'mark-mac'`; both
  arms again → the same typed errors; 1 row.
- **control**: `testApprover("third", 1300)` → nil; 2 rows.

**AS3 `TestApproverLookups`**: enrol `testApprover("a", 1000)` and `testApprover("b", 1100)`.
`Approver(ctx, "a")` → `ok`, nil err, `reflect.DeepEqual(row, ApproverRow{Name: "a", Format: "ssh",
PublicKey: []byte("key-a"), EnrolledAt: 1000})`; `ApproverByKey(ctx, []byte("key-b"))` → B;
`Approver(ctx, "zzz")` → `(ApproverRow{}, false, nil)`; `ApproverByKey(ctx, []byte("key-"))` (a
prefix of both) → `ok == false`, nil err; `ApproverByKey(ctx, []byte("key-a\x00"))` → `ok ==
false`. On a v4 store: both → `*OlderCredentialsStoreError` with `Need == 5`. Kills MUT-LK-OLD,
MUT-LK-NOROWS, MUT-LK-PREFIX.

**RC2 `TestApproverRowMatchesTheTable`**:
```go
// approverRowColumns maps every ApproverRow field to its approvers column.
var approverRowColumns = map[string]string{
	"Name": "name", "Format": "format", "PublicKey": "public_key",
	"EnrolledAt": "enrolled_at", "LastUsedAt": "last_used_at", "RevokedAt": "revoked_at",
}
```
(1) `reflect.TypeOf(ApproverRow{})`'s field names equal the map's keys (count and set). (2)
`PRAGMA table_info(approvers)`'s column names equal the map's values (count and set). (3) Enrol
`ApproverRow{Name: "rc2", Format: "ssh", PublicKey: []byte("rc2-key"), EnrolledAt: 1000}`; raw `UPDATE
approvers SET last_used_at = 1500, revoked_at = 1700 WHERE name = 'rc2'`; `Approver(ctx, "rc2")`
and `ApproverByKey(ctx, []byte("rc2-key"))` each return a row where, field by field through
reflection, every field is non-zero and equals `{rc2, ssh, rc2-key, 1000, 1500, 1700}`. Kills
MUT-SC-SWAP, MUT-SC-DROPCOL, MUT-RC2-FIELD.

**`TestCredentialsMethodsRequireADeadline`** (edit): add
`errs["EnrolApprover"] = c.EnrolApprover(ctx, ApproverRow{Name: "x", Format: "ssh", PublicKey: []byte("k"), EnrolledAt: 1000})`,
`_, _, errs["Approver"] = c.Approver(ctx, "x")`, `_, _, errs["ApproverByKey"] = c.ApproverByKey(ctx, []byte("k"))`.
Kills MUT-EN-DL, MUT-LK-DL.

**Acceptance**
```
test -z "$(gofmt -l host/store)"
go build ./... && go vet ./host/store/
go test ./host/store/ -run '^(TestEnrolApproverRefusals|TestApproverLookups|TestApproverRowMatchesTheTable|TestCredentialsMethodsRequireADeadline)$' -count=1 -v
go test ./host/store/ -count=1
```

### M2c — list and revoke · ~70 prod + ~133 test LOC

**Files**: `host/store/approvers.go` (`---- M2c ----` block), `host/store/approvers_test.go`
(AS3b, AS4), `host/store/credentials_test.go` (deadline test: +2 entries, +1 control).

**AS3b `TestListApprovers`**: v4 store → `*OlderCredentialsStoreError{Need 5}` (MUT-LS-OLD); a
fresh store → `len 0`, nil err; enrol `b`, `a`, `c` (in that order) and raw-revoke `b` → names
`[a b c]`, `b`'s `RevokedAt` set (revoked rows are listed; MUT-LS-ORDER, MUT-LS-ACTIVEONLY); a
second fresh store with 257 approvers inserted raw in one transaction (`n000`…`n256`, keys
`key-n000`…, `enrolled_at` 1000) → `len == 256`, first `n000`, last `n255` (MUT-LS-LIMIT).

**AS4 `TestRevokeApprover`** (one fresh store; enrol `x` and `y` at 1000; `row0 := Approver(x)`):
1. `RevokeApprover(ctx, "x", 999)` → `ErrApproverClockBehind`; `Approver(x)` deep-equals `row0`.
   (MUT-RV-ENROLLED, MUT-CL-BEHIND)
2. `RevokeApprover(ctx, "x", 1000)` → nil; `Approver(x).RevokedAt == 1000`, `LastUsedAt == 0`,
   `approverCount == 2` (MUT-RV-DELETE); `Approver(y).RevokedAt == 0` (MUT-RV-NAME).
3. `RevokeApprover(ctx, "x", 1001)` → `ErrApproverRevoked`; `RevokedAt` still 1000.
   (MUT-RV-REVOKED, MUT-RV-N1, MUT-CL-REVOKED)
4. `RevokeApprover(ctx, "x", 999)` (revoked **and** behind) → `ErrApproverRevoked` (MUT-CL-ORDER).
5. `RevokeApprover(ctx, "nope", 1001)` → `ErrApproverUnknown` (MUT-CL-NOROWS); 2 rows.
6. a v4 store → `*OlderCredentialsStoreError{Need 5}` (MUT-RV-OLD).

**`TestCredentialsMethodsRequireADeadline`** (edit): add `_, errs["ListApprovers"] =
c.ListApprovers(ctx)` and `errs["RevokeApprover"] = c.RevokeApprover(ctx, "x", 1000)`; and a
second positive control after the existing one: a bounded `ListApprovers` returns a nil error.
Kills MUT-LS-DL, MUT-RV-DL.

**Acceptance (and then the sprint gates, §6)**
```
test -z "$(gofmt -l host/store)"
go build ./... && go vet ./host/store/
go test ./host/store/ -run '^(TestListApprovers|TestRevokeApprover|TestCredentialsMethodsRequireADeadline)$' -count=1 -v
go test ./host/store/ -count=1 && go test -race ./host/store/ -run '^(TestEnrolApprover|TestApprover|TestListApprovers|TestRevokeApprover|TestCredentials)' -count=1
```

---

## 4. Mutation table (74 mutants; the controller and judge run them)

**Derivation.** One row per new production `if`, `return`, `switch` arm, SQL predicate and
copied field this sprint adds, enumerated from §2's code, not from the headline properties: the
design's E1–E7 and R1–R2, the `EnrolApprover`/`RevokeApprover`/lookup/list branches, the schema
5 table set and upgrade, the generalised error and its four sites, and CT1's carried branch.
The design's N*, L*, K*, D*, M1 (M3) and P*, B*, C*, A*, F*, G*, V* (sprint 2) rows are **not in
this sprint**; Appendix A lists M3's.

**Procedure.** Save the file (`cp f /tmp/mut.bak`), apply **one** edit, run the command, expect
**FAIL** (red, panic, or a compile error where marked), restore with `cp /tmp/mut.bak f` and
`cmp` it. **Never commit or snapshot a mutant.** Commands assume `export
AILANG_BIN=$HOME/.pinned-ailang/ailang` in the same shell.

Command keys: **A** = `go test ./host/approval/ -run '<test>$' -count=1`; **S** = `go test
./host/store/ -run '<test>$' -count=1`; **C** = `go test ./cmd/ailang-worldd/ -run '<test>$' -count=1`
(each with `^` before the test name).

| Id | Mutant (single edit) | Must fail | Cmd |
|---|---|---|---|
| MUT-CT-TOOMANY | `mint_request.go:142` `… \|\| errors.Is(err, store.ErrTooManyPendingRequests)` → `… \|\| false && errors.Is(…)` (measured surviving today, C18) | CT1 (`rc 2, want exitUsage`) | C `TestEpisodeStartRequestRefusesWhenSixteenArePending` |
| MUT-E1 | E1's `if` deleted | AP1 name arms | A `TestNewApproverRefusals` |
| MUT-E2 | E2's `if` deleted | AP1 webauthn (says "must be ssh") | A `TestNewApproverRefusals` |
| MUT-E2-ORDER | E2 moved after E4 | AP1 webauthn-empty-key | A `TestNewApproverRefusals` |
| MUT-E3 | E3's `if` deleted | AP1 format arms | A `TestNewApproverRefusals` |
| MUT-E4 | E4's `if` deleted | AP1 key-8193 (and key-empty) | A `TestNewApproverRefusals` |
| MUT-E4-EMPTY | E4 → `if len(keyFile) > MaxKeyFileBytes` | AP1 key-empty (`ErrAuthorizedKey` ≠) | A `TestNewApproverRefusals` |
| MUT-E4-GE | E4's `>` → `>=` | AP2 8192-byte boundary | A `TestNewApproverAcceptsTheRealFixtures` |
| MUT-E5 | E5's `if` deleted | AP1 private (`ErrApproverKeyFile` ≠) | A `TestNewApproverRefusals` |
| MUT-E5-TRIM | `bytes.TrimLeft(keyFile, " \t\r\n")` → `keyFile` | AP1 private-leading-space | A `TestNewApproverRefusals` |
| MUT-E6 | E6's `if` deleted | AP1 two-lines, blank-second-line (`ErrAuthorizedKey` ≠) | A `TestNewApproverRefusals` |
| MUT-E7 | E7's `if err != nil { return … }` deleted | AP1 rsa, p384, sk, one-field, bad-base64 (nil error) | A `TestNewApproverRefusals` |
| MUT-NA-NAME | the return's `Name: name` → `Name: ""` | AP2 | A `TestNewApproverAcceptsTheRealFixtures` |
| MUT-NA-NOW | the return's `EnrolledAt: now` deleted | AP2 | A `TestNewApproverAcceptsTheRealFixtures` |
| MUT-NA-NL | `line := bytes.TrimSuffix(keyFile, []byte("\n"))` → `line := keyFile` | AP2 with-newline arms | A `TestNewApproverAcceptsTheRealFixtures` |
| MUT-PAT-48 | `{0,47}` → `{0,48}` | AP4, AP1 name-49 | A `TestApproverNamePatternBoundaries` |
| MUT-PAT-FIRST | the first class `[a-z0-9]` → `[a-z0-9._-]` | AP4 `.a`, `_a`, `-a` | A `TestApproverNamePatternBoundaries` |
| MUT-PAT-ANCHOR | the trailing `$` deleted | AP4 `a b`, `a/b`, `a\n` | A `TestApproverNamePatternBoundaries` |
| MUT-R1 | R1's `if` deleted | AP3 revoked arm | A `TestCheckRevoke` |
| MUT-R2 | R2's `if` deleted | AP3 `(a, 999)` | A `TestCheckRevoke` |
| MUT-R2-LE | `now < a.EnrolledAt` → `now <= a.EnrolledAt` | AP3 `(a, 1000)` | A `TestCheckRevoke` |
| MUT-R-ORDER | R2 before R1 | AP3 revoked-and-behind | A `TestCheckRevoke` |
| MUT-RE-FP | `a.Key.Fingerprint()` → `""` in `RenderEnrolment` | AP5, RC4 `Key.Blob` | A `TestRenderEnrolmentGolden` |
| MUT-RE-R6 | one word of `EnrolmentWarning` changed | AP5 | A `TestRenderEnrolmentGolden` |
| MUT-RE-REC | the `records   approved_by=…` line dropped | AP5 | A `TestRenderEnrolmentGolden` |
| MUT-RE-TYPE | `a.Key.Type` → `"ssh-ed25519"` literal | RC4 `Key.Type` | A `TestEnrolmentPromptShowsEveryPersistedField` |
| MUT-RC4-FIELD | add `Comment string` to `Approver` | RC4 (path count) | A `TestEnrolmentPromptShowsEveryPersistedField` |
| MUT-RC4-KEYFIELD | add `Comment string` to `SSHPublicKey` | RC4 (element type walked) | A `TestEnrolmentPromptShowsEveryPersistedField` |
| MUT-DDL-SCHEMA | drop `CHECK (enrolled_at > 0)` in `credentials_schema.sql` only | AS1 | S `TestCredentialsSchemaMatchesCanonicalDDL` |
| MUT-DDL-CONST | the same in `approversSQL` only | AS1, AS5 | S `TestCredentialsSchemaMatchesCanonicalDDL` |
| MUT-UPG-NOV5 | `upgradeCredentials`' `if version < 5 { … approversSQL }` deleted | AS5 (every subtest: three tables) | S `TestCredentialsUpgradeToV5` |
| MUT-UPG-CLI | `OpenCredentials` → `openCredentialsWritable(path, true)` (the upgrade mode) | AS5 (a) | S `TestCredentialsUpgradeToV5` |
| MUT-CT-V4CASE | `credentialsTables`' `case version == 4:` arm deleted | AS5 v4 (a) (refused) | S `TestCredentialsUpgradeToV5` |
| MUT-TABLESET5 | `case version >= 5:` returns `credentialsTables(4)`'s three names | every fresh open (and AS7) | S `TestCredentialsV5WithoutApproversIsRefused` |
| MUT-OLD-FMT | `Error()`'s `e.Need` → `mintRequestsSchemaVersion` | AS6 (2) | S `TestOlderCredentialsStoreErrorText` |
| MUT-OLD-PUT | `PutMintRequest`'s site omits `Feature` | AS6 (1) PutMintRequest | S `TestOlderCredentialsStoreErrorText` |
| MUT-OLD-GET | `MintRequest`'s site omits `Feature` | AS6 (1) MintRequest | S `TestOlderCredentialsStoreErrorText` |
| MUT-OLD-LIST | `PendingMintRequests`' site omits `Feature` | AS6 (1) PendingMintRequests | S `TestOlderCredentialsStoreErrorText` |
| MUT-OLD-CONSUME | `MintSessionForRequest`'s site omits `Feature` | AS6 (1) MintSessionForRequest | S `TestOlderCredentialsStoreErrorText` |
| MUT-EN-DL | `EnrolApprover`'s `requireDeadline` deleted | deadline test | S `TestCredentialsMethodsRequireADeadline` |
| MUT-EN-OLD | `EnrolApprover`'s version check deleted | AS2 v4 (`no such table`, not `Older`) | S `TestEnrolApproverRefusals` |
| MUT-EN-NAME | the empty-name `if` deleted | AS2 `Name ""` (CHECK error, not `ErrBadApprover`) | S `TestEnrolApproverRefusals` |
| MUT-EN-FORMAT | the format `if` deleted | AS2 `gpg` (inserted) | S `TestEnrolApproverRefusals` |
| MUT-EN-KEY | the empty-key `if` deleted | AS2 `nil` key (untyped NOT NULL/CHECK error) | S `TestEnrolApproverRefusals` |
| MUT-EN-ENROLLED | the `EnrolledAt` `if` deleted | AS2 `0` and `-1` (untyped CHECK error) | S `TestEnrolApproverRefusals` |
| MUT-EN-ENROLLED-LT | `row.EnrolledAt <= 0` → `< 0` | AS2 `0` | S `TestEnrolApproverRefusals` |
| MUT-EN-LASTUSED | the `LastUsedAt` `if` deleted | AS2 `LastUsedAt 1500` (inserted; the INSERT ignores it, so the arm sees nil) | S `TestEnrolApproverRefusals` |
| MUT-EN-REVOKED | the `RevokedAt` `if` deleted | AS2 `RevokedAt 1500` | S `TestEnrolApproverRefusals` |
| MUT-EN-NAMETAKEN | the `if n != 0 { … ErrApproverNameTaken }` deleted | AS2 name taken (untyped PK error, C17) | S `TestEnrolApproverRefusals` |
| MUT-EN-KEYENROLLED | the `if err == nil { … ErrApproverKeyEnrolled }` deleted | AS2 key enrolled (untyped UNIQUE error, C17) | S `TestEnrolApproverRefusals` |
| MUT-EN-HOLDER | `it is approver %q's key` → without `holder` | AS2 `"mark-mac"` text | S `TestEnrolApproverRefusals` |
| MUT-LK-DL | `approverWhere`'s `requireDeadline` deleted | deadline test (`Approver`, `ApproverByKey`) | S `TestCredentialsMethodsRequireADeadline` |
| MUT-LK-OLD | `approverWhere`'s version check deleted | AS3 v4 | S `TestApproverLookups` |
| MUT-LK-NOROWS | the `errors.Is(err, sql.ErrNoRows)` branch deleted | AS3 miss arms (an error instead of `ok == false`) | S `TestApproverLookups` |
| MUT-LK-PREFIX | `public_key = ?` → `instr(public_key, ?) = 1` | AS3 prefix arm | S `TestApproverLookups` |
| MUT-SC-SWAP | `scanApprover` assigns `revoked.Int64` to `LastUsedAt` and `lastUsed.Int64` to `RevokedAt` | RC2 | S `TestApproverRowMatchesTheTable` |
| MUT-SC-DROPCOL | `RevokedAt` never assigned (`r.LastUsedAt = lastUsed.Int64` only) | RC2 (zero field) | S `TestApproverRowMatchesTheTable` |
| MUT-RC2-FIELD | add `Comment string` to `ApproverRow` | RC2 (field set) | S `TestApproverRowMatchesTheTable` |
| MUT-LS-DL | `ListApprovers`' `requireDeadline` deleted | deadline test | S `TestCredentialsMethodsRequireADeadline` |
| MUT-LS-OLD | `ListApprovers`' version check deleted | AS3b v4 | S `TestListApprovers` |
| MUT-LS-ORDER | `ORDER BY name` deleted | AS3b `[a b c]` (rowid order is b, a, c) | S `TestListApprovers` |
| MUT-LS-LIMIT | `LIMIT ?` (and its argument) deleted | AS3b 257 rows | S `TestListApprovers` |
| MUT-LS-ACTIVEONLY | `WHERE revoked_at IS NULL` added | AS3b `b` listed | S `TestListApprovers` |
| MUT-RV-DL | `RevokeApprover`'s `requireDeadline` deleted | deadline test | S `TestCredentialsMethodsRequireADeadline` |
| MUT-RV-OLD | `RevokeApprover`'s version check deleted | AS4 step 6 | S `TestRevokeApprover` |
| MUT-RV-REVOKED | the `UPDATE` drops `AND revoked_at IS NULL` | AS4 step 3 (nil; `RevokedAt` 1001) | S `TestRevokeApprover` |
| MUT-RV-ENROLLED | the `UPDATE` drops `AND enrolled_at <= ?` (and its argument) | AS4 step 1 (nil; and a CHECK error is also not the typed one) | S `TestRevokeApprover` |
| MUT-RV-NAME | the `UPDATE` drops `name = ? AND` (and its argument) | AS4 step 2 (`y` revoked) | S `TestRevokeApprover` |
| MUT-RV-N1 | the `if n != 1 { return classifyRevoke(…) }` deleted | AS4 steps 1, 3, 4, 5 (nil) | S `TestRevokeApprover` |
| MUT-RV-DELETE | the `UPDATE` → `DELETE FROM approvers WHERE name = ? AND revoked_at IS NULL AND enrolled_at <= ?` (arguments `name, now`) | AS4 step 2 (count 1, `Approver(x)` missing) | S `TestRevokeApprover` |
| MUT-CL-NOROWS | `classifyRevoke`'s `ErrNoRows` arm deleted (falls to the generic error) | AS4 step 5 | S `TestRevokeApprover` |
| MUT-CL-REVOKED | the `row.RevokedAt != 0` arm deleted | AS4 step 3 | S `TestRevokeApprover` |
| MUT-CL-BEHIND | the `row.EnrolledAt > now` arm deleted | AS4 step 1 | S `TestRevokeApprover` |
| MUT-CL-ORDER | the clock-behind arm placed before the revoked arm | AS4 step 4 | S `TestRevokeApprover` |

**Propagation-only branches (named, not counted, not mutated).** Every `if err != nil { return
…, fmt.Errorf("store: …: %w", err) }` after `BeginTx`, `QueryRowContext(…).Scan` of the count,
`ExecContext` of the `INSERT`/`UPDATE`, `RowsAffected`, `rows.Err`, `rows.Scan` and `Commit` in
`approvers.go`; `approverWhere`'s non-`ErrNoRows` error return; `EnrolApprover`'s
`!errors.Is(err, sql.ErrNoRows)` return; and `classifyRevoke`'s `err != nil` and `default` arms.
With the decision branches intact, none of these is reachable without fault injection into SQLite,
so a deletion is equivalent under every test here. The judge checks them by reading: each
propagates, none swallows. `scanApprover`'s `Scan` error return is the same.

**Known equivalent:** `Format: FormatSSH` → `Format: format` in `NewApprover`'s return (E3 has
already shown `format == FormatSSH`). Not counted.

---

## 5. Files (final `git status --short` after M2c)

| File | Sub-milestone | Change |
|---|---|---|
| `cmd/ailang-worldd/approve_test.go` | M1a | +CT1 |
| `host/approval/approvaltest/signer.go` | M1a | new (test support) |
| `host/approval/approvaltest/signer_test.go` | M1a | new: TS1, TS3 |
| `host/approval/approver.go` | M1b, M1c | new |
| `host/approval/approver_test.go` | M1b, M1c | new: AP0–AP5, RC4 |
| `host/approval/testdata/approvers/generate.sh` | M1b | new, mode 0755 |
| `host/approval/testdata/approvers/rsa3072.pub`, `p384.pub`, `SOURCE.txt` | M1b | new, generated text |
| `.gitattributes` | M1b | +1 line |
| `host/store/credentials_schema.sql` | M2a | +`approvers` |
| `host/store/credentials.go` | M2a | version 5, `credentialsTables`, upgrade |
| `host/store/mint_requests.go` | M2a | `OlderCredentialsStoreError{Need, Feature}`, 4 sites |
| `host/store/approvers.go` | M2a–M2c | new |
| `host/store/approvers_test.go` | M2a–M2c | new: AS2–AS6, AS3b, RC2, `makeV4Credentials` |
| `host/store/credentials_test.go` | M2a–M2c | AS1 edit, `canonicalApproversDDL`, AS7, deadline edit |
| `host/store/mint_requests_test.go` | M2a | `:132` table list |
| `design_docs/planned/w-world-approvals-by-key-a1b-sprint1-plan.md`, `design_docs/planned/sprint_w-world-approvals-by-key-a1b-s1.json` | planner | this plan |

---

## 6. Sprint gates (controller, after M2c; each was green at the base except the local race leg's load timeouts, §1.1a)

Run in **one** `bash` shell (`bash -c '…'` or a script), from the worktree:

```
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
./scripts/verify_ail.sh                                                    # no .ail change: exactly green
test -z "$(gofmt -l host/approval host/store cmd/ailang-worldd)"
go build ./... && go vet ./...
go test -run '^$' ./...                                                    # compile fence for every _test.go
go test ./... -count=1                                                     # the whole suite
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1         # no new root this sprint (unedited)
go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1
go test ./host/approval/... -count=1 && go test -race ./host/approval/... ./host/store/ -count=1
WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -run '^TestSSHSIGCrossCheckWithSSHKeygen$' -count=1 -v   # unchanged A1a cross-check: parent + 2 subtests PASS
go test ./cmd/ailang-worldd/ -count=1 -v -run '^(TestEpisodeStartRequest|TestApprove|TestSessionMint_|TestSessionNew|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard)' 2>&1 | tee /tmp/w375-cl.log | tail -3
if grep -E '^\s*--- (SKIP|FAIL)' /tmp/w375-cl.log; then exit 1; fi      # no silent skip (matchingHealth t.Skipf's without a tool pin); `! grep` would not trip set -e
# AFTER the controller commits (the tracked-binary gate reads committed blobs):
./scripts/verify_go.sh                                                     # build, go test ./..., evidence manifest, race leg, tracked_binary_gate.sh
git status --short                                                         # only the §5 files
```

**The race leg.** Locally `./scripts/verify_go.sh` was red at the base from three load timeouts
(§1.1a). If it is red after the sprint, read the failures: timeout-only reds in `host/broker`,
`host/daemon` or `host/verifygate` match the baseline; then `VERIFY_GO_LEGS=plain
./scripts/verify_go.sh` must be rc 0 (it holds the tracked-binary gate and `go test ./...`), and
CI's `go-race-verify` job is the race verdict. Any race report, any assertion failure, or any red in
another package is a regression.

**Must stay green, unedited**: every `TestApprove*`, `TestEpisodeStartRequest*` (other than the new
CT1), `TestReadTermLineWithin`, `TestExistingCallersReadUnbounded`, `TestMintRequestFieldMapIsComplete`,
`TestSSHSIG*` (incl. SG1's fixture pins), `TestPutMintRequest*`, `TestMintSessionForRequest*`,
`TestConcurrentApprovalsMintOnce`, `TestCredentialsV4WithoutMintRequestsIsRefused`,
`TestCredentialsUpgradeV1KeepsEveryRowUnboundAndIsIdempotent`, `TestCredentialsSessionVerbsLeaveAVersion1StoreAlone`,
`TestMigrateCredentials*`, `TestPreSplitStoreOpenModes`, `TestProductionContextRoots`,
`TestProductionGoSurface`, `TestDaemonDependencyAllowlist`, `TestEveryVerbHasHelp`,
`TestHelpIsNeverBlockedByTheRegistryKeyGuard`, and the QUICKSTART tests. The only edited existing
tests are `TestCredentialsSchemaMatchesCanonicalDDL`, `TestCredentialsUpgradeToV4` (one line) and
`TestCredentialsMethodsRequireADeadline`.

---

## 7. Executor directive notes

- **No git writes.** The executor never runs `git add`, `commit`, `stash`, `checkout`, `reset`,
  `restore`, `rebase`, `push` or `worktree`. Read-only `git status`/`git diff` are fine.
- **Snapshot after every sub-milestone, cumulatively.** After M1a's acceptance passes, and again
  after each later sub-milestone, copy **every** file created or modified since the base (not just
  this sub-milestone's) to `~/.ailang/state/world-iter375/snap/<M>/`, repo-relative paths preserved,
  modes preserved, plus a manifest:
  ```
  bash -c 'set -euo pipefail; M=M1a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter375-row201-a1b
    D=$HOME/.ailang/state/world-iter375/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/*-a1b-sprint1-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a1b-s1.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M1a`, `M1b`, `M1c`, `M2a`, `M2b`, `M2c`.) No file is ever deleted by this sprint, so the
  porcelain status has only `M` and `??` entries; stop and report if any other status appears.
  The controller rebuilds one commit per sub-milestone from the snapshots, in order.
- **Never snapshot a mutant.** Mutation runs (§4) are the controller's and the judge's, after the
  last snapshot; each restores and `cmp`s its file.
- **One `generate.sh` run.** If AP0's hashes are pinned and `generate.sh` is then re-run, the keys
  change and AP0 reds: re-pin, never hand-edit the `.pub` files.
- **Stop conditions.** A red earlier-milestone test, a `--- SKIP` in a `-v` run of
  `cmd/ailang-worldd`, any need to touch a file outside §5, a `ssh-keygen` prompt, or any reference
  to `~/.ssh` — stop and report; do not work around.
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` (the controller's shell is zsh:
  no word splitting, `echo` eats `\b`). Export `AILANG_BIN` in the same shell as Go tests of
  `cmd/ailang-worldd`/`host/daemon`.

---

## 8. Risks

- **R1 — Size.** ≈1,540 lines in six sub-milestones, each ≤ ~340. If the controller must cut
  further, M2c (list and revoke, ≈200) is the separable tail; it has no caller until M3.
- **R2 — The bump moves the version gates** (design R2/R3). From this sprint on, a v4 store refuses
  bound mints, preset starts and resumes from this build until `serve` of this build restarts and
  upgrades it (`credentials.go:432`, `session_new.go:379`, `episode_preset.go:235`). `episode start`
  already forces that restart through its same-build check. Sprint 1 alone costs the one restart;
  sprint 2 does not cost another. No test outside `host/store` builds a v4 credentials store (C6).
- **R3 — Schema 5 ships before any writer of it.** Until M3 lands, `approvers` is always empty and
  nothing reads it outside tests. That is the intended cut (design §0: enrolment is the root of
  trust and gets its own review).
- **R4 — Generated fixtures.** `rsa3072.pub`/`p384.pub` are made once by `generate.sh` (D-WORLD-109
  allows test keys made in a temp dir). Their private halves never exist outside `mktemp -d`.
- **R5 — `TestApprovalTestIsTestOnly` walks the repo.** It skips dot-directories, so local
  worktrees under `.claude/` or `.wt-*` are not scanned; on CI the checkout has none. Its two
  controls (≥300 files, ≥1 importing test) keep it from passing vacuously.
- **R6 — Haiku and reflection.** RC4's walker (struct recursion, three perturbation kinds) and
  RC2's field↔column map are the most intricate tests. §3 gives their exact algorithms; the
  controller should read them, not only run them.
- **R7 — Design premises P1/P2 are M3's.** They are recorded here so they are not lost; the sprint-2
  planner must carry Appendix A.

---

## 9. Acceptance criteria (the judge scores against these)

1. CT1 exists in `cmd/ailang-worldd/approve_test.go`, passes, and kills MUT-CT-TOOMANY (`exitFatal`
   instead of `exitUsage`).
2. `host/approval/approver.go` matches §2.1: `NewApprover` refuses in the order E1–E7, never echoes
   a private key file, and accepts the real ed25519 and P-256 fixtures with and without a trailing
   newline and at exactly 8 KiB; `CheckRevoke` refuses revoked before clock-behind;
   `RenderEnrolment` prints the §2.1 text byte for byte. `host/approval` imports only the standard
   library.
3. `approvaltest.Signer` produces SSHSIG signatures that `VerifySSHSIG`/`VerifyApproval` accept,
   deterministically, and no non-test Go file imports the package (TS3 with its two controls).
4. `testdata/approvers/` holds exactly the four generated text files, pinned by AP0, NUL-free,
   `generate.sh` executable, `.gitattributes` marks them `-text`; `testdata/sshsig/` unchanged.
5. The credentials store is schema 5 with exactly four tables; `approvers` matches §2.4 in a fresh
   store and in stores upgraded from 1, 3 and 4 (every existing row unchanged); a session verb never
   upgrades; a v5 store without `approvers` is refused unmodified.
6. `OlderCredentialsStoreError` carries `Need` and `Feature`; the mint-request text is byte-identical
   at all four sites; the approver text says `approvers need schema 5` with the restart hint.
7. `EnrolApprover`, `Approver`, `ApproverByKey`, `ListApprovers` and `RevokeApprover` exist with
   §2.7's signatures; each requires a deadline and refuses below v5; enrolment refuses a taken name
   and an enrolled key (revoked rows included) with typed errors and writes nothing; revoke marks,
   never deletes, and classifies unknown / revoked / clock-behind in that order.
8. RC2 and RC4 pass and are red under a new field in `ApproverRow`, `Approver` or `SSHPublicKey`.
9. Each §4 mutant, applied alone, makes its named test fail.
10. Every §6 gate passes, including `verify_go.sh` (with its tracked-binary gate) after the commit;
    no new context root; no test skips; `git status --short` lists only the §5 files.
11. No real key, real store, `~/.ssh` access, mint or publish anywhere in the sprint (D-WORLD-109);
    no bare `ssh-keygen`.

---

## Appendix A — M3 (`approver enroll|list|revoke`), carried to the head of sprint 2

Not executed in this sprint. Recorded so the sprint-2 planner starts from measured ground.

**Decisions the sprint-2 plan should carry (each measured or read at `5e2bf772`):**
- **A-D1 (P1)** One helper, reused by sprint 2's `--signature` (F9, 16 KiB):
  `readRegularFileCapped(path string, max int64) ([]byte, error)` opens with
  `os.OpenFile(path, os.O_RDONLY|syscall.O_NONBLOCK, 0)` (C19), `f.Stat()` the descriptor, refuses
  `!IsRegular()` (`%s is not a regular file`) and `Size() > max` (`%s is %d bytes, over the %d-byte bound`),
  then `readCapped(f, max)`; `readCapped(r io.Reader, max int64)` is `io.ReadAll(io.LimitReader(r,
  max+1))` refusing `len > max` (`grew past the %d-byte bound while being read`). The size-grew
  arm is a unit test of `readCapped` with a `max+1`-byte reader, deterministic, no file race. The
  3 s root is created after the read (it bounds nothing in a file read).
- **A-D2 (P2)** `readTermLineErrWithin(r io.Reader, d time.Duration, after func(time.Duration) <-chan time.Time) (line string, err error, ok bool)`:
  the goroutine sends `bufio.NewReader(r).ReadString('\n')`'s **raw** line and error on a channel
  of size 1. `confirmApproverAct` confirms only on `ok && err == nil && strings.HasSuffix(line, "\n") && isYes(line)`.
  `readTermLine`/`readTermLineWithin` are unchanged (R8).
- **A-D3** One `approverHelp` for the verb and all three subcommands; `ownHelpVerbs["approver"] =
  true`; `ownVerbHelp` `case "approver": return approverHelp, true`; three usage lines after
  `approve --list`. Measured: `usageLinesFor("approve")` matches `cmd == verb || HasPrefix(cmd,
  verb+" ")`, so `approver …` lines do not leak into `approve`'s generic help.
- **A-D4** Test seams for CA6/CA8's store errors: package vars `approverEnrol =
  (*store.Credentials).EnrolApprover` and `approverRevoke = (*store.Credentials).RevokeApprover`.
- **A-D5** A missing or schema-less credentials store at the read-only pre-check → `exitUsage`
  with the store's own text (`… serve or session new creates it`).
- **A-D6** Enrol order: flags (N1) → empty `--name`/`--pubkey` (N2) → `requireExistingStore` (N3) →
  `credentialsFor` → `readRegularFileCapped(…, approval.MaxKeyFileBytes)` (N4) → `NewApprover` (N5)
  → root 1 (3 s) + `OpenCredentialsReadOnly` → N6/N7/N8 → **close** → `confirmApproverAct` (N9a,
  N9, N10a, N10) → root 2 (3 s) + `openCredentialsForSession` (N11) → `approverEnrol` (N12/N13) →
  N14. Roots: `runApproverEnroll` 2, `runApproverList` 1, `runApproverRevoke` 2 (sprint 2's K8 makes
  it 3, design §6 M5 [r2]).
- **A-D7** CA9/CA10 need a terminal safe to poll while the verb writes: `lockedTerm{mu sync.Mutex;
  in io.Reader; out bytes.Buffer; closed bool}` with `Output() string` and `Closed() bool` under the
  mutex (`promptTerm`'s bare `bytes.Buffer` would race; `go test -race` would say so).
- **A-D8** The FIFO arm (CA5) creates the FIFO with `syscall.Mkfifo` in `t.TempDir()` and runs the
  verb under a 5 s guard; its mutants are "`O_NONBLOCK` dropped" (blocks → guard) and "`IsRegular`
  check dropped" (0 bytes read → `ErrApproverKeyFile`, so the arm's `not a regular file` text
  assertion fires) — not the design's "`Lstat` instead of `f.Stat()`" (P1).

**Sub-milestones and estimate (≈1,210 lines):** M3a — dispatch, `approverHelp`, `case "approver"`,
`ownHelpVerbs`/`ownVerbHelp`, usage, `confirmApproverAct`, `readTermLineErrWithin`,
`readRegularFileCapped`/`readCapped`; CA1, the help/keyguard edits, reader and file-helper unit
tests (≈170 + 170). M3b — `runApproverEnroll`; CA2, CA3, CA4 (≈100 + 160). M3c — CA5, CA6, CA9
(≈0 + 240). M3d — `runApproverList`, `runApproverRevoke`; CA7, CA8, CA10, `TestProductionContextRoots`
+3 entries (≈125 + 245).

**M3 mutants (one per branch; killer in brackets):** M1 `--addr` refusal [CA1]; D1 no subcommand
[CA1]; D2 unknown [CA1]; D3 `--help` [CA1]; N1 flag error/positional [CA5]; N2 empty name/pubkey
[CA5]; N3 missing world store [CA5]; N4-NONBLOCK [CA5 FIFO, guard]; N4-REGULAR [CA5 FIFO text];
N4-SIZE (fstat bound) [CA5 8 KiB+1 text]; N4-GREW (`readCapped` `len > max`) [`readCapped` unit
test]; N5 `NewApprover` error [CA5 private key, webauthn]; N5-ECHO (private bytes in stderr) [CA5];
N6 below v5 [CA5]; N7 name taken [CA5]; N8 key enrolled, names holder [CA5]; N9a nil `after` [CA9
arm 2]; N9 no terminal [CA3]; N10a `!ok` timeout [CA9]; N10 non-`y` [CA4]; N10-ERR `err == nil`
dropped [CA4 `y`+error arm]; N10-NL newline-suffix dropped [CA4 `y`+EOF arm]; N10-CLOSE `defer
term.Close()` dropped [CA4/CA9 close-recording terminal]; N10-UNBOUNDED (plain `readTermLine`)
[CA9 guard]; N10-BUDGET (≠120 s) [CA9 recorded duration]; N-LOCK (`openCredentialsForSession`
before the prompt) [CA9 `OpenCredentials` during the wait]; N-READONCE (file read twice) [CA2];
N11 writer busy → `exitFatal` [CA6]; N12 typed race → `exitUsage` [CA6]; N13 other → `exitFatal`
[CA6]; N14 the enrolled line [CA2]; L1 below v5 [CA7]; L2 `no approvers` [CA7]; L3 line fields
[CA7]; L-STATE (always `active`) [CA7]; K1 arity [CA8]; K2 unknown [CA8]; K3 `CheckRevoke` [CA8];
K4a nil `after` [CA10 arm 2]; K4 no terminal [CA8]; K5a timeout [CA10]; K5 non-`y`/EOF/error/partial
[CA8]; K-FENCE (revoke bypasses `confirmApproverAct`) [CA10 guard]; K6-TYPED [CA8]; K6-OTHER
[CA8]; K7 the revoked line [CA8]. **46 mutants.** K8/K9 belong with `DecidedMintRequests` (sprint 2's M4/M5).
