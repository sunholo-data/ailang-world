# w-bounded-z3-report-producer — a bounded, pinned-executable Z3 proof-report producer that mints validator-authenticable envelopes

- **Status**: PLANNED (designed, iteration 251; not implemented)
- **Queue row**: 26, `w-bounded-z3-report-producer` · clause-2 · **v1.0.0 release prerequisite** — ruled routable by attended `D-WORLD-69` = **A** (2026-10-09, exception path 26 → 105 → 172 → 150). Shed from queue item 17 by attended `D-WORLD-24` arm **A** (Mark Edmondson, issue #68, 2026-08-20T16:04:52Z).
- **Filed**: 2026-10-12, iteration 251 (designer role, unattended mission stage)
- **Measurement base**: `ba79d1e` (this worktree HEAD; working tree carries only this document and its evidence directory — V1)
- **Instrument**: `$HOME/.pinned-ailang/ailang`, AILANG **v0.41.0** (`24ee108`), sha256 `1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f` — the CI pin (`ci.yml:134-173` asserts `AILANG v0.41.0`, V32)
- **Estimate**: **~2.0 d** (5 milestones, each ≤ 160 production LOC). The queue row carries 1.0 d from item 17's measured cut and prices the round-12 fixes as "unpriced, pending design" — this design prices them (§Timeline).
- **Scope**: the former §3.4 of `w-validated-proven-evidence-boundary` in full — `host/evidence/proof_producer.go`; pinned executable-byte/version checks; bounded process-group execution; independent stdout/stderr limit+1 caps; strict `ai-check` JSON interpretation; required verified identities; MAC-tagged canonical envelope production; source-object reading; bounded envelope storage. Inherited identifiers **AC5**, **AC20**, **M15**, **M28** keep their numbers; new criteria/mutations number after the parent's sets (AC23+, M31+).

---

## §0 The two inherited round-12 objections, answered by name

`D-WORLD-24` arm A shed the producer with **both round-12 objections unresolved, verbatim** (parent §10.15; charter row 26). They are the first thing this design answers, because the API exists to answer them.

### 0.1 `gpt5-6-sol` — the missing reader + the derived deadline

> *"`Producer.GenerateProof` must read `sourceRef` through `ObjectReader`, but `NewProducer(key, checker, writer, execTimeout, maxOutputBytes)` receives no reader. Even if a reader were added implicitly, both the source read and `WriteObject` use the caller context without requiring or deriving a deadline, so `context.Background()` can block indefinitely on the store's single-connection pool. AC20 only proves propagation of an already-deadlined context, not a bounded operation."*

**Answered in three layers, all in this design:**

1. **The reader is a first-class constructor parameter.** `NewProducer(key, reader, writer, checker, cfg)` takes the *same* `ObjectReader` seam the validator takes (validator.go:50-53, V11). `GenerateProof` reads `sourceRef` through that reader and through nothing else — there is no implicit reader, no ambient store, no package-level handle.
2. **The producer derives its own deadlines, never trusting the caller's.** `ProducerConfig` requires positive `ObjectReadTimeout` and `ObjectWriteTimeout`; the constructor refuses zero/negative values (`ErrInvalidProducerConfig` per-field refusals, AC5). `GenerateProof` derives `readCtx := context.WithTimeout(ctx, ObjectReadTimeout)` before the source read and `writeCtx := context.WithTimeout(ctx, ObjectWriteTimeout)` before the envelope write — **including under `context.Background()`** (§D6, the source-read and envelope-write paragraphs). What each bound buys is stated in §D7: expiry is an explicit operational error, nothing is minted, nothing is written.
3. **The store's live strict guard is defense in depth, not the mechanism.** Since this objection was written, row 23's policy tranche landed the `D-WORLD-37` = A closing move: every exported store I/O method calls `requireDeadline(ctx)` and returns `store.ErrNoDeadline` for a nil or deadline-free context (store.go:236-247; 38 guarded call sites, V8). A producer that forgot the derivation therefore fails **loudly at the store boundary** — `context.Background()` can no longer block indefinitely on the pool at all. The producer still derives its own bounds because the caller's deadline (if any) may be hours: the producer's contract is its own bound, and AC23/M31 make the derivation itself falsifiable (a mutant that widens the derived deadline reds on error identity, not on the store's refusal).

### 0.2 `gemini-3-1-pro` — the unpinned lock wait and AC20's decoy vacuity

> *"Round 11b's additive `WriteObject(ctx, o)` introduces an unpinned lock-contended wait. Because `NewProducer` configures no store timeout, it cannot pin the ordering against the store's busy window at construction. If the caller's context to `GenerateProof` has a deadline shorter than `busy_timeout` (2000 ms), SQLite lock contention will ignore it and overrun the deadline. AC20 vacuously masks this by testing only connection-pool waiting with a decoy, reproducing the exact fixture defect rejected in Round 10 for the read side."*

**Answered in four parts:**

1. **The additive `WriteObject` store surface is not built, and the reason is measured.** Round 11b proposed it when `PutObject` was context-free (parent V51). At `ba79d1e`, `PutObject(ctx, o)` exists, is context-threaded, is deadline-guarded (`requireDeadline` at store.go:547, V8), and carries the ratified three-outcome commit contract (`beginDurable`/`finishDurable`/`UncertainError`, durable.go:10-96, V12). A second write surface would create the parallel-path defect this mission keeps closing — one bounded write path and one nobody re-checks. The producer's `ObjectWriter` seam is an *interface satisfied by the unchanged store*: `PutObject(ctx, store.Object) error` + `BusyTimeout() time.Duration` (§D1's API block). **No store file changes in this design** (V5 records that `WriteObject` reads 0 repo-wide with a firing control).
2. **Both orderings are pinned at construction against run-time values.** `NewProducer` reads `BusyTimeout()` from its reader and writer and refuses a positive window ≥ `ObjectReadTimeout` **or** `ObjectWriteTimeout` with `ErrUnorderedTimeouts`, naming both values (AC30) — the exact mirror of the validator's ratified AC18 pin, extended to the write side. The pin binds the values actually used at run time, never two literals; a wrapper around the real store cannot erase it (AC30's wrapper arm, mirroring the validator's AC22).
3. **The claim is narrowed to exactly what is proven — the `D-WORLD-22` arm-B composition statement, symmetric for the write.** Every wait the producer's OWN code performs (pool acquisition, statements, subprocess) is bounded by its derived deadlines. A LOCK-contended wait inside the store is bounded by `busy_timeout`, which this design does not change; the composition is safe only while `busy_timeout < ObjectReadTimeout` and `busy_timeout < ObjectWriteTimeout` — pinned at construction — and lock-wait deadline dominance remains owned by OPEN queue row **111** (`w-daemon-lock-blocked-read-status-and-deadline-enforcement`, charter line 1923, V30). If the *caller* supplies a deadline shorter than the busy window, a lock-contended write may overrun that caller deadline up to the busy window: that is the caller's authority, the same residual the validator's read side has carried since round 11, and it is declared here rather than absorbed (§D7).
4. **AC20's own text states the decoy's scope, and the stimulus is designed to demonstrably block the pool, not the lock.** The criterion itself (§7, AC20) now says: the decoy is one in-flight `GetObject` of a large decoy object holding the store's **single pooled connection** (`SetMaxOpenConns(1)`, store.go:382, V9); the stimulus **asserts the decoy read is in flight** when the write attempts; it holds **no SQLite write lock**; the blocked write's wait is at connection acquisition, which `database/sql` bounds by the supplied context by stdlib contract (`beginDurable`'s own comment, durable.go:48-51, V12); and the criterion **claims nothing about lock contention** — parent §10.12 records why that distinction is load-bearing (the round-10 read-side rejection was exactly a decoy that exercised the wrong wait).

Both answers are falsifiable: AC5/AC20/AC23-AC30 with mutations M15/M28/M31-M39 (§8), each anchored to the diff this document designs.

---

## Axiom Compliance

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +1 | The envelope is a pure function of (key, canonical report bytes); the report is a pure function of the verified run; the result is content-addressed (`sha256` of canonical envelope bytes) |
| A2: Replayability | 0 | No replay machinery touched; the receipt→`CLAIMED` mapping and replay paths are unchanged |
| A3: Effect Legibility | +1 | Every effectful operation is host-boundary Go in one new file; the AILANG surface is untouched (no `.ail` change, V28) |
| A4: Explicit Authority | +2 | The producer mints DATA (an authenticated envelope), never authority: `ValidatedEvidence` still has no exported fields or public constructor and `Validator.ValidateProof` remains the sole mint (validator.go, V11). Pinned byte/version checks refuse executable substitution explicitly |
| A5: Bounded Verification | +1 | Every wait the producer owns is bounded by a named, constructor-required constant (read/write/exec) plus per-stream output caps and a bounded process group; refusal, never partial trust (§D7) |
| A6: Safe Concurrency | 0 | No new shared production state; test goroutines mirror the existing realstore decoy tests |
| A7: Machines First | +1 | Evidence is machine-checkable canonical JSON bytes with a MAC; refusals are typed sentinels with stable messages |
| A8: Minimal Syntax | 0 | No language surface at all |
| A9: Cost Visibility | +1 | Per-proof cost is visible and bounded: the byte-hash of a ~100 MB binary, one bounded `--version` probe, one bounded `ai-check` run, one read and one write — each with its own refusal on overrun |
| A10: Composability | +2 | Reuses `ObjectReader`, the report/envelope codecs, `procbound`, `proctest`, `childenv`, and the unchanged store; the validator consumes the producer's output with zero validator changes |
| A11: Structured Failure | +1 | Every refusal branch has its own sentinel and its own named test; the write's three outcomes (committed / not-committed / uncertain-reconcile) follow the ratified `D-WORLD-40` contract |
| A12: System Boundary | +2 | The proof producer sits exactly where the parent drew the boundary: solver execution, byte verification, and envelope minting are Go-host effects; the kernel keeps only the untrusted `ProofReceipt → CLAIMED` claim |

**Net Score: +12** → **Decision: move forward.** Hard axioms A1/A3/A4/A7 all ≥ 0.

---

## Problem Statement

**Current state (measured at `ba79d1e`):**

- Item 17's tranche 1 landed the validator, the seal, and the boundary — and **deliberately lost the end-to-end demonstration**: tranche 1 "has no in-repository path from an `ai-check` run to an authenticated envelope; all validator fixtures are hand-authored" (parent §10.16). No production Go code runs `ai-check` at all (V6: 0 non-test occurrences under `host/` and `cmd/`, control fires 13× in `scripts/verify_ail.sh`), and no producer surface exists (V4: `NewProducer`/`GenerateProof`/`ObjectWriter`/`proof_producer` read 0 in production Go, control `NewValidator` fires 26×).
- Every proof envelope in the repository today is hand-authored test data (`host/evidence/authority_test.go:48` `envelopeFor`). Nothing demonstrates that a holder of the host key can actually mint an envelope **from a real `ai-check` result** that the validator then accepts — the property the parent's §3.3 step 11 explicitly left to this row: *"Queue item 26 owns demonstrating that the bounded producer is such a holder and minted them from an `ai-check` result."*
- The strict interpretation this producer must implement is not theoretical. Measured on the pinned binary (V24, V25, V27): a **skipped** contract run exits **rc 0** with `verify.verified == 0`; a **contract-free module** exits **rc 0** with every counter zero and `results: []`. A caller that trusted exit codes — or counters alone — would accept a run that proved nothing. The inherited claim "a Z3 encoding error exits 0 silently" (measured on v0.30.0, parent V10) has *changed shape* on v0.41.0: the encoding-error fixture I could construct exits rc **1** (V25) — which strengthens, not weakens, the JSON-parse-not-exit-code discipline, because which failures are silent is a per-release property of the tool, not a contract.
- The two inherited round-12 objections (§0) name real API holes in the shed spelling: no reader parameter, no derived deadlines, an unpinned write-side lock wait, and a decoy criterion that masked the distinction.

**Impact:** without this row, `PROVEN` evidence cannot be produced in-repo at all — the validator exists to authenticate envelopes nothing can mint. This row is the v1.0.0 release prerequisite that closes that loop under the store's live (weaker-then-guarded) contract.

## Goals

**Primary goal:** a `host/evidence` producer that turns one content-addressed source module plus one pinned-executable `ai-check` run into a MAC-tagged canonical envelope object the landed validator accepts — with every wait, byte, and output stream bounded by the producer's own contract.

**Success metrics:**
1. End-to-end (AC5): real pinned binary + real file-backed store → `GenerateProof` → `ValidateProof` → `Resolve` → `ResolvedGradeProven`, in-repo, no hand-authored bytes anywhere in the chain.
2. Every wait the producer owns is refused-or-bounded by its own derived deadline, proven against the real store on both the read and write sides under `context.Background()` (AC20, AC23), with construction-time ordering pins (AC30).
3. The strict interpretation refuses every real non-proof shape the pinned binary can emit, including the two rc-0 shapes (skipped, contract-free), driven by fixtures generated by the pinned binary itself (AC27).
4. No store surface, no `.ail` file, and no validator behavior changes (V5, V28; conflict surface §9).

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|----------|-----------------|-----------|----------|-------------|
| **No new store write surface.** The `ObjectWriter` seam is an interface satisfied by the unchanged store's `PutObject` + `BusyTimeout()`; round-11b's additive `WriteObject` is not built | Kills gemini's unpinned-write objection without creating a second write path; the original motivation (context-free `PutObject`, parent V51) is gone at HEAD | agent (measured: store.go:543-547, durable.go:10-96, V5/V8/V12) | design | low |
| **`NewProducer` takes `reader`, `writer`, `checker`, `cfg` — reader and writer both first-class** | gpt5-6-sol's objection is an API incompleteness; closing it in the constructor makes the read/write seams injectable and fakeable for the bounded-wait tests | agent (objection verbatim, §0.1) | design | med |
| **`GenerateProof(ctx, sourceRef, requiredIdentities)` keeps the round-9b spelling; required identities are per-call, refused when empty** | The queue row binds this spelling; per-call identities are the caller naming what MUST be verified — a proof of nothing specific is not a proof | agent | design | low |
| **`verify.skipped > 0` is a refusal** (new strictness beyond the inherited text) | Measured on v0.41.0: the skipped shape is the live rc-0 silent class (V25); the report schema has no field that could honestly carry it | agent (fixture-measured, V24/V25) | design | low |
| **The producer byte-verifies and version-probes the pinned executable on EVERY `GenerateProof` call, byte check before first execution** | "before trusting its `ai-check` output" (queue row) means per-run; construction-time-only verification is a TOCTOU — the file can be replaced between construction and use | agent | design | med |
| **Fixtures are committed real-binary output under `host/evidence/testdata/`, re-validated by a fresh-run test** | The directive's fixture rule: never hand-typed; `host/replay/testdata/transition_fixture.ail` is the existing in-`host/` precedent and `verify_ail.sh` does not sweep `host/` (V28) | agent | design | low |
| **Gate change is exactly two literals in `scripts/verify_go.sh`** (`REQUIRED_EVIDENCE_TESTS` + `EXACT_EVIDENCE_TESTS` 37 → new count, same commit as the tests) | The named-manifest gate is the persistent authority; the `.ail` gate is untouched because no `.ail` changes | agent (V22) | design | low |

### Design Freeze

All seven decisions above are made in this document on measured evidence; none requires human input before implementation. Nothing in this design changes a ratified rule, a charter guardrail, or another row's scope. Checked: **all frozen.**

- [x] ObjectWriter-vs-WriteObject (frozen; measured justification V5/V8/V12)
- [x] Constructor parameter set (frozen; objection-driven)
- [x] Per-call required identities + empty-set refusal (frozen)
- [x] skipped>0 refusal (frozen; fixture-measured)
- [x] Per-run executable verification, bytes before execution (frozen)
- [x] Fixture custody: committed real output + fresh-run revalidation (frozen)
- [x] Gate change: two literals, same commit (frozen)

## Solution Design

### Overview

One new production file, `host/evidence/proof_producer.go`, plus tests and pinned-binary-generated fixtures. The producer is a **library-only** sibling of the validator: it holds the same 32-byte MAC key material, reads a source module object through `ObjectReader`, executes the pinned `ailang` binary's `ai-check` subcommand in a bounded process group, strictly interprets the JSON, and stores exactly one object — the MAC-tagged canonical envelope whose `HashRef` it returns. The validator (unchanged) is the only consumer that can turn that envelope into authority.

### Architecture

**Components:**
1. **Constructor** (`NewProducer` + `CheckerConfig` + `ProducerConfig` + the `ObjectWriter` seam): per-field refusals, both ordering pins, key copy (D1).
2. **Executable verifier**: per-run byte hash + bounded `--version` probe, bytes before execution (D2).
3. **Bounded runner**: staging, process group, `procbound` admission/reap, pipe grace, per-stream caps, classification (D3).
4. **Strict interpreter + mint**: the fixed-order JSON interpretation, required identities, anti-vacuity, canonical report, HMAC, envelope bound (D4-D5).
5. **Store legs**: derived-deadline source read with recompute-at-consumption; derived-deadline envelope write with three-outcome classification (D6). What every exceeded bound does is tabulated in D7.

### The API (names binding; field layout implementation detail)

```text
ObjectWriter                  // NEW two-method seam, satisfied by the unchanged *store.Store:
                              //   PutObject(ctx context.Context, o store.Object) error
                              //   BusyTimeout() time.Duration
CheckerConfig                 // the pinned executable and its required identity:
                              //   ExecPath        string          (filesystem path)
                              //   Compiler        hashref.HashRef (byte-hash of the exact executable bytes)
                              //   CompilerVersion string         (exact --version first line, e.g. "AILANG v0.41.0")
ProducerConfig                // the producer's OWN bounds, all required positive:
                              //   ObjectReadTimeout, ObjectWriteTimeout, ExecTimeout time.Duration
                              //   MaxOutputBytes int64          (PER-STREAM cap, applied independently to stdout and stderr)

NewProducer(key [32]byte, reader ObjectReader, writer ObjectWriter, checker CheckerConfig, cfg ProducerConfig) (*Producer, error)

Producer.GenerateProof(ctx context.Context, sourceRef hashref.HashRef, requiredIdentities []string) (hashref.HashRef, error)
```

`ObjectReader` is the validator's existing seam, reused unchanged (validator.go:57-61, V11). The round-9b record's `checker` and `writer` parameters survive as `CheckerConfig` and the `ObjectWriter` seam; the round-12 fix adds `reader`; `execTimeout` and `maxOutputBytes` move into `ProducerConfig` beside their two new timeout siblings. The `GenerateProof` signature is the round-9b spelling verbatim.

### D1 — Constructor, refusals, and the two ordering pins

`NewProducer` refuses, each with `ErrInvalidProducerConfig` naming only its own field (the iteration-106 lesson: the validator's compound constructor condition let three of four disjuncts neuter green — validator.go:64-70 carries the measured comment, V11):

- nil `reader`, nil `writer`;
- `CheckerConfig.ExecPath == ""`, `Checker.IsZero()`, `CheckerConfig.CompilerVersion == ""`;
- non-positive `ObjectReadTimeout`, `ObjectWriteTimeout`, `ExecTimeout`, `MaxOutputBytes` — **this is AC5's constructor arm**, mirroring AC16's validator-side `ObjectReadTimeout` refusal.

Then it reads `window := reader.BusyTimeout()` and `writer.BusyTimeout()` and applies the **two ordering pins** (AC30): a negative window is unknown and refused with `ErrInvalidProducerConfig` (mirroring the validator's unknown-window refusal, validator.go:86-91); a positive window must be strictly below BOTH `ObjectReadTimeout` and `ObjectWriteTimeout`, else the package's existing dedicated `ErrUnorderedTimeouts` (validator.go:30) names both values. The pins read the run-time values from the injected seams, never two literals, so a caller-overridden DSN busy window is pinned at its true value and a wrapper around the real store cannot erase the check (AC30's wrapper arm mirrors the validator's AC22/M30 pattern). The 32-byte MAC key is copied into a fresh heap allocation exactly as the validator copies it — the `keyCopy` shape (validator.go:63-93, V11).

### D2 — Pinned executable byte/version checks (per run, bytes before execution)

Every `GenerateProof` call verifies the executable before trusting any output, in this order:

1. **Byte check FIRST — a tampered binary is never executed.** Read `CheckerConfig.ExecPath`'s bytes and recompute `hashref.Sum(checker.Compiler.Algo(), data)`; refuse unless the digest equals `checker.Compiler.Digest()` — capsule's `verifyExecutable` shape (capsule.go:377-391, V17), at the producer's own call site (recompute-at-consumption: never trust an insert-time or construction-time check). The implementer may choose `os.ReadFile` (capsule's shape; ~100 MB transient) or a streaming hash — both bounded; the property is binding, the mechanism is deferred (§Deferred).
2. **Version probe second, bounded.** Run `<ExecPath> --version` under a context derived from `ctx` with `ExecTimeout` (the same allowance the check run gets), capture stdout **only** (never `2>&1` — the binary emits operational warnings on stderr first, measured 2026-08-17 and recorded in `verify_ail.sh:99-104`, V21), and require the first line to equal `CheckerConfig.CompilerVersion` exactly (`"AILANG v0.41.0"`, V23). Refusal carries `ErrUntrustedTool` naming expected vs got, for both the byte and the version arm (each arm its own named test, AC24).

### D3 — Bounded process-group execution with independent per-stream caps

The `ai-check` run composes the measured shapes this repository already ratified:

- **Staging.** `os.MkdirTemp` scratch root; write the source bytes at `staged/module.ail`; run with `cmd.Dir = root` and the relative module path. `AILANG_RELAX_MODULES=1` — the publish check's own measured rationale (check.go:97-101: without it the verdict depends on whether the interpreter recognises the scratch root as a temp dir, iter-195, V19). A private `AILANG_CACHE_DIR` appended last so an ambient value cannot win (check.go:102, V19). The child environment is `childenv.Scrubbed(os.Environ())` plus those two assignments — the registry-credential scrub (childenv.go:34-48, V20).
- **`TMPDIR` is redirected into the scratch root.** Measured (V26): the pinned binary writes Z3 solver scratch files under `TMPDIR`; with the variable stripped the solver fails with `solver error: creating temp file: … operation not permitted`; with `TMPDIR` pointed at a scratch subdir the run succeeds and the directory is empty afterward. Redirecting confines all child temp state under the removed-on-cleanup root.
- **Process group.** `cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}`; `cmd.Cancel` and the overflow kill both go through ONE package-level `killGroup` seam (`killGroup = func(pgid int) error { return syscall.Kill(-pgid, syscall.SIGKILL) }`), mirroring capsule's D1 shape (capsule.go:83-97, V17) so a forked grandchild holding the pipes cannot outlive the call (row 24's ratified boundary; rows 112/113 own the archive/replay/pkgproj sites, not this one).
- **Bounds.** `execCtx, cancel := context.WithTimeout(ctx, cfg.ExecTimeout)` — the earlier of the caller's deadline and the producer's own allowance. The Z3 per-query flag is the gate's own `ai-check -timeout 5s` spelling (verify_ail.sh:307, V21); `ExecTimeout` bounds the WHOLE run (the reason `GATE_LEG_TIMEOUT_S=120` exists, verify_ail.sh:160-162, V21). `procbound.Admit()` before `Start` (release on Start failure); the pipe closer at `ExecTimeout + pipeCloseGrace` closes both read ends so the drains return; the direct-child `Wait` is bounded through `procbound.Wait` at `ExecTimeout + 2·pipeCloseGrace`; a kill error is joined only when `!procbound.NothingToKill(...)` (the ESRCH/zombie filter, procbound.go:70-88, V18). `pipeCloseGrace = time.Second`, the same value and comment lineage as capsule/broker (15 existing occurrences of the token across those packages — the row-24 precedent, V7).
- **Independent per-stream limit+1 caps.** `stdout` and `stderr` are drained concurrently, each through its own `readCapped(pipe, cfg.MaxOutputBytes)` — `io.ReadAll(io.LimitReader(pipe, limit+1))`, overflow decided when `len(data) > limit`, exactly capsule's spelling (capsule.go:392-403, V17). **`MaxOutputBytes` is a PER-STREAM cap applied independently to each stream**: a single shared cap re-opens the round-8 dual-stream objection the parent closed (§10.11 item (iii)). Overflow on either stream triggers the group kill via a `sync.Once` and refuses the whole proof with `ErrOutputLimit` — the truncated prefix is NEVER parsed.
- **Classification.** Overflow outranks the deadline (the kill is what let the wait return); deadline expiry is `ErrExecTimeout`; both join `ErrCleanupIncomplete` behind the typed error when cleanup outlives the grace; a clean exit proceeds to interpretation. Declared residual, inherited unchanged from row 24: a setsid escapee is disowned, not killed; a direct child surviving a failed kill is reaped later by `procbound`'s background waiter.

### D4 — Strict `ai-check` JSON interpretation, fixture-driven

The interpretation follows `verify_ail.sh`'s JSON-parse-not-exit-code discipline (the gate's own header names the two silent classes; V21) and sources the verified set exactly as the parent specified: **`verify.results[].function` filtered on `status == "verified"`** (parent §10.11: "the verified set remains sourced from `verify.results[].function` filtered on `status == "verified"`"). Fixed order:

1. stdout empty → `ErrVacuousCheck` (no output at all);
2. JSON unparseable → `ErrUnverifiedRun` wrapping the parse error (never an empty-report-plus-success);
3. `check.passed != true` → `ErrUnverifiedRun` (fixture `loaderror`);
4. `verify.available != true` → `ErrUnverifiedRun` (same fixture);
5. `verify.errors > 0` → `ErrUnverifiedRun` (fixtures `encerror`, `mixed` — the mixed fixture proves why: one verified record does not launder an errors>0 run);
6. `verify.counterexample > 0` → `ErrUnverifiedRun` (fixture `refuted`);
7. `verify.skipped > 0` → `ErrUnverifiedRun` (fixture `skippedshape` — **the measured rc-0 silent class on v0.41.0**, V25; the report schema has no honest field for a skipped contract);
8. extract `verified = [r.function for r in results if r.status == "verified"]`; **there is no acceptance dependence on the `verify.verified` counter** — the counters are refusal triggers only (a deliberate cross-check branch on counter-vs-records equality is NOT added: the real binary is self-consistent, so no real fixture could ever pin it — a guard nothing can pin is not a guard, the mission's own measured lesson, types.go:60-68 comment, V13);
9. **anti-vacuity guard**: the extracted `verified` set must be non-empty, else `ErrVacuousCheck`, LOUD — ONE guard, deliberately not split into a zero-records branch plus a zero-verified branch, because two guards with one sentinel on one fixture shape make the first unkillable (the defect this section's own mutation drill caught at design time). The refusal's detail distinguishes its two input shapes: `0 results records from non-empty output` (fixture `nocontracts`: rc 0, every counter zero, `results: []` — a caller keyed on exit codes or counters alone would accept a run that proved nothing, V27) from `0 verified among N results records`. Never an empty report plus success;
10. sort the extracted set; every `requiredIdentities` entry must be present (binary search via the package's existing `containsIdentity`, validator.go:185-196, V11), else `ErrRequiredIdentityMissing` naming the missing identity; empty/nil `requiredIdentities` is itself refused (a proof of nothing specific is not a proof — mirrors the validator's AC21 non-vacuity rule on the minting side).

Unknown JSON fields (e.g. the `duration` and `model` members the fixtures carry) are tolerated: the gate's own parser is field-targeted the same way (verify_ail.sh:355-395, V21), and `ai-check`'s JSON grew members between v0.30.0 and v0.41.0 — refusing unknown fields would couple the producer to one release's shape. The security work is done by the semantic refusals above, not by shape absolutism.

### D5 — MAC-tagged canonical envelope mint

On the only-clean path: build `ProofReportV1{Schema: ProofReportSchemaV1, Subject: sourceRef, Compiler: checker.Compiler, CompilerVersion: checker.CompilerVersion, Verified: sortedVerified, Errors: 0, Counterexamples: 0, CheckPassed: true, ProofSucceeded: true}` — the nine-member canonical report the landed codec already enforces (sorted-unique verified, report_codec.go:86-97, V14). Encode canonically (`EncodeProofReportV1`), compute `HMAC-SHA-256(key, reportBytes)` — the tag covers **exactly the decoded canonical report byte sequence**, matching the validator's recomputation (validator.go:156-159, V11) — and encode the two-member envelope (`EncodeAuthenticatedEnvelope`; members `report` then `mac`, base64url-no-padding, envelope_codec.go:9-33, V15). The envelope byte length must be **≤ `MaxBytes` (256 KiB, types.go:13, V13)** or the mint refuses with `ErrEnvelopeTooLarge`: the validator reads envelopes through `ReadObject(ctx, ref, MaxBytes)` and its pre-materialization probe would refuse anything larger — the producer never writes an envelope its only consumer cannot read. `envelopeRef = hashref.SumSHA256(envelope)` — sha256 specifically, because the validator recomputes `SumSHA256(payload)` and compares to the requested ref (validator.go:137, V11).

### D6 — Source read and envelope storage under the weaker store contract

**The store contract this designs against** (row 26's gate text points at `w-store-deadline-free-residue-owner.md` §3): *the store inherits caller cancellation; zero is a syntax result, not a deadline guarantee.* The row-23 policy tranche then landed the strict guard on top (`D-WORLD-37` = A): every store I/O method refuses nil/deadline-free contexts with `ErrNoDeadline` (store.go:236-247, V8). Neither fact supplies a bound: inheritance gives the producer whatever lifetime the caller chose, and the guard only rejects the deadline-FREE case. **So the producer derives its own bounds at its own call sites:**

- **Source read.** Validate `sourceRef` (canonical `hashref.Parse`); derive `readCtx := context.WithTimeout(ctx, cfg.ObjectReadTimeout)` even under `context.Background()`; read through `reader.ReadObject(readCtx, sourceRef, MaxSourceBytes)` — a NEW named constant, `1 MiB`: the largest `.ail` in the repository today is 25,422 bytes (`packages/se-tools/se_tools/run.ail`, V29), so 1 MiB is ~40× headroom while still refusing a multi-gigabyte object before materialization (the store's length-probe guard, read_object.go:43-83, V11/V8). Over-bound → the store's typed `*store.ObjectTooLargeError` → refusal. Recompute `hashref.Sum(sourceRef.Algo(), sourceBytes)` and require digest equality — recompute-at-consumption, the parent's rule (§3.3 step 4's rationale): the store checked the hash at insert, and the producer still verifies at use, so alternate readers, corruption, and test fakes cannot bypass the integrity boundary. `Subject == sourceRef` follows.
- **Envelope write.** Derive `writeCtx := context.WithTimeout(ctx, cfg.ObjectWriteTimeout)`; write exactly one object through `writer.PutObject(writeCtx, store.Object{Hash: envelopeRef, InterfaceHash: InterfaceHashV1, SemanticID: ProofSemanticID, Provenance: "bounded-ai-check-producer", Payload: envelope})` — the two metadata identities the validator's steps 4-5 require (validator.go:140/:143, V11/V16). `PutObject` is idempotent (`INSERT OR IGNORE`, the SQL at store.go:566, the contract comment at :540-541, V8) and carries the ratified three-outcome contract: a definite not-committed error is an operational failure (the caller may re-run `GenerateProof` — identical inputs mint byte-identical envelope bytes, and a present object is a no-op insert); a `*store.UncertainError` is surfaced **wrapped and un-retried by the producer** — reconciliation is the caller's job, never a blind retry (`D-WORLD-40` = A, durable.go:10-45, V12).

### D7 — What happens when each bound is exceeded (refusal, never partial trust)

| Bound exceeded | Observable | State after |
|---|---|---|
| `ObjectReadTimeout` (source read) | `readCtx` expires; the store returns the ctx-classified error; producer wraps it; zero `HashRef` + error | No mint, no write, no scratch left (`RemoveAll` deferred) |
| `MaxSourceBytes` | The store's pre-materialization probe returns `*ObjectTooLargeError` | Same |
| Byte check | `ErrUntrustedTool` (expected vs got digest) — the binary is never executed | Same |
| Version probe | `ErrUntrustedTool` (expected vs got first line) | Same |
| `ExecTimeout` | `execCtx` fires; `cmd.Cancel` SIGKILLs the whole group; `ErrExecTimeout`, `ErrCleanupIncomplete` joined if the reap outlived the grace | Same |
| `MaxOutputBytes` (either stream) | Group kill at the overflow; `ErrOutputLimit`; the truncated prefix is never parsed | Same |
| Z3 per-query `-timeout 5s` | Whatever the tool emits, every non-`verified` outcome lands in D4's refusals (fixtures prove the shapes) | Same |
| Run unclean (failed/unavailable/errors/counterexample/skipped) | The D4 sentinel for the branch | Same |
| Zero verified records (incl. zero `results` from non-empty output) | `ErrVacuousCheck` — LOUD, detail naming the shape | Same |
| Required identity missing | `ErrRequiredIdentityMissing` naming it | Same |
| Envelope > `MaxBytes` | `ErrEnvelopeTooLarge` — the mint refuses before any write | Same |
| `ObjectWriteTimeout` (envelope write) | `writeCtx` expires; `PutObject` classifies: definite not-committed → operational error; in-flight commit → `*store.UncertainError` (reconcile, never blind retry) | No envelope returned either way; the object may exist under an uncertain outcome, and re-running mints the identical bytes idempotently |

No failure path writes an object and no success path writes two.

### Implementation Plan (milestones; each ≤ 160 production LOC)

- **M1 — constructor, config, ordering pins (~85 prod LOC; 0.25 d).** `proof_producer.go`: `ObjectWriter`, `CheckerConfig`, `ProducerConfig`, `NewProducer` with all per-field refusals and both ordering pins. Tests: nine per-field constructor refusal tests (four non-positive bounds, nil reader, nil writer, three incomplete-checker fields) + the ordering-pin test with its wrapper arm (mirroring `TestConstructorPinsBusyTimeoutBelowObjectReadTimeout` and `TestReaderWaitBoundsCannotBeLostThroughWrapper`). Exit gate: `go test ./host/evidence -count=1` green; the manifest update lands with M5, not before (see Files to Modify/Create and §9).
- **M2 — executable verification (~65 prod LOC; 0.3 d).** Byte-hash + version probe, both derived-deadline-bounded, bytes before execution. Tests: byte-mismatch (fake `Compiler` ref, real binary on disk), version-mismatch (wrong `CompilerVersion` string), absent path.
- **M3 — bounded process-group execution (~160 prod LOC; 0.5 d).** Staging, env composition, `Setpgid`/`killGroup`/`procbound`/pipe-grace rig, per-stream `readCapped`, classification. Tests: tiny-cap stdout overflow (a real `ai-check` run of the committed verifiedpass module under `MaxOutputBytes=16`), tiny-cap stderr overflow arm, exec-timeout arm, kill-seam observation (wrap `killGroup`, assert count).
- **M4 — strict interpretation + mint (~120 prod LOC; 0.35 d).** The D4 parser + D5 mint + required identities + anti-vacuity + envelope bound. Tests: all seven fixture shapes with their expected sentinels; identity arms; the in-process maximal-report oversize-envelope arm; the fresh-run revalidation test (regenerates fixtures with `AILANG_BIN`, fatal-never-skip, and asserts the parser's verdicts on FRESH output match the committed fixtures').
- **M5 — store integration + gates (~70 prod LOC; 0.55 d).** Derived read/write deadlines, recompute-at-consumption, three-outcome classification. Tests: AC5 end-to-end (real binary + real store → validator → `ResolvedGradeProven`), AC20 write-side pool decoy, AC23 read-side pool decoy. Gate: update `scripts/verify_go.sh` (two literals, 37 → observed count, same commit); full mutation re-drill M15/M28/M31-M39; `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` and `./scripts/verify_go.sh` both green.

### Files to Modify/Create

**New files:**
- `host/evidence/proof_producer.go` — the producer (~430 prod LOC total across M1-M5)
- `host/evidence/producer_test.go` — fixture-driven + constructor + verification tests
- `host/evidence/producer_realstore_test.go` — AC5/AC20/AC23 real-store integration tests (mirrors `realstore_test.go` helpers)
- `host/evidence/testdata/aicheck/*.ail` + `*.json` — the seven fixture modules and their real pinned-binary outputs (committed; `host/replay/testdata/transition_fixture.ail` is the in-`host/` .ail precedent, V28; `verify_ail.sh` sweeps only `design_docs`, `world/`, `packages/se-tools`, V21/V28)
- `host/evidence/testdata/aicheck/regen.sh` — the regeneration script (this document's `design_docs/verification/world-iter251/gen_aicheck_fixtures.sh`, adapted to write into the testdata dir)

**Modified files:**
- `scripts/verify_go.sh` — exactly two literals: the banner and the `check_evidence_manifest ... 37` argument (37 → observed count), in the SAME commit as the new tests (V22)

**Unchanged (with receipts):** `host/store/**` (no new surface, V5/V8), `host/evidence/validator.go` + codecs (reused, V11/V14/V15), `world/*.ail` and every `.ail` (V28), `host/procbound`, `host/proctest`, `host/childenv` (reused, V18/V20), `.github/workflows/ci.yml` (both gates already run; the go job already exports the pinned `AILANG_BIN`, V32).

## Examples

### Example 1 — end-to-end mint and validate (the AC5 demonstration)

**Before (today):** the only way to obtain an envelope is to hand-author it in a test (`authority_test.go`'s `envelopeFor`); there is no in-repo path from an `ai-check` run to an authenticated envelope (parent §10.16; verified live: no producer surface V4, no Go `ai-check` runner V6).

**After (illustrative; API names binding, snippets not checkable-in):**

```go
// compose exactly as the validator composes — one store, one key, one pinned binary
s, _ := store.Open(dbPath)                       // unchanged store
reader, writer := s, s                           // *Store satisfies both seams unchanged

checker := evidence.CheckerConfig{
    ExecPath:        os.Getenv("AILANG_BIN"),    // the pinned v0.41.0 release
    Compiler:        compilerRef,                // sha256 of the exact executable bytes
    CompilerVersion: "AILANG v0.41.0",            // the --version first line
}
cfg := evidence.ProducerConfig{
    ObjectReadTimeout:  3 * time.Second,          // must exceed BusyTimeout() when positive
    ObjectWriteTimeout: 3 * time.Second,
    ExecTimeout:        2 * time.Minute,
    MaxOutputBytes:     1 << 20,                 // per stream
}
p, err := evidence.NewProducer(key, reader, writer, checker, cfg)

envelopeRef, err := p.GenerateProof(context.Background(), sourceRef,
    []string{"applyRevision", "plan", "verify"})  // the caller names what MUST be verified
// → envelopeRef addresses the MAC-tagged canonical envelope object in the store

v, _ := evidence.NewValidator(key, s, evidence.CompilerConfig{
    Compiler: compilerRef, CompilerVersion: "AILANG v0.41.0", ObjectReadTimeout: 3 * time.Second,
}, []string{"applyRevision"})
r := v.ValidateProof(ctx, envelopeRef, sourceRef) // unchanged validator
sealed, ok := r.Validated()                      // the ONLY mint of authority
grade, _ := v.Resolve(sealed).Proven()            // ResolvedGradeProven
```

### Example 2 — what the strict interpretation refuses (real fixtures, V24)

| committed fixture (real pinned-binary output) | rc | every counter | statuses | producer verdict |
|---|---|---|---|---|
| `verifiedpass.json` | 0 | clean | `[verified]` | mint (with the right required identities) |
| `refuted.json` | 1 | cex=1 | `[counterexample]` | refuse `ErrUnverifiedRun` |
| `encerror.json` | 1 | errors=1 | `[error]` | refuse `ErrUnverifiedRun` |
| `mixed.json` | 1 | errors=1 | `[verified, error]` | refuse `ErrUnverifiedRun` — one verified record does not launder the run |
| `skippedshape.json` | **0** | skipped=1 | `[skipped]` | refuse `ErrUnverifiedRun` — **the rc-0 silent class** |
| `loaderror.json` | 1 | passed=false | `[]` | refuse `ErrUnverifiedRun` |
| `nocontracts.json` | **0** | all zero | `[]` | refuse `ErrVacuousCheck` — the clean-but-vacuous shape |

## Success Criteria

See §7 (acceptance criteria, AC5/AC20 inherited + AC23-AC30) and §8 (mutations, M15/M28 inherited + M31-M39). Summary gates: both pinned verify gates green at HEAD; the full mutation matrix killed with named tests; the 37 existing evidence tests still green (manifest extended, never truncated).

## Testing Strategy

**Unit/fixture tests (M1-M4):** per-field constructor refusals; byte/version verification arms; per-stream overflow arms with a real `ai-check` run under a tiny cap; the seven-fixture interpretation matrix; the anti-vacuity arm; required-identity arms; the oversize-envelope arm (in-process maximal report through the real codec); the fresh-run revalidation test (regenerates all fixtures with `AILANG_BIN` and asserts verdicts on fresh output — fixtures are never trusted as frozen truth, and a future pin move that changes the JSON shape reds here rather than silently passing).

**Integration tests (M5, real store + real binary, no fakes on the bounded-wait arms):** AC5 end-to-end; AC20 write-side pool decoy and AC23 read-side pool decoy, both mirroring the landed `TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout` stimulus discipline (realstore_test.go:144-210: measured decoy hold, `decoyHoldRatio=20`, race-honest retry ≤ 5, 2× mutant ceiling, 20× watchdog, announced `GOMAXPROCS(2)` precondition — V17/V11).

**Manual/CI:** no manual steps; CI's go job exports the pinned `AILANG_BIN` (V32) and producer tests `t.Fatal` on its absence, never `t.Skip` — a skipped manifest test reds the named-manifest gate by design (`skipped` is a refusal class, verify_go.sh:100-103, V22).

**Binary-dependent arms and the manifest:** every named test is top-level (no subtests) so the manifest's terminal-`Action=pass` counting sees it; binary-dependent arms fatal (never skip) per `check_test.go:79-81`'s discipline (V32).

## Deferred Decisions

- `os.ReadFile` vs streaming SHA-256 for the byte check — agent may choose (both bounded; ReadFile mirrors capsule, streaming avoids the ~100 MB transient).
- The staged module filename under the scratch root — agent may choose (`AILANG_RELAX_MODULES=1` makes it inert; any path works, V19).
- The `Provenance` label text — agent may choose (validator ignores it; tests must not depend on it).
- Test-file split (`producer_test.go` vs more files) — agent may choose; the test NAMES are binding through the manifest.
- Whether the version probe shares the `checkBuffer` pattern or uses bounded pipes — agent may choose, provided both streams stay capped and `2>&1` is never used.

## Non-Goals

- **No new store surface** — no `WriteObject`; `PutObject`/`ReadObject`/`BusyTimeout` are consumed unchanged (gemini's objection is answered by reuse + pins, not by a parallel path).
- **No validator change** — the producer must not touch `ValidatedEvidence`, `ValidateProof`, `Resolve`, or any codec; the validator stays the sole mint (parent §3.5; conflict surface §9).
- **No kernel/AILANG change** — no `.ail` file moves (V28); the receipt stays an untrusted `CLAIMED` claim.
- **No production composition root** — key custody and daemon wiring remain queue row 3 of item 17's decomposition (`w-proven-evidence-production-key-wiring`); the producer is library-only, exactly like tranche 1.
- **No lock-wait deadline enforcement** — row 111 owns lock-contended deadline dominance and the 500-vs-503 classification; this design pins the ordering and declares the residual (§0.2 item 3).
- **No re-pricing or absorption of rows 112/113** — archive/replay/pkgproj subprocess sites are theirs; this design only adds a NEW launch site that follows row 24's ratified shape.
- **No replay/renderer integration** — tranches 3/4 of item 17's decomposition.

## Timeline

| Milestone | Estimate | Production LOC |
|---|---:|---:|
| M1 constructor + pins | 0.25 d | ~85 |
| M2 executable verification | 0.30 d | ~65 |
| M3 bounded process-group execution | 0.50 d | ~160 |
| M4 strict interpretation + mint | 0.35 d | ~120 |
| M5 store integration + gates + mutation re-drill | 0.55 d | ~70 |
| **Total** | **~1.95 d ≈ 2.0 d** | **~500** |

The queue row carries 1.0 d from item 17's measured cut and explicitly leaves the round-12 fixes "unpriced, pending design". This design prices them: the reader parameter and derived read deadline (AC23 + M31), the write deadline and both ordering pins (AC30 + M32), the executable verification rig, the fixture infrastructure, and the decoy tests add ~1.0 d over the carried figure. That is stated, not hidden.

## Risks & Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| `ai-check` JSON shape drifts across releases | Med | Unknown fields tolerated; the fresh-run revalidation test reds on verdict drift when the pin moves; fixtures regenerated in the same commit as any pin bump |
| Wall-clock decoy tests flake | Med | Copy AC16's landed stimulus discipline verbatim (measured hold, ratio, retries, watchdog, GOMAXPROCS precondition) — that test is green in CI today |
| The pinned binary is absent in a local lane | Low | Binary-dependent arms fatal, never skip (check_test.go precedent, V32); CI exports the pin |
| A 100 MB ReadFile per proof is heavy | Low | Deferred to implementer (streaming); the cost is visible and bounded either way |
| The producer accidentally becomes an authority mint | High (design-level) | It cannot: it returns a `HashRef` to DATA; only `ValidateProof` mints `ValidatedEvidence`, and `TestPublicAuthoritySurfaceIsFrozen` still guards the boundary (§9) |
| Scope creep into the store | Med | The no-new-surface decision is frozen; V5 is the receipt |

## Related Documents

- `design_docs/implemented/w-validated-proven-evidence-boundary.md` — the parent: §3.4 (shed marker), §3.5 (boundary), §10.11 (round-9b API record), §10.12 (the load-bearing decoy distinction), §10.15-10.16 (the two objections, the shed ruling)
- `design_docs/implemented/w-store-deadline-free-residue-owner.md` — the weaker store contract (§3) and its explicit instruction to this row (§7: "use the explicit weaker store contract; do not add a twelfth hidden root")
- `design_docs/planned/w-store-bounded-durable-operations.md` — row 23's policy tranche design (the strict guard + three-outcome commit contract this design consumes; its M5a–M7b tranche landed at iter-207, PR #162, while the design document itself remains in planned/)
- `design_docs/implemented/w-host-subprocess-cleanup-boundary.md` — row 24: the process-group/caps/cleanup-boundary shape D3 reuses; its §9 names this row as the next consumer
- `design_docs/implemented/w-daemon-lock-wait-not-deadline-bound.md` — row 22: the composition-claim pattern AC30 mirrors
- Item 17's §9 decomposition names the successor tranche `w-proven-evidence-production-key-wiring` (production key custody; not yet filed as a document — named row, no path). This row stays library-only until it lands

## References

- Charter queue row 26: `design_docs/world-mission.md:4451` (verbatim scope, inherited defects, ordering gate, pricing)
- Rulings: `D-WORLD-24` arm A (the shed), `D-WORLD-69` = A (exception path), `D-WORLD-37` = A + `D-WORLD-40` = A (the store contract consumed here)
- `design_docs/coding-standards.md` — Z3 contracts on the pure core, effects at the host boundary
- Evidence artifacts: `design_docs/verification/world-iter251/` — the fixture generator, the seven committed fixtures, this session's receipts

## Future Work

- Production wiring: who calls `GenerateProof` with a durable key — item 17's tranche-2 row (`w-proven-evidence-production-key-wiring`).
- Renderer/replay consumption — item 17's tranches 4/3.
- If the `ai-check` JSON ever carries a per-identity `duration`-bound the report should attest, that is a new report schema version (`world/proof-report/v2`), not an edit to v1.

---

## §7 Acceptance criteria

AC5 and AC20 are inherited from item 17's removed set with their original identifiers and no renumbering (parent §7 preamble; charter row 26). New criteria number after the parent's last (AC22).

**AC5 — producer constructor refusal of non-positive bounds + end-to-end envelope production (inherited).** `NewProducer` refuses every non-positive bound — `ObjectReadTimeout`, `ObjectWriteTimeout`, `ExecTimeout`, `MaxOutputBytes` — each with `ErrInvalidProducerConfig` naming only its own field, mirroring AC16's validator-side `ObjectReadTimeout` refusal so the bounded-waits axiom is checkable at the API rather than asserted; it likewise refuses nil reader, nil writer, and incomplete `CheckerConfig` (empty path, zero compiler ref, empty version), one observable refusal per field (the iteration-106 compound-condition lesson). The end-to-end arm: against the pinned v0.41.0 binary and a REAL file-backed store, `GenerateProof` over a source object containing a module with a verified contract returns an envelope `HashRef` that a `NewValidator`-constructed validator accepts — `ValidateProof` yields `ValidatedEvidence`, `Resolve` yields `ResolvedGradeProven` — the end-to-end demonstration the shed removed from tranche 1, restored in-repo with no hand-authored bytes anywhere in the chain. Named tests: `TestProducerConstructorRefusesNonPositiveObjectReadTimeout`, `...ObjectWriteTimeout`, `...ExecTimeout`, `...MaxOutputBytes`, `...NilReader`, `...NilWriter`, `...EmptyExecPath`, `...UnsetCompilerIdentity`, `...EmptyCompilerVersion`, and `TestProducerEndToEndMintsEnvelopeValidatorAccepts`. M15 must red the required-identity arm of this criterion's mint path; M35/M36 must red its verification arms; M28 must red its write path.

**AC20 — write-side wait bound proven against the REAL store by a connection-pool decoy, which exercises the POOL wait and NOT the LOCK wait (inherited; scope stated in the criterion itself).** `GenerateProof` derives `writeCtx := context.WithTimeout(ctx, ObjectWriteTimeout)` before the envelope write, even under `context.Background()`. Named integration test `TestRealStoreBlockedEnvelopeWriteReturnsWithinObjectWriteTimeout` (real `host/store`, no fake): the stimulus is one in-flight `GetObject` of a decoy object sized at implementation so a single read outlasts 20× `ObjectWriteTimeout` (hold ASSERTED in-test, mirroring `minDecoyHold`/`decoyHoldRatio`); **the decoy holds the store's single pooled connection (`db.SetMaxOpenConns(1)`, store.go:382) and holds NO SQLite write lock** — the criterion exercises the connection-acquisition wait, which `database/sql` bounds by the supplied context by stdlib contract (`beginDurable`, durable.go:48-51), and **claims nothing about lock contention** (that wait is bounded by `busy_timeout`, is not changed here, and its deadline dominance is OPEN row 111's subject). The blocked `GenerateProof` must return an operational error within the derived bound and write NO object; the no-decoy control (AC5's end-to-end) runs first and must succeed. Classification is race-honest as AC16: a write completing below the bound means the freed connection was won — retry, bounded at 5, exhaustion a loud instrument failure; completion at ≥ 2× the bound is the mutant signature and reds at once; a 20× watchdog makes the hang mode a red. Readings scoped to CODE: `ObjectWriteTimeout` reads 0 in `*.go` today (V7) and ≥ 1 in non-test `host/evidence` at head, by construction. M28 must red the write-error arm; M32 must red the derivation arm.

**AC23 — source-read wait bound, the read-side mirror of AC20 (closes gpt5-6-sol's derived-deadline half for the read).** `GenerateProof` derives `readCtx := context.WithTimeout(ctx, ObjectReadTimeout)` before reading `sourceRef` through the constructor's `ObjectReader`, even under `context.Background()`. Named test `TestRealStoreBlockedSourceReadReturnsWithinObjectReadTimeout`: same real-store pool-decoy stimulus as AC20, applied to the read; the blocked read must return the operational error within the bound and mint nothing; same retry/2×/watchdog classification; the AC5 end-to-end doubles as the no-decoy control. The test asserts the error chain carries the deadline expiry (`context.DeadlineExceeded` wrapped) and NOT `store.ErrNoDeadline` — so a mutant that drops the derivation dies on error identity here, not on the store's guard. M31 must red it.

**AC24 — pinned executable byte/version checks, per run, bytes before execution.** Every `GenerateProof` recomputes the executable's content hash from `CheckerConfig.ExecPath`'s exact bytes under `Checker.Compiler`'s own algorithm and requires digest equality BEFORE executing the binary, then runs a bounded `--version` probe (stdout only, never merged with stderr) and requires the first line to equal `CheckerConfig.CompilerVersion` exactly. A byte mismatch, a version mismatch, and an unreadable path are three separately observable `ErrUntrustedTool` refusals, each naming expected vs got; under each refusal no subprocess runs the binary (byte arm) and no report is minted. Named tests: `TestProducerRefutesTamperedExecutableBytes` (a deliberately wrong `Compiler` ref against the real binary on disk), `TestProducerRefutesWrongCompilerVersion`, `TestProducerRefusesMissingExecutable`. M35 and M36 must red the first two respectively.

**AC25 — bounded process-group execution.** The `ai-check` child runs with `Setpgid: true`; `cmd.Cancel` and the overflow kill both go through one package-level `killGroup` seam and target the whole group; `procbound.Admit()` is called before `Start` (refusing with `ErrCleanupBacklog` at the process-wide cap) and the reservation is released on Start failure or handed to `procbound.Wait`; the drains are unblocked by a pipe closer at `ExecTimeout + pipeCloseGrace` and the direct-child wait is bounded at `ExecTimeout + 2·pipeCloseGrace`; a failed kill is joined behind the typed error only when `!procbound.NothingToKill(...)`; `ErrCleanupIncomplete` joins behind the typed refusal when the reap outlived the grace. Declared residuals, unchanged from row 24: a setsid escapee is disowned, not killed; a child surviving a failed kill is reaped by the background waiter. Named tests: `TestProducerExecTimeoutKillsTheWholeProcessGroup` (the kill observed through the wrapped seam) and `TestProducerAdmissionRefusesBacklog`. M39 must red the exec-timeout arm.

**AC26 — independent stdout/stderr limit+1 caps.** `MaxOutputBytes` is a PER-STREAM cap applied independently: stdout and stderr are each drained through `readCapped(pipe, MaxOutputBytes)` (limit+1 read; overflow decided when `len(data) > limit`); overflow on EITHER stream triggers the group kill once and refuses the proof with `ErrOutputLimit`; the truncated prefix is never parsed; a stream that stays under the cap while its sibling overflows must not change the verdict (independence). Named tests: `TestProducerRefusesStdoutOverflowIndependently` and `TestProducerRefusesStderrOverflowIndependently` — a REAL `ai-check` run of the committed verifiedpass module under a deliberately tiny `MaxOutputBytes` (16 bytes), each arm asserting the typed refusal and no mint. M33 must red both arms.

**AC27 — strict `ai-check` JSON interpretation, driven by fixtures generated by the pinned binary.** The verified set is sourced from `verify.results[].function` filtered on `status == "verified"` — never from the `verify.verified` counter alone — and the producer refuses, each with its own sentinel: unparseable JSON; `check.passed != true`; `verify.available != true`; `verify.errors > 0`; `verify.counterexample > 0`; `verify.skipped > 0` (the measured rc-0 silent class on v0.41.0); empty stdout; and **the anti-vacuity guard: an empty extracted verified set — including the zero-`results` shape — fails LOUDLY with `ErrVacuousCheck`, never an empty report plus success** (the committed `nocontracts` fixture is exactly this shape: rc 0, all counters zero, `results: []`). Exit codes are advisory — the fixtures prove both rc-0 non-proof shapes (`skippedshape`, `nocontracts`) and both rc-1 shapes (`refuted`, `encerror`, `mixed`, `loaderror`). Every fixture in the test plan is REAL pinned-binary output committed under `host/evidence/testdata/aicheck/` with its generating module and a regeneration script; the named test `TestProducerInterpretationMatchesEveryFixtureShape` runs the parser over each fixture and asserts the exact verdict; `TestFixturesRevalidateAgainstFreshBinaryOutput` regenerates with `AILANG_BIN` (fatal, never skip) and asserts the fresh output yields the same verdicts. M34 must red the anti-vacuity arm; M38 must red the errors>0 arm.

**AC28 — required verified identities are per-call and non-vacuous.** `GenerateProof(ctx, sourceRef, requiredIdentities)` refuses nil/empty `requiredIdentities` (a proof of nothing specific is not a proof); every named identity must appear in the extracted verified set, checked by binary search over the sorted set; a missing identity refuses with `ErrRequiredIdentityMissing` naming it, and nothing is minted. Named test `TestProducerRequiresCallerNamedIdentities`: the verifiedpass fixture with the true identity mints (control), with an absent identity refuses, and the empty set refuses. M15 must red the containment arm.

**AC29 — MAC-tagged canonical envelope production and bounded envelope storage.** The report is encoded by `EncodeProofReportV1` (canonical nine-member JSON, sorted-unique verified), the tag is HMAC-SHA-256 over the exact canonical report bytes with the constructor-copied 32-byte key, and the envelope is `EncodeAuthenticatedEnvelope`'s two-member canonical bytes. The envelope length must be ≤ `MaxBytes` or the mint refuses with `ErrEnvelopeTooLarge` BEFORE any write — the producer never stores an envelope the validator's own bounded read cannot retrieve. Exactly one object is written, through `ObjectWriter.PutObject` under the derived `ObjectWriteTimeout`, carrying `InterfaceHashV1` and `ProofSemanticID`; the returned ref is `SumSHA256(envelope)`. Write outcomes follow the ratified three-outcome contract: definite not-committed is an operational error (re-running mints identical bytes idempotently); `*store.UncertainError` is surfaced wrapped and never blindly retried. On any failure the zero `HashRef` is returned and nothing is written. Named tests: `TestProducerMintsValidatorReadableEnvelope` (round-trip through the real validator), `TestProducerRefusesEnvelopeBeyondValidatorReadCap` (in-process maximal report: 256 identities at the codec's 1 KiB string cap, constructed through the real codec — the codec's own bounds are the fixture), `TestProducerSurfacesUncertainWriteOutcome` (via the store's existing `durableAfterCommitHook` test seam, never a fake store). M37 must red the size-bound arm; M28 must red the write-error arm.

**AC30 — construction-time ordering pins, both sides, surviving wrapping.** `NewProducer` reads `BusyTimeout()` from BOTH the reader and the writer and refuses: a negative window (unknown) with `ErrInvalidProducerConfig`, and a positive window ≥ `ObjectReadTimeout` OR ≥ `ObjectWriteTimeout` with `ErrUnorderedTimeouts` naming both values — the pins bind the values actually used at run time, never two literals. Named test `TestProducerConstructorPinsBusyTimeoutBelowBothTimeouts` (real file-backed store): the ordered arm constructs and produces end-to-end; each unordered arm refuses at construction; a wrapper-around-real-store arm exposes the wrapped store's true cached window and cannot erase the pin (mirroring the validator's AC22/M30 pattern). The composition claim it licenses is stated in §0.2 item 3 and §D7: safe while `busy_timeout` is below both producer timeouts; lock-wait deadline dominance stays row 111's. M31/M32 are the derivation mutants this pin makes meaningful.

## §8 Mutation table

M15 and M28 are inherited with their original identifiers (parent §6 preamble; charter row 26). New rows anchor to the diff this document designs and number after the parent's last (M30). Every mutant is a compiling, behaviour-changing edit (`if false && …` or an equivalent in-value widening); a build failure is not a kill. Every observable is returned by or after the mutated mechanism, never a sibling value assigned beside it.

| ID / class | Exact file and neutered edit | Named check that fires | Downstream observable and predicted failure text |
|---|---|---|---|
| **M15 required-identity filter (inherited)** | `host/evidence/proof_producer.go`: the required-identity containment `if !containsIdentity(verified, required)` → `if false && !containsIdentity(verified, required)` | `TestProducerRequiresCallerNamedIdentities` (AC28) | `GenerateProof` mints an envelope over the verifiedpass fixture although the caller named an identity the run never verified; failure: `got envelope ref …; want ErrRequiredIdentityMissing: absent identity "absentIdentity" not verified`. The control arm (true identity) stays green |
| **M28 envelope write path (inherited)** | `host/evidence/proof_producer.go`: the write-error branch `if err := p.writer.PutObject(writeCtx, obj); err != nil` → `if err := p.writer.PutObject(writeCtx, obj); false && err != nil` (the `&& false` rides the condition; the init still runs so the mutant compiles) | `TestProducerEndToEndMintsEnvelopeValidatorAccepts` (AC5) and `TestProducerSurfacesUncertainWriteOutcome` (AC29) | A failed write returns the envelope ref anyway; the validator reads a nonexistent object; failure: `validator refused: got missing; want ResolvedGradeProven` — the end-to-end chain the shed removed is what the mutation breaks |
| **M31 read-deadline derivation** | `host/evidence/proof_producer.go`: `context.WithTimeout(ctx, p.cfg.ObjectReadTimeout)` (source read) → `context.WithTimeout(ctx, time.Hour)` | `TestRealStoreBlockedSourceReadReturnsWithinObjectReadTimeout` (AC23) | The blocked read waits out the pool and mints inside the watchdog; failure: `blocked source read minted after exceeding ObjectReadTimeout: elapsed=53ms; want operational error and no envelope` (the 2× seal-path ceiling reds on attempt 1) |
| **M32 write-deadline derivation** | `host/evidence/proof_producer.go`: `context.WithTimeout(ctx, p.cfg.ObjectWriteTimeout)` → `context.WithTimeout(ctx, time.Hour)` | `TestRealStoreBlockedEnvelopeWriteReturnsWithinObjectWriteTimeout` (AC20) | The blocked write completes after the decoy frees the connection, inside the watchdog; failure: `blocked envelope write stored after exceeding ObjectWriteTimeout: elapsed=53ms; want operational error and no object` |
| **M33 per-stream cap** | `host/evidence/proof_producer.go`: `readCapped`'s overflow decision `if int64(len(data)) > limit` → `if false && int64(len(data)) > limit` | `TestProducerRefusesStdoutOverflowIndependently` and `TestProducerRefusesStderrOverflowIndependently` (AC26) | The 373-byte real `ai-check` JSON passes a 16-byte cap unparsed-to-EOF; the proof mints from over-cap output; failure: `got envelope; want ErrOutputLimit (stdout 373 bytes > cap 16)` |
| **M34 anti-vacuity guard** | `host/evidence/proof_producer.go`: the single vacuity guard `if len(verified) == 0 { return … ErrVacuousCheck }` → `if false && len(verified) == 0 { … }` | `TestProducerInterpretationMatchesEveryFixtureShape` (AC27, the `nocontracts` arm) | The clean-but-vacuous fixture reaches the identity check and refuses with the wrong sentinel; failure: `got ErrRequiredIdentityMissing; want ErrVacuousCheck: 0 results records from non-empty output` — the guard's own identity is the observable |
| **M35 byte verification** | `host/evidence/proof_producer.go`: digest compare `if got.Digest() != checker.Compiler.Digest()` → `if false && got.Digest() != …` | `TestProducerRefutesTamperedExecutableBytes` (AC24) | The wrong-compiler arm executes the unverified binary and mints; failure: `got envelope; want ErrUntrustedTool: executable bytes …; want sha256:<expected>` |
| **M36 version verification** | `host/evidence/proof_producer.go`: first-line compare `if versionLine != checker.CompilerVersion` → `if false && versionLine != …` | `TestProducerRefutesWrongCompilerVersion` (AC24) | The wrong-version arm proceeds and mints a report the configured validator would refuse as tool mismatch; failure: `got envelope; want ErrUntrustedTool: version line "AILANG v9.9.9"; want "AILANG v0.41.0"` |
| **M37 envelope size bound** | `host/evidence/proof_producer.go`: `if len(envelope) > MaxBytes { return … ErrEnvelopeTooLarge }` → `if false && len(envelope) > MaxBytes { … }` | `TestProducerRefusesEnvelopeBeyondValidatorReadCap` (AC29) | The oversize envelope is stored; the validator's bounded read refuses it — the object exists where the test asserts absence; failure: `oversize envelope written: store contains object …; want refusal ErrEnvelopeTooLarge and no object` |
| **M38 run-cleanliness (errors>0)** | `host/evidence/proof_producer.go`: `if d.Verify.Errors > 0 { return … ErrUnverifiedRun }` → `if false && d.Verify.Errors > 0 { … }` | `TestProducerInterpretationMatchesEveryFixtureShape` (AC27, the `mixed`/`encerror` arms) | A run with a Z3 encoding error mints a report claiming a clean proof; failure: `got envelope over a run with verify.errors=1; want ErrUnverifiedRun` — the mixed fixture's one verified record must not launder it |
| **M39 exec-deadline derivation** | `host/evidence/proof_producer.go`: `context.WithTimeout(ctx, p.cfg.ExecTimeout)` (the `ai-check` run) → `context.WithTimeout(ctx, time.Hour)` | `TestProducerExecTimeoutKillsTheWholeProcessGroup` (AC25) | Configured with a tiny `ExecTimeout` against a real multi-second `ai-check` run, the mutant lets the child run to completion; failure: `got clean exit and a mint attempt; want ErrExecTimeout (run exceeded ExecTimeout)` — the group kill the test observes through the wrapped seam never fires |

The constructor's per-field refusals and the two ordering pins are killed by their AC5/AC30 named tests directly (each refusal is a separate observable branch — the iteration-106 lesson — so no additional mutation rows are needed beyond the named tests failing when a branch is neutered; the re-drill plan runs the standard one-mutation-per-branch drill at M5, including `skipped > 0`, `counterexample > 0`, `check.passed`, and `verify.available` neuterings against their fixture arms).

## §9 Conflict Surface

What already lives in the positions this design extends, and what must still pass:

| Neighbor / surface | Collision and ordering |
|---|---|
| **The evidence package's authority invariants** (`validator.go`, V11) | **No collision — the design's central constraint.** `ValidatedEvidence` has no exported fields and no public constructor; `Validator.ValidateProof` is the sole mint; `Resolve` the sole grade resolver; no free resolver exists; serialized receipts never carry authority. The producer mints DATA (an envelope object), never `ValidatedEvidence`; it touches no validator or codec file. `TestPublicAuthoritySurfaceIsFrozen` must still pass unedited |
| **The 37-test evidence manifest** (`scripts/verify_go.sh:27-64`, V22) | Every one of the 37 named tests must still pass; the producer's new tests are ADDED to `REQUIRED_EVIDENCE_TESTS` and `EXACT_EVIDENCE_TESTS` moves 37 → observed count in the same commit. The self-mutation gate `host/verifygate/evidence_manifest_gate_test.go` executes the live script and must still hold (it empties the set and expects refusal — unchanged by an additive edit) |
| **`host/store`** (rows 22/23, landed) | **No file changes.** `PutObject`/`ReadObject`/`BusyTimeout` consumed as-is; the strict guard (`requireDeadline`, store.go:236-247, V8) is the producer's defense-in-depth; the three-outcome contract (durable.go, V12) is consumed, not modified. Row 111's lock-wait work proceeds independently — nothing here touches `busy_timeout` or its DSN application |
| **`host/procbound`, `host/proctest`, `host/childenv`** (row 24, landed) | Reused unchanged. The row-24 design (`design_docs/implemented/w-host-subprocess-cleanup-boundary.md` §9) named this row "the next consumer" of D1's shape, `procbound` and `proctest` — this design is that consumer. The producer's `killGroup` seam is a per-package mirror (capsule/broker each have one), not a shared refactor |
| **`scripts/verify_ail.sh`** | **Untouched** — no `.ail` changes (V28); the producer composes with the gate's JSON discipline rather than competing with it (§D4). Fixture `.ail` files live under `host/evidence/testdata/`, outside the gate's `ROOTS` sweeps (`design_docs`, `world`, `packages/se-tools`; V21), on the existing precedent of `host/replay/testdata/transition_fixture.ail` (V28) |
| **`host/archive/check.go`** | The producer's subprocess-environment composition (RELAX_MODULES + private cache + `childenv.Scrubbed` + `TMPDIR` redirect) is check.go's measured shape extended by one measured assignment (V19/V26); check.go itself is untouched. Row 153's template-build behavior is NOT adopted here — `ai-check` has no compile-cache promotion semantics, and inventing them is out of scope |
| **Row 25 (landed), rows 112/113 (open)** | No overlap: row 25's blocked-child-kill arm is capsule's; 112/113 own archive/replay/pkgproj sites. The producer is a NEW launch site that ships in row-24's bounded shape from day one, adding no new descendant leak for them to inherit |
| **Item 17's decomposition** (§9, parent) | This row is tranche 2 of 5 in that decomposition; tranches 3 (replay), 4 (renderer), and the production key-wiring row are untouched. Ordering holds: nothing downstream may display `PROVEN` until tranche 4; nothing here changes that |
| **The queue row's own ordering gate** | Met before this design started: row 23's 11 → 0 landed at `b9ecabe` (charter row 26's gate note); this design adds NO hidden twelfth deadline-free root — both store calls derive their own deadlines, and the strict guard would refuse them if they did not (V8) |

## §10 Verification Log

All commands run from this worktree root at `ba79d1e` on 2026-10-12 (V1) with `/usr/bin/grep` (this rig's default `grep` is ugrep; recorded commands must be reproducible). Negative-existence rows carry a KNOWN-POSITIVE CONTROL in the same call — a pattern that fires in the same scope, proving the instrument works. Fixture-generation rows name the script that ran them; the script, its inputs, and its outputs are committed under `design_docs/verification/world-iter251/`.

| ID | Claim | Exact command and same-call control | Observed output |
|---|---|---|---|
| V1 | Measurement base; only this doc and its evidence dir are new | `git rev-parse --short HEAD; git status --short` | `ba79d1e`; two lines: `?? design_docs/planned/w-bounded-z3-report-producer.md`, `?? design_docs/verification/world-iter251/` |
| V2 | The parent doc's required sections exist at the cited places | `/usr/bin/grep -n "^### 3.4 SHED\|^### 3.5 What\|^## 10.11 Round 9b\|^## 10.12 Round 10\|^## 10.15 Round 12\|^## 10.16 .*ruling" design_docs/implemented/w-validated-proven-evidence-boundary.md` | hits at 489, 499, 1662, 1731, 1974, 2057 — §3.4 shed marker quotes row 26 verbatim; §10.11 carries the round-9b `NewProducer`/`GenerateProof`/`ObjectWriter` record; §10.15 carries both objections verbatim |
| V3 | Queue row 26's text, inherited defects, and pricing exist as cited | `/usr/bin/grep -n "w-bounded-z3-report-producer" design_docs/world-mission.md` then `awk 'NR==4451' … \| fold -w 200` | line 4451: v1.0.0 prerequisite, `D-WORLD-69` = A on path 26 → 105 → 172 → 150; gate-met note pointing at `w-store-deadline-free-residue-owner.md` §3; the former §3.4 scope list; AC5/M15 + AC20/M28 inheritance; both INHERITED OPEN DEFECT blocks verbatim; `1.0 d carried … round-12 fixes remain unpriced, pending design` |
| V4 | No producer surface exists in production Go | `printf 'check: '; /usr/bin/grep -rn "NewProducer\|GenerateProof\|ObjectWriter\|proof_producer" host/ cmd/ scripts/ --include='*.go' \| wc -l; printf 'control: '; /usr/bin/grep -rn "NewValidator" host/ --include='*.go' \| wc -l` | `check: 0`; control `26` (the validator's constructor fires — the instrument sees package symbols) |
| V5 | No `WriteObject` store surface exists (the additive round-11b surface was never built) | `printf 'check: '; /usr/bin/grep -rn "WriteObject" host/ cmd/ scripts/ --include='*.go' \| wc -l; printf 'control: '; /usr/bin/grep -rn "func (s \*Store) PutObject" host/store/ \| wc -l` | `check: 0`; control `1` (the existing durable write exists and fires) |
| V6 | No production Go runs `ai-check` | `printf 'check: '; /usr/bin/grep -rn "ai-check" host/ cmd/ --include='*.go' \| /usr/bin/grep -v _test \| wc -l; printf 'control: '; /usr/bin/grep -c "ai-check" scripts/verify_ail.sh` | `check: 0`; control `13` (the shell gate runs it; no Go code does) |
| V7 | Every proposed new identifier is unallocated in `host/ cmd/ scripts/` | `for tok in ErrInvalidProducerConfig ErrUntrustedTool ErrVacuousCheck ErrRequiredIdentityMissing ErrUnverifiedRun ErrProducerExecTimeout ErrEnvelopeTooLarge ErrExecTimeout ErrOutputLimit ErrMalformedCheck ProducerConfig CheckerConfig MaxSourceBytes; do printf '%s=%s ' "$tok" "$(/usr/bin/grep -rn "$tok" host/ cmd/ scripts/ 2>/dev/null \| wc -l \| tr -d ' ')"; done; printf '\ncontrol ErrNoDeadline=%s\n' "$(/usr/bin/grep -rn "ErrNoDeadline" host/ \| wc -l \| tr -d ' ')"` | all 12 tokens read `0`; control `16` (an existing sentinel fires, so the empty reads mean absence, not a broken grep). `pipeCloseGrace` reads 15 across `host/capsule`+`host/broker` — the row-24 precedent the producer's per-package const mirrors; no Go cross-package collision is possible |
| V8 | The store's live contract: strict deadline guard on every I/O method | `/usr/bin/grep -rn "requireDeadline(" host/store --include='*.go' \| /usr/bin/grep -v _test \| wc -l; /usr/bin/grep -n "func requireDeadline" -A 9 host/store/store.go; /usr/bin/grep -n "ErrNoDeadline = errors.New\|func (s \*Store) PutObject\|func (s \*Store) GetObject\|INSERT OR IGNORE" host/store/store.go` | 38 guarded call sites incl. `PutObject` (func :543, guard :547), `GetObject` (func :579, guard :583), `ReadObject` (read_object.go:47); `requireDeadline` (store.go:239-247) rejects `ctx == nil` and no-deadline ctx with `ErrNoDeadline` (:237); `PutObject`'s idempotency comment :540-541 and the `INSERT OR IGNORE INTO objects` SQL :566 |
| V9 | `SetMaxOpenConns(1)` — the pool pin the decoy stimuli rely on | `/usr/bin/grep -n "SetMaxOpenConns" host/store/*.go` | `store.go:382: db.SetMaxOpenConns(1)` — note the parent's AC16 cites `store.go:298`; the line has DRIFTED at HEAD and this doc cites the measured line |
| V10 | `BusyTimeout()` is a cached, nonblocking accessor | `/usr/bin/grep -n "func (s \*Store) BusyTimeout\|busyTimeout " host/store/store.go \| head -5` | `:217` field, `:370` accessor returning the field, set at `:299`/`:327`/`:364` from the DSN at Open — never a live query |
| V11 | The validator's live surface this design composes with | `/usr/bin/grep -n "type ObjectReader interface" -A 3 host/evidence/validator.go; sed -n '63,71p' host/evidence/validator.go; /usr/bin/grep -n "WithTimeout\|SumSHA256\|hmac.Equal\|meta.SemanticID\|meta.InterfaceHash\|func containsIdentity\|ErrUnorderedTimeouts\|BusyTimeout is unknown\|window < 0" host/evidence/validator.go` | `ObjectReader` = `ReadObject` + `BusyTimeout()` (:50-53); `NewValidator(key, reader, cfg, requiredIdentities)` at :63 with the measured four-separate-refusals comment (:64-70), per-field refusals and both ordering/unknown-window pins (:72-93, window check :86-91); `ValidateProof` derives `context.WithTimeout(ctx, ObjectReadTimeout)` (:124), recomputes `SumSHA256(payload)` vs ref (:137), checks `meta.SemanticID` (:140) and `meta.InterfaceHash` (:143), recomputes HMAC and `hmac.Equal`s the tag (:156-159); `containsIdentity` binary search (:185-196); `ErrUnorderedTimeouts` sentinel (:30) |
| V12 | `PutObject`'s wait/cutoff contract (three outcomes) | `sed -n '40,100p' host/store/durable.go` | `beginDurable` (:52) acquires under ctx — "database/sql honours ctx while waiting for the connection … a held connection bounds this wait by the caller's deadline, not busy_timeout" (:48-51); `finishDurable` (:73) is the cutoff: pre-cutoff errors are definite not-committed, in-flight commits yield `*UncertainError` with "reconcile before any retry" (:10-45, :73-96) |
| V13 | The codec bounds this design inherits, and the guard-nothing-can-pin lesson's home | `/usr/bin/grep -n "MaxBytes\|MaxStringBytes\|MaxVerifiedIdentities\|ProofReportSchemaV1" host/evidence/types.go \| head -5; sed -n '60,68p' host/evidence/types.go` | `MaxBytes = 256 * 1024`, `MaxStringBytes = 1024`, `MaxVerifiedIdentities = 256`, `ProofReportSchemaV1 = "world/proof-report/v1"` (types.go:13-17); the `decoded`-field comment carries "Iteration 106 measured that branch … A branch nothing can pin is not a guard" (:60-68) |
| V14 | The report codec enforces sorted-unique verified | `sed -n '86,97p' host/evidence/report_codec.go` | the strict pairwise-increase branch refuses `verified identities are not sorted and unique` (:95-96) — with the in-row note that a redundant `sort.StringsAreSorted` belt was removed at iteration 105 as unreachable (:90-94) |
| V15 | The envelope codec's canonical two-member shape | `/usr/bin/grep -n "envelopeFields\|mac\|report" host/evidence/envelope_codec.go \| head -6` | `envelopeFields = {"report","mac"}`; strict decode; mac spelling/width retained as `MACValid` data for the authentication step — the exact bytes the producer must mint |
| V16 | The metadata identities the envelope object must carry | `/usr/bin/grep -n "ProofSemanticID\|InterfaceHashV1" host/evidence/validator.go \| head -4` | `ProofSemanticID = ProofReportSchemaV1` (:15) and `InterfaceHashV1 = SumSHA256("world/authenticated-proof-envelope/v1")` (:18-19); the validator checks `meta.SemanticID` (:140) and `meta.InterfaceHash` (:143) |
| V17 | Capsule's ratified bounded-subprocess shape (D3's template) and the landed decoy-test discipline | `/usr/bin/grep -n "Setpgid\|killGroup\|readCapped\|verifyExecutable" host/capsule/capsule.go \| head -8; sed -n '144,210p' host/evidence/realstore_test.go` | `:246` Setpgid, `:83/:97` killGroup seam, `:392` readCapped limit+1, `:377` verifyExecutable; the AC16 test: `minDecoyHold=20ms`, `decoyHoldRatio=20`, retry ≤ 5, 2× seal ceiling, 20× watchdog, announced GOMAXPROCS(2) precondition |
| V18 | `procbound`'s API surface (Admit/Wait/NothingToKill) | `/usr/bin/grep -n "^func \|^var Err\|^const MaxOutstanding" host/procbound/procbound.go` | `Admit()` nonblocking CAS reservation, `Wait(wait, d, release)` bounded reap, `MaxOutstanding = 8`, `ErrCleanupIncomplete`, `ErrCleanupBacklog`, `NothingToKill`/`GroupEmpty` (the zombie-group ESRCH/EPERM filter) |
| V19 | The publish-check subprocess environment is the measured composition to mirror | `sed -n '95,103p' host/archive/check.go` | `AILANG_RELAX_MODULES=1` with the iter-195 measured rationale; `CacheDirEnv` appended last "so an ambient value cannot win"; `cmd.Dir = root`; `checkTimeout = 10s` |
| V20 | `childenv.Scrubbed` strips registry variables only | `sed -n '34,48p' host/childenv/childenv.go` | strips the four registry assignments, keeps everything else (incl. `TMPDIR`), always non-nil — the nil-Env safety property |
| V21 | The gate's JSON-parse discipline and required-identity manifest are as described | `sed -n '1,16p' scripts/verify_ail.sh; sed -n '315,395p' scripts/verify_ail.sh; /usr/bin/grep -n "GATE_LEG_TIMEOUT_S=\|^ROOTS=(\|-timeout 5s\|never" scripts/verify_ail.sh \| head -6` | header (:5-8) names the two silent classes (Z3-error-exits-0, vanished-contract-exits-0) and why JSON is authoritative; the parser (:315-395) reads `check.passed` (:368), `verify.errors` (:370), `verify.counterexample`, and `verify.results[].function`/`.status` with per-identity required-set checks; `ai-check -timeout 5s` at :307; `GATE_LEG_TIMEOUT_S=120` at :160-162; `ROOTS` sweeps only `design_docs`, `world/`, `packages/se-tools/` (:165-169); the stdout-only/never-`2>&1` stderr-merge warning at :99-104 |
| V22 | The evidence manifest gate's current shape | `python3 -c "import re; raw=open('scripts/verify_go.sh').read(); print(len(re.findall(r'\"(Test\\\\w+)\"', re.search(r'REQUIRED_EVIDENCE_TESTS = \{(.*?)\}', raw, re.S).group(1))))"; /usr/bin/grep -n "37 exact\|check_evidence_manifest \"\$evidence_json\" 37" scripts/verify_go.sh; /usr/bin/grep -n "skipped" scripts/verify_go.sh \| head -2` | 37 members; the literal `37` appears in the banner (:295) and the call (:306); the parser refuses `skipped=` required tests (:100-103) — binary-dependent producer tests must fatal, never skip |
| V23 | The pinned binary's identity | `$HOME/.pinned-ailang/ailang --version 2>/dev/null \| head -1; shasum -a 256 $HOME/.pinned-ailang/ailang` | first line `AILANG v0.41.0`; full: `1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f` |
| V24 | All seven fixtures are REAL pinned-binary output | `bash design_docs/verification/world-iter251/gen_aicheck_fixtures.sh` (the committed generator; scratch modules in a mktemp dir, never in the repo) | banner `AILANG v0.41.0`; per-shape `rc=` and byte counts; the seven committed `.json` files with `rc` sidecars: verifiedpass rc=0 verified=1; refuted rc=1 cex=1; encerror rc=1 errors=1; mixed rc=1 verified=1 errors=1; skippedshape rc=0 skipped=1; loaderror rc=1 passed=false; nocontracts rc=0 results=[] |
| V25 | **The silent class changed shape on v0.41.0.** The inherited claim "a Z3 encoding error exits 0 silently" (parent V10, v0.30.0) does not hold for the encoding-error fixture on v0.41.0; a DIFFERENT rc-0 silent class (skipped) is live | generated in V24; summarised: `python3 - <<…` printing rc/verified/errors/skipped/cex/statuses per fixture | `encerror`: rc **1**, errors=1 (the encoding-error shape now signals); `mixed`: rc 1; `skippedshape`: rc **0**, skipped=1, verified=0 — a clean exit code over an unverified contract. The JSON-authoritative discipline is retained and strengthened, not relaxed |
| V26 | The binary needs a writable `TMPDIR` for solver scratch files; a scratch-root redirect works | two arms: `env -i … (no TMPDIR)` → the run's JSON `reason`; `env -i … TMPDIR=$scratch/tmpdir` → verified count | no-TMPDIR arm: `solver error: creating temp file: open /tmp/ailang-smt-….smt2: operation not permitted` (rc 1, status error); redirect arm: rc 0, `verified=1`, and the tmpdir holds 0 entries after the run — D3's `TMPDIR` redirect is measured, not assumed |
| V27 | The clean-but-vacuous shape exists: rc 0, all counters zero, `results: []` | the V24 `nocontracts` arm (module with one type + one uncontracted function) | rc=0, `check.passed=true`, `verify.available=true`, verified/errors/skipped/counterexample all 0, `results: []` — the anti-vacuity guard's fixture |
| V28 | No `.ail` file changes in this design, and in-`host/` testdata `.ail` is an existing pattern outside the gate's sweep | `find host design_docs/verification -name '*.ail'; printf 'control: '; find world -name '*.ail' \| wc -l; /usr/bin/grep -n "ROOTS=(" -A 5 scripts/verify_ail.sh` | one pre-existing file: `host/replay/testdata/transition_fixture.ail` (the precedent for committed fixture modules under `host/*/testdata/`); control: 4 under `world/`; `ROOTS` sweeps `design_docs`, `./world`, `packages/se-tools/se_tools` only — `host/` is not swept, so committed fixture modules cannot perturb the LEG1 membership compare |
| V29 | `MaxSourceBytes = 1 MiB` has ~40× headroom over the largest repo module | `find world packages design_docs/sketches -name '*.ail' -exec wc -c {} + \| sort -rn \| head -3` | largest: 25,422 bytes (`packages/se-tools/se_tools/run.ail`); then 23,826 and 21,538 |
| V30 | Row 111 (lock-wait owner) is OPEN; rows 22/23/24/25 are LANDED | `/usr/bin/grep -n "^111\. \|^23\. \|^22\. " design_docs/world-mission.md \| cut -c1-120` | 111 has no LANDED marker (owns late-read 503, 500-vs-503, arm A per-request busy_timeout); 22 `[LANDED 2026-09-25 iter-189]`; 23 `[POLICY TRANCHE M5a–M7b LANDED 2026-09-28 iter-207 PR #162]`; 25 `[LANDED 2026-09-27 iter-200 PR #155]` |
| V31 | The `hashref` API the byte check and refs use | `/usr/bin/grep -n "func Sum\|func SumSHA256\|func Parse\|func New\|type HashRef" host/hashref/hashref.go` | `Sum(algo, payload)` (the ref's own algorithm — capsule's byte-check shape), `SumSHA256`, `Parse`, `New`; `HashRef` with `Algo()/Digest()/String()/IsZero()` |
| V32 | Binary-dependent Go tests fatal, never skip; CI exports the pin | `/usr/bin/grep -n "never skip\|AILANG_BIN unset" host/archive/check_test.go \| head -3; /usr/bin/grep -n "AILANG_BIN" .github/workflows/ci.yml \| head -3` | `check_test.go:81` "t.Fatal … never skip"; CI installs the pinned v0.41.0 with a `case` asserting the version token and exports `AILANG_BIN` to `$GITHUB_ENV` (ci.yml:134-173) |
| V33 | The AILANG constructs used by the fixture modules are valid on the pinned binary (every language claim in this doc is live-checked) | the V24 generator run itself: module header, `export func`, `! {}` effect marker, `ensures`, `tests`, string interpolation in body, sibling-call `ensures`, `export type`, uncontracted `export func` | all seven scratch modules loaded and produced `ai-check` JSON on the pinned v0.41.0 (the four rc-0/rc-1 verdicts in V24-V27 are the check results); the load-failure fixture additionally pins the parse-error shape (`PAR_EXPORT_REQUIRES_FUNC`) |

**Sweep note (the directive's own rule):** every sentence in this document claiming the codebase does or does not do X cites a row above. Negative claims (no producer surface: V4; no `WriteObject`: V5; no Go `ai-check` runner: V6; tokens unallocated: V7; no `.ail` change: V28) each carry a same-call firing control. No sentence claims a store property without V8-V12, a validator property without V11/V13-V16, or a tool property without V19-V26/V31-V33.

---

**Document created**: 2026-10-12 · designer role, iteration 251 · measurement base `ba79d1e` · instrument AILANG v0.41.0 (`~/.pinned-ailang/ailang`)