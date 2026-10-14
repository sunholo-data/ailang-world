# W-Marketing-Capture — reproducible, truthful release images

**Status:** REVISED r4 (iteration 266, controller bounded revision under the narrow-refinement carve-out after r3 BLOCKED). Colour forcing (V23), server-HTML stability and script inventory (V24) and real-content two-pass capture (V25) MEASURED by the controller. Product six-shot equality on the seeded demo store remains M7 acceptance, with a named residual path if it fails (M7 fallback).
**Target:** World v1.0.0, queue row 150 / D-WORLD-64
**Priority:** P1 release · **Estimated:** eight executor runs, each ≤30 min, plus attended release capture
**Created:** 2026-10-08 · **Revised:** iteration 266, 2026-10-14
**Base:** `c21f30310e04941afd3ccfb6b8748177db7b06c8` · branch `design/w-marketing-capture-r3`
**Planner-Lane:** codex-ok for the milestones below after design quorum/freeze
**Dependencies:** rows 148 (`38db793`+`9319e6b`), 149 (`b8b6a3f`), 26/105/172/176 and fleet PNG allowance landed; attended publish and mint for release data

## Problem statement

The release needs six reproducible images of real history. The previous design was
blocked on missing producers, wall-clock records, and an unsuccessful Chrome recipe.
The prerequisite code now exists, but it does not automatically populate a demo:
`NewProducer` has no production caller, no shipped skill requests `Human.Approve`,
and the packet page currently resolves only PROVEN evidence or an unavailable reason.
The seed must connect these paths without fabricating receipts, packets or grades.

Chrome 155 exits successfully without a screenshot when `scriptEnabled=false` is
passed. With JavaScript running, even an idle store's footer includes the current
time. Therefore neither Chrome's exit code nor store quiescence proves a stable image.

## Goals and non-goals

Discharge umbrella §6 150(a): normal commits and a coordinator provenance walk;
150(b): two byte-identical captures; 150(c): all six images present and nonblank.
Retain goal, proposals, real packet, all four evidence grades, and effect outcome.
Re-run the capture after row 151, incorporating its actual agent-generated view.

Here **deterministic demo world means repeatable captures of one frozen seeded
store**. Fresh-seed object/log/world hash equality is **NOT claimed**. This explicitly
revises the prior document's stronger sentence and the umbrella D8 shorthand;
it is a residual for Mark, not an implied amendment to the charter. No logical-clock
flag is added. Release acceptance uses the umbrella's literal 150(a)–(c), while
Mark retains ownership of whether fresh-seed equality needs a subsequent row.

No new Go dependency, chromedp, kernel change, shared-driver change, forged envelope,
fixture key, unattended mint, operator approval impersonation, or PNG post-processing.
No remote publish/deploy effect is needed. This revision changes design_docs only.

## Verification log — first-party at iteration 266

Code readings establish existing paths; future behavior is explicitly identified.
Executed results and exact controller commands are in
[revision notes](../verification/world-row150-design/iter266-revision-notes.md).

| ID | Command or file:line | Result |
|---|---|---|
| V1 | `git rev-parse HEAD`; `grep -n '^150\. ' design_docs/world-mission.md`; umbrella `w-world-live-surface.md:202,297` | Base c21f303; charter line 6130 requires full showcase, system Chrome, byte-identical captures and “No Go dependency (no chromedp).” Acceptance is 150(a)–(c). |
| V2 | charter D-WORLD-69/73 ledger `world-mission.md:955,959`; `git show -s --format='%h %s' 0390aa5 9d16a07 d36b148 8511bad 2d14dbd` | Landed: proof producer, PROVEN projection, packet producer/projection, daemon Human.Approve binding, PNG allowance respectively. Charter's old “not landed” annotation is stale; this doc does not edit it. |
| V3 | `host/evidence/proof_producer.go:94,277,472`; `rg -n 'NewProducer\(' host cmd --glob '*.go' --glob '!**/*_test.go'` | Definition only. GenerateProof reads stored source, stages one module, pins bytes/version, runs bounded ai-check, checks named identities, authenticates and writes its envelope. First production caller must be added here; no hand-built envelope. |
| V4 | `host/daemon/workbench_grade.go:138`; `host/evidence/validator.go:135,211` | Object grades currently resolve PROVEN only, through authenticated subject-bound ValidateProof→Resolve; other objects receive an unavailable reason. Existing renderer label constants confer no authority. |
| V5 | `host/coordinator/effectful.go:164,333`; `decisionpacket.go:174,191`; `host/daemon/daemon.go:788` | One allowed successful Human.Approve through normal plan→broker→commit mints a real packet; deadline is the selected grant expiry. Daemon binds it for every episode and serializes it. No operator decision is made. |
| V6 | `host/daemon/workbench.go:479`; `workbench_packet.go:60`; QUICKSTART §11 “Inspect a decision packet”; `daemon.go:925,926,930,939` | Route is `/workbench?object=<packetHash>`; listing is GET `/v1/objects/by-semantic-id/world/decision-packet/v1`. REST POST `/v1/commit` rejects that semantic ID. Packet v1 has two evidence lines, not a four-grade roster. |
| V7 | `packages/se-tools/transitions.json`; `rg -n 'Human.Approve' packages --glob '*.ail'`; `host/daemon/decisionpacket_daemon_test.go:24,40` | Existing skills ailang-read, ailang-check, ailang-cli use registered transitions. No production package requests Human.Approve. Test hu.approve uses a fake runner/interpreter: not a release skill. New package demo-approve must use the existing convention-v2 transition path, not this fixture. |
| V8 | `world/types.ail:20,44`; `host/broker/record.go:36`; `packages/se-tools/se_tools/check.ail:195,233`; `cli.ail:1` | Canonical legacy mapping: TestReport→TESTED, CompilerOutput/RecordedEffect→ATTESTED, AiReview→CLAIMED; ProofReceipt alone→CLAIMED. Real check/test results can come from Ailang.Check/Ailang.CLI; `host/daemon/workspace.go:44,290` requires the v0.52.1 tool pin. Lower-grade projection still needs implementing. |
| V9 | `cmd/ailang-worldd/why.go:574,639,690,700,764`; QUICKSTART §7 | `why <entry> --json` exposes coordinator, ok and named links world/entry/record/input/plan/effect/output. Plain REST objects are not a coordinator walk. Seed must assert the named set, not merely exit 0. |
| V10 | `docs/QUICKSTART.md:17`; `cmd/world-publish/transitions.go:57,89`; `cmd/ailang-worldd/session.go:211,225,250`; `host/authority/mint.go:24` | Publishing descriptors and production mint are attended, TTY-fenced, single-writer operations. Every grant expires with the session. Seed cannot call library Mint. |
| V11 | `cmd/ailang-worldd/session_test.go:29`; `host/daemon/daemon_test.go:1053`; `decisionpacket_daemon_test.go:69` | Test-only seams: testMintEnv, authority.Mint in mintSessionGrants, injected coordinator Now. Permitted for tests only, never the release seed. |
| V12 | `host/daemon/daemon.go:719`; `host/coordinator/coordinator.go:52`; `plan.go:85,122`; `host/broker/broker.go:382` | No CLI logical clock; effect requests hash Now. Task/episode IDs, pins, inputs, plans, outputs, budgets, grant-expiry packet deadline, approval refs and previous head also influence public bytes. Scope is frozen-store captures; fresh-seed equality remains a named residual. |
| V13 | `iter266-chrome-exit-probes.md` (controller); historical `chrome-control.log` | Chrome 155.0.8059.41 base recipe: rc=0 in 1–2s, PNG written, isolated runs byte-identical, IHDR/IDAT/IEND only. scriptEnabled=false: no file in 4/4 rc=0 runs; disable-javascript: script still runs. Supersedes Chrome 154 timeout as current instrument evidence, not product evidence. |
| V14 | `rg -n 'capture|script|Date|theme|prefers-color' host/workbench/render.go host/daemon/workbench*.go`; `host/workbench/live.js:54,109`; `render.go:181,192,254` | Server renders timeline/graph/decisions/packet; footer time comes from JS. No existing still-mode server switch found; seven query keys exclude capture/theme. Propose a serve flag suppressing script emission; do not invent an existing mechanism. Dark CSS uses prefers-color-scheme. |
| V15 | `go test -count=1 -timeout 45s ./host/daemon -run '^TestWorkbench(GraphWiring|DecisionsPane|LiveCursorCrossesGap)$'`; corresponding test bodies | PASS, 0.493s. Recorder-based server regions verified; browser pixels and request-time equality not established. Raw repeat-GET probe in revision notes is PENDING controller. |
| V16 | `go test -count=1 -timeout 45s ./host/coordinator -run '^TestDecisionPacketCodec$'` | PASS, 0.202s. Packet canonical byte shape and schedule refusals established. |
| V17 | pin `--version`, `test --format json world/logepoch.ail`, `ai-check -timeout 5s world/logepoch.ail` (40s supervisor) | v0.41.0/24ee108; test success=true, ten passes, no skips; sameRef and servesEntry verified with no errors/counterexamples/skips. Use this production standalone source, not row26 testdata. |
| V18 | pin `test --format json world/contracts.ail`; `ai-check -timeout 5s world/contracts.ail` (40s supervisor) | BASE FINDING: tests rc=1, four vacuous no_generator skips; ai-check rc=0, five verified. Not passing TESTED evidence; imports also prevent single-source producer staging. Full repo gate was not run or declared red. |
| V19 | `scripts/verify_go.sh:221`; `scripts/tracked_binary_gate.sh:33,64`; allowance command in revision notes | PASS at HEAD: zero binary blobs / 3516 files. Allowed PNG names/chunks/caps are described below. Gate reads committed blobs, not working-tree screenshots. |
| V20 | `go.mod:5`; Go standard packages image/png, hash/crc32; `website/static/AGENTS.md` read | Use Go standard-library PNG decoding and CRC checking; no Python image decoder or implicit installation. Static guide adds no capture authorization. |
| V21 | Exact dark probe and `capture_marketing.sh ... --verify-twice` in revision notes | PENDING (M7 acceptance): actual product file completion/content, two complete byte-equal passes. Freeze cannot claim these as observed. |
| V22 | Protected-path git-show/read comparison in revision notes | OBSERVATION: 12 tracked protected paths; HEAD/current aggregate SHA256 both 05852e6f08c102b42ab5952c4928a75779815512000a319de518d64e2a0f6202; line counts 7363/7363. No protected write. |
| V23 | Controller, iteration 266: the revision-notes dark probe, then a flag matrix (`iter266-colour-scheme-probes.md`), Chrome 155.0.8059.41 on a macOS host with `AppleInterfaceStyle=Dark` | The dark probe FAILED its LIGHT arm: with no flag headless Chrome reports `prefers-color-scheme: dark` (`true|#1f1e1b`), i.e. it follows the host OS appearance. `--force-dark-mode` was indistinguishable from that default. `--blink-settings=preferredColorScheme=1` → light (`false|#f7f6f2`) and `=0` → dark, and with `--screenshot` both exit rc=0 AND write the PNG (the two hashes differ; `=0` equals the `--force-dark-mode` bytes). So EVERY shot pins its scheme explicitly: light shots `preferredColorScheme=1`, the dark shot `preferredColorScheme=0`; `--force-dark-mode` is not used. Behaviour on a Light-OS host is unmeasured; the explicit setting is the host-independent choice, and M7's verifier still asserts the per-shot scheme from pixels. |
| V24 | Controller, iteration 266: `grep -n '<script\|<style\|nonce' host/workbench/render.go`; scratch probe `iter266-real-capture-probe_test.go.txt` (never committed) on a 7-entry store built by the existing test helper `workbenchLiveREST` (`host/daemon/workbench_graph_test.go:36`), two GETs 2 s apart per route, output `iter266-real-capture-probe-output.txt` | render.go emits exactly ONE script tag, `<script src="/workbench/live.js" defer></script>` (render.go:192), and ONE inline `<style>` (render.go:179); NO inline script, NO nonce. With that tag stripped, `/workbench`, `?object=<ref>`, `?object=<ref>&payload=1`, `?from=0&entry=1` and `?from=6&entry=6` returned byte-identical HTML across the two requests with 0 script tags: the complete varying-byte inventory of the server-rendered regions on this store is EMPTY. |
| V25 | Same probe: Chrome 155.0.8059.41, controller-proven base flags + `--force-device-scale-factor=1 --lang=en-US --blink-settings=preferredColorScheme=1|0`, two isolated-profile passes per shot through a loopback proxy that strips the V24 tag (the proposed still mode, simulated); hashes `iter266-real-capture-sha256.txt` | light 1440×1100, dark 1440×1100, object+payload 1440×1400, graph 1440×1100 and phone 390×844 of REAL server-rendered content (timeline, SVG graph, object payload) were each byte-identical across the two passes (sizes 63–137 KB). This is real-content evidence on a test-helper store, not the seeded demo store; M7 re-measures on that. `--screenshot` captures the VIEWPORT only: on 1440×1100 the light shot showed header and timeline, with the graph and decisions panes below the fold, so the shot list must choose routes/viewports that put each required region in view (M7 rectangles). |
| V26 | `host/daemon/workbench.go:303-304,421` | The workbench `payload` flag accepts only `0`/`1`; `payload=true` is `400 BadRequest — malformed payload flag` (measured in V24's first run). The REST `/v1/objects/{ref}?payload=true` form is a different endpoint and stays as written in S4. The provenance shot route is `?object=<ref>&payload=1`. |
| V27 | `host/daemon/workbench_packet.go:75-78`; `host/workbench/render.go:262`; QUICKSTART §11 "Inspect a decision packet" | A validated packet renders `approvals chain: <status>` from the broker's approval record, plus ` by <decider>` only once decided; QUICKSTART documents `approvals chain: pending` until an operator decides. This is the source of packet.png's "pending" content. |
| V28 | `world/types.ail:23-29,44-49` | Evidence kinds mapping to CLAIMED are `AiReview(HashRef, float)` and a raw `ProofReceipt`. An `ai_check` report is compiler output (`CompilerOutput` → ATTESTED), NOT an `AiReview`; the seed produces no AI review, so no real `AiReview` exists to name. |

## High-impact decisions

**H1 RESOLVED:** D-WORLD-69 = A retains the full showcase, ordered 26→105→172→150.
Prerequisites landed at 0390aa5/9d16a07/d36b148, daemon binding at 8511bad and fleet
PNG allowance at 2d14dbd (V2). The former scope/order question is removed.

**D1 — preserve wall-clock authority.** No clock injection, time rewrite or expiry
exception. Random credentials stay private. Packet times, effect request Now,
handler output and budgets remain real. Captures freeze public state, not clocks.
Fresh-seed equality is a documented residual for Mark (V12); no approval is requested here.

**D2 — server-controlled stills, browser JS enabled.** Add explicit
`serve --workbench-still` (default false). It suppresses the single script emission the
renderer has, `<script src="/workbench/live.js" defer></script>` (render.go:192, V24) —
there is no inline script and no nonce to pin, and the stripped server HTML was
byte-identical across repeated requests on every workbench route (V24), so nothing
else needs freezing on the measured store. If M6/M7's repeat-GET on the seeded demo
store ever finds varying bytes, the flag's scope is NOT silently widened: the varying
element is reported and the M7 fallback below applies. All server-rendered
regions and the truthful “Live updates: off. Reload to refresh.” footer remain.
It changes no commits, permissions, expiry, CSP, grade resolution or default live
behavior. The caption says “static capture of recorded history.” Neither broken
JS-disabling Chrome flag is used. This is a new projection option, not a measured
existing option. V24/V25 measured its premise by simulation; M6 asserts it on the real flag.

**D3 — show grades of evidence, not invented object grades.** Add a bounded evidence
roster to the validated packet view, sourced from its carrying invocation's input.
Keep packet v1 bytes and its existing two subject-grade lines. Those lines may
honestly remain unavailable; the roster separately shows actual TESTED and ATTESTED
reports, a truthful demo-authored CLAIMED claim, and a subject-bound PROVEN report. Never promote report text to an
object's grade. Why this is not a package: orchestration, HTTP adapters and rendering
belong to the host; the demo domain transition and evidence-kind law belong in a
small AILANG package. No new world/ exports or Go module requirement.

### Quorum r1/r2 objections — disposition

IDs identify every present reviewer result and both controller results; absent
sonnet in each artifact supplied no objection and is not counted as approval.
“Fixed” means this specification answers the objection; it does not mean future
acceptance measurements passed.

| Objection id / reviewer | Disposition | Evidence / answer |
|---|---|---|
| r1-0 / oc-glm-5-3 | fixed | V13–V15; JS-off recipe removed, server regions verified, still flag specified, bounded success requires real file; V21 product measurement remains explicit. Axiom claims reduced. |
| r1-1 / oc-kimi-k3 | fixed | V12 and D1 remove unreachable fresh-seed AC; varying-byte inventory disclosed, no library-clock fiction. V14 answers rendering premise. |
| r1-3 / gemini-3-1-pro | fixed | V13 supersedes stale Chrome 154 result with successful Chrome 155 base invocation; no kill-after-file success protocol. Product acceptance still required. |
| r1-controller | fixed | H1 ratified and prerequisites measured V2; release mint remains attended, no quorum-supplied human decision. |
| r2-0 / oc-glm-5-3 | fixed | D1 states changed scope/residual; D2 removes JS-off premise; exact raw-GET and browser probes V15/V21; packet scope now ratified; verification table is contiguous. |
| r2-1 / oc-kimi-k3 | refuted | Charter line 6130 explicitly says “No Go dependency (no chromedp).” V1 pins the formerly uncited premise; V13/D2 remove the failing flag. No fallback dependency is authorized. |
| r2-3 / gemini-3-1-pro | fixed | Go image/png decoder and chunk verifier, V20; Python is supervisor only. No Python pixel script or optional installed decoder. |
| r2-controller | fixed | H1 resolved; capture lifecycle explicit below; fresh-seed scope openly residual, no exception extension inferred. |
| r3 / oc-glm-5-3 | fixed (reviewer's verbatim alternative) | D3 + roster paragraph: authored-claim kind via package predicates, drift-tested, rendered "demo-authored"; no AiReview cast; kind-resolution check → UNAVAILABLE + seed refusal; rows render producer/author. The reviewer's FIRST alternative (name S3's ai_check report as an AiReview) rests on a false premise and was not applied: V28. Its two V-row requests: packet pending status V27; 390×844 layout V25 (phone shot rendered real content, byte-identical). |
| r3 / oc-kimi-k3 | fixed | (a) repeat-GET run now, V24: empty varying-byte inventory; D2 now enumerates the one suppressed tag. (b) inline-script V-row: V24 (none). (c) 148/149 SHAs in the Dependencies header. (d) V22/V23 reordered; contingency = M7 fallback (no silent widening of the flag). |
| r3 / claude-sonnet-5@claude-p | fixed | Real-content dual-pass capture V25 (five shots byte-identical); M7 fallback clause added verbatim from the proposed fix. |

## Approach

### 1. Seed composition and attended boundary

`demo seed --db PATH --session FILE --ailang-bin PIN --proof-key-file KEY
--addr URL --manifest FILE` orchestrates an initialized disposable store. Refuse a
nonempty public log/world before seeding; interpreter/epoch/transition registry
bootstrap and attended session rows are expected. Refuse existing output manifest,
missing grant, mismatched db/health/pin, expired session or an unexpected head change.
The seed talks to the sole-writer daemon; it never opens a competing writable store.
All history objects use POST `/v1/commit` or the ordinary coordinator invocation.

Human release steps, in order: provide a disposable episode worktree containing
unchanged world/logepoch.ail; with the daemon stopped, publish the combined se-tools
and demo descriptor manifest using `world-publish transitions --store DB --manifest
MANIFEST --ailang-bin PIN` and type its confirmation phrase. Then run **one**
`session mint --db DB --episode demo --grant world.apply=world:10
--grant Workspace.Read=worktree:1 --grant Ailang.CLI=worktree:2
--grant Human.Approve=worktree:1 --ttl 3600 --out SESSION` and type `y`.
No decision/approve verb is needed: the screenshot shows a real pending request.
A missing worktree is preparation, not permission for seed to bypass mint or create
one through a hidden git operation. QUICKSTART must document its preparation exactly.

The loop can build/test tooling, prepare payloads, generate a fresh private 32-byte
operator key (0600, outside the repo), bootstrap, stop/start, seed after that mint,
read, capture and verify. It cannot type either attended confirmation. Release
capture is a separate attended run; test credentials use only V11's existing seams.
No session token or operator key enters manifests, images, stdout or committed files.

| Sequence | Existing path and actual data |
|---|---|
| S1 | POST `/v1/commit`, Authorization Bearer session; normal observedHead CAS, exact payload hashes, valid world/entry headers per QUICKSTART §2. Commit demo goal and two candidate descriptions (`world/demo/goal/v1`, `world/demo/proposal/v1`), source bytes of unchanged world/logepoch.ail, and explicitly labelled demo commentary. These domain descriptions are not forged kernel Proposal values. |
| S2 | POST `/a2a/`, JSON-RPC tasks/send, task `demo-read`, metadata.skill_id `ailang-read`, data `{path:"world/logepoch.ail"}`. Existing transition packages/se-tools/se_tools/read.ail records Workspace.Read result. Assert read bytes equal committed source; retain invocation and effect record refs. |
| S3 | Same A2A route/task `demo-test`, skill `ailang-cli`, data `{op:"test",path:"world/logepoch.ail",flags:{json:""}}`; then task `demo-check`, skill `ailang-cli`, data `{op:"ai_check",path:"world/logepoch.ail"}`. Existing cli.ail/policy-tool producer returns actual argv, tool/policy pin, exit/stdout/stderr and recorded World effect refs. Fail on null, skipped or failing reports; retain raw bytes. ToolAilangBin is the existing daemon-required v0.52.1 tool pin, distinct from capsule/proof PIN v0.41.0. Never call the CLI report a v0.41.0 proof. |
| S4 | First production `evidence.NewProducer` caller in demo seed. Read source through GET `/v1/objects/{ref}?payload=true`; writer adapter commits the producer-supplied object via session-authenticated POST `/v1/commit`. Invoke GenerateProof(sourceRef,["sameRef","servesEntry"]). No custom envelope encoder, testdata source, fixture key or writable-store shortcut. |
| S5 | POST `/a2a/`, task `demo-approval`, skill **new package** `demo-approve`. Existing coordinator dispatchEffectful→BoundInvoker.Request→HumanHandler→mintApprovalPacket is the transition mechanism (V5). No existing production approval skill is claimed (V7). Input names goal/candidates and four evidence refs. New source implements convention v2 and emits exactly one Human.Approve@worktree cost 1, id e1, requester demo, finish=false, result=null. Descriptor access Human.Approve@worktree cost 0; declared effect cost 1. No extra effect or decision. |
| S6 | GET packet semantic-ID listing, require exactly one packet from S5's carrying entry; fetch packet, request, proposal/effect-plan, input, record, effects and output by hash. `why <approval-entry> --json` and why for read/test/check entries: coordinator=true, ok=true, named world/entry/record/input/plan/effect/output links all true. Serialize public manifest only after all checks. |

Producer adapter BusyTimeout reports 0 (no local lock retry); each HTTP object leg
has 10s timeout, producer read/write bounds 12s, ExecTimeout 30s, MaxOutputBytes 1 MiB,
whole S4 deadline 90s. Pin checker hash/version to health's archived interpreter.
Producer itself stages the standalone module and enforces 5s per solver query.
On uncertain commit outcome, stop and diagnose; never blindly retry/reseed.
Serve seed/capture with `--workspace-root WORKSPACES --tool-ailang-bin TOOLPIN`
(the existing v0.52.1 binary), and `--proof-key-file KEY --proof-required-identity sameRef
--proof-required-identity servesEntry --ailang-bin PIN`; the same private key is
required to authenticate the stored proof. Key loss makes capture fail, not downgrade.

Evidence roster input has only kind/ref pairs (four entries; no caller grade).
TESTED names S3's test invocation record; validate its carrying commit, single
Ailang.CLI effect/request/result, argv test, JSON success, ≥1 named passing test,
zero failed/skipped/vacuous skips; show verdict pass. ATTESTED names the real S2
RecordedEffect; decode RecordConsistent and check request/result hashes. CLAIMED
names a demo-authored claim through a truthful authored-claim kind added by the new
package's kind predicates, drift-tested against world/types.gradeOf's CLAIMED mapping,
and rendered "demo-authored" in its roster row and in captions.md. S1's commentary is
never cast as AiReview (no AI review produced it, V28); the AiReview kind is not reused
for non-AI text. The packet projection verifies that every roster ref resolves to an
object of its declared kind; a kind mismatch yields UNAVAILABLE and seed refusal, and
each roster row renders its producer/author.
PROVEN names S4 envelope plus source subject; only ValidateProof→Resolve can
supply that grade. These are evidence-constructor grades (V8); lower grades do
not assert authenticated source binding. Link all report/subject refs and disclose
compiler/tool pins. Invalid/missing evidence yields UNAVAILABLE and seed refusal.

New package predicates map the three legacy kinds per world/types.gradeOf; contract
and inline tests precede implementation, with drift tests against the existing law.
The host projection follows the validated packet's carrying record.input, verifies
content hash/semantic ID and strict roster shape, and performs the above bounded
reads. It does not scan arbitrary demo semantic IDs to decorate a packet. Roster/read limits: exactly four unique kinds, ≤64 KiB input, ≤1 MiB per result,
≤4 KiB commentary, ≤1 MiB decoded test stdout, ≤256 named tests and a shared 10s
read deadline; any overflow is UNAVAILABLE. No roster appears for other schemas. PROVEN authority stays exclusively in row105.

Manifest: schema/version, frozen head, ordered entry/object refs, source/pins, goal/
candidates, packet/proposal/request, four evidence refs, invocation/effect refs and
why JSON hashes. Credential, absolute paths and capture time stay in private or
separate environment diagnostics. Capture rejects missing/mismatched manifest refs.

### 2. Bounded Chrome 155 capture and complete shot list

`scripts/capture_marketing.sh` serves the seed with proof configuration and
`--workbench-still`, checks health/head with ≤10s requests and ≤30s startup deadline,
and makes two full passes in scratch output directories. No commits, mint, approval
or background agent runs during capture. Check public head/entry/object inventory
before and after both passes; an unexpected mutation invalidates the result.

| Filename | Actual route | CSS viewport | Required content |
|---|---|---|---|
| workbench-light.png | `/workbench` | 1440×1100 | selected head, nonempty live list, graph and pending approval |
| workbench-dark.png | `/workbench` | 1440×1100 | same data, measured dark CSS tokens |
| packet.png | `/workbench?object=<S6 packet hash>` | 1440×1400 | validated packet, real proposal/request links, pending status, four-grade evidence roster |
| graph.png | `/workbench?from=N&entry=N`, N=S5 entry | 1440×1100 | selected entry and checked graph targets; no invented evidence graph |
| provenance.png | `/workbench?object=<S5 record ref>&payload=1` (V26) | 1440×1400 | invocation payload, stored input/plan/effect/output references; accompanying why JSON is the actual walk |
| phone.png | `/workbench` | 390×844 | selected world, readable responsive timeline; no requirement that all desktop panels fit above the fold |

Use controller-proven base flags plus screenshot, viewport, scale=1, lang=en-US;
full argv is in revision notes. HOME stays intact, stdin is /dev/null, each shot
uses an isolated profile and a new absolute destination, fixed Chrome/OS/fonts.
No virtual-time assumption is needed. Every shot pins its colour scheme (V23):
light shots add `--blink-settings=preferredColorScheme=1`, the dark shot
`--blink-settings=preferredColorScheme=0`; an unpinned shot would follow the host OS
appearance. `--force-dark-mode` is not used. M7's verifier asserts each shot's scheme
from its pixels (background luminance in a measured rectangle), so a scheme that
silently failed to apply reds rather than ships. Pin the measured
browser and font environment in sibling text diagnostics; cross-platform equality
is not promised. A failed dark probe blocks capture; do not recolor an image.

Python stdlib supervises subprocesses only: start_new_session, communicate timeout
30s; on expiry SIGKILL the process group, reap within 5s, report failure even if a
PNG exists. Daemon cleanup sends TERM, waits ≤5s, then KILL/reaps ≤5s; cleanup failure
is an error. Shell traps cover normal exit/signals. A six-shot pass is bounded to
300s including startup/cleanup; each wait stays short enough to report progress.
Success requires rc=0 **and a new regular nonempty PNG**; an rc=0 missing file is a
loud error. Isolated destinations prevent a stale screenshot from satisfying it.

Before Chrome, raw HTML assertions establish each required region/ref and no script.
Before publication, the Go verifier decodes every PNG: exact dimensions, >32 RGB
colors, at least 1% pixels with max(R,G,B)-min(R,G,B)>12, and variance in each
shot-specific expected content rectangle. Packet rectangles cover validation and
all four roster rows; baseline rectangles are measured in the controller rehearsal,
not guessed as proof of visible text. Raw DOM checks establish text; crop variance
establishes that those regions are not blank. Full PNG SHA256 must match between
passes. Do not accept sanitized or visually similar bytes. V21 remains unexecuted.

### 3. PNG gate, captions and wiring

Only the six table names, directly under `website/static/img/v1/`, matching
`[a-z0-9][a-z0-9-]*.png`. Verify signature, IHDR first/length 13, valid CRCs, allowed
IHDR/PLTE/IDAT/IEND/sRGB/gAMA/cHRM/iCCP/pHYs chunks, exactly one terminal IEND and zero
trailing bytes. Reject tEXt/zTXt/iTXt/eXIf. ≤1,048,576 bytes each; ≤8,388,608 bytes
for **all** PNGs in that directory. Verify raw output before copying it there.
Oversize/incompatible output fails; no metadata stripping, resizing or recompression.

Sibling text `capture-manifest.json` records hashes/dimensions/shot routes and public
seed refs; `captions.md` gives alt text, static-capture disclosure and evidence limits.
README links light/packet images and reproduction docs. Website getting-started page
uses `/img/v1/<name>.png` with real captions/alt text; no route or secret in alt text.
Only verified output is promoted from scratch. The fleet gate itself is not edited;
its committed-blob check runs after the controller's commit, not falsely on unstaged PNGs.

## Files and planner milestones

Proposed implementation paths below are future work, not files changed by this revision.
Each milestone lands independently in one ≤30-minute codex executor run; release
images require the separate attended run. Split again if a run cannot meet that bound.

| Run | Files / purpose | Acceptance and mutation |
|---|---|---|
| M1 | packages/world-demo/approve.ail, evidence.ail, transitions.json; package documentation | C1: pinned check/tests and named predicate contracts; valid one-effect plan. Mutate effect count to two: exact-one-plan assertion reds. |
| M2 | cmd/ailang-worldd/demo.go, demo_seed.go, demo_test.go; main.go dispatch/help; QUICKSTART seed preparation | C2: session-authenticated goal/candidate commits and existing read/test calls; no hidden Mint. Remove bearer header: successful-seed assertion reds with SessionAbsent. Remove empty-world guard: nonempty-store refusal assertion reds. |
| M3 | cmd/ailang-worldd/demo_proof.go and tests; HTTP producer reader/writer adapters | C3: real pin/source/key producer result validates and is REST-carried. Mutate source bytes: hash check reds before mint; wrong key: PROVEN assertion reds. Null verified set: producer refusal assertion; bypass HTTP writer: REST-carrying-commit assertion reds. |
| M4 | demo_seed.go integration/manifest; CLI why tests | C4: approval invocation produces exactly one anchored packet, effect budget delta and complete why set. Replace coordinator call with plain REST object: coordinator=true assertion reds. Remove effect link: named-link assertion reds. |
| M5 | host/daemon/workbench_packet_evidence.go and tests; host/workbench render structures/template; package law drift test | C5: roster displays actual four evidence grades; existing object-grade lines stay truthful. Mutate commentary kind to PROVEN: claim cannot promote assertion reds; substitute arbitrary test text: TESTED refusal assertion reds; bad record hash: unavailable assertion. Point the CLAIMED roster ref at an object of another kind: kind-resolution assertion reds (UNAVAILABLE + seed refusal). |
| M6 | serve flag in main.go/daemon.go, workbench Page/render flag and tests | C6: still HTML contains all regions/no script/static footer; default retains live script. Force script emission in still mode: no-script assertion reds; hide graph: required-region assertion reds; invert default flag: default-script assertion reds. |
| M7 | scripts/capture_marketing.sh, scripts/marketing_capture.py supervisor, scripts/marketingcheck/*.go standard-library verifier and tests | C7–C9: timeout/cleanup, six files, decoded content, exact bytes and gate rules. Fake Chrome rc=0/no file, hanging child, uniform PNG, missing shot, altered byte, forbidden chunk, bad CRC, trailing byte and cap+1 each trigger their corresponding named assertion. **Fallback:** if raw-byte equality fails on content-bearing shots for a cause other than colour-scheme/script-timestamp, the design documents that as a new residual (parallel to the fresh-seed-equality residual already carved out for Mark) rather than silently loosening the acceptance bar to pixel-similarity; the varying element is named with both hashes. |
| M8 | docs/QUICKSTART.md reproduction; README.md; website/docs/getting-started/marketing-capture.md; captions/manifest text | C10: links resolve and recipe/help match flags. Mutate one filename or serve flag: link/runbook assertion reds. Controller later runs V15/V21, attended seed and promotes six PNGs; an unminted unattended run must refuse. |

## Acceptance criteria

C1–C10 and their mutations above are required, not estimated future passes.
M4 requires `why` coordinator=true and ok=true, ≥1 effect and exact expected link
names; a zero-link report fails. M5 requires all four kinds with valid report refs;
PROVEN requires the operator key and exact source subject, not a fixture badge.
M7 requires six names in both passes, all limits, byte equality, no blank regions,
no false-green rc=0, and bounded child cleanup. These discharge 150(a)–(c).

**C11 OBSERVATION:** before/after SHA256 and line count for world/*.ail, go.mod,
go.sum, tools/launchd/* and design_docs/world-mission.md, plus a scoped diff.
Expected equality is observed without writing any of these paths. Change a copied
observation input: the equality comparator must red. Seed/capture inventory equality
is likewise an OBSERVATION; inject a genuine extra commit in a disposable test store
and the frozen-inventory assertion must red. Do not “restore” protected files.

Tests use real producer/coordinator paths and test-only credential seams V11;
fake phase runners are refusal/unit controls, not release acceptance. Socket tests
that report `bind: operation not permitted` are UNINFORMATIVE UNDER SANDBOX.
Mandatory repository checks belong to executor/controller; no broad-suite pass is
claimed by this documentation revision. A red base check is reported, never hidden.

## Conflict surface

Seed orchestration, package domain, packet projection and still-mode flag are new
work. Existing packet wire/REST prohibition, row105 proof authentication, normal
session resolution, authority expiry, budgets and default live behavior retain their
contracts. Re-run their focused suites after affected milestones. M5 adds evidence
report rows without weakening subject-grade resolution or reinterpreting packet v1.
The PNG allowance is fleet-owned; no verification-gate exception is requested.
No driver/charter/kernel edit is part of row150. Mark owns the fresh-seed residual.

## Axiom alignment

A1/A5: intended benefit from real history and reproducible artifacts; score 0 until
attended seeded captures pass. Explicit evidence-constructor labels preserve the
trust boundary; authentication alone supplies PROVEN. Package-first domain law,
normal broker/coordinator commits and no kernel growth preserve the existing split.
No completed end-to-end or clause movement is claimed by this design.

## Review status

R1/r2 (iteration 246) and r3 (iteration 266, `quorum/w-marketing-capture-2026-10-14T21-13-09Z.json`,
glm/kimi/sonnet reject, OpenAI benched as author vendor, no absentee) remain BLOCKED rounds;
none is retroactively a pass. r3's three objections each carried a concrete reviewer fix and
none disputed the design direction, so r4 is the controller's bounded revision applying those
fixes under the Gate-2 narrow-refinement carve-out (plus measurements V24–V28). Round surfaces:
r3 = evidence-roster kind (D3) and capture byte-stability (D2/M7). Dispositions r1–r3: 10 fixed,
1 refuted, 0 unanswered; absent reviewers are not approvals.
Q1 normal-path composition, Q2 determinism scope/residual, Q4 PNG/wiring and Q5 bounded
milestones are specified with first-party rows. Q3: colour forcing measured (V23); the
V15/V21 product measurements are M6/M7 acceptance, run on the code those milestones add. No additional human scope
question is open; Mark's only residual is fresh-seed hash equality (D1). Tooling/tests and the
attended release capture remain distinct.
