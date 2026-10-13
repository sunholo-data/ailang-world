# Sprint plan — `w-decision-packet-production-projection`

**Iteration 255 · queue row 172 · release prerequisite (moves no clause) · D-WORLD-69 = A path 26 → 105 → 172 → 150**

Design: [w-decision-packet-production-projection.md](w-decision-packet-production-projection.md). Quorum r1 BLOCKED, then r2 BLOCKED, then a narrow-refinement carve-out. Worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter255-row172-design`, branch `design/w-decision-packet-production-projection`, planned at HEAD `b72428a`. The code base is `e869bfe` plus three design-only files (`git diff --stat e869bfe HEAD` lists the design doc, `decisionpacket.ail.txt` and `sketch-check.log` and nothing else). Binding rules: CLAUDE.md and design_docs/coding-standards.md. S1/S2 put the pure laws in the AILANG sketch and the effects in Go at the host boundary. S6 requires non-vacuous gates and a named mutation for every load-bearing criterion. S7 requires usage docs for a usage surface. S8: the `world/` floor does not move.

**~3.35 executor days, 7 milestones C1–C7, with step C4.2 HELD (see the next paragraph).** These are the design's M1–M4, re-cut so that each milestone stays ≤ ~250 production LOC and is cohesive enough for a small-model executor (PD24). Written by iteration 255's planner (killed mid-run, rc=143); reviewed, re-measured and finished by iteration 256's planner (PD28–PD33 and the "iteration 256" rows below).

> **D4 IS HELD — decision `D-WORLD-73` (filed by the controller, default = hold).** This sprint lands the design's M1, M3's validator and projection, the REST reservation and M4, and does **NOT** land M2's daemon binding. Concretely: **step C4.2 is HELD and must not be implemented** (no `humanMu`, no `human` field, no `serializedHandler`, no change to `d.binder`; `host/daemon/daemon.go` is NOT modified). Held with it: AC3's daemon arm (`TestDaemonInvocationMintsDecisionPacket`), AC8 (`TestHumanApproveIsSerializedAcrossEpisodes`), mutations MD2 and MD3. C4.1 (the REST reservation) lands. C5–C7 never touch `d.binder` (PD25), so they land unchanged. With D4 held, **no production path stores a packet** (PD29, measured), so the D2b byte-form freeze is landed and test-pinned but not yet exercised by any production write.

**Plan only.** The planner writes exactly two files: this one and `sprint_w-decision-packet-production-projection.json`. **The EXECUTOR runs no git write** (no add, commit, stash, checkout, restore or branch). Read-only `git diff --stat`, `git diff --quiet` and `git status` are allowed. The executor implements C1→C7 in order, in this worktree. Every gate in a milestone must be green at its own boundary before the next milestone starts. At each boundary, write a report containing: changed files with LOC, the rc and wall time of each gate, the red-first evidence, and the mutation ledger with restore evidence. **The controller commits.** Every boundary is internally consistent, so per-milestone commits are bisectable.

**Controller flags carried from the design (D2b, D4).** D2b, the byte form, is irreversible once a packet is stored; with D4 held no production path stores one (PD29), so the freeze becomes live only when `D-WORLD-73` is answered "bind". D4, the `Human.Approve` daemon binding, is **HELD** for this sprint (`D-WORLD-73`, the design's own default): step **C4.2** is skipped (binding and serialization: AC3's daemon arm, AC8, MD2, MD3). C5–C7 never use `d.binder` (PD25), so they land unchanged. The REST reservation (C4.1) is not the binding and lands. If `D-WORLD-73` is later answered "bind", C4.2 below is the ready-made follow-up sprint step, plus the QUICKSTART sentence swap named in C7.

**Two executor runs (PD30).** **Run A = C1 → C2 → C3**, then the run-A boundary gate. **Run B = C4.1 → C5 → C6 → C7**, then the run-B boundary gate and the C7 re-drill. Each run targets ≤ ~45 min. The executor never runs `scripts/verify_go.sh` (the controller runs it after run B). Run B starts only from a tree on which run A's boundary gate was green and whose run-A files match the sha256 manifest run A reported (PD30).

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang   # AILANG v0.41.0, sha256 1a67b014…b9f — the ONLY binary
export PATH=/opt/homebrew/bin:$PATH
```

Run gate scripts under **bash** (`bash gate.sh`). The Bash tool's shell is zsh, and zsh does not word-split `for t in $X`. In scripts, call `/usr/bin/grep` (the rig's `grep` is ugrep). **The executor must NOT run `scripts/verify_go.sh`** (the controller runs it, outside any sandbox, after run B).

---

## Baseline (planner-measured at `b72428a`, first-party, 2026-10-13, load avg ≈ 3–4.5)

Logs: `~/.ailang/state/iter255-plan/base-*.log`.

| Gate | Result |
|---|---|
| `AILANG_BIN=… ./scripts/verify_ail.sh` | **rc=0**, 31 s: `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 444 se-tools named tests pass`; Leg 1b says `(22 modules)` |
| `go vet ./...` | rc=0 |
| `go test -run '^$' ./...` (compile fence including `_test.go`) | rc=0, 10 s |
| `go test ./host/... ./cmd/... -count=1` (AILANG_BIN set) | **rc=0**, 254 s. Times: host/verifygate 250 s, host/daemon 196 s, host/broker 137 s, host/coordinator 65 s, host/store 18 s, cmd/ailang-worldd 10 s, host/workbench 5 s |
| `go test -race ./host/coordinator ./host/broker ./host/daemon ./host/workbench -count=1 -timeout 9m` | **rc=0**, 225 s. Times: daemon 221 s, broker 163 s, coordinator 53 s, workbench 2 s |
| banked sketch, scratch copy at `sketches/decisionpacket.ail`: `ai-check -timeout 5s` and `test --format json` | check passed, `verify.errors 0`, `counterexample 0`, **9/9 verified** (`lifecycleOk mintable mintedScheduleWellFormed packetRefusal policyCode producerFieldsOk samePair sameRef scheduleWellFormed`); **37** named tests, all `pass` |
| `scripts/verify_go.sh` | **not run by the planner**: its `-race -timeout 8m` leg exceeds the tool's 600 s foreground cap. It is a controller-run final gate (C7), as in row 105. |

No gate is red at base, so every gate below measures the change.

**Load note.** host/daemon `-race` already takes 221 s here (358 s under verify_go's load in row 105). All new daemon tests together must add **< 15 s** to `go test -race ./host/daemon`. The executor records the wall time at C4, C5 and C6. No `t.Parallel()` and no `time.Sleep`. Bounded waits use `select` with `time.After`.

---

## Verification of the design's load-bearing anchors (planner, first-party at `b72428a`)

Every anchor below was read in this worktree. Where a claim was wrong, the PD that resolves it is named.

| Design claim | Measured |
|---|---|
| D1 `dispatchEffectful` at `effectful.go:273`, `now` :284, `planRef` :293, `planEffectInvocation` :333 | exact: `:273 func (c *Coordinator) dispatchEffectful(ctx, surface, id, episodeID, input, bound, d)`, `:284 now := c.cfg.Now()`, `:293 planRef := hashref.SumSHA256(pl.Canonical)`, `:333 p := planEffectInvocation(…, records, now)` |
| D1 `grants` beside the call (`coordinator.go:311-334`) | `:311 grants, err = c.seedGrants(…)`, `:317 c.cfg.Binder(call.EpisodeID, grants)`, `:334 return c.dispatchEffectful(…, bound, d)`. This is the ONLY call (`grep -rn "dispatchEffectful(" host` → 2 hits: this call and the definition) |
| D1 `results[0].Status == StatusOK`; `pl.Effects[0].Requirement.Effect` | `effectful.go:55 StatusOK = "ok"`, `:188 res.Status, res.Output = StatusOK, out`; `effectplan.go:73-77 PlannedEffect{ID, Requirement, Payload}` |
| D1 the authorizing grant is decideOver's choice | `decide.go:66-81`: evaluates `grants[0]` first, then returns the first ALLOWED grant among `grants[1:]`; only if none of those is allowed does it return `(0, Decide(grants[0]))` or the best-ranked denial. **Corrected iteration 256:** it is NOT "the first allowed grant" when `grants[0]` and a later grant are both allowed (it returns the later one). `AuthorizingGrant` delegates to `decideOver`, so producer and session agree by construction; tests must not put an allowed `Human.Approve` grant at index 0 (PD22 keeps `invokeGrant` there). The session decides over a copy of the seeded grants (`broker.go:122-124`; `seedGrants` `effectful.go:105-117` only lowers budgets, never reorders) |
| D1 the pending output is the handler output | `approve.go:121 return pendingBytes(requestObj.Hash)`, `:277-279` `pendingWire{Status:"pending", RequestRef}`, `requestEffect` copies it into `res.Output` (`effectful.go:188`) |
| D4 `d.binder` at `daemon.go:780-782` | exact. `registry()` returns a fresh map on every path (`workspace.go:338-362`); 0 test calls of `d.binder(`, but `workspace_family_test.go:127` and `workspace_test.go:567` pass `Binder: d.binder` to a coordinator (the C4 precedent) |
| D4 `handleCommit` / `decodeCommit` | `handlers.go:714`, `:799`; `SemanticID: object.SemanticID` passthrough at `:847`. This is the ONLY REST write path: `grep -rn "PutObject(\|\.Commit(ctx" host/daemon host/projection cmd/ailang-worldd` (non-test) → `handlers.go:785` alone |
| D5 `readStore` ⊇ PacketReader (`daemon.go:501-513`) | exact; 11 methods, including `GetObject, GetLogEntry, GetRegistryHead, ObjectsBySemanticID, ObjectCommits` |
| D5 `recordV2` fields (`plan.go:39-50`) | exact; decoded elsewhere with `json.Unmarshal` + `wellFormedV2` (`effectful.go:376-378`, `:508`) → PD16 |
| D5 `ObjectCommits(packet, 0, 1)` | **defect** → PD1 (the start sentinel is −1) |
| D5 `RecentApprovals` bounds | `approvals_view.go:61` 40 heads, `:101` truncation, `:103` 20 rows; `ApprovalPage{Summaries, Truncated, MalformedRef}`, `ErrApprovalChainMalformed` |
| D6 object branch `workbench.go:473-483` | `:473 grade, err := d.objectGradeView(ctx, ref)`, `:478 page.Object = &workbench.ObjectView{`, `:482 Grade: grade, …, References: references,`. Object variable is `object` (`:410`, full payload) |
| D6 `ObjectView` `render.go:46-58` | `:46 type ObjectView struct` … `:57 References *ReferenceView`; template `{{with .Object}}` at `:204-210` |
| §9 `verify_ail.sh` manifest | `LEG1_MODULES=(` at `:174`, `REQUIRED_VERIFIED = {` at `:321`, `EXACT_TOTAL_VERIFIED=16` at `:403`, `EXACT_TOTAL_TESTS=40` at `:477` (Leg 2 runs `world/` only). Leg 2c (execargs) is sketch-specific. The sketch adds no test leg |
| §9 "verify_ail.sh: +1 module line, +1 sketch entry" is the whole gate surface | **incomplete** → PD3 (host/verifygate hand-maintained counts) |
| AC9 / M1 "eight identities / eight mirrors" | **stale**: the sketch has 9 exported laws → PD2 |
| §8 MP5 "PutObject instead of Commit.Objects" | **does not compile** → PD7 |
| §8 MV12 killed by the counting wrapper | **survives** → PD6 |
| AC3b "same Result as a packet-free invocation" | **unsatisfiable** → PD9 |
| D5 code 5 "among the newest 40 approval heads" | **imprecise** → PD28: `RecentApprovals` reads ≤ 40 heads but keeps only the first **20** distinct requests (`approvals_view.go:103-106`, `Truncated = true` when cut) |

### Iteration 256 re-measurement (first-party, worktree at `b72428a`, `/usr/bin/grep`)

| # | Anchor | Command | Observed |
|---|---|---|---|
| R1 | producer insertion point | `grep -n "func (c \*Coordinator) dispatchEffectful\|now := c.cfg.Now()\|planRef := hashref.SumSHA256\|p := planEffectInvocation(" host/coordinator/effectful.go`; `sed -n 324,336p …` | `:273`, `:284`, `:293`, `:300` (zero-effect branch, untouched) and `:333`; `:328 composeEffectOutput(fctx, …, results)`, `:330-332 if err != nil { return unrecorded(err) }`, `:333 p := planEffectInvocation(…, records, now)` — the C3 insertion slot exists exactly as written |
| R2 | single `dispatchEffectful` call | `grep -n "grants, err = c.seedGrants\|c.cfg.Binder(call.EpisodeID, grants)\|return c.dispatchEffectful(" host/coordinator/coordinator.go`; `grep -rn "dispatchEffectful(" host --include='*.go' \| wc -l` | `:311`, `:317`, `:334`; `2` (definition + this call) |
| R3 | `d.binder` (C4.2 HELD target) | `sed -n 776,783p host/daemon/daemon.go` | `func (d *Daemon) binder(…)` at `:780`, body `return broker.OpenBinder(d.store, episodeID, grants, d.workspace.registry(episodeID))` at `:781` |
| R4 | C4.1 slot + 400 mapping | `sed -n 836,849p host/daemon/handlers.go`; `sed -n 760,766p …` | `:838 for i, object := range r.Objects {`, `:847 SemanticID: object.SemanticID,`; `:763-764 commit, err := decodeCommit(request); if err != nil { writeAPIError(w, "BadRequest", err.Error(), http.StatusBadRequest)` |
| R5 | PD1 sentinel | `grep -n "genesis\|entry_index > ?" host/store/object_commits.go` | `:20 … entry_index > ? ORDER BY entry_index …`; `:23 "(-1 starts at genesis)"`; and `publishTools` commits genesis at `EntryIndex: 0` (`workspace_test.go:497-498`) — PD1 stands |
| R6 | PD3 verifygate sites | `sed -n '96,97p;221p;231,232p;335,336p;561p' host/verifygate/module_manifest_gate_test.go`; `sed -n 55,104p …` | `files != 24 \|\| ailFiles != 22`; `"… across 22 module(s)"`; `files != 24`; `got != 22`; `` const old = `LEG1_MODULES=( ``; the isolated copy GLOBS `design_docs/sketches/*.ail` (`:62`), so one new sketch moves 24/22 → 25/23. Precedent `git show edadf39` moved 22/20 → 24/22 for two new `.ail` files |
| R7 | PD3 literal = script block | `sed -n 561,584p` verifygate vs `sed -n 174,197p scripts/verify_ail.sh` | identical 22-line lists, alphabetical; `decisionpacket.ail` sorts first (before `effectbroker.ail`) in both |
| R8 | C1 gate edits, dry-run in a SCRATCH COPY (`rsync` of this worktree to `~/.ailang/state/iter256-plan/scratch`, never the worktree) | sketch `cp` + edits (a)(b) applied by script; `bash ./scripts/verify_ail.sh` | **rc=0, 27 s**: `✓ swept .ail module set equals the LEG1_MODULES allowlist (23 modules)`, `ai-check design_docs/sketches/decisionpacket.ail`, banner `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 444 se-tools named tests pass` (unchanged) |
| R9 | MS4 / MS1 in the scratch copy | `sed -i '' '139,140d'` sketch → gate; restore; `0)` → `0 - 1)` on line 73 → gate; restore + `cmp` | MS4 rc=1 `✗ design_docs/sketches/decisionpacket.ail: required identity (decisionpacket, samePair) MISSING from verify.results[] (vanished silently, V20)`; MS1 rc=1 `✗ design_docs/sketches/decisionpacket.ail: verify.counterexample == 1`; restore `cmp` rc=0. Control: the same MS4 edit with the `REQUIRED_VERIFIED` entry REMOVED (LEG1 line kept) → **rc=0**, banner `✓ verify gate PASSED: …` — so only the by-name entry catches the vanished contract |
| R10 | drift-row census and row lines | scratch-only probe test calling `parseSketchTests` on the landed sketch (deleted after) | `laws 9 rows 37`; `policyCode line 28 want int64 -1`; `samePair line 144 want bool false`; **`lifecycleOk line 161 want bool false`** (the iteration-255 text said 162 → fixed in C1's MV7-drift row) |
| R11 | TR.C | `grep -n "const wantCount = 3" host/broker/invoke_boundary_test.go`; `grep -n "func TestRegistryDispatchBindingBoundary" …` | `:274`; `:247` |
| R12 | object branch + render anchors | `sed -n 470,484p host/daemon/workbench.go`; `grep -n "type ObjectView struct\|{{if .PayloadShown}}\|^const partialsHTML" host/workbench/render.go` | `:473 grade, err := d.objectGradeView(ctx, ref)`, `:478 page.Object = &workbench.ObjectView{`, `:482 … References: references,`; render `:46`, `:209` (the `{{if .PayloadShown}}` line), `:239` |
| R13 | PacketReader ⊆ readStore; `RecentApprovals` reader | `sed -n 501,513p host/daemon/daemon.go`; `grep -n "^type ApprovalReader" -A3 host/broker/approvals_view.go` | the five PacketReader signatures appear verbatim in `readStore`; `ApprovalReader{GetRegistryHead, GetObject}` ⊆ PacketReader, so `broker.RecentApprovals(ctx, r, 40)` type-checks with `r PacketReader` |
| R14 | helpers the plan reuses exist | `grep` for each `^func`/`^type` in `host/daemon/*_test.go`, `host/coordinator/*_test.go` | all present: `testCommit(observed store.World, index int64, label string)`, `seedGenesisEmbedded(…) store.World`, `assertErrorClass(…) APIError` (`.Error.Message`), `gradePage(t, d, ref) (int, string)`, `newGradeDaemon(t, ids) (*Daemon, gradeRig)`, `wsCaps` (a type), `publishTools` (`Access {Effect, WorkspaceScope, 0}`, `DeclaredEffects {Effect, WorkspaceScope, 1}`), `fxRig(t, splitAccess, declared...)`, `r.fxCoordinator(st, runner, reg)`, `probe`, `phaseRunner`, `now = int64(100)`, `accessScope = "world"` |
| R15 | QUICKSTART route shape | `grep -n "type objectPageResponse\|type objectResponse" -A6 host/daemon/handlers.go` | `{"items":[{"hash":…}]}` — the C7 commands' JSON path is right |
| R16 | baseline logs | `tail` of `~/.ailang/state/iter255-plan/base-*.log` | verify_ail PASSED banner as above; race coordinator 53 s / broker 163 s / daemon 221 s / workbench 2 s; full `go test` all `ok` (verifygate 250 s); compile fence all `ok` |
| R17 | D4-held production reach (PD29) | see PD29 | 1 production coordinator, bound by `d.binder`, whose registry has 0 `Human` names |

### Round-2 carve-out traceability (design §0.6 → this plan)

| Carve-out item | Where it lands |
|---|---|
| V44 (`validatePublishApproval` enforcement behind D4's authority analysis) | No code: it is D4's premise, and D4's binding is HELD (C4.2). Re-cited in PD29 and in the C7 QUICKSTART text ("decisions stay operator-only") |
| V45 (TR.C census mechanism) + Future Work #5 | Census row "TR.C" below: `wantCount = 3` and `outside = 0` unchanged; with C4.2 held there is not even a `Registry` insertion; FW #5 stays a candidate row, not absorbed |
| V46 (quickstart-test precedent) | C7 `TestQuickstartDecisionPacketSection` (in host/daemon, PD15) |
| V47 (`RecentApprovals` read count: ≤ 80 object reads + 1 head read) | C5 validator step 5 calls `RecentApprovals(ctx, r, 40)` exactly once per page; PD28 states the 20-request row cap that bounds those reads |
| code 7 scoped to the `(requestRef, proposalHash)` pair, `samePair`, 9/9 Z3 | C1 (`samePair` mirror + drift rows, `REQUIRED_VERIFIED` names all 9, MS4 proves the by-name entry is load-bearing), C5 step 7 (`samePair` over `refOf`), sentence `packetRefusedDuplicate` |
| AC5 arm 7b | C6 arm `7b-collapsed-request` (both pages validated, equal `RequestRef`, distinct `ProposalHash`) |
| MP8b | C1 (drift, `sketch line 144`), C6 (arm 7b) |
| two-genuine-same-second-packets AC1/AC3 arm | C3 `TestProducerMintsPacketForApprovalOnly/two-approvals` (2 packets, same `RequestRef`, different `ProposalHash`) + `TestProducerPacketIsCarriedByTheInvocationCommit` two-approvals half (each first-carried by its own entry); the "both render validated" half is C6 arm 7b (PD27). The daemon-binder copy of this arm was in C4.2 and is HELD; C6 7b runs the same scenario through the explicit test binder (PD25) |
| TR.C census quotation | Census table row TR.C (quotes `outside_broker_is_clean` and `inside_broker_exemption`) |
| `RecentApprovals` read-count row (A5/D7) | V47 row above + PD28 |

---

## Acceptance criteria this sprint (D4 held)

| AC | Status this sprint | Test(s) | Milestone | Killing mutants |
|---|---|---|---|---|
| AC1 producer mints exactly one packet, only for an allowed successful `Human.Approve` (incl. the two-approvals arm) | IN | `TestProducerMintsPacketForApprovalOnly` | C3 | MP1, MP5 |
| AC2 packet fields are the commit's facts | IN | `TestProducerPacketBindsInvocationFacts`, `TestDecisionPacketCodec` | C2, C3 | MP2–MP4, MP6, MP7, MB1 |
| AC3 carried by the invocation's own commit — coordinator arm + two-approvals | IN | `TestProducerPacketIsCarriedByTheInvocationCommit`; "both render validated" via C6 arm 7b | C3, C6 | MP5, MP8b |
| AC3 — daemon arm (`TestDaemonInvocationMintsDecisionPacket`) | **HELD (C4.2, `D-WORLD-73`)** | — | — | MD2 (held) |
| AC3b replay/reconcile tolerate the packet | IN (PD9 form) | `TestReplayReconcileToleratesPacketObject` | C3 | MP8 |
| AC4 produced packet validates and renders the validated panel; non-packet pages unchanged | IN | `TestWorkbenchPacketValidated` | C5 | MV12, MP7 (page), MV3-genuine |
| AC5 every forged/mismatched packet renders its own refusal, incl. 7 and 7b | IN | `TestWorkbenchPacketRefusals` | C6 | MV1–MV9, MV11, MP8b |
| AC6 no PROVEN invented | IN | `TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive` | C6 | MV10 |
| AC7 REST cannot store a packet object | IN | `TestRESTCommitRefusesDecisionPacketObjects` | C4.1 | MD1 |
| AC8 `Human.Approve` serialized across episodes | **HELD (C4.2, `D-WORLD-73`)** | — | — | MD3 (held) |
| AC9 laws and mirrors cannot drift; gate names all 9 identities | IN | `TestDecisionPacketSketchRows`, `verify_ail.sh` | C1 | MS1–MS4, MP8b-drift, MV7-drift |
| S7 usage (design M4) | IN (held-variant text, PD32) | `TestQuickstartDecisionPacketSection` | C7 | — (executed doc test) |

## Planner decisions (each with its measurement)

**PD1 — `ObjectCommits` starts at −1, not 0.** `object_commits.go`: "afterEntry … (-1 starts at genesis)", SQL `entry_index > ?`. Entry index 0 is a real carrier: `commitWorkbenchPayload` (`workbench_test.go:256`) commits `testCommit(genesis, 0, …)`. With the design's `0`, a packet first carried at entry 0 would read as uncarried. The validator calls `ObjectCommits(ctx, ref, -1, 1)` for both the packet and the record.

**PD2 — 9 identities, not 8.** `grep -c '^export func' decisionpacket.ail.txt` = 9. The planner's ai-check reports 9/9 verified, and the design's own V24 says `9 verified`. The REQUIRED_VERIFIED entry and the drift map name all nine.

**PD3 — host/verifygate's hand-maintained module counts move 22 → 23.** `module_manifest_gate_test.go` has these sites:
- `:96-97` `files != 24 || ailFiles != 22`
- `:221` `pristineMarker = "✓ 16/16 required world/ identities verified across 22 module(s)"`
- `:231-232` `files != 24`
- `:335-336` `got != 22` ai-check lines
- `:561-584` a verbatim `LEG1_MODULES=(` literal, used as a `mutateCopiedScript` anchor that must occur exactly once

The design's §9 omits this file. Precedent: `edadf39` (row 140) "verifygate's hand-maintained module counts move 20 -> 21". The `world/` total (16), Leg 2 (40) and Leg 2b (444) do not move. `floor_raise_inventory_test.go` is untouched because no `world/` identity is added.

**PD4 — the drift test cannot copy effectplan's "a true and a false row per law" rule.** A scratch copy of `parseSketchTests` run over the banked sketch reports **9 laws, 37 rows**:
- `policyCode` and `packetRefusal` are int-valued, and the parser yields `int64`. Their mirrors must return `int64(...)`, or `got != row.want` is always true.
- `mintedScheduleWellFormed` has only true rows (sketch lines 69–70).

Rule: bool laws need both a true and a false row, except `mintedScheduleWellFormed`, which is an implication law with a named exemption. Int laws need ≥ 2 distinct expected values. The total must be exactly 37.

**PD5 — law mirrors take `lawRef{algo, digest string}`, not `hashref.HashRef`.** `HashRef`'s fields are unexported and `hashref.New` validates hex digests. The sketch rows use `"aa"`, `"q1"`, `"p1"` and algo `"sha512"`. Production converts with `refOf(h) = lawRef{h.Algo(), h.Digest()}`.

**PD6 — MV12 is killed by a body assertion, not by the counting wrapper.** A non-packet object fails code 1 (strict decode, 0 reads; D5 table), so a branch that runs for every object adds **zero** `ObjectsBySemanticID(DecisionPacketV1)` calls and only renders a refusal panel. AC4's non-packet arm asserts that the body contains no `<section aria-label="decision packet">`. The counting wrapper stays as a cost assertion.

**PD7 — concrete MP5.** `coordinator.Store` (`coordinator.go:28-37`) has no `PutObject`, so the design's MP5 does not compile. MP5 becomes: `withPacket` stops appending to `p.Commit.Objects`. The packet is then never stored, and `ObjectCommits(packet)` is empty.

**PD8 — concrete MP8.** The design's MP8 is new hypothetical code. Concrete form: the producer appends the packet ref to `records` before `planEffectInvocation`. Replay then returns `diverged("effects", "2 recorded effect records for 1 planned effects")` (`effectful.go:422-424`), so AC3b reds.

**PD9 — AC3b is satisfiable only without the "packet-free twin".** Every successful `Human.Approve` mints a packet (D1's gate), and a different effect gives a different output, so "the same Result as a packet-free invocation of an otherwise identical skill" cannot be built. AC3b asserts instead:
- `Replay` succeeds and its bytes equal the committed output.
- A resend reconciles (`Reconciled`, equal bytes, zero new handler calls).
- The packet count stays 1.

**PD10 — export the codec.** Package-`daemon` tests must forge packet bytes for AC5, so these become exported: `coordinator.DecisionPacket`, `EncodeDecisionPacket`, `DecodeDecisionPacket`, `DecisionPacketV1`, `ValidatePacket`, `PacketReader`, `PacketResult`. Exporting an encoder confers no authority: REST refuses the label (C4.1), and validation is relational (D5). `mintPacket`, `withPacket`, `mintApprovalPacket` and the law mirrors stay unexported.

**PD11 — [HELD with C4.2; not implemented this sprint] `d.human` is built unconditionally in `New`.** It is typed `broker.Handler`, set right after `d := &Daemon{…}` (daemon.go:666), not inside the `AilangBin` block. Tests use `newHandlerDaemon` (no AilangBin) with `Binder: d.binder`, which is the workspace_family_test precedent. AC8 swaps the inner handler through the field.

**PD12 — sentences are HTML-escaped.** `html/template` renders `'` as `&#39;`. Affected sentences: validated ("producer's"), code 4 ("invocation's") and code 6 ("request's", "producer's"). Every sentence assertion uses `html.EscapeString(sentence)`, as in row 105's PD2.

**PD13 — the code-5 malformed sentence is frozen as** `not a validated decision packet: the approvals chain is malformed at <ref>`. The design writes "… the approvals chain is malformed at <ref>".

**PD14 — forged-chain fixture.** AC5 and AC6 build a **label-forged chain**: plan, effect record, effect result, record v2 and packet, all in ONE `store.Commit` whose `TransitionRef` is the forged record. This is the design's declared D5 residual. Arm 0 is the positive control: the forged chain renders **validated**. Every refusal arm differs from arm 0 by exactly one fact. AC6 uses this chain on a configured daemon (`newGradeDaemon`), so it does not depend on publishing a skill there.

**PD15 — M4's quickstart test lives in host/daemon**, as `TestQuickstartDecisionPacketSection` in `workbench_packet_test.go`, not in `cmd/ailang-worldd/quickstart_packet_test.go`. It needs the package-daemon producer helpers to execute the documented route against a real packet. The QUICKSTART text is appended at end of file as a bold-titled block, not a `### N.` heading. It has **no** line starting with `/tmp/ailang-worldd serve ` and **no** markdown table:
- `TestQuickstartSection11FlagsMatchTheCLI` cuts from `### 11.` to EOF and requires exactly 1 such line (`quickstart_layout_test.go:28-60`).
- `exec_startup_test.go:330-337` parses §10's table between markers.

**PD16 — code 2's record decode** is `json.Unmarshal` + `wellFormedV2(rec)` (the `effectful.go:376-378` idiom), plus `hashref.Parse(rec.TransitionFn)`. The evidence line needs a `HashRef` for `objectGradeView`.

**PD17 — added mutation MS4** (contract vanish), proving the new REQUIRED_VERIFIED entry is load-bearing. Measured: deleting `samePair`'s two `ensures` lines gives `verify.results` 8 entries, `counterexample 0`. Only the by-name entry catches it (`required identity (decisionpacket, samePair) MISSING`).

**PD18 — validator facts short-circuit.** Facts are measured in order, and every fact after the first false one stays false (no reads). `packetRefusal` still computes the code, which is the same code by the law's order. D7's worst case applies to valid packets.

**PD19 — new names for the design's unnamed tests:** `TestDecisionPacketCodec` (C2, "codec tests") and `TestPacketBrokerExports` (C3, "their tests").

**PD20 — template placement.** The packet partial is inlined at the START of the existing `{{if .PayloadShown}}` line, so a nil `Packet` emits 0 bytes. Every non-packet page is therefore byte-identical by construction (AC4's "byte-identical to base").

**PD21 — MP7 mutates ONLY the encoder's hash line.** Encode and decode share `packetBodyHash`. Mutating that helper would change both sides, and the mutant would survive.

**PD22 — AC2's grant ladder** is `[Invoke@world, Human.Approve expired (ExpiresAt=now), live 1760, live 1900]`, and the expected deadline is 1760. This distinguishes first-allowed from first-listed (expired, so the mint is refused) and from last-allowed or max (1900).

**PD23 — added mutation MB1** on `AuthorizingGrant` (it returns the last grant), killed by `TestPacketBrokerExports` and AC2.

**PD24 — milestone re-cut:** M1 → C1 (laws) + C2 (codec); M2 → C3 (broker + producer) + C4 (daemon reach); M3 → C5 (validator + projection) + C6 (forged arms + evidence); M4 → C7.

**PD25 — C5–C7 never use `d.binder`.** Their produced-packet helper builds the coordinator with an explicit test binder that registers `broker.NewHumanHandler(d.store)`. Only C4.2's tests touch `d.binder`, which honours the D4 hold default.

**PD26 — the REST reservation goes inside `decodeCommit`'s object loop** (`handlers.go:838-849`). Its error is answered 400 `BadRequest` with `err.Error()` by the existing `handleCommit` branch (`:762-765`), and that is the only REST write path (see the anchors table).

**PD27 — AC3's "both render validated" half moves to C6 arm 7b.** C4 has no validator yet. The coordinator half (carried by own commits) is in C3. The daemon half (carried by own commits through `d.binder`) was in C4.2 and is HELD with it; C6 arm 7b produces both packets through the explicit test binder and asserts each is first-carried at its own entry AND validated.

**PD28 (iteration 256) — code 5's effective bound is 20 distinct requests within the newest 40 heads.** `RecentApprovals` (`approvals_view.go:56-106`) folds ≤ 40 heads and keeps only the first 20 distinct requests (`if len(rows) > 20 { p.Truncated = true; rows = rows[:20] }`). So a request that is the 21st–40th newest distinct request is reported as `Truncated` and not found. The validator's behaviour is still correct and bounded: it renders the frozen "not established" sentence, never a refutation. But that sentence's text says "not among the newest 40 approval heads", which is false in that window. The sentence is design-frozen (D6) and is NOT changed by this sprint. This is recorded for the controller as a design wording defect (candidate: re-word to "not among the newest 20 approval requests of the newest 40 heads" in a `/v1`-safe follow-up, since sentences are not wire bytes). Arm 5b (41 newer distinct requests) truncates by both rules, so its assertion holds. The C7 QUICKSTART text states the real bound.

**PD29 (iteration 256) — with D4 held, no production path stores a packet (measured).** The packet is minted only inside `dispatchEffectful`, only when the plan's one effect is `Human.Approve` AND the bound registry executed it successfully. Measurements at `b72428a`:
- Production coordinator construction sites (non-test): `grep -rn "coordinator\.New(\|coordinator\.Config{" host cmd --include='*.go' | grep -v _test.go` → exactly **1**, `host/daemon/daemon.go:710`, with `Binder: d.binder` (`:711`).
- `d.binder` binds `d.workspace.registry(episodeID)` (`daemon.go:781`); `grep -rn "Human" host/daemon --include='*.go' | grep -v _test.go | wc -l` → **0**, control `grep -n EffectWorkspaceExec host/daemon/workspace.go` → `:76`, `:360` (the registry's only non-AILANG binding).
- The only other production `Human.Approve` registry is `publish_op.go:135-136` (`Registry{EffectHumanApprove: human}` inside a broker session for `world-publish`); `grep -ln 'host/coordinator"' cmd/world-publish/*.go | grep -v _test` → no file (rc=1), control `grep -ln 'host/broker"' …` → `main.go`, `reconcile_probe.go`. So `world-publish` never reaches the producer.
- After C4.1, REST refuses the label, so the only remaining writers are tests (`store.Commit` directly) and direct DB tampering.
Therefore: on a daemon built from this sprint, `GET /v1/objects/by-semantic-id/world/decision-packet/v1` returns `{"items":[]}` on every store not written before C4.1, and **the D2b byte form is frozen in code and pinned by `TestDecisionPacketCodec`'s golden, but not yet exercised by any production write.** The freeze becomes live the day `D-WORLD-73` is answered "bind". The C7 QUICKSTART text says so.

**PD30 (iteration 256) — two executor runs; fast per-milestone gates, full gates at run boundaries.** The iteration-255 template ran the full suite (~255 s) plus `-race` (~225 s) at EVERY milestone, about 8.5 min × 7 boundaries. That does not fit two ~45-min small-model runs. Replaced by:
- the **fast gate** at every milestone (vet, compile fence, the milestone's NEW tests by name, its touched packages or a named `-run` subset, frozen-tree checks);
- the **run-boundary full gate** after C3 (end of run A) and after C7 (end of run B): `verify_ail.sh`, full `go test ./host/... ./cmd/...`, `-race` on the touched packages, `TestWorkbenchObjectGrade`.
At the end of run A the executor records `shasum -a 256` of every run-A file (the C1–C3 files). At the start of run B it re-computes them and requires equality. Run B never edits a run-A file.

**PD31 (iteration 256) — C7's re-drill is the run-B mutant set on the final tree; run-A mutants are carried by the sha256 manifest.** Run-A mutants (MS1–MS4, MP1–MP6, MP8, MB1, the C1/C2 drift and codec drills) target files AND are killed by test files that run B does not modify. Their run-A ledgers stay valid when PD30's sha256 manifest matches. C7 re-drills every mutant whose target or killer lives in run B: MD1, MV1–MV12, MV3-genuine, MP7-page, MP8b-on-7b. MD2/MD3 are HELD with C4.2. **Rollover:** if run B has used 40 min when the re-drill starts, the executor stops after the run-B boundary gate, reports, and the controller launches the re-drill as a short run "B2" from the same ledger table.

**PD32 (iteration 256) — C7 QUICKSTART text is the D4-held variant, written out in full.** The iteration-255 text described the daemon as binding `Human.Approve` and gave a "replace the sentence" instruction for the held case. A small-model executor must not compose text, so C7 now carries the held text verbatim. It names `D-WORLD-73`, tells the reader the listing is empty on this release, guards `items[0]` against the empty list, and states PD28's 20-request bound. The static test asserts the text does NOT contain `--grant Human.Approve=`.

**PD33 (iteration 256) — the C7 route regexp must stop at whitespace.** In the held text the listing URL ends a line, and `[^ |"']+` matches newlines, so it would swallow the next line. Use `http://127\.0\.0\.1:7644(/v1/objects/by-semantic-id/[^\s|"')]+)`.

---

## Census — every enumerating guard this change can touch (grep-derived) and its value after the sprint

| Guard (file) | Base | After | Milestone |
|---|---|---|---|
| verify_ail.sh Leg 1b `LEG1_MODULES` (`:174`) | 22 modules | **23** (`design_docs/sketches/decisionpacket.ail` added) | C1 |
| verify_ail.sh `REQUIRED_VERIFIED` (`:321`) | 13 keys | **14** (decisionpacket: 9 names) | C1 |
| verify_ail.sh `EXACT_TOTAL_VERIFIED` (world/ only, `:403`) | 16 | 16 (sketches excluded) | — |
| verify_ail.sh Leg 2 `EXACT_TOTAL_TESTS` (`world/` dir, `:477`) | 40 | 40 | — |
| verify_ail.sh Leg 2b se-tools / Leg 2c execargs | 444 / 138 | unchanged | — |
| verify_ail.sh banner | `16 … 40 … 444` | unchanged; Leg-1b line `(23 modules)` | C1 |
| verifygate `module_manifest_gate_test.go` :96-97 / :221 / :231-232 / :335-336 / :561-584 | 24 / 22 / "22 module(s)" / 24 / 22 / 22-line literal | **25 / 23 / "23 module(s)" / 25 / 23 / 23-line literal** (PD3) | C1 |
| verifygate `floor_raise_inventory_test.go` (six-site S8 inventory) | pass | unchanged (no world/ identity) | — |
| TR.C `TestRegistryDispatchBindingBoundary` (`invoke_boundary_test.go:247`, `wantCount = 3` at :274; subtests `outside_broker_is_clean` requires `len(outside) == 0`, `inside_broker_exemption` requires exactly the 3 `publish_op.go` identities) | outside 0, inside 3 | **unchanged: `wantCount` stays 3, outside stays 0, and `invoke_boundary_test.go` is not edited.** New production code has no `.Invoke(` selector and no `broker.Session/NewSession/NewReplaySession`. With C4.2 HELD no 4th binding site lands at all — not even the `reg[...] = d.human` map insertion (which the detector could not see anyway: V45; Future Work #5, the existing `Workspace.Exec` binding has the same shape) | C3, C4.1, C5 |
| Production packet writers (PD29) | 0 (no producer) | **0** with D4 held: the producer exists, but the one production coordinator's binder has no `Human.Approve`, and REST refuses the label | C3, C4.1 |
| AC10 subprocess census (`registry_publish_test.go:1135` `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`; verifygate `TestSubprocessSinkGate`) | current site set | unchanged (no `exec.Command`) | — |
| `TestNoNewDeadlineFreeStoreReads` (`host/store/context_read_test.go:377`, pins `{}`) | 0 | 0 (every new read takes the caller's ctx) | C3, C5 |
| `TestProductionContextRoots` (`context_roots_test.go:71`, 27 named roots) | 27 | 27 (no new `context.Background/TODO/WithoutCancel`) | all |
| `TestProductionGoSurface` | host/cmd scan | +3 production files under host/; no allowlist change | C1, C5 |
| `TestBrokerTestsDoNotCallParallel` (`handlers_parallel_guard_test.go:15`) | 0 calls | 0 (the new broker test has no `t.Parallel`) | C3 |
| Boundary `TestWorldBoundaryDependencyAllowlist`, `TestWorkbenchPackageRemainsTransportFree` | pass | unchanged. host/workbench gains no import; coordinator gains only stdlib `io` | C2, C5 |
| `TestMCPMountMethodSource`, `TestMCPDaemonBudgetWiring` (parse daemon.go) | pass | unchanged (`daemon.go` is not modified: C4.2 HELD) | — |
| `TestWorkbenchNewViewTypesInCensus` (`live_test.go:123`; walks Live/Graph/Decisions only) | pass | unchanged (`PacketView` hangs off `Object`, which is not walked) | C5 |
| `TestSupportedWorkbenchQueryTruthTable` (query keys) | 7 keys | 7 (no new key) | C5 |
| Graph goldens (`host/workbench/testdata/*.golden`) | 2 | unchanged | C5 |
| Workspace exact-key tests (`workspace_test.go:203-213`) | `registry()` keys | unchanged (neither `registry()` nor `d.binder` changes: C4.2 HELD) | — |
| `TestEffectPlanSketchDrift` (effectplan, ≥ 50 rows) | pass | unchanged | — |
| **new** `TestDecisionPacketSketchRows` | — | 9 laws, exactly 37 rows | C1 |
| `TestPublicAuthoritySurfaceIsFrozen` (host/evidence) + verify_go `REQUIRED_EVIDENCE_TESTS` (64) | 64 | unchanged (host/evidence untouched) | — |
| QUICKSTART parsers: §0 box, §8 payloads (3), §9 commands, §10 serve line + refusal table, §11 serve lines (=1, cut to EOF), AGUI §10, live-workbench §11 control | pass | unchanged (PD15 constraints) | C7 |
| `descriptorBaselineSHA` (se-tools), PUB011 ratchet (`uncontracted_exports=3`) | pass | unchanged (no package change) | — |

---

## Frozen / must-not-change

- No hunk in: `world/**`, `packages/**`, `tools/**`, `.github/**`, `scripts/verify_go.sh`, `host/evidence/**`, `host/store/**`, `host/boundary/**`, `host/transitionreg/**`, `host/projection/**`, `host/capsule/**`, `host/replay/**`, `cmd/**`.
- No hunk in: `host/broker/{broker.go, confined.go, binder.go, record.go, approvals_view.go, invoke_boundary_test.go}`, `host/daemon/{daemon.go, workspace.go, workbench_grade.go}`, `host/coordinator/effectplan.go`. (`host/daemon/daemon.go` is frozen because C4.2 is HELD, `D-WORLD-73`.)
- No hunk in any existing `_test.go` **except** `host/verifygate/module_manifest_gate_test.go` (C1, the PD3 sites only). New tests REUSE existing helpers: `newHandlerDaemon`, `requestRecorder`, `requestRecorderAuth`, `authHeader`, `testCommit`, `encodeCommit`, `seedGenesisEmbedded`, `assertErrorClass`, `workbenchRegion`, `gradePage`, `newGradeDaemon`, `mintGradeEnvelope`, `gradeReport`, `storeGradeEnvelope`, `storeGradeSubject`, `gradeTestIdentity`, `publishTools`, `wsTool`, `wsCaps`, `boundedTestContext`, and in coordinator `fxRig`, `phaseRunner`, `probe`, `readPlan`, `r.call`, `r.fxCoordinator`, `r.replayBinder`, `r.recordV2`, `r.effectRecord`, `now`, `accessScope`, `parseSketchTests`, `argInt`, `argString`, `argBool`, `argAt`.
- The design doc and `design_docs/verification/world-row172-design/*` must not change. The sketch lands as a byte-identical copy (`cmp` rc=0).
- `TestWorkbenchObjectGrade` (`workbench_test.go:974`) stays green with its file unmodified.
- **Final changed-file set (C7) — exactly these 21 paths, D4 held:**
  - `design_docs/sketches/decisionpacket.ail` (new) — run A
  - `scripts/verify_ail.sh` — run A
  - `host/verifygate/module_manifest_gate_test.go` — run A
  - `host/coordinator/{plan.go, effectful.go, coordinator.go, decisionpacket.go (new), decisionpacket_drift_test.go (new), decisionpacket_test.go (new), decisionpacket_producer_test.go (new)}` — run A
  - `host/broker/{decide.go, approve.go, packet_exports_test.go (new)}` — run A
  - `host/coordinator/decisionpacket_validate.go (new)` — run B
  - `host/daemon/{handlers.go, workbench.go, workbench_packet.go (new), decisionpacket_daemon_test.go (new), workbench_packet_test.go (new)}` — run B
  - `host/workbench/render.go` — run B
  - `docs/QUICKSTART.md` — run B
  - **NOT in the set:** `host/daemon/daemon.go` (C4.2 HELD). Check: `git diff --quiet -- host/daemon/daemon.go` rc=0 at every boundary.
  - Check command: `git status --porcelain --untracked-files=all | /usr/bin/grep -v '^?? .claude/' | /usr/bin/grep -v 'sprint_w-decision-packet-production-projection.json\|w-decision-packet-production-projection-sprint-plan.md' | awk '{print $2}' | sort` must print exactly the 21 paths above.

## Test write rule

- Every DB and fixture goes under the test's `t.TempDir()` (via `newHandlerDaemon` and `newRig`). Tests never write tracked files, `$HOME` or `~/.ailang`.
- No `t.Parallel()` and no `time.Sleep`.
- Store reads in tests use `boundedTestContext(t)` (the store refuses deadline-free contexts).
- Each named test is written BEFORE its implementation. On the pre-milestone tree it is build-red, which is NOT a kill. The red-first proof is the milestone's mutation drill on compiling code.

## Mutation protocol

1. Apply ONE mutant alone to the boundary tree. Back up the file first: `cp F "$TMPDIR/F.orig"`.
2. The mutant must compile: `go vet ./host/...` rc=0. A build failure is NOT a kill. Exception: MS-mutants, which are `.ail` edits checked by verify_ail.sh.
3. Run the mutant's exact command. The named test must FAIL on its named assertion. These are NOT kills: a build failure, a panic in unrelated code, a SKIP, zero tests matched, an unrelated failure, or a harness timeout.
4. Restore without git: `cp "$TMPDIR/F.orig" F && cmp F "$TMPDIR/F.orig"` must give rc=0. Never use `git checkout --`.
5. Re-run the command and require green.

Ledger row format: `mutant | file:line | command | rc | failing assertion text | restore cmp rc | green rc`. A SURVIVED mutant blocks the milestone. C7 re-drills every run-B mutant once more on the final tree (PD31); run-A mutants are carried by the run-A sha256 manifest. Apply each edit with a short `python3` replace that asserts the old text occurs exactly once (`assert s.count(old) == 1`), never with a hand edit, so the mutant is exactly the named one.

## Gates (PD30): a fast gate at every milestone, a full gate at each run boundary

Write each gate as a script file under `$TMPDIR` and run it with `bash` (the Bash tool's shell is zsh). Each command must finish within the tool's 600 s cap; run every block below as its own foreground call and record each rc and wall time (`s=$(date +%s); …; echo rc=$? wall=$(( $(date +%s)-s ))`). Never background a gate.

**Fast gate (every milestone).** `NEW` is that milestone's new test names joined by `|`. `PKGS` is the package list the milestone names. `SUBSET` is the milestone's `-run` regex for the broader package run (or empty to run the whole package).

```bash
#!/usr/bin/env bash
set -u
export AILANG_BIN=$HOME/.pinned-ailang/ailang; export PATH=/opt/homebrew/bin:$PATH
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter255-row172-design
NEW='…'; PKGS='…'; SUBSET='…'                                              # set per milestone
go vet ./... || exit 1
go test -run '^$' ./... >/dev/null || exit 1                               # compile fence incl. _test.go (NOT go build)
log="$(mktemp)"; go test -count=1 -v $PKGS -run "^($NEW)$" >"$log" 2>&1; rc=$?; tail -40 "$log"; [ $rc -eq 0 ] || exit 1
for t in $(printf '%s' "$NEW" | tr '|' ' '); do /usr/bin/grep -q -- "^--- PASS: $t " "$log" || { echo "$t did not PASS"; exit 1; }; done
if [ -n "$SUBSET" ]; then go test -count=1 $PKGS -run "$SUBSET" || exit 1; else go test -count=1 $PKGS || exit 1; fi
git diff --quiet -- world packages tools .github scripts/verify_go.sh host/evidence host/store host/boundary host/transitionreg host/projection host/capsule host/replay cmd host/daemon/daemon.go host/broker/invoke_boundary_test.go design_docs/planned/w-decision-packet-production-projection.md design_docs/verification || { echo "frozen tree touched"; exit 1; }
git diff --quiet -- host/daemon/workbench_test.go host/daemon/handlers_test.go host/daemon/workbench_grade_test.go host/daemon/workspace_test.go host/coordinator/effectful_test.go host/coordinator/effectplan_drift_test.go || { echo "existing test touched"; exit 1; }
echo FAST GATE GREEN
```

**Run-boundary full gate (after C3 = end of run A; after C7 = end of run B).** Four separate foreground calls:

```bash
# F1 (~30 s)
bash -c 'export AILANG_BIN=$HOME/.pinned-ailang/ailang PATH=/opt/homebrew/bin:$PATH; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter255-row172-design && ./scripts/verify_ail.sh'
# F2 (~255 s; verifygate dominates)
bash -c 'export AILANG_BIN=$HOME/.pinned-ailang/ailang PATH=/opt/homebrew/bin:$PATH; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter255-row172-design && go test ./host/... ./cmd/... -count=1 -timeout 9m'
# F3 (-race; run A: RACE='./host/coordinator ./host/broker' ~165 s; run B: RACE='./host/coordinator ./host/daemon ./host/workbench' ~225 s)
bash -c 'export AILANG_BIN=$HOME/.pinned-ailang/ailang PATH=/opt/homebrew/bin:$PATH; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter255-row172-design && go test -race $RACE -count=1 -timeout 9m'
# F4
bash -c 'export AILANG_BIN=$HOME/.pinned-ailang/ailang PATH=/opt/homebrew/bin:$PATH; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter255-row172-design && go test -count=1 ./host/daemon -run "^TestWorkbenchObjectGrade$" -v | /usr/bin/grep -q "^--- PASS: TestWorkbenchObjectGrade " && git diff --quiet -- host/daemon/workbench_test.go && echo F4 GREEN'
```

A red F2/F3 in a package this sprint does not touch is re-run ONCE alone (`go test -count=1 ./host/<pkg>`); a second red is a STOP-and-report, never a skip. Known base flakes (`TestExecSrtMCP*`, `TestCommitInvocationReceiptIdentity` under `-race`) are reported by name, not fixed.

**Run-A sha256 manifest (PD30).** At the end of run A: `shasum -a 256 design_docs/sketches/decisionpacket.ail scripts/verify_ail.sh host/verifygate/module_manifest_gate_test.go host/coordinator/{plan.go,effectful.go,coordinator.go,decisionpacket.go,decisionpacket_drift_test.go,decisionpacket_test.go,decisionpacket_producer_test.go} host/broker/{decide.go,approve.go,packet_exports_test.go} > ~/.ailang/state/iter256-plan/runA.sha256`. Run B begins with `shasum -a 256 -c ~/.ailang/state/iter256-plan/runA.sha256` (all `OK`, rc=0) and ends with the same check.

---

## Code shape (binding names)

**`host/coordinator/plan.go`**: add one line to the `const` block after `RecordV2`:
```go
	DecisionPacketV1 = "world/decision-packet/v1" // row 172: minted only by dispatchEffectful
```

**`host/coordinator/decisionpacket.go`** (new, package coordinator; C1 adds the laws, C2 the codec, C3 the producer helpers). Imports: `bytes`, `encoding/json`, `errors`, `fmt`, `io`, `host/broker`, `host/hashref`, `host/store`. Add imports only as each milestone uses them.

```go
// C1 — Go mirrors of design_docs/sketches/decisionpacket.ail (bound by decisionpacket_drift_test.go)
type lawRef struct{ algo, digest string }
func refOf(h hashref.HashRef) lawRef { return lawRef{h.Algo(), h.Digest()} }
func policyCode(token string) int        // "Cancel"→0, "EscalateBounded"→1, "ExecuteIfGranted"→2, else -1
func scheduleWellFormed(createdAt, deadlineAt, escalationsRemaining int64) bool // createdAt >= 0 && deadlineAt > createdAt && escalationsRemaining >= 0
func mintable(now, grantExpiresAt int64) bool                    // now >= 0 && now < grantExpiresAt
func mintedScheduleWellFormed(now, grantExpiresAt int64) bool    // scheduleWellFormed(now, grantExpiresAt, 0)
func producerFieldsOk(policy int, escalationsRemaining int64) bool // policy == 0 && escalationsRemaining == 0
func sameRef(a, b lawRef) bool                                    // a.algo == b.algo && a.digest == b.digest
func packetRefusal(bytesOk, carriedOk, proposalOk, requestOk, anchoredOk, lifecycleFact, pairUniqueOk bool) int // 1..7 first false, else 0
func samePair(aRequest, aProposal, bRequest, bProposal lawRef) bool // sameRef(aRequest,bRequest) && sameRef(aProposal,bProposal)
func lifecycleOk(createdAt, requestNow, deadlineAt, escalationsRemaining int64, policy int) bool
	// createdAt == requestNow && scheduleWellFormed(createdAt, deadlineAt, escalationsRemaining) && producerFieldsOk(policy, escalationsRemaining)

// C2 — the v1 wire (D2; IRREVERSIBLE freeze D2b)
type DecisionPacket struct {
	PacketHash, ProposalHash, RequestRef      hashref.HashRef
	CreatedAt, DeadlineAt, EscalationsRemaining int64
	Policy                                    string
}
type packetWire struct { // member order IS the canonical order
	PacketHash           string `json:"packetHash"`
	ProposalHash         string `json:"proposalHash"`
	RequestRef           string `json:"requestRef"`
	CreatedAt            int64  `json:"createdAt"`
	DeadlineAt           int64  `json:"deadlineAt"`
	EscalationsRemaining int64  `json:"escalationsRemaining"`
	Policy               string `json:"policy"`
}
type packetBodyWire struct { /* packetWire minus PacketHash, same order and tags */ }
var errPacketEncoding = errors.New("coordinator: not a canonical world/decision-packet/v1 encoding")
var errPacketSchedule = errors.New("coordinator: the authorizing grant does not admit a well-formed packet schedule")
func packetBodyHash(p DecisionPacket) hashref.HashRef     // SumSHA256(mustJSON(packetBodyWire{…}))
func EncodeDecisionPacket(p DecisionPacket) []byte         // sets p.PacketHash = packetBodyHash(p); returns mustJSON(packetWire{…})
func DecodeDecisionPacket(payload []byte) (DecisionPacket, error)
func mintPacket(now, grantExpiresAt int64, proposal, request hashref.HashRef, by string) (store.Object, error)

// C3 — producer
var errPacketGrant = errors.New("coordinator: no seeded grant authorizes the executed Human.Approve request")
func mintApprovalPacket(grants []broker.Capability, e PlannedEffect, res EffectResult, now int64, planRef hashref.HashRef, surface Surface) (store.Object, error)
func withPacket(p plan, obj store.Object) plan
```

**`host/coordinator/decisionpacket_validate.go`** (new, C5): `PacketReader`, `PacketResult`, `const packetScanLimit = 100`, `func ValidatePacket(ctx context.Context, r PacketReader, obj store.Object) (PacketResult, error)`.

**`host/broker`** (C3): `decide.go`: `func AuthorizingGrant(grants []Capability, r EffectRequest) (Capability, bool)`. `approve.go`: `type ApprovalRequest struct{ Effect, Scope string; Cost int64; Requester string; Now int64 }`, `func ParsePendingApproval(output []byte) (hashref.HashRef, error)`, `func DecodeApprovalRequest(payload []byte) (ApprovalRequest, error)`.

**`host/daemon`** (C4.1, C5):
- `daemon.go`: **HELD (C4.2, `D-WORLD-73`) — not modified this sprint.** (Held shape, for the follow-up: `Daemon.humanMu sync.Mutex`, `Daemon.human broker.Handler`, `type serializedHandler struct{ mu *sync.Mutex; h broker.Handler }` with `Execute`, and the `binder` change.)
- `handlers.go`: the reservation.
- `workbench_packet.go` (new): sentence constants, `packetRefusalSentence`, `(d *Daemon) packetView`.
- `workbench.go`: the branch.

**`host/workbench/render.go`** (C5): `type PacketView struct{…}`, `ObjectView.Packet *PacketView`, template partials `packet` and `packetGrade`.

---

# EXECUTOR RUN A — C1, C2, C3, then the run-A boundary full gate (~45 min)

## C1 — Sketch, gate manifest, Go law mirrors, drift (design M1, laws half)

**Budget:** prod ~45 Go LOC + 13 shell lines; tests ~110 LOC new + ~8 changed lines in verifygate. **~0.4 d.**

Files:
- `design_docs/sketches/decisionpacket.ail` (new)
- `scripts/verify_ail.sh`
- `host/verifygate/module_manifest_gate_test.go`
- `host/coordinator/decisionpacket.go` (new: laws only)
- `host/coordinator/decisionpacket_drift_test.go` (new)

### Tests first

`TestDecisionPacketSketchRows` (AC9) in `decisionpacket_drift_test.go`, package coordinator:
1. Read `filepath.Join("..", "..", "design_docs", "sketches", "decisionpacket.ail")`. Require the prefix `"module sketches/decisionpacket\n"`.
2. `laws := parseSketchTests(t, string(src))` (the existing parser in effectplan_drift_test.go, same package).
3. Define `var packetSketchMirrors = map[string]func(t *testing.T, a []any) any{…}` with exactly these 9 keys:
   - `"policyCode": return int64(policyCode(argString(t, a, 0)))`
   - `"scheduleWellFormed": return scheduleWellFormed(argInt(t,a,0), argInt(t,a,1), argInt(t,a,2))`
   - `"mintable": return mintable(argInt(t,a,0), argInt(t,a,1))`
   - `"mintedScheduleWellFormed": return mintedScheduleWellFormed(argInt(t,a,0), argInt(t,a,1))`
   - `"producerFieldsOk": return producerFieldsOk(int(argInt(t,a,0)), argInt(t,a,1))`
   - `"sameRef": return sameRef(argLawRef(t,a,0), argLawRef(t,a,1))`
   - `"packetRefusal": return int64(packetRefusal(argBool(t,a,0), …, argBool(t,a,6)))`
   - `"samePair": return samePair(argLawRef(t,a,0), argLawRef(t,a,1), argLawRef(t,a,2), argLawRef(t,a,3))`
   - `"lifecycleOk": return lifecycleOk(argInt(t,a,0), argInt(t,a,1), argInt(t,a,2), argInt(t,a,3), int(argInt(t,a,4)))`
4. Add helper `argLawRef(t, a, i) lawRef`. It requires `argAt(t,a,i).(map[string]any)` with exactly 2 keys, `"algo"` and `"digest"`, both strings. Otherwise `t.Fatalf`.
5. Sorted name sets must be equal. On mismatch: `t.Fatalf("sketch laws = %v\nGo mirrors = %v\nevery sketch law needs a Go mirror and vice versa", …)`.
6. For every row: `got := packetSketchMirrors[name](t, row.args)`. If `got != row.want`, use `t.Errorf("sketch line %d %s%v: Go mirror = %v, sketch expects %v", …)`.
7. Shape rules (PD4):
   - Int laws (`policyCode`, `packetRefusal`) need ≥ 2 distinct `want` values.
   - Bool laws need a true and a false row, except `mintedScheduleWellFormed`, which needs ≥ 2 rows all `true`. Keep that exemption in a one-entry map named `implicationLaws`.
8. Final check: `if total != 37 { t.Fatalf("parsed %d sketch rows; the frozen sketch carries exactly 37", total) }`, then `t.Logf("evaluated %d sketch rows over %d laws", total, len(names))`.

### Implementation

1. `cp design_docs/verification/world-row172-design/decisionpacket.ail.txt design_docs/sketches/decisionpacket.ail && cmp design_docs/verification/world-row172-design/decisionpacket.ail.txt design_docs/sketches/decisionpacket.ail` must give rc=0.
2. `scripts/verify_ail.sh`, three edits, all in this one milestone:
   - **(a)** In `LEG1_MODULES=(`, insert `  design_docs/sketches/decisionpacket.ail` as the FIRST line, directly above `  design_docs/sketches/effectbroker.ail`.
   - **(b)** In the Python `REQUIRED_VERIFIED` dict, insert immediately after the execargs entry's closing line (`"consumesNext", … "argCountOk"},`):
     ```python
         # The decision-packet law (queue row 172 M1): Go mirrors in host/coordinator/decisionpacket.go,
         # bound by host/coordinator/decisionpacket_drift_test.go. A sketch, so excluded from the
         # world/ total: gated by name without moving the S8 floor.
         "design_docs/sketches/decisionpacket.ail": {
             "policyCode", "scheduleWellFormed", "mintable", "mintedScheduleWellFormed",
             "producerFieldsOk", "sameRef", "packetRefusal", "samePair", "lifecycleOk"},
     ```
   - **(c)** Comment at `:156-157`: replace `except design_docs/sketches/effectplan.ail and design_docs/sketches/execargs.ail,` with `except design_docs/sketches/effectplan.ail, design_docs/sketches/execargs.ail and` and the next line `# whose law proofs are gated by name.` with `# design_docs/sketches/decisionpacket.ail, whose law proofs are gated by name.` Verifygate has no anchor on this comment (the planner grepped it).
3. `host/verifygate/module_manifest_gate_test.go` (PD3), exactly these edits:
   - `:96` `files != 24 || ailFiles != 22` → `files != 25 || ailFiles != 23`, and `:97` message `want 24 / 22` → `want 25 / 23`
   - `:221` `across 22 module(s)` → `across 23 module(s)`
   - `:231` `files != 24` → `files != 25`, and `:232` `want 24` → `want 25`
   - `:335` `got != 22` → `got != 23`, and `:336` `want 22` → `want 23`
   - In the `const old = \`LEG1_MODULES=(` literal (`:561`), insert the line `  design_docs/sketches/decisionpacket.ail` before `  design_docs/sketches/effectbroker.ail`. The literal must equal the script's block byte-for-byte.
4. `host/coordinator/decisionpacket.go` contains the C1 laws exactly as in "Code shape". `packetRefusal` is written as the sketch's if-chain:
   ```go
   switch { case !bytesOk: return 1; case !carriedOk: return 2; case !proposalOk: return 3; case !requestOk: return 4; case !anchoredOk: return 5; case !lifecycleFact: return 6; case !pairUniqueOk: return 7 }; return 0
   ```

### Gates

```bash
NEW='TestDecisionPacketSketchRows'; PKGS='./host/coordinator'; SUBSET='Drift|Sketch'
# fast gate, plus (separate foreground calls):
./scripts/verify_ail.sh                                                    # ~30 s
go test -count=1 ./host/verifygate -run '^(TestModuleManifestEmptyAllowlistFailsLoudly|TestModuleManifestRejectsStrayModule|TestPristineControlRefusesANonIdenticalRoot)$' -v   # ~2–4 min; together these exercise all five PD3 sites
```

verify_ail.sh must print `swept .ail module set equals the LEG1_MODULES allowlist (23 modules)`, a line `ai-check design_docs/sketches/decisionpacket.ail`, and the unchanged banner `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 444 se-tools named tests pass` (pre-measured in a scratch copy, R8). Red-first evidence for C1: before step 2(b) is applied, MS4's edit is NOT caught (record that run: with no `REQUIRED_VERIFIED` entry the vanished `samePair` contract leaves verify_ail.sh green), and after 2(b) it is (rc=1, the R9 line). This is the proof the by-name entry is load-bearing.

### Mutations (C1)

| ID | Edit | Command | Expected red |
|---|---|---|---|
| MS1 | sketch `  scheduleWellFormed(now, grantExpiresAt, 0)` → `  scheduleWellFormed(now, grantExpiresAt, 0 - 1)` | `./scripts/verify_ail.sh` | rc=1, `✗ design_docs/sketches/decisionpacket.ail: verify.counterexample == 1` (planner-measured: cex 1 on `mintedScheduleWellFormed`) |
| MS2 | sketch `  else if not lifecycleOk then 6` → `  else if not lifecycleOk then 0` | same | rc=1, counterexample == 1 (measured: `packetRefusal`) |
| MS3 | Go `policyCode`: `case "Cancel":` → `case "Cancel", "cancel":` | `go test -count=1 -v ./host/coordinator -run '^TestDecisionPacketSketchRows$'` | `sketch line 28 policyCode[cancel]: Go mirror = 0, sketch expects -1` |
| MS4 (PD17) | sketch: delete the two `samePair` `ensures {…}` lines (`:139-140`) | `./scripts/verify_ail.sh` | rc=1, `required identity (decisionpacket, samePair) MISSING from verify.results[]` |
| MP8b-drift | Go `samePair` body → `return sameRef(aRequest, bRequest)` | drift command | `sketch line 144 samePair…: Go mirror = true, sketch expects false` |
| MV7-drift | Go `lifecycleOk`: drop `createdAt == requestNow &&` | drift command | `sketch line 161 lifecycleOk[10 11 20 0 0]: Go mirror = true, sketch expects false` (line 161 measured R10; iteration 255 said 162) |

Sketch mutants edit the landed `.ail` with a `cp` backup. Line numbers are those of the banked file, which is identical.

---

## C2 — Canonical codec and `mintPacket` (design M1, codec half; D2 freeze)

**Budget:** prod ~95; tests ~120. **~0.4 d.**

Files:
- `host/coordinator/plan.go` (1 const)
- `host/coordinator/decisionpacket.go` (codec)
- `host/coordinator/decisionpacket_test.go` (new)

### Tests first

`TestDecisionPacketCodec` (PD19; covers D2 and the AC2 recompute) is table-driven:
- **golden:** `p := DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan")), RequestRef: hashref.SumSHA256([]byte("req")), CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"}`. Set `body := fmt.Sprintf(\`{"proposalHash":"%s","requestRef":"%s","createdAt":1700,"deadlineAt":1760,"escalationsRemaining":0,"policy":"Cancel"}\`, P, Q)` and `want := fmt.Sprintf(\`{"packetHash":"%s",`, hashref.SumSHA256([]byte(body))) + body[1:]`. Require `string(EncodeDecisionPacket(p)) == want` (the D2 bytes, byte-exact).
- **round trip:** `DecodeDecisionPacket(EncodeDecisionPacket(p))` equals p with `PacketHash` = SumSHA256(body).
- **refusals** (each `errors.Is(err, errPacketEncoding)`):
  - whitespace (`strings.Replace(enc, ",\"createdAt\"", ", \"createdAt\"", 1)`)
  - unknown member (`strings.Replace(enc, "}", ",\"x\":1}", 1)` on the last brace)
  - trailing value (`enc + "{}"`)
  - missing member (built by `json.Marshal` of a map without `escalationsRemaining`)
  - malformed ref (`proposalHash` = `"sha256:zz"`)
  - unknown policy `"cancel"`
  - wrong packetHash (replace its value with `hashref.SumSHA256([]byte("wrong")).String()`)
- **mintPacket:**
  - `mintPacket(1700, 1760, P, Q, "coordinator:a2a")` gives an object with `SemanticID == DecisionPacketV1`, `InterfaceHash == hashref.SumSHA256([]byte(DecisionPacketV1))`, `Provenance == "coordinator:a2a"`, `Hash == SumSHA256(Payload)`, and a decoded `{CreatedAt 1700, DeadlineAt 1760, EscalationsRemaining 0, Policy "Cancel"}`.
  - `mintPacket(1760, 1760, …)` and `mintPacket(-1, 5, …)` return `errors.Is(err, errPacketSchedule)`.

### Implementation

```go
func packetBodyHash(p DecisionPacket) hashref.HashRef {
	return hashref.SumSHA256(mustJSON(packetBodyWire{ProposalHash: p.ProposalHash.String(), RequestRef: p.RequestRef.String(),
		CreatedAt: p.CreatedAt, DeadlineAt: p.DeadlineAt, EscalationsRemaining: p.EscalationsRemaining, Policy: p.Policy}))
}

// EncodeDecisionPacket returns the canonical world/decision-packet/v1 payload
// (row 172 D2: compact JSON in the frozen member order; packetHash is
// SumSHA256 of the same object without packetHash). Amendments are /v2.
func EncodeDecisionPacket(p DecisionPacket) []byte {
	p.PacketHash = packetBodyHash(p)
	return mustJSON(packetWire{PacketHash: p.PacketHash.String(), ProposalHash: p.ProposalHash.String(),
		RequestRef: p.RequestRef.String(), CreatedAt: p.CreatedAt, DeadlineAt: p.DeadlineAt,
		EscalationsRemaining: p.EscalationsRemaining, Policy: p.Policy})
}

func DecodeDecisionPacket(payload []byte) (DecisionPacket, error) {
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	var w packetWire
	if err := dec.Decode(&w); err != nil {
		return DecisionPacket{}, fmt.Errorf("%w: %v", errPacketEncoding, err)
	}
	var extra json.RawMessage
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return DecisionPacket{}, fmt.Errorf("%w: trailing data", errPacketEncoding)
	}
	var refs [3]hashref.HashRef
	for i, text := range []string{w.PacketHash, w.ProposalHash, w.RequestRef} {
		ref, err := hashref.Parse(text)
		if err != nil {
			return DecisionPacket{}, fmt.Errorf("%w: %v", errPacketEncoding, err)
		}
		refs[i] = ref
	}
	if policyCode(w.Policy) < 0 {
		return DecisionPacket{}, fmt.Errorf("%w: unknown policy %q", errPacketEncoding, w.Policy)
	}
	p := DecisionPacket{PacketHash: refs[0], ProposalHash: refs[1], RequestRef: refs[2], CreatedAt: w.CreatedAt,
		DeadlineAt: w.DeadlineAt, EscalationsRemaining: w.EscalationsRemaining, Policy: w.Policy}
	if !bytes.Equal(mustJSON(w), payload) { // MV1 target line
		return DecisionPacket{}, fmt.Errorf("%w: not the canonical byte form", errPacketEncoding)
	}
	if packetBodyHash(p) != p.PacketHash { // MV2 target line
		return DecisionPacket{}, fmt.Errorf("%w: packetHash does not recompute", errPacketEncoding)
	}
	return p, nil
}

func mintPacket(now, grantExpiresAt int64, proposal, request hashref.HashRef, by string) (store.Object, error) {
	if !mintable(now, grantExpiresAt) || !mintedScheduleWellFormed(now, grantExpiresAt) {
		return store.Object{}, errPacketSchedule
	}
	payload := EncodeDecisionPacket(DecisionPacket{ProposalHash: proposal, RequestRef: request, CreatedAt: now,
		DeadlineAt: grantExpiresAt, EscalationsRemaining: 0, Policy: "Cancel"}) // MP6 target: "Cancel"
	return object(by, DecisionPacketV1, payload), nil
}
```

### Gates

Fast gate with `NEW='TestDecisionPacketCodec|TestDecisionPacketSketchRows'`, `PKGS='./host/coordinator'`, `SUBSET=''` (whole coordinator package, ~65 s). Red-first evidence: write `TestDecisionPacketCodec` first and record its build-red run (`undefined: EncodeDecisionPacket`), then implement; the real red-first proof is the C2 mutation ledger.

### Mutations (C2)

All use `go test -count=1 -v ./host/coordinator -run '^TestDecisionPacketCodec$'`.

| ID | Edit | Expected red |
|---|---|---|
| MV1 | `if !bytes.Equal(mustJSON(w), payload) {` → `if false && !bytes.Equal(mustJSON(w), payload) {` | the whitespace arm decodes without error |
| MV2 | `if packetBodyHash(p) != p.PacketHash {` → `if false && packetBodyHash(p) != p.PacketHash {` | the wrong-packetHash arm decodes |
| MP7 (PD21) | in `EncodeDecisionPacket` only: `p.PacketHash = packetBodyHash(p)` → `p.PacketHash = hashref.SumSHA256(mustJSON(packetWire{Policy: p.Policy}))` | golden bytes differ, and the round trip fails `packetHash does not recompute` |
| MP6 | in `mintPacket`: `Policy: "Cancel"` → `Policy: "EscalateBounded"` | the mintPacket arm reports `Policy = EscalateBounded, want Cancel` |

---

## C3 — Broker exports and the coordinator producer (design M2, coordinator half)

**Budget:** prod ~75 (broker ~35, coordinator ~40); tests ~230. **~0.6 d.**

Files:
- `host/broker/decide.go`, `host/broker/approve.go`, `host/broker/packet_exports_test.go` (new, `package broker`, no `t.Parallel`)
- `host/coordinator/effectful.go`, `host/coordinator/coordinator.go`, `host/coordinator/decisionpacket.go`
- `host/coordinator/decisionpacket_producer_test.go` (new)

### Broker implementation (exact)

```go
// decide.go — after Allows:
// AuthorizingGrant returns the grant decideOver selects for r when it allows
// r (row 172 D3: the packet deadline is that grant's expiry); ok is false when
// no grant allows r.
func AuthorizingGrant(grants []Capability, r EffectRequest) (Capability, bool) {
	i, d := decideOver(grants, r)
	if !d.Allowed || i < 0 {
		return Capability{}, false
	}
	return grants[i], true // MB1 target
}
```

```go
// approve.go — after pendingWire:
// ApprovalRequest is a decoded world/approval-request/v1 object.
type ApprovalRequest struct {
	Effect, Scope string
	Cost          int64
	Requester     string
	Now           int64
}

// ParsePendingApproval returns the request ref of a Human.Approve handler
// output ({"status":"pending","requestRef":…}).
func ParsePendingApproval(output []byte) (hashref.HashRef, error) {
	var w pendingWire
	if err := decodeApprovalJSON(output, &w); err != nil {
		return hashref.HashRef{}, err
	}
	if w.Status != "pending" {
		return hashref.HashRef{}, fmt.Errorf("broker: approval output status %q, want pending", w.Status)
	}
	return hashref.Parse(w.RequestRef)
}

// DecodeApprovalRequest strictly decodes a world/approval-request/v1 payload.
func DecodeApprovalRequest(payload []byte) (ApprovalRequest, error) {
	var w approvalRequestWire
	if err := decodeApprovalJSON(payload, &w); err != nil {
		return ApprovalRequest{}, err
	}
	return ApprovalRequest{Effect: w.Effect, Scope: w.Scope, Cost: w.Cost, Requester: w.Requester, Now: w.Now}, nil
}
```

### Coordinator implementation (exact)

1. `effectful.go:273`: append the parameter `, grants []broker.Capability` after `d transitionreg.Descriptor`. `coordinator.go:334`: `return c.dispatchEffectful(ctx, call.Surface, id, call.EpisodeID, input, bound, d, grants)`.
2. `effectful.go`: between `if err != nil { return unrecorded(err) }` (after `composeEffectOutput`) and `p := planEffectInvocation(…)`, insert:
   ```go
   	// Row 172 D1: an allowed, successful Human.Approve mints the decision
   	// packet that rides in this invocation's own commit.
   	var packet *store.Object
   	if len(pl.Effects) == 1 && pl.Effects[0].Requirement.Effect == broker.EffectHumanApprove && results[0].Status == StatusOK { // MP1 target
   		obj, err := mintApprovalPacket(grants, pl.Effects[0], results[0], now, planRef, surface)
   		if err != nil {
   			return unrecorded(err)
   		}
   		packet = &obj
   	}
   ```
   Directly after the `p := planEffectInvocation(…)` line, insert:
   ```go
   	if packet != nil {
   		p = withPacket(p, *packet)
   	}
   ```
   (effectful.go already imports `store` and `broker`; verify with `go vet`.)
3. `decisionpacket.go`:
   ```go
   func mintApprovalPacket(grants []broker.Capability, e PlannedEffect, res EffectResult, now int64,
   	planRef hashref.HashRef, surface Surface) (store.Object, error) {
   	req := broker.EffectRequest{Effect: e.Requirement.Effect, Scope: e.Requirement.Scope, Cost: e.Requirement.Cost, Now: now}
   	grant, ok := broker.AuthorizingGrant(grants, req)
   	if !ok {
   		return store.Object{}, errPacketGrant // structurally unreachable after an allowed dispatch
   	}
   	request, err := broker.ParsePendingApproval(res.Output)
   	if err != nil {
   		return store.Object{}, err
   	}
   	return mintPacket(now, grant.ExpiresAt, planRef, request, writtenBy(surface)) // MP2/MP3/MP4 targets
   }

   func withPacket(p plan, obj store.Object) plan {
   	objects := append(append([]store.Object(nil), p.Objects...), obj)
   	p.Objects, p.Commit.Objects = objects, objects // MP5 target
   	return p
   }
   ```

### Test helpers (in `decisionpacket_producer_test.go`)

```go
var approveReq = transitionreg.EffectRequirement{Effect: broker.EffectHumanApprove, Scope: "demo", Cost: 1}
func approvePlan(id string) string { return `{"plan":"world/effect-plan/v1","effects":[{"id":"` + id + `","effect":"Human.Approve","scope":"demo","cost":1,"payload":{"requester":"coord-test"}}],"finish":false,"result":null}` }
var invokeGrant = broker.Capability{Effect: "Invoke", Scope: accessScope, ExpiresAt: 1000, Budget: 10}
func approveGrant(expires int64) broker.Capability { return broker.Capability{Effect: broker.EffectHumanApprove, Scope: "demo", ExpiresAt: expires, Budget: 5} }
type countingHandler struct { h broker.Handler; n *atomic.Int64 } // Execute: n.Add(1); return h.Execute(...)
func packetsIn(t *testing.T, r *rig) []DecisionPacket   // ObjectsBySemanticID(boundedTestContext(t), DecisionPacketV1, "", 50) → GetObject → DecodeDecisionPacket (t.Fatal on error)
func packetObjects(t *testing.T, r *rig) []store.Object
func approveRig(t *testing.T) (*rig, *Coordinator) // fxRig(t, true, approveReq); r.fxCoordinator(r.st, &phaseRunner{plan: approvePlan("e1")}, broker.Registry{broker.EffectHumanApprove: broker.NewHumanHandler(r.st)})
```

Use the existing `now` (100). It satisfies `0 <= 100 < 1760`.

### Tests first

- `TestPacketBrokerExports` (host/broker):
  - `AuthorizingGrant([expired(ExpiresAt 100), live 1760, live 1900], {Human.Approve, demo, 1, Now 100})` returns `(1760-grant, true)`, and this equals the first `Decide(...).Allowed` in order.
  - All-expired → false. Empty → false.
  - `ParsePendingApproval(pendingBytes(ref)) == ref`. Status `"approve"`, an extra field, and `requestRef:"sha256:zz"` are errors.
  - `DecodeApprovalRequest` of the object `NewHumanHandler(store).Execute(...)` stores gives `{Human.Approve, scope, cost, requester, now}`. A payload with an extra field is an error. Use the package's existing in-test store constructor `openTestStore(t)` (`broker_test.go:15`, an in-memory `store.Open(":memory:")`) and the package's `boundedTestContext(t)` (`m7b_deadline_test.go:10`).
- `TestProducerMintsPacketForApprovalOnly` (AC1), subtests, each with a fresh rig:
  - `approval`: caps `[invokeGrant, approveGrant(1760)]` → 1 packet.
  - `denied`: `[invokeGrant, approveGrant(now)]` → `Dispatch` err nil, `decodeWorld(...).World.Effects[0].Status == StatusDenied`, 0 packets.
  - `non-approval`: `fxRig(t, false)` + `readPlan(false)` + `probe` → 0 packets.
  - `refusal-plan`: the approve rig with plan `{"plan":"world/effect-plan/v1","effects":[],"finish":false,"result":{"ok":false,"refused":"no"}}` → 0 packets.
  - `two-approvals`: one rig, coordinators A (`approvePlan("e1")`) and B (`approvePlan("e2")`) over `r.st` with the same registry. Dispatch `t1` on A and `t2` on B with caps `[invokeGrant, approveGrant(1760)]`. Expect 2 packets, equal `RequestRef` (the V10 collapse) and different `ProposalHash`.
- `TestProducerPacketBindsInvocationFacts` (AC2, PD22):
  - caps `[invokeGrant, approveGrant(now), approveGrant(1760), approveGrant(1900)]`.
  - `wantDeadline` is computed by looping `broker.Decide(g, req)` over caps[1:] and taking the first Allowed. Assert it equals the literal 1760.
  - From the result: `rec := r.recordV2(res.RecordRef)`, `er := r.effectRecord(hashref.MustParse(rec.Effects[0]))`, then `GetObject(er.ResultRef)` → `broker.ParsePendingApproval` → `q`, then `GetObject(q)` → `broker.DecodeApprovalRequest` → `req`, and `rc, _, _ := r.st.GetReceipt(ctx, res.InvocationID)`.
  - Assert: `p.CreatedAt == req.Now == rc.Intent.LogicalTime == now`; `p.DeadlineAt == 1760`; `p.ProposalHash.String() == rec.Plan`; `p.RequestRef == q`; `p.Policy == "Cancel"`; `p.EscalationsRemaining == 0`. `packetsIn` decoding without error proves that packetHash recomputes.
- `TestProducerPacketIsCarriedByTheInvocationCommit` (AC3, coordinator arm):
  - `r.st.ObjectCommits(ctx, k.Hash, -1, 10)` deep-equals `[]int64{res.EntryIndex}`.
  - `r.st.ObjectCommits(ctx, res.RecordRef, -1, 10)[0] == res.EntryIndex`.
  - Two-approvals: each packet's first carrier is its own invocation's `EntryIndex`. Match packets to results via `ProposalHash == r.recordV2(res.RecordRef).Plan`.
- `TestReplayReconcileToleratesPacketObject` (AC3b, PD9):
  - Use the `countingHandler` registry.
  - `out, err := c.Replay(ctx, res.InvocationID, r.replayBinder())` → err nil and `string(out) == string(res.OutputBytes)`.
  - `again, err := c.Dispatch(ctx, r.call("fx", "t1", caps))` → `again.Reconciled`, equal bytes, handler count unchanged, packet count still 1.

### Gates

Fast gate with `NEW='TestPacketBrokerExports|TestProducerMintsPacketForApprovalOnly|TestProducerPacketBindsInvocationFacts|TestProducerPacketIsCarriedByTheInvocationCommit|TestReplayReconcileToleratesPacketObject'`, `PKGS='./host/broker ./host/coordinator'`, `SUBSET='Packet|Producer|Replay|Approv|Decide|Effectful|Dispatch'` (the full broker and coordinator packages run in the run-A boundary gate F2 right after). TR.C stays green with its file unmodified: `go test -count=1 ./host/broker -run '^TestRegistryDispatchBindingBoundary$' -v` → `--- PASS`.

### Run-A boundary full gate (end of run A)

Run F1–F4 (above) with `RACE='./host/coordinator ./host/broker'`. Then write the run-A sha256 manifest (PD30). Report: the C1–C3 ledgers, F1–F4 rc and wall times, and `git status --porcelain --untracked-files=all` (it must list exactly the 13 run-A paths, plus the two plan files if the controller has not committed them yet). **STOP: run A ends here.**

### Mutations (C3)

All use `go test -count=1 -v ./host/coordinator ./host/broker -run '^(TestPacketBrokerExports|TestProducerMintsPacketForApprovalOnly|TestProducerPacketBindsInvocationFacts|TestProducerPacketIsCarriedByTheInvocationCommit|TestReplayReconcileToleratesPacketObject)$'`.

| ID | Edit | Expected red |
|---|---|---|
| MP1 | drop ` && results[0].Status == StatusOK` | AC1 `denied`: `Dispatch` returns `*EffectsUnrecordedError` (errPacketGrant) |
| MP2 | `mintPacket(now, …` → `mintPacket(now+1, …` in `mintApprovalPacket` | AC2 `createdAt 101 != request now 100` |
| MP3 | `grant.ExpiresAt` → `now + 1` in the same call | AC2 `deadlineAt 101, want 1760` |
| MP4 | `planRef` → `hashref.SumSHA256(res.Output)` in the same call | AC2 `proposalHash != record plan` |
| MP5 (PD7) | `p.Objects, p.Commit.Objects = objects, objects` → `p.Objects = objects` | AC1 `approval` 0 packets; AC3 empty ObjectCommits |
| MP6 | re-drill of C2's MP6 | AC2 `policy EscalateBounded` |
| MP7 | re-drill of C2's MP7 | AC1/AC2 `packetsIn` decode `packetHash does not recompute` |
| MP8 (PD8) | immediately **before** `p := planEffectInvocation(…)` insert `if packet != nil { records = append(records, packet.Hash) }` | AC3b `Replay: … 2 recorded effect records for 1 planned effects` |
| MB1 (PD23) | `return grants[i], true` → `return grants[len(grants)-1], true` | `TestPacketBrokerExports` (1900 ≠ 1760) and AC2 |

---

# EXECUTOR RUN B — C4.1, C5, C6, C7, then the run-B boundary full gate and the C7 re-drill (~45 min; PD31 rollover)

Start of run B: from the worktree root, `shasum -a 256 -c ~/.ailang/state/iter256-plan/runA.sha256` must print only `OK` lines (rc=0). If not, STOP and report: run A's tree changed.

## C4 — Daemon reach: REST reservation (C4.1, LANDS) + `Human.Approve` binding and serialization (C4.2, **HELD — `D-WORLD-73`, do NOT implement**) (design M2, daemon half)

**Budget (C4.1 only):** prod ~5; tests ~110 (AC7 + the shared helpers C5–C7 reuse). **~0.25 d.**

Files (this sprint):
- `host/daemon/handlers.go`
- `host/daemon/decisionpacket_daemon_test.go` (new: AC7 + the shared helpers)
- **NOT** `host/daemon/daemon.go` (C4.2 HELD; the fast gate checks `git diff --quiet -- host/daemon/daemon.go`)

### C4.1 REST reservation (LANDS)

In `decodeCommit`'s object loop (`handlers.go:838`), immediately after `for i, object := range r.Objects {`, insert:
```go
		if object.SemanticID == coordinator.DecisionPacketV1 { // row 172 D4: REST may not forge the packet label (MD1 target)
			return store.Commit{}, fmt.Errorf("objects[%d].semanticId: %s is minted only by the invocation coordinator", i, coordinator.DecisionPacketV1)
		}
```
Add the import `"github.com/sunholo-data/ailang-world/host/coordinator"` to handlers.go.

### C4.2 binding — **HELD (`D-WORLD-73`). The executor SKIPS this whole subsection.** Kept verbatim only as the ready-made follow-up step if the decision is answered "bind". Nothing below in C4.2 (steps 1–4, `TestDaemonInvocationMintsDecisionPacket`, `TestHumanApproveIsSerializedAcrossEpisodes`, MD2, MD3) is written, run or reported this sprint.

1. In `Daemon` (daemon.go:376+), after the `reads readStore` field block, add:
   ```go
   	// humanMu serializes Human.Approve across episodes (row 172 D4): the
   	// broker's approvals-head append reads then sets without compare-and-set.
   	humanMu sync.Mutex
   	// human is the Human.Approve handler d.binder binds: serializedHandler
   	// over broker.NewHumanHandler(store). A field so AC8 can wrap a probe.
   	human broker.Handler
   ```
2. In `New`, on the line after the `d := &Daemon{ … }` literal closes (daemon.go:672), add `d.human = serializedHandler{mu: &d.humanMu, h: broker.NewHumanHandler(s)}`. This is unconditional (PD11).
3. Replace the `binder` body (`daemon.go:781`):
   ```go
   	reg := d.workspace.registry(episodeID) // a fresh map per call
   	reg[broker.EffectHumanApprove] = d.human // row 172 D4 (MD2 target)
   	return broker.OpenBinder(d.store, episodeID, grants, reg)
   ```
   Append to the doc comment above it: `Row 172 adds Human.Approve, serialized daemon-wide, so an invocation can request approval and mint its decision packet.`
4. After the `binder` function, add:
   ```go
   // serializedHandler runs one Execute at a time across every episode (row 172 D4).
   type serializedHandler struct {
   	mu *sync.Mutex
   	h  broker.Handler
   }

   func (s serializedHandler) Execute(ctx context.Context, req broker.EffectRequest, payload []byte) ([]byte, error) {
   	s.mu.Lock()         // MD3 target
   	defer s.mu.Unlock() // MD3 target
   	return s.h.Execute(ctx, req, payload)
   }
   ```

### Test helpers (`decisionpacket_daemon_test.go`, package daemon; reused by C5–C7)

```go
const packetNow, packetExpiry = int64(1700), int64(1760) // the design's Examples
func humanPlan(id string) string // `{"plan":"world/effect-plan/v1","effects":[{"id":"`+id+`","effect":"Human.Approve","scope":"worktree","cost":1,"payload":{"requester":"packet-test"}}],"finish":false,"result":null}`
type humanPlanRunner struct{ plan string; plans atomic.Int64 } // RunContext: like execPlanRunner (workspace_family_test.go:81) but returns r.plan
func publishHumanTool(t *testing.T, d *Daemon) // publishTools(t, d, hashref.SumSHA256([]byte("fake-interpreter")), "module host/capsule/main\n\nexport func main(input: string) -> string {\n  input\n}\n", wsTool{ID: "hu.approve", Effect: broker.EffectHumanApprove, Title: "Approve", Description: "decision packet fixture"})
func humanGrants() []broker.Capability // {{Effect: broker.EffectHumanApprove, Scope: broker.WorkspaceScope, ExpiresAt: packetExpiry, Budget: 5}}
func humanCall(t *testing.T, d *Daemon, episode, task string) coordinator.Call // transitionreg.NewRequest(boundedTestContext(t), transitionreg.NewReader(d.store), wsCaps(humanGrants()), packetNow); SkillID "hu.approve", Input map[string]any{}
func testHumanBinder(d *Daemon) coordinator.BinderFor // explicit registry {Human.Approve: broker.NewHumanHandler(d.store)} — PD25
func humanCoordinator(t *testing.T, d *Daemon, binder coordinator.BinderFor, plan string) *coordinator.Coordinator // Now: func() int64 { return packetNow }, MaxInput/MaxOutput 1<<20
func packetForEntry(t *testing.T, d *Daemon, entry int64) store.Object // scan ObjectsBySemanticID(DecisionPacketV1,"",100); the one whose ObjectCommits(-1,1)[0]==entry; full GetObject; t.Fatal unless exactly one
func packetCount(t *testing.T, d *Daemon) int
```

### Tests first

- `TestRESTCommitRefusesDecisionPacketObjects` (AC7):
  - `d := newHandlerDaemon(t)`; `g := seedGenesisEmbedded(t, d, "ac7")`; `c := testCommit(g, 1, "ac7")`; set `c.Objects[0].SemanticID = coordinator.DecisionPacketV1`.
  - POST with `requestRecorderAuth(t, d, authHeader(t, d), http.MethodPost, "/v1/commit", bytes.NewReader(encodeCommit(c)))`.
  - Expect `assertErrorClass(t, rec, 400, "BadRequest").Error.Message == "objects[0].semanticId: world/decision-packet/v1 is minted only by the invocation coordinator"`, `packetCount == 0`, and the selected head still equal to `g.Ref`.
  - Control: `testCommit(g, 1, "ac7-control")` (semantic `test/ac7-control`) → 200.
- **HELD with C4.2 — do NOT write the next two tests this sprint.**
- `TestDaemonInvocationMintsDecisionPacket` (AC3, daemon arm, **C4.2 HELD**):
  - `newHandlerDaemon`; `publishHumanTool`; `c := humanCoordinator(t, d, d.binder, humanPlan("e1"))`; `res, err := c.Dispatch(boundedTestContext(t), humanCall(t, d, "ep1", "t1"))` → err nil.
  - `packetCount == 1`; `k := packetForEntry(t, d, res.EntryIndex)`; `d.store.ObjectCommits(ctx, res.RecordRef, -1, 1)[0] == res.EntryIndex`.
  - `DecodeDecisionPacket(k.Payload)` → `CreatedAt 1700`, `DeadlineAt 1760`.
  - `broker.RecentApprovals(ctx, d.store, 40)` summaries contain `RequestRef` with status `pending`.
  - Wiring: `sh, ok := d.human.(serializedHandler); ok && sh.mu == &d.humanMu`.
  - Two-approvals: a second coordinator with `humanPlan("e2")`, task `t2` → 2 packets, each first-carried by its own entry.
- `TestHumanApproveIsSerializedAcrossEpisodes` (AC8, **C4.2 HELD**):
  - Probe:
    ```go
    type concurrencyProbe struct{ mu sync.Mutex; cur, max int; second chan struct{}; once sync.Once }
    // Execute: lock; cur++; max = max(max,cur); if cur == 2 { once.Do(close(second)) }; unlock;
    //          select { case <-second: case <-time.After(200*time.Millisecond): case <-ctx.Done(): };
    //          lock; cur--; unlock; return []byte(`{"status":"pending","requestRef":"` + hashref.SumSHA256([]byte("ac8")).String() + `"}`), nil
    ```
  - `sh := d.human.(serializedHandler)`; `d.human = serializedHandler{mu: sh.mu, h: probe}`.
  - Two goroutines dispatch through `humanCoordinator(t, d, d.binder, humanPlan("e1"))` on episodes `epA`/`epB` (tasks `a`/`b`), joined by a `sync.WaitGroup`.
  - Both errs nil; `probe.max == 1`.
  - Control in the same test: a FRESH `probe2 := &concurrencyProbe{second: make(chan struct{})}`, set `d.human = probe2` (unwrapped), and dispatch episodes `epC`/`epD` (tasks `c`/`d`) the same way → `probe2.max == 2`. This proves the probe can observe concurrency, so the `max == 1` arm is not vacuous.

### Gates

Fast gate with `NEW='TestRESTCommitRefusesDecisionPacketObjects'`, `PKGS='./host/daemon'`, `SUBSET='Commit|Handler|REST'`. Also, as a separate call: `git diff --quiet -- host/daemon/daemon.go && echo daemon.go untouched` (C4.2 HELD). Red-first evidence: run `TestRESTCommitRefusesDecisionPacketObjects` BEFORE inserting the C4.1 lines (the helpers and test compile against C1–C3 code, since `coordinator.DecisionPacketV1` exists from C2) and record its red: `status = 200, want 400`. Then implement and record green.

### Mutations (C4.1)

Command: `go test -count=1 -v ./host/daemon -run '^TestRESTCommitRefusesDecisionPacketObjects$'`.

| ID | Edit | Expected red |
|---|---|---|
| MD1 | `if object.SemanticID == coordinator.DecisionPacketV1 {` → `if false && object.SemanticID == coordinator.DecisionPacketV1 {` | AC7 `status = 200, want 400` |
| MD2 | **HELD (C4.2)** — not drilled this sprint | — |
| MD3 | **HELD (C4.2)** — not drilled this sprint | — |

---

## C5 — Validator, projection, validated page (design M3, validator half; AC4)

**Budget:** prod ~220 (validator ~140, daemon ~50, render ~30); tests ~140. **~0.7 d.**

Files:
- `host/coordinator/decisionpacket_validate.go` (new)
- `host/workbench/render.go`
- `host/daemon/workbench_packet.go` (new)
- `host/daemon/workbench.go`
- `host/daemon/workbench_packet_test.go` (new)

### Validator (`decisionpacket_validate.go`, exact)

```go
// PacketReader is the read surface ValidatePacket needs; the daemon's readStore satisfies it.
type PacketReader interface {
	GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error)
	GetLogEntry(ctx context.Context, index int64) (store.LogEntry, bool, error)
	GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error)
	ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]store.Object, error)
	ObjectCommits(ctx context.Context, ref hashref.HashRef, afterEntry int64, limit int) ([]int64, error)
}

// PacketResult is ValidatePacket's verdict: Code is packetRefusal's (0 = validated).
type PacketResult struct {
	Code            int
	Entry           int64                  // the packet's first carrying entry (Code 0 or ≥ 3)
	Packet          DecisionPacket         // Code != 1
	TransitionFn    hashref.HashRef        // the record's transitionFn (Code 0 or ≥ 3)
	Approval        broker.ApprovalSummary // the chain's summary for RequestRef (Code 0 or ≥ 6)
	AnchorTruncated bool                   // code 5 failed and the 40-head read was truncated
	MalformedRef    string                 // code 5 failed on a malformed approvals chain
	UniqueTruncated bool                   // code 7's scan saw more than packetScanLimit packets
}

const packetScanLimit = 100
```

`ValidatePacket(ctx, r, obj)` body. Each `return` noted "refuse" is `return refuse()`, where `refuse := func() (PacketResult, error) { res.Code = packetRefusal(f[0], f[1], f[2], f[3], f[4], f[5], f[6]); return res, nil }` and `var f [7]bool`. Every store error returns `PacketResult{}, err` and is never a code. Steps in order (PD18 short-circuit):

1. **bytes:** `p, err := DecodeDecisionPacket(obj.Payload)`. If `err != nil || obj.SemanticID != DecisionPacketV1 || obj.InterfaceHash != hashref.SumSHA256([]byte(DecisionPacketV1))` → refuse. Else `f[0] = true; res.Packet = p`.
2. **carried:**
   - `es, err := r.ObjectCommits(ctx, obj.Hash, -1, 1)` (PD1); `len(es) == 0` → refuse; `e := es[0]`.
   - `entry, ok, err := r.GetLogEntry(ctx, e)`; `!ok` → refuse.
   - `recObj, ok, err := r.GetObject(ctx, entry.TransitionRef)`; `!ok || recObj.SemanticID != RecordV2` → refuse.
   - `var rec recordV2`; `json.Unmarshal(recObj.Payload, &rec) != nil || !wellFormedV2(rec)` → refuse.
   - `tf, err := hashref.Parse(rec.TransitionFn)`; error → refuse (PD16).
   - `rcs, err := r.ObjectCommits(ctx, recObj.Hash, -1, 1)`; `if len(rcs) == 0 || rcs[0] != e {` → refuse (**MV3 target line**).
   - `f[1] = true; res.Entry = e; res.TransitionFn = tf`.
3. **proposal** (**MV4**: the mutant replaces this whole step with `f[2] = true`):
   - `rec.Plan != p.ProposalHash.String()` → refuse.
   - `planObj, ok, err := r.GetObject(ctx, p.ProposalHash)`; `!ok || planObj.SemanticID != EffectPlanV1` → refuse.
   - `f[2] = true`.
4. **request:**
   - `len(rec.Effects) != 1` → refuse.
   - `erObj` = GetObject(MustParse-free: `hashref.Parse(rec.Effects[0])`, refuse on error); `!ok || erObj.SemanticID != broker.EffectRecordV1` → refuse.
   - `er, err := broker.DecodeRecord(erObj.Payload)` → refuse on error. `er.Effect != broker.EffectHumanApprove || !er.Allowed || er.Failed || er.ResultRef.IsZero()` → refuse.
   - `resObj` = GetObject(er.ResultRef); `!ok || resObj.SemanticID != broker.EffectResultV1` → refuse.
   - `pending, err := broker.ParsePendingApproval(resObj.Payload)` → refuse on error. `if pending != p.RequestRef {` → refuse (**MV5 target line**).
   - `reqObj` = GetObject(p.RequestRef); `!ok || reqObj.SemanticID != broker.ApprovalRequestV1` → refuse.
   - `req, err := broker.DecodeApprovalRequest(reqObj.Payload)` → refuse on error. `req.Effect != broker.EffectHumanApprove || req.Scope != er.Scope` → refuse.
   - `f[3] = true`.
5. **anchored:**
   - `page, err := broker.RecentApprovals(ctx, r, 40)`. If `errors.Is(err, broker.ErrApprovalChainMalformed)` { `res.MalformedRef = page.MalformedRef`; refuse }. `err != nil` → return err.
   - `found := false; for _, s := range page.Summaries { if s.RequestRef == p.RequestRef.String() { found, res.Approval = true, s; break } }`.
   - `if !found { res.AnchorTruncated = page.Truncated; return refuse() }`. **MV6** edits `if !found {` → `if false && !found {`. **MV9** edits `res.AnchorTruncated = page.Truncated` → `res.AnchorTruncated = false`.
   - `f[4] = true`.
6. **lifecycle:** `!lifecycleOk(p.CreatedAt, req.Now, p.DeadlineAt, p.EscalationsRemaining, policyCode(p.Policy))` → refuse. Else `f[5] = true`.
7. **pair-unique:**
   - `others, err := r.ObjectsBySemanticID(ctx, DecisionPacketV1, "", packetScanLimit+1)`; `res.UniqueTruncated = len(others) > packetScanLimit`.
   - `for _, o := range others[:min(len(others), packetScanLimit)] { if o.Hash == obj.Hash { continue }; full, ok, err := r.GetObject(ctx, o.Hash); (err → return); if !ok { continue }; q, err := DecodeDecisionPacket(full.Payload); if err != nil { continue }; if samePair(refOf(p.RequestRef), refOf(p.ProposalHash), refOf(q.RequestRef), refOf(q.ProposalHash)) { return refuse() } }`. **MV8** deletes the `if samePair(…) { return refuse() }` statement.
   - `f[6] = true; return refuse()`, which gives code 0.

Imports: `context`, `encoding/json`, `errors`, `host/broker`, `host/hashref`, `host/store`. No `context.Background`.

### render.go (exact)

1. Add to `ObjectView` after `References *ReferenceView`: `Packet *PacketView // row 172: set only on world/decision-packet/v1 pages`.
2. Add a type after `ReferenceView`:
   ```go
   // PacketView is the decision-packet panel (row 172 D6). A refused packet
   // carries only Sentence; links, status and evidence exist only when Validated.
   type PacketView struct {
   	Validated                 bool
   	Sentence                  string
   	EntryIndex                int64
   	EntryHref                 string
   	ProposalRef, ProposalHref string
   	RequestRef, RequestHref   string
   	CreatedAt, DeadlineAt     int64
   	Policy                    string
   	Escalations               int64
   	ApprovalStatus            string
   	TransitionFnRef           string
   	ProposalGrade             GradeView
   	TransitionGrade           GradeView
   }
   ```
3. In `pageHTML`, change the line that begins `{{if .PayloadShown}}` so it begins `{{with .Packet}}{{template "packet" .}}{{end}}{{if .PayloadShown}}` (PD20). Leave the rest of the line unchanged and add no new line.
4. Append to `partialsHTML` (before its closing backquote, each define on its own line):
   ```
   {{define "packet"}}<section aria-label="decision packet"><h3>decision packet</h3><p>{{.Sentence}}</p>{{if .Validated}}<p>createdAt {{.CreatedAt}} · deadlineAt {{.DeadlineAt}} — logical times; this host emits no Timeout transition (item 7 owns the emitter), so this page does not state whether the deadline has passed</p><p>policy {{.Policy}} · escalationsRemaining {{.Escalations}}</p><p>proposal: <a href="{{workbenchHref .ProposalHref}}" class="hash">{{.ProposalRef}}</a></p><p>approval request: <a href="{{workbenchHref .RequestHref}}" class="hash">{{.RequestRef}}</a></p><p>committed at: <a href="{{workbenchHref .EntryHref}}">entry {{.EntryIndex}}</a></p><p>approvals chain: {{.ApprovalStatus}}</p><p>evidence: proposal <span class="hash">{{.ProposalRef}}</span> — {{template "packetGrade" .ProposalGrade}}</p><p>evidence: transition function <span class="hash">{{.TransitionFnRef}}</span> — {{template "packetGrade" .TransitionGrade}}</p>{{end}}</section>{{end}}
   {{define "packetGrade"}}{{if .Available}}<span>{{.Label}}</span>{{if .HasVerdict}} verdict: {{.Verdict}}{{end}}{{else}}GRADE UNAVAILABLE — {{with .Unavailable}}{{.}}{{else}}no grade reason was supplied{{end}}{{end}}{{end}}
   ```
   **MV11** edits `{{if .Validated}}` → `{{if true}}`.

### workbench_packet.go (new, package daemon, exact constants)

```go
const (
	packetValidatedFormat              = "decision packet validated: committed with its invocation record at entry #%d; proposal, approval request and approvals-chain anchor agree; lifecycle fields are the v1 producer's"
	packetUniqueTruncation             = " (duplicate check bounded to the first 100 stored packets)"
	packetRefusedBytes                 = "not a validated decision packet: the bytes are not a canonical world/decision-packet/v1 encoding whose packetHash recomputes"
	packetRefusedCarried               = "not a validated decision packet: it was not committed together with its invocation record"
	packetRefusedProposal              = "not a validated decision packet: proposalHash is not the plan of the invocation record committed with it"
	packetRefusedRequest               = "not a validated decision packet: requestRef is not the approval request that invocation's Human.Approve effect returned"
	packetRefusedAnchor                = "not a validated decision packet: its approval request is not in the broker's approvals chain"
	packetRefusedAnchorBounded         = "not a validated decision packet: its approval request is not among the newest 40 approval heads, so its approvals-chain membership is not established"
	packetRefusedAnchorMalformedPrefix = "not a validated decision packet: the approvals chain is malformed at " // PD13
	packetRefusedLifecycle             = "not a validated decision packet: its lifecycle fields disagree with its request's logical time or are not the v1 producer's"
	packetRefusedDuplicate             = "not a validated decision packet: another stored packet names the same approval request and the same proposal"
)

func packetRefusalSentence(res coordinator.PacketResult) string // switch res.Code 1..7; case 5: MalformedRef != "" → prefix+ref; AnchorTruncated → bounded; else anchor

func (d *Daemon) packetView(ctx context.Context, obj store.Object) (*workbench.PacketView, error) {
	res, err := coordinator.ValidatePacket(ctx, d.reads, obj)
	if err != nil {
		return nil, err
	}
	if res.Code != 0 {
		return &workbench.PacketView{Sentence: packetRefusalSentence(res)}, nil
	}
	sentence := fmt.Sprintf(packetValidatedFormat, res.Entry)
	if res.UniqueTruncated {
		sentence += packetUniqueTruncation
	}
	proposalGrade, err := d.objectGradeView(ctx, res.Packet.ProposalHash)
	if err != nil {
		return nil, err
	}
	tfGrade, err := d.objectGradeView(ctx, res.TransitionFn) // MV10 target
	if err != nil {
		return nil, err
	}
	status := res.Approval.Status
	if res.Approval.DecidedBy != "" {
		status += " by " + res.Approval.DecidedBy
	}
	return &workbench.PacketView{Validated: true, Sentence: sentence, EntryIndex: res.Entry, EntryHref: pageHref(res.Entry, res.Entry),
		ProposalRef: res.Packet.ProposalHash.String(), ProposalHref: "?object=" + res.Packet.ProposalHash.String(),
		RequestRef: res.Packet.RequestRef.String(), RequestHref: "?object=" + res.Packet.RequestRef.String(),
		CreatedAt: res.Packet.CreatedAt, DeadlineAt: res.Packet.DeadlineAt, Policy: res.Packet.Policy,
		Escalations: res.Packet.EscalationsRemaining, ApprovalStatus: status, TransitionFnRef: res.TransitionFn.String(),
		ProposalGrade: proposalGrade, TransitionGrade: tfGrade}, nil
}
```

### workbench.go (exact)

1. Add the import `host/coordinator`.
2. Between the `grade, err := d.objectGradeView(ctx, ref)` error block and `page.Object = &workbench.ObjectView{`, insert:
   ```go
   		var packet *workbench.PacketView
   		if object.SemanticID == coordinator.DecisionPacketV1 { // row 172 D6 (MV12 target)
   			var packetErr error
   			if packet, packetErr = d.packetView(ctx, object); packetErr != nil {
   				d.writeWorkbenchStoreError(w, r, ctx, packetErr)
   				return
   			}
   		}
   ```
3. Add `Packet: packet,` after `References: references,` in the literal.

### Test helpers (`workbench_packet_test.go`)

```go
const packetSectionStart = `<section aria-label="decision packet">`
func packetSection(t *testing.T, body string) string   // workbenchRegion(body, packetSectionStart, "</section>"); t.Fatal if absent
func wantPacketGradeLine(kind, ref string, g workbench.GradeView) string
	// "<p>evidence: " + kind + ` <span class="hash">` + ref + "</span> — " + (g.Available ? "<span>"+string(g.Label)+"</span>" (+" verdict: "+string(g.Verdict) if HasVerdict) : "GRADE UNAVAILABLE — "+html.EscapeString(reason-or-"no grade reason was supplied")) + "</p>"
func producePacket(t *testing.T, d *Daemon, plan, task string) (coordinator.Result, store.Object) // publishHumanTool once per daemon (caller), humanCoordinator(t, d, testHumanBinder(d), plan) → Dispatch(humanCall(t,d,"ep1",task)) → packetForEntry
func recordTransitionFn(t *testing.T, d *Daemon, recordRef hashref.HashRef) hashref.HashRef // GetObject → json.Unmarshal into map[string]any → hashref.MustParse(m["transitionFn"].(string))
type countingPacketScanStore struct { readStore; calls *atomic.Int32 } // counts ObjectsBySemanticID(id == coordinator.DecisionPacketV1)
```

### Tests first

`TestWorkbenchPacketValidated` (AC4):
1. `d := newHandlerDaemon(t)`; `publishHumanTool(t, d)`; `res, k := producePacket(t, d, humanPlan("e1"), "t1")`; `p, _ := coordinator.DecodeDecisionPacket(k.Payload)`.
2. `code, body := gradePage(t, d, k.Hash)` → 200. In `sec := packetSection(t, body)`, require:
   - `strings.Count(sec, "<p>"+html.EscapeString(fmt.Sprintf(packetValidatedFormat, res.EntryIndex))+"</p>") == 1`
   - `href="/workbench?object=` + p.ProposalHash.String() + `"`
   - `href="/workbench?object=` + p.RequestRef.String() + `"`
   - `fmt.Sprintf(\`href="/workbench?from=%d&amp;entry=%d"\`, res.EntryIndex, res.EntryIndex)`
   - `approvals chain: pending`
   - `policy Cancel · escalationsRemaining 0`
   - `createdAt 1700 · deadlineAt 1760 — logical times; this host emits no Timeout transition (item 7 owns the emitter), so this page does not state whether the deadline has passed`
   - `wantPacketGradeLine("proposal", P, pg)` and `wantPacketGradeLine("transition function", TF, tg)`, where `pg, _ := d.objectGradeView(ctx, p.ProposalHash)` and `tg, _ := d.objectGradeView(ctx, recordTransitionFn(t, d, res.RecordRef))`.
3. Non-packet arm (PD6): `var n atomic.Int32; d.reads = countingPacketScanStore{readStore: d.reads, calls: &n}`. Request the page of the plan object `p.ProposalHash` → 200. The body contains none of `packetSectionStart`, `decision packet validated`, `not a validated decision packet`, and `n.Load() == 0`.
4. Control: the packet page with the same wrapper → `n.Load() == 1`.

### Gates

Fast gate with `NEW='TestWorkbenchPacketValidated'`, `PKGS='./host/daemon ./host/workbench'`, `SUBSET='Workbench|Packet|Object'`. Also: `go test -count=1 ./host/workbench -v` must be all PASS (render tests and goldens unchanged), and `go test -count=1 ./host/daemon -run '^TestWorkbenchObjectGrade$' -v` → `--- PASS` with `git diff --quiet -- host/daemon/workbench_test.go`. Red-first evidence: with `workbench_packet.go`, the render type and the workbench.go branch absent, `TestWorkbenchPacketValidated` is build-red (`undefined: packetValidatedFormat`) — record it, but it is not a kill; the C5 mutation ledger is the red-first proof. Record the wall time of the `./host/daemon` subset run.

### Mutations (C5)

All use `go test -count=1 -v ./host/daemon -run '^TestWorkbenchPacketValidated$'`.

| ID | Edit | Expected red |
|---|---|---|
| MV12 (PD6) | `if object.SemanticID == coordinator.DecisionPacketV1 {` → `if true {` | the non-packet arm finds `<section aria-label="decision packet">` (a code-1 refusal) |
| MP7 (page) | re-drill of C2 MP7 | the produced packet's page shows the code-1 sentence; the validated count is 0 |
| MV3-genuine | `if len(rcs) == 0 || rcs[0] != e {` → `if len(rcs) == 0 || rcs[0] == e {` | the genuine packet renders code 2. This mutant is an extra non-vacuity check on the AC4 validated arm. |

---

## C6 — Forged and mismatched packets; evidence only through row 105 (design M3, arms half; AC5, AC6)

**Budget:** prod 0; tests ~300. **~0.6 d.**

Files: `host/daemon/workbench_packet_test.go` (append only).

### Fixture (PD14, exact)

```go
func forgeObject(semanticID string, payload []byte) store.Object // Hash SumSHA256(payload), InterfaceHash SumSHA256([]byte(semanticID)), Provenance "packet-forge"
type forgedRecordV2 struct { /* json tags exactly: invocationId episodeId skillId transitionFn interpreter semanticsEpoch input output plan effects ([]string) — the plan.go:39-50 order */ }
func approvalRequestInChain(t *testing.T, d *Daemon, scope, requester string, now int64) hashref.HashRef
	// out, err := broker.NewHumanHandler(d.store).Execute(boundedTestContext(t), broker.EffectRequest{Effect: broker.EffectHumanApprove, Scope: scope, Cost: 1, Now: now}, []byte(`{"requester":"` + requester + `"}`)); ref, err := broker.ParsePendingApproval(out)
	// request identity is {effect, scope, cost, requester, now} (V10): vary requester to get a DISTINCT request at the same now
type forgedChain struct {
	label       string
	request     hashref.HashRef            // the request the effect result's pending wire names
	scope       string                     // effect-record scope (must equal the request's scope to pass code 4)
	tf          hashref.HashRef            // record transitionFn
	planRef     hashref.HashRef            // record plan (default: the plan object's hash)
	includePlan bool                       // default true
	packet      coordinator.DecisionPacket // default {ProposalHash: planRef, RequestRef: request, CreatedAt: packetNow, DeadlineAt: packetNow+60, Policy: "Cancel"}
	rawPacket   []byte                     // when non-nil, replaces EncodeDecisionPacket(packet)
	extra       []store.Object             // more objects in the same commit
}
func newForgedChain(label string, request hashref.HashRef, scope string) *forgedChain // tf = SumSHA256("forge-tf-"+label); plan = forgeObject(coordinator.EffectPlanV1, `{"forged":"<label>"}`)
func (f *forgedChain) planObject() store.Object
func commitObjects(t *testing.T, d *Daemon, label string, objects []store.Object, transitionRef hashref.HashRef) int64
	// head: d.store.SelectedHead; if !ok → w := seedGenesisEmbedded(t, d, label+"-genesis") else w from d.store.GetWorld(head);
	// c := testCommit(w, w.Revision+1, label); c.Objects = objects; c.Entry.TransitionRef = transitionRef; d.store.Commit(boundedTestContext(t), c) (t.Fatal on error); return c.Entry.Header.EntryIndex
func commitForgedChain(t *testing.T, d *Daemon, f *forgedChain) (store.Object, int64)
	// result := forgeObject(broker.EffectResultV1, []byte(`{"status":"pending","requestRef":"`+f.request.String()+`"}`))
	// er := forgeObject(broker.EffectRecordV1, broker.EncodeRecord(broker.EffectRecord{Effect: broker.EffectHumanApprove, Scope: f.scope, Cost: 1, BudgetBefore: 5, BudgetAfter: 4, Allowed: true, RequestRef: hashref.SumSHA256([]byte("forge-effect-request-"+f.label)), ResultRef: result.Hash}))
	// record := forgeObject(coordinator.RecordV2, json.Marshal(forgedRecordV2{InvocationID: "forge-"+label, EpisodeID: "forge", SkillID: "forge", TransitionFn: f.tf.String(), Interpreter: SumSHA256("forge-interp"), SemanticsEpoch: 1, Input: SumSHA256("forge-in-"+label), Output: SumSHA256("forge-out-"+label), Plan: f.planRef.String(), Effects: []string{er.Hash.String()}}))
	// packet := forgeObject(coordinator.DecisionPacketV1, rawPacket or EncodeDecisionPacket(f.packet))
	// objects := [result, er, record, packet] (+ plan if includePlan) + extra; entry := commitObjects(t, d, f.label, objects, record.Hash)
func wantRefused(t *testing.T, d *Daemon, ref hashref.HashRef, sentence string)
	// 200; sec := packetSection(t, body); strings.Contains(sec, "<p>"+html.EscapeString(sentence)+"</p>"); and none of "approvals chain:", "evidence: proposal", "committed at:" anywhere in body
func wantValidated(t *testing.T, d *Daemon, ref hashref.HashRef, entry int64)
```

### Tests first

`TestWorkbenchPacketRefusals` (AC5). One `t.Run` per arm, each on a FRESH `newHandlerDaemon(t)`. In every forged arm, `scope := "forge"` and `q := approvalRequestInChain(t, d, scope, "forge", packetNow)`, unless the arm says otherwise (arms 2, 5 and 7b create no `q`).

| Arm | Construction | Expectation |
|---|---|---|
| `0-forged-control` | `newForgedChain("c0", q, scope)` | `wantValidated` (the declared residual; positive control) |
| `1-whitespace` | `f.rawPacket = []byte(strings.Replace(string(coordinator.EncodeDecisionPacket(f.packet)), ",\"createdAt\"", ", \"createdAt\"", 1))` | `packetRefusedBytes` |
| `1b-wrong-hash` | rawPacket = the encoding with its packetHash value replaced by `hashref.SumSHA256([]byte("wrong")).String()` | `packetRefusedBytes` |
| `2-recarry` | `publishHumanTool`; `res, k := producePacket(…"e1","t1")`; `p` = decoded k with `DeadlineAt+1`; `k2 := forgeObject(DecisionPacketV1, EncodeDecisionPacket(p))`; `commitObjects(t, d, "recarry", []store.Object{k2}, res.RecordRef)` | page of k2: `packetRefusedCarried` |
| `3-other-plan` | `p2 := forgeObject(coordinator.EffectPlanV1, []byte(`{"forged":"other"}`))`; `f.extra = {p2}`; `f.packet.ProposalHash = p2.Hash` | `packetRefusedProposal` |
| `3b-missing-plan` | `m := hashref.SumSHA256([]byte("missing-plan"))`; `f.includePlan = false`; `f.planRef = m`; `f.packet.ProposalHash = m` | `packetRefusedProposal` |
| `4-other-request` | `q2 := approvalRequestInChain(t, d, scope, "forge-other", packetNow)` (same scope and `now`, distinct requester); `f.packet.RequestRef = q2` | `packetRefusedRequest` |
| `5-unanchored` | NO `q`. `approvalRequestInChain(t, d, "unrelated", "forge", packetNow)` (so the chain exists); `q3 := forgeObject(broker.ApprovalRequestV1, []byte(`{"effect":"Human.Approve","scope":"forge","cost":1,"requester":"forge-unanchored","now":1700}`))` stored by `d.store.PutObject` (never appended); `newForgedChain("c5", q3.Hash, scope)` | `packetRefusedAnchor` |
| `5b-bounded` | chain over q; then 41 × `approvalRequestInChain(t, d, fmt.Sprintf("newer-%d", i), "forge", packetNow)` | `packetRefusedAnchorBounded` |
| `6-created-at` | `f.packet.CreatedAt = packetNow + 1` | `packetRefusedLifecycle` |
| `6b-policy` | `f.packet.Policy = "EscalateBounded"` | `packetRefusedLifecycle` |
| `7-duplicate-pair` | commit the control chain; `dup := f.packet; dup.DeadlineAt++`; `d.store.PutObject(ctx, forgeObject(DecisionPacketV1, EncodeDecisionPacket(dup)))` | control packet's page: `packetRefusedDuplicate` |
| `7b-collapsed-request` | `publishHumanTool`; produce `e1/t1` and `e2/t2` (PD27) | both pages `wantValidated` at their own entries; `RequestRef` equal, `ProposalHash` different |

`TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive` (AC6):
- **configured:** `d, rig := newGradeDaemon(t, []string{gradeTestIdentity})`; `tf := storeGradeSubject(t, d, "packet-tf")`; `storeGradeEnvelope(t, d, mintGradeEnvelope(t, rig.key, gradeReport(rig, tf)))`; `f := newForgedChain("ac6", approvalRequestInChain(t, d, "forge", "forge", packetNow), "forge")`; `f.tf = tf`; commit. Then on the packet page:
  - the section contains `wantPacketGradeLine("transition function", tf.String(), tg)`, and `tg.Available && tg.Label == workbench.GradePROVEN`, i.e. the literal `— <span>PROVEN</span></p>`;
  - the proposal line equals `wantPacketGradeLine("proposal", P, pg)` with `pg.Unavailable == gradeReasonBoundElsewhere`.
- **unconfigured:** the same chain on `newHandlerDaemon` → both lines `GRADE UNAVAILABLE — ` + `html.EscapeString(objectGradeUnavailableReason)`.

### Gates

Fast gate with `NEW='TestWorkbenchPacketValidated|TestWorkbenchPacketRefusals|TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive'`, `PKGS='./host/daemon'`, `SUBSET='Workbench|Packet|Grade'`. Record `go test -count=1 -v ./host/daemon -run '^(TestWorkbenchPacketRefusals|TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive)$'` wall time; the new daemon tests together must add < 15 s (the `-race` figure is measured once, in the run-B boundary gate F3). Red-first evidence: arm `0-forged-control` is the positive control and must render validated on the unmutated tree; every refusal arm is shown red-capable by its own mutant below.

### Mutations (C6)

All use `go test -count=1 -v ./host/daemon -run '^(TestWorkbenchPacketRefusals|TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive)$'`.

| ID | Edit | Killer → expected red |
|---|---|---|
| MV1 | C2's MV1 line | arm 1 renders the validated sentence |
| MV2 | C2's MV2 line | arm 1b renders validated |
| MV3 | `if len(rcs) == 0 || rcs[0] != e {` → `if len(rcs) == 0 {` | arm 2 renders `packetRefusedDuplicate` (code 7, since k2 shares the genuine pair), not code 2 |
| MV4 | replace step 3 with `f[2] = true` | arms 3 and 3b render validated |
| MV5 | `if pending != p.RequestRef {` → `if false && pending != p.RequestRef {` | arm 4 renders **validated** (q2 is a stored, chain-anchored request with the same scope and `now`, so codes 5–7 pass) |
| MV6 | `if !found {` → `if false && !found {` | arm 5 renders validated |
| MV7 | C1's MV7 Go edit | arm 6 renders validated |
| MV8 | delete the `if samePair(…) { return refuse() }` statement | arm 7 renders validated |
| MP8b | C1's MP8b Go edit | arm 7b renders `packetRefusedDuplicate` on both pages |
| MV9 | `res.AnchorTruncated = page.Truncated` → `res.AnchorTruncated = false` | arm 5b renders `packetRefusedAnchor` |
| MV10 | `d.objectGradeView(ctx, res.TransitionFn)` → `d.objectGradeView(ctx, obj.Hash)` | AC6 configured: no `<span>PROVEN</span>` on the transition line |
| MV11 | render.go partial `{{if .Validated}}` → `{{if true}}` | every refused arm: `approvals chain:` present |

---

## C7 — S7 usage, executed doc test, full re-drill, final gates (design M4)

**Budget:** doc ~45 lines; test ~70. **~0.4 d.**

Files:
- `docs/QUICKSTART.md` (append at EOF)
- `host/daemon/workbench_packet_test.go` (append)

### QUICKSTART text (append verbatim after the last line; PD15, PD32 — this is the D4-HELD text; do not compose your own)

````markdown

**Inspect a decision packet** (row 172). The invocation coordinator mints a `world/decision-packet/v1` object only when a skill's plan requests exactly one `Human.Approve` effect, the binder in use binds `Human.Approve`, and a live grant for that exact effect and scope authorizes it. **This release's daemon does not bind `Human.Approve`** (decision `D-WORLD-73` is held), so a daemon on this release mints no packets and the listing below returns `{"items":[]}`. REST commits cannot store an object with this semantic ID either (`400 BadRequest`). A request never decides anything: decisions stay operator-only. When a packet exists, it rides in its invocation's own commit: `createdAt` is the invocation's logical time, `deadlineAt` is the authorizing grant's expiry, `policy` is `Cancel` and `escalationsRemaining` is `0`.

List packets and open the first one's object page while the daemon serves on 7644:

```bash
curl --silent --fail http://127.0.0.1:7644/v1/objects/by-semantic-id/world/decision-packet/v1
PACKET=$(curl --silent --fail http://127.0.0.1:7644/v1/objects/by-semantic-id/world/decision-packet/v1 | python3 -c 'import json,sys; items=json.load(sys.stdin)["items"]; print(items[0]["hash"] if items else "")')
[ -n "$PACKET" ] && curl --silent --fail "http://127.0.0.1:7644/workbench?object=$PACKET" | grep -o 'decision packet validated[^<]*\|not a validated decision packet[^<]*'
```

A produced packet reads `decision packet validated: committed with its invocation record at entry #N; …`. The page links the proposal (the invocation's `world/effect-plan/v1` plan), the approval request and entry N. It shows `approvals chain: pending` until an operator decides, and two evidence lines graded only by the object-grade validator above. Any other object labelled `world/decision-packet/v1` reads `not a validated decision packet: …` with the first failing check and no links. The checks are bytes, carrying commit, proposal, request, approvals-chain anchor, lifecycle fields and duplicate pair. The anchor check reads the newest 40 approval heads and keeps at most the newest 20 distinct requests; beyond that it reports "not established" rather than a verdict. The duplicate check reads the first 100 stored packets. The times are logical; this host emits no Timeout transition, so the page does not say whether the deadline has passed.
````

(For the HELD follow-up only, not this sprint: if `D-WORLD-73` is answered "bind" and C4.2 lands, replace the sentence beginning "**This release's daemon does not bind" with "The daemon binds `Human.Approve` for every episode and serializes it; it never decides a request.", add the grant example "for example `--grant Human.Approve=demo:1` on `session mint` (§0)", and update the static test below.)

### Test first

`TestQuickstartDecisionPacketSection` (S7; PD15, PD32, PD33), in `host/daemon/workbench_packet_test.go`:
1. `guide, err := os.ReadFile("../../docs/QUICKSTART.md")`; `_, section, ok := strings.Cut(string(guide), "**Inspect a decision packet** (row 172).")` → ok.
2. Static requirements on `section`:
   - contains `D-WORLD-73`, `{"items":[]}`, `http://127.0.0.1:7644/v1/objects/by-semantic-id/world/decision-packet/v1`, `/workbench?object=$PACKET`, `decision packet validated`, `not a validated decision packet`, `approvals chain: pending`, `400 BadRequest`;
   - does **NOT** contain `--grant Human.Approve=` (with D4 held a grant cannot produce a packet on the daemon);
   - contains **no** line with prefix `/tmp/ailang-worldd serve ` and **no** line starting with `|`.
3. Execution (the documented commands, run against handler daemons):
   - Extract the route with `regexp.MustCompile(`http://127\.0\.0\.1:7644(/v1/objects/by-semantic-id/[^\s|"')]+)`)` (PD33); require the first submatch to equal `/v1/objects/by-semantic-id/world/decision-packet/v1`.
   - **Held-release arm:** `d0 := newHandlerDaemon(t)`; `publishHumanTool(t, d0)`; `rec := requestRecorder(t, d0, http.MethodGet, route, nil)` → 200 and `strings.TrimSpace(rec.Body.String()) == `{"items":[]}`` (the doc's claim; publishing a `Human.Approve` skill alone mints nothing).
   - **Produced-packet arm:** `d := newHandlerDaemon(t)`; `publishHumanTool(t, d)`; `res, k := producePacket(t, d, humanPlan("e1"), "t1")` (explicit test binder, PD25). `GET route` → 200; decode as `struct{ Items []struct{ Hash string `json:"hash"` } `json:"items"` }`; `Items[0].Hash == k.Hash.String()` (the doc's python path).
   - Require `section` to contain `http://127.0.0.1:7644/workbench?object=$PACKET`; then `page := requestRecorder(t, d, http.MethodGet, strings.Replace("/workbench?object=$PACKET", "$PACKET", Items[0].Hash, 1), nil)` → 200, and `page.Body.String()` contains `html.EscapeString(fmt.Sprintf(packetValidatedFormat, res.EntryIndex))`.
   - Red-first: write the test before appending the doc text and record its red (`section not found`).

### Gates (final milestone, then the run-B boundary)

Fast gate with `NEW='TestQuickstartDecisionPacketSection'`, `PKGS='./host/daemon'`, `SUBSET='Quickstart|Packet'`. Plus, as separate calls:

```bash
go test -count=1 ./cmd/ailang-worldd -run 'Quickstart' -v        # §0/§9/§10/§11 parsers unchanged: every listed test PASSes
go test -count=1 ./host/daemon -run '^(TestWorkbenchLiveQuickstartControl|TestAGUIQuickstartControl|TestExecStartupRefusalTable|TestMCPQuickstartPayloadsVerbatim|TestSeToolsQuickstartPayloadsVerbatim)$' -v
git diff --quiet -- host/daemon/daemon.go && echo "daemon.go untouched (C4.2 HELD)"
git status --porcelain --untracked-files=all    # must list exactly the 21-path "Final changed-file set" (plus the plan files if uncommitted)
```

Then the **run-B boundary full gate**: F1–F4 with `RACE='./host/coordinator ./host/daemon ./host/workbench'`, and `shasum -a 256 -c ~/.ailang/state/iter256-plan/runA.sha256` (all `OK`). Record the F3 `host/daemon` `-race` wall time against the 221 s base (+15 s budget).

**Not run by the executor:** `scripts/verify_go.sh` (controller-run, outside any sandbox; its `-race -timeout 8m` leg exceeds the tool cap).

### Mutation re-drill on the final tree (C7; PD31)

Re-run, one at a time with the mutation protocol, each with its own milestone's command: **MD1; MV1, MV2, MV3, MV3-genuine, MV4, MV5, MV6, MV7, MV8, MV9, MV10, MV11, MV12; MP7 (page); MP8b (arm 7b)** — 16 mutants, ~35 s each. **Not re-drilled:** MD2, MD3 (HELD with C4.2, never applied); MS1–MS4, MS3, MP1–MP6, MP8, MB1, MP8b-drift, MV7-drift (run-A targets and killers, carried by the run-A ledger plus the sha256 manifest check). Rollover per PD31 if run B has used 40 min when the re-drill starts.

---

## Mutation index (design §8 → plan)

| Design | Plan milestone(s) | Note |
|---|---|---|
| MP1–MP4 | C3 | MP2–MP4 edit the `mintPacket` call in `mintApprovalPacket` |
| MP5 | C3 | concrete form PD7 |
| MP6, MP7 | C2 (codec), C3 (AC2), C5 (page, MP7) | MP7 encoder-only (PD21) |
| MP8 | C3 | concrete form PD8 |
| MP8b | C1 (drift), C6 (7b) | |
| MV1, MV2 | C2 (codec), C6 (arms 1, 1b) | |
| MV3–MV9 | C6 (MV3 also C5 as MV3-genuine) | MV3's arm-2 red is the code-7 sentence (shared pair) |
| MV7 | C1 (drift), C6 (arm 6) | |
| MV10 | C6 (AC6) | |
| MV11 | C6 | render.go partial |
| MV12 | C5 | killed by body assertion (PD6) |
| MD1 | C4.1 | lands |
| MD2, MD3 | — | **HELD with C4.2 (`D-WORLD-73`)**: not applied, not drilled this sprint |
| MS1–MS3 | C1 | MS1/MS2 measured by the planner on a scratch copy (cex 1 each) |
| — | MS4 (C1), MB1 (C3) | added: PD17, PD23 |

## Risks

| Risk | Mitigation |
|---|---|
| Haiku misses a verifygate count site (PD3); verifygate reds after 4 minutes | C1 lists all five sites by line number; C1's gate runs the three named verifygate tests explicitly |
| A `publishTools` genesis or entry index collides with forged commits | `commitObjects` always chains from the selected head (`Revision+1`), the coordinator's convention |
| The daemon `-race` budget | Per-arm fresh daemons are cheap (no AilangBin); record wall time at C4/C5/C6; budget +15 s |
| The `·` / `—` characters in the template | They are literal UTF-8 in a Go raw string; `html/template` does not escape them. Tests compare with the same literals |
| D4 (`D-WORLD-73`) | HELD by default: C4.2 is not implemented; PD25 keeps C5–C7 independent of `d.binder`; `daemon.go` is in the frozen list; the QUICKSTART text says the daemon mints no packets (PD29, PD32). If the decision is answered "bind" later, C4.2 is the ready follow-up step |
| D2b freeze not exercised in production (PD29) | Pinned by `TestDecisionPacketCodec`'s golden bytes; no production write exists until D4 binds, so the freeze cannot be violated silently before then |
| Code-5 bounded sentence wording vs the 20-request cap (PD28) | Sentence kept verbatim (design-frozen); finding reported to the controller |
| A run overruns ~45 min | PD30 fast gates; PD31 rollover of the C7 re-drill into a short run B2 |
| D2b's byte form | Frozen by `TestDecisionPacketCodec`'s golden; any change is `/v2` (design D2b) |
| verify_go.sh exceeds the tool cap | The controller runs it (row 105 precedent); the executor never runs it |
