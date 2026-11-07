# Independent evaluation — World iteration 362, row 256 (M1)

**Sprint**: `w-daemon-startup-under-load` (M1, single milestone; AC1–AC9)
**Reviewed commit**: `7655388936b80bea9c4843f0a2e0abc085e52c0f` (branch `sprint/world-iter362-row256`, draft PR 515)
**Implementation base**: `61a8d1c47b45fffc0f36137d29a5285c32da0145` · **Plan parent**: `7bd3f0fc501adc402b09e3755d96066463f2bc69`
**Evaluator worktree (isolated)**: `/Users/voightkampff/dev/sunholo-data/.eval-world-iter362-row256-r1`
**Verdict**: **FAIL** — score **85/100**, on a hard-fail condition (the operative local acceptance criterion is UNMET). The implementation itself is high quality; the blocker is a *measurement/instrument* obligation, not a demonstrated source defect.
**LOCAL_CRITERION: UNMET**

---

## 0. Protocol receipt

| Item | Value |
|---|---|
| `CLAUDE.md` read | yes, in full, in the evaluator worktree (`/Users/voightkampff/dev/sunholo-data/.eval-world-iter362-row256-r1/CLAUDE.md`) |
| Work classification | independent evaluation under the supplied approved design (`design_docs/planned/w-daemon-startup-under-load.md`, r3) and sprint plan (`design_docs/planned/sprint-plan-w-daemon-startup-under-load.md`) |
| `ailang messages` | **not run** — unattended mission stage, no human present (protocol step 2 explicitly skipped in that case) |
| `session_protocol_ack` | called with `{}`; result: `Session protocol acknowledged — mutating tools unlocked for this session.` → **acked = true** |
| Model identity | prescribed lane `pi:openrouter/minimax/minimax-m3`; the session env confirms the evaluator lane as the first entry of `MISSION_EVALUATOR_FALLBACK=pi:openrouter/minimax/minimax-m3,claude:claude-sonnet-4-6,opus`. No per-session model-identity variable is exposed to the agent; the supervisor's NDJSON remains the authority. |
| Judging done in-session | yes — no other agent requested, no other model's output inherited |

## 1. Scope of this evaluation

**What I did**

- Read the authoritative context end to end: `sprint-evaluator/SKILL.md` (full) + `resources/scoring_rubric.md`; shared `mission-control/resources/verification-protocol.md` rules 3d and 3e(a)/(b); repo `CLAUDE.md`; `design_docs/world-mission.md`, `design_docs/coding-standards.md`, `design_docs/DESIGN.md`; the approved r3 design (incl. §3–§5, §7 AC table, §8, §9, §11); the sprint plan (incl. the AC1–AC9 acceptance-ownership block and the gate table); `V/executor-state.json`, `V/controller-authorized-sprint-state.json`, `V/controller-plan-criterion-correction.md`, `V/controller-evaluation-brief.json`.
- Read the **complete production diff** (331 lines, written to `evaluator-r1/production-diff.txt`) rather than controller summaries, and inspected the commit content for omitted hunks: `git diff 61a8d1c4..76553889` touches 6 production files + 1 doc + 4 test files; the only non-source changes are the design/plan/verification artifacts committed by `caee00ca`/`7bd3f0fc`.
- Independently verified the frozen-manifest claim: all **11/11** files in `V/prototype-manifest.json` hash byte-exact against the worktree (recomputed SHA-256, not read from a receipt).
- Independently re-derived the caller census (grep) and re-ran the census test.
- Ran **2 targeted, compiling production mutations** of my own choosing (not a rerun of a controller report), each with control-green → build rc 0 → named-assertion red → byte-exact restore.
- Adjudicated the criterion correction and the disposition of every prior failure against rules 3d/3e.

**What I deliberately did not do** (per the iteration's constraints)

- No broad or full suites: no `go test ./...`, no `scripts/verify_go.sh`, no whole-repo CI rerun, no root-level test invocation.
- No commits, `git add`, `git stash`, `git reset`, `git checkout`; no PR edits; no lifecycle-document moves; no `ailang messages` or other messages; no unrelated source fixes; no real-log, secret or HOME-state access; no `TMPDIR`/`GOTMPDIR` override; no stress/GPU job.
- No mutation outside this isolated worktree. Source was mutated **only** in this worktree and restored from saved copies.

## 2. Implementation integrity (verified, not inherited)

| Check | Result | Method |
|---|---|---|
| 11/11 manifest hashes | **byte-exact** | recomputed SHA-256 of each file vs `V/prototype-manifest.json` |
| `git diff` vs HEAD | empty before and after my mutations | `git diff --stat` |
| Production diff scope | 6 production files + 1 doc + 4 new test files; no row-255 files (`failure.go`, `record.go`, `broker.go`, `handlers_exec.go`, `effectful.go`, `why.go`, `mcpclient.go`, `projection.go`) touched | full diff read |
| New named tests | **13/13 present** (`TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock`, `TestRegistryPhaseClockIsBounded`, `TestServeStartupClocksArePinned`, `TestRegistryPhasePreservesParentCancellation`, `TestExecExclusionLogIsObservationOnly`, `TestHealthExecExclusionIsAdditive`, `TestExecProbeOutOfTimeNamesTheClockNotAnArm`, `TestExecProbeGenuineRefusalRecordedAfterDeadline`, `TestDoctorNamesLeftOutExecProfiles`, `TestDoctorLateRefusalObservationAndReasonCap`, `TestDaemonConstructorCallerInventory`, `TestDaemonCallerInventoryAliasControl`) | read in source |
| Caller census (AC5) | **60 test `New` references, 0 production `New` references, 4 production startup edges** — exactly the design's corrected numbers; the 4 edges are `New→newDaemon`, `startDaemon→startBounded`, `startDaemon→newDaemon`, `Run→startDaemon` | my own `go test -run TestDaemonConstructorCallerInventory` output (`logs/control-C-callers.txt`) **and** my own independent grep |

## 3. Acceptance criteria — independent evidence

Legend: **MET** = evidence on the reviewed bytes; **UNMET** = obligation not discharged.

| AC | Requirement (abridged) | Verdict | Evidence I checked |
|---|---|---|---|
| **AC1** | Fresh, bounded 9 s registry phase derived from serve's parent; startup clock unchanged 9 s; successful Run lifecycle | **MET** | `daemon.go`: `startupBudget`/`registryPhaseBudget` = 9 s each (`TestServeStartupClocksArePinned`, controller run 23/23); `newDaemon` builds `regCtx = WithTimeout(phaseParent, registryBudget)`; `bootstrapRegistry`/`reportTemplates`/`scanIntegrity` all take `regCtx`; `Run` uses `startDaemon(ctx, cfg, startupBudget, registryPhaseBudget)`; `startBounded` unchanged. T1/T2/T3 pass **outside sandbox** (controller V28/23-named receipt). My own socket-free controls green. |
| **AC2** | All failed profiles excluded, incl. late genuine refusal; reasons retained; explicit `--exec-profile` still refuses startup; episode gets the named refusal | **MET** | `applyProbe` deletes every non-nil-result profile **unconditionally**; `OutOfTime` is written only into `st.Excluded`; `profileFor` refuses `excludedWhy` before reading `byProject`; flagged explicit profile still returns the probe refusal (`TestExecProbeOutOfTimeNamesTheClockNotAnArm`, last block). T4/T4b green in my worktree. |
| **AC3** | Observation-only log/health/doctor/help/QUICKSTART; genuine refusal reasons retained; review-before-restart advice; doctor names every exclusion as an exec failure (exit 13) incl. `--quick` | **MET (with a guard gap — see F-3)** | Shared constant `broker.ExecProbeDeadlineObservation` = "Every failed arm was recorded at or after the startup deadline; startup timeout may have contributed." — no causal claim. Log helper adds the observation only when flagged; doctor detail uses the observation + review advice and forbids causal phrasing; `doctorClass("exec profiles") → ("exec", 13)`; `doctor --help` lists the check; `docs/QUICKSTART.md` explains the flag as an observation ("establishes neither a cause nor a scheduling bound"). T5b asserts the absence of causal wording (green in my worktree). |
| **AC4** | `OutOfTime` true only for a **nonempty** failed-arm set **all** recorded at/after the deadline; mixed/no-failure/no-deadline false; late genuine refusal may be true; no scheduling bound | **MET** | `OutOfTime()` = loop requiring every arm `spent` **and** `len(e.Failed) > 0`; `fail()` sets `spent` only when `!deadline.IsZero() && !time.Now().Before(deadline)`; `pp.deadline` is taken from `ctx.Deadline()` (zero ⇒ never spent). I killed an independent mutant that breaks exactly this comparison (M-1 below). |
| **AC5** | Actual caller census + forward guard (imported/local aliases, direct calls, function values, shadow control) | **MET** | My own run: 626 Go files, 60 test `New` refs, 0 production `New` refs, 4 production edges; `TestDaemonCallerInventoryAliasControl` green (alias/direct/local/shadow/nested control); the census allowlist admits only the reviewed internal edges. *Minor note (F-4):* the counts 60/0/4 are logged, not asserted; the guard is enforced through the allowlist instead. |
| **AC6** | Additive health JSON (`omitempty` keeps old bytes; true adds `out_of_time`); doctor reason one line ≤240 bytes incl. the 3-byte ellipsis | **MET** | `HealthExecExclusion.OutOfTime bool \`json:"out_of_time,omitempty"\``; `doctorExecProfiles` uses `oneLine(x.Reason, 240-len("…"))` in **both** branches. I killed an independent mutant that drops the cap in the un-flagged branch (M-2 below). `TestHealthExecExclusionIsAdditive` green (controller receipt). |
| **AC7** | Parent cancellation reaches the registry phase with `context.Canceled`; settlement bounded; writer authority reopens; serial seams; narrow `-race` incl. T1 | **MET** | `newDaemon(ctx, phaseParent, …)`: `regCtx` derives from `phaseParent`, so a canceled serve parent cancels the fresh phase; `TestRegistryPhasePreservesParentCancellation` green in my worktree (socket-free) and in the controller's outside-sandbox race run. Clock vars are test seams restored via `t.Cleanup`. |
| **AC8** | Exact 11-file byte-exact delivery; four presence assertions; 13 named tests; docs/help preserved; focused gates; **controller full gates**; independent evaluation; PR **and merge** Linux plain/race CI | **UNMET** | Byte-exact delivery, 13 named tests, docs/help, focused gates and PR CI (`run 42344789832`: `go host build + test gate`, `go race gate`, `ailang-code verify gate` — all completed/success at the exact head) are satisfied. **The controller's full local Go gate, as corrected, is UNMET** (attempt 1: 2 permitted + 2 additional; attempt 2: 2 permitted + 1 additional). Merge-head CI does not exist yet (branch unmerged; controller-owned post-evaluation step). |
| **AC9** | Serve lifetime owned by Run's original parent, never the spent startup clock | **MET** | `Run`'s lifecycle binds `ctx` (serve's parent); the startup clock is confined to `startBounded`'s construction context; `startDaemon` passes serve's `ctx` as `phaseParent`. V29's `MUT-SERVE-ON-STARTUP-CLOCK` compiles and is killed by T1's health assertion (banked controller evidence, exact source). |

**Acceptance-criteria score**: 8 of 9 met → **26/30** (>50%, no hard fail on this category).

## 4. LOCAL_CRITERION

### **LOCAL_CRITERION: UNMET**

The operative local acceptance criterion is the controller's corrected criterion
(`V/controller-plan-criterion-correction.md`, authority rule 3e(a)): *run the identical full command on
the candidate with the identical inherited environment; require exactly the same two named failures, at the
same diagnostic-truncation layer, and no additional failures; a new/different failure or incomplete run
blocks; plus all Linux CI green.*

| Run | Failures beyond the two permitted | Criterion |
|---|---|---|
| Pristine base `61a8d1c4` | none (2 permitted only) | base red by design → AC repaired |
| Candidate attempt 1 (`V/controller-m1-verify-go-candidate.txt`) | `TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce`, `TestWorktreeSnapshotRefusesAConfiguredFilter` | **blocked** |
| Candidate attempt 2 (`V/controller-m1-verify-go-candidate-attempt2.txt`) | `TestExecProbeNet6ForeignDialIsRetriedNotRefused` | **blocked** |
| PR CI `42344789832` (exact head) | — | all 3 checks success (does **not** erase the local criterion) |

The two *permitted* failures are, on inspection of all three logs, **the same two named tests at the same
diagnostic-truncation layer** — the pinned interpreter's OpenTelemetry trace-exporter deprecation warning
precedes the diagnostic and the retained 400-byte output head truncates before the expected `LDR001`
(`failed …`), in `cmd/world-publish` and `host/archive`; the affected code (`host/archive/check.go`,
`cmd/world-publish/transitions_test.go`) is outside the patch. That half of the corrected criterion is satisfied.
The "no additional failures" half is not.

## 5. CRITERION_CORRECTION_ADJUDICATION

**Form: VALID.** Rule 3e(a) requires an acceptance command that is already red at base to be baselined on a
pristine tree, recorded as part of the criterion, and **fixed as an AC rather than by touching code or quietly
dropping the gate**. The controller did exactly that: pristine base rc 1 with exactly two failures, an
explicit, published correction document, the provider-pinned plan and its hash preserved, no change to the
pin, the warnings, the environment, the test policy or the production code, and the remote gate retained
("Require current PR and merge Linux CI … all present/completed/success"). I found no design-direction override
and no waiver of remote CI. On form, the correction is 3e(a)-compliant and admissible.

**Substance: SOUND BUT NOT ROBUST, and ONE ATTRIBUTION IS MISSING.**

1. *The clause "no additional failures" measures the repository's flake rate as much as the candidate.* The
   paired evidence in `V/` shows one broker test failing **4/100 at pristine base** and **2/100 at the
   candidate** — i.e. the base is not deterministic, and a single full-suite run of ~37 packages can pick up
   an unrelated flake. As written, the repaired criterion is therefore not reliably satisfiable by a correct
   implementation. This is not a reason to ignore it; it is a reason to make the criterion measure the change.
2. *Rule 3e(a) also requires the base to be pinned for each comparison.* Two of the three additional failures
   were dispositioned with paired base/candidate arms (the protocol the controller already used). **The third
   — `TestExecProbeNet6ForeignDialIsRetriedNotRefused` — has no base arm at all**, so its attribution is
   undetermined rather than excluded. That gap, not a proven regression, is what keeps the criterion unmet.
3. *Green CI does not discharge the local criterion.* The correction says so itself, and I agree: Linux CI
   supplies a genuine independent signal that the candidate does not break the Linux suite (plain full suite +
   race job at the exact head), but it does not answer the local comparison the criterion asks for.

**Recommended criterion repair (controller-owned, explicit, not a waiver).** Replace the absolute clause with
an attribution clause: *"no additional failure that is not reproduced at base"*, where each additional failure
must be attributed by a paired base/candidate arm (same test, same environment, N≥100, both arms, outside the
sandbox). A failure appearing in both arms is a base flake and is excluded; a failure appearing only in the
candidate arm is a candidate regression and blocks. This keeps the criterion's teeth (it still blocks a real
regression), removes base-flake noise, and reuses a protocol the controller has already executed twice. The
original correction and the provider-pinned plan stay on the record, as 3e(a) requires.

## 6. Disposition of every prior failure (rule 3d: mechanism, not timing)

| # | Test / package | Where | Disposition | Basis |
|---|---|---|---|---|
| 1 | `TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter` (`cmd/world-publish`) | base + both candidate runs | **BASE-INHERITED (permitted)** — identical failure text and layer in all three runs; pinned-interpreter trace warning consumes the 400-byte head before `LDR001`; code outside the patch | log diff of the three runs |
| 2 | `TestCheckSourceUnderPinnedInterpreter` (`host/archive`) | base + both candidate runs | **BASE-INHERITED (permitted)** — same mechanism and layer | log diff of the three runs |
| 3 | `TestWorktreeSnapshotRefusesAConfiguredFilter` (`host/broker`) | candidate attempt 1 | **BASE FLAKE — not attributable to M1.** The failing arm is the test's *own* positive control ("plain git status did not run the clean filter"). Paired arms: **pristine 4/100 fail, candidate 2/100 fail** — the mechanism is present at base at a *higher* rate (3d(iii): attribution must match on mechanism, not on timing; same test, same mechanism, both arms) | `V/paired-broker-pristine-100.txt`, `V/paired-broker-candidate-100.txt` |
| 4 | `TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce` (`host/broker`) | candidate attempt 1 | **UNATTRIBUTED (honest, no cause and no innocence claimed).** POST=2 vs expected 1. Not reproduced in 100 pristine + 100 candidate + 400 pristine observation-only runs. The patched broker files (`exec_probe.go`, `exec_startup.go`) contain no publish path, but "not reproduced" is not "excluded". I accept the controller's refusal to claim a cause | `V/paired-broker-*.txt`, `V/pristine-publish-request-observation-400.txt`, `V/publish-counter-observation.diff.txt`, `V/publish-observation-restoration.json` |
| 5 | `TestExecProbeNet6ForeignDialIsRetriedNotRefused/raw-twice-proxy` (`host/broker`) | candidate attempt 2 | **UNDETERMINED — attribution gap, not an established regression.** Message: "inconclusive after 3 attempts: every attempt's window held an accept that was not the child's (foreign 1, 2, 1)" — a foreign accept inside the listener's window, consistent with an interfering client. **No baseline arm was run**, so 3d(ii) applies: I may say "consistent with", not "caused by". On code inspection the patch touches only `ExecProbeArmFailure` bookkeeping, `OutOfTime()`, `execProber.deadline`, `fail()` and two lines in `probeExecProfile`; `net6Leg`/`judgeLeg`/`classifyNet6Dial`/`ExecProbeError.Error()` are untouched, so no patched path reaches this verdict. That is *consistent with* non-causation, but it is not the base-pinned measurement the criterion requires | attempt-2 log; source diff; absence of a paired arm |

Net effect: **no failure in this iteration is established as candidate-caused**, but one (row 5) is not
excluded by measurement either. Under the corrected criterion as written, all three additional failures block.

## 7. My own independent attack — 2 production mutations

Both mutations are inside M1's own production diff, compile, and are **distinct from every banked controller
mutant name** (19 historical + 14 revision + V29's serve-parent kill), which I read before choosing them.
Controller green was re-established **before** each edit; each file was restored from a saved copy and
re-hashed against the frozen manifest.

### M-1 — `MUT-EV-SPENT-IGNORES-RECORDING-TIME` (targets AC4 / AC3)

- **Production edit** (`host/broker/exec_probe.go`, `execProber.fail`):
  `spent := !pp.deadline.IsZero() && !time.Now().Before(pp.deadline)` → `spent := !pp.deadline.IsZero()`
- **Why it matters**: this drops the *at-or-after* comparison, so **every** failure recorded under a live
  deadline would be labelled "recorded at or after the startup deadline" — an unsupported timing claim that
  would push operators toward a false timeout diagnosis (the R3 direction). It is the opposite polarity of the
  banked `MUT-FAIL-NEVER-SPENT` (never marks) and is not `MUT-DEADLINE-UNSET` (never sets) nor `MUT-SPENT-ANY`
  (all-vs-any), all of which sit elsewhere.
- **Control before edit**: `go test -count=1 -run 'TestExecProbeOutOfTimeNamesTheClockNotAnArm|TestExecProbeGenuineRefusalRecordedAfterDeadline' ./host/broker/` → **rc 0** (`logs/control-A-broker.txt`)
- **Build**: `go build ./host/broker/` → **rc 0** (`logs/mutA-build.txt`)
- **Red**: `--- FAIL: TestExecProbeOutOfTimeNamesTheClockNotAnArm (0.00s)` — `exec_probe_out_of_time_test.go:25: an arm judged with time remaining is spent` (`logs/mutA-red.txt`)
- **Restore**: copied back from `mutations/exec_probe.go.orig`; SHA-256 back to `e10b1e59…b1` (equal to the frozen manifest hash); re-run → **rc 0** (`logs/mutA-restore-green.txt`)

### M-2 — `MUT-EV-DOCTOR-CAP-UNFLAGGED` (targets AC6)

- **Production edit** (`cmd/ailang-worldd/doctor.go`, `doctorExecProfiles`, non-flagged branch):
  `oneLine(x.Reason, 240-len("…"))` → `x.Reason`
- **Why it matters**: AC6 requires the *one-line, ≤240-byte* reason for **every** exclusion, not only the
  out-of-time one; the banked `MUT-REASON-CAP-240-PREFIX` only exercises the cap's *width* on the flagged
  branch, so the un-flagged branch was unguarded against losing the cap entirely.
- **Control before edit**: `go test -count=1 -run 'TestDoctorLateRefusalObservationAndReasonCap' ./cmd/ailang-worldd/` → **rc 0** (`logs/control-B-doctor.txt`)
- **Build**: `go build ./cmd/ailang-worldd/` → **rc 0** (`logs/mutB-build.txt`)
- **Red**: `--- FAIL: TestDoctorLateRefusalObservationAndReasonCap (0.00s)` — `doctor_exec_profiles_test.go:88: reason must be one line, capped at 240 bytes` (`logs/mutB-red.txt`)
- **Restore**: copied back from `mutations/doctor.go.orig`; SHA-256 back to `759ed7f1…b1` (equal to the frozen manifest hash); re-run → **rc 0** (`logs/mutB-restore-green.txt`)

**Non-vacuity**: both reds name a specific assertion, both are in the production diff of M1, both builds were
rc 0, and neither is a whole-source reversion (no compile failure was used as evidence). The mutation tests are
themselves **non-vacuous in the positive direction**: the same assertions are green on the restored bytes.

**Sandbox note**: `TestDoctorNamesLeftOutExecProfiles` (T5) **cannot be run in my lane** — it binds a loopback
listener and the sandbox denies it (`httptest: failed to listen on a port: listen tcp6 [::1]:0: bind: operation
not permitted`). That result is **UNINFORMATIVE UNDER SANDBOX** — it is neither a product pass nor a product
failure. The controller ran T5/T5b outside the sandbox (23/23 named passes), and I therefore rely on that
receipt for T5 and on my own socket-free run for T5b.

## 8. Findings

| Id | Severity | Finding | Evidence |
|---|---|---|---|
| **F-1** | **BLOCKER** | The corrected local acceptance criterion is **UNMET**: both candidate full runs contain failures beyond the two permitted base failures. | §4, §6 |
| **F-2** | major (instrument, not source) | The corrected criterion's "no additional failures" clause is not robust to the repository's measured base flake rate (4/100 at base for `TestWorktreeSnapshotRefusesAConfiguredFilter`), and the attribution protocol (paired base/candidate arms) was applied to only 2 of the 3 additional failures. | `V/paired-broker-pristine-100.txt` vs `paired-broker-candidate-100.txt`; absence of a paired arm for the net6 test |
| **F-3** | minor | **Guard gap:** no test asserts that `doctor --help` lists the new `exec profiles` check (AC3's "help updated" clause is unguarded — removing the help lines would not redden anything). The production bytes *do* contain the updated help (`doctor.go:57`, `:74`). I did not add a test: this iteration forbids unrelated test rewrites and the clause is not required by the criterion. | grep of `cmd/ailang-worldd/*_test.go` for "exec profiles" / help assertions; `TestDoctorHelpAndAPIKeyGuard` asserts only the usage line |
| **F-4** | minor | The AC5 census numbers (60 test / 0 production / 4 edges) are **logged, not asserted**; the guard is enforced by the production-edge allowlist instead, so a future `New` reference at an allowed (baseline) edge would pass the census silently. The four edges themselves are exactly right and the allowlist still rejects unreviewed edges. | `startup_callers_test.go:168-178` |
| **F-5** | informational | Residual (documented by the design, not hidden): the context-free archives and host scheduling still have **no established total startup bound**, so the 30 s service health wait can still expire; the 9 s + 9 s allowance bounds context-aware work only. The design states this explicitly (R2/R5) and claims no latency guarantee. | design §3.2, §6 R2/R5 |

## 9. Score

| Category | Max | Awarded | Basis |
|---|---|---|---|
| Tests Pass | 20 | **10** | **HARD FAIL condition** — local full Go gate rc 1 and the corrected local criterion UNMET (2 permitted base failures + 3 additional, none established as candidate-caused). Offsetting credit: PR CI full plain + race + ailang gates success at the exact head; 23/23 named focused passes (controller, outside sandbox); my own socket-free controls green; 2/2 of my independent production mutations killed. A base-red gate is not scored 0 (rule 3e(a)), but it cannot be scored 20 while the operative criterion is unmet. |
| Lint Clean | 10 | **10** | `go build ./host/broker/` rc 0 and `go build ./cmd/ailang-worldd/` rc 0 (mine); controller `build`/`vet`/`gofmt` receipts clean for candidate and base; no `TODO`/`HACK`/`FIXME` in the new code. |
| Acceptance Criteria | 30 | **26** | 8 of 9 met (AC8's controller-full-gate leg unmet); >50%, so no hard fail on this category. |
| Code Quality | 15 | **14** | Surgical 11-file diff byte-identical to the approved frozen patch; small focused helpers; no TODO/FIXME; no World file-size rule is in force (`design_docs/coding-standards.md` has none) and the large files are pre-existing. −1 for the unguarded AC3 help clause (F-3). |
| Documentation | 15 | **15** | `docs/QUICKSTART.md` S7 paragraph (observation semantics, "establishes neither a cause nor a scheduling bound", review-then-restart advice); `doctor --help` updated; design r3 + sprint plan complete and consistent with the implementation. |
| Design Fidelity | 10 | **10** | All 11 files byte-exact to the approved r3 patch; every by-construction claim (two clocks from serve's parent, unconditional fail-closed delete, nonempty ALL-after rule, additive JSON, 240-byte one-line cap with ellipsis reservation, exec class 13, observation-only wording, 60/0/4 census) verified on the current bytes. |
| Regression Surface Coverage (conditional) | — | not triggered | No parser/lexer/types/codegen/eval/vm/effects/`cmd/ailang/exec.go` path touched. The design's Conflict Surface is filled in regardless. |
| Performance Verification (conditional) | — | not triggered | The approved r3 design claims **no** performance/wall-time goal and runs no load experiment; it explicitly records that no total startup bound is established (R2/R5). The row's measurable obligations are the two clocks and the diagnostics, not a latency target. |
| **Total** | **100** | **85** | **Hard fail present → automatic rejection regardless of score.** |

## 10. Remedy (in scope; no source change required by any finding)

The blocker is a **measurement** obligation, not a demonstrated source defect. Concrete, bounded, and controller-owned:

1. **Close the one attribution gap (required).** Run `TestExecProbeNet6ForeignDialIsRetriedNotRefused`
   paired — pristine base `61a8d1c4` and candidate `76553889`, same environment and load, N ≥ 100 per arm,
   **outside the sandbox** (socket-dependent). *Fails in both arms* ⇒ base flake, excluded under the
   attribution clause ⇒ the blocker dissolves. *Fails only in the candidate arm* ⇒ a genuine candidate
   regression in `host/broker`; that would change the verdict to a source blocker and require an executor fix
   in the patched probe paths. This is the same protocol already run for failures 3 and 4.
2. **Amend the criterion explicitly (controller-owned, auditable).** Replace "and no additional failures" with
   "and no additional failure that is not reproduced at base", recording the paired-arm attribution rule; keep
   the original correction document and the provider-pinned plan hash on the record as 3e(a) requires. This
   removes base-flake noise without weakening the criterion against real regressions.
3. **Re-measure the corrected local full-Go gate once** on the candidate head with the amended criterion, plus
   the retained requirement that all Linux CI checks stay green. **No production change is needed** unless
   step 1 finds a candidate-arm-only failure.
4. **Optional hardening, not a blocker (F-3/F-4):** add an assertion that `doctor --help` lists
   `exec profiles`, and pin the AC5 census numbers (60/0/4) as assertions rather than log lines. Both are
   one-line test changes inside the sprint's own test files; I did not make them because this iteration
   forbids unrelated test rewrites and neither is required by the criterion.
5. **Merge CI** remains a controller-owned post-evaluation obligation (AC8's merge leg); it cannot exist
   until the branch is merged.

**No human decision is needed** for any of the above (the design states there are no human design-freeze
choices); step 1 is a measurement, step 2 is an explicit criterion repair with the original preserved. If the
controller prefers not to re-measure this iteration, the correct outcome is to **park the row** with F-1 open
rather than to accept on CI green alone.

## 11. Limits of this evaluation (explicit)

- No broad/full suite was run by me; every local verdict in this report is either a focused named test, a build,
  or a source/log reading. The full-suite picture is the controller's and CI's, quoted as evidence, not
  re-derived.
- Socket-dependent tests are **UNINFORMATIVE UNDER SANDBOX** in my lane (loopback bind denied). I ran only
  socket-free tests locally; T5/T1/T2-style verdicts rest on the controller's outside-sandbox receipts, which I
  read but did not reproduce.
- The PR CI receipt (`V/pr-ci-76553889.json`) was read, not re-queried from GitHub.
- The attribution of failure 5 (net6) is *undetermined*, not *excluded*: I did not run a base arm for it, and
  no paired base arm exists in `V/`.
- The design's own documented residuals (no total startup wall-time bound; no causal claim for the timestamp
  flag; the retry loop deferred to a follow-on row) are accepted as declared residuals, not as hidden gaps.

## 12. Final integrity checks

| Check | Result |
|---|---|
| Reviewed HEAD | `7655388936b80bea9c4843f0a2e0abc085e52c0f` (branch `eval/world-iter362-row256-r1`, worktree pinned to the reviewed commit) |
| Working tree after evaluation | `git diff --stat` **empty**; `git status` shows only untracked `design_docs/verification/world-iter362/**` artifacts |
| All 11 frozen-manifest hashes after my mutations | **byte-exact** (recomputed) |
| Mutated files restored | `host/broker/exec_probe.go` → `e10b1e59…b1`; `cmd/ailang-worldd/doctor.go` → `759ed7f1…b1` (both equal to the manifest values); restores performed by `cp` from saved copies — **no `git checkout`/`reset`/`stash` used** |
| Post-restore control runs | green for both mutated areas (`logs/mutA-restore-green.txt`, `logs/mutB-restore-green.txt`) |
| Lingering child jobs | none (`ps` shows no `go test`/`go build`/`ailang`/`worldd` processes owned by this session) |
| Writes outside the evaluator worktree | none beyond ordinary Go build cache / runtime temp files |
| Artifacts written | only under `design_docs/verification/world-iter362/evaluator-r1/` (report, JSON, logs, mutation copies, restoration receipt) |
| Commit / PR / message / lifecycle-document actions | **none** — all such actions are the controller's |

**Machine-readable companion**: `design_docs/verification/world-iter362/evaluator-r1/report.json`
**Restoration receipt**: `design_docs/verification/world-iter362/evaluator-r1/restoration-receipt.txt`
**Logs**: `design_docs/verification/world-iter362/evaluator-r1/logs/`
**Mutation copies + hashes**: `design_docs/verification/world-iter362/evaluator-r1/mutations/`
