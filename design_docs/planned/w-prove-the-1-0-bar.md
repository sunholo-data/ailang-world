# w-prove-the-1-0-bar — the two deliverables that decide whether World 1.0 exists, and the anti-vacuity conditions that make them mean something

- Status: **planned** (design only; nothing implemented) · Date: **2026-09-21** · Rows **92** (clause-5) and **93** (clause-4) · Filed at the attended grooming of 2026-09-21 (`72eace9`)
- Scope: ONE doc for BOTH rows, deliberately. Clause 4 is the floor that proves clause 5's value is not paid for with resident-agent regression — they are two halves of one argument, and splitting them would double the review cost without separating any decision.
- Revision 1 (quorum r1 objection): Conflict Surface section added

---

## §1 Problem — the 1.0 bar named two deliverables and the queue contained neither

Measured at the 2026-09-21 grooming: **34 open rows, 27 of them clause-2 (79%)**, of which
thirteen were mission-harness or bookkeeping and eight more were tests about tests. Clause 5 —
the provenance teeth, *"capability the shell cannot express at any pass-rate"*, the sentence that
says what World is FOR — had three open rows, all workbench render defects. Clause 4 had none.

`grep -c "provenance walk|why did X happen"` over the charter returns **2**, both inside the
clause-5 definition itself and **zero in any queue row**. The two things that decide whether World
1.0 exists were absent from the queue meant to deliver it, for fourteen months of mission time.

That is not a scheduling accident. It is the same failure the fleet's Gate-2 admissibility rule
was written for on the same day, one level up: every row was real, and they summed to a mission
maintaining itself instead of shipping. **Local correctness does not aggregate.**

---

## §2 What clause 5 demands, verbatim — and what it does not

> on ≥3 REAL "why did X happen" questions (arising from actual operation, not synthetic), a
> provenance walk yields the verified answer in ≤5 minutes each, where the pre-World method was a
> grep/log archaeology session.

Four load-bearing words, each of which is an anti-vacuity condition, not decoration:

| Word | What it forbids |
|---|---|
| **REAL** | A question invented because the walk can answer it. Proves the walk, not the capability |
| **≥3** | One lucky question |
| **verified** | An answer whose only source is the walk being assessed. Circular |
| **pre-World method** | A ≤5-minute claim with no baseline. "Fast" is meaningless unarmed |

**The sharpest of these is `verified`.** If the provenance walk is the sole source of the answer,
"verified" reduces to "the walk said so" and the demonstration is self-confirming. The answer must
be checkable against a record the walk did not produce — git history, the session JSONL, the
driver log, the banked eval row. **A question with no independent ground truth is not eligible**,
however real it is.

Clause 5 also says what this is NOT: it is *not* a pass-rate claim. Clause 4 owns that, and clause
4 is explicitly "the FLOOR, deliberately NOT the value proof."

---

## §3 The R1 constraint — the control arm already exists and it is this mission

Clause 5's standing-evidence clause: *"R1 — this mission itself migrating onto World, measured on
incident classes eliminated + human attention per workstream against the 2026-07 markdown-and-
launchd baseline (tonight's operation IS the control arm)."*

This is a gift and a constraint. The questions **must be harvested from this mission's recorded
operation**, not invented — and that record is unusually rich in exactly the right shape, because
the mission has spent months writing down *why* things went wrong.

Three candidates, each already recorded, each with independent ground truth:

1. **Why did iteration 171 have a full log entry and no index row?** Ground truth: the pin-version
   gap — `~/.pinned-ailang/ailang mission rotate-log world` prints `unknown command 'mission'` on
   v0.30.0. Checkable against the binary and the index.
2. **Why did three designer runs execute unfenced while the rulebook said they were sandboxed?**
   Ground truth: `mission_pi_run.sh:155` passes no `-e` flags. Checkable against the script.
3. **Why did a green `shasum -c` faithfulness proof pass over a destroyed commit split?** Ground
   truth: the manifest covers the final tree, and the final tree is identical whether the split is
   correct or collapsed. Checkable against `git show --stat` per commit.

Each took a human a nontrivial archaeology session to find the first time — which is precisely the
pre-World baseline the clause asks to be measured against. **The baseline is not hypothetical; it
is in the log, with its cost.**

---

## §4 Phase A (row 92) — the demonstration

1. **Harvest.** Select ≥3 questions from the mission's own record. Eligibility: arose from real
   operation · has independent ground truth · the original diagnosis cost is recoverable from the
   record. Record why each was selected, and **record any question considered and rejected**, with
   the reason — a harvest with no rejections is evidence the bar was set after the fact.
2. **Baseline.** For each question, state the pre-World method and its measured cost, sourced from
   the record rather than estimated. Where the record does not say, say so — an unmeasured baseline
   is reported as unmeasured, never imputed.
3. **Walk.** Answer each by a provenance walk on World. Time it. The timer starts at the question,
   not at the first successful query.
4. **Verify.** Check each answer against its independent ground truth. An answer the walk produced
   and nothing else confirms is recorded as **unverified**, and does not count toward the ≥3.
5. **Artifact.** The comparison table, not the tooling. Questions × (baseline method, baseline
   cost, walk time, verification source, verdict).

**Honest failure is a result.** If the walk cannot answer a real question, that is clause 5
evidence too, and the row records it rather than substituting an easier question. A demonstration
that cannot fail demonstrates nothing.

---

## §5 Phase B (row 93) — the floor, and why it is second

Clause 4 specifies the gate fully: two reference agents from different providers (Claude Code
agent mode, codex CLI), paired arms *shell* (native tools) vs *World* (MCP transition tools only),
same model, same benchmark set (standard tier), **N≥3 per arm**; the floor holds only if **both**
show World pass-rate ≥ shell − 2pp **and** median wall-clock overhead ≤ +25%. motoko is an
optional third arm, informative, never gate-blocking.

**Why after Phase A:** clause 4 is the floor, not the value. A floor measured before there is
value to protect measures nothing — it would report a number with no decision attached to it. The
ordering is a design choice, not a capability dependency; Phase B could run first and would simply
be less useful.

---

## §6 The precondition decides gate VALIDITY, not gate OUTCOME — and it has a dependency nobody has checked

Clause 4's stability precondition: a reference agent is gate-eligible only if its **shell** arm
completes N≥3 runs with **zero harness-fault failures** (`api_error` / `resource_limit` / harness
classes) and pass-rate range ≤5pp. If both agents are ineligible the gate **PAUSES for
instrumentation** — explicitly *"a measurement problem, explicitly NOT a value-park."*

**This must be reported separately from the outcome**, and in this order: eligibility first, gate
second. A gate outcome reported without its validity check is a number that cannot be acted on.

**The dependency worth naming before anyone runs this:** eligibility is defined in terms of
`error_category`, and the fleet's own operating rule is that **`api_error` is the catch-all
meaning "cause unknown", not "the model failed"**. So a run misclassified as `api_error` makes an
agent *ineligible* and can pause the gate for instrumentation that is not needed — or, worse, a
genuine harness fault classified as something else lets an ineligible agent through and produces a
floor measurement that is metrologically unsound. **Before the arms run, confirm that the
classifier's `api_error` rate on the shell arm is understood**, because the precondition inherits
whatever that classifier gets wrong.

---

## §7 What would make this vacuous

Enumerated because this mission's queue is full of rows about gates that passed by construction,
and this doc should not add another.

- Questions invented to fit the walk → §4.1's rejection log is the guard.
- "Verified" by the walk itself → §2's independent-ground-truth rule.
- A baseline asserted rather than sourced → §4.2 reports unmeasured as unmeasured.
- Timing started after the first successful query → §4.3.
- A floor gate reported without its eligibility check → §6.
- A demonstration with no possible failure mode → §4's honest-failure clause.

---

## Conflict Surface

> Phase B evaluates against the standard tier. This OVERLAPS with the existing mission-harness.
> Phase B will strictly REUSE the existing harness rather than introducing a parallel benchmark
> script.

This is true and needs to be precise, because the word "harness" points at two different machines
and only one of them is Phase B's.

**Phase B reuses the fleet's agent-benchmark machinery — invoked, not forked.** The standard-tier
benchmark set and the reference-agent arms run under `ailang eval --benchmark <name>` (one AI
benchmark) and `ailang eval-suite --models <csv>` (the full suite), the same evaluation tooling
this mission's clause-2 rows and the fleet's eval KPI already lean on. Phase B introduces **no
parallel benchmark script**: nothing new that decides pass/fail or renders wall-clock is written
and owned inside this repo. The pairing of the two arms (shell native tools vs World MCP transition
tools) is an *invocation-layer* concern — which `-model` / `-benchmark` flags and models each arm
passes into the existing suite — not a second runner. Reusing the harness is exactly what keeps
Phase B a measurement of World, not a measurement of new tooling.

**The other harness is out of Phase B's scope.** This repository's own benchmark machinery,
`scripts/bench_worldd.sh` + `bench/BASELINE.md`, is a **non-vacuous daemon smoke benchmark**: it
runs `go test -bench . -benchtime 1x` on `./host/daemon/` and validates the benchmark evidence
blocks, measuring the daemon's wall-clock/load behavior (V7). It is not an agent benchmark runner
and has nothing to say about shell-vs-World pass rates. It is explicitly **out of Phase B's scope**;
Phase B neither edits it nor loads it, so the daemon-smoke gate and the clause-4 floor stay cleanly
separate measurements.

**Phase A (row 92) touches no benchmark machinery at all.** Its only executable surface is the
provenance walk against the worldd daemon query API — harvesting, timing, and verifying answers to
the mission's recorded questions. No `ailang eval`, no `bench_worldd`; nothing Phase B reuses is
even started during the demonstration.

---

## §8 Verification log

| # | Claim | How verified | Result |
|---|---|---|---|
| V1 | The clause-5 bar is not a queue row | `grep -c` on the charter for the bar's phrases | **2 hits, both in the clause definition, 0 in rows** |
| V2 | Clause 4 has no open row | clause tally at grooming | **0 open clause-4 rows** of 34 |
| V3 | The queue was 79% clause-2 catch-all | grooming tally, 2026-09-21 | 27 of 34 open; 13 mission-harness, 8 tests-about-tests |
| V4 | The three candidate questions are real and recorded | read the charter rows and mission log | rows 86 (index row), 78 (unfenced pi), 81 (faithfulness proof) |
| V5 | `api_error` is a catch-all, not a model verdict | the fleet's own operating rule (CLAUDE.md instrument table) | confirms §6's dependency |
| V6 | Stage A of the self-mod publish is green, so clause 7 is one attended step away | ran `build_world_package.sh` + `verify_world_package.sh` 2026-09-21 | **9/9 PASSED**, ready packet equals the committed golden byte-for-byte, projection reproduces with zero git drift |
| V7 | This repo's benchmark machinery is a daemon smoke benchmark, not an agent runner: `scripts/bench_worldd.sh` + `bench/BASELINE.md`, running `go test -bench` on `./host/daemon/` | read `scripts/bench_worldd.sh` + `bench/BASELINE.md` at HEAD `2678c1f` | `bench_worldd.sh` invokes `go test -bench . -benchtime 1x -run '^$' ./host/daemon/` and validates benchmark evidence blocks (`BenchmarkStoreCommit`, `BenchmarkRESTCommit`, ...); BASELINE is the day-1 kernel/broker wall-clock budget — daemon-scoped, no agent scaffolding. **Measures daemon wall-clock/load; not an agent benchmark runner; out of Phase B's scope** |
| V8 | Phase B's agent-benchmark machinery is the fleet's `ailang eval`/`eval-suite` | ran `ailang help` + `ailang eval run --help` at HEAD `2678c1f` | `ailang help` lists `eval — Benchmarks: run one, or a subcommand (suite, report, elo, ...)` and shows `ailang eval --benchmark fizzbuzz --mock` and `ailang eval-suite --models gpt5,claude-sonnet-4-6`; `eval run --help` exposes `-model` and `-benchmark`. Confirms arm pairing is an invocation-layer (flags/models) concern, not a parallel runner |

**Not verified, and named as such:** the pre-World baseline COST for each of the three candidate
questions. The record states what went wrong and how it was found; whether it states how long that
took is unknown until Phase A step 2 reads it. If it does not, that question's baseline is reported
as unmeasured rather than estimated.

---

## §9 Non-goals

- **Building new provenance tooling.** Phase A uses what World has. If the walk needs a feature,
  that is a finding, filed as its own row — not scope here.
- **Tuning for the floor.** Clause 4 says "after honest tuning"; honest means the tuning is
  recorded and applied to both arms.
- **Proving superiority.** Clause 4 is non-inferiority by design.
- **Clause 7.** Row 8's `SM.D` is a separate, attended, irreversible step.
- **Writing a parallel benchmark runner.** §Conflict Surface: Phase B invokes the fleet's `ailang eval`/`eval-suite`; it does not fork them, and `scripts/bench_worldd.sh`/`bench/BASELINE.md` (daemon smoke) is out of scope.
