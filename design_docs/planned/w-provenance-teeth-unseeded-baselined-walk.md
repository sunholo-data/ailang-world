# w-provenance-teeth-unseeded-baselined-walk — clause-5 comparison arm: UNSEEDED question, live-timed pre-World baseline, read-only walk

**Status**: Scoped revision r3 under attended D-WORLD-44=A; fresh full author-excluding quorum REQUIRED. Execution BLOCKED on unmeasured method-level confinement and D2–D5 authority. No approved protocol, walk, or credit claimed.
**Target**: World 1.0 bar, clause 5, queue row **114** (`design_docs/world-mission.md:1891`; regroom position 5, `:1723`)
**Priority**: P0; row 108 prerequisite LANDED iteration221 (current charter row 108).
**Estimated**: ~1 day per question, opportunistic (row text, `:1891`); this r3 authored under the ≤30-minute directive bound
**Dependencies**: row-92 artifact `design_docs/verification/w-1-0-value-demonstration.md`; rows 96/103/104/99/100 landed (`:1855,1861,1863,1869,1871`)
**Designer-Lane (r3 revision author)**: fresh native `gpt-6.1-sol`; preferred GLM5.3 rejected Unknown model, D-WORLD-48 fallback FLAGGED (controller directive). Inherited r2 author: `pi:ollama/kimi-k3:cloud`. Nothing is human-ratified unless cited.
**Planner-Lane**: codex-ok
**Planner-Lane rationale**: protocol-only documentation with no compiler, kernel, package or harness implementation; any later planning remains blocked by freeze/readiness gates. This field classifies the document for the existing derive script, not a routing-policy amendment.
**Base**: `origin/dev` `03ea641884327045ac59754d9fc25ed046066d7d` (HEAD == origin/dev, §16 V1). Historical line anchors in inherited sections identify iteration210, not current positions.
**Scope**: short PROTOCOL doc. No code. No harness changes. Row 108 implementation strictly out of scope; its landed state only discharges ordering.

**r3 scope (D-WORLD-44):** correct the World-only issued-to-verified timer, require prevention and proof of method-level read isolation using already-existing facilities, exclude ancestral diagnosis material, and replace shorthand premise logs. The inherited r1/r2 rejections remain audit evidence, not approvals. No new proxy, filter, recorder, harness, implementation, incident attempt, commit or push is authorized here. If existing enforcement cannot be demonstrated, stop at `NOT-RUN(ISOLATION)` without revealing the question. This design revision does not ratify D2–D5.

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
| D3 Isolation by ENFORCED file/process/network confinement plus GET-only method enforcement on the allowed channel; post-hoc audits are backstops | r1: no-credential does not prevent local mutation/session mint | **designer PROPOSAL** | high |
| D4 All candidates use a `.git`-free allowlisted export excluding diagnosis material regardless of ancestry, with mechanical decontamination proof | Their answers are in the record (P4); worktrees share the object DB (r1) | **designer PROPOSAL** | med |
| D5 Clause-5 ≤300 s uses ONLY World-issued-to-verified wall; baseline and whole-experiment wall separately recorded | Row-92's operative reading; idle time must not hide | **designer PROPOSAL** (timer follows clause text; no added baseline-time threshold) | high |
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

**PR — Isolation prerequisite (before question reveal).** For EACH arm, record the existing enforcement facility identity/version/configuration, enforced root/network/process/method policy, exact commands and raw results. Allowed read must succeed; forbidden repo/state/context reads, alternate network destinations and forbidden process/filesystem operations must fail. W has no direct access to the underlying daemon address or SQLite files. Its only channel MUST enforce GET to the pinned declared read-route cohort, denying all other methods and paths independently of credentials and server authorization. No new filter/proxy/harness may be built for this row.

Before reveal, within the exact W environment, positive-control a permitted GET and negative-control POST `/v1/commit`, POST `/a2a/`, POST `/mcp/`, PUT/PATCH/DELETE on a permitted read path, and a GET outside the cohort. Also demonstrate direct-daemon/bypass access fails. Use an existing sacrificial control target with a known-success mutation request: the same bytes/credentials reach its handler without confinement, but cannot reach it through confinement. Record enforcement provenance plus target-side no-dispatch evidence; a credential-related 401/403, malformed-payload 400, GET-route mux 405, or client rc alone DOES NOT prove method isolation. Controls never target the qualifying corpus; no secret bytes are published. A failed, absent or unverifiable leg yields `NOT-RUN(ISOLATION)` before reveal.

**Current blocking premise:** no such facility/configuration or live controls have been supplied or exercised by this author. Native daemon routing is not that facility (§16 V2/V3/V9). Therefore method confinement is **UNMEASURED**, not available-by-assumption. If an existing authorized facility cannot discharge these controls, retain the blocker and return to the controller; inventing a harness/workaround exceeds D44.

**P2 — Freeze and pins.** Recorded immutably before either arm starts: repo HEAD + `origin/dev` shas; question sha256; packet sha256s; spawn metadata; event/CI refs; `max_attempt_wall_s` (proposal: 3600, all-inclusive); **archival corpus snapshot** (copy of the declared cohort, file-by-file sha256 + byte sizes, taken BEFORE serve and never served); declared-cohort manifest (refs + semanticIds + entry-index bound, no payloads); logical **invariants** (head entry index, ref count, `commit_objects` count) used for the post-walk delta audit. The **served copy is separate and mutable**; its byte-level behaviour under serve+probes is UNKNOWN/PENDING (no byte-inert premise is asserted). Nothing pinned may be edited until verdict; archival shas re-verified at P7.

**P3 — Role packets (whitelist + explicit denials).**
- **B**: question bytes; `.git`-free original-source export selected from the pin by an explicit path allowlist; minimal authentic event artifacts hashed in pins; ordinary shell/grep toolbelt. DENIED: all World surfaces, `design_docs/`, all network, parent checkout/object database, fixes/diagnosis records and unreviewed archives. `git archive <pin>` alone is insufficient: exclude diagnosis-bearing files and content already present in the pinned tree. No bundle/history-clone/worktree alternative is permitted in r3. Controller records an exclusion manifest identifying paths, payloads, commit messages and diagnosis-bearing sources regardless of whether their SHAs are ancestors. Inspect the completed export and packet recursively, including archive members; diagnosis-token grep must have zero hits with an outside-environment positive control. Token scans complement content/manifest review; zero hits alone never proves semantic cleanliness. If authentic archaeology inputs cannot be retained without leaking diagnosis, leave candidate CONTAMINATED-PENDING; do not invent sanitized evidence or mutate source history.
- **W**: IDENTICAL question bytes; enforced GET-only channel to the served corpus's declared read-route cohort, including its recorded method/path policy. DENIED: direct daemon channel, corpus files, all non-GET methods, paths outside the cohort, writes, process spawning and other network destinations. Existing role tooling must satisfy these constraints with PR proof; lack of credentials is not enforcement.
- **V**: the arm's answer text + the E2 ground-truth source list only. V authors neither answer; a distinct fresh verifier per arm is acceptable; the ground truth must be the same.

**P4 — Baseline arm FIRST, live-timed.** PR passes → question issued → B does unrestricted archaeology on its export, verbatim timestamped, declares an answer → V_B verifies against E2 → B's record is SEALED (nothing from B or V_B reaches W).
**P5 — Walk arm SECOND.** Same PR → question issued → W answers from the served copy via declared-cohort reads; every fetched ref logged → V_W verifies → sealed. W never sees B/V_B output.
**P6 — Verification.** Each V checks its arm's answer against the SAME E2 sources, verbatim, timestamped; an answer confirmed by nothing outside its own arm is UNVERIFIED (row-92 §4.4). Question-byte equality re-verified.
**P7 — Accounting.** Pins re-verified; archival hashes compared; delta audit of the served copy against P2 invariants (any accessed object's `committedBy` entry index ≤ freeze head and provenance attributable to operation, else `INVALID(seeding)`); artifact written; failures preserved, never substituted.

### Timer rules (binding)

- **Answer duration (per arm)**: START at question issue → STOP at answer DECLARED. Setup after reveal (serve start, packet staging) is inside; PR/env build are before reveal and recorded separately.
- **Verification duration (per arm)**: logged by V in `verify-{b,w}` as `verification_duration_s`; B/W transcripts contain NO verification blocks (r1 fix).
- **World envelope (controller-recorded)**: question-issued(W) → independently verified-result(W), total continuous WALL including World setup after reveal, handoff/orchestration, verifier role setup and idle. **Only this envelope is compared with 300 s.** Baseline issued(B) → verified-result(B) is recorded separately without a 300 s threshold. Whole-experiment issued(B) → verified-result(W) is also recorded, never used as the clause-5 threshold. Per-arm active + verification subtotals remain separate; idle is not subtracted. Controller timestamps connect isolated transcripts; neither answering role logs verifier commands.
- **Bounds**: exceeding `max_attempt_wall_s` ⇒ `INVALID(clock)`; an aborted/incomplete attempt is never a success and is preserved.
- **Clock integrity**: declared durations vs transcript-derived durations must agree within 5 s; no retrospective timing — an uncaptured duration is UNMEASURED.

## §6 Verdicts, seeding audit, replay-vs-live

**Closed set (per arm)**, extending row 92: `PASS` (verified; for W additionally World envelope ≤300 s) · `FAIL-TIME` (W verified, World envelope >300 s) · `UNVERIFIED` · `UNANSWERABLE(evidence)` — only from a queried, declared-complete cohort · `NOT-RUN(ISOLATION|NO-CORPUS|TRANSPORT)` — labelled, never pass/fail · `INVALID(clock|contamination|seeding|pin)` — integrity failure, preserved, cannot qualify. Candidate state `CONTAMINATED-PENDING` is neither verdict nor pass.

**Integrity audit (per question), attested by V in writing:**
1. Archival snapshot sha256 set identical pre/post (P2); mismatch ⇒ `INVALID(pin)`.
2. Post-walk delta audit of served copy against P2 invariants + full object-access attribution; NO byte-identity demand on the mutable served db absent a measured byte-inert premise (that premise is PENDING and unmeasured).
3. PR transcript quoted: allow rc0, deny rc≠0, both arms, including W method/bypass controls; absent ⇒ the arm is `NOT-RUN(ISOLATION)`.
4. Exposure log shows no diagnosis-material exposure inside any role environment; context manifests clean (E6); question sha256 equality across arms.
5. W transcript shows no write/POST/spawn success (any ⇒ `INVALID(seeding)`).

**Replay vs live.** Qualifying evidence is LIVE. Replays, replayed receipts, and regression harnesses are audit material and NEVER qualifying evidence; row 92's four walks are retrieval-class and do not count (P1/P2).

## §7 Raw transcript artifact schema (`design_docs/verification/w-provenance-teeth-unseeded-baselined-walk/`)

- `index.md` — harvest log (taken + rejection log), completion accounting, comparison table, findings.
- `Q<n>-<slug>/`: `pins` (§P2 fields incl. `env_method{allowlisted-git-archive}` + exclusion manifest, `cohort{path,sha256,producer,completeness}`, `archival_sha256[]`, `invariants{}`, `diagnosis_sha[]`, `max_attempt_wall_s`, spawn/packet shas) · `packets/` (bytes as delivered) · `prereq-transcript` (PR controls, quoted rc) · `baseline-transcript` (`setup{start,stop}`, `question_issued`, `answer_declared`, `answer_duration_s`) · `walk-transcript` (same blocks + `object_access_log` + audit results) · `verify-b`, `verify-w` (`verification_duration_s`, ground-truth refs, verdict per answer) · `envelope` (`t_issue_b`, `t_verified_b`, `baseline_wall_s`, `t_issue_w`, `t_verified_w`, `world_wall_s`, `experiment_wall_s`, `subtotal_b_s`, `subtotal_w_s`), with controller timestamp source and ≤5 s cross-check · `exposure-log` · `verdict` (closed-set, both arms + disqualifiers + V attestation).

## §8 Completion accounting and candidate registry

Row-complete at **qualifying successes ≥ 3**, where qualifying requires ALL of: E1 REAL harvested question; PR passed and quoted; identical question bytes; unseeded integrity audit clean (§6); that question's baseline live-measured in-protocol; the unseeded World answer verified by an independent verifier against the SAME ground truth; World envelope ≤300 s. **Partial baselines (e.g. walk `NOT-RUN(NO-CORPUS)`) and partial attempts never count and never clear pending state.** `attempts_total ≥ qualifying_successes` always; every failure class is preserved in `index.md` forever — no deletion, no substitution, and no speedup statement without its same-question live baseline.

| # | Candidate | Status |
|---|---|---|
| C1 | iter-206 (PR #161 green local → red CI; fix `29e1336`) | **CONTAMINATED-PENDING** — diagnosis is at `:1891` (P4); historical candidate; runs only after PR + P3-B content/manifest decontamination proof |
| C2 | iter-207 (load-class full-suite red, 24/24 clean second reading) | **CONTAMINATED-PENDING** — same posture |
| C3 | iter-210 (LIVE): controller-reported — pristine source HEAD `9b2cfac` == `origin/dev` (locally corroborated, V1), full `verify_go` rc1 locally vs same-SHA remote CI 2/2 green; blind fresh baseline: first sandbox prereq failed rc-6 BEFORE question reveal (PR behaved as designed), permissive host control rc0, corrected sandbox passed allowed rc0 / denied rc1, baseline subsequently RUNNING; no diagnosis disclosed to roles | **HISTORICAL CANDIDATE — author snapshot then in progress**; NOT a measured verified result; World corpus UNKNOWN (E3 open); does not count |

## §9 Acceptance criteria

| AC | Criterion |
|---|---|
| AC1 | Harvest log with ≥1 rejection or explicit none-spotted line |
| AC2 | Identical question bytes to B/W; sha256 equality re-verified at P6 |
| AC3 | B env contains no `.git`, no `design_docs/`, no diagnosis-bearing material at any ancestry; export/packet manifest review plus decontamination grep 0 hits and quoted positive control |
| AC4 | PR controls and method/bypass prevention attribution quoted BEFORE question issue, both arms; else `NOT-RUN(ISOLATION)` |
| AC5 | Archival hashes pre==post; delta audit vs invariants clean; served-byte mutation alone is never `INVALID(pin)` (§P2 invariants govern) |
| AC6 | Answer durations stop at declaration; per-arm verification durations logged by V; World/baseline/experiment wall recorded separately; ≤300 s claims cite only World wall; cross-check ≤5 s |
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
| MU-7 | W attempts POST /v1/commit, /a2a/, /mcp/ or bypasses the read channel | PR attributable method/bypass denial required before reveal; any success ⇒ `INVALID(seeding)` |
| MU-8 | "Approximate" duration | no retrospective timing ⇒ UNMEASURED |
| MU-9 | Cite ≤300 s from answer-duration only (idle hidden) | AC6 World-only wall rule ⇒ `INVALID(clock)` |
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
| Honest World envelope >300 s | FAIL-TIME recorded; loosening the bound is an attended amendment, never a protocol edit |

## §15 Axiom Compliance

A1 +1 (pins/transcripts make verdicts re-derivable) · A2 +1 (frozen archival snapshot; replay class separated) · A3 +1 (evidence standard is operation-produced effect records; the gap must speak) · A4 +1 (authority NARROWED by control-tested confinement, never by credential absence or ambient grants) · A5 +1 (bounded envelope; bounded/keyset reads only) · A6 0 · A7 +1 (hashed packets, fixed schema, closed verdicts) · A8 0 · A9 +1 (durations, subtotals, and idle are recorded fields; nothing imputed) · A10 +1 (read-only over rows 96/103/104/99/100 + row-92 vocabulary) · A11 +1 (NOT-RUN/INVALID named and preserved) · A12 +1 (verification forced outside the answering arm). **Net +10** ✅ proceed to quorum (fresh full author-excluding quorum REQUIRED; execution blockers remain).

## §16 Verification Log — r3 current-base measurements

Executed 2026-10-02 in the iteration222 worktree. Exact commands, exit statuses and raw quoted outputs are in [designer premise transcript](../verification/world-iter222/designer.md). No shorthand command literals are used.

| ID | Measured premise / limit |
|---|---|
| V1 | HEAD and origin/dev both `03ea641884327045ac59754d9fc25ed046066d7d`. |
| V2 | Handler registers GET readers alongside POST commit/A2A/MCP; isProtected matches ONLY POST `/v1/commit`. Route roster is pinned per attempt, never inherited from old line counts. |
| V3 | Middleware calls ResolveContext for protected requests; denial renders session failure, success reaches mux. Unprotected requests pass straight through. This enforces credentials, not GET-only access. |
| V4 | Scoped recorder search returns `docs/QUICKSTART.md` and `scripts/verify_go.sh`; these are documentation/historical verification note, not an established mission capture declaration. Positive instrument control is the known QUICKSTART hit. External fleet/runtime roots not searched; inventory remains UNKNOWN. |
| V6a/V6b | UNMEASURED occurs 10 times; exact comparison rows show 20/13/18/10 s retrieval, baseline UNMEASURED, and text admits prior diagnosis seeding. |
| V7a/V7b | Reverse readers validate deadline/quarantine and use bounded reads. Existence confirmed from full named function bodies' entry paths. |
| V9a–c | MCP resolves a session and can call coordinator.Dispatch; A2A also calls Dispatch. POST projection routes cannot be treated as a read-only channel. |
| V10 | Narrow method-filter search hits only test messages (GET-only card and TARGET-only transport), not a usable filter proof. This is NOT exhaustive absence proof. Existing external enforcement availability is UNKNOWN. |
| V11a–c | Named historical fix `29e1336` is an accessible commit object but merge-base ancestry checks return 1 at BOTH current and inherited tips. Thus this revision does not repeat the reviewer's specific ancestry assertion as a measured fact. Exclusion applies to diagnosis content regardless of SHA ancestry; bundled history is withdrawn entirely. |

**PENDING/UNMEASURED:** exact scratch launch flags, external confinement configuration, method/bypass controls, role manifests, decontamination, every new arm duration, capture coverage and served-byte behaviour. No scratch daemon or qualifying attempt was launched. The future prerequisite must pin actual flags; current session-auth code does not substitute for that proof. No AILANG support claim or new diagnostic code is made.

## §17 References

Row 114 `:1891`; regroom `:1723`; D-WORLD-42/43 `:928/:929`; row-92 artifact + `walk-transcript.md`; row-92 design `design_docs/planned/w-prove-the-1-0-bar.md`; shipped rows 96/103/104/99/100 designs under `design_docs/implemented/`; skill template `design-doc-creator/resources/design_doc_structure.md`.

## §18 Future Work

Optional separately authorized capture-cohort row IF a store/producer is positively declared (gap finding; not built here) · a decontamination-proof harness for C1/C2 (mechanical, control-gated) · use of preserved artifacts for replay/regression only (never qualifying evidence).


## Revision provenance and readiness — iteration222

Inherited r2 and quorum-r1/r2 were read directly from git object `b702dad27cdb96c785e189f6afcc43b419985c41` without changing/checking out that branch. Their exact reviewer objections are preserved verbatim in the designer artifact. Former C3 RUNNING snapshot above is historical: later controller disposition reports roughly 83 s archaeology, narrow PASS84, STOP not transcribed, World NOT-RUN(NO-CORPUS), 0/3 paired. It establishes no qualifying result.

D-WORLD-44=A authorizes this one scoped revision and fresh full author-excluding quorum after row108; current charter records row108 LANDED iteration221. D2–D5 remain proposals; D44 is not a blanket ratification of roles/confinement or execution. Attended triggers 1 (freeze items) and 3 (timing/KPI semantics) require quorum. Exclude the actual r3 author vendor, OpenAI, not merely inherited Kimi; controller records any unavailable seats by name/reason and cannot silently treat old quorum or N−1 as fresh full approval.

**Readiness:** complete scoped protocol revision for review, blocked for execution. Existing method-level enforcement is UNMEASURED; D2–D5 are unfrozen; corpus UNKNOWN. Resolve authority through the controller's existing attended ledger, and confinement only through proof of an already-existing authorized facility. If infeasible, park explicitly rather than constructing a workaround. No implementation, walk, baseline, credit, commit, push, external message or harness edit is part of this deliverable.

**Actor separation:** the r3 designer and diagnosing controller have seen historical diagnoses and cannot serve as blinded B/W actors for those incidents. A fresh candidate does not erase prior context exposure; actor manifests and E6 remain mandatory.
