# W-Marketing-Capture — reproducible, truthful release images

**Status:** BLOCKED discovery — human scope/order gate before planning; capture and seed mechanisms remain unproven
**Target:** World v1.0.0, queue row 150 / D-WORLD-64
**Priority:** P1 release · **Estimated:** NOT READY TO ESTIMATE; prerequisite and capture discovery remains open
**Created:** 2026-10-08 · **Base:** `94fc5ba4d017505849e2adab1e87c8aae8ad0bdc`
**Planner-Lane:** codex-ok after the gate closes
**Dependencies:** landed rows 148/149; production evidence consumption and decision-packet projection if full row-150 scope is retained; an attended session for normal commits

## Problem statement

The attended umbrella D8 correctly asks for reproducible release artifacts, but its seed assumes more semantic machinery than landed row 149 supplies. Row 149 deliberately shipped an approvals-chain pane and an entry→object graph. It did not ship a decision-packet producer or canonical object-grade resolution. Renaming the approvals pane a packet, constructing a fake PROVEN badge, or making `why` print a plain object's payload would misrepresent the release.

This focused design preserves every charter requirement and names the prerequisite gaps. It does **not** claim row 150 can land under a reduced definition. The sections below record prospective requirements and a candidate approach, not a useful deterministic capture recipe already proven. Planning/implementation is gated until the operator decides the scope question and the remaining technical discovery succeeds. No implementation was performed in this design session.

## Goals and non-goals

Produce a repeatable local build artifact: six fixed Chrome shots backed by normally committed world history, with a complete hash manifest and a real coordinator provenance walk. Two independent renders of the same frozen seeded history must produce byte-identical PNGs. Fresh seeds must have identical public object/log/world hashes; credential bytes and SQLite file bytes are deliberately excluded from this definition.

Meet the literal row-150 seed: goal, proposals, decision packet, evidence at every grade, recorded effect outcome, and provenance chain. Missing production support is a prerequisite, never a fabricated object or silently weaker success criterion. Re-run the capture after row 151 adds a real agent-generated view.

No new Go dependency, no chromedp, no kernel edits, no shared driver or verification-gate edits. No remote publication, no browser authentication, no actual network/publish effect as a marketing prop. No new goal composer or general UI protocol. No screenshot-only grade renderer bypass.

## Verification log — first-party at base

Commands ran in this worktree; code readings establish mechanisms rather than inferring them from images. Unexecuted planned assertions remain explicitly pending.

| ID | Premise / method | Result |
|---|---|---|
| V1 | `git rev-parse HEAD`; charter row 150 and D-WORLD-64 read | Base is `94fc5ba4…`; row 150 expressly requires NORMAL commits, every grade, packet, effect, `why`, six shots, byte-identical captures. Exception only changes queue placement. |
| V2 | Skill scaffold dual search | Implemented neural top 0.23, planned top 0.22; below duplicate thresholds. SimHash's 1.00 hits are unrelated sprint plans. Umbrella is deliberately refined here, not duplicated. |
| V3 | `world/types.ail:20–57,91–103,155–168`; pinned `ailang check world/types.ail` | Released pin v0.41.0, commit `24ee108`; check exits 0. Proposal has a string `goal`, evidence list and capability/effect descriptions; DecisionPacket type exists. Neither a type declaration nor arbitrary JSON establishes a host producer. `gradeOf(ProofReceipt)` is CLAIMED; all other current constructors resolve TESTED, ATTESTED or CLAIMED. |
| V4 | `rg -n 'DecisionPacket|decision-packet' host cmd --glob '*.go' --glob '!**/*_test.go'`; positive control same scope `approval-request/v1` | Packet producer search empty; approvals producer exists in broker. This confirms row149 V15 at this base. |
| V5 | `host/daemon/workbench.go:43,477,577–590`; `host/workbench/render.go:264–275` | Object handler always supplies GRADE UNAVAILABLE. Decisions use `broker.RecentApprovals`, not packets. Renderer's grade constructor accepts labels but is not host evidence authority. Screenshot code must not use it to invent PROVEN. |
| V6 | `host/evidence/types.go`, `validator.go:120–216`; implemented evidence-boundary design header and §3.2 | Authenticated, subject-bound validation library exists with sealed resolution; its implemented contract is expressly library-only/non-production. `rg --files host` finds codec/validator, not a production proof producer or renderer consumer. Caller-selected key fixtures cannot certify marketing data. |
| V7 | `cmd/ailang-worldd/why.go:651–664,677–785`; `go test ./cmd/ailang-worldd -run TestWhy -count=1` with 110s alarm | Focused tests pass, 0.455s package time. For ordinary REST transition objects `why` reports no record chain, checks object bytes and labels entry hash client-supplied. Full world/entry/record/input/plan/effect/output checks require a coordinator invocation. Therefore the umbrella's proposed bypass killer (`why` alone) is insufficient. |
| V8 | QUICKSTART §2, `session_new.go:1–25`, `host/authority/mint.go:15–76` | REST commits are session-gated. Production mint is attended and TTY-fenced; core Mint is a primitive, not permission for a demo verb to bypass the fence. Random credentials are private authority state, not deterministic public history. |
| V9 | `host/coordinator/coordinator.go:48–79,354`; `host/daemon/daemon.go:694` | Coordinator accepts injected logical `Now`; daemon passes wall-clock Unix time. An effectful live seed hashes Now into effect request bytes (`broker.go:382`) and downstream effect/record refs; a pure invocation record itself has no time field. Reusing one frozen store permits repeat captures but does not prove two fresh seeds deterministic. A production-safe deterministic composition remains an unverified technical premise, not a separate human clock ask. |
| V10 | `host/daemon/workbench.go:45–96,552–590`; row149 implemented doc V35 | Seven accepted keys; object+payload and from+entry are valid. Graph uses checked entry→transitionFn/interpreter/transitionRef edges. It does not project arbitrary packet evidence edges. There is no packet query key to use for the fixed shot. |
| V11 | `host/daemon/workbench_live_chrome_test.go:58–101` | Existing Chrome drill uses a separate profile, preserves HOME (relocated HOME failed on macOS), and bounds Chrome to 30s. Capture must reuse these operational constraints. |
| V12 | `go.mod:5–8`; `rg --files website`; `website/static/AGENTS.md` | Two direct Go dependencies; output directory is proposed new `website/static/img/v1/`. Existing static AGENTS describes agent daemon usage; image generation adds no new agent-authority behavior. |
| V13 | `design_docs/verification/world-row150-design/chrome-control.log`; two local headless screenshot calls, isolated profiles, 1280×900, each bounded by Python subprocess timeout 30s | Chrome 154.0.8037.98 wrote two byte-identical 16,271-byte PNGs (SHA256 `07a78d63d817d16ac801311e024031a0baeaf93fc56089db317cf2cd04c93576`), but BOTH processes timed out after writing. This is instrument-only evidence, not a successful capture command, not product repeatability. PNGs and source HTML retained beside log. Normal successful exit and reliable process cleanup remain PENDING. |


| V14 | `host/workbench/render.go:212–235,249`; `host/daemon/workbench.go:552–596`; test bodies `TestWorkbenchGraphWiring`, `TestWorkbenchDecisionsPane`, `TestWorkbenchLiveCursorCrossesGap` | `Render` executes html/template server-side. Template ranges emit live list, SVG nodes/edges, decisions and provenance. Handler populates Page before Render. Graph test positive home/selected-entry edges and negative object-page placeholder; decisions empty/denied controls; live positive newest10 and empty/gap controls all inspect raw HTTP response with no JS engine. This proves server-rendered regions, NOT seeded-browser readiness, packet content or six-shot repeatability. Focused raw-response tests executed under110s alarm: `go test ./host/daemon -run 'TestWorkbench(GraphWiring|DecisionsPane|LiveCursorCrossesGap)$' -count=1` PASS, package0.544s; empty/placeholder branches are the negative controls. No browser was involved. |
| V15 | `cmd/ailang-worldd/main.go` flags/dispatch; `host/daemon/daemon.go Config and coordinator construction:694`; middleware:56 | No existing demo/CLI logical-clock injection was found in flags/config/construction. Coordinator construction hardcodes Unix wall time; middleware session resolution independently uses wall time. Library `Config.Now` is not a demonstrated normal CLI mechanism. Deterministic composition retaining expiry remains PENDING. |
| V16 | `host/coordinator/plan.go:25–155`; `effectful.go:284–333`; `host/broker/broker.go:223,382`; `host/authority/mint.go`; read effect record and handler outputs | Public varying inputs include invocation/task ID, episode ID, source/interpreter/epoch refs, input/output bytes, plan bytes and ordered effect refs; entry index/previous hash cascade into world hash. Effect request bytes include Now, effect/scope/cost; budget/result bytes and handler outputs also vary. JSON fixed structs/maps have canonical ordering here, but this does not canonicalize arbitrary handler output. Private random credential and credential/grant expiry may affect admissibility and budget outcomes. Logical time also enters journal intents. Exact whole seed byte inventory is PENDING; time control alone is not claimed sufficient. |
| V17 | charter row26, row105 and D-WORLD-23/24/64 read | 26 is an OPEN ordered producer owner;105 is PARKED/open grade projection owner. D24 explicitly permits hand-authored library fixtures while shedding producer/end-to-end demonstration. D64 exception scope only148–151. Packet-production search V4 has no host owner; controller will file172, not claimed already present. |

No language-limitation claim is inferred from lack of a host producer. V3 checks the actual landed language file; no new AILANG snippet is authored. Negative premises in V4–V6 and V10 have their own scoped reads/searches. Baseline broad suites were deliberately not run for this documentation-only task.

## High-impact decisions and human gate

**H1 — full showcase scope and prerequisite queue order (human decision).** D-WORLD-23 requires this tranche to keep scope, weaken claims to proven facts, and name separate residual owners without asking to absorb their implementation. We apply it: proof creation belongs to OPEN row **26** `w-bounded-z3-report-producer`; production subject-bound grade display belongs to OPEN row **105** `w-workbench-grade-projection`; the absent packet producer/projection needs a new named owner **row 172**, to be filed by the controller. This document implements none of them. D-WORLD-24 permits hand-authored authenticated validator fixtures to demonstrate the library boundary; it explicitly sheds the producer and end-to-end claim. Such fixtures are legitimate synthetic **validator tests**, but cannot establish actual production proof evidence, normal committed packet production, or the workbench's every-grade showcase. Storing a synthetic JSON object does not upgrade this classification.

The further issue is release scope and queue authority: D-WORLD-64 grants the blocked-critical-path exception only to rows 148–151. It does not authorize silently promoting row26/105/172 under that exception. **A (recommended): retain the full showcase and ratify an extension ordered 26→105→172→150 when no critical-path row is routable. B: explicitly narrow row150 to honest unavailable grades and the actual approvals pane, deferring the full showcase to those owners.** B changes a named core deliverable and therefore needs attended ratification (D-WORLD-24's distinction), unlike merely refusing to absorb separately-owned prerequisite implementation. **Default: park150 immediately. Row151 does not bypass the lowest-open-row condition.** Controller owns ledger ID/allocation (D69), queue filing and ratification; this document edits none.

**H2 — fresh-seed determinism is an unverified technical premise, NOT a separate human ask.** The normal daemon has no verified CLI/config logical-clock injection. `Coordinator.Config.Now` exists at library construction; that fact alone does not give `demo seed` a normal-path injection point. AC3 is PENDING and cannot presently be discharged. Planner may design routine demo-local composition only after proving normal credential resolution, grant expiry, budgets and normal commit invariants remain intact. Neither clock injection nor excluding private session rows is proved sufficient for public hash equality. No hidden mint, forged record, post-hoc time rewrite or policy override is authorized. If a concrete proposed implementation actually changes ratified authority semantics, that specific change returns to human review; generic clock wiring is not automatically a permission question.

Quorum is required by the normal mission gate and freeze trigger. R1 was BLOCKED by all three present reviewers; sonnet absent for quota, not a pass. This revision accepts their factual objections and narrows its status/claims. The no-chromedp constraint remains because the charter explicitly requires it; a reviewer cannot reopen that requirement on its own.

### Gates remaining before any design freeze

- [ ] H1 explicitly resolved; full requirements either have production prerequisites or a recorded attended amendment.
- [ ] H2 normal-path mechanism and all varying public byte fields measured; no hidden mint or fixture rewrite.
- [ ] Production packet/evidence ingress and display paths measured with positive and forged-input controls.
- [ ] Fresh-seed determinism measured against approved time semantics, not asserted.
- [ ] Re-quorum passes this narrowed discovery classification; operational measurements still stay pending until actually executed.

## Candidate approach — prospective, not frozen or proven

### 1. Seed orchestration through existing normal paths

`demo seed` is CLI orchestration, not a new storage API. It accepts an explicit empty disposable database, an existing credential, pinned released interpreter, and outputs a public `demo-manifest.json`. Refuse a nonempty world or a store held by a daemon; never overwrite/reseed an existing database. Store/session inspection may read existing state, but history writes must pass through authenticated REST commit or the normal coordinator invocation pipeline. Store.Commit/PutObject must not be called by the seed orchestrator to fabricate history.

Use established commit-envelope construction, content hashes over exact bytes, contiguous entry indexes, real stored interpreter/source objects, observed-head CAS and checked links. Plain descriptive goal/proposal objects are clearly demo-domain descriptions; do not claim arbitrary JSON is a validated kernel Proposal. A packet must use the eventual measured producer and versioned schema, reference its proposal/request, use logical times, and resolve real stored evidence.

Generate evidence by running bounded checks, retaining actual reports, and using approved production validators. TESTED needs an actual pass/fail report and visible verdict; ATTESTED means a recorded compiler/effect statement, not human approval without a person; CLAIMED may use clearly labeled agent/demo commentary. PROVEN requires the trusted production proof path and subject binding. A missing grade is a failing prerequisite, never fallback PROVEN.

Execute a harmless workspace-local effect (for example reading a committed demo text file) through the real registered coordinator transition. Retain actual effect request/result, allowed/failed/denial and budget delta. Capture its real invocation chain with `why <index> --json`; require coordinator=true, overall OK, and all expected check identities. A fake invocation-record object or canned successful result is prohibited.

Manifest binds public world/head, ordered entry hashes, object semantic IDs and hashes, packet ref, evidence subjects/grades/report refs, invocation/effect refs, pinned interpreter hash, and schema/version. Credential, filesystem absolute paths and wall-clock capture time are absent from public deterministic bytes. Environment/version metadata goes in a separate diagnostics file.

### 2. Candidate system-Chrome capture approach (PENDING)

`scripts/capture_marketing.sh --db <seeded-db> --manifest <file> --out <directory> --chrome <system-binary>` starts a loopback-only daemon in a scratch directory, waits on bounded health checks, takes six shots, and tears down on every exit. Each Chrome process gets an isolated profile, original HOME, explicit 30-second timeout, fixed device scale 1, viewport, color scheme, locale, font environment, and no GPU-dependent animation. Missing Chrome, missing manifest fields or unavailable shot content is a loud error. No chromedp or package download.

| ID / filename | Target | viewport | Required visible content |
|---|---|---|---|
| workbench-light.png | `/workbench` | 1440×1100 light | selected world, nonempty live timeline, graph, truthful decisions |
| workbench-dark.png | same | 1440×1100 dark | identical world/data; dark tokens applied |
| packet.png | measured packet projection, ref from manifest | 1440×1100 light | actual packet/proposal + resolved evidence of every required grade |
| graph.png | `/workbench?from=N&entry=N` | 1440×1100 light | checked entry/object graph; do not advertise an evidence layer absent from its source |
| provenance.png | `/workbench?object=<invocationRef>&payload=true` | 1440×1400 light | real invocation object and checked stored links; machine `why` report accompanies image |
| phone.png | `/workbench` | 390×844 light | selected world and readable responsive content |

Packet URL is explicitly PENDING; do not invent a query key. The shotlist must name an implemented route only after H1's prerequisite projection exists. Provenance screenshot plus `why` JSON provide complementary evidence; a payload screenshot alone is not a full walk.

For reproducible stills, disable JavaScript via Chrome's documented existing `--blink-settings=scriptEnabled=false` usage, so wall-clock live-status text and animation cannot race the screenshot. These are static release captures of the live product; caption them accordingly. V14 verifies that the existing server templates emit those product regions without JavaScript. This is a code/server-response fact, not an attended seeded Chrome capture or proof that every shot has its required content. No capture-only production query or forged footer is proposed. Live behavior is separately established by row149's Chrome drill, never claimed proven by a still. Dark forcing and PNG repeatability on this recipe are acceptance measurements, currently pending. V13 shows this local Chrome may write a PNG without exiting before the deadline: capture discovery must establish a bounded successful system-Chrome invocation before this candidate becomes a recipe. No alternative kill-after-file protocol is specified or accepted here; timeout remains failure even if a PNG exists. Product Chrome capture and cleanup remain unmeasured.

Publish only after two complete isolated capture passes against the **same frozen store** have the exact six-file name set, PNG dimensions from shotlist, nonblank content, and byte-equal SHA256s. Do not sanitize PNG metadata or post-process differing pixels to force equality. Require PNG decode, >32 distinct RGB colors and ink occupancy between 0.5% and 95%; supplement with required DOM text/region checks before capture so a nonblank error page cannot pass. Pin Chrome version and OS/fonts in diagnostics; byte identity is promised within that recorded environment, not across arbitrary Chrome versions.

A third capture after a deliberate token change must differ, proving that the equality instrument reads image bytes. Atomic output promotion occurs only after all checks; partial runs do not replace shipped images. The committed output directory contains six PNGs plus public seed/capture manifests and captions. Website and README use only artifacts that passed this recipe.

## Prospective files — contingent census, not an executable sprint

- `cmd/ailang-worldd/main.go` — dispatch/help, ~10 LOC.
- `cmd/ailang-worldd/demo.go` — bounded orchestration, ~220 LOC; host operation, not a package policy or kernel growth.
- `cmd/ailang-worldd/demo_test.go` — authority, commit receipts, empty-store and repeat-seed tests, ~300 LOC.
- `scripts/capture_marketing.sh` — bounded daemon/Chrome lifecycle and promotion, ~150 LOC.
- `scripts/marketing_shots.json` — fixed six-shot contract, ~45 LOC.
- `scripts/check_marketing_capture.py` — stdlib PNG/manifest validation or reuse an already-installed image decoder; no implicit installation, ~150 LOC.
- `scripts/test_capture_marketing.py` — instrument/lifecycle adversarial tests, ~200 LOC.
- `docs/MARKETING_CAPTURE.md` — executed-verbatim seed/mint/capture recipe and limits, ~120 LOC.
- `docs/QUICKSTART.md` — help/runbook link, ~10 LOC.
- `README.md` — image and captions, ~15 LOC.
- `website/src/pages/index.js` — use approved release images, ~20 LOC.
- `website/static/img/v1/` — six PNGs, public manifests and captions; no credentials or database.

Prerequisite producer/projection file census is deliberately **not guessed**. H1 prerequisites need their own focused designs and verified file inventories. H2's routine wiring can be finalized by the planner once its authority invariants are verified. This row does not secretly modify `world/`, `host/evidence`, authority mint or daemon time behavior. Host orchestration belongs in CLI because it starts local processes and produces build artifacts; semantic demo domains belong in packages if they introduce laws/policies (S3).

## Prospective acceptance obligations and mutations — ALL PENDING

All rows below are prospective obligations, not a frozen implementation contract or results. AC3/AC7/AC8 have no demonstrated positive recipe; their mutations cannot presently be drilled. Every mutation must compile/run far enough to reach the named assertion; setup failures are not kills. Planner finalizes test identities after H1/H2 without weakening the assertions.

| AC | Assertion / evidence | Mutation that must make it red |
|---|---|---|
| AC1 normal path | Seed HTTP trace contains authenticated successful commits; coordinator invocation is real; unauthenticated seed refuses before history writes | M1 bypass HTTP with direct Store.Commit: trace assertion fails even if `why` still prints object |
| AC2 authority | No credential/no TTY does not mint or change store; valid supplied credential exercises normal grant checks | M2 internal Mint on missing credential: session/history count assertion fails |
| AC3 deterministic seed | Two new disposable stores with approved logical clock have identical ordered public manifest/log/object hashes, excluding private sessions | M3 use wall clock in one record: exact manifest comparison fails |
| AC4 honest grades | Actual reports resolve subject-bound TESTED/ATTESTED/CLAIMED/PROVEN; forged receipt, wrong subject, missing verifier report never resolve PROVEN; verdict FAIL remains visible | M4 raw receipt→PROVEN, M5 wrong subject accepted, M6 hide FAIL: separate grade/subject/verdict assertions fail |
| AC5 real packet | Packet from production producer links stored proposal, request and evidence, with its actual schema; approvals alone cannot satisfy | M7 substitute approval-request ref: packet semantic/type assertion fails |
| AC6 real effect/provenance | `why --json` overall OK + coordinator flag and required world/entry/record/input/plan/effect/output checks; effect result and budget delta equal recorded execution | M8 remove effect result, M9 plain-object `why` considered complete: required-link/coordinator assertions fail |
| AC7 complete shots | Exact six filenames + routes return 200 + required DOM content + dimensions/PNG/color/occupancy bounds; missing packet blocks | M10 omit phone, M11 capture nonblank 404 page, M12 empty graph: exact-set/status/content assertions fail |
| AC8 repeat images | Each of six files byte-identical across two isolated captures; deliberate stylesheet edit changes image hash | M13 comparison uses filename only: injected changed image passes incorrectly and adversarial assertion fails |
| AC9 bounded cleanup | Chrome timeout/nonzero exit and daemon failure produce nonzero exit, stop all spawned processes, leave previous output unchanged | M14 ignore Chrome failure, M15 skip trap cleanup: exit/output/process assertions fail |
| AC10 docs/site | Run documented commands verbatim after attended setup; README/site refer to exact existing image names and captions state static capture/version | M16 rename referenced image or remove required session option: link/recipe assertion fails |

Targeted affected-package/script checks during implementation; before landing, binding `verify_ail.sh`, Go build/tests and existing required checks must pass. Existing grade/CSP/grammar/mint-fence tests remain regression gates. No broad test run or mutant execution is claimed here. All tests passing and documentation updated are completion conditions, not permission to skip unsatisfied full-scope ACs.

## Conflict surface and deferred decisions

No parser/typechecker/compiler change; required language Conflict Surface does not apply. Operational conflicts: session mint fence (reuse, never override); normal REST authority (reuse); coordinator records and all per-run bytes (unverified technical composition; no separate clock ask); packet versus approvals (prerequisite, never rename); subject-bound proof authority (reuse production only); closed workbench grammar/CSP (reuse existing routes, prerequisite if new packet route); live JS footer (capture static JS-off, clearly caption); single writer (seed and serve sequentially).

Agent latitude after freeze: exact harmless effect, manifest serialization, scratch directory layout and bounds tighter than stated. No latitude to invent human approval, evidence grades, packet schemas, grants, a deterministic credential or a coordinator invocation. Missing Chrome means capture not measured, not skipped-green.

## Axiom compliance

| Axiom | Score | Reason |
|---|---|---|
| A1 Determinism | 0 | PENDING: no fresh-seed/product-capture equality established |
| A2 Replayability | +1 | Real checked coordinator chain and retained manifest |
| A3 Effect Legibility | +1 | Actual request/result/budget records, never simulated outcome |
| A4 Explicit Authority | 0 | Existing session fence retained; policy change requires H2 |
| A5 Bounded Verification | 0 | PENDING: Chrome control timed out; successful bounded product recipe absent |
| A6 Safe Concurrency | 0 | Single writer lifecycle reused |
| A7 Machines First | +1 | Manifest and reproducible commands |
| A8 Minimal Syntax | 0 | No language syntax |
| A9 Cost Visibility | 0 | No cost semantics changes |
| A10 Composability | +1 | Existing CLI/coordinator/workbench paths |
| A11 Structured Failure | +1 | Missing content/authority/browser fail loudly |
| A12 System Boundary | 0 | Build artifact orchestration outside semantic kernel |

Net +5 for the prospective direction, no hard negative. Capture-related A1/A5 are 0/PENDING. Scores do not certify an implementation or discharge H1/H2.

## Timeline, related documents, and review status

Documented source premises are checked; operational discovery is incomplete. There is no reliable product capture recipe or two-fresh-seed measurement, so no implementation schedule is promised. H1 ratification and technical discovery precede sprint estimation. Separate owners 26/105/172 keep their scope and estimates; their implementation is not hidden in row150.

Related: [attended umbrella D8/§6](w-world-live-surface.md) defines this requirement; [landed row149](../implemented/w-workbench-live-and-polish.md) supplies the narrower real surface; [library-only evidence boundary](../implemented/w-validated-proven-evidence-boundary.md) defines why fixtures cannot mint production proof; [coding standards](../coding-standards.md) binds S3/S6/S7. Search's high SimHash matches are not release-capture designs; no neural coverage gate fired.

**Quorum:** R1 BLOCKED (glm/kimi/gemini reject; sonnet absent quota), artifact `../verification/world-row150-design/quorum/w-marketing-capture-2026-10-08T13-05-07Z.json`; one bounded revision here, re-quorum owned by controller. **Human gate:** H1 open. **Planner premise gate:** H2 open; no separate human clock ask. **Implementation:** none. **Fresh full-world repeatability, six product PNGs, dark forcing, semantic packet and all-grade product projection:** PENDING, not measured.

## Controller gate record (iteration 246)

r1 BLOCKED 3/3 external (glm/kimi/gemini reject; sonnet absent quota); r2 BLOCKED 3/3 external (glm/kimi/gemini reject; sonnet absent quota). Artifacts retain every objection verbatim. No design freeze, sprint planning or implementation is authorized. H1 is parked as D-WORLD-69 OPEN in the charter; technical H2 remains unmeasured, not a separate human ask. Kimi r2’s no-chromedp absence premise was re-measured FALSE: `world-mission.md` row150 expressly says "No Go dependency (no chromedp)." Go `image/png` is the proposed future decoder per gemini; no Python dependency or implementation added. Other r2 byte-field/served-content/lifecycle requests remain pending for post-ratification discovery. This transcription does not resolve quorum objections or apply a carve-out.

**Measured shipping conflict (controller, after ADMIN CI):** unchanged `verify_go.sh` binary-artifact gate rejects all Git-binary blobs, including required future release PNGs. ADMIN control PNGs were moved to persistent run storage with hash manifest, not passed as product images. Fleet owns row173/ticket `inbox_1791466110794_90c512a9`; a fleet resolution is another prerequisite before any tracked release image lands. No local allowlist, SVG substitution or hidden gate relaxation. H1 scope/order remains D-WORLD-69 OPEN.
