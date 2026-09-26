# w-store-deadline-free-residue-owner — inherit caller cancellation; zero is a syntax result, not a deadline guarantee

- Status: **planned; plumbing prototyped, policy closing move reserved for Mark**.
- Date: 2026-09-26 · Iteration: 194 · Base: `b5b4a7d` · Owner: **OPEN queue row 23**.
- Scope: clause-2, Go host only · Estimate: **0.75–1 day plumbing**, policy tranche separately priced after ratification.
- Depends on: D-WORLD-18 and D-WORLD-23; follows [daemon-read-cancellation](../implemented/w-daemon-read-cancellation.md), especially §2.8/§10/§13.
- Format/reference: [daemon-lock-wait-not-deadline-bound](../implemented/w-daemon-lock-wait-not-deadline-bound.md). This doc owns that predecessor's explicitly deferred context residue, not its lock-wait enforcement.
- Prototype and evidence: `~/.ailang/state/world-iter194/design/`; patch `proto.patch`. Only this document is the controller's deliverable; prototype source is disposable, uncommitted.

## §1 Problem, measured scope, and ruling

The baseline ratchet pins **8 approval + 2 registry + 1 replay = 11** direct literal
`context.Background()` getter arguments. Its regular expression sees a spelling, not the origin
or deadline of a variable. The baseline also leaves `GetVerifyResult` context-free, and replay
is its sole production caller (F1, F6, Vreads). The owning queue row and the two attended rulings
are recorded in Vrules; row 23 is still OPEN, rather than a landed item used as a residual owner.

**D-WORLD-18 remains binding.** The guard is the follow-on's declared closing move, landable
exactly when the ratchet reads zero; the approval-bound and Commit-atomicity questions return
to Mark. **Zero is a necessary checkpoint, not evidence that a strict guard is compatible.**
Our temporary guard experiment breaks four production-root paths even after zero (Vguard).
Do not silently implement that guard in the plumbing tranche.

**D-WORLD-23 application:** keep the tranche, claim only caller-context propagation, pin its
composition conditions, and retain the policy residue under OPEN row 23. This document does
not close row 23's full inherited guard/Commit contract merely by completing the plumbing.
The controller may record the zero-ratchet subtranche separately; it must not mark the entire
follow-on complete or invent a landed residual owner.

Why this is not a package: this proposal changes Go call signatures and SQLite read context
propagation at the existing host boundary. No world semantics, `.ail`, schema, or compiler
change is proposed. S2/S3 and DESIGN §14 favor leaving this plumbing at that boundary (Vrefs).

The root census covers **host/ and cmd/**. A separate surface guard rejects non-test Go
files elsewhere, except the two explicitly named compiler-regression reproducers (§3, Vsurface).

## §2 Decisions Q1–Q3: thread existing lifetimes, without choosing a duration

### Q1 — approvals: all eight reads receive the real invoking context

Baseline `HumanHandler.Execute(_ context.Context, ...)` discards its broker context. Approve
writes a request and returns `pending`; PollApproval traverses the recorded chain and returns
pending or decided. Neither read waits for a human. The CLI checks attendance **before** calling
the mint operation (F2, Vguardpaths). The first D-WORLD-18 policy question therefore survives
for the **guard and active-I/O budget**, not for passing the caller's existing context.

| Baseline sites | Proposed context chain | Measured consequence |
|---|---|---|
| `approve.go:169` | `DecideApproval(ctx)` → `decideApproval(ctx)` | Exported operator API gains ctx; private mint caller also threads it. |
| `:195` | Execute/decide → `appendApprovalHead(ctx)` | Caller cancellation reaches registry-head read. |
| `:225/:234/:251/:261` | Execute poll / decide → find helpers → `walkApprovalHead(ctx)` | Same ctx reaches head, chain object, decision, and request branches. |
| `:495/:522` | `Session.Invoke(ctx)` → `invoke(ctx)` → `validatePublishApproval(ctx)` | Validation uses the invoking context before dispatch. |

`mintAttendedApproval` is a real private production caller of `decideApproval`; it formerly
created **two** background contexts for its sessions. Add ctx to both mint entry points and
thread it through all three legs. The CLI `runApprove` has no context argument; introduce one
**explicit, named compatibility root at its call** in place of the two hidden roots. This is a
relocation/consolidation, not newly bounded work, and the root census pins it (F3, Vcalls, Vroots).
Never describe it as eliminating deadline-free approval operations.

The live publish path is `cmd/world-publish.runPublish` → `InvokeAttendedPublish(ctx)` → session
→ validation. Its root is **Background, without a deadline**. The handler's existing ExecTimeout
is derived later, inside subprocess execution; it cannot cover the earlier approval reads.
The production Invoke call census has the mint request, mint poll, and attended publish calls;
it does not connect these reads to a capsule ExecTimeout or an HTTP request deadline (F3).
Exported `Session.Invoke` remains usable by other callers; their deadline is their responsibility.

**Measured API cost at base:** `DecideApproval` 0 production / **1 actual test call**;
`decideApproval` 2 production / 7 test calls; `MintAttendedApproval` 1 production / 0 test calls;
private mint 1 production / 4 test calls (Vcalls). F4's count of three exported test *call sites*
is false: the other two textual matches are a comment and a diagnostic string. Its no-production-
caller claim for the **exported** function is true; it says nothing about the private function.

### Q2 — bootstrap: change `New`, not the clock

`registry.Bootstrap` has 1 production / 7 test callers. Its production caller is inside
`daemon.New(cfg)`, which has no ctx at base. `Run(ctx, cfg, announce)` calls New; the CLI calls
Run with `signal.NotifyContext(Background(), ...)`, cancellable but deadline-free (F5, Vcalls).
Thus the predecessor's claim that the startup context was already available *inside New* was
incorrect (Vrefs). It becomes available through an API edit, without a deadline value.

| Option | API cost measured at base | Classification / decision |
|---|---|---|
| **`New(ctx, cfg)`; Run forwards ctx; Bootstrap takes ctx** | New: **1 production + 15 test calls**; Bootstrap: **1 + 7** | **Chosen, policy-free lifetime propagation.** Source-breaking API; no compatibility wrapper that silently manufactures a root. |
| Add `NewWithContext`, retain `New(cfg)` wrapper with Background | Existing calls need not migrate, but wrapper remains a new named deadline-free root | Policy-free only as lifetime compatibility; insufficient for the chosen all-callers migration. Rejected here. |
| Derive a named startup timeout in New/Run | Requires a duration and timeout/failure semantics | **Policy-bearing; Mark's tranche.** |

The prototype keeps the public store argument concrete and delegates to an unexported
`bootstrap(ctx, bootstrapStore, release)` seam solely to witness both reads. The store wrapper
delegates to the real SQLite store; it cancels after the head read, so the object read is proven
to observe that cancellation. New/Run cancellation tests also prove startup failure releases
writer authority by reopening successfully (Vbehavior, M09/M10/M16/M17).

This does **not** make all startup cancellable: store open, archive work, writes, and integrity
scans are outside the promised read-threading surface (Vwrites, Votherroots, Vstartup).

### Q3 — replay and the sixth getter

Add ctx to `ReplayEntry`, `ReplayEpisode`, and `GetVerifyResult`; use `QueryRowContext` in the
cache getter. Episode forwards ctx to entry, then source and cache reads. Base costs: entry
1 production / 1 test call; episode 0 / 18; cache 1 / 13. The entry production call is internal
to episode; there is **no in-repo production root invoking the engine** (F6, Vcalls, Vnocallers).
This is API readiness, not a measured running production deadline.

The prototype uses a private three-method `replayStore` seam while leaving `NewEngine`'s public
store argument unchanged. A test archives a small fake executable, loads a real source object,
cancels immediately after the source read, and observes cancellation at the real cache getter.
It requires no pinned AILANG execution and makes no elapsed-time assertion (Vbehavior).

Leave `runPinnedTransition`'s existing Background→WithTimeout root unchanged and named.
Replay **read** cancellation is the contract here; subprocess cancellation/descendants belong
to OPEN row 112. Do not silently thread the new read ctx through that subprocess as extra scope
(F1, Vrules, Vroots).

## §3 Q4 — empty ratchet plus an independent root census

**Measured baseline:** the controller's grep returns 29 lines, but two are comments; AST parsing
finds **27 executable constructor calls across 58 production files in host/ and cmd/**. In broker/registry/replay,
there are **14** calls: broker 10 (8 reads + 2 mint roots), registry 2, replay 2 (read + subprocess).
After the prototype: **15** roots in host/ and cmd/, **1** in those three packages (F7, Vcalls,
Vroots). These numbers justify strengthening the gate; a direct-argument-only test misses the
very hoisting failure this row is meant to expose.

1. Flip `deadlineFreeReadPins` to `map[string]int{}`. Extend its direct-argument pattern to
   `TODO` and `GetVerifyResult`. Keep its non-vacuity threshold of at least 20 production files.
2. Add `TestProductionContextRoots`: parse all production `.go` files under host/ and cmd/;
   count context `Background`, `TODO`, and `WithoutCancel` selector references by
   **file + enclosing function + constructor**, resolving renamed context imports. Reject a
   dot context import. Scan package-level declarations too. Reference counting also catches
   assigning `context.Background` to a variable before calling it.
3. Exact exception map: the 15 roots in §4. None is permitted in approve.go, registry.go, or
   replay's read entry points. Assert ≥58 scanned files plus seven mandatory anchors:
   approve.go, publish_op.go, registry.go, replay.go, daemon.go, cmd/world-publish/main.go,
   and store.go. Empty or partial scans cannot satisfy the new empty expectation.
   `TestProductionGoSurface` walks the module root (excluding .git and vendor directories),
   rejecting any non-test Go file outside host/ and cmd/ except exactly
   `design_docs/verification/w-race-gate-blindspot/racecontrol/main.go` and
   `design_docs/verification/w-race-gate-blindspot/repro/main.go`. It checks untracked files too
   and requires both allow-listed files to exist. Vsurface measures 58 inside / 2 outside;
   both exceptions have zero `context.` matches. M35 proves a new top-level tree fails.
4. Pair syntax with real-store cancellation and exact-context identity witnesses. The broker
   mint witness requires **5 object + 4 head reads**, covering all six non-publish sites;
   invoke validation requires **2 object reads**. These are dynamic fixture counts, distinct
   from the eight unique source sites (Vbehavior).

**What zero proves:** in the scanned production files, the named getter calls contain no direct
Background/TODO spelling. The root census rejects the drilled introductions in host/ and cmd/:
local-variable Background hoists (M25), renamed `c "context"` imports with `c.Background()`
(M30), `context.WithoutCancel(ctx)` references (M31), function-valued
`f := context.Background; ... f()` references (M32), package-level
`var rootCtx = context.Background()` (M33), and standalone `context.TODO()` roots (M34).
M35 witnesses the separate surface guard against a new non-test Go tree. These are syntax
witnesses, not exhaustive dataflow proofs. With the behavioral tests, the exercised paths pass
their caller's context to real store reads.

**What it does not prove:** a deadline exists; an external caller supplied a non-nil context;
a custom Context implementation obeys the contract; every possible dataflow preserves ctx;
or pool/SQLite lock waits, writes, full startup, replay execution, and human interaction have
an elapsed-time bound. Imports returning contexts and reusing an already-approved root are
not whole-program taint analysis. The composition condition is explicit: *the caller must
supply the intended lifetime, and the read chain must preserve it.* The context-identity and
cancellation assertions pin that condition on the touched chains. The runtime boundary guard
is still needed to reject all deadline-free inputs, including ones source analysis cannot see.

## §4 Every remaining root in host/ and cmd/, with reason

All rows are AST-observed once in the prototype (Vroots); behavior/classification comes from
F3/F5 and Vguardpaths/Votherroots/Vwrites. These are **compatibility exceptions**, not approval
of permanent deadline-free operation. The census pins Background/TODO/WithoutCancel references in host/ and cmd/.

| File / function (`context.Background`) | Why it remains / store implication |
|---|---|
| `cmd/ailang-worldd/cli.go` / `get` | Existing client root; `do` installs client timeout. |
| same / `execute` | Same bounded REST client delegation. |
| same / `executeWithAuth` | `doAuth` installs client timeout. |
| `cmd/ailang-worldd/main.go` / `runServe` | Signal lifetime for daemon Run; no deadline. Reaches bootstrap. |
| `cmd/ailang-worldd/session.go` / `runSessionMint` | Existing CLI mint root; current Mint calls context-free MintSession, not one of the guarded read APIs. |
| same / `runSessionRevoke` | Existing CLI revoke root; forwards to RevokeSession write. |
| `cmd/world-publish/main.go` / `runApprove` | **Explicit relocation of old mint roots**; after attendance check; no deadline. Reaches six unique approval read sites. |
| same / `runPublish` | Existing attended publish root; no deadline before validation's two object reads. |
| same / `runReconcile` | Existing reconcile root; HTTP request derives RequestTimeout; no store read on this root's reconcile path. |
| `host/archive/archive.go` / `probeVersion` | Existing bounded subprocess context, separate lifecycle; row 112. |
| `host/authority/resolver.go` / `Resolve` | Existing compatibility resolver; deadline-free ResolveSession read for commit authorization. |
| `host/capsule/capsule.go` / `Run` | Existing subprocess execution timeout root; not an approval-read parent. |
| `host/daemon/daemon.go` / `drain` | Existing shutdown timeout root; must outlive the serving context's cancellation. |
| `host/replay/replay.go` / `runPinnedTransition` | Existing subprocess timeout root; row 112 owns descendant behavior. |
| `host/store/store.go` / `Commit` | Existing `selectedHeadTx` call in context-free transaction; policy remains with row 23/Mark. |

### What the closing guard would break, measured after threading

Define the audit guard precisely: reject nil or `!ctx.Deadline()` at **eight exported read
methods**: GetObject, GetWorld, GetLogEntry, GetRegistryHead, GetVerifyResult, SelectedHead,
ReadObject, ResolveSession. RevokeSession is a write; selectedHeadTx is internal to Commit.
This audit does not pretend that these eight methods enumerate every database read (Vwrites).

The temporary guard was installed and removed, not shipped. `guard_probe.py` runs four
production-chain probes with their present root lifetime and four bounded controls. **Before:
8/8 pass. With guard: 4 deadline-free arms fail, 4 bounded arms pass** (Vguard).

| Deadline-free production root chain | Static read sites / observed operation | Guard impact |
|---|---|---|
| runApprove → MintAttendedApproval → approve/decide/poll | Six unique sites; full mint fixture performs 9 reads | Approval mint fails at head read; writes can already have happened. |
| runPublish → InvokeAttendedPublish → invoke → validate | Two unique sites; 2 reads when admitted | Fails reading approval decision before dispatch. ExecTimeout offers no protection here. |
| runServe → Run → New → Bootstrap | Two sites; fresh bootstrap head read, existing bootstrap also object read | Daemon startup fails at registry-bootstrap. Signal cancellation is not a deadline. |
| commit middleware → resolver.Resolve → ResolveContext → ResolveSession | One lookup site for well-shaped credentials | Valid token is denied: Resolve collapses read error to DenialUnknown. |

Thus **3 paths / 10 unique sites in the migrated getter family**, or **4 paths / 11 sites**
when the guard includes credential reads. Replay contributes zero production roots, but its
exported API would reject deadline-free callers too. The HTTP GET and projection paths derive
read deadlines from request contexts; they are not these four breakages (Vguardpaths, Vreadsites, Vcomposition).
The probes use real stores and the same lifetime shape; they do not claim to have run the CLI's
interactive TTY or a production deployment. Static root tracing supplies that connection.

A guard on `selectedHeadTx` would additionally reject Commit's internal background read.
That is **not** included in the four-arm result and cannot be smuggled into a read-only guard.
Raw transaction queries, integrity scans, context-free writes, and external API callers require
a separate inventory in the policy tranche (Vwrites, Vstartup). No finite overall wait claim
follows from the audit guard alone.

## §5 Decision for Mark — one word, A or B

**Recommendation: A.** This is a recommendation, not an agent's policy selection. Plumbing in
§2–§3 chooses no deadline value and may land independently. The guard and Commit change must
return to Mark under D-WORLD-18 before their tranche routes (Vrules).

| Answer | One bundled direction for the guard + Commit policy tranche |
|---|---|
| **A — finite active operations** | Preserve the declared strict read guard as the closing move. Require finite budgets for startup/bootstrap, approval active I/O, publish validation, and credential lookup; human think time between commands/polls consumes no active-I/O budget. Design Commit to accept ctx and permit cancellation before durable commit, with all-or-nothing transaction state; an uncertain commit outcome must be reconciled, never blindly retried. |
| **B — preserve admitted work** | Preserve deadline-free attended/startup compatibility and let an admitted Commit finish independently of caller cancellation. This requires Mark to explicitly revise the universal deadline-free-reject closing move into a documented exception policy; do not represent it as satisfying the existing D-WORLD-18 closing move. |

A is recommended because the active approval operations do not wait for a human (F2), and
four measured failures explain exactly what must be bounded before enabling the guard.
Atomicity need not be traded for partial state; the question is cancellation/admission and how
an uncertain durable outcome is reported. The present Commit uses `db.Begin`, a deferred
rollback, transactional writes and `tx.Commit`, without a ctx argument (Vwrites).

**Neither answer authorizes an invented duration in this patch.** After A, the policy designer
must return a concrete bound table (startup, approval, validation/lookup, Commit), cancellation
cutoff and commit-result contract for Mark's ratification. After B, return explicit exceptions
and the amended guard contract. Deadline values and Commit cancellation semantics are
**UNMEASURED/unimplemented here**; a one-word direction must not be laundered into approval of
values Mark has never seen. The guard remains declared, and row 23 remains its named owner.

## §6 Acceptance and executed mutations

No elapsed-time oracle is used by the new behavioral tests. Exact error identity
`errors.Is(err, context.Canceled)`, store-delegating context witnesses, read counts, and live
controls are the oracles (Vbehavior). The harness restores each file in `finally`; all mutants
compiled and failed a named test assertion. Detailed commands/output are `M01.txt`–`M35.txt`;
`mutate.py` and `mutations.json` preserve the mutations and selectors (Vmut).

| Acceptance | Mutation executed on prototype | Assertion / selector | Result |
|---|---|---|---|
| AC1 all eight approval sites preserve ctx | M01–M08 replace each individual getter ctx with Background | `TestApprovalContextIdentity` exact ctx and read census; `TestApprovalCallerCancellation` real error | **8/8 killed** |
| AC2 both registry reads preserve ctx | M09 head and M10 object use Background | `TestBootstrapCallerCancellation`, `TestBootstrapSecondReadCancellation`; exact ctx / canceled real read | **2/2 killed** |
| AC3 replay source and cache preserve ctx | M11 source, M12 cache use Background | `TestReplayCallerCancellation`, `TestReplayCacheCallerCancellation` | **2/2 killed** |
| AC4 cache actually executes with ctx | M13 QueryRowContext→QueryRow | `TestVerifyResultCallerCancellation` expects cancellation, not cache miss | **1/1 killed** |
| AC5 Execute uses its argument on both effects | M14 append helper, M15 poll helper receive Background | Approval identity/cancellation tests | **2/2 killed** |
| AC6 startup caller survives both API hops | M16 New→Bootstrap, M17 Run→New drop ctx | `TestStartupCallerCancellation` expects registry-stage cancellation and successful reopen | **2/2 killed** |
| AC7 episode forwards to entry | M18 ReplayEpisode→ReplayEntry drops ctx | `TestReplayCallerCancellation` | **1/1 killed** |
| AC8 changed API wrappers are not decorative | M19 exported DecideApproval; M20 exported Mint; M21 mint's private decide call; M22 invoke's validation call drop ctx | Approval identity/cancellation tests | **4/4 killed** |
| AC9 empty map cannot pass an empty scan | M23 old scanner's roots list emptied | `TestNoNewDeadlineFreeStoreReads`: scanned-files floor | **1/1 killed** |
| AC10 root census is non-vacuous | M24 new scanner's roots list emptied; M29 skip approve.go | `TestProductionContextRoots`: floor/mandatory anchor | **2/2 killed** |
| AC11 hoisting is visible | M25 assign Background to Execute's local ctx before helpers | `TestProductionContextRoots`: unapproved constructor in Execute | **1/1 killed** |
| AC12 direct literal regression is red at zero | M26 getter gets TODO; M27 gets Background | `TestNoNewDeadlineFreeStoreReads`: count exceeds empty pin | **2/2 killed** |
| AC13 errors from newly threaded getter propagate | M28 cache error arm returns nil | `TestVerifyResultCallerCancellation` | **1/1 killed** |
| AC15 renamed import resolution | M30 rename context to c in approve.go and add c.Background() | `TestProductionContextRoots`: Background=16 | **1/1 killed** |
| AC16 cancellation detachment is visible | M31 add context.WithoutCancel(ctx) | `TestProductionContextRoots`: WithoutCancel=1 | **1/1 killed** |
| AC17 function-valued constructor is visible | M32 assign context.Background to f then call f() | `TestProductionContextRoots`: Background=16 | **1/1 killed** |
| AC18 package-level root is visible | M33 add var rootCtx = context.Background() | `TestProductionContextRoots`: package owner, Background=16 | **1/1 killed** |
| AC19 TODO root is visible independently of getters | M34 add standalone context.TODO() | `TestProductionContextRoots`: TODO=1 | **1/1 killed** |
| AC20 new source tree cannot bypass census | M35 create mutation_surface_probe/root.go | `TestProductionGoSurface`: unexpected non-test Go file | **1/1 killed** |

**Total: 35/35 killed, 0 survived, 0 build-error-only kills.** M13–M22 and M28 are derived from
the shipped signature/wiring/error-handling diff, not merely reinsertion of the original bug.
M30–M34 independently exercise all five added syntactic witnesses. No new mutation survived;
no scanner detection defect was exposed or required a fix. Round 2 adds separate constructor
count logging and the surface guard; every M01–M35 mutation was rerun with restoration in
`finally`, including removal of M35's temporary directory.

AC14 is the release gate: compile all test files, `go vet ./...`, full pinned `go test ./...
-count=1`, and the targeted tests above. Gate results and environmental limitations are in §10.
This row changes no `.ail`; no new AILANG language claim or `verify_ail.sh` run is needed.

## §7 Files and Conflict Surface

Production proposal (under 150 changed production LOC total in the prototype; Vfiles):

- `host/broker/approve.go` — context on Execute and seven approval helper/API functions; eight reads.
- `host/broker/broker.go` — invoke forwards ctx to validation.
- `host/broker/publish_op.go` — ctx-taking mint API and both sessions/private decision.
- `cmd/world-publish/main.go` — documented, pinned CLI mint compatibility root.
- `host/registry/registry.go` — ctx-taking Bootstrap plus private delegating test seam.
- `host/daemon/daemon.go` — ctx-taking New; Run and bootstrap forwarding.
- `host/replay/replay.go` — ctx-taking entry/episode plus private store witness seam.
- `host/store/store.go` — ctx-taking GetVerifyResult and QueryRowContext.

Tests changed/added are enumerated exactly in the patch manifest below; existing call sites
receive explicit test roots. New tests are in broker, daemon, registry, replay, and store;
ratchet/root scans live in store. The temporary guard/probe files are **not** part of the patch.

| Neighbor / surface | Collision and ordering |
|---|---|
| Row 26, bounded Z3 producer | **Semantic/API dependency, not permission to absorb producer work.** Its queue gate requires this row's owned 11→0 work and a known store-wait contract. Publish the weaker contract explicitly; do not infer source/envelope I/O is bounded from zero. Controller should resolve the dependency checkpoint separately from row 23's still-open policy tail (Vrules). |
| Row 111 arm A, per-request busy_timeout | **Yes**, store read methods/cache and daemon constructor tests overlap. Sequence arm A after this threading tranche; it must retain pool-state restoration and its own status/lock-wait proof. Row 23 does not claim deadline dominance (Vrules, Vrefs). |
| Row 112, archive/replay subprocess groups | **Same replay.go file**, different function. Keep runPinnedTransition unchanged; reconcile context signatures independently. Its existing root stays named. Archive has no edit here (Vfiles, Vrules). |
| Row 113, pkgproj descendants | **No proposed file overlap**; its subprocess lifetime remains its own owner (Vfiles, Vrules). |
| Row 25, blocked-child kill coverage | **No proposed file overlap**; no capsule kill/output fixture is altered. Do not count read cancellation as blocked-child coverage (Vfiles, Vrules). |
| Existing callers / S7 | `New`, Bootstrap, DecideApproval, MintAttendedApproval, ReplayEntry/Episode, GetVerifyResult become source-incompatible. Migrate the measured in-repo calls together. External consumers are UNMEASURED. Add package API usage notes naming “read ctx only,” and a runnable cancellation example via the new tests; update any documented API snippets found at implementation time. |

## §8 Milestones, each ≤150 production LOC

1. **M1 approvals (0.3d; ≤80 production LOC).** Thread all eight sites and mint chain; explicit
   CLI root; reduce approval pin by 8. AC1/AC5/AC8; run `go vet ./...`,
   `go test ./host/broker ./cmd/world-publish -count=1` with pinned binary/loopback. Do not add a duration.
2. **M2 bootstrap + replay/cache (0.3d; ≤70 production LOC).** Migrate New/Run/Bootstrap and
   replay/cache APIs and tests, add two small private witness interfaces; reduce remaining
   pins to zero. AC2/AC3/AC4/AC6/AC7/AC13. `go vet ./...`; targeted cancellation tests;
   `go test ./host/daemon ./host/registry ./host/replay ./host/store -count=1` with pinned binary.
3. **M3 honest zero and handoff (0.15–0.4d; 0 additional production logic LOC).** Land empty
   map, root census and root reasons, API usage notes, full mutation matrix and gates.
   AC9–AC12/AC14; `go vet ./...`; `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1`.
   Publish the four-path guard impact and return §5 to Mark. No guard or Commit change.

For individually green milestone commits, pin reductions must accompany the related wiring;
install the final 15-root census only once M2's complete root state exists. The banked prototype
is the combined final state. A separate policy tranche is not sized here because its values
and cancellation contract have not been chosen by Mark.

## §9 Residuals / named owners

| Residual | Named OPEN owner / required next action |
|---|---|
| Strict guard; finite startup/approval/credential budgets; Commit cancellation and outcome contract | **Row 23**, policy tranche, **Mark** ratifies §5 direction and the later concrete bounds. Includes the newly measured authority compatibility lookup; don't hide it behind the original 11-site list. |
| All-write/whole-startup boundedness beyond these getters | **Row 23**, policy inventory before any repo-wide claim; context-free integrity scans/writes were not converted. If split later, controller must create and verify an OPEN owning row first. |
| Actual lock-wait deadline dominance, late 200 and 500/503 behavior | **Row 111**, arm A after threading; no lock-bound claim from these cancellation tests. |
| Producer source-read/envelope-write contract | **Row 26**, use the explicit weaker store contract; do not add a twelfth hidden root. |
| Archive and replay descendant cleanup | **Row 112**. |
| pkgproj descendant lifecycle | **Row 113**. |
| Live child blocked in write coverage | **Row 25**. |

All six numbered owners above are positively located as OPEN queue entries in Vrules. Mark is
the policy decision maker, not a substitute for an open queue owner. Existing subprocess,
transport and other approved compatibility roots in §4 are explicitly preserved; this doc
makes no claim to newly bound them.

## §10 Verification Log

Commands were executed in this worktree, except the explicitly absolute artifact/tool paths.
**F/V source rows before edits describe `b5b4a7d`; prototype rows are labeled.** Artifacts retain
full commands and observed output. Negative caller results are paired with distinct positive
callers, not with matches of the same missing name. AST counts exclude comments and string
literals. No git write command was used.

| ID | Executed command | Observed output / supported claim |
|---|---|---|
| F1 | `grep -n "context.Background()" host/broker/approve.go host/registry/registry.go host/replay/replay.go` | Baseline: approve 169/195/225/234/251/261/495/522 (8); registry 128/136 (2); replay 153 (1), plus subprocess root 323. Confirmed. |
| F2 | `sed -n '101,145p' host/broker/approve.go; sed -n '245,258p' host/broker/broker.go` | Execute signature discards ctx; Approve returns pending after PutObject/append; PollApproval calls find, then pending/decided; broker passes ctx at :258. Confirmed request/poll, no human wait in reads. |
| F3 | `sed -n '120,170p' host/broker/publish_op.go; sed -n '240,290p' host/broker/publish_op.go; sed -n '345,375p' cmd/world-publish/main.go; sed -n '88,105p' host/broker/handlers.go; rg -n '\.Invoke\(' host cmd --glob '*.go' --glob '!**/*_test.go'` | Mint sessions Invoke Background twice; private decide call; InvokeAttendedPublish forwards ctx to session; CLI runPublish supplies Background. handlers.go derives runCtx only inside run subprocess helper. Three non-test Invoke call sites, all publish_op.go; no capsule/HTTP invocation parent. Positive controls: mint and publish calls separately. |
| F4 | `grep -rn "DecideApproval(" --include='*.go' host cmd \| grep -v _test.go; rg -n 'decideApproval\(' host/broker/publish_op.go` | Exported production grep: definition only; independent control: private decideApproval called at publish_op.go:151. No exported production caller confirmed; actual test count corrected in Vcalls. |
| F5 | `rg -n 'func New\|func Run\|registry.Bootstrap\|New\(cfg' host/daemon/daemon.go; sed -n '199,210p' cmd/ailang-worldd/main.go` | New(cfg) :454; registry.Bootstrap :510; Run(ctx) :798 calls New(cfg) :799. CLI :203 signal.NotifyContext(Background), then Run. Confirmed. |
| F6 | `rg -n 'ReplayEntry\(\|ReplayEpisode\(\|GetVerifyResult\(' host cmd --glob '*.go' --glob '!**/*_test.go'; sed -n '781,806p' host/store/store.go` | ReplayEntry/ReplayEpisode definitions and only internal episode→entry call; GetVerifyResult sole caller replay.go:191, definition takes no ctx and calls QueryRow. Distinct positive cache-call control in the same command. Confirmed. |
| F7 | `grep -rn "context.Background()\\|context.TODO()" --include='*.go' host cmd \| grep -v _test.go; grep -rn "context.Background()\\|context.TODO()" --include='*.go' host cmd \| grep -v _test.go \| wc -l` | 29 matching non-test lines, including resolver.go:123 and handlers.go:258 comments. Executable census is 27, not 29. No TODO hits; separate Background matches are positive control. |
| F8 | `head -9 go.mod; go version` | go.mod go 1.26.6; go version go1.26.6 darwin/arm64; modernc.org/sqlite v1.54.0. Confirmed toolchain; compile fence and gates below. |
| Vrules | `sed -n '4251,4285p' design_docs/world-mission.md; grep -n '\| D-WORLD-18 \\|\| D-WORLD-23 ' design_docs/world-mission.md; sed -n '4336,4350p' design_docs/world-mission.md; sed -n '1831,1835p' design_docs/world-mission.md; sed -n '4313,4334p' design_docs/world-mission.md` | Rows 23,25,26,111,112,113 located as open entries; D-WORLD-18 and D-WORLD-23 ratified text read. Row26 gated on row23; row111 armA after23. Row112 archive/replay, row113 pkgproj, row25 capsule blocked writer. No owner is inferred from an absent grep. |
| Vreads | `sed -n '149,277p' host/broker/approve.go; sed -n '485,531p' host/broker/approve.go; sed -n '113,164p' host/registry/registry.go; sed -n '151,263p' host/replay/replay.go; sed -n '361,485p' host/store/context_read_test.go` | Read all 11 sites, helper flows, replay source/cache chain, and original regexp/map. Regexp recognizes direct Background arguments only; map 8/2/1; file floor 20. |
| Vwrites | `rg -n 'func .*Commit\|Begin\(\|BeginTx\|selectedHeadTx\|\.Commit\(' host/store/store.go host/daemon/handlers.go; sed -n '845,955p' host/store/store.go; rg -n 'ResolveContext\|func .*Mint\|func .*Revoke' host/authority; rg -n 'ReadObject\(\|WriteObject\(\|GetObject\(\|GetRegistryHead\(' host/evidence host/transitionreg --glob '*.go' --glob '!**/*_test.go' ` | Commit(c) at :873, Begin(), deferred Rollback(), selectedHeadTx(Background), tx.Commit(); daemon sole Store.Commit call at handlers.go:600. GetVerifyResult QueryRow baseline; authority ResolveSession is separate ctx read, RevokeSession write; ReadObject bounded-reader surface. No Commit ctx migration claimed. |
| Vrefs | `sed -n '111,233p' design_docs/implemented/w-daemon-lock-wait-not-deadline-bound.md; sed -n '610,641p' design_docs/DESIGN.md; sed -n '270,314p' design_docs/implemented/w-daemon-read-cancellation.md; sed -n '554,579p' design_docs/implemented/w-daemon-read-cancellation.md` | Read reference §3–§11 and predecessor §2.8/§10. Predecessor says startup context available, but F5 refutes availability inside New. Reference declares runtime lock bound deferred; DESIGN §14 excludes driver/compiler/live binary self-replacement. |
| Vguardpaths | `sed -n '229,261p' cmd/world-publish/main.go; sed -n '360,374p' cmd/world-publish/main.go; sed -n '199,212p' cmd/ailang-worldd/main.go; sed -n '40,58p' host/daemon/middleware.go; sed -n '123,166p' host/authority/resolver.go; rg -n 'WithTimeout\|ResolveContext' host/projection/projection.go; sed -n '265,274p' host/daemon/handlers.go; sed -n '88,102p' host/broker/handlers.go` | CLI attendance before mint; runPublish Background; signal lifetime for Run; commit middleware Resolve; resolver forwards Background to ResolveSession and converts error to DenialUnknown. HTTP GET readCtx and projection request contexts derive WithTimeout; handler timeout occurs later. |
| Votherroots | `sed -n '30,43p' cmd/ailang-worldd/cli.go; sed -n '95,110p' cmd/ailang-worldd/cli.go; sed -n '465,489p' host/broker/registry_reconcile.go; sed -n '52,74p' host/authority/mint.go; cat host/authority/revoke.go; sed -n '771,786p' host/daemon/daemon.go; sed -n '186,202p' host/capsule/capsule.go; sed -n '465,476p' host/archive/archive.go; sed -n '326,340p' host/replay/replay.go; sed -n '918,940p' host/store/store.go` | REST do/doAuth derive client timeout; reconcile NewRequestWithContext uses RequestTimeout; Mint calls MintSession without ctx, Revoke forwards ctx to write; drain/archive/capsule/replay roots each derive an existing timeout; Commit root persists. |
| Vnocallers | `rg -n 'ReplayEntry\(\|ReplayEpisode\(\|DecideApproval\(\|NewValidator\(' host cmd --glob '*.go' --glob '!**/*_test.go'; rg -n 'registry.Bootstrap\(\|broker.MintAttendedApproval\(' host/daemon/daemon.go cmd/world-publish/main.go` | Prototype production results: ReplayEntry only called by ReplayEpisode; no outside replay invocation; exported DecideApproval definition only; NewValidator definition only. Independent positive controls: registry.Bootstrap in daemon and MintAttendedApproval in CLI. |
| Vpolicy | `git diff --unified=0 -- '*.go' \| grep -E '^\+.*(WithTimeout\|WithDeadline\|func.*ctx\|Bootstrap\(ctx\|MintAttendedApproval\(context.Background)' ` | Added context signatures/Bootstrap forwarding and explicit CLI mint root are present. No added WithTimeout/WithDeadline in tracked Go diff; ctx signatures are distinct positive controls. New tests have deadlines solely as bounded controls. No duration was added to production. |
| Vroots | `go run ~/.ailang/state/world-iter194/design/census.go`; `rg 'root references' ~/.ailang/state/world-iter194/design/M3[0-4].txt` | Re-run from banked source: 58 host/cmd production files; **Background=15, TODO=0, WithoutCancel=0**. Positive mutation controls: M30/M32/M33 **16/0/0**, M31 **15/0/1**, M34 **15/1/0** (Background/TODO/WithoutCancel). All five fail TestProductionContextRoots. Broker/registry/replay retains only runPinnedTransition. Census counts calls, not comments/strings; scanner mutation controls count selector references. |
| Vsurface | `git ls-files -- '*.go' ':!:*_test.go' \| grep -cvE '^(host|cmd)/'`; same with `-cE` and `-vE`; `grep -c "context\." design_docs/verification/w-race-gate-blindspot/racecontrol/main.go design_docs/verification/w-race-gate-blindspot/repro/main.go`; `go test ./host/store -run '^TestProductionGoSurface$' -v -count=1` | **2 outside, 58 inside**. Outside files are exactly the two §3 reproducers; context matches **0 each** (grep exit1 means no matches). Walk guard passes with 58 inside / 2 allowed. M35's new top-level non-test Go file fails the guard; restored in finally. |
| Vfiles | `git diff --numstat; git ls-files --others --exclude-standard` | 8 modified production files; final production diff +67/-46 = 113 changed lines. 28 Go files total including six new tests. Exact final list in §11. No archive/pkgproj/capsule production file in patch; broker/replay/store/daemon positive controls visible. |
| Vbehavior | `go test ./host/broker ./host/registry ./host/replay ./host/daemon ./host/store -run 'Test(ApprovalContextIdentity\|ApprovalCallerCancellation\|Bootstrap.*Cancellation\|Replay.*Cancellation\|StartupCallerCancellation\|VerifyResultCallerCancellation\|ProductionContextRoots\|NoNewDeadlineFreeStoreReads)$' -v -count=1` | 10 top-level named tests execute, with 9 additional subtest RUN lines: all pass across 5 packages. Literal ratchet 0 over 58 files; census 15 roots; cancellation and exact-context witnesses green, including second registry/cache reads. |
| Vstartup | `sed -n '454,520p' host/daemon/daemon.go; sed -n '550,607p' host/daemon/daemon.go; rg -n 'Query\|func ' host/store/scan.go; rg -n 'store\|GetObject\|ReadObject\|GetRegistryHead' host/broker/registry_reconcile.go; rg -n 'NewRequestWithContext\|RequestTimeout' host/broker/registry_reconcile.go` | New still opens store and optionally archives, then Bootstrap and scanIntegrity. scanIntegrity invokes ScanUnreadableLog/Worlds with counters/time checks between calls. Reconcile search has no store/getter hits; distinct positive RequestTimeout/NewRequestWithContext hits show HTTP path. |
| Vcalls | `go run ~/.ailang/state/world-iter194/design/census-round1.go` (banked original tool used before edits) | `census-before.txt`: production/test calls Bootstrap 1/7, DecideApproval 0/1, decideApproval 2/7, MintAttendedApproval 1/0, mintAttendedApproval 1/4, daemon.New 1/15, ReplayEntry 1/1, ReplayEpisode 0/18, GetVerifyResult 1/13; 58 production files; 27 roots. Distinct Bootstrap/mint production calls control negative replay/exported-decide counts. |
| Vguard | `python3 /Users/voightkampff/.ailang/state/world-iter194/design/guard_probe.py` | `guard-control.txt` exit0, eight arms pass. `guard-reject.txt` exit1: mint/free, publish/free, startup/free, resolve/free fail; bounded four pass. Temporary eight-method guard restored in finally; probe sources banked under broker/daemon/authority artifact folders. |
| Vmut | `python3 /Users/voightkampff/.ailang/state/world-iter194/design/mutate.py` | `mutations-run.txt`: M01–M35 KILLED; TALLY 35 35. Each transcript contains `--- FAIL:`; none counted on compile failure. |
| Vreadsites | `rg -n '\.(GetObject\|GetWorld\|GetLogEntry\|GetRegistryHead\|GetVerifyResult\|SelectedHead\|ReadObject\|ResolveSession)\(' host cmd --glob '*.go' --glob '!**/*_test.go'; rg -n 'ValidateProof\(\|ReadSnapshot\(\|\.Publish\(' host cmd --glob '*.go' --glob '!**/*_test.go'; rg -n 'WithTimeout\|objectReadTimeout' host/evidence/validator.go host/projection/projection.go` | Prototype enumerates 34 exported-read call sites: all pass a ctx variable; validator derives ObjectReadTimeout, projection derives maxWait. API-only transition Publish/validator and replay are distinguished from in-repo production roots; positive ReadSnapshot/bounded-read calls control absent root claims. |
| Vcomposition | `rg -n 'readCtx\(\|func .*Workbench\|func .*workbench' host/daemon/workbench.go; rg -n 'NewReplaySession\(\|newSession\(' host/broker/broker.go host/broker/publish_op.go; rg -n 'func .*ScanUnreadable\|Query\(' host/store/scan.go; git rev-parse --short HEAD` | Workbench calls d.readCtx; Replay-mode constructor is separate from the three Live publish sessions. ScanUnreadableLog/Worlds use context-free Query. HEAD b5b4a7d. These are retained composition limits, not converted reads. |
| Vcompile | `go test ./... -run '^$'` | `compile.txt`: 23 ok, 0 FAIL. Deliberately compilation-only, not behavioral evidence. |
| Vvet | `go vet ./...` | `vet.txt`: empty diagnostic output, exit0. Positive compilation control Vcompile has 23 packages. |
| Vvet2 | `go vet ./...` | `vet-round2.txt`: exit0, no diagnostics after round-2 changes. |
| Vtarget2 | `go test ./host/store -run '^(TestProductionContextRoots\|TestNoNewDeadlineFreeStoreReads\|TestProductionGoSurface)$' -v -count=1` | `targeted-round2.txt`: all 3 pass, exit0. Ratchet 0 / 58 files; roots Background=15 / TODO=0 / WithoutCancel=0; surface 58 inside / 2 allowed. These gates were not sandbox-affected. |
| Vtest0 | `go test ./... -count=1` in restricted sandbox, AILANG_BIN unset | `test.txt`: 17 ok, 6 FAIL. Loopback bind refused and pinned-binary-required tests fail. **UNINFORMATIVE UNDER SANDBOX**; controller independently reran the pinned gate outside the sandbox (Vtest). |
| Vpin | `/Users/voightkampff/.pinned-ailang/ailang --version` | AILANG v0.41.0, commit 24ee1088776e21cd06a3781ed18e77f40be06db3. System PATH binary instead reports v0.43.1-8-ga2256b1c5-dirty; not used for full Go validation. |
| Vtest | `AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang go test ./... -count=1` with loopback-enabled execution | `test-pinned.txt`: **23 ok, 0 FAIL, exit0** (round-1 prototype; controller also independently confirmed 23 ok / 0 FAIL outside the sandbox before this revision). Pin is needed by full Go tests even though no standalone verify_ail run is requested. |

The controller's **F4 test-call count is the sole refuted F1–F8 numerical assertion**. F7 is
correct as a textual line count, but would be false if interpreted as 29 executable roots.
The policy-free reading of context threading holds for the prototype: no deadline chosen,
Commit unchanged, guard removed after the experiment. The qualification is behavioral, not
policy: cancellation can now abort these reads; earlier context-free writes may persist. This
is not an atomic approval-operation guarantee. The predecessor's availability-inside-New
premise is corrected explicitly in §2, rather than inherited.

## §11 Exact prototype manifest and usage

Round 2 extends `host/store/context_roots_test.go` with separate constructor counts and
`TestProductionGoSurface`; the manifest remains **28 Go files, including six new test files**.
The banked patch contains these Go files (tracked diff plus explicit `git diff --no-index
/dev/null` additions for new tests; an ordinary git diff alone would omit those tests):

- `cmd/world-publish/main.go`
- `host/broker/approve.go`
- `host/broker/approve_test.go`
- `host/broker/broker.go`
- `host/broker/context_thread_test.go`
- `host/broker/episode_test.go`
- `host/broker/handlers_test.go`
- `host/broker/publish_op.go`
- `host/broker/publish_op_test.go`
- `host/broker/registry_publish_test.go`
- `host/daemon/bench_test.go`
- `host/daemon/context_thread_test.go`
- `host/daemon/daemon.go`
- `host/daemon/daemon_test.go`
- `host/daemon/handlers_test.go`
- `host/daemon/integrity_test.go`
- `host/daemon/lock_wait_order_test.go`
- `host/registry/context_thread_test.go`
- `host/registry/registry.go`
- `host/registry/registry_test.go`
- `host/replay/context_thread_test.go`
- `host/replay/replay.go`
- `host/replay/replay_test.go`
- `host/store/context_read_test.go`
- `host/store/context_roots_test.go`
- `host/store/store.go`
- `host/store/store_test.go`
- `host/store/verify_context_test.go`

The API migration pattern is `daemon.New(ctx, cfg)`, `registry.Bootstrap(ctx, s, release)`,
`broker.DecideApproval(ctx, s, requestRef, decision, operator, now)`,
`broker.MintAttendedApproval(ctx, s, plan)`, `engine.ReplayEpisode(ctx, episode)`,
`engine.ReplayEntry(ctx, episode, index, entry)`, and `s.GetVerifyResult(ctx, fn, interpreter)`.
Use a caller's existing context; these signatures supply no default timeout. The executable
usage/cancellation examples are the new tests; run the Vbehavior command verbatim. Production
package comments and the explicit CLI-root comment state the restricted scope. M3 must carry
these usage notes into the relevant package documentation, not leave public API knowledge only
in tests.

## §12 Axiom / standards check and review status

Adapted from the design skill's rubric (design assessment, not a claim of measured runtime
boundedness): A1 determinism 0; A2 replayability 0; A3 effect legibility +1 (explicit caller
lifetime); A4 authority 0; A5 bounded verification 0 (no global-bound claim); A6 concurrency 0;
A7 machines first +1 (non-vacuous gates); A8 syntax 0; A9 cost visibility 0; A10 composition +1
(context preservation); A11 structured failure +1 (observable cancellation); A12 boundary 0.
Net +4; no proposed A1/A3/A4/A7 violation. S1/S4/S5: no .ail change; S2/S3: host plumbing;
S6: null-case and mutation probes; S7: M3 public API notes are required.

The [design-doc-creator skill](/Users/voightkampff/.claude/skills/design-doc-creator/SKILL.md)
says **“Unattended (mission-loop) docs: ALWAYS run it.”** The freeze trigger also applies
to §5. Quorum status and any absent reviewers are recorded below; no external review verdict
is a substitute for Mark's policy decision.

### Quorum log

Round 1: **BLOCKED 3/3, all present** — oc-glm-5-3 reject, oc-kimi-k3 reject,
claude-sonnet-5@claude-p reject. Author vendor OpenAI benched; gemini-3-1-pro reserve,
not seated. Controller artifact:
`.ailang/state/mission-quorum/w-store-deadline-free-residue-owner-2026-09-26T09-57-25Z.json`;
banked copy `~/.ailang/state/world-iter194/quorum_r1.json`.

- R1: missing scanner witnesses — M30–M34 all killed, per-constructor counts logged in Vroots,
  and §3's claim narrowed to the drilled forms; §6's old hedge removed.
- R2: scan domain did not exhaust non-test Go — Vsurface measures 58 inside / 2 reproducers
  outside; claims scoped to host/ and cmd/, explicit two-file allow-list guarded by
  TestProductionGoSurface, with M35 killed on a new top-level Go file.
- R3: census cited a perishable path — §10 now cites banked census sources; Vroots rerun
  from the durable path, with the original tool preserved as census-round1.go for Vcalls.

This revision answers the recorded objections; no round-2 quorum verdict is claimed.
§5's policy ratification remains separately owed to Mark.

Local completion record: doc, 28-file prototype patch, 35 mutation transcripts, eight guard
probe arms, root/caller censuses, compile/vet and full pinned Go test logs are banked. No runtime deadline guard,
deadline constant, or Commit cancellation implementation remains in the patch. The new
source-surface test guard is distinct from the unshipped runtime deadline guard.

Final artifact check: `patch --dry-run -R -p1 < ~/.ailang/state/world-iter194/design/proto.patch`
returned exit0 (`patch-check.txt`), so the banked patch, including new test files, matches the
prototype tree. `git diff --check` returned exit0 with no diagnostics. These are artifact
checks, not substitutes for Vbehavior/Vtest.
