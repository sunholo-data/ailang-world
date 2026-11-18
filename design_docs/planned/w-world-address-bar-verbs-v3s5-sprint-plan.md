# Sprint plan — w-world-address-bar-verbs-v3, sprint S5 "the CLI and the bar" (iteration 404, row 243)

> **THE SPRINT — row 243 phase V3, sprint S5 (design §0, §1, §2.1, §2.2, §2.3 (`RecheckVerb`), §2.4 (the store calls the
> CLI makes), §2.6 (what the POST writes), §2.7, §2.8, §2.9 P1 and P6, §3, §4, §5, §6 "Sprint 5", §7, §8, §11, §12, §14,
> §16), D-WORLD-118 = A:**
> **M12** — `ailang-worldd episode finish … --confirm-on-surface` (`cmd/ailang-worldd/episode_finish_surface.go`, NEW):
> the flag and its help; refusal 1 (with `--dry-run`, right after F1) and refusal 2 (an unsandboxed pinned check, before
> the check loop: no outcome, no commit); the prompt hook (`else if *onSurface`, the reusable-approval branch unchanged);
> `finishConfirmOnSurface` — the request from `broker.FinishPromptText` of the typed prompt's own fields, `canonicalFinishArgv`
> (pure), id `vr-` + 8 and a 32-byte challenge from `env.random`, `expires = env.now() + 600`, the put under its own 3 s
> root, the prompt to the tty and the card's links to stdout (as `writeAwaitLinks`), `awaitVerbConfirmation` (a poll per
> `awaitPollInterval` through `env.after`, until `expires_at + verbLandGrace` (10 s); the tty is never read), the re-check
> (a) in the poll's own read-only open (the row's digest, challenge and expiry are this process's own; `authority.RecheckVerb`
> over **its own** values; `ApprovedBy()` = `row.approved_by`) and the read-back (b) through the existing
> `worldReader{ctx, c}.object(ref)` — then the unchanged flow from `finishStep("approve")`. `TestProductionContextRoots`
> gains the put and poll roots (W23). **M13** — the bar's finish reason (`episode start` unchanged), the golden regenerated
> by `WORLD_UPDATE_VERB_GOLDEN=1`, `approve vr-<16 hex>` → `navigate /surface/confirm/<id>`, `app.js`'s allowlist gains
> exactly `/surface/confirm/vr-`, the CLI's `approve vr-…` → `exitUsage` with the design's text. **M14** — QUICKSTART §20 +
> CQ1. Plus the carried items: CX21's CLI half, judge-403 note 1 (C14's `now` pinned), S1 Appendix C notes 4 and 5.
> **No schema change, no daemon route, no `.ail`, no *Needs you*, `surfaceCSP` unchanged; no pre-existing
> `episode_finish*_test.go` edited.**
> - **Measured: 1,398 changed lines** — `git diff -M --numstat 18c3d132 HEAD -- cmd host docs` (insertions + deletions,
>   renames on; no rename occurs; **the same 1,398 with no path filter**: `app.js` and the golden are under `host/`, so the
>   filter needs no widening), after a commit of each of the five scratch stages: **323 production** (`episode_finish_surface.go`
>   283, `episode_finish.go` 22, `lookup.go` 13, `approve.go` 3, `app.js` 2), **12 docs** (QUICKSTART §20), **1,063 test**
>   (incl. the regenerated golden, 3). **Cap 1,400: margin 2.**
> - **No cut (S5-D1).** M12 was compiled first and measured **1,277** (design ≈ 730 raw: 1.75×); M13 then measured **59**
>   (≈ 90 raw) and M14 **62** (≈ 90 raw): **1,398 ≤ 1,400** → all three ship; the design's cut (M13 → S6, CQ1 stays) was
>   measured and **not taken**. Design ≈ 910 raw → 1,398 (1.54×).
> - **Five sub-milestones, each green at its commit** (§1.2 C14): **M12a** the unwired production file + its first unit tests
>   + the two context roots (540) → **M12b** the wait, the re-check, the read-back + judge-403 note 1's pin (815) →
>   **M12c** the flag, its two refusals, the prompt hook + the rig with a real daemon (1,277; **the flag is user-visible only
>   from here**) → **M13** the bar and `approve vr-` (1,336) → **M14** QUICKSTART §20 + CQ1 (1,398).
> - **Every line ships as a compiled patch** (Appendix P), extracted from this file and round-tripped onto a fresh
>   `git archive 18c3d132` export, byte-identical to the stage tree (§1.2 C16). The executor applies, never retypes.
> - **Pre-existing test files changed: 5; deleted: 0** (§5, §9 criterion 12): `host/store/context_roots_test.go` (two roots,
>   design §7), `host/verbs/lookup_test.go` (the reason, `vr-`, design §7), `host/surfacepage/testdata/verb-route.golden.json`
>   (regenerated, design §7), `host/daemon/surface_verb_outcomes_test.go` and `host/surfacepage/app_verb_test.go` (outside §7:
>   the golden's generator gains the `approve vr-` row so CB3 replays a real answer, §1.4 F6). **No `episode_finish*_test.go`.**

---

## 0. Binding constraints (each is restated where it governs)

- **D-WORLD-118 = A** (controller, iteration 404: `mission_decisions.sh --open` → no OPEN row; `--check` valid, 105 rows;
  A since iteration 400, attended). S5 is routable; **F6 in `FinishConfirmRefusal`, `DraftOnly` and every "B +≈ 15" branch
  are not built** — CF3 confirms a draft, and nothing in S5 refuses `--ready`.
- **No daemon route, no schema change, no `.ail`**: `host/daemon/daemon.go`, every `host/store` production file,
  `surface_page.go` and `world/`/`packages/` are untouched (the gates' `DAEMON_ROUTES` base = final, `SURFACE_POSTS: 3`,
  `UNTOUCHED_VIOLATIONS: none`). The CLI becomes a credentials-store **writer** (the put, §1.4 F8) through the existing
  `PutVerbRequest`; the daemon gains no writer site (`TestDaemonCredentialWriterSites` unedited, green).
- **The flag is registered only in M12c**, the sub-milestone whose rig proves the whole path (put, wait, re-check,
  read-back, both refusals, every end of the wait); M12a–M12b ship `episode_finish_surface.go` unreachable (no caller).
- **D-WORLD-109**: every store under `t.TempDir()`; keys are `approvaltest`'s (`sigSigner` seed 11 enrolled as `mark-mac`,
  `otherSigner` seed 12 never enrolled; S4's `apSigner` for the daemon pin); the terminal is the rig's `countingTerminal`
  (never `/dev/tty`); `gh` is the rig's `fakeFinishDispatch`; the daemon is `daemon.New` over the rig's stores served by
  `Handler().ServeHTTP` on a loopback request (no listener); no real key, `~/.ssh`, `~/.ailang/world`, mint or publish.
- Coding standards: the request, the argv and the prompt are pure (`finishVerbRequest`, `canonicalFinishArgv`,
  `broker.FinishPromptText`); the CLI's randomness, put, polls and read-back are the host boundary's effects (design §3
  S2); "why not a package" is design §3 S3.
- **Not touched**: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/`, `cmd/world-publish/*`, every
  `host/store`/`host/broker`/`host/approval`/`host/authority`/`host/daemon` production byte, `approve.js`, `decide.js`,
  `confirm.js`, the `/v1` table, every pre-existing `episode_finish*_test.go`.

---

## 1. Pre-flight

### 1.1 Baseline at `18c3d132` (measured by the planner, 2026-11-18)

`tools/gates.sh <base-tree> base-gates base 18c3d132` (the §6 script in `base` mode, on a worktree of the planner's clone at
`18c3d132`; `GATES_SKIP_VERIFY_GO=1`), run beside the planner's compile work (load ≈ 5–8):

| leg | rc | time | note |
|---|---|---|---|
| `TOUCHED_PACKAGES` | — | — | *(empty at the base)* |
| `status` | 0 | 2 s |  |
| `verify_ail` | 0 | 27 s | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the final must match |
| `gofmt`, `build`, `vet` | 0 | 0 / 2 / 1 s |  |
| `fence` | 0 | 20 s | `go test -run '^$' ./...` is the test-file compile fence |
| `pkgs` | 0 | 558 s | `host/verbs` 1.8 s, `host/store` 12 s, `host/surfacepage` 4.8 s, `cmd/ailang-worldd` 196 s, `host/daemon` 341 s |
| `race` | 0 | 94 s | `host/verbs`, `host/store`, `host/surfacepage` whole; the `cmd` and `host/daemon` subsets |
| `node` | 0 | 2 s | the twelve named harnesses `--- PASS`, `WORLD_EXEC_NODE` set |
| `inv` | 0 | 11 s | **53 named `--- PASS`**, 0 SKIP, 0 FAIL |
| `DAEMON_ROUTES` / `SURFACE_POSTS` | — | — | **47** (base 47) / **3** |
| `UNTOUCHED_VIOLATIONS` | — | — | none |
| `chrome` | 0 | 64 s | 24 `--- PASS` (V2b's fourteen, XC1, the two `host/daemon` drills) |
| **`GATES_DONE 782s FAILS: none`** | | | the final tree is compared leg by leg with these (C17) |

### 1.2 Codebase claims checked (at `18c3d132` unless "work" — the planner's stage tree, tag `M14`)

| # | claim | command | observed |
|---|---|---|---|
| C1 | **Input 1 — the design's §2.7/§2.8 cites, re-measured: none of the cited code moved** | `git diff --stat 994ba123 18c3d132 -- cmd/ailang-worldd/{episode_finish,mint_request,approve_await,why_cli,episode_report,cli,approve}.go host/why/why.go host/verbs/lookup.go host/surfacepage/app/app.js host/daemon/episode_finish.go`; control `git diff --stat 994ba123 18c3d132 -- host/daemon/daemon.go cmd/ailang-worldd/episode_finish_operator.go cmd/ailang-worldd/committer.go` | *(empty)*; control `committer.go \| 10`, `episode_finish_operator.go \| 15`, `daemon.go \| 8` (S1's delegations, S3/S4's four registrations) — total `cmd host` since the design: `59 files changed, 5022 insertions(+), 52 deletions(-)` |
| C2 | the cites at HEAD, line by line | `grep -n`/`sed -n` | **unmoved**: `episode_finish.go:52` `var finishClient = newClient`, `:60` the flag set, `:71` `fs.Var(&allowed, "allow-failing", …)`, `:85-89` F1 (`env.openTerminal()` … `defer tty.Close()`), `:96` the 2 h `ctx`, `:102` `c := finishClient(p.Addr)`, `:251` `validateFinishChecks`, `:263` the check loop, `:405-408` the reusable-approval branch, `:410` `finishPrompt`, `:418` `finishStep("approve")`, `:423-430` F6; `mint_request.go:24` `approveReadBudget = 3 * time.Second`, `:102-120` `createMintRequest`'s `env.random` nil refusal and two draws; `approve_await.go:24` `awaitPollInterval = time.Second`, `:29-58` `writeAwaitLinks`; `why_cli.go:181` `type worldReader struct { ctx; c *client; healthCache }`, `:188-190` `why()` via `clientGetter`, `:212` `object(ref)`; `episode_report.go:135` `(worldReader{ctx: ctx, c: c}).health()`; `why.go:66-70` `Object{Hash, SemanticID, Payload []byte}`, `:142-145` `Reader.Object` (`?payload=true`); `cli.go:31` `newClient(base) *client`; `approve.go:228` the `sha256:` arm; `lookup.go:90-92` `case "start", "finish":`, `:124-133` `ApproveTarget` (func at `:125`); `app.js:496` the two-prefix allowlist; `host/daemon/episode_finish.go:186-195` the record phase's `finishDecisionScope` re-check. **Moved**: `daemon.go` `GET /v1/objects/{ref}` `:1113 → :1115` (S3/S4). S1–S4's: `store.OpenCredentials` `credentials.go:210`, `OpenCredentialsReadOnly` `:262`, `PutVerbRequest` `verb_requests.go:121`, `VerbRequestRow` `:72`, `(*Credentials).VerbRequest` `:168`, `authority.RecheckVerb` `judge_verb.go:107`, `broker.FinishPromptText` `finish_prompt.go:24`, `DisplayCommitter` `:14` |
| C3 | **Input 2 — the closest landed analogues, and what was mirrored** | `sed -n 39,262p cmd/ailang-worldd/episode_finish_test.go`; `sed -n 89,175p …_m3_test.go`; `sed -n 95,132p cmd/ailang-worldd/approve_await.go`; `sed -n 20,80p cmd/ailang-worldd/approve_await_test.go`; `sed -n 1,90p cmd/ailang-worldd/approve_surface_sshsig_e2e_test.go` | **Mirrored**: the finish rig (`newFinishFixture`: a git worktree, `f.serve` for report/session/commit/mcp/finish, `fakeFinishDispatch` as gh, `f.run(env, …)`) and m3's re-run pattern (`f.allowRevoked`, a second `f.run` resuming the commit); `createMintRequest`'s draw order (8-byte id, then 32-byte challenge, `env.random == nil` refused) and `approveReadBudget` roots; `awaitSurfaceApproval`'s loop (poll, then the deadline check, then `<-env.after(awaitPollInterval)`) and `pollSurfaceApproval`'s one read-only open per poll under its own root with the judge inside that open; `awaitEnv`'s fake clock (`env.after` advances `now` and calls `onWait`, where a test lands a confirmation mid-wait); the A3 e2e's real daemon (`daemon.New` over a test store, `Handler().ServeHTTP` on a loopback request with `Origin`/`Sec-Fetch-Site`/JSON) to POST a real `approvaltest` SSHSIG. **Differs**: the request is a verb request (`PutVerbRequest`, not `CreateMintRequest`); the wait ends at `expires_at + 10 s` (the mint wait at `expires_at`); the re-check runs over **this process's own** digest/challenge/expiry, never the row's (design §2.7 6(a)); the finish rig's daemon is fake except for `/v1/objects/` and the card's POST (§1.4 F1) |
| C4 | nothing named like S5's code exists | `grep -rn 'confirm-on-surface\|finishConfirmOnSurface\|canonicalFinishArgv\|verbLandGrace\|/surface/confirm/vr-' cmd host --include='*.go' --include='*.js' \| grep -v _test`; control `grep -rn 'awaitPollInterval' cmd --include='*.go' \| grep -v _test \| wc -l` | three **comments** only (`confirm.js:226` "…only at /surface/confirm/vr-<16 hex>", `verb.go:4`, `verb_requests.go:4`) — no code; control `3` |
| C5 | every type, field and function the patches name | `grep -n` at HEAD | **existing**: `sessionEnv{openTerminal, now, random, after}` (`session.go:26-45`), `multiFlag` (`session.go:181`), `verbs.ParseInterleaved`, `presetRefusal` (`preset.go:118`), `approverTime` (`approver.go:343`), `finishSHA`, `finishStep`, `worldReader`/`wireObject = why.Object` (`why_cli.go:152`), `approval.VerbRequest{Schema…Protected}`/`VerbRequestSchema`/`VerbFinish`/`VerbRequestLifetime`/`MaxVerbRequestBytes`/`VerbNamespace`/`CanonicalVerb`/`DecodeVerb`/`VerbMessage`/`VerbMessageHash` (`verb.go`), `authority.VerbJudgeInput`/`RecheckVerb`, `store.VerbRequestRow`/`VerbOutcomeApproved`/`VerbOutcomeSuperseded`/`VerbRequestIDPattern`/`MaxPendingVerbRequests`/`ErrTooManyPendingVerbRequests`/`IsWriterAlreadyActive`/`(*Credentials).ActivePasskeys`/`PendingVerbRequests`, `broker.FinishApprovalScope`/`FormatFinishApprovalScope`/`ParseFinishApprovalScope`/`FinishPromptText`/`DecodeApprovalRequest`/`RecordOperatorApproval`/`ApprovalDecisionV1`/`ApprovalRequestV1`/`EffectHumanApprove`/`ParseExecProfile`/`ExecProfileVersion`, `authority.MintBound`/`ProfileBinding`/`Revoke`; tests: `newFinishFixture`, `memoryDaemon`, `captureFinishTerminal`, `testMintEnv`, `noTerminalEnv`, `awaitEnv`, `seedApprover`, `revokeSeeded`, `seedPasskey`, `pkOrigin`, `sigSigner`, `otherSigner`, `boundedTestContext`; daemon tests: `newConfirmRig`, `confirmSSH`, `apSigner`, `apT0`, `readVerb`. **NEW** (all in S5's files): `verbLandGrace`, `surfaceFinish`, `verbOwn`, `canonicalFinishArgv`, `finishVerbRequest`, `finishConfirmOnSurface`, `surfaceCardLinks`, `putFinishVerbRequest`, `awaitVerbConfirmation`, `pollVerbRequest`, `readBackFinishDecision`, `verbs.verbRequestID`; tests `surfaceUnit`, `storeOwn`, `verbSig`, `credSQL`, `approveByHand`, `pendingVerbRows`, `countingTerminal`, `finishSurfaceRig`, `newSurfaceRig`, `drawnID`, `nextWallSecond`, `protectedProfile` |
| C6 | **Input 3 — fixtures from the real tools** | the test files | requests: the CLI's own `finishVerbRequest` (unit) or the CLI under test (rig), stored by the production `putFinishVerbRequest`; prompts: `broker.FinishPromptText` of the decoded request's own fields, **and** the typed path's captured terminal over the same commit (CF13, `RerunsAndTypedPrompt`); signatures: `sigSigner.Sign(approval.VerbNamespace, approval.VerbMessage(…))` of the stored row (the real signer over the real message); the card's POST: the real `POST /surface/confirm/{id}` of a real `daemon.New`; world objects: the daemon's own C14 writes (rig) or `broker.RecordOperatorApproval` on a test world store (CF7, each refused row the real object with one decoded field re-marshalled); the PR URL: `r.f.finishRecords[0].PR.URL`, the fake gh's own answer; the golden: regenerated by `WORLD_UPDATE_VERB_GOLDEN=1` (C9). No hand-typed bytes or JSON |
| C7 | **Input 4 — tests that enumerate verbs, flags, help, roots, symbols, packages, imports, routes, allowlists, sections** | `grep -rln 'func TestEveryVerbHasHelp\|func TestProductionContextRoots\|…'` (§6 `inv`) | **change**: `TestProductionContextRoots` (+ `episode_finish_surface.go\|putFinishVerbRequest\|Background: 1` and `…\|pollVerbRequest\|Background: 1`; the links read is a child of finish's own ctx, no root); `TestLookup` (the finish reason exact; `episode start`'s exact), `TestApproveTarget` (+ `vr-` rows), `TestSurfaceVerbOutcomes` (+ the `approve vr-` case and golden row), `TestSurfaceAppVerbs` (+ the `vr-` navigate and two refused hrefs). **Unchanged and green** (§6 `inv`): `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch` (no new verb), `TestEpisodeFinishHelp` (it lists required flags and forbids `--yes/--renew/--check`; the new flag is neither), `TestFinishPromptIsTheGolden`, `TestProductionGoSurface` (no store method; counts files, does not pin them), `TestNoNewDeadlineFreeStoreReads`, `TestSurfaceAppPure` (no exported helper), `TestDaemonDependencyAllowlist`/`TestWorldBoundaryDependencyAllowlist` (no new import outside the cmd/host sets already used), `TestFrozenV1RouteTableMatchesMux`, `TestConfirmRoutesRegistered`, `TestDaemonCredentialWriterSites`, `TestQuickstartSection19MatchesTheBar` (§19's block and words unchanged; it reads from §19 to the file's end, which now holds §20 — still green) |
| C8 | **Input 5 — the instruments** | the test files | CF4: `awaitEnv`'s fake clock, each wait run in a goroutine under a 20 s **wall-clock** test deadline (an unbounded wait FAILs, mutant `W-UNBOUNDED`); the exact boundary: a confirmation landing at `now == expires_at + 10` is seen, one at `+ 11` is not; CF8: the busy writer is a held `store.OpenCredentials` (the flock `LOCK_EX\|LOCK_NB`), below v9 is S3/S4's `DROP TABLE verb_requests` + `PRAGMA user_version = 8`, 16 pending are 16 production puts; CF6: `approveByHand` — raw SQL `UPDATE verb_requests SET … outcome = 'approved' …` (design §6); CX21's CLI half: a `BEFORE UPDATE ON verb_requests WHEN NEW.outcome = 'approved'` trigger (S4's CX21 instrument) |
| C9 | the golden's generator, and the exact rows it changes | `grep -n 'WORLD_UPDATE_VERB_GOLDEN' host/daemon/surface_verb_outcomes_test.go`; work: `WORLD_UPDATE_VERB_GOLDEN=1 go test ./host/daemon/ -run '^TestSurfaceVerbOutcomes$'`, then `git diff -U0 -- host/surfacepage/testdata/verb-route.golden.json` | `:17`, `:88` (the repo's own generator: the route's real answers, byte for byte); diff `@@ -4,0 +5 @@` **added** `"approve vr-0123456789abcdef": {…"class":"navigate","outcome":"navigate","cli":"ailang-worldd approve vr-0123456789abcdef","href":"/surface/confirm/vr-0123456789abcdef"}` and `@@ -7 +8 @@` **changed** `"episode finish w1 --preset p"`'s `"reason"` to `"episode finish changes World. The app cannot run it yet; at a terminal, --confirm-on-surface lets you approve its pull request here with your key."`; nothing else (3 numstat lines) |
| C10 | the QUICKSTART test pattern, and §19 to follow | `sed -n 1,55p cmd/ailang-worldd/quickstart_verbs_test.go`; `grep -n '^## 19' docs/QUICKSTART.md` | `TestQuickstartSection19MatchesTheBar` reads `../../docs/QUICKSTART.md`, finds the section by its exact heading after §18, cuts its fenced block and checks each line against the code (`verbs.Split`/`Lookup`) and the page's constants; `1899: ## 19. Operator verbs in the address bar *(attended — pending first verbatim run)*` (the last section). CQ1 follows it for §20 (S5-D3) |
| C11 | the record phase accepts a key-decided approval | `sed -n 180,195p host/daemon/episode_finish.go`; `sed -n 243,278p …` | `finishDecisionScope` requires an `ApprovalDecisionV1` with `decision == "approve"` and an `ApprovalRequestV1` whose scope parses and is unexpired; **`decidedBy` is not read** — a `ssh:`/`webauthn:` decision records (the S5 rig's record phase is the fixture's, so this is the production premise CF3 relies on, unchanged since A4) |
| C12 | `RecheckVerb` has no expiry check (a confirmation seen in the grace window re-checks) | `sed -n 121,146p host/approval/verb_signed.go` | `CheckConfirmation`: `checkApprover(AdmitInput{Now})` (P2–P5: revoked, clock) then the signature over `VerbMessage(in.Digest, in.Challenge, in.Expiry)`; no `now < expiry` test — the grace row of CF4 (`lands at expires_at+grace`) is seen and verifies |
| C13 | the new tests, work | `go test -count=1 -v` by name (`accept.sh`) | all `--- PASS`, none SKIP |
| C14 | **every sub-milestone is green when committed** | `tools/accept_all.sh` (a worktree of `work` at each stage tag running §3's `accept.sh`) | on the final stages (`logs/accept_all.txt`, `logs/accept_all2.txt`; load ≈ 10–19 beside the mutant shards and the final gates): `M12a fast rc=0 ACCEPT_OK` (6 s), `M12b fast rc=0` (6 s), `M12c fast rc=0` (39 s: incl. every pre-existing `TestEpisodeFinish*`, CF11), `M13 fast rc=0` (19 s), `M14 fast rc=0` (18 s), `M14 docs rc=0` (64 s), `M12c cmd rc=0` (308 s, `cmd/ailang-worldd` whole), `M14 cmd rc=0` (232 s), `M13 daemon rc=0` (345 s, `host/daemon` whole), `M12b daemon rc=0` (346 s). Two parts were red on the planner's **first** draft of `accept.sh` and green after it was fixed (no code changed): M13 `fast` named `TestApprove.*` under `v`, which rejects the approve family's Chrome tests that skip without `WORLD_CHROME` (now `p ./cmd/ailang-worldd/ -run '^TestApprove'` plus `v` on CB4); M14 `docs` ran `v` over `TestExecStartupRefusalTable`, which skips its `probe-package-only` subtest by design (now `p` for the `host/daemon` readers) |
| C15 | the snapshot block (§7) on every sub-milestone | `tools/snap_rehearse.sh` (a clone at `18c3d132`, the patches in order, §7's block verbatim) | `M12a 3/0/3`, `M12b 4/0/4`, `M12c 6/0/6`, `M13 14/0/14`, `M14 16/0/16` (`STATUS` lines / `DELETED` lines / MANIFEST lines), each `MANIFEST_OK`, rc 0; dirs `cmd host` for M12a–M13 and `cmd docs host` for M14 (the `find` guard holds both ways) — `logs/snap.txt` |
| C16 | **the round trip** | `tools/roundtrip.sh <this plan>` | `CMP_EQUAL` for all five (each block extracted from this file with §2's awk equals the generated patch byte for byte); each applied with `/usr/bin/patch -p1 -E` after a `--dry-run` to a fresh `git archive 18c3d132` export, no `.orig`/`.rej` under `cmd host docs`, committed and measured: `NUMSTAT M12a 540`, `M12b 815`, `M12c 1277`, `M13 1336`, `M14 1398` (each `unfiltered=` the same); then `diff -r -x .git` against the stage tree (`git archive M14`) → **`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE`** (`logs/roundtrip_last.txt`, re-run on this final plan file) |
| C17 | **the gates on the final tree** | `tools/gates.sh <final-tree> final-gates final <sha>` | `tools/gates.sh <final-tree> final-gates final 740892a05` (`final-tree` = a worktree of the planner's clone at the stage tag `M14`, the five stage commits on `18c3d132`; `diff -r` equal to the round-trip tree, C16; `GATES_SKIP_VERIFY_GO=1`; **load ≈ 13–19**, beside the planner's three mutant shards): **`TOUCHED_PACKAGES: cmd/ailang-worldd host/daemon host/store host/surfacepage host/verbs`** (= the plan's list); `status` 0; **`verify_ail` rc 0** (`✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the base's line, byte for byte); `gofmt`, `build`, `vet`, `fence` rc 0; **`pkgs` rc 1** (947 s: `host/verbs`, `host/store`, `host/surfacepage` ok, `cmd/ailang-worldd` ok 298 s; `host/daemon` red on **two load timings outside S5** — `TestWorkspaceReadThroughProductionWiring` "plan phase exceeded its 4s budget after 6.02s" and `TestSeToolsRunBenchmarkTasksEndToEnd` (the se-tools mock saw a stray GET) — **both `--- PASS` re-run alone** at once (`go test -count=1 -v ./host/daemon/ -run '^(TestWorkspaceReadThroughProductionWiring|TestSeToolsRunBenchmarkTasksEndToEnd)$'` → `ok 12.8s`, `logs/daemon_rerun_alone.txt`); S5 changes no `host/daemon` production byte); **`race` rc 0** (256 s: the three small packages whole; the `cmd` subset 53 `--- PASS`, the `host/daemon` subset 49); **`race3` rc 0** (54 s: the eight CLI wait tests ×3, 45.9 s; the daemon pin ×3); **`node` rc 0** (the twelve named harnesses); **`inv` rc 0: 71 named `--- PASS`** (the base's 53 + S5's 18), 0 SKIP, 0 FAIL, no `NOT PASSED`; **`DAEMON_ROUTES: 47 (base 47)`, `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none`**; **`chrome` rc 0** (89 s: 24 `--- PASS` as at the base — V2b's fourteen, XC1 and the two drills, on the new `app.js`); **`GATES_DONE 1455s FAILS: pkgs`** — the two load reds above, green alone. `verify_go` is the controller's |
| C18 | node, the pinned binaries | `/opt/homebrew/bin/node --version`; `$HOME/.pinned-ailang/ailang --version`; `$HOME/.pinned-ailang-tools/v0.54.15/ailang --version` | `v26.13.0`; `AILANG v0.41.0`; the tools pin is used only by the gates' real-stack legs |
| C19 | **judge-403 note 1 reproduced, then pinned** | work: `J-NOW-C14` (C14's `now` → `now+1` in `host/daemon/surface_confirm.go`) run against S4's tests only (`go test ./host/daemon/ -run '^(TestConfirm\|TestConsumeRefusal)' -skip '^TestConfirmPOSTRecordsAtItsClock$'`), then with S5's pin | **survived** (`ok`) — the judge's finding reproduced; **killed** by `TestConfirmPOSTRecordsAtItsClock` (the decision's `now` = the injected `approveNow` = C15's `consumed_at`) and, through the CLI, by `TestFinishConfirmOnSurfaceOpensThePR` (the read-back decision's `now` = the row's `consumed_at`) — ledger lines `J-NOW-C14(control…)`, `J-NOW-C14`, `J-NOW-C14@cmd` |

### 1.3 Planner decisions where the design or the directive leaves a choice

- **S5-D1 — no cut; M12, M13 and M14 ship, by measurement.** M12 compiled first (with every CF test and the carried
  items): **1,277**. M13 then **59**, M14 **62**: **1,398 ≤ 1,400**. The design's cut ("S5 → CQ1 stays (S7), M13 moves to
  S6") was measured and not taken. The margin is 2 lines: the test files were tightened twice to stay under the cap
  (§1.5) — nothing was dropped that a CF row or a mutant needs.
- **S5-D2 — five sub-milestones; the flag last among M12's.** M12a (540) adds `episode_finish_surface.go` with **no
  caller** (Go allows the unexported functions; nothing is reachable), its pure tests (CF14, the request's fields, CF9,
  the link rule, CF8) and the two context-root pins. M12b (275) adds the wait/re-check/read-back tests (CF4–CF7) and
  judge-403 note 1's daemon pin. M12c (462) registers the flag, the two refusals and the hook **with** the rig that drives the
  whole path through a real daemon (CF1–CF3, CF5, CF6, CF10, CF12, CF13, CX21's CLI half). M12a and M12c exceed the
  directive's ≈ 400 guide (540, 462): M12a's tests exercise only the code it adds, and M12c cannot register the flag
  without its rig — splitting either would ship a test-less file or an unproven flag.
- **S5-D3 — "CQ1 stays (S7)", read by measurement.** `S7` is the charter's QUICKSTART guardrail (the walkthrough is
  "executed-verbatim, S7-maintained", `CLAUDE.md`; `quickstart_layout_test.go:3` cites "S7"), not a sprint: the design means
  CQ1 is kept in S5 even under the cut, because S7 binds the sprint that adds a CLI flag. CQ1 is
  `TestQuickstartSection20MatchesTheFinishFlag`, the `quickstart_verbs_test.go` pattern (C10): §20's fenced command parses
  with `verbs.Split` and ends in `--confirm-on-surface`, the help lists the flag, the three lines finish prints are quoted
  (each format string read from `episode_finish_surface.go` and its placeholder form found in §20), and §20's namespace,
  10 minutes and 10 seconds equal `approval.VerbNamespace`, `VerbRequestLifetime`, `verbLandGrace`. §20's heading carries
  §14–§19's *(attended — pending first verbatim run)*: it needs a key and a browser.
- **S5-D4 — the rig: the finish fixture plus a real daemon for exactly the two routes the path needs** (§1.4 F1).
  `newSurfaceRig` keeps `f.serve` for report/session/commit/mcp/finish (its record phase, its proofs) and routes
  `/v1/objects/` to `daemon.New` over the rig's own world store (`f.db`) and credentials store; the test POSTs the card
  to the same daemon. Its clock starts at the wall clock the daemon reads (each `run` resets it), since the daemon's
  `approveNow` cannot be injected from `cmd/` (§1.4 F3).
- **S5-D5 — the canonical argv.** `fs.Visit` (set flags only, lexical), `--name` for a true bool, `--name value` as two
  words, `--allow-failing` once per value in order; **a bool set false is omitted** (it equals its default; "--name
  false" would parse as a positional). CF14 proves two interleaved spellings (incl. `--title=two words`, `--ready=true`,
  `--dry-run=false`) give one argv and an unset flag never appears; CF3 checks the rig's real argv.
- **S5-D6 — the links are read before the put** (design §2.7 orders the put, step 3, before the print, step 4): a request
  whose links cannot be read is never stored (CF9 row "no credentials store"), so no pending request exists that the CLI
  did not show. The print itself stays after the put.
- **S5-D7 — every refusal is `fail(refusal, "the pull request was not confirmed on the surface; commit kept at <sha>",
  steps(sha))`**: exit 1 (`exitUsage`, design §2.7 "exits 1"), the commit kept, nothing dispatched — including "internal
  error: no randomness source" (the mint path's `exitFatal` is for a request that was never committed; here a commit exists
  and the hand steps are the useful next line).
- **S5-D8 — re-check (a) is two layers**: the row's `digest`/`challenge`/`expires_at` equal this process's own (refusal
  "the stored request is not the one this finish made"), then `RecheckVerb` over **its own** values, then `ApprovedBy()`.
  With the equality guard first, `RecheckVerb` over the row's values is an **equivalent mutant** (`Q-ROW-VALUES`); the pair
  "guard removed and row values" is its own mutant, killed (`Q-ROW-VALUES-NOEQ`).
- **S5-D9 — CF12 (S1 Appendix C note 4): the resume-dispatch reuse skips CF6/CF7, and that is safe.** The reused approval
  is the one the first run already re-checked (a) and read back (b) before its first dispatch; the daemon's record phase
  accepted it only after `finishDecisionScope` re-read the decision and its scope (C11); reuse is the unchanged typed-path
  branch (`finishReusableApproval`: same scope, `ExpiresAt ≥ now + 300`, a record at HEAD with a failed-dispatch outcome),
  and F6's HEAD/remote/title/body re-check still runs before the push. `TestFinishConfirmOnSurfaceResumeReusesTheApproval`
  shows it: no new request, no wait, **zero** `/v1/objects` reads, `resuming under approval <the first decision>`.
- **S5-D10 — S1 Appendix C note 5 / gemini r2 (the ≈ 1 s window after C15)**: fail-safe, tested twice — the unit row
  "the approver revoked since" (`TestPollVerbRequestRechecks`) and the rig (`KeepsTheCommit`: the card's POST returns 200,
  then `revokeSeeded`, then the CLI's poll refuses "does not verify … revoked", commit kept, nothing dispatched, no
  read-back); a re-run asks afresh (CF10).
- **S5-D11 — judge-403 note 1, two pins** (C19): `TestConfirmPOSTRecordsAtItsClock` (host/daemon, new file; S4's files
  unedited) and CF3's read-back assertion `decision.now == row.consumed_at`.
- **S5-D12 — judge-403 notes 2–4 are not taken** (the margin is 2 lines): Appendix C carries them.
- **S5-D13 — CF ids.** The design's §6 names CF1–CF14 as "§2.9 P1's rows"; P1 has no CF11 row: **CF11 = P1's last row**,
  "the typed path (no flag) unchanged" — every pre-existing `TestEpisodeFinish*` test, run unedited at M12c (§3). Mapping:
  §3.1.
- **S5-D14 — timing (directive 8, 12).** The CLI's wait runs on `env.after`'s fake clock everywhere; the unit waits run
  under a 20 s wall-clock deadline. The only real waits are two `nextWallSecond` calls (≤ 1 s each, capped at 3 s): a finish
  approval of the same scope in the same Unix second is the same request object (W59), and the daemon's clock is the wall
  clock (§1.4 F3) — CX21's CLI half needs a second decision, so only `DispatchesTheNamedDecision` waits. Race cost §1.5.
- **S5-D15 — mutants' "package" leg is `suite-nostatic`**: the behavioural tests of the touched package (`cmd`: every new
  test, `TestEpisodeFinish*`, `TestApprove*`, the A3 e2e; 24 s) — the whole `cmd/ailang-worldd` (109 s) per mutant was not
  affordable; no inventory, AST or string-grep test is in it (S4's practice).

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

| # | premise | measured | routing |
|---|---|---|---|
| F1 | design §6 "CF1–CF14 … on the finish rig with a real loopback daemon" | the finish rig has no daemon: `newFinishFixture` sets `finishClient = memoryDaemon(f.serve)`, a hand-written handler (`episode_finish_test.go:110`); no cmd test combines it with `daemon.New` | S5-D4: the rig keeps `f.serve` and adds a real `daemon.New` for `/v1/objects/` and the card's POST, served by `Handler().ServeHTTP` on a loopback request (the A3 e2e's form, `approve_surface_sshsig_e2e_test.go:26`) — no listener (D-WORLD-109) |
| F2 | design §2.7 step 3 "the store F6 proved the daemon serves" + "below v9 → the restart text" | true: `PutVerbRequest` refuses `< 9` with `OlderCredentialsStoreError` naming `restart \`ailang-worldd serve\`` (`mint_requests.go:74`) | CF8 asserts both phrases |
| F3 | design CX21 "the identical body a second later" (and W59) in a CLI rig | the daemon's clock is unexported (`approveNow`, `daemon.go:464`; `daemon.Config` has no clock) — a cmd test cannot step it | `nextWallSecond`: a bounded (≤ 1 s, 3 s cap) wait for the wall second to change, in `DispatchesTheNamedDecision` only (S5-D14) |
| F4 | directive: `daemon.go` `GET /v1/objects/{ref}` at `:1113` | `:1115` (S3/S4's registrations); every other cite unmoved (C1–C2) | cited at HEAD |
| F5 | design §7 "every `episode_finish*_test.go` stays unedited" | holds: `git diff --name-only 18c3d132 M14 -- 'cmd/ailang-worldd/episode_finish*_test.go'` → *(empty)*; S5's tests are three NEW files (`episode_finish_surface{,_rig}_test.go`, and the unit file) | — |
| F6 | design §7's edited-test list for M13 (`lookup_test.go`, the golden) | the golden is regenerated from `TestSurfaceVerbOutcomes`' fixed query list, which has no `vr-` query; CB3 (the page's allowlist) is a node test that replays golden answers (`app_verb_test.go`) | two pre-existing tests outside §7 edited, one line each: `surface_verb_outcomes_test.go` (+ the `approve vr-…` case and query) and `app_verb_test.go` (+ the `vr-` navigate, + `/surface/confirm/x` and `/surface/confirm/mr-…` refused) |
| F7 | design §2.7 step 1 "over `MaxVerbRequestBytes` → refuse" as a separate CLI check | `Validate` already refuses a canonical size over 128 KiB (`verb.go:112`, V10) | the CLI calls `Validate` once; CF9's oversize row ends "; confirm at the terminal" (`V-VALIDATE` killed) |
| F8 | directive: "find any test that enumerates the CLI's credential writer sites" | none exists: `grep -rn 'func Test.*WriterSites' cmd host` → only `host/daemon`'s `TestDaemonCredentialWriterSites` (control: found) | the CLI's new writer (`putFinishVerbRequest`) is pinned by its context-root pin (`TestProductionContextRoots`, a named root) and CF8; a CLI writer-site inventory is an Appendix C candidate, not built (margin) |
| F9 | design §2.7 "the wait ends at `expires_at + verbLandGrace`" vs `RecheckVerb` | C12: the re-check has no expiry test, so a confirmation that lands in the grace window verifies; consistent | the grace row of CF4 confirms it |

**Files outside design §7**: `host/daemon/surface_confirm_now_test.go` (judge-403 note 1), `cmd/ailang-worldd/approve_verb_request_test.go`
(CB4), `cmd/ailang-worldd/quickstart_finish_surface_test.go` (CQ1); §7's "(+ tests)" for `episode_finish_surface.go` are
`episode_finish_surface_test.go` and `episode_finish_surface_rig_test.go`. Pre-existing tests edited outside §7: F6's two.

### 1.5 The size and the race-mode cost

| stage | production | test / docs | changed lines | cumulative (numstat after commit) |
|---|---|---|---|---|
| M12a | 283 (`episode_finish_surface.go`) | 257 (`episode_finish_surface_test.go` 253, context roots 4) | 540 | **540** |
| M12b | 0 | 275 (`episode_finish_surface_test.go` +245, `surface_confirm_now_test.go` 30) | 275 | **815** |
| M12c | 22 (`episode_finish.go`) | 440 (`episode_finish_surface_rig_test.go`) | 462 | **1,277** |
| M13 | 18 (`lookup.go` 13, `approve.go` 3, `app.js` 2) | 41 (`approve_verb_request_test.go` 24, `lookup_test.go` 7, golden 3, outcomes 3, `app_verb_test.go` 4) | 59 | **1,336** |
| M14 | 0 | 50 (`quickstart_finish_surface_test.go`) + 12 docs | 62 | **1,398** |

Design ≈ 910 raw → compiled **1,398 (1.54×)**: M12's tests (CF1–CF14 with every carried item) are 1.75× its raw
estimate, M13/M14 are under theirs. Two tightening passes kept it under the cap (a 46-line daemon pin made 30; a third
rig run, the "afresh" key, folded into CF10; comments shortened); the cost is a 2-line margin, recorded.

**Race-mode wall time of the new tests** (directive 8; measured on the work tree, load ≈ 8, `go test -race -count=1 -v`):
`cmd/ailang-worldd` — `TestFinishConfirmOnSurfaceOpensThePR` 0.43 s, `…RefusesBeforeAnyCheck` 0.23 s, `…RerunsAndTypedPrompt`
0.56 s, `…KeepsTheCommit` 0.47 s, `…DispatchesTheNamedDecision` 1.11 s (≤ 1 s of it is `nextWallSecond`), `…ResumeReusesTheApproval`
0.50 s, `…RefusesBeforeTheStore` 0.10 s, `TestFinishVerbRequestFields` 0.02 s, `TestSurfaceCardLinks` 0.05 s,
`TestPutFinishVerbRequestRefusals` 0.16 s, `TestAwaitVerbConfirmationEnds` 0.35 s, `TestPollVerbRequestRechecks` 0.44 s,
`TestReadBackFinishDecision` 0.05 s, `…RefusesAnUnrecordedDecision` 0.06 s, `TestCanonicalFinishArgv`,
`TestApproveVerbRequestPointsAtItsCard`, `TestQuickstartSection20MatchesTheFinishFlag` 0.00 s — **≈ 4.5 s added to
`cmd/ailang-worldd`'s race run**; `host/daemon` — `TestConfirmPOSTRecordsAtItsClock` **0.11 s** (CI's bound is 8 min per
package; host/daemon takes 446–480 s on slow runners: +0.02 %). `host/verbs`/`host/surfacepage` edits add table rows only.
No sleep but `nextWallSecond` (≤ 1 s, once).

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`; **no
sub-milestone deletes a file**. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M12a   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M12a M12b M12c M13 M14
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter404-row243-v3s5
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v3s5-sprint-plan.md"
T="$HOME/.ailang/state/world-iter404/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) 2>/dev/null | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short --untracked-files=all
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped, reordered or changed.
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan (`awk '$0 == "````bash
  accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/accept.sh"`) and run **each part §3 lists for M, one
  per Bash call**: `bash "$T/accept.sh" "$W" <M> <part> 2>&1 | tee
  "$HOME/.ailang/state/world-iter404/exec-logs/accept_<M>_<part>.txt"` (never only through `tail`). **Never `nohup`,
  never background a test leg; pass a 600000 ms timeout to every long Bash call.**
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every part runs the same `[base]` checks first (gofmt, build, vet of the five touched packages, the
test-file compile fence `go test -run '^$' ./...`). `v` requires every named test to `--- PASS` and none to SKIP or FAIL;
`p` is a whole-package run. Parts: `fast` (named tests; ≈ 0.5–1 min), `cmd` (`cmd/ailang-worldd` whole, ≈ 2 min),
`daemon` (`host/daemon` whole, ≈ 4–8 min at load), `docs` (every test outside `cmd/` that reads QUICKSTART, ≈ 1 min). Each
part is one 600 s call; measured in C14. Commit-time test commands for `commit_ms.sh` are the `fast` parts (§7).

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M> [fast|cmd|daemon|docs]: exactly plan §3's acceptance for sub-milestone M (one bash
# shell; rc 0 = ACCEPT_OK). Each part fits one 600 s call; run every part §3 lists for M, in order.
T=$1; M=$2; P=${3:-fast}
cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang
unset WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }   # -v: every named test PASS, none SKIP/FAIL
p() { go test -count=1 "$@"; }
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M $P: $*"; exit 1; }; }
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
step go vet ./cmd/ailang-worldd/ ./host/daemon/ ./host/store/ ./host/surfacepage/ ./host/verbs/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M $P: FENCE_RED"; exit 1; fi
INV='TestEveryVerbHasHelp|TestVerbTableMatchesTheCLISwitch|TestFinishPromptIsTheGolden|TestEpisodeFinishHelp|TestQuickstartSection19MatchesTheBar'
STORE='TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestStoreSurfaceTakesContext|TestCredentialsMethodsRequireADeadline'
CFA='TestCanonicalFinishArgv|TestFinishConfirmOnSurfaceRefusesBeforeTheStore|TestFinishVerbRequestFields|TestSurfaceCardLinks|TestPutFinishVerbRequestRefusals'
CFB='TestAwaitVerbConfirmationEnds|TestPollVerbRequestRechecks|TestReadBackFinishDecision|TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision'
RIG='TestFinishConfirmOnSurfaceOpensThePR|TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck|TestFinishConfirmOnSurfaceRerunsAndTypedPrompt|TestFinishConfirmOnSurfaceKeepsTheCommit|TestFinishConfirmOnSurfaceDispatchesTheNamedDecision|TestFinishConfirmOnSurfaceResumeReusesTheApproval'
CX='TestConfirmPOSTRecordsAtItsClock|TestConfirmPOSTConfirms|TestConfirmPOSTRetriesAFailedConsume|TestConfirmPOSTRefusedAtC15|TestDaemonCredentialWriterSites|TestConfirmRoutesRegistered'
CB='TestSurfaceVerbOutcomes|TestSurfaceVerbRouteRegistered|TestSurfaceVerbIsReadOnly|TestFrozenV1RouteTableMatchesMux|TestDaemonDependencyAllowlist'
DOCS='TestAGUIQuickstartControl|TestCommitBudgetAndUncertainReconcile|TestExecStartupRefusalTable|TestWorkspaceExecEpisodeSelectsItsProfile|TestMCPQuickstartPayloadsVerbatim|TestSeToolsQuickstartPayloadsVerbatim|TestSeToolsWorkspaceExecWithoutProfileRefuses|TestQuickstartDecisionPacketSection|TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive|TestWorkspacePackageCacheCoversLock'
case "$M:$P" in
M12a:fast) step v ./cmd/ailang-worldd/ -run "^($CFA|$INV)\$"
           step v ./host/store/ -run "^($STORE)\$" ;;
M12b:fast) step v ./cmd/ailang-worldd/ -run "^($CFA|$CFB|$INV)\$"
           step v ./host/daemon/ -run "^($CX)\$" ;;
M12c:fast) step v ./cmd/ailang-worldd/ -run "^($CFA|$CFB|$RIG|$INV|TestEpisodeFinish.*|TestFinishRepository|TestFinishProductionGHArgv)\$"
           step v ./host/store/ -run "^($STORE)\$" ;;
M13:fast)  step p ./host/verbs/
           step v ./host/verbs/ -run '^(TestLookup|TestApproveTarget)$'
           step v ./host/daemon/ -run "^($CB)\$"
           step v ./host/surfacepage/ -run '^(TestSurfaceAppVerbs|TestSurfaceAppPure|TestSurfaceScriptPure)$'
           step v ./cmd/ailang-worldd/ -run "^(TestApproveVerbRequestPointsAtItsCard|$INV)\$"
           step p ./cmd/ailang-worldd/ -run '^TestApprove' ;;   # the approve family (its Chrome tests skip without WORLD_CHROME)
M14:fast)  step v ./cmd/ailang-worldd/ -run "^(TestQuickstartSection20MatchesTheFinishFlag|TestQuickstart.*|$INV)\$" ;;
M14:docs)  step p ./host/daemon/ -run "^($DOCS)\$"   # rc only: TestExecStartupRefusalTable skips a probe subtest by design
           step v ./host/broker/ -run '^(TestQuickstartProtectedContract|TestQuickstartSection10ProfilesLoad)$'
           step v ./cmd/world-publish/ -run '^TestTransitionsWithoutAControllingTerminalRefusesAndWritesNothing$' ;;
M12c:cmd|M14:cmd) step p -timeout 30m ./cmd/ailang-worldd/ ;;
M12b:daemon|M13:daemon) step p -timeout 30m ./host/daemon/ ;;
*) echo "unknown $M $P"; exit 2 ;;
esac
echo "ACCEPT_OK $M $P"
````

### M12a — the unwired production file, its first tests, the two roots · 283 prod + 257 test · Run 1 · part: `fast`

- `cmd/ailang-worldd/episode_finish_surface.go` (NEW, **no caller in M12a**): `verbLandGrace = 10 * time.Second`;
  `surfaceFinish`, `verbOwn`; `canonicalFinishArgv` (pure, S5-D5); `finishVerbRequest` (pure: the request with
  `Prompt = broker.FinishPromptText(scope, title, body, protected, committer)`, `Protected` never nil, `Validate`,
  `CanonicalVerb`); `finishConfirmOnSurface` (design §2.7 steps 1–6: request, `env.random` nil/draw refusals, `vr-` +
  8 bytes, 32-byte challenge, `expires = now + 600`, the links (S5-D6), the put, the prompt to the tty, the stdout block,
  the wait, the read-back, `confirmed by <approved_by> (decision <ref>)`); `surfaceCardLinks` (`writeAwaitLinks`' rule, a
  3 s child of finish's ctx); `putFinishVerbRequest` (own 3 s root, non-waiting writer, busy → "the credentials store is
  busy (another writer holds it); run finish again", `PutVerbRequest`'s refusals as "cannot store request <id>: <err>");
  `awaitVerbConfirmation`; `pollVerbRequest` (own 3 s root, the re-check (a), S5-D8); `readBackFinishDecision`
  (`worldReader.object` twice: the decision — `ApprovalDecisionV1`, `approve`, `decidedBy == approved_by` — and its request —
  `ApprovalRequestV1`, `Human.Approve`, `operator:episode-finish`, `scope == scopeText`; every refusal "… does not approve
  this finish (F6): …").
- `host/store/context_roots_test.go`: `+ "cmd/ailang-worldd/episode_finish_surface.go|putFinishVerbRequest|Background": 1`
  and `…|pollVerbRequest|Background": 1` (W23).
- `cmd/ailang-worldd/episode_finish_surface_test.go` (NEW, first half): helpers; **CF14** `TestCanonicalFinishArgv`;
  **CF9** `TestFinishConfirmOnSurfaceRefusesBeforeTheStore` (over 128 KiB, no randomness, the id's draw fails, the
  challenge's draw fails, no credentials store, busy — each: no tty or stdout byte, no row); `TestFinishVerbRequestFields`
  (every field, committer and protected paths in the prompt, `[]` never `null`); `TestSurfaceCardLinks`; **CF8**
  `TestPutFinishVerbRequestRefusals` (no store, busy, 16 pending, below v9 naming the restart).
- Expected numstat after commit **540**.

### M12b — the wait, the re-check, the read-back; judge-403 note 1 · 0 prod + 275 test · Run 1 · parts: `fast`, `daemon`

- `episode_finish_surface_test.go` (second half): **CF4/CF5** `TestAwaitVerbConfirmationEnds` (bounded at
  `expires_at + 10` with 12 waits under a 20 s wall deadline; the exact boundary seen/not seen; superseded; a vanished row
  → "cannot read request …"); **CF6** `TestPollVerbRequestRechecks` (control + eleven rows: the row's digest, challenge or
  expiry not ours; signed over another digest, challenge or expiry; the mint namespace; a key never enrolled; not an
  artifact; the approver revoked since (note 5); `approved_by` not the key's); **CF7** `TestReadBackFinishDecision`
  (control + eight rows, real objects with one field changed); `TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision`
  (P1's "decision does not approve this scope" at the wiring: a re-checked row whose decision the daemon does not hold →
  "(F6)", no decision returned).
- `host/daemon/surface_confirm_now_test.go` (NEW): `TestConfirmPOSTRecordsAtItsClock` (S5-D11).
- Expected numstat after commit **815**.

### M12c — the flag, its refusals, the hook; the rig · 22 prod + 440 test · Run 1 · parts: `fast`, `cmd`

- `cmd/ailang-worldd/episode_finish.go`: the help's usage line and flag paragraph; `onSurface := fs.Bool("confirm-on-surface",
  …)`; refusal 1 after `defer tty.Close()` ("--confirm-on-surface cannot be combined with --dry-run" / "a dry run opens no pull
  request" / "drop one of them"); refusal 2 after `validateFinishChecks`, before the profile and the check loop ("check
  <name> is unsandboxed: it needs its own typed y" / "--confirm-on-surface replaces only the pull request's y" / "run
  finish without --confirm-on-surface"); `} else if *onSurface {` between the reusable branch and the typed prompt:
  `finishConfirmOnSurface(ctx, c, env, surfaceFinish{credentials, p.Addr, ep, sha, canonicalFinishArgv(fs, allowed, ep),
  scope, scopeText, actualTitle, rendered, pinned.Committer, protected}, tty, stdout)`; a refusal → `fail(…, "…; commit
  kept at "+sha, steps(sha))`; else `approval.Approval = ref; finishStep("approve")`.
- `cmd/ailang-worldd/episode_finish_surface_rig_test.go` (NEW): `newSurfaceRig` (S5-D4); **CF3**
  `TestFinishConfirmOnSurfaceOpensThePR`; **CF1, CF2** `…RefusesBeforeAnyCheck`; **CF4 (rig), CF13, CF10**
  `…RerunsAndTypedPrompt` (with a protected path and a pinned committer); **CF5, CF6 (rig, note 5)** `…KeepsTheCommit`;
  **CX21's CLI half** `…DispatchesTheNamedDecision`; **CF12 (note 4)** `…ResumeReusesTheApproval`.
- **CF11**: every pre-existing `TestEpisodeFinish*` (the typed path), unedited, in this part's `fast` run.
- Expected numstat after commit **1,277**.

### M13 — the bar's reason, `approve vr-` in the bar and the CLI · 18 prod + 41 test · Run 2 · parts: `fast`, `daemon`

- `host/verbs/lookup.go`: `case "start"` keeps "episode start changes World. The app cannot run it yet."; `case "finish"`
  → "episode finish changes World. The app cannot run it yet; at a terminal, --confirm-on-surface lets you approve its pull
  request here with your key." (design §2.8 verbatim); `verbRequestID = ^vr-[0-9a-f]{16}$` and the third `ApproveTarget`
  arm → `/surface/confirm/<id>`.
- `host/surfacepage/app/app.js:496`: `|| b.href.startsWith('/surface/confirm/vr-')` (exactly that prefix).
- `cmd/ailang-worldd/approve.go`: before the `sha256:` arm, `store.VerbRequestIDPattern` → `refuse("request %s is a finish
  confirmation: confirm it on its card, /surface/confirm/%s; the waiting episode finish opens the pull request")` (exit 1).
- Tests: **CB1** `TestLookup` (both reasons exact), **CB2** `TestApproveTarget` (+ `vr-` 16 hex; `vr-0123`, 17 hex, upper
  case, `xvr-` refused) and `TestSurfaceVerbOutcomes` (+ the case; the golden regenerated, C9), **CB3** `TestSurfaceAppVerbs`
  (the `vr-` navigate assigns; `/surface/confirm/x` and `/surface/confirm/mr-…` → "refused: unexpected target"), **CB4**
  `TestApproveVerbRequestPointsAtItsCard` (NEW; exact stderr, exit 1, a short id still the mint refusal).
- Expected numstat after commit **1,336**.

### M14 — QUICKSTART §20 + CQ1 · 12 docs + 50 test · Run 2 · parts: `fast`, `docs`, `cmd`

- `docs/QUICKSTART.md` §20 *Confirm a finish's pull request with your key* (attended — pending first verbatim run): the
  command, what it prints, the card, `approve vr-` in the bar and at the CLI, the 10 minutes + 10 s, every refusal
  keeping the commit, Ctrl-C, re-run.
- `cmd/ailang-worldd/quickstart_finish_surface_test.go` (NEW): **CQ1** (S5-D3).
- Expected numstat after commit **1,398**.

### 3.1 Every input and failure class S5's code reaches, traced (design §2.9 P1 and P6; directive retro 4)

| design row (§2.9) | branch | test (CF id) |
|---|---|---|
| P1 no tty | F1, unchanged | `TestEpisodeFinishFenceAC1` (existing, unedited) |
| P1 with `--dry-run` | refusal 1 | `TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck` (CF1: exit 1, the text, **0 requests read**, tty unread) |
| P1 an unsandboxed pinned check | refusal 2 | same test (CF2: the terminal would say y; **0 outcomes**, no marker file, HEAD = base, no row) |
| P1 mode `finished` | unchanged | `TestEpisodeFinishFinishedRerun` (existing) |
| P1 resume-dispatch + reusable approval | the unchanged branch | `TestFinishConfirmOnSurfaceResumeReusesTheApproval` (CF12) |
| P1 request over 128 KiB | `Validate` | `TestFinishConfirmOnSurfaceRefusesBeforeTheStore` (CF9) |
| P1 `env.random` nil / failing (either draw) | step 2 | same (CF9) |
| P1 links unreadable (S5-D6) | before the put | same (CF9 "no credentials store") |
| P1 writer busy | step 3 | `TestPutFinishVerbRequestRefusals` (CF8) and CF9's "busy" (through `finishConfirmOnSurface`) |
| P1 store below v9 | step 3 | CF8 (the `OlderCredentialsStoreError` text naming the restart) |
| P1 16 pending | step 3 | CF8 |
| P1 card confirms | steps 5–6 | `TestFinishConfirmOnSurfaceOpensThePR` (CF3: the row, the stdout block, `confirmed by ssh:mark-mac (decision …)`, the fake gh's PR URL last, the record names the decision, the tty never read and showing exactly the stored prompt, one sandboxed check run, the argv) |
| P1 decision does not approve this scope | step 6(b) | `TestReadBackFinishDecision` (CF7, eight rows) + `TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision` |
| P1 the re-check fails (digest/challenge/expiry not ours; artifact; revoked; `approved_by`) | step 6(a) | `TestPollVerbRequestRechecks` (CF6, eleven rows) + `KeepsTheCommit` (rig: revoked after C15, note 5) |
| P1 superseded by a newer run | step 5 | `TestAwaitVerbConfirmationEnds/superseded` + `KeepsTheCommit` (rig, CF5) |
| P1 expiry (+ grace), exact boundary | step 5 | `TestAwaitVerbConfirmationEnds` (CF4: bounded, 12 waits; `+10` seen, `+11` not) + `RerunsAndTypedPrompt` (rig) |
| P1 a poll read fails | step 5 | `TestAwaitVerbConfirmationEnds/a poll read fails` (CF4) |
| P1 Ctrl-C | process ends | RV1 (not testable without a signal; the request stays pending until expiry or the next run's supersession, which CF10 shows) |
| P1 re-run after any refusal (incl. after a C14 that succeeded) | `resume`, new request | `RerunsAndTypedPrompt` (CF10: same prompt bytes, new id, the old row superseded at the re-run's clock) + `DispatchesTheNamedDecision` (an orphaned decision is never used) |
| P1 the typed path (no flag) | unchanged | CF11: every `TestEpisodeFinish*`, unedited; CF13: the typed prompt = the stored prompt + the question |
| S4 §1.4 F4 (CX21's CLI half) | step 6 | `TestFinishConfirmOnSurfaceDispatchesTheNamedDecision` (C15 fails → 500, 1 orphan; retry a second later → 200; the CLI dispatches with the row's decision; 2 decisions in the store, the record names the row's) |
| P6 `episode finish w1 --preset p` → terminal | `Lookup` | `TestLookup`, `TestSurfaceVerbOutcomes` (CB1) |
| P6 `approve vr-0123456789abcdef` → navigate; `approve vr-0123` → refused | `ApproveTarget` | `TestApproveTarget`, `TestSurfaceVerbOutcomes` (CB2) |
| P6 a navigate answer with `/surface/confirm/x` → "refused: unexpected target" | `app.js:496` | `TestSurfaceAppVerbs` (CB3) |
| P6 CLI `approve vr-…` → the card pointer | `approve.go` arm | `TestApproveVerbRequestPointsAtItsCard` (CB4) |

---

## 4. Mutation table

Enumerated **from the diff** (every new or changed production branch, every argument of the put and of the call that
builds the request, every time comparison's operator, the grace constant, the request lifetime and the id/challenge sizes,
the two root budgets, the bar's reason, `ApproveTarget`'s `vr-` shape, the page's prefix and the CLI arm), not from the
test plan. Each is one exact replacement (`tools/mutants.py`; a replacement whose text is not found exactly once is "not a
mutant", as is one that does not compile — each of those is re-run in a compiling `-2` form). Each runs against its
**named** test row, then against **`suite-nostatic`** — the touched package's behavioural tests (S5-D15; no inventory,
AST or string-grep test is in it: for `cmd`, every new test, `TestEpisodeFinish*`, `TestApprove*`, the A3 e2e). Ledger:
`$HOME/.ailang/state/world-iter404/plan-scratch/mutants/ledger.txt`, one line per run; **the ids below are the
ledger's**. Classes: **K** killed by both; **E** a named equivalent, with its reason; **CTL** a control that must stay
green. `J-NOW-C14`'s control line (C19) records its survival against S4's tests alone.

| id | file | mutation | named test | named | suite-nostatic (killing tests) | class |
|---|---|---|---|---|---|---|
| `R1-DEL` | `episode_finish.go` | `if *onSurface && *dry {` → `if false && *onSurface && *dry {` | `TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck) | K |
| `R1-OR` | `episode_finish.go` | `if *onSurface && *dry {` → `if *onSurface \|\| *dry {` | `TestEpisodeFinishAttributionCommitDryRunAC5AC6\|TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestEpisodeFinishAttributionCommitDryRunAC5AC6, TestEpisodeFinishFixupReattach, TestEpisodeFinishUnsandboxedAC4, TestFinishConfirmOnSurfaceDispatchesTheNamedDec) | K |
| `R2-DEL` | `episode_finish.go` | `if *onSurface && check.Unsandboxed {` → `if false {` | `TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck) | K |
| `R2-ANY` | `episode_finish.go` | `if *onSurface && check.Unsandboxed {` → `if *onSurface {` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `R2-TYPED` | `episode_finish.go` | `if *onSurface && check.Unsandboxed {` → `if check.Unsandboxed {` | `TestEpisodeFinishUnsandboxedAC4` | killed | killed (TestEpisodeFinishProtected, TestEpisodeFinishUnsandboxedAC4) | K |
| `R2-AFTER` | `episode_finish.go` | `for _, check := range start.Checks { // row 243 V3 §2.7 step 2: before` → `(deleted)`; `// Refresh both evidence and paths after checks, which may produce fil` → `for _, check := range start.Checks { if *onSurface && check.Unsandboxe` | `TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck) | K |
| `H-TYPED` | `episode_finish.go` | `} else if *onSurface { //` → `} else if false { //` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `H-TTY-READ` | `episode_finish.go` | `ref, refusal := finishConfirmOnSurface(` → `_, _ = reader.ReadString('\n') ref, refusal := finishConfirmOnSurface(` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `H-REFUSAL` | `episode_finish.go` | `if refusal != "" { return fail(refusal, "the pull request was not conf` → `if false { return fail(refusal, "the pull request was not confirmed` | `TestFinishConfirmOnSurfaceKeepsTheCommit` | killed | killed (TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceRerunsAndTypedPrompt) | K |
| `H-APPROVAL` | `episode_finish.go` | `approval.Approval = ref` → `(deleted)` | `TestFinishConfirmOnSurfaceOpensThePR` | — | — | NOT A MUTANT (no compile) — re-run as its `-2` form |
| `H-STEP` | `episode_finish.go` | `approval.Approval = ref finishStep("approve")` → `approval.Approval = ref` | `TestFinishConfirmOnSurfaceOpensThePR` | survived | survived (-) | E — `finishStep` is a test seam (`finishAfterStep` is nil in production); dropping the call changes no behaviour |
| `ARG-CREDS` | `episode_finish.go` | `surfaceFinish{credentials: credentials,` → `surfaceFinish{credentials: p.DB,` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `ARG-ADDR` | `episode_finish.go` | `addr: p.Addr,` → `addr: "",` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `ARG-EP` | `episode_finish.go` | `addr: p.Addr, ep: ep,` → `addr: p.Addr, ep: "ep2",` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `ARG-SHA` | `episode_finish.go` | `ep: ep, sha: sha,` → `ep: ep, sha: previous,` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `ARG-ARGV` | `episode_finish.go` | `argv: canonicalFinishArgv(fs, allowed, ep)` → `argv: canonicalFinishArgv(fs, nil, ep)` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `ARG-SCOPE` | `episode_finish.go` | `scope: scope, scopeText` → `scope: broker.FinishApprovalScope{}, scopeText` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `ARG-SCOPETEXT` | `episode_finish.go` | `scopeText: scopeText, title` → `scopeText: scopeText + " ", title` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `ARG-TITLE` | `episode_finish.go` | `title: actualTitle, body: rendered,` → `title: actualTitle + ".", body: rendered,` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `ARG-BODY` | `episode_finish.go` | `body: rendered,` → `body: rendered + "\n",` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `ARG-COMMITTER` | `episode_finish.go` | `committer: pinned.Committer, protected` → `committer: "", protected` | `TestFinishConfirmOnSurfaceRerunsAndTypedPrompt` | killed | killed (TestFinishConfirmOnSurfaceRerunsAndTypedPrompt) | K |
| `ARG-PROTECTED` | `episode_finish.go` | `protected: protected}, tty, stdout)` → `protected: nil}, tty, stdout)` | `TestFinishConfirmOnSurfaceRerunsAndTypedPrompt` | killed | killed (TestFinishConfirmOnSurfaceRerunsAndTypedPrompt) | K |
| `A-ALLOW` | `episode_finish_surface.go` | `case f.Name == "allow-failing":` → `case f.Name == "allow-failing-x":` | `TestCanonicalFinishArgv` | killed | killed (TestCanonicalFinishArgv) | K |
| `A-ALLOW-ORDER` | `episode_finish_surface.go` | `for _, v := range allowed { argv = append(argv, "--allow-failing", v) ` → `for i := range allowed { argv = append(argv, "--allow-failing", allowe` | `TestCanonicalFinishArgv` | killed | killed (TestCanonicalFinishArgv) | K |
| `A-BOOL-TRUE` | `episode_finish_surface.go` | `if f.Value.String() == "true" {` → `if true {` | `TestCanonicalFinishArgv` | killed | killed (TestCanonicalFinishArgv) | K |
| `A-BOOL-FORM` | `episode_finish_surface.go` | `argv = append(argv, "--"+f.Name)` → `argv = append(argv, "--"+f.Name, "true")` | `TestCanonicalFinishArgv` | killed | killed (TestCanonicalFinishArgv, TestFinishConfirmOnSurfaceOpensThePR) | K |
| `A-VALUE` | `episode_finish_surface.go` | `argv = append(argv, "--"+f.Name, f.Value.String())` → `argv = append(argv, "--"+f.Name+"="+f.Value.String())` | `TestCanonicalFinishArgv` | killed | killed (TestCanonicalFinishArgv, TestFinishConfirmOnSurfaceOpensThePR) | K |
| `A-PREFIX` | `episode_finish_surface.go` | `argv := []string{"ailang-worldd", "episode", "finish", ep}` → `argv := []string{"ailang-worldd", "episode", "finish", ep, "--confirm-` | `TestCanonicalFinishArgv` | killed | killed (TestCanonicalFinishArgv, TestFinishConfirmOnSurfaceOpensThePR) | K |
| `A-VISIT` | `episode_finish_surface.go` | `fs.Visit(func(f *flag.Flag) {` → `fs.VisitAll(func(f *flag.Flag) {` | `TestCanonicalFinishArgv` | killed | killed (TestCanonicalFinishArgv, TestFinishConfirmOnSurfaceOpensThePR) | K |
| `V-PROT-COPY` | `episode_finish_surface.go` | `protected := append([]string{}, in.protected...)` → `protected := in.protected` | `TestFinishVerbRequestFields` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR, TestFinishVerbRequestFields) | K |
| `V-PROMPT-COMMITTER` | `episode_finish_surface.go` | `broker.FinishPromptText(in.scope, in.title, in.body, protected, in.com` → `broker.FinishPromptText(in.scope, in.title, in.body, protected, "")` | `TestFinishVerbRequestFields` | killed | killed (TestFinishConfirmOnSurfaceRerunsAndTypedPrompt, TestFinishVerbRequestFields) | K |
| `V-PROMPT-PROTECTED` | `episode_finish_surface.go` | `broker.FinishPromptText(in.scope, in.title, in.body, protected, in.com` → `broker.FinishPromptText(in.scope, in.title, in.body, nil, in.committer` | `TestFinishVerbRequestFields` | killed | killed (TestFinishConfirmOnSurfaceRerunsAndTypedPrompt, TestFinishVerbRequestFields) | K |
| `V-COMMITTER` | `episode_finish_surface.go` | `Committer: in.committer, Protected: protected}` → `Committer: "", Protected: protected}` | `TestFinishVerbRequestFields` | killed | killed (TestFinishConfirmOnSurfaceRerunsAndTypedPrompt, TestFinishVerbRequestFields) | K |
| `V-SCOPE` | `episode_finish_surface.go` | `Scope: in.scopeText,` → `Scope: in.scopeText + "&x=1",` | `TestFinishVerbRequestFields` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `V-VALIDATE` | `episode_finish_surface.go` | `if err := r.Validate(); err != nil { return r, nil, err }` → `(deleted)` | `TestFinishConfirmOnSurfaceRefusesBeforeTheStore` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeTheStore) | K |
| `F-RAND-NIL` | `episode_finish_surface.go` | `if env.random == nil {` → `if false {` | `TestFinishConfirmOnSurfaceRefusesBeforeTheStore` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeTheStore) | K |
| `F-RAND-ID` | `episode_finish_surface.go` | `if err := env.random(idb[:]); err != nil { return "", "randomness: " +` → `_ = env.random(idb[:])` | `TestFinishConfirmOnSurfaceRefusesBeforeTheStore` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeTheStore) | K |
| `F-RAND-CH` | `episode_finish_surface.go` | `if err := env.random(ch[:]); err != nil { return "", "randomness: " + ` → `_ = env.random(ch[:])` | `TestFinishConfirmOnSurfaceRefusesBeforeTheStore` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeTheStore) | K |
| `F-LIFETIME` | `episode_finish_surface.go` | `expires: now + int64(approval.VerbRequestLifetime/time.Second)}` → `expires: now + int64(approval.VerbRequestLifetime/time.Second) + 1}` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `F-ID-PREFIX` | `episode_finish_surface.go` | `own := verbOwn{id: "vr-" +` → `own := verbOwn{id: "vr-0" +` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `F-LINKS-ERR` | `episode_finish_surface.go` | `links, err := surfaceCardLinks(ctx, in.credentials, "/surface/confirm/` → `links, err := surfaceCardLinks(ctx, in.credentials, "/surface/confirm/` | `TestFinishConfirmOnSurfaceRefusesBeforeTheStore` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeTheStore) | K |
| `F-LINKS-PATH` | `episode_finish_surface.go` | `"/surface/confirm/"+own.id, in.addr)` → `"/surface/approve/"+own.id, in.addr)` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `F-PUT-REFUSAL` | `episode_finish_surface.go` | `if refusal := putFinishVerbRequest(in.credentials, row, now); refusal ` → `if refusal := putFinishVerbRequest(in.credentials, row, now); false &&` | `TestFinishConfirmOnSurfaceRefusesBeforeTheStore` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeTheStore) | K |
| `PUT-ID` | `episode_finish_surface.go` | `row := store.VerbRequestRow{RequestID: own.id,` → `row := store.VerbRequestRow{RequestID: "vr-00000000000000aa",` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `PUT-REQUEST` | `episode_finish_surface.go` | `Request: canon, Digest:` → `Request: append(canon, ' '), Digest:` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `PUT-DIGEST` | `episode_finish_surface.go` | `Digest: hex.EncodeToString(own.digest[:]),` → `Digest: hex.EncodeToString(own.challenge[:]),` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `PUT-CHALLENGE` | `episode_finish_surface.go` | `Challenge: hex.EncodeToString(own.challenge[:]), EpisodeID` → `Challenge: hex.EncodeToString(own.digest[:]), EpisodeID` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `PUT-EPISODE` | `episode_finish_surface.go` | `EpisodeID: in.ep, CreatedAt: now,` → `EpisodeID: in.ep + "x", CreatedAt: now,` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR) | K |
| `PUT-CREATED` | `episode_finish_surface.go` | `CreatedAt: now, ExpiresAt: own.expires}` → `CreatedAt: now - 1, ExpiresAt: own.expires}` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `PUT-EXPIRES` | `episode_finish_surface.go` | `CreatedAt: now, ExpiresAt: own.expires}` → `CreatedAt: now, ExpiresAt: own.expires + 1}` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `PUT-NOW` | `episode_finish_surface.go` | `putFinishVerbRequest(in.credentials, row, now)` → `putFinishVerbRequest(in.credentials, row, now+1)` | `TestFinishConfirmOnSurfaceRerunsAndTypedPrompt` | killed | killed (TestFinishConfirmOnSurfaceRerunsAndTypedPrompt) | K |
| `F-TTY-PROMPT` | `episode_finish_surface.go` | `fmt.Fprint(tty, r.Prompt) //` → `_ = tty //` | `TestFinishConfirmOnSurfaceOpensThePR` | — | — | NOT A MUTANT (no compile) — re-run as its `-2` form |
| `F-HASH-8` | `episode_finish_surface.go` | `hex.EncodeToString(h[:8]))` → `hex.EncodeToString(h[:7]))` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `F-KEEPS-LINE` | `episode_finish_surface.go` | `keeps the commit at %s\n", in.sha)` → `keeps the commit at %s\n", in.ep)` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `F-AWAIT-REFUSAL` | `episode_finish_surface.go` | `got, refusal := awaitVerbConfirmation(env, in.credentials, own) if ref` → `got, refusal := awaitVerbConfirmation(env, in.credentials, own) if fal` | `TestFinishConfirmOnSurfaceKeepsTheCommit` | killed | killed (TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceRerunsAndTypedPrompt) | K |
| `F-READBACK` | `episode_finish_surface.go` | `got.DecisionRef, got.ApprovedBy, in.scopeText); refusal != "" {` → `got.DecisionRef, got.ApprovedBy, in.scopeText); false && refusal != ""` | `TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision` | killed | killed (TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision) | K |
| `F-READBACK-SCOPE` | `episode_finish_surface.go` | `got.DecisionRef, got.ApprovedBy, in.scopeText)` → `got.DecisionRef, got.ApprovedBy, got.DecisionRef)` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurfaceRerunsAndTypedPrompt, TestFinishConfirmOnS) | K |
| `F-CONFIRMED-LINE` | `episode_finish_surface.go` | `fmt.Fprintf(stdout, "confirmed by %s (decision %s)\n", got.ApprovedBy,` → `fmt.Fprintf(stdout, "confirmed by %s\n", got.ApprovedBy)` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `F-RETURN-REF` | `episode_finish_surface.go` | `return got.DecisionRef, ""` → `return "sha256:" + strings.Repeat("0", 64), ""` | `TestFinishConfirmOnSurfaceOpensThePR` | — | — | NOT A MUTANT (no compile) — re-run as its `-2` form |
| `L-DEDUPE` | `episode_finish_surface.go` | `if !seen[e.Passkey.Origin] {` → `if true {` | `TestSurfaceCardLinks` | killed | killed (TestSurfaceCardLinks) | K |
| `L-LOOPBACK-SEEN` | `episode_finish_surface.go` | `err == nil && !seen["http://localhost:"+port]` → `err == nil` | `TestSurfaceCardLinks` | killed | killed (TestSurfaceCardLinks) | K |
| `L-NO-PASSKEYS` | `episode_finish_surface.go` | `for _, e := range active {` → `for _, e := range active[:0] {` | `TestSurfaceCardLinks` | killed | killed (TestSurfaceCardLinks) | K |
| `L-ACTIVE-ERR` | `episode_finish_surface.go` | `_ = ro.Close() if err != nil { return nil, err }` → `_ = ro.Close() if false && err != nil { return nil, err }` | `TestFinishConfirmOnSurfaceRefusesBeforeTheStore` | survived | survived (-) | E (class) — `ActivePasskeys` refuses only a store below schema 6, which the put refuses right after (below 9) with the same `serve`-restart remedy; nothing is stored either way, only the refusal's wording differs |
| `L-BUDGET` | `episode_finish_surface.go` | `ctx, cancel := context.WithTimeout(ctx, approveReadBudget)` → `ctx, cancel := context.WithTimeout(ctx, 0)` | `TestSurfaceCardLinks` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurface) | K |
| `P-BUSY` | `episode_finish_surface.go` | `if store.IsWriterAlreadyActive(err) {` → `if false {` | `TestPutFinishVerbRequestRefusals` | killed | killed (TestFinishConfirmOnSurfaceRefusesBeforeTheStore, TestPutFinishVerbRequestRefusals) | K |
| `P-OPEN-ERR` | `episode_finish_surface.go` | `if err != nil { return fmt.Sprintf("cannot open the credentials store ` → `if false { return fmt.Sprintf("cannot open the credentials store %s: %` | `TestPutFinishVerbRequestRefusals` | killed | killed (TestPutFinishVerbRequestRefusals) | K |
| `P-PUT-ERR` | `episode_finish_surface.go` | `if err := st.PutVerbRequest(ctx, row, now); err != nil {` → `if err := st.PutVerbRequest(ctx, row, now); false && err != nil {` | `TestPutFinishVerbRequestRefusals` | killed | killed (TestPutFinishVerbRequestRefusals) | K |
| `P-BUDGET` | `episode_finish_surface.go` | `func putFinishVerbRequest(credPath string, row store.VerbRequestRow, n` → `func putFinishVerbRequest(credPath string, row store.VerbRequestRow, n` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestAwaitVerbConfirmationEnds, TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensTh) | K |
| `W-GRACE` | `episode_finish_surface.go` | `deadline := own.expires + int64(verbLandGrace/time.Second)` → `deadline := own.expires` | `TestAwaitVerbConfirmationEnds` | killed | killed (TestAwaitVerbConfirmationEnds) | K |
| `W-GRACE-CONST` | `episode_finish_surface.go` | `const verbLandGrace = 10 * time.Second` → `const verbLandGrace = 3 * time.Second` | `TestAwaitVerbConfirmationEnds` | killed | killed (TestAwaitVerbConfirmationEnds) | K |
| `J-W-GE-GT` | `episode_finish_surface.go` | `case now >= deadline:` → `case now > deadline:` | `TestAwaitVerbConfirmationEnds` | killed | killed (TestAwaitVerbConfirmationEnds, TestFinishConfirmOnSurfaceRerunsAndTypedPrompt) | K |
| `J-W-GE-EQ` | `episode_finish_surface.go` | `case now >= deadline:` → `case now == deadline+1:` | `TestAwaitVerbConfirmationEnds` | killed | killed (TestAwaitVerbConfirmationEnds, TestFinishConfirmOnSurfaceRerunsAndTypedPrompt) | K |
| `W-UNBOUNDED` | `episode_finish_surface.go` | `case now >= deadline:` → `case false:` | `TestAwaitVerbConfirmationEnds` | — | — | NOT A MUTANT (no compile) — re-run as its `-2` form |
| `W-READ-ERR` | `episode_finish_surface.go` | `case err != nil: return store.VerbRequestRow{}, fmt.Sprintf("cannot re` → `case false: return store.VerbRequestRow{}, fmt.Sprintf("cannot read re` | `TestAwaitVerbConfirmationEnds` | killed | killed (TestAwaitVerbConfirmationEnds) | K |
| `W-REFUSAL` | `episode_finish_surface.go` | `case refusal != "": return store.VerbRequestRow{}, refusal` → `case false: return store.VerbRequestRow{}, refusal` | `TestFinishConfirmOnSurfaceKeepsTheCommit` | killed | killed (TestFinishConfirmOnSurfaceKeepsTheCommit) | K |
| `W-APPROVED` | `episode_finish_surface.go` | `case row.Outcome == store.VerbOutcomeApproved:` → `case false:` | `TestAwaitVerbConfirmationEnds\|TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestAwaitVerbConfirmationEnds, TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurfaceRefusesAnUn) | K |
| `W-SUPERSEDED` | `episode_finish_surface.go` | `case row.Outcome == store.VerbOutcomeSuperseded:` → `case false:` | `TestAwaitVerbConfirmationEnds\|TestFinishConfirmOnSurfaceKeepsTheCommit` | killed | killed (TestAwaitVerbConfirmationEnds, TestFinishConfirmOnSurfaceKeepsTheCommit) | K |
| `W-INTERVAL` | `episode_finish_surface.go` | `<-env.after(awaitPollInterval) } } // pollVerbRequest` → `<-env.after(2 * awaitPollInterval) } } // pollVerbRequest` | `TestAwaitVerbConfirmationEnds` | killed | killed (TestAwaitVerbConfirmationEnds) | K |
| `Q-UNKNOWN` | `episode_finish_surface.go` | `if !ok { return row, "", fmt.Errorf("request %s is no longer in the st` → `if false { return row, "", fmt.Errorf("request %s is no longer in the ` | `TestAwaitVerbConfirmationEnds` | — | — | NOT A MUTANT (no compile) — re-run as its `-2` form |
| `Q-PENDING` | `episode_finish_surface.go` | `if row.Outcome != store.VerbOutcomeApproved { return row, "", nil }` → `if row.Outcome == "superseded" { return row, "", nil }` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestAwaitVerbConfirmationEnds, TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensTh) | K |
| `Q-EQ-DIGEST` | `episode_finish_surface.go` | `if row.Digest != hex.EncodeToString(own.digest[:]) \|\|` → `if` | `TestPollVerbRequestRechecks` | killed | killed (TestPollVerbRequestRechecks) | K |
| `Q-EQ-CHALLENGE` | `episode_finish_surface.go` | `\|\| row.Challenge != hex.EncodeToString(own.challenge[:]) \|\|` → `\|\|` | `TestPollVerbRequestRechecks` | killed | killed (TestPollVerbRequestRechecks) | K |
| `Q-EQ-EXPIRES` | `episode_finish_surface.go` | `\|\| row.ExpiresAt != own.expires {` → `{` | `TestPollVerbRequestRechecks` | killed | killed (TestPollVerbRequestRechecks) | K |
| `Q-RECHECK` | `episode_finish_surface.go` | `signed, err := authority.RecheckVerb(ctx, ro, in, row.Artifact) if err` → `signed, err := authority.RecheckVerb(ctx, ro, in, row.Artifact) if err` | `TestPollVerbRequestRechecks` | killed | killed (TestFinishConfirmOnSurfaceKeepsTheCommit, TestPollVerbRequestRechecks) | K |
| `Q-ROW-VALUES` | `episode_finish_surface.go` | `in := authority.VerbJudgeInput{Digest: own.digest, Challenge: own.chal` → `rd, _ := hex.DecodeString(row.Digest) rc, _ := hex.DecodeString(row.Ch` | `TestPollVerbRequestRechecks` | survived | survived (-) | E — under the equality guard the row's digest, challenge and expiry equal this process's own (S5-D8); its pair `Q-ROW-VALUES-NOEQ` is killed |
| `Q-ROW-VALUES-NOEQ` | `episode_finish_surface.go` | `in := authority.VerbJudgeInput{Digest: own.digest, Challenge: own.chal` → `rd, _ := hex.DecodeString(row.Digest) rc, _ := hex.DecodeString(row.Ch`; `if row.Digest != hex.EncodeToString(own.digest[:]) \|\| row.Challenge !=` → `if false {` | `TestPollVerbRequestRechecks` | killed | killed (TestPollVerbRequestRechecks) | K |
| `Q-NOW` | `episode_finish_surface.go` | `ExpiresAt: own.expires, Now: now}` → `ExpiresAt: own.expires, Now: 0}` | `TestPollVerbRequestRechecks` | killed | killed (TestAwaitVerbConfirmationEnds, TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurfaceRefusesAnUn) | K |
| `Q-RECHECK-EXPIRY` | `episode_finish_surface.go` | `ExpiresAt: own.expires, Now: now}` → `ExpiresAt: own.expires + 1, Now: now}` | `TestPollVerbRequestRechecks` | killed | killed (TestAwaitVerbConfirmationEnds, TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurfaceRefusesAnUn) | K |
| `Q-APPROVEDBY` | `episode_finish_surface.go` | `if signed.ApprovedBy() != row.ApprovedBy {` → `if false {` | `TestPollVerbRequestRechecks` | killed | killed (TestPollVerbRequestRechecks) | K |
| `Q-BUDGET` | `episode_finish_surface.go` | `func pollVerbRequest(credPath string, own verbOwn, now int64) (store.V` → `func pollVerbRequest(credPath string, own verbOwn, now int64) (store.V` | `TestAwaitVerbConfirmationEnds` | killed | killed (TestAwaitVerbConfirmationEnds, TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceKeepsTheCommit, TestFinishConfirmOnSurfaceOpensTh) | K |
| `B-DEC-SEMANTIC` | `episode_finish_surface.go` | `\|\| obj.SemanticID != broker.ApprovalDecisionV1 {` → `{` | `TestReadBackFinishDecision` | killed | killed (TestReadBackFinishDecision) | K |
| `B-DEC-FOUND` | `episode_finish_surface.go` | `if err != nil \|\| !ok \|\| obj.SemanticID != broker.ApprovalDecisionV1 {` → `if err != nil \|\| obj.SemanticID != broker.ApprovalDecisionV1 {` | `TestReadBackFinishDecision` | survived | survived (-) | E — a missing object is the zero `Object`, whose empty `SemanticID` fails the same condition's `ApprovalDecisionV1` test (the 404 row of CF7 is refused either way) |
| `B-APPROVE` | `episode_finish_surface.go` | `if d.Decision != "approve" \|\| d.DecidedBy != approvedBy {` → `if d.DecidedBy != approvedBy {` | `TestReadBackFinishDecision` | killed | killed (TestReadBackFinishDecision) | K |
| `B-DECIDEDBY` | `episode_finish_surface.go` | `if d.Decision != "approve" \|\| d.DecidedBy != approvedBy {` → `if d.Decision != "approve" {` | `TestReadBackFinishDecision` | killed | killed (TestReadBackFinishDecision) | K |
| `B-REQ-SEMANTIC` | `episode_finish_surface.go` | `\|\| qobj.SemanticID != broker.ApprovalRequestV1 {` → `{` | `TestReadBackFinishDecision` | killed | killed (TestReadBackFinishDecision) | K |
| `B-REQ-FOUND` | `episode_finish_surface.go` | `if err != nil \|\| !ok \|\| qobj.SemanticID` → `if err != nil \|\| qobj.SemanticID` | `TestReadBackFinishDecision` | survived | survived (-) | E — as `B-DEC-FOUND`: the zero object fails the `ApprovalRequestV1` test (CF7 "its request is missing" refused either way) |
| `B-EFFECT` | `episode_finish_surface.go` | `if q.Effect != broker.EffectHumanApprove \|\|` → `if` | `TestReadBackFinishDecision` | killed | killed (TestReadBackFinishDecision) | K |
| `B-REQUESTER` | `episode_finish_surface.go` | `\|\| q.Requester != "operator:episode-finish" \|\|` → `\|\|` | `TestReadBackFinishDecision` | killed | killed (TestReadBackFinishDecision) | K |
| `B-SCOPE` | `episode_finish_surface.go` | `\|\| q.Scope != scopeText {` → `{` | `TestReadBackFinishDecision` | killed | killed (TestReadBackFinishDecision) | K |
| `B-JSON` | `episode_finish_surface.go` | `if err := json.Unmarshal(obj.Payload, &d); err != nil { return bad(err` → `_ = json.Unmarshal(obj.Payload, &d)` | `TestReadBackFinishDecision` | survived | survived (-) | E — an undecodable payload leaves `d` zero, and `Decision` "" is refused by the next condition; only the reason text differs |
| `B-DECODE` | `episode_finish_surface.go` | `q, err := broker.DecodeApprovalRequest(qobj.Payload) if err != nil { r` → `q, _ := broker.DecodeApprovalRequest(qobj.Payload)` | `TestReadBackFinishDecision` | survived | survived (-) | E — an undecodable request leaves `q` zero, and `Effect` "" is refused by the next condition; only the reason text differs |
| `B-REQREF` | `episode_finish_surface.go` | `qobj, ok, err := r.object(d.RequestRef)` → `qobj, ok, err := r.object(ref)` | `TestReadBackFinishDecision` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurfaceRerunsAndTypedPrompt, TestFinishConfirmOnS) | K |
| `J-NOW-C14` | `surface_confirm.go` | `confirmRequester, signed.ApprovedBy(), now) // C14` → `confirmRequester, signed.ApprovedBy(), now+1) // C14` | `TestConfirmPOSTRecordsAtItsClock` | killed | killed (TestConfirmPOSTRecordsAtItsClock) | K |
| `J-NOW-C14@cmd` | `surface_confirm.go` | `confirmRequester, signed.ApprovedBy(), now) // C14` → `confirmRequester, signed.ApprovedBy(), now+1) // C14` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `CB-REASON` | `lookup.go` | `The app cannot run it yet; at a terminal, --confirm-on-surface lets yo` → `The app cannot run it yet; at a terminal, --confirm-on-surface lets yo` | `TestLookup` | killed | killed (TestLookup) | K |
| `CB-REASON@golden` | `lookup.go` | `The app cannot run it yet; at a terminal, --confirm-on-surface lets yo` → `The app cannot run it yet; at a terminal, --confirm-on-surface lets yo` | `TestSurfaceVerbOutcomes` | killed | killed (TestSurfaceVerbOutcomes) | K |
| `CB-START` | `lookup.go` | `Reason: "episode start changes World. The app cannot run it yet."}` → `Reason: "episode start changes World. The app cannot run it yet; at a ` | `TestLookup` | killed | killed (TestLookup) | K |
| `CB-VR-ARM` | `lookup.go` | `if len(args) == 1 && verbRequestID.MatchString(args[0]) {` → `if false {` | `TestApproveTarget` | killed | killed (TestApproveTarget) | K |
| `CB-VR-ARM@daemon` | `lookup.go` | `if len(args) == 1 && verbRequestID.MatchString(args[0]) {` → `if false {` | `TestSurfaceVerbOutcomes` | killed | killed (TestSurfaceVerbOutcomes) | K |
| `CB-VR-LEN` | `lookup.go` | `regexp.MustCompile('^vr-[0-9a-f]{16}$')` → `regexp.MustCompile('^vr-[0-9a-f]{16,}$')` | `TestApproveTarget` | killed | killed (TestApproveTarget) | K |
| `CB-VR-SHORT` | `lookup.go` | `regexp.MustCompile('^vr-[0-9a-f]{16}$')` → `regexp.MustCompile('^vr-[0-9a-f]{4,16}$')` | `TestApproveTarget` | killed | killed (TestApproveTarget) | K |
| `CB-VR-HEX` | `lookup.go` | `regexp.MustCompile('^vr-[0-9a-f]{16}$')` → `regexp.MustCompile('^vr-[0-9a-fA-F]{16}$')` | `TestApproveTarget` | killed | killed (TestApproveTarget) | K |
| `CB-VR-ANCHOR` | `lookup.go` | `regexp.MustCompile('^vr-[0-9a-f]{16}$')` → `regexp.MustCompile('vr-[0-9a-f]{16}$')` | `TestApproveTarget` | killed | killed (TestApproveTarget) | K |
| `CB-VR-HREF` | `lookup.go` | `return "/surface/confirm/" + args[0], ""` → `return "/surface/approve/" + args[0], ""` | `TestApproveTarget` | killed | killed (TestApproveTarget) | K |
| `CB-JS-ANY` | `app.js` | `b.href.startsWith('/surface/confirm/vr-')` → `b.href.startsWith('/surface/confirm/')` | `TestSurfaceAppVerbs` | killed | killed (TestSurfaceAppVerbs) | K |
| `CB-JS-NONE` | `app.js` | `\|\| b.href.startsWith('/surface/confirm/vr-')` → `(deleted)` | `TestSurfaceAppVerbs` | killed | killed (TestSurfaceAppVerbs) | K |
| `CB-CLI-ARM` | `approve.go` | `if store.VerbRequestIDPattern.MatchString(id) {` → `if false {` | `TestApproveVerbRequestPointsAtItsCard` | killed | killed (TestApproveVerbRequestPointsAtItsCard) | K |
| `CB-CLI-TEXT` | `approve.go` | `confirm it on its card, /surface/confirm/%s;` → `confirm it on its card, /surface/approve/%s;` | `TestApproveVerbRequestPointsAtItsCard` | killed | killed (TestApproveVerbRequestPointsAtItsCard) | K |
| `H-APPROVAL-2` | `episode_finish.go` | `approval.Approval = ref` → `approval.Approval = "" _ = ref` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurfaceResumeReusesTheApproval) | K |
| `F-TTY-PROMPT-2` | `episode_finish_surface.go` | `fmt.Fprint(tty, r.Prompt) //` → `_, _ = tty, r.Prompt //` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceOpensThePR) | K |
| `F-RETURN-REF-2` | `episode_finish_surface.go` | `return got.DecisionRef, ""` → `return "sha256:" + got.Digest, ""` | `TestFinishConfirmOnSurfaceOpensThePR` | killed | killed (TestFinishConfirmOnSurfaceDispatchesTheNamedDecision, TestFinishConfirmOnSurfaceOpensThePR, TestFinishConfirmOnSurfaceResumeReusesTheApproval) | K |
| `W-UNBOUNDED-2` | `episode_finish_surface.go` | `case now >= deadline:` → `case false && now >= deadline:` | `TestAwaitVerbConfirmationEnds` | killed | killed (TIMEOUT(suite)) | K |
| `Q-UNKNOWN-2` | `episode_finish_surface.go` | `if !ok { return row, "", fmt.Errorf("request %s is no longer in the st` → `if false && !ok { return row, "", fmt.Errorf("request %s is no longer ` | `TestAwaitVerbConfirmationEnds` | killed | killed (TestAwaitVerbConfirmationEnds) | K |
| `CTL-COMMENT` | `episode_finish_surface.go` | `// verbLandGrace is how long` → `// verbLandGrace is the time` | `TestAwaitVerbConfirmationEnds\|TestFinishConfirmOnSurfaceOpensThePR` | survived | survived (-) | CTL — a comment edit: must stay green — green |

**Totals: 118 mutants run (controls excluded), 111 killed (K), 7 named equivalents (E), 0 survivors; 5 not-a-mutant lines re-run as their `-2` forms; 1 control(s) green.**

---

## 5. Files (final `git diff --name-status 18c3d132 HEAD`)

```
M  cmd/ailang-worldd/approve.go
A  cmd/ailang-worldd/approve_verb_request_test.go
M  cmd/ailang-worldd/episode_finish.go
A  cmd/ailang-worldd/episode_finish_surface.go
A  cmd/ailang-worldd/episode_finish_surface_rig_test.go
A  cmd/ailang-worldd/episode_finish_surface_test.go
A  cmd/ailang-worldd/quickstart_finish_surface_test.go
M  docs/QUICKSTART.md
A  host/daemon/surface_confirm_now_test.go
M  host/daemon/surface_verb_outcomes_test.go            (pre-existing: + the approve vr- case and golden query, F6)
M  host/store/context_roots_test.go                     (pre-existing: the put and poll roots, design §7)
M  host/surfacepage/app/app.js
M  host/surfacepage/app_verb_test.go                    (pre-existing: CB3's navigate and refused hrefs, F6)
M  host/surfacepage/testdata/verb-route.golden.json     (pre-existing: regenerated, C9)
M  host/verbs/lookup.go
M  host/verbs/lookup_test.go                            (pre-existing: CB1/CB2, design §7)
```

Plus this plan and its JSON under `design_docs/planned/` (the controller's commit). **Production files**:
`cmd/ailang-worldd/{episode_finish_surface.go,episode_finish.go,approve.go}`, `host/verbs/lookup.go`,
`host/surfacepage/app/app.js` — no `host/store`, `host/broker`, `host/approval`, `host/authority` or `host/daemon` production
byte, no `.ail`.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

One bash shell, never under `nohup`: `bash gates.sh <worktree> $HOME/.ailang/state/world-iter404/gates final <HEAD sha>`
(extract it as `accept.sh` is extracted, with the fence `` ````bash gates.sh ``). Every leg is `tee`'d to
`$HOME/.ailang/state/world-iter404/gates/<leg>.txt`; `ALL_DONE.<sha>` is written last with `FAILS: …`. **Each leg runs in its
own `bash -c` with its own rc; inside a leg, every package run ORs its rc (`r=$(( r | $? ))`) — no `&&` chain can skip a
leg or a package.** The executor sets `GATES_SKIP_VERIFY_GO=1`; the controller runs `verify_go` (`VERIFY_GO_LEGS=plain
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`) after the commit. Legs: the **touched-package derivation**
(printed: `git diff --name-only 18c3d132 HEAD -- cmd host` → the `.go` files' directories; the plan's list is
`cmd/ailang-worldd host/daemon host/store host/surfacepage host/verbs`; a difference FAILS), `status`, `verify_ail`, `gofmt`,
`build`, `vet` (`./...`), `fence` (`go test -run '^$' ./...`), `pkgs` (the five touched packages whole, each its own rc),
**`race`** (`host/verbs`, `host/store`, `host/surfacepage` whole; `cmd/ailang-worldd` on the finish/approve/surface tests;
`host/daemon` on the approval-surface and verb-route tests — the two big packages whole under `-race` exceed the leg's budget
at load; CI's race job runs them whole), **`race3`** (`-race -count=3` over the CLI's wait tests — polling, the fake clock,
SQLite reads under a concurrent POST — and the daemon pin), **`node`** (every page harness incl. `TestSurfaceAppVerbs`,
non-skippable, `--- PASS` per name), **`inv`** (the inventories by name with `-v`, incl. `TestDaemonDependencyAllowlist`,
`TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`, `TestEveryVerbHasHelp`,
`TestFrozenV1RouteTableMatchesMux`, `TestConfirmRoutesRegistered`, `TestDaemonCredentialWriterSites`, `TestSurfaceAppPure` (in
`node`), and every new test), **`DAEMON_ROUTES`** (every `mux.Handle*` in `daemon.go`: base = final) and **`SURFACE_POSTS: 3`**,
**`UNTOUCHED_VIOLATIONS: none`** (no `host/store` file but the roots test, no `world/`, `packages/`, `scripts/`, `tools/`,
`daemon.go`, `surface_page.go`, `.ail`), **`chrome`** (S5 changes `app.js`: V2b's fourteen + XC1 + the two drills), then
`verify_go` unless skipped.

````bash gates.sh
#!/bin/bash
# Row 243 V3 S5 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter404/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e, no leg
# chained to another with &&); every leg is tee'd to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips
# verify_go (the executor sets it; the controller runs it). Never start this under nohup. Writes
# <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
B=18c3d13231b39c7e2b0b437a964bce660d6739ec
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_REQUIRE_SSH_KEYGEN WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.txt, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.txt"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
need() { # need <leg> <test id>...: each must show --- PASS in <leg>.txt
  local lg=$1 t; shift
  for t in "$@"; do grep -qE -- "^ *--- PASS: $t " "$L/$lg.txt" || { echo "NOT PASSED ($lg): $t" | tee -a "$L/legs.txt"; FAILS="$FAILS ${lg}_missing"; }; done
}
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
bar() { echo "$*" | tr ' ' '|'; }
# The touched Go packages, derived from the diff (printed so the controller can diff it). The plan's list is
# cmd/ailang-worldd host/daemon host/store host/surfacepage host/verbs; a difference FAILS.
TOUCHED=$(git diff --name-only "$B" HEAD -- cmd host 2>/dev/null | grep '\.go$' | xargs -n1 dirname 2>/dev/null | sort -u | tr '\n' ' ')
echo "TOUCHED_PACKAGES: $TOUCHED" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ "$TOUCHED" != "cmd/ailang-worldd host/daemon host/store host/surfacepage host/verbs " ]; then FAILS="$FAILS touched_packages_differ"; fi
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
# the five touched packages whole, each its own rc (one red never skips another)
leg pkgs "r=0; for p in ./host/verbs/ ./host/store/ ./host/surfacepage/ ./cmd/ailang-worldd/ ./host/daemon/; do go test -count=1 -timeout 30m \$p; r=\$(( r | \$? )); done; exit \$r"
# -race on every touched package: host/verbs, host/store, host/surfacepage whole; cmd/ailang-worldd on the finish,
# approve and surface tests; host/daemon on the approval-surface and verb-route tests (both whole under -race exceed
# the leg's budget at load; CI's race job runs them whole)
CRACE='^(TestFinishConfirmOnSurface|TestCanonicalFinishArgv|TestFinishVerbRequestFields|TestSurfaceCardLinks|TestPutFinishVerbRequestRefusals|TestAwaitVerbConfirmationEnds|TestPollVerbRequestRechecks|TestReadBackFinishDecision|TestEpisodeFinish|TestApproveVerbRequest|TestSurfaceSSHSIGThenAwaitMints|TestQuickstartSection)'
DRACE='^(TestConfirm|TestConsumeRefusal|TestRequestsJSON|TestApprovePOST|TestApprovalRoutes|TestDecidePOST|TestDecideConcurrent|TestDecideRefuses|TestDecisionRoutes|TestDaemonCredentialWriterSites|TestSurfaceVerb)'
leg race "r=0; go test -race -count=1 ./host/verbs/ ./host/store/ ./host/surfacepage/; r=\$(( r | \$? )); go test -race -count=1 -v ./cmd/ailang-worldd/ -run '$CRACE' > '$L/race_cmd.raw' 2>&1; r=\$(( r | \$? )); grep -E '^(--- |ok|FAIL)' '$L/race_cmd.raw'; go test -race -count=1 -v ./host/daemon/ -run '$DRACE' > '$L/race_daemon.raw' 2>&1; r=\$(( r | \$? )); grep -E '^(--- |ok|FAIL)' '$L/race_daemon.raw'; exit \$r"
# the CLI's wait under -race -count=3 (directive 8, 12): polling, the fake clock, SQLite reads under a concurrent POST
NEW_WAIT="TestAwaitVerbConfirmationEnds TestPollVerbRequestRechecks TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision TestFinishConfirmOnSurfaceOpensThePR TestFinishConfirmOnSurfaceRerunsAndTypedPrompt TestFinishConfirmOnSurfaceKeepsTheCommit TestFinishConfirmOnSurfaceDispatchesTheNamedDecision TestFinishConfirmOnSurfaceResumeReusesTheApproval"
if [ "$MODE" = final ]; then leg race3 "r=0; go test -race -count=3 ./cmd/ailang-worldd/ -run '^($(bar $NEW_WAIT))\$'; r=\$(( r | \$? )); go test -race -count=3 ./host/daemon/ -run '^TestConfirmPOSTRecordsAtItsClock\$'; r=\$(( r | \$? )); exit \$r"; fi
# node: every page harness, non-skippable (WORLD_EXEC_NODE is set above; a harness that cannot run FAILS)
NODE="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceAppVerbs TestNeedsYouListsDecisions TestAddressHintsArePinned TestConfirmPageFailsClosed TestConfirmPageRendersText TestConfirmPageMessage TestConfirmPageStates TestConfirmPageSubmitsAPasskey TestConfirmPageShowsRefusals TestConfirmPageSubmitsAnSSHSIG"
leg node "set -o pipefail; go test -v -count=1 ./host/surfacepage/ -run '^($(bar $NODE))\$'"
need node $NODE
if grep -E '^\s*--- (SKIP|FAIL)' "$L/node.txt"; then FAILS="$FAILS node_skip_or_fail"; fi
NEW_CMD="TestCanonicalFinishArgv TestFinishConfirmOnSurfaceRefusesBeforeTheStore TestFinishVerbRequestFields TestSurfaceCardLinks TestPutFinishVerbRequestRefusals TestAwaitVerbConfirmationEnds TestPollVerbRequestRechecks TestReadBackFinishDecision TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision TestFinishConfirmOnSurfaceOpensThePR TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck TestFinishConfirmOnSurfaceRerunsAndTypedPrompt TestFinishConfirmOnSurfaceKeepsTheCommit TestFinishConfirmOnSurfaceDispatchesTheNamedDecision TestFinishConfirmOnSurfaceResumeReusesTheApproval TestApproveVerbRequestPointsAtItsCard TestQuickstartSection20MatchesTheFinishFlag"
NEW_D="TestConfirmPOSTRecordsAtItsClock"
if [ "$MODE" = base ]; then NEW_CMD=""; NEW_D=""; fi
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestConfirmRoutesRegistered TestDaemonCredentialWriterSites TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes TestDecidePageAssets TestConfirmPageAssetsServed TestConfirmationCardServed TestRequestsJSONConfirmations TestDecideRefusesReservedScopesByKey TestConfirmPOSTConfirms TestConfirmPOSTRetriesAFailedConsume"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsMethodsRequireADeadline TestCredentialsSchema9VerbRequests TestCredentialsNewerThanNineIsRefused TestPutVerbRequestSupersedes TestPutVerbRequestGrammarCapPrune TestVerbRequestExpiryBoundary"
INV_SP="TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestConfirmPageFieldTable TestConfirmCardGolden TestConfirmPageAssets"
INV_VERBS="TestLookup TestApproveTarget"
INV_AUTH="TestJudgeVerbSSH TestJudgeVerbAssertion TestRecheckVerb"
INV_BROKER="TestFinishPromptText TestFinishApprovalScopeStrict TestFinishConfirmRefusal"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestFinishPromptIsTheGolden TestEpisodeFinishHelp TestQuickstartSection19MatchesTheBar"
leg inv "set -o pipefail; r=0
go test -v -count=1 ./host/daemon/ -run '^($(bar $INV_DAEMON)$(alt $NEW_D))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/boundary/ -run '^($(bar $INV_BOUNDARY))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/store/ -run '^($(bar $INV_STORE))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/surfacepage/ -run '^($(bar $INV_SP))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/verbs/ -run '^($(bar $INV_VERBS))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/authority/ -run '^($(bar $INV_AUTH))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/broker/ -run '^($(bar $INV_BROKER))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $INV_CMD)$(alt $NEW_CMD))\$'; r=\$(( r | \$? ))
exit \$r"
need inv $INV_DAEMON $NEW_D $INV_BOUNDARY $INV_STORE $INV_SP $INV_VERBS $INV_AUTH $INV_BROKER $INV_CMD $NEW_CMD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# No new daemon route: every mux registration in daemon.go, counted (base = final), and the /surface POSTs (3)
ROUTES=$(grep -cE 'mux\.(HandleFunc|Handle)\(' host/daemon/daemon.go); POSTS=$(grep -cE 'mux\.HandleFunc\("POST /surface/' host/daemon/daemon.go)
BROUTES=$(git show "$B:host/daemon/daemon.go" | grep -cE 'mux\.(HandleFunc|Handle)\(')
echo "DAEMON_ROUTES: $ROUTES (base $BROUTES) SURFACE_POSTS: $POSTS" | tee -a "$L/legs.txt"
if [ "$ROUTES" != "$BROUTES" ] || [ "$POSTS" != 3 ]; then FAILS="$FAILS routes_changed"; fi
# No schema change, no .ail, the CSP untouched: no file under host/store but its context-roots test, no .ail, surface_page.go unchanged
NOPROD=$(git diff --name-only "$B" HEAD -- host/store world packages scripts tools host/daemon/surface_page.go host/daemon/daemon.go '*.ail' | grep -v '^host/store/context_roots_test.go$' | tr '\n' ' ')
echo "UNTOUCHED_VIOLATIONS: ${NOPROD:-none}" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ -n "$NOPROD" ]; then FAILS="$FAILS untouched_violated"; fi
# live Chrome: S5 changes app.js's navigate allowlist; V2b's fourteen, XC1 and the two drills stay green
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce TestVerbBarEpisodeReportInChrome TestVerbBarWhyInChrome TestConfirmPageShowsTheFinishPromptInChrome"
XD="TestSurfaceLiveChrome TestSurfaceFollowChrome"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; r=0; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(bar $XB))\$' -count=1 -v; r=\$(( r | \$? )); WORLD_CHROME=$CH go test ./host/daemon/ -run '^($(bar $XD))\$' -count=1 -v; r=\$(( r | \$? )); exit \$r"
need chrome $XB $XD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/chrome.txt"; then FAILS="$FAILS chrome_skip_or_fail"; fi
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then echo "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)" | tee -a "$L/legs.txt"
else leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'; fi
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
````

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, rm, mv, commit, stash, checkout, reset, branch). **`patch -p1 -E` (§2) is the
  only writer** of the tree; no hand edit, no `gofmt -w`, no golden regeneration (the golden ships in M13's patch), no
  fuzz. A failed `--dry-run` or a red acceptance → **stop and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance
  part is one call (C14 measured each).
- **Runs**: **Run 1 = M12a → M12b → M12c**, **Run 2 = M13 → M14** (each patch assumes the previous tree; the controller
  may commit between runs). Per sub-milestone the parts, in order: **M12a** `fast`; **M12b** `fast`, `daemon`; **M12c**
  `fast`, `cmd`; **M13** `fast`, `daemon`; **M14** `fast`, `docs`, `cmd`. `accept.sh` exports
  `AILANG_BIN=$HOME/.pinned-ailang/ailang`, `WORLD_EXEC_NODE=/opt/homebrew/bin/node` and
  `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang`, and unsets `WORLD_CHROME`,
  `WORLD_REQUIRE_PASSKEY_BROWSER` and every golden-update variable (no part needs Chrome; the gates' `chrome` leg is the
  controller's).
- **Commit-time test commands** (`commit_ms.sh <M> <msgfile> <test-cmd…>`; each carries its env through `accept.sh` and
  fits one 600 s call — the `cmd`, `daemon` and `docs` parts ran at acceptance and are not repeated at commit):
  `bash $T/accept.sh $W <M> fast` for **M12a, M12b, M12c, M13, M14**, with
  `T=$HOME/.ailang/state/world-iter404/exec-patches W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter404-row243-v3s5`
  (`commit_ms.sh` runs the command in the worktree after copying the snapshot; `accept.sh` takes the tree as its first
  argument). `commit_ms.sh` vets the directories of the changed `.go` files only; `app.js` and the golden (M13) are not
  `.go` — M13's `fast` part runs the node harness over the new `app.js` and the daemon's golden test over the new golden.
- After each sub-milestone's acceptance, snapshot **cumulatively against the base** to
  `$HOME/.ailang/state/world-iter404/snap/<M>/`: every file added or changed since `18c3d132` (repo-relative paths), a
  **`DELETED`** file (one repo path per line; **empty for every S5 stage**), a **`STATUS`** file, and `MANIFEST.sha256`
  (the two plan files, `STATUS` and `DELETED` are not in it). **Never an `rsync` of whole directories.** The `find` takes
  only the top-level dirs that exist (M12a–M13 have `cmd/` and `host/` but no `docs/`; M14 adds `docs/`), rehearsed on all
  five (C15):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M12a   # the sub-milestone just accepted
  B=18c3d13231b39c7e2b0b437a964bce660d6739ec
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter404-row243-v3s5
  D="$HOME/.ailang/state/world-iter404/snap/$M"; rm -rf "${D:?}"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M12a** 3 / 0 / 3; **M12b** 4 / 0 / 4; **M12c** 6 / 0 / 6; **M13**
  14 / 0 / 14; **M14** 16 / 0 / 16. The controller copies each snapshot over the worktree, commits it as one commit, and
  checks the numstat (`git diff -M --numstat 18c3d132 HEAD -- cmd host docs | awk '{s+=$1+$2} END {print s}'`): **M12a
  540, M12b 815, M12c 1,277, M13 1,336, M14 1,398**. A different number → stop.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`);
  gofmt output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal the
  cumulative list of §3's paragraphs); a numstat above 1,400.
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone
  (`TestWorkspaceExecConfiguredBindsTheEpisodeHandler`, `TestAGUISlowReaderCutIsResumable` in the `daemon` parts;
  `TestCLIRealSubprocessEpisode`, `TestCapsuleTemplateIsCopiedNotShared` in the `cmd` parts).
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty` or a real key; every store is a
  `t.TempDir()` file.

---

## 8. Risks

- **R1 — the CLI becomes a credentials writer.** `putFinishVerbRequest` opens the store's non-waiting writer for one
  `PutVerbRequest` (3 s root); a busy store refuses with the commit kept (CF8). A long-held writer elsewhere blocks only
  surface mode; the typed path is unchanged. No inventory pins CLI writer sites (§1.4 F8).
- **R2 — the wait trusts the credentials store's enrolled keys** (design §2.7 "stated trust"): re-check (a) is independent
  of the daemon; the passkey counter's freshness is the daemon's (C11/C15). A daemon that records a confirmation no key gave
  fails (a) (CF6, T9b).
- **R3 — the rig's daemon reads the wall clock** (§1.4 F3): the rig's fake clock restarts at the wall clock each run, and
  one test waits ≤ 1 s for the next second. Under heavy load a run straddling a second boundary is still correct (the
  request's `created_at ≤` the daemon's `now`).
- **R4 — the 2-line margin.** Any executor-side change would breach 1,400; the executor applies the patches only.
- **R5 — host/daemon and cmd under load.** The `daemon` and `cmd` parts run whole packages; a red known-flaky test is
  re-run once alone (§7).

---

## 9. Acceptance criteria (the judge scores against these)

1. **No new daemon route, no schema change, no `.ail`, `surfaceCSP` unchanged**: `DAEMON_ROUTES` base = final,
   `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none`, `verify_ail` identical to the base's, `TestDecidePageAssets` (pins the
   CSP literal) and `TestFrozenV1RouteTableMatchesMux`/`TestConfirmRoutesRegistered`/`TestDaemonCredentialWriterSites` green
   unedited.
2. **The flag is absent until M12c**: M12a and M12b register no flag and call nothing new (`episode_finish.go` unchanged
   until M12c); from M12c the help lists it (`TestQuickstartSection20MatchesTheFinishFlag` reads the help).
3. **The tty is never read in surface mode**: CF3 (`countingTerminal.reads == 0` with a terminal that would answer y),
   CF2, CF12; `H-TTY-READ` killed.
4. **Every refusal keeps the commit and dispatches nothing**: the rig's `kept()` (exit 1, "commit kept at <HEAD>", the
   hand steps, HEAD ≠ base, 0 gh calls, 0 records, 0 typed approve phases) on CF4 (rig), CF5, CF6 (rig); the unit refusals
   leave no row (CF9) and print nothing.
5. **Refusal 2 is before the check loop**: CF2 sees 0 outcomes with a terminal that would say y; `R2-AFTER` killed.
6. **The wait is bounded and its boundary exact**: CF4 (20 s wall deadline; 12 waits; `+10` seen, `+11` not);
   `W-UNBOUNDED-2`, `W-GRACE`, `W-GRACE-CONST`, `J-W-GE-GT`, `J-W-GE-EQ` killed.
7. **The re-check runs over this process's own digest, challenge and expiry**: CF6's eleven rows; `Q-EQ-*`, `Q-RECHECK`,
   `Q-ROW-VALUES-NOEQ`, `Q-APPROVEDBY`, `Q-NOW`, `Q-RECHECK-EXPIRY` killed; `Q-ROW-VALUES` named equivalent (S5-D8).
8. **The read-back requires every field**: CF7's eight rows; every `B-*` killed but four named equivalents (`B-DEC-FOUND`,
   `B-REQ-FOUND`, `B-JSON`, `B-DECODE`: a missing or undecodable object is refused by the same or the next condition, §4).
9. **CX21's CLI half**: the CLI dispatches with the row's decision, never the orphan (`DispatchesTheNamedDecision`).
10. **Judge-403 note 1**: `J-NOW-C14` killed by `TestConfirmPOSTRecordsAtItsClock` and by CF3; its survival against S4's
    tests alone recorded (C19).
11. **Every §2.9 P1/P6 row is driven by a test** (§3.1), incl. notes 4 (CF12 + S5-D9) and 5 (S5-D10).
12. **Exactly these pre-existing test files change**: `host/store/context_roots_test.go` (two roots), `host/verbs/lookup_test.go`
    (CB1/CB2), `host/surfacepage/testdata/verb-route.golden.json` (regenerated: one row added, one reason changed, C9),
    `host/daemon/surface_verb_outcomes_test.go` and `host/surfacepage/app_verb_test.go` (F6). **None deleted; no
    `episode_finish*_test.go`.**
13. **M13**: the finish reason is the design's text and `episode start`'s is unchanged (CB1); `approve vr-<16 hex>` navigates
    to `/surface/confirm/<id>` and nothing else does (CB2); the allowlist is exactly `/surface/confirm/vr-` (CB3,
    `CB-JS-ANY` killed); the CLI arm's text (CB4).
14. **CQ1**: §20 exists after §19, its command parses and carries the flag, and its quoted lines and numbers are the code's.
15. **Size**: `git diff -M --numstat 18c3d132 HEAD -- cmd host docs` sums to **1,398** (≤ 1,400; the same unfiltered); each
    sub-milestone's numstat as §7.
16. **Mutants**: every `K` row killed by its named test and by `suite-nostatic`; every `CTL` green; the `E` rows named with
    reasons; the plan's ids are the ledger's.
17. **Gates**: §6's script `FAILS: none` (or only a directive-listed load-sensitive test, green alone), `verify_go` rc 0.

---

## Appendix C — after this sprint (S5's close; what S6 inherits)

- **Moved out: nothing** (S5-D1: 1,398 ≤ 1,400; the design's cut M13 → S6 measured and not taken). No patch is parked.
- **Closed here**: S4 §1.4 F4 (CX21's CLI half); judge-403 note 1 (`J-NOW-C14`, two pins); S1 Appendix C note 4 (CF12,
  S5-D9) and note 5 / gemini r2 (S5-D10).
- **Residuals carried**: **judge-403 notes 2–4** (not taken, S5-D12): (2) `confirm.js` checks `confirmed:true` without
  status 200; (3) the busy 503's reason text is unasserted; (4) a second POST route would be caught only by the static
  route-table test. **New**: §1.4 F8 — no inventory pins the CLI's credentials-writer sites (the CLI now has
  `putFinishVerbRequest` beside `createMintRequest`'s and the approver verbs'); a `TestCLICredentialWriterSites` mirroring
  the daemon's AST walk is a candidate. S4's residuals unchanged: judge-402 note 2 (a day-expired request of the same
  episode is superseded, not pruned); S3 §1.4 F1 (the DDL accepts a superseded row with a `decision_ref`). Design RV1–RV9
  unchanged (RV1 Ctrl-C: the request stays pending until expiry or the next run's supersession, CF10).
- **What S6 inherits from S5**: the real CLI path — `ailang-worldd episode finish <ep> --preset <p> --confirm-on-surface`
  stores `vr-` + 16 hex, prints `confirm request <id> on the surface (expires <RFC3339>, message hash <16 hex>):`, one
  `  open <origin>/surface/confirm/<id>` per distinct passkey origin and the preset daemon's loopback, `the typed y is not
  read in this mode; …`, then waits (≤ 610 s) and prints `confirmed by <approved_by> (decision <ref>)` and the PR URL. **XC3**
  can drive it as `newSurfaceRig` does (the finish fixture + `daemon.New`), replacing the SSHSIG POST with the virtual
  authenticator on `/surface/confirm/<id>` in live Chrome, and asserting the fake gh's URL; its clock must start at the
  wall clock (§1.4 F3). **XC2** needs only a stored request (`storeOwn`'s form). **CP7** (*Needs you*) reads
  `requests.json`'s `confirmations` (S3) and links `/surface/confirm/<id>`, the prefix the page's allowlist now accepts
  (M13). The bar still answers `terminal` for `episode finish` (V4 runs it).

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M12a — the unwired production file, its first tests, the two context roots

````diff M12a
diff --git a/cmd/ailang-worldd/episode_finish_surface.go b/cmd/ailang-worldd/episode_finish_surface.go
new file mode 100644
index 000000000..9a206d94c
--- /dev/null
+++ b/cmd/ailang-worldd/episode_finish_surface.go
@@ -0,0 +1,283 @@
+package main
+
+// Row 243 V3 sprint 5 M12 (design §2.7): `episode finish … --confirm-on-surface`. After the commit,
+// where the typed y sits, finish stores a world/verb-request/v1 (credentials schema 9), prints its
+// card's links and waits — polling the credentials store, never reading the terminal — for a key's
+// confirmation; before any push it re-verifies the artifact over its own digest, challenge and
+// expiry (a) and reads back the world decision (b). Every refusal keeps the commit.
+
+import (
+	"context"
+	"crypto/sha256"
+	"encoding/hex"
+	"encoding/json"
+	"flag"
+	"fmt"
+	"io"
+	"net"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// verbLandGrace is how long past a request's expiry the wait still polls: above the POST's 3 s
+// commitBudget, so a confirmation admitted just before the expiry is seen (design §2.7 step 5).
+const verbLandGrace = 10 * time.Second
+
+// surfaceFinish is what finishConfirmOnSurface needs from the finish it confirms: the request's
+// fields, the credentials store it is stored in, and the preset daemon's address for the links.
+type surfaceFinish struct {
+	credentials, addr, ep, sha string
+	argv                       []string
+	scope                      broker.FinishApprovalScope
+	scopeText, title, body     string
+	committer                  string
+	protected                  []string
+}
+
+// verbOwn is this process's own request: its id and episode, the SHA-256 of the bytes it stored,
+// the challenge it drew and the expiry it set (design §2.7 step 6(a) re-checks against these).
+type verbOwn struct {
+	id, ep            string
+	digest, challenge [32]byte
+	expires           int64
+}
+
+// canonicalFinishArgv is the request's argv (pure): ailang-worldd episode finish <ep>, then each flag
+// the user set in fs.Visit's lexical order — --name for a true bool, --name value otherwise, and
+// --allow-failing once per value in the order given — so interleaved spellings of one command give
+// one argv, and a flag left at its default never appears (a bool set false is its default).
+func canonicalFinishArgv(fs *flag.FlagSet, allowed []string, ep string) []string {
+	argv := []string{"ailang-worldd", "episode", "finish", ep}
+	fs.Visit(func(f *flag.Flag) {
+		b, isBool := f.Value.(interface{ IsBoolFlag() bool })
+		switch {
+		case f.Name == "allow-failing":
+			for _, v := range allowed {
+				argv = append(argv, "--allow-failing", v)
+			}
+		case isBool && b.IsBoolFlag():
+			if f.Value.String() == "true" {
+				argv = append(argv, "--"+f.Name)
+			}
+		default:
+			argv = append(argv, "--"+f.Name, f.Value.String())
+		}
+	})
+	return argv
+}
+
+// finishVerbRequest is design §2.7 step 1 (pure): the request, validated, and its canonical bytes.
+// The prompt is broker.FinishPromptText of the same fields finishPrompt renders.
+func finishVerbRequest(in surfaceFinish) (approval.VerbRequest, []byte, error) {
+	protected := append([]string{}, in.protected...)
+	r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: in.argv, Episode: in.ep,
+		Prompt: broker.FinishPromptText(in.scope, in.title, in.body, protected, in.committer), Scope: in.scopeText,
+		Title: in.title, Body: in.body, Committer: in.committer, Protected: protected}
+	if err := r.Validate(); err != nil {
+		return r, nil, err
+	}
+	canon, err := approval.CanonicalVerb(r)
+	return r, canon, err
+}
+
+// finishConfirmOnSurface is design §2.7 steps 1-6. It returns the world decision a verified key gave
+// this finish's request, or the refusal; the caller keeps the commit on a refusal.
+func finishConfirmOnSurface(ctx context.Context, c *client, env sessionEnv, in surfaceFinish, tty, stdout io.Writer) (string, string) {
+	r, canon, err := finishVerbRequest(in)
+	if err != nil {
+		return "", "this finish cannot be confirmed on the surface (" + err.Error() + "); confirm at the terminal"
+	}
+	if env.random == nil {
+		return "", "internal error: no randomness source (sessionEnv.random)"
+	}
+	var idb [8]byte
+	var ch [32]byte
+	if err := env.random(idb[:]); err != nil {
+		return "", "randomness: " + err.Error()
+	}
+	if err := env.random(ch[:]); err != nil {
+		return "", "randomness: " + err.Error()
+	}
+	now := env.now()
+	own := verbOwn{id: "vr-" + hex.EncodeToString(idb[:]), ep: in.ep, digest: sha256.Sum256(canon), challenge: ch,
+		expires: now + int64(approval.VerbRequestLifetime/time.Second)}
+	links, err := surfaceCardLinks(ctx, in.credentials, "/surface/confirm/"+own.id, in.addr)
+	if err != nil {
+		return "", fmt.Sprintf("cannot read the credentials store %s: %v", in.credentials, err)
+	}
+	row := store.VerbRequestRow{RequestID: own.id, Request: canon, Digest: hex.EncodeToString(own.digest[:]),
+		Challenge: hex.EncodeToString(own.challenge[:]), EpisodeID: in.ep, CreatedAt: now, ExpiresAt: own.expires}
+	if refusal := putFinishVerbRequest(in.credentials, row, now); refusal != "" {
+		return "", refusal
+	}
+	h, err := approval.VerbMessageHash(own.digest, own.challenge, own.expires)
+	if err != nil {
+		return "", fmt.Sprintf("request %s: %v", own.id, err)
+	}
+	fmt.Fprint(tty, r.Prompt) // the terminal shows what the card shows; nothing is read from it
+	fmt.Fprintf(stdout, "confirm request %s on the surface (expires %s, message hash %s):\n", own.id, approverTime(own.expires), hex.EncodeToString(h[:8]))
+	for _, link := range links {
+		fmt.Fprintf(stdout, "  open %s\n", link)
+	}
+	fmt.Fprintf(stdout, "the typed y is not read in this mode; Ctrl-C stops waiting and keeps the commit at %s\n", in.sha)
+	got, refusal := awaitVerbConfirmation(env, in.credentials, own)
+	if refusal != "" {
+		return "", refusal
+	}
+	if refusal := readBackFinishDecision(worldReader{ctx: ctx, c: c}, got.DecisionRef, got.ApprovedBy, in.scopeText); refusal != "" {
+		return "", refusal
+	}
+	fmt.Fprintf(stdout, "confirmed by %s (decision %s)\n", got.ApprovedBy, got.DecisionRef)
+	return got.DecisionRef, ""
+}
+
+// surfaceCardLinks is writeAwaitLinks' address rule for a card path: one link per distinct active
+// passkey origin, then the preset daemon's loopback unless an origin already is it.
+func surfaceCardLinks(ctx context.Context, credPath, path, addr string) ([]string, error) {
+	ctx, cancel := context.WithTimeout(ctx, approveReadBudget)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		return nil, err
+	}
+	active, err := ro.ActivePasskeys(ctx)
+	_ = ro.Close()
+	if err != nil {
+		return nil, err
+	}
+	var links []string
+	seen := map[string]bool{}
+	for _, e := range active {
+		if !seen[e.Passkey.Origin] {
+			seen[e.Passkey.Origin] = true
+			links = append(links, e.Passkey.Origin+path)
+		}
+	}
+	if _, port, err := net.SplitHostPort(addr); err == nil && !seen["http://localhost:"+port] {
+		links = append(links, "http://localhost:"+port+path)
+	}
+	return links, nil
+}
+
+// putFinishVerbRequest is design §2.7 step 3: one non-waiting writer under its own 3 s root, the
+// put, close. It returns the refusal text, or "".
+func putFinishVerbRequest(credPath string, row store.VerbRequestRow, now int64) string {
+	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
+	defer cancel()
+	st, err := store.OpenCredentials(credPath)
+	if store.IsWriterAlreadyActive(err) {
+		return "the credentials store is busy (another writer holds it); run finish again"
+	}
+	if err != nil {
+		return fmt.Sprintf("cannot open the credentials store %s: %v", credPath, err)
+	}
+	defer func() { _ = st.Close() }()
+	if err := st.PutVerbRequest(ctx, row, now); err != nil {
+		return fmt.Sprintf("cannot store request %s: %v", row.RequestID, err)
+	}
+	return ""
+}
+
+// awaitVerbConfirmation is design §2.7 step 5: a poll every awaitPollInterval through env.after
+// until the row is approved and re-checked (returned), superseded, unreadable, or the wait reaches
+// expires_at + verbLandGrace. It never reads the terminal.
+func awaitVerbConfirmation(env sessionEnv, credPath string, own verbOwn) (store.VerbRequestRow, string) {
+	deadline := own.expires + int64(verbLandGrace/time.Second)
+	for {
+		now := env.now()
+		row, refusal, err := pollVerbRequest(credPath, own, now)
+		switch {
+		case err != nil:
+			return store.VerbRequestRow{}, fmt.Sprintf("cannot read request %s: %v", own.id, err)
+		case refusal != "":
+			return store.VerbRequestRow{}, refusal
+		case row.Outcome == store.VerbOutcomeApproved:
+			return row, ""
+		case row.Outcome == store.VerbOutcomeSuperseded:
+			return store.VerbRequestRow{}, fmt.Sprintf("request %s was superseded by a newer finish of %s", own.id, own.ep)
+		case now >= deadline:
+			return store.VerbRequestRow{}, fmt.Sprintf("no confirmation of request %s by %s", own.id, approverTime(now))
+		}
+		<-env.after(awaitPollInterval)
+	}
+}
+
+// pollVerbRequest is one poll (design §2.7 steps 5 and 6(a)): one read-only open under its own 3 s
+// root, the row, and for an approved row the re-check in the same open — the row's digest,
+// challenge and expiry must be this process's own, authority.RecheckVerb must verify its artifact
+// over this process's own values against the enrolled key, and that key's approver must be the
+// row's. It returns the row and a refusal ("" when none).
+func pollVerbRequest(credPath string, own verbOwn, now int64) (store.VerbRequestRow, string, error) {
+	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		return store.VerbRequestRow{}, "", err
+	}
+	defer func() { _ = ro.Close() }()
+	row, ok, err := ro.VerbRequest(ctx, own.id)
+	if err != nil {
+		return row, "", err
+	}
+	if !ok {
+		return row, "", fmt.Errorf("request %s is no longer in the store", own.id)
+	}
+	if row.Outcome != store.VerbOutcomeApproved {
+		return row, "", nil
+	}
+	bad := func(reason string) (store.VerbRequestRow, string, error) {
+		return row, fmt.Sprintf("the confirmation of request %s does not verify: %s", own.id, reason), nil
+	}
+	if row.Digest != hex.EncodeToString(own.digest[:]) || row.Challenge != hex.EncodeToString(own.challenge[:]) || row.ExpiresAt != own.expires {
+		return bad("the stored request is not the one this finish made")
+	}
+	in := authority.VerbJudgeInput{Digest: own.digest, Challenge: own.challenge, ExpiresAt: own.expires, Now: now}
+	signed, err := authority.RecheckVerb(ctx, ro, in, row.Artifact)
+	if err != nil {
+		return bad(err.Error())
+	}
+	if signed.ApprovedBy() != row.ApprovedBy {
+		return bad(fmt.Sprintf("the key is %s's, the row says %s", signed.ApprovedBy(), row.ApprovedBy))
+	}
+	return row, "", nil
+}
+
+// readBackFinishDecision is design §2.7 step 6(b): through worldReader.object (GET /v1/objects), the
+// decision the row names must approve, be decided by the row's approver, and its request must be
+// this finish scope's Human.Approve for operator:episode-finish. It returns the refusal, or "".
+func readBackFinishDecision(r worldReader, ref, approvedBy, scopeText string) string {
+	bad := func(why string) string {
+		return "the world decision " + ref + " does not approve this finish (F6): " + why
+	}
+	obj, ok, err := r.object(ref)
+	if err != nil || !ok || obj.SemanticID != broker.ApprovalDecisionV1 {
+		return bad(fmt.Sprintf("not a finish approval decision (found %v, %v)", ok, err))
+	}
+	var d struct {
+		RequestRef string `json:"requestRef"`
+		Decision   string `json:"decision"`
+		DecidedBy  string `json:"decidedBy"`
+	}
+	if err := json.Unmarshal(obj.Payload, &d); err != nil {
+		return bad(err.Error())
+	}
+	if d.Decision != "approve" || d.DecidedBy != approvedBy {
+		return bad(fmt.Sprintf("decision %q by %q", d.Decision, d.DecidedBy))
+	}
+	qobj, ok, err := r.object(d.RequestRef)
+	if err != nil || !ok || qobj.SemanticID != broker.ApprovalRequestV1 {
+		return bad(fmt.Sprintf("its request %s is not readable (found %v, %v)", d.RequestRef, ok, err))
+	}
+	q, err := broker.DecodeApprovalRequest(qobj.Payload)
+	if err != nil {
+		return bad(err.Error())
+	}
+	if q.Effect != broker.EffectHumanApprove || q.Requester != "operator:episode-finish" || q.Scope != scopeText {
+		return bad(fmt.Sprintf("it approves %s for %s over another scope", q.Effect, q.Requester))
+	}
+	return ""
+}
diff --git a/cmd/ailang-worldd/episode_finish_surface_test.go b/cmd/ailang-worldd/episode_finish_surface_test.go
new file mode 100644
index 000000000..313795bea
--- /dev/null
+++ b/cmd/ailang-worldd/episode_finish_surface_test.go
@@ -0,0 +1,253 @@
+package main
+
+// Row 243 V3 sprint 5 M12 (design §2.7, §2.9 P1, §6 CF4–CF9, CF14): the confirm-on-surface parts
+// one by one — the canonical argv, the request's refusals, the put's refusals, the wait's every end
+// (a fake clock, a wall-clock test deadline), the re-check (a) and the read-back (b). Every store is
+// under t.TempDir(); every key is approvaltest's (sigSigner enrolled as mark-mac, otherSigner never);
+// an "approved by hand" row is raw SQL, as the design says; world objects are broker's own writes.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"database/sql"
+	"encoding/hex"
+	"errors"
+	"flag"
+	"fmt"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+const sfT0 = int64(1_900_000_000)
+
+// surfaceUnit is a credentials store with sigSigner enrolled as mark-mac, and a real finish input.
+func surfaceUnit(t *testing.T) (string, surfaceFinish) {
+	t.Helper()
+	creds := filepath.Join(t.TempDir(), "credentials.db")
+	seedApprover(t, creds, "mark-mac", sigSigner.PublicKey(), sfT0-100)
+	scope := broker.FinishApprovalScope{Repository: "fixture/repo", Remote: "https://github.com/fixture/repo.git", Head: "custom-branch",
+		SHA: strings.Repeat("a", 40), Base: "dev", TitleSHA256: finishSHA([]byte("Fixture title")), BodySHA256: finishSHA([]byte("body\n")),
+		Draft: true, ExpiresAt: sfT0 + 900}
+	return creds, surfaceFinish{credentials: creds, addr: "127.0.0.1:47330", ep: "ep1", sha: scope.SHA,
+		argv: []string{"ailang-worldd", "episode", "finish", "ep1", "--confirm-on-surface"}, scope: scope,
+		scopeText: broker.FormatFinishApprovalScope(scope), title: "Fixture title", body: "body\n"}
+}
+
+// storeOwn stores in's request as id with challenge bytes ch at now, through the production put.
+func storeOwn(t *testing.T, in surfaceFinish, id string, ch byte, now int64) verbOwn {
+	t.Helper()
+	_, canon, err := finishVerbRequest(in)
+	if err != nil {
+		t.Fatal(err)
+	}
+	own := verbOwn{id: id, ep: in.ep, digest: sha256.Sum256(canon), challenge: [32]byte(bytes.Repeat([]byte{ch}, 32)), expires: now + 600}
+	row := store.VerbRequestRow{RequestID: id, Request: canon, Digest: hex.EncodeToString(own.digest[:]), Challenge: hex.EncodeToString(own.challenge[:]),
+		EpisodeID: in.ep, CreatedAt: now, ExpiresAt: own.expires}
+	if refusal := putFinishVerbRequest(in.credentials, row, now); refusal != "" {
+		t.Fatal(refusal)
+	}
+	return own
+}
+
+// verbSig is s's SSHSIG under namespace over the 86-byte message of digest, challenge and expiry.
+func verbSig(t *testing.T, s *approvaltest.Signer, namespace string, digest, challenge [32]byte, expiry int64) []byte {
+	t.Helper()
+	m, err := approval.VerbMessage(digest, challenge, expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return s.Sign(namespace, m)
+}
+
+// credSQL runs raw statements on the credentials store (the design's "approved by hand").
+func credSQL(t *testing.T, creds string, stmts ...string) {
+	t.Helper()
+	db, err := sql.Open("sqlite", creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	for _, s := range stmts {
+		if _, err := db.Exec(s); err != nil {
+			t.Fatalf("%s: %v", s, err)
+		}
+	}
+}
+
+// approveByHand marks id approved with artifact, approvedBy and a well-formed decision ref; set
+// adds column overrides (raw SQL, as design §6 CF6).
+func approveByHand(t *testing.T, creds, id string, artifact []byte, approvedBy string, at int64, set ...string) {
+	t.Helper()
+	cols := append([]string{fmt.Sprintf("consumed_at = %d", at), "outcome = 'approved'", "approved_by = '" + approvedBy + "'",
+		"artifact = X'" + hex.EncodeToString(artifact) + "'", "decision_ref = 'sha256:" + strings.Repeat("d", 64) + "'"}, set...)
+	credSQL(t, creds, "UPDATE verb_requests SET "+strings.Join(cols, ", ")+" WHERE request_id = '"+id+"'")
+}
+
+// CF14: interleaved spellings of one command give one argv; a flag left at its default (or a bool
+// set false) never appears; --allow-failing once per value in the order given; bools bare.
+func TestCanonicalFinishArgv(t *testing.T) {
+	parse := func(args ...string) []string {
+		fs := flag.NewFlagSet("finish", flag.ContinueOnError)
+		fs.String("preset", "", "")
+		fs.String("presets-dir", "/default", "")
+		fs.String("title", "", "")
+		fs.Bool("ready", false, "")
+		fs.Bool("dry-run", false, "")
+		fs.Bool("confirm-on-surface", false, "")
+		var allowed multiFlag
+		fs.Var(&allowed, "allow-failing", "")
+		pos, err := verbs.ParseInterleaved(fs, args)
+		if err != nil || len(pos) != 1 {
+			t.Fatalf("%v: %v %v", args, pos, err)
+		}
+		return canonicalFinishArgv(fs, allowed, pos[0])
+	}
+	want := []string{"ailang-worldd", "episode", "finish", "ep1", "--allow-failing", "b", "--allow-failing", "a",
+		"--confirm-on-surface", "--preset", "p", "--ready", "--title", "two words"}
+	for _, args := range [][]string{
+		{"ep1", "--preset", "p", "--allow-failing", "b", "--ready", "--title", "two words", "--allow-failing", "a", "--confirm-on-surface"},
+		{"--confirm-on-surface", "--title=two words", "ep1", "--allow-failing=b", "--preset=p", "--allow-failing", "a", "--ready=true", "--dry-run=false"},
+	} {
+		if got := parse(args...); !reflect.DeepEqual(got, want) {
+			t.Fatalf("%q\n got %q\nwant %q", args, got, want)
+		}
+	}
+	if got := parse("ep1"); !reflect.DeepEqual(got, []string{"ailang-worldd", "episode", "finish", "ep1"}) {
+		t.Fatalf("defaults appear: %q", got)
+	}
+}
+
+// CF9: a request the law refuses (over 128 KiB), no randomness, failing randomness: each refuses
+// before anything is stored or shown.
+func TestFinishConfirmOnSurfaceRefusesBeforeTheStore(t *testing.T) {
+	creds, in := surfaceUnit(t)
+	big := in
+	for i := 0; len(strings.Join(big.protected, "")) < approval.MaxVerbRequestBytes; i++ {
+		big.protected = append(big.protected, fmt.Sprintf("p/%04d/%s", i, strings.Repeat("x", 200)))
+	}
+	missing := in
+	missing.credentials = filepath.Join(t.TempDir(), "missing", "credentials.db")
+	random := func(b []byte) error { return nil }
+	for _, c := range []struct {
+		name   string
+		in     surfaceFinish
+		random func([]byte) error
+		want   string
+	}{
+		{"over 128 KiB", big, random, "the request is"},
+		{"no randomness", in, nil, "internal error: no randomness source"},
+		{"the id's draw fails", in, func(b []byte) error { return map[bool]error{true: errors.New("entropy gone")}[len(b) == 8] }, "randomness: entropy gone"},
+		{"the challenge's draw fails", in, func(b []byte) error { return map[bool]error{true: errors.New("entropy gone")}[len(b) == 32] }, "randomness: entropy gone"},
+		{"no credentials store", missing, random, "cannot read the credentials store " + missing.credentials},
+		{"busy", in, random, "the credentials store is busy"},
+	} {
+		env, _ := awaitEnv(sfT0, new(bool), nil)
+		env.random = c.random
+		var tty, out bytes.Buffer
+		if c.name == "busy" {
+			st, err := store.OpenCredentials(creds)
+			if err != nil {
+				t.Fatal(err)
+			}
+			defer func() { _ = st.Close() }()
+		}
+		ref, refusal := finishConfirmOnSurface(boundedTestContext(t), nil, env, c.in, &tty, &out)
+		if ref != "" || !strings.Contains(refusal, c.want) || tty.Len() != 0 || out.Len() != 0 {
+			t.Fatalf("%s: %q %q tty %q out %q", c.name, ref, refusal, tty.String(), out.String())
+		}
+		if c.name == "over 128 KiB" && !strings.HasSuffix(refusal, "; confirm at the terminal") {
+			t.Fatalf("%s: %q", c.name, refusal)
+		}
+		if n := pendingVerbRows(t, creds); n != 0 {
+			t.Fatalf("%s: %d rows stored", c.name, n)
+		}
+	}
+}
+
+// The request's every field is finish's: the prompt is FinishPromptText of the same fields
+// (committer and protected paths included), protected an array even when empty.
+func TestFinishVerbRequestFields(t *testing.T) {
+	_, in := surfaceUnit(t)
+	for _, c := range []struct {
+		committer string
+		protected []string
+	}{{"Mark <m@example.invalid>", []string{"go.mod"}}, {"", nil}} {
+		in.committer, in.protected = c.committer, c.protected
+		r, canon, err := finishVerbRequest(in)
+		want := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: in.argv, Episode: "ep1",
+			Prompt: broker.FinishPromptText(in.scope, in.title, in.body, c.protected, c.committer), Scope: in.scopeText, Title: in.title, Body: in.body,
+			Committer: c.committer, Protected: append([]string{}, c.protected...)}
+		if back, derr := approval.DecodeVerb(canon); err != nil || derr != nil || !reflect.DeepEqual(r, want) || !reflect.DeepEqual(back, want) {
+			t.Fatalf("%+v (%v %v)\nwant %+v", r, err, derr, want)
+		}
+	}
+}
+
+// writeAwaitLinks' address rule: one link per distinct passkey origin, then the daemon's loopback
+// unless an origin already is it.
+func TestSurfaceCardLinks(t *testing.T) {
+	creds, _ := surfaceUnit(t)
+	if got, err := surfaceCardLinks(boundedTestContext(t), creds, "/p", "127.0.0.1:7644"); err != nil || !reflect.DeepEqual(got, []string{"http://localhost:7644/p"}) {
+		t.Fatalf("no passkey: %q %v", got, err)
+	}
+	seedPasskey(t, creds, "pk-a", approvaltest.NewPasskey(approval.AlgEdDSA, 31), sfT0)
+	seedPasskey(t, creds, "pk-b", approvaltest.NewPasskey(approval.AlgEdDSA, 32), sfT0)
+	for addr, want := range map[string][]string{"127.0.0.1:7644": {pkOrigin + "/p"}, "127.0.0.1:47330": {pkOrigin + "/p", "http://localhost:47330/p"}, "no-port": {pkOrigin + "/p"}} {
+		if got, err := surfaceCardLinks(boundedTestContext(t), creds, "/p", addr); err != nil || !reflect.DeepEqual(got, want) {
+			t.Fatalf("%s: %q %v", addr, got, err)
+		}
+	}
+}
+
+// pendingVerbRows counts creds' verb_requests rows.
+func pendingVerbRows(t *testing.T, creds string) int {
+	t.Helper()
+	db, err := sql.Open("sqlite", creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	var n int
+	if err := db.QueryRow(`SELECT count(*) FROM verb_requests`).Scan(&n); err != nil {
+		t.Fatal(err)
+	}
+	return n
+}
+
+// CF8: the put's refusals — a busy writer, a store below schema 9 (naming the restart), 16 pending.
+func TestPutFinishVerbRequestRefusals(t *testing.T) {
+	creds, in := surfaceUnit(t)
+	st, err := store.OpenCredentials(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	_, canon, _ := finishVerbRequest(in)
+	row := store.VerbRequestRow{RequestID: "vr-00000000000000f1", Request: canon, Digest: finishSHA(canon), Challenge: strings.Repeat("c", 64),
+		EpisodeID: "ep1", CreatedAt: sfT0, ExpiresAt: sfT0 + 600}
+	if got := putFinishVerbRequest(filepath.Join(t.TempDir(), "missing", "c.db"), row, sfT0); !strings.HasPrefix(got, "cannot open the credentials store ") {
+		t.Fatalf("no store: %q", got)
+	}
+	if got := putFinishVerbRequest(creds, row, sfT0); got != "the credentials store is busy (another writer holds it); run finish again" {
+		t.Fatalf("busy: %q", got)
+	}
+	_ = st.Close()
+	for i := 0; i < store.MaxPendingVerbRequests; i++ {
+		storeOwn(t, surfaceFinish{credentials: creds, addr: in.addr, ep: fmt.Sprintf("ep-%02d", i), argv: []string{"ailang-worldd", "episode", "finish", fmt.Sprintf("ep-%02d", i)}, scope: in.scope,
+			scopeText: in.scopeText, title: in.title, body: in.body}, fmt.Sprintf("vr-00000000000001%02x", i), 1, sfT0)
+	}
+	if got := putFinishVerbRequest(creds, row, sfT0); got != "cannot store request vr-00000000000000f1: "+store.ErrTooManyPendingVerbRequests.Error() {
+		t.Fatalf("cap: %q", got)
+	}
+	credSQL(t, creds, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
+	if got := putFinishVerbRequest(creds, row, sfT0); !strings.Contains(got, "verb requests need schema 9") || !strings.Contains(got, "restart `ailang-worldd serve`") {
+		t.Fatalf("v8: %q", got)
+	}
+}
diff --git a/host/store/context_roots_test.go b/host/store/context_roots_test.go
index 77fc94b81..8eabd885a 100644
--- a/host/store/context_roots_test.go
+++ b/host/store/context_roots_test.go
@@ -28,6 +28,10 @@ var contextRootPins = map[string]int{
 	"cmd/ailang-worldd/approver_passkey.go|runApproverEnrollPasskey|Background": 3,
 	// Row 201 A3 sprint 1 M3: each approve --await poll's read-only reads and judgement (3 s).
 	"cmd/ailang-worldd/approve_await.go|pollSurfaceApproval|Background": 1,
+	// Row 243 V3 sprint 5 M12: episode finish --confirm-on-surface's put of its verb request (3 s)
+	// and each poll's read-only reads and re-check (3 s).
+	"cmd/ailang-worldd/episode_finish_surface.go|putFinishVerbRequest|Background": 1,
+	"cmd/ailang-worldd/episode_finish_surface.go|pollVerbRequest|Background":      1,
 	// Row 201 A4 sprint 2: approve <request-ref>'s daemon reads before the y/N, and its POST after it (10 s each).
 	"cmd/ailang-worldd/approve_decide.go|runApproveDecide|Background": 1,
 	"cmd/ailang-worldd/approve_decide.go|postDecision|Background":     1,
````

### Patch M12b — the wait, the re-check, the read-back; judge-403 note 1

````diff M12b
diff --git a/cmd/ailang-worldd/episode_finish_surface_test.go b/cmd/ailang-worldd/episode_finish_surface_test.go
index 313795bea..1355699c9 100644
--- a/cmd/ailang-worldd/episode_finish_surface_test.go
+++ b/cmd/ailang-worldd/episode_finish_surface_test.go
@@ -11,17 +11,21 @@ import (
 	"crypto/sha256"
 	"database/sql"
 	"encoding/hex"
+	"encoding/json"
 	"errors"
 	"flag"
 	"fmt"
+	"net/http"
 	"path/filepath"
 	"reflect"
 	"strings"
 	"testing"
+	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
 	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/hashref"
 	"github.com/sunholo-data/ailang-world/host/store"
 	"github.com/sunholo-data/ailang-world/host/verbs"
 )
@@ -251,3 +255,244 @@ func TestPutFinishVerbRequestRefusals(t *testing.T) {
 		t.Fatalf("v8: %q", got)
 	}
 }
+
+// CF4 (+ CF5's and the exact-boundary rows, every way the wait ends): under the fake clock and a
+// wall-clock test deadline, the wait ends at expires_at + verbLandGrace; a confirmation landing at
+// exactly that instant is seen and one a second later is not; a supersession, a vanished row.
+func TestAwaitVerbConfirmationEnds(t *testing.T) {
+	type end struct {
+		row     store.VerbRequestRow
+		refusal string
+		waits   int64
+	}
+	run := func(t *testing.T, creds string, own verbOwn, start int64, onWait func(n, now int64)) end {
+		t.Helper()
+		env, clock := awaitEnv(start, new(bool), onWait)
+		done := make(chan end, 1)
+		go func() {
+			row, refusal := awaitVerbConfirmation(env, creds, own)
+			done <- end{row, refusal, clock.waits}
+		}()
+		select {
+		case e := <-done:
+			return e
+		case <-time.After(20 * time.Second):
+			t.Fatal("the wait did not end: unbounded")
+		}
+		return end{}
+	}
+	t.Run("expiry+grace, bounded", func(t *testing.T) {
+		creds, in := surfaceUnit(t)
+		own := storeOwn(t, in, "vr-00000000000000a1", 0xa1, sfT0)
+		deadline := own.expires + 10
+		e := run(t, creds, own, own.expires-2, nil)
+		if e.refusal != "no confirmation of request vr-00000000000000a1 by "+approverTime(deadline) || e.waits != 12 {
+			t.Fatalf("%q after %d waits", e.refusal, e.waits)
+		}
+	})
+	for _, c := range []struct {
+		name string
+		at   int64 // seconds after expires_at that the confirmation lands
+		seen bool
+	}{{"lands at expires_at+grace", 10, true}, {"lands a second later", 11, false}} {
+		t.Run(c.name, func(t *testing.T) {
+			creds, in := surfaceUnit(t)
+			own := storeOwn(t, in, "vr-00000000000000a2", 0xa2, sfT0)
+			art := verbSig(t, sigSigner, approval.VerbNamespace, own.digest, own.challenge, own.expires)
+			e := run(t, creds, own, own.expires-2, func(_, now int64) {
+				if now == own.expires+c.at {
+					approveByHand(t, creds, own.id, art, "ssh:mark-mac", now)
+				}
+			})
+			if seen := e.refusal == "" && e.row.Outcome == store.VerbOutcomeApproved && e.row.ApprovedBy == "ssh:mark-mac"; seen != c.seen {
+				t.Fatalf("seen %v, want %v: %q", seen, c.seen, e.refusal)
+			}
+		})
+	}
+	t.Run("superseded", func(t *testing.T) {
+		creds, in := surfaceUnit(t)
+		own := storeOwn(t, in, "vr-00000000000000a3", 0xa3, sfT0)
+		e := run(t, creds, own, sfT0, func(n, now int64) {
+			if n == 1 {
+				storeOwn(t, in, "vr-00000000000000a4", 0xa4, now)
+			}
+		})
+		if e.refusal != "request vr-00000000000000a3 was superseded by a newer finish of ep1" || e.waits != 1 {
+			t.Fatalf("%q after %d waits", e.refusal, e.waits)
+		}
+	})
+	t.Run("a poll read fails", func(t *testing.T) {
+		creds, in := surfaceUnit(t)
+		own := storeOwn(t, in, "vr-00000000000000a5", 0xa5, sfT0)
+		e := run(t, creds, own, sfT0, func(n, _ int64) {
+			if n == 2 {
+				credSQL(t, creds, `DELETE FROM verb_requests`)
+			}
+		})
+		if e.refusal != "cannot read request vr-00000000000000a5: request vr-00000000000000a5 is no longer in the store" || e.waits != 2 {
+			t.Fatalf("%q after %d waits", e.refusal, e.waits)
+		}
+	})
+}
+
+// CF6 (design §2.7 step 6(a)): an approved row is re-checked against this process's own digest,
+// challenge and expiry and the enrolled key; each row marked approved by hand (raw SQL) that a key
+// did not give over these bytes refuses "does not verify" — including an approver revoked after the
+// confirmation (S1 Appendix C note 5: fail-safe, the commit is kept and a re-run asks afresh).
+func TestPollVerbRequestRechecks(t *testing.T) {
+	other := [32]byte(bytes.Repeat([]byte{0xee}, 32))
+	for _, c := range []struct {
+		name string
+		art  func(own verbOwn) []byte
+		by   string
+		set  []string
+		pre  func(t *testing.T, creds string)
+		want string
+	}{
+		{"control: the key over these bytes", nil, "ssh:mark-mac", nil, nil, ""},
+		{"the row's digest is not ours", nil, "ssh:mark-mac", []string{"digest = '" + hex.EncodeToString(other[:]) + "'"}, nil, "the stored request is not the one this finish made"},
+		{"the row's challenge is not ours", nil, "ssh:mark-mac", []string{"challenge = '" + hex.EncodeToString(other[:]) + "'"}, nil, "the stored request is not the one this finish made"},
+		{"the row's expiry is not ours", nil, "ssh:mark-mac", []string{"expires_at = expires_at + 1"}, nil, "the stored request is not the one this finish made"},
+		{"signed over another digest", func(o verbOwn) []byte {
+			return verbSig(t, sigSigner, approval.VerbNamespace, other, o.challenge, o.expires)
+		}, "ssh:mark-mac", nil, nil, "does not verify: "},
+		{"signed over another challenge", func(o verbOwn) []byte {
+			return verbSig(t, sigSigner, approval.VerbNamespace, o.digest, other, o.expires)
+		}, "ssh:mark-mac", nil, nil, "does not verify: "},
+		{"signed over another expiry", func(o verbOwn) []byte {
+			return verbSig(t, sigSigner, approval.VerbNamespace, o.digest, o.challenge, o.expires+1)
+		}, "ssh:mark-mac", nil, nil, "does not verify: "},
+		{"the mint namespace", func(o verbOwn) []byte {
+			return verbSig(t, sigSigner, approval.Namespace, o.digest, o.challenge, o.expires)
+		}, "ssh:mark-mac", nil, nil, "does not verify: "},
+		{"a key never enrolled", func(o verbOwn) []byte {
+			return verbSig(t, otherSigner, approval.VerbNamespace, o.digest, o.challenge, o.expires)
+		}, "ssh:mark-mac", nil, nil, "does not verify: "},
+		{"not an artifact", func(verbOwn) []byte { return []byte("y") }, "ssh:mark-mac", nil, nil, "does not verify: "},
+		{"the approver revoked since", nil, "ssh:mark-mac", nil, func(t *testing.T, creds string) { revokeSeeded(t, creds, "mark-mac", sfT0+5) }, "revoked"},
+		{"approved_by is not the key's", nil, "ssh:mark-other", nil, nil, "the key is ssh:mark-mac's, the row says ssh:mark-other"},
+	} {
+		t.Run(c.name, func(t *testing.T) {
+			creds, in := surfaceUnit(t)
+			own := storeOwn(t, in, "vr-00000000000000b1", 0xb1, sfT0)
+			art := verbSig(t, sigSigner, approval.VerbNamespace, own.digest, own.challenge, own.expires)
+			if c.art != nil {
+				art = c.art(own)
+			}
+			approveByHand(t, creds, own.id, art, c.by, sfT0+1, c.set...)
+			if c.pre != nil {
+				c.pre(t, creds)
+			}
+			row, refusal, err := pollVerbRequest(creds, own, sfT0+10)
+			if err != nil || (c.want == "") != (refusal == "") || !strings.Contains(refusal, c.want) {
+				t.Fatalf("%q %v, want %q", refusal, err, c.want)
+			}
+			if c.want != "" && !strings.HasPrefix(refusal, "the confirmation of request vr-00000000000000b1 does not verify: ") {
+				t.Fatalf("%q", refusal)
+			}
+			if c.want == "" && (row.Outcome != store.VerbOutcomeApproved || !bytes.Equal(row.Artifact, art)) {
+				t.Fatalf("control row %+v", row)
+			}
+		})
+	}
+}
+
+// CF7 (design §2.7 step 6(b)): the decision the row names is read back through worldReader.object
+// over GET /v1/objects; the decision and its request are broker's own writes (RecordOperatorApproval
+// on a test world store), each refused row the same object with one field changed.
+func TestReadBackFinishDecision(t *testing.T) {
+	_, in := surfaceUnit(t)
+	st, err := store.Open(filepath.Join(t.TempDir(), "world.db"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = st.Close() }()
+	ctx := boundedTestContext(t)
+	qref, dref, err := broker.RecordOperatorApproval(ctx, st, in.scopeText, "operator:episode-finish", "ssh:mark-mac", sfT0)
+	if err != nil {
+		t.Fatal(err)
+	}
+	real := map[string]wireObject{}
+	for _, ref := range []string{qref.String(), dref.String()} {
+		obj, ok, err := st.GetObject(ctx, hashref.MustParse(ref))
+		if err != nil || !ok {
+			t.Fatal(ref, ok, err)
+		}
+		real[ref] = wireObject{Hash: ref, SemanticID: obj.SemanticID, Payload: obj.Payload}
+	}
+	// edit is the real object at ref with one payload field (or the semantic id) changed.
+	edit := func(ref, field string, value any) map[string]wireObject {
+		objs := map[string]wireObject{}
+		for k, v := range real {
+			objs[k] = v
+		}
+		o := objs[ref]
+		if field == "semanticId" {
+			o.SemanticID = value.(string)
+		} else {
+			var m map[string]any
+			if err := json.Unmarshal(o.Payload, &m); err != nil {
+				t.Fatal(err)
+			}
+			m[field] = value
+			o.Payload, _ = json.Marshal(m)
+		}
+		objs[ref] = o
+		return objs
+	}
+	q, d := qref.String(), dref.String()
+	for _, c := range []struct {
+		name string
+		objs map[string]wireObject
+		want string
+	}{
+		{"control", real, ""},
+		{"a deny", edit(d, "decision", "deny"), `decision "deny"`},
+		{"decided by another", edit(d, "decidedBy", "ssh:mark-other"), `by "ssh:mark-other"`},
+		{"not a decision", edit(d, "semanticId", broker.ApprovalRequestV1), "not a finish approval decision"},
+		{"its request is missing", edit(d, "requestRef", "sha256:"+strings.Repeat("0", 64)), "is not readable"},
+		{"its request is not a request", edit(q, "semanticId", broker.ApprovalDecisionV1), "is not readable"},
+		{"another effect", edit(q, "effect", "Workspace.Write"), "over another scope"},
+		{"another requester", edit(q, "requester", "operator:terminal"), "over another scope"},
+		{"another scope", edit(q, "scope", strings.Replace(in.scopeText, "draft=true", "draft=false", 1)), "over another scope"},
+	} {
+		c := c
+		r := worldReader{ctx: ctx, c: memoryDaemon(func(w http.ResponseWriter, req *http.Request) {
+			o, ok := c.objs[strings.TrimPrefix(req.URL.Path, "/v1/objects/")]
+			if !ok || req.URL.Query().Get("payload") != "true" {
+				http.NotFound(w, req)
+				return
+			}
+			_ = json.NewEncoder(w).Encode(o)
+		})}
+		got := readBackFinishDecision(r, d, "ssh:mark-mac", in.scopeText)
+		if (c.want == "") != (got == "") || !strings.Contains(got, c.want) || (got != "" && !strings.HasPrefix(got, "the world decision "+d+" does not approve this finish (F6): ")) {
+			t.Fatalf("%s: %q, want %q", c.name, got, c.want)
+		}
+	}
+}
+
+// §2.9 P1 "decision does not approve this scope", at the wiring: a confirmation that re-checks but
+// whose decision the daemon does not hold is refused (F6), and no decision is returned.
+func TestFinishConfirmOnSurfaceRefusesAnUnrecordedDecision(t *testing.T) {
+	creds, in := surfaceUnit(t)
+	env, _ := awaitEnv(sfT0, new(bool), func(n, now int64) {
+		if n == 1 {
+			ro, _ := store.OpenCredentialsReadOnly(creds)
+			rows, err := ro.PendingVerbRequests(boundedTestContext(t), now, 16)
+			_ = ro.Close()
+			if err != nil || len(rows) != 1 {
+				t.Fatalf("pending %v %v", rows, err)
+			}
+			d, _ := hex.DecodeString(rows[0].Digest)
+			c, _ := hex.DecodeString(rows[0].Challenge)
+			approveByHand(t, creds, rows[0].RequestID, verbSig(t, sigSigner, approval.VerbNamespace, [32]byte(d), [32]byte(c), rows[0].ExpiresAt), "ssh:mark-mac", now)
+		}
+	})
+	env.random = func(b []byte) error { return nil }
+	var tty, out bytes.Buffer
+	ref, refusal := finishConfirmOnSurface(boundedTestContext(t), memoryDaemon(http.NotFound), env, in, &tty, &out)
+	if ref != "" || !strings.HasPrefix(refusal, "the world decision sha256:"+strings.Repeat("d", 64)+" does not approve this finish (F6): ") || strings.Contains(out.String(), "confirmed by") {
+		t.Fatalf("%q %q\n%s", ref, refusal, out.String())
+	}
+}
diff --git a/host/daemon/surface_confirm_now_test.go b/host/daemon/surface_confirm_now_test.go
new file mode 100644
index 000000000..39be78300
--- /dev/null
+++ b/host/daemon/surface_confirm_now_test.go
@@ -0,0 +1,30 @@
+package daemon
+
+// Row 243 V3 sprint 5 (the sprint-4 judge's note 1, mutant J-NOW-C14): C14 records the finish
+// approval at the handler's own clock — the world decision carries the injected d.approveNow, as
+// C15's consumed_at does: one instant for the whole confirmation.
+
+import (
+	"encoding/json"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+)
+
+func TestConfirmPOSTRecordsAtItsClock(t *testing.T) {
+	a, row := newConfirmRig(t)
+	at := apT0 + 7
+	a.d.approveNow = func() int64 { return at }
+	if rec := a.confirm(row.RequestID, confirmSSH(t, apSigner, approval.VerbNamespace, row)); rec.Code != 200 {
+		t.Fatalf("%d %s", rec.Code, rec.Body.String())
+	}
+	got := a.readVerb(row.RequestID)
+	var decision struct {
+		Now int64 `json:"now"`
+	}
+	obj, ok, err := a.d.store.GetObject(boundedTestContext(t), hashref.MustParse(got.DecisionRef))
+	if err != nil || !ok || json.Unmarshal(obj.Payload, &decision) != nil || decision.Now != at || got.ConsumedAt != at {
+		t.Fatalf("decision at %d (%v %v), consumed at %d; the clock said %d", decision.Now, ok, err, got.ConsumedAt, at)
+	}
+}
````

### Patch M12c — the flag, its two refusals, the prompt hook; the rig with a real daemon

````diff M12c
diff --git a/cmd/ailang-worldd/episode_finish.go b/cmd/ailang-worldd/episode_finish.go
index 01b6c9f98..076ed573d 100644
--- a/cmd/ailang-worldd/episode_finish.go
+++ b/cmd/ailang-worldd/episode_finish.go
@@ -27,6 +27,7 @@ import (
 
 const episodeFinishHelp = `usage: ailang-worldd episode finish <episode> --preset <name> [--presets-dir D]
     [--ready] [--exec-profiles-dir D] [--title <text>] [--allow-failing <check>]... [--include-unattributed] [--dry-run]
+    [--confirm-on-surface]
 
 ATTENDED: requires a controlling terminal before reading files or contacting
 the daemon. Runs the checks pinned at episode start, records each outcome,
@@ -40,6 +41,9 @@ Each unsandboxed check needs its own typed y at that terminal.
   --ready                 open a ready PR (default: draft)
   --exec-profiles-dir D    pinned profile lookup (default: private exec-profiles)
   --dry-run               print git show --stat and body, then reset HEAD
+  --confirm-on-surface    confirm the pull request with an enrolled key on the World
+                          app's card instead of a typed y; waits up to 10 minutes;
+                          the terminal is still required
 
 Without --dry-run, writes <sessions>/<episode>.pr-body.md (0600) and prints
 an exact-byte approval at the terminal. Typed y records approval, pushes and
@@ -67,6 +71,7 @@ func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEn
 	ready := fs.Bool("ready", false, "open a ready PR")
 	profilesDir := fs.String("exec-profiles-dir", defaultExecProfilesDir(), "private exec profiles")
 	dry := fs.Bool("dry-run", false, "print then reset local commit")
+	onSurface := fs.Bool("confirm-on-surface", false, "confirm the pull request with an enrolled key on the surface")
 	var allowed multiFlag
 	fs.Var(&allowed, "allow-failing", "allow this failing check")
 	pos, err := verbs.ParseInterleaved(fs, args)
@@ -87,6 +92,9 @@ func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEn
 		return fail("refusing: no controlling terminal", "finish is an attended act", "run it in a terminal (an agent or CI can't finish an episode)")
 	}
 	defer tty.Close()
+	if *onSurface && *dry { // row 243 V3 §2.7 step 1
+		return fail("--confirm-on-surface cannot be combined with --dry-run", "a dry run opens no pull request", "drop one of them")
+	}
 	if len(pos) != 1 || !workspaceEpisodePattern.MatchString(pos[0]) || *preset == "" {
 		return fail("episode and --preset <name> are required", "finish needs a pinned start record", "ailang-worldd episode finish --help")
 	}
@@ -251,6 +259,11 @@ func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEn
 		if err := validateFinishChecks(start.Checks, allowed, wt); err != nil {
 			return fail(err.Error(), "all checks must be valid before any runs", "fix checks in the next episode's preset")
 		}
+		for _, check := range start.Checks { // row 243 V3 §2.7 step 2: before any check runs
+			if *onSurface && check.Unsandboxed {
+				return fail("check "+check.Name+" is unsandboxed: it needs its own typed y", "--confirm-on-surface replaces only the pull request's y", "run finish without --confirm-on-surface")
+			}
+		}
 		profile, err := finishProfile(ctx, start, credentials, *profilesDir)
 		if err != nil {
 			return fail(err.Error(), "finish must use the start profile binding", "restore the pinned exec profile")
@@ -406,6 +419,15 @@ func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEn
 		scope, scopeText = reused, rep.Finish.Scope
 		approval.Approval = rep.Finish.Approval
 		fmt.Fprintf(tty, "resuming under approval %s, expires in %d s\n", approval.Approval, scope.ExpiresAt-time.Now().Unix())
+	} else if *onSurface { // row 243 V3 §2.7 step 3: a key on the surface instead of the typed y
+		ref, refusal := finishConfirmOnSurface(ctx, c, env, surfaceFinish{credentials: credentials, addr: p.Addr, ep: ep, sha: sha,
+			argv: canonicalFinishArgv(fs, allowed, ep), scope: scope, scopeText: scopeText, title: actualTitle, body: rendered,
+			committer: pinned.Committer, protected: protected}, tty, stdout)
+		if refusal != "" {
+			return fail(refusal, "the pull request was not confirmed on the surface; commit kept at "+sha, steps(sha))
+		}
+		approval.Approval = ref
+		finishStep("approve")
 	} else {
 		finishPrompt(tty, scope, actualTitle, rendered, protected, pinned.Committer)
 		answer, _ := reader.ReadString('\n')
diff --git a/cmd/ailang-worldd/episode_finish_surface_rig_test.go b/cmd/ailang-worldd/episode_finish_surface_rig_test.go
new file mode 100644
index 000000000..a1e5c557b
--- /dev/null
+++ b/cmd/ailang-worldd/episode_finish_surface_rig_test.go
@@ -0,0 +1,440 @@
+package main
+
+// Row 243 V3 sprint 5 M12 (design §2.7, §2.9 P1, §6 CF1–CF3, CF5, CF6, CF10, CF12, CF13; S4 §1.4 F4,
+// CX21's CLI half): episode finish --confirm-on-surface end to end on the finish rig, with a real
+// daemon (daemon.New over the rig's world and credentials stores, served through its Handler on a
+// loopback request) answering POST /surface/confirm/{id} and GET /v1/objects/{ref}. The card is
+// played by POSTing sigSigner's SSHSIG over approval.VerbMessage of the row the CLI stored; the
+// terminal answers y and counts its reads (the surface mode never reads it); the clock is the
+// await rig's fake one, started at the wall clock the daemon reads.
+
+import (
+	"bytes"
+	"context"
+	"database/sql"
+	"encoding/hex"
+	"encoding/json"
+	"fmt"
+	"io"
+	"net"
+	"net/http"
+	"net/http/httptest"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// countingTerminal answers y and counts every read; it records what is written to it.
+type countingTerminal struct {
+	bytes.Buffer
+	in    io.Reader
+	reads int
+}
+
+func (t *countingTerminal) Read(p []byte) (int, error) { t.reads++; return t.in.Read(p) }
+func (t *countingTerminal) Close() error               { return nil }
+
+type finishSurfaceRig struct {
+	t        *testing.T
+	f        *finishFixture
+	d        *daemon.Daemon
+	creds    string
+	tty      *countingTerminal
+	now      int64
+	start    int64
+	waits    int64
+	draws    byte
+	onWait   func(n int64)
+	objReads int
+}
+
+func newSurfaceRig(t *testing.T) *finishSurfaceRig {
+	t.Helper()
+	f := newFinishFixture(t)
+	f.allowRevoked = true // a re-run records after the first run's revocation
+	r := &finishSurfaceRig{t: t, f: f, creds: f.health.CredentialsDB, now: time.Now().Unix()}
+	seedApprover(t, r.creds, "mark-mac", sigSigner.PublicKey(), r.now-100)
+	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: f.db, CredentialsPath: r.creds, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = d.Close() })
+	r.d = d
+	finishClient = func(string) *client {
+		return memoryDaemon(func(w http.ResponseWriter, req *http.Request) {
+			if strings.HasPrefix(req.URL.Path, "/v1/objects/") { // the read-back reaches the real daemon
+				r.objReads++
+				rec := r.handle(req.Method, req.URL.RequestURI(), nil)
+				w.WriteHeader(rec.Code)
+				_, _ = w.Write(rec.Body.Bytes())
+				return
+			}
+			f.serve(w, req)
+		})
+	}
+	return r
+}
+
+// handle serves one loopback request on the real daemon, as the card's page sends it.
+func (r *finishSurfaceRig) handle(method, target string, body []byte) *httptest.ResponseRecorder {
+	req := httptest.NewRequest(method, target, bytes.NewReader(body))
+	req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
+	req.Host = "localhost:7644"
+	if method == http.MethodPost {
+		for k, v := range map[string]string{"Origin": "http://localhost:7644", "Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"} {
+			req.Header.Set(k, v)
+		}
+	}
+	rec := httptest.NewRecorder()
+	r.d.Handler().ServeHTTP(rec, req)
+	return rec
+}
+
+// env is the finish's env: the counting terminal, the fake clock, and a counter as randomness (each
+// draw fills its bytes with the next byte value, so ids and challenges are known to the test).
+func (r *finishSurfaceRig) env(answer string) sessionEnv {
+	r.tty = &countingTerminal{in: strings.NewReader(answer)}
+	return sessionEnv{
+		openTerminal: func() (io.ReadWriteCloser, error) { return r.tty, nil },
+		now:          func() int64 { return r.now },
+		random: func(b []byte) error {
+			r.draws++
+			for i := range b {
+				b[i] = r.draws
+			}
+			return nil
+		},
+		after: func(d time.Duration) <-chan time.Time {
+			r.waits++
+			r.now += int64(d / time.Second)
+			if r.onWait != nil {
+				r.onWait(r.waits)
+			}
+			ch := make(chan time.Time, 1)
+			ch <- time.Unix(r.now, 0)
+			return ch
+		},
+	}
+}
+
+// run is one finish with --confirm-on-surface; its fake clock starts at the wall clock the daemon
+// reads (a request created ahead of it would be "clock behind" at the POST).
+func (r *finishSurfaceRig) run(extra ...string) (int, string, string) {
+	r.waits, r.now = 0, time.Now().Unix()
+	r.start = r.now
+	return r.f.run(r.env("y\n"), append([]string{"--include-unattributed", "--confirm-on-surface"}, extra...)...)
+}
+
+// row reads one verb request from the credentials store.
+func (r *finishSurfaceRig) row(id string) store.VerbRequestRow {
+	r.t.Helper()
+	ro, err := store.OpenCredentialsReadOnly(r.creds)
+	if err != nil {
+		r.t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	row, ok, err := ro.VerbRequest(boundedTestContext(r.t), id)
+	if err != nil || !ok {
+		r.t.Fatalf("row %s: %v %v", id, ok, err)
+	}
+	return row
+}
+
+// card plays the card for id: s's SSHSIG over the stored row's approval.VerbMessage, POSTed to S4's
+// POST /surface/confirm/{id}; it returns the status.
+func (r *finishSurfaceRig) card(s *approvaltest.Signer, id string) int {
+	r.t.Helper()
+	row := r.row(id)
+	dg, _ := hex.DecodeString(row.Digest)
+	ch, _ := hex.DecodeString(row.Challenge)
+	body, _ := json.Marshal(map[string]any{"v": 1, "digest": row.Digest, "sshsig": string(verbSig(r.t, s, approval.VerbNamespace, [32]byte(dg), [32]byte(ch), row.ExpiresAt))})
+	return r.handle(http.MethodPost, "/surface/confirm/"+id, body).Code
+}
+
+// nextWallSecond waits (≤ 1 s, bounded at 3 s) for the wall clock the daemon reads to leave t0: a
+// finish approval of the same scope in the same second is the same request object (design W59).
+func nextWallSecond(t0 int64) {
+	for limit := time.Now().Add(3 * time.Second); time.Now().Unix() == t0 && time.Now().Before(limit); {
+		time.Sleep(20 * time.Millisecond)
+	}
+}
+
+// object is the real daemon's GET /v1/objects/{ref} payload, decoded into out.
+func (r *finishSurfaceRig) object(ref string, out any) {
+	r.t.Helper()
+	var o wireObject
+	rec := r.handle(http.MethodGet, "/v1/objects/"+ref+"?payload=true", nil)
+	if rec.Code != 200 || json.Unmarshal(rec.Body.Bytes(), &o) != nil || json.Unmarshal(o.Payload, out) != nil {
+		r.t.Fatalf("object %s: %d %s", ref, rec.Code, rec.Body.String())
+	}
+}
+
+// id is the request id the n-th run's first draw makes ("vr-" and 8 bytes of value n).
+func drawnID(n byte) string { return "vr-" + strings.Repeat(fmt.Sprintf("%02x", n), 8) }
+
+// kept fails unless a refusal kept the commit and dispatched nothing (design §2.7: every refusal).
+func (r *finishSurfaceRig) kept(rc int, errs, want string) string {
+	r.t.Helper()
+	sha := strings.TrimSpace(r.f.gitOutput("-C", r.f.wt, "rev-parse", "HEAD"))
+	if rc != exitUsage || !strings.Contains(errs, want) || !strings.Contains(errs, "commit kept at "+sha) || !strings.Contains(errs, "next: ailang-worldd episode finish ep1 --preset fixture") ||
+		sha == r.f.before || len(r.f.seam.calls) != 0 || len(r.f.finishRecords) != 0 || len(r.f.scopes) != 0 {
+		r.t.Fatalf("rc %d calls %v records %d\nstderr %s\nwant %q", rc, r.f.seam.calls, len(r.f.finishRecords), errs, want)
+	}
+	return sha
+}
+
+// CF3 (+ CF13, CF14's rig half, judge-403 note 1's CLI half): the card confirms; the CLI re-checks,
+// reads the decision back from the real daemon and opens the pull request with that decision. The
+// terminal is never read; it shows exactly the stored prompt, which is FinishPromptText of the
+// stored fields; the stored row is the CLI's own (id, challenge, times from the env).
+func TestFinishConfirmOnSurfaceOpensThePR(t *testing.T) {
+	r := newSurfaceRig(t)
+	id := drawnID(1)
+	r.onWait = func(n int64) {
+		if n == 1 {
+			if code := r.card(sigSigner, id); code != 200 {
+				t.Errorf("card POST: %d", code)
+			}
+		}
+	}
+	r.f.rep.Start.Checks = []store.StartCheck{{Name: "pass", Argv: []string{"go", "test"}, Cwd: ".", TimeoutS: 1, Exec: "pinned", Args: []string{"test"}}}
+	rc, out, errs := r.run("--allow-failing", "pass")
+	if rc != exitOK {
+		t.Fatalf("rc %d\nstdout %s\nstderr %s", rc, out, errs)
+	}
+	start, row := r.start, r.row(id)
+	req, err := approval.DecodeVerb(row.Request)
+	if err != nil {
+		t.Fatal(err)
+	}
+	sha := strings.TrimSpace(r.f.gitOutput("-C", r.f.wt, "rev-parse", "HEAD"))
+	if row.Challenge != strings.Repeat("02", 32) || row.EpisodeID != "ep1" || row.CreatedAt != start || row.ExpiresAt != start+600 || row.Digest != finishSHA(row.Request) ||
+		row.Outcome != store.VerbOutcomeApproved || row.ApprovedBy != "ssh:mark-mac" {
+		t.Fatalf("row %+v", row)
+	}
+	dg, _ := hex.DecodeString(row.Digest)
+	h, _ := approval.VerbMessageHash([32]byte(dg), [32]byte(bytes.Repeat([]byte{2}, 32)), start+600)
+	pr := r.f.finishRecords[0].PR.URL // the fake gh's own answer
+	want := fmt.Sprintf("committed at %s\nbody: %s\nconfirm request %s on the surface (expires %s, message hash %s):\n  open http://localhost:47330/surface/confirm/%s\n"+
+		"the typed y is not read in this mode; Ctrl-C stops waiting and keeps the commit at %s\nconfirmed by ssh:mark-mac (decision %s)\n", sha, filepath.Join(r.f.sessions, "ep1.pr-body.md"),
+		id, approverTime(start+600), hex.EncodeToString(h[:8]), id, sha, row.DecisionRef)
+	if !strings.HasPrefix(out, want) || !strings.HasSuffix(out, "\n"+pr+"\n") || pr == "" {
+		t.Fatalf("stdout\n%s\nwant prefix\n%s\nand the PR %q last", out, want, pr)
+	}
+	if rec := r.f.finishRecords; len(rec) != 1 || rec[0].Approval != row.DecisionRef || rec[0].Outcome != "done" || len(r.f.scopes) != 0 || r.f.outcomes != 1 {
+		t.Fatalf("records %+v, typed approve phases %d", rec, len(r.f.scopes))
+	}
+	if r.tty.reads != 0 || r.tty.String() != req.Prompt || r.objReads != 2 {
+		t.Fatalf("terminal reads %d, object reads %d, shows\n%q\nwant the stored prompt\n%q", r.tty.reads, r.objReads, r.tty.String(), req.Prompt)
+	}
+	body, _ := os.ReadFile(filepath.Join(r.f.sessions, "ep1.pr-body.md"))
+	scope, err := broker.ParseFinishApprovalScope(req.Scope)
+	if err != nil || req.Prompt != broker.FinishPromptText(scope, req.Title, req.Body, req.Protected, req.Committer) || req.Body != string(body) ||
+		req.Scope != r.f.finishRecords[0].Scope || scope.SHA != sha || !scope.Draft || req.Protected == nil {
+		t.Fatalf("request %+v (%v)", req, err)
+	}
+	if wantArgv := []string{"ailang-worldd", "episode", "finish", "ep1", "--allow-failing", "pass", "--confirm-on-surface", "--include-unattributed", "--preset", "fixture", "--presets-dir", r.f.dir}; !reflect.DeepEqual(req.Argv, wantArgv) {
+		t.Fatalf("argv %q, want %q", req.Argv, wantArgv)
+	}
+	var decision struct {
+		Decision, DecidedBy string
+		Now                 int64
+	}
+	r.object(row.DecisionRef, &decision)
+	if decision.Decision != "approve" || decision.DecidedBy != "ssh:mark-mac" || decision.Now != row.ConsumedAt {
+		t.Fatalf("decision %+v, row consumed at %d", decision, row.ConsumedAt)
+	}
+}
+
+// CF1, CF2: --dry-run refuses after the terminal and before any read; an unsandboxed pinned check
+// refuses before the check loop — no check ran, no outcome recorded, no commit, the terminal unread.
+func TestFinishConfirmOnSurfaceRefusesBeforeAnyCheck(t *testing.T) {
+	r := newSurfaceRig(t)
+	rc, _, errs := r.run("--dry-run")
+	if rc != exitUsage || !strings.Contains(errs, "--confirm-on-surface cannot be combined with --dry-run\n  why: a dry run opens no pull request\n  next: drop one of them") || len(r.f.requests) != 0 || r.tty.reads != 0 {
+		t.Fatalf("CF1: rc %d requests %v reads %d\n%s", rc, r.f.requests, r.tty.reads, errs)
+	}
+	r.f.rep.Start.Checks = []store.StartCheck{{Name: "host", Argv: []string{"/bin/sh", "-c", "echo ran > marker"}, Cwd: ".", TimeoutS: 1, Unsandboxed: true}}
+	rc, _, errs = r.run()
+	_, statErr := os.Stat(filepath.Join(r.f.wt, "marker"))
+	if rc != exitUsage || !strings.Contains(errs, "check host is unsandboxed: it needs its own typed y\n  why: --confirm-on-surface replaces only the pull request's y\n  next: run finish without --confirm-on-surface") ||
+		r.f.outcomes != 0 || !os.IsNotExist(statErr) || r.tty.reads != 0 || r.tty.Len() != 0 || strings.TrimSpace(r.f.gitOutput("-C", r.f.wt, "rev-parse", "HEAD")) != r.f.before || pendingVerbRows(t, r.creds) != 0 {
+		t.Fatalf("CF2: rc %d outcomes %d marker %v reads %d\n%s", rc, r.f.outcomes, statErr, r.tty.reads, errs)
+	}
+}
+
+// protectedProfile pins an exec profile protecting *_test.go and a committer for f's episode (as the
+// typed path's TestEpisodeFinishProtected does) and changes a_test.go; it returns the profiles dir.
+func protectedProfile(t *testing.T, f *finishFixture) string {
+	t.Helper()
+	profiles := filepath.Join(f.root, "profiles")
+	raw, _ := json.Marshal(map[string]any{"profile": broker.ExecProfileVersion, "project": "fixture", "env": map[string]string{}, "path": []string{"/usr/bin", "/bin"},
+		"timeout_ms": 9000, "protected": []string{"*_test.go"}, "commands": map[string]any{"check": map[string]any{"argv": []string{"true"}, "max_args": 1}}})
+	profile, err := broker.ParseExecProfile(raw)
+	if err != nil || os.MkdirAll(profiles, 0o700) != nil || os.WriteFile(filepath.Join(profiles, "fixture.json"), raw, 0o600) != nil {
+		t.Fatal(err)
+	}
+	creds, err := store.OpenCredentials(f.health.CredentialsDB)
+	if err != nil || authority.Revoke(boundedTestContext(t), creds, f.credentialID) != nil {
+		t.Fatal(err)
+	}
+	_, f.credentialID, _, err = authority.MintBound(boundedTestContext(t), creds, "ep1", authority.ProfileBinding{Project: "fixture", Digest: profile.Digest},
+		[]broker.Capability{{Effect: "Workspace.Write", Scope: "world", Budget: 10, ExpiresAt: 4000000000}}, 3600, time.Now().Unix(), nil)
+	if creds.Close(); err != nil {
+		t.Fatal(err)
+	}
+	f.rep.Start.Preset.Bytes = `{"profile":"fixture","committer":"Test Operator <operator@example.invalid>"}`
+	f.write("a_test.go", "protected edit\n")
+	return profiles
+}
+
+// CF4's rig half, CF13, CF10: a request nobody confirms ends at expiry + grace with the commit kept;
+// the typed path's prompt over the same commit is the stored prompt byte for byte (the typed run
+// declines); a re-run stores a new request with the same prompt bytes, supersedes the old one at
+// its own clock, and opens the pull request when that one is confirmed.
+func TestFinishConfirmOnSurfaceRerunsAndTypedPrompt(t *testing.T) {
+	r := newSurfaceRig(t)
+	profiles := protectedProfile(t, r.f) // a protected path and a pinned committer in the prompt
+	first := drawnID(1)
+	r.onWait = func(n int64) {
+		if n == 1 {
+			r.now = r.row(first).ExpiresAt + 10 // the next poll is the last
+		}
+	}
+	rc, _, errs := r.run("--exec-profiles-dir", profiles)
+	sha := r.kept(rc, errs, "no confirmation of request "+first+" by ")
+	old := r.row(first)
+	stored, _ := approval.DecodeVerb(old.Request)
+	if old.Outcome != "" || r.waits != 1 || !strings.HasPrefix(stored.Prompt, "committer: Test Operator <operator@example.invalid>\n") || !strings.Contains(stored.Prompt, "\nprotected paths changed: a_test.go") {
+		t.Fatalf("first: %q %q after %d waits", old.Outcome, stored.Prompt, r.waits)
+	}
+	typed := &captureFinishTerminal{Reader: strings.NewReader("n\n")}
+	env := testMintEnv("", 1)
+	env.openTerminal = func() (io.ReadWriteCloser, error) { return typed, nil }
+	if rc, _, errs := r.f.run(env, "--include-unattributed"); rc != exitUsage || !strings.Contains(errs, "pull request declined") {
+		t.Fatalf("typed: %d %s", rc, errs)
+	}
+	if typed.String() != stored.Prompt+"Open this pull request? [y/N] " {
+		t.Fatalf("CF13: the typed prompt\n%q\nis not the stored prompt\n%q", typed.String(), stored.Prompt)
+	}
+	second := drawnID(3)
+	r.onWait = func(n int64) {
+		if n == 1 && r.card(sigSigner, second) != 200 {
+			t.Error("card POST")
+		}
+	}
+	if rc, out, errs := r.run(); rc != exitOK || strings.TrimSpace(r.f.gitOutput("-C", r.f.wt, "rev-parse", "HEAD")) != sha {
+		t.Fatalf("re-run: %d\n%s\n%s", rc, out, errs)
+	}
+	newer, old := r.row(second), r.row(first)
+	again, _ := approval.DecodeVerb(newer.Request)
+	if again.Prompt != stored.Prompt || newer.Outcome != store.VerbOutcomeApproved || old.Outcome != store.VerbOutcomeSuperseded || old.ConsumedAt != r.start {
+		t.Fatalf("CF10: newer %+v old %+v (re-run at %d)", newer, old, r.start)
+	}
+}
+
+// CF5 and CF6's rig rows (S1 Appendix C note 5): a newer finish supersedes the request; an approver
+// revoked after the confirmation was recorded makes the CLI refuse it (fail-safe) — each keeps the
+// commit and dispatches nothing; a re-run asks afresh (CF10).
+func TestFinishConfirmOnSurfaceKeepsTheCommit(t *testing.T) {
+	r := newSurfaceRig(t)
+	in := surfaceFinish{credentials: r.creds, ep: "ep1", argv: []string{"ailang-worldd", "episode", "finish", "ep1"}, scopeText: "x", title: "t", body: "b"}
+	r.onWait = func(n int64) {
+		if n == 1 {
+			storeOwn(t, in, "vr-00000000000000e1", 0xe1, r.now) // a newer finish of ep1
+		}
+	}
+	rc, _, errs := r.run()
+	r.kept(rc, errs, "request "+drawnID(1)+" was superseded by a newer finish of ep1")
+	revoked := drawnID(3)
+	r.onWait = func(n int64) {
+		if n == 1 {
+			if code := r.card(sigSigner, revoked); code != 200 {
+				t.Errorf("card POST: %d", code)
+			}
+			revokeSeeded(t, r.creds, "mark-mac", r.now)
+		}
+	}
+	rc, _, errs = r.run()
+	r.kept(rc, errs, "the confirmation of request "+revoked+" does not verify: ")
+	if !strings.Contains(errs, "revoked") || r.row(revoked).Outcome != store.VerbOutcomeApproved || r.objReads != 0 {
+		t.Fatalf("revoked: %s (object reads %d)", errs, r.objReads)
+	}
+}
+
+// CX21's CLI half (S4 §1.4 F4): C15 fails after C14 (a verb_requests UPDATE trigger) → 500, one
+// decision no row names; the identical body a second later → 200, a second decision the row names;
+// the waiting finish dispatches with the row's decision, never the orphan.
+func TestFinishConfirmOnSurfaceDispatchesTheNamedDecision(t *testing.T) {
+	r := newSurfaceRig(t)
+	id := drawnID(1)
+	r.onWait = func(n int64) {
+		if n != 1 {
+			return
+		}
+		credSQL(t, r.creds, `CREATE TRIGGER fail_consume BEFORE UPDATE ON verb_requests WHEN NEW.outcome = 'approved' BEGIN SELECT RAISE(ABORT, 'injected consume failure'); END`)
+		t0 := time.Now().Unix()
+		if code := r.card(sigSigner, id); code != 500 || r.row(id).Outcome != "" {
+			t.Errorf("first POST: %d", code)
+		}
+		credSQL(t, r.creds, `DROP TRIGGER fail_consume`)
+		nextWallSecond(t0) // a new decision, not the same one again (W59)
+		if code := r.card(sigSigner, id); code != 200 {
+			t.Errorf("retry POST: %d", code)
+		}
+	}
+	if rc, out, errs := r.run(); rc != exitOK {
+		t.Fatalf("rc %d\n%s\n%s", rc, out, errs)
+	}
+	named := r.row(id).DecisionRef
+	db, err := sql.Open("sqlite", r.f.db)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	var all, isNamed int
+	if err := db.QueryRow(`SELECT count(*), count(*) FILTER (WHERE hash_ref = ?) FROM objects WHERE semantic_id = ?`, named, broker.ApprovalDecisionV1).Scan(&all, &isNamed); err != nil {
+		t.Fatal(err)
+	}
+	if all != 2 || isNamed != 1 || r.f.finishRecords[0].Approval != named {
+		t.Fatalf("%d decisions (%d named %s); the record names %s", all, isNamed, named, r.f.finishRecords[0].Approval)
+	}
+}
+
+// CF12 (S1 Appendix C note 4): a surface-confirmed finish whose push failed is resumed by a re-run
+// under its recorded approval — mode resume-dispatch reuses it (the unchanged branch, design §2.7
+// step 3): no new request, no wait, no re-check or read-back (CF6/CF7 are skipped), F6's re-check of
+// HEAD, remote, title and body still runs; the pull request opens under the first decision.
+func TestFinishConfirmOnSurfaceResumeReusesTheApproval(t *testing.T) {
+	r := newSurfaceRig(t)
+	id := drawnID(1)
+	r.onWait = func(n int64) {
+		if n == 1 && r.card(sigSigner, id) != 200 {
+			t.Error("card POST")
+		}
+	}
+	r.f.seam.fail = "push"
+	if rc, _, errs := r.run(); rc != exitFatal || len(r.f.finishRecords) != 1 || r.f.finishRecords[0].Outcome != "push-failed" {
+		t.Fatalf("first: %d %s", rc, errs)
+	}
+	decision := r.row(id).DecisionRef
+	r.f.seam.fail, r.objReads = "", 0
+	r.onWait = func(int64) { t.Error("the resume waited") }
+	rc, out, errs := r.run()
+	if rc != exitOK || pendingVerbRows(t, r.creds) != 1 || r.objReads != 0 || r.tty.reads != 0 || !strings.HasPrefix(r.tty.String(), "resuming under approval "+decision+", expires in ") {
+		t.Fatalf("resume: %d rows %d object reads %d tty %q\n%s\n%s", rc, pendingVerbRows(t, r.creds), r.objReads, r.tty.String(), out, errs)
+	}
+	if rec := r.f.finishRecords; len(rec) != 2 || rec[1].Approval != decision || rec[1].Outcome != "done" {
+		t.Fatalf("records %+v", rec)
+	}
+}
````

### Patch M13 — the bar's reason, `approve vr-` in the bar and the CLI

````diff M13
diff --git a/cmd/ailang-worldd/approve.go b/cmd/ailang-worldd/approve.go
index b793bf06e..b36041140 100644
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -225,6 +225,9 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 	id := ""
 	if !*list {
 		id = positional[0]
+		if store.VerbRequestIDPattern.MatchString(id) { // row 243 V3 §2.8: the bar's grammar and the CLI's agree on the id
+			return refuse("request %s is a finish confirmation: confirm it on its card, /surface/confirm/%s; the waiting episode finish opens the pull request", id, id)
+		}
 		if decisionRefPattern.MatchString(id) { // row 201 A4: a request ref selects the decision mode (design §2.10)
 			if *await { // CD3
 				return refuse("a request ref is decided at once by the daemon: there is nothing to await")
diff --git a/cmd/ailang-worldd/approve_verb_request_test.go b/cmd/ailang-worldd/approve_verb_request_test.go
new file mode 100644
index 000000000..ddaeba946
--- /dev/null
+++ b/cmd/ailang-worldd/approve_verb_request_test.go
@@ -0,0 +1,24 @@
+package main
+
+// Row 243 V3 sprint 5 M13 (design §2.8, §6 CB4): `approve vr-<16 hex>` at the CLI points at the
+// finish confirmation's card, as the bar's ApproveTarget does; it reads no store and mints nothing.
+
+import (
+	"bytes"
+	"path/filepath"
+	"testing"
+)
+
+func TestApproveVerbRequestPointsAtItsCard(t *testing.T) {
+	db := filepath.Join(t.TempDir(), "world.db")
+	var stdout, stderr bytes.Buffer
+	code := runApprove([]string{"vr-0123456789abcdef", "--db", db}, &stdout, &stderr, noTerminalEnv(1))
+	want := "ailang-worldd approve: request vr-0123456789abcdef is a finish confirmation: confirm it on its card, /surface/confirm/vr-0123456789abcdef; the waiting episode finish opens the pull request\n"
+	if code != exitUsage || stderr.String() != want || stdout.Len() != 0 {
+		t.Fatalf("rc %d\nstderr %q\nwant   %q", code, stderr.String(), want)
+	}
+	stderr.Reset()
+	if code := runApprove([]string{"vr-0123", "--db", db}, &stdout, &stderr, noTerminalEnv(1)); code != exitUsage || !bytes.Contains(stderr.Bytes(), []byte(`"vr-0123" is not a mint request id`)) {
+		t.Fatalf("a short id: rc %d %q", code, stderr.String())
+	}
+}
diff --git a/host/daemon/surface_verb_outcomes_test.go b/host/daemon/surface_verb_outcomes_test.go
index 562a77d11..047c9565d 100644
--- a/host/daemon/surface_verb_outcomes_test.go
+++ b/host/daemon/surface_verb_outcomes_test.go
@@ -35,6 +35,7 @@ func TestSurfaceVerbOutcomes(t *testing.T) {
 		{"episode launch w1 --preset p", "terminal", "host", "episode launch", "ailang-worldd episode launch w1 --preset p", "", "needs a host process"},
 		{"approve " + mr, "navigate", "navigate", "approve", "ailang-worldd approve " + mr, "/surface/approve/" + mr, ""},
 		{"ailang-worldd approve " + ref, "navigate", "navigate", "approve", "ailang-worldd approve " + ref, "/surface/decide/" + ref, ""},
+		{"approve vr-0123456789abcdef", "navigate", "navigate", "approve", "ailang-worldd approve vr-0123456789abcdef", "/surface/confirm/vr-0123456789abcdef", ""},
 		{"approve " + mr + " --deny", "refused", "navigate", "approve", "ailang-worldd approve " + mr + " --deny", "", "the app opens approve"},
 		{"approver enroll --name x", "refused", "never", "approver", "ailang-worldd approver enroll --name x", "", "never run from the app"},
 		{"session frob", "refused", "never", "session", "ailang-worldd session frob", "", "never run from the app"},
@@ -75,7 +76,7 @@ func TestSurfaceVerbOutcomes(t *testing.T) {
 	golden := map[string]string{}
 	for _, q := range []string{"row-189@6", "world shared@0ffbb9c", "back", "world why", "why", "log", "world session",
 		"ailang-worldd world get world-1", "episode", "world episode", "world a b c", "episode finish w1 --preset p",
-		"approve mr-0123456789abcdef", "approve " + ref, "episode report nobody", "episode report w1", "episode report w1 --json",
+		"approve mr-0123456789abcdef", "approve " + ref, "approve vr-0123456789abcdef", "episode report nobody", "episode report w1", "episode report w1 --json",
 		strings.Repeat("x", verbs.MaxLine+1), "why 2", "why head --json", "why 9", "why --result f"} {
 		_, raw := verbGet(a, q)
 		golden[q] = raw
diff --git a/host/surfacepage/app/app.js b/host/surfacepage/app/app.js
index 96f097957..36a1a366a 100644
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -493,7 +493,7 @@
     if (!need || !need.every(k => typeof b[k] === 'string' && b[k] !== '')) return UNREADABLE;
     if (b.outcome === 'address') return { address: true };
     if (b.outcome === 'refused') return typeof b.cli === 'string' && b.cli !== '' ? { panel: { head: '$ ' + b.cli, body: b.reason, copy: b.cli } } : { say: b.reason };
-    if (b.outcome === 'navigate') return b.href.startsWith('/surface/approve/mr-') || b.href.startsWith('/surface/decide/sha256:') ? { assign: b.href } : { say: 'refused: unexpected target' };
+    if (b.outcome === 'navigate') return b.href.startsWith('/surface/approve/mr-') || b.href.startsWith('/surface/decide/sha256:') || b.href.startsWith('/surface/confirm/vr-') ? { assign: b.href } : { say: 'refused: unexpected target' };
     if (b.outcome === 'terminal') return { panel: { head: '$ ' + b.cli, body: b.reason + '\nAt a terminal: ' + b.cli, copy: b.cli } };
     const shown = [b.text, b.json, b.stderr].filter(x => typeof x === 'string' && x !== '');
     if (shown.length !== 1 || (b.exit !== undefined && b.exit !== 1 && b.exit !== 3) || (b.stderr !== undefined) !== (b.exit === 1)) return UNREADABLE;
diff --git a/host/surfacepage/app_verb_test.go b/host/surfacepage/app_verb_test.go
index c6ad5e0c1..16f189557 100644
--- a/host/surfacepage/app_verb_test.go
+++ b/host/surfacepage/app_verb_test.go
@@ -118,12 +118,12 @@ async function run(raw, fetchImpl) {
     assert.deepEqual(log, [['clear'], ['panel', act.panel]], q);
   }
   // Navigate: the two card prefixes (VP5's must-live).
-  for (const [q, href] of [[MR, '/surface/approve/mr-0123456789abcdef'], [REF, '/surface/decide/sha256:' + 'ab'.repeat(32)]]) {
+  for (const [q, href] of [[MR, '/surface/approve/mr-0123456789abcdef'], [REF, '/surface/decide/sha256:' + 'ab'.repeat(32)], ['approve vr-0123456789abcdef', '/surface/confirm/vr-0123456789abcdef']]) {
     const { act, log } = await run(q, golden(q));
     assert.equal(B(q).href, href); assert.deepEqual(act, { assign: href }); assert.deepEqual(log, [['clear'], ['assign', href]]);
   }
   // VP5: any other target is refused.
-  for (const href of ['https://evil/x', '/surface/model.json', '//evil/surface/approve/mr-']) {
+  for (const href of ['https://evil/x', '/surface/model.json', '//evil/surface/approve/mr-', '/surface/confirm/x', '/surface/confirm/mr-0123456789abcdef']) {
     const { log } = await run(MR, reply(200, JSON.stringify({ ...B(MR), href })));
     assert.deepEqual(log, [['say', 'refused: unexpected target']], href);
   }
diff --git a/host/surfacepage/testdata/verb-route.golden.json b/host/surfacepage/testdata/verb-route.golden.json
index 30fa41571..09e97fbb5 100644
--- a/host/surfacepage/testdata/verb-route.golden.json
+++ b/host/surfacepage/testdata/verb-route.golden.json
@@ -2,9 +2,10 @@
   "ailang-worldd world get world-1": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"ailang-worldd world get world-1\",\"argv\":[\"ailang-worldd\",\"world\",\"get\",\"world-1\"],\"verb\":\"world get\",\"class\":\"other\",\"outcome\":\"refused\",\"reason\":\"world get is a CLI command the app does not run; at a terminal: ailang-worldd world get world-1\",\"cli\":\"ailang-worldd world get world-1\"}\n",
   "approve mr-0123456789abcdef": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"approve mr-0123456789abcdef\",\"argv\":[\"approve\",\"mr-0123456789abcdef\"],\"verb\":\"approve\",\"class\":\"navigate\",\"outcome\":\"navigate\",\"cli\":\"ailang-worldd approve mr-0123456789abcdef\",\"href\":\"/surface/approve/mr-0123456789abcdef\"}\n",
   "approve sha256:abababababababababababababababababababababababababababababababab": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"approve sha256:abababababababababababababababababababababababababababababababab\",\"argv\":[\"approve\",\"sha256:abababababababababababababababababababababababababababababababab\"],\"verb\":\"approve\",\"class\":\"navigate\",\"outcome\":\"navigate\",\"cli\":\"ailang-worldd approve sha256:abababababababababababababababababababababababababababababababab\",\"href\":\"/surface/decide/sha256:abababababababababababababababababababababababababababababababab\"}\n",
+  "approve vr-0123456789abcdef": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"approve vr-0123456789abcdef\",\"argv\":[\"approve\",\"vr-0123456789abcdef\"],\"verb\":\"approve\",\"class\":\"navigate\",\"outcome\":\"navigate\",\"cli\":\"ailang-worldd approve vr-0123456789abcdef\",\"href\":\"/surface/confirm/vr-0123456789abcdef\"}\n",
   "back": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"back\",\"argv\":[\"back\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
   "episode": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode\",\"argv\":[\"episode\"],\"verb\":\"episode\",\"class\":\"refused\",\"outcome\":\"refused\",\"reason\":\"ailang-worldd episode needs a subcommand: start | launch | finish | report\",\"cli\":\"ailang-worldd episode\"}\n",
-  "episode finish w1 --preset p": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode finish w1 --preset p\",\"argv\":[\"episode\",\"finish\",\"w1\",\"--preset\",\"p\"],\"verb\":\"episode finish\",\"class\":\"mutate\",\"outcome\":\"terminal\",\"reason\":\"episode finish changes World. The app cannot run it yet.\",\"cli\":\"ailang-worldd episode finish w1 --preset p\"}\n",
+  "episode finish w1 --preset p": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode finish w1 --preset p\",\"argv\":[\"episode\",\"finish\",\"w1\",\"--preset\",\"p\"],\"verb\":\"episode finish\",\"class\":\"mutate\",\"outcome\":\"terminal\",\"reason\":\"episode finish changes World. The app cannot run it yet; at a terminal, --confirm-on-surface lets you approve its pull request here with your key.\",\"cli\":\"ailang-worldd episode finish w1 --preset p\"}\n",
   "episode report nobody": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode report nobody\",\"argv\":[\"episode\",\"report\",\"nobody\"],\"verb\":\"episode report\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd episode report nobody\",\"stderr\":\"ailang-worldd episode report: episode nobody has no committed calls\\n\",\"exit\":1}\n",
   "episode report w1": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode report w1\",\"argv\":[\"episode\",\"report\",\"w1\"],\"verb\":\"episode report\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd episode report w1\",\"text\":\"episode w1: 2 calls, entries 1-2, 2027-01-15T08:00:01Z .. 2027-01-15T08:00:02Z\\ncalls by tool:\\n  \\u003cimg src=x onerror=alert(1)\\u003e      1\\n  ailang-read             1\\nreads: 0 calls, 0 files, 1 outside the worktree\\neffects:\\n  Workspace.Read         allowed               1\\nbudgets (minted / used / left):\\n  Workspace.Read@worktree     50      1     49\\nplan refusals (refused before any effect): 1 (\\u003cimg src=x onerror=alert(1)\\u003e 1)\\n  #2 \\u003cimg src=x onerror=alert(1)\\u003e: \\u003cb\\u003ehostile\\u003c/b\\u003e \\\"reason\\\"\\nhandler refusals (effect allowed, handler ok:false): 0\\n\"}\n",
   "episode report w1 --json": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode report w1 --json\",\"argv\":[\"episode\",\"report\",\"w1\",\"--json\"],\"verb\":\"episode report\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd episode report w1 --json\",\"json\":\"{\\\"episode\\\":\\\"w1\\\",\\\"entries\\\":[1,2],\\\"total_calls\\\":2,\\\"calls_by_tool\\\":{\\\"\\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\":1,\\\"ailang-read\\\":1},\\\"effects\\\":[{\\\"effect\\\":\\\"Workspace.Read\\\",\\\"decision\\\":\\\"allowed\\\",\\\"count\\\":1}],\\\"budgets\\\":[{\\\"effect\\\":\\\"Workspace.Read\\\",\\\"scope\\\":\\\"worktree\\\",\\\"minted\\\":50,\\\"used\\\":1,\\\"left\\\":49}],\\\"plan_refusals\\\":{\\\"count\\\":1,\\\"by_tool\\\":{\\\"\\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\":1},\\\"first\\\":[{\\\"entry\\\":2,\\\"tool\\\":\\\"\\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\",\\\"reason\\\":\\\"\\\\u003cb\\\\u003ehostile\\\\u003c/b\\\\u003e \\\\\\\"reason\\\\\\\"\\\"}]},\\\"handler_refusals\\\":{\\\"count\\\":0,\\\"by_tool\\\":{},\\\"first\\\":[]},\\\"writes\\\":{\\\"written\\\":[],\\\"produced\\\":[],\\\"deltas_truncated\\\":[],\\\"delta_errors\\\":[],\\\"written_truncated\\\":0,\\\"produced_truncated\\\":0,\\\"unlisted\\\":0,\\\"deltas_truncated_omitted\\\":0,\\\"delta_errors_omitted\\\":0},\\\"reads\\\":{\\\"count\\\":0,\\\"files\\\":[],\\\"truncated\\\":0,\\\"unlisted\\\":1},\\\"first_at\\\":\\\"2027-01-15T08:00:01Z\\\",\\\"last_at\\\":\\\"2027-01-15T08:00:02Z\\\",\\\"calls_by_surface\\\":{\\\"mcp\\\":2},\\\"uncommitted_intents\\\":0,\\\"notes\\\":{\\\"budgets\\\":\\\"per (effect, scope) with at least one effect record: left = latest budgetAfter, used = cost of allowed records, minted = left + used\\\",\\\"entries\\\":\\\"first and last log entries of the episode's committed calls\\\",\\\"first_at\\\":\\\"coordinator logical time (unix seconds) of the first and last committed call, UTC\\\",\\\"handler_refusals\\\":\\\"calls whose effect was allowed and ran, but whose handler result is ok:false (the effect record alone shows these as allowed)\\\",\\\"plan_refusals\\\":\\\"calls whose plan named no effect and answered ok:false (refused before any effect ran)\\\",\\\"reads\\\":\\\"committed ailang-read calls with a v2 record whose every effect was allowed and did not fail and whose handler result is not ok:false; the path is the record input's \\\\\\\"path\\\\\\\" (a ranged read counts); a path that is not local and relative counts in unlisted and is never shown; files are ordered by first entry, at most 1000\\\",\\\"writes\\\":\\\"written: successful committed ailang-write/edit/delete calls, from input.path; produced: allowed non-failed effect result worktree_delta.changes; deltas_truncated and delta_errors name evidence gaps; local relative paths only, each path list ordered by first entry and capped at 1000; excess counts and unlisted are explicit\\\"}}\\n\"}\n",
diff --git a/host/verbs/lookup.go b/host/verbs/lookup.go
index c9a82e37d..31cff3be4 100644
--- a/host/verbs/lookup.go
+++ b/host/verbs/lookup.go
@@ -87,9 +87,12 @@ func Lookup(argv []string) (Verb, []string) {
 		switch argv[1] {
 		case "report":
 			return Verb{Words: "episode report", Class: Read, CLI: cli}, argv[2:]
-		case "start", "finish":
-			return Verb{Words: "episode " + argv[1], Class: Mutate, CLI: cli,
-				Reason: "episode " + argv[1] + " changes World. The app cannot run it yet."}, argv[2:]
+		case "start":
+			return Verb{Words: "episode start", Class: Mutate, CLI: cli,
+				Reason: "episode start changes World. The app cannot run it yet."}, argv[2:]
+		case "finish": // row 243 V3 §2.8: still terminal; the reason names the CLI's surface mode
+			return Verb{Words: "episode finish", Class: Mutate, CLI: cli,
+				Reason: "episode finish changes World. The app cannot run it yet; at a terminal, --confirm-on-surface lets you approve its pull request here with your key."}, argv[2:]
 		case "launch":
 			return Verb{Words: "episode launch", Class: Host, CLI: cli,
 				Reason: "episode launch needs a host process. The app cannot run it yet."}, argv[2:]
@@ -119,6 +122,7 @@ func Lookup(argv []string) (Verb, []string) {
 var (
 	mintRequestID = regexp.MustCompile(`^mr-[0-9a-f]{16}$`)
 	decisionRef   = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
+	verbRequestID = regexp.MustCompile(`^vr-[0-9a-f]{16}$`)
 )
 
 // ApproveTarget is approve's navigate arm: one request id or decision ref, no flag, opens its card.
@@ -129,5 +133,8 @@ func ApproveTarget(args []string) (href, reason string) {
 	if len(args) == 1 && decisionRef.MatchString(args[0]) {
 		return "/surface/decide/" + args[0], ""
 	}
+	if len(args) == 1 && verbRequestID.MatchString(args[0]) { // row 243 V3 §2.8: a finish confirmation's card
+		return "/surface/confirm/" + args[0], ""
+	}
 	return "", "the app opens approve <request-id|request-ref> as its card; any flag or other form runs at a terminal"
 }
diff --git a/host/verbs/lookup_test.go b/host/verbs/lookup_test.go
index 90ed58b6f..23cedc9de 100644
--- a/host/verbs/lookup_test.go
+++ b/host/verbs/lookup_test.go
@@ -29,8 +29,8 @@ var lookupCases = []struct {
 	{"episode report w1", "episode report", Read, "", true, []string{"w1"}},
 	{"ailang-worldd episode report w1 --json", "episode report", Read, "", true, []string{"w1", "--json"}},
 	{"approve mr-0123456789abcdef", "approve", Navigate, "", true, []string{"mr-0123456789abcdef"}},
-	{"episode start w1", "episode start", Mutate, "changes World", true, []string{"w1"}},
-	{"episode finish w1 --preset p", "episode finish", Mutate, "changes World", true, []string{"w1", "--preset", "p"}},
+	{"episode start w1", "episode start", Mutate, "episode start changes World. The app cannot run it yet.", true, []string{"w1"}},
+	{"episode finish w1 --preset p", "episode finish", Mutate, "episode finish changes World. The app cannot run it yet; at a terminal, --confirm-on-surface lets you approve its pull request here with your key.", true, []string{"w1", "--preset", "p"}},
 	{"episode launch w1", "episode launch", Host, "needs a host process", true, []string{"w1"}},
 	{"approver enroll --name x", "approver", Never, "approver is never run from the app: enrolling", true, []string{"enroll", "--name", "x"}},
 	{"session frob", "session", Never, "minting and revoking a session", true, []string{"frob"}},
@@ -85,7 +85,8 @@ func TestLookup(t *testing.T) {
 func TestApproveTarget(t *testing.T) {
 	ref := "sha256:" + strings.Repeat("ab", 32)
 	for args, want := range map[string]string{"mr-0123456789abcdef": "/surface/approve/mr-0123456789abcdef", ref: "/surface/decide/" + ref,
-		"mr-0123456789abcdef --deny": "", "--list": "", "": "", ref + "0": "", "sha256:" + strings.Repeat("AB", 32): "", "mr-0123456789abcdef x": "", "mr-0123456789ABCDEF": ""} {
+		"mr-0123456789abcdef --deny": "", "--list": "", "": "", ref + "0": "", "sha256:" + strings.Repeat("AB", 32): "", "mr-0123456789abcdef x": "", "mr-0123456789ABCDEF": "",
+		"vr-0123456789abcdef": "/surface/confirm/vr-0123456789abcdef", "vr-0123": "", "vr-0123456789abcdef0": "", "vr-0123456789ABCDEF": "", "xvr-0123456789abcdef": ""} {
 		href, reason := ApproveTarget(strings.Fields(args))
 		if href != want || (href == "") == (reason == "") {
 			t.Errorf("ApproveTarget(%q) = %q, %q", args, href, reason)
````

### Patch M14 — QUICKSTART §20 + CQ1

````diff M14
diff --git a/cmd/ailang-worldd/quickstart_finish_surface_test.go b/cmd/ailang-worldd/quickstart_finish_surface_test.go
new file mode 100644
index 000000000..23ad51428
--- /dev/null
+++ b/cmd/ailang-worldd/quickstart_finish_surface_test.go
@@ -0,0 +1,50 @@
+package main
+
+// Row 243 V3 sprint 5 M14 (design §6 CQ1; the S7 rule that QUICKSTART describes what the code does):
+// QUICKSTART §20 runs the flag the CLI parses, quotes the lines finish prints and the namespace the
+// card signs under, and states the wait the code waits.
+
+import (
+	"os"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+func TestQuickstartSection20MatchesTheFinishFlag(t *testing.T) {
+	raw, err := os.ReadFile("../../docs/QUICKSTART.md")
+	if err != nil {
+		t.Fatal(err)
+	}
+	guide := string(raw)
+	s19, s20 := strings.Index(guide, "\n## 19. "), strings.Index(guide, "\n## 20. Confirm a finish's pull request with your key *(attended — pending first verbatim run)*\n")
+	if s19 < 0 || s20 < s19 {
+		t.Fatalf("§20 is missing or not after §19 (%d, %d)", s19, s20)
+	}
+	section := guide[s20:]
+	_, block, _ := strings.Cut(section, "```bash\n")
+	block, _, _ = strings.Cut(block, "```")
+	argv, err := verbs.Split(strings.TrimSpace(block))
+	if err != nil || len(argv) < 3 || argv[0] != "ailang-worldd" || argv[1] != "episode" || argv[2] != "finish" || argv[len(argv)-1] != "--confirm-on-surface" {
+		t.Fatalf("§20's command %q (%v)", block, err)
+	}
+	if !strings.Contains(episodeFinishHelp, "--confirm-on-surface") {
+		t.Error("finish's help does not list the flag §20 runs")
+	}
+	src, err := os.ReadFile("episode_finish_surface.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	for _, quoted := range []string{"confirm request %s on the surface (expires %s, message hash %s)", "the typed y is not read in this mode", "confirmed by %s (decision %s)"} {
+		shown := strings.NewReplacer("%s on", "vr-<16 hex> on", "expires %s", "expires <time>", "hash %s", "hash <16 hex>", "by %s", "by <approver>", "decision %s", "decision sha256:…").Replace(quoted)
+		if !strings.Contains(string(src), quoted) || !strings.Contains(section, shown) {
+			t.Errorf("§20 does not quote finish's line %q as %q", quoted, shown)
+		}
+	}
+	if !strings.Contains(section, "-n "+approval.VerbNamespace) || approval.VerbRequestLifetime != 10*time.Minute || !strings.Contains(section, "up to 10 minutes (and 10 seconds' grace)") || verbLandGrace != 10*time.Second {
+		t.Error("§20's namespace, wait or grace is not the code's")
+	}
+}
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index dc1b5e104..7ce33dc75 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1911,3 +1911,15 @@ world shared@0ffbb9c
 ```
 
 The first prints the report exactly as `ailang-worldd episode report w1` prints it at a terminal, and `--json` prints its JSON; an episode with no committed call shows the CLI's own line and `exit 1`. `approve <request-id|request-ref>` opens that request's card (§17, §18), where the decision is your key's. `why head` (or any target the CLI takes: an entry index, `sha256:<ref>`, `a2a:`/`mcp:`/`rest:<id>`, with `--scan N` and `--json`) prints the chain exactly as `ailang-worldd why head` prints it, with `exit 3` under it when a link is broken and the CLI's own line with `exit 1` when the target is not found. The app walks under the daemon's 10 s read deadline (a terminal has 60 s): a walk that needs longer, or a read that fails, is refused with the CLI line to copy, never shown as not found or as a broken link; `why -`, `--result` and `--expect-store` read stdin, a file or another store and run only at a terminal. A refusal that names a CLI line opens the panel with the line to copy. Every other verb shows the exact CLI line and says the app cannot run it yet, with a copy button: `episode start` and `episode finish` today. Global flags, `--follow`, `--markdown` and `--addr` are refused, as are the commands that stay at a terminal (`approver`, `session`, `service`, `serve`, `setup`). A brane named like a CLI verb (`log`, `episode`, …) is reached as `world NAME`, which always navigates. If the daemon does not answer within 15 seconds, or answers something the page cannot read, the bar says so and keeps your line; it never falls back to navigation.
+
+## 20. Confirm a finish's pull request with your key *(attended — pending first verbatim run)*
+
+`episode finish` asks `Open this pull request? [y/N]` after it commits. With `--confirm-on-surface` you answer with an enrolled key instead (a §15 SSH approver or a §16 passkey) on the World app's card; the terminal is still required, and without the flag nothing changes. This guide uses a preset whose daemon serves at `127.0.0.1:7644`, a §15 approver and a §16 passkey `demo-mac`; never a real store.
+
+```bash
+ailang-worldd episode finish w1 --preset se-tools --confirm-on-surface
+```
+
+After the commit, finish shows the prompt's lines at the terminal (without the question) and prints `confirm request vr-<16 hex> on the surface (expires <time>, message hash <16 hex>)`, one `open …/surface/confirm/vr-<16 hex>` line per passkey origin and the daemon's loopback, and `the typed y is not read in this mode`. In the bar, `approve vr-<16 hex>` opens the same card (at a terminal, `ailang-worldd approve vr-<16 hex>` only points at it). The card shows the command, the same prompt and every byte your key signs; check the message hash, then **Approve with passkey**, or sign the saved message with the `ssh-keygen -Y sign -n world-verb@sunholo` line the card shows and paste the signature. Finish waits up to 10 minutes (and 10 seconds' grace), re-checks your signature over its own request, reads back the approval the daemon recorded, prints `confirmed by <approver> (decision sha256:…)` and opens the pull request as a typed `y` would.
+
+Each refusal keeps the commit and opens nothing: with `--dry-run`; a pinned unsandboxed check (it needs its own typed `y`: run finish without the flag); a busy or not-yet-upgraded credentials store (restart `serve` once); no confirmation in time; a newer finish of the same episode; a confirmation that does not verify. Ctrl-C stops waiting and keeps the commit. Run finish again: it resumes the same commit and body and asks afresh.
````
