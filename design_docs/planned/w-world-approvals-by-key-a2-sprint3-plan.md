# Sprint plan — w-world-approvals-by-key-a2, sprint 3 (iteration 381, row 201)

> **THE CUT — measured, and the design's §9 sprint-3 fallback FIRES: M7 becomes its own fourth
> sprint.** Sprint 3 is **M6, the operator paths**: `approver enroll --format webauthn` (§2.12, WN1–WN11,
> with the §2.12a terminal split), `approver list`/`revoke` showing passkeys (§2.12b, L4/L5), and
> `approve <id> --webauthn-link` / `--webauthn-assertion <file>` with F1 widened to five modes (§2.13,
> WF1–WF6) — CW1–CW9 and the named edits.
> - **Measured: M6 is 1,686 changed lines under `cmd/ host/`** — `git diff --numstat <base> HEAD --
>   cmd host`, insertions + deletions, after a commit of each of the 8 scratch stages (467 production,
>   1,219 test). The design estimated M6 at ≈ 1,220 (400 + 820); it measured 1.38× that, in line with
>   sprint 2's M3b (890 vs 450).
> - **That leaves 314 lines under the 2,000 cap, and M7 cannot fit in them** (§1.5): its smallest safe
>   form is the page (index.html + passkey.js: create, get, the 240000 ms AbortController timer, the
>   enrolment code, the blob save), its embed and two routes, PG1, the `TestSurfaceAppRoutes` edit, the
>   Node/DevTools driver, XB1, XB2 and `generate.sh`. The three throwaway probes it would generalise
>   (iteration 379's recording and blob probes, and this iteration's virtual-time probe, §1.2 C18) are
>   already 202 lines of driver before any page, Go test or assertion. A partial M7 (routes without
>   XB1/XB2) is not a safe committed state: it would serve a page no browser test has run.
> - **So M7 — PG1, XB1, XB2, `generate.sh` and QUICKSTART §16 — moves whole to sprint 4 (Appendix C)**,
>   exactly as design §9 names. Every committed state of this sprint is safe: the CLI prints links to
>   `/surface/passkey`, which stays a 404 (`TestSurfaceAppRoutes` "every other name under /surface/ is a
>   404") until sprint 4 lands the page; nothing enrols without the typed code and `y`, and nothing
>   approves without a verified assertion (§8 R2).
> - **The browser was measured anyway, for sprint 4** (§1.2 C18): CDP `Emulation.setVirtualTimePolicy`
>   advances a page's 240000 ms `AbortController` timer to **exactly 240000 ms of page time
>   (`Date.now()` and `performance.now()` deltas) in 53–55 ms of real time, 3 of 3 runs**, under the
>   exact `surfaceCSP`, with an injected `navigator.credentials.create` that settles only on abort. XB2's
>   seam is therefore chosen and measured before sprint 4 starts.
>
> **Every line of code and test is a unified diff in Appendix P**, generated LAST from scratch stages
> that were compiled, vetted, gofmt-clean and green at every sub-milestone, after every mutant ran on
> the final tree. The patches were then extracted **from this file** with §2's `awk`, applied to a
> fresh export of `ab4b2402`, each sub-milestone's acceptance run, and the result compared
> byte-for-byte (`diff -r cmd host`) with the tree the mutants ran on (§1.2 C30–C31). The executor
> applies each patch and runs its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a2.md](w-world-approvals-by-key-a2.md), **r2** (quorum r2 +
  controller carve-out; **not reopened**). Where this plan and the design differ, this plan was
  measured against `ab4b2402` and wins; each difference is a decision (§1.3) or a premise found false
  (§1.4).
- **Starting ground**: sprint 1 (M0–M3a) landed as PR #546 (merge `d5d90e49`); sprint 2 (M3b, M4, M5)
  as PR #548 (merge `3e55451f`), recorded in `d9cb6cbb`. Sprint 2's Appendix C (what sprint 3 inherits)
  and the log-380/judge-380 carry-forwards are folded in (§1.3 S3-D7–S3-D9).
- **Planner base**: `ab4b2402116243c0e0fab9ed5cee8ec80a660a4e` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter381-row201-a2s3`, branch
  `sprint/world-iter381-row201-a2s3`. Planner scratch (outside the repo):
  `$HOME/.ailang/state/world-iter381/plan-scratch/` (`base/` the pristine export, `work/` the tree the
  mutants ran on, `stagebuild/` one commit per sub-milestone, `patches/`, `tools/`, `logs/`, `probe/`).
- **Executor**: as routed by the controller (a small model). **Four runs**: Run 1 (M6a, M6b), Run 2
  (M6c, M6d), Run 3 (M6e, M6f), Run 4 (M6g, M6h). Each sub-milestone is one `patch -p1` (§2) plus its
  acceptance (§3). **Apply, never retype.**
- **Touched**: `cmd/ailang-worldd/{approver.go, approve.go, main.go, help_test.go}`, new
  `cmd/ailang-worldd/{approver_passkey.go, approve_passkey.go, approver_label_test.go,
  approver_passkey_test.go, approver_passkey_ceremony_test.go, approver_passkey_refusals_test.go,
  approve_passkey_test.go, approve_passkey_refusals_test.go}`; `host/approval/cose.go` (one exported
  function, §1.4 P1), new `host/approval/cose_stored_test.go`; `host/store/context_roots_test.go` (the
  design-named +1 root).
- **Not touched**: every other `cmd/ailang-worldd` file (in particular `approve_sig.go`,
  `session*.go`, every existing `_test.go` but `help_test.go`), every other `host/approval` file,
  `host/approval/approvaltest/**`, `host/approval/testdata/**` (`generate.sh` stays the sprint-1 stub
  until sprint 4), every other `host/store` file, `host/authority/**`, `host/daemon/**`,
  `host/surfacepage/**`, `.gitattributes`, `docs/QUICKSTART.md` (§16 moves with M7), `world/*.ail`,
  `packages/**`, `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, any skill, the typed-`y`
  path's messages and exit codes, `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Fakes, temp stores and test keys only (D-WORLD-109 = A).** Every passkey in a committed test is
   `approvaltest.Passkey`'s fixed seed; every store is under `t.TempDir()`; every terminal is a test
   double (`ceremonyTerm`, `lockedTerm`, `noTerminalEnv`). No new test starts a subprocess, a browser,
   Chrome, Node or `ssh-keygen`. Nothing opens `~/.ailang/world`, `~/.ssh`, a real credentials store, a
   real daemon or `/dev/tty`; nothing mints into a real store or publishes. The planner's one browser
   probe (§1.2 C18) used a temp profile, `--headless=new`, `--remote-debugging-port=0` and no
   authenticator at all; it is not committed.
2. **Pure core, effects at the boundary (coding standards S1–S3).** The one `host/approval` change
   (`ParseStoredCOSEKey`) is pure and standard-library-only (`TestApprovalImportsOnlyTheStandardLibrary`
   green: it adds no import). Effects stay at the CLI boundary: `/dev/tty`, the clock, the timer and
   randomness through the existing `sessionEnv` seams; the store through `host/store`'s
   deadline-bound methods. No new Go module, package, exec site, route or `.ail`; **one** new context
   root entry, as the design names (§1.2 C9).
3. **One mutant per new production branch, anchored to the diff** (§4): every `if`, `case`, operand,
   bound and copied field of Appendix P's production hunks, each run against its single named test
   **and** its package suite; named equivalents with their reasons; base controls where a test pins
   pre-existing behaviour. Every run, re-runs included, appends a line to the ledger
   `plan-scratch/logs/mutant_ledger.tsv`.
4. **Exit codes by constant, typed refusals, and nothing written on a refusal.** Every enrol refusal
   arm counts both `approvers` and `webauthn_credentials`; every approve refusal arm checks the request
   pending, no session, no daemon call and the passkey counter unchanged.
5. **No `t.Parallel()`** in any new test (design §6; `TestNoParallelWriteForkPackagesOutsideVerifygate`).
6. **Bisectable.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet` green, the
   `_test.go` compile fence (`go test -run '^$' ./...`) green, with its own tests and every earlier
   sub-milestone's tests green (§1.2 C31).
7. **No git writes by the executor** (§7). The controller commits from snapshots.
8. **The typed-`y` path is unchanged.** §2.12a's split keeps `confirmApproverAct`'s messages and exit
   codes; CA3, CA4, CA9, CA10 and every A1b CLI test pass **unedited** (they are the refactor's oracle,
   M6a).

---

## 1. Pre-flight

### 1.1 Baseline at `ab4b2402` (measured by the planner, 2026-11-11)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`), `go1.26.6
darwin/arm64`, in a `git archive ab4b2402` export at `$HOME/.ailang/state/world-iter381/plan-scratch/base`
(a throwaway git repo; `diff -r cmd host` against the worktree: identical, `BASE_EQUALS_WORKTREE`; the
worktree itself was only read). Full logs: `plan-scratch/logs/baseline/` (`tools/gates.sh`).

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `ab4b2402`) |
| `gofmt -l cmd host` | **0** | 0 s | empty |
| `go build ./... && go vet ./...` | **0** | 4 s | clean |
| `go test -run '^$' ./...` (compile fence) | **0** | 21 s | no FAIL line |
| `go test -count=1 ./host/approval/... ./host/store/ ./host/authority/ ./host/surfacepage/... ./cmd/ailang-worldd/ ./host/daemon/` | **0** | 280 s | all `ok` |
| `go test -race -count=1 ./host/approval/... ./host/surfacepage/... ./cmd/ailang-worldd/ -run 'Passkey\|Approver\|Approve\|ConfirmApprover\|ReadTermLine\|EveryVerbHasHelp'` | **0** | 20 s | no race report |
| the §6 inventories by name with `-v` (design §7's 16 tests + the six reflective tests + `TestEveryVerbHasHelp` + every `Quickstart` test) | **0** | 23 s | every named test `--- PASS`; **one pre-existing `--- SKIP: TestSeToolsQuickstartPayloadsVerbatim`** in the `-run Quickstart` sweep of `./...` (not in a touched package; it skips at the base by its own gate; §6 excludes it by name) |
| `./scripts/verify_ail.sh` | **0** | 45 s | `✓ verify gate PASSED` |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` | **1** | 389 s | **red at the base, `host/verifygate` only**: `TestMissionConfigDefaultFlowIntegration/invalid_ledger` (`default-flow fixture timed out`, a 10 s fixture deadline), the known load-sensitive test the directive lists; every other leg and package green. Re-run alone while the planner's mutant workers loaded the machine it timed out again (`logs/baseline/verifygate_alone.log`, both subtests at 10.1 s); re-run alone after the mutants, at the shared machine's load from other sessions: it fails at the base too (§1.2 C27) |

**Rule (carried from sprint 2):** a red at the base measures the repo, not the sprint. `host/verifygate`
and `host/broker` are not touched by this sprint; the controller re-runs a red package alone once and
does not chase it. **Never start a gate under `nohup`** (sprint 2 §1.1: it reds
`scripts/marketingcheck`'s `shell-hup` test).

### 1.2 Codebase claims checked (at `ab4b2402` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD origin/dev` (worktree) → `ab4b2402116243c0e0fab9ed5cee8ec80a660a4e` twice; `git status --porcelain \| wc -l` → `0` |
| C2 | `confirmApproverAct` is one function holding N9a, N9, the bounded read and the y test (design §2.12a's premise) | `grep -n 'func confirmApproverAct' cmd/ailang-worldd/approver.go` → `33:`; its body (`:33-55`) holds `env.after == nil` (N9a), `env.openTerminal()` (N9), `readTermLineErrWithin(term, approverConfirmBudget, env.after)`, and `rerr != nil \|\| !strings.HasSuffix(line, "\n") \|\| !isYes(line)` |
| C3 | **No exported parser turns stored COSE bytes back into a `COSEKey`** (sprint 2 Appendix C says "parse the stored canonical COSE key") | `grep -n '^func [A-Z][A-Za-z]*(b \[\]byte) (COSEKey' host/approval/*.go` → **empty**; control in the same call: `grep -n '^func ParseAuthData' host/approval/authdata.go` → `36:func ParseAuthData(b []byte, registration bool)`; `ParseCOSEKey(v cborValue)` (`cose.go:36`) takes the unexported `cborValue`, and its only production caller is `authdata.go:67`. Routed by §1.4 P1 |
| C4 | `approverKeyLabel` assumes an SSH blob (design §1.9) | `approver.go:261 func approverKeyLabel(blob []byte) string`, `:262 keyType := "invalid-key"`; a webauthn row would print `invalid-key SHA256:…` |
| C5 | F1 counts three string modes | `approve.go:148 for _, v := range []string{*sigPath, *signWith, *msgPath} {`, `:153 if modes > 1 { // F1` |
| C6 | Nothing of sprint 3 exists | `grep -rln --include='*.go' 'surface/passkey' cmd host` → **empty**; control: `grep -rln --include='*.go' 'surface/app.js' host/daemon` → `daemon.go`, `surface_page.go`. `grep -rn --include='*.go' 'webauthn-link\|webauthn-assertion\|runApproverEnrollPasskey' cmd` → **empty**; control: `grep -rn '"write-message"' cmd` → `approve.go:132` |
| C7 | The CA5 `webauthn` arm and design §6's "edited" list | `approver_enroll_test.go:230 {name: "webauthn", … want: []string{"--registration"}` (sprint 1 already changed `row 201 A2` to `--registration`). On work it stays green **unedited**: WN1 prints its reason and the help, whose new usage line contains `--registration` (§1.4 P4) |
| C8 | `approver enroll`'s and `approve`'s read budget | `mint_request.go:24 const approveReadBudget = 3 * time.Second` |
| C9 | The context-root census | `context_roots_test.go:20-24`: `approve.go\|runApprove\|Background: 1`, `approver.go\|runApproverList: 1`, `\|runApproverEnroll: 2`, `\|runApproverRevoke: 3`. On work: **+1 entry** `approver_passkey.go\|runApproverEnrollPasskey\|Background: 3` (the pre-check, WN11's re-check, the write); the approve paths add none (they take `runApprove`'s root as a parameter) — `TestProductionContextRoots` `--- PASS` (work, `logs/final_gates/inv_store.log`) |
| C10 | `/surface/passkey` is a 404 until sprint 4 | `surface_app_test.go:18 // bytes, and every other name under /surface/ is a 404.`; the daemon registers no `passkey` route (C6) |
| C11 | The AC10 exec-site census skips `_test.go` (for sprint 4's driver) | `registry_publish_test.go:1094 if d.IsDir() \|\| !strings.HasSuffix(path, ".go") \|\| strings.HasSuffix(path, "_test.go") {` |
| C12 | **The registration law accepts an all-zero challenge** (carry-forward, log/judge 380) — a throwaway probe test on work, deleted | `NewPasskeyApprover` with `Challenge: [32]byte{}` and an `approvaltest.Passkey` registration made for it → `PROBE zero challenge: err=<nil> format="webauthn"` (`logs/probe_zero_challenge.log`). Routed by S3-D7 |
| C13 | `bufio.Reader.ReadString`'s contract (the two named equivalents of M6a) | `go doc bufio.Reader.ReadString` → "returns err != nil if and only if the returned data does not end in delim" |
| C14 | `ApproveByAssertion` with a zero `Passkey.Key` fails closed (judge-380 note) | `cose.go:112 func (k COSEKey) verify` has cases `AlgES256`, `AlgEdDSA` and a default returning `ErrCOSEAlg`; this sprint never builds a zero key: `approveAssertion` refuses a stored key `ParseStoredCOSEKey` rejects (WF5, MUT-WF5-KEY killed). No change |
| C15 | The A1b rigs this sprint reuses | `approve_test.go:68 type approveRig`, `:80 newApproveRig`, `:399 approveEnv`, `:1010 approveEffectProblems`; `approver_test.go:26 type lockedTerm`; `session_test.go:38 noTerminalEnv`; `approve_signwith_test.go:56 countDaemonCalls`; `approver_rig_test.go:43 approverStore`; `approvaltest/passkey.go:85 Register(origin, rpID, challenge, …Edit)`, `:96 Assert(origin, rpID, challenge, counter, …Edit)` |
| C16 | Instruments | `go version` → `go1.26.6 darwin/arm64`; `node --version` → `v26.13.0`; `Google Chrome --version` → `Google Chrome 157.0.8103.50`; `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0` |
| C17 | The iteration-379 probes still exist (sprint 4's starting point) | `shasum -a 256 /tmp/w379-designer/browser/probe.mjs` → `01a17c0d…37cf` (= SOURCE.txt's); `/tmp/w379-designer/blob/probe.mjs` (66 lines, the V31 blob-save probe) and `out.json` (`downloadProgress completed`, `cspViolations: []`). They are in `/tmp` and may vanish; Appendix C restates what they measured |
| C18 | **XB2's virtual-time seam works (measured for sprint 4)** | `plan-scratch/probe/vt_probe.mjs` (68 lines, not committed): Chrome 157 headless, temp profile, a page served with the exact `surfaceCSP`, a button whose click starts `setTimeout(abort, 240000)` on an `AbortController` and awaits an injected `navigator.credentials.create` that settles only on `abort`; then `Emulation.setVirtualTimePolicy({policy: "advance", budget: 250000})`. Three runs (`vt_run{1,2,3}.json`): `AbortError`, `expired: true`, `signal.aborted: true`, **`dateDelta: 240000`, `perfDelta: 240000`**, the button disabled while pending and re-enabled after, `virtualTimeBudgetExpired` fired, **53–55 ms real time for the 250 s budget**, 2.3–4.8 s per run in all |
| C19 | M6's sizes, per stage (work) | §1.5's table, from `tools/build_patches.sh` over `stagebuild` (one commit per stage) |
| C20 | The CW4 flake is gone (work, after §1.4 P6's fix) | `go test ./cmd/ailang-worldd -run '^TestPasskey' -count=40` → `ok` (10.5 s); then 15 fresh `-count=1` processes → all `ok`; `go test -race ./cmd/ailang-worldd -run '^(TestPasskey\|TestApprover\|TestConfirmApprover)' -count=3` → `ok`, no race report |
| C25 | **Every repo-wide inventory, on the final tree** (the committed round-trip tree, `-v`, tee'd) | the §6 script (`tools/sprint_gates.sh` = §6 verbatim) on `plan-scratch/roundtrip` at M6h → `PRE_COMMIT_GATES_GREEN` in 518 s: **30 named tests `--- PASS`, 0 SKIP/FAIL** (`logs/final_gates/inv.log`): design §7's 16, the reflective six, `TestEveryVerbHasHelp` and the seven quickstart tests. `TestProductionContextRoots` sees the +1 root and nothing else; `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestFrozenV1RouteTableMatchesMux` (13 routes) and `TestSurfaceAppRoutes` unchanged: no exec site, module, transport import, route or `t.Parallel` was added |
| C26 | Whole packages on the final tree | same run: `go test -count=1 ./host/approval/... ./host/store/ ./host/authority/ ./cmd/ailang-worldd/ ./host/daemon/` → all `ok` (approval 1.0 s, approvaltest 0.5 s, store 13.8 s, authority 0.8 s, `cmd/ailang-worldd` 219.1 s, `host/daemon` 339.9 s; `logs/final_gates/pkgs.log`); every a0a1/A1b/sprint-1/sprint-2 test green with only `help_test.go` and `context_roots_test.go` edited |
| C27 | **The whole repository on the final tree** (committed) — and the machine | `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` on the round-trip tree at M6h → `✓ 0 other binary blobs among 6705 tracked files`, toolchain floor, `✓ all 64 required top-level evidence tests passed exactly once`, every touched package `ok`; **rc 1 on three untouched packages' timing tests**: `host/broker` `TestExecProtectedCheckTimeout/reserve-after-before`, `host/daemon` `TestAGUISlowReaderCutIsResumable`, `host/verifygate` `TestMissionConfigDefaultFlowIntegration` (`logs/final_gates/verify_go.log`). The machine was shared: load average **38–60** from other sessions' processes (`ps`: other agents' `w3-*` scratch binaries, `ailang run`, `godot`, an Ollama runner). Each re-run alone, on both trees, at that load (`logs/rerun_alone.summary`): final tree — broker **PASS**, daemon **PASS**, verifygate FAIL; **base — broker FAIL**, daemon PASS, **verifygate FAIL**. They measure the machine, not the sprint (no file of those packages changes); §6's re-run rule applies. `./scripts/verify_ail.sh` → `✓ verify gate PASSED` |
| C28 | `go vet ./...`, gofmt and `-race` on the final tree | same run: `go vet ./...` empty; `gofmt -l cmd host` empty; `go test -race` approval/... 3.6 s, store 47.3 s, the `cmd/ailang-worldd` passkey/approver tests 13.2 s → `ok`, no race report (`logs/final_gates/race.log`) |
| C29 | Every mutant, on the final tree | 124 rows, the last ledger line of each: **121 KILLED (named and suite), 3 named equivalents SURVIVED (named and suite)**, 0 NOCOMPILE; base controls 3 KILLED; 146 ledger lines in all (§4) |
| C30 | **The patches reproduce the mutant tree** | a fresh `git archive ab4b2402` export + `patch -p1` of Appendix P's eight patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r cmd host` against the tree §4's mutants ran on (`plan-scratch/work`): **identical** (`ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`, `tools/roundtrip.sh`, `logs/roundtrip.log`, 673 s); nothing left uncommitted |
| C31 | **Every sub-milestone is bisectable** | the same fresh export, after **each** patch: `tools/accept.sh <M>` = §3's acceptance verbatim (gofmt empty, build + vet clean, the compile fence, the named `-v` tests with no `--- SKIP`/`--- FAIL` and at least one `--- PASS`, **[cli]** or the whole `cmd/ailang-worldd`, **[pkgs]**) → `ACCEPT <M> rc=0` for all eight; the cumulative numstat after a commit of each = §7's value: **38, 237, 701, 945, 1190, 1397, 1538, 1686** (`NUMSTAT` lines of `logs/roundtrip.log`) |

### 1.3 Planner decisions where the design leaves a choice

- **S3-D1 — The cut** (§1.5): the design's §9 sprint-3 fallback, by measurement. M7 and QUICKSTART §16
  move whole to sprint 4 (Appendix C); QUICKSTART §16 documents the page ceremony, so it moves with the
  page.
- **S3-D2 — Sub-milestones.** Eight, at ≤ ~250 lines where a split leaves each green on its own:
  **M6a** (38) the §2.12a split alone, so the unedited A1b tests are its oracle; **M6b** (199) labels
  and `ParseStoredCOSEKey`; **M6c** (464) the enrol path with its flag and pre-terminal refusals
  (CW1, CW3) — one production function cannot be split from its first tests without a stage carrying
  an untested path, the precedent of sprint 2's M3b1 (345); **M6d** (244) CW2, CW5; **M6e** (245) CW4,
  CW6; **M6f** (209) `--webauthn-link` alone (F1 at four modes); **M6g** (157) `--webauthn-assertion`
  (F1 at five); **M6h** (148) CW8 and the F1 pairs.
- **S3-D3 — The terminal is opened before the link is printed** (design §2.12 orders "6. the link, 7.
  the fence"). The enrol path runs N9a/N9 (`openApproverTerminal`) right after WN5 and only then
  prints the link, so a run without a terminal (an agent) prints no link it could not use; CW3 asserts
  an empty stdout (MUT-N9 killed). The messages and codes are N9a's and N9's.
- **S3-D4 — WN2 is one refusal over three operands** ("--name, --origin and --registration are
  required with --format webauthn"), so N2's ssh text is untouched; one CW1 arm per operand.
- **S3-D5 — WN10's "exactly 26 characters of A–Z2–7" is the comparison itself.** The typed line is
  upper-cased and stripped of `-`, spaces, tabs and the newline, then
  `subtle.ConstantTimeCompare`d with `EnrolmentCode(...)`. `ConstantTimeCompare` returns 0 for unequal
  lengths, and the expected value is always 26 characters of the base32 alphabet, so a separate length
  or alphabet branch could never change an outcome (its mutant would be an equivalent). CW5 still has
  the 25-, 27-character and `0` arms the design names, and the last-character arm kills a prefix
  compare (MUT-WN10-PREFIX).
- **S3-D6 — Two new files hold the CLI** (design §8 allows "or a new `approve_passkey.go`"):
  `approver_passkey.go` (enrol, `ceremonyLink`, the pre-check) and `approve_passkey.go` (links,
  `requestHash`, `approveAssertion`). `approver.go` gains the WN1 dispatch, the help and the label
  helpers; `approve.go` the two flags, F1's five modes and the two dispatches. Both approve paths use
  `runApprove`'s 3 s root (design: "No new context root").
- **S3-D7 — The all-zero challenge is refused (carry-forward, log/judge 380).** WN5 is `env.random ==
  nil || env.random(challenge[:]) != nil || challenge == [32]byte{}` → `exitFatal`, "internal error: no
  fresh challenge (sessionEnv.random is missing, failed, or gave 32 zero bytes)". Why: the law accepts
  a zero challenge (C12), and 32 zero bytes from a source that reported success is the signature of a
  source that never wrote the buffer (probability 2⁻²⁵⁶ otherwise); a predictable challenge would let a
  registration be made before the ceremony. Fail-closed costs one operand and one CW6 arm; the failed
  arm fills `0xff` before erroring, so the error operand is pinned on its own (MUT-WN5-ERR).
- **S3-D8 — `approveAssertion` rebuilds the enrolled passkey from the store row** (sprint 2 Appendix
  C): `Name`, `Format`, `EnrolledAt`, `LastUsedAt`, `RevokedAt` from `approvers`; `RPID`, `Origin`,
  `CredentialID`, `BackupEligible`, `SignCount` (the stored counter, which W15 compares) from
  `webauthn_credentials`; `Key` = `ParseStoredCOSEKey(approvers.public_key)` (§1.4 P1). `AAGUID` and
  `BackupState` are left zero: the assertion law reads neither. Every field the law reads has a
  copied-field mutant (MUT-AS-*), and `CredentialID`'s is a named equivalent (the row was found by
  that exact ID).
- **S3-D9 — `--webauthn-link`'s output** groups the active passkeys (`ActivePasskeys`, by name, at
  most 64) by origin in first-seen order; each link's fragment is `{"v":1,"op":"get","rp","origin",
  "request","expires","hash16","challenge","allow"}` (design §2.13) with `challenge` =
  base64url(`MessageHash`) and `hash16` its first 16 hex. It prints the request, its expiry and
  `hash16` first, so the operator can compare them with the page (design §2.14 step 2), and then the
  `--webauthn-assertion` step with the request's own `--db`/`--credentials-db`.
- **S3-D10 — `approver list` holds its read-only handle across the labels** (`defer ro.Close()`), so
  each webauthn row's `Passkey` read is under the same 3 s root (design §2.12b); revoke computes its
  label on its lookup handle before closing it.
- **S3-D11 — WN4 refuses anything but "does not exist"** (`!errors.Is(err, fs.ErrNotExist)` on
  `Lstat`): an existing file, a dangling symlink, or a path that cannot be checked; the text says
  "already exists or cannot be checked; pass a new path" (CW1's three arms; MUT-WN4-NIL killed by the
  unsearchable-directory arm).
- **S3-D12 — The test terminal reads one line per `Read`**, as a terminal in canonical mode does
  (`ceremonyTerm`), so the code and the `y` are two reads on one terminal; a timed-out prompt's read
  first closes the timer channel and then blocks until the test ends, so no read can run after its
  test completed (the first draft did, 1 run in ~6, §1.4 P6).
- **S3-D13 — The name grammar is judged once, by the registration law** (R1, reached at WN9), as design
  §2.12 lists no earlier name check for the webauthn form ("N2 and N3 as today"; the ssh form's grammar
  check is `NewApprover`'s N5, which the webauthn form does not call). A bad `--name` therefore refuses
  after the browser ceremony, writing nothing. The fragment carries it JSON-encoded and the page will
  show it with `textContent` only (design §2.14), so it is inert; §8 R6 and Appendix C carry the
  operator-time cost to sprint 4, where a one-line early check is cheap.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Sprint 2 Appendix C: "`Key` (parse the stored canonical COSE key)".** No exported function
  can (C3): `ParseCOSEKey` takes the unexported `cborValue`. Routed minimally in M6b: **one exported
  pure function in `host/approval/cose.go`**, `ParseStoredCOSEKey(b []byte) (COSEKey, error)` =
  `decodeCBORExact(b, attestationLimits)` (wrapping a decode error in `ErrCOSEKey`) then `ParseCOSEKey`
  (Q1–Q8) — the body of the existing test helper `parseCOSEBytes` (`cose_test.go:28`), exported. Pinned
  by the new `TestParseStoredCOSEKey` (`cose_stored_test.go`: both Chrome keys round-trip; trailing,
  truncated and empty bytes refused, each `errors.Is` its CBOR error; an RS256 key → `ErrCOSEAlg`) and by
  CW8's stored-key-invalid arm. **This is the one file outside design §8's sprint-3 conflict surface
  that a production hunk touches** (P3 lists the test files).
- **P2 — Design §9: M6 ≈ 1,220 lines (400 + 820), sprint 3 ≈ 1,800.** Measured M6 **1,686** (467 +
  1,219). With M7 the sprint cannot stay under 2,000 (§1.5), so design §9's named fallback fires.
- **P3 — "The design says this sprint touches no file in X" is a claim to measure (sprint 2's P2).**
  Design §8's sprint-3 rows name `cmd/ailang-worldd/{approver_passkey.go, approver.go, approve.go,
  approve_sig.go (or approve_passkey.go), main.go, help_test.go}` and the page files. Files outside
  them that this sprint's patches touch, and why: `host/approval/cose.go` + new `cose_stored_test.go`
  (P1); `host/store/context_roots_test.go` (the +1 root design §6/§7 name as an edited test, but §8
  does not list it); and the new `cmd/ailang-worldd/*_test.go` files (design §6 names
  `approver_passkey_test.go` and `approve_passkey_test.go`; the CW tests are split over six files to
  keep each sub-milestone ≤ ~250 lines). **No existing test file but `help_test.go` and
  `context_roots_test.go` changes**: the a0a1/A1b/sprint-1/sprint-2 tests are green unedited on the
  final tree (C26), including every arm design §6 expected to edit (P4).
- **P4 — Design §6 "edited: `approver_enroll_test.go:230-232`, whose `--format webauthn --pubkey`
  arm now expects WN1's usage refusal".** Not needed: sprint 1 already changed its substring to
  `--registration` (C7), and WN1 prints its reason followed by `approverHelp`, whose new usage line
  contains `--registration`. The arm stays green unedited; CW1's `WN1-webauthn-pubkey` arm asserts
  WN1's own text, which no other branch prints (MUT-WN1, MUT-WN1-PUBKEY killed).
- **P5 — The first mutant pass had 17 drafts that did not compile** (`if false` leaving `ok`, `complete`, `taken`,
  `i` or `err` unused — Go's "declared and not used"). Each was rewritten into a compiling form (`if
  !ok && false`, …) and re-run; the ledger keeps both the NOCOMPILE line and the re-run's (§4). Only
  the final run's verdicts count.
- **P6 — The first CW4 draft read after its test ended**, 1 run in ~6: a timed-out prompt's
  `ceremonyTerm.Read` (the abandoned read goroutine `readTermLineErrWithin` leaves, by design) ran its
  writer-lock probe after `t.TempDir` was removed and called `t.Errorf` → "Fail in goroutine after …
  has completed". Routed in the test (S3-D12): the unanswered prompt's read closes the timer channel
  itself and then blocks on `t.Context()`, so the timer can only fire once the read is parked; 40
  consecutive `-count` runs and 15 fresh processes, then `-race -count=3`, all green.
- **P7 — Design §2.13 lists WF5 as "no approver holds that credential ID".** A held credential whose
  stored key no longer parses (only a hand-edited row) needs a refusal too, or the zero key reaches the
  law: routed as a second WF5 text ("approver %q's stored passkey key is not valid"), pinned by CW8's
  `WF5-stored-key-invalid` arm (MUT-WF5-KEY).
- **P8 — QUICKSTART §16 belongs to M7, not M6.** It walks the page ceremony (`serve`, open the link,
  create, save), so with the page in sprint 4 the section would document a 404. It moves with M7
  (Appendix C); `docs/QUICKSTART.md` is untouched.

### 1.5 The size, and the cut

Lines are insertions + deletions under `cmd/ host/`: **`git diff --numstat <base> HEAD -- cmd host`**
(`.gitattributes` and docs are outside the filter and this sprint changes neither). **Every row is
measured** from the scratch stages, each committed in a throwaway repo (`plan-scratch/stagebuild`, tags
`M6a` … `M6h`; `tools/build_patches.sh`). The cumulative column is not the sum of the steps: a line
edited in two stages (F1's mode list, the help test's list, a comment) counts once in the cumulative
diff.

| Sub-milestone | Production | Test | Step | Cumulative vs `ab4b2402` (after its commit) | Run |
|---|---|---|---|---|---|
| M6a — §2.12a: `openApproverTerminal`, `askApproverLine`, `confirmOnTerminal` (behaviour-preserving) | 38 | 0 | **38** | 38 | 1 |
| M6b — L4/L5 labels in `list` and `revoke`; `ParseStoredCOSEKey`; CW9, `TestParseStoredCOSEKey` | 59 | 140 | **199** | 237 | 1 |
| M6c — `approver enroll --format webauthn` (WN1–WN11, N3, N6, N7, N9–N14), help, +1 root; CW1, CW3 | 219 | 245 | **464** | 701 | 2 |
| M6d — CW2 (both algorithms, BE), CW5 (read-then-substitute) | 0 | 244 | **244** | 945 | 2 |
| M6e — CW4 (unanswered, declined, timers), CW6 (WN5, WN8, WN9, WN11, N11–N13) | 0 | 245 | **245** | 1,190 | 3 |
| M6f — `approve --webauthn-link` (WF1), F1 at four modes, help; CW7 (links) | 86 | 123 | **209** | 1,397 | 3 |
| M6g — `approve --webauthn-assertion` (WF2–WF6), F1 at five modes, help; CW7 (the mint, replay, K8) | 73 | 84 | **157** | 1,538 | 4 |
| M6h — CW8 (WF2–WF6 before any effect), the F1/F2 pairs | 0 | 148 | **148** | **1,686** | 4 |
| **Sprint 3** | **467** | **1,219** | **1,686 ≤ 2,000** (margin 314) | | |
| *M7 (moved, Appendix C)* — page, embed, 2 routes, PG1, `TestSurfaceAppRoutes`, driver, XB1, XB2, `generate.sh` | — | — | *> 314: the probes it generalises are 202 lines alone (C17, C18)* | *> 2,000* | — |

**Why M7 does not fit, and why it moves whole.** Design §9: "Sprint 3 fallback cut: M7 becomes its own
fourth sprint." The margin left by the measured M6 is 314 lines. M7's parts, none optional: the page
(`index.html`, `passkey.js` with the fragment check, create, get, the AbortController/240000 ms timer,
the disabled button, the enrolment code over two SHA-256s, base64url, the blob save — design §2.14),
`host/surfacepage/passkey.go` (embed), two `GET` routes, PG1, the `TestSurfaceAppRoutes` edit, the
Node/DevTools driver `passkey_cdp.mjs` (create, get, download, virtual time — C17's two probes and
C18's are 66 + 68 + 68 lines before generalisation), XB1 and XB2 in Go, and `generate.sh`. The driver
alone, before any page or Go line, spends two thirds of the margin. **No subset of M7 is a safe
committed state**: routes without XB1/XB2 would serve a page no browser has run; XB1/XB2 without the
page test nothing. So M7 moves whole; M6 is complete and safe on its own (§8 R2).

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of
four backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`. Apply it **only** this
way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=M6a   # <- the sub-milestone, exactly as in the Appendix P heading: M6a M6b M6c M6d M6e M6f M6g M6h
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter381-row201-a2s3
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a2-sprint3-plan.md"
T="$HOME/.ailang/state/world-iter381/exec-patches"; mkdir -p "$T"
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
  `export AILANG_BIN=$HOME/.pinned-ailang/ailang`, each command's output `tee`'d to a log under
  `$HOME/.ailang/state/world-iter381/exec-logs/` (never only through `tail`).
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below), and
ends with the package lines (abbreviated **[cli]** and **[pkgs]**):

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/approval/... ./host/store/ ./cmd/ailang-worldd/
go test -run '^$' ./...
...
go test ./cmd/ailang-worldd/ -run '^(TestApprove|TestApprover|TestPasskey|TestConfirmApprover|TestReadTermLine|TestReadRegularFileCapped|TestEveryVerbHasHelp)' -count=1   # [cli]
go test ./host/approval/... ./host/store/ ./host/authority/ -count=1                                                                                                          # [pkgs]
```

M6a, M6c and M6h replace **[cli]** with the whole package, `go test ./cmd/ailang-worldd/ -count=1`
(~2 min): the refactor, the enrol dispatch and the final state reach every CLI test. Each `-v` run below
must show `--- PASS` for every test it names and **no `--- SKIP` or `--- FAIL`**.

### M6a — one terminal, two bounded reads (§2.12a) · 38 prod · Run 1

**Patch** `M6a`: `cmd/ailang-worldd/approver.go` — `confirmApproverAct` becomes
`openApproverTerminal` (N9a, N9) + `confirmOnTerminal` (`askApproverLine` with `approverConfirmBudget`,
then N10a, N10) with the **same messages and exit codes**; `askApproverLine(term, prompt, budget,
after) (line, ok, complete)` writes the prompt and reads one line within the budget (`complete` = no
read error and a trailing newline). No test changes: CA3, CA4, CA9, CA10 and `TestConfirmApproverAct`
are its unedited oracle.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestConfirmApproverAct|TestReadTermLineErrWithin|TestApproverEnrollDeclinedEnrolsNothing|TestApproverEnrollUnansweredPromptTimesOut|TestApproverEnrollWithoutATerminalEnrolsNothing|TestApproverRevokeUnansweredPromptTimesOut|TestApproverRevoke)$' -count=1 -v
go test ./cmd/ailang-worldd/ -count=1                                          # ~2 min
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M6a: 38**.

### M6b — `list` and `revoke` show passkeys (§2.12b); `ParseStoredCOSEKey` · 59 prod + 140 test · Run 1

**Patch** `M6b`:
- `cmd/ailang-worldd/approver.go`: `approverLabel(ctx, ro, row)` — an ssh row's `approverKeyLabel`
  (unchanged), a webauthn row's passkey read on the same handle: **L4** `passkey <ES256|EdDSA> <rpID>
  cred <16 hex>`, **L5** `invalid-passkey`; `passkeyLabel`, `passkeyCred` (first 16 hex of
  SHA-256(credential ID), as the enrolment y/N shows it). `approver list` keeps its read-only handle
  open across the labels (S3-D10); `revoke` labels its prompt on its lookup handle.
- `host/approval/cose.go`: `ParseStoredCOSEKey` (§1.4 P1).
- new `cmd/ailang-worldd/approver_label_test.go`: `pkOrigin`/`pkRPID` (`http://localhost:7644`,
  `localhost`), `seedPasskey` (straight through `EnrolPasskey`, as the CLI stores it); **CW9**
  `TestApproverListAndRevokeShowPasskeys` (an ES256 and a revoked EdDSA passkey and a hand-made
  webauthn row → L4, L4, L5 exactly; revoke's prompt shows the label; a schema-5 store lists its ssh
  approver with no passkey read, and a hand-made webauthn row there names schema 6).
- new `host/approval/cose_stored_test.go`: `TestParseStoredCOSEKey`.

**Acceptance**
```
[base]
go test ./host/approval/ ./cmd/ailang-worldd/ -run '^(TestParseStoredCOSEKey|TestApproverListAndRevokeShowPasskeys|TestApproverList|TestApproverRevoke)$' -count=1 -v
[cli]
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M6b: 237**.

### M6c — `approver enroll --format webauthn` (§2.12) · 219 prod + 245 test · Run 2

**Patch** `M6c`:
- `cmd/ailang-worldd/approver.go`: `--origin`, `--registration`; **WN1** (`webauthn && --pubkey`, or
  not webauthn with `--origin`/`--registration` → its reason and the help, `exitUsage`), then the
  dispatch to `runApproverEnrollPasskey` before N2/N4; the help (usage line, the ceremony paragraph,
  the flags).
- new `cmd/ailang-worldd/approver_passkey.go`: `passkeyCeremonyBudget = 300 s`; the
  `approverEnrolPasskey` seam; `ceremonyLink` and `link()`; `runApproverEnrollPasskey`: **WN2**, N3,
  **WN3**, **WN4**, the read-only pre-check (N6 below schema 6, N7) under root 1, **WN5** (S3-D7),
  N9a/N9 before the link (S3-D3), the link, the code prompt (**WN6** 300 s, **WN7**), **WN8**, **WN9**,
  **WN10** (S3-D5), **WN11** under root 2, the y/N on the same terminal (N10a, N10), the write under
  root 3 (N11, **N12** also on `ErrPasskeyCredentialEnrolled`, N13), N14; `passkeyEnrolPrecheck`.
- `cmd/ailang-worldd/main.go`: the usage line. `cmd/ailang-worldd/help_test.go`: that line in
  `TestEveryVerbHasHelp`'s list. `host/store/context_roots_test.go`: `runApproverEnrollPasskey: 3`.
- new `cmd/ailang-worldd/approver_passkey_test.go`: `enrolRig`, `passkeyCount`, `linkFragment`,
  `challengeOf`, `codeOf`, `enrolEnv`, `dropFlag`, `replaceFlag`; **CW1** `TestPasskeyEnrollFlagRules`
  (15 arms, each under a 5 s guard: WN1 ×3, WN2 ×3, N3, WN3 ×2, WN4 ×3, N6 v5, N7 ×2 — `exitUsage`, no
  terminal, no link, the seeds only) with **CW3** as its `CW3-no-terminal` subtest.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ ./host/store/ -run '^(TestPasskeyEnrollFlagRules|TestApproverEnrollRefusals|TestApproverDispatchAndHelp|TestEveryVerbHasHelp|TestProductionContextRoots)$' -count=1 -v
go test ./cmd/ailang-worldd/ -count=1                                          # ~2 min
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M6c: 701**.

### M6d — the ceremony and the substitute (CW2, CW5) · 244 test · Run 2

**Patch** `M6d`: new `cmd/ailang-worldd/approver_passkey_ceremony_test.go`: `ceremonyTerm` (one line per
read), `browserFor`, `enrolRun`, `enrolRig.enrol` (the browser saves at the code prompt; the writer
lock must be free at both prompts), `typeCode`, `enrolled`, `assertNothingEnrolled`; **CW2**
`TestPasskeyEnrollEnrolsWhatTheBrowserMade` (ES256 and EdDSA: the fragment is exactly `{v 1, create,
localhost, the origin, mark-mac, the challenge env.random drew}`; the terminal shows the code prompt and
`RenderPasskeyEnrolment`, ends with the y/N and is closed; N14's line; the approvers row and the passkey
row equal the registration's, field by field; a backup-eligible registration is stored as one);
**CW5** `TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode` (read-then-substitute → WN10, no
y/N, nothing enrolled; then last character, first group, 25 and 27 characters, a `0` → refused;
lowercase with dashes, spaces and a tab, and one string → enrolled).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestPasskeyEnrollEnrolsWhatTheBrowserMade|TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode)$' -count=1 -v
[cli]
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M6d: 945**.

### M6e — unanswered, declined and refused ceremonies (CW4, CW6) · 245 test · Run 3

**Patch** `M6e`: new `cmd/ailang-worldd/approver_passkey_refusals_test.go`: **CW4**
`TestPasskeyEnrollDeclinedOrUnanswered` (8 arms under a 5 s guard: WN6 code timeout — exactly
`no code within 5m0s`; WN7 EOF, partial line, read error; N10 `n`, EOF, `y`+EOF; N10a y timeout —
`no answer within 2m0s`; each: `exitUsage`, terminal closed, both tables empty, the writer lock free at
each prompt, and the timers asked for exactly `[5m0s]` or `[5m0s, 2m0s]`); **CW6**
`TestPasskeyEnrollRaceAndRefusals` (WN5 nil, failed and zero → `exitFatal` before any terminal; WN8
FIFO and 16 KiB+1; WN9 UV-clear and another challenge; WN11 key and credential stored by another enrol
while the code prompt waited — each refused before the y/N; N12 name, key, credential races →
`exitUsage`; N13 failure and N11 writer busy → `exitFatal`; nothing enrolled); `seedPasskeyRows`.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestPasskeyEnrollDeclinedOrUnanswered|TestPasskeyEnrollRaceAndRefusals)$' -count=1 -v
[cli]
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M6e: 1190**.

### M6f — `approve --webauthn-link` (§2.13, WF1) · 86 prod + 123 test · Run 3

**Patch** `M6f`:
- `cmd/ailang-worldd/approve.go`: `--webauthn-link`; F1 over four modes; after admission and F4, the
  link dispatch (nothing else runs); the help (usage line and paragraph).
- new `cmd/ailang-worldd/approve_passkey.go`: `writePasskeyLinks` (S3-D9; below schema 6 the restart
  text; **WF1**), `requestHash`.
- `main.go`, `help_test.go`: the usage line.
- new `cmd/ailang-worldd/approve_passkey_test.go`: `passkeySigner`, `requestHashOf` (computed from the
  row with `approval.MessageHash`, not the production helper), `passkeyRig`; **CW7** (links)
  `TestApproveByPasskey`: WF1 on an empty store, then on a schema-5 store; four passkeys (two at
  localhost, one at `https://x.ts.net`, one revoked) → exactly two links, each fragment equal to the
  expected one (the active credentials by name), the request/expiry/hash line and the next step; no
  terminal, no daemon call, the request still pending.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveByPasskey|TestEveryVerbHasHelp|TestApproveFlagRules)$' -count=1 -v
[cli]
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M6f: 1397**.

### M6g — `approve --webauthn-assertion` (§2.13, WF2–WF6) · 73 prod + 84 test · Run 4

**Patch** `M6g`:
- `cmd/ailang-worldd/approve.go`: `--webauthn-assertion`; F1 over five modes; the assertion judged
  after admission and F4 and before any network call, setting `signed`; the help.
- `cmd/ailang-worldd/approve_passkey.go`: `approveAssertion` (**WF2** 16 KiB descriptor-checked read,
  **WF3** envelope, **WF4** below schema 6, **WF5** unknown credential and an unparsable stored key,
  **WF6** `ApproveByAssertion` against the store row, S3-D8).
- `main.go`, `help_test.go`: the usage line.
- `approve_passkey_test.go`: `writeAssertion`, `passkeyOf`; **CW7** (the mint) continues
  `TestApproveByPasskey`: the assertion over the link's challenge mints with no terminal,
  `approved_by = webauthn:mark-mac`, the counter 1 → 2, `approved by webauthn:mark-mac` printed,
  `approveEffectProblems` clean but for `approved_by`; the same file again → "was already used"; revoke
  lists the request (K8).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveByPasskey|TestEveryVerbHasHelp|TestApproveFlagRules|TestApproveBySignatureMintsWithNoTerminal)$' -count=1 -v
[cli]
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M6g: 1538**.

### M6h — refusals before any effect (CW8) and the five-mode F1 · 148 test · Run 4

**Patch** `M6h`: new `cmd/ailang-worldd/approve_passkey_refusals_test.go`: **CW8**
`TestApprovePasskeyRefusalsBeforeAnyEffect` (14 arms: WF2 FIFO, 16 KiB+1, missing; WF3; WF4 v5; WF5
unknown credential, invalid stored key; WF6 counter regression, revoked, enrolled after now, another
request's hash, backup eligibility, another origin; expired before the file is read — each `exitUsage`,
no terminal, **0 daemon calls**, the request pending, no session, the counter still 1);
`TestApprovePasskeyFlagRules` (the seven new F1 pairs of the five modes and `--list` with either passkey
mode → exactly the help, no store touched).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApprovePasskeyRefusalsBeforeAnyEffect|TestApprovePasskeyFlagRules|TestApproveFlagRules)$' -count=1 -v
go test ./cmd/ailang-worldd/ -count=1                                          # ~2 min
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M6h: 1686**.

---

## 4. Mutation table (124 rows run on the planner's final tree: 121 KILLED, 3 named equivalents SURVIVED as predicted; plus 3 base controls)

**Derivation.** Anchored to Appendix P's production hunks (`approver.go`, `approver_passkey.go`, `approve.go`, `approve_passkey.go`, `host/approval/cose.go`): one row per new `if`, `case`, condition operand, bound, constant and copied field a test can tell apart — WN1–WN11 (every operand of WN1, WN2 and WN5), N3, N6, N7, N9, N11–N13 on the passkey path, L4/L5, F1's two new modes, WF1–WF6, the link's and the store rows' copied fields, the assertion input's copied fields, the code's normalisation (each replacement) and comparison — plus the M6a split's branches (pinned by the unedited A1b tests) and design §6's named mutants (a 120 s code budget, an unbounded read, a prefix compare, a case-sensitive compare, WN11 deleted, the lookup after the network call, the terminal opened on the passkey path; design's "length check dropped" has no branch to drop, S3-D5, and its 25/27-character arms stay). **Stage** is the sub-milestone whose patch adds the branch; its killing test may land later in the sprint (CW4–CW6 in M6d/M6e kill M6c rows; CW8 in M6h kills M6g rows). **Procedure**: `plan-scratch/tools/mutate.py <tree> tools/mutants.json logs/mutant_ledger.tsv <run> --workers 6` (outside the repo; `tools/mk_mutants.py` generates the table and asserts every anchor occurs exactly once in the final tree): for each row, in its own copy of the tree, apply **one** edit, assert the file changed, `go vet` the package (a mutant that does not compile is not a mutant), run **the single named test** (`go test <pkg> -run '^<test>$' -count=1`) **and then the whole package** (`go test <pkg> -count=1`; `cmd/ailang-worldd` takes ~2–6 min under the parallel load), restore, assert the restore is byte-identical, and **append one line to the ledger** `plan-scratch/logs/mutant_ledger.tsv` (UTC time, run label, id, stage, named, suite, expected, MATCH/MISMATCH). A row is KILLED only when both fail. Run with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` on the final tree (`plan-scratch/work`, byte-identical to the round-tripped tree, §1.2 C30). The ledger holds every run — the first pass `run1-final` (including its NOCOMPILE drafts, §1.4 P5), the re-run of the rewritten drafts `run2-rewritten`, and the base controls `base-controls` — and the table below is the last line of each id. **Never commit or snapshot a mutant.** The exact strings are `tools/mutants.json`; in the table tabs are dropped, a newline prints as ⏎ and a backquote as '.

| Id | Stage | Mutant (single edit; ⏎ is a newline) | Named test (arm) | Named | Suite |
|---|---|---|---|---|---|
| MUT-OPEN-AFTER | M6a | approver.go: if env.after == nil { // N9a / K4a → if false { // N9a / K4a | `TestApproverEnrollUnansweredPromptTimesOut` (no-timer arm (internal error, no terminal)) | KILLED | KILLED |
| MUT-OPEN-TTY | M6a | approver.go: if err != nil { // N9 / K4 → if false { // N9 / K4 | `TestApproverEnrollWithoutATerminalEnrolsNothing` (no terminal (a nil terminal is written)) | KILLED | KILLED |
| MUT-ASK-COMPLETE | M6a | approver.go: return line, ok, err == nil && strings.HasSuffix(line, "\n") → return line, ok, err != nil \|\| strings.HasSuffix(line, "\n") \|\| len(line) >= 0 | `TestApproverEnrollDeclinedEnrolsNothing` (y-then-eof) | KILLED | KILLED |
| MUT-ASK-ERR-EQUIV | M6a | approver.go: return line, ok, err == nil && strings.HasSuffix(line, "\n") → return line, ok, (err == nil \|\| len(line) >= 0) && strings.HasSuffix(line, "\n") | `TestApproverEnrollDeclinedEnrolsNothing` (equivalent: bufio.Reader.ReadString returns a non-nil error only for a line that does not end in the delimiter) | SURVIVED | SURVIVED |
| MUT-ASK-NL-EQUIV | M6a | approver.go: return line, ok, err == nil && strings.HasSuffix(line, "\n") → return line, ok, err == nil && (strings.HasSuffix(line, "\n") \|\| len(line) >= 0) | `TestApproverEnrollDeclinedEnrolsNothing` (equivalent: with a nil error ReadString's line ends in the delimiter, and a timed-out read (line "") is refused by !ok first) | SURVIVED | SURVIVED |
| MUT-CONFIRM-OK | M6a | approver.go: if !ok { // N10a / K5a → if !ok && false { // N10a / K5a | `TestApproverEnrollUnansweredPromptTimesOut` (the timeout message, exactly) | KILLED | KILLED |
| MUT-CONFIRM-COMPLETE | M6a | approver.go: if !complete \|\| !isYes(line) { // N10 / K5 → if !complete && false \|\| !isYes(line) { // N10 / K5 | `TestApproverEnrollDeclinedEnrolsNothing` (y-then-eof) | KILLED | KILLED |
| MUT-CONFIRM-YES | M6a | approver.go: if !complete \|\| !isYes(line) { // N10 / K5 → if !complete \|\| !isYes(line) && false { // N10 / K5 | `TestApproverEnrollDeclinedEnrolsNothing` (n) | KILLED | KILLED |
| MUT-CONFIRM-BUDGET | M6a | approver.go: line, ok, complete := askApproverLine(term, prompt, approverConfirmBudget, env.after) → line, ok, complete := askApproverLine(term, prompt, passkeyCeremonyBudget, env.after) | `TestApproverRevokeUnansweredPromptTimesOut` (the timer is 2m0s) | KILLED | KILLED |
| MUT-CONFIRM-CALL | M6a | approver.go: return confirmOnTerminal(term, env, stderr, prefix, prompt, nothing) → return exitOK | `TestApproverEnrollDeclinedEnrolsNothing` (n) | KILLED | KILLED |
| MUT-L-FORMAT | M6b | approver.go: if row.Format != approval.FormatWebAuthn { → if row.Format != approval.FormatWebAuthn \|\| true { | `TestApproverListAndRevokeShowPasskeys` (list (a passkey labelled as an ssh key)) | KILLED | KILLED |
| MUT-L-FORMAT-ALL | M6b | approver.go: if row.Format != approval.FormatWebAuthn { → if false { | `TestApproverListAndRevokeShowPasskeys` (v5 list of an ssh approver (a passkey asked of schema 5)) | KILLED | KILLED |
| MUT-L-ERR | M6b | approver.go: p, ok, err := ro.Passkey(ctx, row.Name)⏎if err != nil { → p, ok, err := ro.Passkey(ctx, row.Name)⏎if false && err != nil { | `TestApproverListAndRevokeShowPasskeys` (v5 list with a webauthn row) | KILLED | KILLED |
| MUT-L-LIST-ERR | M6b | approver.go: if labels[i], err = approverLabel(ctx, ro, r); err != nil { → if labels[i], err = approverLabel(ctx, ro, r); false && err != nil { | `TestApproverListAndRevokeShowPasskeys` (v5 list with a webauthn row) | KILLED | KILLED |
| MUT-L5 | M6b | approver.go: if !ok { // L5: only a hand-made row → if !ok && false { // L5: only a hand-made row | `TestApproverListAndRevokeShowPasskeys` (list (c: invalid-passkey)) | KILLED | KILLED |
| MUT-L4-ALG | M6b | approver.go: return passkeyLabel(p.Alg, p.RPID, p.CredentialID), nil // L4 → return passkeyLabel(approval.AlgEdDSA, p.RPID, p.CredentialID), nil // L4 | `TestApproverListAndRevokeShowPasskeys` (list (a: ES256)) | KILLED | KILLED |
| MUT-L4-RPID | M6b | approver.go: return passkeyLabel(p.Alg, p.RPID, p.CredentialID), nil // L4 → return passkeyLabel(p.Alg, "", p.CredentialID), nil // L4 | `TestApproverListAndRevokeShowPasskeys` (list) | KILLED | KILLED |
| MUT-L4-CRED | M6b | approver.go: return passkeyLabel(p.Alg, p.RPID, p.CredentialID), nil // L4 → return passkeyLabel(p.Alg, p.RPID, nil), nil // L4 | `TestApproverListAndRevokeShowPasskeys` (list) | KILLED | KILLED |
| MUT-CRED-LEN | M6b | approver.go: return hex.EncodeToString(h[:8]) → return hex.EncodeToString(h[:7]) | `TestApproverListAndRevokeShowPasskeys` (16 hex) | KILLED | KILLED |
| MUT-L-LIST-USE | M6b | approver.go: r.Name, r.Format, labels[i], approverTime(r.EnrolledAt), used, state) → r.Name, r.Format, approverKeyLabel(r.PublicKey)+labels[i][:0], approverTime(r.EnrolledAt), used, state) | `TestApproverListAndRevokeShowPasskeys` (list) | KILLED | KILLED |
| MUT-REVOKE-LABEL | M6b | approver.go: label, err = approverLabel(pre, ro, row) → label = approverKeyLabel(row.PublicKey) | `TestApproverListAndRevokeShowPasskeys` (revoke prompt) | KILLED | KILLED |
| MUT-REVOKE-LABEL-GUARD | M6b | approver.go: if err == nil && ok {⏎label, err → if false && err == nil && ok {⏎label, err | `TestApproverListAndRevokeShowPasskeys` (revoke prompt) | KILLED | KILLED |
| MUT-PSK-DECODE | M6b | approval/cose.go: v, err := decodeCBORExact(b, attestationLimits)⏎if err != nil {⏎return COSEKey{}, fmt.Errorf("%w: %w", ErrCOSEKey, err) → v, err := decodeCBORExact(b, attestationLimits)⏎if false && err != nil {⏎return COSEKey{}, fmt.Errorf("%w: %w", ErrCOSEKey, err) | `TestParseStoredCOSEKey` (trailing (the CBOR error is named)) | KILLED | KILLED |
| MUT-PSK-WRAP | M6b | approval/cose.go: return COSEKey{}, fmt.Errorf("%w: %w", ErrCOSEKey, err)⏎}⏎return ParseCOSEKey(v) → return COSEKey{}, fmt.Errorf("%w: %v", ErrCOSEKey, err)⏎}⏎return ParseCOSEKey(v) | `TestParseStoredCOSEKey` (truncated (errors.Is the CBOR error)) | KILLED | KILLED |
| MUT-WN1 | M6c | approver.go: if webauthn && *pubkey != "" \|\| !webauthn && (*origin != "" \|\| *registration != "") { // WN1 → if false { // WN1 | `TestPasskeyEnrollFlagRules` (WN1-webauthn-pubkey) | KILLED | KILLED |
| MUT-WN1-PUBKEY | M6c | approver.go: if webauthn && *pubkey != "" \|\| !webauthn && (*origin != "" \|\| *registration != "") { // WN1 → if !webauthn && (*origin != "" \|\| *registration != "") { // WN1 | `TestPasskeyEnrollFlagRules` (WN1-webauthn-pubkey) | KILLED | KILLED |
| MUT-WN1-ORIGIN | M6c | approver.go: if webauthn && *pubkey != "" \|\| !webauthn && (*origin != "" \|\| *registration != "") { // WN1 → if webauthn && *pubkey != "" \|\| !webauthn && (*registration != "") { // WN1 | `TestPasskeyEnrollFlagRules` (WN1-ssh-origin) | KILLED | KILLED |
| MUT-WN1-REG | M6c | approver.go: if webauthn && *pubkey != "" \|\| !webauthn && (*origin != "" \|\| *registration != "") { // WN1 → if webauthn && *pubkey != "" \|\| !webauthn && (*origin != "") { // WN1 | `TestPasskeyEnrollFlagRules` (WN1-ssh-registration) | KILLED | KILLED |
| MUT-WN1-AND | M6c | approver.go: if webauthn && *pubkey != "" \|\| !webauthn && (*origin != "" \|\| *registration != "") { // WN1 → if webauthn && *pubkey != "" \|\| !webauthn && (*origin != "" && *registration != "") { // WN1 | `TestPasskeyEnrollFlagRules` (WN1-ssh-origin) | KILLED | KILLED |
| MUT-WN1-SCOPE | M6c | approver.go: if webauthn && *pubkey != "" \|\| !webauthn && (*origin != "" \|\| *registration != "") { // WN1 → if webauthn && *pubkey != "" \|\| (*origin != "" \|\| *registration != "") { // WN1 | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the webauthn enrolment itself) | KILLED | KILLED |
| MUT-WN1-PUBKEY-SCOPE | M6c | approver.go: if webauthn && *pubkey != "" \|\| !webauthn && (*origin != "" \|\| *registration != "") { // WN1 → if *pubkey != "" \|\| !webauthn && (*origin != "" \|\| *registration != "") { // WN1 | `TestApproverEnrollDeclinedEnrolsNothing` (an ssh enrolment (--pubkey)) | KILLED | KILLED |
| MUT-DISPATCH | M6c | approver.go: if webauthn {⏎return runApproverEnrollPasskey( → if false {⏎return runApproverEnrollPasskey( | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the webauthn enrolment itself) | KILLED | KILLED |
| MUT-WN2 | M6c | approver_passkey.go: if name == "" \|\| origin == "" \|\| regPath == "" { // WN2 → if false { // WN2 | `TestPasskeyEnrollFlagRules` (WN2-name) | KILLED | KILLED |
| MUT-WN2-NAME | M6c | approver_passkey.go: if name == "" \|\| origin == "" \|\| regPath == "" { // WN2 → if origin == "" \|\| regPath == "" { // WN2 | `TestPasskeyEnrollFlagRules` (WN2-name) | KILLED | KILLED |
| MUT-WN2-ORIGIN | M6c | approver_passkey.go: if name == "" \|\| origin == "" \|\| regPath == "" { // WN2 → if name == "" \|\| regPath == "" { // WN2 | `TestPasskeyEnrollFlagRules` (WN2-origin) | KILLED | KILLED |
| MUT-WN2-REG | M6c | approver_passkey.go: if name == "" \|\| origin == "" \|\| regPath == "" { // WN2 → if name == "" \|\| origin == "" { // WN2 | `TestPasskeyEnrollFlagRules` (WN2-registration) | KILLED | KILLED |
| MUT-N3 | M6c | approver_passkey.go: credPath, refusal := approverStores(db, creds) // N3⏎if refusal != "" { → credPath, refusal := approverStores(db, creds) // N3⏎if false { | `TestPasskeyEnrollFlagRules` (N3-no-world-store) | KILLED | KILLED |
| MUT-WN3 | M6c | approver_passkey.go: if err != nil { // WN3 → if false { // WN3 | `TestPasskeyEnrollFlagRules` (WN3-ip) | KILLED | KILLED |
| MUT-WN4 | M6c | approver_passkey.go: if _, err := os.Lstat(regPath); !errors.Is(err, fs.ErrNotExist) { // WN4 → if _, err := os.Lstat(regPath); false && !errors.Is(err, fs.ErrNotExist) { // WN4 | `TestPasskeyEnrollFlagRules` (WN4-exists) | KILLED | KILLED |
| MUT-WN4-NIL | M6c | approver_passkey.go: if _, err := os.Lstat(regPath); !errors.Is(err, fs.ErrNotExist) { // WN4 → if _, err := os.Lstat(regPath); err == nil && fs.ErrNotExist != nil { // WN4 | `TestPasskeyEnrollFlagRules` (WN4-unsearchable-dir) | KILLED | KILLED |
| MUT-PRE1 | M6c | approver_passkey.go: refusal = passkeyEnrolPrecheck(pre, credPath, name, nil) // N6, N7, closed before any terminal⏎cancelPre()⏎if refusal != "" { → refusal = passkeyEnrolPrecheck(pre, credPath, name, nil) // N6, N7, closed before any terminal⏎cancelPre()⏎if false { | `TestPasskeyEnrollFlagRules` (N7-name-taken) | KILLED | KILLED |
| MUT-N6 | M6c | approver_passkey.go: if _, _, err := ro.Passkey(ctx, name); err != nil { // N6 → if _, _, err := ro.Passkey(ctx, name); false && err != nil { // N6 | `TestPasskeyEnrollFlagRules` (N6-v5) | KILLED | KILLED |
| MUT-N7 | M6c | approver_passkey.go: if taken { // N7 → if taken && false { // N7 | `TestPasskeyEnrollFlagRules` (N7-name-taken) | KILLED | KILLED |
| MUT-PRE-NIL | M6c | approver_passkey.go: if a == nil {⏎return "" → if false {⏎return "" | `TestPasskeyEnrollFlagRules` (N6-v5 passes; the first check then dereferences nil (any arm past it)) | KILLED | KILLED |
| MUT-WN5 | M6c | approver_passkey.go: if env.random == nil \|\| env.random(challenge[:]) != nil \|\| challenge == [32]byte{} { // WN5 → if false { // WN5 | `TestPasskeyEnrollRaceAndRefusals` (WN5 nil) | KILLED | KILLED |
| MUT-WN5-NIL | M6c | approver_passkey.go: if env.random == nil \|\| env.random(challenge[:]) != nil \|\| challenge == [32]byte{} { // WN5 → if env.random(challenge[:]) != nil \|\| challenge == [32]byte{} { // WN5 | `TestPasskeyEnrollRaceAndRefusals` (WN5 nil (panics)) | KILLED | KILLED |
| MUT-WN5-ERR | M6c | approver_passkey.go: env.random(challenge[:]) != nil \|\| challenge → (env.random(challenge[:]) != nil && false) \|\| challenge | `TestPasskeyEnrollRaceAndRefusals` (WN5 failed (it filled 0xff)) | KILLED | KILLED |
| MUT-WN5-ZERO | M6c | approver_passkey.go: if env.random == nil \|\| env.random(challenge[:]) != nil \|\| challenge == [32]byte{} { // WN5 → if env.random == nil \|\| env.random(challenge[:]) != nil { // WN5 | `TestPasskeyEnrollRaceAndRefusals` (WN5 zero) | KILLED | KILLED |
| MUT-WN5-CODE | M6c | approver_passkey.go: or gave 32 zero bytes); %s\n", prefix, nothing)⏎return exitFatal → or gave 32 zero bytes); %s\n", prefix, nothing)⏎return exitUsage | `TestPasskeyEnrollRaceAndRefusals` (WN5 (exitFatal)) | KILLED | KILLED |
| MUT-N9 | M6c | approver_passkey.go: term, code := openApproverTerminal(env, stderr, prefix, nothing) // N9a, N9: before the link, which only a person can use⏎if code != exitOK { → term, code := openApproverTerminal(env, stderr, prefix, nothing) // N9a, N9: before the link, which only a person can use⏎if false { | `TestPasskeyEnrollFlagRules` (CW3-no-terminal) | KILLED | KILLED |
| MUT-WN6 | M6c | approver_passkey.go: if !ok { // WN6 → if !ok && false { // WN6 | `TestPasskeyEnrollDeclinedOrUnanswered` (WN6-code-timeout (exact message)) | KILLED | KILLED |
| MUT-WN6-BUDGET | M6c | approver_passkey.go: passkeyCeremonyBudget, env.after) → approverConfirmBudget, env.after) | `TestPasskeyEnrollDeclinedOrUnanswered` (WN6-code-timeout (the 300 s timer)) | KILLED | KILLED |
| MUT-WN6-CONST | M6c | approver_passkey.go: const passkeyCeremonyBudget = 300 * time.Second → const passkeyCeremonyBudget = 120 * time.Second | `TestPasskeyEnrollDeclinedOrUnanswered` (WN6-code-timeout (5m0s)) | KILLED | KILLED |
| MUT-WN7 | M6c | approver_passkey.go: if !complete { // WN7 → if !complete && false { // WN7 | `TestPasskeyEnrollDeclinedOrUnanswered` (WN7-code-partial) | KILLED | KILLED |
| MUT-LINK-OP | M6c | approver_passkey.go: ceremonyLink{V: 1, Op: "create", RP: rpID, → ceremonyLink{V: 1, Op: "get", RP: rpID, | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the fragment) | KILLED | KILLED |
| MUT-LINK-RP | M6c | approver_passkey.go: ceremonyLink{V: 1, Op: "create", RP: rpID, → ceremonyLink{V: 1, Op: "create", RP: rpID + ":7644", | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the fragment) | KILLED | KILLED |
| MUT-LINK-NAME | M6c | approver_passkey.go: Origin: origin, Name: name, Challenge: → Origin: origin, Name: "", Challenge: | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the fragment) | KILLED | KILLED |
| MUT-LINK-CHALLENGE | M6c | approver_passkey.go: Challenge: base64.RawURLEncoding.EncodeToString(challenge[:])} → Challenge: base64.RawURLEncoding.EncodeToString(challenge[:16])} | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the fragment (the challenge from env.random)) | KILLED | KILLED |
| MUT-WN8 | M6c | approver_passkey.go: if err != nil { // WN8 → if false { // WN8 | `TestPasskeyEnrollRaceAndRefusals` (WN8-fifo) | KILLED | KILLED |
| MUT-WN8-CAP | M6c | approver_passkey.go: reg, err := readRegularFileCapped(regPath, approval.MaxPasskeyResponseBytes) → reg, err := readRegularFileCapped(regPath, 1<<20) | `TestPasskeyEnrollRaceAndRefusals` (WN8-over-16k) | KILLED | KILLED |
| MUT-WN9 | M6c | approver_passkey.go: if err != nil { // WN9 → if false { // WN9 | `TestPasskeyEnrollRaceAndRefusals` (WN9-uv-clear (panics on the nil passkey)) | KILLED | KILLED |
| MUT-WN9-CHALLENGE | M6c | approver_passkey.go: Origin: origin, Challenge: challenge, Response: reg, → Origin: origin, Challenge: [32]byte{}, Response: reg, | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (enrolment) | KILLED | KILLED |
| MUT-WN9-NOW | M6c | approver_passkey.go: Response: reg, Now: env.now()}) → Response: reg, Now: 0}) | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (enrolment (EnrolledAt)) | KILLED | KILLED |
| MUT-WN10 | M6c | approver_passkey.go: resp.AttestationObject))) != 1 { // WN10 → resp.AttestationObject))) != 1 && false { // WN10 | `TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode` (the substitute) | KILLED | KILLED |
| MUT-WN10-PREFIX | M6c | approver_passkey.go: subtle.ConstantTimeCompare([]byte(typed), []byte(approval.EnrolmentCode(resp.ClientDataJSON, resp.AttestationObject))) != 1 → subtle.ConstantTimeCompare([]byte(typed)[:min(8, len(typed))], []byte(approval.EnrolmentCode(resp.ClientDataJSON, resp.AttestationObject))[:8]) != 1 | `TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode` (last-character) | KILLED | KILLED |
| MUT-WN10-CDJ | M6c | approver_passkey.go: approval.EnrolmentCode(resp.ClientDataJSON, resp.AttestationObject) → approval.EnrolmentCode(resp.AttestationObject, resp.AttestationObject) | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (enrolment) | KILLED | KILLED |
| MUT-WN10-UPPER | M6c | approver_passkey.go: .Replace(strings.ToUpper(line)) → .Replace(line) | `TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode` (lowercase-dashes-spaces-tab) | KILLED | KILLED |
| MUT-WN10-DASH | M6c | approver_passkey.go: strings.NewReplacer("-", "", " ", "", "\t", "", "\n", "") → strings.NewReplacer(" ", "", "\t", "", "\n", "") | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the grouped code) | KILLED | KILLED |
| MUT-WN10-SPACE | M6c | approver_passkey.go: strings.NewReplacer("-", "", " ", "", "\t", "", "\n", "") → strings.NewReplacer("-", "", "\t", "", "\n", "") | `TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode` (lowercase-dashes-spaces-tab) | KILLED | KILLED |
| MUT-WN10-TAB | M6c | approver_passkey.go: strings.NewReplacer("-", "", " ", "", "\t", "", "\n", "") → strings.NewReplacer("-", "", " ", "", "\n", "") | `TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode` (lowercase-dashes-spaces-tab) | KILLED | KILLED |
| MUT-WN10-NL | M6c | approver_passkey.go: strings.NewReplacer("-", "", " ", "", "\t", "", "\n", "") → strings.NewReplacer("-", "", " ", "", "\t", "") | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (enrolment) | KILLED | KILLED |
| MUT-WN11-CALL | M6c | approver_passkey.go: refusal = passkeyEnrolPrecheck(pre, credPath, name, &a) // WN11 → refusal = passkeyEnrolPrecheck(pre, credPath, name, nil) // WN11 | `TestPasskeyEnrollRaceAndRefusals` (WN11-key) | KILLED | KILLED |
| MUT-WN11-REFUSAL | M6c | approver_passkey.go: // WN11⏎cancelPre()⏎if refusal != "" { → // WN11⏎cancelPre()⏎if false { | `TestPasskeyEnrollRaceAndRefusals` (WN11-key) | KILLED | KILLED |
| MUT-WN11-KEY | M6c | approver_passkey.go: if enrolled { // WN11: the key → if false { // WN11: the key | `TestPasskeyEnrollRaceAndRefusals` (WN11-key (the y/N is not shown)) | KILLED | KILLED |
| MUT-WN11-CRED | M6c | approver_passkey.go: if enrolled { // WN11: the credential → if false { // WN11: the credential | `TestPasskeyEnrollRaceAndRefusals` (WN11-credential) | KILLED | KILLED |
| MUT-YN | M6c | approver_passkey.go: Enrol approver %s? [y/N] ", name), nothing); code != exitOK { → Enrol approver %s? [y/N] ", name), nothing); false { | `TestPasskeyEnrollDeclinedOrUnanswered` (N10-n) | KILLED | KILLED |
| MUT-YN-RENDER | M6c | approver_passkey.go: approval.RenderPasskeyEnrolment(a)+fmt.Sprintf( → fmt.Sprintf( | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the y/N shows the passkey) | KILLED | KILLED |
| MUT-N11 | M6c | approver_passkey.go: if err != nil { // N11 → if false { // N11 | `TestPasskeyEnrollRaceAndRefusals` (N11-writer-busy) | KILLED | KILLED |
| MUT-N12-NAME | M6c | approver_passkey.go: case errors.Is(err, store.ErrApproverNameTaken), errors.Is(err, store.ErrApproverKeyEnrolled), errors.Is(err, store.ErrPasskeyCredentialEnrolled): // N12 → case errors.Is(err, store.ErrApproverKeyEnrolled), errors.Is(err, store.ErrPasskeyCredentialEnrolled): // N12 | `TestPasskeyEnrollRaceAndRefusals` (N12-name) | KILLED | KILLED |
| MUT-N12-KEY | M6c | approver_passkey.go: case errors.Is(err, store.ErrApproverNameTaken), errors.Is(err, store.ErrApproverKeyEnrolled), errors.Is(err, store.ErrPasskeyCredentialEnrolled): // N12 → case errors.Is(err, store.ErrApproverNameTaken), errors.Is(err, store.ErrPasskeyCredentialEnrolled): // N12 | `TestPasskeyEnrollRaceAndRefusals` (N12-key) | KILLED | KILLED |
| MUT-N12-CRED | M6c | approver_passkey.go: case errors.Is(err, store.ErrApproverNameTaken), errors.Is(err, store.ErrApproverKeyEnrolled), errors.Is(err, store.ErrPasskeyCredentialEnrolled): // N12 → case errors.Is(err, store.ErrApproverNameTaken), errors.Is(err, store.ErrApproverKeyEnrolled): // N12 | `TestPasskeyEnrollRaceAndRefusals` (N12-credential) | KILLED | KILLED |
| MUT-N13 | M6c | approver_passkey.go: case err != nil: // N13 → case false: // N13 | `TestPasskeyEnrollRaceAndRefusals` (N13-failure) | KILLED | KILLED |
| MUT-N13-CODE | M6c | approver_passkey.go: fmt.Fprintf(stderr, "%scannot enrol approver %q: %v\n", prefix, a.Name, err)⏎return exitFatal → fmt.Fprintf(stderr, "%scannot enrol approver %q: %v\n", prefix, a.Name, err)⏎return exitUsage | `TestPasskeyEnrollRaceAndRefusals` (N13-failure (exitFatal)) | KILLED | KILLED |
| MUT-ROW-KEY | M6c | approver_passkey.go: PublicKey: p.Key.Canonical(), EnrolledAt: a.EnrolledAt} → PublicKey: p.CredentialID, EnrolledAt: a.EnrolledAt} | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the stored row) | KILLED | KILLED |
| MUT-ROW-AAGUID | M6c | approver_passkey.go: AAGUID: p.AAGUID[:], → AAGUID: make([]byte, 16), | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the stored passkey) | KILLED | KILLED |
| MUT-ROW-ALG | M6c | approver_passkey.go: Alg: p.Key.Alg, → Alg: approval.AlgES256, | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the stored passkey (EdDSA)) | KILLED | KILLED |
| MUT-ROW-BE | M6c | approver_passkey.go: BackupEligible: p.BackupEligible, SignCount: int64(p.SignCount)}) → BackupEligible: false, SignCount: int64(p.SignCount)}) | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the backup-eligible passkey) | KILLED | KILLED |
| MUT-ROW-COUNT | M6c | approver_passkey.go: BackupEligible: p.BackupEligible, SignCount: int64(p.SignCount)}) → BackupEligible: p.BackupEligible, SignCount: 0}) | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` (the stored passkey (counter 1)) | KILLED | KILLED |
| MUT-F1-LINK | M6f | approve.go: for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *passkeyLink, *assertion != ""} { → for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *assertion != ""} { | `TestApprovePasskeyFlagRules` (the --webauthn-link pairs) | KILLED | KILLED |
| MUT-F1-ASSERT | M6g | approve.go: for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *passkeyLink, *assertion != ""} { → for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *passkeyLink} { | `TestApprovePasskeyFlagRules` (the --webauthn-assertion pairs) | KILLED | KILLED |
| MUT-WF-LINK-DISPATCH | M6f | approve.go: if *passkeyLink { // WF1 → if false { // WF1 | `TestApproveByPasskey` (link (no terminal)) | KILLED | KILLED |
| MUT-WF-LINK-REFUSAL | M6f | approve.go: if refusal := writePasskeyLinks(ctx, credPath, stdout, row, mr); refusal != "" { → if refusal := writePasskeyLinks(ctx, credPath, stdout, row, mr); false && refusal != "" { | `TestApproveByPasskey` (WF1-none-enrolled) | KILLED | KILLED |
| MUT-WF-ASSERT-DISPATCH | M6g | approve.go: if *assertion != "" { // the passkey → if false { // the passkey | `TestApproveByPasskey` (assertion (no terminal)) | KILLED | KILLED |
| MUT-WF-ASSERT-REFUSAL | M6g | approve.go: if signed, refusal = approveAssertion(ctx, credPath, *assertion, row, now); refusal != "" { → if signed, refusal = approveAssertion(ctx, credPath, *assertion, row, now); false { | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF6-revoked (the terminal opens)) | KILLED | KILLED |
| MUT-LINK-V6 | M6f | approve_passkey.go: if err != nil { // below schema 6 names the serve restart → if false && err != nil { // below schema 6 names the serve restart | `TestApproveByPasskey` (WF1-none-enrolled (v5)) | KILLED | KILLED |
| MUT-WF1 | M6f | approve_passkey.go: if len(active) == 0 { // WF1 → if false { // WF1 | `TestApproveByPasskey` (WF1-none-enrolled) | KILLED | KILLED |
| MUT-HASH16 | M6f | approve_passkey.go: hash16 := hex.EncodeToString(h[:8]) → hash16 := hex.EncodeToString(h[:7]) | `TestApproveByPasskey` (the links) | KILLED | KILLED |
| MUT-LINK-GROUP | M6f | approve_passkey.go: if l == nil { → if true { | `TestApproveByPasskey` (two links, not three) | KILLED | KILLED |
| MUT-LINK-ALLOW | M6f | approve_passkey.go: l.Allow = append(l.Allow, base64.RawURLEncoding.EncodeToString(e.Passkey.CredentialID)) → _ = l | `TestApproveByPasskey` (the allow lists) | KILLED | KILLED |
| MUT-LINK-GET-RP | M6f | approve_passkey.go: l = &ceremonyLink{V: 1, Op: "get", RP: e.Passkey.RPID, → l = &ceremonyLink{V: 1, Op: "get", RP: "localhost", | `TestApproveByPasskey` (the tailnet link) | KILLED | KILLED |
| MUT-LINK-EXPIRES | M6f | approve_passkey.go: Request: row.RequestID, Expires: row.ExpiresAt, → Request: row.RequestID, Expires: row.CreatedAt, | `TestApproveByPasskey` (the links) | KILLED | KILLED |
| MUT-LINK-REQUEST | M6f | approve_passkey.go: Request: row.RequestID, Expires: row.ExpiresAt, → Request: "", Expires: row.ExpiresAt, | `TestApproveByPasskey` (the links) | KILLED | KILLED |
| MUT-LINK-CHALLENGE-H | M6f | approve_passkey.go: Hash16: hash16, Challenge: base64.RawURLEncoding.EncodeToString(h[:])} → Hash16: hash16, Challenge: base64.RawURLEncoding.EncodeToString(h[:16])} | `TestApproveByPasskey` (the links) | KILLED | KILLED |
| MUT-HASH-EXPIRY | M6f | approve_passkey.go: return approval.MessageHash([32]byte(stored), [32]byte(challenge), row.ExpiresAt) → return approval.MessageHash([32]byte(stored), [32]byte(challenge), row.ExpiresAt+1) | `TestApproveByPasskey` (the links (the test hashes the row itself)) | KILLED | KILLED |
| MUT-WF2 | M6g | approve_passkey.go: if err != nil { // WF2 → if false { // WF2 | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF2-fifo) | KILLED | KILLED |
| MUT-WF2-CAP | M6g | approve_passkey.go: file, err := readRegularFileCapped(path, approval.MaxPasskeyResponseBytes) → file, err := readRegularFileCapped(path, 1<<20) | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF2-over-16k) | KILLED | KILLED |
| MUT-WF3 | M6g | approve_passkey.go: if err != nil { // WF3 → if false { // WF3 | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF3-not-an-assertion) | KILLED | KILLED |
| MUT-WF4 | M6g | approve_passkey.go: if err != nil { // WF4 → if false { // WF4 | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF4-v5-store) | KILLED | KILLED |
| MUT-WF5 | M6g | approve_passkey.go: if !ok { // WF5 → if !ok && false { // WF5 | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF5-unknown-credential) | KILLED | KILLED |
| MUT-WF5-KEY | M6g | approve_passkey.go: if err != nil { // WF5: only a hand-edited row → if false { // WF5: only a hand-edited row | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF5-stored-key-invalid) | KILLED | KILLED |
| MUT-WF6 | M6g | approve_passkey.go: if err != nil { // WF6 → if false { // WF6 | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF6-revoked (a nil approval: the terminal opens)) | KILLED | KILLED |
| MUT-AS-NOW | M6g | approve_passkey.go: Admit: approval.AdmitInput{Now: now, → Admit: approval.AdmitInput{Now: row.ExpiresAt + 1, | `TestApproveByPasskey` (assertion) | KILLED | KILLED |
| MUT-AS-CHALLENGE | M6g | approve_passkey.go: Challenge: [32]byte(challenge),⏎Response:  file, → Challenge: [32]byte{challenge[0] ^ challenge[0]},⏎Response:  file, | `TestApproveByPasskey` (assertion) | KILLED | KILLED |
| MUT-AS-RESPONSE | M6g | approve_passkey.go: Response:  file, → Response:  nil, | `TestApproveByPasskey` (assertion) | KILLED | KILLED |
| MUT-AS-NAME | M6g | approve_passkey.go: Approver: approval.Approver{Name: a.Name, Format: a.Format, → Approver: approval.Approver{Name: "x", Format: a.Format, | `TestApproveByPasskey` (assertion (approved_by)) | KILLED | KILLED |
| MUT-AS-FORMAT | M6g | approve_passkey.go: Approver: approval.Approver{Name: a.Name, Format: a.Format, → Approver: approval.Approver{Name: a.Name, Format: "ssh", | `TestApproveByPasskey` (assertion) | KILLED | KILLED |
| MUT-AS-ENROLLED | M6g | approve_passkey.go: EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt, → EnrolledAt: 0, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt, | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF6-enrolled-after-now) | KILLED | KILLED |
| MUT-AS-REVOKED | M6g | approve_passkey.go: EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt, → EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: 0, | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF6-revoked (refused before the network)) | KILLED | KILLED |
| MUT-AS-RPID | M6g | approve_passkey.go: Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, → Passkey: &approval.Passkey{RPID: "x.ts.net", Origin: p.Origin, | `TestApproveByPasskey` (assertion) | KILLED | KILLED |
| MUT-AS-ORIGIN | M6g | approve_passkey.go: Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, → Passkey: &approval.Passkey{RPID: p.RPID, Origin: "http://localhost:7644", | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF6-another-origin) | KILLED | KILLED |
| MUT-AS-CRED-EQUIV | M6g | approve_passkey.go: CredentialID: p.CredentialID, Key: key, → CredentialID: resp.RawID, Key: key, | `TestApproveByPasskey` (equivalent: the row was found by resp.RawID exactly, so p.CredentialID is the same bytes) | SURVIVED | SURVIVED |
| MUT-AS-KEY | M6g | approve_passkey.go: CredentialID: p.CredentialID, Key: key, → CredentialID: p.CredentialID, Key: approval.COSEKey{Alg: key.Alg}, | `TestApproveByPasskey` (assertion) | KILLED | KILLED |
| MUT-AS-BE | M6g | approve_passkey.go: BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}, → BackupEligible: false, SignCount: uint32(p.SignCount)}}, | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF6-backup-eligible) | KILLED | KILLED |
| MUT-AS-COUNT | M6g | approve_passkey.go: BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}}, → BackupEligible: p.BackupEligible, SignCount: 0}}, | `TestApprovePasskeyRefusalsBeforeAnyEffect` (WF6-counter-regression (refused before the network)) | KILLED | KILLED |
| BC-CONFIRM-COMPLETE | base `ab4b2402` | approver.go: if rerr != nil \|\| !strings.HasSuffix(line, "\n") \|\| !isYes(line) { → if !isYes(line) && (rerr == nil \|\| strings.HasSuffix(line, "\n") \|\| true) { | `TestApproverEnrollDeclinedEnrolsNothing` (y-then-eof: the A1b test pinned the complete-line rule before the split, as it pins MUT-CONFIRM-COMPLETE after it) | KILLED | KILLED |
| BC-CONFIRM-OK | base `ab4b2402` | approver.go: if !ok { // N10a / K5a → if !ok && false { // N10a / K5a | `TestApproverEnrollUnansweredPromptTimesOut` (the timeout message, as MUT-CONFIRM-OK after the split) | KILLED | KILLED |
| BC-OPEN-AFTER | base `ab4b2402` | approver.go: if env.after == nil { // N9a / K4a → if false { // N9a / K4a | `TestApproverEnrollUnansweredPromptTimesOut` (the no-timer arm, as MUT-OPEN-AFTER after the split) | KILLED | KILLED |

**Named equivalents (SURVIVED, as predicted):**
- **MUT-ASK-ERR-EQUIV** (M6a): equivalent because bufio.Reader.ReadString returns a non-nil error only for a line that does not end in the delimiter.
- **MUT-ASK-NL-EQUIV** (M6a): equivalent because with a nil error ReadString's line ends in the delimiter, and a timed-out read (line "") is refused by !ok first.
- **MUT-AS-CRED-EQUIV** (M6g): equivalent because the row was found by resp.RawID exactly, so p.CredentialID is the same bytes.

**Propagation-only (named, not mutated):** the `err` returns that hand a store or law error back unchanged (`passkeyEnrolPrecheck`'s `Approver`, `ApproverByKey` and `PasskeyByCredentialID` errors, which only a driver failure reaches once N6 has shown schema 6; `OpenCredentialsReadOnly` failures, the pre-existing pattern of A1b; `requestHash`'s `MessageHash` error, which needs a negative expiry that admission has already refused); `ceremonyLink.link()`'s `json.Marshal` (it cannot fail on strings, integers and a string slice). The help texts are pinned by `TestEveryVerbHasHelp`, `TestApproverDispatchAndHelp` and the F1 tests' exact-help comparisons, not by mutants.

**Design §6's mutant list (copied), and where each lands:** K1–K16, Q1–Q8, D1–D8, J1–J7, O1–O3, E2–E4, R1–R15, P3w, W1–W15, B7–B8, S1–S9, the three `EnrolPasskey` pre-checks, C13–C18, A6 and the M0 carry-forwards — sprints 1 and 2 (landed). **WN1–WN11, WF1–WF6, the widened F1, L4–L5 — this sprint (M6a–M6h, above).** PG1's method/CSP/asset branches and the page's origin/rp, abort, timer and save branches (XB1, XB2) — sprint 4 (Appendix C).

---

## 5. Files (final `git status --short` after M6h)

| File | Sub-milestone | Change |
|---|---|---|
| `cmd/ailang-worldd/approver.go` | M6a, M6b, M6c | §2.12a split; `approverLabel`, `passkeyLabel`, `passkeyCred`, list/revoke labels; WN1, the dispatch, the help |
| `host/approval/cose.go`, `cose_stored_test.go` | M6b | `ParseStoredCOSEKey` (§1.4 P1); its test |
| `cmd/ailang-worldd/approver_label_test.go` | M6b | new: `seedPasskey`; CW9 |
| `cmd/ailang-worldd/approver_passkey.go` | M6c | new: the enrol path |
| `cmd/ailang-worldd/approver_passkey_test.go` | M6c | new: the enrol rig; CW1, CW3 |
| `cmd/ailang-worldd/main.go`, `help_test.go` | M6c, M6f, M6g | one usage line each |
| `host/store/context_roots_test.go` | M6c | +1 root entry |
| `cmd/ailang-worldd/approver_passkey_ceremony_test.go` | M6d | new: CW2, CW5 |
| `cmd/ailang-worldd/approver_passkey_refusals_test.go` | M6e | new: CW4, CW6 |
| `cmd/ailang-worldd/approve.go`, `approve_passkey.go` | M6f, M6g | the two flags, F1, the dispatches; new: links, `approveAssertion` |
| `cmd/ailang-worldd/approve_passkey_test.go` | M6f, M6g | new: CW7 |
| `cmd/ailang-worldd/approve_passkey_refusals_test.go` | M6h | new: CW8, the F1 pairs |
| this plan and `sprint_w-world-approvals-by-key-a2-s3.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each measured at the base, §1.1)

In **one** `bash` shell, from the worktree, **not under `nohup`** (§1.1). Every leg `tee`s its full
output to a log; a leg's verdict is read from its log, never from `tail` alone:

```
bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter381-row201-a2s3; cd "$W"
L=$HOME/.ailang/state/world-iter381/gates; mkdir -p "$L"
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                      # no .ail change: exactly green
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                           # compile fence for every _test.go
go test -count=1 ./host/approval/... ./host/store/ ./host/authority/ ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
# -race on the touched packages (host/daemon: no touched test this sprint, M7 moved):
go test -race -count=1 ./host/approval/... ./host/store/ 2>&1 | tee "$L/race.log"
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestPasskey|TestApproveByPasskey|TestApprovePasskey|TestApproverListAndRevokeShowPasskeys|TestApprover|TestConfirmApproverAct|TestReadTermLine|TestEveryVerbHasHelp|TestApproveFlagRules)' 2>&1 | tee -a "$L/race.log"
# the repo-wide inventories (design §7), the reflective tests and the quickstart tests, by name, -v:
go test -v -count=1 ./host/broker/ -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^TestNoParallelWriteForkPackagesOutsideVerifygate$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/ ./host/store/ ./host/authority/ -run '^(TestSignedApprovalFieldsAreAllBound|TestEnrolmentPromptShowsEveryPersistedField|TestRequestConsumeCompleteness|TestApproverRowMatchesTheTable|TestPasskeyRowMatchesTheTable|TestRequestMintFieldsAreSourced)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run '^(TestEveryVerbHasHelp|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/broker/ -run '^(TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux \
  TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup \
  TestEveryCommandDirectoryIsAProtectedGroup TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads \
  TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestNoParallelWriteForkPackagesOutsideVerifygate \
  TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestSignedApprovalFieldsAreAllBound \
  TestEnrolmentPromptShowsEveryPersistedField TestRequestConsumeCompleteness TestApproverRowMatchesTheTable \
  TestPasskeyRowMatchesTheTable TestRequestMintFieldsAreSourced TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 \
  TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI \
  TestQuickstartServiceSectionMatchesTheCLIHelp TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
echo PRE_COMMIT_GATES_GREEN | tee "$L/pre_commit.ok"
EOF
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter381-row201-a2s3 && VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter381/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter381/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter381-row201-a2s3 status --short   # only the two plan files, until the controller commits them
```

**At the base**: every pre-commit leg above is green at `ab4b2402` (§1.1; the 30 named inventory,
reflective, help and quickstart tests all `--- PASS` there). The after-commit `verify_go.sh` was red at
the base on `host/verifygate`'s `TestMissionConfigDefaultFlowIntegration` alone (a 10 s fixture
deadline under load, §1.1). On the final tree (§1.2 C25–C28) the pre-commit script is green in 518 s with all 30 named tests `--- PASS`; the after-commit `verify_go.sh` is green on every touched package and red only on three untouched packages' timing tests, which fail alone at the base under the same machine load (C27). The same script re-run at the base during that load (load average 23–60, `logs/base_sprint_gates2`, 619 s) was red in its package leg on `host/daemon` timing tests (`TestSessionSelfGit`, `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock`) with `cmd/ailang-worldd` at 430 s, while the first baseline (§1.1, lower load) had every package of that leg green in 280 s: **run the gates on a quiet machine, and re-run a red untouched package alone once.**
**Known load- or environment-sensitive (re-run once, alone; do not chase)**:
`TestMissionConfigDefaultFlowIntegration/invalid_ledger` (and `/valid_ledger` under heavy load),
`TestSeedRestartIgnoresASurvivingOldCopy/survivor` (row 259), `host/broker`'s stray-loopback tests and
its local `-race`, the local `-race` of `host/daemon` and `host/verifygate`,
`TestWorktreeSnapshotRefusesAConfiguredFilter`, and anything run under `nohup`. None is in a package this
sprint changes a production file of. Any other red, any race report, or any red in
`cmd/ailang-worldd`, `host/approval` or `host/store` is a regression. CI's `go-race-verify` is the race
verdict.

**Must stay green, unedited**: every a0a1, A1b, sprint-1 and sprint-2 test, except the two named
edits (`help_test.go`'s usage list, `context_roots_test.go`'s +1 root). In particular CA1–CA10,
CR1–CR13, `TestConfirmApproverAct`, `TestApproverEnrollRefusals` (its `webauthn` arm unedited, §1.4 P4),
`TestApproverList`, `TestApproverRevoke`, `TestApproveFlagRules`, every `TestApproveBySignature*` and
`TestApproveSignWith*`, and every §7 inventory.

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **Runs.** Run 1 = M6a, M6b; Run 2 = M6c, M6d; Run 3 = M6e, M6f; Run 4 = M6g, M6h. Each run reads §0,
  §2 and its own §3 sections, then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an
  uncommitted tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter381-row201-a2s3 && git diff --numstat ab4b2402116243c0e0fab9ed5cee8ec80a660a4e HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M6a **38**, M6b **237**, M6c **701**, M6d **945**, M6e **1190**,
  M6f **1397**, M6g **1538**, M6h **1686**. **Any other value means a patch was not applied as written:
  stop and report.**
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M6a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter381-row201-a2s3
    D=$HOME/.ailang/state/world-iter381/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a2-sprint3-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a2-s3.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M6a` … `M6h`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any
  other status. The controller rebuilds one commit per sub-milestone (iteration 377's `commit_ms.sh`
  with `WT` and the snapshot root pointed here).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- SKIP` or `--- FAIL` in any `-v` run; any need to touch a file outside §5;
  running Chrome, Node, `ssh-keygen` or any browser; any reference to `~/.ssh`, `~/.ailang/world` or a
  real `/dev/tty`. Stop and report; do not work around.
- **Timing.** Every new test runs in under 1 s (the CW7/CW8 rigs make a request each: ~0.1 s); M6a's,
  M6c's and M6h's whole `cmd/ailang-worldd` ~2 min; **[cli]** ~1 min; the compile fence ~25 s.
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does
  not word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a
  file; never pipe a gate's only copy through `tail`. **Never start a gate under `nohup`** (§1.1). Poll
  any background command's artifact in a bounded `date +%s` loop in the same turn; never end a turn
  announcing a wait.

---

## 8. Risks

- **R1 — Size.** 1,686 of 2,000, measured; the per-step expected numstat makes drift visible at once.
- **R2 — Between this sprint and sprint 4 the printed links lead to a 404.** `approver enroll --format
  webauthn` and `approve --webauthn-link` print `<origin>/surface/passkey#…`, which the daemon serves
  only from sprint 4 (C10). This is safe, not merely inert: an operator cannot produce a registration
  for the printed challenge at the daemon's origin without the page, so no passkey can enrol; nothing
  enrols without the typed code and `y` (CW3–CW6) and nothing approves without a verified assertion
  (CW7, CW8). The help says what the page does; QUICKSTART §16 waits for the page (§1.4 P8). The
  attended exit run (design §10) needs sprint 4.
- **R3 — The enrolment code is typed against a file the operator's browser saved** (design R1, §2.8):
  unchanged by this plan; CW5 pins the read-then-substitute attacker.
- **R4 — The test terminal is a model of a canonical-mode terminal** (one line per read). A real
  `/dev/tty` delivers one line per `read(2)` in canonical mode, so two `bufio.Reader`s in sequence
  (one per prompt) see one line each; a paste of several lines at the code prompt would feed its next
  line to the y/N. Only the operator's own paste can do that, and the y/N's text is on the terminal
  before that line is read; the attended run is the real check.
- **R6 — A malformed `--name` refuses only after the ceremony** (S3-D13): `ailang-worldd approver enroll
  --name Mark …` prints the link, waits for the code, and then refuses with `ErrApproverName` (WN9),
  enrolling nothing. Safe, but it costs the operator a browser ceremony; sprint 4 may add the grammar
  check beside WN2 (Appendix C).
- **R5 — Haiku and patches.** The executor transcribes nothing; a skipped patch fails `--dry-run` or the
  numstat; a hand "fix" after a red test is forbidden by §7.

---

## 9. Acceptance criteria (the judge scores against these)

1. **§2.12a**: `confirmApproverAct` is `openApproverTerminal` + `confirmOnTerminal` (`askApproverLine`)
   with the same messages and exit codes; CA3, CA4, CA9, CA10 and `TestConfirmApproverAct` pass
   **unedited**; its mutants die by those tests (MUT-OPEN-*, MUT-ASK-COMPLETE, MUT-CONFIRM-*).
2. **WN1–WN4, N3, N6, N7** refuse before any terminal, print no link and write nothing, one arm per
   operand of WN1 and WN2 (CW1); with no terminal enrol refuses with the tty fix and prints no link
   (CW3).
3. **WN5** refuses a missing, failed or all-zero challenge source with `exitFatal` before any terminal
   (CW6; S3-D7 records the zero-challenge decision and why).
4. **The ceremony**: the link's fragment carries exactly the drawn challenge, the origin's rpID and the
   name; the code prompt is bounded by exactly 300 s and the y/N by 120 s, on one terminal; the writer
   lock is free at both prompts; the right code and `y` store exactly the registration's key (canonical),
   credential, AAGUID, algorithm, BE and counter, and N14 names the passkey (CW2, CW4).
5. **WN10** refuses a substituted registration even when the genuine code is typed, and a code wrong in
   its last character, first group, length or alphabet, with the y/N never shown; it accepts the genuine
   code typed in lowercase with dashes, spaces and a tab, or as one string (CW5).
6. **WN8, WN9, WN11, N11–N13**: a FIFO or oversize file, a registration the law refuses, a key or
   credential stored by another enrol while the code prompt waited, each refused before the y/N; the
   write's races refuse and its failures are fatal; nothing enrolled each time (CW6).
7. **L4, L5**: `approver list` labels passkeys (revoked included) and a hand-made row exactly; `revoke`
   shows the label; a schema-5 store is asked for no passkey (CW9).
8. **`--webauthn-link`**: one link per origin with the request, expiry, `hash16`, challenge and the
   active credentials; no terminal, no row, no daemon call; WF1 and the schema-6 refusal (CW7).
9. **`--webauthn-assertion`**: mints with no terminal, records `approved_by = webauthn:<name>`, advances
   the counter, leaves the effects the typed `y` leaves; a replay is refused; revoke lists it (CW7).
   **WF2–WF6** each refuse before any daemon call, with the request pending and the counter unchanged
   (CW8, 14 arms); F1 refuses every pair of the five modes and F2 `--list` with either (CW8).
10. **`ParseStoredCOSEKey`** reads back exactly what `Canonical()` stored and refuses anything else
    (`TestParseStoredCOSEKey`); `host/approval` stays standard-library-only and pure.
11. Each §4 mutant, applied alone, makes its named test **and** its package suite fail; the named
    equivalents survive for the stated reasons; the base controls behave as stated; the ledger has a line
    per run.
12. Every §6 gate passes, including every inventory, reflective, help and quickstart test by name, the
    `-race` legs, and `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after the commit (modulo §6's
    re-run rule); no test skips; `git status --short` lists only §5.
13. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000; M7 and
    QUICKSTART §16 are in Appendix C, not in this sprint.
14. No real key, browser profile, authenticator, store, `/dev/tty`, `~/.ssh` access, mint or publish
    anywhere in the sprint (D-WORLD-109).

---

## Appendix C — sprint 4 (moved from this sprint: M7 and QUICKSTART §16)

Design §9's sprint-3 fallback fired by measurement (§1.5). **Sprint 4 = M7, "the page, the live browser
and QUICKSTART"**, exactly as design §2.14, §2.15 and §6 (M7) specify, unchanged:
- **The page**: `host/surfacepage/passkey/{index.html, passkey.js}` (`go:embed` in
  `host/surfacepage/passkey.go`), served by `host/daemon` at `GET /surface/passkey` and `GET
  /surface/passkey.js` with `setSurfaceHeaders` and **the unchanged `surfaceCSP`**; no store, no
  request, no `fetch`/`XMLHttpRequest`/form `action`, `textContent` only; the [r2] AbortController with
  the fixed 240000 ms timer, `signal` beside `publicKey`, the button disabled while pending, a late
  result ignored, saving disabled after expiry; the §2.7 JSON built by the page itself; the grouped
  enrolment code; the `blob:` Save link (design V31).
- **The tests**: PG1 `TestPasskeyPageRoutes`; the `TestSurfaceAppRoutes` edit (+2 names); **XB1**
  `TestPasskeyPageInChrome` and **XB2** `TestPasskeyPageAbortsAPendingCeremony` (gated: `WORLD_CHROME`
  + `node` on `PATH`; `t.Skip("SKIP COUNTED: …")` unless `WORLD_REQUIRE_PASSKEY_BROWSER=1`, then
  `t.Fatal`), run ≥ 3× in a row by that sprint's planner with timings; `testdata/passkey_cdp.mjs`;
  `host/approval/testdata/webauthn/generate.sh` made real (it is the sprint-1 stub today).
- **QUICKSTART §16** "Approve with a passkey" *(attended — pending first verbatim run)* (§1.4 P8).

What sprint 4 inherits from this sprint, to use as built:
- The links it must render: `ceremonyLink` (`cmd/ailang-worldd/approver_passkey.go`) —
  `{"v":1,"op":"create","rp","origin","name","challenge"}` for enrolment and
  `{"v":1,"op":"get","rp","origin","request","expires","hash16","challenge","allow":[…]}` for approval,
  JSON in that member order, base64url (no padding) in the fragment of `<origin>/surface/passkey#…`;
  `challenge` and each `allow` entry base64url; `expires` Unix seconds; `hash16` the first 16 hex of the
  message hash. Sprint 4's XB1 should build its links **with the production code** (call
  `runApproverEnrollPasskey`/`writePasskeyLinks` through the CLI, or `ceremonyLink.link()` from a
  `cmd` test) rather than re-encoding them, so the page and the CLI cannot drift.
- The file names the page suggests are the CLI's: `world-registration-<name>.json` (CW tests use it)
  and `world-assertion-<id>.json`.
- **XB2's seam is measured** (§1.2 C18): `Emulation.setVirtualTimePolicy({policy: "advance", budget:
  250000})` after the click advances the page's 240000 ms timer to exactly 240000 ms of page time
  (`Date.now()` and `performance.now()` deltas both 240000), fires `Emulation.virtualTimeBudgetExpired`,
  and costs 53–55 ms of real time (3 of 3, Chrome 157, the exact `surfaceCSP`); the pending operation is
  injected with `Page.addScriptToEvaluateOnNewDocument` replacing `navigator.credentials.create` by a
  promise that rejects with `signal.reason` on `abort`. The probe is
  `$HOME/.ailang/state/world-iter381/plan-scratch/probe/vt_probe.mjs` (68 lines, temp profile, no
  authenticator); copy it into that plan's scratch if this one is rotated.
- **XB1's starting points** (C17): the iteration-379 recording probe (`/tmp/w379-designer/browser/
  probe.mjs`, sha256 `01a17c0d…37cf`: `WebAuthn.enable`, `WebAuthn.addVirtualAuthenticator` with
  `ctap2_1`, `internal`, resident key, UV; create and two gets) and the V31 blob probe
  (`/tmp/w379-designer/blob/probe.mjs`: `Browser.setDownloadBehavior` with `eventsEnabled`, a
  DevTools-dispatched mouse click on the Save link, `downloadWillBegin`/`downloadProgress completed`, a
  `securitypolicyviolation` listener). Both are in `/tmp` and may be gone; their behaviour is recorded
  in design V20–V22, V31 and here.
- The AC10 census skips `_test.go` files (C11), so a test-only Chrome/Node launcher adds no production
  exec site; `TestFrozenV1RouteTableMatchesMux` stays at 13 (the page adds no `/v1` route).
- **A cheap improvement sprint 4 may take** (§8 R6, S3-D13): refuse a `--name` that fails
  `approval.ApproverNamePattern` beside WN2, before the link is printed (one branch, one CW1 arm).
- The size: M7 was not built here; its parts are listed in §1.5. Sprint 4 has the whole 2,000-line cap.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M6a — one terminal, two bounded reads (§2.12a)

````diff M6a
--- a/cmd/ailang-worldd/approver.go
+++ b/cmd/ailang-worldd/approver.go
@@ -29,25 +29,49 @@
 // confirmApproverAct shows prompt at /dev/tty and returns exitOK only on a y typed within
 // approverConfirmBudget, as a complete line with no read error (design §2.10a [r2]). Every other
 // outcome returns its refusal code having written nothing. It is called with no store handle
-// open, so no writer lock is held while it waits.
+// open, so no writer lock is held while it waits. Row 201 A2 (§2.12a) splits it into
+// openApproverTerminal and askApproverLine, so passkey enrolment asks twice on one terminal.
 func confirmApproverAct(env sessionEnv, stderr io.Writer, prefix, prompt, nothing string) int {
+	term, code := openApproverTerminal(env, stderr, prefix, nothing)
+	if code != exitOK {
+		return code
+	}
+	defer func() { _ = term.Close() }() // closed on every return
+	return confirmOnTerminal(term, env, stderr, prefix, prompt, nothing)
+}
+
+// openApproverTerminal opens /dev/tty for an approver act (N9a, N9 / K4a, K4): the terminal, or
+// the refusal code having opened nothing. The caller closes the terminal.
+func openApproverTerminal(env sessionEnv, stderr io.Writer, prefix, nothing string) (io.ReadWriteCloser, int) {
 	if env.after == nil { // N9a / K4a
 		fmt.Fprintf(stderr, "%sinternal error: no prompt timer (sessionEnv.after)\n", prefix)
-		return exitFatal
+		return nil, exitFatal
 	}
 	term, err := env.openTerminal()
 	if err != nil { // N9 / K4
 		fmt.Fprintf(stderr, "%srefusing: no controlling terminal (%v); %s. %s\n", prefix, err, nothing, sessionTTYFix)
-		return exitUsage
+		return nil, exitUsage
 	}
-	defer func() { _ = term.Close() }() // closed on every return
-	fmt.Fprint(term, prompt)            // the budget starts here
-	line, rerr, ok := readTermLineErrWithin(term, approverConfirmBudget, env.after)
+	return term, exitOK
+}
+
+// askApproverLine writes prompt to term and reads one line within budget, counted from the
+// prompt: ok is false when after(budget) fires first; complete is true only for a line that ends
+// in a newline and was read with no error.
+func askApproverLine(term io.ReadWriter, prompt string, budget time.Duration, after func(time.Duration) <-chan time.Time) (line string, ok, complete bool) {
+	fmt.Fprint(term, prompt) // the budget starts here
+	line, err, ok := readTermLineErrWithin(term, budget, after)
+	return line, ok, err == nil && strings.HasSuffix(line, "\n")
+}
+
+// confirmOnTerminal is the y/N on an open terminal (N10a, N10 / K5a, K5).
+func confirmOnTerminal(term io.ReadWriter, env sessionEnv, stderr io.Writer, prefix, prompt, nothing string) int {
+	line, ok, complete := askApproverLine(term, prompt, approverConfirmBudget, env.after)
 	if !ok { // N10a / K5a
 		fmt.Fprintf(stderr, "%sno answer within %s; %s\n", prefix, approverConfirmBudget, nothing)
 		return exitUsage
 	}
-	if rerr != nil || !strings.HasSuffix(line, "\n") || !isYes(line) { // N10 / K5: EOF, read error, partial line, not y
+	if !complete || !isYes(line) { // N10 / K5: EOF, read error, partial line, not y
 		fmt.Fprintf(stderr, "%saborted (not confirmed); %s. fix: run it again and type y at the prompt\n", prefix, nothing)
 		return exitUsage
 	}
````

### Patch M6b — `list` and `revoke` show passkeys; `ParseStoredCOSEKey`

````diff M6b
--- a/cmd/ailang-worldd/approver.go
+++ b/cmd/ailang-worldd/approver.go
@@ -7,6 +7,8 @@
 import (
 	"bufio"
 	"context"
+	"crypto/sha256"
+	"encoding/hex"
 	"errors"
 	"flag"
 	"fmt"
@@ -256,8 +258,8 @@
 	if err != nil {
 		return approverRefuse(stderr, prefix, "cannot read the credentials store %s: %v", credPath, err)
 	}
+	defer func() { _ = ro.Close() }()
 	rows, err := ro.ListApprovers(ctx)
-	_ = ro.Close()
 	if err != nil { // L1: below schema 5 names the serve restart
 		return approverRefuse(stderr, prefix, "%v", err)
 	}
@@ -265,7 +267,13 @@
 		fmt.Fprintln(stdout, "no approvers")
 		return exitOK
 	}
-	for _, r := range rows { // L3
+	labels := make([]string, len(rows))
+	for i, r := range rows {
+		if labels[i], err = approverLabel(ctx, ro, r); err != nil {
+			return approverRefuse(stderr, prefix, "%v", err)
+		}
+	}
+	for i, r := range rows { // L3
 		used := "never"
 		if r.LastUsedAt != 0 {
 			used = approverTime(r.LastUsedAt)
@@ -275,7 +283,7 @@
 			state = "revoked " + approverTime(r.RevokedAt)
 		}
 		fmt.Fprintf(stdout, "%s  %s  %s  enrolled %s  last used %s  %s\n",
-			r.Name, r.Format, approverKeyLabel(r.PublicKey), approverTime(r.EnrolledAt), used, state)
+			r.Name, r.Format, labels[i], approverTime(r.EnrolledAt), used, state)
 	}
 	return exitOK
 }
@@ -288,8 +296,36 @@
 		keyType = k.Type
 	}
 	return keyType + " " + approval.SSHPublicKey{Blob: blob}.Fingerprint()
+}
+
+// approverLabel is row's key label (row 201 A2 §2.12b): an ssh row's approverKeyLabel, or a
+// webauthn row's passkey read on ro, revoked ones included. A webauthn row exists only at schema 6,
+// so a schema-5 store is never asked for a passkey.
+func approverLabel(ctx context.Context, ro *store.Credentials, row store.ApproverRow) (string, error) {
+	if row.Format != approval.FormatWebAuthn {
+		return approverKeyLabel(row.PublicKey), nil
+	}
+	p, ok, err := ro.Passkey(ctx, row.Name)
+	if err != nil {
+		return "", err
+	}
+	if !ok { // L5: only a hand-made row
+		return "invalid-passkey", nil
+	}
+	return passkeyLabel(p.Alg, p.RPID, p.CredentialID), nil // L4
 }
 
+// passkeyLabel is "passkey <ES256|EdDSA> <rpID> cred <first 16 hex of SHA-256(credential ID)>".
+func passkeyLabel(alg int64, rpID string, credentialID []byte) string {
+	return "passkey " + approval.COSEKey{Alg: alg}.AlgName() + " " + rpID + " cred " + passkeyCred(credentialID)
+}
+
+// passkeyCred is the first 16 hex of SHA-256(credential ID), as the enrolment y/N shows it.
+func passkeyCred(credentialID []byte) string {
+	h := sha256.Sum256(credentialID)
+	return hex.EncodeToString(h[:8])
+}
+
 // approverTime is a Unix time in RFC 3339, UTC.
 func approverTime(v int64) string { return time.Unix(v, 0).UTC().Format(time.RFC3339) }
 
@@ -428,6 +464,10 @@
 		return approverRefuse(stderr, prefix, "cannot read the credentials store %s: %v", credPath, err)
 	}
 	row, ok, err := ro.Approver(pre, name)
+	label := ""
+	if err == nil && ok {
+		label, err = approverLabel(pre, ro, row)
+	}
 	_ = ro.Close()
 	cancelPre()
 	if err != nil { // below schema 5 names the serve restart
@@ -442,7 +482,7 @@
 
 	// The fence (K4a, K4, K5a, K5).
 	prompt := fmt.Sprintf("Revoke approver %s (%s)? It will never approve again, and its name and key can never be enrolled again. [y/N] ",
-		name, approverKeyLabel(row.PublicKey))
+		name, label)
 	if code := confirmApproverAct(env, stderr, prefix, prompt, "nothing was revoked"); code != exitOK {
 		return code
 	}
--- a/cmd/ailang-worldd/approver_label_test.go
+++ b/cmd/ailang-worldd/approver_label_test.go
@@ -0,0 +1,104 @@
+package main
+
+// Row 201 A2 sprint 3 M6b (design_docs/planned/w-world-approvals-by-key-a2.md §2.12b, CW9): approver
+// list and revoke show a passkey approver's passkey. Every store is under t.TempDir(); every key is
+// approvaltest.Passkey's fixed seed (D-WORLD-109).
+
+import (
+	"bytes"
+	"context"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// The test passkeys' origin and relying party ID: the daemon's own page on this machine.
+const pkOrigin, pkRPID = "http://localhost:7644", "localhost"
+
+// seedPasskey enrols name with pk's key at `at` straight through the store (no verb, no terminal),
+// as approver enroll --format webauthn stores it.
+func seedPasskey(t *testing.T, creds, name string, pk *approvaltest.Passkey, at int64) {
+	t.Helper()
+	st, err := store.OpenCredentials(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = st.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	k := pk.Key()
+	if err := st.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: at},
+		store.PasskeyRow{Name: name, RPID: pkRPID, Origin: pkOrigin, CredentialID: pk.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: 1}); err != nil {
+		t.Fatal(err)
+	}
+}
+
+// CW9: list labels a passkey (L4) and a hand-made webauthn row with no passkey (L5); revoke shows
+// the passkey's label; a schema-5 store lists its ssh approvers and is asked for no passkey.
+func TestApproverListAndRevokeShowPasskeys(t *testing.T) {
+	es, ed := approvaltest.NewPasskey(approval.AlgES256, 21), approvaltest.NewPasskey(approval.AlgEdDSA, 22)
+	esLabel := "passkey ES256 localhost cred " + passkeyCred(es.CredentialID())
+	db, creds := approverStore(t)
+	seedPasskey(t, creds, "a", es, approveT0)
+	seedPasskey(t, creds, "b", ed, approveT0)
+	revokeSeeded(t, creds, "b", approveT0+200)
+	raw := rawCredentials(t, creds)
+	if _, err := raw.Exec(`INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES ('c', 'webauthn', x'a0', 1800000000)`); err != nil {
+		t.Fatal(err)
+	}
+	_ = raw.Close()
+	var out, errw bytes.Buffer
+	if code := runApprover([]string{"list", "--db", db}, &out, &errw, sessionEnv{}); code != exitOK {
+		t.Fatalf("list: code %d; stderr %q", code, errw.String())
+	}
+	want := "a  webauthn  " + esLabel + "  enrolled 2027-01-15T08:00:00Z  last used never  active\n" +
+		"b  webauthn  passkey EdDSA localhost cred " + passkeyCred(ed.CredentialID()) + "  enrolled 2027-01-15T08:00:00Z  last used never  revoked 2027-01-15T08:03:20Z\n" +
+		"c  webauthn  invalid-passkey  enrolled 2027-01-15T08:00:00Z  last used never  active\n"
+	if out.String() != want {
+		t.Fatalf("list:\n%q\nwant:\n%q", out.String(), want)
+	}
+	if len(passkeyCred(es.CredentialID())) != 16 || passkeyCred(es.CredentialID()) == passkeyCred(ed.CredentialID()) {
+		t.Fatalf("cred labels %q %q: want 16 distinct hex", passkeyCred(es.CredentialID()), passkeyCred(ed.CredentialID()))
+	}
+
+	term := &lockedTerm{in: strings.NewReader("y\n")}
+	out.Reset()
+	if code := runApprover([]string{"revoke", "a", "--db", db}, &out, &errw, approveEnv(term, nil, atTime(approveT0+300))); code != exitOK {
+		t.Fatalf("revoke: code %d; stderr %q", code, errw.String())
+	}
+	if wantPrompt := "Revoke approver a (" + esLabel + ")? "; !strings.HasPrefix(term.Output(), wantPrompt) {
+		t.Fatalf("revoke prompt %q, want the prefix %q", term.Output(), wantPrompt)
+	}
+	if out.String() != "revoked approver a\nno session was approved by it\n" {
+		t.Fatalf("revoke stdout %q", out.String())
+	}
+
+	// A schema-5 store: list asks for no passkey of its ssh approvers, and a hand-made webauthn row
+	// there names the restart.
+	db5, creds5 := approverStore(t)
+	seedApprover(t, creds5, "s", fixtureKeyOf(t, "ed25519.pub"), approveT0)
+	raw = rawCredentials(t, creds5)
+	for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+		if _, err := raw.Exec(stmt); err != nil {
+			t.Fatal(err)
+		}
+	}
+	_ = raw.Close()
+	out.Reset()
+	if code := runApprover([]string{"list", "--db", db5}, &out, &errw, sessionEnv{}); code != exitOK || !strings.HasPrefix(out.String(), "s  ssh  ssh-ed25519 "+edFingerprint) {
+		t.Fatalf("v5 list: code %d, stdout %q; stderr %q", code, out.String(), errw.String())
+	}
+	raw = rawCredentials(t, creds5)
+	if _, err := raw.Exec(`INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES ('w', 'webauthn', x'a0', 1800000000)`); err != nil {
+		t.Fatal(err)
+	}
+	_ = raw.Close()
+	errw.Reset()
+	if code := runApprover([]string{"list", "--db", db5}, &out, &errw, sessionEnv{}); code != exitUsage || !strings.Contains(errw.String(), "passkeys need schema 6") {
+		t.Fatalf("v5 list with a webauthn row: code %d, stderr %q; want exitUsage naming schema 6", code, errw.String())
+	}
+}
--- a/host/approval/cose.go
+++ b/host/approval/cose.go
@@ -85,6 +85,17 @@
 	return k, nil
 }
 
+// ParseStoredCOSEKey parses a key as approvers.public_key holds it (row 201 A2 sprint 3: approve
+// --webauthn-assertion rebuilds the enrolled passkey from its store row): exactly one COSE key,
+// judged by Q1-Q8, with nothing after it.
+func ParseStoredCOSEKey(b []byte) (COSEKey, error) {
+	v, err := decodeCBORExact(b, attestationLimits)
+	if err != nil {
+		return COSEKey{}, fmt.Errorf("%w: %w", ErrCOSEKey, err)
+	}
+	return ParseCOSEKey(v)
+}
+
 func (k COSEKey) point() []byte { return append(append([]byte{4}, k.X...), k.Y...) }
 
 // Canonical is the key's CTAP2 canonical encoding (labels 1, 3, -1, -2[, -3]): what is stored.
--- a/host/approval/cose_stored_test.go
+++ b/host/approval/cose_stored_test.go
@@ -0,0 +1,36 @@
+package approval
+
+// Row 201 A2 sprint 3 M6b (plan §1.4): ParseStoredCOSEKey, which approve --webauthn-assertion
+// rebuilds an enrolled passkey's key with, reads back exactly what Canonical() stored.
+
+import (
+	"bytes"
+	"errors"
+	"reflect"
+	"testing"
+)
+
+func TestParseStoredCOSEKey(t *testing.T) {
+	for _, alg := range []string{"es256", "eddsa"} {
+		k, raw := chromeKey(t, alg)
+		got, err := ParseStoredCOSEKey(k.Canonical())
+		if err != nil || !reflect.DeepEqual(got, k) {
+			t.Fatalf("%s: %+v, %v; want Chrome's key back", alg, got, err)
+		}
+		for name, arm := range map[string]struct {
+			b    []byte
+			want error
+		}{
+			"trailing":  {append(bytes.Clone(raw), 0), ErrCBORTrailing},
+			"truncated": {raw[:len(raw)-1], ErrCBORTruncated},
+			"empty":     {nil, ErrCBORTruncated},
+		} {
+			if _, err := ParseStoredCOSEKey(arm.b); !errors.Is(err, ErrCOSEKey) || !errors.Is(err, arm.want) {
+				t.Fatalf("%s %s: %v, want ErrCOSEKey wrapping %v", alg, name, err, arm.want)
+			}
+		}
+	}
+	if _, err := ParseStoredCOSEKey(cb(cbMap{1, 3, 3, -257, -1, 1})); !errors.Is(err, ErrCOSEAlg) {
+		t.Fatalf("an RS256 key: %v, want ErrCOSEAlg", err)
+	}
+}
````

### Patch M6c — `approver enroll --format webauthn`

````diff M6c
--- a/cmd/ailang-worldd/approver.go
+++ b/cmd/ailang-worldd/approver.go
@@ -142,19 +142,27 @@
 // ---- M3c ----
 
 const approverHelp = `usage: ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
+       ailang-worldd approver enroll --name <name> --format webauthn --origin <origin> --registration <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approver list [--db <path>] [--credentials-db <path>]
        ailang-worldd approver revoke <name> [--db <path>] [--credentials-db <path>]
 
-ATTENDED (enroll, revoke). An approver is a name and one SSH public key, kept in
-the credentials store (row 201). An enrolled, unrevoked key approves a mint
-request with its signature: 'ailang-worldd approve <request-id> --signature
-<file>' or '--sign-with <key>'. enroll and revoke each ask y at /dev/tty and wait at most 2m0s for it;
+ATTENDED (enroll, revoke). An approver is a name and one SSH public key or one
+passkey, kept in the credentials store (row 201). An enrolled, unrevoked key
+approves a mint request with its signature: 'ailang-worldd approve <request-id>
+--signature <file>' or '--sign-with <key>', or for a passkey '--webauthn-link'
+then '--webauthn-assertion <file>'. enroll and revoke each ask y at /dev/tty and wait at most 2m0s for it;
 an agent cannot run them.
 
 enroll   reads --pubkey once (one ssh-ed25519 or ecdsa-sha2-nistp256 key line,
          at most 8 KiB; a private key file is refused and never shown), shows
          the key's SHA256 fingerprint, and stores it after y. A name and a key
          are each enrolled at most once, ever, revoked ones included.
+enroll --format webauthn prints a link to the daemon's page /surface/passkey
+         at --origin; open it in a browser, create the passkey (Face ID,
+         Touch ID, a PIN or a security key), and save the file it offers as
+         --registration, which must not exist yet. Then type the 26-character
+         enrolment code the page shows, within 5m0s (a fingerprint of that
+         file, not a secret), check what the y/N shows, and type y.
 list     prints every approver, read-only: name, format, key type,
          fingerprint, enrolled, last used, and active or revoked.
 revoke   marks the approver revoked after y. The row is kept: it never
@@ -164,8 +172,12 @@
 
   --name <name>             1-48 of a-z 0-9 . _ - (the first a letter or digit);
                             its approvals record approved_by=ssh:<name>
-  --format ssh              the only format until row 201 A2 (webauthn)
+  --format ssh|webauthn     an SSH key (--pubkey) or a passkey (--origin, --registration);
+                            a passkey's approvals record approved_by=webauthn:<name>
   --pubkey <file>           the .pub file (compare: ssh-keygen -l -f <file>)
+  --origin <origin>         http://localhost:<port> or https://<host>.ts.net, as the
+                            daemon is reached (never an IP address); the passkey is bound to it
+  --registration <file>     where the browser saves the registration (a new file)
   --db <path>               the world store (default ~/.ailang/world/world.db)
   --credentials-db <path>   default credentials.db beside --db
 The global --addr is refused: approvers live in the credentials store.
@@ -338,11 +350,13 @@
 // runApproverEnroll is `approver enroll` (design §2.10 N1-N14, plan A-D6 order).
 func runApproverEnroll(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 	const prefix = "ailang-worldd approver enroll: "
-	var name, format, pubkey *string
+	var name, format, pubkey, origin, registration *string
 	db, creds, positional, code, done := approverFlags("enroll", args, stdout, stderr, func(fs *flag.FlagSet) {
 		name = fs.String("name", "", "the approver's name")
-		format = fs.String("format", "", "ssh")
+		format = fs.String("format", "", "ssh or webauthn")
 		pubkey = fs.String("pubkey", "", "the public key file (.pub)")
+		origin = fs.String("origin", "", "webauthn: the ceremony page's origin")
+		registration = fs.String("registration", "", "webauthn: the registration file the browser saves")
 	})
 	if done {
 		return code
@@ -351,6 +365,14 @@
 		fmt.Fprint(stderr, approverHelp)
 		return exitUsage
 	}
+	webauthn := *format == approval.FormatWebAuthn
+	if webauthn && *pubkey != "" || !webauthn && (*origin != "" || *registration != "") { // WN1 (row 201 A2 §2.12)
+		fmt.Fprintf(stderr, "%s--pubkey is for --format ssh; --origin and --registration are for --format webauthn\n\n%s", prefix, approverHelp)
+		return exitUsage
+	}
+	if webauthn {
+		return runApproverEnrollPasskey(*name, *origin, *registration, db, creds, stdout, stderr, env)
+	}
 	if *name == "" || *pubkey == "" { // N2
 		return approverRefuse(stderr, prefix, "--name and --pubkey are required")
 	}
--- a/cmd/ailang-worldd/approver_passkey.go
+++ b/cmd/ailang-worldd/approver_passkey.go
@@ -0,0 +1,182 @@
+package main
+
+// Row 201 A2 sprint 3 (design_docs/planned/w-world-approvals-by-key-a2.md §2.12): `approver enroll
+// --format webauthn`. The browser's registration reaches the CLI as a file; the operator types, at
+// /dev/tty, the 26-character enrolment code the page showed (a fingerprint of that file, not a
+// secret), then y. No store handle is open while either prompt waits.
+
+import (
+	"context"
+	"crypto/subtle"
+	"encoding/base64"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"io/fs"
+	"os"
+	"strings"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// passkeyCeremonyBudget bounds the enrolment code's read, counted from its prompt (WN6): the
+// operator opens a browser, passes Face ID or Touch ID (or uses a phone), and saves a file. A
+// constant, not a flag, so a script cannot widen it.
+const passkeyCeremonyBudget = 300 * time.Second
+
+// approverEnrolPasskey is the passkey enrol write, a test seam (CW6 replaces it to inject a race);
+// production never reassigns it.
+var approverEnrolPasskey = (*store.Credentials).EnrolPasskey
+
+// ceremonyLink is the ceremony page's fragment (design §2.12 step 6, §2.13): what one create or
+// get asks. The page reads it from location.hash; it never reaches a server.
+type ceremonyLink struct {
+	V         int      `json:"v"`
+	Op        string   `json:"op"`
+	RP        string   `json:"rp"`
+	Origin    string   `json:"origin"`
+	Name      string   `json:"name,omitempty"`
+	Request   string   `json:"request,omitempty"`
+	Expires   int64    `json:"expires,omitempty"`
+	Hash16    string   `json:"hash16,omitempty"`
+	Challenge string   `json:"challenge"`
+	Allow     []string `json:"allow,omitempty"`
+}
+
+// link is <origin>/surface/passkey#<base64url(JSON)>.
+func (c ceremonyLink) link() string {
+	b, _ := json.Marshal(c) // strings, integers and a string slice: Marshal cannot fail
+	return c.Origin + "/surface/passkey#" + base64.RawURLEncoding.EncodeToString(b)
+}
+
+// runApproverEnrollPasskey is approver enroll --format webauthn (design §2.12 WN2-WN11 with N3,
+// N6, N7 and N9-N14), after runApproverEnroll's WN1.
+func runApproverEnrollPasskey(name, origin, regPath, db, creds string, stdout, stderr io.Writer, env sessionEnv) int {
+	const prefix, nothing = "ailang-worldd approver enroll: ", "nothing was enrolled"
+	if name == "" || origin == "" || regPath == "" { // WN2
+		return approverRefuse(stderr, prefix, "--name, --origin and --registration are required with --format webauthn")
+	}
+	credPath, refusal := approverStores(db, creds) // N3
+	if refusal != "" {
+		return approverRefuse(stderr, prefix, "%s", refusal)
+	}
+	rpID, err := approval.ParseCeremonyOrigin(origin)
+	if err != nil { // WN3
+		return approverRefuse(stderr, prefix, "--origin %q: %v", origin, err)
+	}
+	if _, err := os.Lstat(regPath); !errors.Is(err, fs.ErrNotExist) { // WN4
+		return approverRefuse(stderr, prefix, "--registration %s already exists or cannot be checked; pass a new path: the browser saves the new registration there", regPath)
+	}
+	pre, cancelPre := context.WithTimeout(context.Background(), approveReadBudget)
+	refusal = passkeyEnrolPrecheck(pre, credPath, name, nil) // N6, N7, closed before any terminal
+	cancelPre()
+	if refusal != "" {
+		return approverRefuse(stderr, prefix, "%s", refusal)
+	}
+	var challenge [32]byte
+	if env.random == nil || env.random(challenge[:]) != nil || challenge == [32]byte{} { // WN5
+		fmt.Fprintf(stderr, "%sinternal error: no fresh challenge (sessionEnv.random is missing, failed, or gave 32 zero bytes); %s\n", prefix, nothing)
+		return exitFatal
+	}
+	term, code := openApproverTerminal(env, stderr, prefix, nothing) // N9a, N9: before the link, which only a person can use
+	if code != exitOK {
+		return code
+	}
+	defer func() { _ = term.Close() }()
+	fmt.Fprintf(stdout, "open this link in a browser (the daemon serving %s shows the page; nothing is sent to it), create the passkey, and save the file it offers as %s:\n  %s\n",
+		origin, regPath, ceremonyLink{V: 1, Op: "create", RP: rpID, Origin: origin, Name: name, Challenge: base64.RawURLEncoding.EncodeToString(challenge[:])}.link())
+	line, ok, complete := askApproverLine(term, "Enrolment code shown by the browser (XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX): ", passkeyCeremonyBudget, env.after)
+	if !ok { // WN6
+		return approverRefuse(stderr, prefix, "no code within %s; %s", passkeyCeremonyBudget, nothing)
+	}
+	if !complete { // WN7: end of input, a read error or a partial line
+		return approverRefuse(stderr, prefix, "aborted (no complete code line); %s. fix: run it again and type the code the page shows", nothing)
+	}
+	reg, err := readRegularFileCapped(regPath, approval.MaxPasskeyResponseBytes)
+	if err != nil { // WN8
+		return approverRefuse(stderr, prefix, "--registration: %v; %s", err, nothing)
+	}
+	a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: name, Origin: origin, Challenge: challenge, Response: reg, Now: env.now()})
+	if err != nil { // WN9
+		return approverRefuse(stderr, prefix, "%s does not enrol: %v; %s", regPath, err, nothing)
+	}
+	resp, _ := approval.ParseRegistrationResponse(reg) // R3 has parsed it
+	typed := strings.NewReplacer("-", "", " ", "", "\t", "", "\n", "").Replace(strings.ToUpper(line))
+	if subtle.ConstantTimeCompare([]byte(typed), []byte(approval.EnrolmentCode(resp.ClientDataJSON, resp.AttestationObject))) != 1 { // WN10: one attempt
+		return approverRefuse(stderr, prefix, "the code you typed is not the code of %s; %s", regPath, nothing)
+	}
+	pre, cancelPre = context.WithTimeout(context.Background(), approveReadBudget)
+	refusal = passkeyEnrolPrecheck(pre, credPath, name, &a) // WN11
+	cancelPre()
+	if refusal != "" {
+		return approverRefuse(stderr, prefix, "%s; %s", refusal, nothing)
+	}
+	if code := confirmOnTerminal(term, env, stderr, prefix, approval.RenderPasskeyEnrolment(a)+fmt.Sprintf("Enrol approver %s? [y/N] ", name), nothing); code != exitOK {
+		return code // N10a, N10
+	}
+
+	// The write, under the credentials writer lock taken only now (N11-N13).
+	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
+	defer cancel()
+	st, err := openCredentialsForSession(ctx, db, credPath)
+	if err != nil { // N11
+		fmt.Fprintf(stderr, "%scannot open the credentials store for writing: %v\n", prefix, err)
+		return exitFatal
+	}
+	defer func() { _ = st.Close() }()
+	p := a.Passkey
+	err = approverEnrolPasskey(st, ctx, store.ApproverRow{Name: a.Name, Format: a.Format, PublicKey: p.Key.Canonical(), EnrolledAt: a.EnrolledAt},
+		store.PasskeyRow{Name: a.Name, RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, AAGUID: p.AAGUID[:], Alg: p.Key.Alg,
+			BackupEligible: p.BackupEligible, SignCount: int64(p.SignCount)})
+	switch {
+	case errors.Is(err, store.ErrApproverNameTaken), errors.Is(err, store.ErrApproverKeyEnrolled), errors.Is(err, store.ErrPasskeyCredentialEnrolled): // N12
+		return approverRefuse(stderr, prefix, "%v (another enrol ran first); %s", err, nothing)
+	case err != nil: // N13
+		fmt.Fprintf(stderr, "%scannot enrol approver %q: %v\n", prefix, a.Name, err)
+		return exitFatal
+	}
+	fmt.Fprintf(stdout, "enrolled approver %s: %s; approvals by it record approved_by=webauthn:%s\n", a.Name, passkeyLabel(p.Key.Alg, p.RPID, p.CredentialID), a.Name) // N14
+	return exitOK
+}
+
+// passkeyEnrolPrecheck is enrol's read-only refusal on its own handle: "" when the store is at
+// schema 6 (N6) and no approver holds the name (N7) and, given the verified registration a,
+// neither its key nor its credential ID (WN11). Revoked rows count: none is ever reused.
+func passkeyEnrolPrecheck(ctx context.Context, credPath, name string, a *approval.Approver) string {
+	ro, err := store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		return fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
+	}
+	defer func() { _ = ro.Close() }()
+	if _, _, err := ro.Passkey(ctx, name); err != nil { // N6: below schema 6 names the serve restart
+		return err.Error()
+	}
+	_, taken, err := ro.Approver(ctx, name)
+	if err != nil {
+		return err.Error()
+	}
+	if taken { // N7
+		return fmt.Sprintf("an approver named %q already exists (a name is never reused, even after revoke); fix: choose another --name", name)
+	}
+	if a == nil {
+		return ""
+	}
+	holder, enrolled, err := ro.ApproverByKey(ctx, a.Passkey.Key.Canonical())
+	if err != nil {
+		return err.Error()
+	}
+	if enrolled { // WN11: the key
+		return fmt.Sprintf("this passkey's key is already enrolled as approver %q (a key is never enrolled twice, even after revoke)", holder.Name)
+	}
+	holder, _, enrolled, err = ro.PasskeyByCredentialID(ctx, a.Passkey.CredentialID)
+	if err != nil {
+		return err.Error()
+	}
+	if enrolled { // WN11: the credential
+		return fmt.Sprintf("credential cred %s is already enrolled as approver %q (a credential is never enrolled twice, even after revoke)", passkeyCred(a.Passkey.CredentialID), holder.Name)
+	}
+	return ""
+}
--- a/cmd/ailang-worldd/approver_passkey_test.go
+++ b/cmd/ailang-worldd/approver_passkey_test.go
@@ -0,0 +1,240 @@
+package main
+
+// Row 201 A2 sprint 3 M6c (design_docs/planned/w-world-approvals-by-key-a2.md §2.12, CW1, CW3):
+// approver enroll --format webauthn refuses before any terminal. Every store is under t.TempDir();
+// every passkey is approvaltest.Passkey's fixed seed (D-WORLD-109).
+
+import (
+	"bytes"
+	"encoding/base64"
+	"encoding/json"
+	"io"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+// enrolRig is a fresh store and the (absent) path the browser saves the registration to.
+type enrolRig struct{ db, creds, reg string }
+
+func newEnrolRig(t *testing.T) enrolRig {
+	t.Helper()
+	db, creds := approverStore(t)
+	return enrolRig{db: db, creds: creds, reg: filepath.Join(t.TempDir(), "world-registration-mark-mac.json")}
+}
+
+// args is approver enroll --name mark-mac --format webauthn --origin pkOrigin --registration reg --db db.
+func (r enrolRig) args() []string {
+	return []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--origin", pkOrigin, "--registration", r.reg, "--db", r.db}
+}
+
+// passkeyCount is the number of webauthn_credentials rows in creds.
+func passkeyCount(t *testing.T, creds string) int {
+	t.Helper()
+	raw := rawCredentials(t, creds)
+	defer func() { _ = raw.Close() }()
+	var n int
+	if err := raw.QueryRow(`SELECT count(*) FROM sqlite_master WHERE name = 'webauthn_credentials'`).Scan(&n); err != nil || n == 0 {
+		return 0 // a schema-5 store has none
+	}
+	if err := raw.QueryRow(`SELECT count(*) FROM webauthn_credentials`).Scan(&n); err != nil {
+		t.Fatal(err)
+	}
+	return n
+}
+
+// linkFragment decodes the fragment of the first ceremony link in out.
+func linkFragment(t *testing.T, out string) ceremonyLink {
+	t.Helper()
+	_, rest, ok := strings.Cut(out, "/surface/passkey#")
+	if !ok {
+		t.Fatalf("no ceremony link in %q", out)
+	}
+	frag, _, _ := strings.Cut(rest, "\n")
+	b, err := base64.RawURLEncoding.Strict().DecodeString(frag)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var c ceremonyLink
+	if err := json.Unmarshal(b, &c); err != nil {
+		t.Fatal(err)
+	}
+	return c
+}
+
+// challengeOf is a fragment's challenge.
+func challengeOf(t *testing.T, c ceremonyLink) [32]byte {
+	t.Helper()
+	b, err := base64.RawURLEncoding.Strict().DecodeString(c.Challenge)
+	if err != nil || len(b) != 32 {
+		t.Fatalf("challenge %q: %v", c.Challenge, err)
+	}
+	return [32]byte(b)
+}
+
+// codeOf is the grouped enrolment code of a registration file, as the page shows it.
+func codeOf(t *testing.T, file []byte) string {
+	t.Helper()
+	r, err := approval.ParseRegistrationResponse(file)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return approval.GroupEnrolmentCode(approval.EnrolmentCode(r.ClientDataJSON, r.AttestationObject))
+}
+
+// enrolEnv is approveEnv with seqRandom as the challenge source.
+func enrolEnv(term io.ReadWriteCloser, opened *bool) sessionEnv {
+	env := approveEnv(term, opened, atTime(approveT0))
+	env.random = seqRandom()
+	return env
+}
+
+// CW1, CW3: every refusal before the terminal is exitUsage, names its reason, opens no terminal,
+// prints no link and enrols nothing; each arm runs under a 5 s guard. Without a terminal, enrol
+// refuses with the tty fix and prints no link.
+func TestPasskeyEnrollFlagRules(t *testing.T) {
+	ed := fixturePub("ed25519.pub")
+	for _, arm := range []struct {
+		name  string
+		args  func(t *testing.T, r enrolRig) []string
+		setup func(t *testing.T, r enrolRig)
+		want  []string
+		rows  int // approvers seeded
+		pks   int // passkeys seeded
+	}{
+		{name: "WN1-webauthn-pubkey", args: func(t *testing.T, r enrolRig) []string { return append(r.args(), "--pubkey", ed) },
+			want: []string{"--pubkey is for --format ssh; --origin and --registration are for --format webauthn", "usage: ailang-worldd approver"}},
+		{name: "WN1-ssh-origin", args: func(t *testing.T, r enrolRig) []string {
+			return append(enrollArgs(r.db, "mark-mac", ed), "--origin", pkOrigin)
+		}, want: []string{"--origin and --registration are for --format webauthn"}},
+		{name: "WN1-ssh-registration", args: func(t *testing.T, r enrolRig) []string {
+			return append(enrollArgs(r.db, "mark-mac", ed), "--registration", r.reg)
+		}, want: []string{"--origin and --registration are for --format webauthn"}},
+		{name: "WN2-name", args: func(t *testing.T, r enrolRig) []string { return dropFlag(r.args(), "--name") },
+			want: []string{"--name, --origin and --registration are required with --format webauthn"}},
+		{name: "WN2-origin", args: func(t *testing.T, r enrolRig) []string { return dropFlag(r.args(), "--origin") },
+			want: []string{"--name, --origin and --registration are required with --format webauthn"}},
+		{name: "WN2-registration", args: func(t *testing.T, r enrolRig) []string { return dropFlag(r.args(), "--registration") },
+			want: []string{"--name, --origin and --registration are required with --format webauthn"}},
+		{name: "N3-no-world-store", args: func(t *testing.T, r enrolRig) []string {
+			return replaceFlag(r.args(), "--db", filepath.Join(t.TempDir(), "absent.db"))
+		}, want: []string{"does not exist; refusing to create one"}},
+		{name: "WN3-ip", args: func(t *testing.T, r enrolRig) []string {
+			return replaceFlag(r.args(), "--origin", "http://127.0.0.1:7644")
+		},
+			want: []string{`--origin "http://127.0.0.1:7644": `, "use http://localhost:<port>"}},
+		{name: "WN3-other-host", args: func(t *testing.T, r enrolRig) []string {
+			return replaceFlag(r.args(), "--origin", "https://evil.example")
+		}, want: []string{`--origin "https://evil.example": `, "localhost or a tailnet *.ts.net name"}},
+		{name: "WN4-exists", setup: func(t *testing.T, r enrolRig) {
+			if err := os.WriteFile(r.reg, []byte("{}"), 0o600); err != nil {
+				t.Fatal(err)
+			}
+		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{"world-registration-mark-mac.json already exists or cannot be checked; pass a new path"}},
+		{name: "WN4-dangling-symlink", setup: func(t *testing.T, r enrolRig) {
+			if err := os.Symlink(filepath.Join(t.TempDir(), "absent"), r.reg); err != nil {
+				t.Fatal(err)
+			}
+		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{"already exists or cannot be checked"}},
+		{name: "WN4-unsearchable-dir", args: func(t *testing.T, r enrolRig) []string {
+			dir := filepath.Join(t.TempDir(), "locked")
+			if err := os.Mkdir(dir, 0o000); err != nil {
+				t.Fatal(err)
+			}
+			t.Cleanup(func() { _ = os.Chmod(dir, 0o700) })
+			return replaceFlag(r.args(), "--registration", filepath.Join(dir, "r.json"))
+		}, want: []string{"r.json already exists or cannot be checked"}},
+		{name: "N6-v5", setup: func(t *testing.T, r enrolRig) {
+			raw := rawCredentials(t, r.creds)
+			defer func() { _ = raw.Close() }()
+			for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+				if _, err := raw.Exec(stmt); err != nil {
+					t.Fatal(err)
+				}
+			}
+		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{"passkeys need schema 6", "restart `ailang-worldd serve`"}},
+		{name: "N7-name-taken", rows: 1, setup: func(t *testing.T, r enrolRig) {
+			seedApprover(t, r.creds, "mark-mac", fixtureKeyOf(t, "ed25519.pub"), approveT0-10)
+		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{`an approver named "mark-mac" already exists`}},
+		{name: "N7-name-taken-revoked-passkey", rows: 1, pks: 1, setup: func(t *testing.T, r enrolRig) {
+			seedPasskey(t, r.creds, "mark-mac", approvaltest.NewPasskey(approval.AlgEdDSA, 31), approveT0-10)
+			revokeSeeded(t, r.creds, "mark-mac", approveT0-5)
+		}, args: func(t *testing.T, r enrolRig) []string { return r.args() }, want: []string{`an approver named "mark-mac" already exists`}},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			r := newEnrolRig(t)
+			if arm.setup != nil {
+				arm.setup(t, r)
+			}
+			args := arm.args(t, r)
+			var out, errw bytes.Buffer
+			var opened bool
+			done := make(chan int, 1)
+			go func() { done <- runApprover(args, &out, &errw, enrolEnv(nil, &opened)) }()
+			var code int
+			select {
+			case code = <-done:
+			case <-time.After(5 * time.Second):
+				t.Fatal("enrol did not refuse within 5 s")
+			}
+			if code != exitUsage || opened || strings.Contains(out.String(), "/surface/passkey") {
+				t.Fatalf("code %d, terminal opened %v, stdout %q; want exitUsage, none, no link\nstderr %q", code, opened, out.String(), errw.String())
+			}
+			for _, w := range arm.want {
+				if !strings.Contains(errw.String(), w) {
+					t.Fatalf("stderr %q lacks %q", errw.String(), w)
+				}
+			}
+			if n, m := len(approverRows(t, r.creds)), passkeyCount(t, r.creds); n != arm.rows || m != arm.pks {
+				t.Fatalf("%d approvers and %d passkeys, want %d and %d (the seeds)", n, m, arm.rows, arm.pks)
+			}
+		})
+	}
+
+	t.Run("CW3-no-terminal", func(t *testing.T) {
+		r := newEnrolRig(t)
+		env := noTerminalEnv(approveT0)
+		env.after, env.random = time.After, seqRandom()
+		var out, errw bytes.Buffer
+		if code := runApprover(r.args(), &out, &errw, env); code != exitUsage || out.Len() != 0 {
+			t.Fatalf("code %d, stdout %q; want exitUsage and no link", code, out.String())
+		}
+		for _, w := range []string{"refusing: no controlling terminal (", "nothing was enrolled", sessionTTYFix} {
+			if !strings.Contains(errw.String(), w) {
+				t.Fatalf("stderr %q lacks %q", errw.String(), w)
+			}
+		}
+		if n, m := len(approverRows(t, r.creds)), passkeyCount(t, r.creds); n != 0 || m != 0 {
+			t.Fatalf("%d approvers, %d passkeys; want 0", n, m)
+		}
+	})
+}
+
+// dropFlag removes flag and its value from args.
+func dropFlag(args []string, flag string) []string {
+	var out []string
+	for i := 0; i < len(args); i++ {
+		if args[i] == flag {
+			i++
+			continue
+		}
+		out = append(out, args[i])
+	}
+	return out
+}
+
+// replaceFlag sets flag's value in args.
+func replaceFlag(args []string, flag, value string) []string {
+	out := append([]string(nil), args...)
+	for i := range out {
+		if out[i] == flag {
+			out[i+1] = value
+		}
+	}
+	return out
+}
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -41,7 +41,7 @@
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
 		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
-		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve --list [--decided]", "approver enroll --name", "approver list", "approver revoke <name>",
+		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve --list [--decided]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
 			t.Errorf("top-level usage lacks %q", want)
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -127,6 +127,7 @@
   ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
+  ailang-worldd approver enroll --name <name> --format webauthn --origin <origin> --registration <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approver list [--db <path>] [--credentials-db <path>]
   ailang-worldd approver revoke <name> [--db <path>] [--credentials-db <path>]
   ailang-worldd [--addr <url>] episode report <episode> [--json | --markdown] [--follow] [--interval 2s] [--addr <url>]
--- a/host/store/context_roots_test.go
+++ b/host/store/context_roots_test.go
@@ -23,6 +23,9 @@
 	"cmd/ailang-worldd/approver.go|runApproverEnroll|Background":      2, // row 201 A1b M3: the read-only pre-check (3 s), then the write after y (3 s)
 	"cmd/ailang-worldd/approver.go|runApproverRevoke|Background":      3, // row 201 A1b M3: the read-only lookup (3 s), the write after y (3 s); sprint 4 K8: the read-only listing (3 s)
 	"cmd/ailang-worldd/episode_launch.go|launchEpisode|Background":    3, // row 217 M2: HTTP preflight (8 s), agent (wall cap + margin), append (3 s)
+	// Row 201 A2 sprint 3: the read-only pre-check (3 s), the re-check after the enrolment code (3 s)
+	// and the write after y (3 s).
+	"cmd/ailang-worldd/approver_passkey.go|runApproverEnrollPasskey|Background": 3,
 
 	"cmd/ailang-worldd/cli.go|get|Background":                                1,
 	"cmd/ailang-worldd/cli.go|execute|Background":                            1,
````

### Patch M6d — the ceremony and the substitute (CW2, CW5)

````diff M6d
--- a/cmd/ailang-worldd/approver_passkey_ceremony_test.go
+++ b/cmd/ailang-worldd/approver_passkey_ceremony_test.go
@@ -0,0 +1,244 @@
+package main
+
+// Row 201 A2 sprint 3 M6d (design_docs/planned/w-world-approvals-by-key-a2.md §2.12, §2.8 [r1];
+// CW2, CW5): the passkey enrolment ceremony at the terminal, the code then y, and the substitute it
+// refuses. Every store is under t.TempDir(); every passkey is approvaltest.Passkey's fixed seed
+// (D-WORLD-109).
+
+import (
+	"bytes"
+	"context"
+	"encoding/base64"
+	"os"
+	"reflect"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// ceremonyTerm is the operator's terminal in a passkey enrolment: its n-th Read returns answer(n),
+// one line per read, as a terminal in canonical mode does; writes and Close are recorded.
+type ceremonyTerm struct {
+	lockedTerm
+	reads  int
+	answer func(n int) (string, error)
+}
+
+func (c *ceremonyTerm) Read(p []byte) (int, error) {
+	s, err := c.answer(c.reads)
+	c.reads++
+	return copy(p, s), err
+}
+
+// browserFor is what the test's browser saves for a fragment: pk's registration for its challenge
+// at its origin and rp.
+func browserFor(pk *approvaltest.Passkey, e ...approvaltest.Edit) func(t *testing.T, f ceremonyLink) []byte {
+	return func(t *testing.T, f ceremonyLink) []byte { return pk.Register(f.Origin, f.RP, challengeOf(t, f), e...) }
+}
+
+// enrolRun is one approver enroll --format webauthn run.
+type enrolRun struct {
+	code           int
+	stdout, stderr string
+	term           *ceremonyTerm
+	file           []byte
+}
+
+// enrol runs approver enroll --format webauthn on r. At the code prompt the test's browser saves
+// save(fragment) to the registration path (nil: it made the path itself) and the terminal answers
+// typed(file); at the y/N it answers yes, after running atYes when not nil. The credentials writer
+// lock must be free at both prompts.
+func (r enrolRig) enrol(t *testing.T, atYes func(), save func(t *testing.T, f ceremonyLink) []byte, typed func(file []byte) string, yes string) enrolRun {
+	t.Helper()
+	var out, errw bytes.Buffer
+	run := enrolRun{term: &ceremonyTerm{}}
+	free := func(prompt string) {
+		if w, err := store.OpenCredentials(r.creds); err != nil {
+			t.Errorf("the credentials writer lock is held at the %s prompt: %v", prompt, err)
+		} else {
+			_ = w.Close()
+		}
+	}
+	run.term.answer = func(n int) (string, error) {
+		switch n {
+		case 0:
+			free("code")
+			if run.file = save(t, linkFragment(t, out.String())); run.file != nil {
+				if err := os.WriteFile(r.reg, run.file, 0o600); err != nil {
+					t.Error(err)
+				}
+			}
+			return typed(run.file), nil
+		case 1:
+			free("y/N")
+			if atYes != nil {
+				atYes()
+			}
+			return yes, nil
+		}
+		return "", os.ErrClosed
+	}
+	run.code = runApprover(r.args(), &out, &errw, enrolEnv(run.term, nil))
+	run.stdout, run.stderr = out.String(), errw.String()
+	return run
+}
+
+// typeCode answers the code prompt with the file's grouped code.
+func typeCode(t *testing.T) func(file []byte) string {
+	return func(file []byte) string { return codeOf(t, file) + "\n" }
+}
+
+// enrolled is the store's only approver and passkey, or a test failure.
+func enrolled(t *testing.T, creds string) (store.ApproverRow, store.PasskeyRow) {
+	t.Helper()
+	rows := approverRows(t, creds)
+	if len(rows) != 1 || passkeyCount(t, creds) != 1 {
+		t.Fatalf("%d approvers and %d passkeys, want 1 and 1", len(rows), passkeyCount(t, creds))
+	}
+	ro, err := store.OpenCredentialsReadOnly(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	p, ok, err := ro.Passkey(ctx, rows[0].Name)
+	if err != nil || !ok {
+		t.Fatalf("passkey of %q: %v %v", rows[0].Name, ok, err)
+	}
+	return rows[0], p
+}
+
+// assertNothingEnrolled fails unless the run refused with code, wrote want to stderr, closed the
+// terminal and left both tables empty.
+func assertNothingEnrolled(t *testing.T, r enrolRig, run enrolRun, code int, want ...string) {
+	t.Helper()
+	if run.code != code {
+		t.Fatalf("code %d, want %d; stderr %q", run.code, code, run.stderr)
+	}
+	for _, w := range want {
+		if !strings.Contains(run.stderr, w) {
+			t.Fatalf("stderr %q lacks %q", run.stderr, w)
+		}
+	}
+	if run.term.reads > 0 && !run.term.Closed() {
+		t.Fatal("the terminal was not closed")
+	}
+	if n, m := len(approverRows(t, r.creds)), passkeyCount(t, r.creds); n != 0 || m != 0 {
+		t.Fatalf("%d approvers and %d passkeys, want none", n, m)
+	}
+}
+
+// CW2: the right code and y enrol exactly the registration the browser made for the printed
+// challenge: one approvers row (webauthn, the canonical COSE key) and one passkey row; the y/N
+// shows RenderPasskeyEnrolment; N14 names the passkey.
+func TestPasskeyEnrollEnrolsWhatTheBrowserMade(t *testing.T) {
+	for _, alg := range []int64{approval.AlgES256, approval.AlgEdDSA} {
+		pk := approvaltest.NewPasskey(alg, 41)
+		r := newEnrolRig(t)
+		run := r.enrol(t, nil, browserFor(pk), typeCode(t), "y\n")
+		if run.code != exitOK {
+			t.Fatalf("alg %d: code %d; stderr %q", alg, run.code, run.stderr)
+		}
+		f := linkFragment(t, run.stdout)
+		seq := make([]byte, 32)
+		_ = seqRandom()(seq)
+		if want := (ceremonyLink{V: 1, Op: "create", RP: pkRPID, Origin: pkOrigin, Name: "mark-mac", Challenge: base64.RawURLEncoding.EncodeToString(seq)}); !reflect.DeepEqual(f, want) {
+			t.Fatalf("fragment %+v, want %+v (the challenge from env.random)", f, want)
+		}
+		if !strings.Contains(run.stdout, "\n  "+pkOrigin+"/surface/passkey#") {
+			t.Fatalf("stdout %q lacks the link at the origin", run.stdout)
+		}
+		a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: pkOrigin, Challenge: [32]byte(seq), Response: run.file, Now: approveT0})
+		if err != nil {
+			t.Fatal(err)
+		}
+		shown := run.term.Output()
+		for _, w := range []string{"Enrolment code shown by the browser (XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX): ", approval.RenderPasskeyEnrolment(a)} {
+			if !strings.Contains(shown, w) {
+				t.Fatalf("the terminal lacks %q:\n%s", w, shown)
+			}
+		}
+		if !strings.HasSuffix(shown, "Enrol approver mark-mac? [y/N] ") || !run.term.Closed() {
+			t.Fatalf("the terminal does not end with the y/N, or was not closed: %q", shown)
+		}
+		label := passkeyLabel(alg, pkRPID, pk.CredentialID())
+		if want := "enrolled approver mark-mac: " + label + "; approvals by it record approved_by=webauthn:mark-mac\n"; !strings.HasSuffix(run.stdout, want) {
+			t.Fatalf("stdout %q lacks %q", run.stdout, want)
+		}
+		row, p := enrolled(t, r.creds)
+		wantRow := store.ApproverRow{Name: "mark-mac", Format: "webauthn", PublicKey: pk.Key().Canonical(), EnrolledAt: approveT0}
+		wantP := store.PasskeyRow{Name: "mark-mac", RPID: pkRPID, Origin: pkOrigin, CredentialID: pk.CredentialID(),
+			AAGUID: bytes.Repeat([]byte{1, 2, 3, 4, 5, 6, 7, 8}, 2), Alg: alg, SignCount: 1}
+		if !reflect.DeepEqual(row, wantRow) || !reflect.DeepEqual(p, wantP) {
+			t.Fatalf("stored %+v %+v, want %+v %+v", row, p, wantRow, wantP)
+		}
+	}
+	// A backup-eligible passkey is stored as one (W13 compares against it).
+	r := newEnrolRig(t)
+	if run := r.enrol(t, nil, browserFor(approvaltest.NewPasskey(approval.AlgEdDSA, 42), approvaltest.Edit{Set: 0x08}), typeCode(t), "y\n"); run.code != exitOK {
+		t.Fatalf("BE: code %d; stderr %q", run.code, run.stderr)
+	}
+	if _, p := enrolled(t, r.creds); !p.BackupEligible {
+		t.Fatal("a backup-eligible passkey was stored as not")
+	}
+}
+
+// CW5 [r1]: read-then-substitute. The attacker reads the genuine registration G, computes its code,
+// and saves its own registration X for the same printed challenge; the operator types G's code.
+// WN10 refuses, the y/N is never shown, nothing is enrolled. A code differing from G's only in its
+// last character, its first group or its length, or holding a 0, is refused on G's own file; G's
+// code in lowercase with dashes, spaces and a tab, or as one string, enrols G.
+func TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode(t *testing.T) {
+	genuine, attacker := approvaltest.NewPasskey(approval.AlgES256, 44), approvaltest.NewPasskey(approval.AlgES256, 45)
+	var genuineCode string
+	substitute := func(t *testing.T, f ceremonyLink) []byte {
+		genuineCode = codeOf(t, genuine.Register(f.Origin, f.RP, challengeOf(t, f)))
+		return attacker.Register(f.Origin, f.RP, challengeOf(t, f))
+	}
+	r := newEnrolRig(t)
+	run := r.enrol(t, nil, substitute, func([]byte) string { return genuineCode + "\n" }, "y\n")
+	assertNothingEnrolled(t, r, run, exitUsage, "the code you typed is not the code of "+r.reg+"; nothing was enrolled")
+	if genuineCode == codeOf(t, run.file) || strings.Contains(run.term.Output(), "[y/N]") {
+		t.Fatalf("the substitute's code equals the genuine one, or the y/N was shown: %q", run.term.Output())
+	}
+
+	other := func(c byte) byte {
+		if c == 'A' {
+			return 'B'
+		}
+		return 'A'
+	}
+	for _, arm := range []struct {
+		name  string
+		typed func(code string) string // code is the grouped genuine code
+		ok    bool
+	}{
+		{name: "last-character", typed: func(c string) string { return c[:len(c)-1] + string(other(c[len(c)-1])) }},
+		{name: "first-group", typed: func(c string) string { return string(other(c[0])) + c[1:] }},
+		{name: "25-characters", typed: func(c string) string { return c[:len(c)-1] }},
+		{name: "27-characters", typed: func(c string) string { return c + "A" }},
+		{name: "a-zero", typed: func(c string) string { return "0" + c[1:] }},
+		{name: "lowercase-dashes-spaces-tab", ok: true, typed: func(c string) string {
+			return strings.ToLower(strings.Replace(strings.Replace(c, "-", " - ", 1), "-", "\t", 2))
+		}},
+		{name: "one-string", ok: true, typed: func(c string) string { return strings.ReplaceAll(c, "-", "") }},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			r := newEnrolRig(t)
+			run := r.enrol(t, nil, browserFor(genuine), func(file []byte) string { return arm.typed(codeOf(t, file)) + "\n" }, "y\n")
+			if arm.ok {
+				if run.code != exitOK {
+					t.Fatalf("code %d; stderr %q", run.code, run.stderr)
+				}
+				enrolled(t, r.creds)
+				return
+			}
+			assertNothingEnrolled(t, r, run, exitUsage, "the code you typed is not the code of ")
+		})
+	}
+}
````

### Patch M6e — unanswered, declined and refused ceremonies (CW4, CW6)

````diff M6e
--- a/cmd/ailang-worldd/approver_passkey_refusals_test.go
+++ b/cmd/ailang-worldd/approver_passkey_refusals_test.go
@@ -0,0 +1,245 @@
+package main
+
+// Row 201 A2 sprint 3 M6e (design_docs/planned/w-world-approvals-by-key-a2.md §2.12; CW4, CW6): the
+// passkey enrolment's unanswered, declined and refused ceremonies, and its write's races. Every
+// store is under t.TempDir(); every passkey is approvaltest.Passkey's fixed seed (D-WORLD-109).
+
+import (
+	"bytes"
+	"context"
+	"errors"
+	"fmt"
+	"io"
+	"os"
+	"reflect"
+	"strings"
+	"sync"
+	"syscall"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// CW4: an unanswered code (WN6, after exactly 300 s), a code cut short (WN7), and a y/N declined,
+// cut short or unanswered (N10, N10a, after 120 s) each refuse, close the terminal and enrol
+// nothing; the writer lock is free at both prompts (enrol's helper checks it). Each arm runs under
+// a 5 s guard, so an unbounded read fails the arm rather than the run.
+func TestPasskeyEnrollDeclinedOrUnanswered(t *testing.T) {
+	pk := approvaltest.NewPasskey(approval.AlgES256, 43)
+	const aborted = "aborted (not confirmed); nothing was enrolled. fix: run it again and type y at the prompt"
+	for _, arm := range []struct {
+		name  string
+		code  func(t *testing.T, file []byte) (string, error) // the code prompt's answer; nil: the file's code
+		yes   func() (string, error)                          // the y/N's answer
+		block int                                             // the prompt left unanswered (its timer fires): 0 none, 1 the code, 2 the y/N
+		want  string
+	}{
+		{name: "WN6-code-timeout", block: 1, want: "ailang-worldd approver enroll: no code within 5m0s; nothing was enrolled\n"},
+		{name: "WN7-code-eof", code: func(*testing.T, []byte) (string, error) { return "", io.EOF },
+			want: "ailang-worldd approver enroll: aborted (no complete code line); nothing was enrolled. fix: run it again and type the code the page shows\n"},
+		{name: "WN7-code-partial", code: func(t *testing.T, file []byte) (string, error) { return codeOf(t, file), io.EOF },
+			want: "ailang-worldd approver enroll: aborted (no complete code line)"},
+		{name: "WN7-code-read-error", code: func(t *testing.T, file []byte) (string, error) { return codeOf(t, file)[:7], errInjectedRead },
+			want: "ailang-worldd approver enroll: aborted (no complete code line)"},
+		{name: "N10-n", yes: func() (string, error) { return "n\n", nil }, want: aborted},
+		{name: "N10-eof", yes: func() (string, error) { return "", io.EOF }, want: aborted},
+		{name: "N10-y-then-eof", yes: func() (string, error) { return "y", io.EOF }, want: aborted},
+		{name: "N10a-y-timeout", block: 2, want: "ailang-worldd approver enroll: no answer within 2m0s; nothing was enrolled\n"},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			r := newEnrolRig(t)
+			var mu sync.Mutex
+			var asked []time.Duration
+			fire := make(chan time.Time) // closed by the unanswered prompt's read, once it waits
+			env := enrolEnv(nil, nil)
+			env.after = func(d time.Duration) <-chan time.Time {
+				mu.Lock()
+				defer mu.Unlock()
+				asked = append(asked, d)
+				if arm.block != 0 && d == []time.Duration{passkeyCeremonyBudget, approverConfirmBudget}[arm.block-1] {
+					return fire
+				}
+				return make(chan time.Time)
+			}
+			var out, errw bytes.Buffer
+			term := &ceremonyTerm{}
+			term.answer = func(n int) (string, error) {
+				if w, err := store.OpenCredentials(r.creds); err != nil { // the writer lock is free at each prompt
+					t.Errorf("prompt %d: the credentials writer lock is held: %v", n, err)
+				} else {
+					_ = w.Close()
+				}
+				if n+1 == arm.block {
+					close(fire)
+					<-t.Context().Done()
+					return "", io.EOF
+				}
+				if n == 0 {
+					file := browserFor(pk)(t, linkFragment(t, out.String()))
+					if err := os.WriteFile(r.reg, file, 0o600); err != nil {
+						t.Error(err)
+					}
+					if arm.code != nil {
+						return arm.code(t, file)
+					}
+					return codeOf(t, file) + "\n", nil
+				}
+				return arm.yes()
+			}
+			env.openTerminal = func() (io.ReadWriteCloser, error) { return term, nil }
+			done := make(chan int, 1)
+			go func() { done <- runApprover(r.args(), &out, &errw, env) }()
+			var code int
+			select {
+			case code = <-done:
+			case <-time.After(5 * time.Second):
+				t.Fatal("enrol did not return within 5 s (an unbounded read)")
+			}
+			assertNothingEnrolled(t, r, enrolRun{code: code, stderr: errw.String(), term: term}, exitUsage, arm.want)
+			mu.Lock()
+			defer mu.Unlock()
+			want := []time.Duration{passkeyCeremonyBudget}
+			if arm.yes != nil || arm.block == 2 {
+				want = append(want, approverConfirmBudget)
+			}
+			if !reflect.DeepEqual(asked, want) || passkeyCeremonyBudget != 300*time.Second {
+				t.Fatalf("the prompt timers were %v, want %v (the code 5m0s, the y/N 2m0s)", asked, want)
+			}
+		})
+	}
+}
+
+// CW6: a missing, failed or all-zero challenge source is an internal error before any terminal
+// (WN5); a FIFO at the registration path (WN8), a registration the law refuses (WN9), and a key or
+// credential another enrol stored while the code prompt waited (WN11) refuse before the y/N; and
+// the write's races (N12) refuse, its failures (N11, N13) are fatal, each enrolling nothing.
+func TestPasskeyEnrollRaceAndRefusals(t *testing.T) {
+	pk := approvaltest.NewPasskey(approval.AlgEdDSA, 46)
+	t.Run("WN5", func(t *testing.T) {
+		for name, random := range map[string]func([]byte) error{
+			"nil": nil,
+			"failed": func(b []byte) error {
+				copy(b, bytes.Repeat([]byte{0xff}, 32))
+				return errors.New("entropy unavailable")
+			},
+			"zero": func(b []byte) error { clear(b); return nil },
+		} {
+			r := newEnrolRig(t)
+			var opened bool
+			env := enrolEnv(nil, &opened)
+			env.random = random
+			var out, errw bytes.Buffer
+			code := runApprover(r.args(), &out, &errw, env)
+			if code != exitFatal || opened || out.Len() != 0 || !strings.Contains(errw.String(), "internal error: no fresh challenge") {
+				t.Fatalf("%s: code %d, terminal opened %v, stdout %q, stderr %q; want exitFatal before any terminal", name, code, opened, out.String(), errw.String())
+			}
+		}
+	})
+	for _, arm := range []struct {
+		name string
+		save func(t *testing.T, r enrolRig, f ceremonyLink) []byte // the browser's step; nil: it made the path itself
+		want string
+		pre  int // approvers another enrol stored
+	}{
+		{name: "WN8-fifo", want: "is not a regular file; nothing was enrolled", save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
+			if err := syscall.Mkfifo(r.reg, 0o600); err != nil {
+				t.Error(err)
+			}
+			return nil
+		}},
+		{name: "WN8-over-16k", want: "is 16385 bytes, over the 16384-byte bound; nothing was enrolled", save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
+			if err := os.WriteFile(r.reg, bytes.Repeat([]byte(" "), 16<<10+1), 0o600); err != nil {
+				t.Error(err)
+			}
+			return nil
+		}},
+		{name: "WN9-uv-clear", want: "does not enrol: " + approval.ErrUserVerification.Error(), save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
+			return browserFor(pk, approvaltest.Edit{Clear: 0x04})(t, f)
+		}},
+		{name: "WN9-another-challenge", want: "does not enrol: " + approval.ErrChallenge.Error(), save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
+			return pk.Register(f.Origin, f.RP, [32]byte{9})
+		}},
+		{name: "WN11-key", pre: 1, want: `this passkey's key is already enrolled as approver "other"`, save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
+			seedPasskeyRows(t, r.creds, "other", pk.Key().Canonical(), approvaltest.NewPasskey(approval.AlgEdDSA, 47).CredentialID())
+			return browserFor(pk)(t, f)
+		}},
+		{name: "WN11-credential", pre: 1, want: "credential cred " + passkeyCred(pk.CredentialID()) + ` is already enrolled as approver "other"`,
+			save: func(t *testing.T, r enrolRig, f ceremonyLink) []byte {
+				seedPasskeyRows(t, r.creds, "other", approvaltest.NewPasskey(approval.AlgEdDSA, 47).Key().Canonical(), pk.CredentialID())
+				return browserFor(pk)(t, f)
+			}},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			r := newEnrolRig(t)
+			run := r.enrol(t, nil, func(t *testing.T, f ceremonyLink) []byte { return arm.save(t, r, f) }, func(file []byte) string {
+				if file == nil { // the browser saved no registration to show a code for
+					return "AAAAA-AAAAA-AAAAA-AAAAA-AAAAAA\n"
+				}
+				return codeOf(t, file) + "\n"
+			}, "y\n")
+			if run.code != exitUsage || !strings.Contains(run.stderr, arm.want) || strings.Contains(run.term.Output(), "[y/N]") {
+				t.Fatalf("code %d, stderr %q, terminal %q; want exitUsage naming %q before the y/N", run.code, run.stderr, run.term.Output(), arm.want)
+			}
+			if n, m := len(approverRows(t, r.creds)), passkeyCount(t, r.creds); n != arm.pre || m != arm.pre {
+				t.Fatalf("%d approvers and %d passkeys, want %d (another enrol's)", n, m, arm.pre)
+			}
+		})
+	}
+	for _, arm := range []struct {
+		name string
+		err  error // the write's error (the seam)
+		hold bool  // another writer holds the lock at the write
+		code int
+		want string
+	}{
+		{name: "N12-name", err: store.ErrApproverNameTaken, code: exitUsage, want: "(another enrol ran first); nothing was enrolled"},
+		{name: "N12-key", err: store.ErrApproverKeyEnrolled, code: exitUsage, want: "(another enrol ran first); nothing was enrolled"},
+		{name: "N12-credential", err: store.ErrPasskeyCredentialEnrolled, code: exitUsage, want: "(another enrol ran first); nothing was enrolled"},
+		{name: "N13-failure", err: errors.New("disk on fire"), code: exitFatal, want: `cannot enrol approver "mark-mac": disk on fire`},
+		{name: "N11-writer-busy", hold: true, code: exitFatal, want: "cannot open the credentials store for writing"},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			if arm.err != nil {
+				orig := approverEnrolPasskey
+				t.Cleanup(func() { approverEnrolPasskey = orig })
+				err := fmt.Errorf("%w: approver %q holds it", arm.err, "x")
+				approverEnrolPasskey = func(*store.Credentials, context.Context, store.ApproverRow, store.PasskeyRow) error { return err }
+			}
+			r := newEnrolRig(t)
+			var hold *store.Credentials
+			atYes := func() {
+				if arm.hold {
+					var err error
+					if hold, err = store.OpenCredentials(r.creds); err != nil {
+						t.Error(err)
+					}
+				}
+			}
+			run := r.enrol(t, atYes, browserFor(pk), typeCode(t), "y\n")
+			if hold != nil {
+				_ = hold.Close()
+			}
+			assertNothingEnrolled(t, r, run, arm.code, arm.want)
+		})
+	}
+}
+
+// seedPasskeyRows stores approver name with key and credential straight through the store.
+func seedPasskeyRows(t *testing.T, creds, name string, key, credentialID []byte) {
+	t.Helper()
+	st, err := store.OpenCredentials(creds)
+	if err != nil {
+		t.Error(err)
+		return
+	}
+	defer func() { _ = st.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	if err := st.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: key, EnrolledAt: approveT0 - 10},
+		store.PasskeyRow{Name: name, RPID: pkRPID, Origin: pkOrigin, CredentialID: credentialID, AAGUID: make([]byte, 16), Alg: approval.AlgEdDSA}); err != nil {
+		t.Error(err)
+	}
+}
````

### Patch M6f — `approve --webauthn-link`

````diff M6f
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -25,6 +25,7 @@
        ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --sign-with <key> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-id> --webauthn-link [--db <path>] [--credentials-db <path>]
        ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
 
 ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
@@ -61,6 +62,11 @@
 --write-message <file> writes that message (90 bytes) to a NEW file and prints
   the ssh-keygen -Y sign line and the --signature step. It opens no terminal
   and writes no store row: the request stays pending.
+--webauthn-link prints, for each origin of the enrolled passkeys, a link to
+  the daemon's page /surface/passkey that asks the browser's passkey (with
+  Face ID, Touch ID, a PIN or a security key) to sign the request's message
+  hash, and the --webauthn-assertion step. Check the page shows the request,
+  expiry and hash it prints. It opens no terminal and writes no store row.
 --list prints the pending requests, read-only. A request expires 10 minutes
 after it is made; at most 16 are pending per credentials store.
 --list --decided prints the used requests instead, newest first (at most 64):
@@ -131,6 +137,7 @@
 	sigPath := fs.String("signature", "", "approve with an enrolled approver's SSHSIG over the request (no terminal)")
 	msgPath := fs.String("write-message", "", "write the request's approval message to this new file, to sign elsewhere")
 	signWith := fs.String("sign-with", "", "approve by signing the request with ssh-keygen -Y sign -f <key> (no terminal)")
+	passkeyLink := fs.Bool("webauthn-link", false, "print the ceremony page's links for the request (no terminal, no write)")
 	positional, err := parseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approveHelp)
@@ -145,12 +152,12 @@
 		return exitUsage
 	}
 	modes := 0
-	for _, v := range []string{*sigPath, *signWith, *msgPath} {
-		if v != "" {
+	for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *passkeyLink} {
+		if on {
 			modes++
 		}
 	}
-	if modes > 1 { // F1: at most one of --signature, --sign-with and --write-message
+	if modes > 1 { // F1: at most one of --signature, --sign-with, --write-message and --webauthn-link
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
@@ -232,6 +239,12 @@
 		}
 		return exitOK
 	}
+	if *passkeyLink { // WF1: the ceremony links, and nothing else (row 201 A2 §2.13)
+		if refusal := writePasskeyLinks(ctx, credPath, stdout, row, mr); refusal != "" {
+			return refuse("%s", refusal)
+		}
+		return exitOK
+	}
 	var signed *approval.SignedApproval
 	if *sigPath != "" { // the signature, judged before any network call (design §2.7 F9-F14)
 		var refusal string
--- a/cmd/ailang-worldd/approve_passkey.go
+++ b/cmd/ailang-worldd/approve_passkey.go
@@ -0,0 +1,66 @@
+package main
+
+// Row 201 A2 sprint 3 (design_docs/planned/w-world-approvals-by-key-a2.md §2.13): `approve
+// <request-id> --webauthn-link` prints the ceremony page's links for the request. It runs after
+// admission and F4, under runApprove's 3 s read root.
+
+import (
+	"context"
+	"encoding/base64"
+	"encoding/hex"
+	"fmt"
+	"io"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// writePasskeyLinks is approve --webauthn-link (design §2.13): one get link per distinct origin
+// among the active passkeys, each allowing that origin's credentials, then the next step. It opens
+// no terminal and writes no row: the request stays pending. It returns the refusal text, or "".
+func writePasskeyLinks(ctx context.Context, credPath string, stdout io.Writer, row store.MintRequestRow, mr approval.MintRequest) string {
+	ro, err := store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		return fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
+	}
+	active, err := ro.ActivePasskeys(ctx)
+	_ = ro.Close()
+	if err != nil { // below schema 6 names the serve restart
+		return err.Error()
+	}
+	if len(active) == 0 { // WF1
+		return "no passkey approver is enrolled; see ailang-worldd approver list"
+	}
+	h, err := requestHash(row)
+	if err != nil {
+		return fmt.Sprintf("request %s: %v", row.RequestID, err)
+	}
+	hash16 := hex.EncodeToString(h[:8])
+	fmt.Fprintf(stdout, "approve request %s (expires %s, message hash %s) with a passkey: open the link for its origin, check the page shows the same request, expiry and hash, and save the file it offers:\n",
+		row.RequestID, approverTime(row.ExpiresAt), hash16)
+	links := map[string]*ceremonyLink{}
+	var order []string
+	for _, e := range active { // ActivePasskeys is ordered by name
+		l := links[e.Passkey.Origin]
+		if l == nil {
+			l = &ceremonyLink{V: 1, Op: "get", RP: e.Passkey.RPID, Origin: e.Passkey.Origin, Request: row.RequestID, Expires: row.ExpiresAt,
+				Hash16: hash16, Challenge: base64.RawURLEncoding.EncodeToString(h[:])}
+			links[e.Passkey.Origin] = l
+			order = append(order, e.Passkey.Origin)
+		}
+		l.Allow = append(l.Allow, base64.RawURLEncoding.EncodeToString(e.Passkey.CredentialID))
+	}
+	for _, o := range order {
+		fmt.Fprintf(stdout, "  %s\n", links[o].link())
+	}
+	fmt.Fprintf(stdout, "then approve it:\n  ailang-worldd approve %s --webauthn-assertion <the saved file> --db %s --credentials-db %s\n",
+		row.RequestID, shellQuote(mr.WorldDB), shellQuote(mr.CredentialsDB))
+	return ""
+}
+
+// requestHash is the admitted row's message hash, the challenge a passkey signs (design §2.1).
+func requestHash(row store.MintRequestRow) ([32]byte, error) {
+	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
+	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
+	return approval.MessageHash([32]byte(stored), [32]byte(challenge), row.ExpiresAt)
+}
--- a/cmd/ailang-worldd/approve_passkey_test.go
+++ b/cmd/ailang-worldd/approve_passkey_test.go
@@ -0,0 +1,121 @@
+package main
+
+// Row 201 A2 sprint 3 M6f, M6g (design_docs/planned/w-world-approvals-by-key-a2.md §2.13, CW7):
+// approve --webauthn-link and --webauthn-assertion. Every store is under t.TempDir(); every passkey is
+// approvaltest.Passkey's fixed seed (D-WORLD-109).
+
+import (
+	"bytes"
+	"encoding/base64"
+	"encoding/hex"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+// passkeySigner is the enrolled passkey of the approve tests (mark-mac, stored counter 1).
+var passkeySigner = approvaltest.NewPasskey(approval.AlgEdDSA, 51)
+
+// requestHashOf is the message hash a passkey signs to approve the rig's stored row.
+func requestHashOf(t *testing.T, a approveRig) [32]byte {
+	t.Helper()
+	row := a.row(t)
+	d, err := hex.DecodeString(row.Digest)
+	if err != nil {
+		t.Fatal(err)
+	}
+	c, err := hex.DecodeString(row.Challenge)
+	if err != nil {
+		t.Fatal(err)
+	}
+	h, err := approval.MessageHash([32]byte(d), [32]byte(c), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return h
+}
+
+// passkeyRig is a request for ep1 at approveT0 with passkeySigner enrolled as mark-mac at
+// approveT0-100.
+func passkeyRig(t *testing.T) approveRig {
+	t.Helper()
+	a := newApproveRig(t, "ep1", approveT0)
+	seedPasskey(t, a.credPath, "mark-mac", passkeySigner, approveT0-100)
+	return a
+}
+
+// CW7: --webauthn-link prints one get link per origin (the request, its expiry and message hash,
+// and that origin's active credentials, by name), opens no terminal and writes nothing; with no
+// passkey it refuses (WF1), and below schema 6 it names the restart.
+func TestApproveByPasskey(t *testing.T) {
+	t.Run("WF1-none-enrolled", func(t *testing.T) {
+		a := newApproveRig(t, "ep1", approveT0)
+		calls := countDaemonCalls(t)
+		var opened bool
+		var stdout, stderr bytes.Buffer
+		if code := runApprove(a.args("--webauthn-link"), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitUsage || opened || *calls != 0 ||
+			stderr.String() != "ailang-worldd approve: no passkey approver is enrolled; see ailang-worldd approver list\n" {
+			t.Fatalf("rc %d, terminal %v, daemon calls %d, stderr %q", code, opened, *calls, stderr.String())
+		}
+		raw := rawCredentials(t, a.credPath)
+		for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`} {
+			if _, err := raw.Exec(stmt); err != nil {
+				t.Fatal(err)
+			}
+		}
+		_ = raw.Close()
+		stderr.Reset()
+		if code := runApprove(a.args("--webauthn-link"), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitUsage || !strings.Contains(stderr.String(), "passkeys need schema 6") {
+			t.Fatalf("v5: rc %d, stderr %q; want exitUsage naming schema 6", code, stderr.String())
+		}
+	})
+
+	a := passkeyRig(t)
+	other, third, revoked := approvaltest.NewPasskey(approval.AlgES256, 52), approvaltest.NewPasskey(approval.AlgEdDSA, 53), approvaltest.NewPasskey(approval.AlgEdDSA, 54)
+	seedPasskey(t, a.credPath, "b-other", other, approveT0-100)
+	seedPasskey(t, a.credPath, "c-tailnet", third, approveT0-100)
+	seedPasskey(t, a.credPath, "d-revoked", revoked, approveT0-100)
+	revokeSeeded(t, a.credPath, "d-revoked", approveT0-50)
+	raw := rawCredentials(t, a.credPath)
+	if _, err := raw.Exec(`UPDATE webauthn_credentials SET origin = 'https://x.ts.net', rp_id = 'x.ts.net' WHERE name = 'c-tailnet'`); err != nil {
+		t.Fatal(err)
+	}
+	_ = raw.Close()
+	calls := countDaemonCalls(t)
+	var opened bool
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--webauthn-link"), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitOK || opened || *calls != 0 {
+		t.Fatalf("link: rc %d, terminal %v, daemon calls %d; stderr %q", code, opened, *calls, stderr.String())
+	}
+	h := requestHashOf(t, a)
+	row := a.row(t)
+	b64 := base64.RawURLEncoding.EncodeToString
+	get := func(rp, origin string, allow ...[]byte) ceremonyLink {
+		l := ceremonyLink{V: 1, Op: "get", RP: rp, Origin: origin, Request: a.id, Expires: row.ExpiresAt, Hash16: hex.EncodeToString(h[:8]), Challenge: b64(h[:])}
+		for _, id := range allow {
+			l.Allow = append(l.Allow, b64(id))
+		}
+		return l
+	}
+	want := []ceremonyLink{get(pkRPID, pkOrigin, other.CredentialID(), passkeySigner.CredentialID()), get("x.ts.net", "https://x.ts.net", third.CredentialID())}
+	out := stdout.String()
+	if n := strings.Count(out, "/surface/passkey#"); n != 2 {
+		t.Fatalf("%d links, want 2:\n%s", n, out)
+	}
+	_, second, _ := strings.Cut(out, "/surface/passkey#")
+	if got := []ceremonyLink{linkFragment(t, out), linkFragment(t, second)}; !reflect.DeepEqual(got, want) {
+		t.Fatalf("links %+v, want %+v", got, want)
+	}
+	for _, w := range []string{"approve request " + a.id + " (expires " + approverTime(row.ExpiresAt) + ", message hash " + hex.EncodeToString(h[:8]) + ")",
+		"\n  " + pkOrigin + "/surface/passkey#", "\n  https://x.ts.net/surface/passkey#", "ailang-worldd approve " + a.id + " --webauthn-assertion <the saved file> --db "} {
+		if !strings.Contains(out, w) {
+			t.Fatalf("stdout lacks %q:\n%s", w, out)
+		}
+	}
+	if row := a.row(t); row.Consumed || len(a.sessionRows(t)) != 0 {
+		t.Fatal("--webauthn-link consumed the request or minted")
+	}
+}
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -41,7 +41,7 @@
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
 		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
-		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve --list [--decided]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
+		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve <request-id> --webauthn-link", "approve --list [--decided]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
 			t.Errorf("top-level usage lacks %q", want)
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -125,6 +125,7 @@
   ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --sign-with <key> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
+  ailang-worldd approve <request-id> --webauthn-link [--db <path>] [--credentials-db <path>]
   ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format webauthn --origin <origin> --registration <file> [--db <path>] [--credentials-db <path>]
````

### Patch M6g — `approve --webauthn-assertion`

````diff M6g
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -26,11 +26,13 @@
        ailang-worldd approve <request-id> --sign-with <key> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --webauthn-link [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-id> --webauthn-assertion <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
 
 ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
 201), after a y/N typed at /dev/tty, or, with --signature or --sign-with,
-after an enrolled approver's SSH signature over the request verifies, with no
+after an enrolled approver's SSH signature over the request verifies, or, with
+--webauthn-assertion, after an enrolled passkey's assertion verifies, with no
 terminal.
 'ailang-worldd episode start <episode> ... --request' makes the request without
 a terminal; nothing is minted until approve.
@@ -67,6 +69,10 @@
   Face ID, Touch ID, a PIN or a security key) to sign the request's message
   hash, and the --webauthn-assertion step. Check the page shows the request,
   expiry and hash it prints. It opens no terminal and writes no store row.
+--webauthn-assertion <file> replaces step 3's y/N with the page's saved file
+  (at most 16 KiB): the assertion is checked against the passkey enrolled for
+  its credential before any network call, its sign counter must move forward,
+  and the request records approved_by=webauthn:<name>.
 --list prints the pending requests, read-only. A request expires 10 minutes
 after it is made; at most 16 are pending per credentials store.
 --list --decided prints the used requests instead, newest first (at most 64):
@@ -138,6 +144,7 @@
 	msgPath := fs.String("write-message", "", "write the request's approval message to this new file, to sign elsewhere")
 	signWith := fs.String("sign-with", "", "approve by signing the request with ssh-keygen -Y sign -f <key> (no terminal)")
 	passkeyLink := fs.Bool("webauthn-link", false, "print the ceremony page's links for the request (no terminal, no write)")
+	assertion := fs.String("webauthn-assertion", "", "approve with an enrolled passkey's assertion file (no terminal)")
 	positional, err := parseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approveHelp)
@@ -152,12 +159,12 @@
 		return exitUsage
 	}
 	modes := 0
-	for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *passkeyLink} {
+	for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *passkeyLink, *assertion != ""} {
 		if on {
 			modes++
 		}
 	}
-	if modes > 1 { // F1: at most one of --signature, --sign-with, --write-message and --webauthn-link
+	if modes > 1 { // F1: at most one of --signature, --sign-with, --write-message, --webauthn-link and --webauthn-assertion
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
@@ -252,6 +259,12 @@
 			return refuse("%s", refusal)
 		}
 	}
+	if *assertion != "" { // the passkey's assertion, judged before any network call (design §2.13 WF2-WF6)
+		var refusal string
+		if signed, refusal = approveAssertion(ctx, credPath, *assertion, row, now); refusal != "" {
+			return refuse("%s", refusal)
+		}
+	}
 	if *signWith != "" { // the signer's output, judged as --signature's (design §2.7 step 3, F7a, F8)
 		armored, refusal := signForApproval(env, *signWith, row, stderr)
 		if refusal != "" {
--- a/cmd/ailang-worldd/approve_passkey.go
+++ b/cmd/ailang-worldd/approve_passkey.go
@@ -1,11 +1,13 @@
 package main
 
 // Row 201 A2 sprint 3 (design_docs/planned/w-world-approvals-by-key-a2.md §2.13): `approve
-// <request-id> --webauthn-link` prints the ceremony page's links for the request. It runs after
-// admission and F4, under runApprove's 3 s read root.
+// <request-id> --webauthn-link` prints the ceremony page's links for the request, and
+// `--webauthn-assertion <file>` approves it with an enrolled passkey's assertion, with no terminal.
+// Both run after admission and F4, under runApprove's 3 s read root.
 
 import (
 	"context"
+	"crypto/sha256"
 	"encoding/base64"
 	"encoding/hex"
 	"fmt"
@@ -64,3 +66,50 @@
 	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
 	return approval.MessageHash([32]byte(stored), [32]byte(challenge), row.ExpiresAt)
 }
+
+// approveAssertion is approve --webauthn-assertion's check, before any network call (design §2.13
+// WF2-WF6): the file is read once through a descriptor check (WF2) and parsed (WF3), the passkey
+// its credential ID names is looked up on a read-only handle (WF4, WF5), and the assertion law
+// judges it against the stored row and that passkey (WF6). It returns the verified approval, or
+// the refusal text. The caller has already admitted row at now and decoded its challenge (F4).
+func approveAssertion(ctx context.Context, credPath, path string, row store.MintRequestRow, now int64) (*approval.SignedApproval, string) {
+	file, err := readRegularFileCapped(path, approval.MaxPasskeyResponseBytes)
+	if err != nil { // WF2
+		return nil, fmt.Sprintf("--webauthn-assertion: %v", err)
+	}
+	resp, err := approval.ParseAssertionResponse(file)
+	if err != nil { // WF3
+		return nil, fmt.Sprintf("--webauthn-assertion %s is not a passkey assertion: %v", path, err)
+	}
+	ro, err := store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		return nil, fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
+	}
+	a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)
+	_ = ro.Close()
+	if err != nil { // WF4: below schema 6 names the serve restart
+		return nil, err.Error()
+	}
+	if !ok { // WF5
+		return nil, fmt.Sprintf("the assertion's credential cred %s is not an enrolled passkey; see ailang-worldd approver list", passkeyCred(resp.RawID))
+	}
+	key, err := approval.ParseStoredCOSEKey(a.PublicKey)
+	if err != nil { // WF5: only a hand-edited row
+		return nil, fmt.Sprintf("approver %q's stored passkey key is not valid (%v); see ailang-worldd approver list", a.Name, err)
+	}
+	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
+	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
+	signed, err := approval.ApproveByAssertion(approval.AssertionInput{
+		Admit: approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Consumed,
+			StoredDigest: [32]byte(stored), BytesDigest: sha256.Sum256(row.Request)},
+		Challenge: [32]byte(challenge),
+		Response:  file,
+		Approver: approval.Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, LastUsedAt: a.LastUsedAt, RevokedAt: a.RevokedAt,
+			Passkey: &approval.Passkey{RPID: p.RPID, Origin: p.Origin, CredentialID: p.CredentialID, Key: key,
+				BackupEligible: p.BackupEligible, SignCount: uint32(p.SignCount)}},
+	})
+	if err != nil { // WF6
+		return nil, fmt.Sprintf("the passkey assertion does not approve request %s: %v", row.RequestID, err)
+	}
+	return signed, ""
+}
--- a/cmd/ailang-worldd/approve_passkey_test.go
+++ b/cmd/ailang-worldd/approve_passkey_test.go
@@ -1,19 +1,24 @@
 package main
 
-// Row 201 A2 sprint 3 M6f, M6g (design_docs/planned/w-world-approvals-by-key-a2.md §2.13, CW7):
-// approve --webauthn-link and --webauthn-assertion. Every store is under t.TempDir(); every passkey is
+// Row 201 A2 sprint 3 M6f (design_docs/planned/w-world-approvals-by-key-a2.md §2.13, CW7): approve
+// --webauthn-link and --webauthn-assertion. Every store is under t.TempDir(); every passkey is
 // approvaltest.Passkey's fixed seed (D-WORLD-109).
 
 import (
 	"bytes"
+	"context"
 	"encoding/base64"
 	"encoding/hex"
+	"os"
+	"path/filepath"
 	"reflect"
 	"strings"
 	"testing"
+	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
 )
 
 // passkeySigner is the enrolled passkey of the approve tests (mark-mac, stored counter 1).
@@ -38,6 +43,33 @@
 	return h
 }
 
+// writeAssertion writes an assertion file under t.TempDir() and returns its path.
+func writeAssertion(t *testing.T, file []byte) string {
+	t.Helper()
+	p := filepath.Join(t.TempDir(), "world-assertion.json")
+	if err := os.WriteFile(p, file, 0o600); err != nil {
+		t.Fatal(err)
+	}
+	return p
+}
+
+// passkeyOf reads name's passkey row.
+func passkeyOf(t *testing.T, creds, name string) store.PasskeyRow {
+	t.Helper()
+	ro, err := store.OpenCredentialsReadOnly(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	p, ok, err := ro.Passkey(ctx, name)
+	if err != nil || !ok {
+		t.Fatalf("passkey %q: %v %v", name, ok, err)
+	}
+	return p
+}
+
 // passkeyRig is a request for ep1 at approveT0 with passkeySigner enrolled as mark-mac at
 // approveT0-100.
 func passkeyRig(t *testing.T) approveRig {
@@ -49,7 +81,10 @@
 
 // CW7: --webauthn-link prints one get link per origin (the request, its expiry and message hash,
 // and that origin's active credentials, by name), opens no terminal and writes nothing; with no
-// passkey it refuses (WF1), and below schema 6 it names the restart.
+// passkey it refuses (WF1), and below schema 6 it names the restart. --webauthn-assertion with the
+// passkey's assertion over the link's challenge mints with no terminal: approved_by
+// webauthn:mark-mac, the counter advanced, the effects the typed y makes. The same file again is
+// refused as already used, and revoke lists the approval (K8).
 func TestApproveByPasskey(t *testing.T) {
 	t.Run("WF1-none-enrolled", func(t *testing.T) {
 		a := newApproveRig(t, "ep1", approveT0)
@@ -118,4 +153,45 @@
 	if row := a.row(t); row.Consumed || len(a.sessionRows(t)) != 0 {
 		t.Fatal("--webauthn-link consumed the request or minted")
 	}
+
+	file := writeAssertion(t, passkeySigner.Assert(pkOrigin, pkRPID, h, 2))
+	stdout.Reset()
+	if code := runApprove(a.args("--webauthn-assertion", file), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitOK || opened {
+		t.Fatalf("assertion: rc %d, terminal opened %v; want exitOK with no terminal\nstderr %s", code, opened, stderr.String())
+	}
+	row = a.row(t)
+	if !row.Consumed || row.ApprovedBy != "webauthn:mark-mac" {
+		t.Fatalf("request consumed %v by %q, want webauthn:mark-mac", row.Consumed, row.ApprovedBy)
+	}
+	if p := passkeyOf(t, a.credPath, "mark-mac"); p.SignCount != 2 {
+		t.Fatalf("sign count %d, want 2", p.SignCount)
+	}
+	if !strings.Contains(stdout.String(), "approved by webauthn:mark-mac\n") {
+		t.Fatalf("stdout lacks the approver:\n%s", stdout.String())
+	}
+	mr, err := approval.Decode(row.Request)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var probs []string
+	for _, p := range approveEffectProblems(mr, collectApproveFacts(t, a, mr, stdout.String())) {
+		if !strings.HasPrefix(p, "approved_by:") { // the typed check wants tty; this path records webauthn:mark-mac (above)
+			probs = append(probs, p)
+		}
+	}
+	if len(probs) != 0 {
+		t.Fatalf("the effects differ from the request:\n%s", strings.Join(probs, "\n"))
+	}
+
+	stderr.Reset()
+	if code := runApprove(a.args("--webauthn-assertion", file), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+2))); code != exitUsage || !strings.Contains(stderr.String(), "was already used") || opened {
+		t.Fatalf("replay: rc %d, stderr %q; want exitUsage, already used", code, stderr.String())
+	}
+
+	term := &lockedTerm{in: strings.NewReader("y\n")}
+	stdout.Reset()
+	if code := runApprover([]string{"revoke", "mark-mac", "--db", a.r.f.db}, &stdout, &stderr, approveEnv(term, nil, atTime(approveT0+3))); code != exitOK ||
+		!strings.Contains(stdout.String(), "it approved 1 session(s)") || !strings.Contains(stdout.String(), "  "+a.id+"  episode ep1") {
+		t.Fatalf("revoke: rc %d, stdout %q; want the passkey's approval listed (K8)", code, stdout.String())
+	}
 }
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -41,7 +41,7 @@
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
 		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
-		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve <request-id> --webauthn-link", "approve --list [--decided]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
+		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve <request-id> --webauthn-link", "approve <request-id> --webauthn-assertion <file>", "approve --list [--decided]", "approver enroll --name", "approver enroll --name <name> --format webauthn --origin <origin> --registration <file>", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
 			t.Errorf("top-level usage lacks %q", want)
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -126,6 +126,7 @@
   ailang-worldd approve <request-id> --sign-with <key> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --webauthn-link [--db <path>] [--credentials-db <path>]
+  ailang-worldd approve <request-id> --webauthn-assertion <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format webauthn --origin <origin> --registration <file> [--db <path>] [--credentials-db <path>]
````

### Patch M6h — refusals before any effect (CW8) and the five-mode F1

````diff M6h
--- a/cmd/ailang-worldd/approve_passkey_refusals_test.go
+++ b/cmd/ailang-worldd/approve_passkey_refusals_test.go
@@ -0,0 +1,148 @@
+package main
+
+// Row 201 A2 sprint 3 M6g (design_docs/planned/w-world-approvals-by-key-a2.md §2.13, CW8): approve
+// --webauthn-assertion refuses before any terminal and any network call, and F1 takes one of five
+// modes. Every store is under t.TempDir(); every passkey is approvaltest.Passkey's fixed seed
+// (D-WORLD-109).
+
+import (
+	"bytes"
+	"errors"
+	"os"
+	"path/filepath"
+	"strings"
+	"syscall"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+// CW8: each refusal of --webauthn-assertion (WF2-WF6, and admission before the file is read) exits
+// exitUsage with its text, opens no terminal, makes no daemon call, and leaves the request pending,
+// no session and the counter at 1.
+func TestApprovePasskeyRefusalsBeforeAnyEffect(t *testing.T) {
+	good := func(t *testing.T, a approveRig) []byte {
+		return passkeySigner.Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
+	}
+	sql := func(stmts ...string) func(t *testing.T, a approveRig) {
+		return func(t *testing.T, a approveRig) {
+			raw := rawCredentials(t, a.credPath)
+			defer func() { _ = raw.Close() }()
+			for _, s := range stmts {
+				if _, err := raw.Exec(s); err != nil {
+					t.Fatal(err)
+				}
+			}
+		}
+	}
+	for _, arm := range []struct {
+		name  string
+		setup func(t *testing.T, a approveRig)
+		file  func(t *testing.T, a approveRig) []byte // nil: path makes the path
+		path  func(t *testing.T) string
+		now   int64 // 0: approveT0+1
+		want  string
+	}{
+		{name: "WF2-fifo", path: func(t *testing.T) string {
+			p := filepath.Join(t.TempDir(), "assertion.fifo")
+			if err := syscall.Mkfifo(p, 0o600); err != nil {
+				t.Fatal(err)
+			}
+			return p
+		}, want: "assertion.fifo is not a regular file"},
+		{name: "WF2-over-16k", file: func(*testing.T, approveRig) []byte { return bytes.Repeat([]byte(" "), 16<<10+1) },
+			want: "is 16385 bytes, over the 16384-byte bound"},
+		{name: "WF2-missing", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "absent.json") }, want: "--webauthn-assertion: open "},
+		{name: "WF3-not-an-assertion", file: func(*testing.T, approveRig) []byte { return []byte(`{"id":"x"}`) },
+			want: "is not a passkey assertion: " + approval.ErrPasskeyResponse.Error()},
+		{name: "WF4-v5-store", setup: sql(`DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`), file: good,
+			want: "passkeys need schema 6"},
+		{name: "WF5-unknown-credential", file: func(t *testing.T, a approveRig) []byte {
+			return approvaltest.NewPasskey(approval.AlgEdDSA, 55).Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
+		}, want: "the assertion's credential cred " + passkeyCred(approvaltest.NewPasskey(approval.AlgEdDSA, 55).CredentialID()) + " is not an enrolled passkey"},
+		{name: "WF5-stored-key-invalid", setup: sql(`UPDATE approvers SET public_key = x'a0' WHERE name = 'mark-mac'`), file: good,
+			want: `approver "mark-mac"'s stored passkey key is not valid`},
+		{name: "WF6-counter-regression", file: func(t *testing.T, a approveRig) []byte {
+			return passkeySigner.Assert(pkOrigin, pkRPID, requestHashOf(t, a), 1)
+		}, want: approval.ErrSignCountRegression.Error() + ": enrolled 1, this assertion 1"},
+		{name: "WF6-revoked", setup: func(t *testing.T, a approveRig) { revokeSeeded(t, a.credPath, "mark-mac", approveT0-50) }, file: good,
+			want: approval.ErrApproverRevoked.Error()},
+		{name: "WF6-enrolled-after-now", setup: sql(`UPDATE approvers SET enrolled_at = 1800000050 WHERE name = 'mark-mac'`), file: good,
+			want: approval.ErrApproverClockBehind.Error()},
+		{name: "WF6-other-request", file: func(t *testing.T, a approveRig) []byte {
+			h := requestHashOf(t, a)
+			h[31] ^= 1
+			return passkeySigner.Assert(pkOrigin, pkRPID, h, 2)
+		}, want: approval.ErrChallenge.Error()},
+		{name: "WF6-backup-eligible", setup: sql(`UPDATE webauthn_credentials SET backup_eligible = 1 WHERE name = 'mark-mac'`), file: good,
+			want: approval.ErrBackupEligibility.Error()},
+		{name: "WF6-another-origin", setup: sql(`UPDATE webauthn_credentials SET origin = 'http://localhost:7645' WHERE name = 'mark-mac'`), file: good,
+			want: approval.ErrOriginMismatch.Error()},
+		{name: "expired-first", path: func(t *testing.T) string { return filepath.Join(t.TempDir(), "absent.json") }, now: approveT0 + 600, want: "expired at"},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			a := passkeyRig(t)
+			if arm.setup != nil {
+				arm.setup(t, a)
+			}
+			var path string
+			if arm.path != nil {
+				path = arm.path(t)
+			} else {
+				path = writeAssertion(t, arm.file(t, a))
+			}
+			calls := countDaemonCalls(t)
+			now := arm.now
+			if now == 0 {
+				now = approveT0 + 1
+			}
+			var opened bool
+			var stdout, stderr bytes.Buffer
+			code := runApprove(a.args("--webauthn-assertion", path), &stdout, &stderr, approveEnv(nil, &opened, atTime(now)))
+			if code != exitUsage || opened || *calls != 0 {
+				t.Fatalf("rc %d, terminal opened %v, daemon calls %d; want exitUsage, none, 0\nstderr %s", code, opened, *calls, stderr.String())
+			}
+			if !strings.Contains(stderr.String(), arm.want) {
+				t.Fatalf("stderr %q lacks %q", stderr.String(), arm.want)
+			}
+			if row := a.row(t); row.Consumed || len(a.sessionRows(t)) != 0 {
+				t.Fatal("a refused assertion consumed the request or minted")
+			}
+			if arm.name != "WF4-v5-store" {
+				if p := passkeyOf(t, a.credPath, "mark-mac"); p.SignCount != 1 {
+					t.Fatalf("sign count %d, want 1", p.SignCount)
+				}
+			}
+		})
+	}
+}
+
+// CW8 (F1, F2): any two of the five modes, and --list with a passkey mode, print the help and exit
+// exitUsage before any store is read.
+func TestApprovePasskeyFlagRules(t *testing.T) {
+	db := dbPathOf(t)
+	modes := [][]string{{"--signature", "s"}, {"--sign-with", "k"}, {"--write-message", "m"}, {"--webauthn-link"}, {"--webauthn-assertion", "f"}}
+	var cases [][]string
+	for i := range modes {
+		for j := i + 1; j < len(modes); j++ {
+			if i < 3 && j < 3 {
+				continue // TestApproveFlagRules has the three SSH pairs
+			}
+			cases = append(cases, append(append(append([]string{firstRequestID}, modes[i]...), modes[j]...), "--db", db))
+		}
+	}
+	cases = append(cases, []string{"--list", "--webauthn-link", "--db", db}, []string{"--list", "--webauthn-assertion", "f", "--db", db})
+	if len(cases) != 9 {
+		t.Fatalf("%d cases, want 7 pairs and 2 --list arms", len(cases))
+	}
+	for _, args := range cases {
+		var stdout, stderr bytes.Buffer
+		if code := runApprove(args, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0))); code != exitUsage || stderr.String() != approveHelp {
+			t.Fatalf("%q: rc %d, stderr %q; want exitUsage and exactly the help", args, code, stderr.String())
+		}
+		if _, err := os.Lstat(filepath.Join(filepath.Dir(db), "credentials.db")); !errors.Is(err, os.ErrNotExist) {
+			t.Fatalf("%q: a credentials store was touched (err %v)", args, err)
+		}
+	}
+}
````
