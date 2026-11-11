# Sprint plan — w-world-approvals-by-key-a3a, sprint 1 (iteration 384, row 201)

> **THE SPRINT — A3a sprint 1, "the record and the minter": M1 credentials schema 7
> `surface_approvals`, M2 the shared judge `host/authority/judge.go`, M3 `approve <request-id> --await
> [--await-for <duration>]`, exactly as the A3 child design r2 §9 lists them.** Nothing writes the new
> table outside tests; no route, page or `/v1` change; the daemon is not touched.
> - **Measured: 1,543 changed lines under `cmd/ host/`** — `git diff --numstat <base> HEAD -- cmd host`,
>   insertions + deletions, after a commit of each of the seven scratch stages (579 production: the
>   store 181, the judge 129, the CLI wrappers 107, the minter 162; 964 test). `docs/` and `design_docs/`
>   are outside that filter. Margin to the 2,000 cap: **457**. **The design's fallback cut does not
>   fire**: M3 stays in this sprint (Appendix C names what sprint 2 starts with).
> - **The judge is shared, and the CLI's texts did not move**: `approve --webauthn-assertion` and
>   `--signature`/`--sign-with` now call `authority.JudgeAssertion`/`JudgeSSHSignature` through thin
>   wrappers, and every refusal text assertion of the A1b/A2 CLI tests is unedited and green. The only
>   edits to existing `TestApprove*`/`TestApprover*`/`TestPasskey*` tests are one `DROP TABLE
>   surface_approvals` token in each of six downgrade setups (§1.4 P1: at schema 7 a store relabelled 4,
>   5 or 6 that still holds the new table is refused as "not a credentials store").
> - **The minter re-judges, never trusts**: `approve --await` polls read-only every second under a 3 s
>   root, judges the recorded artifact with the same judge against the row read in the same poll, and
>   only a verified approval reaches approve's steps 2–6 (`Signed` set, no terminal). AW7 inserts rows
>   the law refuses (an unenrolled key's valid signature, an assertion over another request, an
>   `approver` column naming `webauthn:mark-mac` over garbage): each is refused once, no terminal opens,
>   nothing mints.
> - **Every line of code and test is a unified diff in Appendix P**, generated LAST from compiled,
>   vetted, gofmt-clean, green scratch stages; the patches were then extracted **from this file** with
>   §2's `awk`, applied to a fresh `git archive 0a25c0d6` export, each sub-milestone's acceptance run,
>   and the result compared byte-for-byte with the mutant tree (§1.2 C17–C18). The executor applies
>   each patch and runs its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a3.md](w-world-approvals-by-key-a3.md), **r2** (quorum r1+r2 +
  controller carve-out + independent audit VALID 84; **not reopened**). D-WORLD-114 (no Tauri half in
  A3a; the CLI minter re-judges and mints) is OPEN with default A = keep both, which applies from this
  iteration: built as designed. Where this plan and the design differ, this plan was measured against
  `0a25c0d6` and wins; each difference is a decision (§1.3) or a premise found false (§1.4).
- **Starting ground**: A2 sprints 1–4 merged (sprint 4 = `88ff1565`); the A3 design r2 landed in PR
  #553 (`0a25c0d6`). The controller's V32 base-drift run at `0a25c0d6`: all 12 V1 sha256 prefixes equal
  the design's, `cmd/` and `host/` byte-identical to the design's base `0140bd7f`, V8, V18 and V15 hold,
  so the design's line numbers are current.
- **Planner base**: `0a25c0d6c6c9af6cdf81da75270f41af41a44bad` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter384-a3a`, branch
  `sprint/world-iter384-row201-a3a-s1`. Planner scratch (outside the repo, never under `/tmp`):
  `$HOME/.ailang/state/world-iter384/plan-scratch/` (`base/` the pristine export, `work/` the stage
  repo (tags `M1a` … `M3b`) and the tree the mutants were cut from, `patches/`, `tools/`, `logs/`,
  `roundtrip/`, `mut/` the mutant workers' copies).
- **Executor**: as routed by the controller (a small model). **Three runs**: Run 1 (M1a, M1b, M1c), Run
  2 (M2a, M2b), Run 3 (M3a, M3b). Each sub-milestone is one `patch -p1` (§2) plus its acceptance (§3).
  **Apply, never retype.**
- **Touched**: `host/store/{credentials.go, credentials_schema.sql, credentials_test.go,
  approvers_test.go, mint_requests_test.go}`, `cmd/ailang-worldd/{approver_rig_test.go,
  approver_label_test.go, approver_passkey_test.go, approve_passkey_test.go,
  approve_passkey_refusals_test.go, approve_sig_test.go}` (M1a); new `host/store/surface_approvals.go`,
  new `host/store/surface_approvals_test.go`, `host/store/credentials_test.go` (M1b);
  `host/store/mint_requests.go`, `surface_approvals_test.go` (M1c); new `host/authority/{judge.go,
  judge_test.go}` (M2a); `cmd/ailang-worldd/{approve_passkey.go, approve_sig.go}`, new `cmd/ailang-worldd/approve_passkey_now_test.go` (M2b);
  `cmd/ailang-worldd/{approve.go, main.go, help_test.go}`, new `cmd/ailang-worldd/{approve_await.go,
  approve_await_test.go}`, `host/store/context_roots_test.go` (M3a); new
  `cmd/ailang-worldd/approve_await_refusals_test.go` (M3b).
- **Not touched**: every `host/approval` file (the law stays pure and standard-library-only),
  `host/daemon/**` (the daemon keeps no writer and gains no route; `surfaceCSP` unchanged),
  `host/surfacepage/**`, `session_new.go` (`provisionEpisode` is reused as is), `world/*.ail`,
  `packages/**`, `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, `docs/**`, any skill, any `/v1`
  route, `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys and temp stores only (D-WORLD-109 = A).** Every store is under `t.TempDir()`; every key
   is `approvaltest.NewSigner`/`NewPasskey`'s fixed seed; every terminal is a test double
   (`promptTerm`, `approveEnv`); the minter's clock and timer are the `sessionEnv.now`/`after` fakes
   (`awaitEnv`). Nothing opens `~/.ailang/world`, `~/.ssh`, a real credentials store or `/dev/tty`;
   nothing mints into a real store, publishes, runs a browser or `ssh-keygen`. The "surface" in the
   minter's tests is `store.RecordSurfaceApproval` called by the test mid-wait: no route writes the
   table in this sprint (design §0 sprint boundaries).
2. **Pure core, effects at the boundary (coding standards S1–S3).** No `host/approval` file changes
   (`TestApprovalImportsOnlyTheStandardLibrary` unchanged). The judge does store reads only, on a
   read-only handle its caller opened and closes; the minter's effects are its read-only polls and
   `provisionEpisode`'s existing mint. No `.ail` (design §3: no invariant Z3 could state that is not
   already the law's). No new Go module, package, production exec site, `/v1` route; **one** new context
   root (`approve_await.go|pollSurfaceApproval`).
3. **One mutant per new production branch, anchored to the diff** (§4): every guard arm, bound,
   upgrade step, prune, judge branch, wrapper case, flag rule and minter branch of Appendix P's
   production hunks, each run against its single named test **and** its package suite, named
   equivalents with their reasons. Every run appends a line to `plan-scratch/logs/mutant_ledger.tsv`.
4. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
5. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and
   `go vet` green, the `_test.go` compile fence (`go test -run '^$' ./...`) green, its own tests and every
   earlier one green. **Nothing mints from the new table before the shared judge exists and is tested**
   (M3 lands after M2), and the table has no production writer at all in this sprint.
6. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `0a25c0d6` (measured by the planner, 2026-11-11)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`), `go1.26.6
darwin/arm64`, in a `git archive 0a25c0d6` export at `$HOME/.ailang/state/world-iter384/plan-scratch/base`
(a throwaway git repo; `diff -r cmd host` against the worktree: identical, `BASE_EQUALS_WORKTREE`; the
worktree itself was only read). The pre-commit legs are `tools/gates.sh <tree> <logs> base` = §6 without
the five new tests (they do not exist at the base). Logs: `plan-scratch/logs/baseline/`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `0a25c0d6`) |
| `./scripts/verify_ail.sh` | **0** | — | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `gofmt -l cmd host`; `go build ./...`; `go vet ./...`; `go test -run '^$' ./...` | **0** | — | empty / clean / no FAIL line |
| `go test -count=1 -timeout 20m ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/` | **0** | — | all `ok` (`cmd/ailang-worldd` 185.2 s, `host/daemon` 330.2 s, `host/store` 15.5 s) |
| `-race` legs (§6, the base's tests of the same patterns) | **0** | — | `ok`, no race report (`host/store` 6.5 s, `host/authority` 2.4 s, `cmd/ailang-worldd` 21.0 s) |
| §6's 28 named inventory, reflective, help and quickstart tests, `-v` | **0** | — | 28 `--- PASS`, 0 SKIP, 0 FAIL |
| the whole pre-commit script | **0** | **475 s** | `PRE_COMMIT_GATES_GREEN 475s` (load 10–17) |

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red
untouched package alone once and does not chase it. **Never start a gate under `nohup`.**
`VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` reads committed blobs and runs `go test ./...`; it is the
after-commit leg (§6) and was run on the planner's committed round-trip tree (§1.2 C19).

### 1.2 Codebase claims checked (at `0a25c0d6` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD origin/dev` (worktree) → `0a25c0d6c6c9af6cdf81da75270f41af41a44bad` twice; `git status --porcelain \| wc -l` → `0` |
| C2 | **Credentials schema is 6, and its table set is pinned exactly** | `grep -n "const credentialsSchemaVersion = \|case version >= 6\|if version < 6" host/store/credentials.go` → `43:const credentialsSchemaVersion = 6`, `74:	case version >= 6:`, `364:	if version < 6 {`; `credentialsTables(6)` = `approvers, episode_starts, mint_requests, session_credentials, webauthn_credentials` (`:75`) |
| C3 | **`enforceSchema` accepts a store only when its tables equal `credentialsTables(version)` exactly** | `host/store/credentials.go:323 case known && slices.Equal(tables, credentialsTables(version)):`, else `:339` "user_version %d with tables %v is not a credentials store (is it the world store?)" — the premise of §1.4 P1 |
| C4 | The existing schema pins | `grep -rn 'version != 6\|want 6 (row 201 A2)' --include='*_test.go' host cmd` → `host/store/credentials_test.go:115`, `:116` only; the exact table list `[]string{"approvers", "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"}` at `credentials_test.go:79`, `approvers_test.go:134`, `mint_requests_test.go:132` |
| C5 | **The judges live in `package main`** | `cmd/ailang-worldd/approve_passkey.go:1 package main`, `:75 func approveAssertion`; `approve_sig.go:1 package main`, `:42 func approveArmored` — the daemon cannot call them (design V9) |
| C6 | `host/authority` already imports the law and the store | `go list -f '{{join .Imports " "}}' ./host/authority \| tr ' ' '\n' \| grep sunholo` → `host/approval`, `host/broker`, `host/store`: `judge.go` adds no package edge |
| C7 | The minter's timer seam exists (design V29) | `sed -n 36,37p cmd/ailang-worldd/session.go` → `// after is time.After in production: the timer that bounds approve's y/N (row 201 §2.6a).` / `after func(time.Duration) <-chan time.Time` |
| C8 | Both artifact bounds are 16 KiB | `host/approval/passkey_json.go:15 const MaxPasskeyResponseBytes = 16 << 10`; `cmd/ailang-worldd/approve_sig.go:23 const approveMaxSignatureBytes = 16 << 10` → `MaxSurfaceApprovalArtifactBytes = 16 << 10`, the DDL's `BETWEEN 1 AND 16384` |
| C9 | **No test inventories `host/authority`'s files** | `grep -rln 'host/authority/' --include='*_test.go' host cmd \| wc -l` → `0`; control in the same call: `grep -rln 'host/approval/' --include='*_test.go' host cmd \| wc -l` → `15` |
| C10 | The context-root pin has no `--await` entry at the base | `grep -c 'approve_await\|runApproveAwait' host/store/context_roots_test.go` → `0`; control in the same call: `grep -c 'approve.go\|runApprove\|' host/store/context_roots_test.go` → `1` (`:20`, the existing `runApprove` root) |
| C11 | **A CLI refuses bound mints into a store below the current schema** | `cmd/ailang-worldd/session_new.go:380 if ro.Version() < store.CredentialsSchemaVersion {` → `UnboundCredentialsStoreError`; `episode_preset.go:235` the same. So after this sprint a schema-6 store refuses `session new`/`episode start`'s bound mint until `serve` restarts once (design R6, §1.4 P5) |
| C12 | **At schema 7 the six downgrade setups break eight CLI tests** (work, M1a without the six `DROP TABLE surface_approvals` tokens) | `go test ./cmd/ailang-worldd/ -run '^(TestApprovePasskeyRefusalsBeforeAnyEffect\|TestApproveByPasskey\|TestApproveBySignatureRefusalsBeforeAnyEffect\|TestApproverEnrollRefusals\|TestApproverListAndRevokeShowPasskeys\|TestPasskeyEnrollFlagRules\|TestApproverRevoke\|TestApproverList)$' -v` → **8 `--- FAIL`**, e.g. `approve_passkey_refusals_test.go:26: store: credentials store "…": user_version 5 with tables [approvers epi…` (`logs/p1_downgrade.log`); with the tokens: 8 `--- PASS` (M1a's acceptance) |
| C13 | The request's loopback daemon in the CLI rigs is the passkeys' origin | work: `TestApproveAwaitMintsARecordedPasskeyApproval` logged `mr.Daemon.Addr` = `127.0.0.1:7644`, and `pkOrigin` = `http://localhost:7644` (`approver_label_test.go:20`): so the minter prints one card address for both (S1-D6), and AW6 (no passkey) is the loopback line's own check |
| C14 | **Every repo-wide inventory, on the final tree** | the §6 script (`tools/gates.sh <tree> <logs> final` = §6's pre-commit block) on the round-trip tree (byte-identical to the mutant tree, C17) → `PRE_COMMIT_GATES_GREEN 409s`; pkgs: `host/store` ok 13.684s, `host/authority` ok 1.402s, `host/approval` ok 1.210s, `host/approval/approvaltest` ok 1.153s, `cmd/ailang-worldd` ok 129.880s, `host/daemon` ok 288.222s (`logs/final_gates/`): design §7's inventories by name — `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (no exec site added), `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestProductionContextRoots` (+1 root), `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestCredentialsMethodsRequireADeadline` (+2), `TestFrozenV1RouteTableMatchesMux` (13), `TestApprovalImportsOnlyTheStandardLibrary`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestSubprocessSinkGate`, `TestEveryVerbHasHelp`, the six quickstart tests, the schema pin `TestCredentialsSchemaMatchesCanonicalDDL` and the five new SA/JG tests — **33 named `--- PASS`, 0 SKIP/FAIL**; the AW leg 8 `--- PASS`. An earlier run of the same script on the tree before PK-NOW (`logs/final_gates_work_pre_pknow/`): `PRE_COMMIT_GATES_GREEN 517s`, 33 named PASS, `cmd/ailang-worldd` 129.5 s, `host/daemon` 328.3 s, `-race` `ok` ×3 |
| C15 | Every mutant, on the final tree | 104 rows, the last ledger line of each: **100 KILLED (named and suite), 4 named equivalents SURVIVED (named and suite)**, 0 not as expected; 157 ledger lines in all (§4; `logs/mutant_ledger.tsv`) |
| C16 | **A mutant run found a gap, fixed in the tests** | run1: `MUT-WP-JUDGE-NOW` (the passkey wrapper judging at the request's creation instead of `now`) SURVIVED named and suite — no passkey test enrolled an approver between the request and the approve. PK-NOW (`TestApproveByPasskeyJudgesTheApproverAtNow`, M2b, the twin of A1b's SIG-NOW) was added; run3: KILLED by it and by the suite. `MUT-POLL-SA-ERR` SURVIVED: a named equivalent (§4) |
| C17 | **The patches reproduce the mutant tree** | a fresh `git archive 0a25c0d6` export + `patch -p1` of Appendix P's seven patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r cmd host` against `plan-scratch/work` (the tree §4's mutants were cut from): **identical** (`ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`, `tools/roundtrip.sh`, `logs/roundtrip.log`); nothing uncommitted after M3b |
| C18 | **Every sub-milestone is bisectable** | the same fresh export, after **each** patch: `tools/accept.sh <tree> <M>` = §3 verbatim → `ACCEPT <M> rc=0` for all seven (M1a 539 s at load ~30: `host/store` 18.9 s, `cmd/ailang-worldd` 421.3 s, `host/daemon` 503.4 s; M1b 51 s; M1c 43 s; M2a 13 s; M2b 68 s; M3a 172 s; M3b 135 s; `logs/rt/`); the cumulative numstat after a commit of each = §7: **70, 404, 463, 839, 971, 1313, 1543** (`NUMSTAT` lines of `logs/roundtrip.log`). A first round trip (`logs/roundtrip_try1_stopped_for_pk_now.log`: M1a and M1b `ACCEPT rc=0`, numstat 70, 404) was stopped to add PK-NOW |
| C19 | **The whole repository on the final tree** (committed) | `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` on the committed round-trip tree at M3b → **rc 0**: `✓ 0 other binary blobs among 6751 tracked files`; `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`; `✓ all 64 required top-level evidence tests passed exactly once`; `✓ go gate PASSED (plain): build clean; legs run pass with pinned AILANG_BIN (AILANG v0.41.0` (`logs/final_verify_go.log`) |
| C20 | `go vet ./...`, gofmt, `-race`, `verify_ail.sh` on the final tree | same §6 run: `go vet ./...` empty; `gofmt -l cmd host` empty; `-race` `host/store` (credentials/surface/mint/passkey/approver/roots subset), `host/authority`, `cmd/ailang-worldd` approve subset → `ok host/store`, `ok host/authority`, `ok cmd/ailang-worldd`; `./scripts/verify_ail.sh` → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` (no `.ail` change) |

### 1.3 Planner decisions where the design leaves a choice

- **S1-D1 — The judge's steps.** `JudgeError{Step, RawID, Fingerprint, Approver, Err}` with `Step` one of
  `request` (the row's digest or challenge is not 64 hex: refused before any conversion; the CLI never
  reaches it after admission and F4, the daemon will after R8/R10), `parse`, `store` (the lookup's own
  error: below schema 6 or 5, the CLI's WF4/F12 restart text), `unknown-key`, `stored-key` and `law`.
  The design lists four (`parse|unknown-key|stored-key|law`); `request` and `store` are added so that no
  conversion can panic and so the CLI keeps its WF4/F12 texts. `RawID` (assertion) and `Fingerprint`
  (SSHSIG) carry what the CLI's unknown-key text names; `Approver` what its stored-key text names.
  `JudgeError` unwraps to `Err`, so `errors.Is(err, approval.ErrX)` holds through it.
- **S1-D2 — The wrappers open the read-only handle, then call the judge.** `approveAssertion` reads the
  file (WF2, unchanged), opens `OpenCredentialsReadOnly`, calls `JudgeAssertion`, closes, and maps the
  step to its old text (`passkeyRefusal`); `approveArmored` likewise (`sshRefusal`). The one order change
  — the open now precedes the parse — is not observable by any test: the row was read from the same
  store moments earlier in the same run (every WF/F text assertion green, unedited).
- **S1-D3 — JG1/JG2 judge the judge, the CLI tests judge the texts.** A `host/authority` test cannot
  import `package main`, so JG1/JG2 assert, for each CLI refusal arm (WF3–WF6, F11–F14), the judge's
  `Step`, its cause (`errors.Is`) and the field the CLI prints; the CLI's own texts stay pinned by
  `TestApprovePasskeyRefusalsBeforeAnyEffect` and `TestApproveBySignatureRefusalsBeforeAnyEffect`,
  whose assertions are unedited. Both judges also get `request-digest`, `request-challenge`,
  `law-consumed`, `law-bytes-changed` and `law-at-now-expired` arms, and a verified control that
  `Binds` the row.
- **S1-D4 — The context root is `approve_await.go|pollSurfaceApproval|Background: 1`** (design §2.6:
  "`approve.go|runApproveAwait|Background: 3`"). One root per poll, in the function that opens the
  handle, as every other approve root is pinned (one entry, count = the `context.Background()` calls in
  that function). The minter's start lines reuse `runApprove`'s existing 3 s root (`writeAwaitLinks`,
  as `--webauthn-link` does).
- **S1-D5 — `--await-for`** is a Go duration of at least 1s and only with `--await` (else the help and
  `exitUsage`, AW5); the deadline is `min(expires_at, start + whole seconds)`. `--await` is the sixth
  F1 mode; `--list --await` is F2.
- **S1-D6 — The start lines**: `waiting for an approval of request <id> (expires <RFC3339>, message
  hash <hash16>) from the surface; open its card and approve it there:` then one `  open
  <origin>/surface/approve/<id>` per distinct active passkey origin, then `http://localhost:<port>` of
  the request's own daemon address unless already printed (C13). Nothing secret, no fragment.
- **S1-D7 — M2 is admission.** Each poll re-reads the request and runs `approveAdmission` on it: used
  elsewhere, expired, clock behind or bytes ≠ digest end the wait with the CLI's own texts. So at the
  expiry the wait ends with `expiredRefusal`'s text, and M6's "no approval for request <id> by <time>;
  it stays pending until <expiry>" appears only at `--await-for`.
- **S1-D8 — M4**: `ailang-worldd approve: refused the surface's approval: <the CLI's text for that
  step>; waiting for another` on stderr, once per distinct `received_at`.
- **S1-D9 — Sprint-1 posture of the card address.** `/surface/approve/<id>` is served from sprint 3
  and nothing but tests writes `surface_approvals` until sprint 2's POST, so in production this sprint's
  `--await` waits until expiry or `--await-for`. This is design §0's sprint boundary (the table and the
  minter are reviewed before anything can reach them), the same posture as A2's links to
  `/surface/passkey` before its M7e.
- **S1-D10 — `received_at` is `RecordSurfaceApproval`'s `now`**; a caller's `ReceivedAt` is ignored.
  S2 also refuses an `approver` that is exactly `<format>:` (no name).
- **S1-D11 — Format dispatch**: `ssh` → `JudgeSSHSignature`, otherwise `JudgeAssertion` (the table's
  CHECK admits only `webauthn` and `ssh`).
- **S1-D12 — The prune is guarded by `version >= 7`**, so a schema 4–6 store (which a CLI never
  upgrades) still prunes and puts as before (SA3's v6 arm).
- **S1-D13 — Planner-added tests**: AW8 `TestApproveAwaitNeedsSchema7AndItsRequest` (a v6 store names
  the schema-7 restart at the first poll; a v5 store names the passkeys' restart before it; a request
  deleted mid-wait is named), and the judge arms of S1-D3. AW1 enrols a second passkey at the same
  origin (one address printed); AW6 checks the loopback address alone.
- **S1-D14 — Sub-milestones** (§1.5): M1a the table and its pins (with the six downgrade tokens, P1);
  M1b the record and read; M1c the prune; M2a the judge and JG1/JG2 (inert: nothing calls it yet);
  M2b the CLI wrappers (and PK-NOW, the passkey twin of A1b's SIG-NOW, added for the mutant MUT-WP-JUDGE-NOW); M3a the minter, its help, root and AW1/AW5; M3b AW2–AW4, AW6–AW8. M1b (334),
  M2a (376) and M3a (342) are above ~250: each is one production unit with the tests that pin it, and
  splitting them would land code without its tests.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §2.5/§6 "every existing `TestApprove*`/`TestPasskey*` stays unedited" and "must stay
  green, unedited: every a0a1, A1b and A2 test".** False at schema 7: six test setups make an older store
  by dropping tables from a current one and relabelling `user_version` (4, 5); the new table remains, and
  `enforceSchema` (C3) refuses the store as "not a credentials store" instead of naming the restart —
  eight tests red (C12). Routed: each of the six statement lists gains `DROP TABLE surface_approvals`
  first (one token per site); no assertion changes. **Files outside the design's §8**:
  `cmd/ailang-worldd/{approver_rig_test.go (approverStoreV4), approver_label_test.go,
  approver_passkey_test.go, approve_passkey_test.go, approve_passkey_refusals_test.go,
  approve_sig_test.go}`.
- **P2 — The schema pins** (design §6 "except the named edits (… the schema pins …)"):
  `TestCredentialsSchemaMatchesCanonicalDDL` (table list, `version != 6`), `TestCredentialsUpgradeToV4`
  and `TestCredentialsUpgradeToV6` (table list) are edited to schema 7 (C4); the new table's DDL is
  compared in SA1 (`host/store/{credentials_test.go, approvers_test.go, mint_requests_test.go}`).
- **P3 — Design §2.6's context-root entry** (`approve.go|runApproveAwait|Background: 3`) does not match
  the pin's grammar (one entry per function, counted calls); routed by S1-D4.
- **P4 — Design §2.5's four steps** cannot carry the CLI's WF4/F12 texts or refuse a non-hex row
  without a panic; routed by S1-D1.
- **P5 — Design R6 "schema 7 costs one `serve` restart"** is wider than the surface: C11's two checks
  compare with `CredentialsSchemaVersion`, so until `serve` restarts with this build a schema-6 store
  also refuses `session new` and `episode start`'s **bound** mints (`UnboundCredentialsStoreError`,
  which names the restart). Same as every earlier schema bump (A1b's 5, A2's 6); recorded for the
  controller, not routed.
- **P6 — Design §6 JG1 "the CLI wrapper's text agree"** cannot be asserted from `host/authority`;
  routed by S1-D3.

### 1.5 The size

Lines are insertions + deletions under `cmd/ host/`: **`git diff --numstat <base> HEAD -- cmd host`**
(`docs/` and `design_docs/` are outside the filter). **Every row is measured** from the scratch stages,
each committed in a throwaway repo (`plan-scratch/work`, tags `M1a` … `M3b`; `tools/build_patches.py`).

| Sub-milestone | Production | Test | Step | Cumulative vs `0a25c0d6` (after its commit) | Run |
|---|---|---|---|---|---|
| M1a — schema 7: the table, the upgrade, the pins, the six downgrade tokens | 48 | 22 | **70** | 70 | 1 |
| M1b — `RecordSurfaceApproval` (S1–S5), `SurfaceApproval`; SA1, SA2, SA4 | 123 | 211 | **334** | 404 | 1 |
| M1c — the prune extension; SA3 | 10 | 49 | **59** | 463 | 1 |
| M2a — `host/authority/judge.go`; JG1, JG2 (inert) | 129 | 247 | **376** | 839 | 2 |
| M2b — the CLI wrappers call the judge; PK-NOW | 107 | 25 | **132** | 971 | 2 |
| M3a — `approve --await`: flags, help, the minter, the root; AW1, AW5 | 162 | 180 | **342** | 1,313 | 3 |
| M3b — AW2, AW3, AW4, AW6, AW7, AW8 | 0 | 230 | **230** | 1,543 | 3 |
| **Sprint 1** | **579** | **964** | **1,543 ≤ 2,000** (margin 457) | | |

The design estimated sprint 1 at ≈ 1,650 (M1 ≈ 600, M2 ≈ 440, M3 ≈ 610); it measured **1,543**: M1 463,
M2 508, M3 572. The fallback cut (M3 to sprint 2) is not needed.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of
four backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has an
`@@ -0,0 …` hunk. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=M1a   # <- the sub-milestone, exactly as in the Appendix P heading: M1a M1b M1c M2a M2b M3a M3b
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter384-a3a
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a3a-sprint1-plan.md"
T="$HOME/.ailang/state/world-iter384/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an
  earlier sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with
  `export AILANG_BIN=$HOME/.pinned-ailang/ailang` and **`unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
  WORLD_PASSKEY_RECORD_DIR`**, each command's output `tee`'d to a log under
  `$HOME/.ailang/state/world-iter384/exec-logs/` (never only through `tail`).
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/store/ ./host/authority/ ./cmd/ailang-worldd/
go test -run '^$' ./...
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- FAIL` and no `--- SKIP`**;
each package run must end `ok` with no `FAIL` line. The planner's `tools/accept.sh <tree> <M> <logs>` is
exactly these lines (§1.2 C18).

### M1a — schema 7: the table and its pins · 48 prod + 22 test · Run 1

**Patch** `M1a`:
- `host/store/credentials_schema.sql`: the `surface_approvals` table (design §2.2), column for column.
- `host/store/credentials.go`: `credentialsSchemaVersion = 7`; `credentialsTables(7)` adds
  `surface_approvals` (sorted), `credentialsTables(6)` keeps the A2 set; `upgradeCredentials` gains `if
  version < 7 { surfaceApprovalsSQL }`.
- new `host/store/surface_approvals.go`: `surfaceApprovalsTable`, `surfaceApprovalsSchemaVersion = 7`,
  `surfaceApprovalsSQL` (the upgrade's DDL).
- `host/store/{credentials_test.go, approvers_test.go, mint_requests_test.go}`: the schema pins (P2).
- `cmd/ailang-worldd/{approver_rig_test.go, approver_label_test.go, approver_passkey_test.go,
  approve_passkey_test.go, approve_passkey_refusals_test.go, approve_sig_test.go}`: one `DROP TABLE
  surface_approvals` in each downgrade setup (P1).

**Acceptance** (~8–10 min: the whole `host/daemon` runs here, because `serve` upgrades the store)
```
[base]
go test ./host/store/ -run '^(TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsUpgradeToV4|TestCredentialsUpgradeToV6)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestApprovePasskeyRefusalsBeforeAnyEffect|TestApproveByPasskey|TestApproveBySignatureRefusalsBeforeAnyEffect|TestApproverEnrollRefusals|TestApproverListAndRevokeShowPasskeys|TestPasskeyEnrollFlagRules|TestApproverRevoke|TestApproverList)$' -count=1 -v
go test -count=1 -timeout 20m ./host/store/ ./cmd/ailang-worldd/ ./host/daemon/
```

**Expected cumulative numstat after the controller's commit of M1a: 70**.

### M1b — the record and the read · 123 prod + 211 test · Run 1

**Patch** `M1b`: `host/store/surface_approvals.go` gains `MaxSurfaceApprovalArtifactBytes` (16 KiB),
`ErrBadSurfaceApproval`, `SurfaceApprovalRow`, `RecordSurfaceApproval` (S1 deadline and `< 7` →
`OlderCredentialsStoreError{Feature: "surface approvals", Need: 7}`; S2 grammar; S3 unknown; S4
consumed / clock-behind / expired / digest-changed → the existing sentinels; S5 the upsert, `received_at
= now`) and `SurfaceApproval` (deadline, `< 7`, read-only, `ok=false` when none). New
`host/store/surface_approvals_test.go`: **SA1** `TestCredentialsSchema7SurfaceApprovals`, **SA2**
`TestRecordSurfaceApprovalGuards`; **SA4** `TestCredentialsMethodsRequireADeadline` +2 methods.

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestCredentialsSchema7SurfaceApprovals|TestRecordSurfaceApprovalGuards|TestCredentialsMethodsRequireADeadline|TestNoNewDeadlineFreeStoreReads)$' -count=1 -v
go test ./host/store/ -count=1
```

**Expected cumulative numstat after the controller's commit of M1b: 404**.

### M1c — the prune · 10 prod + 49 test · Run 1

**Patch** `M1c`: `host/store/mint_requests.go`: `PutMintRequest`'s prune also deletes
`surface_approvals` rows whose request is gone, in the same transaction, from schema 7 (S1-D12).
`surface_approvals_test.go`: **SA3** `TestSurfaceApprovalPrunedWithItsRequest`.

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestSurfaceApprovalPrunedWithItsRequest|TestRecordSurfaceApprovalGuards|TestPutMintRequestCapsAndPrunes|TestPutMintRequestNeedsVersion4)$' -count=1 -v
go test ./host/store/ -count=1
```

**Expected cumulative numstat after the controller's commit of M1c: 463**.

### M2a — the shared judge · 129 prod + 247 test · Run 2

**Patch** `M2a`: new `host/authority/judge.go` (`JudgeAssertion`, `JudgeSSHSignature`, `JudgeError`,
the step constants, `ErrNotEnrolled`; S1-D1). New `host/authority/judge_test.go`: **JG1**
`TestJudgeAssertionIsTheCLIJudge`, **JG2** `TestJudgeSSHSignatureIsTheCLIJudge` (S1-D3). Inert: nothing
calls the judge until M2b.

**Acceptance**
```
[base]
go test ./host/authority/ -run '^(TestJudgeAssertionIsTheCLIJudge|TestJudgeSSHSignatureIsTheCLIJudge)$' -count=1 -v
go test ./host/authority/ -count=1
```

**Expected cumulative numstat after the controller's commit of M2a: 839**.

### M2b — the CLI wrappers call the judge · 107 prod + 25 test · Run 2

**Patch** `M2b`: `cmd/ailang-worldd/approve_passkey.go`: `approveAssertion` = WF2's read, the read-only
open, `authority.JudgeAssertion`, `passkeyRefusal` (WF3–WF6's texts, byte-identical);
`cmd/ailang-worldd/approve_sig.go`: `approveArmored` = the open, `authority.JudgeSSHSignature`,
`sshRefusal` (F11–F14's texts). No existing test changes: the A1b/A2 CLI tests are the check. New
`cmd/ailang-worldd/approve_passkey_now_test.go`: **PK-NOW** `TestApproveByPasskeyJudgesTheApproverAtNow` (a
passkey enrolled after the request and before the approve approves it: the wrapper judges at `now`;
the passkey twin of A1b's `TestApproveBySignatureJudgesTheApproverAtNow`, added when the first mutant
run found MUT-WP-JUDGE-NOW alive, §4).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApprovePasskeyRefusalsBeforeAnyEffect|TestApproveByPasskey|TestApproveBySignatureRefusalsBeforeAnyEffect|TestApproveBySignatureMintsWithNoTerminal|TestApproveBySignatureJudgesTheApproverAtNow|TestApproveRevokeRaceIsRefusedAtTheWrite|TestApproveSignWithFakeSigner|TestApproveSignWithDeadline|TestApproveByPasskeyJudgesTheApproverAtNow)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestApprove|TestApprover|TestPasskeyEnroll)' -count=1
```

**Expected cumulative numstat after the controller's commit of M2b: 971**.

### M3a — `approve --await`, the minter · 162 prod + 180 test · Run 3

**Patch** `M3a`:
- `cmd/ailang-worldd/approve.go`: the usage line and the `--await` paragraph; `--await`,
  `--await-for`; F1's sixth mode; the `--await-for` rule (S1-D5); the `--await` branch after F4
  (`writeAwaitLinks` under the existing root, then `awaitSurfaceApproval`), which sets `signed` for
  steps 2–6.
- new `cmd/ailang-worldd/approve_await.go`: `writeAwaitLinks` (S1-D6), `pollSurfaceApproval` (one
  read-only open under its own 3 s root: the row, the recorded approval, the judgement; S1-D11),
  `awaitSurfaceApproval` (M1–M6; S1-D7, S1-D8).
- `cmd/ailang-worldd/main.go`: the top-level usage line. `cmd/ailang-worldd/help_test.go`
  (`TestEveryVerbHasHelp`): the line. `host/store/context_roots_test.go`: the new root (S1-D4).
- new `cmd/ailang-worldd/approve_await_test.go`: `awaitEnv`/`awaitClock` (the fake clock and timer),
  `recordSurface`, `requireMintedBy`; **AW1** `TestApproveAwaitMintsARecordedPasskeyApproval`, **AW5**
  `TestApproveAwaitFlagRules`.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveAwaitMintsARecordedPasskeyApproval|TestApproveAwaitFlagRules|TestEveryVerbHasHelp|TestQuickstartServiceSectionMatchesTheCLIHelp)$' -count=1 -v
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1 -v
go test ./cmd/ailang-worldd/ -count=1 -timeout 20m                                                     # ~3-5 min
```

**Expected cumulative numstat after the controller's commit of M3a: 1313**.

### M3b — the minter's refusals and stops · 230 test · Run 3

**Patch** `M3b`: new `cmd/ailang-worldd/approve_await_refusals_test.go`: **AW2**
`TestApproveAwaitRejudgesAndWaits`, **AW3** `TestApproveAwaitStopsAtExpiry` (expiry; `--await-for`),
**AW4** `TestApproveAwaitStopsWhenUsedElsewhere`, **AW6** `TestApproveAwaitSSHSIG`, **AW7**
`TestApproveAwaitNeverTrustsTheRow` (three arms), **AW8** `TestApproveAwaitNeedsSchema7AndItsRequest`
(S1-D13).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveAwaitMintsARecordedPasskeyApproval|TestApproveAwaitFlagRules|TestApproveAwaitRejudgesAndWaits|TestApproveAwaitStopsAtExpiry|TestApproveAwaitStopsWhenUsedElsewhere|TestApproveAwaitSSHSIG|TestApproveAwaitNeverTrustsTheRow|TestApproveAwaitNeedsSchema7AndItsRequest)$' -count=1 -v   # exactly 8 top-level --- PASS
go test ./cmd/ailang-worldd/ -count=1 -timeout 20m                                                     # ~3-5 min
```

**Expected cumulative numstat after the controller's commit of M3b: 1543**.

---

## 4. Mutation table

**104 rows run on the planner's final tree: 100 KILLED (named and suite), 4 named equivalents SURVIVED as predicted.**

**Derivation.** Anchored to Appendix P's production hunks: `credentials.go` (the version, both table sets, the upgrade step and its bound), the two DDL copies, `surface_approvals.go` (S1's deadline and version guards, every S2 arm and the artifact bound, S3, each S4 arm and its bound, S5's upsert, column and `received_at`, the read's deadline, version and no-row branches, the restart error's `Need` and `Feature`), `mint_requests.go` (the prune's guard, its bound and its statement), `judge.go` (the row-hex guards, the admission inputs, every step of both judges, every copied field the law reads), the two wrappers (each step's text and field, the default, the `errors.As` guard, the clock), `approve.go` (F1's sixth mode, each `--await-for` rule, the branch and its two refusals) and `approve_await.go` (the start lines' error, hash and both de-duplications, each poll branch and judgement input, each loop branch: the read error, admission, the verified approval, **AW7's "trust the row" mutant** `MUT-LOOP-TRUST-THE-ROW`, the print-once rule, the deadline's bound, guard and `min`, the poll interval). **Procedure**: `plan-scratch/tools/mutate.py <run> --workers 4` (outside the repo; the table is `tools/mutants.py`, and every anchor is asserted to occur exactly once, or at a named occurrence, in the final tree): for each row, in its own worker copy of the tree, apply **one** edit, build the package's test binary with `-vet=off` (a mutant that does not compile is recorded `NOCOMPILE`, never counted), run **the single named test** (`-test.run '^<test>$'`) **and then its package suite** — `host/store` and `host/authority` whole; for `cmd/ailang-worldd` the approve suite `-test.run '^(TestApprove|TestApprover|TestPasskeyEnroll|TestEveryVerbHasHelp|TestEpisodeStartRequest)'` (the whole package is ~2–3 min per mutant and holds the §6 load-sensitive tests, so its verdict would measure the machine) — restore, assert the restore is byte-identical, and **append one line to the ledger** `plan-scratch/logs/mutant_ledger.tsv` (UTC time, run label, id, stage, named, suite, expected, MATCH/MISMATCH). A row is KILLED only when both fail. **Never commit or snapshot a mutant.** In the table tabs are dropped, a newline prints as ⏎, a backquote as ' and a pipe as \|.

| Id | Stage | Mutant (single edit) | Named test | Named | Suite |
|---|---|---|---|---|---|
| MUT-SCHEMA-VERSION | M1a | credentials.go: const credentialsSchemaVersion = 7 → const credentialsSchemaVersion = 6 | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-TABLES-7 | M1a | credentials.go: credentialsTable, surfaceApprovalsTable, passkeysTable} → credentialsTable, passkeysTable} | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-TABLES-6 | M1a | credentials.go: case version == 6: → case version == 66: | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-UPGRADE-STEP | M1a | credentials.go: if version < 7 {⏎statements = append(statements, surfaceApprovalsSQL) → if false {⏎statements = append(statements, surfaceApprovalsSQL) | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-UPGRADE-BOUND | M1a | credentials.go: if version < 7 { → if version <= 5 { | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-DDL-UPGRADE | M1a | surface_approvals.go: artifact    BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384), → artifact    BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16385), | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-DDL-SCHEMA-FILE | M1a | credentials_schema.sql: format      TEXT NOT NULL CHECK (format IN ('webauthn', 'ssh')), → format      TEXT NOT NULL CHECK (format IN ('webauthn', 'ssh', 'tty')), | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-REC-DEADLINE | M1b | surface_approvals.go: if err := requireDeadline(ctx); err != nil { // S1 → if err := requireDeadline(ctx); err != nil && false { // S1 | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED |
| MUT-REC-VERSION | M1b | surface_approvals.go: if c.version < surfaceApprovalsSchemaVersion { // S1 → if false { // S1 | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-S2-ID | M1b | surface_approvals.go: case !MintRequestIDPattern.MatchString(row.RequestID): → case false: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S2-FORMAT | M1b | surface_approvals.go: case row.Format != "webauthn" && row.Format != "ssh": → case false: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S2-FORMAT-SSH | M1b | surface_approvals.go: case row.Format != "webauthn" && row.Format != "ssh": → case row.Format != "webauthn": | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S2-APPROVER-PREFIX | M1b | surface_approvals.go: case !strings.HasPrefix(row.Approver, row.Format+":") \|\| len(row.Approver) == len(row.Format)+1: → case !strings.HasPrefix(row.Approver, "") \|\| len(row.Approver) == len(row.Format)+1: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S2-APPROVER-NAME | M1b | surface_approvals.go: case !strings.HasPrefix(row.Approver, row.Format+":") \|\| len(row.Approver) == len(row.Format)+1: → case !strings.HasPrefix(row.Approver, row.Format+":"): | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S2-ARTIFACT-MIN | M1b | surface_approvals.go: case len(row.Artifact) < 1 \|\| len(row.Artifact) > MaxSurfaceApprovalArtifactBytes: → case len(row.Artifact) > MaxSurfaceApprovalArtifactBytes: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S2-ARTIFACT-MAX | M1b | surface_approvals.go: case len(row.Artifact) < 1 \|\| len(row.Artifact) > MaxSurfaceApprovalArtifactBytes: → case len(row.Artifact) < 1: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S2-ARTIFACT-CONST | M1b | surface_approvals.go: const MaxSurfaceApprovalArtifactBytes = 16 << 10 → const MaxSurfaceApprovalArtifactBytes = 16<<10 - 1 | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S2-DIGEST | M1b | surface_approvals.go: case !isLowerHex64(row.Digest): → case false: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S3-UNKNOWN | M1b | surface_approvals.go: if errors.Is(err, sql.ErrNoRows) { // S3 → if false { // S3 | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S4-CONSUMED | M1b | surface_approvals.go: case consumed:⏎ → case false:⏎ | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S4-CLOCK | M1b | surface_approvals.go: case now < createdAt: → case false: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S4-CLOCK-BOUND | M1b | surface_approvals.go: case now < createdAt: → case now <= createdAt: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S4-EXPIRED | M1b | surface_approvals.go: case now >= expAt: → case false: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S4-EXPIRED-BOUND | M1b | surface_approvals.go: case now >= expAt: → case now > expAt: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S4-DIGEST | M1b | surface_approvals.go: case digest != row.Digest: → case false: | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S5-UPSERT | M1b | surface_approvals.go: ON CONFLICT(request_id) DO UPDATE SET format = excluded.format, approver = excluded.approver, artifact = excluded.artifact, digest = excluded.digest, received_at = excluded.received_at → ON CONFLICT(request_id) DO NOTHING | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S5-ARTIFACT-COLUMN | M1b | surface_approvals.go: artifact = excluded.artifact, → artifact = surface_approvals.artifact, | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-S5-RECEIVED | M1b | surface_approvals.go: row.Digest, now); err != nil { // S5 → row.Digest, row.ReceivedAt); err != nil { // S5 | `TestRecordSurfaceApprovalGuards` | KILLED | KILLED |
| MUT-READ-DEADLINE | M1b | surface_approvals.go: if err := requireDeadline(ctx); err != nil {⏎return SurfaceApprovalRow{}, false, err → if err := requireDeadline(ctx); err != nil && false {⏎return SurfaceApprovalRow{}, false, err | `TestCredentialsMethodsRequireADeadline` | KILLED | KILLED |
| MUT-READ-VERSION | M1b | surface_approvals.go: if c.version < surfaceApprovalsSchemaVersion {⏎return SurfaceApprovalRow{}, false, c.olderForSurfaceApprovals() → if false {⏎return SurfaceApprovalRow{}, false, c.olderForSurfaceApprovals() | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-READ-NOROWS | M1b | surface_approvals.go: if errors.Is(err, sql.ErrNoRows) {⏎return SurfaceApprovalRow{}, false, nil → if false {⏎return SurfaceApprovalRow{}, false, nil | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-OLDER-NEED | M1b | surface_approvals.go: Need: surfaceApprovalsSchemaVersion, Feature → Need: passkeysSchemaVersion, Feature | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-OLDER-FEATURE | M1b | surface_approvals.go: Feature: "surface approvals"} → Feature: "surface approval"} | `TestCredentialsSchema7SurfaceApprovals` | KILLED | KILLED |
| MUT-PRUNE-GUARD | M1c | mint_requests.go: if c.version >= surfaceApprovalsSchemaVersion { // row 201 A3 → if true { // row 201 A3 | `TestSurfaceApprovalPrunedWithItsRequest` | KILLED | KILLED |
| MUT-PRUNE-GUARD-BOUND | M1c | mint_requests.go: if c.version >= surfaceApprovalsSchemaVersion { // row 201 A3 → if c.version > surfaceApprovalsSchemaVersion { // row 201 A3 | `TestSurfaceApprovalPrunedWithItsRequest` | KILLED | KILLED |
| MUT-PRUNE-STMT | M1c | mint_requests.go: NOT IN (SELECT request_id FROM mint_requests)' → NOT IN (SELECT request_id FROM mint_requests WHERE consumed_at IS NULL)' | `TestSurfaceApprovalPrunedWithItsRequest` | KILLED | KILLED |
| MUT-JG-DIGEST | M2a | judge.go: stored, err := hex.DecodeString(row.Digest)⏎if err != nil \|\| len(stored) != 32 { → stored, err := hex.DecodeString(row.Digest)⏎if false { | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JG-CHALLENGE | M2a | judge.go: challenge, err := hex.DecodeString(row.Challenge)⏎if err != nil \|\| len(challenge) != 32 { → challenge, err := hex.DecodeString(row.Challenge)⏎if false { | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-JG-NOW | M2a | judge.go: approval.AdmitInput{Now: now, → approval.AdmitInput{Now: row.CreatedAt, | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JG-CONSUMED | M2a | judge.go: Consumed: row.Consumed, → Consumed: false, | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-JG-BYTES | M2a | judge.go: BytesDigest: sha256.Sum256(row.Request)} → BytesDigest: [sha256.Size]byte(stored)} | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-PARSE | M2a | judge.go: resp, err := approval.ParseAssertionResponse(response)⏎if err != nil { → resp, err := approval.ParseAssertionResponse(response)⏎if err != nil && false { | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-STORE | M2a | judge.go: a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)⏎if err != nil { → a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)⏎if err != nil && false { | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-UNKNOWN | M2a | judge.go: if !ok {⏎return nil, &JudgeError{Step: JudgeStepUnknownKey, RawID → if !ok && false {⏎return nil, &JudgeError{Step: JudgeStepUnknownKey, RawID | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-RAWID | M2a | judge.go: RawID: resp.RawID, → RawID: nil, | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-STOREDKEY | M2a | judge.go: key, err := approval.ParseStoredCOSEKey(a.PublicKey)⏎if err != nil { → key, err := approval.ParseStoredCOSEKey(a.PublicKey)⏎if err != nil && false { | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-STOREDKEY-NAME | M2a | judge.go: Approver: a.Name, → Approver: "", | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-CHALLENGE | M2a | judge.go: Challenge: challenge,⏎Response:  response, → Challenge: [32]byte{0: challenge[0]},⏎Response:  response, | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-REVOKED | M2a | judge.go: RevokedAt: a.RevokedAt, → RevokedAt: 0, | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-ORIGIN | M2a | judge.go: Origin: p.Origin, → Origin: "http://localhost:7644", | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-COUNTER | M2a | judge.go: SignCount: uint32(p.SignCount)}}, → SignCount: 0}}, | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JA-LAW | M2a | judge.go: if err != nil {⏎return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎}⏎return signed, nil⏎}⏎⏎// JudgeSSHSignature → if err != nil && false {⏎return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎}⏎return signed, nil⏎}⏎⏎// JudgeSSHSignature | `TestJudgeAssertionIsTheCLIJudge` | KILLED | KILLED |
| MUT-JS-PARSE | M2a | judge.go: sig, err := approval.ParseSSHSIG(armored)⏎if err != nil { → sig, err := approval.ParseSSHSIG(armored)⏎if err != nil && sig == nil && false { | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-JS-STORE | M2a | judge.go: enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)⏎if err != nil { → enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)⏎if err != nil && false { | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-JS-UNKNOWN | M2a | judge.go: if !ok {⏎return nil, &JudgeError{Step: JudgeStepUnknownKey, Fingerprint → if !ok && false {⏎return nil, &JudgeError{Step: JudgeStepUnknownKey, Fingerprint | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-JS-FPR | M2a | judge.go: Fingerprint: sig.PublicKey.Fingerprint(), → Fingerprint: "", | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-JS-CHALLENGE | M2a | judge.go: Challenge: challenge,⏎Signature: armored, → Challenge: [32]byte{0: challenge[0]},⏎Signature: armored, | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-JS-FORMAT | M2a | judge.go: Approver: approval.Approver{Name: enrolled.Name, Format: enrolled.Format, → Approver: approval.Approver{Name: enrolled.Name, Format: "ssh", | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-JS-REVOKED | M2a | judge.go: RevokedAt: enrolled.RevokedAt}, → RevokedAt: 0}, | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-JS-LAW | M2a | judge.go: if err != nil {⏎return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎}⏎return signed, nil⏎}⏎ → if err != nil && false {⏎return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎}⏎return signed, nil⏎}⏎ (2nd occurrence) | `TestJudgeSSHSignatureIsTheCLIJudge` | KILLED | KILLED |
| MUT-WP-PARSE-TEXT | M2b | approve_passkey.go: "%s is not a passkey assertion: %v" → "%s is not an assertion: %v" | `TestApprovePasskeyRefusalsBeforeAnyEffect` | KILLED | KILLED |
| MUT-WP-UNKNOWN | M2b | approve_passkey.go: passkeyCred(je.RawID)) → passkeyCred(nil)) | `TestApprovePasskeyRefusalsBeforeAnyEffect` | KILLED | KILLED |
| MUT-WP-STOREDKEY | M2b | approve_passkey.go: see ailang-worldd approver list", je.Approver, je.Err) → see ailang-worldd approver list", "", je.Err) | `TestApprovePasskeyRefusalsBeforeAnyEffect` | KILLED | KILLED |
| MUT-WP-LAW | M2b | approve_passkey.go: "the passkey assertion does not approve request %s: %v", row.RequestID, je.Err) → "the passkey assertion does not approve request %s: %v", row.EpisodeID, je.Err) | `TestApproveAwaitRejudgesAndWaits` | KILLED | KILLED |
| MUT-WP-DEFAULT | M2b | approve_passkey.go: return je.Err.Error() // WF4 → return "refused" // WF4 | `TestApprovePasskeyRefusalsBeforeAnyEffect` | KILLED | KILLED |
| MUT-WP-NOT-JE | M2b | approve_passkey.go: if !errors.As(err, &je) { → if !errors.As(err, &je) && false { | `TestApprovePasskeyRefusalsBeforeAnyEffect` | SURVIVED | SURVIVED |
| MUT-WP-JUDGE-NOW | M2b | approve_passkey.go: authority.JudgeAssertion(ctx, ro, row, file, now) → authority.JudgeAssertion(ctx, ro, row, file, row.CreatedAt) | `TestApproveByPasskeyJudgesTheApproverAtNow` | KILLED | KILLED |
| MUT-WS-PARSE-TEXT | M2b | approve_sig.go: "%s is not an SSH signature: %v" → "%s is not a signature: %v" | `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED | KILLED |
| MUT-WS-UNKNOWN | M2b | approve_sig.go: see ailang-worldd approver list", je.Fingerprint) → see ailang-worldd approver list", "") | `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED | KILLED |
| MUT-WS-LAW | M2b | approve_sig.go: "the signature does not approve request %s: %v", row.RequestID, je.Err) → "the signature does not approve %s: %v", row.RequestID, je.Err) | `TestApproveSignWithFakeSigner` | KILLED | KILLED |
| MUT-WS-DEFAULT | M2b | approve_sig.go: return je.Err.Error() // F12 → return "refused" // F12 | `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED | KILLED |
| MUT-WS-NOT-JE | M2b | approve_sig.go: if !errors.As(err, &je) { → if !errors.As(err, &je) && false { | `TestApproveBySignatureRefusalsBeforeAnyEffect` | SURVIVED | SURVIVED |
| MUT-WS-JUDGE-NOW | M2b | approve_sig.go: authority.JudgeSSHSignature(ctx, ro, row, armored, now) → authority.JudgeSSHSignature(ctx, ro, row, armored, row.CreatedAt) | `TestApproveBySignatureJudgesTheApproverAtNow` | KILLED | KILLED |
| MUT-F1-AWAIT | M3a | approve.go: *assertion != "", *await} { → *assertion != ""} { | `TestApproveAwaitFlagRules` | KILLED | KILLED |
| MUT-AWAITFOR-NEEDS-AWAIT | M3a | approve.go: if awaitForSet && (!*await \|\| *awaitFor < time.Second) { → if awaitForSet && (*awaitFor < time.Second) { | `TestApproveAwaitFlagRules` | KILLED | KILLED |
| MUT-AWAITFOR-MIN | M3a | approve.go: if awaitForSet && (!*await \|\| *awaitFor < time.Second) { → if awaitForSet && (!*await) { | `TestApproveAwaitFlagRules` | KILLED | KILLED |
| MUT-AWAITFOR-BOUND | M3a | approve.go: *awaitFor < time.Second) { → *awaitFor < time.Millisecond) { | `TestApproveAwaitFlagRules` | KILLED | KILLED |
| MUT-AWAITFOR-SET | M3a | approve.go: f.Name == "await-for" }) → f.Name == "await" }) | `TestApproveAwaitFlagRules` | KILLED | KILLED |
| MUT-AWAIT-BRANCH | M3a | approve.go: if *await { // the minter → if false { // the minter | `TestApproveAwaitMintsARecordedPasskeyApproval` | KILLED | KILLED |
| MUT-AWAIT-LINKS-REFUSAL | M3a | approve.go: if refusal := writeAwaitLinks(ctx, credPath, stdout, row, mr); refusal != "" { → if refusal := writeAwaitLinks(ctx, credPath, stdout, row, mr); false && refusal != "" { | `TestApproveAwaitNeedsSchema7AndItsRequest` | KILLED | KILLED |
| MUT-AWAIT-REFUSAL | M3a | approve.go: if signed, refusal = awaitSurfaceApproval(env, credPath, row, *awaitFor, stderr); refusal != "" { → if signed, refusal = awaitSurfaceApproval(env, credPath, row, *awaitFor, stderr); false && refusal != "" { | `TestApproveAwaitStopsAtExpiry` | KILLED | KILLED |
| MUT-AL-ERR | M3a | approve_await.go: if err != nil { // below schema 6 names the serve restart → if err != nil && false { // below schema 6 names the serve restart | `TestApproveAwaitNeedsSchema7AndItsRequest` | KILLED | KILLED |
| MUT-AL-HASH | M3a | approve_await.go: hex.EncodeToString(h[:8])) → hex.EncodeToString(h[:7])) | `TestApproveAwaitMintsARecordedPasskeyApproval` | KILLED | KILLED |
| MUT-AL-ORIGIN-DEDUPE | M3a | approve_await.go: if !seen[e.Passkey.Origin] { → if true { | `TestApproveAwaitMintsARecordedPasskeyApproval` | KILLED | KILLED |
| MUT-AL-LOOPBACK | M3a | approve_await.go: net.SplitHostPort(mr.Daemon.Addr); err == nil && → net.SplitHostPort(mr.Daemon.Addr); err != nil && | `TestApproveAwaitSSHSIG` | KILLED | KILLED |
| MUT-AL-LOOPBACK-DEDUPE | M3a | approve_await.go: err == nil && !seen["http://localhost:"+port] { → err == nil { | `TestApproveAwaitMintsARecordedPasskeyApproval` | KILLED | KILLED |
| MUT-POLL-NOTFOUND | M3a | approve_await.go: if !ok {⏎return awaitPoll{}, fmt.Errorf → if !ok && false {⏎return awaitPoll{}, fmt.Errorf | `TestApproveAwaitNeedsSchema7AndItsRequest` | KILLED | KILLED |
| MUT-POLL-SA-ERR | M3a | approve_await.go: ro.SurfaceApproval(ctx, id); err != nil \|\| !p.recorded { → ro.SurfaceApproval(ctx, id); !p.recorded { | `TestApproveAwaitNeedsSchema7AndItsRequest` | SURVIVED | SURVIVED |
| MUT-POLL-FORMAT | M3a | approve_await.go: if p.rec.Format == "ssh" { → if p.rec.Format != "ssh" { | `TestApproveAwaitSSHSIG` | KILLED | KILLED |
| MUT-POLL-SSH-REFUSAL | M3a | approve_await.go: p.refusal = sshRefusal(jerr, "the surface's approval", p.row) → p.refusal = "refused" | `TestApproveAwaitNeverTrustsTheRow` | KILLED | KILLED |
| MUT-POLL-PK-REFUSAL | M3a | approve_await.go: p.refusal = passkeyRefusal(jerr, "the surface's approval", p.row) → p.refusal = "refused" | `TestApproveAwaitNeverTrustsTheRow` | KILLED | KILLED |
| MUT-POLL-JUDGE-NOW | M3a | approve_await.go: authority.JudgeAssertion(ctx, ro, p.row, p.rec.Artifact, now) → authority.JudgeAssertion(ctx, ro, p.row, p.rec.Artifact, now+3600) | `TestApproveAwaitMintsARecordedPasskeyApproval` | KILLED | KILLED |
| MUT-POLL-JUDGE-ARTIFACT | M3a | approve_await.go: authority.JudgeSSHSignature(ctx, ro, p.row, p.rec.Artifact, now) → authority.JudgeSSHSignature(ctx, ro, p.row, []byte(p.rec.Approver), now) | `TestApproveAwaitSSHSIG` | KILLED | KILLED |
| MUT-LOOP-ERR | M3a | approve_await.go: p, err := pollSurfaceApproval(credPath, row.RequestID, now)⏎if err != nil { → p, err := pollSurfaceApproval(credPath, row.RequestID, now)⏎if err != nil && false { | `TestApproveAwaitNeedsSchema7AndItsRequest` | KILLED | KILLED |
| MUT-LOOP-ADMISSION | M3a | approve_await.go: if _, refusal := approveAdmission(p.row, now); refusal != "" { // M2 → if _, refusal := approveAdmission(p.row, now); false && refusal != "" { // M2 | `TestApproveAwaitStopsWhenUsedElsewhere` | KILLED | KILLED |
| MUT-LOOP-TRUST-THE-ROW | M3a | approve_await.go: if p.signed != nil { // M5 → if p.recorded { // M5 | `TestApproveAwaitNeverTrustsTheRow` | KILLED | KILLED |
| MUT-LOOP-SIGNED | M3a | approve_await.go: if p.signed != nil { // M5 → if false { // M5 | `TestApproveAwaitMintsARecordedPasskeyApproval` | KILLED | KILLED |
| MUT-LOOP-PRINT-ONCE | M3a | approve_await.go: if p.recorded && p.rec.ReceivedAt != shown { // M4 → if p.recorded && (p.rec.ReceivedAt != shown \|\| true) { // M4 | `TestApproveAwaitRejudgesAndWaits` | KILLED | KILLED |
| MUT-LOOP-PRINT-RECORDED | M3a | approve_await.go: if p.recorded && p.rec.ReceivedAt != shown { // M4 → if p.rec.ReceivedAt != shown { // M4 | `TestApproveAwaitMintsARecordedPasskeyApproval` | KILLED | KILLED |
| MUT-LOOP-SHOWN | M3a | approve_await.go: shown = p.rec.ReceivedAt⏎ → _ = p.rec.ReceivedAt⏎ | `TestApproveAwaitRejudgesAndWaits` | KILLED | KILLED |
| MUT-LOOP-DEADLINE-CMP | M3a | approve_await.go: if now >= deadline { // M6 → if now > deadline { // M6 | `TestApproveAwaitStopsAtExpiry` | KILLED | KILLED |
| MUT-LOOP-AWAITFOR-GUARD | M3a | approve_await.go: if awaitFor > 0 { → if true { | `TestApproveAwaitMintsARecordedPasskeyApproval` | KILLED | KILLED |
| MUT-LOOP-DEADLINE-MIN | M3a | approve_await.go: deadline = min(deadline, env.now()+int64(awaitFor/time.Second)) → deadline = env.now() + int64(awaitFor/time.Second) | `TestApproveAwaitStopsAtExpiry` | SURVIVED | SURVIVED |
| MUT-LOOP-INTERVAL | M3a | approve_await.go: const awaitPollInterval = time.Second → const awaitPollInterval = 2 * time.Second | `TestApproveAwaitStopsAtExpiry` | KILLED | KILLED |

**Named equivalents** (each SURVIVED by the named test and the suite, as predicted):
- **MUT-WP-NOT-JE** — every error JudgeAssertion returns is a *JudgeError (each return site constructs one), so errors.As always succeeds and the guard is unreachable; it stays as the fallback for a future error kind.
- **MUT-WS-NOT-JE** — the same for JudgeSSHSignature.
- **MUT-POLL-SA-ERR** — SurfaceApproval returns ok = false with every error it returns, so the `!p.recorded` arm alone already returns that error (AW8's v6 arm is refused either way); the `err != nil` arm only states the contract.
- **MUT-LOOP-DEADLINE-MIN** — admission (M2, S1-D7) refuses at expires_at before the deadline check runs, so a deadline later than the expiry is never reached; min() only keeps the printed bound honest.

**Ledger.** `plan-scratch/logs/mutant_ledger.tsv`, 157 lines: run1-final 104, run2-final 50, run3-reanchored 3. The table shows each id's last line. **run1-final** ran all 104 on the tree before PK-NOW: 7 did not compile (an unused import or variable: MUT-S2-APPROVER-PREFIX, MUT-JG-BYTES, MUT-JA-UNKNOWN, MUT-JA-CHALLENGE, MUT-JS-UNKNOWN, MUT-JS-CHALLENGE, and in run2 MUT-POLL-NOTFOUND and MUT-LOOP-PRINT-ONCE) and were re-anchored to compiling edits of the same branch; **MUT-WP-JUDGE-NOW SURVIVED** (no passkey test enrolled the approver between the request and the approve), so PK-NOW was added to M2b; MUT-POLL-SA-ERR SURVIVED and is a named equivalent (below). **run2-final** re-ran every `cmd/ailang-worldd` mutant on the final tree (with PK-NOW) and the six re-anchored `host/store`/`host/authority` ones; **run3-reanchored** ran the last three (MUT-WP-JUDGE-NOW with PK-NOW as its named test). No `host/store` or `host/authority` file changed after run1 started, so their run1 lines are verdicts on the final code.

---

## 5. Files (final `git status --short` after M3b)

| File | Sub-milestone | Change |
|---|---|---|
| `host/store/credentials_schema.sql`, `credentials.go` | M1a | schema 7: the table, the version, the table set, the upgrade step |
| `host/store/surface_approvals.go` | M1a, M1b | new: the DDL (M1a); the row, the guards, the record and the read (M1b) |
| `host/store/credentials_test.go`, `approvers_test.go`, `mint_requests_test.go` | M1a (+M1b) | the schema pins (P2); SA4's two methods |
| `cmd/ailang-worldd/approver_rig_test.go`, `approver_label_test.go`, `approver_passkey_test.go`, `approve_passkey_test.go`, `approve_passkey_refusals_test.go`, `approve_sig_test.go` | M1a | one `DROP TABLE surface_approvals` token each (P1) |
| `host/store/surface_approvals_test.go` | M1b, M1c | new: SA1, SA2 (M1b); SA3 (M1c) |
| `host/store/mint_requests.go` | M1c | the prune extension |
| `host/authority/judge.go`, `judge_test.go` | M2a | new: the shared judge; JG1, JG2 |
| `cmd/ailang-worldd/approve_passkey.go`, `approve_sig.go`; `approve_passkey_now_test.go` | M2b | the wrappers; new PK-NOW |
| `cmd/ailang-worldd/approve.go`, `main.go`, `help_test.go`; `host/store/context_roots_test.go` | M3a | the mode, its help, the root |
| `cmd/ailang-worldd/approve_await.go`, `approve_await_test.go` | M3a | new: the minter; AW1, AW5 |
| `cmd/ailang-worldd/approve_await_refusals_test.go` | M3b | new: AW2–AW4, AW6–AW8 |
| this plan and `sprint_w-world-approvals-by-key-a3a-s1.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each measured at the base, §1.1)

In **one** `bash` shell, from the worktree, **not under `nohup`**. Every leg `tee`s its full output to a
log under `$HOME/.ailang/state/world-iter384/gates/`; a verdict is read from its log, never from `tail`
alone. No browser leg: nothing in this sprint serves or drives a page.

```
bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter384-a3a; cd "$W"
L=$HOME/.ailang/state/world-iter384/gates; mkdir -p "$L"
s0=$(date +%s)
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                      # no .ail change: exactly green
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                           # compile fence for every _test.go
go test -count=1 -timeout 20m ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
go test -race -count=1 ./host/store/ -run 'Credentials|SurfaceApproval|MintRequest|Passkey|Approver|ContextRoots' 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/authority/ 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestApprove|TestPasskeyEnroll|TestEveryVerbHasHelp)' 2>&1 | tee -a "$L/race.log"
# the repo-wide inventories (design §7), the schema pins, the reflective, help and quickstart tests and the new tests, by name, -v:
go test -v -count=1 ./host/broker/ -run '^(TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory|TestQuickstartDecisionPacketSection|TestPasskeyPageRoutes)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline|TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsSchema7SurfaceApprovals|TestRecordSurfaceApprovalGuards|TestSurfaceApprovalPrunedWithItsRequest)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run '^(TestEveryVerbHasHelp|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/authority/ -run '^(TestJudgeAssertionIsTheCLIJudge|TestJudgeSSHSignatureIsTheCLIJudge)$' 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract \
  TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestQuickstartDecisionPacketSection TestPasskeyPageRoutes \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL \
  TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly \
  TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI \
  TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp \
  TestCredentialsSchema7SurfaceApprovals TestRecordSurfaceApprovalGuards TestSurfaceApprovalPrunedWithItsRequest \
  TestJudgeAssertionIsTheCLIJudge TestJudgeSSHSignatureIsTheCLIJudge; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
go test -v -count=1 ./cmd/ailang-worldd/ -run '^TestApproveAwait' 2>&1 | tee "$L/aw.log"
test "$(grep -cE '^--- PASS: TestApproveAwait' "$L/aw.log")" = 8
if grep -E '^\s*--- (SKIP|FAIL)' "$L/aw.log"; then exit 1; fi
echo "PRE_COMMIT_GATES_GREEN $(( $(date +%s) - s0 ))s" | tee "$L/pre_commit.ok"
EOF
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter384-a3a && VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter384/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter384/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter384-a3a status --short   # only the two plan files, until the controller commits them
```

**At the base**: every pre-commit leg but the new tests is green at `0a25c0d6` in 475 s (§1.1; 28 named
tests `--- PASS`). On the final tree: `PRE_COMMIT_GATES_GREEN 409s`; pkgs: `host/store` ok 13.684s, `host/authority` ok 1.402s, `host/approval` ok 1.210s, `host/approval/approvaltest` ok 1.153s, `cmd/ailang-worldd` ok 129.880s, `host/daemon` ok 288.222s (33 named tests `--- PASS`, the AW leg 8 `--- PASS`, `-race` ok host/store, ok host/authority, ok cmd/ailang-worldd; §1.2 C14, C20), and the after-commit `verify_go.sh` **rc 0**: `✓ 0 other binary blobs among 6751 tracked files`; `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`; `✓ all 64 required top-level evidence tests passed exactly once`; `✓ go gate PASSED (plain): build clean; legs run pass with pinned AILANG_BIN (AILANG v0.41.0` (§1.2 C19). Re-run a red untouched package alone once; never under `nohup`.

**Known load- or environment-sensitive (re-run once, alone; do not chase)**:
`TestSeedRestartIgnoresASurvivingOldCopy` (host/broker), `TestHostsLaunchersExecHost`, `host/broker`'s
exec-probe timeouts and its local `-race`, the local `-race` of `host/daemon` and `host/verifygate`,
`TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
`TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck`'s supervisor tests,
`host/daemon`'s `TestAGUISlowReaderCutIsResumable`, `TestSessionSelfGit`,
`TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock`,
`cmd/ailang-worldd`'s `TestWorldContextAndStatusline` and `TestCLIRealSubprocessEpisode`, and anything
run under `nohup`. Any other red, any race report, any red in `host/store`, `host/authority`,
`host/approval`, or in `cmd/ailang-worldd`'s `TestApprove*`/`TestApprover*`/`TestPasskeyEnroll*`, is a
regression. CI's `go-race-verify` is the race verdict.

**Must stay green**: every a0a1, A1b and A2 test, with only the named edits: the six downgrade tokens
(P1), the three schema pins (P2), `TestCredentialsMethodsRequireADeadline` (+2), `TestEveryVerbHasHelp`
(+1 line) and `TestProductionContextRoots` (+1 root). In particular every WF2–WF6 and F9–F14 text
assertion, `TestFrozenV1RouteTableMatchesMux` (13 routes, untouched), `TestSurfaceAppRoutes` (untouched)
and every §7 inventory.

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **No browser, no `ssh-keygen`, no real store.** Nothing in this sprint needs Chrome or node; `go test`
  with `WORLD_CHROME` unset (§2) leaves the A2 browser tests as their counted skips, which only the
  whole-package runs include (not a §3 `-v` run).
- **Runs.** Run 1 = M1a, M1b, M1c; Run 2 = M2a, M2b; Run 3 = M3a, M3b. Each run reads §0, §2 and its own §3
  sections, then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an
  uncommitted tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter384-a3a && git diff --numstat 0a25c0d6c6c9af6cdf81da75270f41af41a44bad HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M1a **70**, M1b **404**, M1c **463**, M2a **839**, M2b **971**, M3a
  **1313**, M3b **1543**. **Any other value means a patch was not applied as written: stop and report.**
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M1a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter384-a3a
    D=$HOME/.ailang/state/world-iter384/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a3a-sprint1-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a3a-s1.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M1a` … `M3b`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any
  other status. The controller rebuilds one commit per sub-milestone (`$HOME/.ailang/state/world-iter384/commit_ms.sh`
  with `WT` and the snapshot root pointed here).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- FAIL`, or a `--- SKIP` in a §3 `-v` run; any need to touch a file
  outside §5; running Chrome, node, `ssh-keygen` or any browser; any reference to `~/.ssh`,
  `~/.ailang/world` or a real `/dev/tty`. Stop and report; do not work around.
- **Timing.** M1a's acceptance ~8–10 min (whole `host/store`, `cmd/ailang-worldd` and `host/daemon`);
  M3a's and M3b's whole `cmd/ailang-worldd` ~3–5 min each; every other acceptance under 1 min plus the
  compile fence (~25–60 s).
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does
  not word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a
  file; never pipe a gate's only copy through `tail`. **Never start a gate under `nohup`.** Poll any
  background command's artifact in a bounded `date +%s` loop in the same turn; never end a turn
  announcing a wait.

---

## 8. Risks

- **R1 — The schema bump's reach (P5).** A schema-6 store refuses `RecordSurfaceApproval` (by design)
  and also bound mints, until `serve` restarts with this build; every refusal names the restart. No
  daemon code changes: `ensureCredentialsStore` upgrades any store below `CredentialsSchemaVersion`.
- **R2 — The minter waits for nothing in production this sprint** (S1-D9): no route writes
  `surface_approvals` until sprint 2 and the card page lands in sprint 3. `--await` is correct and
  bounded (expiry, `--await-for`), just not yet reachable from a browser.
- **R3 — A refused artifact is replaced, not appended** (S5): a later verified POST replaces an earlier
  refused one while the request is pending; the minter prints each refusal once (per `received_at`).
- **R4 — Load-sensitive whole-package runs** (§6's list): the M1a acceptance runs the whole
  `host/daemon` (330 s at the base); re-run a red untouched test alone once.
- **R5 — Haiku and patches.** The executor transcribes nothing; a skipped patch fails `--dry-run` or the
  numstat; a hand "fix" after a red test is forbidden by §7.

---

## 9. Acceptance criteria (the judge scores against these)

1. **Schema 7**: a fresh credentials store is user_version 7 with exactly `approvers, episode_starts,
   mint_requests, session_credentials, surface_approvals, webauthn_credentials`; `surface_approvals` in
   `credentials_schema.sql` and in `surfaceApprovalsSQL` equal design §2.2 column for column; a schema-6
   store is upgraded by `UpgradeCredentials` only, once, its rows kept; a session open never upgrades
   (SA1).
2. **`RecordSurfaceApproval`**: each of S1 (no deadline; below 7 → `OlderCredentialsStoreError` naming
   schema 7 and the restart), S2 (id, format, `approver` = `<format>:<name>`, 1..16384-byte artifact,
   64 lowercase hex digest), S3 (unknown) and S4 (consumed, clock behind, expired, digest changed)
   refuses with its sentinel and writes **0 rows**; the bounds are admitted; S5 replaces a pending
   request's row and sets `received_at = now` (SA2). `SurfaceApproval` is deadline-required, refuses
   below 7, and reads back exactly what was recorded.
3. **The prune**: `PutMintRequest` deletes the recorded approval of every request it prunes, in the
   same transaction, keeps a consumed request's, and still works on a schema-6 store (SA3).
4. **One judge**: `host/authority/judge.go` exports `JudgeAssertion` and `JudgeSSHSignature` with the
   design §2.5 signatures; for every CLI refusal (WF3–WF6, F11–F14) the step, cause and printed field
   are pinned (JG1, JG2); a verified artifact `Binds` the row.
5. **The CLI's texts did not move**: `approve --webauthn-assertion`, `--signature` and `--sign-with` call
   the shared judge at approve's own `now` (PK-NOW, SIG-NOW); every A1b/A2 refusal-text assertion is unedited and green (the only edits to those
   tests are P1's six `DROP TABLE surface_approvals` tokens).
6. **`approve --await`** mints a recorded, verified passkey approval (AW1) or SSHSIG (AW6) with **no
   terminal**, `approved_by` `webauthn:<name>`/`ssh:<name>`, and the effects the typed `y` makes; it
   prints the request, its message hash and the card's address per distinct origin plus the loopback
   daemon (no secret).
7. **The minter never trusts the row** (AW7): an unenrolled key's valid signature, an assertion over
   another request, and an `approver` column naming `webauthn:mark-mac` over garbage are each judged,
   refused once, and mint nothing; no terminal opens.
8. **The wait**: a refused artifact is printed once per `received_at` and the wait goes on; a replacement
   mints (AW2); the wait ends exactly at the expiry (admission's text) or at `--await-for` (exact poll
   count and time), `exitUsage`, nothing minted (AW3); a consume elsewhere ends it with "already used"
   (AW4); a store below 7 names the restart, a deleted request is named (AW8).
9. **Flags**: `--await` is F1's sixth exclusive mode; `--await-for` only with `--await` and ≥ 1s; `--list
   --await` refused; each prints the help and exits `exitUsage` before any store read (AW5); the help
   and the top-level usage carry the line (`TestEveryVerbHasHelp`).
10. Each §4 mutant, applied alone, makes its named test **and** its package suite fail; the named
    equivalents survive for the stated reasons; the ledger has a line per run.
11. Every §6 gate passes, including every inventory, schema pin, reflective, help and quickstart test by
    name, the `-race` legs, and `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after the commit (modulo
    §6's re-run rule); `git status --short` lists only §5.
12. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000.
13. `host/approval` unchanged (pure, standard-library-only); `host/daemon` and `host/surfacepage`
    unchanged (no route, no writer kept, `surfaceCSP` unchanged); no new production exec site, module,
    package, `/v1` route; exactly one new context root.
14. No real key, store, `/dev/tty`, `~/.ssh` access, browser, mint, publish or `ssh-keygen` anywhere in
    the sprint (D-WORLD-109).

---

## Appendix C — after this sprint

- **Sprint 2, "the routes" (M4–M5), unchanged from design §9**: the surface origin set `O` and `serve
  --surface-origin`, the read routes, the guards, `POST /surface/approve/{id}` (which calls this
  sprint's `authority.JudgeAssertion`/`JudgeSSHSignature` and `store.RecordSurfaceApproval` through a
  per-POST non-waiting writer), RT12 `TestDaemonCredentialWriterSites`, and E2E1, which drives this
  sprint's `approve --await` end to end. Nothing of M3 moves (the fallback cut did not fire).
- **Carry into sprint 2's plan**: the daemon's R8/R10 run before the judge, so `JudgeStepRequest` stays
  unreachable from the routes too; the route maps `JudgeError.Step` to its 422 reason (the CLI's
  `passkeyRefusal`/`sshRefusal` are `package main` and cannot be shared — a pure text table in
  `host/authority` or `host/surfacepage` is the sprint-2 planner's choice).
- **D-WORLD-114** stays open (default A applied); if Mark vetoes (b), the minter's mint location returns
  to design.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M1a — schema 7: the table and its pins

````diff M1a
--- a/cmd/ailang-worldd/approve_passkey_refusals_test.go
+++ b/cmd/ailang-worldd/approve_passkey_refusals_test.go
@@ -56,7 +56,7 @@ func TestApprovePasskeyRefusalsBeforeAnyEffect(t *testing.T) {
 		{name: "WF2-missing", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "absent.json") }, want: "--webauthn-assertion: open "},
 		{name: "WF3-not-an-assertion", file: func(*testing.T, approveRig) []byte { return []byte(`{"id":"x"}`) },
 			want: "is not a passkey assertion: " + approval.ErrPasskeyResponse.Error()},
-		{name: "WF4-v5-store", setup: sql(`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
+		{name: "WF4-v5-store", setup: sql(`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
 			want: "passkeys need schema 6"},
 		{name: "WF5-unknown-credential", file: func(t *testing.T, a approveRig) []byte {
 			return approvaltest.NewPasskey(approval.AlgEdDSA, 55).Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
--- a/cmd/ailang-worldd/approve_passkey_test.go
+++ b/cmd/ailang-worldd/approve_passkey_test.go
@@ -96,7 +96,7 @@ func TestApproveByPasskey(t *testing.T) {
 			t.Fatalf("rc %d, terminal %v, daemon calls %d, stderr %q", code, opened, *calls, stderr.String())
 		}
 		raw := rawCredentials(t, a.credPath)
-		for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+		for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 			if _, err := raw.Exec(stmt); err != nil {
 				t.Fatal(err)
 			}
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -216,7 +216,7 @@ func TestApproveBySignatureRefusalsBeforeAnyEffect(t *testing.T) {
 		{name: "v4-store", setup: func(t *testing.T, a approveRig) {
 			raw := rawCredentials(t, a.credPath)
 			defer func() { _ = raw.Close() }()
-			for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+			for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 				if _, err := raw.Exec(stmt); err != nil {
 					t.Fatal(err)
 				}
--- a/cmd/ailang-worldd/approver_label_test.go
+++ b/cmd/ailang-worldd/approver_label_test.go
@@ -82,7 +82,7 @@ func TestApproverListAndRevokeShowPasskeys(t *testing.T) {
 	db5, creds5 := approverStore(t)
 	seedApprover(t, creds5, "s", fixtureKeyOf(t, "ed25519.pub"), approveT0)
 	raw = rawCredentials(t, creds5)
-	for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+	for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 		if _, err := raw.Exec(stmt); err != nil {
 			t.Fatal(err)
 		}
--- a/cmd/ailang-worldd/approver_passkey_test.go
+++ b/cmd/ailang-worldd/approver_passkey_test.go
@@ -154,7 +154,7 @@ func TestPasskeyEnrollFlagRules(t *testing.T) {
 		{name: "N6-v5", setup: func(t *testing.T, r enrolRig) {
 			raw := rawCredentials(t, r.creds)
 			defer func() { _ = raw.Close() }()
-			for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+			for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 				if _, err := raw.Exec(stmt); err != nil {
 					t.Fatal(err)
 				}
--- a/cmd/ailang-worldd/approver_rig_test.go
+++ b/cmd/ailang-worldd/approver_rig_test.go
@@ -64,7 +64,7 @@ func approverStoreV4(t *testing.T) (db, creds string) {
 	db, creds = approverStore(t)
 	raw := rawCredentials(t, creds)
 	defer func() { _ = raw.Close() }()
-	for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+	for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 		if _, err := raw.Exec(stmt); err != nil {
 			t.Fatal(err)
 		}
--- a/host/store/approvers_test.go
+++ b/host/store/approvers_test.go
@@ -131,7 +131,7 @@ func TestCredentialsUpgradeToV6(t *testing.T) {
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["approvers"]), normalizeDDL(canonicalApproversDDL); got != want {
 				t.Fatalf("upgraded approvers DDL:\n got: %s\nwant: %s", got, want)
 			}
--- a/host/store/credentials.go
+++ b/host/store/credentials.go
@@ -39,8 +39,8 @@ var credentialsSchemaSQL string
 // independent of the world store's currentSchemaVersion. Version 2 added
 // profile bindings; version 3 adds immutable attended episode start history;
 // version 4 adds mint requests (row 201); version 5 adds approvers (row 201 A1b); version 6 adds
-// passkeys (row 201 A2).
-const credentialsSchemaVersion = 6
+// passkeys (row 201 A2); version 7 adds the surface's recorded approvals (row 201 A3).
+const credentialsSchemaVersion = 7
 
 // CredentialsSchemaVersion exports credentialsSchemaVersion: serve compares a
 // store's version with it to decide whether to upgrade.
@@ -71,7 +71,9 @@ const credentialsTable = "session_credentials"
 // credentialsTables is the exact, sorted application table set of a credentials store at version.
 func credentialsTables(version int) []string {
 	switch {
-	case version >= 6:
+	case version >= 7:
+		return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable, surfaceApprovalsTable, passkeysTable}
+	case version == 6:
 		return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable, passkeysTable}
 	case version == 5:
 		return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable}
@@ -364,6 +366,9 @@ func upgradeCredentials(display string, db *sql.DB, version int) error {
 	if version < 6 {
 		statements = append(statements, passkeysSQL)
 	}
+	if version < 7 {
+		statements = append(statements, surfaceApprovalsSQL)
+	}
 	for _, stmt := range statements {
 		if _, err := tx.Exec(stmt); err != nil {
 			return fmt.Errorf("store: upgrade credentials %q: %w", display, err)
--- a/host/store/credentials_schema.sql
+++ b/host/store/credentials_schema.sql
@@ -83,3 +83,16 @@ CREATE TABLE IF NOT EXISTS webauthn_credentials (
     backup_eligible INTEGER NOT NULL CHECK (backup_eligible IN (0, 1)),
     sign_count      INTEGER NOT NULL CHECK (sign_count BETWEEN 0 AND 4294967295)
 );
+
+-- Approvals POSTed to the surface (version 7, row 201 A3). One row per request: the artifact the
+-- daemon verified, kept for the minter, which re-judges it. Never a credential, never trusted.
+-- A consumed request's row stays as the record of which artifact approved it; a row whose request
+-- was pruned is pruned with it.
+CREATE TABLE IF NOT EXISTS surface_approvals (
+    request_id  TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'mr-[0-9a-f]*'),
+    format      TEXT NOT NULL CHECK (format IN ('webauthn', 'ssh')),
+    approver    TEXT NOT NULL CHECK (approver <> ''),
+    artifact    BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
+    digest      TEXT NOT NULL CHECK (length(digest) = 64),
+    received_at INTEGER NOT NULL
+);
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -76,7 +76,7 @@ var canonicalPasskeysDDL = `CREATE TABLE webauthn_credentials (
 func TestCredentialsSchemaMatchesCanonicalDDL(t *testing.T) {
 	c := openTestCredentials(t)
 	ddl := tableDDL(t, c.db)
-	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"})
+	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
 	for name, got := range map[string]string{"schema file": ddl["webauthn_credentials"], "passkeysSQL": passkeysSQL} {
 		if normalizeDDL(got) != normalizeDDL(canonicalPasskeysDDL) {
 			t.Fatalf("webauthn_credentials DDL (%s):\n got: %s\nwant: %s", name, got, canonicalPasskeysDDL)
@@ -112,8 +112,8 @@ func TestCredentialsSchemaMatchesCanonicalDDL(t *testing.T) {
 		t.Fatalf("v4 world table differs from the version-1 credentials table: %s", got)
 	}
 	var version int
-	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != 6 {
-		t.Fatalf("credentials user_version = %d %v, want 6 (row 201 A2)", version, err)
+	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != 7 {
+		t.Fatalf("credentials user_version = %d %v, want 7 (row 201 A3)", version, err)
 	}
 }
 
--- a/host/store/mint_requests_test.go
+++ b/host/store/mint_requests_test.go
@@ -129,7 +129,7 @@ func TestCredentialsUpgradeToV4(t *testing.T) {
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["mint_requests"]), normalizeDDL(canonicalMintRequestsDDL); got != want {
 				t.Fatalf("upgraded mint_requests DDL:\n got: %s\nwant: %s", got, want)
 			}
--- a/host/store/surface_approvals.go
+++ b/host/store/surface_approvals.go
@@ -0,0 +1,24 @@
+package store
+
+// Row 201 A3 (design_docs/planned/w-world-approvals-by-key-a3.md §2.2): the surface's recorded
+// approvals (credentials schema 7). The daemon stores a POSTed signed artifact only after the A1b/A2
+// law verified it; the minter (approve --await) reads it and judges it again. A row is never a
+// credential and never trusted: its approver and digest columns are display only.
+
+// ---- M1a ----
+
+const surfaceApprovalsTable = "surface_approvals"
+
+// surfaceApprovalsSchemaVersion is the first credentials schema with surface_approvals.
+const surfaceApprovalsSchemaVersion = 7
+
+// surfaceApprovalsSQL is the version 6 -> 7 upgrade: credentials_schema.sql's surface_approvals
+// table, column for column (TestCredentialsSchemaMatchesCanonicalDDL compares the two).
+const surfaceApprovalsSQL = `CREATE TABLE surface_approvals (
+    request_id  TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'mr-[0-9a-f]*'),
+    format      TEXT NOT NULL CHECK (format IN ('webauthn', 'ssh')),
+    approver    TEXT NOT NULL CHECK (approver <> ''),
+    artifact    BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
+    digest      TEXT NOT NULL CHECK (length(digest) = 64),
+    received_at INTEGER NOT NULL
+)`
````

### Patch M1b — the record and the read

````diff M1b
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -298,6 +298,8 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	_, _, _, errs["PasskeyByCredentialID"] = c.PasskeyByCredentialID(ctx, p.CredentialID)
 	_, _, errs["Passkey"] = c.Passkey(ctx, "x")
 	_, errs["ActivePasskeys"] = c.ActivePasskeys(ctx)
+	errs["RecordSurfaceApproval"] = c.RecordSurfaceApproval(ctx, SurfaceApprovalRow{}, 1000)
+	_, _, errs["SurfaceApproval"] = c.SurfaceApproval(ctx, "mr-0000000000000001")
 	for name, err := range errs {
 		if !errors.Is(err, ErrNoDeadline) {
 			t.Errorf("%s(Background) = %v, want ErrNoDeadline", name, err)
--- a/host/store/surface_approvals.go
+++ b/host/store/surface_approvals.go
@@ -5,6 +5,14 @@ package store
 // law verified it; the minter (approve --await) reads it and judges it again. A row is never a
 // credential and never trusted: its approver and digest columns are display only.
 
+import (
+	"context"
+	"database/sql"
+	"errors"
+	"fmt"
+	"strings"
+)
+
 // ---- M1a ----
 
 const surfaceApprovalsTable = "surface_approvals"
@@ -22,3 +30,118 @@ const surfaceApprovalsSQL = `CREATE TABLE surface_approvals (
     digest      TEXT NOT NULL CHECK (length(digest) = 64),
     received_at INTEGER NOT NULL
 )`
+
+// ---- M1b ----
+
+// MaxSurfaceApprovalArtifactBytes bounds a recorded artifact: the A2 assertion envelope's and the
+// --signature file's own 16 KiB bound (approval.MaxPasskeyResponseBytes, approveMaxSignatureBytes).
+const MaxSurfaceApprovalArtifactBytes = 16 << 10
+
+// ErrBadSurfaceApproval refuses a row that breaks the table's grammar (S2).
+var ErrBadSurfaceApproval = errors.New("store: invalid surface approval row")
+
+// SurfaceApprovalRow is one surface_approvals row: the artifact the daemon verified for a request.
+// Approver ("<format>:<name>") and Digest (the request's digest when it was recorded) are display
+// only; the minter judges Artifact again against the request it reads.
+type SurfaceApprovalRow struct {
+	RequestID  string
+	Format     string // "webauthn" or "ssh"
+	Approver   string // "<format>:<name>"
+	Artifact   []byte // the assertion JSON or the armored SSHSIG, 1..16384 bytes
+	Digest     string // 64 lowercase hex
+	ReceivedAt int64  // set by RecordSurfaceApproval to its now
+}
+
+func (c *Credentials) olderForSurfaceApprovals() error {
+	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: surfaceApprovalsSchemaVersion, Feature: "surface approvals"}
+}
+
+// checkSurfaceApprovalRow is RecordSurfaceApproval's grammar (S2): each refusal is ErrBadSurfaceApproval.
+func checkSurfaceApprovalRow(row SurfaceApprovalRow) error {
+	switch {
+	case !MintRequestIDPattern.MatchString(row.RequestID):
+		return fmt.Errorf("%w: request id %q", ErrBadSurfaceApproval, row.RequestID)
+	case row.Format != "webauthn" && row.Format != "ssh":
+		return fmt.Errorf("%w: format %q is not webauthn or ssh", ErrBadSurfaceApproval, row.Format)
+	case !strings.HasPrefix(row.Approver, row.Format+":") || len(row.Approver) == len(row.Format)+1:
+		return fmt.Errorf("%w: approver %q is not %s:<name>", ErrBadSurfaceApproval, row.Approver, row.Format)
+	case len(row.Artifact) < 1 || len(row.Artifact) > MaxSurfaceApprovalArtifactBytes:
+		return fmt.Errorf("%w: an artifact is 1 to %d bytes, not %d", ErrBadSurfaceApproval, MaxSurfaceApprovalArtifactBytes, len(row.Artifact))
+	case !isLowerHex64(row.Digest):
+		return fmt.Errorf("%w: digest is not 64 lowercase hex", ErrBadSurfaceApproval)
+	}
+	return nil
+}
+
+// RecordSurfaceApproval stores row for a pending request in one transaction, received at now: the
+// request must exist (S3), be unconsumed, created at or before now, unexpired at now, and carry
+// row's digest (S4), or nothing is written and the request's own sentinel is returned. A later row
+// for the same request replaces the earlier one (S5): the minter may have refused it.
+func (c *Credentials) RecordSurfaceApproval(ctx context.Context, row SurfaceApprovalRow, now int64) error {
+	if err := requireDeadline(ctx); err != nil { // S1
+		return err
+	}
+	if c.version < surfaceApprovalsSchemaVersion { // S1
+		return c.olderForSurfaceApprovals()
+	}
+	if err := checkSurfaceApprovalRow(row); err != nil { // S2
+		return err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return fmt.Errorf("store: record surface approval: begin: %w", err)
+	}
+	defer func() { _ = tx.Rollback() }()
+	var (
+		digest           string
+		createdAt, expAt int64
+		consumed         bool
+	)
+	err = tx.QueryRowContext(ctx, `SELECT digest, created_at, expires_at, consumed_at IS NOT NULL FROM mint_requests WHERE request_id = ?`, row.RequestID).
+		Scan(&digest, &createdAt, &expAt, &consumed)
+	if errors.Is(err, sql.ErrNoRows) { // S3
+		return fmt.Errorf("%w: %s", ErrRequestUnknown, row.RequestID)
+	}
+	if err != nil {
+		return fmt.Errorf("store: record surface approval: %w", err)
+	}
+	switch { // S4
+	case consumed:
+		return fmt.Errorf("%w: %s", ErrRequestConsumed, row.RequestID)
+	case now < createdAt:
+		return fmt.Errorf("%w: %s", ErrRequestClockBehind, row.RequestID)
+	case now >= expAt:
+		return fmt.Errorf("%w: %s", ErrRequestExpired, row.RequestID)
+	case digest != row.Digest:
+		return fmt.Errorf("%w: %s", ErrRequestDigestChanged, row.RequestID)
+	}
+	if _, err := tx.ExecContext(ctx, `INSERT INTO surface_approvals (request_id, format, approver, artifact, digest, received_at) VALUES (?, ?, ?, ?, ?, ?)
+ON CONFLICT(request_id) DO UPDATE SET format = excluded.format, approver = excluded.approver, artifact = excluded.artifact, digest = excluded.digest, received_at = excluded.received_at`,
+		row.RequestID, row.Format, row.Approver, row.Artifact, row.Digest, now); err != nil { // S5
+		return fmt.Errorf("store: record surface approval: %w", err)
+	}
+	if err := tx.Commit(); err != nil {
+		return fmt.Errorf("store: record surface approval: commit: %w", err)
+	}
+	return nil
+}
+
+// SurfaceApproval reads the recorded approval for request id. ok is false when there is none.
+func (c *Credentials) SurfaceApproval(ctx context.Context, id string) (SurfaceApprovalRow, bool, error) {
+	if err := requireDeadline(ctx); err != nil {
+		return SurfaceApprovalRow{}, false, err
+	}
+	if c.version < surfaceApprovalsSchemaVersion {
+		return SurfaceApprovalRow{}, false, c.olderForSurfaceApprovals()
+	}
+	var row SurfaceApprovalRow
+	err := c.db.QueryRowContext(ctx, `SELECT request_id, format, approver, artifact, digest, received_at FROM surface_approvals WHERE request_id = ?`, id).
+		Scan(&row.RequestID, &row.Format, &row.Approver, &row.Artifact, &row.Digest, &row.ReceivedAt)
+	if errors.Is(err, sql.ErrNoRows) {
+		return SurfaceApprovalRow{}, false, nil
+	}
+	if err != nil {
+		return SurfaceApprovalRow{}, false, fmt.Errorf("store: read surface approval: %w", err)
+	}
+	return row, true, nil
+}
--- a/host/store/surface_approvals_test.go
+++ b/host/store/surface_approvals_test.go
@@ -0,0 +1,209 @@
+package store
+
+// Row 201 A3 sprint 1 M1 (design_docs/planned/w-world-approvals-by-key-a3.md §2.2, §6 SA1-SA3):
+// credentials schema 7's surface_approvals. Every store is under t.TempDir() (D-WORLD-109).
+
+import (
+	"bytes"
+	"errors"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+)
+
+// canonicalSurfaceApprovalsDDL is the current (version 7, row 201 A3) surface_approvals table.
+var canonicalSurfaceApprovalsDDL = `CREATE TABLE surface_approvals (
+    request_id  TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'mr-[0-9a-f]*'),
+    format      TEXT NOT NULL CHECK (format IN ('webauthn', 'ssh')),
+    approver    TEXT NOT NULL CHECK (approver <> ''),
+    artifact    BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
+    digest      TEXT NOT NULL CHECK (length(digest) = 64),
+    received_at INTEGER NOT NULL
+)`
+
+// makeV6Credentials writes a schema-6 credentials store at path: makeV5Credentials, the
+// webauthn_credentials table and user_version 6.
+func makeV6Credentials(t *testing.T, path string, rows []SessionRow) {
+	t.Helper()
+	makeV5Credentials(t, path, rows)
+	db := rawDB(t, path)
+	if _, err := db.Exec(passkeysSQL); err != nil {
+		t.Fatal(err)
+	}
+	setPositiveFixtureVersion(t, db, 6)
+}
+
+// testSurfaceApproval is a well-formed row for testMintRequestRow(n, ...).
+func testSurfaceApproval(n int) SurfaceApprovalRow {
+	r := testMintRequestRow(n, 1000)
+	return SurfaceApprovalRow{RequestID: r.RequestID, Format: "webauthn", Approver: "webauthn:mark", Artifact: []byte(`{"id":"x"}`), Digest: r.Digest}
+}
+
+// surfaceApprovalCount counts c's surface_approvals rows.
+func surfaceApprovalCount(t *testing.T, c *Credentials) int {
+	t.Helper()
+	var n int
+	if err := c.db.QueryRow(`SELECT count(*) FROM surface_approvals`).Scan(&n); err != nil {
+		t.Fatal(err)
+	}
+	return n
+}
+
+// SA1: a fresh store is schema 7 with surface_approvals as credentials_schema.sql and the upgrade
+// write it; a schema-6 store is upgraded once by UpgradeCredentials only, its passkey kept and the
+// new table empty; a session open of the v6 store neither upgrades it nor records an approval,
+// which names the serve restart.
+func TestCredentialsSchema7SurfaceApprovals(t *testing.T) {
+	c := openTestCredentials(t)
+	ddl := tableDDL(t, c.db)
+	for name, got := range map[string]string{"schema file": ddl["surface_approvals"], "surfaceApprovalsSQL": surfaceApprovalsSQL} {
+		if normalizeDDL(got) != normalizeDDL(canonicalSurfaceApprovalsDDL) {
+			t.Fatalf("surface_approvals DDL (%s):\n got: %s\nwant: %s", name, got, canonicalSurfaceApprovalsDDL)
+		}
+	}
+	if got, want := credentialsTables(7), []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"}; !reflect.DeepEqual(got, want) {
+		t.Fatalf("credentialsTables(7) = %v, want %v", got, want)
+	}
+	if got, want := credentialsTables(6), []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"}; !reflect.DeepEqual(got, want) {
+		t.Fatalf("credentialsTables(6) = %v, want %v", got, want)
+	}
+
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	makeV6Credentials(t, path, []SessionRow{oneBoundRow()})
+	a, p := testPasskey("kept", 3)
+	raw := rawDB(t, path)
+	for _, ins := range []struct {
+		q    string
+		args []any
+	}{
+		{`INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES (?, ?, ?, ?)`, []any{a.Name, a.Format, a.PublicKey, a.EnrolledAt}},
+		{`INSERT INTO webauthn_credentials (name, credential_id, rp_id, origin, alg, aaguid, backup_eligible, sign_count) VALUES (?, ?, ?, ?, ?, ?, 0, 1)`,
+			[]any{p.Name, p.CredentialID, p.RPID, p.Origin, p.Alg, p.AAGUID}},
+		{`INSERT INTO mint_requests (request_id, request, digest, challenge, episode_id, created_at, expires_at) VALUES (?, x'7b7d', ?, ?, 'ep', 1000, 1600)`,
+			[]any{testSurfaceApproval(1).RequestID, testSurfaceApproval(1).Digest, strings.Repeat("cd", 32)}},
+	} {
+		if _, err := raw.Exec(ins.q, ins.args...); err != nil {
+			t.Fatal(err)
+		}
+	}
+	v6, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	err = v6.RecordSurfaceApproval(testCtx(t), testSurfaceApproval(1), 1010)
+	var oe *OlderCredentialsStoreError
+	if !errors.As(err, &oe) || oe.Version != 6 || oe.Need != 7 || oe.Feature != "surface approvals" ||
+		!strings.Contains(err.Error(), "surface approvals need schema 7: restart `ailang-worldd serve`") {
+		t.Fatalf("v6 record: %v, want the restart refusal naming schema 7", err)
+	}
+	if _, _, err := v6.SurfaceApproval(testCtx(t), testSurfaceApproval(1).RequestID); !errors.As(err, &oe) || oe.Need != 7 {
+		t.Fatalf("v6 read: %v, want OlderCredentialsStoreError need 7", err)
+	}
+	_ = v6.Close()
+	if got := userVersion(t, path); got != 6 {
+		t.Fatalf("a session open changed user_version to %d, want 6", got)
+	}
+	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded {
+		t.Fatalf("UpgradeCredentials: upgraded %v, err %v; want true, nil", upgraded, err)
+	}
+	if again, err := UpgradeCredentials(path); err != nil || again {
+		t.Fatalf("second UpgradeCredentials: upgraded %v, err %v; want false, nil", again, err)
+	}
+	if got := userVersion(t, path); got != 7 {
+		t.Fatalf("user_version after the upgrade = %d, want 7", got)
+	}
+	up := tableDDL(t, rawDB(t, path))
+	requireExactTableNames(t, "upgraded v6 store", up, credentialsTables(7))
+	if got, want := normalizeDDL(up["surface_approvals"]), normalizeDDL(canonicalSurfaceApprovalsDDL); got != want {
+		t.Fatalf("upgraded surface_approvals DDL:\n got: %s\nwant: %s", got, want)
+	}
+	ro, err := OpenCredentialsReadOnly(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	if _, gotP, ok, err := ro.PasskeyByCredentialID(testCtx(t), p.CredentialID); err != nil || !ok || gotP.SignCount != 1 {
+		t.Fatalf("the passkey after the upgrade: ok %v err %v %+v", ok, err, gotP)
+	}
+	if _, ok, err := ro.SurfaceApproval(testCtx(t), testSurfaceApproval(1).RequestID); err != nil || ok {
+		t.Fatalf("a recorded approval right after the upgrade: ok %v err %v; want none", ok, err)
+	}
+}
+
+// SA2: each guard of RecordSurfaceApproval refuses with its sentinel and writes nothing; a valid row
+// is read back with received_at = now; a later valid row replaces it while the request is pending.
+func TestRecordSurfaceApprovalGuards(t *testing.T) {
+	c := openTestCredentials(t)
+	ctx := testCtx(t)
+	pending := testMintRequestRow(1, 1000)
+	if err := c.PutMintRequest(ctx, pending, 1000); err != nil {
+		t.Fatal(err)
+	}
+	used := testMintRequestRow(2, 1000)
+	if err := c.PutMintRequest(ctx, used, 1000); err != nil {
+		t.Fatal(err)
+	}
+	if err := c.MintSessionForRequest(ctx, consumeSession(1), okConsume(used, 1005, new(int), new(int))); err != nil {
+		t.Fatal(err)
+	}
+	edit := func(f func(*SurfaceApprovalRow)) SurfaceApprovalRow {
+		r := testSurfaceApproval(1)
+		f(&r)
+		return r
+	}
+	for _, arm := range []struct {
+		name string
+		row  SurfaceApprovalRow
+		now  int64
+		want error
+	}{
+		{"S2-id", edit(func(r *SurfaceApprovalRow) { r.RequestID = "mr-00000000000000A1" }), 1010, ErrBadSurfaceApproval},
+		{"S2-format", edit(func(r *SurfaceApprovalRow) { r.Format, r.Approver = "tty", "tty:mark" }), 1010, ErrBadSurfaceApproval},
+		{"S2-approver-other-format", edit(func(r *SurfaceApprovalRow) { r.Approver = "ssh:mark" }), 1010, ErrBadSurfaceApproval},
+		{"S2-approver-no-name", edit(func(r *SurfaceApprovalRow) { r.Approver = "webauthn:" }), 1010, ErrBadSurfaceApproval},
+		{"S2-artifact-empty", edit(func(r *SurfaceApprovalRow) { r.Artifact = nil }), 1010, ErrBadSurfaceApproval},
+		{"S2-artifact-over", edit(func(r *SurfaceApprovalRow) { r.Artifact = bytes.Repeat([]byte("a"), 16385) }), 1010, ErrBadSurfaceApproval},
+		{"S2-digest", edit(func(r *SurfaceApprovalRow) { r.Digest = strings.Repeat("AB", 32) }), 1010, ErrBadSurfaceApproval},
+		{"S3-unknown", edit(func(r *SurfaceApprovalRow) { r.RequestID = testMintRequestRow(9, 1000).RequestID }), 1010, ErrRequestUnknown},
+		{"S4-consumed", edit(func(r *SurfaceApprovalRow) { r.RequestID = used.RequestID }), 1010, ErrRequestConsumed},
+		{"S4-clock-behind", testSurfaceApproval(1), 999, ErrRequestClockBehind},
+		{"S4-expired", testSurfaceApproval(1), 1600, ErrRequestExpired},
+		{"S4-digest-changed", edit(func(r *SurfaceApprovalRow) { r.Digest = strings.Repeat("ef", 32) }), 1010, ErrRequestDigestChanged},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			if err := c.RecordSurfaceApproval(ctx, arm.row, arm.now); !errors.Is(err, arm.want) {
+				t.Fatalf("RecordSurfaceApproval = %v, want %v", err, arm.want)
+			}
+			if n := surfaceApprovalCount(t, c); n != 0 {
+				t.Fatalf("%d rows after a refusal, want 0", n)
+			}
+		})
+	}
+	// Controls: the bounds themselves are admitted (an artifact of exactly 16384 bytes, the
+	// request's creation second, the last second before expiry).
+	for _, ok := range []struct {
+		row SurfaceApprovalRow
+		now int64
+	}{
+		{edit(func(r *SurfaceApprovalRow) { r.Artifact = bytes.Repeat([]byte("a"), 16384) }), 1000},
+		{edit(func(r *SurfaceApprovalRow) { r.Artifact = []byte("b") }), 1599},
+	} {
+		if err := c.RecordSurfaceApproval(ctx, ok.row, ok.now); err != nil {
+			t.Fatalf("control at %d: %v", ok.now, err)
+		}
+	}
+	want := SurfaceApprovalRow{RequestID: pending.RequestID, Format: "ssh", Approver: "ssh:mark-mac", Artifact: []byte("-----BEGIN SSH SIGNATURE-----\n"), Digest: pending.Digest, ReceivedAt: 1010}
+	in := want
+	in.ReceivedAt = 0
+	if err := c.RecordSurfaceApproval(ctx, in, 1010); err != nil {
+		t.Fatal(err)
+	}
+	got, ok, err := c.SurfaceApproval(ctx, pending.RequestID)
+	if err != nil || !ok || !reflect.DeepEqual(got, want) || surfaceApprovalCount(t, c) != 1 {
+		t.Fatalf("after S5's replace: %+v ok %v err %v (%d rows); want %+v alone", got, ok, err, surfaceApprovalCount(t, c), want)
+	}
+	if _, ok, err := c.SurfaceApproval(ctx, used.RequestID); err != nil || ok {
+		t.Fatalf("the consumed request has a recorded approval: ok %v err %v", ok, err)
+	}
+}
````

### Patch M1c — the prune

````diff M1c
--- a/host/store/mint_requests.go
+++ b/host/store/mint_requests.go
@@ -92,8 +92,9 @@ type MintRequestRow struct {
 	CredentialID string // "" unless Consumed
 }
 
-// PutMintRequest stores one new, unconsumed request. It prunes requests expired a day ago, and
-// refuses when MaxPendingMintRequests live ones are already held.
+// PutMintRequest stores one new, unconsumed request. It prunes requests expired a day ago (and,
+// from schema 7, their recorded surface approvals), and refuses when MaxPendingMintRequests live
+// ones are already held.
 func (c *Credentials) PutMintRequest(ctx context.Context, row MintRequestRow, now int64) error {
 	if err := requireDeadline(ctx); err != nil {
 		return err
@@ -130,6 +131,11 @@ func (c *Credentials) PutMintRequest(ctx context.Context, row MintRequestRow, no
 	if _, err := tx.ExecContext(ctx, `DELETE FROM mint_requests WHERE consumed_at IS NULL AND expires_at <= ?`, now-mintRequestPruneAfter); err != nil {
 		return fmt.Errorf("store: put mint request: prune: %w", err)
 	}
+	if c.version >= surfaceApprovalsSchemaVersion { // row 201 A3: a pruned request's recorded approval goes with it
+		if _, err := tx.ExecContext(ctx, `DELETE FROM surface_approvals WHERE request_id NOT IN (SELECT request_id FROM mint_requests)`); err != nil {
+			return fmt.Errorf("store: put mint request: prune surface approvals: %w", err)
+		}
+	}
 	var live int
 	if err := tx.QueryRowContext(ctx, `SELECT count(*) FROM mint_requests WHERE consumed_at IS NULL AND expires_at > ?`, now).Scan(&live); err != nil {
 		return fmt.Errorf("store: put mint request: count: %w", err)
--- a/host/store/surface_approvals_test.go
+++ b/host/store/surface_approvals_test.go
@@ -207,3 +207,52 @@ func TestRecordSurfaceApprovalGuards(t *testing.T) {
 		t.Fatalf("the consumed request has a recorded approval: ok %v err %v", ok, err)
 	}
 }
+
+// SA3: PutMintRequest's prune deletes the recorded approval of a request it prunes, in the same
+// transaction, and keeps a consumed request's; on a schema-6 store it prunes and puts as before.
+func TestSurfaceApprovalPrunedWithItsRequest(t *testing.T) {
+	c := openTestCredentials(t)
+	ctx := testCtx(t)
+	gone, kept := testMintRequestRow(1, 1000), testMintRequestRow(2, 1000)
+	for _, r := range []MintRequestRow{gone, kept} {
+		if err := c.PutMintRequest(ctx, r, 1000); err != nil {
+			t.Fatal(err)
+		}
+		row := testSurfaceApproval(1)
+		row.RequestID = r.RequestID
+		if err := c.RecordSurfaceApproval(ctx, row, 1010); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if err := c.MintSessionForRequest(ctx, consumeSession(1), okConsume(kept, 1020, new(int), new(int))); err != nil {
+		t.Fatal(err)
+	}
+	if err := c.PutMintRequest(ctx, testMintRequestRow(3, 1600+86400), 1600+86400); err != nil {
+		t.Fatal(err)
+	}
+	if _, ok, err := c.SurfaceApproval(ctx, gone.RequestID); err != nil || ok {
+		t.Fatalf("the pruned request's approval: ok %v err %v; want gone", ok, err)
+	}
+	if _, ok, err := c.SurfaceApproval(ctx, kept.RequestID); err != nil || !ok {
+		t.Fatalf("the consumed request's approval: ok %v err %v; want kept", ok, err)
+	}
+	if n := surfaceApprovalCount(t, c); n != 1 {
+		t.Fatalf("%d rows after the prune, want 1", n)
+	}
+
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	makeV6Credentials(t, path, nil)
+	v6, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = v6.Close() }()
+	for _, r := range []MintRequestRow{gone, testMintRequestRow(3, 1600+86400)} {
+		if err := v6.PutMintRequest(ctx, r, r.CreatedAt); err != nil {
+			t.Fatalf("v6 put %s: %v", r.RequestID, err)
+		}
+	}
+	if _, ok, err := v6.MintRequest(ctx, gone.RequestID); err != nil || ok {
+		t.Fatalf("v6: the expired request was not pruned: ok %v err %v", ok, err)
+	}
+}
````

### Patch M2a — the shared judge

````diff M2a
--- a/host/authority/judge.go
+++ b/host/authority/judge.go
@@ -0,0 +1,129 @@
+package authority
+
+// Row 201 A3 sprint 1 M2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.5): one judge for
+// the CLI and the daemon. approve --webauthn-assertion and --signature (cmd/ailang-worldd), the
+// minter approve --await, and from sprint 2 the surface's POST all judge an artifact here: the
+// lookups on a read-only credentials handle, then the A1b/A2 law, unchanged.
+
+import (
+	"context"
+	"crypto/sha256"
+	"encoding/hex"
+	"errors"
+	"fmt"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// The steps a JudgeError names, in the order the judges run them.
+const (
+	JudgeStepRequest    = "request"     // the row's digest or challenge is not 64 hex
+	JudgeStepParse      = "parse"       // the artifact is not an assertion or an SSHSIG
+	JudgeStepStore      = "store"       // the lookup failed (below the schema its format needs)
+	JudgeStepUnknownKey = "unknown-key" // no enrolled approver holds the artifact's key or credential
+	JudgeStepStoredKey  = "stored-key"  // the enrolled passkey's stored key does not parse
+	JudgeStepLaw        = "law"         // the law refused it
+)
+
+// ErrNotEnrolled is an unknown-key JudgeError's Err.
+var ErrNotEnrolled = errors.New("authority: no enrolled approver holds that key")
+
+// JudgeError is a refused judgement: the step that refused, what the CLI names for it, and the
+// cause. The CLI prints its own text per step (cmd/ailang-worldd passkeyRefusal, sshRefusal).
+type JudgeError struct {
+	Step        string
+	RawID       []byte // unknown-key, assertion: the credential ID the response names
+	Fingerprint string // unknown-key, SSHSIG: the signature key's fingerprint
+	Approver    string // stored-key: the approver whose stored key is invalid
+	Err         error
+}
+
+func (e *JudgeError) Error() string { return fmt.Sprintf("authority: judge: %s: %v", e.Step, e.Err) }
+
+func (e *JudgeError) Unwrap() error { return e.Err }
+
+// judgeAdmit is the law's admission input and the challenge for row at now. The CLI has admitted
+// row and checked its challenge before (F4); the daemon checks them too (R8, R10). A row that is
+// not 64 hex in either is refused here, never converted.
+func judgeAdmit(row store.MintRequestRow, now int64) (approval.AdmitInput, [32]byte, error) {
+	stored, err := hex.DecodeString(row.Digest)
+	if err != nil || len(stored) != 32 {
+		return approval.AdmitInput{}, [32]byte{}, &JudgeError{Step: JudgeStepRequest, Err: fmt.Errorf("request %s's stored digest is not 64 hex", row.RequestID)}
+	}
+	challenge, err := hex.DecodeString(row.Challenge)
+	if err != nil || len(challenge) != 32 {
+		return approval.AdmitInput{}, [32]byte{}, &JudgeError{Step: JudgeStepRequest, Err: fmt.Errorf("request %s's stored challenge is not 64 hex", row.RequestID)}
+	}
+	return approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Consumed,
+		StoredDigest: [32]byte(stored), BytesDigest: sha256.Sum256(row.Request)}, [32]byte(challenge), nil
+}
+
+// JudgeAssertion judges a passkey assertion response for row at now: ParseAssertionResponse, the
+// passkey its credential ID names (PasskeyByCredentialID on ro), ParseStoredCOSEKey, then
+// approval.ApproveByAssertion. It returns the verified approval or a *JudgeError.
+func JudgeAssertion(ctx context.Context, ro *store.Credentials, row store.MintRequestRow, response []byte, now int64) (*approval.SignedApproval, error) {
+	admit, challenge, err := judgeAdmit(row, now)
+	if err != nil {
+		return nil, err
+	}
+	resp, err := approval.ParseAssertionResponse(response)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
+	}
+	a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepStore, Err: err}
+	}
+	if !ok {
+		return nil, &JudgeError{Step: JudgeStepUnknownKey, RawID: resp.RawID, Err: ErrNotEnrolled}
+	}
+	key, err := approval.ParseStoredCOSEKey(a.PublicKey)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepStoredKey, Approver: a.Name, Err: err}
+	}
+	signed, err := approval.ApproveByAssertion(approval.AssertionInput{
+		Admit:     admit,
+		Challenge: challenge,
+		Response:  response,
+		Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,
+			Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,
+				BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}},
+	})
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
+	}
+	return signed, nil
+}
+
+// JudgeSSHSignature judges an armored SSHSIG for row at now: ParseSSHSIG, the approver its key
+// names (ApproverByKey on ro), then approval.ApproveBySignature. It returns the verified approval
+// or a *JudgeError.
+func JudgeSSHSignature(ctx context.Context, ro *store.Credentials, row store.MintRequestRow, armored []byte, now int64) (*approval.SignedApproval, error) {
+	admit, challenge, err := judgeAdmit(row, now)
+	if err != nil {
+		return nil, err
+	}
+	sig, err := approval.ParseSSHSIG(armored)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepParse, Err: err}
+	}
+	enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepStore, Err: err}
+	}
+	if !ok {
+		return nil, &JudgeError{Step: JudgeStepUnknownKey, Fingerprint: sig.PublicKey.Fingerprint(), Err: ErrNotEnrolled}
+	}
+	signed, err := approval.ApproveBySignature(approval.SignatureInput{
+		Admit:     admit,
+		Challenge: challenge,
+		Signature: armored,
+		Approver: approval.Approver{Name: enrolled.Name, Format: enrolled.Format, Key: approval.SSHPublicKey{Blob: enrolled.PublicKey},
+			EnrolledAt: enrolled.EnrolledAt, LastUsedAt: enrolled.LastUsedAt, RevokedAt: enrolled.RevokedAt},
+	})
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
+	}
+	return signed, nil
+}
--- a/host/authority/judge_test.go
+++ b/host/authority/judge_test.go
@@ -0,0 +1,247 @@
+package authority
+
+// Row 201 A3 sprint 1 M2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.5, §6 JG1, JG2): the
+// shared judge, over the CLI's refusal matrix (WF3-WF6, F11-F14). The CLI's texts for each step are
+// pinned, unedited, by cmd/ailang-worldd's TestApprovePasskeyRefusalsBeforeAnyEffect and
+// TestApproveBySignatureRefusalsBeforeAnyEffect. Every key is approvaltest's (D-WORLD-109).
+
+import (
+	"database/sql"
+	"encoding/hex"
+	"errors"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+const jgOrigin, jgRPID = "http://localhost:7644", "localhost"
+
+var (
+	jgPasskey = approvaltest.NewPasskey(approval.AlgEdDSA, 61)
+	jgSigner  = approvaltest.NewSigner(62)
+)
+
+// jgRig is a schema-7 store at path holding au1Request's row, jgPasskey enrolled as mark-mac
+// (counter 1) and jgSigner as mark-ssh, both at 900.
+func jgRig(t *testing.T) (string, store.MintRequestRow) {
+	t.Helper()
+	path := filepath.Join(t.TempDir(), store.CredentialsFileName)
+	st, err := store.OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = st.Close() }()
+	_, row := au1Request(t, st)
+	k := jgPasskey.Key()
+	ctx := boundedTestContext(t)
+	if err := st.EnrolPasskey(ctx, store.ApproverRow{Name: "mark-mac", Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: 900},
+		store.PasskeyRow{Name: "mark-mac", RPID: jgRPID, Origin: jgOrigin, CredentialID: jgPasskey.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: 1}); err != nil {
+		t.Fatal(err)
+	}
+	if err := st.EnrolApprover(ctx, store.ApproverRow{Name: "mark-ssh", Format: "ssh", PublicKey: jgSigner.PublicKey().Blob, EnrolledAt: 900}); err != nil {
+		t.Fatal(err)
+	}
+	return path, row
+}
+
+// jgExec runs stmts on the store at path.
+func jgExec(t *testing.T, path string, stmts ...string) {
+	t.Helper()
+	db, err := sql.Open("sqlite", path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	for _, s := range stmts {
+		if _, err := db.Exec(s); err != nil {
+			t.Fatal(err)
+		}
+	}
+}
+
+// jgHash and jgMessage are what an approval of row signs.
+func jgHash(t *testing.T, row store.MintRequestRow) [32]byte {
+	t.Helper()
+	d, _ := hex.DecodeString(row.Digest)
+	c, _ := hex.DecodeString(row.Challenge)
+	h, err := approval.MessageHash([32]byte(d), [32]byte(c), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return h
+}
+
+func jgMessage(t *testing.T, row store.MintRequestRow) []byte {
+	t.Helper()
+	d, _ := hex.DecodeString(row.Digest)
+	c, _ := hex.DecodeString(row.Challenge)
+	m, err := approval.SignedMessage([32]byte(d), [32]byte(c), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return m
+}
+
+// judgeArm is one refusal arm: what to change, the artifact, and the step and cause it must name.
+type judgeArm struct {
+	name     string
+	setup    []string                    // SQL on the store before the judgement
+	edit     func(*store.MintRequestRow) // the row the judge is handed
+	artifact func(t *testing.T, row store.MintRequestRow) []byte
+	now      int64 // 0: 1010
+	step     string
+	cause    error  // errors.Is target; nil: only the step
+	field    string // the JudgeError field the CLI prints, when the step has one
+}
+
+// runJudgeArms judges each arm on a fresh rig through judge, and requires the step, the cause and
+// the CLI's field; a verified control on the same rig binds the row.
+func runJudgeArms(t *testing.T, judge func(ro *store.Credentials, row store.MintRequestRow, artifact []byte, now int64) (*approval.SignedApproval, error),
+	arms []judgeArm, field func(*JudgeError) string, good func(t *testing.T, row store.MintRequestRow) []byte, approvedBy string) {
+	t.Helper()
+	for _, arm := range arms {
+		t.Run(arm.name, func(t *testing.T) {
+			path, row := jgRig(t)
+			jgExec(t, path, arm.setup...)
+			art := arm.artifact(t, row)
+			if arm.edit != nil {
+				arm.edit(&row)
+			}
+			now := arm.now
+			if now == 0 {
+				now = 1010
+			}
+			ro, err := store.OpenCredentialsReadOnly(path)
+			if err != nil {
+				t.Fatal(err)
+			}
+			defer func() { _ = ro.Close() }()
+			signed, err := judge(ro, row, art, now)
+			var je *JudgeError
+			if signed != nil || !errors.As(err, &je) || je.Step != arm.step {
+				t.Fatalf("judge = %v, %v; want a %q JudgeError", signed, err, arm.step)
+			}
+			if arm.cause != nil && !errors.Is(err, arm.cause) {
+				t.Fatalf("judge = %v, want it to wrap %v", err, arm.cause)
+			}
+			if got := field(je); got != arm.field {
+				t.Fatalf("the CLI's field = %q, want %q", got, arm.field)
+			}
+		})
+	}
+	t.Run("verified", func(t *testing.T) {
+		path, row := jgRig(t)
+		ro, err := store.OpenCredentialsReadOnly(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer func() { _ = ro.Close() }()
+		signed, err := judge(ro, row, good(t, row), 1010)
+		if err != nil || signed == nil || signed.ApprovedBy() != approvedBy {
+			t.Fatalf("judge = %v, %v; want approved by %s", signed, err, approvedBy)
+		}
+		d, _ := hex.DecodeString(row.Digest)
+		c, _ := hex.DecodeString(row.Challenge)
+		if err := signed.Binds([32]byte(d), [32]byte(c), row.ExpiresAt); err != nil {
+			t.Fatalf("the verified approval does not bind the row: %v", err)
+		}
+	})
+}
+
+// requestArms are the row arms both judges share: a row whose digest or challenge is not 64 hex is
+// refused before any parse or lookup; a consumed row, or one whose bytes are not its digest, is
+// refused by the law's admission.
+func requestArms(good func(t *testing.T, row store.MintRequestRow) []byte) []judgeArm {
+	return []judgeArm{
+		{name: "request-digest", edit: func(r *store.MintRequestRow) { r.Digest = strings.Repeat("zz", 32) }, artifact: good, step: JudgeStepRequest},
+		{name: "request-challenge", edit: func(r *store.MintRequestRow) { r.Challenge = "abcd" }, artifact: good, step: JudgeStepRequest},
+		{name: "law-consumed", edit: func(r *store.MintRequestRow) { r.Consumed = true }, artifact: good, step: JudgeStepLaw},
+		{name: "law-bytes-changed", edit: func(r *store.MintRequestRow) { r.Request = append(append([]byte(nil), r.Request...), ' ') }, artifact: good, step: JudgeStepLaw},
+	}
+}
+
+// JG1: JudgeAssertion's step, cause and field for each of the CLI's WF3-WF6 refusals, and a
+// verified assertion that binds the row.
+func TestJudgeAssertionIsTheCLIJudge(t *testing.T) {
+	good := func(t *testing.T, row store.MintRequestRow) []byte {
+		return jgPasskey.Assert(jgOrigin, jgRPID, jgHash(t, row), 2)
+	}
+	stranger := approvaltest.NewPasskey(approval.AlgEdDSA, 63)
+	arms := append(requestArms(good),
+		judgeArm{name: "WF3-not-an-assertion", artifact: func(*testing.T, store.MintRequestRow) []byte { return []byte(`{"id":"x"}`) },
+			step: JudgeStepParse, cause: approval.ErrPasskeyResponse},
+		judgeArm{name: "WF4-v5-store", setup: []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, artifact: good, step: JudgeStepStore},
+		judgeArm{name: "WF5-unknown-credential", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
+			return stranger.Assert(jgOrigin, jgRPID, jgHash(t, row), 2)
+		}, step: JudgeStepUnknownKey, cause: ErrNotEnrolled, field: string(stranger.CredentialID())},
+		judgeArm{name: "WF5-stored-key-invalid", setup: []string{`UPDATE approvers SET public_key = x'a0' WHERE name = 'mark-mac'`}, artifact: good,
+			step: JudgeStepStoredKey, field: "mark-mac"},
+		judgeArm{name: "WF6-counter-regression", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
+			return jgPasskey.Assert(jgOrigin, jgRPID, jgHash(t, row), 1)
+		}, step: JudgeStepLaw, cause: approval.ErrSignCountRegression},
+		judgeArm{name: "WF6-other-request", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
+			h := jgHash(t, row)
+			h[31] ^= 1
+			return jgPasskey.Assert(jgOrigin, jgRPID, h, 2)
+		}, step: JudgeStepLaw, cause: approval.ErrChallenge},
+		judgeArm{name: "WF6-revoked", setup: []string{`UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-mac'`}, artifact: good,
+			step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
+		judgeArm{name: "WF6-another-origin", setup: []string{`UPDATE webauthn_credentials SET origin = 'http://localhost:7645' WHERE name = 'mark-mac'`}, artifact: good,
+			step: JudgeStepLaw, cause: approval.ErrOriginMismatch},
+		judgeArm{name: "law-at-now-expired", artifact: good, now: 1600, step: JudgeStepLaw},
+	)
+	field := func(je *JudgeError) string {
+		switch je.Step {
+		case JudgeStepUnknownKey:
+			return string(je.RawID)
+		case JudgeStepStoredKey:
+			return je.Approver
+		}
+		return ""
+	}
+	runJudgeArms(t, func(ro *store.Credentials, row store.MintRequestRow, artifact []byte, now int64) (*approval.SignedApproval, error) {
+		return JudgeAssertion(boundedTestContext(t), ro, row, artifact, now)
+	}, arms, field, good, "webauthn:mark-mac")
+}
+
+// JG2: JudgeSSHSignature's step, cause and field for each of the CLI's F11-F14 refusals, and a
+// verified signature that binds the row.
+func TestJudgeSSHSignatureIsTheCLIJudge(t *testing.T) {
+	good := func(t *testing.T, row store.MintRequestRow) []byte {
+		return jgSigner.Sign(approval.Namespace, jgMessage(t, row))
+	}
+	stranger := approvaltest.NewSigner(64)
+	arms := append(requestArms(good),
+		judgeArm{name: "F11-garbage", artifact: func(*testing.T, store.MintRequestRow) []byte { return []byte("not a signature\n") },
+			step: JudgeStepParse, cause: approval.ErrArmor},
+		judgeArm{name: "F12-v4-store", setup: []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
+			artifact: good, step: JudgeStepStore},
+		judgeArm{name: "F13-unenrolled-key", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
+			return stranger.Sign(approval.Namespace, jgMessage(t, row))
+		}, step: JudgeStepUnknownKey, cause: ErrNotEnrolled, field: stranger.PublicKey().Fingerprint()},
+		judgeArm{name: "F14-other-request", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
+			m := jgMessage(t, row)
+			m[len(m)-9] ^= 1
+			return jgSigner.Sign(approval.Namespace, m)
+		}, step: JudgeStepLaw, cause: approval.ErrBadSignature},
+		judgeArm{name: "F14-git-namespace", artifact: func(t *testing.T, row store.MintRequestRow) []byte { return jgSigner.Sign("git", jgMessage(t, row)) },
+			step: JudgeStepLaw, cause: approval.ErrNamespace},
+		judgeArm{name: "F14-revoked", setup: []string{`UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-ssh'`}, artifact: good,
+			step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
+		judgeArm{name: "F14-format-not-ssh", setup: []string{`UPDATE approvers SET format = 'webauthn' WHERE name = 'mark-ssh'`}, artifact: good, step: JudgeStepLaw},
+		judgeArm{name: "law-at-now-expired", artifact: good, now: 1600, step: JudgeStepLaw},
+	)
+	field := func(je *JudgeError) string {
+		if je.Step == JudgeStepUnknownKey {
+			return je.Fingerprint
+		}
+		return ""
+	}
+	runJudgeArms(t, func(ro *store.Credentials, row store.MintRequestRow, artifact []byte, now int64) (*approval.SignedApproval, error) {
+		return JudgeSSHSignature(boundedTestContext(t), ro, row, artifact, now)
+	}, arms, field, good, "ssh:mark-ssh")
+}
````

### Patch M2b — the CLI wrappers call the judge

````diff M2b
--- a/cmd/ailang-worldd/approve_passkey.go
+++ b/cmd/ailang-worldd/approve_passkey.go
@@ -7,13 +7,14 @@ package main
 
 import (
 	"context"
-	"crypto/sha256"
 	"encoding/base64"
 	"encoding/hex"
+	"errors"
 	"fmt"
 	"io"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
 	"github.com/sunholo-data/ailang-world/host/store"
 )
 
@@ -68,48 +69,42 @@ func requestHash(row store.MintRequestRow) ([32]byte, error) {
 }
 
 // approveAssertion is approve --webauthn-assertion's check, before any network call (design §2.13
-// WF2-WF6): the file is read once through a descriptor check (WF2) and parsed (WF3), the passkey
-// its credential ID names is looked up on a read-only handle (WF4, WF5), and the assertion law
-// judges it against the stored row and that passkey (WF6). It returns the verified approval, or
-// the refusal text. The caller has already admitted row at now and decoded its challenge (F4).
+// WF2-WF6): the file is read once through a descriptor check (WF2), then judged on a read-only
+// handle by authority.JudgeAssertion (row 201 A3 §2.5: the parse WF3, the lookup WF4, WF5, the law
+// WF6), the judge the minter and the daemon share. It returns the verified approval, or the refusal
+// text. The caller has already admitted row at now and decoded its challenge (F4).
 func approveAssertion(ctx context.Context, credPath, path string, row store.MintRequestRow, now int64) (*approval.SignedApproval, string) {
 	file, err := readRegularFileCapped(path, approval.MaxPasskeyResponseBytes)
 	if err != nil { // WF2
 		return nil, fmt.Sprintf("--webauthn-assertion: %v", err)
 	}
-	resp, err := approval.ParseAssertionResponse(file)
-	if err != nil { // WF3
-		return nil, fmt.Sprintf("--webauthn-assertion %s is not a passkey assertion: %v", path, err)
-	}
 	ro, err := store.OpenCredentialsReadOnly(credPath)
 	if err != nil {
 		return nil, fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
 	}
-	a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)
+	signed, err := authority.JudgeAssertion(ctx, ro, row, file, now)
 	_ = ro.Close()
-	if err != nil { // WF4: below schema 6 names the serve restart
-		return nil, err.Error()
-	}
-	if !ok { // WF5
-		return nil, fmt.Sprintf("the assertion's credential cred %s is not an enrolled passkey; see ailang-worldd approver list", passkeyCred(resp.RawID))
+	if err != nil {
+		return nil, passkeyRefusal(err, "--webauthn-assertion "+path, row)
 	}
-	key, err := approval.ParseStoredCOSEKey(a.PublicKey)
-	if err != nil { // WF5: only a hand-edited row
-		return nil, fmt.Sprintf("approver %q's stored passkey key is not valid (%v); see ailang-worldd approver list", a.Name, err)
+	return signed, ""
+}
+
+// passkeyRefusal is the text of a JudgeAssertion refusal of source's assertion (WF3-WF6).
+func passkeyRefusal(err error, source string, row store.MintRequestRow) string {
+	var je *authority.JudgeError
+	if !errors.As(err, &je) {
+		return err.Error()
 	}
-	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
-	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
-	signed, err := approval.ApproveByAssertion(approval.AssertionInput{
-		Admit: approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Consumed,
-			StoredDigest: [32]byte(stored), BytesDigest: sha256.Sum256(row.Request)},
-		Challenge: [32]byte(challenge),
-		Response:  file,
-		Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,
-			Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,
-				BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}},
-	})
-	if err != nil { // WF6
-		return nil, fmt.Sprintf("the passkey assertion does not approve request %s: %v", row.RequestID, err)
+	switch je.Step {
+	case authority.JudgeStepParse: // WF3
+		return fmt.Sprintf("%s is not a passkey assertion: %v", source, je.Err)
+	case authority.JudgeStepUnknownKey: // WF5
+		return fmt.Sprintf("the assertion's credential cred %s is not an enrolled passkey; see ailang-worldd approver list", passkeyCred(je.RawID))
+	case authority.JudgeStepStoredKey: // WF5: only a hand-edited row
+		return fmt.Sprintf("approver %q's stored passkey key is not valid (%v); see ailang-worldd approver list", je.Approver, je.Err)
+	case authority.JudgeStepLaw: // WF6
+		return fmt.Sprintf("the passkey assertion does not approve request %s: %v", row.RequestID, je.Err)
 	}
-	return signed, ""
+	return je.Err.Error() // WF4: below schema 6 names the serve restart
 }
--- a/cmd/ailang-worldd/approve_passkey_now_test.go
+++ b/cmd/ailang-worldd/approve_passkey_now_test.go
@@ -0,0 +1,25 @@
+package main
+
+// Row 201 A3 sprint 1 M2b (planner, the mutant MUT-WP-JUDGE-NOW): through the shared judge, the
+// passkey law judges the approver at approve's own clock reading, not at the request's creation —
+// the passkey twin of TestApproveBySignatureJudgesTheApproverAtNow (A1b sprint 4 SIG-NOW).
+
+import (
+	"bytes"
+	"testing"
+)
+
+// PK-NOW: a passkey enrolled after the request was made (+5) and before the approve (+10) approves it.
+func TestApproveByPasskeyJudgesTheApproverAtNow(t *testing.T) {
+	a := newApproveRig(t, "ep1", approveT0)
+	seedPasskey(t, a.credPath, "mark-mac", passkeySigner, approveT0+5)
+	file := writeAssertion(t, passkeySigner.Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2))
+	var opened bool
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--webauthn-assertion", file), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+10))); code != exitOK || opened {
+		t.Fatalf("rc %d, terminal opened %v; want exitOK with no terminal for a passkey enrolled at +5 and an approve at +10\nstderr %s", code, opened, stderr.String())
+	}
+	if row := a.row(t); !row.Consumed || row.ApprovedBy != "webauthn:mark-mac" {
+		t.Fatalf("request consumed %v by %q, want webauthn:mark-mac", row.Consumed, row.ApprovedBy)
+	}
+}
--- a/cmd/ailang-worldd/approve_sig.go
+++ b/cmd/ailang-worldd/approve_sig.go
@@ -6,14 +6,15 @@ package main
 
 import (
 	"context"
-	"crypto/sha256"
 	"encoding/hex"
+	"errors"
 	"fmt"
 	"io"
 	"os"
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
 	"github.com/sunholo-data/ailang-world/host/sshsign"
 	"github.com/sunholo-data/ailang-world/host/store"
 )
@@ -37,39 +38,36 @@ func approveSignature(ctx context.Context, credPath, path string, row store.Mint
 }
 
 // approveArmored judges an armored signature from source (--signature's file or --sign-with's
-// signer): parsed (F11), its key looked up among the enrolled approvers (F12, F13), and judged by
-// the signature law against the stored row (F14).
+// signer) on a read-only handle by authority.JudgeSSHSignature (row 201 A3 §2.5: the parse F11, the
+// lookup F12, F13, the law F14), the judge the minter and the daemon share.
 func approveArmored(ctx context.Context, credPath, source string, armored []byte, row store.MintRequestRow, now int64) (*approval.SignedApproval, string) {
-	sig, err := approval.ParseSSHSIG(armored)
-	if err != nil { // F11
-		return nil, fmt.Sprintf("%s is not an SSH signature: %v", source, err)
-	}
 	ro, err := store.OpenCredentialsReadOnly(credPath)
 	if err != nil {
 		return nil, fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
 	}
-	enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)
+	signed, err := authority.JudgeSSHSignature(ctx, ro, row, armored, now)
 	_ = ro.Close()
-	if err != nil { // F12: below schema 5 names the serve restart
-		return nil, err.Error()
+	if err != nil {
+		return nil, sshRefusal(err, source, row)
 	}
-	if !ok { // F13
-		return nil, fmt.Sprintf("the signature's key %s is not an enrolled approver; see ailang-worldd approver list", sig.PublicKey.Fingerprint())
+	return signed, ""
+}
+
+// sshRefusal is the text of a JudgeSSHSignature refusal of source's signature (F11-F14).
+func sshRefusal(err error, source string, row store.MintRequestRow) string {
+	var je *authority.JudgeError
+	if !errors.As(err, &je) {
+		return err.Error()
 	}
-	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
-	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
-	signed, err := approval.ApproveBySignature(approval.SignatureInput{
-		Admit: approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Consumed,
-			StoredDigest: [32]byte(stored), BytesDigest: sha256.Sum256(row.Request)},
-		Challenge: [32]byte(challenge),
-		Signature: armored,
-		Approver: approval.Approver{Name: enrolled.Name, Format: enrolled.Format, Key: approval.SSHPublicKey{Blob: enrolled.PublicKey},
-			EnrolledAt: enrolled.EnrolledAt, LastUsedAt: enrolled.LastUsedAt, RevokedAt: enrolled.RevokedAt},
-	})
-	if err != nil { // F14
-		return nil, fmt.Sprintf("the signature does not approve request %s: %v", row.RequestID, err)
+	switch je.Step {
+	case authority.JudgeStepParse: // F11
+		return fmt.Sprintf("%s is not an SSH signature: %v", source, je.Err)
+	case authority.JudgeStepUnknownKey: // F13
+		return fmt.Sprintf("the signature's key %s is not an enrolled approver; see ailang-worldd approver list", je.Fingerprint)
+	case authority.JudgeStepLaw: // F14
+		return fmt.Sprintf("the signature does not approve request %s: %v", row.RequestID, je.Err)
 	}
-	return signed, ""
+	return je.Err.Error() // F12: below schema 5 names the serve restart
 }
 
 // ---- M5c ----
````

### Patch M3a — approve --await, the minter

````diff M3a
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -27,6 +27,7 @@ const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--
        ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --webauthn-link [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --webauthn-assertion <file> [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-id> --await [--await-for <duration>] [--db <path>] [--credentials-db <path>]
        ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
 
 ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
@@ -73,6 +74,13 @@ approve <request-id> (mr- and 16 hex):
   (at most 16 KiB): the assertion is checked against the passkey enrolled for
   its credential before any network call, its sign counter must move forward,
   and the request records approved_by=webauthn:<name>.
+--await replaces step 3's y/N with an approval made on the daemon's surface
+  (the request's card, /surface/approve/<id>): it prints the card's address,
+  then waits, with no terminal, until the surface has recorded a passkey
+  assertion or SSH signature for the request. That artifact is judged again
+  here, as --webauthn-assertion and --signature judge theirs; a refused one is
+  shown and the wait goes on. It ends when the request is approved, used,
+  expired, or --await-for (at least 1s) has passed; the request stays pending.
 --list prints the pending requests, read-only. A request expires 10 minutes
 after it is made; at most 16 are pending per credentials store.
 --list --decided prints the used requests instead, newest first (at most 64):
@@ -145,6 +153,8 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 	signWith := fs.String("sign-with", "", "approve by signing the request with ssh-keygen -Y sign -f <key> (no terminal)")
 	passkeyLink := fs.Bool("webauthn-link", false, "print the ceremony page's links for the request (no terminal, no write)")
 	assertion := fs.String("webauthn-assertion", "", "approve with an enrolled passkey's assertion file (no terminal)")
+	await := fs.Bool("await", false, "wait for an approval recorded by the daemon's surface, judge it, and mint (no terminal)")
+	awaitFor := fs.Duration("await-for", 0, "with --await: stop waiting after this long (default: the request's expiry)")
 	positional, err := parseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approveHelp)
@@ -159,12 +169,12 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 		return exitUsage
 	}
 	modes := 0
-	for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *passkeyLink, *assertion != ""} {
+	for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *passkeyLink, *assertion != "", *await} {
 		if on {
 			modes++
 		}
 	}
-	if modes > 1 { // F1: at most one of --signature, --sign-with, --write-message, --webauthn-link and --webauthn-assertion
+	if modes > 1 { // F1: at most one of --signature, --sign-with, --write-message, --webauthn-link, --webauthn-assertion and --await
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
@@ -176,6 +186,12 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
+	awaitForSet := false
+	fs.Visit(func(f *flag.Flag) { awaitForSet = awaitForSet || f.Name == "await-for" })
+	if awaitForSet && (!*await || *awaitFor < time.Second) { // AW5: --await-for is an --await, of at least 1s
+		fmt.Fprint(stderr, approveHelp)
+		return exitUsage
+	}
 	id := ""
 	if !*list {
 		id = positional[0]
@@ -253,6 +269,15 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 		return exitOK
 	}
 	var signed *approval.SignedApproval
+	if *await { // the minter: a surface approval, judged again (row 201 A3 §2.6)
+		if refusal := writeAwaitLinks(ctx, credPath, stdout, row, mr); refusal != "" {
+			return refuse("%s", refusal)
+		}
+		var refusal string
+		if signed, refusal = awaitSurfaceApproval(env, credPath, row, *awaitFor, stderr); refusal != "" {
+			return refuse("%s", refusal)
+		}
+	}
 	if *sigPath != "" { // the signature, judged before any network call (design §2.7 F9-F14)
 		var refusal string
 		if signed, refusal = approveSignature(ctx, credPath, *sigPath, row, now); refusal != "" {
--- a/cmd/ailang-worldd/approve_await.go
+++ b/cmd/ailang-worldd/approve_await.go
@@ -0,0 +1,132 @@
+package main
+
+// Row 201 A3 sprint 1 M3 (design_docs/planned/w-world-approvals-by-key-a3.md §2.6): `approve
+// <request-id> --await [--await-for <duration>]`, the minter. It opens no terminal: it waits for an
+// approval the daemon's surface recorded (credentials schema 7, surface_approvals), judges that
+// artifact again with the same law (authority.JudgeAssertion, JudgeSSHSignature) against the
+// request it reads in the same poll, and only a verified one goes on to approve's steps 2-6. The
+// row's approver and digest columns are never trusted.
+
+import (
+	"context"
+	"encoding/hex"
+	"fmt"
+	"io"
+	"net"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// awaitPollInterval is the minter's pause between polls (M1), through env.after.
+const awaitPollInterval = time.Second
+
+// writeAwaitLinks prints the minter's start (design §2.6): the request, its message hash and
+// expiry, and the card's address for each distinct active passkey origin and the request's own
+// loopback daemon. Nothing it prints is secret. It returns the refusal text, or "".
+func writeAwaitLinks(ctx context.Context, credPath string, stdout io.Writer, row store.MintRequestRow, mr approval.MintRequest) string {
+	ro, err := store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		return fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
+	}
+	active, err := ro.ActivePasskeys(ctx)
+	_ = ro.Close()
+	if err != nil { // below schema 6 names the serve restart
+		return err.Error()
+	}
+	h, err := requestHash(row)
+	if err != nil {
+		return fmt.Sprintf("request %s: %v", row.RequestID, err)
+	}
+	fmt.Fprintf(stdout, "waiting for an approval of request %s (expires %s, message hash %s) from the surface; open its card and approve it there:\n",
+		row.RequestID, approverTime(row.ExpiresAt), hex.EncodeToString(h[:8]))
+	seen := map[string]bool{}
+	for _, e := range active {
+		if !seen[e.Passkey.Origin] {
+			seen[e.Passkey.Origin] = true
+			fmt.Fprintf(stdout, "  open %s/surface/approve/%s\n", e.Passkey.Origin, row.RequestID)
+		}
+	}
+	if _, port, err := net.SplitHostPort(mr.Daemon.Addr); err == nil && !seen["http://localhost:"+port] {
+		fmt.Fprintf(stdout, "  open http://localhost:%s/surface/approve/%s\n", port, row.RequestID)
+	}
+	return ""
+}
+
+// awaitPoll is one poll's reading: the request row, its recorded approval if any, and that
+// approval's judgement (a verified approval, or the refusal text).
+type awaitPoll struct {
+	row      store.MintRequestRow
+	rec      store.SurfaceApprovalRow
+	recorded bool
+	signed   *approval.SignedApproval
+	refusal  string
+}
+
+// pollSurfaceApproval is one poll (M1, M3): one read-only open under its own 3 s root, the request
+// row and its recorded approval, and the approval's artifact judged against that row at now.
+func pollSurfaceApproval(credPath, id string, now int64) (awaitPoll, error) {
+	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		return awaitPoll{}, err
+	}
+	defer func() { _ = ro.Close() }()
+	var p awaitPoll
+	var ok bool
+	if p.row, ok, err = ro.MintRequest(ctx, id); err != nil {
+		return awaitPoll{}, err
+	}
+	if !ok {
+		return awaitPoll{}, fmt.Errorf("request %s is no longer in the store", id)
+	}
+	if p.rec, p.recorded, err = ro.SurfaceApproval(ctx, id); err != nil || !p.recorded {
+		return p, err
+	}
+	var jerr error
+	if p.rec.Format == "ssh" { // M3: the artifact, judged again; never the approver column
+		if p.signed, jerr = authority.JudgeSSHSignature(ctx, ro, p.row, p.rec.Artifact, now); jerr != nil {
+			p.refusal = sshRefusal(jerr, "the surface's approval", p.row)
+		}
+	} else if p.signed, jerr = authority.JudgeAssertion(ctx, ro, p.row, p.rec.Artifact, now); jerr != nil {
+		p.refusal = passkeyRefusal(jerr, "the surface's approval", p.row)
+	}
+	return p, nil
+}
+
+// awaitSurfaceApproval is the minter's wait (M1-M6), after admission and F4: a poll every second
+// until the request has a recorded approval that the law verifies (returned), or the request is
+// used elsewhere or expires (admission's own text), or the wait reaches now + awaitFor (0: the
+// expiry). A refused approval is printed once per received_at, and the wait goes on.
+func awaitSurfaceApproval(env sessionEnv, credPath string, row store.MintRequestRow, awaitFor time.Duration, stderr io.Writer) (*approval.SignedApproval, string) {
+	deadline := row.ExpiresAt
+	if awaitFor > 0 {
+		deadline = min(deadline, env.now()+int64(awaitFor/time.Second))
+	}
+	shown := int64(-1)
+	for {
+		now := env.now()
+		p, err := pollSurfaceApproval(credPath, row.RequestID, now)
+		if err != nil {
+			return nil, fmt.Sprintf("cannot read request %s: %v", row.RequestID, err)
+		}
+		if _, refusal := approveAdmission(p.row, now); refusal != "" { // M2: used elsewhere, expired
+			return nil, refusal
+		}
+		if p.signed != nil { // M5
+			return p.signed, ""
+		}
+		if p.recorded && p.rec.ReceivedAt != shown { // M4
+			fmt.Fprintf(stderr, "ailang-worldd approve: refused the surface's approval: %s; waiting for another\n", p.refusal)
+			shown = p.rec.ReceivedAt
+		}
+		if now >= deadline { // M6
+			return nil, fmt.Sprintf("no approval for request %s by %s; it stays pending until %s",
+				row.RequestID, approverTime(now), approverTime(row.ExpiresAt))
+		}
+		<-env.after(awaitPollInterval)
+	}
+}
--- a/cmd/ailang-worldd/approve_await_test.go
+++ b/cmd/ailang-worldd/approve_await_test.go
@@ -0,0 +1,176 @@
+package main
+
+// Row 201 A3 sprint 1 M3a (design_docs/planned/w-world-approvals-by-key-a3.md §2.6, §6 AW1, AW5):
+// approve --await, the minter. Every store is under t.TempDir(); every key is approvaltest's;
+// the surface's recording is the store's RecordSurfaceApproval, called by the test mid-wait
+// (D-WORLD-109: no daemon route writes the table in this sprint).
+
+import (
+	"bytes"
+	"context"
+	"encoding/hex"
+	"net"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// awaitClock is approve --await's fake clock: env.after advances now by its duration and counts
+// the waits, then calls onWait(waits, now) before it fires.
+type awaitClock struct {
+	now, waits int64
+	onWait     func(n, now int64)
+}
+
+// awaitEnv is the minter's env at start: the fake clock, and a terminal that records into *opened
+// and would answer y, so any path that reaches the typed y is seen (AW1, AW7).
+func awaitEnv(start int64, opened *bool, onWait func(n, now int64)) (sessionEnv, *awaitClock) {
+	c := &awaitClock{now: start, onWait: onWait}
+	env := approveEnv(&promptTerm{in: strings.NewReader("y\n")}, opened, func() int64 { return c.now })
+	env.after = func(d time.Duration) <-chan time.Time {
+		c.waits++
+		c.now += int64(d / time.Second)
+		if c.onWait != nil {
+			c.onWait(c.waits, c.now)
+		}
+		ch := make(chan time.Time, 1)
+		ch <- time.Unix(c.now, 0)
+		return ch
+	}
+	return env, c
+}
+
+// recordSurface records artifact for the rig's request as the surface would, at now.
+func recordSurface(t *testing.T, a approveRig, format, approver string, artifact []byte, now int64) {
+	t.Helper()
+	st, err := store.OpenCredentials(a.credPath)
+	if err != nil {
+		t.Error(err)
+		return
+	}
+	defer func() { _ = st.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	if err := st.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: a.id, Format: format, Approver: approver, Artifact: artifact, Digest: a.row(t).Digest}, now); err != nil {
+		t.Error(err)
+	}
+}
+
+// requireMintedBy fails unless the rig's request was consumed by approvedBy with the effects the
+// typed y makes (CL16), and stdout shows the approver.
+func requireMintedBy(t *testing.T, a approveRig, approvedBy, stdout string) {
+	t.Helper()
+	row := a.row(t)
+	if !row.Consumed || row.ApprovedBy != approvedBy || len(a.sessionRows(t)) != 1 {
+		t.Fatalf("request consumed %v by %q, %d sessions; want one, by %s", row.Consumed, row.ApprovedBy, len(a.sessionRows(t)), approvedBy)
+	}
+	if !strings.Contains(stdout, "approved by "+approvedBy+"\n") {
+		t.Fatalf("stdout lacks the approver:\n%s", stdout)
+	}
+	mr, err := approval.Decode(row.Request)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var probs []string
+	for _, p := range approveEffectProblems(mr, collectApproveFacts(t, a, mr, stdout)) {
+		if !strings.HasPrefix(p, "approved_by:") { // the typed check wants tty
+			probs = append(probs, p)
+		}
+	}
+	if len(probs) != 0 {
+		t.Fatalf("the effects differ from the request:\n%s", strings.Join(probs, "\n"))
+	}
+}
+
+// AW1: a passkey assertion recorded mid-wait (after the third wait) is judged and mints with no
+// terminal: one session, approved_by webauthn:mark-mac, the counter advanced, the token in the
+// request's session file; the minter printed the request's message hash and the card's address
+// once for the two passkeys' one origin, which is also the request's loopback daemon, and polled
+// exactly four times.
+func TestApproveAwaitMintsARecordedPasskeyApproval(t *testing.T) {
+	a := passkeyRig(t)
+	seedPasskey(t, a.credPath, "b-other", approvaltest.NewPasskey(approval.AlgES256, 52), approveT0-100)
+	h := requestHashOf(t, a)
+	var opened bool
+	env, c := awaitEnv(approveT0+1, &opened, func(n, now int64) {
+		if n == 3 {
+			recordSurface(t, a, "webauthn", "webauthn:mark-mac", passkeySigner.Assert(pkOrigin, pkRPID, h, 2), now)
+		}
+	})
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitOK || opened {
+		t.Fatalf("rc %d, terminal opened %v; want exitOK with no terminal\nstderr %s", code, opened, stderr.String())
+	}
+	if c.waits != 3 {
+		t.Fatalf("%d waits, want 3 (the fourth poll finds the approval)", c.waits)
+	}
+	requireMintedBy(t, a, "webauthn:mark-mac", stdout.String())
+	if p := passkeyOf(t, a.credPath, "mark-mac"); p.SignCount != 2 {
+		t.Fatalf("sign count %d, want 2", p.SignCount)
+	}
+	mr, err := approval.Decode(a.row(t).Request)
+	if err != nil {
+		t.Fatal(err)
+	}
+	_, port, err := net.SplitHostPort(mr.Daemon.Addr)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if pkOrigin != "http://localhost:"+port {
+		t.Fatalf("the rig's passkey origin %s is not the request's loopback daemon %s", pkOrigin, mr.Daemon.Addr)
+	}
+	for _, want := range []string{"waiting for an approval of request " + a.id + " (expires " + approverTime(a.row(t).ExpiresAt) + ", message hash " + hex.EncodeToString(h[:8]) + ") from the surface",
+		"\n  open " + pkOrigin + "/surface/approve/" + a.id + "\n"} {
+		if !strings.Contains(stdout.String(), want) {
+			t.Fatalf("stdout lacks %q:\n%s", want, stdout.String())
+		}
+	}
+	if n := strings.Count(stdout.String(), "/surface/approve/"); n != 1 {
+		t.Fatalf("%d card addresses, want 1 (one origin):\n%s", n, stdout.String())
+	}
+	if strings.Contains(stderr.String(), "refused the surface's approval") {
+		t.Fatalf("a verified approval was refused:\n%s", stderr.String())
+	}
+}
+
+// AW5: --await is a sixth exclusive mode (F1); --await-for needs --await and at least 1s; --list
+// refuses --await (F2). Each prints the help and exits exitUsage before any store is read.
+func TestApproveAwaitFlagRules(t *testing.T) {
+	a := newApproveRig(t, "ep1", approveT0)
+	for _, extra := range [][]string{
+		{"--await", "--signature", "x.sig"}, {"--await", "--sign-with", "k"}, {"--await", "--write-message", "m"},
+		{"--await", "--webauthn-link"}, {"--await", "--webauthn-assertion", "a.json"},
+		{"--await-for", "30s"}, {"--await", "--await-for", "0s"}, {"--await", "--await-for", "999ms"}, {"--await", "--await-for", "-5s"},
+	} {
+		var opened bool
+		var stdout, stderr bytes.Buffer
+		env, _ := awaitEnv(approveT0+1, &opened, nil) // a fake clock: a rule that let a mode through waits at most to the expiry
+		if code := runApprove(a.args(extra...), &stdout, &stderr, env); code != exitUsage || opened ||
+			!strings.HasPrefix(stderr.String(), "usage: ailang-worldd approve") {
+			t.Errorf("%v: rc %d, terminal %v, stderr %q; want exitUsage and the help", extra, code, opened, stderr.String())
+		}
+	}
+	var stdout, stderr bytes.Buffer
+	if code := runApprove([]string{"--list", "--await", "--db", a.r.f.db}, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+1))); code != exitUsage ||
+		!strings.HasPrefix(stderr.String(), "usage: ailang-worldd approve") {
+		t.Errorf("--list --await: rc %d, stderr %q; want exitUsage and the help", code, stderr.String())
+	}
+	if !strings.Contains(approveHelp, "ailang-worldd approve <request-id> --await [--await-for <duration>] [--db <path>] [--credentials-db <path>]\n") {
+		t.Error("the help lacks the --await usage line")
+	}
+	if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
+		t.Fatal("a refused flag combination consumed the request or minted")
+	}
+	// Control: --await-for 1s with --await is accepted (no approval is recorded, so it waits one
+	// second and stops, naming the time).
+	env, c := awaitEnv(approveT0+1, nil, nil)
+	stderr.Reset()
+	if code := runApprove(a.args("--await", "--await-for", "1s"), &stdout, &stderr, env); code != exitUsage || c.waits != 1 ||
+		!strings.Contains(stderr.String(), "no approval for request "+a.id+" by "+approverTime(approveT0+2)+"; it stays pending until ") {
+		t.Fatalf("control: rc %d, %d waits, stderr %q", code, c.waits, stderr.String())
+	}
+}
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -41,7 +41,7 @@ func TestEveryVerbHasHelp(t *testing.T) {
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
 		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
-		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve <request-id> --webauthn-link", "approve <request-id> --webauthn-assertion <file>", "approve --list [--decided]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
+		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve <request-id> --webauthn-link", "approve <request-id> --webauthn-assertion <file>", "approve <request-id> --await [--await-for <duration>]", "approve --list [--decided]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
 			t.Errorf("top-level usage lacks %q", want)
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -127,6 +127,7 @@ Usage:
   ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --webauthn-link [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --webauthn-assertion <file> [--db <path>] [--credentials-db <path>]
+  ailang-worldd approve <request-id> --await [--await-for <duration>] [--db <path>] [--credentials-db <path>]
   ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format webauthn --origin <origin> --registration <file> [--db <path>] [--credentials-db <path>]
--- a/host/store/context_roots_test.go
+++ b/host/store/context_roots_test.go
@@ -26,6 +26,8 @@ var contextRootPins = map[string]int{
 	// Row 201 A2 sprint 3: the read-only pre-check (3 s), the re-check after the enrolment code (3 s)
 	// and the write after y (3 s).
 	"cmd/ailang-worldd/approver_passkey.go|runApproverEnrollPasskey|Background": 3,
+	// Row 201 A3 sprint 1 M3: each approve --await poll's read-only reads and judgement (3 s).
+	"cmd/ailang-worldd/approve_await.go|pollSurfaceApproval|Background": 1,
 
 	"cmd/ailang-worldd/cli.go|get|Background":                                1,
 	"cmd/ailang-worldd/cli.go|execute|Background":                            1,
````

### Patch M3b — the minter's refusals and stops

````diff M3b
--- a/cmd/ailang-worldd/approve_await_refusals_test.go
+++ b/cmd/ailang-worldd/approve_await_refusals_test.go
@@ -0,0 +1,230 @@
+package main
+
+// Row 201 A3 sprint 1 M3b (design_docs/planned/w-world-approvals-by-key-a3.md §2.6, §6 AW2-AW4,
+// AW6, AW7): the minter re-judges every recorded approval, stops at the request's expiry, at
+// --await-for and when the request is used elsewhere, and mints an SSH signature as ssh:<name>.
+
+import (
+	"bytes"
+	"net"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// AW2: an assertion whose counter did not move is refused once (per received_at, however many
+// polls see it), the wait goes on, and a replacement recorded later mints.
+func TestApproveAwaitRejudgesAndWaits(t *testing.T) {
+	a := passkeyRig(t)
+	h := requestHashOf(t, a)
+	env, c := awaitEnv(approveT0+1, nil, func(n, now int64) {
+		switch n {
+		case 1:
+			recordSurface(t, a, "webauthn", "webauthn:mark-mac", passkeySigner.Assert(pkOrigin, pkRPID, h, 1), now)
+		case 4:
+			recordSurface(t, a, "webauthn", "webauthn:mark-mac", passkeySigner.Assert(pkOrigin, pkRPID, h, 2), now)
+		}
+	})
+	var opened bool
+	env.openTerminal = approveEnv(nil, &opened, nil).openTerminal
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitOK || opened || c.waits != 4 {
+		t.Fatalf("rc %d, terminal %v, %d waits; want exitOK, no terminal, 4 waits\nstderr %s", code, opened, c.waits, stderr.String())
+	}
+	want := "ailang-worldd approve: refused the surface's approval: the passkey assertion does not approve request " + a.id + ": " +
+		approval.ErrSignCountRegression.Error() + ": enrolled 1, this assertion 1; if you did not cause this, run: ailang-worldd approver revoke mark-mac; waiting for another\n"
+	if !strings.HasPrefix(stderr.String(), want) || strings.Count(stderr.String(), "refused the surface's approval") != 1 {
+		t.Fatalf("stderr %q, want exactly one refusal, first: %q", stderr.String(), want)
+	}
+	requireMintedBy(t, a, "webauthn:mark-mac", stdout.String())
+}
+
+// AW3: with nothing recorded the wait ends exactly at the request's expiry (admission's text) or
+// exactly at --await-for, exitUsage, nothing minted, the request still pending in the second case.
+func TestApproveAwaitStopsAtExpiry(t *testing.T) {
+	t.Run("expiry", func(t *testing.T) {
+		a := newApproveRig(t, "ep1", approveT0)
+		row := a.row(t)
+		env, c := awaitEnv(approveT0+1, nil, nil)
+		var stdout, stderr bytes.Buffer
+		if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitUsage {
+			t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
+		}
+		if c.now != row.ExpiresAt || c.waits != row.ExpiresAt-approveT0-1 {
+			t.Fatalf("stopped at %d after %d waits, want exactly the expiry %d after %d", c.now, c.waits, row.ExpiresAt, row.ExpiresAt-approveT0-1)
+		}
+		if want := "ailang-worldd approve: " + expiredRefusal(row) + "\n"; stderr.String() != want {
+			t.Fatalf("stderr %q, want %q", stderr.String(), want)
+		}
+		if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
+			t.Fatal("an expired wait consumed the request or minted")
+		}
+	})
+	t.Run("await-for", func(t *testing.T) {
+		a := newApproveRig(t, "ep1", approveT0)
+		env, c := awaitEnv(approveT0+1, nil, nil)
+		var stdout, stderr bytes.Buffer
+		if code := runApprove(a.args("--await", "--await-for", "30s"), &stdout, &stderr, env); code != exitUsage {
+			t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
+		}
+		if c.now != approveT0+31 || c.waits != 30 {
+			t.Fatalf("stopped at %d after %d waits, want %d after 30", c.now, c.waits, approveT0+31)
+		}
+		want := "ailang-worldd approve: no approval for request " + a.id + " by " + approverTime(approveT0+31) +
+			"; it stays pending until " + approverTime(a.row(t).ExpiresAt) + "\n"
+		if stderr.String() != want {
+			t.Fatalf("stderr %q, want %q", stderr.String(), want)
+		}
+		if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
+			t.Fatal("a timed-out wait consumed the request or minted")
+		}
+	})
+}
+
+// AW4: a consume by another process mid-wait ends the wait with admission's "already used" text.
+func TestApproveAwaitStopsWhenUsedElsewhere(t *testing.T) {
+	a := newApproveRig(t, "ep1", approveT0)
+	env, c := awaitEnv(approveT0+1, nil, func(n, now int64) {
+		if n == 2 {
+			raw := rawCredentials(t, a.credPath)
+			defer func() { _ = raw.Close() }()
+			if _, err := raw.Exec(`UPDATE mint_requests SET consumed_at = ?, approved_by = 'tty', credential_id = 'elsewhere' WHERE request_id = ?`, now, a.id); err != nil {
+				t.Error(err)
+			}
+		}
+	})
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitUsage || c.waits != 2 {
+		t.Fatalf("rc %d after %d waits, want exitUsage after 2\nstderr %s", code, c.waits, stderr.String())
+	}
+	if !strings.Contains(stderr.String(), "was already used at "+approverTime(approveT0+3)+" by tty (credential elsewhere)") {
+		t.Fatalf("stderr %q lacks the already-used refusal", stderr.String())
+	}
+	if len(a.sessionRows(t)) != 0 {
+		t.Fatal("the minter minted a request used elsewhere")
+	}
+}
+
+// AW6: a recorded SSHSIG by an enrolled key mints as ssh:mark-mac, with no terminal; with no
+// passkey enrolled the card's only address is the request's loopback daemon.
+func TestApproveAwaitSSHSIG(t *testing.T) {
+	a, _ := signedRig(t)
+	sig := sigSigner.Sign(approval.Namespace, requestMessage(t, a))
+	var opened bool
+	env, _ := awaitEnv(approveT0+1, &opened, func(n, now int64) {
+		if n == 1 {
+			recordSurface(t, a, "ssh", "ssh:mark-mac", sig, now)
+		}
+	})
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--await"), &stdout, &stderr, env); code != exitOK || opened {
+		t.Fatalf("rc %d, terminal opened %v; want exitOK with no terminal\nstderr %s", code, opened, stderr.String())
+	}
+	requireMintedBy(t, a, "ssh:mark-mac", stdout.String())
+	mr, err := approval.Decode(a.row(t).Request)
+	if err != nil {
+		t.Fatal(err)
+	}
+	_, port, err := net.SplitHostPort(mr.Daemon.Addr)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if n := strings.Count(stdout.String(), "/surface/approve/"); n != 1 || !strings.Contains(stdout.String(), "\n  open http://localhost:"+port+"/surface/approve/"+a.id+"\n") {
+		t.Fatalf("want exactly the loopback card address:\n%s", stdout.String())
+	}
+}
+
+// AW8 (planner): a store below schema 7 ends the wait at the first poll naming the serve restart
+// (below schema 6, before it, naming the passkeys' restart), and a request deleted mid-wait ends
+// it naming that; nothing is minted.
+func TestApproveAwaitNeedsSchema7AndItsRequest(t *testing.T) {
+	for _, arm := range []struct {
+		name  string
+		setup []string // before the run
+		gone  bool     // delete the request after the first wait
+		want  string
+	}{
+		{name: "v6-store", setup: []string{`DROP TABLE surface_approvals`, `PRAGMA user_version = 6`}, want: "surface approvals need schema 7: restart `ailang-worldd serve`"},
+		{name: "v5-store", setup: []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, want: "passkeys need schema 6"},
+		{name: "deleted", gone: true, want: "is no longer in the store"},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			a := newApproveRig(t, "ep1", approveT0)
+			raw := rawCredentials(t, a.credPath)
+			for _, s := range arm.setup {
+				if _, err := raw.Exec(s); err != nil {
+					t.Fatal(err)
+				}
+			}
+			_ = raw.Close()
+			var opened bool
+			env, c := awaitEnv(approveT0+1, &opened, func(n, now int64) {
+				if n == 1 && arm.gone {
+					raw := rawCredentials(t, a.credPath)
+					defer func() { _ = raw.Close() }()
+					if _, err := raw.Exec(`DELETE FROM mint_requests WHERE request_id = ?`, a.id); err != nil {
+						t.Error(err)
+					}
+				}
+			})
+			var stdout, stderr bytes.Buffer
+			code := runApprove(a.args("--await"), &stdout, &stderr, env)
+			if code != exitUsage || opened || !strings.Contains(stderr.String(), arm.want) {
+				t.Fatalf("rc %d, terminal %v, stderr %q; want exitUsage naming %q", code, opened, stderr.String(), arm.want)
+			}
+			if want := map[bool]int64{false: 0, true: 1}[arm.gone]; c.waits != want {
+				t.Fatalf("%d waits, want %d", c.waits, want)
+			}
+			if len(a.sessionRows(t)) != 0 {
+				t.Fatal("minted")
+			}
+		})
+	}
+}
+
+// AW7: a row the minter did not get from a verified POST — an unenrolled key's valid signature, a
+// passkey assertion over another request, an approver column naming mark-mac over garbage — is
+// judged, refused once, and never mints: no terminal, no session, the request pending.
+func TestApproveAwaitNeverTrustsTheRow(t *testing.T) {
+	for _, arm := range []struct {
+		name, format, approver string
+		artifact               func(t *testing.T, a approveRig) []byte
+		want                   string
+	}{
+		{"unenrolled-key", "ssh", "ssh:mark-mac", func(t *testing.T, a approveRig) []byte {
+			return otherSigner.Sign(approval.Namespace, requestMessage(t, a))
+		}, "the signature's key " + otherSigner.PublicKey().Fingerprint() + " is not an enrolled approver"},
+		{"other-request", "webauthn", "webauthn:mark-mac", func(t *testing.T, a approveRig) []byte {
+			h := requestHashOf(t, a)
+			h[31] ^= 1
+			return passkeySigner.Assert(pkOrigin, pkRPID, h, 2)
+		}, "the passkey assertion does not approve request "},
+		{"garbage-named-mark", "webauthn", "webauthn:mark-mac", func(*testing.T, approveRig) []byte { return []byte(`{"id":"x"}`) },
+			"the surface's approval is not a passkey assertion: " + approval.ErrPasskeyResponse.Error()},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			a := passkeyRig(t)
+			seedApprover(t, a.credPath, "mark-ssh", sigSigner.PublicKey(), approveT0-100)
+			art := arm.artifact(t, a)
+			var opened bool
+			env, _ := awaitEnv(approveT0+1, &opened, func(n, now int64) {
+				if n == 1 {
+					recordSurface(t, a, arm.format, arm.approver, art, now)
+				}
+			})
+			var stdout, stderr bytes.Buffer
+			code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, env)
+			if code != exitUsage || opened {
+				t.Fatalf("rc %d, terminal opened %v; want exitUsage, no terminal\nstderr %s", code, opened, stderr.String())
+			}
+			if n := strings.Count(stderr.String(), "refused the surface's approval: "); n != 1 || !strings.Contains(stderr.String(), arm.want) ||
+				!strings.Contains(stderr.String(), "no approval for request "+a.id) {
+				t.Fatalf("stderr %q: want one refusal naming %q, then the timeout", stderr.String(), arm.want)
+			}
+			if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
+				t.Fatal("the minter trusted the row")
+			}
+		})
+	}
+}
````
