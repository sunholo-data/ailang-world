# w-ailang-semantic-share: keep World's laws in AILANG, and the v1.0.0 pin — r0

**Status:** r0, **attended 2026-11-28**. Not quorum-reviewed. The loop's pick-time designer
re-measures §1 and §4 at its base and runs the quorum as usual. Nothing here waives a gate.
**Queue row:** 275. **Rulings:** D-WORLD-126 = A (gate-primary and the teaching prompt to
v0.64.1), D-WORLD-127 = A (measure now, ratchet after S5), D-WORLD-128 = A (S1 next, ahead of row
243 V4 S9c; S2–S5 after row 243 V4), all RESOLVED attended 2026-11-28 — see §7.

## §0 What Mark asked (2026-11-28)

> *"we are looking at making sure ailang world is using ailang where it makes sense and not
> slipping to go"*

Then, on the measurement and recommendations: *"can you take care of this?"*

## §1 Problem (measured at `ffd70358`)

By the charter's own rule this is drift. Coding standard S3 says: *"A FALLING `.ail` share … is
drift against S3."* The share fell by more than half, and nothing noticed, because the S3 ledger
has one row (2026-08-04) and S3 is scored by the evaluator rather than checked by a machine.

| Measure (ledger repro: `wc -l`, non-test vs test `.go` under `host/ cmd/`) | 2026-08-04 | 2026-11-28 |
|---|---|---|
| `.ail` core (`world/`) | 434 | 696 |
| `.ail` checked docs (`design_docs/**`, all sketches) | 809 | 1,749 |
| `.ail` packages (unique; excludes the 725-line `world-core` copy of `world/`) | — | 4,096 |
| Go host, non-test (`host/ cmd/`) | 6,773 | 89,283 |
| Go tests (`host/ cmd/`) | 10,506 | 170,967 |
| `.ail` share, like for like (core + docs vs Go host) | 15.5% | **2.7%** |
| `.ail` share, all unique `.ail` (6,645) vs all non-test Go (91,359) | — | **6.8%** |

Activity since 2026-09-01: **2** commits touch `world/*.ail`, **19** touch package or sketch
`.ail`, **594** touch `host/` or `cmd/` Go.

**Landed after this measurement:** row 274 M1 (`8acbb393`) added `design_docs/sketches/appmanifest.ail`, 247 lines and 23 contracted laws, with `host/appmanifest/drift_test.go` reading the sketch. That is the pattern §2 asks for. The loop's designer re-measures §1 at its base.

The `.ail` that exists is in good shape. Every contract gives the same result on v0.41.0,
v0.54.15 and v0.64.1: 156 verified (70 in `world/` and the sketches, 16 in `world-core`, 70 in
`se-tools`) and 41 uncontracted. So the problem is not rotting AILANG. New decision logic has
been going into Go.

### §1.1 Line share is the wrong thing to steer by

S2 puts effects in Go on purpose. The daemon, the store, the sandbox, HTTP, crypto and the CLI
plumbing are Go and should stay Go. The audit (§4) puts the realistic law-shaped Go at
**3–5k lines** of 91k. Giving every one of those an `.ail` source would lift the line share to
roughly 9–10%, no higher. A line-share target would push the wrong things into AILANG.

**The measure this doc proposes instead: law coverage.** Keep an inventory of World's laws (pure
decisions over typed facts: refusal orders, state transitions, admission, validation,
canonical keys). For each one, record whether an `.ail` source with contracts exists and whether
a drift test that **parses** that source binds the Go copy. Line share stays in the ledger as a
secondary column.

## §2 What belongs where — the rule

| Class | What it is | Where it lives |
|---|---|---|
| **A — effects** | HTTP handlers, SQL and store I/O, process and sandbox spawning, filesystem, CDP/Chrome, crypto codecs (SSHSIG, CBOR/COSE, passkey), strict JSON/regex/UTF-8 decoding, CLI flags | Go. Never ported. |
| **B — laws** | A pure decision over already-decoded facts: refusal order and codes, state-machine transitions, admission and eligibility, validation rules, canonical key render/parse round trips, outcome classification | **`.ail` is the source of truth**, with `requires`/`ensures` (S1). Go executes a transcription under a drift test that parses the `.ail`. |
| **C — projections** | Rendering, layout, view models, sentences | An AILANG package where it is a reusable domain; Go otherwise. Lowest priority. |

How a law runs: **the transcription pattern**, as execargs, effectplan and decisionpacket already
do. AILANG holds the law and its proofs, and Go executes a copy that a drift test keeps
byte-for-byte in step. This doc does not propose running the interpreter on the host's hot path.
New *domains* (S3's "everything else") still land as packages that run in AILANG.

**The design-doc question.** S3 asks every addition to `world/` or `host/` "why is this not a
package?". This doc adds a second question for any new class-B function in Go: **"where is
its `.ail` law?"**. The answer is either a sketch path plus its drift test, or a reason it cannot
be written in the pinned verifier's fragment, with the upstream issue number.

## §3 Why it slipped

1. **Nothing measured it.** The ledger has one row. The watch-item in S3 (a kernel-surface
   manifest gate) was never built. The record step of each iteration doesn't measure the split.
2. **The pinned binary teaches an old fragment.** Agents brief themselves from the pinned
   binary's `ailang prompt`. The v0.41.0 prompt is headed *"AILANG v0.16.6"* and says, at line
   2189, *"What gets skipped: Recursive functions, higher-order functions, effectful functions,
   float parameters"*. The v0.54.15 prompt (tool pin) says the same at line 2243. The v0.64.1
   prompt (lines 2619–2627) verifies recursion and list HOFs when `requires` bounds the depth,
   along with float, `%`, nested record patterns and the `std/list`/`std/option` helpers.
3. **The gate makes a new law in the old fragment impossible to land.** `verify_ail.sh` asserts a
   hardcoded manifest of verified identities. A recursive law is SKIPPED on v0.41.0 (§5.2), so it
   cannot join the manifest. The honest move is then to write the law in Go.
4. **Adding one verified law costs six files.** `verify_ail.sh`'s "floor-raise coupling
   inventory" lists them: the module, the regenerated `world-core` copy, two constants, the ready-
   packet golden, `SELF_MOD_PUBLISH.md` and the manifest gate's control marker. Each is there for
   a reason, but together they make Go the cheaper path for a 30-line rule.
5. **Since September, the work has been mostly host surface.** Row 243 (the operator surface)
   is mostly class A by nature, which is fine. The laws that came with it, such as the
   operator-action lifecycle and the allowance rules, went into Go next to their I/O, because
   nothing asked "where is its `.ail` law?".

## §4 Audit: where Go has absorbed laws (at `ffd70358`)

Method: a Go AST pass classifying every non-test function by signature and body (A if it touches
ctx, http, sql, store, io, os, exec, net, time or IO-shaped methods; C if pure and render-named; B
otherwise), then the top candidates read by hand. Estimates are ±30%. The "decl" column (types,
consts, embedded assets) was not split.

| Package | Non-test lines | A | B | C | decl |
|---|---|---|---|---|---|
| `cmd/ailang-worldd` | 22,603 | ~15.5k | ~1.2k | ~0.6k | ~4.6k |
| `host/daemon` | 14,482 | ~8.9k | ~1.3k | ~0.8k | ~3.5k |
| `host/broker` | 13,238 | ~5.9k | ~3.0k (~0.6k already mirrored, ~0.5k profile parsers) | ~0.7k | ~3.6k |
| `host/store` | 9,516 | ~6.2k | ~0.7k, plus a state machine in SQL | ~0.3k | ~2.3k |
| `host/approval` | 3,561 | ~1.6k (of which ~1.3k crypto codecs: stay Go) | ~0.7k | ~0.3k | ~1.2k |
| `host/surfacepage` | 3,462 | <0.1k | ~0.8k (mostly layout, really C) | ~0.9k | ~1.7k |
| `host/coordinator` | 2,191 | ~0.9k | ~0.5k (~60% already mirrored) | ~0.1k | ~0.7k |
| `host/surfacefollow` | 2,136 | ~1.0k | ~0.5k | ~0.1k | ~0.5k |
| `cmd/world-publish` | 2,012 | ~1.0k | ~0.3k | <0.1k | ~0.6k |
| `host/surface` | 1,878 | <0.1k | ~1.05k (validators, git-output parsers) | ~0.2k | ~0.6k |

### §4.1 Laws that live only in Go (ranked)

1. **The operator-action lifecycle.** The states requested → confirmed → claimed → running or
   uncertain → succeeded, failed, expired or cancelled exist only as about 40 `state = '…'` and
   `SET state = '…'` literals across `host/store/operator_*.go`, with lazy expiry at
   `operator_actions.go:68`. No single artifact holds the transition table. This is the strongest
   candidate: a transition law with contracts, with Go keeping the SQL and a drift test checking
   that every SQL transition is one the law allows.
2. **Decision law D1–D10**: `host/broker/decision_law.go:63` `SubjectRefusal`, `:93`
   `DecisionRefusal`, `ReservedApprovalScope`, `DecisionWindowMax` (~90 lines). Step-coded refusal
   order. `decisionpacket.ail` covers the packet, not this decision.
3. **Allowance draw D1–D8**: `host/approval/allowance.go:369` `AdmitAllowanceDraw`, `:107`
   `Validate` (~130 lines). The D-WORLD-96 bounds: 30-day validity, 12 h TTL, 10 draws a day,
   mission-loop exclusion.
4. **Admission**: `host/approval/admit.go:44` `Admit` and `:92` `MintMatches` (~60 lines). The
   verdict order and the mint-row binding.
5. **Operator action validation O1–O13**: `host/approval/operator_action.go:118` (~68 lines).
6. **Request validation**: `host/approval/request.go:95` `MintRequest.Validate` and
   `host/approval/verb.go:69` `VerbRequest.Validate` (~117 lines).
7. **Approval scope grammar**: `host/broker/approve.go:468` `ParsePublishApprovalScope` with its
   formatter, and `:704` `ParseFinishApprovalScope`. The law is the round trip and the term order.
8. **Finish confirmation F1–F5**: `host/broker/finish_confirm.go:23` (~20 lines).
9. **Publish outcome**: `host/broker/registry_publish.go:659` `classifyPublisherResult` (~30 lines).
10. **Preset check policy**: `cmd/ailang-worldd/preset.go:204` `validateGraderChecks` and `:257`
    `pinPresetChecks` (~90 lines), and the lexical half of `episode_finish.go:583`.
11. **Launch and session policy**: `episode_launch.go:44` `launchMode`, `session_new.go:705`
    `sessionNewGrants`, `episode_finish.go:741` `finishOutcome` (~95 lines).
12. **Journal intent validation**: `host/store/journal.go:318`, `:349`, `:371` (~75 lines),
    partly covered by `storejournal.ail`.

### §4.2 Defects: Go that says it follows a law and doesn't, or isn't held to it

- **D1: a divergent copy of a proven law.** `design_docs/sketches/worlddapi.ail:23` proves
  `isLoopbackHost(h) == (h == "127.0.0.1" || h == "::1" || h == "localhost")`. The copy
  `host/broker/registry_publish.go:404` has the same name but also returns true for any host
  starting `127.`, so the loopback-only test constructor in `validatePublishOrigin` (`:387`)
  accepts plain `http` to `127.example.com`. The production constructor only over-refuses, so the
  exposure is in test-constructor paths. `host/daemon/daemon.go:211` matches the law. No test
  binds either copy to the sketch.
- **D2: drift tests that restate rows by hand.** `effectbroker.ail` ↔
  `host/broker/decide_test.go` `TestSketchRows`; `worlddapi` `isLoopbackHost`/`clampLimit` ↔
  `host/daemon` tests; `storejournal` `mayReportNotStarted`/`retryAllowed` ↔ `recover_test.go`.
  These test files copy the sketch's rows instead of parsing the sketch, so editing the sketch
  turns nothing red. The execargs, effectplan and decisionpacket drift tests parse the sketch
  and are the model.
- **D3: a law applied with no row test.** `storejournal.ail` (`receiptState`,
  `outcomeMatchesIntent`, `effectDispatchLawful`, `writableRefText`, …) is applied in
  `host/store/journal.go`, and `GetReceipt` (`:857`) says it "mirrors receiptState". The only test
  that names the sketch is the manifest gate, which checks that it compiles.
- **D4: a transcription with no drift test.** `host/coordinator/plan.go:99` says
  `planInvocation` "transcribes world/transitions.ail's laws"; no drift test binds it.
- **D5: one verdict rendered three ways.** `cmd/ailang-worldd/approve.go:122`
  `approveAdmission`, `host/surfacepage/approve_model.go:245` `AdmissionRefusal` and
  `host/daemon/surface_confirm.go:43` `confirmAdmission` each re-render `approval.Admit`'s
  verdicts, with slightly different bad-digest handling.
- **D6: overlapping kernels.** The sketches `transitions`, `worldkernel` and `logepoch` overlap
  `world/` (and its `world-core` copy); they read as predecessors and could go stale. Retire them
  or mark them as history.

## §5 The pin

### §5.1 State

| Role (`pins/ailang-releases.json`) | Release |
|---|---|
| interpreter, proof, frozen-compiler, gate-primary | v0.41.0 |
| tool | v0.54.15 |
| observed | v0.60.4, v0.41.0, v0.33.0, v0.30.0 |
| retired | v0.52.1 |
| `go.mod` (`serveapi/protocol` and friends, `internal/apiserver`) | v0.47.2 |

- Upstream, `v1.0.0-rc.1` is tagged (`58f045ba6`, 2026-11-28) but has no GitHub release and no
  binaries. "Build and Release" runs on `dev` failed repeatedly on 2026-11-28.
- The newest published release is **v0.64.1** (2026-11-26, commit `f43d12e`). Its
  `darwin.arm64` tarball's digest matches the release's `.sha256` file.
- World pins only released binaries with recorded digests, so the rc cannot be the pin.
- `.mcp.json`'s `ailang-docs` server serves *current* docs, so an agent can read about a feature
  the gate binary does not have. A third skew.

### §5.2 Probe: recursive list laws, on the published binaries

All three binaries were downloaded from the GitHub releases. The v0.41.0 and v0.54.15 digests
match the pins; v0.64.1 matches its `.sha256`.

| Probe (`[]`/`x :: rest` recursion) | v0.41.0 | v0.54.15 | v0.64.1 |
|---|---|---|---|
| `allPositive: [int] -> bool`, trivial ensures | SKIPPED (list patterns outside the fragment) | SKIPPED | **VERIFIED (bounded, depth 2)** |
| count over `[int]`, `result >= 0`, no `requires` | SKIPPED | SKIPPED | INCONCLUSIVE (`bounded_unrolling_inconclusive`) |
| same, `requires { length(xs) <= 2 }`, default depth 2 | SKIPPED | SKIPPED | INCONCLUSIVE |
| same, `requires { length(xs) <= 2 }`, `--verify-recursive-depth 3` | — | — | **VERIFIED** |
| `countAllowed: [string] -> int` (`x == "ok"`), `result >= 0`, no `requires` | SKIPPED | SKIPPED | **ERROR (`solver_unknown`)** |
| `countAllowed`, `requires { length(xs) <= 2 }`, depth 3 | — | — | **VERIFIED** |

Reading:
- v0.64.1 makes recursive list laws writable, as its prompt says, *if the depth covers the
  `requires` bound*. A `requires { length(xs) <= N }` needs depth > N.
- The string-element case returns `solver_unknown`, an ERROR that fails the file, where the int
  analogue returns INCONCLUSIVE. That inconsistency goes upstream (§8).
- Every one of these is a bounded proof, not an induction. A law over lists of any length still
  needs either a bound World can justify, or an inductive proof upstream doesn't have yet.

### §5.3 What moving the pin gives and costs

- **Gives:** new laws can use recursion, list HOFs, floats and `std/list`/`std/option` with real
  proofs. The teaching prompt stops telling agents these are skipped.
- **Costs nothing in existing proofs:** all 156 verify identically on v0.41.0 and v0.64.1 (§1).
- **Costs re-verification for the other roles.** gate-primary carries the replay goldens
  (re-verified byte-exact when it moved to v0.41.0). frozen-compiler carries `world-core`'s
  `contentHash`/`interfaceHash`. Interpreter is what `serve` runs. The tool role already moved on
  its own (row 139 D7), so moving one role at a time is precedented.

## §6 Plan

Sprints, each ≤ 1,400 changed lines, in order.

- **S1: measure and bind what exists (no new law).**
  - `scripts/semantic_share.sh`: the ledger repro, emitting JSON (line split plus law coverage
    against a checked-in law inventory, `design_docs/semantic-laws.json`, seeded from §4.1). The
    iteration record cites it.
  - Fix D1: the broker's `isLoopbackHost` follows the proven law, with a drift test parsing
    `worlddapi.ail`. If accepting `127.0.0.0/8` is wanted, change the law and prove it first.
  - Convert D2's hand-copied tests to parse the sketch. Add drift tests for D3 and D4.
  - Collapse D5 to one renderer of `approval.Admit` verdicts.
  - Mark or retire D6's predecessor sketches.
  - The "where is its `.ail` law?" question is already in S3 (landed with this filing, D-WORLD-127).
- **S2: the pin (per D-WORLD-126).** Move gate-primary and the teaching prompt to the ruled
  release: re-verify the replay goldens, swap `AILANG_BIN` in CI and the verify scripts, add the
  release row with digests. The other roles follow at v1.0.0.
- **S3: the operator-action lifecycle law** (§4.1 #1): `design_docs/sketches/operatoraction.ail`
  with the states, the allowed transitions, expiry and the refusal per illegal transition, with
  contracts. A drift test extracts every `SET state` transition from `host/store/operator_*.go` and
  checks it against the law.
- **S4: the approval laws** (§4.1 #2–#8): one `approvallaw.ail` (or a few) for the decision law,
  allowance draw, admission, operator-action validation, request validation, scope grammar and
  finish confirmation. Drift tests parse it.
- **S5: preset, launch and journal policy** (§4.1 #9–#12). Decide per item whether it is a sketch
  law or belongs in a package (presets are package-shaped).
- **After S1–S5, the ratchet (per D-WORLD-127).** A CI check that fails when a function in the law
  inventory loses its source or drift test, or when a new exported pure function appears in the
  policy packages (`host/approval`, `host/broker` decision files, `host/coordinator`) without an
  inventory entry or a written exemption.

S1 needs nothing. S2 needs D-WORLD-126. S3–S5 are easier after S2, but can proceed on v0.41.0
where a law fits its fragment (no list-pattern recursion).

## §7 Rulings (all RESOLVED attended 2026-11-28; ledger rows in the charter)

- **D-WORLD-126: which release does World pin next?**
  - **A:** move gate-primary and the teaching prompt to v0.64.1 now; the other roles move at
    v1.0.0.
  - **B:** wait for v1.0.0 and move every role together.
  - **C:** move every role to v0.64.1 now.
  - **Recommendation:** A. v1.0.0's date is unknown and its release build is failing.
  - **RESOLVED A** (Mark, attended 2026-11-28).
- **D-WORLD-127: how is S3 enforced?**
  - **A:** the measurement script plus the inventory, cited in every iteration record, and the
    "where is its `.ail` law?" question scored by the evaluator; the CI ratchet after S5.
  - **B:** the CI ratchet from S1.
  - **C:** measurement only.
  - **Recommendation:** A. A ratchet over an unseeded inventory would mostly produce exemptions.
  - **RESOLVED A** (Mark, attended 2026-11-28).
- **D-WORLD-128: where does row 275 sit?**
  - **A:** S1 next, ahead of row 243 V4 S9c (it is small and fixes a live divergence); S2 onward
    after row 243 V4.
  - **B:** the whole row after row 243 V4.
  - **C:** after rows 272–274.
  - **Recommendation:** A.
  - **RESOLVED A** (Mark, attended 2026-11-28). S2–S5's order against rows 272–274 is the next
    attended groom's call.

## §8 Upstream

Filed as **sunholo-data/ailang#2494**, with a mission-control note, per the frozen-core rule (no local workaround): the
`solver_unknown` ERROR for a string-element recursive count, where the int analogue is
INCONCLUSIVE (§5.2), with the reproducer, and a feature ask for inductive proofs of structural
list recursion (the bounded proofs cover only lists up to the depth). The `v1.0.0` release-build
failure is upstream's own and is not re-filed here.

## §9 Not in scope

- Rewriting effects, I/O, crypto codecs or decoders in AILANG.
- Running the interpreter on the host's hot path.
- A line-share target.
- Any change to `tools/launchd/*`, the shared skills, or a vendored AILANG.

## §10 Verification log (2026-11-28, at `ffd70358`)

| # | Claim | How measured |
|---|---|---|
| V1 | Line split (§1) | `git ls-files` + `wc -l` per path set; packages excluding `packages/world-core/world/*` |
| V2 | Commit activity since 2026-09-01 | `git log --since=2026-09-01 --format=%h -- <paths> \| wc -l` |
| V3 | 156 contracts identical across pins | `ailang verify --json` per module (`world/`, `design_docs/sketches/`, `packages/se-tools/se_tools/`) and `--package packages/world-core`, on each binary, isolated `AILANG_STATE_DIR` |
| V4 | Prompt text per pin | `ailang prompt` on each binary; line numbers as quoted |
| V5 | Probe (§5.2) | `ailang verify --json [--verify-recursive-depth N]` on scratch modules, published binaries |
| V6 | Binary digests | `shasum -a 256` of the release tarballs vs `pins/ailang-releases.json` and v0.64.1's `.sha256` |
| V7 | rc.1 has no release | `gh release view v1.0.0-rc.1` → not found; `git log -1 v1.0.0-rc.1` in the AILANG checkout |
| V8 | Audit (§4) | Go AST classifier plus hand reading; ±30% |
| V9 | D1 divergence | `git grep -n -A8 'func isLoopbackHost'` |
| V10 | Operator state literals | `git grep -oE "SET state *= *'[a-z_]+'" -- 'host/store/operator_*.go'` |
