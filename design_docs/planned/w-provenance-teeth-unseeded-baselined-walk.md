# w-provenance-teeth-unseeded-baselined-walk — clause-5 comparison arm: UNSEEDED question, live-timed pre-World baseline, read-only walk

**Status**: **PARKED — needs-human-review** after author-excluding quorum r1/r2 BLOCKED (iteration210). D-WORLD-44 must authorize any further revision; no plan/implementation or approved protocol is claimed. Author text below is the preserved r2 snapshot.
**Target**: World 1.0 bar, clause 5, queue row **114** (`design_docs/world-mission.md:1891`; regroom position 5, `:1723`)
**Priority**: P0 (row 108 parked under D-WORLD-43 — `:929`, `:1879`)
**Estimated**: ~1 day per question, opportunistic (row text, `:1891`); this r2 authored under the ≤30-minute directive bound
**Dependencies**: row-92 artifact `design_docs/verification/w-1-0-value-demonstration.md`; rows 96/103/104/99/100 landed (`:1855,1861,1863,1869,1871`)
**Designer-Lane (sole author)**: `pi:ollama/kimi-k3:cloud` — design authorship only; nothing here is human-ratified unless explicitly cited
**Planner-Lane**: `codex:gpt-6.1-sol` — routing field only; no planning/authorship by that lane is claimed
**Base**: `origin/dev` `9b2cfac` (HEAD == origin/dev, §16 V1)
**Scope**: short PROTOCOL doc. No code. No harness changes. Row 108 strictly out of scope.

**r2 revision summary (binding changes from r1):** (1) B's environment is an original-source export WITHOUT `.git`; all git-worktree isolation claims are withdrawn (shared object DB lets B read post-pin commits — r1 objection, root-verified). (2) Capture inventory downgraded to UNKNOWN; a missing corpus is `NOT-RUN(NO-CORPUS)`, never `UNANSWERABLE`-by-absence. (3) Confinement is an enforced, pre-reveal, control-tested prerequisite; all "no credential ⇒ physically impossible" claims removed. (4) The diagnosis embargo is REMOVED — normal incident recording is never prohibited; isolation is achieved by packets+confinement plus an exposure log; D2–D5 relabelled designer PROPOSALS. (5) Author header corrected to Designer-Lane. (6) Completion requires REAL + live-measured baseline + unseeded World answer verified within a ≤300 s wall envelope + integrity, ≥3 times. (7) Timer split: per-arm answer duration, per-arm verification duration, and the controller's total-wall envelope are separate fields. (8) Frozen archival corpus snapshot is distinct from the mutable served copy; byte-inertness of serving is PENDING, not assumed. (9) Verification Log rebuilt with exact commands, quoted results, controls, and declared missing roots. (10) iter210 entered as live candidate C3 (partial baseline in progress; not a result).

---

## §1 Problem — row 92 proved SEEDED RETRIEVAL with UNMEASURED baselines; clause 5's comparison arm is missing

Clause 5 (`world-mission.md:1000-1006`): *"on ≥3 REAL 'why did X happen' questions (arising from actual operation, not synthetic), a provenance walk yields the verified answer in ≤5 minutes each, where the pre-World method was a grep/log archaeology session."* Row 92 landed 4/4 PASS at 10–20 s, but its artifact's own honest reading states the walks measure **retrieval of a recorded diagnosis**, and every pre-World baseline is `UNMEASURED (record states method, not cost)` (§16 V6). The attended review ratified the gap (`:1891`): *"ROW 92 PROVED RETRIEVAL, NOT DIAGNOSIS… the comparison arm is missing."*

**Premises (verified while authoring → §16; controller-attributed where marked):**

| # | Premise | Evidence |
|---|---|---|
| P1 | Row 92's four walked answers were hand-seeded by the sprint, then retrieved | `:1891`; row-92 artifact comparison rows (V6) |
| P2 | All four pre-World baselines are UNMEASURED; no speedup ratio may be claimed from them | row-92 artifact ×10, V6 |
| P3 | Standing authority is ACTIVE harvest: mission operation IS real operation; the next Gate-2 "why did X happen" is taken; baseline timed live FIRST | `:1723` (attended, Mark, 2026-09-28) |
| P4 | iter-206/207 candidate diagnoses (fix `29e1336`; load-sensitivity) are **already recorded in the row text itself** ⇒ retrospective attempts are CONTAMINATED until mechanically proven clean; a pin alone does NOT suffice | `:1891` (read directly, V5) |
| P5 | Machine read surfaces exist, but the roster **drifted since r1** (`daemon.go:713-729`, incl. projection routes r1's census missed): every attempt must pin a DECLARED COHORT; no global route-census claim is load-bearing here | V2, V3 |
| P6 | Reverse readers exist: `Store.ObjectReferences`, `Store.ObjectCommits`, workbench edge tables | V7 |
| P7 | **Capture inventory UNKNOWN**: a grep over this repo's `scripts/ tools/ docs/` finds no mission→World recorder, but the agent-runtime/fleet roots that would own one were NOT searched (declared missing roots, V4) — absence is asserted nowhere beyond that exact scope | V4 |
| P8 | Row-92 operational caveats for any walker: scratch db + non-default port; single-writer dance; `--addr` needs a scheme; registry-key refusal | row-92 artifact F-5 |
| P9 | `isProtected` gates `POST /v1/commit` only; GET routes are open (declared residual R1) — a routing fact, NOT a confinement mechanism | V3 |

**Unmeasured/PENDING:** every baseline cost and walk duration (protocol unrun); served-copy byte behaviour under serve+probes (PENDING — no byte-inert measurement exists, §6 step 2); capture coverage (UNKNOWN, P7); C1/C2 decontamination (PENDING mechanical proof, E4); C3's baseline result (RUNNING, controller-reported, §8).

## §2 The anti-vacuity words, re-armed

| Word | Forbids HERE |
|---|---|
| **UNSEEDED** | Any evidence object hand-created for the walk; the walker reads only what OPERATION produced, from a declared cohort (P7) |
| **BASELINED** | Any ≤300 s verdict without that question's own live-measured grep/log archaeology duration beside it; partial baselines are not value proof (§8) |
| **CONFINED** | Any isolation claim not backed by control-tested enforcement (allow rc0 / deny rc≠0, quoted, BEFORE question reveal) — capability-absence prose is not isolation (§5 PR) |

## §3 High-Impact Decisions

| Decision | Why High Impact | Chosen By | Change Cost |
|---|---|---|---|
| D1 Active harvest at Gate 2 (the next "why did X happen" becomes a row-114 question) | Attended regroom orders it | **human** (`:1723`) | high |
| D2 Three fresh isolated roles B (baseline), W (walker), V (verifier); none is the diagnosing controller | Row 92's single-executor shape enabled seeding | **designer PROPOSAL** — needs ratification | high |
| D3 Isolation by ENFORCED confinement (file+process+network, control-tested pre-reveal), never by credential absence | r1: no-credential does not prevent local mutation/session mint | **designer PROPOSAL** | high |
| D4 Retrospective candidates (C1/C2) run only from a `.git`-free export with a mechanical decontamination proof | Their answers are in the record (P4); worktrees share the object DB (r1) | **designer PROPOSAL** | med |
| D5 ≤300 s claims use the total-wall envelope (§5 Timer); speedup ratios only against same-question live baselines | Row-92's operative reading; idle time must not hide | **designer PROPOSAL** (clause text + row-92 plan) | high |
| D6 Zero code, zero harness edits; deliverable is protocol + artifacts | Row text: "the protocol, not new code" | **human** (`:1891`) | high |

### Design Freeze

- [x] D1 (human authority exists) — [ ] D2/D3/D4/D5 are designer proposals, unfrozen pending attended ratification — [x] D6

## §4 Deferred Decisions

Scratch db path/port/export location per run; transcript container format (§7 fields fixed); deny-manifest token list per candidate (controller at harvest, audited by V); whether a NO-CORPUS gap files its optional row same-iteration or next (≤1 later).

## §5 The protocol

**P0 — Candidate event.** At Gate 2 the controller writes ONE harvest-log line: question text, event ref, ground-truth candidates, and starts the attempt's **exposure log** (immutable timestamps: question recorded; any later exposure of diagnosis-bearing material inside a role environment). **There is no embargo**: the controller's normal diagnosis recording continues unimpeded; isolation is delivered by packet environments and confinement, not by what the record contains. Abandonment is bounded: any NOT-RUN/INVALID/aborted attempt is preserved, its packet environments and holding copies are released, and normal incident recording proceeds regardless.

**P1 — Eligibility preflight** (every leg quoted in the artifact):

| Leg | Test |
|---|---|
| E1 REAL | Arose from this operation; event ref exists; not invented to fill the quota (`:1891`) |
| E2 GROUND TRUTH | An answer-checking source exists that NEITHER arm produces (CI log, git history, surviving binary, external system); identical for both arms at P6 |
| E3 CORPUS (walk arm) | A corpus cohort must first be positively DECLARED: store path + archival sha256 + producer + explicit completeness scope. None declared ⇒ walk arm `NOT-RUN(NO-CORPUS)` (baseline unaffected; NOT `UNANSWERABLE`). A declared-complete cohort that is actually probed and cannot answer ⇒ `UNANSWERABLE(evidence)`. The capture gap is itself a finding and may file ONE optional future scoped row (`:1891`: "that is the finding") |
| E4 SOURCE ISOLATION | Retrospective candidates are CONTAMINATED by default (P4). They run only if the P3-B decontamination proof passes; else they stay `CONTAMINATED-PENDING` (a candidate state, never a silent pass) |
| E5 QUESTION FROZEN | Exact bytes + sha256; identical bytes to B and W, re-verified at P6 |
| E6 ROLES FRESH | B/W/V spawned fresh; per-role context manifests recorded; memory/context injection disabled or audited; diagnosis tokens in any manifest ⇒ `INVALID(contamination)` |

**PR — Isolation prerequisite (gates question reveal).** For EACH arm environment, BEFORE the question is issued: (a) **allow control** — a read inside the allowed root returns rc0; (b) **deny controls** — reads of named forbidden paths (parent repo checkout incl. its `.git`, `design_docs/verification/`, `design_docs/world-mission*.md`, `~/.ailang/state/`, quorum/sprint dirs) and, per arm, named methods (network egress for B; network except loopback scratch port + any write outside scratch for W) return rc≠0. All commands and rc values are quoted in `prereq-transcript`. Failure or absence ⇒ `NOT-RUN(ISOLATION)`; the question is never revealed to that environment.

**P2 — Freeze and pins.** Recorded immutably before either arm starts: repo HEAD + `origin/dev` shas; question sha256; packet sha256s; spawn metadata; event/CI refs; `max_attempt_wall_s` (proposal: 3600, all-inclusive); **archival corpus snapshot** (copy of the declared cohort, file-by-file sha256 + byte sizes, taken BEFORE serve and never served); declared-cohort manifest (refs + semanticIds + entry-index bound, no payloads); logical **invariants** (head entry index, ref count, `commit_objects` count) used for the post-walk delta audit. The **served copy is separate and mutable**; its byte-level behaviour under serve+probes is UNKNOWN/PENDING (no byte-inert premise is asserted). Nothing pinned may be edited until verdict; archival shas re-verified at P7.

**P3 — Role packets (whitelist + explicit denials).**
- **B**: question bytes; an original-source export at the pin — `git archive <pin> | tar -x` (NO `.git`; git worktrees are REJECTED as isolation — shared object DB admits post-pin reads, r1-verified); minimal authentic event artifacts (controller-captured gate/CI output, hashed in pins); ordinary toolbelt (shell/grep). DENIED: all World surfaces, all `design_docs/`, all network. Only permitted alternative: `git bundle` history truncated at the pin, THEN positive proofs post-clone — `git rev-list --all --count` equals the pin-time count, `git cat-file -e <diagnosis-sha>` fails for every diagnosis-bearing sha named in pins, and the same PR confinement envelope. For C1/C2 the decontamination control additionally greps the finished environment for diagnosis tokens (0 hits required, positive control quoted).
- **W**: IDENTICAL question bytes; the served corpus copy on a scratch port (P8 caveats); the declared-cohort route list. DENIED: writes outside scratch, process spawning, network except loopback scratch port; no session credential is issued (belt-and-suspenders, not the mechanism).
- **V**: the arm's answer text + the E2 ground-truth source list only. V authors neither answer; a distinct fresh verifier per arm is acceptable; the ground truth must be the same.

**P4 — Baseline arm FIRST, live-timed.** PR passes → question issued → B does unrestricted archaeology on its export, verbatim timestamped, declares an answer → V_B verifies against E2 → B's record is SEALED (nothing from B or V_B reaches W).
**P5 — Walk arm SECOND.** Same PR → question issued → W answers from the served copy via declared-cohort reads; every fetched ref logged → V_W verifies → sealed. W never sees B/V_B output.
**P6 — Verification.** Each V checks its arm's answer against the SAME E2 sources, verbatim, timestamped; an answer confirmed by nothing outside its own arm is UNVERIFIED (row-92 §4.4). Question-byte equality re-verified.
**P7 — Accounting.** Pins re-verified; archival hashes compared; delta audit of the served copy against P2 invariants (any accessed object's `committedBy` entry index ≤ freeze head and provenance attributable to operation, else `INVALID(seeding)`); artifact written; failures preserved, never substituted.

### Timer rules (binding)

- **Answer duration (per arm)**: START at question issue → STOP at answer DECLARED. Setup after reveal (serve start, packet staging) is inside; PR/env build are before reveal and recorded separately.
- **Verification duration (per arm)**: logged by V in `verify-{b,w}` as `verification_duration_s`; B/W transcripts contain NO verification blocks (r1 fix).
- **Envelope (per attempt, controller-recorded)**: question-issued(B) → verified-result(last arm) total WALL, including setup-after-reveal, orchestration, role setup, and idle. **All ≤300 s claims use the envelope.** The answer-active + verification subtotal per arm is reported alongside; idle is never hidden.
- **Bounds**: exceeding `max_attempt_wall_s` ⇒ `INVALID(clock)`; an aborted/incomplete attempt is never a success and is preserved.
- **Clock integrity**: declared durations vs transcript-derived durations must agree within 5 s; no retrospective timing — an uncaptured duration is UNMEASURED.

## §6 Verdicts, seeding audit, replay-vs-live

**Closed set (per arm)**, extending row 92: `PASS` (verified AND envelope ≤300 s) · `FAIL-TIME` (verified, >300 s) · `UNVERIFIED` · `UNANSWERABLE(evidence)` — only from a queried, declared-complete cohort · `NOT-RUN(ISOLATION|NO-CORPUS|TRANSPORT)` — labelled, never pass/fail · `INVALID(clock|contamination|seeding|pin)` — integrity failure, preserved, cannot qualify. Candidate state `CONTAMINATED-PENDING` is neither verdict nor pass.

**Integrity audit (per question), attested by V in writing:**
1. Archival snapshot sha256 set identical pre/post (P2); mismatch ⇒ `INVALID(pin)`.
2. Post-walk delta audit of served copy against P2 invariants + full object-access attribution; NO byte-identity demand on the mutable served db absent a measured byte-inert premise (that premise is PENDING and unmeasured).
3. PR transcript quoted: allow rc0, deny rc≠0, both arms; absent ⇒ the arm is `NOT-RUN(ISOLATION)`.
4. Exposure log shows no diagnosis-material exposure inside any role environment; context manifests clean (E6); question sha256 equality across arms.
5. W transcript shows no write/POST/spawn success (any ⇒ `INVALID(seeding)`).

**Replay vs live.** Qualifying evidence is LIVE. Replays, replayed receipts, and regression harnesses are audit material and NEVER qualifying evidence; row 92's four walks are retrieval-class and do not count (P1/P2).

## §7 Raw transcript artifact schema (`design_docs/verification/w-provenance-teeth-unseeded-baselined-walk/`)

- `index.md` — harvest log (taken + rejection log), completion accounting, comparison table, findings.
- `Q<n>-<slug>/`: `pins` (§P2 fields incl. `env_method{git-archive|bundle}`, `cohort{path,sha256,producer,completeness}`, `archival_sha256[]`, `invariants{}`, `diagnosis_sha[]`, `max_attempt_wall_s`, spawn/packet shas) · `packets/` (bytes as delivered) · `prereq-transcript` (PR controls, quoted rc) · `baseline-transcript` (`setup{start,stop}`, `question_issued`, `answer_declared`, `answer_duration_s`) · `walk-transcript` (same blocks + `object_access_log` + audit results) · `verify-b`, `verify-w` (`verification_duration_s`, ground-truth refs, verdict per answer) · `envelope` (`t_issue_b`, `t_verified_last`, `wall_s`, `subtotal_b_s`, `subtotal_w_s`) · `exposure-log` · `verdict` (closed-set, both arms + disqualifiers + V attestation).

## §8 Completion accounting and candidate registry

Row-complete at **qualifying successes ≥ 3**, where qualifying requires ALL of: E1 REAL harvested question; PR passed and quoted; identical question bytes; unseeded integrity audit clean (§6); that question's baseline live-measured in-protocol; the unseeded World answer verified by an independent verifier against the SAME ground truth; envelope ≤300 s. **Partial baselines (e.g. walk `NOT-RUN(NO-CORPUS)`) and partial attempts never count and never clear pending state.** `attempts_total ≥ qualifying_successes` always; every failure class is preserved in `index.md` forever — no deletion, no substitution, and no speedup statement without its same-question live baseline.

| # | Candidate | Status |
|---|---|---|
| C1 | iter-206 (PR #161 green local → red CI; fix `29e1336`) | **CONTAMINATED-PENDING** — diagnosis is at `:1891` (P4); runs only after PR + P3-B decontamination proof |
| C2 | iter-207 (load-class full-suite red, 24/24 clean second reading) | **CONTAMINATED-PENDING** — same posture |
| C3 | iter-210 (LIVE): controller-reported — pristine source HEAD `9b2cfac` == `origin/dev` (locally corroborated, V1), full `verify_go` rc1 locally vs same-SHA remote CI 2/2 green; blind fresh baseline: first sandbox prereq failed rc-6 BEFORE question reveal (PR behaved as designed), permissive host control rc0, corrected sandbox passed allowed rc0 / denied rc1, baseline subsequently RUNNING; no diagnosis disclosed to roles | **CANDIDATE — partial baseline observation in progress**; NOT a measured verified result; World corpus UNKNOWN (E3 open); does not count |

## §9 Acceptance criteria

| AC | Criterion |
|---|---|
| AC1 | Harvest log with ≥1 rejection or explicit none-spotted line |
| AC2 | Identical question bytes to B/W; sha256 equality re-verified at P6 |
| AC3 | B env contains no `.git`, no `design_docs/`, no post-pin material; decontamination grep 0 hits with quoted positive control |
| AC4 | PR controls quoted (allow rc0, deny rc≠0) BEFORE question issue, both arms; else `NOT-RUN(ISOLATION)` |
| AC5 | Archival hashes pre==post; delta audit vs invariants clean; served-byte mutation alone is never `INVALID(pin)` (§P2 invariants govern) |
| AC6 | Answer durations stop at declaration; per-arm verification durations logged by V; envelope wall recorded; ≤300 s claims cite the envelope; cross-check ≤5 s |
| AC7 | Both answers verified against the SAME E2 ground truth; verifier independent; B/V_B sealed before W starts |
| AC8 | Verdicts from the §6 closed set; NOT-RUN/INVALID preserved as attempts |
| AC9 | ≥3 qualifying successes per §8; partials never count |
| AC10 | No speedup ratio without that question's live-measured baseline |
| AC11 | Wall order: B first, W second; no B/W transcript contains verification commands |
| AC12 | Artifact conforms to §7 with all fixed fields |
| AC13 | C1/C2 carried CONTAMINATED-PENDING until decontamination is mechanically demonstrated; C3 partial-not-counted |
| AC14 | Zero repo code/harness diffs |
| AC15 | No embargo: normal diagnosis recording never prohibited anywhere in this doc; isolation via packets+confinement+exposure log |
| AC16 | Missing corpus ⇒ `NOT-RUN(NO-CORPUS)`; `UNANSWERABLE(evidence)` only from an actually queried declared-complete cohort; no invented corpus, no seeded answer |

## §10 Adversarial controls (audit mutations)

| MU | Mutation | Expected kill |
|---|---|---|
| MU-1 | Plant hand-created object in served copy | delta-attribution fails ⇒ `INVALID(seeding)` |
| MU-2 | Change question bytes between arms | sha256 mismatch ⇒ `INVALID(contamination)` |
| MU-3 | Move start to first successful query | transcript contradiction ⇒ `INVALID(clock)` |
| MU-4 | Diagnosis-bearing file in a packet/env | AC3 grep + manifest ⇒ `INVALID(contamination)` |
| MU-5 | Re-run a recorded attempt as "live" | §6 replay rule; pins consumed ⇒ cannot qualify |
| MU-6 | Substitute an easier question after failure | candidate ids immutable ⇒ attempt `INVALID`, rejection logged |
| MU-7 | W writes/mints anyway | confinement deny rc quoted; any success ⇒ `INVALID(seeding)` |
| MU-8 | "Approximate" duration | no retrospective timing ⇒ UNMEASURED |
| MU-9 | Cite ≤300 s from answer-duration only (idle hidden) | AC6 envelope rule ⇒ `INVALID(clock)` |
| MU-10 | B reads post-pin objects via shared `.git`/external path | no `.git` in env + deny-control rc≠0 quoted (MU-4) ⇒ `INVALID(contamination)` |
| MU-11 | Issue question with failed/absent prereq | sequencing audit ⇒ `NOT-RUN(ISOLATION)`, preserved |
| MU-12 | Present serve/boot mutation of the served copy as seeding or as `INVALID(pin)` | archival/served separation + invariants (AC5) ⇒ correctly classed; only archival mismatch or invariant breach is `INVALID(pin)` |
| MU-13 | Backdate the exposure log | controller timestamps cross-checked against pins ⇒ `INVALID(contamination)` |

## §11 r1 objection disposition (binding map)

| r1 objection | Resolved by |
|---|---|
| gemini: B/W transcripts cannot carry V commands; cross-isolated clock impossible | Timer split (§5 Timer), schema §7, AC6/AC11, MU-9 |
| oc-glm: worktree shares object DB ⇒ C1/C2 degenerate to retrieval | §P3 env export, PR, AC3/AC4, MU-10, §8 candidate states |
| oc-glm: unproven byte-inert serve premise self-voids attempts | §P2 archival/served split + delta audit, AC5, MU-12 |
| oc-glm: embargo unbounded, blocks incident recording | embargo removed (§P0), AC15, bounded abandonment |
| controller: designer-labelled human decisions / ratified embargo | §3 relabelled PROPOSAL/human, AC15 |
| controller: absent corpus ≠ unanswerable complete corpus | E3, §6, AC16 |
| controller: completion must require measured baseline + unseeded integrity | §8, AC9/AC10 |
| controller: no-credential alone is not confinement | PR, §P3, AC4, MU-7 |
| controller: Verification Log lacks controls for absence claims | §16 controls + missing-root declaration, AC16 |

## §12 Conflict Surface

No parser/lexer/typechecker/codegen/eval positions; no code at all. Protocol-level surfaces, enumerated: the route cohort (read-only use; roster declared per attempt, P5) · the charter/log recording habit (**no change** — r1's embargo exception is retracted; diagnosis recording is never gated) · row-114 text (in every deny list; B/W environments provably exclude it) · Gate-2 STATUS stamp (untouched; protocol adds only a harvest line) · row-92 artifact vocabulary (extended, not redefined). **Must still work:** row-92's artifact reads unchanged; the mux route tests; the Gate-2 stamp. **Deliberately changes:** only this protocol's artifact layout.

## §13 Non-goals

No code/harness — including NOT building a capture recorder; the gap is a finding with an optional future scoped row · row 108 / D-WORLD-43 · no AILANG-language claims · no speedup claims without live baselines · no synthetic dry-runs as evidence · no git-worktree-based isolation anywhere in this protocol · no outer-loop recording-policy changes.

## §14 Risks & mitigations

| Risk | Mitigation |
|---|---|
| No declared corpus ever exists; walk arm stuck at NOT-RUN(NO-CORPUS) | That IS the anticipated first finding (`:1891`); baselines still get measured; ONE optional scoped row per gap |
| Retrospective contamination via the record (P4) | CONTAMINATED-PENDING default; PR + decontamination proof gate (AC3/AC4) |
| Confinement enforcement turns out unavailable (sandbox rc-6-class failures) | `NOT-RUN(ISOLATION)` is honest, preserved; C3's prereq sequence shows the control catching exactly this |
| Few real incidents; temptation to invent | E1 rejects invented questions; rejection log exposes scarcity |
| Honest envelope >300 s | FAIL-TIME recorded; loosening the bound is an attended amendment, never a protocol edit |

## §15 Axiom Compliance

A1 +1 (pins/transcripts make verdicts re-derivable) · A2 +1 (frozen archival snapshot; replay class separated) · A3 +1 (evidence standard is operation-produced effect records; the gap must speak) · A4 +1 (authority NARROWED by control-tested confinement, never by credential absence or ambient grants) · A5 +1 (bounded envelope; bounded/keyset reads only) · A6 0 · A7 +1 (hashed packets, fixed schema, closed verdicts) · A8 0 · A9 +1 (durations, subtotals, and idle are recorded fields; nothing imputed) · A10 +1 (read-only over rows 96/103/104/99/100 + row-92 vocabulary) · A11 +1 (NOT-RUN/INVALID named and preserved) · A12 +1 (verification forced outside the answering arm). **Net +10** ✅ proceed to quorum (r1 REJECT outstanding until re-run).

## §16 Verification Log (r2 — every load-bearing repo claim: exact command, quoted result, control)

Authored 2026-09-30, lane `pi:ollama/kimi-k3:cloud`, worktree on `origin/dev` `9b2cfac`. Doc-only: no `.ail`/Go diffs; no protocol execution is claimed from this seat. C3's CI states and sandbox rc values are **controller-reported in the authoring directive, not independently re-run here**.

| # | Claim | Command | Quoted result / control |
|---|---|---|---|
| V1 | Base HEAD == origin/dev | `git rev-parse HEAD; git rev-parse origin/dev` | both `9b2cfac20f770811b113d62e987530f7e6299660` |
| V2 | Route roster at base (P5) | `grep -n 'mux.Handle' host/daemon/daemon.go` | `:713-729` — GET /v1/health, /v1/head, /v1/worlds, /v1/objects, by-semantic-id, /v1/log/{index}, /v1/log, /v1/registry, POST /v1/commit, /v1/receipts, /workbench, `/.well-known/agent.json`, `POST /a2a/`; **r1's "ten-route 705-732" census had drifted and omitted the projection routes ⇒ cohort declared per attempt (P5)**. Positive control: the same grep returns the known `GET /v1/health` at :713 |
| V3 | Protected mutation surface (P9) | `sed -n '733,741p' host/daemon/daemon.go` | `isProtected` = `POST /v1/commit` only; comment: "The eight GET routes pass through unauthenticated as declared residual R1" |
| V4 | Capture-recorder search (P7) | `grep -rln 'ailang-worldd\|/v1/commit\|worldd commit' scripts/ tools/ docs/` | hits: `scripts/verify_go.sh` (:186, a historical note about a committed binary), `docs/QUICKSTART.md`; neither records mission operation into a World store. Positive control: `grep -l ailang-worldd docs/QUICKSTART.md` hits. **Declared missing roots:** agent-runtime/hook/fleet sources outside this repo (e.g. `~/.claude`, the fleet's launchd source) NOT searched ⇒ inventory UNKNOWN beyond this scope; absence not asserted |
| V5 | Row 114 text + recorded C1/C2 diagnoses (P3/P4) | `grep -n 'unseeded' design_docs/world-mission.md` | `:1723` regroom position 5 "ACTIVE, not passive (Mark, attended 2026-09-28)"; `:1891` full row incl. "Controller answer: a test-budget leak… fixed `29e1336`" and the C2 load answer ⇒ contamination confirmed by direct read |
| V6 | Row-92 seeded/UNMEASURED anchors (P1/P2) | `grep -c UNMEASURED …/w-1-0-value-demonstration.md`; `sed -n '40p;143,148p' …` | UNMEASURED ×10 (incl. `:54`); `:40` rejection-log header; comparison rows show baselines `UNMEASURED`, walks 20/13/18/10 s |
| V7 | Reverse readers (P6) | `grep -n 'func (s \*Store) ObjectReferences' host/store/object_references.go`; `…ObjectCommits' host/store/object_commits.go` | `:79`, `:26`; `host/daemon/workbench.go` present |
| V8 | Row-108 parked refs | `grep -n 'D-WORLD-42' design_docs/world-mission.md`; `grep -n 'D-WORLD-43' design_docs/world-mission.md` | first hit `:928`; first hit `:929` |

**Unmeasured here, by design:** all durations; served-copy byte-inertness (PENDING); capture coverage (UNKNOWN, V4 scope); C1/C2 decontamination (unrun); C3 baseline outcome (running, controller-reported). Nothing imputed; nothing retrospectively timed.

## §17 References

Row 114 `:1891`; regroom `:1723`; D-WORLD-42/43 `:928/:929`; row-92 artifact + `walk-transcript.md`; row-92 design `design_docs/planned/w-prove-the-1-0-bar.md`; shipped rows 96/103/104/99/100 designs under `design_docs/implemented/`; skill template `design-doc-creator/resources/design_doc_structure.md`.

## §18 Future Work

Optional scoped capture-cohort row IF a store/producer is positively declared (gap finding; not built here) · a decontamination-proof harness for C1/C2 (mechanical, control-gated) · use of preserved artifacts for replay/regression only (never qualifying evidence).


## Controller disposition — iteration210 (bookkeeping, not a designer revision)

Two rounds BLOCKED; one designer revision spent. Gemini and GLM external seats present both rounds; Kimi author vendor excluded. Sonnet absent due mission-ration quota measurement; GPT6.1 API seat unreachable429/nocredits. Subscription Agent/CLI transports were separate measured surfaces. No N−1 proceed claimed.

R2 blockers: V6/V7 command literals use shorthand ellipses, not exact reproducible commands; the allowed loopback channel does not enforce GET-only access, so the write/seeding prevention claim remains unsupported; the git-bundle alternative cannot exclude diagnosis commits already ancestral to its pin; combined baseline+walk300s envelope changes the clause5 World-only verified-walk bound; D2–D5 are proposals, not ratifications. These exceed the narrow verbatim-fix carve-out. No third design revision or sprint.

C3 historical RUNNING author snapshot is superseded only by separate observation: executor-reported≈83s archaeology, independently judged narrowPASS84, STOP stamp not in raw transcript. WalkNOT-RUN(NO-CORPUS), captureinventoryUNKNOWN, qualifyingpairedquestions0/3. No controller ruling replaces the independent judge.

D-WORLD-44: authorize one further scoped rotation-designer revision and full quorum (recommendedA), or defer row114 for attended regroom(B). Default: retain this draft unmerged, no walk/sprint under this rejected protocol; existing active harvest may still preserve partial real-incident observations without bar credit. No capture recorder or method-filter harness is authorized by this design park.

Quorum artifacts accompany this draft under `design_docs/verification/world-iter210-quorum/`; raw author sessions remain `/tmp/world-row114-kimi*.ndjson`.
