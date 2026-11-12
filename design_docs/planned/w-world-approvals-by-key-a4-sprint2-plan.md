# Sprint plan — w-world-approvals-by-key-a4, sprint 2 "the keys", headed by the typed CLI (iteration 390, row 201)

> **THE SPRINT — A4 sprint 2 (design §0 "Sprint 2", §2.7, §2.10, §6 "Sprint 2", §9 "Sprint 2 = M4–M5"), headed by
> sprint 1's moved typed CLI (sprint-1 plan Appendix C):** `ailang-worldd approve <request-ref> [--deny]` at a
> terminal (CD1–CD4) — the shipped client that makes sprint 1's `POST /surface/decide/{ref}` reachable by an
> operator; the sprint-1 judge's test gaps J07, J31/J32, J33; the pure decision-signature law
> (`DecideBySignature`, `SignedDecision`); credentials **schema 8** `key_decisions` with the design's DDL verbatim
> (both `sha256:*` CHECKs) and its writer (`RecordKeyDecision` K0–K5, `LandKeyDecision`); the decision judge
> `JudgeDecisionSSH`; the POST's **`sshsig` arm** in sprint 1's arm slot (judge → non-waiting writer →
> `RecordKeyDecision` → `DecideApproval` → `LandKeyDecision`); `approve <request-ref> --sign-with <key>`.
> - **Measured: 1,983 changed lines under `cmd/ host/`** — `git diff --numstat <base> HEAD -- cmd host`,
>   insertions + deletions, after a commit of each of the ten scratch stages (**738 production** incl. the 15-line
>   schema file, **1,245 test**). Margin to the 2,000 cap: **17**.
> - **The design's priced fallback fired, from the measured size (S2-D1):** the **webauthn arm** (RD10 passkey
>   half, CD8–CD9, DM7, the passkey halves of K0/KD7/KD8/RD13/RD14, and KD2 — the SQL counter is passkey-only)
>   moves to sprint 3; **SSHSIG ships alone.** Measured at the design's price the full sprint is ≈ 2,200; with the
>   passkey arm out the built set measured 2,019 *with* `--write-message`/`--signature` on a request ref (CD6–CD7),
>   so those two flags (−63 measured, prod + test) move to sprint 3 too. Never cut and built: CD1–CD4, R4e, RD13 (a),
>   KD5/KD6, schema 8, the judge, the J-gap tests.
> - **Schema 8 is forward-compatible (S2-D3):** the table is the design's DDL verbatim and format-agnostic
>   (`format IN ('ssh','webauthn')`, `artifact BLOB`); the passkey counter lives in the existing
>   `webauthn_credentials.sign_count`. Sprint 3's passkey arm adds Go code only — **no schema 9**.
> - **How the CLI finds the daemon for a request ref (S2-D2):** `--daemon <host:port>` (default: the installed
>   service's address, `defaultServiceAddr()`), then `GET /v1/health` must name the same `--db` and the same
>   credentials store (`store.SameDatabase`), because the typed proof is written beside the served credentials
>   store — `episode finish`'s F6 rule.
> - **Four premises of the design are false at HEAD and are routed minimally** (§1.4): the production SSH signer
>   hard-codes the mint namespace (P1 — found by this plan's real-`ssh-keygen` sibling test, routed by
>   `sshsign.SignIn`); schema 8 breaks ten existing downgrade fixtures across three packages, several in the
>   design's "unedited" list (P2, the A3a-sprint-1 precedent); sprint 1's body grammar must change a struct tag
>   (P4); `approve` refuses the global `--addr` (P3).
> - **Every line of code and test is a unified diff in Appendix P**, generated from compiled, vetted, gofmt-clean,
>   green scratch stages; extracted **from this file** with §2's `awk`, applied to a fresh `git archive 8338eebb`
>   export and compared byte-for-byte with the stage tree (§1.2 C25). The executor applies each patch and runs
>   its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a4.md](w-world-approvals-by-key-a4.md), **r2** (reviewed in iteration 388;
  **not reopened**). **D-WORLD-115** (the decision message with a derived challenge) is OPEN; its default **A = keep
  it** is applied: built as designed. Where sprint 2 depends on the answer: §1.3 S2-D9. Where this plan and the
  design differ, this plan was measured against `8338eebb` and wins; each difference is a decision (§1.3) or a
  premise found false (§1.4).
- **Starting ground**: A4 sprint 1 landed (PR #563, merge `338ff821`; judge sonnet PASS 88) and is recorded
  (`8338eebb`, PR #564). Sprint-1 plan: [w-world-approvals-by-key-a4-sprint1-plan.md](w-world-approvals-by-key-a4-sprint1-plan.md)
  (its Appendix C is this sprint's head).
- **Planner base**: `8338eebb173290168f07e6b426d461c0122f67a8` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter390-a4-s2`, branch `sprint/world-iter390-row201-a4-s2`.
  Planner scratch (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter390/plan-scratch/`
  (`base/` the pristine export, `work/` the tree the mutants were cut from (tag `FINAL`), `work2/` the stage repo
  (tags `M4a` … `M5c`), `patches/`, `tools/` (`gates.sh`, `accept.sh`, `stages.sh`, `mutants.py`, `mutate.py`,
  `roundtrip.sh`), `logs/`, `roundtrip/`, `mut/`).
- **Executor**: as routed by the controller. **Three runs**: Run 1 (M4a, M4b, M4c — the J gaps and the typed
  CLI), Run 2 (M4d, M4e, M4f, M4g — the law, schema 8, the judge), Run 3 (M5a, M5b, M5c — the POST's key arm and
  `--sign-with`). Each sub-milestone is one `patch -p1` (§2) plus its acceptance (§3). **Apply, never retype.**
- **Touched** (§5): new `cmd/ailang-worldd/approve_decide.go` (+ `approve_decide_rig_test.go`,
  `approve_decide_test.go`, `approve_decide_key_test.go`), `cmd/ailang-worldd/{approve.go, main.go, session.go,
  help_test.go}`, new `host/approval/decision_signed.go` (+test), new `host/authority/judge_decision.go` (+test),
  new `host/store/key_decisions.go` (+ `key_decisions_test.go`, `key_decisions_record_test.go`),
  `host/store/{credentials.go, credentials_schema.sql}`, new `host/daemon/surface_decide_key.go` (+
  `surface_decide_key_test.go`, `surface_decide_key_retry_test.go`, `surface_decide_jgaps_test.go`),
  `host/daemon/surface_decide.go`, `host/surfacepage/decide_model.go`, `host/sshsign/sshsign.go`; edited tests
  (§1.4 P2, the inventories): `host/store/{credentials_test.go, approvers_test.go, mint_requests_test.go,
  surface_approvals_test.go, context_roots_test.go}`, `host/authority/judge_test.go`,
  `host/daemon/{surface_approve_post_test.go, surface_approve_read_test.go, surface_approve_post_guard_test.go,
  surface_decide_discovery_test.go}`, `host/broker/decision_law_test.go`, seven `cmd/ailang-worldd/*_test.go`
  downgrade fixtures. **Outside the design's §8 file list**: `host/sshsign/sshsign.go` and
  `cmd/ailang-worldd/session.go` (P1), every `_test.go` above, `cmd/ailang-worldd/main.go` (one usage line).
- **Not touched**: `host/approval/{request,signed,sshsig,passkey_assert}.go` (the mint law), `host/authority/judge.go`
  (the mint judge), `host/broker/**` production (sprint 1's law and `approve.go`), `host/coordinator/**`,
  `host/daemon/{daemon.go, episode_finish.go, surface_approve.go, surface_approve_guard.go, surface_page.go}`,
  `host/surfacepage/{approve_model.go, app/**, approve/**, passkey/**}`, `docs/QUICKSTART.md` (§18 is sprint 3's),
  `world/*.ail`, `packages/**`, `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, any skill, `~/.ssh`,
  `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys and temp stores only (D-WORLD-109 = A).** Every store is under `t.TempDir()`; every SSH key is
   `approvaltest.NewSigner(seed)`, except CD5's real-tool sibling, which makes an ed25519 key with `ssh-keygen -q
   -N ""` in a `t.TempDir()` under a temp `HOME` with `SSH_AUTH_SOCK=` (the A1b XS1 pattern) and signs only that
   key. The CLI tests talk to an `httptest` loopback server; nothing opens `~/.ailang/world`, `~/.ssh`, a real
   store or a real `/dev/tty`; no mint into a real store; no publish.
2. **Pure core, effects at the boundary (coding standards S1–S3).** `host/approval/decision_signed.go` is pure
   and standard-library only (`TestApprovalImportsOnlyTheStandardLibrary`); `now` is an input. The judge reads the
   credentials store through a read-only handle; the daemon opens the writer for two transactions and one world
   write only, under `humanMu`, and never keeps it (`TestDaemonCredentialWriterSites` pins the two sites). No
   `.ail` (design §3 S1); no package (S3).
3. **`decidedBy` is derived, never accepted.** The grammar's strict decode still refuses `decidedBy` (RD11);
   the key arm passes `SignedDecision.DecidedBy()` (`ssh:<name>`), the typed arm the constant.
4. **No `/v1` route, no new route at all, no CSP change, no CORS header, no new Go dependency.**
5. **One mutant per new production branch, anchored to the diff** (§4), each run with its named test and its
   package suite (inventory and static tests excluded); named equivalents and survivors with reasons; every run
   appends a ledger line.
6. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
7. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet`
   green, the `_test.go` compile fence green, its own tests and every earlier one green (§1.2 C23).
8. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `8338eebb` (measured by the planner, 2026-11-12)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0), `go1.26.6 darwin/arm64`,
in a `git archive 8338eebb` export at `plan-scratch/base` (a throwaway git repo; `diff -r cmd` and `diff -r host`
against the worktree: identical, `BASE_EQUALS_WORKTREE`; the worktree itself was only read). The pre-commit legs
are `tools/gates.sh <tree> <logs> base` = §6's script with the sprint's new-test lists emptied. Load average
13–23 on 16 cores at the start (`uptime` 22:05Z), with this planner's own mutant runs alongside later.
Logs: `plan-scratch/logs/base_gates/` (`legs.txt`, one log per leg).

| Leg (`tools/gates.sh <base export> … base`) | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `8338eebb`) |
| `verify_ail` (`./scripts/verify_ail.sh`) | **0** | 59 s | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `gofmt`, `build`, `vet`, `fence` | **0** | 0 / 5 / 3 / 32 s | empty / clean / clean / no FAIL line |
| `pkgs` (`host/approval/...`, `host/authority`, `host/store`, `host/sshsign`, `host/surfacepage`, `host/daemon`, `cmd/ailang-worldd`) | **0** | 405 s | every package `ok` (`host/daemon 404.1s`, `cmd/ailang-worldd 193.7s`) |
| `broker_named` | **0** | 33 s | every named test PASS |
| `broker_whole_reported` (reported, not gating) | **1** | 602 s | only `--- FAIL: TestExecProtectedCheckTimeout (42.28s)` (known list); **re-run alone: `ok host/broker 39.6s`** |
| `race` (three `-race` runs) | **0** | 244 s | `ok` ×6 packages, no race report |
| `inv` (every named inventory, `-v`, base mode: no sprint test) | **0** | 74 s | **49** named `--- PASS`, 0 SKIP, 0 FAIL, none missing |
| whole script | — | 1,459 s | `GATES_DONE … FAILS: broker_whole_reported` (the reported leg only) |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` (after commit) | — | — | **not measured by this planner** (it reads committed blobs; the planner makes no commit in the worktree); the controller's after-commit leg (§6) with the rule below |

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red untouched
package alone once and does not chase it. **Never start a gate under `nohup`.** `VERIFY_GO_LEGS=plain
./scripts/verify_go.sh` reads committed blobs and runs `go test ./...`; it is the after-commit leg (§6).

### 1.2 Codebase claims checked (at `8338eebb` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD` (worktree) → `8338eebb173290168f07e6b426d461c0122f67a8`; `git status --porcelain \| wc -l` → `0` |
| C2 | Controller-verified anchors hold | `grep -n "const credentialsSchemaVersion = " host/store/credentials.go` → `43:const credentialsSchemaVersion = 7`; `grep -n "is not a mint request id" cmd/ailang-worldd/approve.go` → `199:` (the design said `:197-201`); `grep -n "The grammar admits typed only in sprint 1" host/daemon/surface_decide.go` → `179:`; `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0` |
| C3 | **Drift since the design's base is exactly sprint 1** | `git diff --stat b0fbcd3e 8338eebb -- cmd host` → `17 files changed, 1982 insertions(+), 3 deletions(-)`, every file sprint 1's (§5 of its plan); the design's V1 hashes of the files sprint 2 reads re-measure identical: `shasum -a 256 host/authority/judge.go cmd/ailang-worldd/approve.go host/approval/sshsig.go host/approval/request.go cmd/ailang-worldd/episode_finish_operator.go host/daemon/surface_approve.go` → `d3a6620f0aa00022`, `6a38843161c7d75f`, `a40850200782b3e4`, `cf9897ec3f10881e`, `0adb363562679cd5`, `6f3cafda604f4c23` — equal to the design header; so the design's line numbers in those files hold |
| C4 | **The production SSH signer signs only under the mint namespace** (P1) | `grep -n '"-n", approval.Namespace' host/sshsign/sshsign.go` → `47: cmd := exec.CommandContext(ctx, bin, "-Y", "sign", "-n", approval.Namespace, "-f", key)`; `grep -n "sign:   signWithSSHKeygen" cmd/ailang-worldd/session.go` → `57:`; work: with the decision mode calling `env.sign`, the real-`ssh-keygen` sibling CD5 (`TestApproveDecideSignWithSSHKeygen`) would post a signature in `world-approval@sunholo`, which the daemon refuses `ErrNamespace`; the fake-signer test CD5 cannot see it (it never runs ssh-keygen) |
| C5 | **`approve` refuses the global `--addr`** (P3) | `grep -n "the daemon.s address is part of the mint request" cmd/ailang-worldd/main.go` → `465:`; so the decision mode takes its own `--daemon <host:port>` (S2-D2) |
| C6 | **The credentials store refuses any table set but its version's exactly** (P2) | `grep -n "slices.Equal(tables, credentialsTables(version))" host/store/credentials.go` → `325:`; `grep -rl '\`DROP TABLE surface_approvals\`' --include="*_test.go" host cmd \| wc -l` → `10` downgrade fixtures (each drops later tables and sets an older `user_version`); work, with schema 8 and no fixture edit, `go test ./host/store/` → `TestCredentialsUpgradeToV6`, `TestCredentialsSchemaMatchesCanonicalDDL`, `TestCredentialsUpgradeToV4`, `TestCredentialsSchema7SurfaceApprovals` FAIL (`unexpected table "key_decisions"`, `user_version after the upgrade = 8, want 7`); `git log -L` of `approve_passkey_test.go:99` → `38a6d903 feat(row 201): A3a sprint 1 M1a — credentials schema 7 …` edited the same fixtures for schema 7 (precedent) |
| C7 | The writer-site inventory | `surface_approve_post_guard_test.go:168` `want := map[string]string{"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials"}` keyed by function name; work: the new site is `keyArm`, so the map gains one entry (design §7 "1 → 2") |
| C8 | Sprint 1's grammar admits typed only, and RD2 already pins two-format and no-format bodies | `decide_model.go:63 Typed bool \`json:"typed"\``, `:84 case !body.Typed:`; `surface_decide_post_test.go:62 {"R6 a key format", … field(\`"typed":true,"sshsig":"-----BEGIN SSH SIGNATURE-----"\`) … 400}` and `:64 "R6 no format"`; work: with "exactly one of typed and sshsig" both arms still answer 400, so RD2 stays unedited |
| C9 | The decide caller census keys on the calling function (S2-D5) | `decide_once_test.go:161 "host/daemon/surface_decide.go:handleSurfaceDecide"`; work: the key arm reaches `broker.DecideApproval` through the same call in `handleSurfaceDecide` (`decidedBy` is a variable), so `TestDecideApprovalCallerCensus` stays unedited and PASS |
| C10 | The daemon's health names both stores | `daemon.go:654 CredentialsDB string \`json:"credentials_db"\``, `:629 DBPath string \`json:"db_path"\``; `episode_finish.go` (cmd) `:106`, `:113` compare both with `store.SameDatabase` (F6) |
| C11 | The judge's step names (JD1–JD3 agree) | `judge.go:22 JudgeStepParse = "parse"`, `:23 JudgeStepStore = "store"`, `:24 JudgeStepUnknownKey = "unknown-key"`, `:26 JudgeStepLaw = "law"`; `JudgeDecisionSSH` returns the same `*JudgeError` steps (JD tests) |
| C12 | `useApprover` is the consume's in-transaction approver binding (K3) | `approvers.go:245 func useApprover(ctx context.Context, tx *sql.Tx, a ApproverUse, at int64) error`: `UPDATE approvers SET last_used_at = ? WHERE name = ? AND format = ? AND public_key = ? AND revoked_at IS NULL AND enrolled_at <= ?`, else `classifyApproverUse` (C9–C12 sentinels) |
| C13 | A world-write fault can be injected without a production seam (RD13 a) | `schema.sql:16 semantic_id TEXT NOT NULL` on `objects`; `object_commits_test.go:92`, `journal_test.go:246` already use `CREATE TRIGGER … RAISE(ABORT, …)`; work: a trigger on `objects` for `world/approval-decision/v1` makes `broker.DecideApproval` fail after `RecordKeyDecision` committed (`TestDecidePOSTRetryAfterAFailedWorldWrite`) |
| C14 | **The J07 gap** | `decision_law_test.go:203 {"worktree#effect=Registry.Publish", "publish"}` — the only effect-arm row ends with the substring, so `strings.Contains` → `strings.HasSuffix` (`decision_law.go:118`) survived (iteration-389 judge); work: the new row `worktree#effect=Registry.Publish&x=1` kills it (`M4a`, ledger `J07-HASSUFFIX`) |
| C15 | **The J31/J32 gap** | `surface_decide.go:29 var decisionRefPattern = regexp.MustCompile(\`^sha256:[0-9a-f]{64}$\`)`; no daemon test posts a 65-hex or upper-hex ref (iteration-389 controller reproduced `{64}$` → `{64}` surviving); work: `TestDecideRefGrammarIsAnchored` posts and GETs 65, 63 and upper-case hex (ledger `J31-ANCHOR`, `J32-CASE`) |
| C16 | **The J33 gap** | `surface_decide.go:112 case "D1", "D2", "D3", "D4a", "D4b":`; no daemon test drives GET to D4a; work: `TestDecisionSubjectRefusesAnUndecodablePacket` (an undecodable stored packet; also sprint 1's unit-only D4a-undecodable row) (ledger `J33-GET-D4A`) |
| C17 | `host/approval` stays standard-library only | work: `decision_signed.go` imports `bytes` only; `TestApprovalImportsOnlyTheStandardLibrary` PASS |
| C18 | The context-root inventory | `context_roots_test.go` keys `file\|func\|Background`; work: `approve_decide.go\|runApproveDecide`, `\|postDecision`, `\|decideBySignWith` (+3, each with its bound in the comment); `TestProductionContextRoots` PASS |
| C19 | The top-level usage and help inventory | `help_test.go:44` lists the approve lines; work: + `"approve <request-ref> [--deny] [--daemon <host:port>]"` (one edited line), `main.go` usage + one line; `TestEveryVerbHasHelp` PASS |
| C20 | `ssh-keygen` on the rig | `command -v ssh-keygen` → `/usr/bin/ssh-keygen`; CD5's sibling ran (no SKIP) |
| C21 | **Nothing outside §5 moved** | stage repo: `git diff --stat <base> M5c -- world packages scripts tools docs design_docs \| wc -l` → `0`; production: `git diff --stat <base> M5c -- '*.go' '*.sql' ':!*_test.go'` → `13 files changed, 710 insertions(+), 28 deletions(-)`; added lines matching `t.Parallel\|Access-Control\|OpenCredentials(\|context.Background()\|/dev/tty\|.ssh` → the three CLI roots (pinned, C18), `keyArm`'s writer (C7), two test writers (`TestDecidePOSTBusyWriter`'s held writer, KD1's v7 open), and two comments naming `/dev/tty`; no `t.Parallel`, no CORS header |
| C22 | **Every mutant, on the final tree** | 125 rows, the last ledger line of each: 109 KILLED (named and suite), 14 `E`, 1 `I`, 1 `L`, 0 not as expected, 0 `NOCOMPILE` (§4, `logs/mutant_ledger.tsv`, generated by `tools/mtable.py`) |
| C23 | **Every sub-milestone is green when committed** | `tools/accept.sh` on a clone of the stage repo checked out at each tag in turn (`logs/accept_<M>.log`): `ACCEPT_OK` for **every one of the ten** (`M4a` … `M5c`): each stage gofmt-clean, `go build ./...` and `go vet` of the touched packages clean, the `_test.go` compile fence without a FAIL line, its §3 `-v` tests `--- PASS` with no SKIP/FAIL (CD5's real-`ssh-keygen` sibling ran: `--- PASS: TestApproveDecideSignWithSSHKeygen`), and its package runs `ok` (M5b's whole `host/daemon`, M5c's whole `cmd/ailang-worldd`) |
| C24 | **Every sub-milestone is committed and measured** | the round-trip tree commits each patch in turn: `NUMSTAT M4a 55`, `M4b 450`, `M4c 559`, `M4d 777`, `M4e 1142`, `M4f 1258`, `M4g 1380`, `M5a 1658`, `M5b 1804`, `M5c 1983` (`logs/roundtrip_final.log`) |
| C25 | **The patches reproduce the stage tree, which is the mutant tree** | `tools/roundtrip.sh <this plan>`: a fresh `git archive 8338eebb` export + `patch -p1` of Appendix P's ten patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r` of `cmd`, `host`, `docs` against `plan-scratch/work2` (`M5c`): **identical** (`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE`); `work2` is built from `work` (tag `FINAL`, the tree §4's mutants were cut from) by `tools/stages.sh`, which ends `WORK2_EQUALS_FINAL` (`diff -r cmd host`) |
| C26 | **The gates, on the final tree** (the round-trip tree) | `tools/gates.sh <round-trip tree> logs/final_gates final`: `verify_ail` → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` (no `.ail` changed); `gofmt`, `build`, `vet`, `fence` rc 0; `pkgs` rc 0 (`host/daemon 370.4s`, `cmd/ailang-worldd 152.0s`, every package `ok`); `broker_named` rc 0; `broker_whole_reported` rc 0 (`ok host/broker 221.1s`: green on this run, the base's known-list red did not recur); `race` rc 0 (`ok` ×6: host/approval, host/authority, host/store, host/sshsign, host/daemon (the decide, approve and finish tests), cmd/ailang-worldd (the approve tests); no race report); `inv` rc 0: **73** named `--- PASS` (49 carried, as at the base, + 24 new: `TestDecideBySignatureApprover` … `TestApproveDecideSignWithSSHKeygen`), 0 SKIP, 0 FAIL, none missing — incl. `TestFrozenV1RouteTableMatchesMux` (13), `TestDaemonCredentialWriterSites` (two sites), `TestProductionContextRoots` (+3), `TestCredentialsMethodsRequireADeadline` (+3), `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `TestEveryVerbHasHelp`, the quickstart tests; `GATES_DONE 858s FAILS: none` (`logs/final_gates/legs.txt`) |

### 1.3 Planner decisions where the design leaves a choice

- **S2-D1 — The fallback cut fires, from the measured size: the webauthn arm moves to sprint 3, and so do
  `--write-message`/`--signature` on a request ref.** Measured: the built set *with* CD6–CD7 was **2,019** (over
  the cap); with the webauthn arm it would have been the design's ≈ +500 more. Moved (Appendix C): DM7, RD10's
  passkey half, CD8–CD9, KD2 (the SQL counter is passkey-only: `usePasskey`), the passkey halves of K0, KD7, KD8,
  RD13 ((b), (c)) and RD14, `DecisionSubject.Passkeys` (sprint-1 P7), CD6 (`--write-message`) and CD7
  (`--signature` + `--expiry`). In this sprint a request ref with `--signature`, `--write-message`,
  `--webauthn-*` or `--await` is refused before any daemon call (CD3, CD5). Built, never cut: CD1–CD4, R4e,
  RD13 (a), KD5/KD6, schema 8, the judge, J07/J31/J32/J33.
- **S2-D2 — How the CLI finds the daemon for a request ref.** `--daemon <host:port>` (a local flag, because
  `approve` refuses the global `--addr`, P3), default `defaultServiceAddr()` (the installed service's address,
  else `127.0.0.1:7644`). Before any decision read, `GET /v1/health` must name `--db` and the credentials store
  beside it (or `--credentials-db`) — `store.SameDatabase`, as `episode finish` F6 — because the typed proof is
  written beside the served credentials store and the daemon consumes it there. The POST carries `Origin:
  http://<host:port>` (in `O` by construction, design §2.8) and no `Sec-Fetch-Site`.
- **S2-D3 — Schema 8 is the design's DDL verbatim and forward-compatible.** No passkey-only column: the
  passkey counter is `webauthn_credentials.sign_count` (schema 6) and `usePasskey` already moves it; a passkey
  key decision is a `format = 'webauthn'` row with the assertion as `artifact`. Sprint 3 adds Go only. `KeyDecision`
  (the read) ships for the tests and sprint 3's listing; `RecordKeyDecision(ctx, row, use ApproverUse, now)` takes
  the store's own `ApproverUse` (the store cannot import `host/approval`'s `SignedDecision`); sprint 3 extends the
  same call with the passkey use (`ApproverUse.Passkey`, already a field).
- **S2-D4 — K0 for SSHSIG.** An unlanded row whose artifact, packet ref, verb, format, approver and expiry equal
  the posted ones is the exact match: `useApprover` still runs (revocation is re-checked on both paths, design §2.7
  [r2]); the row is kept as first recorded (no rewrite, its `decided_at` unchanged). Anything else replaces the
  unlanded row (K5). A landed row refuses (`ErrKeyDecisionLanded`) — unreachable through the POST, where D6
  refuses first. For SSHSIG the counter does not exist, so the exact-match path differs from replacement only in
  keeping the first record (KD7's SSH half asserts it).
- **S2-D5 — N4 / K0's slot, as Appendix C fixed it; sprint-1 code edited only at the places it names.** In
  `handleSurfaceDecide`, at the arm slot: `if body.SSHSIG != nil { keyArm }` else the typed proof; then the one
  `broker.DecideApproval` call with the derived `decidedBy`; then, for the key arm, `LandKeyDecision`. `keyArm`
  (new file): `JudgeDecisionSSH` on a read-only handle (step `store` → `approvalStoreError`: 503 `Restart`; other
  steps → 422 `NotApproved`) → `Binds` → `store.OpenCredentials` non-waiting (busy → 503 `Busy`, `Retry-After: 1`)
  → `RecordKeyDecision` (an approver race → 422; an older store → 503 `Restart`; else 500). The writer is closed by
  the handler's `defer` after `LandKeyDecision`. Sprint-1 lines edited: the grammar (`ParseDecideBody`'s format
  cases, `DecideBody` gains `SSHSIG` and `typed` becomes `omitempty`), the arm slot, the decide call's `decidedBy`,
  the land call, the response's `decidedBy`, RD9b/RD9c's SSHSIG arms (a shared helper `requireKeyRefused`).
- **S2-D6 — The decision signer is a second seam.** `sessionEnv.signDecision` (production
  `signDecisionWithSSHKeygen` → `sshsign.SignIn(ctx, "ssh-keygen", approval.DecisionNamespace, …)`); `sshsign.Sign`
  becomes `SignIn` under the mint namespace, so its one driven exec site and every mint caller are unchanged
  (`TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` counts sites per file: still one). The signer's
  root is bounded by the signed expiry (`window.max − now`, ≤ 600 s).
- **S2-D7 — The CLI's fake daemon and its real sibling.** CD1–CD5 run against an `httptest` loopback server whose
  subject is built with the daemon's own encoders and types (`coordinator.EncodeDecisionPacket`,
  `surfacepage.DecisionSubject`, `DecideBody`, `DecisionMade`); the daemon side of the same wire is host/daemon's
  RD2–RD13 through `d.Handler()`. The "real tool" of CD5 is `ssh-keygen`: `TestApproveDecideSignWithSSHKeygen`
  runs the production decision signer on a temp key and checks the posted signature with the daemon's own
  `VerifyDecisionSSHSIG`, and that it is **not** a mint signature (`ErrNamespace`). A live run instead of a captured
  file: it is what found P1; with `WORLD_REQUIRE_SSH_KEYGEN=1` a missing `ssh-keygen` fails (CI).
- **S2-D8 — Design Q2 (`approve --list` decisions):** waits for sprint 3's `requests.json` `decisions` view;
  `approve --list` is unchanged (it lists mint requests). An operator reads a request ref from the packet page
  or the agent's output.
- **S2-D9 — D-WORLD-115.** Under B (a stored daemon challenge), the 32-byte P slot of `DecisionMessage` carries the
  challenge: `DecideBySignature`'s message input, `JudgeDecisionSSH`'s `DecisionJudgeInput`, the CLI's
  `decideBySignWith` message and the GET subject (a challenge field) change, plus an issuing step; schema 8 would
  need an issuance table or columns (the design prices ≈ +550). Built under **A**; nothing else in this sprint
  depends on it.
- **S2-D10 — Carried items.** R4e is a daemon test of L1 for the key format; its must-die half ("with L1 compiled
  out the key decision lands but `validatePublishApproval` refuses it") is the mutant `LAW-D5-KEY` (L1 deleted):
  R4e is KILLED and sprint 1's R4b (`TestPublishApprovalRefusesAKeyDecision`) still PASSES on that mutant (ledger,
  §4). H-LOCK: `TestDecideConcurrentPOSTsDecideOnce` (eight concurrent POSTs of one key decision: exactly one 200,
  seven 409 D6). `TestDaemonCredentialWriterSites` 1 → 2.
- **S2-D11 — Sub-milestones** (§1.5). M4b (395) is over the 250–350 band: it is one coherent unit (the decision
  mode's production code with its rig and its flag/no-terminal refusals); its tests CD1/CD4 are M4c.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §2.10 "`--sign-with` … the A1b/A2 flag names, the decision message instead of the mint
  message".** The production signer cannot sign a decision: `sshsign.Sign` hard-codes `-n world-approval@sunholo`
  (C4). Routed by `sshsign.SignIn` (namespace parameter; `Sign` delegates) and `sessionEnv.signDecision`
  (S2-D6). Found by the real-tool sibling, not by the fake signer.
- **P2 — Design §6 "Must stay green, unedited: every `TestApprove*`, `TestSurfaceApprove*`" with schema 8.**
  Ten downgrade fixtures (C6) must also drop `key_decisions`, or a "schema 5/6 store" is refused as not a
  credentials store: `host/daemon/surface_approve_{read,post}_test.go` (TestSurfaceApprove*), `host/authority/
  judge_test.go`, `cmd/ailang-worldd/{approve_await_refusals, approve_passkey_refusals, approve_passkey,
  approve_sig, approver_label, approver_passkey, approver_rig}_test.go` (TestApprove*/TestApprover*) — one token
  per site, `\`DROP TABLE key_decisions\`, ` before `\`DROP TABLE surface_approvals\``. Four store tests pin the
  table set or the version (`credentials_test.go`, `approvers_test.go`, `mint_requests_test.go`,
  `surface_approvals_test.go`). Precedent: A3a sprint 1 M1a (`38a6d903`) did exactly this for schema 7.
- **P3 — Design §2.10 assumes the CLI can reach "the daemon"**; `approve` refuses the global `--addr` (C5) and a
  request ref carries no address (sprint-1 Appendix C). Routed by `--daemon` (S2-D2).
- **P4 — Sprint-1 `DecideBody.Typed` is `json:"typed"`** (C8): a CLI marshalling an `sshsig` body would send
  `"typed":false`. Routed by `omitempty` (decoding is unchanged; RD2's typed-false body is still "no format").
- **P5 — Design §6 KD2 (`TestRecordKeyDecisionCounter`) is in "Sprint 2" but tests only the passkey counter**;
  moved with the webauthn arm (S2-D1).
- **P6 — Design §6 RD11 "a body with a forged `decidedBy` → 400"** is already sprint 1's grammar (RD2 "R6
  decidedBy"); RD11 here re-asserts it on an `sshsig` body and pins the landed `decidedBy` (`TestDecidePOSTSSHSIG`).
- **P7 — Design §6 KD9 "the approver revoked after the read-only judge" is store-level only**, as the design
  says; its HTTP mapping (`isApproverRace` → 422) is unreachable end to end because the read-only judge refuses a
  revoked approver first (§3.1; named survivor `D-RACE`).
- **P8 — Design §7 "`TestSurfaceAppRoutes` +2 (sprint 3)" / no route this sprint**: unchanged (no route).

### 1.5 The size

Lines are insertions + deletions under `cmd/ host/`: **`git diff --numstat <base> HEAD -- cmd host`**
(`docs/` and `design_docs/` are outside the filter). **Every row is measured** from the scratch stages, each
committed in a throwaway repo (`plan-scratch/work2`, tags `M4a` … `M5c`), and re-measured on the round-trip tree
(`roundtrip.sh` `NUMSTAT` lines, §1.2 C24). "Own" is the stage's own diff against the previous stage.

| Sub-milestone | Production | Test | Own | Cumulative vs `8338eebb` (after its commit) | Run |
|---|---|---|---|---|---|
| M4a — J07, J31/J32, J33 (test only) | 0 | 55 | 55 | **55** | 1 |
| M4b — the typed decision mode (`approve_decide.go`, dispatch, `--deny`, `--daemon`, help, roots); CD2, CD3 | 211 | 184 | 395 | **450** | 1 |
| M4c — CD1, CD4 | 0 | 109 | 109 | **559** | 1 |
| M4d — `DecideBySignature`, `SignedDecision` (`decision_signed.go`); DM3–DM6, DM8 (SSH) | 91 | 127 | 218 | **777** | 2 |
| M4e — schema 8 (`key_decisions.go`, the DDL, the upgrade); KD1, KD5, KD6; P2's fixtures; the deadline inventory | 196 | 169 | 365 | **1,142** | 2 |
| M4f — KD3, KD4/KD7 (SSH), K2, KD9's store half | 0 | 116 | 116 | **1,258** | 2 |
| M4g — `JudgeDecisionSSH`; JD1–JD3 | 45 | 77 | 122 | **1,380** | 2 |
| M5a — the POST's `sshsig` arm (`surface_decide_key.go`, the slot, the grammar); RD10, RD11, the judge's refusals; writer sites 1 → 2 | 128 | 150 | 278 | **1,658** | 3 |
| M5b — RD12, RD13 (a) (key and typed), R4e, H-LOCK, RD9b/RD9c's SSHSIG arms | 0 | 146 | 146 | **1,804** | 3 |
| M5c — `--sign-with` on a request ref, `sshsign.SignIn`, `signDecision`; CD5 + its real-tool sibling | 77 | 112 | 189 | **1,983** | 3 |
| **Sprint 2** (cumulative, `git diff --numstat 8338eebb M5c -- cmd host`) | **738** | **1,245** | | **1,983 ≤ 2,000** (margin 17) | |

Per file, cumulative (+/−): `approve.go` +29, `approve_decide.go` +229, `main.go` +1, `session.go` +5,
`approve_decide_rig_test.go` +179, `approve_decide_test.go` +109, `approve_decide_key_test.go` +111, `help_test.go`
+1 −1, seven cmd fixtures +1 −1 each, `decision_signed.go` +91, `decision_signed_test.go` +127,
`judge_decision.go` +45, `judge_decision_test.go` +77, `judge_test.go` +2 −2, `decision_law_test.go` +3 −2,
`surface_decide.go` +18 −4, `surface_decide_key.go` +70, `surface_decide_key_test.go` +148,
`surface_decide_key_retry_test.go` +129, `surface_decide_jgaps_test.go` +50, `surface_decide_discovery_test.go`
+17, `surface_approve_post_guard_test.go` +1 −1, `surface_approve_{post,read}_test.go` +2 −2 each,
`sshsign.go` +9 −5, `credentials.go` +9 −3, `credentials_schema.sql` +15, `key_decisions.go` +169,
`key_decisions_test.go` +122, `key_decisions_record_test.go` +116, `credentials_test.go` +6 −3,
`context_roots_test.go` +4, `approvers_test.go`/`mint_requests_test.go` +1 −1 each, `surface_approvals_test.go`
+3 −3, `decide_model.go` +20 −16. Against the design's sprint 2 ≈ 1,800 + the CLI ≈ 400: the built part with the
passkey arm and CD6–CD7 out came in at 1,983.

---
## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four
backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`
and an `@@ -0,0 …` hunk. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M4a   # <- the sub-milestone, exactly as in the Appendix P heading: M4a M4b M4c M4d M4e M4f M4g M5a M5b M5c
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter390-a4-s2
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a4-sprint2-plan.md"
T="$HOME/.ailang/state/world-iter390/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with `export
  AILANG_BIN=$HOME/.pinned-ailang/ailang`, each command's output `tee`'d to a log under
  `$HOME/.ailang/state/world-iter390/exec-logs/` (never only through `tail`). **Never `nohup`, never background a
  test leg.** The planner's `tools/accept.sh <tree> <M>` is exactly §3's lines; it printed `ACCEPT_OK` for every
  stage on the stage repo (§1.2 C23).
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/approval/... ./host/authority/ ./host/store/ ./host/sshsign/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/ ./host/broker/
go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)' && echo FENCE_RED   # must print nothing
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- FAIL` and no `--- SKIP`**; each
package run must end `ok` with no `FAIL` line.

### M4a — the sprint-1 judge's test gaps · 55 test · Run 1

**Patch** `M4a`: `host/broker/decision_law_test.go` — one R4a row `worktree#effect=Registry.Publish&x=1` → publish
(J07; gofmt realigns the two comment lines above it). New `host/daemon/surface_decide_jgaps_test.go`:
`TestDecideRefGrammarIsAnchored` (J31/J32: POST and GET with 65 hex, 63 hex, upper-case hex → 400 D1, 0 decisions,
proof unconsumed; the 64-hex GET as control), `TestDecisionSubjectRefusesAnUndecodablePacket` (J33: one undecodable
stored packet → GET 409 D4a "cannot be counted completely", POST 409, 0 decisions).

**Acceptance**
```
[base]
go test ./host/broker/ -run '^TestReservedApprovalScope$' -count=1 -v
go test ./host/daemon/ -run '^(TestDecideRefGrammarIsAnchored|TestDecisionSubjectRefusesAnUndecodablePacket)$' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M4a: 55**.

### M4b — the typed decision mode · 211 prod + 184 test · Run 1

**Patch** `M4b`: new `cmd/ailang-worldd/approve_decide.go` — `decisionRefPattern` (`^sha256:[0-9a-f]{64}$`),
`approveDecideBudget = 10 s`, `decideOptions`, `decisionSubject`, `checkDecisionSubject` (strict decode; both byte
fields hash to their refs; the request and the packet decode strictly; the packet carries the ref),
`sha256Ref`, `refDigest`, `renderDecision` (every field of both objects, entry, proposal page, the verb in
capitals, the deadline), `refusalReason`, `runApproveDecide` (credentials path; `--daemon` or the service's
address; `/v1/health` must name both stores; GET; the checks; the display refusal; `openApproverTerminal` +
`confirmOnTerminal` (A1b's 120 s y/N); `writeOperatorProof(credentials, "decision-"+hex(Q), "decide-"+verb)`;
POST), `postDecision` (`Origin` = the daemon's own base URL, the proof header, the answer printed).
`approve.go`: two flags `--deny`, `--daemon`; CD3's rule (`--deny`/`--daemon` need a request ref); the id-shape
dispatch beside the `mr-` check; in this stage any key flag or `--await` with a request ref is refused (M5c narrows
it); the help lines. `main.go`: one usage line. `help_test.go`: the line in `TestEveryVerbHasHelp`.
`context_roots_test.go`: `runApproveDecide`, `postDecision` (+2 roots, 10 s each). New
`approve_decide_rig_test.go`: `decideTerm`, `fakeDecider` (an `httptest` loopback daemon), `decisionObjects`,
`run`, `proofFile`; CD2 `TestApproveDecideNoTerminal`, CD3 `TestApproveDecideFlagRules`.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveDecideNoTerminal|TestApproveDecideFlagRules|TestEveryVerbHasHelp)$' -count=1 -v
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M4b: 450**.

### M4c — the typed decision end to end, the subject checks · 109 test · Run 1

**Patch** `M4c`: new `cmd/ailang-worldd/approve_decide_test.go`: CD1 `TestApproveDecideTyped` (approve y, deny y,
approve n; the POST body, `Origin`, the proof file's phase/name/hash; every field in the prompt; a refused POST
reported with the daemon's reason), CD4 `TestApproveDecideRefusesAMismatchedSubject` (eleven arms, each: no terminal
opened, nothing posted).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveDecideTyped|TestApproveDecideRefusesAMismatchedSubject)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^TestApprove' -count=1
```
**Expected cumulative numstat after the controller's commit of M4c: 559**. **After M4c an operator can decide a
packet at a terminal by a typed `y`** (sprint 1's typed POST now has its shipped client).

### M4d — the decision signature law · 91 prod + 127 test · Run 2

**Patch** `M4d`: new `host/approval/decision_signed.go` — `SignedDecision` (unexported fields), `DecisionSignatureInput`,
`VerifyDecisionSSHSIG` (the 91-byte message under `DecisionNamespace`), `DecideBySignature` (`checkApprover` P2–P5 at
`now`, then the signature), `Binds` (verb, packet, request, expiry), `DecidedBy`, `Use`. New
`decision_signed_test.go` (external package `approval_test`, approvaltest keys): DM8 `TestDecideBySignatureApprover`,
DM3/DM4 `TestMintSignatureIsNotADecision`, DM5 `TestDecisionBindsTheVerb`, DM6 `TestSignedDecisionBinds`.

**Acceptance**
```
[base]
go test ./host/approval/... -run '^(TestDecideBySignatureApprover|TestMintSignatureIsNotADecision|TestDecisionBindsTheVerb|TestSignedDecisionBinds|TestApprovalImportsOnlyTheStandardLibrary)$' -count=1 -v
go test ./host/approval/... -count=1
```
**Expected cumulative numstat after the controller's commit of M4d: 777**.

### M4e — credentials schema 8 · 196 prod + 169 test · Run 2

**Patch** `M4e`: new `host/store/key_decisions.go` — `keyDecisionsSQL` (the design's DDL verbatim),
`KeyDecisionRow`, `ErrBadKeyDecision`, `ErrKeyDecisionLanded`, `ErrKeyDecisionNotPending`, `RecordKeyDecision`
(K1 deadline + schema 8, K2 grammar, K0 landed/exact match, K3 `useApprover` on both paths, K5 upsert of the
unlanded row), `LandKeyDecision`, `KeyDecision`. `credentials.go`: version 8, `credentialsTables(8)`, the 7 → 8
upgrade. `credentials_schema.sql`: the table. Edited store tests (P2): `credentials_test.go` (table set, version 8,
`TestCredentialsMethodsRequireADeadline` +3 methods), `approvers_test.go`, `mint_requests_test.go`,
`surface_approvals_test.go`. P2's ten downgrade fixtures (`host/daemon` ×2, `host/authority` ×1, `cmd` ×7). New
`key_decisions_test.go`: KD1 `TestCredentialsSchema8KeyDecisions`, KD5/KD6 `TestKeyDecisionsRefChecks`.

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestCredentialsSchema8KeyDecisions|TestKeyDecisionsRefChecks|TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsSchema7SurfaceApprovals|TestCredentialsUpgradeToV4|TestCredentialsUpgradeToV6|TestCredentialsMethodsRequireADeadline)$' -count=1 -v
go test ./host/store/ ./host/authority/ -count=1
go test ./host/daemon/ -run '^(TestSurfaceApprove|TestApprov|TestRequest)' -count=1
go test ./cmd/ailang-worldd/ -run '^(TestApprove|TestApprover)' -count=1
```
**Expected cumulative numstat after the controller's commit of M4e: 1142**. RK7 recurs: schema 8 costs one
`serve` restart, named by the existing `OlderCredentialsStoreError` text.

### M4f — the evidence writer's transaction · 116 test · Run 2

**Patch** `M4f`: new `host/store/key_decisions_record_test.go`: KD3/KD9 `TestRecordKeyDecisionRevokedSince`, KD4/KD7
`TestRecordKeyDecisionReplaceUnlanded`, K2 `TestRecordKeyDecisionGrammar`.

**Acceptance**
```
[base]
go test ./host/store/ -run '^TestRecordKeyDecision' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M4f: 1258**.

### M4g — the decision judge · 45 prod + 77 test · Run 2

**Patch** `M4g`: new `host/authority/judge_decision.go` — `DecisionJudgeInput`, `JudgeDecisionSSH` (parse → store →
unknown-key → law, `judge.go`'s step names); new `judge_decision_test.go` (JD1–JD3: eight arms and the control,
over A3's `jgRig`).

**Acceptance**
```
[base]
go test ./host/authority/ -run '^(TestJudgeDecisionSSH|TestJudgeSSHSignatureIsTheCLIJudge|TestJudgeAssertionIsTheCLIJudge)$' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M4g: 1380**.

### M5a — the POST's `sshsig` arm · 128 prod + 150 test · Run 3

**Patch** `M5a`: new `host/daemon/surface_decide_key.go` — `isApproverRace`, `keyArm` (S2-D5). `surface_decide.go`:
the arm slot, the derived `decidedBy` in the decide call and the answer, `LandKeyDecision`. `decide_model.go`:
`DecideBody.SSHSIG`, `typed` `omitempty` (P4), the grammar's two format cases, comments. 
`surface_approve_post_guard_test.go`: the writer-site map gains `keyArm`. New `surface_decide_key_test.go`: the
rig's `enrolDecider`, `decideSSHBody`, `keyRow`, `postKey`; RD10/RD11 `TestDecidePOSTSSHSIG`;
`TestDecidePOSTSSHSIGRefusals` (seven 422/409 arms, then a v7 store's 503 `Restart` from the writer and a v4
store's from the judge).

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestDecidePOSTSSHSIG|TestDecidePOSTSSHSIGRefusals|TestDecidePOSTGuards|TestDaemonCredentialWriterSites|TestDecisionSubjectJSON|TestDecideTypedLands)$' -count=1 -v
go test ./host/surfacepage/ -count=1
```
**Expected cumulative numstat after the controller's commit of M5a: 1658**.

### M5b — busy, retry, R4 by key, the lock, the discovery by key · 146 test · Run 3

**Patch** `M5b`: new `surface_decide_key_retry_test.go`: RD12 `TestDecidePOSTBusyWriter`, RD13 (a)
`TestDecidePOSTRetryAfterAFailedWorldWrite` (a SQLite trigger on the world store's `objects`; a typed POST under the
same fault is 500 too), R4e
`TestDecideRefusesReservedScopesByKey`, H-LOCK `TestDecideConcurrentPOSTsDecideOnce`. `surface_decide_discovery_test.go`:
`requireKeyRefused` and its calls in RD9b and RD9c (a valid SSHSIG → 409, 0 decisions, 0 `key_decisions` rows).

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestDecidePOSTBusyWriter|TestDecidePOSTRetryAfterAFailedWorldWrite|TestDecideRefusesReservedScopesByKey|TestDecideConcurrentPOSTsDecideOnce|TestDecideDiscovery|TestDecideRefusesASharedRequest)' -count=1 -v
go test ./host/daemon/ -count=1 -timeout 30m
```
**Expected cumulative numstat after the controller's commit of M5b: 1804**.

### M5c — `--sign-with` on a request ref · 77 prod + 112 test · Run 3

**Patch** `M5c`: `host/sshsign/sshsign.go` — `SignIn` (namespace parameter), `Sign` delegates (P1).
`cmd/ailang-worldd/session.go`: `sessionEnv.signDecision`, set to `signDecisionWithSSHKeygen` in `realSessionEnv`.
`approve_decide.go`: `decideBySignWith`, `signDecisionWithSSHKeygen`, `decideOptions.signWith`. `approve.go`: the
key-flag rule narrowed to "`--sign-with` only", the help line. `context_roots_test.go`: `decideBySignWith` (+1 root,
bounded by the signed expiry). New `approve_decide_key_test.go`: CD5 `TestApproveDecideSignWith` (the 91-byte message,
the signer's deadline, the POST, a signer error, the production seam, `--signature`/`--write-message` refused) and
its real-tool sibling `TestApproveDecideSignWithSSHKeygen`.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveDecide|TestEveryVerbHasHelp|TestRealSessionEnvSignsWithSSHKeygen|TestSignWithSSHKeygenEndToEnd)' -count=1 -v
go test ./host/sshsign/ ./host/store/ -count=1
go test ./host/broker/ -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' -count=1 -v
go test ./cmd/ailang-worldd/ -count=1 -timeout 30m
```
**Expected cumulative numstat after the controller's commit of M5c: 1983**. **After M5c an operator can decide a
packet by an enrolled SSH key**: `ailang-worldd approve <request-ref> [--deny] --sign-with <key>`.

### 3.1 Reachability: every new refusal and branch, from the entry point (directive retro 4)

Entry points: the CLI `ailang-worldd approve <arg> …` → `runApprove` → (a request ref) `runApproveDecide`; the
daemon's `POST /surface/decide/{ref}` → `approvalRoute` (R2) → `handleSurfaceDecide` → R0, R3–R6 → `humanMu` → the
law D1–D10 → the arm slot. "E2E" = a test drives it from the entry point (the CLI against the fake daemon, or
`d.Handler()`); "unit" = only a package test reaches it.

| Branch | Input class → trace | Driven by |
|---|---|---|
| CLI dispatch | `sha256:` + 64 hex → decision mode; 65 hex → the `mr-` check refuses | E2E CD1, CD3 |
| `--deny`/`--daemon` rule | with an `mr-` id, with `--list` → refused before any read | E2E CD3 |
| key-flag rule | `--await`, `--webauthn-link`, `--signature`, `--write-message` with a ref → refused, 0 GETs | E2E CD3, CD5 |
| health checks | daemon serving another store / credentials store → refused, no terminal | E2E CD4 |
| health error | daemon down → refused | **not driven**: with the mutant the zero health names no store and the next check refuses anyway (named equivalent `C-HEALTH-ERR`) |
| GET status | the daemon's 409 → its reason printed | E2E CD4 |
| subject checks | unknown field; bad base64; request bytes of another request; packet bytes not hashing to the packet ref; another request's packet; undecodable request; undecodable packet | E2E CD4 (each arm) |
| display refusal | `refusal` set (D5–D7 at the daemon's now) → refused before the terminal | E2E CD4 |
| no terminal | `openTerminal` fails → refused naming `--sign-with` | E2E CD2 |
| y/N | `n` → nothing posted, no proof | E2E CD1 |
| proof write error | an unwritable proofs directory | **not driven** (named survivor `C-PROOF-ERR`): `writeOperatorProof`'s refusals are row 236's, tested there |
| POST refused / answered | 409 → the reason; 200 → printed | E2E CD1 |
| `--sign-with` | the decision signer seam; a signer error; the production seam; the real ssh-keygen | E2E CD5 (+ sibling) |
| grammar: no format / two formats | `typed` absent and no `sshsig`; both → 400 | E2E RD2 (sprint 1, unedited) |
| key arm: judge parse / unknown key / law (revoked, mint namespace, verb, expiry, garbage) | 422 `NotApproved`, 0 decisions, 0 rows | E2E `TestDecidePOSTSSHSIGRefusals` |
| key arm: another packet's P | D8 refuses first (409) — the law runs before the arm | E2E same test |
| key arm: judge store step | a v4 credentials store → 503 `Restart` | E2E same test |
| key arm: `Binds` | unreachable: the judge verifies exactly the subject `Binds` checks (named equivalent `D-BINDS`) | — |
| key arm: busy writer | a held writer → 503 `Busy`, `Retry-After: 1` | E2E RD12 |
| key arm: `RecordKeyDecision` older store | a v7 store → 503 `Restart` | E2E `TestDecidePOSTSSHSIGRefusals` |
| key arm: approver race | revoked between the judge and the writer | **unit only** (KD3/KD9 at the store): through the POST the read-only judge refuses a revoked approver first (named survivor `D-RACE`) |
| `DecideApproval` error | the world write fails (trigger) → 500, the row unlanded; typed → 500 | E2E RD13 (a) (sprint 1's `H-DECIDE-ERR` survivor is now killed, `D-DECIDE-ERR`) |
| exact-match retry | the identical body after the failure → 200, the row kept and landed | E2E RD13 (a); unit KD7 |
| land error | a store fault between the decision and the land | **not driven** (named survivor `D-LAND-ERR`): the chain is the authority; `approve --list --decided`'s "not linked" view is sprint 3's |
| R4 by key | publish- and finish-shaped requests → 403, 0 rows | E2E R4e; L2 alone: R4b on the L1-deleted mutant (`LAW-D5-L2`) |
| discovery by key | hidden second packet, overflow → 409, 0 rows | E2E RD9b, RD9c |
| `humanMu` | eight concurrent POSTs → one 200, seven 409 | E2E H-LOCK |
| store K0 landed | a landed row → `ErrKeyDecisionLanded` | **unit only** (KD4): through the POST D6 refuses a landed request first |

---
## 4. Mutation table

**125 rows on the planner's final tree (`work`, tag `FINAL` = M5c), the last ledger line of each: 109 KILLED by the named test and the package suite, 14 `E` rows (named equivalents and survivors, each with its reason), 1 `I` row (killed only by its inventory, the suite survives — recorded both ways), 1 `L` row (a must-live half that lives); 0 not as expected; 0 `NOCOMPILE` among the last lines.** 135 ledger lines in all (`plan-scratch/logs/mutant_ledger.tsv`; runs: run1 118, run2 16, run3 1). 5 ledger lines are `NOCOMPILE` (run1: an unused variable or import left by the edit) and are **not counted**: each row was rewritten to compile (`_ = …`, `n < 0`, another field) and re-run in run2; the NOCOMPILE lines stay in the ledger. Four run1 rows that survived exposed missing test arms (`C-PKT-HASH`, `C-REQ-DECODE`, `S-ISREF-HEX`, `D-DECIDE-ERR`) and one whose named test was not the killing one (`A-NAMESPACE`); the arms were added to CD4, the store grammar test, RD13 (a) and DM3 (the patches in Appendix P are the corrected ones) and the rows re-run in run2 on the final tree; run1's other rows ran on the tree before those arms existed (arms only add kills; every run1 row that survived or did not compile was re-run on the final tree). run2 also ran the J-gap rows and R4e's pair; run3 the carried `H-LOCK` (eight concurrent POSTs; killed in its one run, a probabilistic kill, §8 R6). **No kill rests on an inventory or static test** except the one `I` row, recorded as such: the daemon suite skips RD1 and the writer-site inventory, the store suite runs no deadline or context-root inventory, the broker suite skips the census.

**Derivation.** Anchored to Appendix P's production hunks, one row per new production branch, enumerated from the
diff (`git diff 8338eebb M5c -- cmd host ':!*_test.go'`): the CLI's id-shape dispatch and its anchor, the
`--deny`/`--daemon` rule (both operands), the key-flag rule (both ways), the verb, the two health checks and the
health error, the GET status, the strict decode, each subject check (base64, request hash, packet hash, request
decode, packet decode, packet carries the ref), the display refusal, the terminal and the y/N, the proof's name,
phase and write error, `typed`, the expiry, `Origin`, the proof header, the POST status and the answer's request,
two prompt fields (M4b); the `--sign-with` dispatch, the signer seam's nil check, the signed verb and expiry, the
signer's deadline, the signer error, the posted `sshsig`, both namespaces (`signDecisionWithSSHKeygen`,
`sshsign.SignIn`) (M5c); the law's approver checks, format and clock, the signature, its namespace and its verb,
each `Binds` field and the unverified value, `DecidedBy`/`Use` on the unverified value, the format in `DecidedBy`,
the key copy (M4d); the schema's table set and upgrade, K1 (deadline, version on each of the three methods), each
K2 operand, the lower-case ref grammar, the not-found read, K0 (landed, exact, and each of its five field
conjuncts), K3 (both paths, deleted), K5 (the write and its `WHERE`), `decided_at`, the land's ref check, `IS
NULL` and row count, the read's not-found (M4e); each judge step, the fingerprint, and each approver field passed
to the law (M4g); the arm switch, the arm's refusal return, the derived `decidedBy` (call and answer), the land
call and its error, the writer close, the decide error, the verb byte, the read-only open, `Binds`, the judge's
store-step mapping, the 422, the writer open, the evidence row's packet ref, artifact, expiry and time, the
record error and the approver-race mapping, the grammar's two format cases (M5a). Plus the J gaps on sprint-1 code
(`J07-HASSUFFIX`, `J31-ANCHOR`, `J32-CASE`, `J33-GET-D4A`) and R4e's must-die pair on sprint-1's L1 (`LAW-D5-KEY`
killed by R4e; `LAW-D5-L2`: the same mutant leaves sprint 1's L2 test `TestPublishApprovalRefusesAKeyDecision`
passing — L2 alone holds).
**Procedure**: `plan-scratch/tools/mutate.py <run> <tree> <workers> [ids]` (the table is `tools/mutants.py`;
every anchor is asserted to occur exactly once in the tree it edits): for each row, in its own APFS clone of the
tree (`cp -cR`), apply **one** edit, run the named test alone, then the package suite, delete the clone, and append
one ledger line (UTC time, run, tree, id, stage, file, from → to, named test, named, suite, the failing tests,
seconds, expected). `GOFLAGS=-vet=off`; a mutant that does not compile is `NOCOMPILE`, never counted. Suites
(inventories and static tests excluded): `host/approval` all `^Test` minus `TestApprovalImportsOnlyTheStandardLibrary`;
`host/store` `^(TestCredentialsSchema|TestKeyDecisions|TestRecordKeyDecision|TestCredentialsUpgrade)` (no deadline or
context-root inventory); `host/authority` all; `host/daemon` `^(TestDecision|TestDecide|TestSurfaceApprove|TestApprov)`
minus `TestDecisionRoutesRegistered`, `TestApprovalRoutesRegistered`, `TestDaemonCredentialWriterSites`;
`cmd/ailang-worldd` `^(TestApproveDecide|TestApprove|TestApprover|TestSignWith|TestRealSessionEnv)` minus
`TestEveryVerbHasHelp`; `host/broker` `^(TestDecisionRefusal|TestReservedApprovalScope|TestPublishApproval|TestDecideApproval)`
minus the census. Expected: `K` both fail; `E` a named equivalent or survivor (reason in the row); `I` killed only by
its named **inventory** test, the behavioural suite survives (both recorded, directive retro 3); `L` the named test
must **pass** (a must-live half). **Never commit or snapshot a mutant.** In the table a newline prints as ⏎ and a
pipe as \|.

| Id | Stage | file: from → to | Named test | Named | Suite | Failing tests (suite, else named) | Expected |
|---|---|---|---|---|---|---|---|
| C-DISPATCH | M4b | approve.go: if decisionRefPattern.MatchString(id) { → if false { | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideFlagRules,TestApproveDecideNoTerminal,TestApproveDecideRefusesAMismatchedSubject,TestApproveDecideSignWith,TestApproveDecideSignWithSSHKeygen,T | K |
| C-REF-ANCHOR | M4b | approve_decide.go: `^sha256:[0-9a-f]{64}$` → `^sha256:[0-9a-f]{64}` | `TestApproveDecideFlagRules` | KILLED | KILLED | TestApproveDecideFlagRules | K |
| C-DENY-RULE | M4b | approve.go: if (*deny \|\| *daemonAddr != "") && (*list → if false && (*list | `TestApproveDecideFlagRules` | KILLED | KILLED | TestApproveDecideFlagRules | K |
| C-DENY-RULE-DAEMON | M4b | approve.go: if (*deny \|\| *daemonAddr != "") && (*list → if (*deny) && (*list | `TestApproveDecideFlagRules` | KILLED | KILLED | TestApproveDecideFlagRules | K |
| C-MODES | M5c | approve.go: if modes != 0 && *signWith == "" { → if false { | `TestApproveDecideFlagRules` | KILLED | KILLED | TestApproveDecideSignWith | K |
| C-MODES-SIGNWITH | M5c | approve.go: if modes != 0 && *signWith == "" { → if modes != 0 { | `TestApproveDecideSignWith` | KILLED | KILLED | TestApproveDecideSignWith,TestApproveDecideSignWithSSHKeygen | K |
| C-VERB-DENY | M4b | approve.go:     verb = "deny" →     verb = "approve" | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideSignWith,TestApproveDecideTyped | K |
| C-HEALTH-ERR | M4b | approve_decide.go: h, err := fetchDaemonHealth(hostPort)⏎ if err != nil { → h, err := fetchDaemonHealth(hostPort)⏎ if false && err != nil { | `TestApproveDecideRefusesAMismatchedSubject` | SURVIVED | SURVIVED | — | E: equivalent: with the health error ignored the zero health names no store, so the next check (`SameDatabase("", --db)` is false) refuses anyway; only the text differs |
| C-SAME-DB | M4b | approve_decide.go: store.SameDatabase(h.DBPath, o.db); !same { → store.SameDatabase(h.DBPath, o.db); false && !same { | `TestApproveDecideRefusesAMismatchedSubject` | KILLED | KILLED | TestApproveDecideRefusesAMismatchedSubject | K |
| C-SAME-CREDS | M4b | approve_decide.go: store.SameDatabase(h.CredentialsDB, credPath); !same { → store.SameDatabase(h.CredentialsDB, credPath); false && !same { | `TestApproveDecideRefusesAMismatchedSubject` | KILLED | KILLED | TestApproveDecideRefusesAMismatchedSubject | K |
| C-GET-STATUS | M4b | approve_decide.go:  if status != http.StatusOK {⏎  return refuse("the daemon will not show →  if false {⏎  return refuse("the daemon will not show | `TestApproveDecideRefusesAMismatchedSubject` | KILLED | KILLED | TestApproveDecideRefusesAMismatchedSubject | K |
| C-STRICT | M4b | approve_decide.go:  dec.DisallowUnknownFields()⏎ → ⏎ | `TestApproveDecideRefusesAMismatchedSubject` | KILLED | KILLED | TestApproveDecideRefusesAMismatchedSubject | K |
| C-REQ-HASH | M4b | approve_decide.go: \|\| sha256Ref(reqBytes) != ref  →  | `TestApproveDecideRefusesAMismatchedSubject` | KILLED | KILLED | TestApproveDecideRefusesAMismatchedSubject | K |
| C-PKT-HASH | M4b | approve_decide.go: \|\| sha256Ref(pktBytes) != s.Packet { → { | `TestApproveDecideRefusesAMismatchedSubject` | KILLED | KILLED | TestApproveDecideRefusesAMismatchedSubject | K |
| C-BASE64 | M4b | approve_decide.go: if err1 != nil \|\| err2 != nil \|\|  → if _, _ = err1, err2;  | `TestApproveDecideRefusesAMismatchedSubject` | SURVIVED | SURVIVED | — | E: equivalent: bytes that do not decode as base64url come back empty or truncated, and those cannot hash to the ref, so the hash check refuses them (a hash preimage aside) |
| C-REQ-DECODE | M4b | approve_decide.go: broker.DecodeApprovalRequest(reqBytes); err != nil { → broker.DecodeApprovalRequest(reqBytes); false && err != nil { | `TestApproveDecideRefusesAMismatchedSubject` | KILLED | KILLED | TestApproveDecideRefusesAMismatchedSubject | K |
| C-PKT-DECODE | M4b | approve_decide.go: coordinator.DecodeDecisionPacket(pktBytes); err != nil \|\| s.packet → coordinator.DecodeDecisionPacket(pktBytes); s.packet | `TestApproveDecideRefusesAMismatchedSubject` | SURVIVED | SURVIVED | — | E: equivalent: a packet that does not decode is the zero packet, whose `RequestRef` never equals the ref, so the carries check refuses it |
| C-PKT-CARRIES | M4b | approve_decide.go: err != nil \|\| s.packet.RequestRef.String() != ref { → err != nil { | `TestApproveDecideRefusesAMismatchedSubject` | KILLED | KILLED | TestApproveDecideRefusesAMismatchedSubject | K |
| C-DISPLAY-REFUSAL | M4b | approve_decide.go: if s.Refusal != "" { → if false { | `TestApproveDecideRefusesAMismatchedSubject` | KILLED | KILLED | TestApproveDecideRefusesAMismatchedSubject | K |
| C-TTY | M4b | approve_decide.go: term, code := openApproverTerminal(env, stderr, "ailang-worldd approve: ", nothing)⏎ if code != exitOK { → term, code := openApproverTerminal(env, stderr, "ailang-worldd approve: ", nothing)⏎ if false && code != exitOK { | `TestApproveDecideNoTerminal` | KILLED | KILLED | TestApproveDecideNoTerminal | K |
| C-CONFIRM | M4b | approve_decide.go: _ = term.Close()⏎ if code != exitOK { → _ = term.Close()⏎ if false && code != exitOK { | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-PROOF-NAME | M4b | approve_decide.go: "decision-"+strings.TrimPrefix(o.ref, "sha256:") → "decision-"+o.ref | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-PROOF-PHASE | M4b | approve_decide.go: "decide-"+o.verb) → "decide-approve") | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-PROOF-ERR | M4b | approve_decide.go: proof, err := writeOperatorProof(credPath, "decision-"+strings.TrimPrefix(o.ref, "sha256:"), "decide-"+o.verb)⏎ if err != nil { → proof, err := writeOperatorProof(credPath, "decision-"+strings.TrimPrefix(o.ref, "sha256:"), "decide-"+o.verb)⏎ if false && err != nil { | `TestApproveDecideTyped` | SURVIVED | SURVIVED | — | E: survivor (untested defensive arm): `writeOperatorProof` fails only on an unusable 0700 proofs directory or a `crypto/rand` failure (row 236's refusals, tested there); mutated, the POST carries no valid proof and the daemon refuses 403 |
| C-TYPED | M4b | approve_decide.go: Expiry: s.Window.Max, Typed: true} → Expiry: s.Window.Max, Typed: false} | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-EXPIRY | M4b | approve_decide.go: Expiry: s.Window.Max, Typed: true} → Expiry: s.Window.Min, Typed: true} | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-ORIGIN | M4b | approve_decide.go: hdr := map[string]string{"Origin": origin} → hdr := map[string]string{"Origin": "http://localhost"} | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-PROOF-HDR | M4b | approve_decide.go: if proof != "" { → if false { | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-POST-STATUS | M4b | approve_decide.go:  if status != http.StatusOK {⏎  return refuse("the daemon refused →  if false {⏎  return refuse("the daemon refused | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-MADE-REQUEST | M4b | approve_decide.go: err != nil \|\| made.Request != ref { → err != nil { | `TestApproveDecideTyped` | SURVIVED | SURVIVED | — | E: survivor (defensive): the daemon's 200 always names its request (`DecisionMade`); the check guards a foreign answer the fake never sends |
| C-RENDER-REQUESTER | M4b | approve_decide.go: r.Effect, r.Scope, r.Cost, r.Requester, r.Now → r.Effect, r.Scope, r.Cost, r.Effect, r.Now | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-RENDER-VERB | M4b | approve_decide.go: p.Policy, strings.ToUpper(verb)) → p.Policy, verb) | `TestApproveDecideTyped` | KILLED | KILLED | TestApproveDecideTyped | K |
| C-SIGNWITH-DISPATCH | M5c | approve_decide.go: if o.signWith != "" { → if false { | `TestApproveDecideSignWith` | KILLED | KILLED | TestApproveDecideSignWith | K |
| C-SIGN-NIL | M5c | approve_decide.go: if env.signDecision == nil { → if false { | `TestApproveDecideSignWith` | SURVIVED | SURVIVED | — | E: survivor (defensive): `realSessionEnv` sets `signDecision` (CD5 pins it by reflection) and every test env reaching `--sign-with` sets it |
| C-SIGN-VERB | M5c | approve_decide.go: verb, _ := approval.DecisionVerb(o.verb) → verb := approval.VerbApprove | `TestApproveDecideSignWith` | KILLED | KILLED | TestApproveDecideSignWith | K |
| C-SIGN-EXPIRY | M5c | approve_decide.go: Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max} → Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Min} | `TestApproveDecideSignWith` | KILLED | KILLED | TestApproveDecideSignWith,TestApproveDecideSignWithSSHKeygen | K |
| C-SIGN-DEADLINE | M5c | approve_decide.go: time.Duration(max(body.Expiry-env.now(), 1))*time.Second → time.Hour | `TestApproveDecideSignWith` | KILLED | KILLED | TestApproveDecideSignWith | K |
| C-SIGN-ERR | M5c | approve_decide.go: armored, err := env.signDecision(ctx, o.signWith, msg, stderr)⏎ if err != nil { → armored, err := env.signDecision(ctx, o.signWith, msg, stderr)⏎ if false && err != nil { | `TestApproveDecideSignWith` | KILLED | KILLED | TestApproveDecideSignWith | K |
| C-SIGN-SSHSIG | M5c | approve_decide.go:  body.SSHSIG = &sig⏎ →  _ = sig⏎ | `TestApproveDecideSignWith` | KILLED | KILLED | TestApproveDecideSignWith,TestApproveDecideSignWithSSHKeygen | K |
| C-SIGN-NAMESPACE | M5c | approve_decide.go: sshsign.SignIn(ctx, "ssh-keygen", approval.DecisionNamespace, key, msg, stderr) → sshsign.SignIn(ctx, "ssh-keygen", approval.Namespace, key, msg, stderr) | `TestApproveDecideSignWithSSHKeygen` | KILLED | KILLED | TestApproveDecideSignWithSSHKeygen | K |
| C-SSHSIGN-NS | M5c | sshsign.go: "-Y", "sign", "-n", namespace, "-f", key) → "-Y", "sign", "-n", approval.Namespace, "-f", key) | `TestApproveDecideSignWithSSHKeygen` | KILLED | KILLED | TestApproveDecideSignWithSSHKeygen | K |
| A-CHECKAPPROVER | M4d | decision_signed.go: in.Approver, FormatSSH); err != nil { → in.Approver, FormatSSH); false && err != nil { | `TestDecideBySignatureApprover` | KILLED | KILLED | TestDecideBySignatureApprover | K |
| A-FORMAT | M4d | decision_signed.go: in.Approver, FormatSSH); err → in.Approver, in.Approver.Format); err | `TestDecideBySignatureApprover` | KILLED | KILLED | TestDecideBySignatureApprover | K |
| A-NOW | M4d | decision_signed.go: AdmitInput{Now: in.Now} → AdmitInput{Now: in.Expiry} | `TestDecideBySignatureApprover` | KILLED | KILLED | TestDecideBySignatureApprover | K |
| A-VERIFY | M4d | decision_signed.go: in.Expiry, in.Approver.Key); err != nil { → in.Expiry, in.Approver.Key); false && err != nil { | `TestDecisionBindsTheVerb` | KILLED | KILLED | TestDecisionBindsTheVerb,TestMintSignatureIsNotADecision | K |
| A-NAMESPACE | M4d | decision_signed.go: return VerifySSHSIG(s, DecisionNamespace, msg, key) → return VerifySSHSIG(s, Namespace, msg, key) | `TestMintSignatureIsNotADecision` | KILLED | KILLED | TestDecideBySignatureApprover,TestDecisionBindsTheVerb,TestMintSignatureIsNotADecision,TestSignedDecisionBinds | K |
| A-MSG-VERB | M4d | decision_signed.go: msg, err := DecisionMessage(verb, packet, request, expiry) → msg, err := DecisionMessage(VerbApprove, packet, request, expiry) | `TestDecisionBindsTheVerb` | KILLED | KILLED | TestDecisionBindsTheVerb,TestSignedDecisionBinds | K |
| A-BINDS-NIL | M4d | decision_signed.go:  case s == nil \|\| !s.ok:⏎ →  case s == nil:⏎ | `TestSignedDecisionBinds` | KILLED | KILLED | TestSignedDecisionBinds | K |
| A-BINDS-VERB | M4d | decision_signed.go:  case s.verb != verb:⏎ →  case false:⏎ | `TestSignedDecisionBinds` | KILLED | KILLED | TestSignedDecisionBinds | K |
| A-BINDS-PACKET | M4d | decision_signed.go:  case s.packet != packet:⏎ →  case false:⏎ | `TestSignedDecisionBinds` | KILLED | KILLED | TestSignedDecisionBinds | K |
| A-BINDS-REQUEST | M4d | decision_signed.go:  case s.request != request:⏎ →  case false:⏎ | `TestSignedDecisionBinds` | KILLED | KILLED | TestSignedDecisionBinds | K |
| A-BINDS-EXPIRY | M4d | decision_signed.go:  case s.expiry != expiry:⏎ →  case false:⏎ | `TestSignedDecisionBinds` | KILLED | KILLED | TestSignedDecisionBinds | K |
| A-DECIDEDBY-OK | M4d | decision_signed.go: func (s *SignedDecision) DecidedBy() string {⏎ if s == nil \|\| !s.ok { → func (s *SignedDecision) DecidedBy() string {⏎ if s == nil { | `TestSignedDecisionBinds` | KILLED | KILLED | TestSignedDecisionBinds | K |
| A-USE-OK | M4d | decision_signed.go: func (s *SignedDecision) Use() (name, format string, key []byte, ok bool) {⏎ if s == nil \|\| !s.ok { → func (s *SignedDecision) Use() (name, format string, key []byte, ok bool) {⏎ if s == nil { | `TestSignedDecisionBinds` | KILLED | KILLED | TestSignedDecisionBinds | K |
| A-DECIDEDBY-FORMAT | M4d | decision_signed.go: return s.format + ":" + s.name → return "ssh:" + s.name | `TestDecideBySignatureApprover` | SURVIVED | SURVIVED | — | E: equivalent in sprint 2: the only decision format is ssh (P3 refuses others), so `format + ":"` is `"ssh:"`; it goes live with sprint 3's webauthn arm |
| A-KEY-CLONE | M4d | decision_signed.go: key: bytes.Clone(in.Approver.Key.Blob)} → key: nil} | `TestDecideBySignatureApprover` | KILLED | KILLED | TestDecideBySignatureApprover | K |
| S-TABLES8 | M4e | credentials.go:  case version >= 8:⏎ →  case version >= 9:⏎ | `TestCredentialsSchema8KeyDecisions` | KILLED | KILLED | TestCredentialsSchema7SurfaceApprovals,TestCredentialsSchema8KeyDecisions,TestCredentialsUpgradeToV4,TestCredentialsUpgradeToV4/v1,TestCredentialsUpgradeToV4/v2 | K |
| S-UPGRADE8 | M4e | credentials.go:  if version < 8 {⏎ →  if false {⏎ | `TestCredentialsSchema8KeyDecisions` | KILLED | KILLED | TestCredentialsSchema7SurfaceApprovals,TestCredentialsSchema8KeyDecisions,TestCredentialsUpgradeToV4,TestCredentialsUpgradeToV4/v1,TestCredentialsUpgradeToV4/v2 | K |
| S-K1-DEADLINE | M4e | key_decisions.go: func (c *Credentials) RecordKeyDecision(ctx context.Context, row KeyDecisionRow, use ApproverUse, now int64) error {⏎ if err := requireDeadline(ctx); err != nil { → func (c *Credentials) RecordKeyDecision(ctx context.Context, row KeyDecisionRow, use ApproverUse, now int64) error {⏎ if err := requireDeadline(ctx); false && err != nil { | `TestCredentialsMethodsRequireADeadline` | KILLED | SURVIVED | TestCredentialsMethodsRequireADeadline | I: killed only by the inventory `TestCredentialsMethodsRequireADeadline` (named); the behavioural suite survives: every caller passes a deadline |
| S-K1-VERSION | M4e | key_decisions.go:  if c.version < keyDecisionsSchemaVersion { // K1⏎ →  if false { // K1⏎ | `TestCredentialsSchema8KeyDecisions` | KILLED | KILLED | TestCredentialsSchema8KeyDecisions | K |
| S-LAND-VERSION | M4e | key_decisions.go:   return c.olderForKeyDecisions()⏎ }⏎ if !isDecisionRef(decisionRef) { →   _ = 0⏎ }⏎ if !isDecisionRef(decisionRef) { | `TestCredentialsSchema8KeyDecisions` | KILLED | KILLED | TestCredentialsSchema8KeyDecisions | K |
| S-READ-VERSION | M4e | key_decisions.go:   return KeyDecisionRow{}, false, c.olderForKeyDecisions() →   _ = 0 | `TestCredentialsSchema8KeyDecisions` | KILLED | KILLED | TestCredentialsSchema8KeyDecisions | K |
| S-K2-REQUEST | M4e | key_decisions.go: case !isDecisionRef(row.RequestRef) \|\| !isDecisionRef(row.PacketRef): → case !isDecisionRef(row.PacketRef): | `TestRecordKeyDecisionGrammar` | KILLED | KILLED | TestRecordKeyDecisionGrammar | K |
| S-K2-PACKET | M4e | key_decisions.go: case !isDecisionRef(row.RequestRef) \|\| !isDecisionRef(row.PacketRef): → case !isDecisionRef(row.RequestRef): | `TestRecordKeyDecisionGrammar` | KILLED | KILLED | TestRecordKeyDecisionGrammar | K |
| S-ISREF-HEX | M4e | key_decisions.go: s[:7] == "sha256:" && isLowerHex64(s[7:]) → s[:7] == "sha256:" | `TestRecordKeyDecisionGrammar` | KILLED | KILLED | TestRecordKeyDecisionGrammar | K |
| S-K2-VERB | M4e | key_decisions.go: case row.Verb != "approve" && row.Verb != "deny": → case false: | `TestRecordKeyDecisionGrammar` | KILLED | KILLED | TestRecordKeyDecisionGrammar | K |
| S-K2-FORMAT | M4e | key_decisions.go: case row.Format != use.Format \|\| row.Approver → case row.Approver | `TestRecordKeyDecisionGrammar` | KILLED | KILLED | TestRecordKeyDecisionGrammar | K |
| S-K2-APPROVER | M4e | key_decisions.go: case row.Format != use.Format \|\| row.Approver != use.Format+":"+use.Name: → case row.Format != use.Format: | `TestRecordKeyDecisionGrammar` | KILLED | KILLED | TestRecordKeyDecisionGrammar | K |
| S-K2-ARTIFACT | M4e | key_decisions.go: case len(row.Artifact) < 1 \|\| len(row.Artifact) > MaxSurfaceApprovalArtifactBytes: → case false: | `TestRecordKeyDecisionGrammar` | KILLED | KILLED | TestRecordKeyDecisionGrammar | K |
| S-SCAN-NOROWS | M4e | key_decisions.go: if err != nil && !errors.Is(err, sql.ErrNoRows) { → if err != nil { | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded,TestRecordKeyDecisionRevokedSince | K |
| S-K0-LANDED | M4e | key_decisions.go: if found && old.DecisionRef != "" { // K0: landed → if false { // K0: landed | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded | K |
| S-K0-EXACT | M4e | key_decisions.go: exact := found && bytes.Equal(old.Artifact, row.Artifact) → exact := false && bytes.Equal(old.Artifact, row.Artifact) | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded | K |
| S-K0-PACKET | M4e | key_decisions.go: && old.PacketRef == row.PacketRef  →  | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded | K |
| S-K0-VERB | M4e | key_decisions.go: && old.Verb == row.Verb && → && | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded | K |
| S-K0-FORMAT | M4e | key_decisions.go: old.Format == row.Format && old.Approver → old.Approver | `TestRecordKeyDecisionReplaceUnlanded` | SURVIVED | SURVIVED | — | E: equivalent: K2 holds `row.Format == use.Format`, and one artifact (its SSHSIG names its key) is judged to exactly one approver, so the same artifact never arrives with another format |
| S-K0-APPROVER | M4e | key_decisions.go: && old.Approver == row.Approver && old.Expiry → && old.Expiry | `TestRecordKeyDecisionReplaceUnlanded` | SURVIVED | SURVIVED | — | E: equivalent, the same reason: the judge looks the approver up by the signature's own key, so the same artifact derives the same approver |
| S-K0-EXPIRY | M4e | key_decisions.go: && old.Expiry == row.Expiry // K0 → // K0 | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded | K |
| S-K3-BOTH | M4e | key_decisions.go: if err := useApprover(ctx, tx, use, now); err != nil { // K3 → if err := useApprover(ctx, tx, use, now); err != nil && !exact { // K3 | `TestRecordKeyDecisionRevokedSince` | KILLED | KILLED | TestRecordKeyDecisionRevokedSince | K |
| S-K3-DEL | M4e | key_decisions.go: if err := useApprover(ctx, tx, use, now); err != nil { // K3 → if err := useApprover(ctx, tx, use, now); false && err != nil { // K3 | `TestRecordKeyDecisionRevokedSince` | KILLED | KILLED | TestRecordKeyDecisionRevokedSince | K |
| S-K5-NOWRITE | M4e | key_decisions.go:  if !exact { // K5⏎ →  if exact { // K5⏎ | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded | K |
| S-K5-WHERE | M4e | key_decisions.go:  WHERE key_decisions.decision_ref IS NULL` → ` | `TestRecordKeyDecisionReplaceUnlanded` | SURVIVED | SURVIVED | — | E: equivalent: K0 refuses a landed row before the upsert; the `WHERE` is a second fence |
| S-DECIDED-AT | M4e | key_decisions.go: row.Artifact, row.Expiry, now); err != nil { → row.Artifact, row.Expiry, row.DecidedAt); err != nil { | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded | K |
| S-LAND-REF | M4e | key_decisions.go:  if !isDecisionRef(decisionRef) { →  if false { | `TestRecordKeyDecisionGrammar` | KILLED | KILLED | TestRecordKeyDecisionGrammar | K |
| S-LAND-NULL | M4e | key_decisions.go: WHERE request_ref = ? AND decision_ref IS NULL`, decisionRef, request) → WHERE request_ref = ?`, decisionRef, request) | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded | K |
| S-LAND-N | M4e | key_decisions.go: err != nil \|\| n != 1 { → err != nil \|\| n < 0 { | `TestRecordKeyDecisionReplaceUnlanded` | KILLED | KILLED | TestRecordKeyDecisionReplaceUnlanded | K |
| S-READ-NOROWS | M4e | key_decisions.go:  if errors.Is(err, sql.ErrNoRows) {⏎  return KeyDecisionRow{}, false, nil →  if false {⏎  return KeyDecisionRow{}, false, nil | `TestRecordKeyDecisionGrammar` | KILLED | KILLED | TestRecordKeyDecisionGrammar | K |
| J-PARSE-ERR | M4g | judge_decision.go: sig, err := approval.ParseSSHSIG(armored)⏎ if err != nil { → sig, err := approval.ParseSSHSIG(armored)⏎ if false && err != nil { | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| J-STORE-ERR | M4g | judge_decision.go: ro.ApproverByKey(ctx, sig.PublicKey.Blob)⏎ if err != nil { → ro.ApproverByKey(ctx, sig.PublicKey.Blob)⏎ if false && err != nil { | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| J-UNKNOWN | M4g | judge_decision.go:  if !ok { →  if _ = ok; false { | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| J-FINGERPRINT | M4g | judge_decision.go: Fingerprint: sig.PublicKey.Fingerprint(), → Fingerprint: "", | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| J-LAW-ERR | M4g | judge_decision.go: RevokedAt: enrolled.RevokedAt}})⏎ if err != nil { → RevokedAt: enrolled.RevokedAt}})⏎ if false && err != nil { | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| J-NOW | M4g | judge_decision.go: Expiry: in.Expiry, Now: in.Now, → Expiry: in.Expiry, Now: in.Expiry, | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| J-VERB | M4g | judge_decision.go: DecisionSignatureInput{Verb: in.Verb, → DecisionSignatureInput{Verb: approval.VerbDeny, | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| J-FORMAT | M4g | judge_decision.go: Format: enrolled.Format, → Format: approval.FormatSSH, | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| J-REVOKED | M4g | judge_decision.go: RevokedAt: enrolled.RevokedAt} → RevokedAt: 0} | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| J-ENROLLED | M4g | judge_decision.go: EnrolledAt: enrolled.EnrolledAt, → EnrolledAt: 0, | `TestJudgeDecisionSSH` | KILLED | KILLED | TestJudgeDecisionSSH | K |
| D-ARM-SWITCH | M5a | surface_decide.go:  if body.SSHSIG != nil {⏎ →  if false {⏎ | `TestDecidePOSTSSHSIG` | KILLED | KILLED | TestDecideConcurrentPOSTsDecideOnce,TestDecidePOSTBusyWriter,TestDecidePOSTRetryAfterAFailedWorldWrite,TestDecidePOSTSSHSIG,TestDecidePOSTSSHSIGRefusals | K |
| D-KEYARM-OK | M5a | surface_decide.go: d.keyArm(ctx, w, r, f, body, now); !ok { → d.keyArm(ctx, w, r, f, body, now); false && !ok { | `TestDecidePOSTSSHSIGRefusals` | KILLED | KILLED | TestDecidePOSTBusyWriter | K |
| D-DECIDEDBY | M5a | surface_decide.go: body.Verb, decidedBy, now) // 3. → body.Verb, typedDecider, now) // 3. | `TestDecidePOSTSSHSIG` | KILLED | KILLED | TestDecidePOSTSSHSIG | K |
| D-DECIDEDBY-RESP | M5a | surface_decide.go: DecidedBy: decidedBy}) → DecidedBy: typedDecider}) | `TestDecidePOSTSSHSIG` | KILLED | KILLED | TestDecidePOSTSSHSIG | K |
| D-LAND | M5a | surface_decide.go:  if evidence != nil { // 4. →  if false { // 4. | `TestDecidePOSTSSHSIG` | KILLED | KILLED | TestDecidePOSTRetryAfterAFailedWorldWrite,TestDecidePOSTSSHSIG | K |
| D-LAND-ERR | M5a | surface_decide.go: evidence.LandKeyDecision(ctx, ref, decision.String()); err != nil { → evidence.LandKeyDecision(ctx, ref, decision.String()); false && err != nil { | `TestDecidePOSTSSHSIG` | SURVIVED | SURVIVED | — | E: survivor (unreachable without a store fault between `RecordKeyDecision`'s commit and `LandKeyDecision` on the same open writer); the chain stays authoritative |
| D-WRITER-CLOSE | M5a | surface_decide.go:   defer func() { _ = evidence.Close() }()⏎ → ⏎ | `TestDecidePOSTRetryAfterAFailedWorldWrite` | KILLED | KILLED | TestDecidePOSTRetryAfterAFailedWorldWrite,TestDecidePOSTSSHSIG | K |
| D-DECIDE-ERR | M5a | surface_decide.go: body.Verb, decidedBy, now) // 3.⏎ if err != nil { → body.Verb, decidedBy, now) // 3.⏎ if false && err != nil { | `TestDecidePOSTRetryAfterAFailedWorldWrite` | KILLED | KILLED | TestDecidePOSTRetryAfterAFailedWorldWrite | K |
| D-VERB-BYTE | M5a | surface_decide_key.go: verb, _ := approval.DecisionVerb(body.Verb) // D9 held → verb := approval.VerbApprove // D9 held | `TestDecidePOSTSSHSIG` | KILLED | KILLED | TestDecidePOSTRetryAfterAFailedWorldWrite,TestDecidePOSTSSHSIG,TestDecidePOSTSSHSIGRefusals | K |
| D-RO-ERR | M5a | surface_decide_key.go: ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)⏎ if err != nil { → ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)⏎ if false && err != nil { | `TestDecidePOSTSSHSIG` | SURVIVED | SURVIVED | — | E: survivor (unreachable: the read-only open fails only on an unreadable credentials file, which `serve` created) |
| D-BINDS | M5a | surface_decide_key.go:  if err == nil {⏎  err = signed.Binds(verb, f.Packet, f.Request, body.Expiry)⏎ }⏎ →  | `TestDecidePOSTSSHSIGRefusals` | SURVIVED | SURVIVED | — | E: equivalent: `JudgeDecisionSSH` verifies exactly (verb, P, Q, expiry) and the `SignedDecision` carries what it verified, so `Binds` over the same values never refuses (kept, as the design's §2.5 step) |
| D-JUDGE-STORE | M5a | surface_decide_key.go: case errors.As(err, &je) && je.Step == authority.JudgeStepStore: → case errors.As(err, &je) && false: | `TestDecidePOSTSSHSIGRefusals` | KILLED | KILLED | TestDecidePOSTSSHSIGRefusals | K |
| D-422 | M5a | surface_decide_key.go:  case err != nil:⏎  writeApprovalRefusal →  case false:⏎  writeApprovalRefusal | `TestDecidePOSTSSHSIGRefusals` | KILLED | KILLED | TestDecidePOSTSSHSIGRefusals | K |
| D-WRITER-ERR | M5a | surface_decide_key.go: wr, err := store.OpenCredentials(d.credentialsPath) // non-waiting: a busy writer is 503 Busy⏎ if err != nil { → wr, err := store.OpenCredentials(d.credentialsPath) // non-waiting: a busy writer is 503 Busy⏎ if false && err != nil { | `TestDecidePOSTBusyWriter` | KILLED | KILLED | TestDecidePOSTBusyWriter | K |
| D-PACKETREF | M5a | surface_decide_key.go: PacketRef: "sha256:" + hex.EncodeToString(f.Packet[:]), → PacketRef: "sha256:" + hex.EncodeToString(f.Request[:]), | `TestDecidePOSTSSHSIG` | KILLED | KILLED | TestDecidePOSTSSHSIG | K |
| D-ARTIFACT | M5a | surface_decide_key.go: Approver: signed.DecidedBy(), Artifact: artifact, → Approver: signed.DecidedBy(), Artifact: []byte("x"), | `TestDecidePOSTSSHSIG` | KILLED | KILLED | TestDecidePOSTSSHSIG | K |
| D-EXPIRY-ROW | M5a | surface_decide_key.go: Artifact: artifact, Expiry: body.Expiry} → Artifact: artifact, Expiry: 0} | `TestDecidePOSTSSHSIG` | KILLED | KILLED | TestDecidePOSTSSHSIG | K |
| D-RECORD-NOW | M5a | surface_decide_key.go: PublicKey: key}, now) → PublicKey: key}, now+1) | `TestDecidePOSTSSHSIG` | KILLED | KILLED | TestDecidePOSTSSHSIG | K |
| D-RECORD-ERR | M5a | surface_decide_key.go: PublicKey: key}, now)⏎ if err != nil { → PublicKey: key}, now)⏎ if false && err != nil { | `TestDecidePOSTSSHSIGRefusals` | KILLED | KILLED | TestDecidePOSTSSHSIGRefusals | K |
| D-RACE | M5a | surface_decide_key.go:   if isApproverRace(err) { →   if false { | `TestDecidePOSTSSHSIGRefusals` | SURVIVED | SURVIVED | — | E: survivor (unreachable through the POST: the read-only judge refuses a revoked or re-keyed approver first; KD3/KD9 prove the writer's refusal at the store, §1.4 P7) |
| M-NOFORMAT | M5a | decide_model.go: case !body.Typed && body.SSHSIG == nil: → case false: | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| M-TWOFORMATS | M5a | decide_model.go: case body.Typed && body.SSHSIG != nil: → case false: | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| J07-HASSUFFIX | M4a | decision_law.go: strings.Contains(scope, "effect="+EffectRegistryPublish) → strings.HasSuffix(scope, "effect="+EffectRegistryPublish) | `TestReservedApprovalScope` | KILLED | KILLED | TestReservedApprovalScope | K |
| J31-ANCHOR | M4a | surface_decide.go: `^sha256:[0-9a-f]{64}$` → `^sha256:[0-9a-f]{64}` | `TestDecideRefGrammarIsAnchored` | KILLED | KILLED | TestDecideRefGrammarIsAnchored | K |
| J32-CASE | M4a | surface_decide.go: `^sha256:[0-9a-f]{64}$` → `^sha256:[0-9a-fA-F]{64}$` | `TestDecideRefGrammarIsAnchored` | KILLED | KILLED | TestDecideRefGrammarIsAnchored | K |
| J33-GET-D4A | M4a | surface_decide.go: case "D1", "D2", "D3", "D4a", "D4b": → case "D1", "D2", "D3", "D4b": | `TestDecisionSubjectRefusesAnUndecodablePacket` | KILLED | KILLED | TestDecisionSubjectRefusesAnUndecodablePacket | K |
| LAW-D5-KEY | M5b | decision_law.go: ReservedApprovalScope(f.Scope); reserved { → ReservedApprovalScope(f.Scope); reserved && false { | `TestDecideRefusesReservedScopesByKey` | KILLED | KILLED | TestDecideRefusesReservedScopesByKey,TestDecideRefusesReservedScopesTyped,TestDecisionSubjectJSON | K |
| LAW-D5-L2 | M5b | decision_law.go: ReservedApprovalScope(f.Scope); reserved { → ReservedApprovalScope(f.Scope); reserved && false { | `TestPublishApprovalRefusesAKeyDecision` | SURVIVED | KILLED | TestDecisionRefusalD5,TestDecisionRefusalOrder | L: must live: on the L1-deleted mutant, sprint 1's L2 test still PASSES (R4e's must-die pair with `LAW-D5-KEY`) |
| H-LOCK | M5b | surface_decide.go:  d.humanMu.Lock()⏎ defer d.humanMu.Unlock()⏎ →  | `TestDecideConcurrentPOSTsDecideOnce` | KILLED | KILLED | TestDecideConcurrentPOSTsDecideOnce | K |

---

## 5. Files (final `git status --short` after M5c)

| File | Sub-milestone | Change |
|---|---|---|
| `host/broker/decision_law_test.go` | M4a | J07's row |
| `host/daemon/surface_decide_jgaps_test.go` | M4a | new: J31/J32, J33 |
| `cmd/ailang-worldd/approve_decide.go` | M4b, M5c | new: the decision mode; M5c `--sign-with` |
| `cmd/ailang-worldd/approve.go` | M4b, M5c | `--deny`, `--daemon`, the dispatch, the rules, help |
| `cmd/ailang-worldd/main.go`, `help_test.go` | M4b | one usage line; its inventory line |
| `cmd/ailang-worldd/approve_decide_rig_test.go` | M4b | new: the fake daemon; CD2, CD3 |
| `host/store/context_roots_test.go` | M4b, M5c | +3 roots |
| `cmd/ailang-worldd/approve_decide_test.go` | M4c | new: CD1, CD4 |
| `host/approval/decision_signed.go`, `decision_signed_test.go` | M4d | new: the law; DM3–DM6, DM8 |
| `host/store/key_decisions.go`, `credentials.go`, `credentials_schema.sql` | M4e | new table, schema 8 |
| `host/store/key_decisions_test.go`; `credentials_test.go`, `approvers_test.go`, `mint_requests_test.go`, `surface_approvals_test.go` | M4e | new: KD1, KD5, KD6; P2 edits |
| `host/daemon/surface_approve_{post,read}_test.go`, `host/authority/judge_test.go`, seven `cmd/ailang-worldd/*_test.go` | M4e | P2: drop `key_decisions` in downgrade fixtures |
| `host/store/key_decisions_record_test.go` | M4f | new: KD3/KD9, KD4/KD7, K2 |
| `host/authority/judge_decision.go`, `judge_decision_test.go` | M4g | new: the judge; JD1–JD3 |
| `host/daemon/surface_decide_key.go`; `surface_decide.go`; `host/surfacepage/decide_model.go` | M5a | new key arm; the slot; the grammar |
| `host/daemon/surface_decide_key_test.go`; `surface_approve_post_guard_test.go` | M5a | new: RD10, RD11, refusals; writer sites 1 → 2 |
| `host/daemon/surface_decide_key_retry_test.go`; `surface_decide_discovery_test.go` | M5b | new: RD12, RD13 (a), R4e, H-LOCK; RD9b/RD9c by key |
| `host/sshsign/sshsign.go`; `cmd/ailang-worldd/session.go`; `approve_decide_key_test.go` | M5c | `SignIn`; `signDecision`; new: CD5 + sibling |
| this plan and `sprint_w-world-approvals-by-key-a4-s2.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each measured at the base, §1.1)

In **one** `bash` shell (macOS `/bin/bash` 3.2: no `declare -A`), from the worktree, **not under `nohup`**, nothing
else running. **Every leg records its own rc and the script runs every leg** (no `set -e`); every leg `tee`s its
full output to `$HOME/.ailang/state/world-iter390/gates/<leg>.log`, and `legs.txt` holds one `LEG <name> rc=<n>`
line per leg plus the inventory verdict. The planner's `plan-scratch/tools/gates.sh` is exactly this script; copy
it, never retype it:

```
mkdir -p $HOME/.ailang/state/world-iter390/gates
/bin/bash $HOME/.ailang/state/world-iter390/plan-scratch/tools/gates.sh /Users/voightkampff/dev/sunholo-data/.wt-world-iter390-a4-s2 $HOME/.ailang/state/world-iter390/gates final
cat $HOME/.ailang/state/world-iter390/gates/legs.txt
```

The script (verbatim, `tools/gates.sh`):

```bash
#!/bin/bash
# A4 sprint 2 pre-commit gates (plan §6). usage: gates.sh <tree> <logdir> <base|final>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e).
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.log, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.log"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs 'go test -count=1 -timeout 30m ./host/approval/... ./host/authority/ ./host/store/ ./host/sshsign/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/'
leg broker_named 'go test -count=1 ./host/broker/ -run "^(TestReservedApprovalScope|TestDecisionRefusal|TestDecideApproval|TestPublishApproval|TestApprovalContextIdentity|TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential)" -v'
leg broker_whole_reported 'go test -count=1 -timeout 30m ./host/broker/'
leg race 'go test -race -count=1 ./host/approval/ ./host/authority/ ./host/store/ ./host/sshsign/ && go test -race -count=1 ./host/daemon/ -run "^(TestDecision|TestDecide|TestApprov|TestSurfaceApprove|TestDaemonCredentialWriterSites|TestEpisodeFinish)" && go test -race -count=1 ./cmd/ailang-worldd/ -run "^(TestApproveDecide|TestApprove|TestApprover)"'
NEW_A="TestDecideBySignatureApprover TestMintSignatureIsNotADecision TestDecisionBindsTheVerb TestSignedDecisionBinds"
NEW_AU="TestJudgeDecisionSSH"
NEW_S="TestCredentialsSchema8KeyDecisions TestKeyDecisionsRefChecks TestRecordKeyDecisionRevokedSince TestRecordKeyDecisionReplaceUnlanded TestRecordKeyDecisionGrammar"
NEW_D="TestDecideRefGrammarIsAnchored TestDecisionSubjectRefusesAnUndecodablePacket TestDecidePOSTSSHSIG TestDecidePOSTSSHSIGRefusals TestDecidePOSTBusyWriter TestDecidePOSTRetryAfterAFailedWorldWrite TestDecideRefusesReservedScopesByKey TestDecideConcurrentPOSTsDecideOnce"
NEW_C="TestApproveDecideTyped TestApproveDecideNoTerminal TestApproveDecideFlagRules TestApproveDecideRefusesAMismatchedSubject TestApproveDecideSignWith TestApproveDecideSignWithSSHKeygen"
if [ "$MODE" = base ]; then NEW_A=""; NEW_AU=""; NEW_S=""; NEW_D=""; NEW_C=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
INV_APPROVAL="TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestDecisionMessageLayout TestDecisionMessageRefusesBadInput"
INV_AUTH="TestJudgeSSHSignatureIsTheCLIJudge TestJudgeAssertionIsTheCLIJudge"
INV_BROKER="TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract TestApprovalContextIdentity TestReservedApprovalScope TestDecideApprovalCallerCensus"
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestReadCtxCancelledAfterHandler TestMCPMountMethodSource TestDaemonCredentialWriterSites TestQuickstartDecisionPacketSection TestDaemonConstructorCallerInventory TestDecidePOSTGuards TestDecideDiscoveryFindsAHiddenSecondPacket TestDecideDiscoveryRefusesAnOverflow"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsSchema7SurfaceApprovals"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_CMD="TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp TestRealSessionEnvSignsWithSSHKeygen TestSignWithSSHKeygenEndToEnd"
INV_RUNBOOK="TestNoCIStepOrScriptReachesThePublishEntrypoint"
leg inv "set -o pipefail
go test -v -count=1 ./host/approval/... -run '^($(echo $INV_APPROVAL | tr ' ' '|')$(alt $NEW_A))\$' &&
go test -v -count=1 ./host/authority/ -run '^($(echo $INV_AUTH | tr ' ' '|')$(alt $NEW_AU))\$' &&
go test -v -count=1 ./host/broker/ -run '^($(echo $INV_BROKER | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|')$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|')$(alt $NEW_S))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/runbook/ -run '^($(echo $INV_RUNBOOK | tr ' ' '|'))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$'"
# every named test must PASS (none SKIP, none FAIL)
MISSING=""
for t in $INV_APPROVAL $INV_AUTH $INV_BROKER $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_VG $INV_CMD $INV_RUNBOOK $NEW_A $NEW_AU $NEW_S $NEW_D $NEW_C; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || MISSING="$MISSING $t"
done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then FAILS="$FAILS inv_skip_or_fail"; fi
[ -z "$MISSING" ] || { echo "NOT PASSED:$MISSING" | tee -a "$L/legs.txt"; FAILS="$FAILS inv_missing"; }
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
```

Then, **after the controller commits** (the tracked-binary gate reads committed blobs; this leg runs `go test ./...`):

```
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter390-a4-s2 && VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter390/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter390/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter390-a4-s2 status --short   # only the two plan files, until the controller commits them
```

**Reading `legs.txt`.** Gating: `verify_ail`, `gofmt`, `build`, `vet`, `fence`, `pkgs`, `broker_named`, `race`, `inv`
(and `inv_missing`/`inv_skip_or_fail` must not appear). **Reported, not gating:** `broker_whole_reported` — any red
there outside the known list below, or any red among `TestDecide*`, `TestDecision*`, `TestReserved*`,
`TestPublishApproval*`, `TestApprovalContextIdentity`, `TestEverySubprocessSite*`, is a regression; a known-list red
is re-run alone once. A red `pkgs` or `race` leg on a known-list test is re-run alone once (the script does not stop,
so every other leg's verdict is already in its log).

**At the base and on the final tree:** §1.1's table and §1.2 C21–C26.

**Known load- or environment-sensitive (re-run once, alone; do not chase)**: `TestSeedRestartIgnoresASurvivingOldCopy`,
`TestHostsLaunchersExecHost`, `host/broker`'s exec-probe timeouts (incl. `TestExecProtectedCheckTimeout/reserve-after-before`,
`TestExecProbeNet6*`, `TestPipeCloseBoundsEscapedDescendant`, `TestFailedKillBoundsCallAndChildIsReapedLater`,
`TestReservationIsExactUnderConcurrency`), `TestPublishApprovalRefusalSetWithALandedPositiveControl`, local `-race` of
host/broker, host/daemon, host/verifygate, `TestMCPPostBudgetProductionConstants/batch`, `TestBackoff`
(host/surfacefollow), `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
`TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck` supervisor tests,
`TestMarketingCaptureInventoryMutation`, `TestSessionSelfGit`, `TestServeStartsWhenTheProbeSpendsTheStartupClock`
(red at the base), `TestRegistryPhaseHasItsOwnClock`, `TestWorkspaceProtectedPolicyWiring`,
`TestConvertedStoreMethodsM5b`, host/pkgproj under load, `host/coordinator`'s `TestCapsuleWarmEqualsColdForSeTools`,
`cmd/world-publish`'s interpreter-timeout reds.

**Must stay green, unedited** (the design's list, minus P2's one-token fixture edits): every `TestApprove*`,
`TestPasskey*`, `TestApprovePage*`, `TestSurfaceApprove*`, `registry_publish_test.go`, `episode_finish*_test.go`,
`decisionpacket*_test.go`, `TestFrozenV1RouteTableMatchesMux` (13), `TestDaemonDependencyAllowlist`,
`TestSurfaceAppRoutes`, `TestApprovalRoutesRegistered`, `TestDecideApprovalCallerCensus`, sprint 1's RD2–RD8, RD9,
RD9d. The only edited existing tests are listed in §5.

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **Never `nohup`, never background a test leg.** Every acceptance command runs in the foreground of one bash
  shell, `tee`'d to `$HOME/.ailang/state/world-iter390/exec-logs/`. **Every Go Bash call passes `timeout: 600000`,
  foreground, no Monitor** (iteration 389: an executor hit the Bash tool's 120 s default, backgrounded, ended its
  turn waiting, and the shell was lost). The longest single legs are `go test ./host/daemon/ -count=1` (M5b) and
  `go test ./cmd/ailang-worldd/ -count=1` (M5c), each ≈ 4–7 min under load: run each as its own Bash call.
- **Runs.** Run 1 = M4a, M4b, M4c; Run 2 = M4d, M4e, M4f, M4g; Run 3 = M5a, M5b, M5c. Each run reads §0, §2 and its
  own §3 sections, then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an uncommitted
  tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter390-a4-s2 && git diff --numstat 8338eebb173290168f07e6b426d461c0122f67a8 HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M4a **55**, M4b **450**, M4c **559**, M4d **777**, M4e **1142**, M4f **1258**, M4g
  **1380**, M5a **1658**, M5b **1804**, M5c **1983**. **Any other value means a patch was not applied as written:
  stop and report.**
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M4a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter390-a4-s2
    D=$HOME/.ailang/state/world-iter390/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a4-sprint2-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a4-s2.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M4a` … `M5c`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any other
  status. The controller rebuilds one commit per sub-milestone (`commit_ms.sh` with `WT` and the snapshot root
  pointed here; **commit the plan BEFORE running it — it deletes untracked files**).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- FAIL`, or a `--- SKIP` in a §3 `-v` run (`TestApproveDecideSignWithSSHKeygen`
  SKIPs only without `ssh-keygen`; on the rig it is at `/usr/bin/ssh-keygen`, so a SKIP is a stop); any need to
  touch a file outside §5; any reference to `~/.ssh`, `~/.ailang/world`, a real `/dev/tty`, a real key or a real
  store. Stop and report; do not work around. A red in an untouched package that is on §6's known list is re-run
  alone once.
- **Timing (under the planner's load, 13–53).** Each `go test` of one package compiles in ~20–60 s; the §3 `-v`
  runs take seconds once compiled; M4e's four package runs ≈ 6 min; M5b's and M5c's whole-package runs ≈ 4–7 min
  each (`TestDecideDiscoveryRefusesAnOverflow` stores 2,001 packets twice in ~3 s).
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does not
  word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a file; never
  pipe a gate's only copy through `tail`. macOS has no `timeout` command.

---

## 8. Risks

- **R1 — No passkey decides a packet yet** (S2-D1): the webauthn arm, `--write-message`/`--signature` on a request
  ref and the counter (KD2) are sprint 3's. An operator decides by the typed `y` or `--sign-with`; a key on another
  machine has no out-of-band path until sprint 3.
- **R2 — The CLI's daemon tests use a fake daemon** (S2-D7): the wire is built from the daemon's own types on both
  sides and the daemon side is tested through `d.Handler()`, but no test runs the CLI against a real daemon with a
  committed packet (it needs `host/daemon`'s test-only forging helpers). The attended run (design §10 steps 2–3)
  is the first CLI ↔ daemon run.
- **R3 — `approvalStoreError`'s 503 `Restart` text names "surface approvals"** for a store below schema 8 on the
  decide route (it is A3a's shared mapping, not edited). The instruction (restart `serve` once) is right; the noun
  is A3a's. `serve` upgrades the store at start, so the arm is reachable only if the writer was busy at start.
- **R4 — The writer is held across one world write** (design §2.7 allows it): a terminal `session`/`approver`
  verb during a decide gets its existing busy refusal.
- **R5 — RK3b, RK4, RK6 (design) are unchanged.**
- **R6 — H-LOCK's test is probabilistic as a mutant kill** (eight concurrent POSTs): with the lock it is
  deterministic (one 200, seven 409); without it the interleaving decides whether a 503/500 appears (§4 records
  what the run saw).

---

## 9. Acceptance criteria (the judge scores against these)

1. **The typed CLI** (M4b, M4c): `ailang-worldd approve <sha256:…64 hex> [--deny] [--daemon h:p]` finds the daemon
   (`--daemon` or the service address), refuses one that does not serve `--db` and its credentials store, GETs the
   subject, refuses before any prompt when a byte field does not hash to its ref, the packet does not carry the
   ref, either object does not decode strictly, or the daemon's display refusal is set; shows every field of both
   objects with the verb in capitals and the deadline; on `y` at `/dev/tty` writes the proof
   `operator-proofs/decision-<hex Q>.json` with phase `decide-<verb>` and POSTs `typed: true` with the window's
   max as expiry and the daemon's loopback `Origin`; `n`, no terminal, `--deny`/`--daemon` with an `mr-` id,
   `--await` or a passkey/`--signature`/`--write-message` flag with a request ref → refused, nothing posted.
2. **The J gaps** (M4a): J07, J31/J32, J33 each have a test that kills the iteration-389 survivor (§4).
3. **The decision law** (M4d): `DecideBySignature` runs P2–P5 then the SSHSIG over the 91-byte message under
   `world-decision@sunholo`; a mint signature never decides and a decision signature never mints; the verb and the
   expiry are signed; `Binds` checks verb, packet, request, expiry; `DecidedBy` is `ssh:<name>`; `host/approval`
   stays stdlib-only and pure; the mint law is untouched.
4. **Schema 8** (M4e, M4f): a fresh store is version 8 with `key_decisions` byte-equal to the design's DDL (both
   `sha256:*` CHECKs, KD5/KD6); a v7 store is upgraded only by `serve` (`UpgradeCredentials`), and before that every
   key-decision method refuses naming the restart; `RecordKeyDecision` re-checks the approver in its transaction
   on the exact-match path too (KD3/KD9), keeps an exact retry's first record, replaces an unlanded row that differs
   in artifact, packet, verb or expiry, refuses a landed one; `LandKeyDecision` lands once; every method requires a
   deadline.
5. **The judge** (M4g): `JudgeDecisionSSH` refuses with `judge.go`'s step names (parse, store, unknown-key with the
   fingerprint, law) and returns a decision bound to its subject; `judge.go` is untouched.
6. **The key arm** (M5a, M5b): the POST admits exactly one of `typed` and `sshsig`; an enrolled key's SSHSIG decides
   approve or deny as `ssh:<name>`, derived, with one landed `key_decisions` row naming the decision; the law runs
   first; every judge refusal is 422 with 0 decisions and 0 rows; a v4/v7 credentials store is 503 `Restart`; a
   busy writer is 503 `Busy` with `Retry-After: 1` and nothing written; a failed world write leaves one unlanded
   row and 0 decisions, and the identical retry lands exactly one decision and the row (RD13 a); a replay is 409;
   publish- and finish-shaped requests by key are 403 with no row (R4e) and L2 alone still refuses a key-decided
   publish (`LAW-D5-L2`); a hidden second packet and an overflow refuse by key with no row; eight concurrent POSTs
   decide once; the daemon has exactly two credentials-writer sites.
7. **`--sign-with`** (M5c): signs the 91-byte decision message under `world-decision@sunholo` (the real `ssh-keygen`
   sibling verifies with the daemon's check and is not a mint signature), bounded by the signed expiry; a signer
   error posts nothing; no terminal, no proof.
8. **No `/v1` route, no new route, no `.ail`, no CSP change, no CORS header, no new dependency**; `decidedBy` is
   never read from a caller; `TestFrozenV1RouteTableMatchesMux` 13; `TestDecideApprovalCallerCensus` unedited.
9. Each §4 mutant, applied alone, makes its named test **and** its package suite fail (inventories excluded), except
   the `E`/`I`/`L` rows with their reasons; the ledger has a line per run.
10. Every §6 gate passes (modulo §6's re-run rule), including every inventory by name, the `-race` leg and
    `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after the commit; `git status --short` lists only §5.
11. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000 (1,983).
12. **The fallback cut is recorded with measured numbers** (S2-D1, Appendix C).
13. **What an operator can do after this sprint**: approve or deny a pending, uniquely-packeted `Human.Approve`
    request at a terminal by a typed `y` (`approve <request-ref> [--deny]`), or by an enrolled SSH key
    (`approve <request-ref> [--deny] --sign-with <key>`), with `decidedBy` `operator:terminal` / `ssh:<name>`.
    What the agent can do: nothing new (R0 refuses its session; it has no enrolled key and cannot write the proof
    directory; XG1 pins the callers). **The attended exit run (design §10) still waits for**: sprint 3 (the passkey
    arm and its counter, the decision card and *Needs you*, `--write-message`/`--signature` on a request ref,
    QUICKSTART §18); steps 2 (typed) and 3 (SSH key) of §10 are runnable after this sprint.
14. No real key, store, `/dev/tty`, `~/.ssh` access, mint into a real store or publish anywhere in the sprint
    (D-WORLD-109); `ssh-keygen -Y sign` only on a key made in a test's temp dir.

---

## Appendix C — after this sprint (sprint 3's head and handoff)

- **Moved to sprint 3 (S2-D1), measured where built, priced where not:**
  - **The webauthn arm** (design price ≈ 500): `DecideByAssertion` (DM7, the A2 verifier's defer-W15 form),
    `JudgeDecisionAssertion`, the POST's `webauthn` grammar member and arm (W15 + `usePasskey` in
    `RecordKeyDecision`'s transaction on the non-exact path; the exact-match path skips only the counter), KD2,
    the passkey halves of KD7, KD8, RD13 ((b) counter N, (c) counter 0) and RD14 (across a restart),
    `DecisionSubject.Passkeys` (sprint-1 P7), CD8 `--webauthn-link`, CD9 `--webauthn-assertion` on a request ref.
    `RecordKeyDecision`'s signature carries the passkey use through `ApproverUse.Passkey` (already a field).
  - **`--write-message` / `--signature --expiry` on a request ref** (CD6, CD7): built and green in this planner's
    scratch, **63 changed lines measured** (that variant measured 2,019 total against 1,956 without it), cut for the
    cap; it is tagged `CD67-CUT` (`afbc1291`) in `plan-scratch/work` for sprint 3's planner (its sign seam predates
    `signDecision`: rebase it on M5c).
  - Sprint 3's own: the listing (`requests.json` `decisions`), the card, *Needs you*, P1 (Chrome's 91-byte
    message), QUICKSTART §18, design Q2 (`approve --list` decisions).
  - Sprint 3's planner must price its own fallback: ≈ 1,500 (design) + ≈ 500 (webauthn arm) + 63 ≈ 2,060.
- **Named survivors to close when the code they guard changes**: `D-RACE` (unreachable through the POST),
  `D-LAND-ERR`, `C-PROOF-ERR`, `C-SIGN-NIL` (§4 reasons).
- **For the attended exit run** (design §10, Mark; NOT run by the loop, D-WORLD-109): steps 2 (typed) and 3
  (`--sign-with` with an enrolled key) are runnable after this sprint, with `--daemon` if the service is not
  installed; step 4 (passkey) and step 5's saved body (`--write-message`/`--signature`) wait for sprint 3.
- **D-WORLD-115** stays open (default A applied); S2-D9 names what B would change.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M4a — the sprint-1 judge's test gaps (J07, J31/J32, J33)

````diff M4a
diff --git a/host/broker/decision_law_test.go b/host/broker/decision_law_test.go
index fb53b00..f186275 100644
--- a/host/broker/decision_law_test.go
+++ b/host/broker/decision_law_test.go
@@ -199,8 +199,9 @@ func TestReservedApprovalScope(t *testing.T) {
 	for _, tc := range []struct{ scope, owner string }{
 		{PublishApprovalScopeFor(id, PublishHashes{TarballSHA256: "t", ContentHash: "c", InterfaceHash: "i"}, 9), "publish"},
 		{"pkg#effect=Other&manifest=m&tarball=t&content=c&interface=i&expires=9", "publish"}, // the parse arm alone
-		{"registry:x", "publish"},                       // the prefix arm alone
-		{"worktree#effect=Registry.Publish", "publish"}, // the effect arm alone
+		{"registry:x", "publish"},                           // the prefix arm alone
+		{"worktree#effect=Registry.Publish", "publish"},     // the effect arm alone
+		{"worktree#effect=Registry.Publish&x=1", "publish"}, // the effect arm, not at the end (sprint-1 judge J07)
 		{FormatFinishApprovalScope(FinishApprovalScope{Repository: "o/r", Remote: "https://github.com/o/r.git", Head: "world/e", SHA: strings.Repeat("a", 40),
 			Base: "dev", TitleSHA256: strings.Repeat("b", 64), BodySHA256: strings.Repeat("c", 64), ExpiresAt: 9}), "finish"},
 		{"github-pr:x", "finish"}, // the finish prefix arm alone
diff --git a/host/daemon/surface_decide_jgaps_test.go b/host/daemon/surface_decide_jgaps_test.go
new file mode 100644
index 0000000..ee1cf62
--- /dev/null
+++ b/host/daemon/surface_decide_jgaps_test.go
@@ -0,0 +1,50 @@
+package daemon
+
+// Row 201 A4 sprint 2 M4a: the sprint-1 judge's test gaps J31-J33 (iteration 389), end to end.
+
+import (
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+)
+
+// J31/J32: D1's grammar is anchored and lower-case on both routes: 65 hex, 63 hex and upper-case
+// hex are each 400 D1 (a 65-hex ref that slipped past D1 would make hashref.MustParse panic), 0
+// decisions; the 64-hex ref is the control.
+func TestDecideRefGrammarIsAnchored(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("anchored", broker.WorkspaceScope, "agent")
+	proof := r.proof(q, "decide-approve")
+	_, head := r.decisions()
+	for name, ref := range map[string]string{"65 hex": q.String() + "0", "63 hex": q.String()[:70], "upper hex": "sha256:" + strings.ToUpper(q.Digest())} {
+		got := requireRefusal(t, "POST "+name, r.post(ref, decideBody("approve", p, q, decNow+50), proof), 400, "BadRequest")
+		if !strings.Contains(got.Reason, "is not a request ref") {
+			t.Errorf("POST %s: %q", name, got.Reason)
+		}
+		requireRefusal(t, "GET "+name, r.serve("GET", "/surface/decisions/"+ref, apLoopback, nil), 400, "BadRequest")
+		r.requireNothingDecided(name, q, head)
+	}
+	if rec := r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil); rec.Code != 200 {
+		t.Fatalf("control GET: %d %s", rec.Code, rec.Body.String())
+	}
+}
+
+// J33 (and sprint 1's D4a-undecodable row, unit-only until now): a stored packet that does not
+// decode makes the discovery incomplete, so GET answers D4a's 409 with no subject and the POST
+// refuses with 0 decisions.
+func TestDecisionSubjectRefusesAnUndecodablePacket(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("undecodable", broker.WorkspaceScope, "agent")
+	if err := r.d.store.PutObject(boundedTestContext(t), forgeObject(coordinator.DecisionPacketV1, []byte(`{"not":"a packet"}`))); err != nil {
+		t.Fatal(err)
+	}
+	_, head := r.decisions()
+	got := requireRefusal(t, "GET", r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil), 409, "Conflict")
+	if !strings.Contains(got.Reason, "cannot be counted completely") {
+		t.Errorf("GET: %q", got.Reason)
+	}
+	requireRefusal(t, "POST", r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 409, "Conflict")
+	r.requireNothingDecided("undecodable", q, head)
+}
````

### Patch M4b — the typed decision mode (CD2, CD3)

````diff M4b
diff --git a/cmd/ailang-worldd/approve.go b/cmd/ailang-worldd/approve.go
index f3ff92f..debe434 100644
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -29,6 +29,7 @@ const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--
        ailang-worldd approve <request-id> --webauthn-assertion <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --await [--await-for <duration>] [--db <path>] [--credentials-db <path>]
        ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-ref> [--deny] [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
 
 ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
 201), after a y/N typed at /dev/tty, or, with --signature or --sign-with,
@@ -91,6 +92,16 @@ after it is made; at most 16 are pending per credentials store.
                             (default ~/.ailang/world/world.db)
   --credentials-db <path>   default credentials.db beside --db
 The global --addr is refused: the daemon's address is part of the request.
+
+approve <request-ref> (sha256: and 64 hex, row 201 A4) decides a pending
+Human.Approve request an agent raised in an episode, whose decision packet the
+daemon committed: approve, or with --deny deny. It reads the request and its
+packet from the daemon serving --db (--daemon <host:port>, default the
+installed service's address), checks that both hash to their refs, shows every
+field, asks y at /dev/tty, and has the daemon record the decision (decided by
+operator:terminal). A request is decided once. --await is refused (the daemon
+decides at once); publish and finish approvals stay with world-publish
+approve and episode finish.
 `
 
 // approveAdmission is approve's admission of a stored row, before any terminal: the decoded
@@ -155,6 +166,8 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 	assertion := fs.String("webauthn-assertion", "", "approve with an enrolled passkey's assertion file (no terminal)")
 	await := fs.Bool("await", false, "wait for an approval recorded by the daemon's surface, judge it, and mint (no terminal)")
 	awaitFor := fs.Duration("await-for", 0, "with --await: stop waiting after this long (default: the request's expiry)")
+	deny := fs.Bool("deny", false, "with a request ref: deny it instead of approving it")
+	daemonAddr := fs.String("daemon", "", "with a request ref: the daemon serving --db (default: the installed service's address)")
 	positional, err := parseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approveHelp)
@@ -192,9 +205,22 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
+	if (*deny || *daemonAddr != "") && (*list || !decisionRefPattern.MatchString(positional[0])) { // CD3: --deny and --daemon decide a request ref
+		return refuse("--deny and --daemon take a Human.Approve request ref (sha256: and 64 hex); a mint request has no deny: it expires")
+	}
 	id := ""
 	if !*list {
 		id = positional[0]
+		if decisionRefPattern.MatchString(id) { // row 201 A4: a request ref selects the decision mode (design §2.10)
+			if modes != 0 { // CD3: --await, and the key flags until they decide a request ref
+				return refuse("a request ref is decided at once by the daemon: there is nothing to await, and no key flag decides it yet")
+			}
+			verb := "approve"
+			if *deny {
+				verb = "deny"
+			}
+			return runApproveDecide(decideOptions{ref: id, verb: verb, db: *db, credsDB: *credsDB, daemon: *daemonAddr}, stdout, stderr, env, refuse)
+		}
 		if !store.MintRequestIDPattern.MatchString(id) {
 			return refuse("%q is not a mint request id (mr- and 16 lowercase hex)", id)
 		}
diff --git a/cmd/ailang-worldd/approve_decide.go b/cmd/ailang-worldd/approve_decide.go
new file mode 100644
index 0000000..91e38b0
--- /dev/null
+++ b/cmd/ailang-worldd/approve_decide.go
@@ -0,0 +1,184 @@
+package main
+
+// Row 201 A4 sprint 2 (design_docs/planned/w-world-approvals-by-key-a4.md §2.10): `approve
+// <request-ref> [--deny]` decides a pending Human.Approve request whose decision packet the
+// coordinator committed. The CLI never opens the world store: it reads the subject from the
+// daemon serving --db (GET /surface/decisions/{ref}), checks both byte fields against their refs
+// and the packet's requestRef here, shows every field, and POSTs the decision to
+// /surface/decide/{ref}, where the daemon runs the decision law under humanMu and writes. The typed
+// arm writes the row-236 operator proof beside the served credentials store only after a y typed
+// at /dev/tty.
+
+import (
+	"bytes"
+	"context"
+	"encoding/base64"
+	"encoding/json"
+	"fmt"
+	"io"
+	"net/http"
+	"regexp"
+	"strings"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// decisionRefPattern is a Human.Approve request ref: approve's decision mode (design §2.10).
+var decisionRefPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
+
+// approveDecideBudget bounds the decision mode's reads before the y/N, and its POST after it.
+const approveDecideBudget = 10 * time.Second
+
+// decideOptions are approve's flags in the decision mode.
+type decideOptions struct {
+	ref, verb, db, credsDB, daemon string
+}
+
+// decisionSubject is a subject the CLI checked: its request and packet, decoded from bytes that
+// hash to their refs.
+type decisionSubject struct {
+	surfacepage.DecisionSubject
+	request broker.ApprovalRequest
+	packet  coordinator.DecisionPacket
+}
+
+// checkDecisionSubject decodes the daemon's GET answer for ref and checks it here (design §2.8:
+// the daemon sends bytes, never a message or a hash): the request bytes hash to ref, the packet
+// bytes hash to the packet ref, both decode strictly, and the packet carries ref. It returns the
+// refusal text, or "".
+func checkDecisionSubject(ref string, body []byte) (decisionSubject, string) {
+	var s decisionSubject
+	dec := json.NewDecoder(bytes.NewReader(body))
+	dec.DisallowUnknownFields()
+	if err := dec.Decode(&s.DecisionSubject); err != nil {
+		return s, fmt.Sprintf("the daemon's answer for request %s is not a decision subject (%v); nothing was decided", ref, err)
+	}
+	mismatch := fmt.Sprintf("the daemon's bytes for request %s do not match their refs; refusing before any prompt (nothing was decided)", ref)
+	reqBytes, err1 := base64.RawURLEncoding.DecodeString(s.RequestBytes)
+	pktBytes, err2 := base64.RawURLEncoding.DecodeString(s.PacketBytes)
+	if err1 != nil || err2 != nil || sha256Ref(reqBytes) != ref || sha256Ref(pktBytes) != s.Packet {
+		return s, mismatch
+	}
+	var err error
+	if s.request, err = broker.DecodeApprovalRequest(reqBytes); err != nil {
+		return s, mismatch
+	}
+	if s.packet, err = coordinator.DecodeDecisionPacket(pktBytes); err != nil || s.packet.RequestRef.String() != ref {
+		return s, mismatch
+	}
+	return s, ""
+}
+
+// sha256Ref is the sha256: ref of b.
+func sha256Ref(b []byte) string { return "sha256:" + finishSHA(b) }
+
+// renderDecision is the typed y/N: every field of both signed objects (design §2.8's field
+// table), the entry and the proposal's page, then the verb in capitals and the deadline.
+func renderDecision(ref, verb string, s decisionSubject) string {
+	r, p := s.request, s.packet
+	return fmt.Sprintf("request %s\n  effect               %s\n  scope                %s\n  cost                 %d\n  requester            %s\n  now                  %d\n"+
+		"packet %s (log entry %d)\n  packetHash           %s\n  proposalHash         %s (/workbench?object=%s)\n  requestRef           %s\n"+
+		"  createdAt            %d\n  deadlineAt           %d (%s)\n  escalationsRemaining %d\n  policy               %s\n"+
+		"%s this request? A request is decided once; the decision is final. [y/N] ",
+		ref, r.Effect, r.Scope, r.Cost, r.Requester, r.Now, s.Packet, s.Entry, p.PacketHash, p.ProposalHash, p.ProposalHash, p.RequestRef,
+		p.CreatedAt, p.DeadlineAt, approverTime(p.DeadlineAt), p.EscalationsRemaining, p.Policy, strings.ToUpper(verb))
+}
+
+// refusalReason is the reason of a daemon refusal body, or the body itself.
+func refusalReason(body []byte) string {
+	var r surfacepage.ApproveRefusal
+	if json.Unmarshal(body, &r) == nil && r.Reason != "" {
+		return r.Reason
+	}
+	return strings.TrimSpace(string(body))
+}
+
+// runApproveDecide is approve's decision mode (design §2.10). The daemon is --daemon, or the
+// home service's address; it must serve --db and the credentials store beside it (the proof is
+// written there). Every refusal before the POST decides nothing.
+func runApproveDecide(o decideOptions, stdout, stderr io.Writer, env sessionEnv, refuse func(string, ...any) int) int {
+	credPath, err := credentialsFor(o.db, o.credsDB)
+	if err != nil {
+		return refuse("%v", err)
+	}
+	if o.daemon == "" {
+		o.daemon = defaultServiceAddr()
+	}
+	base := clientBaseURL(o.daemon)
+	hostPort := strings.TrimSuffix(strings.TrimPrefix(base, "http://"), "/")
+	h, err := fetchDaemonHealth(hostPort)
+	if err != nil {
+		return refuse("%v; nothing was decided; fix: start the daemon serving %s, or name it with --daemon <host:port>", err, o.db)
+	}
+	if same, _ := store.SameDatabase(h.DBPath, o.db); !same {
+		return refuse("the daemon on %s serves store %s, not %s; nothing was decided; fix: --daemon <host:port> of the daemon serving it", hostPort, h.DBPath, o.db)
+	}
+	if same, _ := store.SameDatabase(h.CredentialsDB, credPath); !same {
+		return refuse("the daemon on %s serves credentials store %s, not %s; nothing was decided (the operator proof is written beside the served one)", hostPort, h.CredentialsDB, credPath)
+	}
+	c := newClient(base)
+	ctx, cancel := context.WithTimeout(context.Background(), approveDecideBudget)
+	defer cancel()
+	status, body, err := c.do(ctx, http.MethodGet, "/surface/decisions/"+o.ref, nil)
+	if err != nil {
+		return refuse("%v; nothing was decided", err)
+	}
+	if status != http.StatusOK {
+		return refuse("the daemon will not show request %s (HTTP %d): %s", o.ref, status, refusalReason(body))
+	}
+	s, refusal := checkDecisionSubject(o.ref, body)
+	if refusal != "" {
+		return refuse("%s", refusal)
+	}
+	if s.Refusal != "" {
+		return refuse("request %s cannot be decided now: %s", o.ref, s.Refusal)
+	}
+	const nothing = "nothing was decided; decide with an enrolled key instead (--sign-with <key> or --signature <file>)"
+	term, code := openApproverTerminal(env, stderr, "ailang-worldd approve: ", nothing)
+	if code != exitOK {
+		return code
+	}
+	code = confirmOnTerminal(term, env, stderr, "ailang-worldd approve: ", renderDecision(o.ref, o.verb, s), "nothing was decided")
+	_ = term.Close()
+	if code != exitOK {
+		return code
+	}
+	proof, err := writeOperatorProof(credPath, "decision-"+strings.TrimPrefix(o.ref, "sha256:"), "decide-"+o.verb)
+	if err != nil {
+		return refuse("cannot write the operator proof beside %s: %v; nothing was decided", credPath, err)
+	}
+	return postDecision(c, base, o.ref, surfacepage.DecideBody{V: 1, Verb: o.verb, Packet: strings.TrimPrefix(s.Packet, "sha256:"),
+		Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max, Typed: true}, proof, stdout, refuse)
+}
+
+// postDecision POSTs body to /surface/decide/{ref} with the daemon's own loopback Origin (in O by
+// construction; design §2.8) and, for the typed arm, the operator proof; it prints the decision.
+func postDecision(c *client, origin, ref string, body surfacepage.DecideBody, proof string, stdout io.Writer, refuse func(string, ...any) int) int {
+	data, err := json.Marshal(body)
+	if err != nil {
+		return refuse("%v", err)
+	}
+	hdr := map[string]string{"Origin": origin}
+	if proof != "" {
+		hdr["X-World-Operator-Proof"] = proof
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), approveDecideBudget)
+	defer cancel()
+	status, _, resp, err := c.doRequest(ctx, http.MethodPost, "/surface/decide/"+ref, bytes.NewReader(data), hdr)
+	if err != nil {
+		return refuse("%v; run approve %s again to see whether it was decided", err, ref)
+	}
+	if status != http.StatusOK {
+		return refuse("the daemon refused to decide request %s (HTTP %d): %s", ref, status, refusalReason(resp))
+	}
+	var made surfacepage.DecisionMade
+	if err := json.Unmarshal(resp, &made); err != nil || made.Request != ref {
+		return refuse("the daemon's answer for request %s is not a decision (%s)", ref, strings.TrimSpace(string(resp)))
+	}
+	fmt.Fprintf(stdout, "decided %s: request %s\ndecision %s by %s\n", made.Decided, made.Request, made.Decision, made.DecidedBy)
+	return exitOK
+}
diff --git a/cmd/ailang-worldd/approve_decide_rig_test.go b/cmd/ailang-worldd/approve_decide_rig_test.go
new file mode 100644
index 0000000..4fbf220
--- /dev/null
+++ b/cmd/ailang-worldd/approve_decide_rig_test.go
@@ -0,0 +1,179 @@
+package main
+
+// Row 201 A4 sprint 2 M4b (design §2.10, §6 CD2, CD3): approve <request-ref> [--deny], typed,
+// against a fake daemon (httptest on the loopback) whose subject is built with the daemon's own
+// encoders and JSON types (coordinator.EncodeDecisionPacket, surfacepage.DecisionSubject); the
+// real daemon side of the same wire is host/daemon's RD2-RD9 (D-WORLD-109: temp dirs only, no
+// real store, no /dev/tty).
+
+import (
+	"bytes"
+	"encoding/base64"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"net/http"
+	"net/http/httptest"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// decideTerm is a terminal answering answer and keeping what was written to it.
+type decideTerm struct {
+	r   io.Reader
+	out bytes.Buffer
+}
+
+func (d *decideTerm) Read(p []byte) (int, error)  { return d.r.Read(p) }
+func (d *decideTerm) Write(p []byte) (int, error) { return d.out.Write(p) }
+func (d *decideTerm) Close() error                { return nil }
+
+// fakeDecider is a daemon serving /v1/health, GET /surface/decisions/{ref} and POST
+// /surface/decide/{ref} for one request; tests tamper with subject, health and the answers.
+type fakeDecider struct {
+	srv        *httptest.Server
+	dir, db    string
+	creds      string
+	ref        string
+	subject    surfacepage.DecisionSubject
+	getStatus  int
+	getBody    string
+	health     daemon.HealthResponse
+	gets       int
+	posts      []*http.Request
+	postBodies []surfacepage.DecideBody
+	postStatus int
+	postBody   string
+}
+
+// decisionObjects are a Human.Approve request of scope and a packet carrying it, as stored bytes.
+func decisionObjects(scope string) (req, pkt []byte, q hashref.HashRef) {
+	req = []byte(fmt.Sprintf(`{"effect":"Human.Approve","scope":%q,"cost":1,"requester":"agent","now":1700}`, scope))
+	q = hashref.SumSHA256(req)
+	pkt = coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan " + scope)), RequestRef: q,
+		CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"})
+	return req, pkt, q
+}
+
+func newFakeDecider(t *testing.T) *fakeDecider {
+	t.Helper()
+	f := &fakeDecider{dir: t.TempDir(), getStatus: 200, postStatus: 200}
+	f.db, f.creds = filepath.Join(f.dir, "world.db"), filepath.Join(f.dir, "credentials.db")
+	for _, p := range []string{f.db, f.creds} {
+		if err := os.WriteFile(p, nil, 0o600); err != nil {
+			t.Fatal(err)
+		}
+	}
+	req, pkt, q := decisionObjects("worktree")
+	f.ref = q.String()
+	f.subject = surfacepage.DecisionSubject{V: 1, Request: f.ref, RequestBytes: base64.RawURLEncoding.EncodeToString(req),
+		Packet: hashref.SumSHA256(pkt).String(), PacketBytes: base64.RawURLEncoding.EncodeToString(pkt), Entry: 7,
+		Proposal: hashref.SumSHA256([]byte("plan worktree")).String(), Now: 1710, Window: surfacepage.DecisionWindow{Min: 1711, Max: 1760},
+		SSH: surfacepage.DecisionSSH{Namespace: "world-decision@sunholo"}}
+	f.health = daemon.HealthResponse{DBPath: f.db, CredentialsDB: f.creds}
+	mux := http.NewServeMux()
+	mux.HandleFunc("GET /v1/health", func(w http.ResponseWriter, _ *http.Request) { _ = json.NewEncoder(w).Encode(f.health) })
+	mux.HandleFunc("GET /surface/decisions/{ref}", func(w http.ResponseWriter, r *http.Request) {
+		f.gets++
+		w.WriteHeader(f.getStatus)
+		if f.getStatus != 200 {
+			_, _ = w.Write([]byte(`{"v":1,"error":"Conflict","reason":"the packet's deadline passed"}`))
+			return
+		}
+		if f.getBody != "" {
+			_, _ = w.Write([]byte(f.getBody))
+			return
+		}
+		_ = json.NewEncoder(w).Encode(f.subject)
+	})
+	mux.HandleFunc("POST /surface/decide/{ref}", func(w http.ResponseWriter, r *http.Request) {
+		var b surfacepage.DecideBody
+		raw, _ := io.ReadAll(r.Body)
+		_ = json.Unmarshal(raw, &b)
+		f.posts, f.postBodies = append(f.posts, r), append(f.postBodies, b)
+		w.WriteHeader(f.postStatus)
+		if f.postBody != "" {
+			_, _ = w.Write([]byte(f.postBody))
+			return
+		}
+		_ = json.NewEncoder(w).Encode(surfacepage.DecisionMade{V: 1, Decided: b.Verb, Request: "sha256:" + b.Request, Decision: "sha256:" + strings.Repeat("d", 64), DecidedBy: "operator:terminal"})
+	})
+	f.srv = httptest.NewServer(mux)
+	t.Cleanup(f.srv.Close)
+	return f
+}
+
+func (f *fakeDecider) addr() string { return strings.TrimPrefix(f.srv.URL, "http://") }
+
+// run is approve <ref> --db --daemon extra... with term (nil: no terminal); it returns the code,
+// stdout and stderr, and whether a terminal was opened.
+func (f *fakeDecider) run(t *testing.T, term io.ReadWriteCloser, extra ...string) (int, string, string, bool) {
+	t.Helper()
+	opened := false
+	env := approveEnv(term, &opened, atTime(1710))
+	if term == nil {
+		env.openTerminal = func() (io.ReadWriteCloser, error) { opened = true; return nil, errors.New("no tty") }
+	}
+	var stdout, stderr bytes.Buffer
+	code := runApprove(append([]string{f.ref, "--db", f.db, "--daemon", f.addr()}, extra...), &stdout, &stderr, env)
+	return code, stdout.String(), stderr.String(), opened
+}
+
+// proofFile is the operator proof the CLI wrote for the request, or ok false.
+func (f *fakeDecider) proofFile(t *testing.T) (p struct {
+	SHA256    string `json:"sha256"`
+	Phase     string `json:"phase"`
+	ExpiresAt int64  `json:"expires_at"`
+}, ok bool) {
+	t.Helper()
+	data, err := os.ReadFile(filepath.Join(f.dir, "operator-proofs", "decision-"+strings.TrimPrefix(f.ref, "sha256:")+".json"))
+	if err != nil {
+		return p, false
+	}
+	if err := json.Unmarshal(data, &p); err != nil {
+		t.Fatal(err)
+	}
+	return p, true
+}
+
+// CD2: no terminal → refused naming the key flags, nothing posted, no proof written.
+func TestApproveDecideNoTerminal(t *testing.T) {
+	f := newFakeDecider(t)
+	code, _, stderr, opened := f.run(t, nil)
+	if _, wrote := f.proofFile(t); code != exitUsage || !opened || len(f.posts) != 0 || wrote || !strings.Contains(stderr, "--sign-with <key>") {
+		t.Fatalf("no terminal: rc %d, opened %v, %d POSTs, proof %v\nstderr %s", code, opened, len(f.posts), wrote, stderr)
+	}
+}
+
+// CD3: --deny and --daemon need a request ref (an mr- id, --list); --await, a passkey flag and a
+// key flag are refused with a request ref; two modes are F1. Each before any daemon call.
+func TestApproveDecideFlagRules(t *testing.T) {
+	f := newFakeDecider(t)
+	for _, tc := range []struct {
+		args []string
+		want string
+	}{
+		{[]string{"mr-0123456789abcdef", "--deny"}, "a mint request has no deny"},
+		{[]string{f.ref + "0"}, "is not a mint request id"}, // the decision mode's ref is anchored
+		{[]string{"mr-0123456789abcdef", "--daemon", f.addr()}, "a mint request has no deny"},
+		{[]string{"--list", "--deny"}, "a mint request has no deny"},
+		{[]string{f.ref, "--await"}, "there is nothing to await"},
+		{[]string{f.ref, "--webauthn-link"}, "there is nothing to await"},
+		{[]string{f.ref, "--signature", "x", "--sign-with", "y"}, "usage: ailang-worldd approve"},
+	} {
+		var stdout, stderr bytes.Buffer
+		opened := false
+		code := runApprove(append(tc.args, "--db", f.db), &stdout, &stderr, approveEnv(nil, &opened, atTime(1710)))
+		if code != exitUsage || opened || f.gets != 0 || !strings.Contains(stderr.String(), tc.want) {
+			t.Errorf("%v: rc %d, opened %v, %d GETs\nstderr %s", tc.args, code, opened, f.gets, stderr.String())
+		}
+	}
+}
diff --git a/cmd/ailang-worldd/help_test.go b/cmd/ailang-worldd/help_test.go
index acdc65a..e64d356 100644
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -41,7 +41,7 @@ func TestEveryVerbHasHelp(t *testing.T) {
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
 		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
-		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve <request-id> --webauthn-link", "approve <request-id> --webauthn-assertion <file>", "approve <request-id> --await [--await-for <duration>]", "approve --list [--decided]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
+		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve <request-id> --webauthn-link", "approve <request-id> --webauthn-assertion <file>", "approve <request-id> --await [--await-for <duration>]", "approve --list [--decided]", "approve <request-ref> [--deny] [--daemon <host:port>]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
 			t.Errorf("top-level usage lacks %q", want)
diff --git a/cmd/ailang-worldd/main.go b/cmd/ailang-worldd/main.go
index 6189476..6541bfa 100644
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -132,6 +132,7 @@ Usage:
   ailang-worldd approve <request-id> --webauthn-assertion <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --await [--await-for <duration>] [--db <path>] [--credentials-db <path>]
   ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
+  ailang-worldd approve <request-ref> [--deny] [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format webauthn --origin <origin> --registration <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approver list [--db <path>] [--credentials-db <path>]
diff --git a/host/store/context_roots_test.go b/host/store/context_roots_test.go
index 42949b0..1ce516b 100644
--- a/host/store/context_roots_test.go
+++ b/host/store/context_roots_test.go
@@ -28,6 +28,9 @@ var contextRootPins = map[string]int{
 	"cmd/ailang-worldd/approver_passkey.go|runApproverEnrollPasskey|Background": 3,
 	// Row 201 A3 sprint 1 M3: each approve --await poll's read-only reads and judgement (3 s).
 	"cmd/ailang-worldd/approve_await.go|pollSurfaceApproval|Background": 1,
+	// Row 201 A4 sprint 2: approve <request-ref>'s daemon reads before the y/N, and its POST after it (10 s each).
+	"cmd/ailang-worldd/approve_decide.go|runApproveDecide|Background": 1,
+	"cmd/ailang-worldd/approve_decide.go|postDecision|Background":     1,
 
 	"cmd/ailang-worldd/cli.go|get|Background":                                1,
 	"cmd/ailang-worldd/cli.go|execute|Background":                            1,
````

### Patch M4c — the typed decision end to end (CD1, CD4)

````diff M4c
diff --git a/cmd/ailang-worldd/approve_decide_test.go b/cmd/ailang-worldd/approve_decide_test.go
new file mode 100644
index 0000000..ac58f2d
--- /dev/null
+++ b/cmd/ailang-worldd/approve_decide_test.go
@@ -0,0 +1,109 @@
+package main
+
+// Row 201 A4 sprint 2 M4c (design §2.10, §6 CD1, CD4): the typed decision end to end against the
+// fake daemon (approve_decide_rig_test.go), and the subject checks before any prompt.
+
+import (
+	"encoding/base64"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// CD1: y at the terminal → one POST: typed, the verb, both digests, the window's max as expiry,
+// the daemon's loopback Origin, the proof whose file (phase decide-<verb>, named for the request)
+// hashes to it; every field was shown with the verb in capitals; --deny denies; n posts nothing.
+func TestApproveDecideTyped(t *testing.T) {
+	for _, tc := range []struct{ verb, answer string }{{"approve", "y\n"}, {"deny", "y\n"}, {"approve", "n\n"}} {
+		f := newFakeDecider(t)
+		term := &decideTerm{r: strings.NewReader(tc.answer)}
+		var extra []string
+		if tc.verb == "deny" {
+			extra = []string{"--deny"}
+		}
+		code, stdout, stderr, _ := f.run(t, term, extra...)
+		prompt := term.out.String()
+		for _, want := range []string{f.ref, "effect               Human.Approve", "scope                worktree", "cost                 1",
+			"requester            agent", "now                  1700", f.subject.Packet + " (log entry 7)", "proposalHash         " + f.subject.Proposal,
+			"requestRef           " + f.ref, "createdAt            1700", "deadlineAt           1760 (1970-01-01T00:29:20Z)", "escalationsRemaining 0",
+			"policy               Cancel", strings.ToUpper(tc.verb) + " this request?"} {
+			if !strings.Contains(prompt, want) {
+				t.Errorf("%s %q: the prompt lacks %q:\n%s", tc.verb, tc.answer, want, prompt)
+			}
+		}
+		proof, wrote := f.proofFile(t)
+		if tc.answer == "n\n" {
+			if code != exitUsage || len(f.posts) != 0 || wrote {
+				t.Fatalf("n: rc %d, %d POSTs, proof written %v", code, len(f.posts), wrote)
+			}
+			continue
+		}
+		if code != exitOK || len(f.posts) != 1 {
+			t.Fatalf("%s: rc %d, %d POSTs\nstderr %s", tc.verb, code, len(f.posts), stderr)
+		}
+		b, r := f.postBodies[0], f.posts[0]
+		want := surfacepage.DecideBody{V: 1, Verb: tc.verb, Packet: strings.TrimPrefix(f.subject.Packet, "sha256:"), Request: strings.TrimPrefix(f.ref, "sha256:"), Expiry: 1760, Typed: true}
+		if b != want || r.Header.Get("Origin") != f.srv.URL || r.Header.Get("Content-Type") != "application/json" {
+			t.Fatalf("%s: POST %+v (Origin %q), want %+v from %s", tc.verb, b, r.Header.Get("Origin"), want, f.srv.URL)
+		}
+		if !wrote || proof.Phase != "decide-"+tc.verb || proof.SHA256 != finishSHA([]byte(r.Header.Get("X-World-Operator-Proof"))) {
+			t.Fatalf("%s: proof %+v (written %v) does not prove the POST's header", tc.verb, proof, wrote)
+		}
+		if !strings.Contains(stdout, "decided "+tc.verb+": request "+f.ref) || !strings.Contains(stdout, "by operator:terminal") {
+			t.Errorf("%s: stdout %q", tc.verb, stdout)
+		}
+	}
+	f := newFakeDecider(t) // a refused POST is reported with the daemon's reason
+	f.postStatus, f.postBody = 409, `{"v":1,"error":"Conflict","reason":"request was already decided deny by operator:terminal"}`
+	if code, _, stderr, _ := f.run(t, &decideTerm{r: strings.NewReader("y\n")}); code != exitUsage || !strings.Contains(stderr, "(HTTP 409): request was already decided deny") {
+		t.Fatalf("refused POST: rc %d, stderr %s", code, stderr)
+	}
+}
+
+// CD4: the CLI checks the daemon's answer before any prompt: request bytes of another request,
+// packet bytes that do not hash to the packet ref, a packet carrying another request, a request
+// that does not decode, an unknown field, a daemon refusal, a display refusal, a daemon serving
+// another store or credentials store → refused, no terminal opened, nothing posted.
+func TestApproveDecideRefusesAMismatchedSubject(t *testing.T) {
+	otherReq, otherPkt, _ := decisionObjects("world")
+	_, _, q := decisionObjects("worktree")
+	samePkt := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("another plan")), RequestRef: q, CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"})
+	garbage := []byte(`{"effect":"Human.Approve","extra":1}`)
+	garbagePkt := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: q, RequestRef: hashref.SumSHA256(garbage), CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"})
+	for _, tc := range []struct {
+		name string
+		edit func(f *fakeDecider)
+		want string
+	}{
+		{"request bytes", func(f *fakeDecider) { f.subject.RequestBytes = base64.RawURLEncoding.EncodeToString(otherReq) }, "do not match their refs"},
+		{"packet bytes", func(f *fakeDecider) { f.subject.PacketBytes = base64.RawURLEncoding.EncodeToString(samePkt) }, "do not match their refs"},
+		{"another request's packet", func(f *fakeDecider) {
+			f.subject.PacketBytes, f.subject.Packet = base64.RawURLEncoding.EncodeToString(otherPkt), hashref.SumSHA256(otherPkt).String()
+		}, "do not match their refs"},
+		{"undecodable request", func(f *fakeDecider) {
+			f.ref = hashref.SumSHA256(garbage).String()
+			f.subject.Request, f.subject.RequestBytes = f.ref, base64.RawURLEncoding.EncodeToString(garbage)
+			f.subject.PacketBytes, f.subject.Packet = base64.RawURLEncoding.EncodeToString(garbagePkt), hashref.SumSHA256(garbagePkt).String()
+		}, "do not match their refs"},
+		{"undecodable packet", func(f *fakeDecider) {
+			f.subject.PacketBytes, f.subject.Packet = base64.RawURLEncoding.EncodeToString(garbage), hashref.SumSHA256(garbage).String()
+		}, "do not match their refs"},
+		{"bad base64", func(f *fakeDecider) { f.subject.RequestBytes = "!" }, "do not match their refs"},
+		{"unknown field", func(f *fakeDecider) { f.getBody = `{"v":1,"message":"AAAA"}` }, "is not a decision subject"},
+		{"daemon refusal", func(f *fakeDecider) { f.getStatus = 409 }, "will not show request"},
+		{"display refusal", func(f *fakeDecider) { f.subject.Refusal = "the packet's deadline passed" }, "cannot be decided now: the packet's deadline passed"},
+		{"another store", func(f *fakeDecider) { f.health.DBPath = filepath.Join(f.dir, "other.db") }, "serves store"},
+		{"another credentials store", func(f *fakeDecider) { f.health.CredentialsDB = filepath.Join(f.dir, "other-credentials.db") }, "serves credentials store"},
+	} {
+		f := newFakeDecider(t)
+		tc.edit(f)
+		code, _, stderr, opened := f.run(t, &decideTerm{r: strings.NewReader("y\n")})
+		if code != exitUsage || opened || len(f.posts) != 0 || !strings.Contains(stderr, tc.want) {
+			t.Errorf("%s: rc %d, opened %v, %d POSTs\nstderr %s", tc.name, code, opened, len(f.posts), stderr)
+		}
+	}
+}
````

### Patch M4d — the decision signature law (DM3–DM6, DM8)

````diff M4d
diff --git a/host/approval/decision_signed.go b/host/approval/decision_signed.go
new file mode 100644
index 0000000..4f6e8fc
--- /dev/null
+++ b/host/approval/decision_signed.go
@@ -0,0 +1,91 @@
+package approval
+
+// Row 201 A4 sprint 2 (design §2.3, §2.5 D11): the pure law of a decision signed by an enrolled SSH
+// key: the mint law's P2-P5 and VerifySSHSIG over the 91-byte message under DecisionNamespace (the
+// mint law is unchanged). The passkey half is sprint 3's (plan S2-D1). No I/O; now is an input.
+
+import "bytes"
+
+// SignedDecision is a verified key decision. Its fields are unexported, so only DecideBySignature
+// can return one with ok set; the zero value binds nothing.
+type SignedDecision struct {
+	ok              bool
+	verb            byte
+	packet, request [32]byte
+	expiry          int64
+	name, format    string
+	key             []byte
+}
+
+// DecisionSignatureInput is what DecideBySignature judges: the decision's subject (verb, packet
+// digest P, request digest Q, expiry), the clock, the armored SSHSIG and the enrolled approver
+// its key names.
+type DecisionSignatureInput struct {
+	Verb            byte
+	Packet, Request [32]byte
+	Expiry, Now     int64
+	Signature       []byte
+	Approver        Approver
+}
+
+// VerifyDecisionSSHSIG parses an armored SSHSIG and verifies it over DecisionMessage(verb, packet,
+// request, expiry) under DecisionNamespace with key. It returns the first error.
+func VerifyDecisionSSHSIG(armored []byte, verb byte, packet, request [32]byte, expiry int64, key SSHPublicKey) error {
+	s, err := ParseSSHSIG(armored)
+	if err != nil {
+		return err
+	}
+	msg, err := DecisionMessage(verb, packet, request, expiry)
+	if err != nil {
+		return err
+	}
+	return VerifySSHSIG(s, DecisionNamespace, msg, key)
+}
+
+// DecideBySignature is the decision signature law: the approver's state first (P2-P5, the mint
+// law's own checks, so a revoked key says "revoked" even when its signature is good), then the
+// signature over the decision message. The expiry window and the subject are the decision law's
+// (D8-D10), checked before this runs.
+func DecideBySignature(in DecisionSignatureInput) (*SignedDecision, error) {
+	if err := checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatSSH); err != nil { // P2-P5
+		return nil, err
+	}
+	if err := VerifyDecisionSSHSIG(in.Signature, in.Verb, in.Packet, in.Request, in.Expiry, in.Approver.Key); err != nil {
+		return nil, err
+	}
+	return &SignedDecision{ok: true, verb: in.Verb, packet: in.Packet, request: in.Request, expiry: in.Expiry,
+		name: in.Approver.Name, format: in.Approver.Format, key: bytes.Clone(in.Approver.Key.Blob)}, nil
+}
+
+// Binds reports whether s is a verified decision of exactly this verb, packet, request and expiry.
+func (s *SignedDecision) Binds(verb byte, packet, request [32]byte, expiry int64) error {
+	switch {
+	case s == nil || !s.ok:
+		return ErrNotVerified
+	case s.verb != verb:
+		return &BindingError{Field: "verb"}
+	case s.packet != packet:
+		return &BindingError{Field: "packet"}
+	case s.request != request:
+		return &BindingError{Field: "request"}
+	case s.expiry != expiry:
+		return &BindingError{Field: "expiry"}
+	}
+	return nil
+}
+
+// DecidedBy is the decision's decidedBy, format ":" name ("ssh:mark-mac"); "" unless verified.
+func (s *SignedDecision) DecidedBy() string {
+	if s == nil || !s.ok {
+		return ""
+	}
+	return s.format + ":" + s.name
+}
+
+// Use is what the evidence row binds in SQL: the approver's name, format and key (a copy).
+func (s *SignedDecision) Use() (name, format string, key []byte, ok bool) {
+	if s == nil || !s.ok {
+		return "", "", nil, false
+	}
+	return s.name, s.format, bytes.Clone(s.key), true
+}
diff --git a/host/approval/decision_signed_test.go b/host/approval/decision_signed_test.go
new file mode 100644
index 0000000..fd3cdae
--- /dev/null
+++ b/host/approval/decision_signed_test.go
@@ -0,0 +1,127 @@
+package approval_test
+
+// Row 201 A4 sprint 2 M4d (design §2.3, §6 DM3-DM6, DM8; the passkey halves are sprint 3's): the
+// decision signature law with approvaltest keys (D-WORLD-109). An external test package because
+// approvaltest imports host/approval.
+
+import (
+	"errors"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+var dmP, dmQ = [32]byte{0x50}, [32]byte{0x51}
+
+const dmNow, dmExpiry = int64(1710), int64(1760)
+
+// dmInput is a valid decision input: mark-mac's approve of (P, Q, expiry), signed by s's key over
+// the decision message, enrolled before now.
+func dmInput(t *testing.T, s *approvaltest.Signer, verb byte) approval.DecisionSignatureInput {
+	t.Helper()
+	msg, err := approval.DecisionMessage(verb, dmP, dmQ, dmExpiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return approval.DecisionSignatureInput{Verb: verb, Packet: dmP, Request: dmQ, Expiry: dmExpiry, Now: dmNow,
+		Signature: s.Sign(approval.DecisionNamespace, msg),
+		Approver:  approval.Approver{Name: "mark-mac", Format: approval.FormatSSH, Key: s.PublicKey(), EnrolledAt: dmNow - 100}}
+}
+
+// DM8 + the control: P2-P5 each refuse a good signature (name, format, revoked, clock behind);
+// the valid input decides as ssh:mark-mac, binds its subject and hands the key to the evidence row.
+func TestDecideBySignatureApprover(t *testing.T) {
+	s := approvaltest.NewSigner(41)
+	for _, tc := range []struct {
+		name string
+		edit func(a *approval.Approver)
+		want error
+	}{
+		{"P2 name", func(a *approval.Approver) { a.Name = "Mark" }, approval.ErrApproverName},
+		{"P3 format", func(a *approval.Approver) { a.Format = approval.FormatWebAuthn }, approval.ErrApproverFormat},
+		{"P4 revoked", func(a *approval.Approver) { a.RevokedAt = dmNow - 1 }, approval.ErrApproverRevoked},
+		{"P5 clock behind", func(a *approval.Approver) { a.EnrolledAt = dmNow + 1 }, approval.ErrApproverClockBehind},
+	} {
+		in := dmInput(t, s, approval.VerbApprove)
+		tc.edit(&in.Approver)
+		if d, err := approval.DecideBySignature(in); d != nil || !errors.Is(err, tc.want) {
+			t.Errorf("%s: %v, %v; want %v", tc.name, d, err, tc.want)
+		}
+	}
+	d, err := approval.DecideBySignature(dmInput(t, s, approval.VerbApprove))
+	if err != nil || d.DecidedBy() != "ssh:mark-mac" || d.Binds(approval.VerbApprove, dmP, dmQ, dmExpiry) != nil {
+		t.Fatalf("control: %v, %v", d, err)
+	}
+	if name, format, key, ok := d.Use(); !ok || name != "mark-mac" || format != approval.FormatSSH || string(key) != string(s.PublicKey().Blob) {
+		t.Fatalf("Use = %q %q %x %v", name, format, key, ok)
+	}
+}
+
+// DM3 (SSH half): the mint message signed under the mint namespace, and the mint message signed
+// under the decision namespace, are refused by the decision law; DM4: a decision signature is
+// refused by the mint law (namespace), so neither artifact crosses.
+func TestMintSignatureIsNotADecision(t *testing.T) {
+	s := approvaltest.NewSigner(42)
+	mint, err := approval.SignedMessage(dmQ, dmP, dmExpiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	decision, _ := approval.DecisionMessage(approval.VerbApprove, dmP, dmQ, dmExpiry)
+	for name, sig := range map[string][]byte{"mint namespace": s.Sign(approval.Namespace, mint), "mint message": s.Sign(approval.DecisionNamespace, mint),
+		"decision message, mint namespace": s.Sign(approval.Namespace, decision)} {
+		in := dmInput(t, s, approval.VerbApprove)
+		in.Signature = sig
+		if _, err := approval.DecideBySignature(in); err == nil {
+			t.Errorf("%s: a mint signature decided", name)
+		}
+	}
+	in := dmInput(t, s, approval.VerbApprove)
+	if err := approval.VerifyApproval(in.Signature, dmQ, dmP, dmExpiry, s.PublicKey()); !errors.Is(err, approval.ErrNamespace) {
+		t.Fatalf("DM4: a decision signature checked as a mint: %v, want ErrNamespace", err)
+	}
+}
+
+// DM5: an approve signature posted as a deny (and the reverse) does not verify.
+func TestDecisionBindsTheVerb(t *testing.T) {
+	s := approvaltest.NewSigner(43)
+	for _, v := range [][2]byte{{approval.VerbApprove, approval.VerbDeny}, {approval.VerbDeny, approval.VerbApprove}} {
+		in := dmInput(t, s, v[0])
+		in.Verb = v[1]
+		if _, err := approval.DecideBySignature(in); !errors.Is(err, approval.ErrBadSignature) {
+			t.Errorf("signed %d, posted %d: %v, want ErrBadSignature", v[0], v[1], err)
+		}
+	}
+	in := dmInput(t, s, approval.VerbApprove)
+	in.Expiry++ // the expiry is signed too
+	if _, err := approval.DecideBySignature(in); !errors.Is(err, approval.ErrBadSignature) {
+		t.Errorf("another expiry: %v", err)
+	}
+}
+
+// DM6: Binds refuses each field that differs and the zero value; DecidedBy and Use are empty for
+// an unverified decision.
+func TestSignedDecisionBinds(t *testing.T) {
+	d, err := approval.DecideBySignature(dmInput(t, approvaltest.NewSigner(44), approval.VerbDeny))
+	if err != nil {
+		t.Fatal(err)
+	}
+	for field, err := range map[string]error{
+		"verb":    d.Binds(approval.VerbApprove, dmP, dmQ, dmExpiry),
+		"packet":  d.Binds(approval.VerbDeny, dmQ, dmQ, dmExpiry),
+		"request": d.Binds(approval.VerbDeny, dmP, dmP, dmExpiry),
+		"expiry":  d.Binds(approval.VerbDeny, dmP, dmQ, dmExpiry+1),
+	} {
+		var be *approval.BindingError
+		if !errors.As(err, &be) || be.Field != field {
+			t.Errorf("%s: %v", field, err)
+		}
+	}
+	var zero *approval.SignedDecision
+	if err := zero.Binds(approval.VerbDeny, dmP, dmQ, dmExpiry); !errors.Is(err, approval.ErrNotVerified) || zero.DecidedBy() != "" {
+		t.Fatalf("nil decision: %v %q", err, zero.DecidedBy())
+	}
+	if _, _, _, ok := (&approval.SignedDecision{}).Use(); ok || (&approval.SignedDecision{}).DecidedBy() != "" || !errors.Is((&approval.SignedDecision{}).Binds(0, dmP, dmQ, 0), approval.ErrNotVerified) {
+		t.Fatal("the zero decision is usable")
+	}
+}
````

### Patch M4e — credentials schema 8 (KD1, KD5, KD6; P2's fixtures)

````diff M4e
diff --git a/cmd/ailang-worldd/approve_await_refusals_test.go b/cmd/ailang-worldd/approve_await_refusals_test.go
index 6fd989c..496f68e 100644
--- a/cmd/ailang-worldd/approve_await_refusals_test.go
+++ b/cmd/ailang-worldd/approve_await_refusals_test.go
@@ -145,8 +145,8 @@ func TestApproveAwaitNeedsSchema7AndItsRequest(t *testing.T) {
 		gone  bool     // delete the request after the first wait
 		want  string
 	}{
-		{name: "v6-store", setup: []string{`DROP TABLE surface_approvals`, `PRAGMA user_version = 6`}, want: "surface approvals need schema 7: restart `ailang-worldd serve`"},
-		{name: "v5-store", setup: []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, want: "passkeys need schema 6"},
+		{name: "v6-store", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`}, want: "surface approvals need schema 7: restart `ailang-worldd serve`"},
+		{name: "v5-store", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, want: "passkeys need schema 6"},
 		{name: "deleted", gone: true, want: "is no longer in the store"},
 	} {
 		t.Run(arm.name, func(t *testing.T) {
diff --git a/cmd/ailang-worldd/approve_passkey_refusals_test.go b/cmd/ailang-worldd/approve_passkey_refusals_test.go
index 28c4ca1..de84d00 100644
--- a/cmd/ailang-worldd/approve_passkey_refusals_test.go
+++ b/cmd/ailang-worldd/approve_passkey_refusals_test.go
@@ -56,7 +56,7 @@ func TestApprovePasskeyRefusalsBeforeAnyEffect(t *testing.T) {
 		{name: "WF2-missing", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "absent.json") }, want: "--webauthn-assertion: open "},
 		{name: "WF3-not-an-assertion", file: func(*testing.T, approveRig) []byte { return []byte(`{"id":"x"}`) },
 			want: "is not a passkey assertion: " + approval.ErrPasskeyResponse.Error()},
-		{name: "WF4-v5-store", setup: sql(`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
+		{name: "WF4-v5-store", setup: sql(`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
 			want: "passkeys need schema 6"},
 		{name: "WF5-unknown-credential", file: func(t *testing.T, a approveRig) []byte {
 			return approvaltest.NewPasskey(approval.AlgEdDSA, 55).Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
diff --git a/cmd/ailang-worldd/approve_passkey_test.go b/cmd/ailang-worldd/approve_passkey_test.go
index e9d37d4..cf2dfc3 100644
--- a/cmd/ailang-worldd/approve_passkey_test.go
+++ b/cmd/ailang-worldd/approve_passkey_test.go
@@ -96,7 +96,7 @@ func TestApproveByPasskey(t *testing.T) {
 			t.Fatalf("rc %d, terminal %v, daemon calls %d, stderr %q", code, opened, *calls, stderr.String())
 		}
 		raw := rawCredentials(t, a.credPath)
-		for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+		for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 			if _, err := raw.Exec(stmt); err != nil {
 				t.Fatal(err)
 			}
diff --git a/cmd/ailang-worldd/approve_sig_test.go b/cmd/ailang-worldd/approve_sig_test.go
index 2a85b55..3bdfcbf 100644
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -216,7 +216,7 @@ func TestApproveBySignatureRefusalsBeforeAnyEffect(t *testing.T) {
 		{name: "v4-store", setup: func(t *testing.T, a approveRig) {
 			raw := rawCredentials(t, a.credPath)
 			defer func() { _ = raw.Close() }()
-			for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+			for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 				if _, err := raw.Exec(stmt); err != nil {
 					t.Fatal(err)
 				}
diff --git a/cmd/ailang-worldd/approver_label_test.go b/cmd/ailang-worldd/approver_label_test.go
index 738a82c..40ac7ea 100644
--- a/cmd/ailang-worldd/approver_label_test.go
+++ b/cmd/ailang-worldd/approver_label_test.go
@@ -82,7 +82,7 @@ func TestApproverListAndRevokeShowPasskeys(t *testing.T) {
 	db5, creds5 := approverStore(t)
 	seedApprover(t, creds5, "s", fixtureKeyOf(t, "ed25519.pub"), approveT0)
 	raw = rawCredentials(t, creds5)
-	for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+	for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 		if _, err := raw.Exec(stmt); err != nil {
 			t.Fatal(err)
 		}
diff --git a/cmd/ailang-worldd/approver_passkey_test.go b/cmd/ailang-worldd/approver_passkey_test.go
index 1bdb201..4f95f65 100644
--- a/cmd/ailang-worldd/approver_passkey_test.go
+++ b/cmd/ailang-worldd/approver_passkey_test.go
@@ -154,7 +154,7 @@ func TestPasskeyEnrollFlagRules(t *testing.T) {
 		{name: "N6-v5", setup: func(t *testing.T, r enrolRig) {
 			raw := rawCredentials(t, r.creds)
 			defer func() { _ = raw.Close() }()
-			for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+			for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
 				if _, err := raw.Exec(stmt); err != nil {
 					t.Fatal(err)
 				}
diff --git a/cmd/ailang-worldd/approver_rig_test.go b/cmd/ailang-worldd/approver_rig_test.go
index 08f340c..b3a176a 100644
--- a/cmd/ailang-worldd/approver_rig_test.go
+++ b/cmd/ailang-worldd/approver_rig_test.go
@@ -64,7 +64,7 @@ func approverStoreV4(t *testing.T) (db, creds string) {
 	db, creds = approverStore(t)
 	raw := rawCredentials(t, creds)
 	defer func() { _ = raw.Close() }()
-	for _, stmt := range []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+	for _, stmt := range []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 		if _, err := raw.Exec(stmt); err != nil {
 			t.Fatal(err)
 		}
diff --git a/host/authority/judge_test.go b/host/authority/judge_test.go
index c72748b..2af94a6 100644
--- a/host/authority/judge_test.go
+++ b/host/authority/judge_test.go
@@ -174,7 +174,7 @@ func TestJudgeAssertionIsTheCLIJudge(t *testing.T) {
 	arms := append(requestArms(good),
 		judgeArm{name: "WF3-not-an-assertion", artifact: func(*testing.T, store.MintRequestRow) []byte { return []byte(`{"id":"x"}`) },
 			step: JudgeStepParse, cause: approval.ErrPasskeyResponse},
-		judgeArm{name: "WF4-v5-store", setup: []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, artifact: good, step: JudgeStepStore},
+		judgeArm{name: "WF4-v5-store", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`}, artifact: good, step: JudgeStepStore},
 		judgeArm{name: "WF5-unknown-credential", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
 			return stranger.Assert(jgOrigin, jgRPID, jgHash(t, row), 2)
 		}, step: JudgeStepUnknownKey, cause: ErrNotEnrolled, field: string(stranger.CredentialID())},
@@ -218,7 +218,7 @@ func TestJudgeSSHSignatureIsTheCLIJudge(t *testing.T) {
 	arms := append(requestArms(good),
 		judgeArm{name: "F11-garbage", artifact: func(*testing.T, store.MintRequestRow) []byte { return []byte("not a signature\n") },
 			step: JudgeStepParse, cause: approval.ErrArmor},
-		judgeArm{name: "F12-v4-store", setup: []string{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
+		judgeArm{name: "F12-v4-store", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
 			artifact: good, step: JudgeStepStore},
 		judgeArm{name: "F13-unenrolled-key", artifact: func(t *testing.T, row store.MintRequestRow) []byte {
 			return stranger.Sign(approval.Namespace, jgMessage(t, row))
diff --git a/host/daemon/surface_approve_post_test.go b/host/daemon/surface_approve_post_test.go
index 773e24a..0a4f3a0 100644
--- a/host/daemon/surface_approve_post_test.go
+++ b/host/daemon/surface_approve_post_test.go
@@ -140,7 +140,7 @@ func TestApprovePOSTRefusalArms(t *testing.T) {
 		}},
 		{"R11-store-v5", "", "Restart", "restart `ailang-worldd serve`", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
 			body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
-			setVersion(t, a.d.credentialsPath, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
+			setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
 			return body, nil
 		}},
 		{"R13-busy", "", "Busy", "busy", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
@@ -152,7 +152,7 @@ func TestApprovePOSTRefusalArms(t *testing.T) {
 			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
 		}},
 		{"R14-v6", "", "Restart", "restart `ailang-worldd serve`", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
-			setVersion(t, a.d.credentialsPath, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`)
+			setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`)
 			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
 		}},
 		{"R14-race", "", "Conflict", "changed while it was approved", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
diff --git a/host/daemon/surface_approve_read_test.go b/host/daemon/surface_approve_read_test.go
index a45ac67..becea07 100644
--- a/host/daemon/surface_approve_read_test.go
+++ b/host/daemon/surface_approve_read_test.go
@@ -63,7 +63,7 @@ func TestRequestsJSON(t *testing.T) {
 			t.Errorf("requests.json carries %q", leak)
 		}
 	}
-	setVersion(t, a.d.credentialsPath, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`) // below 7: listed, nothing recorded
+	setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`) // below 7: listed, nothing recorded
 	if body := a.serve("GET", "/surface/requests.json", apLoopback, nil).Body.String(); strings.Count(body, `"recorded":null`) != 3 || !strings.Contains(body, `"unavailable":""`) {
 		t.Fatalf("a schema-6 store: %s", body)
 	}
@@ -128,7 +128,7 @@ func TestRequestCardJSON(t *testing.T) {
 	if err := json.Unmarshal(a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil).Body.Bytes(), &c); err != nil || !c.Consumed || c.ApprovedBy != "tty" {
 		t.Fatalf("consumed card = %+v, %v", c, err)
 	}
-	setVersion(t, a.d.credentialsPath, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
+	setVersion(t, a.d.credentialsPath, `DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
 	if rec := a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"passkeys":[]`) {
 		t.Fatalf("a schema-5 store's card = %d %s, want 200 with no passkeys", rec.Code, rec.Body.String())
 	}
diff --git a/host/store/approvers_test.go b/host/store/approvers_test.go
index 432dd12..3d1baea 100644
--- a/host/store/approvers_test.go
+++ b/host/store/approvers_test.go
@@ -131,7 +131,7 @@ func TestCredentialsUpgradeToV6(t *testing.T) {
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["approvers"]), normalizeDDL(canonicalApproversDDL); got != want {
 				t.Fatalf("upgraded approvers DDL:\n got: %s\nwant: %s", got, want)
 			}
diff --git a/host/store/credentials.go b/host/store/credentials.go
index 6a41b8b..092bb0f 100644
--- a/host/store/credentials.go
+++ b/host/store/credentials.go
@@ -39,8 +39,9 @@ var credentialsSchemaSQL string
 // independent of the world store's currentSchemaVersion. Version 2 added
 // profile bindings; version 3 adds immutable attended episode start history;
 // version 4 adds mint requests (row 201); version 5 adds approvers (row 201 A1b); version 6 adds
-// passkeys (row 201 A2); version 7 adds the surface's recorded approvals (row 201 A3).
-const credentialsSchemaVersion = 7
+// passkeys (row 201 A2); version 7 adds the surface's recorded approvals (row 201 A3); version 8
+// adds the key decisions of Human.Approve requests (row 201 A4).
+const credentialsSchemaVersion = 8
 
 // CredentialsSchemaVersion exports credentialsSchemaVersion: serve compares a
 // store's version with it to decide whether to upgrade.
@@ -71,7 +72,9 @@ const credentialsTable = "session_credentials"
 // credentialsTables is the exact, sorted application table set of a credentials store at version.
 func credentialsTables(version int) []string {
 	switch {
-	case version >= 7:
+	case version >= 8:
+		return []string{approversTable, "episode_starts", keyDecisionsTable, mintRequestsTable, credentialsTable, surfaceApprovalsTable, passkeysTable}
+	case version == 7:
 		return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable, surfaceApprovalsTable, passkeysTable}
 	case version == 6:
 		return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable, passkeysTable}
@@ -369,6 +372,9 @@ func upgradeCredentials(display string, db *sql.DB, version int) error {
 	if version < 7 {
 		statements = append(statements, surfaceApprovalsSQL)
 	}
+	if version < 8 {
+		statements = append(statements, keyDecisionsSQL)
+	}
 	for _, stmt := range statements {
 		if _, err := tx.Exec(stmt); err != nil {
 			return fmt.Errorf("store: upgrade credentials %q: %w", display, err)
diff --git a/host/store/credentials_schema.sql b/host/store/credentials_schema.sql
index 9e82f28..f3c4bd4 100644
--- a/host/store/credentials_schema.sql
+++ b/host/store/credentials_schema.sql
@@ -96,3 +96,18 @@ CREATE TABLE IF NOT EXISTS surface_approvals (
     digest      TEXT NOT NULL CHECK (length(digest) = 64),
     received_at INTEGER NOT NULL
 );
+
+-- Key decisions of Human.Approve requests (version 8, row 201 A4). The evidence of what verified,
+-- kept beside the approvers; the authority is the world store's decision, never this row.
+-- decision_ref is NULL until the world decision landed.
+CREATE TABLE IF NOT EXISTS key_decisions (
+    request_ref   TEXT PRIMARY KEY CHECK (length(request_ref) = 71 AND request_ref GLOB 'sha256:*'),
+    packet_ref    TEXT NOT NULL CHECK (length(packet_ref) = 71 AND packet_ref GLOB 'sha256:*'),
+    verb          TEXT NOT NULL CHECK (verb IN ('approve', 'deny')),
+    format        TEXT NOT NULL CHECK (format IN ('ssh', 'webauthn')),
+    approver      TEXT NOT NULL CHECK (approver <> ''),
+    artifact      BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
+    expiry        INTEGER NOT NULL,
+    decided_at    INTEGER NOT NULL,
+    decision_ref  TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*'))
+);
diff --git a/host/store/credentials_test.go b/host/store/credentials_test.go
index b995084..39f7afd 100644
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -76,7 +76,7 @@ var canonicalPasskeysDDL = `CREATE TABLE webauthn_credentials (
 func TestCredentialsSchemaMatchesCanonicalDDL(t *testing.T) {
 	c := openTestCredentials(t)
 	ddl := tableDDL(t, c.db)
-	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
+	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
 	for name, got := range map[string]string{"schema file": ddl["webauthn_credentials"], "passkeysSQL": passkeysSQL} {
 		if normalizeDDL(got) != normalizeDDL(canonicalPasskeysDDL) {
 			t.Fatalf("webauthn_credentials DDL (%s):\n got: %s\nwant: %s", name, got, canonicalPasskeysDDL)
@@ -112,8 +112,8 @@ func TestCredentialsSchemaMatchesCanonicalDDL(t *testing.T) {
 		t.Fatalf("v4 world table differs from the version-1 credentials table: %s", got)
 	}
 	var version int
-	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != 7 {
-		t.Fatalf("credentials user_version = %d %v, want 7 (row 201 A3)", version, err)
+	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != 8 {
+		t.Fatalf("credentials user_version = %d %v, want 8 (row 201 A4)", version, err)
 	}
 }
 
@@ -300,6 +300,9 @@ func TestCredentialsMethodsRequireADeadline(t *testing.T) {
 	_, errs["ActivePasskeys"] = c.ActivePasskeys(ctx)
 	errs["RecordSurfaceApproval"] = c.RecordSurfaceApproval(ctx, SurfaceApprovalRow{}, 1000)
 	_, _, errs["SurfaceApproval"] = c.SurfaceApproval(ctx, "mr-0000000000000001")
+	errs["RecordKeyDecision"] = c.RecordKeyDecision(ctx, KeyDecisionRow{}, ApproverUse{}, 1000)
+	errs["LandKeyDecision"] = c.LandKeyDecision(ctx, "x", "y")
+	_, _, errs["KeyDecision"] = c.KeyDecision(ctx, "x")
 	for name, err := range errs {
 		if !errors.Is(err, ErrNoDeadline) {
 			t.Errorf("%s(Background) = %v, want ErrNoDeadline", name, err)
diff --git a/host/store/key_decisions.go b/host/store/key_decisions.go
new file mode 100644
index 0000000..15adb54
--- /dev/null
+++ b/host/store/key_decisions.go
@@ -0,0 +1,169 @@
+package store
+
+// Row 201 A4 sprint 2 (design §2.7): credentials schema 8's key_decisions, the evidence of a key
+// decision. The authority is the world store's decision, never this row; decision_ref links the two
+// once it landed (NULL until then). Format-agnostic: sprint 3's passkey arm adds no column (S2-D3).
+
+import (
+	"bytes"
+	"context"
+	"database/sql"
+	"errors"
+	"fmt"
+)
+
+const keyDecisionsTable = "key_decisions"
+
+// keyDecisionsSchemaVersion is the first credentials schema with key_decisions.
+const keyDecisionsSchemaVersion = 8
+
+// keyDecisionsSQL is the version 7 -> 8 upgrade: credentials_schema.sql's key_decisions table,
+// column for column (TestCredentialsSchema8KeyDecisions compares the two), the design's DDL verbatim.
+const keyDecisionsSQL = `CREATE TABLE key_decisions (
+    request_ref   TEXT PRIMARY KEY CHECK (length(request_ref) = 71 AND request_ref GLOB 'sha256:*'),
+    packet_ref    TEXT NOT NULL CHECK (length(packet_ref) = 71 AND packet_ref GLOB 'sha256:*'),
+    verb          TEXT NOT NULL CHECK (verb IN ('approve', 'deny')),
+    format        TEXT NOT NULL CHECK (format IN ('ssh', 'webauthn')),
+    approver      TEXT NOT NULL CHECK (approver <> ''),
+    artifact      BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
+    expiry        INTEGER NOT NULL,
+    decided_at    INTEGER NOT NULL,
+    decision_ref  TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*'))
+)`
+
+var (
+	// ErrBadKeyDecision refuses a row or a ref that breaks the table's grammar (K2).
+	ErrBadKeyDecision = errors.New("store: invalid key decision row")
+	// ErrKeyDecisionLanded refuses recording a request whose key decision already landed (K0).
+	ErrKeyDecisionLanded = errors.New("store: the request's key decision already landed")
+	// ErrKeyDecisionNotPending refuses landing a request with no unlanded key decision.
+	ErrKeyDecisionNotPending = errors.New("store: no unlanded key decision for the request")
+)
+
+// KeyDecisionRow is one key_decisions row.
+type KeyDecisionRow struct {
+	RequestRef  string // sha256: and 64 lowercase hex, the request decided
+	PacketRef   string // the one packet carrying it
+	Verb        string // "approve" or "deny"
+	Format      string // "ssh" (sprint 2) or "webauthn" (sprint 3)
+	Approver    string // "<format>:<name>", the decision's decidedBy
+	Artifact    []byte // the armored SSHSIG (or the assertion), 1..16384 bytes
+	Expiry      int64  // the signed expiry
+	DecidedAt   int64  // set by RecordKeyDecision to its now
+	DecisionRef string // "" until LandKeyDecision
+}
+
+const keyDecisionColumns = `request_ref, packet_ref, verb, format, approver, artifact, expiry, decided_at, coalesce(decision_ref, '')`
+
+func (c *Credentials) olderForKeyDecisions() error {
+	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: keyDecisionsSchemaVersion, Feature: "key decisions"}
+}
+
+// isDecisionRef is the table's ref grammar, held in code to lower-case hex (the CHECKs hold length and prefix).
+func isDecisionRef(s string) bool { return len(s) == 71 && s[:7] == "sha256:" && isLowerHex64(s[7:]) }
+
+// checkKeyDecisionRow is RecordKeyDecision's grammar (K2): each refusal is ErrBadKeyDecision.
+func checkKeyDecisionRow(row KeyDecisionRow, use ApproverUse) error {
+	switch {
+	case !isDecisionRef(row.RequestRef) || !isDecisionRef(row.PacketRef):
+		return fmt.Errorf("%w: request %q, packet %q", ErrBadKeyDecision, row.RequestRef, row.PacketRef)
+	case row.Verb != "approve" && row.Verb != "deny":
+		return fmt.Errorf("%w: verb %q", ErrBadKeyDecision, row.Verb)
+	case row.Format != use.Format || row.Approver != use.Format+":"+use.Name:
+		return fmt.Errorf("%w: %s by %q is not the approver %s:%s", ErrBadKeyDecision, row.Format, row.Approver, use.Format, use.Name)
+	case len(row.Artifact) < 1 || len(row.Artifact) > MaxSurfaceApprovalArtifactBytes:
+		return fmt.Errorf("%w: an artifact is 1 to %d bytes, not %d", ErrBadKeyDecision, MaxSurfaceApprovalArtifactBytes, len(row.Artifact))
+	}
+	return nil
+}
+
+func scanKeyDecision(scan func(...any) error) (KeyDecisionRow, error) {
+	var r KeyDecisionRow
+	err := scan(&r.RequestRef, &r.PacketRef, &r.Verb, &r.Format, &r.Approver, &r.Artifact, &r.Expiry, &r.DecidedAt, &r.DecisionRef)
+	return r, err
+}
+
+// RecordKeyDecision records row for the approver use at now, in one transaction (design §2.7): K1
+// deadline, schema 8; K2 grammar; K0 a landed row refuses, an unlanded one equal in artifact and
+// signed fields is kept (the exact-match retry); K3 useApprover on both paths; K5 insert or replace
+// the unlanded row. The passkey counter (K4) is sprint 3's.
+func (c *Credentials) RecordKeyDecision(ctx context.Context, row KeyDecisionRow, use ApproverUse, now int64) error {
+	if err := requireDeadline(ctx); err != nil { // K1
+		return err
+	}
+	if c.version < keyDecisionsSchemaVersion { // K1
+		return c.olderForKeyDecisions()
+	}
+	if err := checkKeyDecisionRow(row, use); err != nil { // K2
+		return err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return fmt.Errorf("store: record key decision: begin: %w", err)
+	}
+	defer func() { _ = tx.Rollback() }()
+	old, err := scanKeyDecision(tx.QueryRowContext(ctx, `SELECT `+keyDecisionColumns+` FROM key_decisions WHERE request_ref = ?`, row.RequestRef).Scan)
+	found := err == nil
+	if err != nil && !errors.Is(err, sql.ErrNoRows) {
+		return fmt.Errorf("store: record key decision: %w", err)
+	}
+	if found && old.DecisionRef != "" { // K0: landed
+		return fmt.Errorf("%w: %s by %s", ErrKeyDecisionLanded, row.RequestRef, old.Approver)
+	}
+	exact := found && bytes.Equal(old.Artifact, row.Artifact) && old.PacketRef == row.PacketRef && old.Verb == row.Verb &&
+		old.Format == row.Format && old.Approver == row.Approver && old.Expiry == row.Expiry // K0: the exact match
+	if err := useApprover(ctx, tx, use, now); err != nil { // K3, on both paths
+		return err
+	}
+	if !exact { // K5
+		if _, err := tx.ExecContext(ctx, `INSERT INTO key_decisions (request_ref, packet_ref, verb, format, approver, artifact, expiry, decided_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
+ON CONFLICT(request_ref) DO UPDATE SET packet_ref = excluded.packet_ref, verb = excluded.verb, format = excluded.format, approver = excluded.approver,
+artifact = excluded.artifact, expiry = excluded.expiry, decided_at = excluded.decided_at WHERE key_decisions.decision_ref IS NULL`,
+			row.RequestRef, row.PacketRef, row.Verb, row.Format, row.Approver, row.Artifact, row.Expiry, now); err != nil {
+			return fmt.Errorf("store: record key decision: %w", err)
+		}
+	}
+	if err := tx.Commit(); err != nil {
+		return fmt.Errorf("store: record key decision: commit: %w", err)
+	}
+	return nil
+}
+
+// LandKeyDecision links request's unlanded key decision to the world decision decisionRef.
+func (c *Credentials) LandKeyDecision(ctx context.Context, request, decisionRef string) error {
+	if err := requireDeadline(ctx); err != nil {
+		return err
+	}
+	if c.version < keyDecisionsSchemaVersion {
+		return c.olderForKeyDecisions()
+	}
+	if !isDecisionRef(decisionRef) {
+		return fmt.Errorf("%w: decision %q", ErrBadKeyDecision, decisionRef)
+	}
+	res, err := c.db.ExecContext(ctx, `UPDATE key_decisions SET decision_ref = ? WHERE request_ref = ? AND decision_ref IS NULL`, decisionRef, request)
+	if err != nil {
+		return fmt.Errorf("store: land key decision: %w", err)
+	}
+	if n, err := res.RowsAffected(); err != nil || n != 1 {
+		return fmt.Errorf("%w: %s (%v)", ErrKeyDecisionNotPending, request, err)
+	}
+	return nil
+}
+
+// KeyDecision reads request's key decision; ok is false when there is none.
+func (c *Credentials) KeyDecision(ctx context.Context, request string) (KeyDecisionRow, bool, error) {
+	if err := requireDeadline(ctx); err != nil {
+		return KeyDecisionRow{}, false, err
+	}
+	if c.version < keyDecisionsSchemaVersion {
+		return KeyDecisionRow{}, false, c.olderForKeyDecisions()
+	}
+	row, err := scanKeyDecision(c.db.QueryRowContext(ctx, `SELECT `+keyDecisionColumns+` FROM key_decisions WHERE request_ref = ?`, request).Scan)
+	if errors.Is(err, sql.ErrNoRows) {
+		return KeyDecisionRow{}, false, nil
+	}
+	if err != nil {
+		return KeyDecisionRow{}, false, fmt.Errorf("store: read key decision: %w", err)
+	}
+	return row, true, nil
+}
diff --git a/host/store/key_decisions_test.go b/host/store/key_decisions_test.go
new file mode 100644
index 0000000..d05d736
--- /dev/null
+++ b/host/store/key_decisions_test.go
@@ -0,0 +1,122 @@
+package store
+
+// Row 201 A4 sprint 2 M4e (design §2.7, §6 KD1, KD5, KD6; KD3, KD4, KD7 and KD9 are M4f's, in
+// key_decisions_record_test.go; KD2's counter is the passkey arm's, sprint 3): credentials schema
+// 8's key_decisions. Every store is under t.TempDir() (D-WORLD-109).
+
+import (
+	"errors"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+)
+
+// canonicalKeyDecisionsDDL is the current (version 8, row 201 A4) key_decisions table: the
+// design's DDL, both sha256 CHECKs verbatim (KD5, KD6).
+var canonicalKeyDecisionsDDL = `CREATE TABLE key_decisions (
+    request_ref   TEXT PRIMARY KEY CHECK (length(request_ref) = 71 AND request_ref GLOB 'sha256:*'),
+    packet_ref    TEXT NOT NULL CHECK (length(packet_ref) = 71 AND packet_ref GLOB 'sha256:*'),
+    verb          TEXT NOT NULL CHECK (verb IN ('approve', 'deny')),
+    format        TEXT NOT NULL CHECK (format IN ('ssh', 'webauthn')),
+    approver      TEXT NOT NULL CHECK (approver <> ''),
+    artifact      BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
+    expiry        INTEGER NOT NULL,
+    decided_at    INTEGER NOT NULL,
+    decision_ref  TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*'))
+)`
+
+var kdQ, kdP, kdD = "sha256:" + strings.Repeat("a1", 32), "sha256:" + strings.Repeat("b2", 32), "sha256:" + strings.Repeat("c3", 32)
+
+// kdStore is a schema-8 store with mark-ssh (an ssh approver, key "k1") enrolled at 900.
+func kdStore(t *testing.T) *Credentials {
+	t.Helper()
+	c := openTestCredentials(t)
+	if err := c.EnrolApprover(testCtx(t), ApproverRow{Name: "mark-ssh", Format: "ssh", PublicKey: []byte("k1"), EnrolledAt: 900}); err != nil {
+		t.Fatal(err)
+	}
+	return c
+}
+
+var kdUse = ApproverUse{Name: "mark-ssh", Format: "ssh", PublicKey: []byte("k1")}
+
+func kdRow() KeyDecisionRow {
+	return KeyDecisionRow{RequestRef: kdQ, PacketRef: kdP, Verb: "approve", Format: "ssh", Approver: "ssh:mark-ssh", Artifact: []byte("sig-1"), Expiry: 1760}
+}
+
+func kdRead(t *testing.T, c *Credentials) (KeyDecisionRow, bool) {
+	t.Helper()
+	r, ok, err := c.KeyDecision(testCtx(t), kdQ)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return r, ok
+}
+
+// KD1: a fresh store is schema 8 with key_decisions as credentials_schema.sql and the upgrade
+// write it; a schema-7 store is upgraded by UpgradeCredentials only, and before that every key
+// decision method refuses naming the serve restart.
+func TestCredentialsSchema8KeyDecisions(t *testing.T) {
+	c := openTestCredentials(t)
+	for name, got := range map[string]string{"schema file": tableDDL(t, c.db)["key_decisions"], "keyDecisionsSQL": keyDecisionsSQL} {
+		if normalizeDDL(got) != normalizeDDL(canonicalKeyDecisionsDDL) {
+			t.Fatalf("key_decisions DDL (%s):\n got: %s\nwant: %s", name, got, canonicalKeyDecisionsDDL)
+		}
+	}
+	if got, want := credentialsTables(8), []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"}; !reflect.DeepEqual(got, want) {
+		t.Fatalf("credentialsTables(8) = %v", got)
+	}
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	makeV6Credentials(t, path, []SessionRow{oneBoundRow()})
+	raw := rawDB(t, path)
+	if _, err := raw.Exec(surfaceApprovalsSQL); err != nil {
+		t.Fatal(err)
+	}
+	setPositiveFixtureVersion(t, raw, 7)
+	v7, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var oe *OlderCredentialsStoreError
+	_, _, readErr := v7.KeyDecision(testCtx(t), kdQ)
+	for name, err := range map[string]error{"record": v7.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1000), "land": v7.LandKeyDecision(testCtx(t), kdQ, kdD), "read": readErr} {
+		if !errors.As(err, &oe) || oe.Version != 7 || oe.Need != 8 || !strings.Contains(err.Error(), "key decisions need schema 8: restart `ailang-worldd serve`") {
+			t.Errorf("v7 %s: %v, want the restart refusal naming schema 8", name, err)
+		}
+	}
+	_ = v7.Close()
+	if upgraded, err := UpgradeCredentials(path); err != nil || !upgraded || userVersion(t, path) != 8 {
+		t.Fatalf("UpgradeCredentials: %v %v, version %d", upgraded, err, userVersion(t, path))
+	}
+	up := tableDDL(t, rawDB(t, path))
+	requireExactTableNames(t, "upgraded v7 store", up, credentialsTables(8))
+	if normalizeDDL(up["key_decisions"]) != normalizeDDL(canonicalKeyDecisionsDDL) {
+		t.Fatalf("upgraded key_decisions DDL: %s", up["key_decisions"])
+	}
+}
+
+// KD5, KD6: the DDL's ref CHECKs, by raw INSERTs: packet_ref and decision_ref of length 71 not
+// starting sha256:, of length 70, and decision_ref "" are constraint errors; sha256: + 64 hex and
+// a NULL decision_ref are accepted (the controls).
+func TestKeyDecisionsRefChecks(t *testing.T) {
+	c := openTestCredentials(t)
+	ins := func(request, packet string, decision any) error {
+		_, err := c.db.Exec(`INSERT INTO key_decisions (request_ref, packet_ref, verb, format, approver, artifact, expiry, decided_at, decision_ref) VALUES (?, ?, 'approve', 'ssh', 'ssh:m', x'01', 1, 1, ?)`, request, packet, decision)
+		return err
+	}
+	sha512 := "sha512:" + strings.Repeat("ab", 32)
+	for i, tc := range []struct {
+		packet   string
+		decision any
+		ok       bool
+	}{
+		{kdP, nil, true}, {kdP, kdD, true}, // controls
+		{sha512, nil, false}, {kdP[:70], nil, false}, // KD5
+		{kdP, sha512, false}, {kdP, kdD[:70], false}, {kdP, "", false}, // KD6
+	} {
+		request := "sha256:" + strings.Repeat(string(rune('0'+i)), 64)
+		if err := ins(request, tc.packet, tc.decision); (err == nil) != tc.ok {
+			t.Errorf("packet %q decision %v: err %v, want ok %v", tc.packet, tc.decision, err, tc.ok)
+		}
+	}
+}
diff --git a/host/store/mint_requests_test.go b/host/store/mint_requests_test.go
index c2ca693..3c3bb8e 100644
--- a/host/store/mint_requests_test.go
+++ b/host/store/mint_requests_test.go
@@ -129,7 +129,7 @@ func TestCredentialsUpgradeToV4(t *testing.T) {
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "key_decisions", "mint_requests", "session_credentials", "surface_approvals", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["mint_requests"]), normalizeDDL(canonicalMintRequestsDDL); got != want {
 				t.Fatalf("upgraded mint_requests DDL:\n got: %s\nwant: %s", got, want)
 			}
diff --git a/host/store/surface_approvals_test.go b/host/store/surface_approvals_test.go
index b0bbad1..55991cf 100644
--- a/host/store/surface_approvals_test.go
+++ b/host/store/surface_approvals_test.go
@@ -110,11 +110,11 @@ func TestCredentialsSchema7SurfaceApprovals(t *testing.T) {
 	if again, err := UpgradeCredentials(path); err != nil || again {
 		t.Fatalf("second UpgradeCredentials: upgraded %v, err %v; want false, nil", again, err)
 	}
-	if got := userVersion(t, path); got != 7 {
-		t.Fatalf("user_version after the upgrade = %d, want 7", got)
+	if got := userVersion(t, path); got != credentialsSchemaVersion {
+		t.Fatalf("user_version after the upgrade = %d, want %d", got, credentialsSchemaVersion)
 	}
 	up := tableDDL(t, rawDB(t, path))
-	requireExactTableNames(t, "upgraded v6 store", up, credentialsTables(7))
+	requireExactTableNames(t, "upgraded v6 store", up, credentialsTables(credentialsSchemaVersion))
 	if got, want := normalizeDDL(up["surface_approvals"]), normalizeDDL(canonicalSurfaceApprovalsDDL); got != want {
 		t.Fatalf("upgraded surface_approvals DDL:\n got: %s\nwant: %s", got, want)
 	}
````

### Patch M4f — the evidence writer's transaction (KD3, KD4, KD7, KD9, K2)

````diff M4f
diff --git a/host/store/key_decisions_record_test.go b/host/store/key_decisions_record_test.go
new file mode 100644
index 0000000..b62ac86
--- /dev/null
+++ b/host/store/key_decisions_record_test.go
@@ -0,0 +1,116 @@
+package store
+
+// Row 201 A4 sprint 2 M4f (design §2.7, §6 KD3, KD4, KD7 and KD9 in their SSH halves, K2):
+// RecordKeyDecision's transaction and grammar, LandKeyDecision. Every store is under t.TempDir()
+// (D-WORLD-109); the fixtures are key_decisions_test.go's.
+
+import (
+	"errors"
+	"reflect"
+	"strings"
+	"testing"
+)
+
+// KD3 (and KD9's store half): an approver revoked after the read-only judge refuses inside the
+// transaction on the first record and on an exact-match retry alike; nothing is written.
+func TestRecordKeyDecisionRevokedSince(t *testing.T) {
+	c := kdStore(t)
+	if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1000); err != nil {
+		t.Fatal(err)
+	}
+	first, _ := kdRead(t, c)
+	if err := c.RevokeApprover(testCtx(t), "mark-ssh", 1001); err != nil {
+		t.Fatal(err)
+	}
+	for _, row := range []KeyDecisionRow{kdRow(), func() KeyDecisionRow { r := kdRow(); r.Artifact = []byte("sig-2"); return r }()} {
+		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, 1002); !errors.Is(err, ErrApproverRevoked) {
+			t.Errorf("artifact %s after the revoke: %v, want ErrApproverRevoked", row.Artifact, err)
+		}
+	}
+	if got, _ := kdRead(t, c); !reflect.DeepEqual(got, first) {
+		t.Fatalf("the row changed: %+v, was %+v", got, first)
+	}
+	other := kdUse
+	other.PublicKey = []byte("k2")
+	if err := c.RecordKeyDecision(testCtx(t), func() KeyDecisionRow { r := kdRow(); r.RequestRef = kdD; return r }(), other, 1003); !errors.Is(err, ErrApproverRevoked) {
+		t.Fatalf("another key for the name: %v", err)
+	}
+}
+
+// KD4, KD7 (SSH half): the exact artifact retried keeps the unlanded row as it was (decided_at
+// of the first record); a row differing in the artifact or any signed field replaces it; once
+// landed, no row replaces it (ErrKeyDecisionLanded) and it cannot land twice.
+func TestRecordKeyDecisionReplaceUnlanded(t *testing.T) {
+	c := kdStore(t)
+	if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1000); err != nil {
+		t.Fatal(err)
+	}
+	if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, 1010); err != nil {
+		t.Fatal(err)
+	}
+	if got, ok := kdRead(t, c); !ok || got.DecidedAt != 1000 || got.DecisionRef != "" {
+		t.Fatalf("exact retry: %+v %v, want the first record kept", got, ok)
+	}
+	at := int64(1020)
+	for name, edit := range map[string]func(r *KeyDecisionRow){
+		"artifact": func(r *KeyDecisionRow) { r.Artifact = []byte("sig-2") },
+		"packet":   func(r *KeyDecisionRow) { r.PacketRef = kdD },
+		"verb":     func(r *KeyDecisionRow) { r.Verb = "deny" },
+		"expiry":   func(r *KeyDecisionRow) { r.Expiry++ },
+	} {
+		row := kdRow()
+		edit(&row)
+		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, at); err != nil {
+			t.Fatal(err)
+		}
+		row.DecidedAt = at
+		if got, _ := kdRead(t, c); !reflect.DeepEqual(got, row) {
+			t.Errorf("%s: %+v, want the row replaced by %+v", name, got, row)
+		}
+		if err := c.RecordKeyDecision(testCtx(t), kdRow(), kdUse, at+1); err != nil { // back to the first, for the next arm
+			t.Fatal(err)
+		}
+		at += 10
+	}
+	if err := c.LandKeyDecision(testCtx(t), kdQ, kdD); err != nil {
+		t.Fatal(err)
+	}
+	if got, _ := kdRead(t, c); got.DecisionRef != kdD {
+		t.Fatalf("landed: %+v", got)
+	}
+	for _, row := range []KeyDecisionRow{kdRow(), func() KeyDecisionRow { r := kdRow(); r.Artifact = []byte("sig-3"); return r }()} {
+		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, 2000); !errors.Is(err, ErrKeyDecisionLanded) {
+			t.Errorf("after landing, %s: %v, want ErrKeyDecisionLanded", row.Artifact, err)
+		}
+	}
+	if err := c.LandKeyDecision(testCtx(t), kdQ, kdD); !errors.Is(err, ErrKeyDecisionNotPending) {
+		t.Fatalf("second land: %v", err)
+	}
+}
+
+// K2: the grammar refuses a bad ref, an upper-case ref, a verb, an approver that is not the use's,
+// a format that is not the use's, and an empty or oversized artifact; nothing is written.
+func TestRecordKeyDecisionGrammar(t *testing.T) {
+	c := kdStore(t)
+	for name, edit := range map[string]func(r *KeyDecisionRow){
+		"request ref":    func(r *KeyDecisionRow) { r.RequestRef = kdQ[:70] },
+		"upper-case ref": func(r *KeyDecisionRow) { r.PacketRef = "sha256:" + strings.ToUpper(kdP[7:]) },
+		"verb":           func(r *KeyDecisionRow) { r.Verb = "maybe" },
+		"approver":       func(r *KeyDecisionRow) { r.Approver = "ssh:mallory" },
+		"format":         func(r *KeyDecisionRow) { r.Format = "webauthn" },
+		"empty artifact": func(r *KeyDecisionRow) { r.Artifact = nil },
+		"large artifact": func(r *KeyDecisionRow) { r.Artifact = make([]byte, 16385) },
+	} {
+		row := kdRow()
+		edit(&row)
+		if err := c.RecordKeyDecision(testCtx(t), row, kdUse, 1000); !errors.Is(err, ErrBadKeyDecision) {
+			t.Errorf("%s: %v, want ErrBadKeyDecision", name, err)
+		}
+	}
+	if err := c.LandKeyDecision(testCtx(t), kdQ, "sha256:x"); !errors.Is(err, ErrBadKeyDecision) {
+		t.Errorf("land a bad decision ref: %v", err)
+	}
+	if _, ok := kdRead(t, c); ok {
+		t.Fatal("a refused row was written")
+	}
+}
````

### Patch M4g — the decision judge (JD1–JD3)

````diff M4g
diff --git a/host/authority/judge_decision.go b/host/authority/judge_decision.go
new file mode 100644
index 0000000..f24c36f
--- /dev/null
+++ b/host/authority/judge_decision.go
@@ -0,0 +1,45 @@
+package authority
+
+// Row 201 A4 sprint 2 (design §2.7): the decision judge beside the mint judge (judge.go, unchanged),
+// on a read-only handle, with judge.go's step names, over a decision's subject rather than a
+// MintRequestRow. The passkey judge (W15 deferred to the writer) is sprint 3's (plan S2-D1).
+
+import (
+	"context"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// DecisionJudgeInput is a decision's subject at now: the verb byte, the packet digest P, the
+// request digest Q and the signed expiry.
+type DecisionJudgeInput struct {
+	Verb            byte
+	Packet, Request [32]byte
+	Expiry, Now     int64
+}
+
+// JudgeDecisionSSH judges an armored SSHSIG decision: ParseSSHSIG, the approver its key names
+// (ApproverByKey on ro), then approval.DecideBySignature. It returns the verified decision or a
+// *JudgeError.
+func JudgeDecisionSSH(ctx context.Context, ro *store.Credentials, in DecisionJudgeInput, armored []byte) (*approval.SignedDecision, error) {
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
+	signed, err := approval.DecideBySignature(approval.DecisionSignatureInput{Verb: in.Verb, Packet: in.Packet, Request: in.Request,
+		Expiry: in.Expiry, Now: in.Now, Signature: armored,
+		Approver: approval.Approver{Name: enrolled.Name, Format: enrolled.Format, Key: approval.SSHPublicKey{Blob: enrolled.PublicKey},
+			EnrolledAt: enrolled.EnrolledAt, LastUsedAt: enrolled.LastUsedAt, RevokedAt: enrolled.RevokedAt}})
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
+	}
+	return signed, nil
+}
diff --git a/host/authority/judge_decision_test.go b/host/authority/judge_decision_test.go
new file mode 100644
index 0000000..18d127b
--- /dev/null
+++ b/host/authority/judge_decision_test.go
@@ -0,0 +1,77 @@
+package authority
+
+// Row 201 A4 sprint 2 M4f (design §2.7, §6 JD1-JD3, SSH half): JudgeDecisionSSH's steps, causes and
+// fields agree with the mint judge's (JudgeSSHSignature, JG2): parse, store, unknown key, law; a
+// verified decision binds its subject. Every key is approvaltest's (D-WORLD-109).
+
+import (
+	"errors"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+var jdIn = DecisionJudgeInput{Verb: approval.VerbDeny, Packet: [32]byte{0x50}, Request: [32]byte{0x51}, Expiry: 1760, Now: 1710}
+
+// jdSign is s's SSHSIG over in's decision message under namespace.
+func jdSign(t *testing.T, s *approvaltest.Signer, namespace string, in DecisionJudgeInput) []byte {
+	t.Helper()
+	msg, err := approval.DecisionMessage(in.Verb, in.Packet, in.Request, in.Expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return s.Sign(namespace, msg)
+}
+
+func TestJudgeDecisionSSH(t *testing.T) {
+	good := jdSign(t, jgSigner, approval.DecisionNamespace, jdIn)
+	stranger := approvaltest.NewSigner(65)
+	mintMsg, _ := approval.SignedMessage(jdIn.Request, jdIn.Packet, jdIn.Expiry)
+	for _, tc := range []struct {
+		name     string
+		setup    []string
+		artifact []byte
+		in       DecisionJudgeInput
+		step     string
+		cause    error
+		field    string
+	}{
+		{name: "JD-parse", artifact: []byte("not a signature\n"), in: jdIn, step: JudgeStepParse, cause: approval.ErrArmor},
+		{name: "JD-store-v4", setup: []string{`DROP TABLE key_decisions`, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`},
+			artifact: good, in: jdIn, step: JudgeStepStore},
+		{name: "JD-unknown-key", artifact: jdSign(t, stranger, approval.DecisionNamespace, jdIn), in: jdIn, step: JudgeStepUnknownKey, cause: ErrNotEnrolled,
+			field: stranger.PublicKey().Fingerprint()},
+		{name: "JD-law-mint-signature", artifact: jgSigner.Sign(approval.Namespace, mintMsg), in: jdIn, step: JudgeStepLaw, cause: approval.ErrNamespace},
+		{name: "JD-law-other-verb", artifact: good, in: DecisionJudgeInput{Verb: approval.VerbApprove, Packet: jdIn.Packet, Request: jdIn.Request, Expiry: jdIn.Expiry, Now: jdIn.Now},
+			step: JudgeStepLaw, cause: approval.ErrBadSignature},
+		{name: "JD-law-format", setup: []string{`UPDATE approvers SET format = 'webauthn' WHERE name = 'mark-ssh'`}, artifact: good, in: jdIn, step: JudgeStepLaw},
+		{name: "JD-law-revoked", setup: []string{`UPDATE approvers SET revoked_at = 950 WHERE name = 'mark-ssh'`}, artifact: good, in: jdIn, step: JudgeStepLaw, cause: approval.ErrApproverRevoked},
+		{name: "JD-law-clock-behind", artifact: good, in: DecisionJudgeInput{Verb: jdIn.Verb, Packet: jdIn.Packet, Request: jdIn.Request, Expiry: jdIn.Expiry, Now: 899},
+			step: JudgeStepLaw, cause: approval.ErrApproverClockBehind},
+	} {
+		path, _ := jgRig(t)
+		jgExec(t, path, tc.setup...)
+		ro, err := store.OpenCredentialsReadOnly(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		d, err := JudgeDecisionSSH(boundedTestContext(t), ro, tc.in, tc.artifact)
+		_ = ro.Close()
+		var je *JudgeError
+		if d != nil || !errors.As(err, &je) || je.Step != tc.step || (tc.cause != nil && !errors.Is(err, tc.cause)) || je.Fingerprint != tc.field {
+			t.Errorf("%s: %v, %v; want step %s cause %v field %q", tc.name, d, err, tc.step, tc.cause, tc.field)
+		}
+	}
+	path, _ := jgRig(t)
+	ro, err := store.OpenCredentialsReadOnly(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	d, err := JudgeDecisionSSH(boundedTestContext(t), ro, jdIn, good)
+	if err != nil || d.DecidedBy() != "ssh:mark-ssh" || d.Binds(jdIn.Verb, jdIn.Packet, jdIn.Request, jdIn.Expiry) != nil {
+		t.Fatalf("control: %v, %v", d, err)
+	}
+}
````

### Patch M5a — the POST's sshsig arm (RD10, RD11, the judge's refusals)

````diff M5a
diff --git a/host/daemon/surface_approve_post_guard_test.go b/host/daemon/surface_approve_post_guard_test.go
index 6c04f38..4003f68 100644
--- a/host/daemon/surface_approve_post_guard_test.go
+++ b/host/daemon/surface_approve_post_guard_test.go
@@ -165,7 +165,7 @@ func TestDaemonCredentialWriterSites(t *testing.T) {
 			})
 		}
 	}
-	if want := map[string]string{"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials"}; !reflect.DeepEqual(sites, want) || readOnly < 5 {
+	if want := map[string]string{"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials", "keyArm": "OpenCredentials"}; !reflect.DeepEqual(sites, want) || readOnly < 5 {
 		t.Fatalf("writer sites %v (read-only opens %d), want exactly %v and at least 5 read-only opens", sites, readOnly, want)
 	}
 }
diff --git a/host/daemon/surface_decide.go b/host/daemon/surface_decide.go
index 7f2ccf4..df32a1b 100644
--- a/host/daemon/surface_decide.go
+++ b/host/daemon/surface_decide.go
@@ -176,15 +176,29 @@ func (d *Daemon) handleSurfaceDecide(w http.ResponseWriter, r *http.Request, o [
 		return
 	}
 	q := hashref.MustParse(ref) // D1 held
-	// 2. the arm (D11). The grammar admits typed only in sprint 1; the phase carries the verb.
-	if err := d.consumeOperatorProof("decision-"+q.Digest(), "decide-"+body.Verb, r.Header.Get("X-World-Operator-Proof")); err != nil {
+	// 2. the arm (D11): a key's judge and evidence (sprint 2, surface_decide_key.go), or the typed
+	// proof, whose phase carries the verb.
+	decidedBy, evidence := typedDecider, (*store.Credentials)(nil)
+	if body.SSHSIG != nil {
+		var ok bool
+		if evidence, decidedBy, ok = d.keyArm(ctx, w, r, f, body, now); !ok {
+			return
+		}
+		defer func() { _ = evidence.Close() }()
+	} else if err := d.consumeOperatorProof("decision-"+q.Digest(), "decide-"+body.Verb, r.Header.Get("X-World-Operator-Proof")); err != nil {
 		writeApprovalRefusal(w, http.StatusForbidden, "Forbidden", "the typed decision is not proven: "+err.Error()+"; decide at a terminal with ailang-worldd approve "+ref)
 		return
 	}
-	decision, err := broker.DecideApproval(ctx, d.store, q, body.Verb, typedDecider, now) // 3.
+	decision, err := broker.DecideApproval(ctx, d.store, q, body.Verb, decidedBy, now) // 3.
 	if err != nil {
 		d.approvalInternal(w, r, err)
 		return
 	}
-	writeApprovalJSON(w, http.StatusOK, surfacepage.DecisionMade{V: 1, Decided: body.Verb, Request: ref, Decision: decision.String(), DecidedBy: typedDecider})
+	if evidence != nil { // 4. the key decision's evidence names the decision that landed
+		if err := evidence.LandKeyDecision(ctx, ref, decision.String()); err != nil {
+			d.approvalInternal(w, r, err)
+			return
+		}
+	}
+	writeApprovalJSON(w, http.StatusOK, surfacepage.DecisionMade{V: 1, Decided: body.Verb, Request: ref, Decision: decision.String(), DecidedBy: decidedBy})
 }
diff --git a/host/daemon/surface_decide_key.go b/host/daemon/surface_decide_key.go
new file mode 100644
index 0000000..29a7c0a
--- /dev/null
+++ b/host/daemon/surface_decide_key.go
@@ -0,0 +1,70 @@
+package daemon
+
+// Row 201 A4 sprint 2 (design §2.5 D11, §2.7): the POST's key arm, in sprint 1's arm slot (S1-D9):
+// after the law, under humanMu, the decision judge on a read-only handle, Binds, then a non-waiting
+// writer for RecordKeyDecision (K0-K5). SSHSIG only; the passkey arm is sprint 3's (plan S2-D1).
+
+import (
+	"context"
+	"encoding/hex"
+	"errors"
+	"net/http"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// isApproverRace is a key decision's approver that changed since the read-only judge: revoked,
+// re-keyed, gone or enrolled after now (K3, in the writer's transaction).
+func isApproverRace(err error) bool {
+	return errors.Is(err, store.ErrApproverRevoked) || errors.Is(err, store.ErrApproverKeyChanged) ||
+		errors.Is(err, store.ErrApproverUnknown) || errors.Is(err, store.ErrApproverClockBehind)
+}
+
+// keyArm is D11 for a key format. It returns the open credentials writer (the caller closes it
+// after LandKeyDecision) and the decidedBy the verified key derives, or ok false when it has
+// answered a refusal: nothing is decided, and nothing is recorded unless RecordKeyDecision committed.
+func (d *Daemon) keyArm(ctx context.Context, w http.ResponseWriter, r *http.Request, f broker.DecisionFacts, body surfacepage.DecideBody, now int64) (*store.Credentials, string, bool) {
+	verb, _ := approval.DecisionVerb(body.Verb) // D9 held
+	artifact := []byte(*body.SSHSIG)
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return nil, "", false
+	}
+	signed, err := authority.JudgeDecisionSSH(ctx, ro, authority.DecisionJudgeInput{Verb: verb, Packet: f.Packet, Request: f.Request, Expiry: body.Expiry, Now: now}, artifact)
+	_ = ro.Close()
+	if err == nil {
+		err = signed.Binds(verb, f.Packet, f.Request, body.Expiry)
+	}
+	var je *authority.JudgeError
+	switch {
+	case errors.As(err, &je) && je.Step == authority.JudgeStepStore:
+		d.approvalStoreError(w, r, err)
+		return nil, "", false
+	case err != nil:
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the SSH signature does not decide request "+f.Ref+": "+err.Error())
+		return nil, "", false
+	}
+	wr, err := store.OpenCredentials(d.credentialsPath) // non-waiting: a busy writer is 503 Busy
+	if err != nil {
+		d.approvalStoreError(w, r, err)
+		return nil, "", false
+	}
+	name, format, key, _ := signed.Use()
+	err = wr.RecordKeyDecision(ctx, store.KeyDecisionRow{RequestRef: f.Ref, PacketRef: "sha256:" + hex.EncodeToString(f.Packet[:]), Verb: body.Verb,
+		Format: format, Approver: signed.DecidedBy(), Artifact: artifact, Expiry: body.Expiry}, store.ApproverUse{Name: name, Format: format, PublicKey: key}, now)
+	if err != nil {
+		_ = wr.Close()
+		if isApproverRace(err) {
+			writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the approver changed since the signature was judged: "+err.Error())
+		} else {
+			d.approvalStoreError(w, r, err)
+		}
+		return nil, "", false
+	}
+	return wr, signed.DecidedBy(), true
+}
diff --git a/host/daemon/surface_decide_key_test.go b/host/daemon/surface_decide_key_test.go
new file mode 100644
index 0000000..0218dc8
--- /dev/null
+++ b/host/daemon/surface_decide_key_test.go
@@ -0,0 +1,148 @@
+package daemon
+
+// Row 201 A4 sprint 2 M5a (design §2.5 D11, §2.7, §6 RD10, RD11 and the judge's refusals; the
+// passkey halves are sprint 3's): POST /surface/decide/{ref} with an SSHSIG by an enrolled
+// approvaltest key (D-WORLD-109), through d.Handler().
+
+import (
+	"context"
+	"encoding/hex"
+	"encoding/json"
+	"fmt"
+	"net/http/httptest"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+var kdSigner = approvaltest.NewSigner(71)
+
+// enrolDecider enrols s as name at packetNow, before the rig's clock (decNow).
+func (r decideRig) enrolDecider(name string, s *approvaltest.Signer) {
+	r.t.Helper()
+	r.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.EnrolApprover(ctx, store.ApproverRow{Name: name, Format: approval.FormatSSH, PublicKey: s.PublicKey().Blob, EnrolledAt: packetNow})
+	})
+}
+
+// decideSSHBody is a POST body deciding verb on (p, q, expiry) with s's SSHSIG under namespace.
+func decideSSHBody(t *testing.T, s *approvaltest.Signer, namespace, verb string, p store.Object, q hashref.HashRef, expiry int64) []byte {
+	t.Helper()
+	v, err := approval.DecisionVerb(verb)
+	if err != nil {
+		v = approval.VerbApprove
+	}
+	msg, err := approval.DecisionMessage(v, [32]byte(mustDigest(t, p.Hash)), [32]byte(mustDigest(t, q)), expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	b, _ := json.Marshal(map[string]any{"v": 1, "verb": verb, "packet": p.Hash.Digest(), "request": q.Digest(), "expiry": expiry, "sshsig": string(s.Sign(namespace, msg))})
+	return b
+}
+
+func mustDigest(t *testing.T, h hashref.HashRef) []byte {
+	t.Helper()
+	d, err := hex.DecodeString(h.Digest())
+	if err != nil || len(d) != 32 {
+		t.Fatal(h, err)
+	}
+	return d
+}
+
+// keyRow is q's key_decisions row, read on a fresh read-only handle.
+func (r decideRig) keyRow(q hashref.HashRef) (store.KeyDecisionRow, bool) {
+	r.t.Helper()
+	ro, err := store.OpenCredentialsReadOnly(r.d.credentialsPath)
+	if err != nil {
+		r.t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	row, ok, err := ro.KeyDecision(boundedTestContext(r.t), q.String())
+	if err != nil {
+		r.t.Fatal(err)
+	}
+	return row, ok
+}
+
+// postKey is a key POST: the loopback headers and no proof.
+func (r decideRig) postKey(q hashref.HashRef, body []byte) *httptest.ResponseRecorder {
+	return r.post(q.String(), body, "")
+}
+
+// RD10, RD11: an enrolled key's SSHSIG decides approve and deny as ssh:<name>, derived from the
+// verified key (a body naming decidedBy is 400); the decision is on the chain and the evidence row
+// is landed with the decision's ref.
+func TestDecidePOSTSSHSIG(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolDecider("mark-ssh", kdSigner)
+	for _, verb := range []string{"approve", "deny"} {
+		q, p := r.packet("ssh-"+verb, broker.WorkspaceScope, "agent-"+verb)
+		forged := strings.Replace(string(decideSSHBody(t, kdSigner, approval.DecisionNamespace, verb, p, q, decNow+50)), `"v":1`, `"v":1,"decidedBy":"ssh:root"`, 1)
+		requireRefusal(t, "RD11 decidedBy", r.postKey(q, []byte(forged)), 400, "BadRequest")
+		rec := r.postKey(q, decideSSHBody(t, kdSigner, approval.DecisionNamespace, verb, p, q, decNow+50))
+		var made surfacepage.DecisionMade
+		if err := json.Unmarshal(rec.Body.Bytes(), &made); err != nil || rec.Code != 200 || made.Decided != verb || made.DecidedBy != "ssh:mark-ssh" {
+			t.Fatalf("%s: %d %s", verb, rec.Code, rec.Body.String())
+		}
+		_, decided, gotVerb, by, err := broker.ApprovalChainState(boundedTestContext(t), r.d.store, q)
+		row, ok := r.keyRow(q)
+		if err != nil || !decided || gotVerb != verb || by != "ssh:mark-ssh" || !ok || row.DecisionRef != made.Decision || row.Verb != verb || !strings.HasPrefix(string(row.Artifact), "-----BEGIN SSH SIGNATURE-----") ||
+			row.PacketRef != p.Hash.String() || row.Approver != "ssh:mark-ssh" || row.Expiry != decNow+50 || row.DecidedAt != decNow {
+			t.Fatalf("%s: chain %v %q %q %v; row %+v %v; made %+v", verb, decided, gotVerb, by, err, row, ok, made)
+		}
+	}
+}
+
+// D11 (sshsig): the judge's refusals are 422 with nothing decided and nothing recorded: an
+// unenrolled key, a revoked approver, a mint-namespace signature, an approve signature posted as
+// deny, a signature over another packet, another expiry, and garbage; and the law runs first: a
+// valid signature with an expiry past the window is D10's 409.
+func TestDecidePOSTSSHSIGRefusals(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolDecider("mark-ssh", kdSigner)
+	r.enrolDecider("gone", approvaltest.NewSigner(72))
+	r.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.RevokeApprover(ctx, "gone", packetNow+1)
+	})
+	q, p := r.packet("ssh-refusals", broker.WorkspaceScope, "agent")
+	_, other := r.packet("ssh-other", broker.WorkspaceScope, "other")
+	_, head := r.decisions()
+	good := decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, decNow+50)
+	for _, tc := range []struct {
+		name   string
+		body   []byte
+		status int
+		reason string
+	}{
+		{"unenrolled", decideSSHBody(t, approvaltest.NewSigner(73), approval.DecisionNamespace, "approve", p, q, decNow+50), 422, "no enrolled approver holds that key"},
+		{"revoked", decideSSHBody(t, approvaltest.NewSigner(72), approval.DecisionNamespace, "approve", p, q, decNow+50), 422, "the approver is revoked"},
+		{"mint namespace", decideSSHBody(t, kdSigner, approval.Namespace, "approve", p, q, decNow+50), 422, "wrong namespace"},
+		{"verb", []byte(strings.Replace(string(good), `"verb":"approve"`, `"verb":"deny"`, 1)), 422, "signature does not verify"},
+		{"another packet", decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", other, q, decNow+50), 409, "reload it"},
+		{"another expiry", []byte(strings.Replace(string(good), fmt.Sprintf(`"expiry":%d`, decNow+50), fmt.Sprintf(`"expiry":%d`, decNow+49), 1)), 422, "signature does not verify"},
+		{"garbage", []byte(strings.Replace(string(good), `"sshsig":"-----BEGIN`, `"sshsig":"x-----BEGIN`, 1)), 422, "bad armor"},
+		{"D10 first", decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, packetNow+61), 409, "stale or too far ahead"},
+	} {
+		got := requireRefusal(t, tc.name, r.postKey(q, tc.body), tc.status, map[int]string{422: "NotApproved", 409: "Conflict"}[tc.status])
+		if !strings.Contains(got.Reason, tc.reason) {
+			t.Errorf("%s: %q", tc.name, got.Reason)
+		}
+		if n, h := r.decisions(); n != 0 || h != head {
+			t.Fatalf("%s: %d decisions", tc.name, n)
+		}
+		if _, ok := r.keyRow(q); ok {
+			t.Fatalf("%s: evidence was recorded", tc.name)
+		}
+	}
+	for _, drop := range [][]string{{`DROP TABLE key_decisions`, `PRAGMA user_version = 7`}, // v7: the writer's K1, then v4: the judge's store step
+		{`DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`}} {
+		setVersion(t, r.d.credentialsPath, drop...)
+		requireRefusal(t, drop[len(drop)-1], r.postKey(q, good), 503, "Restart")
+	}
+}
diff --git a/host/surfacepage/decide_model.go b/host/surfacepage/decide_model.go
index f6d2135..ffada1d 100644
--- a/host/surfacepage/decide_model.go
+++ b/host/surfacepage/decide_model.go
@@ -1,8 +1,8 @@
 package surfacepage
 
 // Row 201 A4 sprint 1 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decision
-// routes' JSON shapes and the POST body's grammar, pure. Sprint 1 decides with the typed arm only;
-// the key formats (sshsig, webauthn) arrive with sprint 2, so until then they are unknown fields.
+// routes' JSON shapes and the POST body's grammar, pure. Sprint 2 admits exactly one of typed and
+// sshsig; webauthn arrives with sprint 3's passkey arm, so until then it is an unknown field.
 
 import (
 	"bytes"
@@ -52,23 +52,25 @@ type DecisionMade struct {
 	DecidedBy string `json:"decidedBy"`
 }
 
-// DecideBody is the POST body (design §2.8), sprint 1's typed arm: the verb, the packet digest P
-// and request digest Q the caller saw (64 lowercase hex), the expiry, and typed: true.
+// DecideBody is the POST body (design §2.8): the verb, the packet digest P and request digest Q
+// the caller saw (64 lowercase hex), the expiry, and exactly one format: typed: true (sprint 1) or
+// sshsig, an armored SSH signature over the decision message (sprint 2).
 type DecideBody struct {
-	V       int    `json:"v"`
-	Verb    string `json:"verb"`
-	Packet  string `json:"packet"`
-	Request string `json:"request"`
-	Expiry  int64  `json:"expiry"`
-	Typed   bool   `json:"typed"`
+	V       int     `json:"v"`
+	Verb    string  `json:"verb"`
+	Packet  string  `json:"packet"`
+	Request string  `json:"request"`
+	Expiry  int64   `json:"expiry"`
+	Typed   bool    `json:"typed,omitempty"`
+	SSHSIG  *string `json:"sshsig,omitempty"`
 }
 
 // ErrDecideBody refuses a body that is not DecideBody's grammar.
-var ErrDecideBody = errors.New(`the body is not {"v":1,"verb":…,"packet":<64 hex>,"request":<64 hex>,"expiry":…,"typed":true}`)
+var ErrDecideBody = errors.New(`the body is not {"v":1,"verb":…,"packet":<64 hex>,"request":<64 hex>,"expiry":…} with exactly one of "typed":true and "sshsig":<armored SSHSIG>`)
 
-// ParseDecideBody is the POST's grammar: strict JSON (no unknown field, so no decidedBy and no key
-// format yet; exactly one value), v 1, typed true, both digests 64 lowercase hex. The verb and the
-// expiry are the law's (D9, D10), not the grammar's.
+// ParseDecideBody is the POST's grammar: strict JSON (no unknown field, so no decidedBy; exactly
+// one value), v 1, exactly one format, both digests 64 lowercase hex. The verb and the expiry are
+// the law's (D9, D10), not the grammar's.
 func ParseDecideBody(b []byte) (body DecideBody, packet, request [32]byte, err error) {
 	dec := json.NewDecoder(bytes.NewReader(b))
 	dec.DisallowUnknownFields()
@@ -81,8 +83,10 @@ func ParseDecideBody(b []byte) (body DecideBody, packet, request [32]byte, err e
 	switch {
 	case body.V != 1:
 		return DecideBody{}, packet, request, fmt.Errorf("%w: v is %d", ErrDecideBody, body.V)
-	case !body.Typed:
-		return DecideBody{}, packet, request, fmt.Errorf("%w: it has no format (this build decides with typed only)", ErrDecideBody)
+	case !body.Typed && body.SSHSIG == nil:
+		return DecideBody{}, packet, request, fmt.Errorf("%w: it has no format", ErrDecideBody)
+	case body.Typed && body.SSHSIG != nil:
+		return DecideBody{}, packet, request, fmt.Errorf("%w: it has two formats", ErrDecideBody)
 	case !lowerHex64(body.Packet) || !lowerHex64(body.Request):
 		return DecideBody{}, packet, request, fmt.Errorf("%w: packet and request are each 64 lowercase hex", ErrDecideBody)
 	}
````

### Patch M5b — busy, retry, R4 by key, the lock, the discovery by key (RD12, RD13 a, R4e, H-LOCK)

````diff M5b
diff --git a/host/daemon/surface_decide_discovery_test.go b/host/daemon/surface_decide_discovery_test.go
index cdccc0c..b02ca9b 100644
--- a/host/daemon/surface_decide_discovery_test.go
+++ b/host/daemon/surface_decide_discovery_test.go
@@ -8,8 +8,11 @@ import (
 	"testing"
 	"time"
 
+	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/broker"
 	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
 )
 
 // RD9: two invocations whose requests collide (one ref) and two packets → 409 D4b, 0 decisions.
@@ -48,6 +51,19 @@ func TestDecideDiscoveryFindsAHiddenSecondPacket(t *testing.T) {
 		t.Errorf("hidden: %q", got.Reason)
 	}
 	r.requireNothingDecided("hidden", q, head)
+	r.requireKeyRefused("hidden", q, p, head) // sprint 2: the SSHSIG arm too
+}
+
+// requireKeyRefused fails unless a valid SSHSIG decision of q is 409 with nothing decided and no
+// key_decisions row: the discovery refuses before the key arm runs.
+func (r decideRig) requireKeyRefused(what string, q hashref.HashRef, p store.Object, head hashref.HashRef) {
+	r.t.Helper()
+	r.enrolDecider("mark-ssh", kdSigner)
+	requireRefusal(r.t, what+" ssh", r.postKey(q, decideSSHBody(r.t, kdSigner, approval.DecisionNamespace, "approve", p, q, decNow+50)), 409, "Conflict")
+	if _, ok := r.keyRow(q); ok {
+		r.t.Fatalf("%s ssh: evidence was recorded", what)
+	}
+	r.requireNothingDecided(what+" ssh", q, head)
 }
 
 // RD9c: 2,001 packets, exactly one carrying the request → 409 D4a (the overflow sentinel), 0
@@ -63,6 +79,7 @@ func TestDecideDiscoveryRefusesAnOverflow(t *testing.T) {
 		t.Errorf("overflow: %q", got.Reason)
 	}
 	r.requireNothingDecided("overflow", q, head)
+	r.requireKeyRefused("overflow", q, p, head) // sprint 2: the SSHSIG arm too
 }
 
 // RD9d (must live): 150 packets, exactly one carrying the request, sorting after the first 101 →
diff --git a/host/daemon/surface_decide_key_retry_test.go b/host/daemon/surface_decide_key_retry_test.go
new file mode 100644
index 0000000..ee882f4
--- /dev/null
+++ b/host/daemon/surface_decide_key_retry_test.go
@@ -0,0 +1,129 @@
+package daemon
+
+// Row 201 A4 sprint 2 M5b (design §2.7, §6 RD12, RD13 (a), R4e; sprint 1's carried H-LOCK): the key
+// arm's busy writer, the retry after a failed world write, the reserved scopes by key, and humanMu
+// under concurrent POSTs. The rig's key helpers are surface_decide_key_test.go's.
+
+import (
+	"database/sql"
+	"encoding/json"
+	"strings"
+	"sync"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// RD12: a held credentials writer is 503 Busy with Retry-After, nothing decided or recorded; once
+// it is released the same body decides.
+func TestDecidePOSTBusyWriter(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolDecider("mark-ssh", kdSigner)
+	q, p := r.packet("busy", broker.WorkspaceScope, "agent")
+	body := decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, decNow+50)
+	held, err := store.OpenCredentials(r.d.credentialsPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	rec := r.postKey(q, body)
+	requireRefusal(t, "busy", rec, 503, "Busy")
+	if n, _ := r.decisions(); n != 0 || rec.Header().Get("Retry-After") != "1" {
+		t.Fatalf("busy: %d decisions, Retry-After %q", n, rec.Header().Get("Retry-After"))
+	}
+	_ = held.Close()
+	if _, ok := r.keyRow(q); ok {
+		t.Fatal("busy: evidence was recorded")
+	}
+	if rec := r.postKey(q, body); rec.Code != 200 {
+		t.Fatalf("released: %d %s", rec.Code, rec.Body.String())
+	}
+}
+
+// RD13 (a): the world write fails after RecordKeyDecision (a trigger refuses the decision object)
+// → 500, one unlanded row, 0 decisions; the identical body retried → 200, one decision, the row
+// landed and kept as first recorded (K0's exact match); a replay after that is D6's 409.
+func TestDecidePOSTRetryAfterAFailedWorldWrite(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolDecider("mark-ssh", kdSigner)
+	q, p := r.packet("retry", broker.WorkspaceScope, "agent")
+	tq, tp := r.packet("retry-typed", broker.WorkspaceScope, "typed")
+	body := decideSSHBody(t, kdSigner, approval.DecisionNamespace, "deny", p, q, decNow+50)
+	db, err := sql.Open("sqlite", r.d.cfg.DBPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	if _, err := db.Exec(`CREATE TRIGGER fail_decision BEFORE INSERT ON objects WHEN NEW.semantic_id = '` + broker.ApprovalDecisionV1 + `' BEGIN SELECT RAISE(ABORT, 'injected world-write failure'); END`); err != nil {
+		t.Fatal(err)
+	}
+	requireRefusal(t, "failed world write", r.postKey(q, body), 500, "Internal")
+	requireRefusal(t, "failed world write, typed", r.post(tq.String(), decideBody("approve", tp, tq, decNow+50), r.proof(tq, "decide-approve")), 500, "Internal")
+	first, ok := r.keyRow(q)
+	if n, _ := r.decisions(); n != 0 || !ok || first.DecisionRef != "" {
+		t.Fatalf("after the failure: %d decisions, row %+v %v", n, first, ok)
+	}
+	if _, err := db.Exec(`DROP TRIGGER fail_decision`); err != nil {
+		t.Fatal(err)
+	}
+	r.d.approveNow = func() int64 { return decNow + 5 }
+	rec := r.postKey(q, body)
+	var made surfacepage.DecisionMade
+	_ = json.Unmarshal(rec.Body.Bytes(), &made)
+	row, _ := r.keyRow(q)
+	if n, _ := r.decisions(); rec.Code != 200 || n != 1 || row.DecisionRef != made.Decision || row.DecidedAt != first.DecidedAt {
+		t.Fatalf("retry: %d %s; %d decisions; row %+v (first %+v)", rec.Code, rec.Body.String(), n, row, first)
+	}
+	requireRefusal(t, "replay", r.postKey(q, body), 409, "Conflict")
+}
+
+// R4e: publish- and finish-shaped requests with a valid key are refused by L1 (403 with the
+// owner), nothing decided or recorded. With L1 compiled out the key decision would land; the
+// publish consumer's L2 refuses it (R4b), the must-die pair the mutant ledger records (LAW-D5).
+func TestDecideRefusesReservedScopesByKey(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolDecider("mark-ssh", kdSigner)
+	for scope, owner := range map[string]string{"registry:pkg#effect=Registry.Publish": "world-publish approve", "github-pr:o/r#x": "ailang-worldd episode finish"} {
+		q, p := r.packet("reserved-"+owner, scope, "agent")
+		got := requireRefusal(t, scope, r.postKey(q, decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, decNow+50)), 403, "Forbidden")
+		if !strings.Contains(got.Reason, owner) {
+			t.Errorf("%s: %q", scope, got.Reason)
+		}
+		if _, ok := r.keyRow(q); ok {
+			t.Fatalf("%s: evidence was recorded", scope)
+		}
+	}
+	if n, _ := r.decisions(); n != 0 {
+		t.Fatalf("%d decisions", n)
+	}
+}
+
+// H-LOCK (carried from sprint 1): eight concurrent POSTs of one valid key decision are serialized
+// by humanMu: exactly one decides and every other sees it decided (D6's 409), never a busy writer
+// or a refused second write.
+func TestDecideConcurrentPOSTsDecideOnce(t *testing.T) {
+	r := newDecideRig(t)
+	r.enrolDecider("mark-ssh", kdSigner)
+	q, p := r.packet("concurrent", broker.WorkspaceScope, "agent")
+	body := decideSSHBody(t, kdSigner, approval.DecisionNamespace, "approve", p, q, decNow+50)
+	codes := make([]int, 8)
+	var wg sync.WaitGroup
+	for i := range codes {
+		wg.Add(1)
+		go func(i int) { defer wg.Done(); codes[i] = r.postKey(q, body).Code }(i)
+	}
+	wg.Wait()
+	ok := 0
+	for _, c := range codes {
+		if c == 200 {
+			ok++
+		} else if c != 409 {
+			t.Errorf("a concurrent POST answered %d (codes %v)", c, codes)
+		}
+	}
+	if n, _ := r.decisions(); ok != 1 || n != 1 {
+		t.Fatalf("%d decided, %d decisions (codes %v)", ok, n, codes)
+	}
+}
````

### Patch M5c — --sign-with on a request ref (CD5 and its real-tool sibling)

````diff M5c
diff --git a/cmd/ailang-worldd/approve.go b/cmd/ailang-worldd/approve.go
index debe434..36b9662 100644
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -30,6 +30,7 @@ const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--
        ailang-worldd approve <request-id> --await [--await-for <duration>] [--db <path>] [--credentials-db <path>]
        ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-ref> [--deny] [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-ref> [--deny] --sign-with <key> [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
 
 ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
 201), after a y/N typed at /dev/tty, or, with --signature or --sign-with,
@@ -101,7 +102,9 @@ installed service's address), checks that both hash to their refs, shows every
 field, asks y at /dev/tty, and has the daemon record the decision (decided by
 operator:terminal). A request is decided once. --await is refused (the daemon
 decides at once); publish and finish approvals stay with world-publish
-approve and episode finish.
+approve and episode finish. --sign-with <key> opens no terminal: ssh-keygen
+-Y sign -n world-decision@sunholo signs the 91-byte decision message, and the
+decision is by ssh:<name>, an approver enrolled with 'approver enroll'.
 `
 
 // approveAdmission is approve's admission of a stored row, before any terminal: the decoded
@@ -212,14 +215,14 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 	if !*list {
 		id = positional[0]
 		if decisionRefPattern.MatchString(id) { // row 201 A4: a request ref selects the decision mode (design §2.10)
-			if modes != 0 { // CD3: --await, and the key flags until they decide a request ref
-				return refuse("a request ref is decided at once by the daemon: there is nothing to await, and no key flag decides it yet")
+			if modes != 0 && *signWith == "" { // CD3: --sign-with is the key flag that decides a request ref (plan S2-D1)
+				return refuse("a request ref is decided at once by the daemon, by the typed y or --sign-with <key>: there is nothing to await, and the other key flags do not decide it")
 			}
 			verb := "approve"
 			if *deny {
 				verb = "deny"
 			}
-			return runApproveDecide(decideOptions{ref: id, verb: verb, db: *db, credsDB: *credsDB, daemon: *daemonAddr}, stdout, stderr, env, refuse)
+			return runApproveDecide(decideOptions{ref: id, verb: verb, db: *db, credsDB: *credsDB, daemon: *daemonAddr, signWith: *signWith}, stdout, stderr, env, refuse)
 		}
 		if !store.MintRequestIDPattern.MatchString(id) {
 			return refuse("%q is not a mint request id (mr- and 16 lowercase hex)", id)
diff --git a/cmd/ailang-worldd/approve_decide.go b/cmd/ailang-worldd/approve_decide.go
index 91e38b0..77fcc2a 100644
--- a/cmd/ailang-worldd/approve_decide.go
+++ b/cmd/ailang-worldd/approve_decide.go
@@ -13,6 +13,7 @@ import (
 	"bytes"
 	"context"
 	"encoding/base64"
+	"encoding/hex"
 	"encoding/json"
 	"fmt"
 	"io"
@@ -21,8 +22,10 @@ import (
 	"strings"
 	"time"
 
+	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/broker"
 	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/sshsign"
 	"github.com/sunholo-data/ailang-world/host/store"
 	"github.com/sunholo-data/ailang-world/host/surfacepage"
 )
@@ -35,7 +38,7 @@ const approveDecideBudget = 10 * time.Second
 
 // decideOptions are approve's flags in the decision mode.
 type decideOptions struct {
-	ref, verb, db, credsDB, daemon string
+	ref, verb, db, credsDB, daemon, signWith string
 }
 
 // decisionSubject is a subject the CLI checked: its request and packet, decoded from bytes that
@@ -73,6 +76,12 @@ func checkDecisionSubject(ref string, body []byte) (decisionSubject, string) {
 	return s, ""
 }
 
+// refDigest is the 32 bytes of a sha256: ref the subject check has held to its bytes.
+func refDigest(ref string) [32]byte {
+	d, _ := hex.DecodeString(strings.TrimPrefix(ref, "sha256:"))
+	return [32]byte(d)
+}
+
 // sha256Ref is the sha256: ref of b.
 func sha256Ref(b []byte) string { return "sha256:" + finishSHA(b) }
 
@@ -137,6 +146,9 @@ func runApproveDecide(o decideOptions, stdout, stderr io.Writer, env sessionEnv,
 	if s.Refusal != "" {
 		return refuse("request %s cannot be decided now: %s", o.ref, s.Refusal)
 	}
+	if o.signWith != "" {
+		return decideBySignWith(o, s, c, base, env, stdout, stderr, refuse)
+	}
 	const nothing = "nothing was decided; decide with an enrolled key instead (--sign-with <key> or --signature <file>)"
 	term, code := openApproverTerminal(env, stderr, "ailang-worldd approve: ", nothing)
 	if code != exitOK {
@@ -155,6 +167,39 @@ func runApproveDecide(o decideOptions, stdout, stderr io.Writer, env sessionEnv,
 		Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max, Typed: true}, proof, stdout, refuse)
 }
 
+// decideBySignWith is approve <request-ref> --sign-with (design §2.10, CD5): the 91-byte decision
+// message for the window's latest expiry, signed by env.signDecision (ssh-keygen -Y sign under
+// DecisionNamespace) under a root bounded by that expiry, and posted as sshsig. The daemon judges
+// it (design §2.7); no terminal opens and no proof is written.
+func decideBySignWith(o decideOptions, s decisionSubject, c *client, base string, env sessionEnv, stdout, stderr io.Writer, refuse func(string, ...any) int) int {
+	if env.signDecision == nil {
+		fmt.Fprintln(stderr, "ailang-worldd approve: internal error: no signer (sessionEnv.signDecision)")
+		return exitFatal
+	}
+	verb, _ := approval.DecisionVerb(o.verb)
+	body := surfacepage.DecideBody{V: 1, Verb: o.verb, Packet: strings.TrimPrefix(s.Packet, "sha256:"), Request: strings.TrimPrefix(o.ref, "sha256:"), Expiry: s.Window.Max}
+	msg, err := approval.DecisionMessage(verb, refDigest(s.Packet), refDigest(o.ref), body.Expiry)
+	if err != nil {
+		return refuse("%v", err)
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(max(body.Expiry-env.now(), 1))*time.Second)
+	defer cancel()
+	armored, err := env.signDecision(ctx, o.signWith, msg, stderr)
+	if err != nil {
+		return refuse("--sign-with %s: %v; nothing was decided", o.signWith, err)
+	}
+	sig := string(armored)
+	body.SSHSIG = &sig
+	return postDecision(c, base, o.ref, body, "", stdout, refuse)
+}
+
+// signDecisionWithSSHKeygen is the production decision signer: ssh-keygen from PATH through
+// host/sshsign's one driven, environment-scrubbed site, under DecisionNamespace. PATH is not a
+// trust boundary: the daemon verifies the output against the enrolled key over World's message.
+func signDecisionWithSSHKeygen(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error) {
+	return sshsign.SignIn(ctx, "ssh-keygen", approval.DecisionNamespace, key, msg, stderr)
+}
+
 // postDecision POSTs body to /surface/decide/{ref} with the daemon's own loopback Origin (in O by
 // construction; design §2.8) and, for the typed arm, the operator proof; it prints the decision.
 func postDecision(c *client, origin, ref string, body surfacepage.DecideBody, proof string, stdout io.Writer, refuse func(string, ...any) int) int {
diff --git a/cmd/ailang-worldd/approve_decide_key_test.go b/cmd/ailang-worldd/approve_decide_key_test.go
new file mode 100644
index 0000000..02de42f
--- /dev/null
+++ b/cmd/ailang-worldd/approve_decide_key_test.go
@@ -0,0 +1,111 @@
+package main
+
+// Row 201 A4 sprint 2 M5b (design §2.10, §6 CD5; CD6-CD9 are sprint 3's, plan S2-D1): approve
+// <request-ref> --sign-with against the fake daemon (approve_decide_rig_test.go): the 91-byte
+// decision message through the signer seam, and the real ssh-keygen on a key made in a temp dir
+// (D-WORLD-109: no real key, no agent).
+
+import (
+	"bytes"
+	"errors"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// keyRun is approve <ref> --db --daemon extra... with env.
+func (f *fakeDecider) keyRun(t *testing.T, env sessionEnv, extra ...string) (int, string, string) {
+	t.Helper()
+	var stdout, stderr bytes.Buffer
+	code := runApprove(append([]string{f.ref, "--db", f.db, "--daemon", f.addr()}, extra...), &stdout, &stderr, env)
+	return code, stdout.String(), stderr.String()
+}
+
+// decisionMessageOf is the message deciding verb on f's subject at expiry.
+func (f *fakeDecider) decisionMessageOf(t *testing.T, verb byte, expiry int64) []byte {
+	t.Helper()
+	m, err := approval.DecisionMessage(verb, refDigest(f.subject.Packet), refDigest(f.ref), expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return m
+}
+
+// CD5: --sign-with signs the 91-byte decision message for the window's max through the decision
+// signer, bounded by that expiry, and posts its output as sshsig (no terminal, no proof); a signer
+// error posts nothing; the production env's decision signer is signDecisionWithSSHKeygen.
+func TestApproveDecideSignWith(t *testing.T) {
+	f := newFakeDecider(t)
+	fs := &fakeSigner{ret: func([]byte) ([]byte, error) { return []byte("ARMORED"), nil }}
+	opened := false
+	env := approveEnv(nil, &opened, atTime(1710))
+	env.signDecision = fs.sign
+	code, stdout, stderr := f.keyRun(t, env, "--deny", "--sign-with", "/k/id")
+	if want := f.decisionMessageOf(t, approval.VerbDeny, 1760); code != exitOK || opened || fs.key != "/k/id" || !bytes.Equal(fs.msg, want) || len(fs.msg) != 91 ||
+		fs.left <= 0 || fs.left > 50*time.Second || len(f.posts) != 1 {
+		t.Fatalf("rc %d, opened %v, key %q, msg %x (left %v), %d POSTs\nstderr %s", code, opened, fs.key, fs.msg, fs.left, len(f.posts), stderr)
+	}
+	if b := f.postBodies[0]; b.Typed || b.SSHSIG == nil || *b.SSHSIG != "ARMORED" || b.Expiry != 1760 || b.Verb != "deny" || f.posts[0].Header.Get("X-World-Operator-Proof") != "" {
+		t.Fatalf("POST %+v", b)
+	}
+	if _, wrote := f.proofFile(t); wrote || !strings.Contains(stdout, "decided deny") {
+		t.Fatalf("proof written %v; stdout %q", wrote, stdout)
+	}
+	g := newFakeDecider(t)
+	env.signDecision = (&fakeSigner{ret: func([]byte) ([]byte, error) { return nil, errors.New("agent refused") }}).sign
+	if code, _, stderr := g.keyRun(t, env, "--sign-with", "/k/id"); code != exitUsage || len(g.posts) != 0 || !strings.Contains(stderr, "--sign-with /k/id: agent refused; nothing was decided") {
+		t.Fatalf("signer error: rc %d, %d POSTs\nstderr %s", code, len(g.posts), stderr)
+	}
+	if got, want := reflect.ValueOf(realSessionEnv().signDecision).Pointer(), reflect.ValueOf(signDecisionWithSSHKeygen).Pointer(); got != want {
+		t.Fatal("realSessionEnv().signDecision is not signDecisionWithSSHKeygen")
+	}
+	for _, flag := range []string{"--signature", "--write-message"} { // not a request ref's flags in this sprint
+		if code, _, stderr := g.keyRun(t, env, flag, "/x"); code != exitUsage || g.gets != 1 || !strings.Contains(stderr, "the other key flags do not decide it") {
+			t.Errorf("%s: rc %d, GETs %d\nstderr %s", flag, code, g.gets, stderr)
+		}
+	}
+}
+
+// CD5 with the real tool: the production decision signer runs the real ssh-keygen on a key made
+// in a temp dir (no agent, temp HOME); what it posts verifies with the daemon's own check
+// (VerifyDecisionSSHSIG) over the message for the window's max, and not over the mint namespace.
+// With WORLD_REQUIRE_SSH_KEYGEN=1 a missing ssh-keygen is a failure, not a skip (CI).
+func TestApproveDecideSignWithSSHKeygen(t *testing.T) {
+	if _, err := exec.LookPath("ssh-keygen"); err != nil {
+		if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
+			t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
+		}
+		t.Skip("SKIP COUNTED: ssh-keygen not found; the decision --sign-with end-to-end test did NOT run")
+	}
+	t.Setenv("HOME", t.TempDir())
+	t.Setenv("SSH_AUTH_SOCK", "")
+	key := filepath.Join(t.TempDir(), "k")
+	sshKeygen(t, "-q", "-t", "ed25519", "-N", "", "-C", "cd5", "-f", key)
+	pub, err := os.ReadFile(key + ".pub")
+	if err != nil {
+		t.Fatal(err)
+	}
+	k, err := approval.ParseAuthorizedKey(pub)
+	if err != nil {
+		t.Fatal(err)
+	}
+	f := newFakeDecider(t)
+	env := approveEnv(nil, nil, atTime(1710))
+	env.signDecision = signDecisionWithSSHKeygen
+	if code, _, stderr := f.keyRun(t, env, "--sign-with", key); code != exitOK || len(f.posts) != 1 {
+		t.Fatalf("rc %d, %d POSTs\nstderr %s", code, len(f.posts), stderr)
+	}
+	armored := []byte(*f.postBodies[0].SSHSIG)
+	if err := approval.VerifyDecisionSSHSIG(armored, approval.VerbApprove, refDigest(f.subject.Packet), refDigest(f.ref), 1760, k); err != nil {
+		t.Fatalf("the real signature does not decide: %v", err)
+	}
+	if err := approval.VerifyApproval(armored, refDigest(f.ref), refDigest(f.subject.Packet), 1760, k); !errors.Is(err, approval.ErrNamespace) {
+		t.Fatalf("checked as a mint: %v, want ErrNamespace", err)
+	}
+}
diff --git a/cmd/ailang-worldd/session.go b/cmd/ailang-worldd/session.go
index 90eb940..549ed3c 100644
--- a/cmd/ailang-worldd/session.go
+++ b/cmd/ailang-worldd/session.go
@@ -38,6 +38,9 @@ type sessionEnv struct {
 	// sign is approve --sign-with's signer, signWithSSHKeygen in production (row 201 A1b §2.9):
 	// the armored SSHSIG of msg by key. World never reads the private key; ssh-keygen does.
 	sign func(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error)
+	// signDecision is approve <request-ref> --sign-with's signer, signDecisionWithSSHKeygen in
+	// production (row 201 A4): the same, under the decision namespace world-decision@sunholo.
+	signDecision func(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error)
 }
 
 // sessionTTYFix ends every terminal refusal of the attended session verbs. It
@@ -55,6 +58,8 @@ func realSessionEnv() sessionEnv {
 		random: func(b []byte) error { _, err := rand.Read(b); return err },
 		after:  time.After,
 		sign:   signWithSSHKeygen,
+
+		signDecision: signDecisionWithSSHKeygen,
 	}
 }
 
diff --git a/host/sshsign/sshsign.go b/host/sshsign/sshsign.go
index 859f56f..69d1baf 100644
--- a/host/sshsign/sshsign.go
+++ b/host/sshsign/sshsign.go
@@ -39,12 +39,16 @@ var (
 	ErrEmptySignature    = errors.New("ssh-keygen -Y sign printed nothing")
 )
 
-// Sign runs `<bin> -Y sign -n world-approval@sunholo -f <key>` under ctx with msg on stdin and
-// returns its standard output, the armored signature. stderr (where ssh-keygen reports, and
-// from where a passphrase prompt goes to /dev/tty) is the operator's. A failed run (G1), over
-// MaxSignatureBytes of output (G2) and no output (G3) are errors.
+// Sign is SignIn under the mint namespace, world-approval@sunholo.
 func Sign(ctx context.Context, bin, key string, msg []byte, stderr io.Writer) ([]byte, error) {
-	cmd := exec.CommandContext(ctx, bin, "-Y", "sign", "-n", approval.Namespace, "-f", key)
+	return SignIn(ctx, bin, approval.Namespace, key, msg, stderr)
+}
+
+// SignIn runs `<bin> -Y sign -n <namespace> -f <key>` under ctx with msg on stdin and returns the
+// armored signature it prints (row 201 A4: decisions under world-decision@sunholo). stderr is the
+// operator's. A failed run (G1), over MaxSignatureBytes of output (G2) and none (G3) are errors.
+func SignIn(ctx context.Context, bin, namespace, key string, msg []byte, stderr io.Writer) ([]byte, error) {
+	cmd := exec.CommandContext(ctx, bin, "-Y", "sign", "-n", namespace, "-f", key)
 	cmd.Env = childenv.Scrubbed(os.Environ())
 	cmd.Stdin = bytes.NewReader(msg)
 	out := &cappedBuffer{max: MaxSignatureBytes}
diff --git a/host/store/context_roots_test.go b/host/store/context_roots_test.go
index 1ce516b..77fc94b 100644
--- a/host/store/context_roots_test.go
+++ b/host/store/context_roots_test.go
@@ -31,6 +31,7 @@ var contextRootPins = map[string]int{
 	// Row 201 A4 sprint 2: approve <request-ref>'s daemon reads before the y/N, and its POST after it (10 s each).
 	"cmd/ailang-worldd/approve_decide.go|runApproveDecide|Background": 1,
 	"cmd/ailang-worldd/approve_decide.go|postDecision|Background":     1,
+	"cmd/ailang-worldd/approve_decide.go|decideBySignWith|Background": 1, // --sign-with's signer, bounded by the signed expiry (≤ 600 s)
 
 	"cmd/ailang-worldd/cli.go|get|Background":                                1,
 	"cmd/ailang-worldd/cli.go|execute|Background":                            1,
````
